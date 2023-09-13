import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsSDKType, Claim, ClaimSDKType } from "./claimsmanager";
import * as _m0 from "protobufjs/minimal";
import { isSet } from "../../../helpers";
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
export interface QueryClaimsRequest {
  chainId: string;
  address: string;
  pagination?: PageRequest;
}
export interface QueryClaimsRequestSDKType {
  chain_id: string;
  address: string;
  pagination?: PageRequestSDKType;
}
export interface QueryClaimsResponse {
  claims: Claim[];
  pagination?: PageResponse;
}
export interface QueryClaimsResponseSDKType {
  claims: ClaimSDKType[];
  pagination?: PageResponseSDKType;
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

function createBaseQueryClaimsRequest(): QueryClaimsRequest {
  return {
    chainId: "",
    address: "",
    pagination: undefined
  };
}

export const QueryClaimsRequest = {
  encode(message: QueryClaimsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }

    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClaimsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        case 2:
          message.address = reader.string();
          break;

        case 3:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryClaimsRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      address: isSet(object.address) ? String(object.address) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryClaimsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.address !== undefined && (obj.address = message.address);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryClaimsRequest>): QueryClaimsRequest {
    const message = createBaseQueryClaimsRequest();
    message.chainId = object.chainId ?? "";
    message.address = object.address ?? "";
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryClaimsResponse(): QueryClaimsResponse {
  return {
    claims: [],
    pagination: undefined
  };
}

export const QueryClaimsResponse = {
  encode(message: QueryClaimsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.claims) {
      Claim.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClaimsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.claims.push(Claim.decode(reader, reader.uint32()));
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

  fromJSON(object: any): QueryClaimsResponse {
    return {
      claims: Array.isArray(object?.claims) ? object.claims.map((e: any) => Claim.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryClaimsResponse): unknown {
    const obj: any = {};

    if (message.claims) {
      obj.claims = message.claims.map(e => e ? Claim.toJSON(e) : undefined);
    } else {
      obj.claims = [];
    }

    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryClaimsResponse>): QueryClaimsResponse {
    const message = createBaseQueryClaimsResponse();
    message.claims = object.claims?.map(e => Claim.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};