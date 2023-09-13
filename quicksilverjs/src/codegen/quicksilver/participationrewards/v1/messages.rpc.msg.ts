import { Rpc } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSubmitClaim, MsgSubmitClaimResponse } from "./messages";
/** Msg defines the participationrewards Msg service. */

export interface Msg {
  submitClaim(request: MsgSubmitClaim): Promise<MsgSubmitClaimResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.submitClaim = this.submitClaim.bind(this);
  }

  submitClaim(request: MsgSubmitClaim): Promise<MsgSubmitClaimResponse> {
    const data = MsgSubmitClaim.encode(request).finish();
    const promise = this.rpc.request("quicksilver.participationrewards.v1.Msg", "SubmitClaim", data);
    return promise.then(data => MsgSubmitClaimResponse.decode(new _m0.Reader(data)));
  }

}