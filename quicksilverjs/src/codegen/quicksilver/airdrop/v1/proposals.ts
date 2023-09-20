import { ZoneDrop, ZoneDropAmino, ZoneDropSDKType } from "./airdrop";
import * as _m0 from "protobufjs/minimal";
import { isSet, bytesFromBase64, base64FromBytes, DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.airdrop.v1";
export interface RegisterZoneDropProposal {
  title: string;
  description: string;
  zoneDrop: ZoneDrop;
  claimRecords: Uint8Array;
}
export interface RegisterZoneDropProposalProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.RegisterZoneDropProposal";
  value: Uint8Array;
}
export interface RegisterZoneDropProposalAmino {
  title: string;
  description: string;
  zone_drop?: ZoneDropAmino;
  claim_records: Uint8Array;
}
export interface RegisterZoneDropProposalAminoMsg {
  type: "/quicksilver.airdrop.v1.RegisterZoneDropProposal";
  value: RegisterZoneDropProposalAmino;
}
export interface RegisterZoneDropProposalSDKType {
  title: string;
  description: string;
  zone_drop: ZoneDropSDKType;
  claim_records: Uint8Array;
}
function createBaseRegisterZoneDropProposal(): RegisterZoneDropProposal {
  return {
    title: "",
    description: "",
    zoneDrop: ZoneDrop.fromPartial({}),
    claimRecords: new Uint8Array()
  };
}
export const RegisterZoneDropProposal = {
  typeUrl: "/quicksilver.airdrop.v1.RegisterZoneDropProposal",
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
    const obj = createBaseRegisterZoneDropProposal();
    if (isSet(object.title)) obj.title = String(object.title);
    if (isSet(object.description)) obj.description = String(object.description);
    if (isSet(object.zoneDrop)) obj.zoneDrop = ZoneDrop.fromJSON(object.zoneDrop);
    if (isSet(object.claimRecords)) obj.claimRecords = bytesFromBase64(object.claimRecords);
    return obj;
  },
  toJSON(message: RegisterZoneDropProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.zoneDrop !== undefined && (obj.zoneDrop = message.zoneDrop ? ZoneDrop.toJSON(message.zoneDrop) : undefined);
    message.claimRecords !== undefined && (obj.claimRecords = base64FromBytes(message.claimRecords !== undefined ? message.claimRecords : new Uint8Array()));
    return obj;
  },
  fromPartial(object: DeepPartial<RegisterZoneDropProposal>): RegisterZoneDropProposal {
    const message = createBaseRegisterZoneDropProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    if (object.zoneDrop !== undefined && object.zoneDrop !== null) {
      message.zoneDrop = ZoneDrop.fromPartial(object.zoneDrop);
    }
    message.claimRecords = object.claimRecords ?? new Uint8Array();
    return message;
  },
  fromSDK(object: RegisterZoneDropProposalSDKType): RegisterZoneDropProposal {
    return {
      title: object?.title,
      description: object?.description,
      zoneDrop: object.zone_drop ? ZoneDrop.fromSDK(object.zone_drop) : undefined,
      claimRecords: object?.claim_records
    };
  },
  toSDK(message: RegisterZoneDropProposal): RegisterZoneDropProposalSDKType {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    message.zoneDrop !== undefined && (obj.zone_drop = message.zoneDrop ? ZoneDrop.toSDK(message.zoneDrop) : undefined);
    obj.claim_records = message.claimRecords;
    return obj;
  },
  fromAmino(object: RegisterZoneDropProposalAmino): RegisterZoneDropProposal {
    return {
      title: object.title,
      description: object.description,
      zoneDrop: object?.zone_drop ? ZoneDrop.fromAmino(object.zone_drop) : undefined,
      claimRecords: object.claim_records
    };
  },
  toAmino(message: RegisterZoneDropProposal): RegisterZoneDropProposalAmino {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.zone_drop = message.zoneDrop ? ZoneDrop.toAmino(message.zoneDrop) : undefined;
    obj.claim_records = message.claimRecords;
    return obj;
  },
  fromAminoMsg(object: RegisterZoneDropProposalAminoMsg): RegisterZoneDropProposal {
    return RegisterZoneDropProposal.fromAmino(object.value);
  },
  fromProtoMsg(message: RegisterZoneDropProposalProtoMsg): RegisterZoneDropProposal {
    return RegisterZoneDropProposal.decode(message.value);
  },
  toProto(message: RegisterZoneDropProposal): Uint8Array {
    return RegisterZoneDropProposal.encode(message).finish();
  },
  toProtoMsg(message: RegisterZoneDropProposal): RegisterZoneDropProposalProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.RegisterZoneDropProposal",
      value: RegisterZoneDropProposal.encode(message).finish()
    };
  }
};