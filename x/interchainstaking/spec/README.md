# Interchain Staking

## Abstract

Module, `x/interchainstaking`, defines and implements the core Quicksilver
protocol.

## Contents

1. [Concepts](#Concepts)
1. [State](#State)
1. [Messages](#Messages)
1. [Transactions](#Transactions)
1. [Proposals](#Proposals)
1. [Events](#Events)
1. [Hooks](#Hooks)
1. [Queries](#Queries)
1. [Keepers](#Keepers)
1. [Parameters](#Parameters)
1. [Begin Block](#Begin-Block)
1. [After Epoch End](#After-Epoch-End)
1. [IBC](#IBC)

## Concepts

Key concepts, mechanisms and core logic for the module;

### Registered Zone

A `RegisteredZone` is the core record kept by the Interchain Staking module. It
is created when a `CreateRegisteredZone` proposal has been passed by
governance.

It keeps track of the chain ID, and related IBC connection, Interchain Accounts
managed by Quicksilver of the host chain (for the purposes of Deposit,
Delegation, Withdrawal and Performance monitoring), the host chain's
validatorset, `aggregate_intent` and the `redemption_rate` used when minting
and burning qAssets for native assets.

### Redemption Rate

The `redemption_rate` is the ratio between qAsset supply, tracked by the Bank
module, and the total number of native assets staked against a given zone. This
ratio is used to determine how many qAssets to mint when staking with
Quicksilver, to ensure everyone joins the pool with the correct proportion.
Additionally, the previous epoch's redemption_rate is tracked, and used to
calculate the number of tokens to unbond when redeeming qAssets. The minimum of
the the current and last rates are used to negate an attack whereby a user can
enter the protocol immediately before the epoch, and exit immediately after to
claim a disproportionate amount of rewards for a very short exposure to the
protocol.

### Intent Signalling

Intent Signalling is the mechanism by which users of the protocol are able to
Signal to which validators they will their proportion of the stake pool is
delegated. In order to maintain fungibility of qAssets, we must pool assets and
delegate them as a single entity. Users are able to signal to which validators,
and with what weightings, they wish their proportion of stake to be delegated.
This is aggregated on an epochly basis.

### Aggregate Intent

The Aggregate Intent is calculated epochly, based upon the Signaled Intent from
each user, and the weight given to that intent based upon the assets the user
holds (this information is drawn from the Claimsmanager module). This aggregate
itnent is used as a target, for the protocol to use when determining where to
allocate assets during delegation, rebalance and undelegation processes.

### Interchain Accounts

## State

### Zone

A `Zone` represents a Cosmos based blockchain that integrates with the
Quicksilver protocol via Interchain Accounts (ICS) and Interblockchain
Communication (IBC).

```go
type Zone struct {
	ConnectionId                 string                                 `protobuf:"bytes,1,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	ChainId                      string                                 `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	DepositAddress               *ICAAccount                            `protobuf:"bytes,3,opt,name=deposit_address,json=depositAddress,proto3" json:"deposit_address,omitempty"`
	WithdrawalAddress            *ICAAccount                            `protobuf:"bytes,4,opt,name=withdrawal_address,json=withdrawalAddress,proto3" json:"withdrawal_address,omitempty"`
	PerformanceAddress           *ICAAccount                            `protobuf:"bytes,5,opt,name=performance_address,json=performanceAddress,proto3" json:"performance_address,omitempty"`
	DelegationAddress            *ICAAccount                            `protobuf:"bytes,6,opt,name=delegation_address,json=delegationAddress,proto3" json:"delegation_address,omitempty"`
	AccountPrefix                string                                 `protobuf:"bytes,7,opt,name=account_prefix,json=accountPrefix,proto3" json:"account_prefix,omitempty"`
	LocalDenom                   string                                 `protobuf:"bytes,8,opt,name=local_denom,json=localDenom,proto3" json:"local_denom,omitempty"`
	BaseDenom                    string                                 `protobuf:"bytes,9,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	RedemptionRate               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=redemption_rate,json=redemptionRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"redemption_rate"`
	LastRedemptionRate           github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,11,opt,name=last_redemption_rate,json=lastRedemptionRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"last_redemption_rate"`
	Validators                   []*Validator                           `protobuf:"bytes,12,rep,name=validators,proto3" json:"validators,omitempty"`
	AggregateIntent              ValidatorIntents                       `protobuf:"bytes,13,rep,name=aggregate_intent,json=aggregateIntent,proto3,castrepeated=ValidatorIntents" json:"aggregate_intent,omitempty"`
	MultiSend                    bool                                   `protobuf:"varint,14,opt,name=multi_send,json=multiSend,proto3" json:"multi_send,omitempty"`
	LiquidityModule              bool                                   `protobuf:"varint,15,opt,name=liquidity_module,json=liquidityModule,proto3" json:"liquidity_module,omitempty"`
	WithdrawalWaitgroup          uint32                                 `protobuf:"varint,16,opt,name=withdrawal_waitgroup,json=withdrawalWaitgroup,proto3" json:"withdrawal_waitgroup,omitempty"`
	IbcNextValidatorsHash        []byte                                 `protobuf:"bytes,17,opt,name=ibc_next_validators_hash,json=ibcNextValidatorsHash,proto3" json:"ibc_next_validators_hash,omitempty"`
	ValidatorSelectionAllocation uint64                                 `protobuf:"varint,18,opt,name=validator_selection_allocation,json=validatorSelectionAllocation,proto3" json:"validator_selection_allocation,omitempty"`
	HoldingsAllocation           uint64                                 `protobuf:"varint,19,opt,name=holdings_allocation,json=holdingsAllocation,proto3" json:"holdings_allocation,omitempty"`
	LastEpochHeight              int64                                  `protobuf:"varint,20,opt,name=last_epoch_height,json=lastEpochHeight,proto3" json:"last_epoch_height,omitempty"`
	Tvl                          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,21,opt,name=tvl,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"tvl"`
	UnbondingPeriod              int64                                  `protobuf:"varint,22,opt,name=unbonding_period,json=unbondingPeriod,proto3" json:"unbonding_period,omitempty"`
}
```

- **ConnectionId** - remote zone connection identifier;
- **ChainId** - remote zone identifier;
- **DepositAddress** - remote zone deposit address;
- **WithdrawalAddress** - remote zone withdrawal address;
- **PerformanceAddress** - remote zone performance address (each validator gets
  an exact equal delegation from this account to measure performance);
- **DelegationAddresses** - remote zone delegation addresses to represent
  granular voting power;
- **AccountPrefix** - remote zone account address prefix;
- **LocalDenom** - protocol denomination (qAsset), e.g. uqatom;
- **BaseDenom** - remote zone denomination (uStake), e.g. uatom;
- **RedemptionRate** - redemption rate between protocol qAsset and native
  remote asset;
- **LastRedemptionRate** - redemption rate as at previous epoch boundary
  (used to prevent epoch boundary gaming);
- **Validators** - list of validators on the remote zone;
- **AggregateIntent** - the aggregated delegation intent of the protocol for
  this remote zone. The map key is the corresponding validator address
  contained in the `ValidatorIntent`;
- **MultiSend** - multisend support on remote zone; deprecated;
- **LiquidityModule** - liquidity module enabled on remote zone;
- **WithdrawalWaitgroup** - tally of pending withdrawal transactions;
- **IbcNextValidatorHash** -
- **ValidatorSelectionAllocation** - proportional zone rewards allocation for
  validator selection;
- **HoldingsAllocation** - proportional zone rewards allocation for asset
  holdings;
- **LastEpochHeight** - the height of this chain at the last Quicksilver epoch
  boundary;
- **Tvl** - the Total Value Locked for this zone (in terms of Atom value);
- **UnbondingPeriod** - this zone's unbonding period;

### ICAAccount

An `ICAAccount` represents an account on an remote zone under the control of
the protocol.

```go
type ICAAccount struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// balance defines the different coins this balance holds.
	Balance           github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=balance,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"balance"`
	PortName          string                                   `protobuf:"bytes,3,opt,name=port_name,json=portName,proto3" json:"port_name,omitempty"`
	WithdrawalAddress string                                   `protobuf:"bytes,4,opt,name=withdrawal_address,json=withdrawalAddress,proto3" json:"withdrawal_address,omitempty"`
	BalanceWaitgroup  uint32                                   `protobuf:"varint,5,opt,name=balance_waitgroup,json=balanceWaitgroup,proto3" json:"balance_waitgroup,omitempty"`
}
```

- **Address** - the account address on the remote zone;
- **Balance** - the account balance on the remote zone;
- **PortName** - the port name to access the remote zone;
- **WithdrawalAddress** - the address withdrawals are sent to for this account;
- **BalanceWaitgroup** - the tally of pending balance query transactions sent
  to the remote zone;

### Distribution

```go
type Distribution struct {
	Valoper string `protobuf:"bytes,1,opt,name=valoper,proto3" json:"valoper,omitempty"`
	Amount  uint64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}
```

### WithdrawalRecord

```go
type WithdrawalRecord struct {
	ChainId        string                                   `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Delegator      string                                   `protobuf:"bytes,2,opt,name=delegator,proto3" json:"delegator,omitempty"`
	Distribution   []*Distribution                          `protobuf:"bytes,3,rep,name=distribution,proto3" json:"distribution,omitempty"`
	Recipient      string                                   `protobuf:"bytes,4,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount         github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,5,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	BurnAmount     github_com_cosmos_cosmos_sdk_types.Coin  `protobuf:"bytes,6,opt,name=burn_amount,json=burnAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"burn_amount"`
	Txhash         string                                   `protobuf:"bytes,7,opt,name=txhash,proto3" json:"txhash,omitempty"`
	Status         int32                                    `protobuf:"varint,8,opt,name=status,proto3" json:"status,omitempty"`
	CompletionTime time.Time                                `protobuf:"bytes,9,opt,name=completion_time,json=completionTime,proto3,stdtime" json:"completion_time"`
}
```

### UnbondingRecord

```go
type UnbondingRecord struct {
	ChainId       string   `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	EpochNumber   int64    `protobuf:"varint,2,opt,name=epoch_number,json=epochNumber,proto3" json:"epoch_number,omitempty"`
	Validator     string   `protobuf:"bytes,3,opt,name=validator,proto3" json:"validator,omitempty"`
	RelatedTxhash []string `protobuf:"bytes,4,rep,name=related_txhash,json=relatedTxhash,proto3" json:"related_txhash,omitempty"`
}
```

### RedelegationRecord

```go
type RedelegationRecord struct {
	ChainId        string    `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	EpochNumber    int64     `protobuf:"varint,2,opt,name=epoch_number,json=epochNumber,proto3" json:"epoch_number,omitempty"`
	Source         string    `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	Destination    string    `protobuf:"bytes,4,opt,name=destination,proto3" json:"destination,omitempty"`
	Amount         int64     `protobuf:"varint,5,opt,name=amount,proto3" json:"amount,omitempty"`
	CompletionTime time.Time `protobuf:"bytes,6,opt,name=completion_time,json=completionTime,proto3,stdtime" json:"completion_time"`
}
```

### TransferRecord

```go
type TransferRecord struct {
	Sender    string                                  `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Recipient string                                  `protobuf:"bytes,2,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount    github_com_cosmos_cosmos_sdk_types.Coin `protobuf:"bytes,3,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"amount"`
}
```

### Validator

`Validator` represents relevant meta data of a validator within a zone.

```go
type Validator struct {
	ValoperAddress  string                                 `protobuf:"bytes,1,opt,name=valoper_address,json=valoperAddress,proto3" json:"valoper_address,omitempty"`
	CommissionRate  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=commission_rate,json=commissionRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"commission_rate"`
	DelegatorShares github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=delegator_shares,json=delegatorShares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegator_shares"`
	VotingPower     github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=voting_power,json=votingPower,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"voting_power"`
	Score           github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=score,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"score"`
	Status          string                                 `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	Jailed          bool                                   `protobuf:"varint,7,opt,name=jailed,proto3" json:"jailed,omitempty"`
	Tombstoned      bool                                   `protobuf:"varint,8,opt,name=tombstoned,proto3" json:"tombstoned,omitempty"`
	JailedSince     time.Time                              `protobuf:"bytes,9,opt,name=jailed_since,json=jailedSince,proto3,stdtime" json:"jailed_since"`
}
```

- **ValoperAddress** - the validator address;
- **CommissionRate** - the validator commission rate;
- **DelegatorShares** -
- **VotingPower** - the validator voting power on the remote zone;
- **Score** - the validator Quicksilver protocol overall score;
- **Status** -
- **Jailed** - is this validator currently jailed;
- **Tombstoned** - is this validator tombstoned;
- **JailedSince** - blocktime timestamp when this validator was jailed;

### ValidatorIntent

`ValidatorIntent` represents the weighted delegation intent to a particular
validator.

```go
type ValidatorIntent struct {
	ValoperAddress string                                 `protobuf:"bytes,1,opt,name=valoper_address,proto3" json:"valoper_address,omitempty"`
	Weight         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=weight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"weight" yaml:"weight"`
}
```

- **ValoperAddress** - the remote zone validator address;
- **Weight** - the weight of intended delegation to this validator;

### DelegatorIntent

`DelegatorIntent` represents the current delegation intent for this zone.
Delegations are incrementally adjusted towards the `Zone.AggregateIntent`.

```go
type DelegatorIntent struct {
	Delegator string           `protobuf:"bytes,1,opt,name=delegator,proto3" json:"delegator,omitempty"`
	Intents   ValidatorIntents `protobuf:"bytes,2,rep,name=intents,proto3,castrepeated=ValidatorIntents" json:"intents,omitempty"`
}
```

- **Delegator** - the delegation account address on the remote zone;
- **Intents** - the delegation intents to individual validators on the remote
  zone;

### Delegation

`Delegation` represents the actual delegations made by
`RegisteredZone.DelegationAddresses` to validators on the remote zone;

```go
type Delegation struct {
	DelegationAddress string                                  `protobuf:"bytes,1,opt,name=delegation_address,json=delegationAddress,proto3" json:"delegation_address,omitempty"`
	ValidatorAddress  string                                  `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Amount            github_com_cosmos_cosmos_sdk_types.Coin `protobuf:"bytes,3,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"amount"`
	Height            int64                                   `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
	RedelegationEnd   int64                                   `protobuf:"varint,5,opt,name=redelegation_end,json=redelegationEnd,proto3" json:"redelegation_end,omitempty"`
}
```

- **DelegationAddress** - the delegator address on the remote zone;
- **ValidatorAddress** - the validator address on the remote zone;
- **Amount** - the amount delegated;
- **Height** - the block height at which the delegation occured;
- **RedelegationEnd** - ;

### PortConnectionTuple

```go
type PortConnectionTuple struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	PortId       string `protobuf:"bytes,2,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
}
```

### Receipt

```go
type Receipt struct {
	ChainId string                                   `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Sender  string                                   `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Txhash  string                                   `protobuf:"bytes,3,opt,name=txhash,proto3" json:"txhash,omitempty"`
	Amount  github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,4,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}
```

## Messages

```protobuf
// Msg defines the interchainstaking Msg service.
service Msg {
  // RequestRedemption defines a method for requesting burning of qAssets for
  // native assets.
  rpc RequestRedemption(MsgRequestRedemption)
      returns (MsgRequestRedemptionResponse) {
    option (google.api.http) = {
      post : "/quicksilver/tx/v1/interchainstaking/redeem"
      body : "*"
    };
  };
  // SignalIntent defines a method for signalling voting intent for one or more
  // validators.
  rpc SignalIntent(MsgSignalIntent) returns (MsgSignalIntentResponse) {
    option (google.api.http) = {
      post : "/quicksilver/tx/v1/interchainstaking/intent"
      body : "*"
    };
  };
}
```

### MsgRequestRedemption

Redeems the indicated qAsset coin amount from the protocol, converting the
qAsset back to the native asset at the appropriate redemption rate.

```go
// MsgRequestRedemption represents a message type to request a burn of qAssets
// for native assets.
type MsgRequestRedemption struct {
	Value              types.Coin `protobuf:"bytes,1,opt,name=value,proto3" json:"value" yaml:"coin"`
	DestinationAddress string     `protobuf:"bytes,2,opt,name=destination_address,json=destinationAddress,proto3" json:"destination_address,omitempty"`
	FromAddress        string     `protobuf:"bytes,3,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
}
```

- **Value** - qAsset as standard cosmos sdk cli coin string, {amount}{denomination};
- **DestinationAddress** - standard cosmos sdk bech32 address string;
- **FromAddress** - standard cosmos sdk bech32 address string;

**Transaction**: [`redeem`](#redeem)

### MsgSignalIntent

Signal validator delegation intent for a given zone by weight.

```go
// MsgSignalIntent represents a message type for signalling voting intent for
// one or more validators.
type MsgSignalIntent struct {
	ChainId     string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
	Intents     string `protobuf:"bytes,2,opt,name=intents,proto3" json:"intents,omitempty" yaml:"intents"`
	FromAddress string `protobuf:"bytes,3,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
}
```

- **ChainId** - zone identifier string;
- **Intents** - list of validator intents according to weight;
- **FromAddress** - standard cosmos sdk bech32 address string;

**Transaction**: [`signal-intent`](#signal-intent)

## Transactions

### signal-intent

Signal validator delegation intent by providing a comma seperated string
containing a decimal weight and the bech32 validator address.

`quicksilverd signal-intent [chain_id] [delegation_intent]`

Example:

`quicksilverd signal-intent cosmoshub-4 0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz`

### redeem

Redeem qAssets for native tokens.

`quicksilverd redeem [coins] [destination_address]`

Example:

`quicksilverd redeem 2500000uatom cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w`

## Proposals

### register-zone

Submit a zone registration proposal.

`quicksilverd register-zone [proposal-file]`

The proposal must include an initial deposit and the details must be provided
as a json file, e.g.

```json
{
  "title": "Register cosmoshub-4",
  "description": "Onboard the cosmoshub-4 zone to Quicksilver",
  "connection_id": "connection-3",
  "base_denom": "uatom",
  "local_denom": "uqatom",
  "account_prefix": "cosmos",
  "multi_send": true,
  "liquidity_module": false,
  "deposit": "512000000uqck"
}
```

### update-zone

Submit a zone update proposal.

`quicksilverd update-zone [proposal-file]`

The proposal must include a deposit and the details must be provided as a json
file, e.g.

```json
{
  "title": "Enable liquidity module for cosmoshub-4",
  "description": "Update cosmoshub-4 to enable liquidity module",
  "chain_id": "cosmoshub-4",
  "changes": [
    {
      "key": "liquidity_module",
      "value": "true"
    }
  ],
  "deposit": "512000000uqck"
}
```

## Events

Events emitted by module for tracking messages and index transactions;

### RegisterZone

| Type          | Attribute Key | Attribute Value   |
| :------------ | :------------ | :---------------- |
| message       | module        | interchainstaking |
| register_zone | connection_id | {connection_id}   |
| register_zone | chain_id      | {chain_id}        |

### MsgClaim

| Type               | Attribute Key | Attribute Value   |
| :----------------- | :------------ | :---------------- |
| message            | module        | interchainstaking |
| request_redemption | burn_amount   | {burn_amount}     |
| request_redemption | redeem_amount | {redeem_amount}   |
| request_redemption | recipient     | {recipient}       |
| request_redemption | chain_id      | {chain_id}        |
| request_redemption | connection_id | {connection_id}   |

## Hooks

N/A

## Queries

```protobuf
service Query {
  // ZoneInfos provides meta data on connected zones.
  rpc ZoneInfos(QueryZonesInfoRequest) returns (QueryZonesInfoResponse) {
    option (google.api.http).get = "/quicksilver/interchainstaking/v1/zones";
  }
  // DepositAccount provides data on the deposit address for a connected zone.
  rpc DepositAccount(QueryDepositAccountForChainRequest)
      returns (QueryDepositAccountForChainResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/deposit_address";
  }
  // DelegatorIntent provides data on the intent of the delegator for the given
  // zone.
  rpc DelegatorIntent(QueryDelegatorIntentRequest)
      returns (QueryDelegatorIntentResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/delegator_intent/"
        "{delegator_address}";
  }

  // Delegations provides data on the delegations for the given zone.
  rpc Delegations(QueryDelegationsRequest) returns (QueryDelegationsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/delegations";
  }

   // Delegations provides data on the delegations for the given zone.
   rpc Receipts(QueryReceiptsRequest) returns (QueryReceiptsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/receipts";
  }

  // WithdrawalRecords provides data on the active withdrawals.
  rpc ZoneWithdrawalRecords(QueryWithdrawalRecordsRequest)
      returns (QueryWithdrawalRecordsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/withdrawal_records/{delegator_address}";
  }

  // WithdrawalRecords provides data on the active withdrawals.
  rpc WithdrawalRecords(QueryWithdrawalRecordsRequest)
      returns (QueryWithdrawalRecordsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/withdrawal_records";
  }

  // UnbondingRecords provides data on the active unbondings.
  rpc UnbondingRecords(QueryUnbondingRecordsRequest)
      returns (QueryUnbondingRecordsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/unbonding_records";
  }

  // RedelegationRecords provides data on the active unbondings.
  rpc RedelegationRecords(QueryRedelegationRecordsRequest)
      returns (QueryRedelegationRecordsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/redelegation_records";
  }
}
```

### zones

Query registered zones.

`quicksilverd query interchainstaking zones`

Example response:

```yaml
pagination:
  next_key: null
  total: "2"
zones:
- account_prefix: cosmos
  aggregate_intent: {}
  base_denom: uatom
  chain_id: lstest-1
  connection_id: connection-0
  delegation_addresses:
  - address: cosmos12hww50r7q7xyhspt72c9c8n3uyknqhv208sxuq9mcqdqjv0mcreq62maa2
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "25083333"
      denom: uatom
    port_name: icacontroller-lstest-1.delegate.9
  ...
  deposit_address:
    address: cosmos146xjrj2tass9fvtcw30dtl9v8f4t26z7cjxlxw0paxyyxmx2hqcq73vk6p
    balance:
    - amount: "25000000"
      denom: cosmosvaloper16pxh2v4hr28h2gkntgfk8qgh47pfmjfhvcamkc3
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uatom
    port_name: icacontroller-lstest-1.deposit
  holdings_allocation: []
  ibc_next_validators_hash: Qn4t+8M6bod6ewSYwnPScdWwwbSE7mc47GlMpuo15d0=
  last_redemption_rate: "1.000000000000000000"
  liquidity_module: true
  local_denom: uqatom
  multi_send: true
  performance_address:
    address: cosmos1yp64sfc5d4g4xtemptachyd2jaraxz8c5vptp7swgnv86l3ll3yqzz72wk
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uatom
    port_name: icacontroller-lstest-1.performance
  redemption_rate: "1.000000000000000000"
  validator_selection_allocation: []
  validators:
  - commission_rate: "0.030000000000000000"
    delegator_shares: "4000093333.000000000000000000"
    score: "0.000000000000000000"
    valoper_address: cosmosvaloper12evgzwsc2av7nfc5x7p74g9ppmfwm30xug6pwv
    voting_power: "4000093333"
  ...
  withdrawal_address:
    address: cosmos1w7x78xu4ms3qwspryl8jjy57l3esns8ayh6mj9g3544wmgnfnzrs86lr9p
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uatom
    port_name: icacontroller-lstest-1.withdrawal
  withdrawal_waitgroup: 12
- account_prefix: osmo
  aggregate_intent: {}
  base_denom: uosmo
  chain_id: lstest-2
  connection_id: connection-1
  delegation_addresses: []
  deposit_address:
    address: osmo14s68pery7n8s9cm6lzwxv4s0ppucctv28fcmtqg852965hfgpuvsmq5edm
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uosmo
    port_name: icacontroller-lstest-2.deposit
  holdings_allocation: []
  ibc_next_validators_hash: TyRByZjTIrfQ81mMlvoRyg1crPx4Kk9Lur+Kkty06h8=
  last_redemption_rate: "1.000000000000000000"
  liquidity_module: true
  local_denom: uqosmo
  multi_send: true
  performance_address:
    address: osmo1v799s4plwuyux8xunmzxcw6y2g8t5u373ravsy2zxg7k0x8g7pdsaa8ve9
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uosmo
    port_name: icacontroller-lstest-2.performance
  redemption_rate: "1.000000000000000000"
  validator_selection_allocation: []
  validators: []
  withdrawal_address:
    address: osmo1xcs9r3ssmjndgr09jww29cs60ygck8g3udyl6savg3nkercfhl2qtp3lwv
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uosmo
    port_name: icacontroller-lstest-2.withdrawal
  withdrawal_waitgroup: 0
```

### intent

Query delegation intent for a given chain.

`quicksilverd query interchainstaking intent [chain_id] [delegator_addr]`

### deposit-account

Query deposit account address for a given chain.

`quicksilverd query interchainstaking deposit-account [chain_id]`

## Keepers

https://pkg.go.dev/github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper

## Parameters

Module parameters:

| Key                   | Type    | Default |
| :-------------------- | :------ | :------ |
| deposit_interval      | uint64  | 20      |
| validatorset_interval | uint64  | 200     |
| commission_rate       | sdk.Dec | "0.025" |
| unbonding_enabled     | bool    | false   |

Description of parameters:

- `deposit_interval` - monitoring and handling interval of registered zones' deposit accounts;
- `validatorset_interval` - monitoring and updating interval of registered zones' validator sets;
- `commission_rate` - default commission rate for Quicksilver validators;
- `unbonding_enabled` - flag to indicate if unbondings are enabled for the Quicksilver protocol;

## Begin Block

Iterate through all registered zones and check validator set status. If the
status has changed, requery the validator set and update zone state.

## After Epoch End

The following is performed at the end of every epoch for each registered zone:

- Aggregate Intents:
  1. Iterate through all stored instances of `DelegatorIntent` for each zone
     and obtain the **delegator account balance**;
  2. Compute the **base balance** using the account balance and `RedpemtionRate`;
  3. Ordinalize the delegator's validator intents by `Weight`;
  4. Set the zone `AggregateIntent` and update zone state;
- Query delegator delegations for each zone and update delegation records:
  1. Query delegator delegations `cosmos.staking.v1beta1.Query/DelegatorDelegations`;
  2. For each response (per delegator `DelegationsCallback`), verify every
     delegation record (via IBC `DelegationCallback`) and update delegation
     record accordingly (add, update or remove);
  3. Update validator set;
  4. Update zone;
- Withdraw delegation rewards for each zone and distribute:

  1. Query delegator rewards `cosmos.distribution.v1beta1.Query/DelegationTotalRewards`;
  2. For each response (per delegator `RewardsCallback`), send withdrawal
     messages for each of its validator delegations and add tally to
     `WithdrawalWaitgroup`;
  3. For each IBC acknowledgement decrement the `WithdrawalWaitgroup`. Once
     all responses are collected (`WithdrawalWaitgroup == 0`) query the balance
     of `WithdrawalAddress` (`cosmos.bank.v1beta1.Query/AllBalances`), then
     distribute rewards (`DistributeRewardsFromWithdrawAccount`).

     This approach ensures the exact rewards amount is known at the time of
     distribution.

## IBC

### Messages, Acknowledgements & Handlers

#### MsgWithdrawDelegatorReward

Triggered at the end of every epoch if delegator accounts have accrued rewards.
Collects rewards to zone withdrawal account `WithdrawalAddress` and distributes
rewards once all delegator rewards withdrawals have been acknowledged.

- **Endpoint:** `/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward`
- **Handler:** `HandleWithdrawRewards`

#### MsgRedeemTokensforShares

Triggered during execution of `Delegate` for delegation allocations that are
not in the zone `BaseDenom`. During callback the relevant delegation record is
updated.

- **Endpoint:** `/cosmos.staking.v1beta1.MsgRedeemTokensforShares`
- **Handler:** `HandleRedeemTokens`

#### MsgTokenizeShares

Triggered by `RequestRedemption` when a user redeems qAssets. Withdrawal
records are set or updated accordingly.  
See [MsgRequestRedemption](#MsgRequestRedemption).

- **Endpoint:** `/cosmos.staking.v1beta1.MsgTokenizeShares`
- **Handler:** `HandleTokenizedShares`

#### MsgDelegate

Triggered by `Delegate` whenever delagtions are made by the protocol to zone
validators. `HandleDelegate` distinguishes `DelegationAddresses` and updates
delegation records for these delegation accounts.

- **Endpoint:** `/cosmos.staking.v1beta1.MsgDelegate`
- **Handler:** `HandleDelegate`

#### MsgBeginRedelegate

Triggered at the end of every epoch during `Rebalance`.

- **Endpoint:** `/cosmos.staking.v1beta1.MsgBeginRedelegate`
- **Handler:** `HandleBeginRedelegate`

#### MsgSend

Triggered by `TransferToDelegate` during `HandleReceiptTransaction`.  
See [Deposit Interval](#Deposit-Interval).

- **Endpoint:** `/cosmos.bank.v1beta1.MsgSend`
- **Handler:** `HandleCompleteSend`

`HandleCompleteSend` executes one of the following options based on the
`FromAddress` and `ToAddress` of the msg:

1. **Delegate rewards accoring to global intents.**  
   (If `FromAddress` is the zone's `WithdrawalAddress`);
2. **Withdraw native assets for user.**  
   (If `FromAddress` is one of zone's `DelegationAddresses`);
3. **Delegate amount according to delegation plan.**  
   (If `FromAddress` is `DepositAddress` and `ToAddress` is one of zone's `DelegationAddresses`);

#### MsgSetWithdrawAddress

Triggered during zone initialization for every `DelegationAddresses` and
for the `PerformanceAddress`. The purpose of using a dedicated withdrawal
account allows for accurate rewards withdrawal accounting, that would otherwise
be impossible as the rewards amount will only be known at the time the msg is
triggered, and not at the time it was executed by the remote zone (due to network
latency and different zone block times, etc).

- **Endpoint:** `/cosmos.distribution.v1beta1.MsgSetWithdrawAddress`
- **Handler:** `HandleUpdatedWithdrawAddress`

#### MsgTransfer

Triggered by `DistributeRewardsFromWithdrawAccount` to distribute rewards
across the zone delegation accounts and collect fees to the module fee account.
The `RedemptionRate` is updated accordingly.  
See [WithdrawalAddress Balances](#WithdrawalAddress-Balances).

- **Endpoint:** `/ibc.applications.transfer.v1.MsgTransfer`
- **Handler:** `HandleMsgTransfer`

### Queries, Requests & Callbacks

This module registeres the following queries, requests and callbacks.

#### DepositAddress Balances

For every registered zone a periodic `AllBalances` query is run against the
`DepositAddress`. The query is proven by utilizing provable KV queries that
update the individual account balances `AccountBalanceCallback`, trigger the
`depositInterval` and finally update the zone state.

- **Query:** `cosmos.bank.v1beta1.Query/AllBalances`
- **Callback:** `AllBalancesCallback`

#### Delegator Delegations

Query delegator delegations for each zone and update delegation records.  
See [After Epoch End](#After-Epoch-End).

- **Query:** `cosmos.staking.v1beta1.Query/DelegatorDelegations`
- **Callback:** `DelegationsCallback`

#### Delegate Total Rewards

Withdraw delegation rewards for each zone and distribute.  
See [After Epoch End](#After-Epoch-End).

- **Query:** `cosmos.distribution.v1beta1.Query/DelegationTotalRewards`
- **Callback:** `RewardsCallback`

#### WithdrawalAddress Balances

Triggered by `HandleWithdrawRewards`.  
See [MsgWithdrawDelegatorReward](#MsgWithdrawDelegatorReward).

- **Query:** `cosmos.bank.v1beta1.Query/AllBalances`
- **Callback:** `DistributeRewardsFromWithdrawAccount`

#### Deposit Interval

Monitors transaction events of the zone `DepositAddress` on the remote chain
for receipt transactions that are then handled by `HandleReceiptTransaction`.
On valid receipts the delegation intent is updated (`UpdateIntent`) and new
qAssets minted and transferred to the sender (`MintQAsset`). A delegation
plan is computed (`DeterminePlanForDelegation`) and then executed
(`TransferToDelegate`). Successfully executed receipts are recorded to state.

- **Query:** `cosmos.tx.v1beta1.Service/GetTxsEvent`
- **Callback:** `DepositIntervalCallback`

#### Performance Balance Query

Triggered at zone registration when the zone performance account
`PerformanceAddress` is created. It monitors the performance account balance
until sufficient funds are available to execute the performance delegations.  
See [x/participationrewards/spec](../../participationrewards/spec/README.md).

- **Query:** `cosmos.bank.v1beta1.Query/AllBalances`
- **Callback:** `PerfBalanceCallback`

#### Validator Set Query

An essential query to ensure that the registred zone state accurately reflects
the validator set of the remote zone for bonded, unbonded and unbonding
validators.

- **Query:** `cosmos.staking.v1beta1.Query/Validators`
- **Callback:** `ValsetCallback`
