import { Coin, CoinAmino, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import * as _m0 from "protobufjs/minimal";
import { isSet, DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.interchainstaking.v1";
/**
 * MsgRequestRedemption represents a message type to request a burn of qAssets
 * for native assets.
 */
export interface MsgRequestRedemption {
  value: Coin;
  destinationAddress: string;
  fromAddress: string;
}
export interface MsgRequestRedemptionProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption";
  value: Uint8Array;
}
/**
 * MsgRequestRedemption represents a message type to request a burn of qAssets
 * for native assets.
 */
export interface MsgRequestRedemptionAmino {
  value?: CoinAmino;
  destination_address: string;
  from_address: string;
}
export interface MsgRequestRedemptionAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgRequestRedemption";
  value: MsgRequestRedemptionAmino;
}
/**
 * MsgRequestRedemption represents a message type to request a burn of qAssets
 * for native assets.
 */
export interface MsgRequestRedemptionSDKType {
  value: CoinSDKType;
  destination_address: string;
  from_address: string;
}
/**
 * MsgSignalIntent represents a message type for signalling voting intent for
 * one or more validators.
 */
export interface MsgSignalIntent {
  chainId: string;
  intents: string;
  fromAddress: string;
}
export interface MsgSignalIntentProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent";
  value: Uint8Array;
}
/**
 * MsgSignalIntent represents a message type for signalling voting intent for
 * one or more validators.
 */
export interface MsgSignalIntentAmino {
  chain_id: string;
  intents: string;
  from_address: string;
}
export interface MsgSignalIntentAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgSignalIntent";
  value: MsgSignalIntentAmino;
}
/**
 * MsgSignalIntent represents a message type for signalling voting intent for
 * one or more validators.
 */
export interface MsgSignalIntentSDKType {
  chain_id: string;
  intents: string;
  from_address: string;
}
/** MsgRequestRedemptionResponse defines the MsgRequestRedemption response type. */
export interface MsgRequestRedemptionResponse {}
export interface MsgRequestRedemptionResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemptionResponse";
  value: Uint8Array;
}
/** MsgRequestRedemptionResponse defines the MsgRequestRedemption response type. */
export interface MsgRequestRedemptionResponseAmino {}
export interface MsgRequestRedemptionResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgRequestRedemptionResponse";
  value: MsgRequestRedemptionResponseAmino;
}
/** MsgRequestRedemptionResponse defines the MsgRequestRedemption response type. */
export interface MsgRequestRedemptionResponseSDKType {}
/** MsgSignalIntentResponse defines the MsgSignalIntent response type. */
export interface MsgSignalIntentResponse {}
export interface MsgSignalIntentResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntentResponse";
  value: Uint8Array;
}
/** MsgSignalIntentResponse defines the MsgSignalIntent response type. */
export interface MsgSignalIntentResponseAmino {}
export interface MsgSignalIntentResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgSignalIntentResponse";
  value: MsgSignalIntentResponseAmino;
}
/** MsgSignalIntentResponse defines the MsgSignalIntent response type. */
export interface MsgSignalIntentResponseSDKType {}
function createBaseMsgRequestRedemption(): MsgRequestRedemption {
  return {
    value: Coin.fromPartial({}),
    destinationAddress: "",
    fromAddress: ""
  };
}
export const MsgRequestRedemption = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
  encode(message: MsgRequestRedemption, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.value !== undefined) {
      Coin.encode(message.value, writer.uint32(10).fork()).ldelim();
    }
    if (message.destinationAddress !== "") {
      writer.uint32(18).string(message.destinationAddress);
    }
    if (message.fromAddress !== "") {
      writer.uint32(26).string(message.fromAddress);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestRedemption {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRequestRedemption();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.value = Coin.decode(reader, reader.uint32());
          break;
        case 2:
          message.destinationAddress = reader.string();
          break;
        case 3:
          message.fromAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgRequestRedemption {
    const obj = createBaseMsgRequestRedemption();
    if (isSet(object.value)) obj.value = Coin.fromJSON(object.value);
    if (isSet(object.destinationAddress)) obj.destinationAddress = String(object.destinationAddress);
    if (isSet(object.fromAddress)) obj.fromAddress = String(object.fromAddress);
    return obj;
  },
  toJSON(message: MsgRequestRedemption): unknown {
    const obj: any = {};
    message.value !== undefined && (obj.value = message.value ? Coin.toJSON(message.value) : undefined);
    message.destinationAddress !== undefined && (obj.destinationAddress = message.destinationAddress);
    message.fromAddress !== undefined && (obj.fromAddress = message.fromAddress);
    return obj;
  },
  fromPartial(object: DeepPartial<MsgRequestRedemption>): MsgRequestRedemption {
    const message = createBaseMsgRequestRedemption();
    if (object.value !== undefined && object.value !== null) {
      message.value = Coin.fromPartial(object.value);
    }
    message.destinationAddress = object.destinationAddress ?? "";
    message.fromAddress = object.fromAddress ?? "";
    return message;
  },
  fromSDK(object: MsgRequestRedemptionSDKType): MsgRequestRedemption {
    return {
      value: object.value ? Coin.fromSDK(object.value) : undefined,
      destinationAddress: object?.destination_address,
      fromAddress: object?.from_address
    };
  },
  toSDK(message: MsgRequestRedemption): MsgRequestRedemptionSDKType {
    const obj: any = {};
    message.value !== undefined && (obj.value = message.value ? Coin.toSDK(message.value) : undefined);
    obj.destination_address = message.destinationAddress;
    obj.from_address = message.fromAddress;
    return obj;
  },
  fromAmino(object: MsgRequestRedemptionAmino): MsgRequestRedemption {
    return {
      value: object?.value ? Coin.fromAmino(object.value) : undefined,
      destinationAddress: object.destination_address,
      fromAddress: object.from_address
    };
  },
  toAmino(message: MsgRequestRedemption): MsgRequestRedemptionAmino {
    const obj: any = {};
    obj.value = message.value ? Coin.toAmino(message.value) : undefined;
    obj.destination_address = message.destinationAddress;
    obj.from_address = message.fromAddress;
    return obj;
  },
  fromAminoMsg(object: MsgRequestRedemptionAminoMsg): MsgRequestRedemption {
    return MsgRequestRedemption.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRequestRedemptionProtoMsg): MsgRequestRedemption {
    return MsgRequestRedemption.decode(message.value);
  },
  toProto(message: MsgRequestRedemption): Uint8Array {
    return MsgRequestRedemption.encode(message).finish();
  },
  toProtoMsg(message: MsgRequestRedemption): MsgRequestRedemptionProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
      value: MsgRequestRedemption.encode(message).finish()
    };
  }
};
function createBaseMsgSignalIntent(): MsgSignalIntent {
  return {
    chainId: "",
    intents: "",
    fromAddress: ""
  };
}
export const MsgSignalIntent = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
  encode(message: MsgSignalIntent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (message.intents !== "") {
      writer.uint32(18).string(message.intents);
    }
    if (message.fromAddress !== "") {
      writer.uint32(26).string(message.fromAddress);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSignalIntent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSignalIntent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.intents = reader.string();
          break;
        case 3:
          message.fromAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgSignalIntent {
    const obj = createBaseMsgSignalIntent();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.intents)) obj.intents = String(object.intents);
    if (isSet(object.fromAddress)) obj.fromAddress = String(object.fromAddress);
    return obj;
  },
  toJSON(message: MsgSignalIntent): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.intents !== undefined && (obj.intents = message.intents);
    message.fromAddress !== undefined && (obj.fromAddress = message.fromAddress);
    return obj;
  },
  fromPartial(object: DeepPartial<MsgSignalIntent>): MsgSignalIntent {
    const message = createBaseMsgSignalIntent();
    message.chainId = object.chainId ?? "";
    message.intents = object.intents ?? "";
    message.fromAddress = object.fromAddress ?? "";
    return message;
  },
  fromSDK(object: MsgSignalIntentSDKType): MsgSignalIntent {
    return {
      chainId: object?.chain_id,
      intents: object?.intents,
      fromAddress: object?.from_address
    };
  },
  toSDK(message: MsgSignalIntent): MsgSignalIntentSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.intents = message.intents;
    obj.from_address = message.fromAddress;
    return obj;
  },
  fromAmino(object: MsgSignalIntentAmino): MsgSignalIntent {
    return {
      chainId: object.chain_id,
      intents: object.intents,
      fromAddress: object.from_address
    };
  },
  toAmino(message: MsgSignalIntent): MsgSignalIntentAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.intents = message.intents;
    obj.from_address = message.fromAddress;
    return obj;
  },
  fromAminoMsg(object: MsgSignalIntentAminoMsg): MsgSignalIntent {
    return MsgSignalIntent.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSignalIntentProtoMsg): MsgSignalIntent {
    return MsgSignalIntent.decode(message.value);
  },
  toProto(message: MsgSignalIntent): Uint8Array {
    return MsgSignalIntent.encode(message).finish();
  },
  toProtoMsg(message: MsgSignalIntent): MsgSignalIntentProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
      value: MsgSignalIntent.encode(message).finish()
    };
  }
};
function createBaseMsgRequestRedemptionResponse(): MsgRequestRedemptionResponse {
  return {};
}
export const MsgRequestRedemptionResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemptionResponse",
  encode(_: MsgRequestRedemptionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestRedemptionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRequestRedemptionResponse();
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
  fromJSON(_: any): MsgRequestRedemptionResponse {
    const obj = createBaseMsgRequestRedemptionResponse();
    return obj;
  },
  toJSON(_: MsgRequestRedemptionResponse): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<MsgRequestRedemptionResponse>): MsgRequestRedemptionResponse {
    const message = createBaseMsgRequestRedemptionResponse();
    return message;
  },
  fromSDK(_: MsgRequestRedemptionResponseSDKType): MsgRequestRedemptionResponse {
    return {};
  },
  toSDK(_: MsgRequestRedemptionResponse): MsgRequestRedemptionResponseSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: MsgRequestRedemptionResponseAmino): MsgRequestRedemptionResponse {
    return {};
  },
  toAmino(_: MsgRequestRedemptionResponse): MsgRequestRedemptionResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgRequestRedemptionResponseAminoMsg): MsgRequestRedemptionResponse {
    return MsgRequestRedemptionResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRequestRedemptionResponseProtoMsg): MsgRequestRedemptionResponse {
    return MsgRequestRedemptionResponse.decode(message.value);
  },
  toProto(message: MsgRequestRedemptionResponse): Uint8Array {
    return MsgRequestRedemptionResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgRequestRedemptionResponse): MsgRequestRedemptionResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemptionResponse",
      value: MsgRequestRedemptionResponse.encode(message).finish()
    };
  }
};
function createBaseMsgSignalIntentResponse(): MsgSignalIntentResponse {
  return {};
}
export const MsgSignalIntentResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntentResponse",
  encode(_: MsgSignalIntentResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSignalIntentResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSignalIntentResponse();
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
  fromJSON(_: any): MsgSignalIntentResponse {
    const obj = createBaseMsgSignalIntentResponse();
    return obj;
  },
  toJSON(_: MsgSignalIntentResponse): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<MsgSignalIntentResponse>): MsgSignalIntentResponse {
    const message = createBaseMsgSignalIntentResponse();
    return message;
  },
  fromSDK(_: MsgSignalIntentResponseSDKType): MsgSignalIntentResponse {
    return {};
  },
  toSDK(_: MsgSignalIntentResponse): MsgSignalIntentResponseSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: MsgSignalIntentResponseAmino): MsgSignalIntentResponse {
    return {};
  },
  toAmino(_: MsgSignalIntentResponse): MsgSignalIntentResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgSignalIntentResponseAminoMsg): MsgSignalIntentResponse {
    return MsgSignalIntentResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSignalIntentResponseProtoMsg): MsgSignalIntentResponse {
    return MsgSignalIntentResponse.decode(message.value);
  },
  toProto(message: MsgSignalIntentResponse): Uint8Array {
    return MsgSignalIntentResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSignalIntentResponse): MsgSignalIntentResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntentResponse",
      value: MsgSignalIntentResponse.encode(message).finish()
    };
  }
};