// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ostracon/libs/bits/types.proto

package bits

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type BitArray struct {
	Bits  int64    `protobuf:"varint,1,opt,name=bits,proto3" json:"bits,omitempty"`
	Elems []uint64 `protobuf:"varint,2,rep,packed,name=elems,proto3" json:"elems,omitempty"`
}

func (m *BitArray) Reset()         { *m = BitArray{} }
func (m *BitArray) String() string { return proto.CompactTextString(m) }
func (*BitArray) ProtoMessage()    {}
func (*BitArray) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbd48d316ef10ed8, []int{0}
}
func (m *BitArray) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BitArray) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BitArray.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BitArray) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BitArray.Merge(m, src)
}
func (m *BitArray) XXX_Size() int {
	return m.Size()
}
func (m *BitArray) XXX_DiscardUnknown() {
	xxx_messageInfo_BitArray.DiscardUnknown(m)
}

var xxx_messageInfo_BitArray proto.InternalMessageInfo

func (m *BitArray) GetBits() int64 {
	if m != nil {
		return m.Bits
	}
	return 0
}

func (m *BitArray) GetElems() []uint64 {
	if m != nil {
		return m.Elems
	}
	return nil
}

func init() {
	proto.RegisterType((*BitArray)(nil), "ostracon.libs.bits.BitArray")
}

func init() { proto.RegisterFile("ostracon/libs/bits/types.proto", fileDescriptor_cbd48d316ef10ed8) }

var fileDescriptor_cbd48d316ef10ed8 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcb, 0x2f, 0x2e, 0x29,
	0x4a, 0x4c, 0xce, 0xcf, 0xd3, 0xcf, 0xc9, 0x4c, 0x2a, 0xd6, 0x4f, 0xca, 0x2c, 0x29, 0xd6, 0x2f,
	0xa9, 0x2c, 0x48, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x82, 0xc9, 0xeb, 0x81,
	0xe4, 0xf5, 0x40, 0xf2, 0x4a, 0x26, 0x5c, 0x1c, 0x4e, 0x99, 0x25, 0x8e, 0x45, 0x45, 0x89, 0x95,
	0x42, 0x42, 0x5c, 0x2c, 0x20, 0x31, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xe6, 0x20, 0x30, 0x5b, 0x48,
	0x84, 0x8b, 0x35, 0x35, 0x27, 0x35, 0xb7, 0x58, 0x82, 0x49, 0x81, 0x59, 0x83, 0x25, 0x08, 0xc2,
	0x71, 0xf2, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27,
	0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xc3, 0xf4, 0xcc,
	0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0x9c, 0xcc, 0xbc, 0x54, 0x7d, 0xb8, 0x9b,
	0xc0, 0x6e, 0xd0, 0xc7, 0x74, 0x62, 0x12, 0x1b, 0x58, 0xc6, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff,
	0x73, 0x45, 0xe0, 0x5d, 0xbf, 0x00, 0x00, 0x00,
}

func (m *BitArray) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BitArray) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BitArray) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Elems) > 0 {
		dAtA2 := make([]byte, len(m.Elems)*10)
		var j1 int
		for _, num := range m.Elems {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintTypes(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x12
	}
	if m.Bits != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Bits))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BitArray) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Bits != 0 {
		n += 1 + sovTypes(uint64(m.Bits))
	}
	if len(m.Elems) > 0 {
		l = 0
		for _, e := range m.Elems {
			l += sovTypes(uint64(e))
		}
		n += 1 + sovTypes(uint64(l)) + l
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BitArray) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BitArray: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BitArray: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bits", wireType)
			}
			m.Bits = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Bits |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Elems = append(m.Elems, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTypes
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTypes
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Elems) == 0 {
					m.Elems = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTypes
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Elems = append(m.Elems, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Elems", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)