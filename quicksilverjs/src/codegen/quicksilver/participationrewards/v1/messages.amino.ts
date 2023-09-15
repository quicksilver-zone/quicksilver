import { ClaimType, ClaimTypeSDKType, Proof, ProofSDKType } from "../../claimsmanager/v1/claimsmanager";
import { MsgSubmitClaim, MsgSubmitClaimSDKType } from "./messages";
export const AminoConverter = {
  "/quicksilver.participationrewards.v1.MsgSubmitClaim": {
    aminoType: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
    toAmino: MsgSubmitClaim.toAmino,
    fromAmino: MsgSubmitClaim.fromAmino
  }
};