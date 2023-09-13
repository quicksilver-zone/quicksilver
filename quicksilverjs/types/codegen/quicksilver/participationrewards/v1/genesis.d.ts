import { Params, ParamsSDKType, KeyedProtocolData, KeyedProtocolDataSDKType } from "./participationrewards";
import * as _m0 from "protobufjs/minimal";
/** GenesisState defines the participationrewards module's genesis state. */
export interface GenesisState {
    params?: Params;
    protocolData: KeyedProtocolData[];
}
/** GenesisState defines the participationrewards module's genesis state. */
export interface GenesisStateSDKType {
    params?: ParamsSDKType;
    protocol_data: KeyedProtocolDataSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
