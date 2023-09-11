import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgRequestRedemption, MsgSignalIntent } from "./messages";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/quicksilver.interchainstaking.v1.MsgRequestRedemption", MsgRequestRedemption], ["/quicksilver.interchainstaking.v1.MsgSignalIntent", MsgSignalIntent]];
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
    }

  }
};