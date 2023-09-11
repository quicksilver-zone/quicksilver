import { ClaimType, ClaimTypeSDKType } from "../../claimsmanager/v1/claimsmanager";
import { ProofOps, ProofOpsSDKType } from "../../../tendermint/crypto/proof";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/**
 * MsgSubmitClaim represents a message type for submitting a participation
 * claim regarding the given zone (chain).
 */
export interface MsgSubmitClaim {
    userAddress: string;
    zone: string;
    srcZone: string;
    claimType: ClaimType;
    proofs: Proof[];
}
/**
 * MsgSubmitClaim represents a message type for submitting a participation
 * claim regarding the given zone (chain).
 */
export interface MsgSubmitClaimSDKType {
    user_address: string;
    zone: string;
    src_zone: string;
    claim_type: ClaimTypeSDKType;
    proofs: ProofSDKType[];
}
/** MsgSubmitClaimResponse defines the MsgSubmitClaim response type. */
export interface MsgSubmitClaimResponse {
}
/** MsgSubmitClaimResponse defines the MsgSubmitClaim response type. */
export interface MsgSubmitClaimResponseSDKType {
}
/** Proof defines a type used to cryptographically prove a claim. */
export interface Proof {
    key: Uint8Array;
    data: Uint8Array;
    proofOps?: ProofOps;
    height: Long;
    proofType: string;
}
/** Proof defines a type used to cryptographically prove a claim. */
export interface ProofSDKType {
    key: Uint8Array;
    data: Uint8Array;
    proof_ops?: ProofOpsSDKType;
    height: Long;
    proof_type: string;
}
export declare const MsgSubmitClaim: {
    encode(message: MsgSubmitClaim, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitClaim;
    fromJSON(object: any): MsgSubmitClaim;
    toJSON(message: MsgSubmitClaim): unknown;
    fromPartial(object: Partial<MsgSubmitClaim>): MsgSubmitClaim;
};
export declare const MsgSubmitClaimResponse: {
    encode(_: MsgSubmitClaimResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitClaimResponse;
    fromJSON(_: any): MsgSubmitClaimResponse;
    toJSON(_: MsgSubmitClaimResponse): unknown;
    fromPartial(_: Partial<MsgSubmitClaimResponse>): MsgSubmitClaimResponse;
};
export declare const Proof: {
    encode(message: Proof, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Proof;
    fromJSON(object: any): Proof;
    toJSON(message: Proof): unknown;
    fromPartial(object: Partial<Proof>): Proof;
};
