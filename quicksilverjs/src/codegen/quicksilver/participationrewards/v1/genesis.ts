import { Params, ParamsSDKType, KeyedProtocolData, KeyedProtocolDataSDKType } from "./participationrewards";
import * as _m0 from "protobufjs/minimal";
import { isSet } from "../../../helpers";
/** GenesisState defines the participationrewards module's genesis state. */

export interface GenesisState {
  params?: Params;
  protocolData: KeyedProtocolData[];
}
/** GenesisState defines the participationrewards module's genesis state. */

export interface GenesisStateSDKType {
  params?: ParamsSDKType;
  protocol_data: KeyedProtocolDataSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    protocolData: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.protocolData) {
      KeyedProtocolData.encode(v!, writer.uint32(18).fork()).ldelim();
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
          message.protocolData.push(KeyedProtocolData.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): GenesisState {
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
      protocolData: Array.isArray(object?.protocolData) ? object.protocolData.map((e: any) => KeyedProtocolData.fromJSON(e)) : []
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);

    if (message.protocolData) {
      obj.protocolData = message.protocolData.map(e => e ? KeyedProtocolData.toJSON(e) : undefined);
    } else {
      obj.protocolData = [];
    }

    return obj;
  },

  fromPartial(object: Partial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.protocolData = object.protocolData?.map(e => KeyedProtocolData.fromPartial(e)) || [];
    return message;
  }

};