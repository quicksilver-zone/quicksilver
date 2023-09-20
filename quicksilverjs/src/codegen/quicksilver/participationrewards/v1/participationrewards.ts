import * as _m0 from "protobufjs/minimal";
import { isSet, DeepPartial, bytesFromBase64, base64FromBytes } from "../../../helpers";
export const protobufPackage = "quicksilver.participationrewards.v1";
export enum ProtocolDataType {
  /** ProtocolDataTypeUndefined - Undefined action (per protobuf spec) */
  ProtocolDataTypeUndefined = 0,
  ProtocolDataTypeConnection = 1,
  ProtocolDataTypeOsmosisParams = 2,
  ProtocolDataTypeLiquidToken = 3,
  ProtocolDataTypeOsmosisPool = 4,
  ProtocolDataTypeCrescentPool = 5,
  ProtocolDataTypeSifchainPool = 6,
  UNRECOGNIZED = -1,
}
export const ProtocolDataTypeSDKType = ProtocolDataType;
export const ProtocolDataTypeAmino = ProtocolDataType;
export function protocolDataTypeFromJSON(object: any): ProtocolDataType {
  switch (object) {
    case 0:
    case "ProtocolDataTypeUndefined":
      return ProtocolDataType.ProtocolDataTypeUndefined;
    case 1:
    case "ProtocolDataTypeConnection":
      return ProtocolDataType.ProtocolDataTypeConnection;
    case 2:
    case "ProtocolDataTypeOsmosisParams":
      return ProtocolDataType.ProtocolDataTypeOsmosisParams;
    case 3:
    case "ProtocolDataTypeLiquidToken":
      return ProtocolDataType.ProtocolDataTypeLiquidToken;
    case 4:
    case "ProtocolDataTypeOsmosisPool":
      return ProtocolDataType.ProtocolDataTypeOsmosisPool;
    case 5:
    case "ProtocolDataTypeCrescentPool":
      return ProtocolDataType.ProtocolDataTypeCrescentPool;
    case 6:
    case "ProtocolDataTypeSifchainPool":
      return ProtocolDataType.ProtocolDataTypeSifchainPool;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ProtocolDataType.UNRECOGNIZED;
  }
}
export function protocolDataTypeToJSON(object: ProtocolDataType): string {
  switch (object) {
    case ProtocolDataType.ProtocolDataTypeUndefined:
      return "ProtocolDataTypeUndefined";
    case ProtocolDataType.ProtocolDataTypeConnection:
      return "ProtocolDataTypeConnection";
    case ProtocolDataType.ProtocolDataTypeOsmosisParams:
      return "ProtocolDataTypeOsmosisParams";
    case ProtocolDataType.ProtocolDataTypeLiquidToken:
      return "ProtocolDataTypeLiquidToken";
    case ProtocolDataType.ProtocolDataTypeOsmosisPool:
      return "ProtocolDataTypeOsmosisPool";
    case ProtocolDataType.ProtocolDataTypeCrescentPool:
      return "ProtocolDataTypeCrescentPool";
    case ProtocolDataType.ProtocolDataTypeSifchainPool:
      return "ProtocolDataTypeSifchainPool";
    case ProtocolDataType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * DistributionProportions defines the proportions of minted QCK that is to be
 * allocated as participation rewards.
 */
export interface DistributionProportions {
  validatorSelectionAllocation: string;
  holdingsAllocation: string;
  lockupAllocation: string;
}
export interface DistributionProportionsProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.DistributionProportions";
  value: Uint8Array;
}
/**
 * DistributionProportions defines the proportions of minted QCK that is to be
 * allocated as participation rewards.
 */
export interface DistributionProportionsAmino {
  validator_selection_allocation: string;
  holdings_allocation: string;
  lockup_allocation: string;
}
export interface DistributionProportionsAminoMsg {
  type: "/quicksilver.participationrewards.v1.DistributionProportions";
  value: DistributionProportionsAmino;
}
/**
 * DistributionProportions defines the proportions of minted QCK that is to be
 * allocated as participation rewards.
 */
export interface DistributionProportionsSDKType {
  validator_selection_allocation: string;
  holdings_allocation: string;
  lockup_allocation: string;
}
/** Params holds parameters for the participationrewards module. */
export interface Params_v1 {
  /**
   * distribution_proportions defines the proportions of the minted
   * participation rewards;
   */
  distributionProportions: DistributionProportions;
}
export interface Params_v1ProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.Params_v1";
  value: Uint8Array;
}
/** Params holds parameters for the participationrewards module. */
export interface Params_v1Amino {
  /**
   * distribution_proportions defines the proportions of the minted
   * participation rewards;
   */
  distribution_proportions?: DistributionProportionsAmino;
}
export interface Params_v1AminoMsg {
  type: "/quicksilver.participationrewards.v1.Params_v1";
  value: Params_v1Amino;
}
/** Params holds parameters for the participationrewards module. */
export interface Params_v1SDKType {
  distribution_proportions: DistributionProportionsSDKType;
}
/** Params holds parameters for the participationrewards module. */
export interface Params {
  /**
   * distribution_proportions defines the proportions of the minted
   * participation rewards;
   */
  distributionProportions: DistributionProportions;
  claimsEnabled: boolean;
}
export interface ParamsProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.Params";
  value: Uint8Array;
}
/** Params holds parameters for the participationrewards module. */
export interface ParamsAmino {
  /**
   * distribution_proportions defines the proportions of the minted
   * participation rewards;
   */
  distribution_proportions?: DistributionProportionsAmino;
  claims_enabled: boolean;
}
export interface ParamsAminoMsg {
  type: "/quicksilver.participationrewards.v1.Params";
  value: ParamsAmino;
}
/** Params holds parameters for the participationrewards module. */
export interface ParamsSDKType {
  distribution_proportions: DistributionProportionsSDKType;
  claims_enabled: boolean;
}
export interface KeyedProtocolData {
  key: string;
  protocolData: ProtocolData;
}
export interface KeyedProtocolDataProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.KeyedProtocolData";
  value: Uint8Array;
}
export interface KeyedProtocolDataAmino {
  key: string;
  protocol_data?: ProtocolDataAmino;
}
export interface KeyedProtocolDataAminoMsg {
  type: "/quicksilver.participationrewards.v1.KeyedProtocolData";
  value: KeyedProtocolDataAmino;
}
export interface KeyedProtocolDataSDKType {
  key: string;
  protocol_data: ProtocolDataSDKType;
}
/**
 * Protocol Data is an arbitrary data type held against a given zone for the
 * determination of rewards.
 */
export interface ProtocolData {
  type: string;
  data: Uint8Array;
}
export interface ProtocolDataProtoMsg {
  typeUrl: "/quicksilver.participationrewards.v1.ProtocolData";
  value: Uint8Array;
}
/**
 * Protocol Data is an arbitrary data type held against a given zone for the
 * determination of rewards.
 */
export interface ProtocolDataAmino {
  type: string;
  data: Uint8Array;
}
export interface ProtocolDataAminoMsg {
  type: "/quicksilver.participationrewards.v1.ProtocolData";
  value: ProtocolDataAmino;
}
/**
 * Protocol Data is an arbitrary data type held against a given zone for the
 * determination of rewards.
 */
export interface ProtocolDataSDKType {
  type: string;
  data: Uint8Array;
}
function createBaseDistributionProportions(): DistributionProportions {
  return {
    validatorSelectionAllocation: "",
    holdingsAllocation: "",
    lockupAllocation: ""
  };
}
export const DistributionProportions = {
  typeUrl: "/quicksilver.participationrewards.v1.DistributionProportions",
  encode(message: DistributionProportions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.validatorSelectionAllocation !== "") {
      writer.uint32(10).string(message.validatorSelectionAllocation);
    }
    if (message.holdingsAllocation !== "") {
      writer.uint32(18).string(message.holdingsAllocation);
    }
    if (message.lockupAllocation !== "") {
      writer.uint32(26).string(message.lockupAllocation);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): DistributionProportions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDistributionProportions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.validatorSelectionAllocation = reader.string();
          break;
        case 2:
          message.holdingsAllocation = reader.string();
          break;
        case 3:
          message.lockupAllocation = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): DistributionProportions {
    const obj = createBaseDistributionProportions();
    if (isSet(object.validatorSelectionAllocation)) obj.validatorSelectionAllocation = String(object.validatorSelectionAllocation);
    if (isSet(object.holdingsAllocation)) obj.holdingsAllocation = String(object.holdingsAllocation);
    if (isSet(object.lockupAllocation)) obj.lockupAllocation = String(object.lockupAllocation);
    return obj;
  },
  toJSON(message: DistributionProportions): unknown {
    const obj: any = {};
    message.validatorSelectionAllocation !== undefined && (obj.validatorSelectionAllocation = message.validatorSelectionAllocation);
    message.holdingsAllocation !== undefined && (obj.holdingsAllocation = message.holdingsAllocation);
    message.lockupAllocation !== undefined && (obj.lockupAllocation = message.lockupAllocation);
    return obj;
  },
  fromPartial(object: DeepPartial<DistributionProportions>): DistributionProportions {
    const message = createBaseDistributionProportions();
    message.validatorSelectionAllocation = object.validatorSelectionAllocation ?? "";
    message.holdingsAllocation = object.holdingsAllocation ?? "";
    message.lockupAllocation = object.lockupAllocation ?? "";
    return message;
  },
  fromSDK(object: DistributionProportionsSDKType): DistributionProportions {
    return {
      validatorSelectionAllocation: object?.validator_selection_allocation,
      holdingsAllocation: object?.holdings_allocation,
      lockupAllocation: object?.lockup_allocation
    };
  },
  toSDK(message: DistributionProportions): DistributionProportionsSDKType {
    const obj: any = {};
    obj.validator_selection_allocation = message.validatorSelectionAllocation;
    obj.holdings_allocation = message.holdingsAllocation;
    obj.lockup_allocation = message.lockupAllocation;
    return obj;
  },
  fromAmino(object: DistributionProportionsAmino): DistributionProportions {
    return {
      validatorSelectionAllocation: object.validator_selection_allocation,
      holdingsAllocation: object.holdings_allocation,
      lockupAllocation: object.lockup_allocation
    };
  },
  toAmino(message: DistributionProportions): DistributionProportionsAmino {
    const obj: any = {};
    obj.validator_selection_allocation = message.validatorSelectionAllocation;
    obj.holdings_allocation = message.holdingsAllocation;
    obj.lockup_allocation = message.lockupAllocation;
    return obj;
  },
  fromAminoMsg(object: DistributionProportionsAminoMsg): DistributionProportions {
    return DistributionProportions.fromAmino(object.value);
  },
  fromProtoMsg(message: DistributionProportionsProtoMsg): DistributionProportions {
    return DistributionProportions.decode(message.value);
  },
  toProto(message: DistributionProportions): Uint8Array {
    return DistributionProportions.encode(message).finish();
  },
  toProtoMsg(message: DistributionProportions): DistributionProportionsProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.DistributionProportions",
      value: DistributionProportions.encode(message).finish()
    };
  }
};
function createBaseParams_v1(): Params_v1 {
  return {
    distributionProportions: DistributionProportions.fromPartial({})
  };
}
export const Params_v1 = {
  typeUrl: "/quicksilver.participationrewards.v1.Params_v1",
  encode(message: Params_v1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.distributionProportions !== undefined) {
      DistributionProportions.encode(message.distributionProportions, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Params_v1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams_v1();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.distributionProportions = DistributionProportions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Params_v1 {
    const obj = createBaseParams_v1();
    if (isSet(object.distributionProportions)) obj.distributionProportions = DistributionProportions.fromJSON(object.distributionProportions);
    return obj;
  },
  toJSON(message: Params_v1): unknown {
    const obj: any = {};
    message.distributionProportions !== undefined && (obj.distributionProportions = message.distributionProportions ? DistributionProportions.toJSON(message.distributionProportions) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<Params_v1>): Params_v1 {
    const message = createBaseParams_v1();
    if (object.distributionProportions !== undefined && object.distributionProportions !== null) {
      message.distributionProportions = DistributionProportions.fromPartial(object.distributionProportions);
    }
    return message;
  },
  fromSDK(object: Params_v1SDKType): Params_v1 {
    return {
      distributionProportions: object.distribution_proportions ? DistributionProportions.fromSDK(object.distribution_proportions) : undefined
    };
  },
  toSDK(message: Params_v1): Params_v1SDKType {
    const obj: any = {};
    message.distributionProportions !== undefined && (obj.distribution_proportions = message.distributionProportions ? DistributionProportions.toSDK(message.distributionProportions) : undefined);
    return obj;
  },
  fromAmino(object: Params_v1Amino): Params_v1 {
    return {
      distributionProportions: object?.distribution_proportions ? DistributionProportions.fromAmino(object.distribution_proportions) : undefined
    };
  },
  toAmino(message: Params_v1): Params_v1Amino {
    const obj: any = {};
    obj.distribution_proportions = message.distributionProportions ? DistributionProportions.toAmino(message.distributionProportions) : undefined;
    return obj;
  },
  fromAminoMsg(object: Params_v1AminoMsg): Params_v1 {
    return Params_v1.fromAmino(object.value);
  },
  fromProtoMsg(message: Params_v1ProtoMsg): Params_v1 {
    return Params_v1.decode(message.value);
  },
  toProto(message: Params_v1): Uint8Array {
    return Params_v1.encode(message).finish();
  },
  toProtoMsg(message: Params_v1): Params_v1ProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.Params_v1",
      value: Params_v1.encode(message).finish()
    };
  }
};
function createBaseParams(): Params {
  return {
    distributionProportions: DistributionProportions.fromPartial({}),
    claimsEnabled: false
  };
}
export const Params = {
  typeUrl: "/quicksilver.participationrewards.v1.Params",
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.distributionProportions !== undefined) {
      DistributionProportions.encode(message.distributionProportions, writer.uint32(10).fork()).ldelim();
    }
    if (message.claimsEnabled === true) {
      writer.uint32(16).bool(message.claimsEnabled);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.distributionProportions = DistributionProportions.decode(reader, reader.uint32());
          break;
        case 2:
          message.claimsEnabled = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Params {
    const obj = createBaseParams();
    if (isSet(object.distributionProportions)) obj.distributionProportions = DistributionProportions.fromJSON(object.distributionProportions);
    if (isSet(object.claimsEnabled)) obj.claimsEnabled = Boolean(object.claimsEnabled);
    return obj;
  },
  toJSON(message: Params): unknown {
    const obj: any = {};
    message.distributionProportions !== undefined && (obj.distributionProportions = message.distributionProportions ? DistributionProportions.toJSON(message.distributionProportions) : undefined);
    message.claimsEnabled !== undefined && (obj.claimsEnabled = message.claimsEnabled);
    return obj;
  },
  fromPartial(object: DeepPartial<Params>): Params {
    const message = createBaseParams();
    if (object.distributionProportions !== undefined && object.distributionProportions !== null) {
      message.distributionProportions = DistributionProportions.fromPartial(object.distributionProportions);
    }
    message.claimsEnabled = object.claimsEnabled ?? false;
    return message;
  },
  fromSDK(object: ParamsSDKType): Params {
    return {
      distributionProportions: object.distribution_proportions ? DistributionProportions.fromSDK(object.distribution_proportions) : undefined,
      claimsEnabled: object?.claims_enabled
    };
  },
  toSDK(message: Params): ParamsSDKType {
    const obj: any = {};
    message.distributionProportions !== undefined && (obj.distribution_proportions = message.distributionProportions ? DistributionProportions.toSDK(message.distributionProportions) : undefined);
    obj.claims_enabled = message.claimsEnabled;
    return obj;
  },
  fromAmino(object: ParamsAmino): Params {
    return {
      distributionProportions: object?.distribution_proportions ? DistributionProportions.fromAmino(object.distribution_proportions) : undefined,
      claimsEnabled: object.claims_enabled
    };
  },
  toAmino(message: Params): ParamsAmino {
    const obj: any = {};
    obj.distribution_proportions = message.distributionProportions ? DistributionProportions.toAmino(message.distributionProportions) : undefined;
    obj.claims_enabled = message.claimsEnabled;
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
      typeUrl: "/quicksilver.participationrewards.v1.Params",
      value: Params.encode(message).finish()
    };
  }
};
function createBaseKeyedProtocolData(): KeyedProtocolData {
  return {
    key: "",
    protocolData: ProtocolData.fromPartial({})
  };
}
export const KeyedProtocolData = {
  typeUrl: "/quicksilver.participationrewards.v1.KeyedProtocolData",
  encode(message: KeyedProtocolData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.protocolData !== undefined) {
      ProtocolData.encode(message.protocolData, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): KeyedProtocolData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseKeyedProtocolData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.protocolData = ProtocolData.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): KeyedProtocolData {
    const obj = createBaseKeyedProtocolData();
    if (isSet(object.key)) obj.key = String(object.key);
    if (isSet(object.protocolData)) obj.protocolData = ProtocolData.fromJSON(object.protocolData);
    return obj;
  },
  toJSON(message: KeyedProtocolData): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.protocolData !== undefined && (obj.protocolData = message.protocolData ? ProtocolData.toJSON(message.protocolData) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<KeyedProtocolData>): KeyedProtocolData {
    const message = createBaseKeyedProtocolData();
    message.key = object.key ?? "";
    if (object.protocolData !== undefined && object.protocolData !== null) {
      message.protocolData = ProtocolData.fromPartial(object.protocolData);
    }
    return message;
  },
  fromSDK(object: KeyedProtocolDataSDKType): KeyedProtocolData {
    return {
      key: object?.key,
      protocolData: object.protocol_data ? ProtocolData.fromSDK(object.protocol_data) : undefined
    };
  },
  toSDK(message: KeyedProtocolData): KeyedProtocolDataSDKType {
    const obj: any = {};
    obj.key = message.key;
    message.protocolData !== undefined && (obj.protocol_data = message.protocolData ? ProtocolData.toSDK(message.protocolData) : undefined);
    return obj;
  },
  fromAmino(object: KeyedProtocolDataAmino): KeyedProtocolData {
    return {
      key: object.key,
      protocolData: object?.protocol_data ? ProtocolData.fromAmino(object.protocol_data) : undefined
    };
  },
  toAmino(message: KeyedProtocolData): KeyedProtocolDataAmino {
    const obj: any = {};
    obj.key = message.key;
    obj.protocol_data = message.protocolData ? ProtocolData.toAmino(message.protocolData) : undefined;
    return obj;
  },
  fromAminoMsg(object: KeyedProtocolDataAminoMsg): KeyedProtocolData {
    return KeyedProtocolData.fromAmino(object.value);
  },
  fromProtoMsg(message: KeyedProtocolDataProtoMsg): KeyedProtocolData {
    return KeyedProtocolData.decode(message.value);
  },
  toProto(message: KeyedProtocolData): Uint8Array {
    return KeyedProtocolData.encode(message).finish();
  },
  toProtoMsg(message: KeyedProtocolData): KeyedProtocolDataProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.KeyedProtocolData",
      value: KeyedProtocolData.encode(message).finish()
    };
  }
};
function createBaseProtocolData(): ProtocolData {
  return {
    type: "",
    data: new Uint8Array()
  };
}
export const ProtocolData = {
  typeUrl: "/quicksilver.participationrewards.v1.ProtocolData",
  encode(message: ProtocolData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): ProtocolData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProtocolData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.string();
          break;
        case 2:
          message.data = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): ProtocolData {
    const obj = createBaseProtocolData();
    if (isSet(object.type)) obj.type = String(object.type);
    if (isSet(object.data)) obj.data = bytesFromBase64(object.data);
    return obj;
  },
  toJSON(message: ProtocolData): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.data !== undefined && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    return obj;
  },
  fromPartial(object: DeepPartial<ProtocolData>): ProtocolData {
    const message = createBaseProtocolData();
    message.type = object.type ?? "";
    message.data = object.data ?? new Uint8Array();
    return message;
  },
  fromSDK(object: ProtocolDataSDKType): ProtocolData {
    return {
      type: object?.type,
      data: object?.data
    };
  },
  toSDK(message: ProtocolData): ProtocolDataSDKType {
    const obj: any = {};
    obj.type = message.type;
    obj.data = message.data;
    return obj;
  },
  fromAmino(object: ProtocolDataAmino): ProtocolData {
    return {
      type: object.type,
      data: object.data
    };
  },
  toAmino(message: ProtocolData): ProtocolDataAmino {
    const obj: any = {};
    obj.type = message.type;
    obj.data = message.data;
    return obj;
  },
  fromAminoMsg(object: ProtocolDataAminoMsg): ProtocolData {
    return ProtocolData.fromAmino(object.value);
  },
  fromProtoMsg(message: ProtocolDataProtoMsg): ProtocolData {
    return ProtocolData.decode(message.value);
  },
  toProto(message: ProtocolData): Uint8Array {
    return ProtocolData.encode(message).finish();
  },
  toProtoMsg(message: ProtocolData): ProtocolDataProtoMsg {
    return {
      typeUrl: "/quicksilver.participationrewards.v1.ProtocolData",
      value: ProtocolData.encode(message).finish()
    };
  }
};