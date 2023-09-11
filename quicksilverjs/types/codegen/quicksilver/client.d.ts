import { OfflineSigner, GeneratedType, Registry } from "@cosmjs/proto-signing";
import { AminoTypes, SigningStargateClient } from "@cosmjs/stargate";
export declare const quicksilverAminoConverters: {
    "/quicksilver.tokenfactory.v1beta1.MsgCreateDenom": {
        aminoType: string;
        toAmino: ({ sender, subdenom }: import("./tokenfactory/v1beta1/tx").MsgCreateDenom) => {
            sender: string;
            subdenom: string;
        };
        fromAmino: ({ sender, subdenom }: {
            sender: string;
            subdenom: string;
        }) => import("./tokenfactory/v1beta1/tx").MsgCreateDenom;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgMint": {
        aminoType: string;
        toAmino: ({ sender, amount }: import("./tokenfactory/v1beta1/tx").MsgMint) => {
            sender: string;
            amount: {
                denom: string;
                amount: string;
            };
        };
        fromAmino: ({ sender, amount }: {
            sender: string;
            amount: {
                denom: string;
                amount: string;
            };
        }) => import("./tokenfactory/v1beta1/tx").MsgMint;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgBurn": {
        aminoType: string;
        toAmino: ({ sender, amount }: import("./tokenfactory/v1beta1/tx").MsgBurn) => {
            sender: string;
            amount: {
                denom: string;
                amount: string;
            };
        };
        fromAmino: ({ sender, amount }: {
            sender: string;
            amount: {
                denom: string;
                amount: string;
            };
        }) => import("./tokenfactory/v1beta1/tx").MsgBurn;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgChangeAdmin": {
        aminoType: string;
        toAmino: ({ sender, denom, newAdmin }: import("./tokenfactory/v1beta1/tx").MsgChangeAdmin) => {
            sender: string;
            denom: string;
            new_admin: string;
        };
        fromAmino: ({ sender, denom, new_admin }: {
            sender: string;
            denom: string;
            new_admin: string;
        }) => import("./tokenfactory/v1beta1/tx").MsgChangeAdmin;
    };
    "/quicksilver.tokenfactory.v1beta1.MsgSetDenomMetadata": {
        aminoType: string;
        toAmino: ({ sender, metadata }: import("./tokenfactory/v1beta1/tx").MsgSetDenomMetadata) => {
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
        fromAmino: ({ sender, metadata }: {
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
        }) => import("./tokenfactory/v1beta1/tx").MsgSetDenomMetadata;
    };
    "/quicksilver.participationrewards.v1.MsgSubmitClaim": {
        aminoType: string;
        toAmino: ({ userAddress, zone, srcZone, claimType, proofs }: import("./participationrewards/v1/messages").MsgSubmitClaim) => {
            user_address: string;
            zone: string;
            src_zone: string;
            claim_type: number;
            proofs: {
                key: Uint8Array;
                data: Uint8Array;
                proof_ops: {
                    ops: {
                        type: string;
                        key: Uint8Array;
                        data: Uint8Array;
                    }[];
                };
                height: string;
                proof_type: string;
            }[];
        };
        fromAmino: ({ user_address, zone, src_zone, claim_type, proofs }: {
            user_address: string;
            zone: string;
            src_zone: string;
            claim_type: number;
            proofs: {
                key: Uint8Array;
                data: Uint8Array;
                proof_ops: {
                    ops: {
                        type: string;
                        key: Uint8Array;
                        data: Uint8Array;
                    }[];
                };
                height: string;
                proof_type: string;
            }[];
        }) => import("./participationrewards/v1/messages").MsgSubmitClaim;
    };
    "/quicksilver.interchainstaking.v1.MsgRequestRedemption": {
        aminoType: string;
        toAmino: ({ value, destinationAddress, fromAddress }: import("./interchainstaking/v1/messages").MsgRequestRedemption) => {
            value: {
                denom: string;
                amount: string;
            };
            destination_address: string;
            from_address: string;
        };
        fromAmino: ({ value, destination_address, from_address }: {
            value: {
                denom: string;
                amount: string;
            };
            destination_address: string;
            from_address: string;
        }) => import("./interchainstaking/v1/messages").MsgRequestRedemption;
    };
    "/quicksilver.interchainstaking.v1.MsgSignalIntent": {
        aminoType: string;
        toAmino: ({ chainId, intents, fromAddress }: import("./interchainstaking/v1/messages").MsgSignalIntent) => {
            chain_id: string;
            intents: {
                valoper_address: string;
                weight: string;
            }[];
            from_address: string;
        };
        fromAmino: ({ chain_id, intents, from_address }: {
            chain_id: string;
            intents: {
                valoper_address: string;
                weight: string;
            }[];
            from_address: string;
        }) => import("./interchainstaking/v1/messages").MsgSignalIntent;
    };
    "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse": {
        aminoType: string;
        toAmino: ({ chainId, queryId, result, proofOps, height, fromAddress }: import("./interchainquery/v1/messages").MsgSubmitQueryResponse) => {
            chain_id: string;
            query_id: string;
            result: Uint8Array;
            proof_ops: {
                ops: {
                    type: string;
                    key: Uint8Array;
                    data: Uint8Array;
                }[];
            };
            height: string;
            from_address: string;
        };
        fromAmino: ({ chain_id, query_id, result, proof_ops, height, from_address }: {
            chain_id: string;
            query_id: string;
            result: Uint8Array;
            proof_ops: {
                ops: {
                    type: string;
                    key: Uint8Array;
                    data: Uint8Array;
                }[];
            };
            height: string;
            from_address: string;
        }) => import("./interchainquery/v1/messages").MsgSubmitQueryResponse;
    };
    "/quicksilver.airdrop.v1.MsgClaim": {
        aminoType: string;
        toAmino: ({ chainId, action, address, proofs }: import("./airdrop/v1/messages").MsgClaim) => {
            chain_id: string;
            action: string;
            address: string;
            proofs: {
                key: Uint8Array;
                data: Uint8Array;
                proof_ops: {
                    ops: {
                        type: string;
                        key: Uint8Array;
                        data: Uint8Array;
                    }[];
                };
                height: string;
            }[];
        };
        fromAmino: ({ chain_id, action, address, proofs }: {
            chain_id: string;
            action: string;
            address: string;
            proofs: {
                key: Uint8Array;
                data: Uint8Array;
                proof_ops: {
                    ops: {
                        type: string;
                        key: Uint8Array;
                        data: Uint8Array;
                    }[];
                };
                height: string;
            }[];
        }) => import("./airdrop/v1/messages").MsgClaim;
    };
};
export declare const quicksilverProtoRegistry: ReadonlyArray<[string, GeneratedType]>;
export declare const getSigningQuicksilverClientOptions: ({ defaultTypes }?: {
    defaultTypes?: ReadonlyArray<[string, GeneratedType]>;
}) => {
    registry: Registry;
    aminoTypes: AminoTypes;
};
export declare const getSigningQuicksilverClient: ({ rpcEndpoint, signer, defaultTypes }: {
    rpcEndpoint: string;
    signer: OfflineSigner;
    defaultTypes?: ReadonlyArray<[string, GeneratedType]>;
}) => Promise<SigningStargateClient>;
