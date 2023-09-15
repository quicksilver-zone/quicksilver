import { Status, StatusSDKType, ZoneDrop, ZoneDropAmino, ZoneDropSDKType, ClaimRecord, ClaimRecordAmino, ClaimRecordSDKType, statusFromJSON, statusToJSON } from "./airdrop";
import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { Coin, CoinAmino, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, isSet } from "../../../helpers";
export const protobufPackage = "quicksilver.airdrop.v1";
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequest {}
export interface QueryParamsRequestProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryParamsRequest";
  value: Uint8Array;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequestAmino {}
export interface QueryParamsRequestAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryParamsRequest";
  value: QueryParamsRequestAmino;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params defines the parameters of the module. */
  params: Params;
}
export interface QueryParamsResponseProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryParamsResponse";
  value: Uint8Array;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponseAmino {
  /** params defines the parameters of the module. */
  params?: ParamsAmino;
}
export interface QueryParamsResponseAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryParamsResponse";
  value: QueryParamsResponseAmino;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponseSDKType {
  params: ParamsSDKType;
}
/** QueryZoneDropRequest is the request type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropRequest {
  /** chain_id identifies the zone. */
  chainId: string;
}
export interface QueryZoneDropRequestProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropRequest";
  value: Uint8Array;
}
/** QueryZoneDropRequest is the request type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropRequestAmino {
  /** chain_id identifies the zone. */
  chain_id: string;
}
export interface QueryZoneDropRequestAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryZoneDropRequest";
  value: QueryZoneDropRequestAmino;
}
/** QueryZoneDropRequest is the request type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropRequestSDKType {
  chain_id: string;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropResponse {
  zoneDrop: ZoneDrop;
}
export interface QueryZoneDropResponseProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropResponse";
  value: Uint8Array;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropResponseAmino {
  zone_drop?: ZoneDropAmino;
}
export interface QueryZoneDropResponseAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryZoneDropResponse";
  value: QueryZoneDropResponseAmino;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropResponseSDKType {
  zone_drop: ZoneDropSDKType;
}
/**
 * QueryAccountBalanceRequest is the request type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceRequest {
  /** chain_id identifies the zone. */
  chainId: string;
}
export interface QueryAccountBalanceRequestProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryAccountBalanceRequest";
  value: Uint8Array;
}
/**
 * QueryAccountBalanceRequest is the request type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceRequestAmino {
  /** chain_id identifies the zone. */
  chain_id: string;
}
export interface QueryAccountBalanceRequestAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryAccountBalanceRequest";
  value: QueryAccountBalanceRequestAmino;
}
/**
 * QueryAccountBalanceRequest is the request type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceRequestSDKType {
  chain_id: string;
}
/**
 * QueryAccountBalanceResponse is the response type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceResponse {
  accountBalance: Coin;
}
export interface QueryAccountBalanceResponseProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryAccountBalanceResponse";
  value: Uint8Array;
}
/**
 * QueryAccountBalanceResponse is the response type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceResponseAmino {
  account_balance?: CoinAmino;
}
export interface QueryAccountBalanceResponseAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryAccountBalanceResponse";
  value: QueryAccountBalanceResponseAmino;
}
/**
 * QueryAccountBalanceResponse is the response type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceResponseSDKType {
  account_balance: CoinSDKType;
}
/** QueryZoneDropsRequest is the request type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsRequest {
  /**
   * status enables to query zone airdrops matching a given status:
   *  - Active
   *  - Future
   *  - Expired
   */
  status: Status;
  pagination: PageRequest;
}
export interface QueryZoneDropsRequestProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropsRequest";
  value: Uint8Array;
}
/** QueryZoneDropsRequest is the request type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsRequestAmino {
  /**
   * status enables to query zone airdrops matching a given status:
   *  - Active
   *  - Future
   *  - Expired
   */
  status: Status;
  pagination?: PageRequestAmino;
}
export interface QueryZoneDropsRequestAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryZoneDropsRequest";
  value: QueryZoneDropsRequestAmino;
}
/** QueryZoneDropsRequest is the request type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsRequestSDKType {
  status: Status;
  pagination: PageRequestSDKType;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsResponse {
  zoneDrops: ZoneDrop[];
  pagination: PageResponse;
}
export interface QueryZoneDropsResponseProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropsResponse";
  value: Uint8Array;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsResponseAmino {
  zone_drops: ZoneDropAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryZoneDropsResponseAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryZoneDropsResponse";
  value: QueryZoneDropsResponseAmino;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsResponseSDKType {
  zone_drops: ZoneDropSDKType[];
  pagination: PageResponseSDKType;
}
/** QueryClaimRecordRequest is the request type for Query/ClaimRecord RPC method. */
export interface QueryClaimRecordRequest {
  chainId: string;
  address: string;
}
export interface QueryClaimRecordRequestProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordRequest";
  value: Uint8Array;
}
/** QueryClaimRecordRequest is the request type for Query/ClaimRecord RPC method. */
export interface QueryClaimRecordRequestAmino {
  chain_id: string;
  address: string;
}
export interface QueryClaimRecordRequestAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryClaimRecordRequest";
  value: QueryClaimRecordRequestAmino;
}
/** QueryClaimRecordRequest is the request type for Query/ClaimRecord RPC method. */
export interface QueryClaimRecordRequestSDKType {
  chain_id: string;
  address: string;
}
/**
 * QueryClaimRecordResponse is the response type for Query/ClaimRecord RPC
 * method.
 */
export interface QueryClaimRecordResponse {
  claimRecord: ClaimRecord;
}
export interface QueryClaimRecordResponseProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordResponse";
  value: Uint8Array;
}
/**
 * QueryClaimRecordResponse is the response type for Query/ClaimRecord RPC
 * method.
 */
export interface QueryClaimRecordResponseAmino {
  claim_record?: ClaimRecordAmino;
}
export interface QueryClaimRecordResponseAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryClaimRecordResponse";
  value: QueryClaimRecordResponseAmino;
}
/**
 * QueryClaimRecordResponse is the response type for Query/ClaimRecord RPC
 * method.
 */
export interface QueryClaimRecordResponseSDKType {
  claim_record: ClaimRecordSDKType;
}
/**
 * QueryClaimRecordsRequest is the request type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsRequest {
  chainId: string;
  pagination: PageRequest;
}
export interface QueryClaimRecordsRequestProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordsRequest";
  value: Uint8Array;
}
/**
 * QueryClaimRecordsRequest is the request type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsRequestAmino {
  chain_id: string;
  pagination?: PageRequestAmino;
}
export interface QueryClaimRecordsRequestAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryClaimRecordsRequest";
  value: QueryClaimRecordsRequestAmino;
}
/**
 * QueryClaimRecordsRequest is the request type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsRequestSDKType {
  chain_id: string;
  pagination: PageRequestSDKType;
}
/**
 * QueryClaimRecordsResponse is the response type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsResponse {
  claimRecords: ClaimRecord[];
  pagination: PageResponse;
}
export interface QueryClaimRecordsResponseProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordsResponse";
  value: Uint8Array;
}
/**
 * QueryClaimRecordsResponse is the response type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsResponseAmino {
  claim_records: ClaimRecordAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryClaimRecordsResponseAminoMsg {
  type: "/quicksilver.airdrop.v1.QueryClaimRecordsResponse";
  value: QueryClaimRecordsResponseAmino;
}
/**
 * QueryClaimRecordsResponse is the response type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsResponseSDKType {
  claim_records: ClaimRecordSDKType[];
  pagination: PageResponseSDKType;
}
function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}
export const QueryParamsRequest = {
  typeUrl: "/quicksilver.airdrop.v1.QueryParamsRequest",
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(_: any): QueryParamsRequest {
    const obj = createBaseQueryParamsRequest();
    return obj;
  },
  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
  fromSDK(_: QueryParamsRequestSDKType): QueryParamsRequest {
    return {};
  },
  toSDK(_: QueryParamsRequest): QueryParamsRequestSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: QueryParamsRequestAmino): QueryParamsRequest {
    return {};
  },
  toAmino(_: QueryParamsRequest): QueryParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryParamsRequestAminoMsg): QueryParamsRequest {
    return QueryParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryParamsRequestProtoMsg): QueryParamsRequest {
    return QueryParamsRequest.decode(message.value);
  },
  toProto(message: QueryParamsRequest): Uint8Array {
    return QueryParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryParamsRequest): QueryParamsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryParamsRequest",
      value: QueryParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    params: Params.fromPartial({})
  };
}
export const QueryParamsResponse = {
  typeUrl: "/quicksilver.airdrop.v1.QueryParamsResponse",
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryParamsResponse {
    const obj = createBaseQueryParamsResponse();
    if (isSet(object.params)) obj.params = Params.fromJSON(object.params);
    return obj;
  },
  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    }
    return message;
  },
  fromSDK(object: QueryParamsResponseSDKType): QueryParamsResponse {
    return {
      params: object.params ? Params.fromSDK(object.params) : undefined
    };
  },
  toSDK(message: QueryParamsResponse): QueryParamsResponseSDKType {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toSDK(message.params) : undefined);
    return obj;
  },
  fromAmino(object: QueryParamsResponseAmino): QueryParamsResponse {
    return {
      params: object?.params ? Params.fromAmino(object.params) : undefined
    };
  },
  toAmino(message: QueryParamsResponse): QueryParamsResponseAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryParamsResponseAminoMsg): QueryParamsResponse {
    return QueryParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryParamsResponseProtoMsg): QueryParamsResponse {
    return QueryParamsResponse.decode(message.value);
  },
  toProto(message: QueryParamsResponse): Uint8Array {
    return QueryParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryParamsResponse): QueryParamsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryParamsResponse",
      value: QueryParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryZoneDropRequest(): QueryZoneDropRequest {
  return {
    chainId: ""
  };
}
export const QueryZoneDropRequest = {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropRequest",
  encode(message: QueryZoneDropRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZoneDropRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryZoneDropRequest {
    const obj = createBaseQueryZoneDropRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    return obj;
  },
  toJSON(message: QueryZoneDropRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryZoneDropRequest>): QueryZoneDropRequest {
    const message = createBaseQueryZoneDropRequest();
    message.chainId = object.chainId ?? "";
    return message;
  },
  fromSDK(object: QueryZoneDropRequestSDKType): QueryZoneDropRequest {
    return {
      chainId: object?.chain_id
    };
  },
  toSDK(message: QueryZoneDropRequest): QueryZoneDropRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAmino(object: QueryZoneDropRequestAmino): QueryZoneDropRequest {
    return {
      chainId: object.chain_id
    };
  },
  toAmino(message: QueryZoneDropRequest): QueryZoneDropRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAminoMsg(object: QueryZoneDropRequestAminoMsg): QueryZoneDropRequest {
    return QueryZoneDropRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZoneDropRequestProtoMsg): QueryZoneDropRequest {
    return QueryZoneDropRequest.decode(message.value);
  },
  toProto(message: QueryZoneDropRequest): Uint8Array {
    return QueryZoneDropRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryZoneDropRequest): QueryZoneDropRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropRequest",
      value: QueryZoneDropRequest.encode(message).finish()
    };
  }
};
function createBaseQueryZoneDropResponse(): QueryZoneDropResponse {
  return {
    zoneDrop: ZoneDrop.fromPartial({})
  };
}
export const QueryZoneDropResponse = {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropResponse",
  encode(message: QueryZoneDropResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.zoneDrop !== undefined) {
      ZoneDrop.encode(message.zoneDrop, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZoneDropResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.zoneDrop = ZoneDrop.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryZoneDropResponse {
    const obj = createBaseQueryZoneDropResponse();
    if (isSet(object.zoneDrop)) obj.zoneDrop = ZoneDrop.fromJSON(object.zoneDrop);
    return obj;
  },
  toJSON(message: QueryZoneDropResponse): unknown {
    const obj: any = {};
    message.zoneDrop !== undefined && (obj.zoneDrop = message.zoneDrop ? ZoneDrop.toJSON(message.zoneDrop) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryZoneDropResponse>): QueryZoneDropResponse {
    const message = createBaseQueryZoneDropResponse();
    if (object.zoneDrop !== undefined && object.zoneDrop !== null) {
      message.zoneDrop = ZoneDrop.fromPartial(object.zoneDrop);
    }
    return message;
  },
  fromSDK(object: QueryZoneDropResponseSDKType): QueryZoneDropResponse {
    return {
      zoneDrop: object.zone_drop ? ZoneDrop.fromSDK(object.zone_drop) : undefined
    };
  },
  toSDK(message: QueryZoneDropResponse): QueryZoneDropResponseSDKType {
    const obj: any = {};
    message.zoneDrop !== undefined && (obj.zone_drop = message.zoneDrop ? ZoneDrop.toSDK(message.zoneDrop) : undefined);
    return obj;
  },
  fromAmino(object: QueryZoneDropResponseAmino): QueryZoneDropResponse {
    return {
      zoneDrop: object?.zone_drop ? ZoneDrop.fromAmino(object.zone_drop) : undefined
    };
  },
  toAmino(message: QueryZoneDropResponse): QueryZoneDropResponseAmino {
    const obj: any = {};
    obj.zone_drop = message.zoneDrop ? ZoneDrop.toAmino(message.zoneDrop) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryZoneDropResponseAminoMsg): QueryZoneDropResponse {
    return QueryZoneDropResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZoneDropResponseProtoMsg): QueryZoneDropResponse {
    return QueryZoneDropResponse.decode(message.value);
  },
  toProto(message: QueryZoneDropResponse): Uint8Array {
    return QueryZoneDropResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryZoneDropResponse): QueryZoneDropResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropResponse",
      value: QueryZoneDropResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAccountBalanceRequest(): QueryAccountBalanceRequest {
  return {
    chainId: ""
  };
}
export const QueryAccountBalanceRequest = {
  typeUrl: "/quicksilver.airdrop.v1.QueryAccountBalanceRequest",
  encode(message: QueryAccountBalanceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAccountBalanceRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAccountBalanceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryAccountBalanceRequest {
    const obj = createBaseQueryAccountBalanceRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    return obj;
  },
  toJSON(message: QueryAccountBalanceRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryAccountBalanceRequest>): QueryAccountBalanceRequest {
    const message = createBaseQueryAccountBalanceRequest();
    message.chainId = object.chainId ?? "";
    return message;
  },
  fromSDK(object: QueryAccountBalanceRequestSDKType): QueryAccountBalanceRequest {
    return {
      chainId: object?.chain_id
    };
  },
  toSDK(message: QueryAccountBalanceRequest): QueryAccountBalanceRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAmino(object: QueryAccountBalanceRequestAmino): QueryAccountBalanceRequest {
    return {
      chainId: object.chain_id
    };
  },
  toAmino(message: QueryAccountBalanceRequest): QueryAccountBalanceRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAminoMsg(object: QueryAccountBalanceRequestAminoMsg): QueryAccountBalanceRequest {
    return QueryAccountBalanceRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAccountBalanceRequestProtoMsg): QueryAccountBalanceRequest {
    return QueryAccountBalanceRequest.decode(message.value);
  },
  toProto(message: QueryAccountBalanceRequest): Uint8Array {
    return QueryAccountBalanceRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAccountBalanceRequest): QueryAccountBalanceRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryAccountBalanceRequest",
      value: QueryAccountBalanceRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAccountBalanceResponse(): QueryAccountBalanceResponse {
  return {
    accountBalance: Coin.fromPartial({})
  };
}
export const QueryAccountBalanceResponse = {
  typeUrl: "/quicksilver.airdrop.v1.QueryAccountBalanceResponse",
  encode(message: QueryAccountBalanceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accountBalance !== undefined) {
      Coin.encode(message.accountBalance, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAccountBalanceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAccountBalanceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accountBalance = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryAccountBalanceResponse {
    const obj = createBaseQueryAccountBalanceResponse();
    if (isSet(object.accountBalance)) obj.accountBalance = Coin.fromJSON(object.accountBalance);
    return obj;
  },
  toJSON(message: QueryAccountBalanceResponse): unknown {
    const obj: any = {};
    message.accountBalance !== undefined && (obj.accountBalance = message.accountBalance ? Coin.toJSON(message.accountBalance) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryAccountBalanceResponse>): QueryAccountBalanceResponse {
    const message = createBaseQueryAccountBalanceResponse();
    if (object.accountBalance !== undefined && object.accountBalance !== null) {
      message.accountBalance = Coin.fromPartial(object.accountBalance);
    }
    return message;
  },
  fromSDK(object: QueryAccountBalanceResponseSDKType): QueryAccountBalanceResponse {
    return {
      accountBalance: object.account_balance ? Coin.fromSDK(object.account_balance) : undefined
    };
  },
  toSDK(message: QueryAccountBalanceResponse): QueryAccountBalanceResponseSDKType {
    const obj: any = {};
    message.accountBalance !== undefined && (obj.account_balance = message.accountBalance ? Coin.toSDK(message.accountBalance) : undefined);
    return obj;
  },
  fromAmino(object: QueryAccountBalanceResponseAmino): QueryAccountBalanceResponse {
    return {
      accountBalance: object?.account_balance ? Coin.fromAmino(object.account_balance) : undefined
    };
  },
  toAmino(message: QueryAccountBalanceResponse): QueryAccountBalanceResponseAmino {
    const obj: any = {};
    obj.account_balance = message.accountBalance ? Coin.toAmino(message.accountBalance) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAccountBalanceResponseAminoMsg): QueryAccountBalanceResponse {
    return QueryAccountBalanceResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAccountBalanceResponseProtoMsg): QueryAccountBalanceResponse {
    return QueryAccountBalanceResponse.decode(message.value);
  },
  toProto(message: QueryAccountBalanceResponse): Uint8Array {
    return QueryAccountBalanceResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAccountBalanceResponse): QueryAccountBalanceResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryAccountBalanceResponse",
      value: QueryAccountBalanceResponse.encode(message).finish()
    };
  }
};
function createBaseQueryZoneDropsRequest(): QueryZoneDropsRequest {
  return {
    status: 0,
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryZoneDropsRequest = {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropsRequest",
  encode(message: QueryZoneDropsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZoneDropsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = (reader.int32() as any);
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryZoneDropsRequest {
    const obj = createBaseQueryZoneDropsRequest();
    if (isSet(object.status)) obj.status = statusFromJSON(object.status);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryZoneDropsRequest): unknown {
    const obj: any = {};
    message.status !== undefined && (obj.status = statusToJSON(message.status));
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryZoneDropsRequest>): QueryZoneDropsRequest {
    const message = createBaseQueryZoneDropsRequest();
    message.status = object.status ?? 0;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryZoneDropsRequestSDKType): QueryZoneDropsRequest {
    return {
      status: isSet(object.status) ? statusFromJSON(object.status) : -1,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryZoneDropsRequest): QueryZoneDropsRequestSDKType {
    const obj: any = {};
    message.status !== undefined && (obj.status = statusToJSON(message.status));
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryZoneDropsRequestAmino): QueryZoneDropsRequest {
    return {
      status: isSet(object.status) ? statusFromJSON(object.status) : -1,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryZoneDropsRequest): QueryZoneDropsRequestAmino {
    const obj: any = {};
    obj.status = message.status;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryZoneDropsRequestAminoMsg): QueryZoneDropsRequest {
    return QueryZoneDropsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZoneDropsRequestProtoMsg): QueryZoneDropsRequest {
    return QueryZoneDropsRequest.decode(message.value);
  },
  toProto(message: QueryZoneDropsRequest): Uint8Array {
    return QueryZoneDropsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryZoneDropsRequest): QueryZoneDropsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropsRequest",
      value: QueryZoneDropsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryZoneDropsResponse(): QueryZoneDropsResponse {
  return {
    zoneDrops: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryZoneDropsResponse = {
  typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropsResponse",
  encode(message: QueryZoneDropsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.zoneDrops) {
      ZoneDrop.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZoneDropsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.zoneDrops.push(ZoneDrop.decode(reader, reader.uint32()));
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
  fromJSON(object: any): QueryZoneDropsResponse {
    const obj = createBaseQueryZoneDropsResponse();
    if (Array.isArray(object?.zoneDrops)) obj.zoneDrops = object.zoneDrops.map((e: any) => ZoneDrop.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryZoneDropsResponse): unknown {
    const obj: any = {};
    if (message.zoneDrops) {
      obj.zoneDrops = message.zoneDrops.map(e => e ? ZoneDrop.toJSON(e) : undefined);
    } else {
      obj.zoneDrops = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryZoneDropsResponse>): QueryZoneDropsResponse {
    const message = createBaseQueryZoneDropsResponse();
    message.zoneDrops = object.zoneDrops?.map(e => ZoneDrop.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryZoneDropsResponseSDKType): QueryZoneDropsResponse {
    return {
      zoneDrops: Array.isArray(object?.zone_drops) ? object.zone_drops.map((e: any) => ZoneDrop.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryZoneDropsResponse): QueryZoneDropsResponseSDKType {
    const obj: any = {};
    if (message.zoneDrops) {
      obj.zone_drops = message.zoneDrops.map(e => e ? ZoneDrop.toSDK(e) : undefined);
    } else {
      obj.zone_drops = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryZoneDropsResponseAmino): QueryZoneDropsResponse {
    return {
      zoneDrops: Array.isArray(object?.zone_drops) ? object.zone_drops.map((e: any) => ZoneDrop.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryZoneDropsResponse): QueryZoneDropsResponseAmino {
    const obj: any = {};
    if (message.zoneDrops) {
      obj.zone_drops = message.zoneDrops.map(e => e ? ZoneDrop.toAmino(e) : undefined);
    } else {
      obj.zone_drops = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryZoneDropsResponseAminoMsg): QueryZoneDropsResponse {
    return QueryZoneDropsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZoneDropsResponseProtoMsg): QueryZoneDropsResponse {
    return QueryZoneDropsResponse.decode(message.value);
  },
  toProto(message: QueryZoneDropsResponse): Uint8Array {
    return QueryZoneDropsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryZoneDropsResponse): QueryZoneDropsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryZoneDropsResponse",
      value: QueryZoneDropsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryClaimRecordRequest(): QueryClaimRecordRequest {
  return {
    chainId: "",
    address: ""
  };
}
export const QueryClaimRecordRequest = {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordRequest",
  encode(message: QueryClaimRecordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClaimRecordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryClaimRecordRequest {
    const obj = createBaseQueryClaimRecordRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.address)) obj.address = String(object.address);
    return obj;
  },
  toJSON(message: QueryClaimRecordRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryClaimRecordRequest>): QueryClaimRecordRequest {
    const message = createBaseQueryClaimRecordRequest();
    message.chainId = object.chainId ?? "";
    message.address = object.address ?? "";
    return message;
  },
  fromSDK(object: QueryClaimRecordRequestSDKType): QueryClaimRecordRequest {
    return {
      chainId: object?.chain_id,
      address: object?.address
    };
  },
  toSDK(message: QueryClaimRecordRequest): QueryClaimRecordRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.address = message.address;
    return obj;
  },
  fromAmino(object: QueryClaimRecordRequestAmino): QueryClaimRecordRequest {
    return {
      chainId: object.chain_id,
      address: object.address
    };
  },
  toAmino(message: QueryClaimRecordRequest): QueryClaimRecordRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: QueryClaimRecordRequestAminoMsg): QueryClaimRecordRequest {
    return QueryClaimRecordRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryClaimRecordRequestProtoMsg): QueryClaimRecordRequest {
    return QueryClaimRecordRequest.decode(message.value);
  },
  toProto(message: QueryClaimRecordRequest): Uint8Array {
    return QueryClaimRecordRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryClaimRecordRequest): QueryClaimRecordRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordRequest",
      value: QueryClaimRecordRequest.encode(message).finish()
    };
  }
};
function createBaseQueryClaimRecordResponse(): QueryClaimRecordResponse {
  return {
    claimRecord: ClaimRecord.fromPartial({})
  };
}
export const QueryClaimRecordResponse = {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordResponse",
  encode(message: QueryClaimRecordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.claimRecord !== undefined) {
      ClaimRecord.encode(message.claimRecord, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClaimRecordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.claimRecord = ClaimRecord.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryClaimRecordResponse {
    const obj = createBaseQueryClaimRecordResponse();
    if (isSet(object.claimRecord)) obj.claimRecord = ClaimRecord.fromJSON(object.claimRecord);
    return obj;
  },
  toJSON(message: QueryClaimRecordResponse): unknown {
    const obj: any = {};
    message.claimRecord !== undefined && (obj.claimRecord = message.claimRecord ? ClaimRecord.toJSON(message.claimRecord) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryClaimRecordResponse>): QueryClaimRecordResponse {
    const message = createBaseQueryClaimRecordResponse();
    if (object.claimRecord !== undefined && object.claimRecord !== null) {
      message.claimRecord = ClaimRecord.fromPartial(object.claimRecord);
    }
    return message;
  },
  fromSDK(object: QueryClaimRecordResponseSDKType): QueryClaimRecordResponse {
    return {
      claimRecord: object.claim_record ? ClaimRecord.fromSDK(object.claim_record) : undefined
    };
  },
  toSDK(message: QueryClaimRecordResponse): QueryClaimRecordResponseSDKType {
    const obj: any = {};
    message.claimRecord !== undefined && (obj.claim_record = message.claimRecord ? ClaimRecord.toSDK(message.claimRecord) : undefined);
    return obj;
  },
  fromAmino(object: QueryClaimRecordResponseAmino): QueryClaimRecordResponse {
    return {
      claimRecord: object?.claim_record ? ClaimRecord.fromAmino(object.claim_record) : undefined
    };
  },
  toAmino(message: QueryClaimRecordResponse): QueryClaimRecordResponseAmino {
    const obj: any = {};
    obj.claim_record = message.claimRecord ? ClaimRecord.toAmino(message.claimRecord) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryClaimRecordResponseAminoMsg): QueryClaimRecordResponse {
    return QueryClaimRecordResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryClaimRecordResponseProtoMsg): QueryClaimRecordResponse {
    return QueryClaimRecordResponse.decode(message.value);
  },
  toProto(message: QueryClaimRecordResponse): Uint8Array {
    return QueryClaimRecordResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryClaimRecordResponse): QueryClaimRecordResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordResponse",
      value: QueryClaimRecordResponse.encode(message).finish()
    };
  }
};
function createBaseQueryClaimRecordsRequest(): QueryClaimRecordsRequest {
  return {
    chainId: "",
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryClaimRecordsRequest = {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordsRequest",
  encode(message: QueryClaimRecordsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClaimRecordsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryClaimRecordsRequest {
    const obj = createBaseQueryClaimRecordsRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryClaimRecordsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryClaimRecordsRequest>): QueryClaimRecordsRequest {
    const message = createBaseQueryClaimRecordsRequest();
    message.chainId = object.chainId ?? "";
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryClaimRecordsRequestSDKType): QueryClaimRecordsRequest {
    return {
      chainId: object?.chain_id,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryClaimRecordsRequest): QueryClaimRecordsRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryClaimRecordsRequestAmino): QueryClaimRecordsRequest {
    return {
      chainId: object.chain_id,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryClaimRecordsRequest): QueryClaimRecordsRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryClaimRecordsRequestAminoMsg): QueryClaimRecordsRequest {
    return QueryClaimRecordsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryClaimRecordsRequestProtoMsg): QueryClaimRecordsRequest {
    return QueryClaimRecordsRequest.decode(message.value);
  },
  toProto(message: QueryClaimRecordsRequest): Uint8Array {
    return QueryClaimRecordsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryClaimRecordsRequest): QueryClaimRecordsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordsRequest",
      value: QueryClaimRecordsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryClaimRecordsResponse(): QueryClaimRecordsResponse {
  return {
    claimRecords: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryClaimRecordsResponse = {
  typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordsResponse",
  encode(message: QueryClaimRecordsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.claimRecords) {
      ClaimRecord.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClaimRecordsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.claimRecords.push(ClaimRecord.decode(reader, reader.uint32()));
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
  fromJSON(object: any): QueryClaimRecordsResponse {
    const obj = createBaseQueryClaimRecordsResponse();
    if (Array.isArray(object?.claimRecords)) obj.claimRecords = object.claimRecords.map((e: any) => ClaimRecord.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryClaimRecordsResponse): unknown {
    const obj: any = {};
    if (message.claimRecords) {
      obj.claimRecords = message.claimRecords.map(e => e ? ClaimRecord.toJSON(e) : undefined);
    } else {
      obj.claimRecords = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryClaimRecordsResponse>): QueryClaimRecordsResponse {
    const message = createBaseQueryClaimRecordsResponse();
    message.claimRecords = object.claimRecords?.map(e => ClaimRecord.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryClaimRecordsResponseSDKType): QueryClaimRecordsResponse {
    return {
      claimRecords: Array.isArray(object?.claim_records) ? object.claim_records.map((e: any) => ClaimRecord.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryClaimRecordsResponse): QueryClaimRecordsResponseSDKType {
    const obj: any = {};
    if (message.claimRecords) {
      obj.claim_records = message.claimRecords.map(e => e ? ClaimRecord.toSDK(e) : undefined);
    } else {
      obj.claim_records = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryClaimRecordsResponseAmino): QueryClaimRecordsResponse {
    return {
      claimRecords: Array.isArray(object?.claim_records) ? object.claim_records.map((e: any) => ClaimRecord.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryClaimRecordsResponse): QueryClaimRecordsResponseAmino {
    const obj: any = {};
    if (message.claimRecords) {
      obj.claim_records = message.claimRecords.map(e => e ? ClaimRecord.toAmino(e) : undefined);
    } else {
      obj.claim_records = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryClaimRecordsResponseAminoMsg): QueryClaimRecordsResponse {
    return QueryClaimRecordsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryClaimRecordsResponseProtoMsg): QueryClaimRecordsResponse {
    return QueryClaimRecordsResponse.decode(message.value);
  },
  toProto(message: QueryClaimRecordsResponse): Uint8Array {
    return QueryClaimRecordsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryClaimRecordsResponse): QueryClaimRecordsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.QueryClaimRecordsResponse",
      value: QueryClaimRecordsResponse.encode(message).finish()
    };
  }
};