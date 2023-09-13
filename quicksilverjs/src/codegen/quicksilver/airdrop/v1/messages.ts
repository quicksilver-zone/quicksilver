import { ProofOps, ProofOpsSDKType } from "../../../tendermint/crypto/proof";
import * as _m0 from "protobufjs/minimal";
import { Long, isSet, bytesFromBase64, base64FromBytes } from "../../../helpers";
export interface MsgClaim {
  chainId: string;
  action: Long;
  address: string;
  proofs: Proof[];
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
export interface MsgClaimResponseSDKType {
  amount: Long;
}
export interface Proof {
  key: Uint8Array;
  data: Uint8Array;
  proofOps?: ProofOps;
  height: Long;
}
export interface ProofSDKType {
  key: Uint8Array;
  data: Uint8Array;
  proof_ops?: ProofOpsSDKType;
  height: Long;
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
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      action: isSet(object.action) ? Long.fromValue(object.action) : Long.ZERO,
      address: isSet(object.address) ? String(object.address) : "",
      proofs: Array.isArray(object?.proofs) ? object.proofs.map((e: any) => Proof.fromJSON(e)) : []
    };
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

  fromPartial(object: Partial<MsgClaim>): MsgClaim {
    const message = createBaseMsgClaim();
    message.chainId = object.chainId ?? "";
    message.action = object.action !== undefined && object.action !== null ? Long.fromValue(object.action) : Long.ZERO;
    message.address = object.address ?? "";
    message.proofs = object.proofs?.map(e => Proof.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMsgClaimResponse(): MsgClaimResponse {
  return {
    amount: Long.UZERO
  };
}

export const MsgClaimResponse = {
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
    return {
      amount: isSet(object.amount) ? Long.fromValue(object.amount) : Long.UZERO
    };
  },

  toJSON(message: MsgClaimResponse): unknown {
    const obj: any = {};
    message.amount !== undefined && (obj.amount = (message.amount || Long.UZERO).toString());
    return obj;
  },

  fromPartial(object: Partial<MsgClaimResponse>): MsgClaimResponse {
    const message = createBaseMsgClaimResponse();
    message.amount = object.amount !== undefined && object.amount !== null ? Long.fromValue(object.amount) : Long.UZERO;
    return message;
  }

};

function createBaseProof(): Proof {
  return {
    key: new Uint8Array(),
    data: new Uint8Array(),
    proofOps: undefined,
    height: Long.ZERO
  };
}

export const Proof = {
  encode(message: Proof, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key.length !== 0) {
      writer.uint32(10).bytes(message.key);
    }

    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }

    if (message.proofOps !== undefined) {
      ProofOps.encode(message.proofOps, writer.uint32(26).fork()).ldelim();
    }

    if (!message.height.isZero()) {
      writer.uint32(32).int64(message.height);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Proof {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProof();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.key = reader.bytes();
          break;

        case 2:
          message.data = reader.bytes();
          break;

        case 3:
          message.proofOps = ProofOps.decode(reader, reader.uint32());
          break;

        case 4:
          message.height = (reader.int64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): Proof {
    return {
      key: isSet(object.key) ? bytesFromBase64(object.key) : new Uint8Array(),
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(),
      proofOps: isSet(object.proofOps) ? ProofOps.fromJSON(object.proofOps) : undefined,
      height: isSet(object.height) ? Long.fromValue(object.height) : Long.ZERO
    };
  },

  toJSON(message: Proof): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = base64FromBytes(message.key !== undefined ? message.key : new Uint8Array()));
    message.data !== undefined && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    message.proofOps !== undefined && (obj.proofOps = message.proofOps ? ProofOps.toJSON(message.proofOps) : undefined);
    message.height !== undefined && (obj.height = (message.height || Long.ZERO).toString());
    return obj;
  },

  fromPartial(object: Partial<Proof>): Proof {
    const message = createBaseProof();
    message.key = object.key ?? new Uint8Array();
    message.data = object.data ?? new Uint8Array();
    message.proofOps = object.proofOps !== undefined && object.proofOps !== null ? ProofOps.fromPartial(object.proofOps) : undefined;
    message.height = object.height !== undefined && object.height !== null ? Long.fromValue(object.height) : Long.ZERO;
    return message;
  }

};