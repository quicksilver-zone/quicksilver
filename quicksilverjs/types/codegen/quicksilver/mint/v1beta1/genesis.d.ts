import { Minter, MinterSDKType, Params, ParamsSDKType } from "./mint";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/** GenesisState defines the mint module's genesis state. */
export interface GenesisState {
    /** minter is a space for holding current rewards information. */
    minter?: Minter;
    /** params defines all the paramaters of the module. */
    params?: Params;
    /** current reduction period start epoch */
    reductionStartedEpoch: Long;
}
/** GenesisState defines the mint module's genesis state. */
export interface GenesisStateSDKType {
    /** minter is a space for holding current rewards information. */
    minter?: MinterSDKType;
    /** params defines all the paramaters of the module. */
    params?: ParamsSDKType;
    /** current reduction period start epoch */
    reduction_started_epoch: Long;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
