//@ts-nocheck
import { AminoMsg } from "@cosmjs/amino";
import { Long } from "../../../helpers";
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
export const AminoConverter = {
  "/quicksilver.airdrop.v1.MsgClaim": {
    aminoType: "quicksilver/MsgClaim",
    toAmino: ({
      chainId,
      action,
      address,
      proofs
    }: MsgClaim): AminoMsgClaim["value"] => {
      return {
        chain_id: chainId,
        action: action.toString(),
        address,
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
          height: el0.height.toString()
        }))
      };
    },
    fromAmino: ({
      chain_id,
      action,
      address,
      proofs
    }: AminoMsgClaim["value"]): MsgClaim => {
      return {
        chainId: chain_id,
        action: Long.fromString(action),
        address,
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
          height: Long.fromString(el0.height)
        }))
      };
    }
  }
};