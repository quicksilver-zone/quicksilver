import { Timestamp, TimestampSDKType } from "../../../google/protobuf/timestamp";
import { Duration, DurationSDKType } from "../../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
export interface EpochInfo {
    identifier: string;
    startTime?: Timestamp;
    duration?: Duration;
    currentEpoch: Long;
    currentEpochStartTime?: Timestamp;
    epochCountingStarted: boolean;
    currentEpochStartHeight: Long;
}
export interface EpochInfoSDKType {
    identifier: string;
    start_time?: TimestampSDKType;
    duration?: DurationSDKType;
    current_epoch: Long;
    current_epoch_start_time?: TimestampSDKType;
    epoch_counting_started: boolean;
    current_epoch_start_height: Long;
}
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisState {
    epochs: EpochInfo[];
}
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisStateSDKType {
    epochs: EpochInfoSDKType[];
}
export declare const EpochInfo: {
    encode(message: EpochInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): EpochInfo;
    fromJSON(object: any): EpochInfo;
    toJSON(message: EpochInfo): unknown;
    fromPartial(object: Partial<EpochInfo>): EpochInfo;
};
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
