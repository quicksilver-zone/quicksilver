import { AminoMsg } from "@cosmjs/amino";
import { MsgCreateDenom, MsgMint, MsgBurn, MsgChangeAdmin, MsgSetDenomMetadata } from "./tx";
export interface AminoMsgCreateDenom extends AminoMsg {
    type: "/quicksilver.tokenfactory.v1beta1.MsgCreateDenom";
    value: {
        sender: string;
        subdenom: string;
    };
}
export interface AminoMsgMint extends AminoMsg {
    type: "/quicksilver.tokenfactory.v1beta1.MsgMint";
    value: {
        sender: string;
        amount: {
            denom: string;
            amount: string;
        };
    };
}
export interface AminoMsgBurn extends AminoMsg {
    type: "/quicksilver.tokenfactory.v1beta1.MsgBurn";
    value: {
        sender: string;
        amount: {
            denom: string;
            amount: string;
        };
    };
}
export interface AminoMsgChangeAdmin extends AminoMsg {
    type: "/quicksilver.tokenfactory.v1beta1.MsgChangeAdmin";
    value: {
        sender: string;
        denom: string;
        new_admin: string;
    };
}
export interface AminoMsgSetDenomMetadata extends AminoMsg {
    type: "/quicksilver.tokenfactory.v1beta1.MsgSetDenomMetadata";
    value: {
        sender: string;
        metadata: {
            description: string;
            denom_units: {
                denom: string;
                exponent: number;
                aliases: string[];
            }[];
            base: string;
            display: string;
            name: string;
            symbol: string;
        };
    };
}
export declare const AminoConverter: {
    "/quicksilver.tokenfactory.v1beta1.MsgCreateDenom": {
        aminoType: string;
        toAmino: ({ sender, subdenom }: MsgCreateDenom) => AminoMsgCreateDenom["value"];
        fromAmino: ({ sender, subdenom }: AminoMsgCreateDenom["value"]) => MsgCreateDenom;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgMint": {
        aminoType: string;
        toAmino: ({ sender, amount }: MsgMint) => AminoMsgMint["value"];
        fromAmino: ({ sender, amount }: AminoMsgMint["value"]) => MsgMint;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgBurn": {
        aminoType: string;
        toAmino: ({ sender, amount }: MsgBurn) => AminoMsgBurn["value"];
        fromAmino: ({ sender, amount }: AminoMsgBurn["value"]) => MsgBurn;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgChangeAdmin": {
        aminoType: string;
        toAmino: ({ sender, denom, newAdmin }: MsgChangeAdmin) => AminoMsgChangeAdmin["value"];
        fromAmino: ({ sender, denom, new_admin }: AminoMsgChangeAdmin["value"]) => MsgChangeAdmin;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgSetDenomMetadata": {
        aminoType: string;
        toAmino: ({ sender, metadata }: MsgSetDenomMetadata) => AminoMsgSetDenomMetadata["value"];
        fromAmino: ({ sender, metadata }: AminoMsgSetDenomMetadata["value"]) => MsgSetDenomMetadata;
    };
};
