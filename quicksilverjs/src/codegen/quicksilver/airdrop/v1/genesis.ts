import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { ZoneDrop, ZoneDropAmino, ZoneDropSDKType, ClaimRecord, ClaimRecordAmino, ClaimRecordSDKType } from "./airdrop";
import * as _m0 from "protobufjs/minimal";
import { isSet, DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.airdrop.v1";
/** GenesisState defines the airdrop module's genesis state. */
export interface GenesisState {
  params: Params;
  zoneDrops: ZoneDrop[];
  claimRecords: ClaimRecord[];
}
export interface GenesisStateProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the airdrop module's genesis state. */
export interface GenesisStateAmino {
  params?: ParamsAmino;
  zone_drops: ZoneDropAmino[];
  claim_records: ClaimRecordAmino[];
}
export interface GenesisStateAminoMsg {
  type: "/quicksilver.airdrop.v1.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the airdrop module's genesis state. */
export interface GenesisStateSDKType {
  params: ParamsSDKType;
  zone_drops: ZoneDropSDKType[];
  claim_records: ClaimRecordSDKType[];
}
function createBaseGenesisState(): GenesisState {
  return {
    params: Params.fromPartial({}),
    zoneDrops: [],
    claimRecords: []
  };
}
export const GenesisState = {
  typeUrl: "/quicksilver.airdrop.v1.GenesisState",
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.zoneDrops) {
      ZoneDrop.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.claimRecords) {
      ClaimRecord.encode(v!, writer.uint32(26).fork()).ldelim();
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
          message.zoneDrops.push(ZoneDrop.decode(reader, reader.uint32()));
          break;
        case 3:
          message.claimRecords.push(ClaimRecord.decode(reader, reader.uint32()));
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
    if (Array.isArray(object?.zoneDrops)) obj.zoneDrops = object.zoneDrops.map((e: any) => ZoneDrop.fromJSON(e));
    if (Array.isArray(object?.claimRecords)) obj.claimRecords = object.claimRecords.map((e: any) => ClaimRecord.fromJSON(e));
    return obj;
  },
  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.zoneDrops) {
      obj.zoneDrops = message.zoneDrops.map(e => e ? ZoneDrop.toJSON(e) : undefined);
    } else {
      obj.zoneDrops = [];
    }
    if (message.claimRecords) {
      obj.claimRecords = message.claimRecords.map(e => e ? ClaimRecord.toJSON(e) : undefined);
    } else {
      obj.claimRecords = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    }
    message.zoneDrops = object.zoneDrops?.map(e => ZoneDrop.fromPartial(e)) || [];
    message.claimRecords = object.claimRecords?.map(e => ClaimRecord.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: GenesisStateSDKType): GenesisState {
    return {
      params: object.params ? Params.fromSDK(object.params) : undefined,
      zoneDrops: Array.isArray(object?.zone_drops) ? object.zone_drops.map((e: any) => ZoneDrop.fromSDK(e)) : [],
      claimRecords: Array.isArray(object?.claim_records) ? object.claim_records.map((e: any) => ClaimRecord.fromSDK(e)) : []
    };
  },
  toSDK(message: GenesisState): GenesisStateSDKType {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toSDK(message.params) : undefined);
    if (message.zoneDrops) {
      obj.zone_drops = message.zoneDrops.map(e => e ? ZoneDrop.toSDK(e) : undefined);
    } else {
      obj.zone_drops = [];
    }
    if (message.claimRecords) {
      obj.claim_records = message.claimRecords.map(e => e ? ClaimRecord.toSDK(e) : undefined);
    } else {
      obj.claim_records = [];
    }
    return obj;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    return {
      params: object?.params ? Params.fromAmino(object.params) : undefined,
      zoneDrops: Array.isArray(object?.zone_drops) ? object.zone_drops.map((e: any) => ZoneDrop.fromAmino(e)) : [],
      claimRecords: Array.isArray(object?.claim_records) ? object.claim_records.map((e: any) => ClaimRecord.fromAmino(e)) : []
    };
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    if (message.zoneDrops) {
      obj.zone_drops = message.zoneDrops.map(e => e ? ZoneDrop.toAmino(e) : undefined);
    } else {
      obj.zone_drops = [];
    }
    if (message.claimRecords) {
      obj.claim_records = message.claimRecords.map(e => e ? ClaimRecord.toAmino(e) : undefined);
    } else {
      obj.claim_records = [];
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
      typeUrl: "/quicksilver.airdrop.v1.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};