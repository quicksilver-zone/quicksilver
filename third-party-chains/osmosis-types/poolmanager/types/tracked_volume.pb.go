// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: osmosis/poolmanager/v1beta1/tracked_volume.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type TrackedVolume struct {
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func (m *TrackedVolume) Reset()         { *m = TrackedVolume{} }
func (m *TrackedVolume) String() string { return proto.CompactTextString(m) }
func (*TrackedVolume) ProtoMessage()    {}
func (*TrackedVolume) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a2e3e91de3baf1a, []int{0}
}
func (m *TrackedVolume) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TrackedVolume) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TrackedVolume.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TrackedVolume) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TrackedVolume.Merge(m, src)
}
func (m *TrackedVolume) XXX_Size() int {
	return m.Size()
}
func (m *TrackedVolume) XXX_DiscardUnknown() {
	xxx_messageInfo_TrackedVolume.DiscardUnknown(m)
}

var xxx_messageInfo_TrackedVolume proto.InternalMessageInfo

func (m *TrackedVolume) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func init() {
	proto.RegisterType((*TrackedVolume)(nil), "osmosis.poolmanager.v1beta1.TrackedVolume")
}

func init() {
	proto.RegisterFile("osmosis/poolmanager/v1beta1/tracked_volume.proto", fileDescriptor_0a2e3e91de3baf1a)
}

var fileDescriptor_0a2e3e91de3baf1a = []byte{
	// 275 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x31, 0x4e, 0xc3, 0x30,
	0x14, 0x86, 0x13, 0x21, 0x75, 0x28, 0x62, 0xa9, 0x18, 0xa0, 0x48, 0x2e, 0x62, 0xea, 0x12, 0xbb,
	0x85, 0x1b, 0x94, 0x1b, 0x20, 0xc4, 0xd0, 0x05, 0x39, 0x8e, 0x95, 0x58, 0x49, 0xfc, 0x82, 0xfd,
	0x12, 0xa9, 0x9c, 0x82, 0x73, 0x70, 0x92, 0x8e, 0x1d, 0x99, 0x00, 0x25, 0x17, 0x41, 0xb5, 0x0d,
	0x0a, 0x93, 0x9f, 0x7e, 0xeb, 0xfb, 0x7e, 0xfb, 0x4d, 0x57, 0x60, 0x6b, 0xb0, 0xca, 0xb2, 0x06,
	0xa0, 0xaa, 0xb9, 0xe6, 0xb9, 0x34, 0xac, 0x5b, 0xa7, 0x12, 0xf9, 0x9a, 0xa1, 0xe1, 0xa2, 0x94,
	0xd9, 0x73, 0x07, 0x55, 0x5b, 0x4b, 0xda, 0x18, 0x40, 0x98, 0x5d, 0x05, 0x82, 0x8e, 0x08, 0x1a,
	0x88, 0x39, 0x11, 0xee, 0x96, 0xa5, 0xdc, 0xca, 0x3f, 0x8d, 0x00, 0xa5, 0x3d, 0x3c, 0x3f, 0xcf,
	0x21, 0x07, 0x37, 0xb2, 0xe3, 0xe4, 0xd3, 0x1b, 0x9c, 0x9e, 0x3d, 0xfa, 0xaa, 0x27, 0xd7, 0x34,
	0x13, 0xd3, 0x09, 0xaf, 0xa1, 0xd5, 0x78, 0x11, 0x5f, 0x9f, 0x2c, 0x4f, 0x6f, 0x2f, 0xa9, 0xf7,
	0xd2, 0xa3, 0xf7, 0xb7, 0x8c, 0xde, 0x83, 0xd2, 0x9b, 0xd5, 0xfe, 0x73, 0x11, 0xbd, 0x7f, 0x2d,
	0x96, 0xb9, 0xc2, 0xa2, 0x4d, 0xa9, 0x80, 0x9a, 0x85, 0x47, 0xf8, 0x23, 0xb1, 0x59, 0xc9, 0x70,
	0xd7, 0x48, 0xeb, 0x00, 0xfb, 0x10, 0xd4, 0x1b, 0xdc, 0xf7, 0x24, 0x3e, 0xf4, 0x24, 0xfe, 0xee,
	0x49, 0xfc, 0x36, 0x90, 0xe8, 0x30, 0x90, 0xe8, 0x63, 0x20, 0xd1, 0x76, 0x3b, 0x72, 0xbd, 0xb4,
	0x4a, 0x94, 0x56, 0x55, 0x9d, 0x34, 0xc9, 0x2b, 0x68, 0x39, 0x0e, 0x18, 0x16, 0xca, 0x64, 0x49,
	0xc3, 0x0d, 0xee, 0x12, 0x51, 0x70, 0xa5, 0x2d, 0x0b, 0xdb, 0x49, 0x5c, 0xe3, 0xbf, 0xad, 0xba,
	0x24, 0x9d, 0xb8, 0x2f, 0xdf, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5a, 0x93, 0x84, 0x52, 0x79,
	0x01, 0x00, 0x00,
}

func (m *TrackedVolume) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TrackedVolume) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TrackedVolume) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrackedVolume(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTrackedVolume(dAtA []byte, offset int, v uint64) int {
	offset -= sovTrackedVolume(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TrackedVolume) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovTrackedVolume(uint64(l))
		}
	}
	return n
}

func sovTrackedVolume(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTrackedVolume(x uint64) (n int) {
	return sovTrackedVolume(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TrackedVolume) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrackedVolume
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
			return fmt.Errorf("proto: TrackedVolume: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TrackedVolume: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrackedVolume
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrackedVolume
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrackedVolume
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrackedVolume(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrackedVolume
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
func skipTrackedVolume(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTrackedVolume
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
					return 0, ErrIntOverflowTrackedVolume
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
					return 0, ErrIntOverflowTrackedVolume
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
				return 0, ErrInvalidLengthTrackedVolume
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTrackedVolume
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTrackedVolume
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTrackedVolume        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTrackedVolume          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTrackedVolume = fmt.Errorf("proto: unexpected end of group")
)
