import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/** Minter represents the minting state. */
export interface Minter {
    /** current annual inflation rate */
    inflation: string;
    /** current annual expected provisions */
    annualProvisions: string;
}
/** Minter represents the minting state. */
export interface MinterSDKType {
    /** current annual inflation rate */
    inflation: string;
    /** current annual expected provisions */
    annual_provisions: string;
}
/** Params holds parameters for the mint module. */
export interface Params {
    /** type of coin to mint */
    mintDenom: string;
    /** maximum annual change in inflation rate */
    inflationRateChange: string;
    /** maximum inflation rate */
    inflationMax: string;
    /** minimum inflation rate */
    inflationMin: string;
    /** goal of percent bonded atoms */
    goalBonded: string;
    /** expected blocks per year */
    blocksPerYear: Long;
}
/** Params holds parameters for the mint module. */
export interface ParamsSDKType {
    /** type of coin to mint */
    mint_denom: string;
    /** maximum annual change in inflation rate */
    inflation_rate_change: string;
    /** maximum inflation rate */
    inflation_max: string;
    /** minimum inflation rate */
    inflation_min: string;
    /** goal of percent bonded atoms */
    goal_bonded: string;
    /** expected blocks per year */
    blocks_per_year: Long;
}
export declare const Minter: {
    encode(message: Minter, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Minter;
    fromJSON(object: any): Minter;
    toJSON(message: Minter): unknown;
    fromPartial(object: Partial<Minter>): Minter;
};
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial(object: Partial<Params>): Params;
};
