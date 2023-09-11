import { ZoneDrop, ZoneDropSDKType } from "./airdrop";
import * as _m0 from "protobufjs/minimal";
import { isSet, bytesFromBase64, base64FromBytes } from "../../../helpers";
export interface RegisterZoneDropProposal {
  title: string;
  description: string;
  zoneDrop?: ZoneDrop;
  claimRecords: Uint8Array;
}
export interface RegisterZoneDropProposalSDKType {
  title: string;
  description: string;
  zone_drop?: ZoneDropSDKType;
  claim_records: Uint8Array;
}

function createBaseRegisterZoneDropProposal(): RegisterZoneDropProposal {
  return {
    title: "",
    description: "",
    zoneDrop: undefined,
    claimRecords: new Uint8Array()
  };
}

export const RegisterZoneDropProposal = {
  encode(message: RegisterZoneDropProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    if (message.zoneDrop !== undefined) {
      ZoneDrop.encode(message.zoneDrop, writer.uint32(26).fork()).ldelim();
    }

    if (message.claimRecords.length !== 0) {
      writer.uint32(34).bytes(message.claimRecords);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegisterZoneDropProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegisterZoneDropProposal();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;

        case 2:
          message.description = reader.string();
          break;

        case 3:
          message.zoneDrop = ZoneDrop.decode(reader, reader.uint32());
          break;

        case 4:
          message.claimRecords = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): RegisterZoneDropProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      zoneDrop: isSet(object.zoneDrop) ? ZoneDrop.fromJSON(object.zoneDrop) : undefined,
      claimRecords: isSet(object.claimRecords) ? bytesFromBase64(object.claimRecords) : new Uint8Array()
    };
  },

  toJSON(message: RegisterZoneDropProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.zoneDrop !== undefined && (obj.zoneDrop = message.zoneDrop ? ZoneDrop.toJSON(message.zoneDrop) : undefined);
    message.claimRecords !== undefined && (obj.claimRecords = base64FromBytes(message.claimRecords !== undefined ? message.claimRecords : new Uint8Array()));
    return obj;
  },

  fromPartial(object: Partial<RegisterZoneDropProposal>): RegisterZoneDropProposal {
    const message = createBaseRegisterZoneDropProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.zoneDrop = object.zoneDrop !== undefined && object.zoneDrop !== null ? ZoneDrop.fromPartial(object.zoneDrop) : undefined;
    message.claimRecords = object.claimRecords ?? new Uint8Array();
    return message;
  }

};