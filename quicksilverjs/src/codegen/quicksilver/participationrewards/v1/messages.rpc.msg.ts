import { ClaimType, ClaimTypeSDKType, Proof, ProofSDKType } from "../../claimsmanager/v1/claimsmanager";
import * as fm from "../../../grpc-gateway";
import { MsgSubmitClaim, MsgSubmitClaimSDKType, MsgSubmitClaimResponse, MsgSubmitClaimResponseSDKType } from "./messages";
export class Msg {
  static submitClaim(request: MsgSubmitClaim, initRequest?: fm.InitReq): Promise<MsgSubmitClaimResponse> {
    return fm.fetchReq(`/quicksilver.participationrewards.v1/submitClaim`, {
      ...initRequest,
      method: "POST",
      body: JSON.stringify(request, fm.replacer)
    });
  }
}