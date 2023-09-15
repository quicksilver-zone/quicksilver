import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import * as fm from "../../../grpc-gateway";
import { MsgRequestRedemption, MsgRequestRedemptionSDKType, MsgRequestRedemptionResponse, MsgRequestRedemptionResponseSDKType, MsgSignalIntent, MsgSignalIntentSDKType, MsgSignalIntentResponse, MsgSignalIntentResponseSDKType } from "./messages";
import { MsgGovCloseChannel, MsgGovCloseChannelSDKType, MsgGovCloseChannelResponse, MsgGovCloseChannelResponseSDKType, MsgGovReopenChannel, MsgGovReopenChannelSDKType, MsgGovReopenChannelResponse, MsgGovReopenChannelResponseSDKType } from "./proposals";
export class Msg {
  /**
   * RequestRedemption defines a method for requesting burning of qAssets for
   * native assets.
   */
  static requestRedemption(request: MsgRequestRedemption, initRequest?: fm.InitReq): Promise<MsgRequestRedemptionResponse> {
    return fm.fetchReq(`/quicksilver.interchainstaking.v1/requestRedemption`, {
      ...initRequest,
      method: "POST",
      body: JSON.stringify(request, fm.replacer)
    });
  }
  /**
   * SignalIntent defines a method for signalling voting intent for one or more
   * validators.
   */
  static signalIntent(request: MsgSignalIntent, initRequest?: fm.InitReq): Promise<MsgSignalIntentResponse> {
    return fm.fetchReq(`/quicksilver.interchainstaking.v1/signalIntent`, {
      ...initRequest,
      method: "POST",
      body: JSON.stringify(request, fm.replacer)
    });
  }
  static govCloseChannel(request: MsgGovCloseChannel, initRequest?: fm.InitReq): Promise<MsgGovCloseChannelResponse> {
    return fm.fetchReq(`/quicksilver.interchainstaking.v1/govCloseChannel`, {
      ...initRequest,
      method: "POST",
      body: JSON.stringify(request, fm.replacer)
    });
  }
  static govReopenChannel(request: MsgGovReopenChannel, initRequest?: fm.InitReq): Promise<MsgGovReopenChannelResponse> {
    return fm.fetchReq(`/quicksilver.interchainstaking.v1/govReopenChannel`, {
      ...initRequest,
      method: "POST",
      body: JSON.stringify(request, fm.replacer)
    });
  }
}