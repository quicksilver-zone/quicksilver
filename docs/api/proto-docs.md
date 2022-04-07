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
    - [PeriodicQuery](#quicksilver.interchainquery.v1.PeriodicQuery)
    - [PeriodicQuery.QueryParametersEntry](#quicksilver.interchainquery.v1.PeriodicQuery.QueryParametersEntry)
    - [SingleQuery](#quicksilver.interchainquery.v1.SingleQuery)
    - [SingleQuery.QueryParametersEntry](#quicksilver.interchainquery.v1.SingleQuery.QueryParametersEntry)
  
- [quicksilver/interchainquery/v1/messages.proto](#quicksilver/interchainquery/v1/messages.proto)
    - [MsgSubmitQueryResponse](#quicksilver.interchainquery.v1.MsgSubmitQueryResponse)
    - [MsgSubmitQueryResponseResponse](#quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse)
  
    - [Msg](#quicksilver.interchainquery.v1.Msg)
  
- [quicksilver/interchainstaking/v1/genesis.proto](#quicksilver/interchainstaking/v1/genesis.proto)
    - [Delegation](#quicksilver.interchainstaking.v1.Delegation)
    - [DelegatorIntent](#quicksilver.interchainstaking.v1.DelegatorIntent)
    - [GenesisState](#quicksilver.interchainstaking.v1.GenesisState)
    - [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount)
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
    - [QueryDelegatorIntentRequest](#quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest)
    - [QueryDelegatorIntentResponse](#quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse)
    - [QueryDepositAccountForChainRequest](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest)
    - [QueryDepositAccountForChainResponse](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse)
    - [QueryRegisteredZonesInfoRequest](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoRequest)
    - [QueryRegisteredZonesInfoResponse](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoResponse)
  
    - [Query](#quicksilver.interchainstaking.v1.Query)
  
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
| `queries` | [PeriodicQuery](#quicksilver.interchainquery.v1.PeriodicQuery) | repeated |  |






<a name="quicksilver.interchainquery.v1.PeriodicQuery"></a>

### PeriodicQuery



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `connection_id` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `query_type` | [string](#string) |  |  |
| `query_parameters` | [PeriodicQuery.QueryParametersEntry](#quicksilver.interchainquery.v1.PeriodicQuery.QueryParametersEntry) | repeated |  |
| `period` | [string](#string) |  |  |
| `last_height` | [string](#string) |  |  |






<a name="quicksilver.interchainquery.v1.PeriodicQuery.QueryParametersEntry"></a>

### PeriodicQuery.QueryParametersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="quicksilver.interchainquery.v1.SingleQuery"></a>

### SingleQuery



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `connection_id` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `query_type` | [string](#string) |  |  |
| `query_parameters` | [SingleQuery.QueryParametersEntry](#quicksilver.interchainquery.v1.SingleQuery.QueryParametersEntry) | repeated |  |
| `emit_height` | [string](#string) |  |  |






<a name="quicksilver.interchainquery.v1.SingleQuery.QueryParametersEntry"></a>

### SingleQuery.QueryParametersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |





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
| `amount` | [string](#string) |  | TODO: determine whether this is Dec (shares) or Coins (tokens) |
| `rewards` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `redelegation_end` | [int64](#int64) |  | Delegations here? or against validator? |






<a name="quicksilver.interchainstaking.v1.DelegatorIntent"></a>

### DelegatorIntent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator` | [string](#string) |  |  |
| `intents` | [ValidatorIntent](#quicksilver.interchainstaking.v1.ValidatorIntent) | repeated |  |






<a name="quicksilver.interchainstaking.v1.GenesisState"></a>

### GenesisState
GenesisState defines the interchainstaking module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `zones` | [RegisteredZone](#quicksilver.interchainstaking.v1.RegisteredZone) | repeated |  |






<a name="quicksilver.interchainstaking.v1.ICAAccount"></a>

### ICAAccount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | balance defines the different coins this balance holds. |
| `delegated_balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `port_name` | [string](#string) |  | Delegations here? or against validator? |






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






<a name="quicksilver.interchainstaking.v1.RegisteredZone"></a>

### RegisteredZone



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `connection_id` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `deposit_address` | [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount) |  |  |
| `delegation_addresses` | [ICAAccount](#quicksilver.interchainstaking.v1.ICAAccount) | repeated |  |
| `local_denom` | [string](#string) |  |  |
| `base_denom` | [string](#string) |  |  |
| `redemption_rate` | [string](#string) |  |  |
| `validators` | [Validator](#quicksilver.interchainstaking.v1.Validator) | repeated |  |
| `delegator_intent` | [RegisteredZone.DelegatorIntentEntry](#quicksilver.interchainstaking.v1.RegisteredZone.DelegatorIntentEntry) | repeated |  |
| `aggregate_intent` | [RegisteredZone.AggregateIntentEntry](#quicksilver.interchainstaking.v1.RegisteredZone.AggregateIntentEntry) | repeated |  |
| `multi_send` | [bool](#bool) |  |  |






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
| `voting_power` | [string](#string) |  |  |
| `delegations` | [Delegation](#quicksilver.interchainstaking.v1.Delegation) | repeated |  |






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
| `identifier` | [string](#string) |  |  |
| `connection_id` | [string](#string) |  |  |
| `base_denom` | [string](#string) |  |  |
| `local_denom` | [string](#string) |  |  |
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



<a name="quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest"></a>

### QueryDelegatorIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |
| `from_address` | [string](#string) |  |  |






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





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="quicksilver.interchainstaking.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RegisteredZoneInfos` | [QueryRegisteredZonesInfoRequest](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoRequest) | [QueryRegisteredZonesInfoResponse](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoResponse) | RegisteredZoneInfos provides meta data on connected zones. | GET|/quicksilver/interchainstaking/v1/zones|
| `DepositAccountFromAddress` | [QueryDepositAccountForChainRequest](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest) | [QueryDepositAccountForChainResponse](#quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse) | DepositAccountFromAddress provides data on the deposit address for a connected zone. | GET|/quicksilver/interchainstaking/v1/zones/deposit_address|
| `DelegatorIntent` | [QueryDelegatorIntentRequest](#quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest) | [QueryDelegatorIntentResponse](#quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse) | DelegatorIntent provides data on the intent of the delegator for the given zone. | GET|/quicksilver/interchainstaking/v1/zones/delegator_intent|

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

