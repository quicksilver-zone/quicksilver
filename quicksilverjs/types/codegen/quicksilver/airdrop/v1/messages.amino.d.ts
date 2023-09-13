import { AminoMsg } from "@cosmjs/amino";
import { MsgClaim } from "./messages";
export interface AminoMsgClaim extends AminoMsg {
    type: "quicksilver/MsgClaim";
    value: {
        chain_id: string;
        action: string;
        address: string;
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
        }[];
    };
}
export declare const AminoConverter: {
    "/quicksilver.airdrop.v1.MsgClaim": {
        aminoType: string;
        toAmino: ({ chainId, action, address, proofs }: MsgClaim) => AminoMsgClaim["value"];
        fromAmino: ({ chain_id, action, address, proofs }: AminoMsgClaim["value"]) => MsgClaim;
    };
};
