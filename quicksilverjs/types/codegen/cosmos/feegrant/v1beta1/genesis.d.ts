import { Grant, GrantSDKType } from "./feegrant";
import * as _m0 from "protobufjs/minimal";
/** GenesisState contains a set of fee allowances, persisted from the store */
export interface GenesisState {
    allowances: Grant[];
}
/** GenesisState contains a set of fee allowances, persisted from the store */
export interface GenesisStateSDKType {
    allowances: GrantSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
