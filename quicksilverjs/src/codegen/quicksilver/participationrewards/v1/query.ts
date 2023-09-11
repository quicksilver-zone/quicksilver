import { Params, ParamsSDKType } from "./participationrewards";
import * as _m0 from "protobufjs/minimal";
import { isSet, bytesFromBase64, base64FromBytes } from "../../../helpers";
/** QueryParamsRequest is the request type for the Query/Params RPC method. */

export interface QueryParamsRequest {}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */

export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */

export interface QueryParamsResponse {
  /** params defines the parameters of the module. */
  params?: Params;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */

export interface QueryParamsResponseSDKType {
  /** params defines the parameters of the module. */
  params?: ParamsSDKType;
}
/** QueryProtocolDataRequest is the request type for querying Protocol Data. */

export interface QueryProtocolDataRequest {
  type: string;
  key: string;
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
/** QueryProtocolDataResponse is the response type for querying Protocol Data. */

export interface QueryProtocolDataResponseSDKType {
  /** data defines the underlying protocol data. */
  data: Uint8Array[];
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
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
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: Partial<QueryParamsRequest>): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  }

};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    params: undefined
  };
}

export const QueryParamsResponse = {
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
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined
    };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryParamsResponse>): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryProtocolDataRequest(): QueryProtocolDataRequest {
  return {
    type: "",
    key: ""
  };
}

export const QueryProtocolDataRequest = {
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
    return {
      type: isSet(object.type) ? String(object.type) : "",
      key: isSet(object.key) ? String(object.key) : ""
    };
  },

  toJSON(message: QueryProtocolDataRequest): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.key !== undefined && (obj.key = message.key);
    return obj;
  },

  fromPartial(object: Partial<QueryProtocolDataRequest>): QueryProtocolDataRequest {
    const message = createBaseQueryProtocolDataRequest();
    message.type = object.type ?? "";
    message.key = object.key ?? "";
    return message;
  }

};

function createBaseQueryProtocolDataResponse(): QueryProtocolDataResponse {
  return {
    data: []
  };
}

export const QueryProtocolDataResponse = {
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
    return {
      data: Array.isArray(object?.data) ? object.data.map((e: any) => bytesFromBase64(e)) : []
    };
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

  fromPartial(object: Partial<QueryProtocolDataResponse>): QueryProtocolDataResponse {
    const message = createBaseQueryProtocolDataResponse();
    message.data = object.data?.map(e => e) || [];
    return message;
  }

};