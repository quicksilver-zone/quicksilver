import { Rpc } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgRequestRedemption, MsgRequestRedemptionResponse, MsgSignalIntent, MsgSignalIntentResponse } from "./messages";
/** Msg defines the interchainstaking Msg service. */

export interface Msg {
  /**
   * RequestRedemption defines a method for requesting burning of qAssets for
   * native assets.
   */
  requestRedemption(request: MsgRequestRedemption): Promise<MsgRequestRedemptionResponse>;
  /**
   * SignalIntent defines a method for signalling voting intent for one or more
   * validators.
   */

  signalIntent(request: MsgSignalIntent): Promise<MsgSignalIntentResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.requestRedemption = this.requestRedemption.bind(this);
    this.signalIntent = this.signalIntent.bind(this);
  }

  requestRedemption(request: MsgRequestRedemption): Promise<MsgRequestRedemptionResponse> {
    const data = MsgRequestRedemption.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Msg", "RequestRedemption", data);
    return promise.then(data => MsgRequestRedemptionResponse.decode(new _m0.Reader(data)));
  }

  signalIntent(request: MsgSignalIntent): Promise<MsgSignalIntentResponse> {
    const data = MsgSignalIntent.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Msg", "SignalIntent", data);
    return promise.then(data => MsgSignalIntentResponse.decode(new _m0.Reader(data)));
  }

}