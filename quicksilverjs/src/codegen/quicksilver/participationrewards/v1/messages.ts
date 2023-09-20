import { ClaimType, ClaimTypeSDKType, Proof, ProofAmino, ProofSDKType, claimTypeFromJSON, claimTypeToJSON } from "../../claimsmanager/v1/claimsmanager";
import * as _m0 from "protobufjs/minimal";
import { isSet, DeepPartial } from "../../../helpers";
export const protobufPackage = "quicksilver.participationrewards.v1";
/**
 * MsgSubmitClaim represents a message type for submitting a participation
 * claim regarding the given zone (chain).
 */
export interface MsgSubmitClaim {
  userAddress: string;
  zone: string;
  srcZone: string;
  claimType: ClaimType;
  proofs: Proof[];
}
export interface MsgSubmitClaimProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim";
  value: Uint8Array;
}
/**
 * MsgSubmitClaim represents a message type for submitting a participation
 * claim regarding the given zone (chain).
 */
export interface MsgSubmitClaimAmino {
  user_address: string;
  zone: string;
  src_zone: string;
  claim_type: ClaimType;
  proofs: ProofAmino[];
}
export interface MsgSubmitClaimAminoMsg {
  type: "/quicksilver.participationrewards.v1.MsgSubmitClaim";
  value: MsgSubmitClaimAmino;
}
/**
 * MsgSubmitClaim represents a message type for submitting a participation
 * claim regarding the given zone (chain).
 */
export interface MsgSubmitClaimSDKType {
  user_address: string;
  zone: string;
  src_zone: string;
  claim_type: ClaimType;
  proofs: ProofSDKType[];
}
/** MsgSubmitClaimResponse defines the MsgSubmitClaim response type. */
export interface MsgSubmitClaimResponse {}
export interface MsgSubmitClaimResponseProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaimResponse";
  value: Uint8Array;
}
/** MsgSubmitClaimResponse defines the MsgSubmitClaim response type. */
export interface MsgSubmitClaimResponseAmino {}
export interface MsgSubmitClaimResponseAminoMsg {
  type: "/quicksilver.participationrewards.v1.MsgSubmitClaimResponse";
  value: MsgSubmitClaimResponseAmino;
}
/** MsgSubmitClaimResponse defines the MsgSubmitClaim response type. */
export interface MsgSubmitClaimResponseSDKType {}
function createBaseMsgSubmitClaim(): MsgSubmitClaim {
  return {
    userAddress: "",
    zone: "",
    srcZone: "",
    claimType: 0,
    proofs: []
  };
}
export const MsgSubmitClaim = {
  typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
  encode(message: MsgSubmitClaim, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userAddress !== "") {
      writer.uint32(10).string(message.userAddress);
    }
    if (message.zone !== "") {
      writer.uint32(18).string(message.zone);
    }
    if (message.srcZone !== "") {
      writer.uint32(26).string(message.srcZone);
    }
    if (message.claimType !== 0) {
      writer.uint32(32).int32(message.claimType);
    }
    for (const v of message.proofs) {
      Proof.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitClaim {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSubmitClaim();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userAddress = reader.string();
          break;
        case 2:
          message.zone = reader.string();
          break;
        case 3:
          message.srcZone = reader.string();
          break;
        case 4:
          message.claimType = (reader.int32() as any);
          break;
        case 5:
          message.proofs.push(Proof.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgSubmitClaim {
    const obj = createBaseMsgSubmitClaim();
    if (isSet(object.user_address)) obj.userAddress = String(object.user_address);
    if (isSet(object.zone)) obj.zone = String(object.zone);
    if (isSet(object.src_zone)) obj.srcZone = String(object.src_zone);
    if (isSet(object.claim_type)) obj.claimType = claimTypeFromJSON(object.claim_type);
    if (Array.isArray(object?.proofs)) obj.proofs = object.proofs.map((e: any) => Proof.fromJSON(e));
    return obj;
  },
  toJSON(message: MsgSubmitClaim): unknown {
    const obj: any = {};
    message.userAddress !== undefined && (obj.user_address = message.userAddress);
    message.zone !== undefined && (obj.zone = message.zone);
    message.srcZone !== undefined && (obj.src_zone = message.srcZone);
    message.claimType !== undefined && (obj.claim_type = claimTypeToJSON(message.claimType));
    if (message.proofs) {
      obj.proofs = message.proofs.map(e => e ? Proof.toJSON(e) : undefined);
    } else {
      obj.proofs = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<MsgSubmitClaim>): MsgSubmitClaim {
    const message = createBaseMsgSubmitClaim();
    message.userAddress = object.userAddress ?? "";
    message.zone = object.zone ?? "";
    message.srcZone = object.srcZone ?? "";
    message.claimType = object.claimType ?? 0;
    message.proofs = object.proofs?.map(e => Proof.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: MsgSubmitClaimSDKType): MsgSubmitClaim {
    return {
      userAddress: object?.user_address,
      zone: object?.zone,
      srcZone: object?.src_zone,
      claimType: isSet(object.claim_type) ? claimTypeFromJSON(object.claim_type) : -1,
      proofs: Array.isArray(object?.proofs) ? object.proofs.map((e: any) => Proof.fromSDK(e)) : []
    };
  },
  toSDK(message: MsgSubmitClaim): MsgSubmitClaimSDKType {
    const obj: any = {};
    obj.user_address = message.userAddress;
    obj.zone = message.zone;
    obj.src_zone = message.srcZone;
    message.claimType !== undefined && (obj.claim_type = claimTypeToJSON(message.claimType));
    if (message.proofs) {
      obj.proofs = message.proofs.map(e => e ? Proof.toSDK(e) : undefined);
    } else {
      obj.proofs = [];
    }
    return obj;
  },
  fromAmino(object: MsgSubmitClaimAmino): MsgSubmitClaim {
    return {
      userAddress: object.user_address,
      zone: object.zone,
      srcZone: object.src_zone,
      claimType: isSet(object.claim_type) ? claimTypeFromJSON(object.claim_type) : -1,
      proofs: Array.isArray(object?.proofs) ? object.proofs.map((e: any) => Proof.fromAmino(e)) : []
    };
  },
  toAmino(message: MsgSubmitClaim): MsgSubmitClaimAmino {
    const obj: any = {};
    obj.user_address = message.userAddress;
    obj.zone = message.zone;
    obj.src_zone = message.srcZone;
    obj.claim_type = message.claimType;
    if (message.proofs) {
      obj.proofs = message.proofs.map(e => e ? Proof.toAmino(e) : undefined);
    } else {
      obj.proofs = [];
    }
    return obj;
  },
  fromAminoMsg(object: MsgSubmitClaimAminoMsg): MsgSubmitClaim {
    return MsgSubmitClaim.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSubmitClaimProtoMsg): MsgSubmitClaim {
    return MsgSubmitClaim.decode(message.value);
  },
  toProto(message: MsgSubmitClaim): Uint8Array {
    return MsgSubmitClaim.encode(message).finish();
  },
  toProtoMsg(message: MsgSubmitClaim): MsgSubmitClaimProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
      value: MsgSubmitClaim.encode(message).finish()
    };
  }
};
function createBaseMsgSubmitClaimResponse(): MsgSubmitClaimResponse {
  return {};
}
export const MsgSubmitClaimResponse = {
  typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaimResponse",
  encode(_: MsgSubmitClaimResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitClaimResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSubmitClaimResponse();
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
  fromJSON(_: any): MsgSubmitClaimResponse {
    const obj = createBaseMsgSubmitClaimResponse();
    return obj;
  },
  toJSON(_: MsgSubmitClaimResponse): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<MsgSubmitClaimResponse>): MsgSubmitClaimResponse {
    const message = createBaseMsgSubmitClaimResponse();
    return message;
  },
  fromSDK(_: MsgSubmitClaimResponseSDKType): MsgSubmitClaimResponse {
    return {};
  },
  toSDK(_: MsgSubmitClaimResponse): MsgSubmitClaimResponseSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: MsgSubmitClaimResponseAmino): MsgSubmitClaimResponse {
    return {};
  },
  toAmino(_: MsgSubmitClaimResponse): MsgSubmitClaimResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgSubmitClaimResponseAminoMsg): MsgSubmitClaimResponse {
    return MsgSubmitClaimResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSubmitClaimResponseProtoMsg): MsgSubmitClaimResponse {
    return MsgSubmitClaimResponse.decode(message.value);
  },
  toProto(message: MsgSubmitClaimResponse): Uint8Array {
    return MsgSubmitClaimResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSubmitClaimResponse): MsgSubmitClaimResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaimResponse",
      value: MsgSubmitClaimResponse.encode(message).finish()
    };
  }
};