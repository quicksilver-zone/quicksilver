import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Query, QueryAmino, QuerySDKType } from "./interchainquery";
import { Tx, TxAmino, TxSDKType } from "../../../cosmos/tx/v1beta1/tx";
import { TxResponse, TxResponseAmino, TxResponseSDKType } from "../../../cosmos/base/abci/v1beta1/abci";
import { TxProof, TxProofAmino, TxProofSDKType } from "../../../tendermint/types/types";
import { Header, HeaderAmino, HeaderSDKType } from "../../../ibc/lightclients/tendermint/v1/tendermint";
import * as _m0 from "protobufjs/minimal";
import { isSet, DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.interchainquery.v1";
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryRequestsRequest {
  pagination: PageRequest;
  chainId: string;
}
export interface QueryRequestsRequestProtoMsg {
  typeUrl: "/quicksilver.interchainquery.v1.QueryRequestsRequest";
  value: Uint8Array;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryRequestsRequestAmino {
  pagination?: PageRequestAmino;
  chain_id: string;
}
export interface QueryRequestsRequestAminoMsg {
  type: "/quicksilver.interchainquery.v1.QueryRequestsRequest";
  value: QueryRequestsRequestAmino;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryRequestsRequestSDKType {
  pagination: PageRequestSDKType;
  chain_id: string;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryRequestsResponse {
  /** params defines the parameters of the module. */
  queries: Query[];
  pagination: PageResponse;
}
export interface QueryRequestsResponseProtoMsg {
  typeUrl: "/quicksilver.interchainquery.v1.QueryRequestsResponse";
  value: Uint8Array;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryRequestsResponseAmino {
  /** params defines the parameters of the module. */
  queries: QueryAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryRequestsResponseAminoMsg {
  type: "/quicksilver.interchainquery.v1.QueryRequestsResponse";
  value: QueryRequestsResponseAmino;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryRequestsResponseSDKType {
  queries: QuerySDKType[];
  pagination: PageResponseSDKType;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxWithProofResponse {
  /** tx is the queried transaction; deprecated. */
  tx: Tx;
  /**
   * deprecated, v1.2.13
   * proof is the tmproto.TxProof for the queried tx
   */
  txResponse: TxResponse;
  proof: TxProof;
  /** ibc-go header to validate txs */
  header: Header;
}
export interface GetTxWithProofResponseProtoMsg {
  typeUrl: "/quicksilver.interchainquery.v1.GetTxWithProofResponse";
  value: Uint8Array;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxWithProofResponseAmino {
  /** tx is the queried transaction; deprecated. */
  tx?: TxAmino;
  /**
   * deprecated, v1.2.13
   * proof is the tmproto.TxProof for the queried tx
   */
  tx_response?: TxResponseAmino;
  proof?: TxProofAmino;
  /** ibc-go header to validate txs */
  header?: HeaderAmino;
}
export interface GetTxWithProofResponseAminoMsg {
  type: "/quicksilver.interchainquery.v1.GetTxWithProofResponse";
  value: GetTxWithProofResponseAmino;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxWithProofResponseSDKType {
  tx: TxSDKType;
  tx_response: TxResponseSDKType;
  proof: TxProofSDKType;
  header: HeaderSDKType;
}
function createBaseQueryRequestsRequest(): QueryRequestsRequest {
  return {
    pagination: PageRequest.fromPartial({}),
    chainId: ""
  };
}
export const QueryRequestsRequest = {
  typeUrl: "/quicksilver.interchainquery.v1.QueryRequestsRequest",
  encode(message: QueryRequestsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    if (message.chainId !== "") {
      writer.uint32(18).string(message.chainId);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRequestsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRequestsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        case 2:
          message.chainId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryRequestsRequest {
    const obj = createBaseQueryRequestsRequest();
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    return obj;
  },
  toJSON(message: QueryRequestsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    message.chainId !== undefined && (obj.chainId = message.chainId);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryRequestsRequest>): QueryRequestsRequest {
    const message = createBaseQueryRequestsRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    message.chainId = object.chainId ?? "";
    return message;
  },
  fromSDK(object: QueryRequestsRequestSDKType): QueryRequestsRequest {
    return {
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined,
      chainId: object?.chain_id
    };
  },
  toSDK(message: QueryRequestsRequest): QueryRequestsRequestSDKType {
    const obj: any = {};
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAmino(object: QueryRequestsRequestAmino): QueryRequestsRequest {
    return {
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined,
      chainId: object.chain_id
    };
  },
  toAmino(message: QueryRequestsRequest): QueryRequestsRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAminoMsg(object: QueryRequestsRequestAminoMsg): QueryRequestsRequest {
    return QueryRequestsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryRequestsRequestProtoMsg): QueryRequestsRequest {
    return QueryRequestsRequest.decode(message.value);
  },
  toProto(message: QueryRequestsRequest): Uint8Array {
    return QueryRequestsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryRequestsRequest): QueryRequestsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainquery.v1.QueryRequestsRequest",
      value: QueryRequestsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryRequestsResponse(): QueryRequestsResponse {
  return {
    queries: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryRequestsResponse = {
  typeUrl: "/quicksilver.interchainquery.v1.QueryRequestsResponse",
  encode(message: QueryRequestsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.queries) {
      Query.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRequestsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRequestsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.queries.push(Query.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryRequestsResponse {
    const obj = createBaseQueryRequestsResponse();
    if (Array.isArray(object?.queries)) obj.queries = object.queries.map((e: any) => Query.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryRequestsResponse): unknown {
    const obj: any = {};
    if (message.queries) {
      obj.queries = message.queries.map(e => e ? Query.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryRequestsResponse>): QueryRequestsResponse {
    const message = createBaseQueryRequestsResponse();
    message.queries = object.queries?.map(e => Query.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryRequestsResponseSDKType): QueryRequestsResponse {
    return {
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => Query.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryRequestsResponse): QueryRequestsResponseSDKType {
    const obj: any = {};
    if (message.queries) {
      obj.queries = message.queries.map(e => e ? Query.toSDK(e) : undefined);
    } else {
      obj.queries = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryRequestsResponseAmino): QueryRequestsResponse {
    return {
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => Query.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryRequestsResponse): QueryRequestsResponseAmino {
    const obj: any = {};
    if (message.queries) {
      obj.queries = message.queries.map(e => e ? Query.toAmino(e) : undefined);
    } else {
      obj.queries = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryRequestsResponseAminoMsg): QueryRequestsResponse {
    return QueryRequestsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryRequestsResponseProtoMsg): QueryRequestsResponse {
    return QueryRequestsResponse.decode(message.value);
  },
  toProto(message: QueryRequestsResponse): Uint8Array {
    return QueryRequestsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryRequestsResponse): QueryRequestsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainquery.v1.QueryRequestsResponse",
      value: QueryRequestsResponse.encode(message).finish()
    };
  }
};
function createBaseGetTxWithProofResponse(): GetTxWithProofResponse {
  return {
    tx: Tx.fromPartial({}),
    txResponse: TxResponse.fromPartial({}),
    proof: TxProof.fromPartial({}),
    header: Header.fromPartial({})
  };
}
export const GetTxWithProofResponse = {
  typeUrl: "/quicksilver.interchainquery.v1.GetTxWithProofResponse",
  encode(message: GetTxWithProofResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tx !== undefined) {
      Tx.encode(message.tx, writer.uint32(10).fork()).ldelim();
    }
    if (message.txResponse !== undefined) {
      TxResponse.encode(message.txResponse, writer.uint32(18).fork()).ldelim();
    }
    if (message.proof !== undefined) {
      TxProof.encode(message.proof, writer.uint32(26).fork()).ldelim();
    }
    if (message.header !== undefined) {
      Header.encode(message.header, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): GetTxWithProofResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetTxWithProofResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tx = Tx.decode(reader, reader.uint32());
          break;
        case 2:
          message.txResponse = TxResponse.decode(reader, reader.uint32());
          break;
        case 3:
          message.proof = TxProof.decode(reader, reader.uint32());
          break;
        case 4:
          message.header = Header.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): GetTxWithProofResponse {
    const obj = createBaseGetTxWithProofResponse();
    if (isSet(object.tx)) obj.tx = Tx.fromJSON(object.tx);
    if (isSet(object.txResponse)) obj.txResponse = TxResponse.fromJSON(object.txResponse);
    if (isSet(object.proof)) obj.proof = TxProof.fromJSON(object.proof);
    if (isSet(object.header)) obj.header = Header.fromJSON(object.header);
    return obj;
  },
  toJSON(message: GetTxWithProofResponse): unknown {
    const obj: any = {};
    message.tx !== undefined && (obj.tx = message.tx ? Tx.toJSON(message.tx) : undefined);
    message.txResponse !== undefined && (obj.txResponse = message.txResponse ? TxResponse.toJSON(message.txResponse) : undefined);
    message.proof !== undefined && (obj.proof = message.proof ? TxProof.toJSON(message.proof) : undefined);
    message.header !== undefined && (obj.header = message.header ? Header.toJSON(message.header) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<GetTxWithProofResponse>): GetTxWithProofResponse {
    const message = createBaseGetTxWithProofResponse();
    if (object.tx !== undefined && object.tx !== null) {
      message.tx = Tx.fromPartial(object.tx);
    }
    if (object.txResponse !== undefined && object.txResponse !== null) {
      message.txResponse = TxResponse.fromPartial(object.txResponse);
    }
    if (object.proof !== undefined && object.proof !== null) {
      message.proof = TxProof.fromPartial(object.proof);
    }
    if (object.header !== undefined && object.header !== null) {
      message.header = Header.fromPartial(object.header);
    }
    return message;
  },
  fromSDK(object: GetTxWithProofResponseSDKType): GetTxWithProofResponse {
    return {
      tx: object.tx ? Tx.fromSDK(object.tx) : undefined,
      txResponse: object.tx_response ? TxResponse.fromSDK(object.tx_response) : undefined,
      proof: object.proof ? TxProof.fromSDK(object.proof) : undefined,
      header: object.header ? Header.fromSDK(object.header) : undefined
    };
  },
  toSDK(message: GetTxWithProofResponse): GetTxWithProofResponseSDKType {
    const obj: any = {};
    message.tx !== undefined && (obj.tx = message.tx ? Tx.toSDK(message.tx) : undefined);
    message.txResponse !== undefined && (obj.tx_response = message.txResponse ? TxResponse.toSDK(message.txResponse) : undefined);
    message.proof !== undefined && (obj.proof = message.proof ? TxProof.toSDK(message.proof) : undefined);
    message.header !== undefined && (obj.header = message.header ? Header.toSDK(message.header) : undefined);
    return obj;
  },
  fromAmino(object: GetTxWithProofResponseAmino): GetTxWithProofResponse {
    return {
      tx: object?.tx ? Tx.fromAmino(object.tx) : undefined,
      txResponse: object?.tx_response ? TxResponse.fromAmino(object.tx_response) : undefined,
      proof: object?.proof ? TxProof.fromAmino(object.proof) : undefined,
      header: object?.header ? Header.fromAmino(object.header) : undefined
    };
  },
  toAmino(message: GetTxWithProofResponse): GetTxWithProofResponseAmino {
    const obj: any = {};
    obj.tx = message.tx ? Tx.toAmino(message.tx) : undefined;
    obj.tx_response = message.txResponse ? TxResponse.toAmino(message.txResponse) : undefined;
    obj.proof = message.proof ? TxProof.toAmino(message.proof) : undefined;
    obj.header = message.header ? Header.toAmino(message.header) : undefined;
    return obj;
  },
  fromAminoMsg(object: GetTxWithProofResponseAminoMsg): GetTxWithProofResponse {
    return GetTxWithProofResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: GetTxWithProofResponseProtoMsg): GetTxWithProofResponse {
    return GetTxWithProofResponse.decode(message.value);
  },
  toProto(message: GetTxWithProofResponse): Uint8Array {
    return GetTxWithProofResponse.encode(message).finish();
  },
  toProtoMsg(message: GetTxWithProofResponse): GetTxWithProofResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainquery.v1.GetTxWithProofResponse",
      value: GetTxWithProofResponse.encode(message).finish()
    };
  }
};