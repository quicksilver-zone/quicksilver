import { ProofOps, ProofOpsSDKType } from "../../../tendermint/crypto/proof";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/** MsgSubmitQueryResponse represents a message type to fulfil a query request. */
export interface MsgSubmitQueryResponse {
    chainId: string;
    queryId: string;
    result: Uint8Array;
    proofOps?: ProofOps;
    height: Long;
    fromAddress: string;
}
/** MsgSubmitQueryResponse represents a message type to fulfil a query request. */
export interface MsgSubmitQueryResponseSDKType {
    chain_id: string;
    query_id: string;
    result: Uint8Array;
    proof_ops?: ProofOpsSDKType;
    height: Long;
    from_address: string;
}
/**
 * MsgSubmitQueryResponseResponse defines the MsgSubmitQueryResponse response
 * type.
 */
export interface MsgSubmitQueryResponseResponse {
}
/**
 * MsgSubmitQueryResponseResponse defines the MsgSubmitQueryResponse response
 * type.
 */
export interface MsgSubmitQueryResponseResponseSDKType {
}
export declare const MsgSubmitQueryResponse: {
    encode(message: MsgSubmitQueryResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitQueryResponse;
    fromJSON(object: any): MsgSubmitQueryResponse;
    toJSON(message: MsgSubmitQueryResponse): unknown;
    fromPartial(object: Partial<MsgSubmitQueryResponse>): MsgSubmitQueryResponse;
};
export declare const MsgSubmitQueryResponseResponse: {
    encode(_: MsgSubmitQueryResponseResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitQueryResponseResponse;
    fromJSON(_: any): MsgSubmitQueryResponseResponse;
    toJSON(_: MsgSubmitQueryResponseResponse): unknown;
    fromPartial(_: Partial<MsgSubmitQueryResponseResponse>): MsgSubmitQueryResponseResponse;
};
