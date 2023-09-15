import { ProofOps, ProofOpsAmino, ProofOpsSDKType } from "../../../tendermint/crypto/proof";
import { Long, DeepPartial, isSet, bytesFromBase64, base64FromBytes } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.claimsmanager.v1";
export enum ClaimType {
  /** ClaimTypeUndefined - Undefined action (per protobuf spec) */
  ClaimTypeUndefined = 0,
  ClaimTypeLiquidToken = 1,
  ClaimTypeOsmosisPool = 2,
  ClaimTypeCrescentPool = 3,
  ClaimTypeSifchainPool = 4,
  UNRECOGNIZED = -1,
}
export const ClaimTypeSDKType = ClaimType;
export const ClaimTypeAmino = ClaimType;
export function claimTypeFromJSON(object: any): ClaimType {
  switch (object) {
    case 0:
    case "ClaimTypeUndefined":
      return ClaimType.ClaimTypeUndefined;
    case 1:
    case "ClaimTypeLiquidToken":
      return ClaimType.ClaimTypeLiquidToken;
    case 2:
    case "ClaimTypeOsmosisPool":
      return ClaimType.ClaimTypeOsmosisPool;
    case 3:
    case "ClaimTypeCrescentPool":
      return ClaimType.ClaimTypeCrescentPool;
    case 4:
    case "ClaimTypeSifchainPool":
      return ClaimType.ClaimTypeSifchainPool;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ClaimType.UNRECOGNIZED;
  }
}
export function claimTypeToJSON(object: ClaimType): string {
  switch (object) {
    case ClaimType.ClaimTypeUndefined:
      return "ClaimTypeUndefined";
    case ClaimType.ClaimTypeLiquidToken:
      return "ClaimTypeLiquidToken";
    case ClaimType.ClaimTypeOsmosisPool:
      return "ClaimTypeOsmosisPool";
    case ClaimType.ClaimTypeCrescentPool:
      return "ClaimTypeCrescentPool";
    case ClaimType.ClaimTypeSifchainPool:
      return "ClaimTypeSifchainPool";
    case ClaimType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** Params holds parameters for the claimsmanager module. */
export interface Params {}
export interface ParamsProtoMsg {
  typeUrl: "/quicksilver.claimsmanager.v1.Params";
  value: Uint8Array;
}
/** Params holds parameters for the claimsmanager module. */
export interface ParamsAmino {}
export interface ParamsAminoMsg {
  type: "/quicksilver.claimsmanager.v1.Params";
  value: ParamsAmino;
}
/** Params holds parameters for the claimsmanager module. */
export interface ParamsSDKType {}
/** Claim define the users claim for holding assets within a given zone. */
export interface Claim {
  userAddress: string;
  chainId: string;
  module: ClaimType;
  sourceChainId: string;
  amount: Long;
}
export interface ClaimProtoMsg {
  typeUrl: "/quicksilver.claimsmanager.v1.Claim";
  value: Uint8Array;
}
/** Claim define the users claim for holding assets within a given zone. */
export interface ClaimAmino {
  user_address: string;
  chain_id: string;
  module: ClaimType;
  source_chain_id: string;
  amount: string;
}
export interface ClaimAminoMsg {
  type: "/quicksilver.claimsmanager.v1.Claim";
  value: ClaimAmino;
}
/** Claim define the users claim for holding assets within a given zone. */
export interface ClaimSDKType {
  user_address: string;
  chain_id: string;
  module: ClaimType;
  source_chain_id: string;
  amount: Long;
}
/** Proof defines a type used to cryptographically prove a claim. */
export interface Proof {
  key: Uint8Array;
  data: Uint8Array;
  proofOps: ProofOps;
  height: Long;
  proofType: string;
}
export interface ProofProtoMsg {
  typeUrl: "/quicksilver.claimsmanager.v1.Proof";
  value: Uint8Array;
}
/** Proof defines a type used to cryptographically prove a claim. */
export interface ProofAmino {
  key: Uint8Array;
  data: Uint8Array;
  proof_ops?: ProofOpsAmino;
  height: string;
  proof_type: string;
}
export interface ProofAminoMsg {
  type: "/quicksilver.claimsmanager.v1.Proof";
  value: ProofAmino;
}
/** Proof defines a type used to cryptographically prove a claim. */
export interface ProofSDKType {
  key: Uint8Array;
  data: Uint8Array;
  proof_ops: ProofOpsSDKType;
  height: Long;
  proof_type: string;
}
function createBaseParams(): Params {
  return {};
}
export const Params = {
  typeUrl: "/quicksilver.claimsmanager.v1.Params",
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
      typeUrl: "/quicksilver.claimsmanager.v1.Params",
      value: Params.encode(message).finish()
    };
  }
};
function createBaseClaim(): Claim {
  return {
    userAddress: "",
    chainId: "",
    module: 0,
    sourceChainId: "",
    amount: Long.UZERO
  };
}
export const Claim = {
  typeUrl: "/quicksilver.claimsmanager.v1.Claim",
  encode(message: Claim, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userAddress !== "") {
      writer.uint32(10).string(message.userAddress);
    }
    if (message.chainId !== "") {
      writer.uint32(18).string(message.chainId);
    }
    if (message.module !== 0) {
      writer.uint32(24).int32(message.module);
    }
    if (message.sourceChainId !== "") {
      writer.uint32(34).string(message.sourceChainId);
    }
    if (!message.amount.isZero()) {
      writer.uint32(40).uint64(message.amount);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Claim {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClaim();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userAddress = reader.string();
          break;
        case 2:
          message.chainId = reader.string();
          break;
        case 3:
          message.module = (reader.int32() as any);
          break;
        case 4:
          message.sourceChainId = reader.string();
          break;
        case 5:
          message.amount = (reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Claim {
    const obj = createBaseClaim();
    if (isSet(object.userAddress)) obj.userAddress = String(object.userAddress);
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.module)) obj.module = claimTypeFromJSON(object.module);
    if (isSet(object.sourceChainId)) obj.sourceChainId = String(object.sourceChainId);
    if (isSet(object.amount)) obj.amount = Long.fromValue(object.amount);
    return obj;
  },
  toJSON(message: Claim): unknown {
    const obj: any = {};
    message.userAddress !== undefined && (obj.userAddress = message.userAddress);
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.module !== undefined && (obj.module = claimTypeToJSON(message.module));
    message.sourceChainId !== undefined && (obj.sourceChainId = message.sourceChainId);
    message.amount !== undefined && (obj.amount = (message.amount || Long.UZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<Claim>): Claim {
    const message = createBaseClaim();
    message.userAddress = object.userAddress ?? "";
    message.chainId = object.chainId ?? "";
    message.module = object.module ?? 0;
    message.sourceChainId = object.sourceChainId ?? "";
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Long.fromValue(object.amount);
    }
    return message;
  },
  fromSDK(object: ClaimSDKType): Claim {
    return {
      userAddress: object?.user_address,
      chainId: object?.chain_id,
      module: isSet(object.module) ? claimTypeFromJSON(object.module) : -1,
      sourceChainId: object?.source_chain_id,
      amount: object?.amount
    };
  },
  toSDK(message: Claim): ClaimSDKType {
    const obj: any = {};
    obj.user_address = message.userAddress;
    obj.chain_id = message.chainId;
    message.module !== undefined && (obj.module = claimTypeToJSON(message.module));
    obj.source_chain_id = message.sourceChainId;
    obj.amount = message.amount;
    return obj;
  },
  fromAmino(object: ClaimAmino): Claim {
    return {
      userAddress: object.user_address,
      chainId: object.chain_id,
      module: isSet(object.module) ? claimTypeFromJSON(object.module) : -1,
      sourceChainId: object.source_chain_id,
      amount: Long.fromString(object.amount)
    };
  },
  toAmino(message: Claim): ClaimAmino {
    const obj: any = {};
    obj.user_address = message.userAddress;
    obj.chain_id = message.chainId;
    obj.module = message.module;
    obj.source_chain_id = message.sourceChainId;
    obj.amount = message.amount ? message.amount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: ClaimAminoMsg): Claim {
    return Claim.fromAmino(object.value);
  },
  fromProtoMsg(message: ClaimProtoMsg): Claim {
    return Claim.decode(message.value);
  },
  toProto(message: Claim): Uint8Array {
    return Claim.encode(message).finish();
  },
  toProtoMsg(message: Claim): ClaimProtoMsg {
    return {
      typeUrl: "/quicksilver.claimsmanager.v1.Claim",
      value: Claim.encode(message).finish()
    };
  }
};
function createBaseProof(): Proof {
  return {
    key: new Uint8Array(),
    data: new Uint8Array(),
    proofOps: ProofOps.fromPartial({}),
    height: Long.ZERO,
    proofType: ""
  };
}
export const Proof = {
  typeUrl: "/quicksilver.claimsmanager.v1.Proof",
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
    if (message.proofType !== "") {
      writer.uint32(42).string(message.proofType);
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
        case 5:
          message.proofType = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Proof {
    const obj = createBaseProof();
    if (isSet(object.key)) obj.key = bytesFromBase64(object.key);
    if (isSet(object.data)) obj.data = bytesFromBase64(object.data);
    if (isSet(object.proof_ops)) obj.proofOps = ProofOps.fromJSON(object.proof_ops);
    if (isSet(object.height)) obj.height = Long.fromValue(object.height);
    if (isSet(object.proof_type)) obj.proofType = String(object.proof_type);
    return obj;
  },
  toJSON(message: Proof): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = base64FromBytes(message.key !== undefined ? message.key : new Uint8Array()));
    message.data !== undefined && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    message.proofOps !== undefined && (obj.proof_ops = message.proofOps ? ProofOps.toJSON(message.proofOps) : undefined);
    message.height !== undefined && (obj.height = (message.height || Long.ZERO).toString());
    message.proofType !== undefined && (obj.proof_type = message.proofType);
    return obj;
  },
  fromPartial(object: DeepPartial<Proof>): Proof {
    const message = createBaseProof();
    message.key = object.key ?? new Uint8Array();
    message.data = object.data ?? new Uint8Array();
    if (object.proofOps !== undefined && object.proofOps !== null) {
      message.proofOps = ProofOps.fromPartial(object.proofOps);
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = Long.fromValue(object.height);
    }
    message.proofType = object.proofType ?? "";
    return message;
  },
  fromSDK(object: ProofSDKType): Proof {
    return {
      key: object?.key,
      data: object?.data,
      proofOps: object.proof_ops ? ProofOps.fromSDK(object.proof_ops) : undefined,
      height: object?.height,
      proofType: object?.proof_type
    };
  },
  toSDK(message: Proof): ProofSDKType {
    const obj: any = {};
    obj.key = message.key;
    obj.data = message.data;
    message.proofOps !== undefined && (obj.proof_ops = message.proofOps ? ProofOps.toSDK(message.proofOps) : undefined);
    obj.height = message.height;
    obj.proof_type = message.proofType;
    return obj;
  },
  fromAmino(object: ProofAmino): Proof {
    return {
      key: object.key,
      data: object.data,
      proofOps: object?.proof_ops ? ProofOps.fromAmino(object.proof_ops) : undefined,
      height: Long.fromString(object.height),
      proofType: object.proof_type
    };
  },
  toAmino(message: Proof): ProofAmino {
    const obj: any = {};
    obj.key = message.key;
    obj.data = message.data;
    obj.proof_ops = message.proofOps ? ProofOps.toAmino(message.proofOps) : undefined;
    obj.height = message.height ? message.height.toString() : undefined;
    obj.proof_type = message.proofType;
    return obj;
  },
  fromAminoMsg(object: ProofAminoMsg): Proof {
    return Proof.fromAmino(object.value);
  },
  fromProtoMsg(message: ProofProtoMsg): Proof {
    return Proof.decode(message.value);
  },
  toProto(message: Proof): Uint8Array {
    return Proof.encode(message).finish();
  },
  toProtoMsg(message: Proof): ProofProtoMsg {
    return {
      typeUrl: "/quicksilver.claimsmanager.v1.Proof",
      value: Proof.encode(message).finish()
    };
  }
};