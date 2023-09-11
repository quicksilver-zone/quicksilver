import { AminoMsg } from "@cosmjs/amino";
import { MsgCreateVestingAccount } from "./tx";
export interface AminoMsgCreateVestingAccount extends AminoMsg {
    type: "cosmos-sdk/MsgCreateVestingAccount";
    value: {
        from_address: string;
        to_address: string;
        amount: {
            denom: string;
            amount: string;
        }[];
        end_time: string;
        delayed: boolean;
    };
}
export declare const AminoConverter: {
    "/cosmos.vesting.v1beta1.MsgCreateVestingAccount": {
        aminoType: string;
        toAmino: ({ fromAddress, toAddress, amount, endTime, delayed }: MsgCreateVestingAccount) => AminoMsgCreateVestingAccount["value"];
        fromAmino: ({ from_address, to_address, amount, end_time, delayed }: AminoMsgCreateVestingAccount["value"]) => MsgCreateVestingAccount;
    };
};
