import * as _m0 from "protobufjs/minimal";
export declare enum ProtocolDataType {
    /** ProtocolDataTypeUndefined - Undefined action (per protobuf spec) */
    ProtocolDataTypeUndefined = 0,
    ProtocolDataTypeConnection = 1,
    ProtocolDataTypeOsmosisParams = 2,
    ProtocolDataTypeLiquidToken = 3,
    ProtocolDataTypeOsmosisPool = 4,
    ProtocolDataTypeCrescentPool = 5,
    ProtocolDataTypeSifchainPool = 6,
    UNRECOGNIZED = -1
}
export declare enum ProtocolDataTypeSDKType {
    /** ProtocolDataTypeUndefined - Undefined action (per protobuf spec) */
    ProtocolDataTypeUndefined = 0,
    ProtocolDataTypeConnection = 1,
    ProtocolDataTypeOsmosisParams = 2,
    ProtocolDataTypeLiquidToken = 3,
    ProtocolDataTypeOsmosisPool = 4,
    ProtocolDataTypeCrescentPool = 5,
    ProtocolDataTypeSifchainPool = 6,
    UNRECOGNIZED = -1
}
export declare function protocolDataTypeFromJSON(object: any): ProtocolDataType;
export declare function protocolDataTypeToJSON(object: ProtocolDataType): string;
/**
 * DistributionProportions defines the proportions of minted QCK that is to be
 * allocated as participation rewards.
 */
export interface DistributionProportions {
    validatorSelectionAllocation: string;
    holdingsAllocation: string;
    lockupAllocation: string;
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
export interface Params {
    /**
     * distribution_proportions defines the proportions of the minted
     * participation rewards;
     */
    distributionProportions?: DistributionProportions;
}
/** Params holds parameters for the participationrewards module. */
export interface ParamsSDKType {
    /**
     * distribution_proportions defines the proportions of the minted
     * participation rewards;
     */
    distribution_proportions?: DistributionProportionsSDKType;
}
export interface KeyedProtocolData {
    key: string;
    protocolData?: ProtocolData;
}
export interface KeyedProtocolDataSDKType {
    key: string;
    protocol_data?: ProtocolDataSDKType;
}
/**
 * Protocol Data is an arbitrary data type held against a given zone for the
 * determination of rewards.
 */
export interface ProtocolData {
    type: string;
    data: Uint8Array;
}
/**
 * Protocol Data is an arbitrary data type held against a given zone for the
 * determination of rewards.
 */
export interface ProtocolDataSDKType {
    type: string;
    data: Uint8Array;
}
export declare const DistributionProportions: {
    encode(message: DistributionProportions, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DistributionProportions;
    fromJSON(object: any): DistributionProportions;
    toJSON(message: DistributionProportions): unknown;
    fromPartial(object: Partial<DistributionProportions>): DistributionProportions;
};
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial(object: Partial<Params>): Params;
};
export declare const KeyedProtocolData: {
    encode(message: KeyedProtocolData, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): KeyedProtocolData;
    fromJSON(object: any): KeyedProtocolData;
    toJSON(message: KeyedProtocolData): unknown;
    fromPartial(object: Partial<KeyedProtocolData>): KeyedProtocolData;
};
export declare const ProtocolData: {
    encode(message: ProtocolData, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ProtocolData;
    fromJSON(object: any): ProtocolData;
    toJSON(message: ProtocolData): unknown;
    fromPartial(object: Partial<ProtocolData>): ProtocolData;
};
