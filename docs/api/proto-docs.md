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
  
- [quicksilver/interchainstaking/v1/genesis.proto](#quicksilver/interchainstaking/v1/genesis.proto)
    - [GenesisState](#quicksilver.interchainstaking.v1.GenesisState)
    - [RegisteredZone](#quicksilver.interchainstaking.v1.RegisteredZone)
  
- [quicksilver/interchainstaking/v1/messages.proto](#quicksilver/interchainstaking/v1/messages.proto)
    - [MsgRegisterZone](#quicksilver.interchainstaking.v1.MsgRegisterZone)
    - [MsgRegisterZoneResponse](#quicksilver.interchainstaking.v1.MsgRegisterZoneResponse)
  
    - [Msg](#quicksilver.interchainstaking.v1.Msg)
  
- [quicksilver/interchainstaking/v1/query.proto](#quicksilver/interchainstaking/v1/query.proto)
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



<a name="quicksilver/interchainstaking/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainstaking/v1/genesis.proto



<a name="quicksilver.interchainstaking.v1.GenesisState"></a>

### GenesisState
GenesisState defines the epochs module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `zones` | [RegisteredZone](#quicksilver.interchainstaking.v1.RegisteredZone) | repeated |  |






<a name="quicksilver.interchainstaking.v1.RegisteredZone"></a>

### RegisteredZone



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `deposit_address` | [string](#string) |  |  |
| `delegation_addresses` | [string](#string) | repeated |  |
| `local_denom` | [string](#string) |  |  |
| `remote_denom` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="quicksilver/interchainstaking/v1/messages.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainstaking/v1/messages.proto



<a name="quicksilver.interchainstaking.v1.MsgRegisterZone"></a>

### MsgRegisterZone
MsgRegisterZone represents a message to send coins from one account to
another.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `local_denom` | [string](#string) |  |  |
| `remote_denom` | [string](#string) |  |  |
| `from_address` | [string](#string) |  |  |






<a name="quicksilver.interchainstaking.v1.MsgRegisterZoneResponse"></a>

### MsgRegisterZoneResponse
MsgRegisterZoneResponse defines the Msg/Send response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="quicksilver.interchainstaking.v1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RegisterZone` | [MsgRegisterZone](#quicksilver.interchainstaking.v1.MsgRegisterZone) | [MsgRegisterZoneResponse](#quicksilver.interchainstaking.v1.MsgRegisterZoneResponse) | RegisterZone defines a method for sending coins from one account to another account. | |

 <!-- end services -->



<a name="quicksilver/interchainstaking/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## quicksilver/interchainstaking/v1/query.proto



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
| `RegisteredZoneInfos` | [QueryRegisteredZonesInfoRequest](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoRequest) | [QueryRegisteredZonesInfoResponse](#quicksilver.interchainstaking.v1.QueryRegisteredZonesInfoResponse) | RegisteredZoneInfos provide running epochInfos | GET|/quicksilver/interchainstaking/v1/zones|

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

