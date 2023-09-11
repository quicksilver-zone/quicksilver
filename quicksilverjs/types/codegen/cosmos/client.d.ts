import { OfflineSigner, GeneratedType, Registry } from "@cosmjs/proto-signing";
import { AminoTypes, SigningStargateClient } from "@cosmjs/stargate";
export declare const cosmosAminoConverters: {
    "/cosmos.vesting.v1beta1.MsgCreateVestingAccount": {
        aminoType: string;
        toAmino: ({ fromAddress, toAddress, amount, endTime, delayed }: import("./vesting/v1beta1/tx").MsgCreateVestingAccount) => {
            from_address: string;
            to_address: string;
            amount: {
                denom: string;
                amount: string;
            }[];
            end_time: string;
            delayed: boolean;
        };
        fromAmino: ({ from_address, to_address, amount, end_time, delayed }: {
            from_address: string;
            to_address: string;
            amount: {
                denom: string;
                amount: string;
            }[];
            end_time: string;
            delayed: boolean;
        }) => import("./vesting/v1beta1/tx").MsgCreateVestingAccount;
    };
    "/cosmos.staking.v1beta1.MsgCreateValidator": {
        aminoType: string;
        toAmino: ({ description, commission, minSelfDelegation, delegatorAddress, validatorAddress, pubkey, value }: import("./staking/v1beta1/tx").MsgCreateValidator) => {
            description: {
                moniker: string;
                identity: string;
                website: string;
                security_contact: string;
                details: string;
            };
            commission: {
                rate: string;
                max_rate: string;
                max_change_rate: string;
            };
            min_self_delegation: string;
            delegator_address: string;
            validator_address: string;
            pubkey: {
                type_url: string;
                value: Uint8Array;
            };
            value: {
                denom: string;
                amount: string;
            };
        };
        fromAmino: ({ description, commission, min_self_delegation, delegator_address, validator_address, pubkey, value }: {
            description: {
                moniker: string;
                identity: string;
                website: string;
                security_contact: string;
                details: string;
            };
            commission: {
                rate: string;
                max_rate: string;
                max_change_rate: string;
            };
            min_self_delegation: string;
            delegator_address: string;
            validator_address: string;
            pubkey: {
                type_url: string;
                value: Uint8Array;
            };
            value: {
                denom: string;
                amount: string;
            };
        }) => import("./staking/v1beta1/tx").MsgCreateValidator;
    };
    "/cosmos.staking.v1beta1.MsgEditValidator": {
        aminoType: string;
        toAmino: ({ description, validatorAddress, commissionRate, minSelfDelegation }: import("./staking/v1beta1/tx").MsgEditValidator) => {
            description: {
                moniker: string;
                identity: string;
                website: string;
                security_contact: string;
                details: string;
            };
            validator_address: string;
            commission_rate: string;
            min_self_delegation: string;
        };
        fromAmino: ({ description, validator_address, commission_rate, min_self_delegation }: {
            description: {
                moniker: string;
                identity: string;
                website: string;
                security_contact: string;
                details: string;
            };
            validator_address: string;
            commission_rate: string;
            min_self_delegation: string;
        }) => import("./staking/v1beta1/tx").MsgEditValidator;
    };
    "/cosmos.staking.v1beta1.MsgDelegate": {
        aminoType: string;
        toAmino: ({ delegatorAddress, validatorAddress, amount }: import("./staking/v1beta1/tx").MsgDelegate) => {
            delegator_address: string;
            validator_address: string;
            amount: {
                denom: string;
                amount: string;
            };
        };
        fromAmino: ({ delegator_address, validator_address, amount }: {
            delegator_address: string;
            validator_address: string;
            amount: {
                denom: string;
                amount: string;
            };
        }) => import("./staking/v1beta1/tx").MsgDelegate;
    };
    "/cosmos.staking.v1beta1.MsgBeginRedelegate": {
        aminoType: string;
        toAmino: ({ delegatorAddress, validatorSrcAddress, validatorDstAddress, amount }: import("./staking/v1beta1/tx").MsgBeginRedelegate) => {
            delegator_address: string;
            validator_src_address: string;
            validator_dst_address: string;
            amount: {
                denom: string;
                amount: string;
            };
        };
        fromAmino: ({ delegator_address, validator_src_address, validator_dst_address, amount }: {
            delegator_address: string;
            validator_src_address: string;
            validator_dst_address: string;
            amount: {
                denom: string;
                amount: string;
            };
        }) => import("./staking/v1beta1/tx").MsgBeginRedelegate;
    };
    "/cosmos.staking.v1beta1.MsgUndelegate": {
        aminoType: string;
        toAmino: ({ delegatorAddress, validatorAddress, amount }: import("./staking/v1beta1/tx").MsgUndelegate) => {
            delegator_address: string;
            validator_address: string;
            amount: {
                denom: string;
                amount: string;
            };
        };
        fromAmino: ({ delegator_address, validator_address, amount }: {
            delegator_address: string;
            validator_address: string;
            amount: {
                denom: string;
                amount: string;
            };
        }) => import("./staking/v1beta1/tx").MsgUndelegate;
    };
    "/cosmos.slashing.v1beta1.MsgUnjail": {
        aminoType: string;
        toAmino: ({ validatorAddr }: import("./slashing/v1beta1/tx").MsgUnjail) => {
            validator_addr: string;
        };
        fromAmino: ({ validator_addr }: {
            validator_addr: string;
        }) => import("./slashing/v1beta1/tx").MsgUnjail;
    };
    "/cosmos.gov.v1beta1.MsgSubmitProposal": {
        aminoType: string;
        toAmino: ({ content, initialDeposit, proposer }: import("./gov/v1beta1/tx").MsgSubmitProposal) => {
            content: {
                type_url: string;
                value: Uint8Array;
            };
            initial_deposit: {
                denom: string;
                amount: string;
            }[];
            proposer: string;
        };
        fromAmino: ({ content, initial_deposit, proposer }: {
            content: {
                type_url: string;
                value: Uint8Array;
            };
            initial_deposit: {
                denom: string;
                amount: string;
            }[];
            proposer: string;
        }) => import("./gov/v1beta1/tx").MsgSubmitProposal;
    };
    "/cosmos.gov.v1beta1.MsgVote": {
        aminoType: string;
        toAmino: ({ proposalId, voter, option }: import("./gov/v1beta1/tx").MsgVote) => {
            proposal_id: string;
            voter: string;
            option: number;
        };
        fromAmino: ({ proposal_id, voter, option }: {
            proposal_id: string;
            voter: string;
            option: number;
        }) => import("./gov/v1beta1/tx").MsgVote;
    };
    "/cosmos.gov.v1beta1.MsgVoteWeighted": {
        aminoType: string;
        toAmino: ({ proposalId, voter, options }: import("./gov/v1beta1/tx").MsgVoteWeighted) => {
            proposal_id: string;
            voter: string;
            options: {
                option: number;
                weight: string;
            }[];
        };
        fromAmino: ({ proposal_id, voter, options }: {
            proposal_id: string;
            voter: string;
            options: {
                option: number;
                weight: string;
            }[];
        }) => import("./gov/v1beta1/tx").MsgVoteWeighted;
    };
    "/cosmos.gov.v1beta1.MsgDeposit": {
        aminoType: string;
        toAmino: ({ proposalId, depositor, amount }: import("./gov/v1beta1/tx").MsgDeposit) => {
            proposal_id: string;
            depositor: string;
            amount: {
                denom: string;
                amount: string;
            }[];
        };
        fromAmino: ({ proposal_id, depositor, amount }: {
            proposal_id: string;
            depositor: string;
            amount: {
                denom: string;
                amount: string;
            }[];
        }) => import("./gov/v1beta1/tx").MsgDeposit;
    };
    "/cosmos.feegrant.v1beta1.MsgGrantAllowance": {
        aminoType: string;
        toAmino: ({ granter, grantee, allowance }: import("./feegrant/v1beta1/tx").MsgGrantAllowance) => {
            granter: string;
            grantee: string;
            allowance: {
                type_url: string;
                value: Uint8Array;
            };
        };
        fromAmino: ({ granter, grantee, allowance }: {
            granter: string;
            grantee: string;
            allowance: {
                type_url: string;
                value: Uint8Array;
            };
        }) => import("./feegrant/v1beta1/tx").MsgGrantAllowance;
    };
    "/cosmos.feegrant.v1beta1.MsgRevokeAllowance": {
        aminoType: string;
        toAmino: ({ granter, grantee }: import("./feegrant/v1beta1/tx").MsgRevokeAllowance) => {
            granter: string;
            grantee: string;
        };
        fromAmino: ({ granter, grantee }: {
            granter: string;
            grantee: string;
        }) => import("./feegrant/v1beta1/tx").MsgRevokeAllowance;
    };
    "/cosmos.evidence.v1beta1.MsgSubmitEvidence": {
        aminoType: string;
        toAmino: ({ submitter, evidence }: import("./evidence/v1beta1/tx").MsgSubmitEvidence) => {
            submitter: string;
            evidence: {
                type_url: string;
                value: Uint8Array;
            };
        };
        fromAmino: ({ submitter, evidence }: {
            submitter: string;
            evidence: {
                type_url: string;
                value: Uint8Array;
            };
        }) => import("./evidence/v1beta1/tx").MsgSubmitEvidence;
    };
    "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress": {
        aminoType: string;
        toAmino: ({ delegatorAddress, withdrawAddress }: import("./distribution/v1beta1/tx").MsgSetWithdrawAddress) => {
            delegator_address: string;
            withdraw_address: string;
        };
        fromAmino: ({ delegator_address, withdraw_address }: {
            delegator_address: string;
            withdraw_address: string;
        }) => import("./distribution/v1beta1/tx").MsgSetWithdrawAddress;
    };
    "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward": {
        aminoType: string;
        toAmino: ({ delegatorAddress, validatorAddress }: import("./distribution/v1beta1/tx").MsgWithdrawDelegatorReward) => {
            delegator_address: string;
            validator_address: string;
        };
        fromAmino: ({ delegator_address, validator_address }: {
            delegator_address: string;
            validator_address: string;
        }) => import("./distribution/v1beta1/tx").MsgWithdrawDelegatorReward;
    };
    "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission": {
        aminoType: string;
        toAmino: ({ validatorAddress }: import("./distribution/v1beta1/tx").MsgWithdrawValidatorCommission) => {
            validator_address: string;
        };
        fromAmino: ({ validator_address }: {
            validator_address: string;
        }) => import("./distribution/v1beta1/tx").MsgWithdrawValidatorCommission;
    };
    "/cosmos.distribution.v1beta1.MsgFundCommunityPool": {
        aminoType: string;
        toAmino: ({ amount, depositor }: import("./distribution/v1beta1/tx").MsgFundCommunityPool) => {
            amount: {
                denom: string;
                amount: string;
            }[];
            depositor: string;
        };
        fromAmino: ({ amount, depositor }: {
            amount: {
                denom: string;
                amount: string;
            }[];
            depositor: string;
        }) => import("./distribution/v1beta1/tx").MsgFundCommunityPool;
    };
    "/cosmos.crisis.v1beta1.MsgVerifyInvariant": {
        aminoType: string;
        toAmino: ({ sender, invariantModuleName, invariantRoute }: import("./crisis/v1beta1/tx").MsgVerifyInvariant) => {
            sender: string;
            invariant_module_name: string;
            invariant_route: string;
        };
        fromAmino: ({ sender, invariant_module_name, invariant_route }: {
            sender: string;
            invariant_module_name: string;
            invariant_route: string;
        }) => import("./crisis/v1beta1/tx").MsgVerifyInvariant;
    };
    "/cosmos.bank.v1beta1.MsgSend": {
        aminoType: string;
        toAmino: ({ fromAddress, toAddress, amount }: import("./bank/v1beta1/tx").MsgSend) => {
            from_address: string;
            to_address: string;
            amount: {
                denom: string;
                amount: string;
            }[];
        };
        fromAmino: ({ from_address, to_address, amount }: {
            from_address: string;
            to_address: string;
            amount: {
                denom: string;
                amount: string;
            }[];
        }) => import("./bank/v1beta1/tx").MsgSend;
    };
    "/cosmos.bank.v1beta1.MsgMultiSend": {
        aminoType: string;
        toAmino: ({ inputs, outputs }: import("./bank/v1beta1/tx").MsgMultiSend) => {
            inputs: {
                address: string;
                coins: {
                    denom: string;
                    amount: string;
                }[];
            }[];
            outputs: {
                address: string;
                coins: {
                    denom: string;
                    amount: string;
                }[];
            }[];
        };
        fromAmino: ({ inputs, outputs }: {
            inputs: {
                address: string;
                coins: {
                    denom: string;
                    amount: string;
                }[];
            }[];
            outputs: {
                address: string;
                coins: {
                    denom: string;
                    amount: string;
                }[];
            }[];
        }) => import("./bank/v1beta1/tx").MsgMultiSend;
    };
    "/cosmos.authz.v1beta1.MsgGrant": {
        aminoType: string;
        toAmino: ({ granter, grantee, grant }: import("./authz/v1beta1/tx").MsgGrant) => {
            granter: string;
            grantee: string;
            grant: {
                authorization: {
                    type_url: string;
                    value: Uint8Array;
                };
                expiration: {
                    seconds: string;
                    nanos: number;
                };
            };
        };
        fromAmino: ({ granter, grantee, grant }: {
            granter: string;
            grantee: string;
            grant: {
                authorization: {
                    type_url: string;
                    value: Uint8Array;
                };
                expiration: {
                    seconds: string;
                    nanos: number;
                };
            };
        }) => import("./authz/v1beta1/tx").MsgGrant;
    };
    "/cosmos.authz.v1beta1.MsgExec": {
        aminoType: string;
        toAmino: ({ grantee, msgs }: import("./authz/v1beta1/tx").MsgExec) => {
            grantee: string;
            msgs: {
                type_url: string;
                value: Uint8Array;
            }[];
        };
        fromAmino: ({ grantee, msgs }: {
            grantee: string;
            msgs: {
                type_url: string;
                value: Uint8Array;
            }[];
        }) => import("./authz/v1beta1/tx").MsgExec;
    };
    "/cosmos.authz.v1beta1.MsgRevoke": {
        aminoType: string;
        toAmino: ({ granter, grantee, msgTypeUrl }: import("./authz/v1beta1/tx").MsgRevoke) => {
            granter: string;
            grantee: string;
            msg_type_url: string;
        };
        fromAmino: ({ granter, grantee, msg_type_url }: {
            granter: string;
            grantee: string;
            msg_type_url: string;
        }) => import("./authz/v1beta1/tx").MsgRevoke;
    };
};
export declare const cosmosProtoRegistry: ReadonlyArray<[string, GeneratedType]>;
export declare const getSigningCosmosClientOptions: () => {
    registry: Registry;
    aminoTypes: AminoTypes;
};
export declare const getSigningCosmosClient: ({ rpcEndpoint, signer }: {
    rpcEndpoint: string;
    signer: OfflineSigner;
}) => Promise<SigningStargateClient>;
