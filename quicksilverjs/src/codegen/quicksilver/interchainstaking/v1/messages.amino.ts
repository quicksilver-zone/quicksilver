import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import { MsgRequestRedemption, MsgRequestRedemptionSDKType, MsgSignalIntent, MsgSignalIntentSDKType } from "./messages";
import { MsgGovCloseChannel, MsgGovCloseChannelSDKType, MsgGovReopenChannel, MsgGovReopenChannelSDKType } from "./proposals";
export const AminoConverter = {
  "/quicksilver.interchainstaking.v1.MsgRequestRedemption": {
    aminoType: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
    toAmino: MsgRequestRedemption.toAmino,
    fromAmino: MsgRequestRedemption.fromAmino
  },
  "/quicksilver.interchainstaking.v1.MsgSignalIntent": {
    aminoType: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
    toAmino: MsgSignalIntent.toAmino,
    fromAmino: MsgSignalIntent.fromAmino
  },
  "/quicksilver.interchainstaking.v1.MsgGovCloseChannel": {
    aminoType: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
    toAmino: MsgGovCloseChannel.toAmino,
    fromAmino: MsgGovCloseChannel.fromAmino
  },
  "/quicksilver.interchainstaking.v1.MsgGovReopenChannel": {
    aminoType: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
    toAmino: MsgGovReopenChannel.toAmino,
    fromAmino: MsgGovReopenChannel.fromAmino
  }
};