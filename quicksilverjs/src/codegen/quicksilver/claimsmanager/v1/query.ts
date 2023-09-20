import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Claim, ClaimAmino, ClaimSDKType } from "./claimsmanager";
import * as _m0 from "protobufjs/minimal";
import { isSet, DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.claimsmanager.v1";
/** QueryClaimsRequest is the request type for the Query/Claims RPC method. */
export interface QueryClaimsRequest {
  chainId: string;
  address: string;
  pagination: PageRequest;
}
export interface QueryClaimsRequestProtoMsg {
  typeUrl: "/quicksilver.claimsmanager.v1.QueryClaimsRequest";
  value: Uint8Array;
}
/** QueryClaimsRequest is the request type for the Query/Claims RPC method. */
export interface QueryClaimsRequestAmino {
  chain_id: string;
  address: string;
  pagination?: PageRequestAmino;
}
export interface QueryClaimsRequestAminoMsg {
  type: "/quicksilver.claimsmanager.v1.QueryClaimsRequest";
  value: QueryClaimsRequestAmino;
}
/** QueryClaimsRequest is the request type for the Query/Claims RPC method. */
export interface QueryClaimsRequestSDKType {
  chain_id: string;
  address: string;
  pagination: PageRequestSDKType;
}
/** QueryClaimsResponse is the response type for the Query/Claims RPC method. */
export interface QueryClaimsResponse {
  claims: Claim[];
  pagination: PageResponse;
}
export interface QueryClaimsResponseProtoMsg {
  typeUrl: "/quicksilver.claimsmanager.v1.QueryClaimsResponse";
  value: Uint8Array;
}
/** QueryClaimsResponse is the response type for the Query/Claims RPC method. */
export interface QueryClaimsResponseAmino {
  claims: ClaimAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryClaimsResponseAminoMsg {
  type: "/quicksilver.claimsmanager.v1.QueryClaimsResponse";
  value: QueryClaimsResponseAmino;
}
/** QueryClaimsResponse is the response type for the Query/Claims RPC method. */
export interface QueryClaimsResponseSDKType {
  claims: ClaimSDKType[];
  pagination: PageResponseSDKType;
}
function createBaseQueryClaimsRequest(): QueryClaimsRequest {
  return {
    chainId: "",
    address: "",
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryClaimsRequest = {
  typeUrl: "/quicksilver.claimsmanager.v1.QueryClaimsRequest",
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
    const obj = createBaseQueryClaimsRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.address)) obj.address = String(object.address);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryClaimsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.address !== undefined && (obj.address = message.address);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryClaimsRequest>): QueryClaimsRequest {
    const message = createBaseQueryClaimsRequest();
    message.chainId = object.chainId ?? "";
    message.address = object.address ?? "";
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryClaimsRequestSDKType): QueryClaimsRequest {
    return {
      chainId: object?.chain_id,
      address: object?.address,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryClaimsRequest): QueryClaimsRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.address = message.address;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryClaimsRequestAmino): QueryClaimsRequest {
    return {
      chainId: object.chain_id,
      address: object.address,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryClaimsRequest): QueryClaimsRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.address = message.address;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryClaimsRequestAminoMsg): QueryClaimsRequest {
    return QueryClaimsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryClaimsRequestProtoMsg): QueryClaimsRequest {
    return QueryClaimsRequest.decode(message.value);
  },
  toProto(message: QueryClaimsRequest): Uint8Array {
    return QueryClaimsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryClaimsRequest): QueryClaimsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.claimsmanager.v1.QueryClaimsRequest",
      value: QueryClaimsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryClaimsResponse(): QueryClaimsResponse {
  return {
    claims: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryClaimsResponse = {
  typeUrl: "/quicksilver.claimsmanager.v1.QueryClaimsResponse",
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
    const obj = createBaseQueryClaimsResponse();
    if (Array.isArray(object?.claims)) obj.claims = object.claims.map((e: any) => Claim.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
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
  fromPartial(object: DeepPartial<QueryClaimsResponse>): QueryClaimsResponse {
    const message = createBaseQueryClaimsResponse();
    message.claims = object.claims?.map(e => Claim.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryClaimsResponseSDKType): QueryClaimsResponse {
    return {
      claims: Array.isArray(object?.claims) ? object.claims.map((e: any) => Claim.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryClaimsResponse): QueryClaimsResponseSDKType {
    const obj: any = {};
    if (message.claims) {
      obj.claims = message.claims.map(e => e ? Claim.toSDK(e) : undefined);
    } else {
      obj.claims = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryClaimsResponseAmino): QueryClaimsResponse {
    return {
      claims: Array.isArray(object?.claims) ? object.claims.map((e: any) => Claim.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryClaimsResponse): QueryClaimsResponseAmino {
    const obj: any = {};
    if (message.claims) {
      obj.claims = message.claims.map(e => e ? Claim.toAmino(e) : undefined);
    } else {
      obj.claims = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryClaimsResponseAminoMsg): QueryClaimsResponse {
    return QueryClaimsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryClaimsResponseProtoMsg): QueryClaimsResponse {
    return QueryClaimsResponse.decode(message.value);
  },
  toProto(message: QueryClaimsResponse): Uint8Array {
    return QueryClaimsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryClaimsResponse): QueryClaimsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.claimsmanager.v1.QueryClaimsResponse",
      value: QueryClaimsResponse.encode(message).finish()
    };
  }
};