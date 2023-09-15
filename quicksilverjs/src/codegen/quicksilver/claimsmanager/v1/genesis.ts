import { Params, ParamsAmino, ParamsSDKType, Claim, ClaimAmino, ClaimSDKType } from "./claimsmanager";
import * as _m0 from "protobufjs/minimal";
import { isSet, DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.claimsmanager.v1";
/** GenesisState defines the claimsmanager module's genesis state. */
export interface GenesisState {
  params: Params;
  claims: Claim[];
}
export interface GenesisStateProtoMsg {
  typeUrl: "/quicksilver.claimsmanager.v1.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the claimsmanager module's genesis state. */
export interface GenesisStateAmino {
  params?: ParamsAmino;
  claims: ClaimAmino[];
}
export interface GenesisStateAminoMsg {
  type: "/quicksilver.claimsmanager.v1.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the claimsmanager module's genesis state. */
export interface GenesisStateSDKType {
  params: ParamsSDKType;
  claims: ClaimSDKType[];
}
function createBaseGenesisState(): GenesisState {
  return {
    params: Params.fromPartial({}),
    claims: []
  };
}
export const GenesisState = {
  typeUrl: "/quicksilver.claimsmanager.v1.GenesisState",
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.claims) {
      Claim.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.claims.push(Claim.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): GenesisState {
    const obj = createBaseGenesisState();
    if (isSet(object.params)) obj.params = Params.fromJSON(object.params);
    if (Array.isArray(object?.claims)) obj.claims = object.claims.map((e: any) => Claim.fromJSON(e));
    return obj;
  },
  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.claims) {
      obj.claims = message.claims.map(e => e ? Claim.toJSON(e) : undefined);
    } else {
      obj.claims = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    }
    message.claims = object.claims?.map(e => Claim.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: GenesisStateSDKType): GenesisState {
    return {
      params: object.params ? Params.fromSDK(object.params) : undefined,
      claims: Array.isArray(object?.claims) ? object.claims.map((e: any) => Claim.fromSDK(e)) : []
    };
  },
  toSDK(message: GenesisState): GenesisStateSDKType {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toSDK(message.params) : undefined);
    if (message.claims) {
      obj.claims = message.claims.map(e => e ? Claim.toSDK(e) : undefined);
    } else {
      obj.claims = [];
    }
    return obj;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    return {
      params: object?.params ? Params.fromAmino(object.params) : undefined,
      claims: Array.isArray(object?.claims) ? object.claims.map((e: any) => Claim.fromAmino(e)) : []
    };
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    if (message.claims) {
      obj.claims = message.claims.map(e => e ? Claim.toAmino(e) : undefined);
    } else {
      obj.claims = [];
    }
    return obj;
  },
  fromAminoMsg(object: GenesisStateAminoMsg): GenesisState {
    return GenesisState.fromAmino(object.value);
  },
  fromProtoMsg(message: GenesisStateProtoMsg): GenesisState {
    return GenesisState.decode(message.value);
  },
  toProto(message: GenesisState): Uint8Array {
    return GenesisState.encode(message).finish();
  },
  toProtoMsg(message: GenesisState): GenesisStateProtoMsg {
    return {
      typeUrl: "/quicksilver.claimsmanager.v1.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};