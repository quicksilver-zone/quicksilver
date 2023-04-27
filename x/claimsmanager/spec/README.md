# Claims Manager

## Abstract

Module, `x/claimsmanager`, provides storage and retrieval mechanisms for proof
based claims utilized in other modules.

## Contents

1. [Concepts](#concepts)
1. [State](#state)
1. [Messages](#messages)
1. [Transactions](#transactions)
1. [Events](#events)
1. [Hooks](#hooks)
1. [Queries](#queries)
1. [Keepers](#keepers)
1. [Parameters](#parameters)
1. [Begin Block](#begin-block)
1. [End Block](#end-block)

## Concepts

`x/claimsmanager` is simply a data store for use by other modules to avoid
unnecessary or circular dependencies.

## State

### ClaimType

```go
const (
	// Undefined action (per protobuf spec)
	ClaimTypeUndefined    ClaimType = 0
	ClaimTypeLiquidToken  ClaimType = 1
	ClaimTypeOsmosisPool  ClaimType = 2
	ClaimTypeCrescentPool ClaimType = 3
	ClaimTypeOsmosisLiquidGamm ClaimType = 4
	ClaimTypeUmeeCollateral ClaimType = 5
)

var ClaimType_name = map[int32]string{
	0: "ClaimTypeUndefined",
	1: "ClaimTypeLiquidToken",
	2: "ClaimTypeOsmosisPool",
	3: "ClaimTypeCrescentPool",
	4: "ClaimTypeSifchainPool",
}

var ClaimType_value = map[string]int32{
	"ClaimTypeUndefined":    0,
	"ClaimTypeLiquidToken":  1,
	"ClaimTypeOsmosisPool":  2,
	"ClaimTypeCrescentPool": 3,
	"ClaimTypeSifchainPool": 4,
}
```

### Claim

```go
var (
	KeyPrefixClaim          = []byte{0x00}
	KeyPrefixLastEpochClaim = []byte{0x01}
)

// ClaimKey returns the key for storing a given claim.
func GetGenericKeyClaim(key []byte, chainID string, address string, module ClaimType, srcChainID string) []byte {
	typeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(typeBytes, uint32(module))
	key = append(key, []byte(chainID)...)
	key = append(key, byte(0x00))
	key = append(key, []byte(address)...)
	key = append(key, typeBytes...)
	return append(key, []byte(srcChainID)...)
}

func GetKeyClaim(chainID string, address string, module ClaimType, srcChainID string) []byte {
	return GetGenericKeyClaim(KeyPrefixClaim, chainID, address, module, srcChainID)
}

func GetPrefixClaim(chainID string) []byte {
	return append(KeyPrefixClaim, []byte(chainID)...)
}

func GetPrefixUserClaim(chainID string, address string) []byte {
	key := KeyPrefixClaim
	key = append(key, []byte(chainID)...)
	key = append(key, byte(0x00))
	key = append(key, []byte(address)...)
	return key
}

func GetKeyLastEpochClaim(chainID string, address string, module ClaimType, srcChainID string) []byte {
	return GetGenericKeyClaim(KeyPrefixLastEpochClaim, chainID, address, module, srcChainID)
}

func GetPrefixLastEpochClaim(chainID string) []byte {
	return append(KeyPrefixLastEpochClaim, []byte(chainID)...)
}

func GetPrefixLastEpochUserClaim(chainID string, address string) []byte {
	key := KeyPrefixLastEpochClaim
	key = append(key, []byte(chainID)...)
	key = append(key, byte(0x00))
	key = append(key, []byte(address)...)
	return key
}

// Claim define the users claim for holding assets within a given zone.
type Claim struct {
	UserAddress   string    `protobuf:"bytes,1,opt,name=user_address,json=userAddress,proto3" json:"user_address,omitempty"`
	ChainId       string    `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Module        ClaimType `protobuf:"varint,3,opt,name=module,proto3,enum=quicksilver.claimsmanager.v1.ClaimType" json:"module,omitempty"`
	SourceChainId string    `protobuf:"bytes,4,opt,name=source_chain_id,json=sourceChainId,proto3" json:"source_chain_id,omitempty"`
	Amount        uint64    `protobuf:"varint,5,opt,name=amount,proto3" json:"amount,omitempty"`
}
```

### Proof

```go
// Proof defines a type used to cryptographically prove a claim.
type Proof struct {
	Key       []byte           `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Data      []byte           `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	ProofOps  *crypto.ProofOps `protobuf:"bytes,3,opt,name=proof_ops,proto3" json:"proof_ops,omitempty"`
	Height    int64            `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
	ProofType string           `protobuf:"bytes,5,opt,name=proof_type,proto3" json:"proof_type,omitempty"`
}
```

### ProtocolData

```
// KeyedProtocolData is a keyed representation of protocoldata used in the
// genesis file.
message KeyedProtocolData {
  string key = 1;
  ProtocolData protocol_data = 3;
}

// ProtocolData is an arbitrary data type held against a given zone for the
// determination of claim allocaton.
message ProtocolData {
  string type = 1;
  bytes data = 2 [ (gogoproto.casttype) = "encoding/json.RawMessage" ];
  bool is_latest = 3;
}

// ProtocolDataType is an enumerated list of types of protocol data.
enum ProtocolDataType {
  option (gogoproto.goproto_enum_prefix) = false;

  // Undefined action (per protobuf spec)
  ProtocolDataTypeUndefined = 0;
  ProtocolDataTypeConnection = 1;
  ProtocolDataTypeOsmosisParams = 2;
  ProtocolDataTypeLiquidToken = 3;
  ProtocolDataTypeOsmosisPool = 4;
  ProtocolDataTypeCrescentPool = 5;
}
```

### Events

```
Message ClaimableEvent {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  
  string event_module = 0;              // module name doing the claim, “govbyproxy”, “participationrewards”, “airdrop” (types.ModuleName)
  string event_name = 1;                // “epoch”, “proposal/cosmoshub-4/764”, “airdrop/cosmoshub-4”
  Map<string, int64> heights = 2;       // heights across each chain for this event - when initialised create with the chainID and 0, for each relevant zone. For all zones, use all zones we have a connectionData for; airdrop proofs will require _latest_ (pool contents, etc.) to be submitted by the user.
  max_claim_time time.Time = 3;         // when can we claim until? (Proposal end vote time, next epoch boundary, airdrop end)
}
```

## Messages

Description of message types that trigger state transitions;

### MsgSubmitClaim

```
// MsgSubmitClaim represents a message type for submitting a participation
// claim regarding the given zone (chain).
message MsgSubmitClaim {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;

  string user_address = 1 [ json_name = "user_address", (cosmos_proto.scalar) = "cosmos.AddressString" ];
  string zone = 2 [ json_name = "zone" ];
  string src_zone = 3 [ json_name = "src_zone" ];
  quicksilver.claimsmanager.v1.ClaimType claim_type = 4 [ json_name = "claim_type" ];
  repeated quicksilver.claimsmanager.v1.Proof proofs = 5 [ json_name = "proofs" ];
  string event = 6 [ json_name = "event" ];
}
```

### MsgUpdateProtocolData

// MsgUpdateProtocolData represents a message type for updating a 'latest' style protocol data.
message MsgUpdateProtocolData {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;

  string type = 1 [ (gogoproto.moretags) = "yaml:\"type\"" ];
  string data = 2 [
    (gogoproto.moretags) = "yaml:\"data\"",
    (gogoproto.casttype) = "encoding/json.RawMessage"
  ];
  string key = 3 [ (gogoproto.moretags) = "yaml:\"key\"" ];
  quicksilver.claimsmanager.v1.Proof proof = 4 [ json_name = "proof" ];
  string event = 5 [ (gogoproto.moretags) = "yaml:\"event\"" ];
}

## Transactions

Description of transactions that collect messages in specific contexts to trigger state transitions;

`x/claimsmanager` does not provide any transactions, it is the responsibility of calling modules to do so.

## Events

N/A

## Hooks

N/A

## Queries

Description of available information request queries;

```protobuf
// Query provides defines the gRPC querier service.
service Query {
  // Claims returns all zone claims from the current epoch.
  rpc Claims(QueryClaimsRequest) returns (QueryClaimsResponse) {
    option (google.api.http).get = "/quicksilver/claimsmanager/v1/claims/{chain_id}";
  }

  // LastEpochClaims returns all zone claims from the last epoch.
  rpc LastEpochClaims(QueryClaimsRequest) returns (QueryClaimsResponse) {
    option (google.api.http).get = "/quicksilver/claimsmanager/v1/previous_epoch_claims/{chain_id}";
  }

  // UserClaims returns all zone claims for a given address from the current epoch.
  rpc UserClaims(QueryClaimsRequest) returns (QueryClaimsResponse) {
    option (google.api.http).get = "/quicksilver/claimsmanager/v1/user/{address}/claims";
  }

  // UserLastEpochClaims returns all zone claims for a given address from the last epoch.
  rpc UserLastEpochClaims(QueryClaimsRequest) returns (QueryClaimsResponse) {
    option (google.api.http).get = "/quicksilver/claimsmanager/v1/user/{address}/previous_epoch_claims";
  }
}
```

The above queries use the following `Request` and `Response` types:

```go
// QueryClaimsRequest is the request type for the Query/Claims RPC method.
type QueryClaimsRequest struct {
	ChainId    string             `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
	Address    string             `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Pagination *query.PageRequest `protobuf:"bytes,3,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryClaimsResponse is the response type for the Query/Claims RPC method.
type QueryClaimsResponse struct {
	Claims     []Claim             `protobuf:"bytes,1,rep,name=claims,proto3" json:"claims"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
```

## Keepers

<https://pkg.go.dev/github.com/ingenuity-build/quicksilver/x/claimsmanager/keeper>

## Parameters

N/A

## Begin Block

N/A

## End Block

N/A

