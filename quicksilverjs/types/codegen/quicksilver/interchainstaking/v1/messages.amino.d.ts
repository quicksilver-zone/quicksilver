import { AminoMsg } from "@cosmjs/amino";
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
export declare const AminoConverter: {
    "/quicksilver.interchainstaking.v1.MsgRequestRedemption": {
        aminoType: string;
        toAmino: ({ value, destinationAddress, fromAddress }: MsgRequestRedemption) => AminoMsgRequestRedemption["value"];
        fromAmino: ({ value, destination_address, from_address }: AminoMsgRequestRedemption["value"]) => MsgRequestRedemption;
    };
    "/quicksilver.interchainstaking.v1.MsgSignalIntent": {
        aminoType: string;
        toAmino: ({ chainId, intents, fromAddress }: MsgSignalIntent) => AminoMsgSignalIntent["value"];
        fromAmino: ({ chain_id, intents, from_address }: AminoMsgSignalIntent["value"]) => MsgSignalIntent;
    };
};
