import { Proof, ProofSDKType } from "../../claimsmanager/v1/claimsmanager";
import { MsgClaim, MsgClaimSDKType } from "./messages";
export const AminoConverter = {
  "/quicksilver.airdrop.v1.MsgClaim": {
    aminoType: "/quicksilver.airdrop.v1.MsgClaim",
    toAmino: MsgClaim.toAmino,
    fromAmino: MsgClaim.fromAmino
  }
};