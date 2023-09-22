import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgRequestRedemption, MsgRequestRedemptionSDKType, MsgSignalIntent, MsgSignalIntentSDKType } from "./messages";
import { MsgGovCloseChannel, MsgGovCloseChannelSDKType, MsgGovReopenChannel, MsgGovReopenChannelSDKType } from "./proposals";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/quicksilver.interchainstaking.v1.MsgRequestRedemption", MsgRequestRedemption], ["/quicksilver.interchainstaking.v1.MsgSignalIntent", MsgSignalIntent], ["/quicksilver.interchainstaking.v1.MsgGovCloseChannel", MsgGovCloseChannel], ["/quicksilver.interchainstaking.v1.MsgGovReopenChannel", MsgGovReopenChannel]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    requestRedemption(value: MsgRequestRedemption) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
        value: MsgRequestRedemption.encode(value).finish()
      };
    },
    signalIntent(value: MsgSignalIntent) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
        value: MsgSignalIntent.encode(value).finish()
      };
    },
    govCloseChannel(value: MsgGovCloseChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
        value: MsgGovCloseChannel.encode(value).finish()
      };
    },
    govReopenChannel(value: MsgGovReopenChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
        value: MsgGovReopenChannel.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    requestRedemption(value: MsgRequestRedemption) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
        value
      };
    },
    signalIntent(value: MsgSignalIntent) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
        value
      };
    },
    govCloseChannel(value: MsgGovCloseChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
        value
      };
    },
    govReopenChannel(value: MsgGovReopenChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
        value
      };
    }
  },
  toJSON: {
    requestRedemption(value: MsgRequestRedemption) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
        value: MsgRequestRedemption.toJSON(value)
      };
    },
    signalIntent(value: MsgSignalIntent) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
        value: MsgSignalIntent.toJSON(value)
      };
    },
    govCloseChannel(value: MsgGovCloseChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
        value: MsgGovCloseChannel.toJSON(value)
      };
    },
    govReopenChannel(value: MsgGovReopenChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
        value: MsgGovReopenChannel.toJSON(value)
      };
    }
  },
  fromJSON: {
    requestRedemption(value: any) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
        value: MsgRequestRedemption.fromJSON(value)
      };
    },
    signalIntent(value: any) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
        value: MsgSignalIntent.fromJSON(value)
      };
    },
    govCloseChannel(value: any) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
        value: MsgGovCloseChannel.fromJSON(value)
      };
    },
    govReopenChannel(value: any) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
        value: MsgGovReopenChannel.fromJSON(value)
      };
    }
  },
  fromPartial: {
    requestRedemption(value: MsgRequestRedemption) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgRequestRedemption",
        value: MsgRequestRedemption.fromPartial(value)
      };
    },
    signalIntent(value: MsgSignalIntent) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgSignalIntent",
        value: MsgSignalIntent.fromPartial(value)
      };
    },
    govCloseChannel(value: MsgGovCloseChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
        value: MsgGovCloseChannel.fromPartial(value)
      };
    },
    govReopenChannel(value: MsgGovReopenChannel) {
      return {
        typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
        value: MsgGovReopenChannel.fromPartial(value)
      };
    }
  }
};