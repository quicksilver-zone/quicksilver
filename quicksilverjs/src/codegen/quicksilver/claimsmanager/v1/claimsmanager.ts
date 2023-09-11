import * as _m0 from "protobufjs/minimal";
import { Long, isSet } from "../../../helpers";
export enum ClaimType {
  /** ClaimTypeUndefined - Undefined action (per protobuf spec) */
  ClaimTypeUndefined = 0,
  ClaimTypeLiquidToken = 1,
  ClaimTypeOsmosisPool = 2,
  ClaimTypeCrescentPool = 3,
  ClaimTypeSifchainPool = 4,
  UNRECOGNIZED = -1,
}
export enum ClaimTypeSDKType {
  /** ClaimTypeUndefined - Undefined action (per protobuf spec) */
  ClaimTypeUndefined = 0,
  ClaimTypeLiquidToken = 1,
  ClaimTypeOsmosisPool = 2,
  ClaimTypeCrescentPool = 3,
  ClaimTypeSifchainPool = 4,
  UNRECOGNIZED = -1,
}
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
/** Claim define the users claim for holding assets within a given zone. */

export interface ClaimSDKType {
  user_address: string;
  chain_id: string;
  module: ClaimTypeSDKType;
  source_chain_id: string;
  amount: Long;
}

function createBaseParams(): Params {
  return {};
}

export const Params = {
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
    return {};
  },

  toJSON(_: Params): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: Partial<Params>): Params {
    const message = createBaseParams();
    return message;
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
    return {
      userAddress: isSet(object.userAddress) ? String(object.userAddress) : "",
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      module: isSet(object.module) ? claimTypeFromJSON(object.module) : 0,
      sourceChainId: isSet(object.sourceChainId) ? String(object.sourceChainId) : "",
      amount: isSet(object.amount) ? Long.fromValue(object.amount) : Long.UZERO
    };
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

  fromPartial(object: Partial<Claim>): Claim {
    const message = createBaseClaim();
    message.userAddress = object.userAddress ?? "";
    message.chainId = object.chainId ?? "";
    message.module = object.module ?? 0;
    message.sourceChainId = object.sourceChainId ?? "";
    message.amount = object.amount !== undefined && object.amount !== null ? Long.fromValue(object.amount) : Long.UZERO;
    return message;
  }

};