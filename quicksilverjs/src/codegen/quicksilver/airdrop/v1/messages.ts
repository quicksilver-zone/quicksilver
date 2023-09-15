import { Proof, ProofAmino, ProofSDKType } from "../../claimsmanager/v1/claimsmanager";
import { Long, isSet, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.airdrop.v1";
export interface MsgClaim {
  chainId: string;
  action: Long;
  address: string;
  proofs: Proof[];
}
export interface MsgClaimProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.MsgClaim";
  value: Uint8Array;
}
export interface MsgClaimAmino {
  chain_id: string;
  action: string;
  address: string;
  proofs: ProofAmino[];
}
export interface MsgClaimAminoMsg {
  type: "/quicksilver.airdrop.v1.MsgClaim";
  value: MsgClaimAmino;
}
export interface MsgClaimSDKType {
  chain_id: string;
  action: Long;
  address: string;
  proofs: ProofSDKType[];
}
export interface MsgClaimResponse {
  amount: Long;
}
export interface MsgClaimResponseProtoMsg {
  typeUrl: "/quicksilver.airdrop.v1.MsgClaimResponse";
  value: Uint8Array;
}
export interface MsgClaimResponseAmino {
  amount: string;
}
export interface MsgClaimResponseAminoMsg {
  type: "/quicksilver.airdrop.v1.MsgClaimResponse";
  value: MsgClaimResponseAmino;
}
export interface MsgClaimResponseSDKType {
  amount: Long;
}
function createBaseMsgClaim(): MsgClaim {
  return {
    chainId: "",
    action: Long.ZERO,
    address: "",
    proofs: []
  };
}
export const MsgClaim = {
  typeUrl: "/quicksilver.airdrop.v1.MsgClaim",
  encode(message: MsgClaim, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (!message.action.isZero()) {
      writer.uint32(16).int64(message.action);
    }
    if (message.address !== "") {
      writer.uint32(26).string(message.address);
    }
    for (const v of message.proofs) {
      Proof.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaim {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaim();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.action = (reader.int64() as Long);
          break;
        case 3:
          message.address = reader.string();
          break;
        case 4:
          message.proofs.push(Proof.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgClaim {
    const obj = createBaseMsgClaim();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.action)) obj.action = Long.fromValue(object.action);
    if (isSet(object.address)) obj.address = String(object.address);
    if (Array.isArray(object?.proofs)) obj.proofs = object.proofs.map((e: any) => Proof.fromJSON(e));
    return obj;
  },
  toJSON(message: MsgClaim): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.action !== undefined && (obj.action = (message.action || Long.ZERO).toString());
    message.address !== undefined && (obj.address = message.address);
    if (message.proofs) {
      obj.proofs = message.proofs.map(e => e ? Proof.toJSON(e) : undefined);
    } else {
      obj.proofs = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<MsgClaim>): MsgClaim {
    const message = createBaseMsgClaim();
    message.chainId = object.chainId ?? "";
    if (object.action !== undefined && object.action !== null) {
      message.action = Long.fromValue(object.action);
    }
    message.address = object.address ?? "";
    message.proofs = object.proofs?.map(e => Proof.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: MsgClaimSDKType): MsgClaim {
    return {
      chainId: object?.chain_id,
      action: object?.action,
      address: object?.address,
      proofs: Array.isArray(object?.proofs) ? object.proofs.map((e: any) => Proof.fromSDK(e)) : []
    };
  },
  toSDK(message: MsgClaim): MsgClaimSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.action = message.action;
    obj.address = message.address;
    if (message.proofs) {
      obj.proofs = message.proofs.map(e => e ? Proof.toSDK(e) : undefined);
    } else {
      obj.proofs = [];
    }
    return obj;
  },
  fromAmino(object: MsgClaimAmino): MsgClaim {
    return {
      chainId: object.chain_id,
      action: Long.fromString(object.action),
      address: object.address,
      proofs: Array.isArray(object?.proofs) ? object.proofs.map((e: any) => Proof.fromAmino(e)) : []
    };
  },
  toAmino(message: MsgClaim): MsgClaimAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.action = message.action ? message.action.toString() : undefined;
    obj.address = message.address;
    if (message.proofs) {
      obj.proofs = message.proofs.map(e => e ? Proof.toAmino(e) : undefined);
    } else {
      obj.proofs = [];
    }
    return obj;
  },
  fromAminoMsg(object: MsgClaimAminoMsg): MsgClaim {
    return MsgClaim.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgClaimProtoMsg): MsgClaim {
    return MsgClaim.decode(message.value);
  },
  toProto(message: MsgClaim): Uint8Array {
    return MsgClaim.encode(message).finish();
  },
  toProtoMsg(message: MsgClaim): MsgClaimProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.MsgClaim",
      value: MsgClaim.encode(message).finish()
    };
  }
};
function createBaseMsgClaimResponse(): MsgClaimResponse {
  return {
    amount: Long.UZERO
  };
}
export const MsgClaimResponse = {
  typeUrl: "/quicksilver.airdrop.v1.MsgClaimResponse",
  encode(message: MsgClaimResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.amount.isZero()) {
      writer.uint32(8).uint64(message.amount);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.amount = (reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgClaimResponse {
    const obj = createBaseMsgClaimResponse();
    if (isSet(object.amount)) obj.amount = Long.fromValue(object.amount);
    return obj;
  },
  toJSON(message: MsgClaimResponse): unknown {
    const obj: any = {};
    message.amount !== undefined && (obj.amount = (message.amount || Long.UZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<MsgClaimResponse>): MsgClaimResponse {
    const message = createBaseMsgClaimResponse();
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Long.fromValue(object.amount);
    }
    return message;
  },
  fromSDK(object: MsgClaimResponseSDKType): MsgClaimResponse {
    return {
      amount: object?.amount
    };
  },
  toSDK(message: MsgClaimResponse): MsgClaimResponseSDKType {
    const obj: any = {};
    obj.amount = message.amount;
    return obj;
  },
  fromAmino(object: MsgClaimResponseAmino): MsgClaimResponse {
    return {
      amount: Long.fromString(object.amount)
    };
  },
  toAmino(message: MsgClaimResponse): MsgClaimResponseAmino {
    const obj: any = {};
    obj.amount = message.amount ? message.amount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgClaimResponseAminoMsg): MsgClaimResponse {
    return MsgClaimResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgClaimResponseProtoMsg): MsgClaimResponse {
    return MsgClaimResponse.decode(message.value);
  },
  toProto(message: MsgClaimResponse): Uint8Array {
    return MsgClaimResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgClaimResponse): MsgClaimResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.airdrop.v1.MsgClaimResponse",
      value: MsgClaimResponse.encode(message).finish()
    };
  }
};