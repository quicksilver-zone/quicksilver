import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgSubmitClaim } from "./messages";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/quicksilver.participationrewards.v1.MsgSubmitClaim", MsgSubmitClaim]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    submitClaim(value: MsgSubmitClaim) {
      return {
        typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
        value: MsgSubmitClaim.encode(value).finish()
      };
    }

  },
  withTypeUrl: {
    submitClaim(value: MsgSubmitClaim) {
      return {
        typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
        value
      };
    }

  },
  toJSON: {
    submitClaim(value: MsgSubmitClaim) {
      return {
        typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
        value: MsgSubmitClaim.toJSON(value)
      };
    }

  },
  fromJSON: {
    submitClaim(value: any) {
      return {
        typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
        value: MsgSubmitClaim.fromJSON(value)
      };
    }

  },
  fromPartial: {
    submitClaim(value: MsgSubmitClaim) {
      return {
        typeUrl: "/quicksilver.participationrewards.v1.MsgSubmitClaim",
        value: MsgSubmitClaim.fromPartial(value)
      };
    }

  }
};