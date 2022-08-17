# Airdrop

## Abstract

The purpose of this module is to distribute QCK airdrops to users for engaging in activities related to newly onboarded zones.

**Objectives:**

* Provide airdrop to qualifying users;
* Couple airdrop to specific actions or tasks;
* Airdrops configured on a per zone basis;
* Must decay to zero over a specified period, starting at a specific time / block height;

## Contents

1. [Concepts](#Concepts)
1. [State](#State)
1. [Messages](#Messages)
1. [Transactions](#Transactions)
1. [Events](#Events)
1. [Hooks](#Hooks)
1. [Queries](#Queries)
1. [Keepers](#Keepers)
1. [Parameters](#Parameters)
1. [Proposals](#Proposals)
1. [Begin Block](#Begin-Block)
1. [End Block](#End-Block)

## Concepts

Key concepts, mechanisms and core logic for the module;

### Module Accounts

The airdrop module utilizes a module account for every zone airdrop identified by the zone's chain ID to manage airdrop claims accounting. It also has a main module account where all unclaimed amounts will be collected on conclusion of a zone airdrop.

### ZoneDrops

A `ZoneDrop` refers to a zone airdrop and is identified by the zone's chain ID.

#### Status

A `ZoneDrop` is considered `Active` if, and only if, the current `BlockTime` is between the time window of the airdrop `StartTime` and the added airdrop `Duration` and `Decay` times.

**Formula:**

* active:  `BlockTime` > `StartTime` && `BlockTime` < `StartTime`+`Duration`+`Decay`;
* future:  `BlockTime` < `StartTime`;
* expired: `BlockTime` > `StartTime`+`Duration`+`Decay`;

#### Duration & Decay

Airdrop `Duration` refers to the time period an airdrop is active at its full reward potential. Every action claimed during this period will receive rewards directly proportional to the qualifying allocation and the weight of the particular action.

Airdrop `Decay` refers to the time period after the `Duration` during which the airdrop is still active, but rewards are discounted according to the decay proportion. Thus, any action claimed at the half way mark of the decay duration will only receive 50% of that action's qualifying allocation and weight.

### Actions

Airdrop rewards are coupled to specific actions or tasks that users are to perform to unlock airdrop rewards, that they may then claim. Of note here is that the deposit action is subdivided into tiers, where each subsequent tier is unlocked by reaching a particular threshold of the `BaseValue` in deposits.

### Claim Records

A `ClaimRecord` represents an individual user's full potential airdrop rewards and is set at as part of the airdrop proposal. Individual rewards are scalable according to the `BaseValue` which may represent particular snapshot data in accordance with the airdrop proposal.

Any claims completed are recorded against the `ClaimRecord` and claimed amounts may never exceed the defined `MaxAllocation`.

## State

### Action

```
// Action is used as an enum to denote specific actions or tasks.
type Action int32

const (
	// Initial claim action
	ActionInitialClaim Action = 0
	// Deposit tier 1 (e.g. > 5% of base_value)
	ActionDepositT1 Action = 1
	// Deposit tier 2 (e.g. > 10% of base_value)
	ActionDepositT2 Action = 2
	// Deposit tier 3 (e.g. > 15% of base_value)
	ActionDepositT3 Action = 3
	// Deposit tier 4 (e.g. > 22% of base_value)
	ActionDepositT4 Action = 4
	// Deposit tier 5 (e.g. > 30% of base_value)
	ActionDepositT5 Action = 5
	// Active QCK delegation
	ActionStakeQCK Action = 6
	// Intent is set
	ActionSignalIntent Action = 7
	// Cast governance vote on QS
	ActionQSGov Action = 8
	// Governance By Proxy (GbP): cast vote on remote zone
	ActionGbP Action = 9
	// Provide liquidity on Osmosis
	ActionOsmosis Action = 10
)

var Action_name = map[int32]string{
	0:  "ActionInitialClaim",
	1:  "ActionDepositT1",
	2:  "ActionDepositT2",
	3:  "ActionDepositT3",
	4:  "ActionDepositT4",
	5:  "ActionDepositT5",
	6:  "ActionStakeQCK",
	7:  "ActionSignalIntent",
	8:  "ActionQSGov",
	9:  "ActionGbP",
	10: "ActionOsmosis",
}

var Action_value = map[string]int32{
	"ActionInitialClaim": 0,
	"ActionDepositT1":    1,
	"ActionDepositT2":    2,
	"ActionDepositT3":    3,
	"ActionDepositT4":    4,
	"ActionDepositT5":    5,
	"ActionStakeQCK":     6,
	"ActionSignalIntent": 7,
	"ActionQSGov":        8,
	"ActionGbP":          9,
	"ActionOsmosis":      10,
}
```

### Status

```
// Status is used as an enum to denote zone status.
type Status int32

const (
	StatusActive  Status = 0
	StatusFuture  Status = 1
	StatusExpired Status = 2
)

var Status_name = map[int32]string{
	0: "StatusActive",
	1: "StatusFuture",
	2: "StatusExpired",
}

var Status_value = map[string]int32{
	"StatusActive":  0,
	"StatusFuture":  1,
	"StatusExpired": 2,
}
```

### ZoneDrop

```
KeyPrefixZoneDrop    = []byte{0x01}

func GetKeyZoneDrop(chainID string) []byte {
	return append(KeyPrefixZoneDrop, []byte(chainID)...)
}

// ZoneDrop represents an airdrop for a specific zone.
type ZoneDrop struct {
	ChainId     string                                   `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	StartTime   time.Time                                `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3,stdtime" json:"start_time" yaml:"start_time"`
	Duration    time.Duration                            `protobuf:"bytes,3,opt,name=duration,proto3,stdduration" json:"duration,omitempty" yaml:"duration"`
	Decay       time.Duration                            `protobuf:"bytes,4,opt,name=decay,proto3,stdduration" json:"decay,omitempty" yaml:"decay"`
	Allocation  uint64                                   `protobuf:"varint,5,opt,name=allocation,proto3" json:"allocation,omitempty"`
	Actions     []github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,rep,name=actions,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"actions"`
	IsConcluded bool                                     `protobuf:"varint,7,opt,name=is_concluded,json=isConcluded,proto3" json:"is_concluded,omitempty"`
}
```

### ClaimRecord

```
KeyPrefixClaimRecord = []byte{0x02}

func GetKeyClaimRecord(chainID string, addr sdk.AccAddress) []byte {
	return append(append(KeyPrefixClaimRecord, []byte(chainID)...), addr...)
}

func GetPrefixClaimRecord(chainID string) []byte {
	return append(KeyPrefixClaimRecord, []byte(chainID)...)
}

// ClaimRecord represents a users' claim (including completed claims) for a
// given zone.
type ClaimRecord struct {
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// Protobuf3 does not allow enum as map key
	ActionsCompleted map[int32]*CompletedAction `protobuf:"bytes,3,rep,name=actions_completed,json=actionsCompleted,proto3" json:"actions_completed,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	MaxAllocation    uint64                     `protobuf:"varint,4,opt,name=max_allocation,json=maxAllocation,proto3" json:"max_allocation,omitempty"`
	BaseValue        uint64                     `protobuf:"varint,5,opt,name=base_value,json=baseValue,proto3" json:"base_value,omitempty"`
}
```

### CompletedAction

```
// CompletedAction represents a claim action completed by the user.
type CompletedAction struct {
	CompleteTime time.Time `protobuf:"bytes,1,opt,name=complete_time,json=completeTime,proto3,stdtime" json:"complete_time" yaml:"complete_time"`
	ClaimAmount  uint64    `protobuf:"varint,2,opt,name=claim_amount,json=claimAmount,proto3" json:"claim_amount,omitempty"`
}
```

## Messages

Description of message types that trigger state transitions;

```
// Msg defines the airdrop Msg service.
service Msg {
  rpc Claim(MsgClaim) returns (MsgClaimResponse) {
    option (google.api.http) = {
      post : "/quicksilver/tx/v1/airdrop/claim"
      body : "*"
    };
  }
}
```

### claim

Claim the airdrop for the given action in the given zone.

Use: `claim [chainID] [action]`

```
type MsgClaim struct {
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
	Action  int32  `protobuf:"varint,2,opt,name=action,proto3" json:"action,omitempty" yaml:"action"`
	Address string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	Proof   []byte `protobuf:"bytes,4,opt,name=proof,proto3" json:"proof,omitempty" yaml:"proof"`
}

type MsgClaimResponse struct {
	Amount uint64 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty" yaml:"amount"`
}
```

## Transactions

Description of transactions that collect messages in specific contexts to trigger state transitions;

## Events

Events emitted by module for tracking messages and index transactions;

### RegisterZoneDropProposal

| Type              | Attribute Key | Attribute Value |
|:------------------|:--------------|:----------------|
| message           | module        | airdrop         |
| register_zonedrop | chain_id      | {chain_id}      |

### MsgClaim

| Type              | Attribute Key | Attribute Value |
|:------------------|:--------------|:----------------|
| airdrop_claim     | sender        | {address}       |
| airdrop_claim     | zone          | {chain_id}      |
| airdrop_claim     | action        | {action}        |
| airdrop_claim     | amount        | {amount}        |

## Hooks

Description of hook functions that may be used by other modules to execute operations at specific points within this module;

## Queries

Description of available information request queries;

```
service Query {
  // Params returns the total set of airdrop parameters.
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/quicksilver/airdrop/v1/params";
  }
  // ZoneDrop returns the details of the specified zone airdrop.
  rpc ZoneDrop(QueryZoneDropRequest) returns (QueryZoneDropResponse) {
    option (google.api.http).get =
        "/quicksilver/airdrop/v1/zonedrop/{chain_id}";
  }
  // AccountBalance returns the module account balance of the specified zone.
  rpc AccountBalance(QueryAccountBalanceRequest)
      returns (QueryAccountBalanceResponse) {
    option (google.api.http).get =
        "/quicksilver/airdrop/v1/accountbalance/{chain_id}";
  }
  // ZoneDrops returns all zone airdrops of the specified status.
  rpc ZoneDrops(QueryZoneDropsRequest) returns (QueryZoneDropsResponse) {
    option (google.api.http).get = "/quicksilver/airdrop/v1/zonedrops/{status}";
  }
  // ClaimRecord returns the claim record that corresponds to the given zone and
  // address.
  rpc ClaimRecord(QueryClaimRecordRequest) returns (QueryClaimRecordResponse) {
    option (google.api.http).get =
        "/quicksilver/airdrop/v1/claimrecord/{chain_id}/{address}";
  }
  // ClaimRecords returns all the claim records of the given zone.
  rpc ClaimRecords(QueryClaimRecordsRequest)
      returns (QueryClaimRecordsResponse) {
    option (google.api.http).get =
        "/quicksilver/airdrop/v1/claimrecords/{chain_id}";
  }
}
```

### params

Query the current airdrop module parameters.

Use: `params`

```
// QueryParamsRequest is the request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

// QueryParamsResponse is the response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params defines the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}
```

### zone

Query the airdrop details of the specified zone.

Use: `zone [chain_id]`

```
// QueryZoneDropRequest is the request type for Query/ZoneDrop RPC method.
type QueryZoneDropRequest struct {
	// chain_id identifies the zone.
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
}

// QueryZoneDropResponse is the response type for Query/ZoneDrop RPC method.
type QueryZoneDropResponse struct {
	ZoneDrop ZoneDrop `protobuf:"bytes,1,opt,name=zone_drop,json=zoneDrop,proto3" json:"zone_drop"`
}
```

### account-balance

Returns the airdrop module account balance of the specified zone.

Use: `account-balance [chain_id]`

```
// QueryAccountBalanceRequest is the request type for Query/AccountBalance RPC
// method.
type QueryAccountBalanceRequest struct {
	// chain_id identifies the zone.
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
}

// QueryAccountBalanceResponse is the response type for Query/AccountBalance RPC
// method.
type QueryAccountBalanceResponse struct {
	AccountBalance *types.Coin `protobuf:"bytes,1,opt,name=account_balance,json=accountBalance,proto3" json:"account_balance,omitempty" yaml:"account_balance"`
}
```

### zone-drops

Query all airdrops of the specified status.

Use: `zone-drops [status]`

```
// QueryZoneDropsRequest is the request type for Query/ZoneDrops RPC method.
type QueryZoneDropsRequest struct {
	// status enables to query zone airdrops matching a given status:
	//  - Active
	//  - Future
	//  - Expired
	Status     Status             `protobuf:"varint,1,opt,name=status,proto3,enum=quicksilver.airdrop.v1.Status" json:"status,omitempty"`
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryZoneDropResponse is the response type for Query/ZoneDrops RPC method.
type QueryZoneDropsResponse struct {
	ZoneDrops  []ZoneDrop          `protobuf:"bytes,1,rep,name=zone_drops,json=zoneDrops,proto3" json:"zone_drops"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
```

### claim-record

Query airdrop claim record details of the given address for the given zone.

Use: `claim-record [chain_id] [address]`

```
// QueryClaimRecordRequest is the request type for Query/ClaimRecord RPC method.
type QueryClaimRecordRequest struct {
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
}

// QueryClaimRecordResponse is the response type for Query/ClaimRecord RPC
// method.
type QueryClaimRecordResponse struct {
	ClaimRecord *ClaimRecord `protobuf:"bytes,1,opt,name=claim_record,json=claimRecord,proto3" json:"claim_record,omitempty" yaml:"claim_record"`
}
```

## Keepers

Keepers exposed by module;

## Parameters

Module parameters:

| Key   | Type  | Example |
|:-- ---|:-- ---|:--   ---|

Description of parameters:

* `param_name` - short description;

## Proposals

Register a zone airdrop proposal.

```
type RegisterZoneDropProposal struct {
	Title        string    `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description  string    `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	ZoneDrop     *ZoneDrop `protobuf:"bytes,3,opt,name=zone_drop,json=zoneDrop,proto3" json:"zone_drop,omitempty" yaml:"zone_drop"`
	ClaimRecords []byte    `protobuf:"bytes,4,opt,name=claim_records,json=claimRecords,proto3" json:"claim_records,omitempty" yaml:"claim_records"`
}
```

## Begin Block

Description of logic executed with optional methods or external hooks;

## End Block

Description of logic executed with optional methods or external hooks;

At the end of every block the module iterates through all unconcluded airdrops (expired but not yet concluded) and calls `EndZoneDrop` for each instance, that deletes all associated `ClaimRecord`s.

