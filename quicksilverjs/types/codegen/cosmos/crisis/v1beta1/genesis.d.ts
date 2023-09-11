import { Coin, CoinSDKType } from "../../base/v1beta1/coin";
import * as _m0 from "protobufjs/minimal";
/** GenesisState defines the crisis module's genesis state. */
export interface GenesisState {
    /**
     * constant_fee is the fee used to verify the invariant in the crisis
     * module.
     */
    constantFee?: Coin;
}
/** GenesisState defines the crisis module's genesis state. */
export interface GenesisStateSDKType {
    /**
     * constant_fee is the fee used to verify the invariant in the crisis
     * module.
     */
    constant_fee?: CoinSDKType;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
