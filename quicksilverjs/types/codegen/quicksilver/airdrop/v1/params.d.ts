import * as _m0 from "protobufjs/minimal";
/** Params holds parameters for the airdrop module. */
export interface Params {
}
/** Params holds parameters for the airdrop module. */
export interface ParamsSDKType {
}
export declare const Params: {
    encode(_: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromJSON(_: any): Params;
    toJSON(_: Params): unknown;
    fromPartial(_: Partial<Params>): Params;
};
