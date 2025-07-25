// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gaia/liquid/v1beta1/liquid.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
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

// TokenizeShareLockStatus indicates whether the address is able to tokenize
// shares
type TokenizeShareLockStatus int32

const (
	// UNSPECIFIED defines an empty tokenize share lock status
	TOKENIZE_SHARE_LOCK_STATUS_UNSPECIFIED TokenizeShareLockStatus = 0
	// LOCKED indicates the account is locked and cannot tokenize shares
	TOKENIZE_SHARE_LOCK_STATUS_LOCKED TokenizeShareLockStatus = 1
	// UNLOCKED indicates the account is unlocked and can tokenize shares
	TOKENIZE_SHARE_LOCK_STATUS_UNLOCKED TokenizeShareLockStatus = 2
	// LOCK_EXPIRING indicates the account is unable to tokenize shares, but
	// will be able to tokenize shortly (after 1 unbonding period)
	TOKENIZE_SHARE_LOCK_STATUS_LOCK_EXPIRING TokenizeShareLockStatus = 3
)

var TokenizeShareLockStatus_name = map[int32]string{
	0: "TOKENIZE_SHARE_LOCK_STATUS_UNSPECIFIED",
	1: "TOKENIZE_SHARE_LOCK_STATUS_LOCKED",
	2: "TOKENIZE_SHARE_LOCK_STATUS_UNLOCKED",
	3: "TOKENIZE_SHARE_LOCK_STATUS_LOCK_EXPIRING",
}

var TokenizeShareLockStatus_value = map[string]int32{
	"TOKENIZE_SHARE_LOCK_STATUS_UNSPECIFIED":   0,
	"TOKENIZE_SHARE_LOCK_STATUS_LOCKED":        1,
	"TOKENIZE_SHARE_LOCK_STATUS_UNLOCKED":      2,
	"TOKENIZE_SHARE_LOCK_STATUS_LOCK_EXPIRING": 3,
}

func (x TokenizeShareLockStatus) String() string {
	return proto.EnumName(TokenizeShareLockStatus_name, int32(x))
}

func (TokenizeShareLockStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7b1e248decf35ce8, []int{0}
}

// Params defines the parameters for the x/liquid module.
type Params struct {
	// global_liquid_staking_cap represents a cap on the portion of stake that
	// comes from liquid staking providers
	GlobalLiquidStakingCap cosmossdk_io_math.LegacyDec `protobuf:"bytes,8,opt,name=global_liquid_staking_cap,json=globalLiquidStakingCap,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"global_liquid_staking_cap" yaml:"global_liquid_staking_cap"`
	// validator_liquid_staking_cap represents a cap on the portion of stake that
	// comes from liquid staking providers for a specific validator
	ValidatorLiquidStakingCap cosmossdk_io_math.LegacyDec `protobuf:"bytes,9,opt,name=validator_liquid_staking_cap,json=validatorLiquidStakingCap,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"validator_liquid_staking_cap" yaml:"validator_liquid_staking_cap"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b1e248decf35ce8, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

// TokenizeShareRecord represents a tokenized delegation
type TokenizeShareRecord struct {
	Id            uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Owner         string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	ModuleAccount string `protobuf:"bytes,3,opt,name=module_account,json=moduleAccount,proto3" json:"module_account,omitempty"`
	Validator     string `protobuf:"bytes,4,opt,name=validator,proto3" json:"validator,omitempty"`
}

func (m *TokenizeShareRecord) Reset()         { *m = TokenizeShareRecord{} }
func (m *TokenizeShareRecord) String() string { return proto.CompactTextString(m) }
func (*TokenizeShareRecord) ProtoMessage()    {}
func (*TokenizeShareRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b1e248decf35ce8, []int{1}
}
func (m *TokenizeShareRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TokenizeShareRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TokenizeShareRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TokenizeShareRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenizeShareRecord.Merge(m, src)
}
func (m *TokenizeShareRecord) XXX_Size() int {
	return m.Size()
}
func (m *TokenizeShareRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenizeShareRecord.DiscardUnknown(m)
}

var xxx_messageInfo_TokenizeShareRecord proto.InternalMessageInfo

func (m *TokenizeShareRecord) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *TokenizeShareRecord) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *TokenizeShareRecord) GetModuleAccount() string {
	if m != nil {
		return m.ModuleAccount
	}
	return ""
}

func (m *TokenizeShareRecord) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

// PendingTokenizeShareAuthorizations stores a list of addresses that have their
// tokenize share enablement in progress
type PendingTokenizeShareAuthorizations struct {
	Addresses []string `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (m *PendingTokenizeShareAuthorizations) Reset()         { *m = PendingTokenizeShareAuthorizations{} }
func (m *PendingTokenizeShareAuthorizations) String() string { return proto.CompactTextString(m) }
func (*PendingTokenizeShareAuthorizations) ProtoMessage()    {}
func (*PendingTokenizeShareAuthorizations) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b1e248decf35ce8, []int{2}
}
func (m *PendingTokenizeShareAuthorizations) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PendingTokenizeShareAuthorizations) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PendingTokenizeShareAuthorizations.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PendingTokenizeShareAuthorizations) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingTokenizeShareAuthorizations.Merge(m, src)
}
func (m *PendingTokenizeShareAuthorizations) XXX_Size() int {
	return m.Size()
}
func (m *PendingTokenizeShareAuthorizations) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingTokenizeShareAuthorizations.DiscardUnknown(m)
}

var xxx_messageInfo_PendingTokenizeShareAuthorizations proto.InternalMessageInfo

func (m *PendingTokenizeShareAuthorizations) GetAddresses() []string {
	if m != nil {
		return m.Addresses
	}
	return nil
}

// TokenizeShareRecordReward represents the properties of tokenize share
type TokenizeShareRecordReward struct {
	RecordId uint64                                      `protobuf:"varint,1,opt,name=record_id,json=recordId,proto3" json:"record_id,omitempty"`
	Reward   github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,2,rep,name=reward,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"reward"`
}

func (m *TokenizeShareRecordReward) Reset()         { *m = TokenizeShareRecordReward{} }
func (m *TokenizeShareRecordReward) String() string { return proto.CompactTextString(m) }
func (*TokenizeShareRecordReward) ProtoMessage()    {}
func (*TokenizeShareRecordReward) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b1e248decf35ce8, []int{3}
}
func (m *TokenizeShareRecordReward) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TokenizeShareRecordReward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TokenizeShareRecordReward.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TokenizeShareRecordReward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenizeShareRecordReward.Merge(m, src)
}
func (m *TokenizeShareRecordReward) XXX_Size() int {
	return m.Size()
}
func (m *TokenizeShareRecordReward) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenizeShareRecordReward.DiscardUnknown(m)
}

var xxx_messageInfo_TokenizeShareRecordReward proto.InternalMessageInfo

// LiquidValidator is the storage layout for details about a validator's liquid
// stake.
type LiquidValidator struct {
	// operator_address defines the address of the validator's operator; bech
	// encoded in JSON.
	OperatorAddress string `protobuf:"bytes,1,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty"`
	// Number of shares either tokenized or owned by a liquid staking provider
	LiquidShares cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=liquid_shares,json=liquidShares,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"liquid_shares" yaml:"liquid_shares"`
}

func (m *LiquidValidator) Reset()         { *m = LiquidValidator{} }
func (m *LiquidValidator) String() string { return proto.CompactTextString(m) }
func (*LiquidValidator) ProtoMessage()    {}
func (*LiquidValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b1e248decf35ce8, []int{4}
}
func (m *LiquidValidator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LiquidValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LiquidValidator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LiquidValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LiquidValidator.Merge(m, src)
}
func (m *LiquidValidator) XXX_Size() int {
	return m.Size()
}
func (m *LiquidValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_LiquidValidator.DiscardUnknown(m)
}

var xxx_messageInfo_LiquidValidator proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("gaia.liquid.v1beta1.TokenizeShareLockStatus", TokenizeShareLockStatus_name, TokenizeShareLockStatus_value)
	proto.RegisterType((*Params)(nil), "gaia.liquid.v1beta1.Params")
	proto.RegisterType((*TokenizeShareRecord)(nil), "gaia.liquid.v1beta1.TokenizeShareRecord")
	proto.RegisterType((*PendingTokenizeShareAuthorizations)(nil), "gaia.liquid.v1beta1.PendingTokenizeShareAuthorizations")
	proto.RegisterType((*TokenizeShareRecordReward)(nil), "gaia.liquid.v1beta1.TokenizeShareRecordReward")
	proto.RegisterType((*LiquidValidator)(nil), "gaia.liquid.v1beta1.LiquidValidator")
}

func init() { proto.RegisterFile("gaia/liquid/v1beta1/liquid.proto", fileDescriptor_7b1e248decf35ce8) }

var fileDescriptor_7b1e248decf35ce8 = []byte{
	// 778 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x4f, 0x4f, 0x33, 0x45,
	0x1c, 0xee, 0xb4, 0x95, 0xd0, 0xd1, 0xf7, 0xa5, 0x2e, 0x8d, 0x6e, 0x01, 0xb7, 0x75, 0x09, 0xda,
	0x54, 0xdb, 0x0d, 0x72, 0xab, 0xa7, 0xfe, 0x53, 0x1b, 0x1a, 0x68, 0x76, 0x8b, 0x31, 0x24, 0x66,
	0x9d, 0xee, 0x4e, 0xb6, 0x93, 0x6e, 0x77, 0xca, 0xcc, 0x16, 0x2c, 0x1f, 0xc0, 0x10, 0x4f, 0x5e,
	0x4c, 0x3c, 0x19, 0x12, 0x2f, 0xc6, 0xc4, 0x84, 0x03, 0x9f, 0xc0, 0x13, 0xde, 0x08, 0x27, 0xe3,
	0x01, 0x09, 0x1c, 0xf0, 0xec, 0x27, 0x30, 0xbb, 0xb3, 0x2d, 0x60, 0x90, 0x97, 0x4b, 0x3b, 0xbf,
	0xe7, 0xf7, 0xcc, 0xd3, 0xe7, 0xf7, 0x67, 0x0a, 0xf3, 0x0e, 0x22, 0x48, 0x73, 0xc9, 0xde, 0x98,
	0xd8, 0xda, 0xfe, 0x7a, 0x0f, 0xfb, 0x68, 0x3d, 0x0a, 0xcb, 0x23, 0x46, 0x7d, 0x2a, 0x2d, 0x06,
	0x8c, 0x72, 0x04, 0x45, 0x8c, 0xa5, 0x8c, 0x43, 0x1d, 0x1a, 0xe6, 0xb5, 0xe0, 0x24, 0xa8, 0x4b,
	0x6f, 0xa2, 0x21, 0xf1, 0xa8, 0x16, 0x7e, 0x46, 0x90, 0x62, 0x51, 0x3e, 0xa4, 0x5c, 0xeb, 0x21,
	0x8e, 0x67, 0xfa, 0x16, 0x25, 0x5e, 0x94, 0xcf, 0x8a, 0xbc, 0x29, 0xb4, 0x44, 0x20, 0x52, 0xea,
	0x55, 0x1c, 0xce, 0x75, 0x10, 0x43, 0x43, 0x2e, 0x7d, 0x0f, 0x60, 0xd6, 0x71, 0x69, 0x0f, 0xb9,
	0xa6, 0x30, 0x62, 0x72, 0x1f, 0x0d, 0x88, 0xe7, 0x98, 0x16, 0x1a, 0xc9, 0xf3, 0x79, 0x50, 0x48,
	0xd5, 0x76, 0xcf, 0x2e, 0x73, 0xb1, 0x3f, 0x2f, 0x73, 0xcb, 0x42, 0x84, 0xdb, 0x83, 0x32, 0xa1,
	0xda, 0x10, 0xf9, 0xfd, 0x72, 0x1b, 0x3b, 0xc8, 0x9a, 0x34, 0xb0, 0xf5, 0xcf, 0x65, 0x2e, 0x3f,
	0x41, 0x43, 0xb7, 0xa2, 0xfe, 0xaf, 0x9a, 0x7a, 0x71, 0x5a, 0x82, 0x91, 0x8f, 0x06, 0xb6, 0x7e,
	0xbe, 0x3d, 0x29, 0x02, 0xfd, 0x2d, 0x41, 0x6f, 0x87, 0x6c, 0x43, 0x90, 0xeb, 0x68, 0x24, 0xfd,
	0x08, 0xe0, 0xca, 0x3e, 0x72, 0x89, 0x8d, 0x7c, 0xca, 0x1e, 0xb3, 0x96, 0x0a, 0xad, 0x7d, 0xf9,
	0x3c, 0x6b, 0xab, 0xc2, 0xda, 0x53, 0x82, 0x8f, 0xba, 0xcb, 0xce, 0x6e, 0xfc, 0xd7, 0x60, 0xe5,
	0x9d, 0xbf, 0x8f, 0x73, 0xe0, 0xdb, 0xdb, 0x93, 0x62, 0x26, 0x9c, 0xf3, 0xd7, 0xd3, 0x49, 0x8b,
	0xbe, 0xaa, 0xdf, 0x00, 0xb8, 0xd8, 0xa5, 0x03, 0xec, 0x91, 0x43, 0x6c, 0xf4, 0x11, 0xc3, 0x3a,
	0xb6, 0x28, 0xb3, 0xa5, 0x97, 0x30, 0x4e, 0x6c, 0x19, 0xe4, 0x41, 0x21, 0xa9, 0xc7, 0x89, 0x2d,
	0x65, 0xe0, 0x6b, 0xf4, 0xc0, 0xc3, 0x4c, 0x8e, 0x07, 0xf5, 0xe8, 0x22, 0x90, 0xd6, 0xe0, 0xcb,
	0x21, 0xb5, 0xc7, 0x2e, 0x36, 0x91, 0x65, 0xd1, 0xb1, 0xe7, 0xcb, 0x89, 0x30, 0xfd, 0x42, 0xa0,
	0x55, 0x01, 0x4a, 0x2b, 0x30, 0x35, 0x33, 0x28, 0x27, 0x43, 0xc6, 0x1d, 0x50, 0x49, 0x06, 0x0e,
	0xd5, 0x1a, 0x54, 0x3b, 0xd8, 0xb3, 0x89, 0xe7, 0x3c, 0xb0, 0x53, 0x1d, 0xfb, 0x7d, 0xca, 0xc8,
	0x21, 0xf2, 0x09, 0xf5, 0x78, 0xa0, 0x84, 0x6c, 0x9b, 0x61, 0xce, 0x31, 0x97, 0x41, 0x3e, 0x11,
	0x28, 0xcd, 0x00, 0xf5, 0x57, 0x00, 0xb3, 0x8f, 0x14, 0xa3, 0xe3, 0x03, 0xc4, 0x6c, 0x69, 0x19,
	0xa6, 0x58, 0x18, 0x9b, 0xb3, 0xca, 0xe6, 0x05, 0xd0, 0xb2, 0x25, 0x02, 0xe7, 0x58, 0x48, 0x93,
	0xe3, 0xf9, 0x44, 0xe1, 0xf5, 0x8f, 0x56, 0xca, 0x51, 0x8f, 0x83, 0xb5, 0x9d, 0x2e, 0x7d, 0xd0,
	0xf0, 0x3a, 0x25, 0x5e, 0x6d, 0x23, 0x18, 0xe7, 0x2f, 0x7f, 0xe5, 0x3e, 0x70, 0x88, 0xdf, 0x1f,
	0xf7, 0xca, 0x16, 0x1d, 0x46, 0x9b, 0x1b, 0x7d, 0x95, 0xb8, 0x3d, 0xd0, 0xfc, 0xc9, 0x08, 0xf3,
	0xe9, 0x1d, 0xae, 0x47, 0x3f, 0x50, 0x99, 0x3f, 0x3a, 0xce, 0xc5, 0x7e, 0x08, 0x6a, 0xfe, 0x0d,
	0xc0, 0x05, 0x31, 0xb0, 0xcf, 0xa7, 0xdd, 0x90, 0xea, 0x30, 0x4d, 0x47, 0x98, 0x85, 0xd3, 0x8f,
	0x2a, 0x0b, 0xcd, 0xa6, 0x6a, 0xf2, 0xc5, 0x69, 0x29, 0x13, 0xb9, 0xaa, 0x8a, 0x8c, 0xe1, 0x33,
	0xe2, 0x39, 0xfa, 0xc2, 0xf4, 0x46, 0x04, 0x4b, 0x5f, 0xc1, 0x17, 0xd3, 0xcd, 0x09, 0xda, 0xc0,
	0xc5, 0x58, 0x6a, 0x1f, 0x3f, 0x6f, 0x0b, 0x33, 0x62, 0x0b, 0x1f, 0x28, 0xa8, 0xfa, 0x1b, 0x22,
	0x0e, 0xfb, 0xca, 0xef, 0x8a, 0x28, 0xfe, 0x0e, 0xe0, 0xdb, 0x0f, 0x9a, 0xde, 0xa6, 0xd6, 0xc0,
	0xf0, 0x91, 0x3f, 0xe6, 0x52, 0x11, 0xbe, 0xd7, 0xdd, 0xde, 0x6c, 0x6e, 0xb5, 0x76, 0x9b, 0xa6,
	0xf1, 0x59, 0x55, 0x6f, 0x9a, 0xed, 0xed, 0xfa, 0xa6, 0x69, 0x74, 0xab, 0xdd, 0x1d, 0xc3, 0xdc,
	0xd9, 0x32, 0x3a, 0xcd, 0x7a, 0xeb, 0x93, 0x56, 0xb3, 0x91, 0x8e, 0x49, 0x6b, 0xf0, 0xdd, 0x27,
	0xb8, 0xc1, 0xb9, 0xd9, 0x48, 0x03, 0xe9, 0x7d, 0xb8, 0xfa, 0xa4, 0x64, 0x44, 0x8c, 0x4b, 0x1f,
	0xc2, 0xc2, 0x2b, 0xf4, 0xcc, 0xe6, 0x17, 0x9d, 0x96, 0xde, 0xda, 0xfa, 0x34, 0x9d, 0x58, 0x4a,
	0x1e, 0xfd, 0xa4, 0xc4, 0x6a, 0xee, 0xd9, 0xb5, 0x02, 0xce, 0xaf, 0x15, 0x70, 0x75, 0xad, 0x80,
	0xef, 0x6e, 0x94, 0xd8, 0xf9, 0x8d, 0x12, 0xfb, 0xe3, 0x46, 0x89, 0xed, 0xea, 0xf7, 0x26, 0xbd,
	0x37, 0x26, 0xd6, 0x80, 0x13, 0x77, 0x1f, 0xb3, 0xd2, 0x21, 0xf5, 0xf0, 0x7d, 0x40, 0xf3, 0xfb,
	0x84, 0xd9, 0xa5, 0x11, 0x62, 0xfe, 0xa4, 0x64, 0xf5, 0x11, 0xf1, 0xb8, 0x16, 0x3c, 0xbc, 0x52,
	0xb8, 0x0c, 0xd3, 0xc7, 0x17, 0x06, 0xbd, 0xb9, 0xf0, 0x5f, 0x6e, 0xe3, 0xdf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x22, 0x2f, 0x0b, 0x2e, 0x82, 0x05, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.GlobalLiquidStakingCap.Equal(that1.GlobalLiquidStakingCap) {
		return false
	}
	if !this.ValidatorLiquidStakingCap.Equal(that1.ValidatorLiquidStakingCap) {
		return false
	}
	return true
}
func (this *TokenizeShareRecord) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TokenizeShareRecord)
	if !ok {
		that2, ok := that.(TokenizeShareRecord)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Id != that1.Id {
		return false
	}
	if this.Owner != that1.Owner {
		return false
	}
	if this.ModuleAccount != that1.ModuleAccount {
		return false
	}
	if this.Validator != that1.Validator {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ValidatorLiquidStakingCap.Size()
		i -= size
		if _, err := m.ValidatorLiquidStakingCap.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLiquid(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.GlobalLiquidStakingCap.Size()
		i -= size
		if _, err := m.GlobalLiquidStakingCap.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLiquid(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	return len(dAtA) - i, nil
}

func (m *TokenizeShareRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TokenizeShareRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TokenizeShareRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintLiquid(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ModuleAccount) > 0 {
		i -= len(m.ModuleAccount)
		copy(dAtA[i:], m.ModuleAccount)
		i = encodeVarintLiquid(dAtA, i, uint64(len(m.ModuleAccount)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintLiquid(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintLiquid(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PendingTokenizeShareAuthorizations) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PendingTokenizeShareAuthorizations) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PendingTokenizeShareAuthorizations) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Addresses) > 0 {
		for iNdEx := len(m.Addresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Addresses[iNdEx])
			copy(dAtA[i:], m.Addresses[iNdEx])
			i = encodeVarintLiquid(dAtA, i, uint64(len(m.Addresses[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *TokenizeShareRecordReward) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TokenizeShareRecordReward) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TokenizeShareRecordReward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Reward) > 0 {
		for iNdEx := len(m.Reward) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Reward[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLiquid(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.RecordId != 0 {
		i = encodeVarintLiquid(dAtA, i, uint64(m.RecordId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LiquidValidator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LiquidValidator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LiquidValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintLiquid(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.OperatorAddress) > 0 {
		i -= len(m.OperatorAddress)
		copy(dAtA[i:], m.OperatorAddress)
		i = encodeVarintLiquid(dAtA, i, uint64(len(m.OperatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintLiquid(dAtA []byte, offset int, v uint64) int {
	offset -= sovLiquid(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.GlobalLiquidStakingCap.Size()
	n += 1 + l + sovLiquid(uint64(l))
	l = m.ValidatorLiquidStakingCap.Size()
	n += 1 + l + sovLiquid(uint64(l))
	return n
}

func (m *TokenizeShareRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovLiquid(uint64(m.Id))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovLiquid(uint64(l))
	}
	l = len(m.ModuleAccount)
	if l > 0 {
		n += 1 + l + sovLiquid(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovLiquid(uint64(l))
	}
	return n
}

func (m *PendingTokenizeShareAuthorizations) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Addresses) > 0 {
		for _, s := range m.Addresses {
			l = len(s)
			n += 1 + l + sovLiquid(uint64(l))
		}
	}
	return n
}

func (m *TokenizeShareRecordReward) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RecordId != 0 {
		n += 1 + sovLiquid(uint64(m.RecordId))
	}
	if len(m.Reward) > 0 {
		for _, e := range m.Reward {
			l = e.Size()
			n += 1 + l + sovLiquid(uint64(l))
		}
	}
	return n
}

func (m *LiquidValidator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OperatorAddress)
	if l > 0 {
		n += 1 + l + sovLiquid(uint64(l))
	}
	l = m.LiquidShares.Size()
	n += 1 + l + sovLiquid(uint64(l))
	return n
}

func sovLiquid(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLiquid(x uint64) (n int) {
	return sovLiquid(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquid
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GlobalLiquidStakingCap", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GlobalLiquidStakingCap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorLiquidStakingCap", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ValidatorLiquidStakingCap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLiquid(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquid
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
func (m *TokenizeShareRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquid
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
			return fmt.Errorf("proto: TokenizeShareRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TokenizeShareRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ModuleAccount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ModuleAccount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLiquid(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquid
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
func (m *PendingTokenizeShareAuthorizations) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquid
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
			return fmt.Errorf("proto: PendingTokenizeShareAuthorizations: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PendingTokenizeShareAuthorizations: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addresses", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addresses = append(m.Addresses, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLiquid(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquid
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
func (m *TokenizeShareRecordReward) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquid
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
			return fmt.Errorf("proto: TokenizeShareRecordReward: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TokenizeShareRecordReward: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RecordId", wireType)
			}
			m.RecordId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RecordId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reward", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Reward = append(m.Reward, types.DecCoin{})
			if err := m.Reward[len(m.Reward)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLiquid(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquid
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
func (m *LiquidValidator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquid
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
			return fmt.Errorf("proto: LiquidValidator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LiquidValidator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OperatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OperatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidShares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquid
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
				return ErrInvalidLengthLiquid
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLiquid
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
			skippy, err := skipLiquid(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquid
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
func skipLiquid(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLiquid
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
					return 0, ErrIntOverflowLiquid
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
					return 0, ErrIntOverflowLiquid
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
				return 0, ErrInvalidLengthLiquid
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLiquid
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLiquid
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLiquid        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLiquid          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLiquid = fmt.Errorf("proto: unexpected end of group")
)
