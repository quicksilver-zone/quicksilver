# Participation Rewards

## Abstract

Module, `x/participatiorewards`, defines and implements the mechanisms to track,
allocate and distribute protocol participation rewards to users.

## Contents

1. [Concepts](#concepts)
1. [State](#state)
1. [Messages](#messages)
1. [Transactions](#transactions)
1. [Proposals](#proposals)
1. [Events](#events)
1. [Hooks](#hooks)
1. [Queries](#queries)
1. [Keepers](#keepers)
1. [Parameters](#parameters)
1. [Begin Block](#begin-block)
1. [End Block](#end-block)
1. [After Epoch End](#after-epoch-end)

## Concepts

The purpose of the participation rewards module is to reward users for protocol participation.

Specifically, we want to reward users for:

1. Staking and locking of QCK on the Quicksilver chain;
2. Positive validator selection, validators are ranked equally on performance and decentralization;
3. Holdings of off-chain assets (qAssets);

The total inflation allocation for participation rewards is divided
proportionally for each of the above according to the module [parameters](#parameters).

### 1. Lockup Rewards

The staking and lockup rewards allocation is moved to the fee collector account
to be distributed by the staking module on the next begin blocker. Thus, the
**user rewards allocation** will be proportional to their stake of the overall
staked pool.

### 2. Validator Selection Rewards

Validators are ranked on two aspects with equal weighting, namely
decentralization and performance.

The **decentralilzation scores** are based on the normalized voting power of the
validators within a given zone, favouring smaller validators.

The **performance scores** are based on the validator rewards earned by a
special performance account that delegates an exact amount to each validator.
The total rewards earned by the performance account is divided by the number of
active validators to obtain the expected rewards. The performance score for
each validator is then simply the percentage of actual rewards compared to the
expected rewards (capped at 100%).

The overall **validator scores** are simply the multiple of their
decentralization score and their performance score.

Individual **users scores** are based on their validator selection intent
signalled at the previous epoch boundary. The user intent weights are
multiplied by the corresponding validator scores for the given zone and an
overall user score is calculated for the given zone along with an
**overall zone score**.

Each zone receives a **proportional rewards allocation** based on the total
value locked (TVL) for the zone relative to the TVL of all zones across the
protocol.

The overall zone score and the proportional rewards allocation determines the
amount of **tokens per point (TPP)** to be allocated for the given zone. Thus,
the **user rewards allocation** is the user's score multiplied by the TPP.

### 3. Holdings Rewards

Each zone receives a **proportional rewards allocation** based on the total
value locked (TVL) for the zone relative to the TVL of all zones across the
protocol.

Thus, the **user rewards allocation** is proportional to their holdings of
qAssets across all zones, capped at 2% per account.

### 4. Protocol Data

Rewarding for protocol participation, specifically for off-chain assets,
requires provable claims. **Protocol Data** describes any protocol specific
state that must be tracked in order to obtain provable claims. The module
defines a `Submodule` interface that must be implemented to obtain provable
claims against any specific protocol.

The following standrad sub-modules are implemented:

* `LiquidTokenModule` - to track off-chain liquid qAssets.
* `OsmosisModule` - to track qAssets locked in Osmosis pools.

## State

A `Score` is maintained for every `Validator` within a `Zone`. `Score` is
initially set to zero and is updated at the end of every epoch to reflect the
**overall score** for the validator (decntralization_score * performance_score).

A `ValidatorSelectionAllocation` and `HoldingsAllocation` are maintained for
every `Zone`. These are calculated and set at the end of every epoch according
to the rewards allocation proportions that are distributed to zones based on
their Total Value Locked (TVL) relative to the TVL of the overall protocol.

### ProtocolData

#### Types

```go
type ProtocolDataType int32

const (
	// Undefined action (per protobuf spec)
	ProtocolDataTypeUndefined     ProtocolDataType = 0
	ProtocolDataTypeConnection    ProtocolDataType = 1
	ProtocolDataTypeOsmosisParams ProtocolDataType = 2
	ProtocolDataTypeLiquidToken   ProtocolDataType = 3
	ProtocolDataTypeOsmosisPool   ProtocolDataType = 4
	ProtocolDataTypeCrescentPool  ProtocolDataType = 5
	ProtocolDataTypeSifchainPool  ProtocolDataType = 6
)

var ProtocolDataType_name = map[int32]string{
	0: "ProtocolDataTypeUndefined",
	1: "ProtocolDataTypeConnection",
	2: "ProtocolDataTypeOsmosisParams",
	3: "ProtocolDataTypeLiquidToken",
	4: "ProtocolDataTypeOsmosisPool",
	5: "ProtocolDataTypeCrescentPool",
	6: "ProtocolDataTypeSifchainPool",
}

var ProtocolDataType_value = map[string]int32{
	"ProtocolDataTypeUndefined":     0,
	"ProtocolDataTypeConnection":    1,
	"ProtocolDataTypeOsmosisParams": 2,
	"ProtocolDataTypeLiquidToken":   3,
	"ProtocolDataTypeOsmosisPool":   4,
	"ProtocolDataTypeCrescentPool":  5,
	"ProtocolDataTypeSifchainPool":  6,
}
```

#### Connection

```go
// ConnectionProtocolData defines state for connection tracking.
type ConnectionProtocolData struct {
	ConnectionID string
	ChainID      string
	LastEpoch    int64
	Prefix       string
}
```

#### Liquid

```go
// LiquidAllowedDenomProtocolData defines protocol state to track off-chain
// liquid qAssets.
type LiquidAllowedDenomProtocolData struct {
	// The chain on which the qAssets reside currently.
	ChainID string
	// The chain for which the qAssets were issued.
	RegisteredZoneChainID string
	// The IBC denom.
	IbcDenom string
	// The qAsset denom.
	QAssetDenom string
}
```

#### Osmosis

```go
// OsmosisPoolProtocolData defines protocol state to track qAssets locked in
// Osmosis pools.
type OsmosisPoolProtocolData struct {
	PoolID      uint64
	PoolName    string
	LastUpdated time.Time
	PoolData    json.RawMessage
	PoolType    string
	Zones       map[string]string // chainID: IBC/denom
}

type OsmosisParamsProtocolData struct {
	ChainID string
}
```

## Messages

Description of message types that trigger state transitions;

```protobuf
// Msg defines the participationrewards Msg service.
service Msg {
  rpc SubmitClaim(MsgSubmitClaim) returns (MsgSubmitClaimResponse) {
    option (google.api.http) = {
      post : "/quicksilver/tx/v1/participationrewards/claim"
      body : "*"
    };
  };
}
```

### MsgSubmitClaim

SubmitClaim is used to verify, by proof, that the given user address has a claim against the given asset type for the given zone.

```go
// MsgSubmitClaim represents a message type for submitting a participation
// claim regarding the given zone (chain).
type MsgSubmitClaim struct {
	UserAddress string          `protobuf:"bytes,1,opt,name=user_address,proto3" json:"user_address,omitempty"`
	Zone        string          `protobuf:"bytes,2,opt,name=zone,proto3" json:"zone,omitempty"`
	SrcZone     string          `protobuf:"bytes,3,opt,name=src_zone,proto3" json:"src_zone,omitempty"`
	ClaimType   types.ClaimType `protobuf:"varint,4,opt,name=claim_type,proto3,enum=quicksilver.claimsmanager.v1.ClaimType" json:"claim_type,omitempty"`
	Proofs      []*types.Proof  `protobuf:"bytes,5,rep,name=proofs,proto3" json:"proofs,omitempty"`
}
```

* **UserAddress** - the address of the claimant account on the native chain;
* **Zone** - the native zone related to the qAsset;
* **SrcZone** - the zone on which the qAsset is used (from where the proof originates);
* **ClaimType** - see [`x/claimsmanager/spec/README.md#ClaimType`](../../claimsmanager/spec/README.md#ClaimType);
* **Proofs** - see [`x/claimsmanager/spec/README.md#Proof`](../../claimsmanager/spec/README.md#Proof);

**Transaction**: [`claim`](#claim)

## Transactions

Description of transactions that collect messages in specific contexts to trigger state transitions;

### claim

Submit proof of assets held in the given zone.

`claim [zone] [src-zone] [claim-type] [payload-file].json`

## Proposals

### add-protocol-data

Submit an add protocol data proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

`add-protocol-data [proposal-file]`

Example:

`quicksilverd tx gov submit-proposal add-protocol-data <path/to/proposal.json> --from=<key_or_address>`

Where proposal.json contains:

```json
{
  "title": "Add Osmosis Atom/qAtom Pool",
  "description": "Add Osmosis Atom/qAtom Pool to support participation rewards",
  "protocol": "osmosis",
  "key": "pools/XXX",
  "type": "osmosispool",
  "data": {
	"poolID": "596",
	"ibcToken": "27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
	"localDenom": "uqatom"
  },
  "deposit": "512000000uqck"
}
```

## Events

N/A

[? this should probably emit some events for monitoring and tracking purposes ?]

## Hooks

N/A

## Queries

Participation Rewards module provides the below queries to check the module's state:

```protobuf
// Query provides defines the gRPC querier service.
service Query {
  // Params returns the total set of participation rewards parameters.
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get =
        "/quicksilver/participationrewards/v1/params";
  }

  rpc ProtocolData(QueryProtocolDataRequest)
      returns (QueryProtocolDataResponse) {
    option (google.api.http).get =
        "/quicksilver/participationrewards/v1/protocoldata/{type}/{key}";
  }
}
```

### params

Query the current airdrop module parameters.

```go
// QueryParamsRequest is the request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

// QueryParamsResponse is the response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params defines the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}
```

### protocoldata

Query the specified protocol data.

```go
// QueryProtocolDataRequest is the request type for querying Protocol Data.
type QueryProtocolDataRequest struct {
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Key  string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

// QueryProtocolDataResponse is the response type for querying Protocol Data.
type QueryProtocolDataResponse struct {
	// data defines the underlying protocol data.
	Data []encoding_json.RawMessage `protobuf:"bytes,1,rep,name=data,proto3,casttype=encoding/json.RawMessage" json:"data,omitempty" yaml:"data"`
}
```

## Keepers

<https://pkg.go.dev/github.com/ingenuity-build/quicksilver/x/participationrewards/keeper>

## Parameters

Module parameters:

| Key                                                     | Type         | Example |
|:--------------------------------------------------------|:-------------|:--------|
| distribution_proportions.validator_selection_allocation | string (dec) | "0.34"  |
| distribution_proportions.holdings_allocation            | string (dec) | "0.33"  |
| distribution_proportions.lockup_allocation              | string (dec) | "0.33"  |

Description of parameters:

* `validator_selection_allocation` - the percentage of inflation rewards allocated to validator selection rewards;
* `holdings_allocation` - the percentage of inflation rewards allocated to qAssets hoildings rewards;
* `lockup_allocation` - the percentage of inflation rewards allocated to staking and locking of QCK;

## Begin Block

N/A

## End Block

N/A

## After Epoch End

The following is performed at the end of every epoch:

* Obtains the rewards allocations according to the module balances and
  distribution proportions parameters;
* Allocate zone rewards according to the proportional zone Total Value Locked
  (TVL) for both **Validator Selection** and **qAsset Holdings**;
* Calculate validator selection scores and allocations for every zone:
  1. Obtain performance account delegation rewards (`performanceScores`);
  2. Calculate decentralization scores (`distributionScores`);
  3. Calculate overall validator scores;
  4. Calculate user validator selection rewards;
  5. Distribute validator selection rewards;
* Calculate qAsset holdings:
  1. Obtain qAssets held by account (locally and off-chain via claims / Proof of
     Posession);
  2. Calculate user proportion (cap at 2%);
  3. Normalize and distribute allocation;
* Allocate lockup rewards by sending portion to `feeCollector` for distribution
  by Staking Module;
* Update protocol data with the epoch boundary block height;
* Update osmosis pools protocol data;

## IBC

### Messages, Acknowledgements & Handlers

### Queries, Requests & Callbacks

This module registeres the following queries, requests and callbacks.

#### Performance Delegation Rewards

Queries the performance delegation rewards of the zone and computes the
validator scores based on the performance rewards.

* **Query:** `cosmos.distribution.v1beta1.Query/DelegationTotalRewards`
* **Callback:** `ValidatorSelectionRewardsCallback`

#### Osmosis Pool Update

Updates the registered Osmosis pools at the end of each epoch.

* **Query:** `store/gamm/key`
* **Callback:** `OsmosisPoolUpdateCallback`

#### Epoch Block

Queries and records the block height of the registered zone at the epoch
boundary.

* **Query:** `cosmos.base.tendermint.v1beta1.Service/GetLatestBlock`
* **Callback:** `SetEpochBlockCallback`
