package evidence

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	dbm "github.com/tendermint/tm-db"
)

func TestMain(m *testing.M) {
	types.RegisterMockEvidences(cdc)

	code := m.Run()
	os.Exit(code)
}

func initializeValidatorState(valAddr []byte, height int64) dbm.DB {
	stateDB := dbm.NewMemDB()

	// create validator set and state
	vals := []*types.Validator{
		{Address: valAddr, VotingPower: 1},
	}
	state := sm.State{
		LastBlockHeight:             0,
		LastBlockTime:               tmtime.Now(),
		Validators:                  types.NewRandomValidatorSet(vals, types.MakeRoundHash([]byte{}, 1, 0)),
		NextValidators:              types.NewRandomValidatorSet(vals, types.MakeRoundHash([]byte{}, 2, 0)),
		LastHeightValidatorsChanged: 1,
		ConsensusParams: types.ConsensusParams{
			Evidence: types.EvidenceParams{
				MaxAgeNumBlocks: 10000,
				MaxAgeDuration:  48 * time.Hour,
			},
		},
	}

	// save all states up to height
	for i := int64(0); i < height; i++ {
		state.LastBlockHeight = i
		sm.SaveState(stateDB, state)
	}

	return stateDB
}

func TestEvidencePool(t *testing.T) {

	var (
		valAddr      = []byte("val1")
		height       = int64(5)
		stateDB      = initializeValidatorState(valAddr, height)
		evidenceDB   = dbm.NewMemDB()
		pool         = NewPool(stateDB, evidenceDB)
		evidenceTime = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	goodEvidence := types.NewMockEvidence(height, time.Now(), 0, valAddr)
	badEvidence := types.NewMockEvidence(height, evidenceTime, 0, valAddr)

	// bad evidence
	err := pool.AddEvidence(badEvidence)
	assert.NotNil(t, err)
	// err: evidence created at 2019-01-01 00:00:00 +0000 UTC has expired. Evidence can not be older than: ...

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-pool.EvidenceWaitChan()
		wg.Done()
	}()

	err = pool.AddEvidence(goodEvidence)
	assert.Nil(t, err)
	wg.Wait()

	assert.Equal(t, 1, pool.evidenceList.Len())

	// if we send it again, it shouldnt change the size
	err = pool.AddEvidence(goodEvidence)
	assert.Nil(t, err)
	assert.Equal(t, 1, pool.evidenceList.Len())
}

func TestEvidencePoolIsCommitted(t *testing.T) {
	// Initialization:
	var (
		valAddr       = []byte("validator_address")
		height        = int64(42)
		lastBlockTime = time.Now()
		stateDB       = initializeValidatorState(valAddr, height)
		evidenceDB    = dbm.NewMemDB()
		pool          = NewPool(stateDB, evidenceDB)
	)

	// evidence not seen yet:
	evidence := types.NewMockEvidence(height, time.Now(), 0, valAddr)
	assert.False(t, pool.IsCommitted(evidence))

	// evidence seen but not yet committed:
	assert.NoError(t, pool.AddEvidence(evidence))
	assert.False(t, pool.IsCommitted(evidence))

	// evidence seen and committed:
	pool.MarkEvidenceAsCommitted(height, lastBlockTime, []types.Evidence{evidence})
	assert.True(t, pool.IsCommitted(evidence))
}

func TestAddEvidence(t *testing.T) {

	var (
		valAddr      = []byte("val1")
		height       = int64(100002)
		stateDB      = initializeValidatorState(valAddr, height)
		evidenceDB   = dbm.NewMemDB()
		pool         = NewPool(stateDB, evidenceDB)
		evidenceTime = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	testCases := []struct {
		evHeight      int64
		evTime        time.Time
		expErr        bool
		evDescription string
	}{
		{height, time.Now(), false, "valid evidence"},
		{height, evidenceTime, true, "evidence created at 2019-01-01 00:00:00 +0000 UTC has expired"},
		{int64(1), time.Now(), true, "evidence from height 1 is too old"},
		{int64(1), evidenceTime, true,
			"evidence from height 1 is too old & evidence created at 2019-01-01 00:00:00 +0000 UTC has expired"},
	}

	for _, tc := range testCases {
		tc := tc
		ev := types.NewMockEvidence(tc.evHeight, tc.evTime, 0, valAddr)
		err := pool.AddEvidence(ev)
		if tc.expErr {
			assert.Error(t, err)
		}
	}
}
