//@ts-nocheck
import { AminoMsg } from "@cosmjs/amino";
import { Long } from "../../../helpers";
import { MsgRequestRedemption, MsgSignalIntent } from "./messages";
export interface AminoMsgRequestRedemption extends AminoMsg {
  type: "quicksilver/MsgRequestRedemption";
  value: {
    value: {
      denom: string;
      amount: string;
    };
    destination_address: string;
    from_address: string;
  };
}
export interface AminoMsgSignalIntent extends AminoMsg {
  type: "quicksilver/MsgSignalIntent";
  value: {
    chain_id: string;
    intents: {
      valoper_address: string;
      weight: string;
    }[];
    from_address: string;
  };
}
export const AminoConverter = {
  "/quicksilver.interchainstaking.v1.MsgRequestRedemption": {
    aminoType: "quicksilver/MsgRequestRedemption",
    toAmino: ({
      value,
      destinationAddress,
      fromAddress
    }: MsgRequestRedemption): AminoMsgRequestRedemption["value"] => {
      return {
        value: {
          denom: value.denom,
          amount: Long.fromValue(value.amount).toString()
        },
        destination_address: destinationAddress,
        from_address: fromAddress
      };
    },
    fromAmino: ({
      value,
      destination_address,
      from_address
    }: AminoMsgRequestRedemption["value"]): MsgRequestRedemption => {
      return {
        value: {
          denom: value.denom,
          amount: value.amount
        },
        destinationAddress: destination_address,
        fromAddress: from_address
      };
    }
  },
  "/quicksilver.interchainstaking.v1.MsgSignalIntent": {
    aminoType: "quicksilver/MsgSignalIntent",
    toAmino: ({
      chainId,
      intents,
      fromAddress
    }: MsgSignalIntent): AminoMsgSignalIntent["value"] => {
      return {
        chain_id: chainId,
        intents: intents.map(el0 => ({
          valoper_address: el0.valoperAddress,
          weight: el0.weight
        })),
        from_address: fromAddress
      };
    },
    fromAmino: ({
      chain_id,
      intents,
      from_address
    }: AminoMsgSignalIntent["value"]): MsgSignalIntent => {
      return {
        chainId: chain_id,
        intents: intents.map(el0 => ({
          valoperAddress: el0.valoper_address,
          weight: el0.weight
        })),
        fromAddress: from_address
      };
    }
  }
};