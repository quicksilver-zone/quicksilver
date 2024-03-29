// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: quicksilver/lsm-types/v1/validator.proto

package lsmtypes

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	types1 "github.com/cosmos/cosmos-sdk/x/staking/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Validator defines a validator, together with the total amount of the
// Validator's bond shares and their exchange rate to coins. Slashing results in
// a decrease in the exchange rate, allowing correct calculation of future
// undelegations without iterating over delegators. When coins are delegated to
// this validator, the validator is credited with a delegation whose number of
// bond shares is based on the amount of coins delegated divided by the current
// exchange rate. Voting power can be calculated as total bonded shares
// multiplied by exchange rate.
type Validator struct {
	// operator_address defines the address of the validator's operator; bech
	// encoded in JSON.
	OperatorAddress string `protobuf:"bytes,1,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty" yaml:"operator_address"`
	// consensus_pubkey is the consensus public key of the validator, as a
	// Protobuf Any.
	ConsensusPubkey *types.Any `protobuf:"bytes,2,opt,name=consensus_pubkey,json=consensusPubkey,proto3" json:"consensus_pubkey,omitempty" yaml:"consensus_pubkey"`
	// jailed defined whether the validator has been jailed from bonded status or
	// not.
	Jailed bool `protobuf:"varint,3,opt,name=jailed,proto3" json:"jailed,omitempty"`
	// status is the validator status (bonded/unbonding/unbonded).
	Status types1.BondStatus `protobuf:"varint,4,opt,name=status,proto3,enum=cosmos.staking.v1beta1.BondStatus" json:"status,omitempty"`
	// tokens define the delegated tokens (incl. self-delegation).
	Tokens github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=tokens,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"tokens"`
	// delegator_shares defines total shares issued to a validator's delegators.
	DelegatorShares github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=delegator_shares,json=delegatorShares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegator_shares" yaml:"delegator_shares"`
	// description defines the description terms for the validator.
	Description types1.Description `protobuf:"bytes,7,opt,name=description,proto3" json:"description"`
	// unbonding_height defines, if unbonding, the height at which this validator
	// has begun unbonding.
	UnbondingHeight int64 `protobuf:"varint,8,opt,name=unbonding_height,json=unbondingHeight,proto3" json:"unbonding_height,omitempty" yaml:"unbonding_height"`
	// unbonding_time defines, if unbonding, the min time for the validator to
	// complete unbonding.
	UnbondingTime time.Time `protobuf:"bytes,9,opt,name=unbonding_time,json=unbondingTime,proto3,stdtime" json:"unbonding_time" yaml:"unbonding_time"`
	// commission defines the commission parameters.
	Commission types1.Commission `protobuf:"bytes,10,opt,name=commission,proto3" json:"commission"`
	// Deprecated: This field has been deprecated with LSM in favor of the
	// validator bond
	MinSelfDelegation github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,11,opt,name=min_self_delegation,json=minSelfDelegation,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_self_delegation" yaml:"min_self_delegation"` // Deprecated: Do not use.
	// strictly positive if this validator's unbonding has been stopped by
	// external modules
	UnbondingOnHoldRefCount int64 `protobuf:"varint,12,opt,name=unbonding_on_hold_ref_count,json=unbondingOnHoldRefCount,proto3" json:"unbonding_on_hold_ref_count,omitempty"`
	// list of unbonding ids, each uniquely identifing an unbonding of this
	// validator
	UnbondingIds []uint64 `protobuf:"varint,13,rep,packed,name=unbonding_ids,json=unbondingIds,proto3" json:"unbonding_ids,omitempty"`
	// Number of shares self bonded from the validator
	ValidatorBondShares github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,14,opt,name=validator_bond_shares,json=validatorBondShares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"validator_bond_shares" yaml:"validator_bond_shares"`
	// Number of shares either tokenized or owned by a liquid staking provider
	LiquidShares github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,15,opt,name=liquid_shares,json=liquidShares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquid_shares" yaml:"liquid_shares"`
}

func (m *Validator) Reset()      { *m = Validator{} }
func (*Validator) ProtoMessage() {}
func (*Validator) Descriptor() ([]byte, []int) {
	return fileDescriptor_93ed251ae1dc13ed, []int{0}
}
func (m *Validator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Validator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Validator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Validator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Validator.Merge(m, src)
}
func (m *Validator) XXX_Size() int {
	return m.Size()
}
func (m *Validator) XXX_DiscardUnknown() {
	xxx_messageInfo_Validator.DiscardUnknown(m)
}

var xxx_messageInfo_Validator proto.InternalMessageInfo

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
	return fileDescriptor_93ed251ae1dc13ed, []int{1}
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
	proto.RegisterType((*Validator)(nil), "cosmos.lsmstaking.v1beta1.Validator")
	proto.RegisterType((*QueryValidatorsResponse)(nil), "cosmos.lsmstaking.v1beta1.QueryValidatorsResponse")
}

func init() {
	proto.RegisterFile("quicksilver/lsm-types/v1/validator.proto", fileDescriptor_93ed251ae1dc13ed)
}

var fileDescriptor_93ed251ae1dc13ed = []byte{
	// 870 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x3f, 0x6f, 0xe4, 0x44,
	0x14, 0x5f, 0x93, 0xb0, 0x5c, 0x26, 0x7f, 0xf6, 0xf0, 0xe5, 0x88, 0xc9, 0xa1, 0xf5, 0xe2, 0x3b,
	0xc1, 0x0a, 0x29, 0xb6, 0x92, 0xeb, 0x22, 0x9a, 0xdb, 0x8b, 0x42, 0xc2, 0x49, 0x10, 0x26, 0x88,
	0x82, 0xc6, 0xf2, 0xda, 0xb3, 0xde, 0x61, 0xed, 0x19, 0xc7, 0x33, 0x8e, 0x30, 0x05, 0x05, 0x34,
	0x94, 0x57, 0x52, 0xa6, 0xe4, 0x03, 0xdc, 0x87, 0x38, 0x51, 0x5d, 0x89, 0x28, 0x16, 0x94, 0x34,
	0xd4, 0xf9, 0x04, 0x68, 0xc6, 0xe3, 0xb1, 0x59, 0x92, 0x22, 0xd5, 0xee, 0xbc, 0xf7, 0x7b, 0xbf,
	0xf7, 0x67, 0x7e, 0x9e, 0x07, 0x86, 0x67, 0x05, 0x0e, 0x67, 0x0c, 0x27, 0xe7, 0x28, 0xf7, 0x12,
	0x96, 0xee, 0xf0, 0x32, 0x43, 0xcc, 0x3b, 0xdf, 0xf5, 0xce, 0x83, 0x04, 0x47, 0x01, 0xa7, 0xb9,
	0x9b, 0xe5, 0x94, 0x53, 0xf3, 0xfd, 0x90, 0xb2, 0x94, 0x32, 0x37, 0x61, 0x29, 0xe3, 0xc1, 0x0c,
	0x93, 0xd8, 0x3d, 0xdf, 0x1d, 0x23, 0x1e, 0xec, 0x6e, 0x7f, 0x52, 0xb9, 0xbc, 0x71, 0xc0, 0x90,
	0x77, 0x56, 0xa0, 0xbc, 0xf4, 0x94, 0xcb, 0xcb, 0x82, 0x18, 0x93, 0x80, 0x63, 0x4a, 0x2a, 0x9a,
	0xed, 0x7e, 0x1b, 0x5b, 0xa3, 0x42, 0x8a, 0x6b, 0xff, 0x13, 0xe5, 0x57, 0x39, 0x34, 0xa4, 0xce,
	0x59, 0xa1, 0x54, 0x31, 0xbe, 0x3c, 0x79, 0xaa, 0xb2, 0xca, 0xb5, 0x19, 0xd3, 0x98, 0x56, 0x76,
	0xf1, 0xaf, 0x0e, 0x88, 0x29, 0x8d, 0x13, 0xe4, 0xc9, 0xd3, 0xb8, 0x98, 0x78, 0x01, 0x29, 0x95,
	0xcb, 0x5e, 0x74, 0x71, 0x9c, 0x22, 0xc6, 0x83, 0x34, 0xab, 0x00, 0xce, 0xcf, 0x00, 0xac, 0x7c,
	0x53, 0x4f, 0xc3, 0x3c, 0x04, 0xf7, 0x69, 0x86, 0x72, 0xf1, 0xdf, 0x0f, 0xa2, 0x28, 0x47, 0x8c,
	0x59, 0xc6, 0xc0, 0x18, 0xae, 0x8c, 0x1e, 0x5d, 0xcf, 0xed, 0xad, 0x32, 0x48, 0x93, 0x7d, 0x67,
	0x11, 0xe1, 0xc0, 0x5e, 0x6d, 0x7a, 0x56, 0x59, 0x4c, 0x0e, 0xee, 0x87, 0x94, 0x30, 0x44, 0x58,
	0xc1, 0xfc, 0xac, 0x18, 0xcf, 0x50, 0x69, 0xbd, 0x35, 0x30, 0x86, 0xab, 0x7b, 0x9b, 0x6e, 0x55,
	0x91, 0x5b, 0x57, 0xe4, 0x3e, 0x23, 0xe5, 0xe8, 0x69, 0xc3, 0xbe, 0x18, 0xe7, 0xfc, 0xfe, 0x6a,
	0x67, 0x53, 0x0d, 0x21, 0xcc, 0xcb, 0x8c, 0x53, 0xf7, 0xa4, 0x18, 0xbf, 0x40, 0x25, 0xec, 0x69,
	0xe8, 0x89, 0x44, 0x9a, 0xef, 0x81, 0xee, 0x77, 0x01, 0x4e, 0x50, 0x64, 0x2d, 0x0d, 0x8c, 0xe1,
	0x3d, 0xa8, 0x4e, 0xe6, 0x3e, 0xe8, 0x32, 0x1e, 0xf0, 0x82, 0x59, 0xcb, 0x03, 0x63, 0xb8, 0xb1,
	0xe7, 0xb8, 0x8a, 0x6f, 0xe1, 0xae, 0xdd, 0x11, 0x25, 0xd1, 0xa9, 0x44, 0x42, 0x15, 0x61, 0x1e,
	0x82, 0x2e, 0xa7, 0x33, 0x44, 0x98, 0xf5, 0xb6, 0x9c, 0x83, 0xfb, 0x7a, 0x6e, 0x77, 0xfe, 0x9c,
	0xdb, 0x1f, 0xc5, 0x98, 0x4f, 0x8b, 0xb1, 0x1b, 0xd2, 0x54, 0x5d, 0x91, 0xfa, 0xd9, 0x61, 0xd1,
	0xcc, 0x93, 0x42, 0x73, 0x8f, 0x09, 0x87, 0x2a, 0x5a, 0x4c, 0x24, 0x42, 0x09, 0x8a, 0xe5, 0xe0,
	0xd8, 0x34, 0xc8, 0x11, 0xb3, 0xba, 0x92, 0xf1, 0xf8, 0x0e, 0x8c, 0x07, 0x28, 0x6c, 0x26, 0xb5,
	0xc8, 0xe7, 0xc0, 0x9e, 0x36, 0x9d, 0x4a, 0x8b, 0xf9, 0x02, 0xac, 0x46, 0x88, 0x85, 0x39, 0xce,
	0x84, 0x4a, 0xad, 0x77, 0xe4, 0x15, 0x3c, 0xbe, 0xad, 0xfd, 0x83, 0x06, 0x3a, 0x5a, 0x16, 0x55,
	0xc1, 0x76, 0xb4, 0x10, 0x47, 0x41, 0xc6, 0x94, 0x44, 0x98, 0xc4, 0xfe, 0x14, 0xe1, 0x78, 0xca,
	0xad, 0x7b, 0x03, 0x63, 0xb8, 0xd4, 0x16, 0xc7, 0x22, 0xc2, 0x81, 0x3d, 0x6d, 0x3a, 0x92, 0x16,
	0x33, 0x02, 0x1b, 0x0d, 0x4a, 0xe8, 0xd1, 0x5a, 0x91, 0x75, 0x6d, 0xff, 0x4f, 0x1a, 0x5f, 0xd7,
	0x62, 0x1d, 0x7d, 0x28, 0xca, 0xb9, 0x9e, 0xdb, 0x0f, 0x17, 0xb3, 0x88, 0x78, 0xe7, 0xe5, 0x5f,
	0xb6, 0x01, 0xd7, 0xb5, 0x51, 0x84, 0x99, 0x47, 0x00, 0x84, 0x34, 0x4d, 0x31, 0x63, 0xa2, 0x73,
	0x20, 0x33, 0xdc, 0x7a, 0xf1, 0xcf, 0x35, 0x52, 0x35, 0xde, 0x8a, 0x35, 0x7f, 0x04, 0x0f, 0x52,
	0x4c, 0x7c, 0x86, 0x92, 0x89, 0xaf, 0x06, 0x2c, 0x28, 0x57, 0xe5, 0xed, 0x7d, 0x71, 0x37, 0x3d,
	0x5c, 0xcf, 0xed, 0xed, 0xaa, 0x85, 0x1b, 0x28, 0x1d, 0xcb, 0x80, 0xef, 0xa6, 0x98, 0x9c, 0xa2,
	0x64, 0x72, 0xa0, 0xad, 0xe6, 0xa7, 0xe0, 0x51, 0xd3, 0x2f, 0x25, 0xfe, 0x94, 0x26, 0x91, 0x9f,
	0xa3, 0x89, 0x1f, 0xd2, 0x82, 0x70, 0x6b, 0x4d, 0x5c, 0x01, 0xdc, 0xd2, 0x90, 0x2f, 0xc9, 0x11,
	0x4d, 0x22, 0x88, 0x26, 0xcf, 0x85, 0xdb, 0x7c, 0x0c, 0x9a, 0xc1, 0xf8, 0x38, 0x62, 0xd6, 0xfa,
	0x60, 0x69, 0xb8, 0x0c, 0xd7, 0xb4, 0xf1, 0x38, 0x62, 0xe6, 0x4f, 0x06, 0x78, 0xa8, 0xdf, 0x44,
	0x5f, 0x38, 0x6a, 0x8d, 0x6e, 0xdc, 0xb9, 0xcb, 0x4a, 0xa3, 0x1f, 0x54, 0x5d, 0xde, 0x48, 0xea,
	0xc0, 0x07, 0xda, 0x2e, 0x3f, 0xb9, 0x4a, 0xac, 0x33, 0xb0, 0x9e, 0xe0, 0xb3, 0x02, 0xeb, 0xdc,
	0x3d, 0x99, 0xfb, 0xf0, 0xce, 0xb9, 0x37, 0xab, 0xdc, 0xff, 0x21, 0x73, 0xe0, 0x5a, 0x75, 0xae,
	0x92, 0xed, 0xaf, 0xfd, 0x72, 0x61, 0x77, 0x7e, 0xbd, 0xb0, 0x3b, 0xff, 0x5c, 0xd8, 0x1d, 0xe7,
	0x95, 0x01, 0xb6, 0xbe, 0x12, 0x6f, 0xbb, 0x7e, 0x0a, 0x19, 0x44, 0x2c, 0x13, 0xef, 0x8b, 0xf9,
	0x39, 0x00, 0xba, 0x5a, 0xf1, 0x1a, 0x2e, 0x0d, 0x57, 0xf7, 0x9e, 0xb8, 0xb7, 0x2e, 0x0c, 0x57,
	0x53, 0xd4, 0x52, 0x6a, 0xa2, 0xcd, 0xcf, 0x00, 0x68, 0x96, 0x86, 0x7a, 0x11, 0x3f, 0xae, 0xb9,
	0xc4, 0xd6, 0x70, 0xe5, 0x86, 0xd1, 0x5c, 0x27, 0x41, 0x8c, 0xea, 0x42, 0x60, 0x2b, 0x74, 0x7f,
	0x59, 0x94, 0x3d, 0x3a, 0xf9, 0xed, 0xb2, 0x6f, 0xbc, 0xbe, 0xec, 0x1b, 0x6f, 0x2e, 0xfb, 0xc6,
	0xdf, 0x97, 0x7d, 0xe3, 0xe5, 0x55, 0xbf, 0xf3, 0xe6, 0xaa, 0xdf, 0xf9, 0xe3, 0xaa, 0xdf, 0xf9,
	0x76, 0xaf, 0x35, 0xb0, 0xd6, 0x26, 0xdc, 0xf9, 0x81, 0x12, 0xd4, 0x36, 0x78, 0xdf, 0x8b, 0xe5,
	0x28, 0x07, 0x38, 0xee, 0xca, 0x6f, 0xef, 0xe9, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x76, 0xdf,
	0x58, 0x7a, 0x3b, 0x07, 0x00, 0x00,
}

func (m *Validator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Validator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Validator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.LiquidShares.Size()
		i -= size
		if _, err := m.LiquidShares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintValidator(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x7a
	{
		size := m.ValidatorBondShares.Size()
		i -= size
		if _, err := m.ValidatorBondShares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintValidator(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x72
	if len(m.UnbondingIds) > 0 {
		dAtA2 := make([]byte, len(m.UnbondingIds)*10)
		var j1 int
		for _, num := range m.UnbondingIds {
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
		i = encodeVarintValidator(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x6a
	}
	if m.UnbondingOnHoldRefCount != 0 {
		i = encodeVarintValidator(dAtA, i, uint64(m.UnbondingOnHoldRefCount))
		i--
		dAtA[i] = 0x60
	}
	{
		size := m.MinSelfDelegation.Size()
		i -= size
		if _, err := m.MinSelfDelegation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintValidator(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	{
		size, err := m.Commission.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintValidator(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	n4, err4 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UnbondingTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.UnbondingTime):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintValidator(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x4a
	if m.UnbondingHeight != 0 {
		i = encodeVarintValidator(dAtA, i, uint64(m.UnbondingHeight))
		i--
		dAtA[i] = 0x40
	}
	{
		size, err := m.Description.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintValidator(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.DelegatorShares.Size()
		i -= size
		if _, err := m.DelegatorShares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintValidator(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.Tokens.Size()
		i -= size
		if _, err := m.Tokens.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintValidator(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.Status != 0 {
		i = encodeVarintValidator(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x20
	}
	if m.Jailed {
		i--
		if m.Jailed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.ConsensusPubkey != nil {
		{
			size, err := m.ConsensusPubkey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintValidator(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.OperatorAddress) > 0 {
		i -= len(m.OperatorAddress)
		copy(dAtA[i:], m.OperatorAddress)
		i = encodeVarintValidator(dAtA, i, uint64(len(m.OperatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
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
			i = encodeVarintValidator(dAtA, i, uint64(size))
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
				i = encodeVarintValidator(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintValidator(dAtA []byte, offset int, v uint64) int {
	offset -= sovValidator(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Validator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OperatorAddress)
	if l > 0 {
		n += 1 + l + sovValidator(uint64(l))
	}
	if m.ConsensusPubkey != nil {
		l = m.ConsensusPubkey.Size()
		n += 1 + l + sovValidator(uint64(l))
	}
	if m.Jailed {
		n += 2
	}
	if m.Status != 0 {
		n += 1 + sovValidator(uint64(m.Status))
	}
	l = m.Tokens.Size()
	n += 1 + l + sovValidator(uint64(l))
	l = m.DelegatorShares.Size()
	n += 1 + l + sovValidator(uint64(l))
	l = m.Description.Size()
	n += 1 + l + sovValidator(uint64(l))
	if m.UnbondingHeight != 0 {
		n += 1 + sovValidator(uint64(m.UnbondingHeight))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UnbondingTime)
	n += 1 + l + sovValidator(uint64(l))
	l = m.Commission.Size()
	n += 1 + l + sovValidator(uint64(l))
	l = m.MinSelfDelegation.Size()
	n += 1 + l + sovValidator(uint64(l))
	if m.UnbondingOnHoldRefCount != 0 {
		n += 1 + sovValidator(uint64(m.UnbondingOnHoldRefCount))
	}
	if len(m.UnbondingIds) > 0 {
		l = 0
		for _, e := range m.UnbondingIds {
			l += sovValidator(uint64(e))
		}
		n += 1 + sovValidator(uint64(l)) + l
	}
	l = m.ValidatorBondShares.Size()
	n += 1 + l + sovValidator(uint64(l))
	l = m.LiquidShares.Size()
	n += 1 + l + sovValidator(uint64(l))
	return n
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
			n += 1 + l + sovValidator(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovValidator(uint64(l))
	}
	return n
}

func sovValidator(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozValidator(x uint64) (n int) {
	return sovValidator(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Validator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowValidator
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
			return fmt.Errorf("proto: Validator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Validator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OperatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OperatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsensusPubkey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
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
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ConsensusPubkey == nil {
				m.ConsensusPubkey = &types.Any{}
			}
			if err := m.ConsensusPubkey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Jailed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Jailed = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= types1.BondStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tokens", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Tokens.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorShares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DelegatorShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
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
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Description.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingHeight", wireType)
			}
			m.UnbondingHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnbondingHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
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
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UnbondingTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Commission", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
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
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Commission.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSelfDelegation", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinSelfDelegation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingOnHoldRefCount", wireType)
			}
			m.UnbondingOnHoldRefCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnbondingOnHoldRefCount |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowValidator
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
				m.UnbondingIds = append(m.UnbondingIds, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowValidator
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
					return ErrInvalidLengthValidator
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthValidator
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
				if elementCount != 0 && len(m.UnbondingIds) == 0 {
					m.UnbondingIds = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowValidator
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
					m.UnbondingIds = append(m.UnbondingIds, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingIds", wireType)
			}
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorBondShares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ValidatorBondShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidShares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipValidator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthValidator
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
func (m *QueryValidatorsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowValidator
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
					return ErrIntOverflowValidator
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
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
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
					return ErrIntOverflowValidator
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
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
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
			skippy, err := skipValidator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthValidator
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
func skipValidator(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowValidator
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
					return 0, ErrIntOverflowValidator
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
					return 0, ErrIntOverflowValidator
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
				return 0, ErrInvalidLengthValidator
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupValidator
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthValidator
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthValidator        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowValidator          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupValidator = fmt.Errorf("proto: unexpected end of group")
)
