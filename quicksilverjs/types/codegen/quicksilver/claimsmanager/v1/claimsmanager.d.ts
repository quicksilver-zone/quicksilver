import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
export declare enum ClaimType {
    /** ClaimTypeUndefined - Undefined action (per protobuf spec) */
    ClaimTypeUndefined = 0,
    ClaimTypeLiquidToken = 1,
    ClaimTypeOsmosisPool = 2,
    ClaimTypeCrescentPool = 3,
    ClaimTypeSifchainPool = 4,
    UNRECOGNIZED = -1
}
export declare enum ClaimTypeSDKType {
    /** ClaimTypeUndefined - Undefined action (per protobuf spec) */
    ClaimTypeUndefined = 0,
    ClaimTypeLiquidToken = 1,
    ClaimTypeOsmosisPool = 2,
    ClaimTypeCrescentPool = 3,
    ClaimTypeSifchainPool = 4,
    UNRECOGNIZED = -1
}
export declare function claimTypeFromJSON(object: any): ClaimType;
export declare function claimTypeToJSON(object: ClaimType): string;
/** Params holds parameters for the claimsmanager module. */
export interface Params {
}
/** Params holds parameters for the claimsmanager module. */
export interface ParamsSDKType {
}
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
export declare const Params: {
    encode(_: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromJSON(_: any): Params;
    toJSON(_: Params): unknown;
    fromPartial(_: Partial<Params>): Params;
};
export declare const Claim: {
    encode(message: Claim, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Claim;
    fromJSON(object: any): Claim;
    toJSON(message: Claim): unknown;
    fromPartial(object: Partial<Claim>): Claim;
};
