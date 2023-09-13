import { Params, ParamsSDKType, Claim, ClaimSDKType } from "./claimsmanager";
import * as _m0 from "protobufjs/minimal";
/** GenesisState defines the claimsmanager module's genesis state. */
export interface GenesisState {
    params?: Params;
    claims: Claim[];
}
/** GenesisState defines the claimsmanager module's genesis state. */
export interface GenesisStateSDKType {
    params?: ParamsSDKType;
    claims: ClaimSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
