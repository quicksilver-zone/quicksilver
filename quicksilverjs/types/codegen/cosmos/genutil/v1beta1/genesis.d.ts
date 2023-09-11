import * as _m0 from "protobufjs/minimal";
/** GenesisState defines the raw genesis transaction in JSON. */
export interface GenesisState {
    /** gen_txs defines the genesis transactions. */
    genTxs: Uint8Array[];
}
/** GenesisState defines the raw genesis transaction in JSON. */
export interface GenesisStateSDKType {
    /** gen_txs defines the genesis transactions. */
    gen_txs: Uint8Array[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
