import { ProofOps, ProofOpsSDKType } from "../../../tendermint/crypto/proof";
import * as fm from "../../../grpc-gateway";
import { MsgSubmitQueryResponse, MsgSubmitQueryResponseSDKType, MsgSubmitQueryResponseResponse, MsgSubmitQueryResponseResponseSDKType } from "./messages";
export class Msg {
  /** SubmitQueryResponse defines a method for submit query responses. */
  static submitQueryResponse(request: MsgSubmitQueryResponse, initRequest?: fm.InitReq): Promise<MsgSubmitQueryResponseResponse> {
    return fm.fetchReq(`/quicksilver.interchainquery.v1/submitQueryResponse`, {
      ...initRequest,
      method: "POST",
      body: JSON.stringify(request, fm.replacer)
    });
  }
}