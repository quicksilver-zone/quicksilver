import { AminoMsg } from "@cosmjs/amino";
import { MsgSubmitQueryResponse } from "./messages";
export interface AminoMsgSubmitQueryResponse extends AminoMsg {
    type: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse";
    value: {
        chain_id: string;
        query_id: string;
        result: Uint8Array;
        proof_ops: {
            ops: {
                type: string;
                key: Uint8Array;
                data: Uint8Array;
            }[];
        };
        height: string;
        from_address: string;
    };
}
export declare const AminoConverter: {
    "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse": {
        aminoType: string;
        toAmino: ({ chainId, queryId, result, proofOps, height, fromAddress }: MsgSubmitQueryResponse) => AminoMsgSubmitQueryResponse["value"];
        fromAmino: ({ chain_id, query_id, result, proof_ops, height, from_address }: AminoMsgSubmitQueryResponse["value"]) => MsgSubmitQueryResponse;
    };
};
