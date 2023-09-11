import { Query, QuerySDKType } from "./interchainquery";
import * as _m0 from "protobufjs/minimal";
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisState {
    queries: Query[];
}
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisStateSDKType {
    queries: QuerySDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
