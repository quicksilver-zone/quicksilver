import { AminoMsg } from "@cosmjs/amino";
import { MsgSubmitClaim } from "./messages";
export interface AminoMsgSubmitClaim extends AminoMsg {
    type: "quicksilver/MsgSubmitClaim";
    value: {
        user_address: string;
        zone: string;
        src_zone: string;
        claim_type: number;
        proofs: {
            key: Uint8Array;
            data: Uint8Array;
            proof_ops: {
                ops: {
                    type: string;
                    key: Uint8Array;
                    data: Uint8Array;
                }[];
            };
            height: string;
            proof_type: string;
        }[];
    };
}
export declare const AminoConverter: {
    "/quicksilver.participationrewards.v1.MsgSubmitClaim": {
        aminoType: string;
        toAmino: ({ userAddress, zone, srcZone, claimType, proofs }: MsgSubmitClaim) => AminoMsgSubmitClaim["value"];
        fromAmino: ({ user_address, zone, src_zone, claim_type, proofs }: AminoMsgSubmitClaim["value"]) => MsgSubmitClaim;
    };
};
