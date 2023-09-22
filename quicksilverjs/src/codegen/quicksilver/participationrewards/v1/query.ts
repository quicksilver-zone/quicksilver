import { Params, ParamsAmino, ParamsSDKType } from "./participationrewards";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, isSet, bytesFromBase64, base64FromBytes } from "../../../helpers";
export const protobufPackage = "quicksilver.participationrewards.v1";
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequest {}
export interface QueryParamsRequestProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.QueryParamsRequest";
  value: Uint8Array;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequestAmino {}
export interface QueryParamsRequestAminoMsg {
  type: "/quicksilver.participationrewards.v1.QueryParamsRequest";
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
  typeUrl: "/quicksilver.participationrewards.v1.QueryParamsResponse";
  value: Uint8Array;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponseAmino {
  /** params defines the parameters of the module. */
  params?: ParamsAmino;
}
export interface QueryParamsResponseAminoMsg {
  type: "/quicksilver.participationrewards.v1.QueryParamsResponse";
  value: QueryParamsResponseAmino;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponseSDKType {
  params: ParamsSDKType;
}
/** QueryProtocolDataRequest is the request type for querying Protocol Data. */
export interface QueryProtocolDataRequest {
  type: string;
  key: string;
}
export interface QueryProtocolDataRequestProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.QueryProtocolDataRequest";
  value: Uint8Array;
}
/** QueryProtocolDataRequest is the request type for querying Protocol Data. */
export interface QueryProtocolDataRequestAmino {
  type: string;
  key: string;
}
export interface QueryProtocolDataRequestAminoMsg {
  type: "/quicksilver.participationrewards.v1.QueryProtocolDataRequest";
  value: QueryProtocolDataRequestAmino;
}
/** QueryProtocolDataRequest is the request type for querying Protocol Data. */
export interface QueryProtocolDataRequestSDKType {
  type: string;
  key: string;
}
/** QueryProtocolDataResponse is the response type for querying Protocol Data. */
export interface QueryProtocolDataResponse {
  /** data defines the underlying protocol data. */
  data: Uint8Array[];
}
export interface QueryProtocolDataResponseProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.QueryProtocolDataResponse";
  value: Uint8Array;
}
/** QueryProtocolDataResponse is the response type for querying Protocol Data. */
export interface QueryProtocolDataResponseAmino {
  /** data defines the underlying protocol data. */
  data: Uint8Array[];
}
export interface QueryProtocolDataResponseAminoMsg {
  type: "/quicksilver.participationrewards.v1.QueryProtocolDataResponse";
  value: QueryProtocolDataResponseAmino;
}
/** QueryProtocolDataResponse is the response type for querying Protocol Data. */
export interface QueryProtocolDataResponseSDKType {
  data: Uint8Array[];
}
function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}
export const QueryParamsRequest = {
  typeUrl: "/quicksilver.participationrewards.v1.QueryParamsRequest",
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
      typeUrl: "/quicksilver.participationrewards.v1.QueryParamsRequest",
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
  typeUrl: "/quicksilver.participationrewards.v1.QueryParamsResponse",
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
      typeUrl: "/quicksilver.participationrewards.v1.QueryParamsResponse",
      value: QueryParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryProtocolDataRequest(): QueryProtocolDataRequest {
  return {
    type: "",
    key: ""
  };
}
export const QueryProtocolDataRequest = {
  typeUrl: "/quicksilver.participationrewards.v1.QueryProtocolDataRequest",
  encode(message: QueryProtocolDataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.key !== "") {
      writer.uint32(18).string(message.key);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryProtocolDataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProtocolDataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.string();
          break;
        case 2:
          message.key = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryProtocolDataRequest {
    const obj = createBaseQueryProtocolDataRequest();
    if (isSet(object.type)) obj.type = String(object.type);
    if (isSet(object.key)) obj.key = String(object.key);
    return obj;
  },
  toJSON(message: QueryProtocolDataRequest): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.key !== undefined && (obj.key = message.key);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryProtocolDataRequest>): QueryProtocolDataRequest {
    const message = createBaseQueryProtocolDataRequest();
    message.type = object.type ?? "";
    message.key = object.key ?? "";
    return message;
  },
  fromSDK(object: QueryProtocolDataRequestSDKType): QueryProtocolDataRequest {
    return {
      type: object?.type,
      key: object?.key
    };
  },
  toSDK(message: QueryProtocolDataRequest): QueryProtocolDataRequestSDKType {
    const obj: any = {};
    obj.type = message.type;
    obj.key = message.key;
    return obj;
  },
  fromAmino(object: QueryProtocolDataRequestAmino): QueryProtocolDataRequest {
    return {
      type: object.type,
      key: object.key
    };
  },
  toAmino(message: QueryProtocolDataRequest): QueryProtocolDataRequestAmino {
    const obj: any = {};
    obj.type = message.type;
    obj.key = message.key;
    return obj;
  },
  fromAminoMsg(object: QueryProtocolDataRequestAminoMsg): QueryProtocolDataRequest {
    return QueryProtocolDataRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryProtocolDataRequestProtoMsg): QueryProtocolDataRequest {
    return QueryProtocolDataRequest.decode(message.value);
  },
  toProto(message: QueryProtocolDataRequest): Uint8Array {
    return QueryProtocolDataRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryProtocolDataRequest): QueryProtocolDataRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.QueryProtocolDataRequest",
      value: QueryProtocolDataRequest.encode(message).finish()
    };
  }
};
function createBaseQueryProtocolDataResponse(): QueryProtocolDataResponse {
  return {
    data: []
  };
}
export const QueryProtocolDataResponse = {
  typeUrl: "/quicksilver.participationrewards.v1.QueryProtocolDataResponse",
  encode(message: QueryProtocolDataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.data) {
      writer.uint32(10).bytes(v!);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): QueryProtocolDataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProtocolDataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data.push(reader.bytes());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): QueryProtocolDataResponse {
    const obj = createBaseQueryProtocolDataResponse();
    if (Array.isArray(object?.data)) obj.data = object.data.map((e: any) => bytesFromBase64(e));
    return obj;
  },
  toJSON(message: QueryProtocolDataResponse): unknown {
    const obj: any = {};
    if (message.data) {
      obj.data = message.data.map(e => base64FromBytes(e !== undefined ? e : new Uint8Array()));
    } else {
      obj.data = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<QueryProtocolDataResponse>): QueryProtocolDataResponse {
    const message = createBaseQueryProtocolDataResponse();
    message.data = object.data?.map(e => e) || [];
    return message;
  },
  fromSDK(object: QueryProtocolDataResponseSDKType): QueryProtocolDataResponse {
    return {
      data: Array.isArray(object?.data) ? object.data.map((e: any) => e) : []
    };
  },
  toSDK(message: QueryProtocolDataResponse): QueryProtocolDataResponseSDKType {
    const obj: any = {};
    if (message.data) {
      obj.data = message.data.map(e => e);
    } else {
      obj.data = [];
    }
    return obj;
  },
  fromAmino(object: QueryProtocolDataResponseAmino): QueryProtocolDataResponse {
    return {
      data: Array.isArray(object?.data) ? object.data.map((e: any) => e) : []
    };
  },
  toAmino(message: QueryProtocolDataResponse): QueryProtocolDataResponseAmino {
    const obj: any = {};
    if (message.data) {
      obj.data = message.data.map(e => e);
    } else {
      obj.data = [];
    }
    return obj;
  },
  fromAminoMsg(object: QueryProtocolDataResponseAminoMsg): QueryProtocolDataResponse {
    return QueryProtocolDataResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryProtocolDataResponseProtoMsg): QueryProtocolDataResponse {
    return QueryProtocolDataResponse.decode(message.value);
  },
  toProto(message: QueryProtocolDataResponse): Uint8Array {
    return QueryProtocolDataResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryProtocolDataResponse): QueryProtocolDataResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.QueryProtocolDataResponse",
      value: QueryProtocolDataResponse.encode(message).finish()
    };
  }
};