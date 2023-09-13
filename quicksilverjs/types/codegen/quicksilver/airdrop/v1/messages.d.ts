import { ProofOps, ProofOpsSDKType } from "../../../tendermint/crypto/proof";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
export interface MsgClaim {
    chainId: string;
    action: Long;
    address: string;
    proofs: Proof[];
}
export interface MsgClaimSDKType {
    chain_id: string;
    action: Long;
    address: string;
    proofs: ProofSDKType[];
}
export interface MsgClaimResponse {
    amount: Long;
}
export interface MsgClaimResponseSDKType {
    amount: Long;
}
export interface Proof {
    key: Uint8Array;
    data: Uint8Array;
    proofOps?: ProofOps;
    height: Long;
}
export interface ProofSDKType {
    key: Uint8Array;
    data: Uint8Array;
    proof_ops?: ProofOpsSDKType;
    height: Long;
}
export declare const MsgClaim: {
    encode(message: MsgClaim, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaim;
    fromJSON(object: any): MsgClaim;
    toJSON(message: MsgClaim): unknown;
    fromPartial(object: Partial<MsgClaim>): MsgClaim;
};
export declare const MsgClaimResponse: {
    encode(message: MsgClaimResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimResponse;
    fromJSON(object: any): MsgClaimResponse;
    toJSON(message: MsgClaimResponse): unknown;
    fromPartial(object: Partial<MsgClaimResponse>): MsgClaimResponse;
};
export declare const Proof: {
    encode(message: Proof, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Proof;
    fromJSON(object: any): Proof;
    toJSON(message: Proof): unknown;
    fromPartial(object: Partial<Proof>): Proof;
};
