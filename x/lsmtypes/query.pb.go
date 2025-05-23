// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/lsmstaking/v1beta1/query.proto

package lsmtypes

import (
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// QueryValidatorsResponse is response type for the Query/Validators RPC method
type QueryValidatorsResponse struct {
	// validators contains all the queried validators.
	Validators []Validator `protobuf:"bytes,1,rep,name=validators,proto3" json:"validators"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryValidatorsResponse) Reset()         { *m = QueryValidatorsResponse{} }
func (m *QueryValidatorsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryValidatorsResponse) ProtoMessage()    {}
func (*QueryValidatorsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c06978e9e1a8aa4, []int{0}
}
func (m *QueryValidatorsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryValidatorsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryValidatorsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryValidatorsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryValidatorsResponse.Merge(m, src)
}
func (m *QueryValidatorsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryValidatorsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryValidatorsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryValidatorsResponse proto.InternalMessageInfo

func (m *QueryValidatorsResponse) GetValidators() []Validator {
	if m != nil {
		return m.Validators
	}
	return nil
}

func (m *QueryValidatorsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryValidatorsResponse)(nil), "cosmos.lsmstaking.v1beta1.QueryValidatorsResponse")
}

func init() {
	proto.RegisterFile("cosmos/lsmstaking/v1beta1/query.proto", fileDescriptor_0c06978e9e1a8aa4)
}

var fileDescriptor_0c06978e9e1a8aa4 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0x1b, 0x05, 0xc1, 0xec, 0xe4, 0x10, 0x9c, 0x43, 0xe2, 0x10, 0xc5, 0x31, 0x30, 0x61,
	0xf3, 0x1b, 0xec, 0xe2, 0x45, 0x50, 0x77, 0xf0, 0xe0, 0x2d, 0x9d, 0x21, 0x86, 0xb5, 0x79, 0x5d,
	0x93, 0x15, 0xe7, 0xa7, 0xf0, 0x63, 0x78, 0x14, 0xfc, 0x12, 0x3b, 0xee, 0xe8, 0x49, 0xa4, 0x3d,
	0xf8, 0x35, 0xa4, 0x6d, 0xda, 0xf5, 0xa2, 0x97, 0xf0, 0x78, 0xfc, 0xde, 0x2f, 0xef, 0xfd, 0xf1,
	0xd9, 0x14, 0x4c, 0x08, 0x86, 0x05, 0x26, 0x34, 0x96, 0xcf, 0x94, 0x96, 0x2c, 0x19, 0xfa, 0xc2,
	0xf2, 0x21, 0x9b, 0x2f, 0x44, 0xbc, 0xa4, 0x51, 0x0c, 0x16, 0xda, 0x87, 0x25, 0x46, 0x37, 0x18,
	0x75, 0x58, 0x77, 0x8f, 0x87, 0x4a, 0x03, 0x2b, 0xde, 0x92, 0xee, 0x0e, 0x9c, 0xd4, 0xe7, 0x46,
	0x94, 0x9a, 0x5a, 0x1a, 0x71, 0xa9, 0x34, 0xb7, 0x0a, 0xb4, 0x63, 0xff, 0x59, 0xc0, 0x2e, 0x23,
	0x61, 0x1c, 0xb6, 0x2f, 0x41, 0x42, 0x51, 0xb2, 0xbc, 0x72, 0xdd, 0x23, 0x09, 0x20, 0x03, 0xc1,
	0x78, 0xa4, 0x18, 0xd7, 0x1a, 0x6c, 0x61, 0x76, 0x33, 0x27, 0x1f, 0x08, 0x1f, 0xdc, 0xe5, 0xbf,
	0xdf, 0xf3, 0x40, 0x3d, 0x72, 0x0b, 0xb1, 0x99, 0x08, 0x13, 0x81, 0x36, 0xa2, 0x7d, 0x83, 0x71,
	0x52, 0x77, 0x3b, 0xa8, 0xb7, 0xdd, 0x6f, 0x8d, 0x4e, 0xe9, 0x9f, 0x57, 0xd2, 0x5a, 0x31, 0xde,
	0x5d, 0x7d, 0x1d, 0x7b, 0x6f, 0x3f, 0xef, 0x03, 0x34, 0x69, 0x28, 0xda, 0x57, 0x18, 0x6f, 0x6e,
	0xeb, 0x6c, 0xf5, 0x50, 0xbf, 0x35, 0x3a, 0xaf, 0x84, 0x79, 0x10, 0xb4, 0xcc, 0xb3, 0x12, 0xde,
	0x72, 0x29, 0xaa, 0x6d, 0x26, 0x8d, 0xd1, 0xf1, 0xf5, 0x2a, 0x25, 0x68, 0x9d, 0x12, 0xf4, 0x9d,
	0x12, 0xf4, 0x9a, 0x11, 0x6f, 0x9d, 0x11, 0xef, 0x33, 0x23, 0xde, 0xc3, 0x48, 0x2a, 0xfb, 0xb4,
	0xf0, 0xe9, 0x14, 0x42, 0x36, 0x5f, 0xa8, 0xe9, 0xcc, 0xa8, 0x20, 0x11, 0xf1, 0xc5, 0x0b, 0x68,
	0xd1, 0x6c, 0xb0, 0xe7, 0x3c, 0xcd, 0x22, 0x3d, 0x7f, 0xa7, 0x88, 0xe2, 0xf2, 0x37, 0x00, 0x00,
	0xff, 0xff, 0xd5, 0x4f, 0x93, 0xec, 0xe8, 0x01, 0x00, 0x00,
}

func (m *QueryValidatorsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryValidatorsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryValidatorsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Validators) > 0 {
		for iNdEx := len(m.Validators) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Validators[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryValidatorsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Validators) > 0 {
		for _, e := range m.Validators {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryValidatorsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryValidatorsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryValidatorsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validators", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validators = append(m.Validators, Validator{})
			if err := m.Validators[len(m.Validators)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
