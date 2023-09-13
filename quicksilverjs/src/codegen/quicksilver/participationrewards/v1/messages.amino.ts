//@ts-nocheck
import { claimTypeFromJSON } from "../../claimsmanager/v1/claimsmanager";
import { AminoMsg } from "@cosmjs/amino";
import { Long } from "../../../helpers";
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
export const AminoConverter = {
  "/quicksilver.participationrewards.v1.MsgSubmitClaim": {
    aminoType: "quicksilver/MsgSubmitClaim",
    toAmino: ({
      userAddress,
      zone,
      srcZone,
      claimType,
      proofs
    }: MsgSubmitClaim): AminoMsgSubmitClaim["value"] => {
      return {
        user_address: userAddress,
        zone,
        src_zone: srcZone,
        claim_type: claimType,
        proofs: proofs.map(el0 => ({
          key: el0.key,
          data: el0.data,
          proof_ops: {
            ops: el0.proofOps.ops.map(el1 => ({
              type: el1.type,
              key: el1.key,
              data: el1.data
            }))
          },
          height: el0.height.toString(),
          proof_type: el0.proofType
        }))
      };
    },
    fromAmino: ({
      user_address,
      zone,
      src_zone,
      claim_type,
      proofs
    }: AminoMsgSubmitClaim["value"]): MsgSubmitClaim => {
      return {
        userAddress: user_address,
        zone,
        srcZone: src_zone,
        claimType: claimTypeFromJSON(claim_type),
        proofs: proofs.map(el0 => ({
          key: el0.key,
          data: el0.data,
          proofOps: {
            ops: el0.proof_ops.ops.map(el2 => ({
              type: el2.type,
              key: el2.key,
              data: el2.data
            }))
          },
          height: Long.fromString(el0.height),
          proofType: el0.proof_type
        }))
      };
    }
  }
};