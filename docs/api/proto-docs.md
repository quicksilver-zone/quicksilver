<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [quicksilver/epochs/v1/genesis.proto](#quicksilver/epochs/v1/genesis.proto)
    - [EpochInfo](#quicksilver.epochs.v1.EpochInfo)
    - [GenesisState](#quicksilver.epochs.v1.GenesisState)
  
- [quicksilver/epochs/v1/query.proto](#quicksilver/epochs/v1/query.proto)
    - [QueryCurrentEpochRequest](#quicksilver.epochs.v1.QueryCurrentEpochRequest)
    - [QueryCurrentEpochResponse](#quicksilver.epochs.v1.QueryCurrentEpochResponse)
    - [QueryEpochsInfoRequest](#quicksilver.epochs.v1.QueryEpochsInfoRequest)
    - [QueryEpochsInfoResponse](#quicksilver.epochs.v1.QueryEpochsInfoResponse)
  
    - [Query](#quicksilver.epochs.v1.Query)
  
- [quicksilver/interchainquery/v1/genesis.proto](#quicksilver/interchainquery/v1/genesis.proto)
    - [DataPoint](#quicksilver.interchainquery.v1.DataPoint)
    - [GenesisState](#quicksilver.interchainquery.v1.GenesisState)
    - [Query](#quicksilver.interchainquery.v1.Query)
  
- [quicksilver/interchainquery/v1/messages.proto](#quicksilver/interchainquery/v1/messages.proto)
    - [MsgSubmitQueryResponse](#quicksilver.interchainquery.v1.MsgSubmitQueryResponse)
    - [MsgSubmitQueryResponseResponse](#quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse)
  
    - [Msg](#quicksilver.interchainquery.v1.Msg)
  
- [quicksilver/interchainstaking/v1/genesis.proto](#quicksilver/interchainstaking/v1/genesis.proto)
    - [Delegation](#quicksilver.interchainstaking.v1.Delegation)
    - [DelegationPlan](#quicksilver.interchainstaking.v1.DelegationPlan)
    - [DelegationPlan.DelegationPlanItem](#quicksilver.interchainstaking.v1.DelegationPlan.DelegationPlanItem)
    - [DelegationPlan.ValueEntry](#quicksilver.interchainstaking.v1.DelegationPlan.ValueEntry)
    - [DelegatorIntent](#quicksilver.interchainstaking.v1.DelegatorIntent)
    - [DistributionPlan](#quicksilver.interchainstaking.v1.DistributionPlan)
    - [DistributionPlan.ValueEntry](#quicksilver.interchainstaking.v1.DistributionPlan.ValueEntry)
    - [GenesisState](#quicksilver.interchainstaking.v1.GenesisState)
    - [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount)
    - [Params](#quicksilver.interchainstaking.v1.Params)
    - [PortConnectionTuple](#quicksilver.interchainstaking.v1.PortConnectionTuple)
    - [Receipt](#quicksilver.interchainstaking.v1.Receipt)
    - [RegisteredZone](#quicksilver.interchainstaking.v1.RegisteredZone)
    - [RegisteredZone.AggregateIntentEntry](#quicksilver.interchainstaking.v1.RegisteredZone.AggregateIntentEntry)
    - [RegisteredZone.DelegatorIntentEntry](#quicksilver.interchainstaking.v1.RegisteredZone.DelegatorIntentEntry)
    - [TransferRecord](#quicksilver.interchainstaking.v1.TransferRecord)
    - [Validator](#quicksilver.interchainstaking.v1.Validator)
    - [ValidatorIntent](#quicksilver.interchainstaking.v1.ValidatorIntent)
    - [WithdrawalRecord](#quicksilver.interchainstaking.v1.WithdrawalRecord)
  
- [quicksilver/interchainstaking/v1/messages.proto](#quicksilver/interchainstaking/v1/messages.proto)
    - [MsgRegisterZone](#quicksilver.interchainstaking.v1.MsgRegisterZone)
    - [MsgRegisterZoneResponse](#quicksilver.interchainstaking.v1.MsgRegisterZoneResponse)
    - [MsgRequestRedemption](#quicksilver.interchainstaking.v1.MsgRequestRedemption)
    - [MsgRequestRedemptionResponse](#quicksilver.interchainstaking.v1.MsgRequestRedemptionResponse)
    - [MsgSignalIntent](#quicksilver.interchainstaking.v1.MsgSignalIntent)
    - [MsgSignalIntentResponse](#quicksilver.interchainstaking.v1.MsgSignalIntentResponse)
  
    - [Msg](#quicksilver.interchainstaking.v1.Msg)
  
- [quicksilver/interchainstaking/v1/query.proto](#quicksilver/interchainstaking/v1/query.proto)
    - [QueryDelegationsRequest](#quicksilver.interchainstaking.v1.QueryDelegationsRequest)
    - [QueryDelegationsResponse](#quicksilver.interchainstaking.v1.QueryDelegationsResponse)
    - [QueryDelegatorDelegationsRequest](#quicksilver.interchainstaking.v1.QueryDelegatorDelegationsRequest)
    - [QueryDelegatorDelegationsResponse](#quicksilver.interchainstaking.v1.QueryDelegatorDelegationsResponse)
    - [QueryDelegatorIntentRequest](#quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest)
    - [QueryDelegatorIntentResponse](#quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse)
    - [QueryDepositAccountForChainRequest](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest)
    - [QueryDepositAccountForChainResponse](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse)
    - [QueryRegisteredZonesInfoRequest](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoRequest)
    - [QueryRegisteredZonesInfoResponse](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoResponse)
    - [QueryValidatorDelegationsRequest](#quicksilver.interchainstaking.v1.QueryValidatorDelegationsRequest)
    - [QueryValidatorDelegationsResponse](#quicksilver.interchainstaking.v1.QueryValidatorDelegationsResponse)
  
    - [Query](#quicksilver.interchainstaking.v1.Query)
  
- [quicksilver/mint/v1beta1/mint.proto](#quicksilver/mint/v1beta1/mint.proto)
    - [DistributionProportions](#quicksilver.mint.v1beta1.DistributionProportions)
    - [Minter](#quicksilver.mint.v1beta1.Minter)
    - [Params](#quicksilver.mint.v1beta1.Params)
  
- [quicksilver/mint/v1beta1/genesis.proto](#quicksilver/mint/v1beta1/genesis.proto)
    - [GenesisState](#quicksilver.mint.v1beta1.GenesisState)
  
- [quicksilver/mint/v1beta1/query.proto](#quicksilver/mint/v1beta1/query.proto)
    - [QueryEpochProvisionsRequest](#quicksilver.mint.v1beta1.QueryEpochProvisionsRequest)
    - [QueryEpochProvisionsResponse](#quicksilver.mint.v1beta1.QueryEpochProvisionsResponse)
    - [QueryParamsRequest](#quicksilver.mint.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#quicksilver.mint.v1beta1.QueryParamsResponse)
  
    - [Query](#quicksilver.mint.v1beta1.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="quicksilver/epochs/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/epochs/v1/genesis.proto



<a name="quicksilver.epochs.v1.EpochInfo"></a>

### EpochInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `start_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `current_epoch` | [int64](#int64) |  |  |
| `current_epoch_start_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `epoch_counting_started` | [bool](#bool) |  |  |
| `current_epoch_start_height` | [int64](#int64) |  |  |






<a name="quicksilver.epochs.v1.GenesisState"></a>

### GenesisState
GenesisState defines the epochs module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `epochs` | [EpochInfo](#quicksilver.epochs.v1.EpochInfo) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="quicksilver/epochs/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/epochs/v1/query.proto



<a name="quicksilver.epochs.v1.QueryCurrentEpochRequest"></a>

### QueryCurrentEpochRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |






<a name="quicksilver.epochs.v1.QueryCurrentEpochResponse"></a>

### QueryCurrentEpochResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current_epoch` | [int64](#int64) |  |  |






<a name="quicksilver.epochs.v1.QueryEpochsInfoRequest"></a>

### QueryEpochsInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="quicksilver.epochs.v1.QueryEpochsInfoResponse"></a>

### QueryEpochsInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `epochs` | [EpochInfo](#quicksilver.epochs.v1.EpochInfo) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="quicksilver.epochs.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `EpochInfos` | [QueryEpochsInfoRequest](#quicksilver.epochs.v1.QueryEpochsInfoRequest) | [QueryEpochsInfoResponse](#quicksilver.epochs.v1.QueryEpochsInfoResponse) | EpochInfos provide running epochInfos | GET|/quicksilver/epochs/v1/epochs|
| `CurrentEpoch` | [QueryCurrentEpochRequest](#quicksilver.epochs.v1.QueryCurrentEpochRequest) | [QueryCurrentEpochResponse](#quicksilver.epochs.v1.QueryCurrentEpochResponse) | CurrentEpoch provide current epoch of specified identifier | GET|/quicksilver/epochs/v1/current_epoch|

 <!-- end services -->



<a name="quicksilver/interchainquery/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainquery/v1/genesis.proto



<a name="quicksilver.interchainquery.v1.DataPoint"></a>

### DataPoint



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `remote_height` | [string](#string) |  |  |
| `local_height` | [string](#string) |  |  |
| `value` | [bytes](#bytes) |  |  |






<a name="quicksilver.interchainquery.v1.GenesisState"></a>

### GenesisState
GenesisState defines the epochs module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `queries` | [Query](#quicksilver.interchainquery.v1.Query) | repeated |  |






<a name="quicksilver.interchainquery.v1.Query"></a>

### Query



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `connection_id` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `query_type` | [string](#string) |  |  |
| `request` | [bytes](#bytes) |  |  |
| `period` | [string](#string) |  |  |
| `last_height` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="quicksilver/interchainquery/v1/messages.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainquery/v1/messages.proto



<a name="quicksilver.interchainquery.v1.MsgSubmitQueryResponse"></a>

### MsgSubmitQueryResponse
MsgSubmitQueryResponse represents a message type to fulfil a query request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |
| `query_id` | [string](#string) |  |  |
| `result` | [bytes](#bytes) |  |  |
| `proof_ops` | [tendermint.crypto.ProofOps](#tendermint.crypto.ProofOps) |  |  |
| `height` | [int64](#int64) |  |  |
| `from_address` | [string](#string) |  |  |






<a name="quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse"></a>

### MsgSubmitQueryResponseResponse
MsgSubmitQueryResponseResponse defines the MsgSubmitQueryResponse response
type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="quicksilver.interchainquery.v1.Msg"></a>

### Msg
Msg defines the interchainquery Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitQueryResponse` | [MsgSubmitQueryResponse](#quicksilver.interchainquery.v1.MsgSubmitQueryResponse) | [MsgSubmitQueryResponseResponse](#quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse) | SubmitQueryResponse defines a method for submit query responses. | POST|/interchainquery/tx/v1beta1/submitquery|

 <!-- end services -->



<a name="quicksilver/interchainstaking/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainstaking/v1/genesis.proto



<a name="quicksilver.interchainstaking.v1.Delegation"></a>

### Delegation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `height` | [int64](#int64) |  |  |
| `redelegation_end` | [int64](#int64) |  |  |






<a name="quicksilver.interchainstaking.v1.DelegationPlan"></a>

### DelegationPlan



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [DelegationPlan.ValueEntry](#quicksilver.interchainstaking.v1.DelegationPlan.ValueEntry) | repeated |  |






<a name="quicksilver.interchainstaking.v1.DelegationPlan.DelegationPlanItem"></a>

### DelegationPlan.DelegationPlanItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="quicksilver.interchainstaking.v1.DelegationPlan.ValueEntry"></a>

### DelegationPlan.ValueEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [DelegationPlan.DelegationPlanItem](#quicksilver.interchainstaking.v1.DelegationPlan.DelegationPlanItem) |  |  |






<a name="quicksilver.interchainstaking.v1.DelegatorIntent"></a>

### DelegatorIntent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator` | [string](#string) |  |  |
| `intents` | [ValidatorIntent](#quicksilver.interchainstaking.v1.ValidatorIntent) | repeated |  |






<a name="quicksilver.interchainstaking.v1.DistributionPlan"></a>

### DistributionPlan



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [DistributionPlan.ValueEntry](#quicksilver.interchainstaking.v1.DistributionPlan.ValueEntry) | repeated |  |






<a name="quicksilver.interchainstaking.v1.DistributionPlan.ValueEntry"></a>

### DistributionPlan.ValueEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [DelegationPlan](#quicksilver.interchainstaking.v1.DelegationPlan) |  |  |






<a name="quicksilver.interchainstaking.v1.GenesisState"></a>

### GenesisState
GenesisState defines the interchainstaking module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#quicksilver.interchainstaking.v1.Params) |  |  |
| `zones` | [RegisteredZone](#quicksilver.interchainstaking.v1.RegisteredZone) | repeated |  |






<a name="quicksilver.interchainstaking.v1.ICAAccount"></a>

### ICAAccount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | balance defines the different coins this balance holds. |
| `delegated_balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `port_name` | [string](#string) |  |  |
| `balance_waitgroup` | [uint32](#uint32) |  | Delegations here? or against validator? |






<a name="quicksilver.interchainstaking.v1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_account_count` | [uint64](#uint64) |  |  |
| `delegation_account_split` | [uint64](#uint64) |  |  |
| `deposit_interval` | [uint64](#uint64) |  |  |
| `delegate_interval` | [uint64](#uint64) |  |  |
| `delegations_interval` | [uint64](#uint64) |  |  |
| `validatorset_interval` | [uint64](#uint64) |  |  |
| `commission_rate` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.PortConnectionTuple"></a>

### PortConnectionTuple



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `port_id` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.Receipt"></a>

### Receipt



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `zone` | [RegisteredZone](#quicksilver.interchainstaking.v1.RegisteredZone) |  |  |
| `sender` | [string](#string) |  |  |
| `txhash` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `distribution_plan` | [DistributionPlan](#quicksilver.interchainstaking.v1.DistributionPlan) |  |  |






<a name="quicksilver.interchainstaking.v1.RegisteredZone"></a>

### RegisteredZone



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `deposit_address` | [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount) |  |  |
| `withdrawal_address` | [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount) |  |  |
| `performance_address` | [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount) |  |  |
| `delegation_addresses` | [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount) | repeated |  |
| `account_prefix` | [string](#string) |  |  |
| `local_denom` | [string](#string) |  |  |
| `base_denom` | [string](#string) |  |  |
| `redemption_rate` | [string](#string) |  |  |
| `validators` | [Validator](#quicksilver.interchainstaking.v1.Validator) | repeated |  |
| `delegator_intent` | [RegisteredZone.DelegatorIntentEntry](#quicksilver.interchainstaking.v1.RegisteredZone.DelegatorIntentEntry) | repeated |  |
| `aggregate_intent` | [RegisteredZone.AggregateIntentEntry](#quicksilver.interchainstaking.v1.RegisteredZone.AggregateIntentEntry) | repeated |  |
| `multi_send` | [bool](#bool) |  |  |
| `last_redemption_rate` | [string](#string) |  |  |
| `withdrawal_waitgroup` | [uint32](#uint32) |  |  |
| `ibc_next_validators_hash` | [bytes](#bytes) |  |  |






<a name="quicksilver.interchainstaking.v1.RegisteredZone.AggregateIntentEntry"></a>

### RegisteredZone.AggregateIntentEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [ValidatorIntent](#quicksilver.interchainstaking.v1.ValidatorIntent) |  |  |






<a name="quicksilver.interchainstaking.v1.RegisteredZone.DelegatorIntentEntry"></a>

### RegisteredZone.DelegatorIntentEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [DelegatorIntent](#quicksilver.interchainstaking.v1.DelegatorIntent) |  |  |






<a name="quicksilver.interchainstaking.v1.TransferRecord"></a>

### TransferRecord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="quicksilver.interchainstaking.v1.Validator"></a>

### Validator



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valoper_address` | [string](#string) |  |  |
| `commission_rate` | [string](#string) |  |  |
| `delegator_shares` | [string](#string) |  |  |
| `voting_power` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.ValidatorIntent"></a>

### ValidatorIntent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valoper_address` | [string](#string) |  |  |
| `weight` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.WithdrawalRecord"></a>

### WithdrawalRecord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `status` | [int32](#int32) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="quicksilver/interchainstaking/v1/messages.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainstaking/v1/messages.proto



<a name="quicksilver.interchainstaking.v1.MsgRegisterZone"></a>

### MsgRegisterZone
MsgRegisterZone represents a message type to register a new zone. TODO:
deprecate in favour of governance vote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `base_denom` | [string](#string) |  |  |
| `local_denom` | [string](#string) |  |  |
| `account_prefix` | [string](#string) |  |  |
| `from_address` | [string](#string) |  |  |
| `multi_send` | [bool](#bool) |  |  |






<a name="quicksilver.interchainstaking.v1.MsgRegisterZoneResponse"></a>

### MsgRegisterZoneResponse
MsgRegisterZoneResponse defines the MsgRegisterZone response type.






<a name="quicksilver.interchainstaking.v1.MsgRequestRedemption"></a>

### MsgRequestRedemption
MsgRegisterZone represents a message type to request a burn of qAssets for
native assets.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coin` | [string](#string) |  |  |
| `destination_address` | [string](#string) |  |  |
| `from_address` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.MsgRequestRedemptionResponse"></a>

### MsgRequestRedemptionResponse
MsgRequestRedemptionResponse defines the MsgRequestRedemption response type.






<a name="quicksilver.interchainstaking.v1.MsgSignalIntent"></a>

### MsgSignalIntent
MsgSignalIntent represents a message type for signalling voting intent for
one or more validators.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |
| `intents` | [ValidatorIntent](#quicksilver.interchainstaking.v1.ValidatorIntent) | repeated |  |
| `from_address` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.MsgSignalIntentResponse"></a>

### MsgSignalIntentResponse
MsgSignalIntentResponse defines the MsgSignalIntent response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="quicksilver.interchainstaking.v1.Msg"></a>

### Msg
Msg defines the interchainstaking Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RegisterZone` | [MsgRegisterZone](#quicksilver.interchainstaking.v1.MsgRegisterZone) | [MsgRegisterZoneResponse](#quicksilver.interchainstaking.v1.MsgRegisterZoneResponse) | RegisterZone defines a method for registering a new zone. TODO: deprecate in favour of governance vote. | POST|/quicksilver/tx/v1/interchainstaking/zone|
| `RequestRedemption` | [MsgRequestRedemption](#quicksilver.interchainstaking.v1.MsgRequestRedemption) | [MsgRequestRedemptionResponse](#quicksilver.interchainstaking.v1.MsgRequestRedemptionResponse) | RequestRedemption defines a method for requesting burning of qAssets for native assets. | POST|/quicksilver/tx/v1/interchainstaking/redeem|
| `SignalIntent` | [MsgSignalIntent](#quicksilver.interchainstaking.v1.MsgSignalIntent) | [MsgSignalIntentResponse](#quicksilver.interchainstaking.v1.MsgSignalIntentResponse) | SignalIntent defines a method for signalling voting intent for one or more validators. | POST|/quicksilver/tx/v1/interchainstaking/intent|

 <!-- end services -->



<a name="quicksilver/interchainstaking/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainstaking/v1/query.proto



<a name="quicksilver.interchainstaking.v1.QueryDelegationsRequest"></a>

### QueryDelegationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryDelegationsResponse"></a>

### QueryDelegationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegations` | [Delegation](#quicksilver.interchainstaking.v1.Delegation) | repeated |  |






<a name="quicksilver.interchainstaking.v1.QueryDelegatorDelegationsRequest"></a>

### QueryDelegatorDelegationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryDelegatorDelegationsResponse"></a>

### QueryDelegatorDelegationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegations` | [Delegation](#quicksilver.interchainstaking.v1.Delegation) | repeated |  |






<a name="quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest"></a>

### QueryDelegatorIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |
| `delegator_address` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse"></a>

### QueryDelegatorIntentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `intent` | [DelegatorIntent](#quicksilver.interchainstaking.v1.DelegatorIntent) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest"></a>

### QueryDepositAccountForChainRequest
QueryDepositAccountForChainRequest is the request type for the
Query/InterchainAccountAddress RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse"></a>

### QueryDepositAccountForChainResponse
QueryDepositAccountForChainResponse the response type for the
Query/InterchainAccountAddress RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposit_account_address` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoRequest"></a>

### QueryRegisteredZonesInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoResponse"></a>

### QueryRegisteredZonesInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `zones` | [RegisteredZone](#quicksilver.interchainstaking.v1.RegisteredZone) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryValidatorDelegationsRequest"></a>

### QueryValidatorDelegationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="quicksilver.interchainstaking.v1.QueryValidatorDelegationsResponse"></a>

### QueryValidatorDelegationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegations` | [Delegation](#quicksilver.interchainstaking.v1.Delegation) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="quicksilver.interchainstaking.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RegisteredZoneInfos` | [QueryRegisteredZonesInfoRequest](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoRequest) | [QueryRegisteredZonesInfoResponse](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoResponse) | RegisteredZoneInfos provides meta data on connected zones. | GET|/quicksilver/interchainstaking/v1/zones|
| `DepositAccount` | [QueryDepositAccountForChainRequest](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest) | [QueryDepositAccountForChainResponse](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse) | DepositAccount provides data on the deposit address for a connected zone. | GET|/quicksilver/interchainstaking/v1/zones/{chain_id}/deposit_address|
| `DelegatorIntent` | [QueryDelegatorIntentRequest](#quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest) | [QueryDelegatorIntentResponse](#quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse) | DelegatorIntent provides data on the intent of the delegator for the given zone. | GET|/quicksilver/interchainstaking/v1/zones/{chain_id}/delegator_intent/{delegator_address}|
| `Delegations` | [QueryDelegationsRequest](#quicksilver.interchainstaking.v1.QueryDelegationsRequest) | [QueryDelegationsResponse](#quicksilver.interchainstaking.v1.QueryDelegationsResponse) | Delegations provides data on the delegations for the given zone. | GET|/quicksilver/interchainstaking/v1/zones/{chain_id}/delegations|
| `DelegatorDelegations` | [QueryDelegatorDelegationsRequest](#quicksilver.interchainstaking.v1.QueryDelegatorDelegationsRequest) | [QueryDelegatorDelegationsResponse](#quicksilver.interchainstaking.v1.QueryDelegatorDelegationsResponse) | DelegatorDelegations provides data on the delegations from a given delegator for the given zone. | GET|/quicksilver/interchainstaking/v1/zones/{chain_id}/delegator_delegations/{delegator_address}|
| `ValidatorDelegations` | [QueryValidatorDelegationsRequest](#quicksilver.interchainstaking.v1.QueryValidatorDelegationsRequest) | [QueryValidatorDelegationsResponse](#quicksilver.interchainstaking.v1.QueryValidatorDelegationsResponse) | ValidatorDelegations provides data on the delegations to a given validator for the given zone. | GET|/quicksilver/interchainstaking/v1/zones/{chain_id}/validator_delegations/{validator_address}|

 <!-- end services -->



<a name="quicksilver/mint/v1beta1/mint.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/mint/v1beta1/mint.proto



<a name="quicksilver.mint.v1beta1.DistributionProportions"></a>

### DistributionProportions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `staking` | [string](#string) |  | staking defines the proportion of the minted minted_denom that is to be allocated as staking rewards. |
| `pool_incentives` | [string](#string) |  | pool_incentives defines the proportion of the minted minted_denom that is to be allocated as pool incentives. |
| `participation_rewards` | [string](#string) |  | participation_rewards defines the proportion of the minted minted_denom that is to be allocated to participation rewards address. |
| `community_pool` | [string](#string) |  | community_pool defines the proportion of the minted minted_denom that is to be allocated to the community pool. |






<a name="quicksilver.mint.v1beta1.Minter"></a>

### Minter
Minter represents the minting state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `epoch_provisions` | [string](#string) |  | current epoch provisions |






<a name="quicksilver.mint.v1beta1.Params"></a>

### Params
Params holds parameters for the mint module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mint_denom` | [string](#string) |  | type of coin to mint |
| `genesis_epoch_provisions` | [string](#string) |  | epoch provisions from the first epoch |
| `epoch_identifier` | [string](#string) |  | mint epoch identifier |
| `reduction_period_in_epochs` | [int64](#int64) |  | number of epochs take to reduce rewards |
| `reduction_factor` | [string](#string) |  | reduction multiplier to execute on each period |
| `distribution_proportions` | [DistributionProportions](#quicksilver.mint.v1beta1.DistributionProportions) |  | distribution_proportions defines the proportion of the minted denom |
| `minting_rewards_distribution_start_epoch` | [int64](#int64) |  | start epoch to distribute minting rewards |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="quicksilver/mint/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/mint/v1beta1/genesis.proto



<a name="quicksilver.mint.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the mint module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minter` | [Minter](#quicksilver.mint.v1beta1.Minter) |  | minter is a space for holding current rewards information. |
| `params` | [Params](#quicksilver.mint.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `reduction_started_epoch` | [int64](#int64) |  | current reduction period start epoch |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="quicksilver/mint/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/mint/v1beta1/query.proto



<a name="quicksilver.mint.v1beta1.QueryEpochProvisionsRequest"></a>

### QueryEpochProvisionsRequest
QueryEpochProvisionsRequest is the request type for the
Query/EpochProvisions RPC method.






<a name="quicksilver.mint.v1beta1.QueryEpochProvisionsResponse"></a>

### QueryEpochProvisionsResponse
QueryEpochProvisionsResponse is the response type for the
Query/EpochProvisions RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `epoch_provisions` | [bytes](#bytes) |  | epoch_provisions is the current minting per epoch provisions value. |






<a name="quicksilver.mint.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="quicksilver.mint.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#quicksilver.mint.v1beta1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="quicksilver.mint.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#quicksilver.mint.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#quicksilver.mint.v1beta1.QueryParamsResponse) | Params returns the total set of minting parameters. | GET|/quicksilver/mint/v1beta1/params|
| `EpochProvisions` | [QueryEpochProvisionsRequest](#quicksilver.mint.v1beta1.QueryEpochProvisionsRequest) | [QueryEpochProvisionsResponse](#quicksilver.mint.v1beta1.QueryEpochProvisionsResponse) | EpochProvisions current minting epoch provisions value. | GET|/quicksilver/mint/v1beta1/epoch_provisions|

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

