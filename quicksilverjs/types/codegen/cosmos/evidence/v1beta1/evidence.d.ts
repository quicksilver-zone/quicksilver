import { Timestamp, TimestampSDKType } from "../../../google/protobuf/timestamp";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/**
 * Equivocation implements the Evidence interface and defines evidence of double
 * signing misbehavior.
 */
export interface Equivocation {
    height: Long;
    time?: Timestamp;
    power: Long;
    consensusAddress: string;
}
/**
 * Equivocation implements the Evidence interface and defines evidence of double
 * signing misbehavior.
 */
export interface EquivocationSDKType {
    height: Long;
    time?: TimestampSDKType;
    power: Long;
    consensus_address: string;
}
export declare const Equivocation: {
    encode(message: Equivocation, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Equivocation;
    fromJSON(object: any): Equivocation;
    toJSON(message: Equivocation): unknown;
    fromPartial(object: Partial<Equivocation>): Equivocation;
};
