import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.airdrop.v1";
/** Params holds parameters for the airdrop module. */
export interface Params {}
export interface ParamsProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.Params";
  value: Uint8Array;
}
/** Params holds parameters for the airdrop module. */
export interface ParamsAmino {}
export interface ParamsAminoMsg {
  type: "/quicksilver.airdrop.v1.Params";
  value: ParamsAmino;
}
/** Params holds parameters for the airdrop module. */
export interface ParamsSDKType {}
function createBaseParams(): Params {
  return {};
}
export const Params = {
  typeUrl: "/quicksilver.airdrop.v1.Params",
  encode(_: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
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
  fromJSON(_: any): Params {
    const obj = createBaseParams();
    return obj;
  },
  toJSON(_: Params): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<Params>): Params {
    const message = createBaseParams();
    return message;
  },
  fromSDK(_: ParamsSDKType): Params {
    return {};
  },
  toSDK(_: Params): ParamsSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: ParamsAmino): Params {
    return {};
  },
  toAmino(_: Params): ParamsAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: ParamsAminoMsg): Params {
    return Params.fromAmino(object.value);
  },
  fromProtoMsg(message: ParamsProtoMsg): Params {
    return Params.decode(message.value);
  },
  toProto(message: Params): Uint8Array {
    return Params.encode(message).finish();
  },
  toProtoMsg(message: Params): ParamsProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.Params",
      value: Params.encode(message).finish()
    };
  }
};