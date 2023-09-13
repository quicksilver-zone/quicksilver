import { Rpc } from "../helpers";
export declare const createRPCMsgClient: ({ rpc }: {
    rpc: Rpc;
}) => Promise<{
    cosmos: {
        authz: {
            v1beta1: import("../cosmos/authz/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        bank: {
            v1beta1: import("../cosmos/bank/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        crisis: {
            v1beta1: import("../cosmos/crisis/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        distribution: {
            v1beta1: import("../cosmos/distribution/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        evidence: {
            v1beta1: import("../cosmos/evidence/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        feegrant: {
            v1beta1: import("../cosmos/feegrant/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        gov: {
            v1beta1: import("../cosmos/gov/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        slashing: {
            v1beta1: import("../cosmos/slashing/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        staking: {
            v1beta1: import("../cosmos/staking/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
        vesting: {
            v1beta1: import("../cosmos/vesting/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
    };
    quicksilver: {
        airdrop: {
            v1: import("./airdrop/v1/messages.rpc.msg").MsgClientImpl;
        };
        interchainquery: {
            v1: import("./interchainquery/v1/messages.rpc.msg").MsgClientImpl;
        };
        interchainstaking: {
            v1: import("./interchainstaking/v1/messages.rpc.msg").MsgClientImpl;
        };
        participationrewards: {
            v1: import("./participationrewards/v1/messages.rpc.msg").MsgClientImpl;
        };
        tokenfactory: {
            v1beta1: import("./tokenfactory/v1beta1/tx.rpc.msg").MsgClientImpl;
        };
    };
}>;
