package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/pkg/errors"

	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmmath "github.com/tendermint/tendermint/libs/math"
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

var MaxVoters = 20

// VoterSet represent a set of *Validator at a given height.
type VoterSet struct {
	// NOTE: persisted via reflect, must be exported.
	Voters []*Validator `json:"voters"`

	// cached (unexported)
	totalVotingPower int64
}

func NewVoterSet(valz []*Validator) *VoterSet {
	sort.Sort(ValidatorsByAddress(valz))
	vals := &VoterSet{Voters: copyValidatorListShallow(valz), totalVotingPower: 0}
	vals.updateTotalVotingPower()
	return vals
}

func (voters *VoterSet) ValidateBasic() error {
	if voters.IsNilOrEmpty() {
		return errors.New("voter set is nil or empty")
	}

	for idx, val := range voters.Voters {
		if err := val.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid validator #%d: %w", idx, err)
		}
	}

	return nil
}

// IsNilOrEmpty returns true if validator set is nil or empty.
func (voters *VoterSet) IsNilOrEmpty() bool {
	return voters == nil || len(voters.Voters) == 0
}

// HasAddress returns true if address given is in the validator set, false -
// otherwise.
func (voters *VoterSet) HasAddress(address []byte) bool {
	idx := sort.Search(len(voters.Voters), func(i int) bool {
		return bytes.Compare(address, voters.Voters[i].Address) <= 0
	})
	return idx < len(voters.Voters) && bytes.Equal(voters.Voters[idx].Address, address)
}

// GetByAddress returns an index of the validator with address and validator
// itself if found. Otherwise, -1 and nil are returned.
func (voters *VoterSet) GetByAddress(address []byte) (index int32, val *Validator) {
	idx := sort.Search(len(voters.Voters), func(i int) bool {
		return bytes.Compare(address, voters.Voters[i].Address) <= 0
	})
	if idx < len(voters.Voters) && bytes.Equal(voters.Voters[idx].Address, address) {
		return int32(idx), voters.Voters[idx].Copy()
	}
	return -1, nil
}

// GetByIndex returns the validator's address and validator itself by index.
// It returns nil values if index is less than 0 or greater or equal to
// len(VoterSet.Validators).
func (voters *VoterSet) GetByIndex(index int32) (address []byte, val *Validator) {
	if index < 0 || int(index) >= len(voters.Voters) {
		return nil, nil
	}
	val = voters.Voters[index]
	return val.Address, val.Copy()
}

// Size returns the length of the validator set.
func (voters *VoterSet) Size() int {
	return len(voters.Voters)
}

func copyValidatorListShallow(vals []*Validator) []*Validator {
	result := make([]*Validator, len(vals))
	copy(result, vals)
	return result
}

// VoterSet.Copy() copies validator list shallow
func (voters *VoterSet) Copy() *VoterSet {
	return &VoterSet{
		Voters:           copyValidatorListShallow(voters.Voters),
		totalVotingPower: voters.totalVotingPower,
	}
}

// Forces recalculation of the set's total voting power.
// Panics if total voting power is bigger than MaxTotalVotingPower.
func (voters *VoterSet) updateTotalVotingPower() {
	sum := int64(0)
	for _, val := range voters.Voters {
		// mind overflow
		sum = safeAddClip(sum, val.VotingPower)
		if sum > MaxTotalVotingPower {
			panic(fmt.Sprintf(
				"Total voting power should be guarded to not exceed %v; got: %v",
				MaxTotalVotingPower,
				sum))
		}
	}

	voters.totalVotingPower = sum
}

func (voters *VoterSet) TotalVotingPower() int64 {
	if voters.totalVotingPower == 0 {
		voters.updateTotalVotingPower()
	}
	return voters.totalVotingPower
}

// Hash returns the Merkle root hash build using voters (as leaves) in the
// set.
func (voters *VoterSet) Hash() []byte {
	if len(voters.Voters) == 0 {
		return nil
	}
	bzs := make([][]byte, len(voters.Voters))
	for i, voter := range voters.Voters {
		bzs[i] = voter.Bytes()
	}
	return merkle.HashFromByteSlices(bzs)
}

// VerifyCommit verifies +2/3 of the set had signed the given commit.
//
// It checks all the signatures! While it's safe to exit as soon as we have
// 2/3+ signatures, doing so would impact incentivization logic in the ABCI
// application that depends on the LastCommitInfo sent in BeginBlock, which
// includes which voters signed. For instance, Gaia incentivizes proposers
// with a bonus for including more than +2/3 of the signatures.
func (voters *VoterSet) VerifyCommit(chainID string, blockID BlockID,
	height int64, commit *Commit) error {

	if voters.Size() != len(commit.Signatures) {
		return NewErrInvalidCommitSignatures(voters.Size(), len(commit.Signatures))
	}

	// Validate Height and BlockID.
	if height != commit.Height {
		return NewErrInvalidCommitHeight(height, commit.Height)
	}
	if !blockID.Equals(commit.BlockID) {
		return fmt.Errorf("invalid commit -- wrong block ID: want %v, got %v",
			blockID, commit.BlockID)
	}

	talliedVotingPower := int64(0)
	votingPowerNeeded := voters.TotalVotingPower() * 2 / 3
	for idx, commitSig := range commit.Signatures {
		if commitSig.Absent() {
			continue // OK, some signatures can be absent.
		}

		// The voters and commit have a 1-to-1 correspondance.
		// This means we don't need the voter address or to do any lookup.
		voter := voters.Voters[idx]

		// Validate signature.
		voteSignBytes := commit.VoteSignBytes(chainID, int32(idx))
		if !voter.PubKey.VerifySignature(voteSignBytes, commitSig.Signature) {
			return fmt.Errorf("wrong signature (#%d): %X", idx, commitSig.Signature)
		}
		// Good!
		if commitSig.ForBlock() {
			talliedVotingPower += voter.VotingPower
		}
		// else {
		// It's OK. We include stray signatures (~votes for nil) to measure
		// voter availability.
		// }
	}

	if got, needed := talliedVotingPower, votingPowerNeeded; got <= needed {
		return ErrNotEnoughVotingPowerSigned{Got: got, Needed: needed}
	}

	return nil
}

// LIGHT CLIENT VERIFICATION METHODS

// VerifyCommitLight verifies +2/3 of the set had signed the given commit.
//
// This method is primarily used by the light client and does not check all the
// signatures.
func (voters *VoterSet) VerifyCommitLight(chainID string, blockID BlockID,
	height int64, commit *Commit) error {

	if voters.Size() != len(commit.Signatures) {
		return NewErrInvalidCommitSignatures(voters.Size(), len(commit.Signatures))
	}

	// Validate Height and BlockID.
	if height != commit.Height {
		return NewErrInvalidCommitHeight(height, commit.Height)
	}
	if !blockID.Equals(commit.BlockID) {
		return fmt.Errorf("invalid commit -- wrong block ID: want %v, got %v",
			blockID, commit.BlockID)
	}

	talliedVotingPower := int64(0)
	votingPowerNeeded := voters.TotalVotingPower() * 2 / 3
	for idx, commitSig := range commit.Signatures {
		// No need to verify absent or nil votes.
		if !commitSig.ForBlock() {
			continue
		}

		// The vals and commit have a 1-to-1 correspondance.
		// This means we don't need the voter address or to do any lookup.
		voter := voters.Voters[idx]

		// Validate signature.
		voteSignBytes := commit.VoteSignBytes(chainID, int32(idx))
		if !voter.PubKey.VerifySignature(voteSignBytes, commitSig.Signature) {
			return fmt.Errorf("wrong signature (#%d): %X", idx, commitSig.Signature)
		}

		talliedVotingPower += voter.VotingPower

		// return as soon as +2/3 of the signatures are verified
		if talliedVotingPower > votingPowerNeeded {
			return nil
		}
	}

	return ErrNotEnoughVotingPowerSigned{Got: talliedVotingPower, Needed: votingPowerNeeded}
}

// VerifyCommitLightTrusting verifies that trustLevel of the voter set signed
// this commit.
//
// NOTE the given voters do not necessarily correspond to the voter set
// for this commit, but there may be some intersection.
//
// This method is primarily used by the light client and does not check all the
// signatures.
func (voters *VoterSet) VerifyCommitLightTrusting(chainID string, commit *Commit, trustLevel tmmath.Fraction) error {
	// sanity check
	if trustLevel.Denominator == 0 {
		return errors.New("trustLevel has zero Denominator")
	}

	var (
		talliedVotingPower int64
		seenVoters         = make(map[int32]int, len(commit.Signatures)) // voter index -> commit index
	)

	// Safely calculate voting power needed.
	totalVotingPowerMulByNumerator, overflow := safeMul(voters.TotalVotingPower(), int64(trustLevel.Numerator))
	if overflow {
		return errors.New("int64 overflow while calculating voting power needed. please provide smaller trustLevel numerator")
	}
	votingPowerNeeded := totalVotingPowerMulByNumerator / int64(trustLevel.Denominator)

	for idx, commitSig := range commit.Signatures {
		// No need to verify absent or nil votes.
		if !commitSig.ForBlock() {
			continue
		}

		// We don't know the voters that committed this block, so we have to
		// check for each vote if its voter is already known.
		voterIdx, voter := voters.GetByAddress(commitSig.ValidatorAddress)

		if voter != nil {
			// check for double vote of voter on the same commit
			if firstIndex, ok := seenVoters[voterIdx]; ok {
				secondIndex := idx
				return fmt.Errorf("double vote from %v (%d and %d)", voter, firstIndex, secondIndex)
			}
			seenVoters[voterIdx] = idx

			// Validate signature.
			voteSignBytes := commit.VoteSignBytes(chainID, int32(idx))
			if !voter.PubKey.VerifySignature(voteSignBytes, commitSig.Signature) {
				return fmt.Errorf("wrong signature (#%d): %X", idx, commitSig.Signature)
			}

			talliedVotingPower += voter.VotingPower

			if talliedVotingPower > votingPowerNeeded {
				return nil
			}
		}
	}

	return ErrNotEnoughVotingPowerSigned{Got: talliedVotingPower, Needed: votingPowerNeeded}
}

//-----------------

// IsErrNotEnoughVotingPowerSigned returns true if err is
// ErrNotEnoughVotingPowerSigned.
func IsErrNotEnoughVotingPowerSigned(err error) bool {
	_, ok := errors.Cause(err).(ErrNotEnoughVotingPowerSigned)
	return ok
}

// ErrNotEnoughVotingPowerSigned is returned when not enough voters signed
// a commit.
type ErrNotEnoughVotingPowerSigned struct {
	Got    int64
	Needed int64
}

func (e ErrNotEnoughVotingPowerSigned) Error() string {
	return fmt.Sprintf("invalid commit -- insufficient voting power: got %d, needed more than %d", e.Got, e.Needed)
}

//----------------

// Iterate will run the given function over the set.
func (voters *VoterSet) Iterate(fn func(index int, val *Validator) bool) {
	for i, val := range voters.Voters {
		stop := fn(i, val)
		if stop {
			break
		}
	}
}

func (voters *VoterSet) String() string {
	return voters.StringIndented("")
}

// StringIndented returns an intended string representation of VoterSet.
func (voters *VoterSet) StringIndented(indent string) string {
	if voters == nil {
		return "nil-VoterSet"
	}
	var voterStrings []string
	voters.Iterate(func(index int, voter *Validator) bool {
		voterStrings = append(voterStrings, voter.String())
		return false
	})
	return fmt.Sprintf(`VoterSet{
%s  Voters:
%s    %v
%s}`,
		indent, indent, strings.Join(voterStrings, "\n"+indent+"    "),
		indent)

}

// ToProto converts VoterSet to protobuf
func (voters *VoterSet) ToProto() (*tmproto.VoterSet, error) {
	if voters.IsNilOrEmpty() {
		return &tmproto.VoterSet{}, nil // voter set should never be nil
	}

	vp := new(tmproto.VoterSet)
	valsProto := make([]*tmproto.Validator, len(voters.Voters))
	for i := 0; i < len(voters.Voters); i++ {
		valp, err := voters.Voters[i].ToProto()
		if err != nil {
			return nil, err
		}
		valsProto[i] = valp
	}
	vp.Validators = valsProto

	vp.TotalVotingPower = voters.totalVotingPower

	return vp, nil
}

// VoterSetFromProto sets a protobuf VoterSet to the given pointer.
// It returns an error if any of the validators from the set or the proposer
// is invalid
func VoterSetFromProto(vp *tmproto.VoterSet) (*VoterSet, error) {
	if vp == nil {
		return nil, errors.New("nil voter set") // voter set should never be nil, bigger issues are at play if empty
	}
	voters := new(VoterSet)

	valsProto := make([]*Validator, len(vp.Validators))
	for i := 0; i < len(vp.Validators); i++ {
		v, err := ValidatorFromProto(vp.Validators[i])
		if err != nil {
			return nil, err
		}
		valsProto[i] = v
	}
	voters.Voters = valsProto

	voters.totalVotingPower = vp.GetTotalVotingPower()

	return voters, voters.ValidateBasic()
}

func SelectVoter(validators *ValidatorSet, proofHash []byte) *VoterSet {
	// TODO: decide MaxVoters; make it to config
	if len(proofHash) == 0 || validators.Size() <= MaxVoters {
		// height 1 has voter set that is same to validator set
		result := &VoterSet{Voters: copyValidatorListShallow(validators.Validators), totalVotingPower: 0}
		result.updateTotalVotingPower()
		return result
	}

	seed := hashToSeed(proofHash)
	candidates := make([]tmrand.Candidate, len(validators.Validators))
	for i, val := range validators.Validators {
		candidates[i] = &candidate{idx: i, win: 0, val: val}
	}
	totalSampling := tmrand.RandomSamplingToMax(seed, candidates, MaxVoters, uint64(validators.TotalVotingPower()))
	voters := 0
	for _, candi := range candidates {
		if candi.(*candidate).win > 0 {
			voters++
		}
	}

	vals := make([]*Validator, voters)
	index := 0
	for _, candi := range candidates {
		if candi.(*candidate).win > 0 {
			vals[index] = &Validator{Address: candi.(*candidate).val.Address,
				PubKey: candi.(*candidate).val.PubKey,
				// VotingPower = TotalVotingPower * win / totalSampling : can be overflow
				VotingPower: validators.TotalVotingPower()/int64(totalSampling)*int64(candi.(*candidate).win) +
					int64(math.Ceil(float64(validators.TotalVotingPower()%int64(totalSampling))/float64(int64(totalSampling))*
						float64(candi.(*candidate).win)))}
			index++
		}
	}
	return NewVoterSet(vals)
}

// This should be used in only test
func ToVoterAll(validators *ValidatorSet) *VoterSet {
	return NewVoterSet(validators.Validators)
}

// candidate save simple validator data for selecting proposer
type candidate struct {
	idx int
	win uint64
	val *Validator
}

func (c *candidate) Priority() uint64 {
	// TODO Is it possible to have a negative VotingPower?
	if c.val.VotingPower < 0 {
		return 0
	}
	return uint64(c.val.VotingPower)
}

func (c *candidate) LessThan(other tmrand.Candidate) bool {
	o, ok := other.(*candidate)
	if !ok {
		panic("incompatible type")
	}
	return bytes.Compare(c.val.Address, o.val.Address) < 0
}

func (c *candidate) IncreaseWin() {
	c.win++
}

func hashToSeed(hash []byte) uint64 {
	for len(hash) < 8 {
		hash = append(hash, byte(0))
	}
	return binary.LittleEndian.Uint64(hash[:8])
}

// MakeRoundHash combines the VRF hash, block height, and round to create a hash value for each round. This value is
// used for random sampling of the Proposer.
func MakeRoundHash(proofHash []byte, height int64, round int32) []byte {
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b, uint64(height))
	binary.LittleEndian.PutUint64(b[8:], uint64(round))
	hash := tmhash.New()
	hash.Write(proofHash)
	hash.Write(b[:8])
	hash.Write(b[8:16])
	return hash.Sum(nil)
}

// RandVoterSet returns a randomized validator/voter set, useful for testing.
// NOTE: PrivValidator are in order.
// UNSTABLE
func RandVoterSet(numVoters int, votingPower int64) (*ValidatorSet, *VoterSet, []PrivValidator) {
	valz := make([]*Validator, numVoters)
	privValidators := make([]PrivValidator, numVoters)
	for i := 0; i < numVoters; i++ {
		val, privValidator := RandValidator(false, votingPower)
		valz[i] = val
		privValidators[i] = privValidator
	}
	vals := NewValidatorSet(valz)
	sort.Sort(PrivValidatorsByAddress(privValidators))
	return vals, SelectVoter(vals, []byte{}), privValidators
}