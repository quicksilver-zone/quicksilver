# Interchain Query

## Abstract

Module, `x/interchainquery`, defines and implements the mechanisms to
facilitate provable cross-chain queries.

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

[needs ellaboration]

## State

### Query

```go
type Query struct {
	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ConnectionId string `protobuf:"bytes,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	ChainId      string `protobuf:"bytes,3,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	QueryType    string `protobuf:"bytes,4,opt,name=query_type,json=queryType,proto3" json:"query_type,omitempty"`
	Request      []byte `protobuf:"bytes,5,opt,name=request,proto3" json:"request,omitempty"`
	// change these to uint64 in v0.5.0
	Period       github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=period,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"period"`
	LastHeight   github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,7,opt,name=last_height,json=lastHeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"last_height"`
	CallbackId   string                                 `protobuf:"bytes,8,opt,name=callback_id,json=callbackId,proto3" json:"callback_id,omitempty"`
	Ttl          uint64                                 `protobuf:"varint,9,opt,name=ttl,proto3" json:"ttl,omitempty"`
	LastEmission github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,10,opt,name=last_emission,json=lastEmission,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"last_emission"`
}
```

### DataPoint

```go
type DataPoint struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// change these to uint64 in v0.5.0
	RemoteHeight github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=remote_height,json=remoteHeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"remote_height"`
	LocalHeight  github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=local_height,json=localHeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"local_height"`
	Value        []byte                                 `protobuf:"bytes,4,opt,name=value,proto3" json:"result,omitempty"`
}
```

## Messages

Description of message types that trigger state transitions;

```protobuf
service Msg {
  // SubmitQueryResponse defines a method for submit query responses.
  rpc SubmitQueryResponse(MsgSubmitQueryResponse)
      returns (MsgSubmitQueryResponseResponse) {
    option (google.api.http) = {
      post : "/interchainquery/tx/v1beta1/submitquery"
      body : "*"
    };
  };
}
```

### MsgSubmitQueryResponse

MsgSubmitQueryResponse is used to signal a response from a remote chain.

```go
// MsgSubmitQueryResponse represents a message type to fulfil a query request.
type MsgSubmitQueryResponse struct {
	ChainId     string           `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
	QueryId     string           `protobuf:"bytes,2,opt,name=query_id,json=queryId,proto3" json:"query_id,omitempty" yaml:"query_id"`
	Result      []byte           `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty" yaml:"result"`
	ProofOps    *crypto.ProofOps `protobuf:"bytes,4,opt,name=proof_ops,json=proofOps,proto3" json:"proof_ops,omitempty" yaml:"proof_ops"`
	Height      int64            `protobuf:"varint,5,opt,name=height,proto3" json:"height,omitempty" yaml:"height"`
	FromAddress string           `protobuf:"bytes,6,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
}
```

* **ChainId** - the chain id of the remote chain;
* **QueryId** - the query id that solicited this response;
* **Result** - the encoded query response from the remote chain;
* **ProofOps** - the cryptographic proofs related to this response;
* **Height** - the block height of the remote chain at the time of response;
* **FromAddress** - ;

## Transactions

N/A

## Events

Events emitted by module for tracking messages and index transactions;

### EndBlocker

| Type    | Attribute Key | Attribute Value   |
|:--------|:--------------|:------------------|
| message | module        | interchainquery   |
| message | query_id      | {query_id}        |
| message | chain_id      | {chain_id}        |
| message | connection_id | {connection_id}   |
| message | type          | {query_type}      |
| message | height        | "0"               |
| message | request       | {request}         |

## Hooks

N/A

## Queries

Description of available information request queries;

```protobuf
service QuerySrvr {
  // Params returns the total set of minting parameters.
  rpc Queries(QueryRequestsRequest) returns (QueryRequestsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainquery/v1/queries/{chain_id}";
  }
}
```

### queries

Query the existing IBC queries of the module.

```go
type QueryRequestsRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	ChainId    string             `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

type QueryRequestsResponse struct {
	// params defines the parameters of the module.
	Queries    []Query             `protobuf:"bytes,1,rep,name=queries,proto3" json:"queries"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
```

## Keepers

<https://pkg.go.dev/github.com/quicksilver-zone/quicksilver/x/interchainquery/keeper>

## Parameters

N/A

## Begin Block

N/A

## End Block

* Iterate through all queries and emit events for periodic queries.
* Iterate through all data points to perform garbage collection.
