import { Proof, ProofSDKType } from "../../claimsmanager/v1/claimsmanager";
import * as fm from "../../../grpc-gateway";
import { MsgClaim, MsgClaimSDKType, MsgClaimResponse, MsgClaimResponseSDKType } from "./messages";
export class Msg {
  static claim(request: MsgClaim, initRequest?: fm.InitReq): Promise<MsgClaimResponse> {
    return fm.fetchReq(`/quicksilver.airdrop.v1/claim`, {
      ...initRequest,
      method: "POST",
      body: JSON.stringify(request, fm.replacer)
    });
  }
}