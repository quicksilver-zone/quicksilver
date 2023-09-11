import { Rpc } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgClaim, MsgClaimResponse } from "./messages";
/** Msg defines the airdrop Msg service. */

export interface Msg {
  claim(request: MsgClaim): Promise<MsgClaimResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.claim = this.claim.bind(this);
  }

  claim(request: MsgClaim): Promise<MsgClaimResponse> {
    const data = MsgClaim.encode(request).finish();
    const promise = this.rpc.request("quicksilver.airdrop.v1.Msg", "Claim", data);
    return promise.then(data => MsgClaimResponse.decode(new _m0.Reader(data)));
  }

}