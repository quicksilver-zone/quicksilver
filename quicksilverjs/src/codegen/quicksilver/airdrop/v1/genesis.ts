import { Params, ParamsSDKType } from "./params";
import { ZoneDrop, ZoneDropSDKType, ClaimRecord, ClaimRecordSDKType } from "./airdrop";
import * as _m0 from "protobufjs/minimal";
import { isSet } from "../../../helpers";
/** GenesisState defines the airdrop module's genesis state. */

export interface GenesisState {
  params?: Params;
  zoneDrops: ZoneDrop[];
  claimRecords: ClaimRecord[];
}
/** GenesisState defines the airdrop module's genesis state. */

export interface GenesisStateSDKType {
  params?: ParamsSDKType;
  zone_drops: ZoneDropSDKType[];
  claim_records: ClaimRecordSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    zoneDrops: [],
    claimRecords: []
  };
}

export const GenesisState = {
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
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
      zoneDrops: Array.isArray(object?.zoneDrops) ? object.zoneDrops.map((e: any) => ZoneDrop.fromJSON(e)) : [],
      claimRecords: Array.isArray(object?.claimRecords) ? object.claimRecords.map((e: any) => ClaimRecord.fromJSON(e)) : []
    };
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

  fromPartial(object: Partial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.zoneDrops = object.zoneDrops?.map(e => ZoneDrop.fromPartial(e)) || [];
    message.claimRecords = object.claimRecords?.map(e => ClaimRecord.fromPartial(e)) || [];
    return message;
  }

};