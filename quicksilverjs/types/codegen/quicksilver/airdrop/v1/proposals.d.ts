import { ZoneDrop, ZoneDropSDKType } from "./airdrop";
import * as _m0 from "protobufjs/minimal";
export interface RegisterZoneDropProposal {
    title: string;
    description: string;
    zoneDrop?: ZoneDrop;
    claimRecords: Uint8Array;
}
export interface RegisterZoneDropProposalSDKType {
    title: string;
    description: string;
    zone_drop?: ZoneDropSDKType;
    claim_records: Uint8Array;
}
export declare const RegisterZoneDropProposal: {
    encode(message: RegisterZoneDropProposal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): RegisterZoneDropProposal;
    fromJSON(object: any): RegisterZoneDropProposal;
    toJSON(message: RegisterZoneDropProposal): unknown;
    fromPartial(object: Partial<RegisterZoneDropProposal>): RegisterZoneDropProposal;
};
