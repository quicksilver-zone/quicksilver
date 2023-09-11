import { Params, ParamsSDKType } from "./params";
import { ZoneDrop, ZoneDropSDKType, ClaimRecord, ClaimRecordSDKType } from "./airdrop";
import * as _m0 from "protobufjs/minimal";
/** GenesisState defines the airdrop module's genesis state. */
export interface GenesisState {
    params?: Params;
    zoneDrops: ZoneDrop[];
    claimRecords: ClaimRecord[];
}
/** GenesisState defines the airdrop module's genesis state. */
export interface GenesisStateSDKType {
    params?: ParamsSDKType;
    zone_drops: ZoneDropSDKType[];
    claim_records: ClaimRecordSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
