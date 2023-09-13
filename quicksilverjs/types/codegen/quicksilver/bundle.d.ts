import * as _83 from "./airdrop/v1/airdrop";
import * as _84 from "./airdrop/v1/genesis";
import * as _85 from "./airdrop/v1/messages";
import * as _86 from "./airdrop/v1/params";
import * as _87 from "./airdrop/v1/proposals";
import * as _88 from "./airdrop/v1/query";
import * as _89 from "./claimsmanager/v1/claimsmanager";
import * as _90 from "./claimsmanager/v1/genesis";
import * as _92 from "./claimsmanager/v1/query";
import * as _93 from "./epochs/v1/genesis";
import * as _94 from "./epochs/v1/query";
import * as _95 from "./interchainquery/v1/genesis";
import * as _96 from "./interchainquery/v1/interchainquery";
import * as _97 from "./interchainquery/v1/messages";
import * as _98 from "./interchainquery/v1/query";
import * as _99 from "./interchainstaking/v1/genesis";
import * as _100 from "./interchainstaking/v1/interchainstaking";
import * as _101 from "./interchainstaking/v1/messages";
import * as _102 from "./interchainstaking/v1/proposals";
import * as _103 from "./interchainstaking/v1/query";
import * as _104 from "./participationrewards/v1/proposals";
import * as _105 from "./mint/v1beta1/genesis";
import * as _106 from "./mint/v1beta1/mint";
import * as _107 from "./mint/v1beta1/query";
import * as _108 from "./participationrewards/v1/genesis";
import * as _109 from "./participationrewards/v1/messages";
import * as _110 from "./participationrewards/v1/participationrewards";
import * as _111 from "./participationrewards/v1/query";
import * as _112 from "./tokenfactory/v1beta1/authorityMetadata";
import * as _113 from "./tokenfactory/v1beta1/genesis";
import * as _114 from "./tokenfactory/v1beta1/params";
import * as _115 from "./tokenfactory/v1beta1/query";
import * as _116 from "./tokenfactory/v1beta1/tx";
import * as _182 from "./airdrop/v1/query.rpc.Query";
import * as _183 from "./claimsmanager/v1/query.rpc.Query";
import * as _184 from "./epochs/v1/query.rpc.Query";
import * as _185 from "./interchainstaking/v1/query.rpc.Query";
import * as _186 from "./mint/v1beta1/query.rpc.Query";
import * as _187 from "./participationrewards/v1/query.rpc.Query";
import * as _188 from "./tokenfactory/v1beta1/query.rpc.Query";
import * as _189 from "./airdrop/v1/messages.rpc.msg";
import * as _190 from "./interchainquery/v1/messages.rpc.msg";
import * as _191 from "./interchainstaking/v1/messages.rpc.msg";
import * as _192 from "./participationrewards/v1/messages.rpc.msg";
import * as _193 from "./tokenfactory/v1beta1/tx.rpc.msg";
export declare namespace quicksilver {
    namespace airdrop {
        const v1: {
            MsgClientImpl: typeof _189.MsgClientImpl;
            QueryClientImpl: typeof _182.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _88.QueryParamsRequest): Promise<_88.QueryParamsResponse>;
                zoneDrop(request: _88.QueryZoneDropRequest): Promise<_88.QueryZoneDropResponse>;
                accountBalance(request: _88.QueryAccountBalanceRequest): Promise<_88.QueryAccountBalanceResponse>;
                zoneDrops(request: _88.QueryZoneDropsRequest): Promise<_88.QueryZoneDropsResponse>;
                claimRecord(request: _88.QueryClaimRecordRequest): Promise<_88.QueryClaimRecordResponse>;
                claimRecords(request: _88.QueryClaimRecordsRequest): Promise<_88.QueryClaimRecordsResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    claim(value: _85.MsgClaim): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    claim(value: _85.MsgClaim): {
                        typeUrl: string;
                        value: _85.MsgClaim;
                    };
                };
                toJSON: {
                    claim(value: _85.MsgClaim): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    claim(value: any): {
                        typeUrl: string;
                        value: _85.MsgClaim;
                    };
                };
                fromPartial: {
                    claim(value: _85.MsgClaim): {
                        typeUrl: string;
                        value: _85.MsgClaim;
                    };
                };
            };
            AminoConverter: {
                "/quicksilver.airdrop.v1.MsgClaim": {
                    aminoType: string;
                    toAmino: ({ chainId, action, address, proofs }: _85.MsgClaim) => {
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
                    }) => _85.MsgClaim;
                };
            };
            QueryParamsRequest: {
                encode(_: _88.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryParamsRequest;
                fromJSON(_: any): _88.QueryParamsRequest;
                toJSON(_: _88.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_88.QueryParamsRequest>): _88.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _88.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryParamsResponse;
                fromJSON(object: any): _88.QueryParamsResponse;
                toJSON(message: _88.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_88.QueryParamsResponse>): _88.QueryParamsResponse;
            };
            QueryZoneDropRequest: {
                encode(message: _88.QueryZoneDropRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryZoneDropRequest;
                fromJSON(object: any): _88.QueryZoneDropRequest;
                toJSON(message: _88.QueryZoneDropRequest): unknown;
                fromPartial(object: Partial<_88.QueryZoneDropRequest>): _88.QueryZoneDropRequest;
            };
            QueryZoneDropResponse: {
                encode(message: _88.QueryZoneDropResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryZoneDropResponse;
                fromJSON(object: any): _88.QueryZoneDropResponse;
                toJSON(message: _88.QueryZoneDropResponse): unknown;
                fromPartial(object: Partial<_88.QueryZoneDropResponse>): _88.QueryZoneDropResponse;
            };
            QueryAccountBalanceRequest: {
                encode(message: _88.QueryAccountBalanceRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryAccountBalanceRequest;
                fromJSON(object: any): _88.QueryAccountBalanceRequest;
                toJSON(message: _88.QueryAccountBalanceRequest): unknown;
                fromPartial(object: Partial<_88.QueryAccountBalanceRequest>): _88.QueryAccountBalanceRequest;
            };
            QueryAccountBalanceResponse: {
                encode(message: _88.QueryAccountBalanceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryAccountBalanceResponse;
                fromJSON(object: any): _88.QueryAccountBalanceResponse;
                toJSON(message: _88.QueryAccountBalanceResponse): unknown;
                fromPartial(object: Partial<_88.QueryAccountBalanceResponse>): _88.QueryAccountBalanceResponse;
            };
            QueryZoneDropsRequest: {
                encode(message: _88.QueryZoneDropsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryZoneDropsRequest;
                fromJSON(object: any): _88.QueryZoneDropsRequest;
                toJSON(message: _88.QueryZoneDropsRequest): unknown;
                fromPartial(object: Partial<_88.QueryZoneDropsRequest>): _88.QueryZoneDropsRequest;
            };
            QueryZoneDropsResponse: {
                encode(message: _88.QueryZoneDropsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryZoneDropsResponse;
                fromJSON(object: any): _88.QueryZoneDropsResponse;
                toJSON(message: _88.QueryZoneDropsResponse): unknown;
                fromPartial(object: Partial<_88.QueryZoneDropsResponse>): _88.QueryZoneDropsResponse;
            };
            QueryClaimRecordRequest: {
                encode(message: _88.QueryClaimRecordRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryClaimRecordRequest;
                fromJSON(object: any): _88.QueryClaimRecordRequest;
                toJSON(message: _88.QueryClaimRecordRequest): unknown;
                fromPartial(object: Partial<_88.QueryClaimRecordRequest>): _88.QueryClaimRecordRequest;
            };
            QueryClaimRecordResponse: {
                encode(message: _88.QueryClaimRecordResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryClaimRecordResponse;
                fromJSON(object: any): _88.QueryClaimRecordResponse;
                toJSON(message: _88.QueryClaimRecordResponse): unknown;
                fromPartial(object: Partial<_88.QueryClaimRecordResponse>): _88.QueryClaimRecordResponse;
            };
            QueryClaimRecordsRequest: {
                encode(message: _88.QueryClaimRecordsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryClaimRecordsRequest;
                fromJSON(object: any): _88.QueryClaimRecordsRequest;
                toJSON(message: _88.QueryClaimRecordsRequest): unknown;
                fromPartial(object: Partial<_88.QueryClaimRecordsRequest>): _88.QueryClaimRecordsRequest;
            };
            QueryClaimRecordsResponse: {
                encode(message: _88.QueryClaimRecordsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _88.QueryClaimRecordsResponse;
                fromJSON(object: any): _88.QueryClaimRecordsResponse;
                toJSON(message: _88.QueryClaimRecordsResponse): unknown;
                fromPartial(object: Partial<_88.QueryClaimRecordsResponse>): _88.QueryClaimRecordsResponse;
            };
            RegisterZoneDropProposal: {
                encode(message: _87.RegisterZoneDropProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _87.RegisterZoneDropProposal;
                fromJSON(object: any): _87.RegisterZoneDropProposal;
                toJSON(message: _87.RegisterZoneDropProposal): unknown;
                fromPartial(object: Partial<_87.RegisterZoneDropProposal>): _87.RegisterZoneDropProposal;
            };
            Params: {
                encode(_: _86.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _86.Params;
                fromJSON(_: any): _86.Params;
                toJSON(_: _86.Params): unknown;
                fromPartial(_: Partial<_86.Params>): _86.Params;
            };
            MsgClaim: {
                encode(message: _85.MsgClaim, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _85.MsgClaim;
                fromJSON(object: any): _85.MsgClaim;
                toJSON(message: _85.MsgClaim): unknown;
                fromPartial(object: Partial<_85.MsgClaim>): _85.MsgClaim;
            };
            MsgClaimResponse: {
                encode(message: _85.MsgClaimResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _85.MsgClaimResponse;
                fromJSON(object: any): _85.MsgClaimResponse;
                toJSON(message: _85.MsgClaimResponse): unknown;
                fromPartial(object: Partial<_85.MsgClaimResponse>): _85.MsgClaimResponse;
            };
            Proof: {
                encode(message: _85.Proof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _85.Proof;
                fromJSON(object: any): _85.Proof;
                toJSON(message: _85.Proof): unknown;
                fromPartial(object: Partial<_85.Proof>): _85.Proof;
            };
            GenesisState: {
                encode(message: _84.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _84.GenesisState;
                fromJSON(object: any): _84.GenesisState;
                toJSON(message: _84.GenesisState): unknown;
                fromPartial(object: Partial<_84.GenesisState>): _84.GenesisState;
            };
            actionFromJSON(object: any): _83.Action;
            actionToJSON(object: _83.Action): string;
            statusFromJSON(object: any): _83.Status;
            statusToJSON(object: _83.Status): string;
            Action: typeof _83.Action;
            ActionSDKType: typeof _83.ActionSDKType;
            Status: typeof _83.Status;
            StatusSDKType: typeof _83.StatusSDKType;
            ZoneDrop: {
                encode(message: _83.ZoneDrop, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _83.ZoneDrop;
                fromJSON(object: any): _83.ZoneDrop;
                toJSON(message: _83.ZoneDrop): unknown;
                fromPartial(object: Partial<_83.ZoneDrop>): _83.ZoneDrop;
            };
            ClaimRecord_ActionsCompletedEntry: {
                encode(message: _83.ClaimRecord_ActionsCompletedEntry, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _83.ClaimRecord_ActionsCompletedEntry;
                fromJSON(object: any): _83.ClaimRecord_ActionsCompletedEntry;
                toJSON(message: _83.ClaimRecord_ActionsCompletedEntry): unknown;
                fromPartial(object: Partial<_83.ClaimRecord_ActionsCompletedEntry>): _83.ClaimRecord_ActionsCompletedEntry;
            };
            ClaimRecord: {
                encode(message: _83.ClaimRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _83.ClaimRecord;
                fromJSON(object: any): _83.ClaimRecord;
                toJSON(message: _83.ClaimRecord): unknown;
                fromPartial(object: Partial<_83.ClaimRecord>): _83.ClaimRecord;
            };
            CompletedAction: {
                encode(message: _83.CompletedAction, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _83.CompletedAction;
                fromJSON(object: any): _83.CompletedAction;
                toJSON(message: _83.CompletedAction): unknown;
                fromPartial(object: Partial<_83.CompletedAction>): _83.CompletedAction;
            };
        };
    }
    namespace claimsmanager {
        const v1: {
            QueryClientImpl: typeof _183.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _92.QueryParamsRequest): Promise<_92.QueryParamsResponse>;
                claims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
                lastEpochClaims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
                userClaims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
                userLastEpochClaims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
            };
            QueryParamsRequest: {
                encode(_: _92.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _92.QueryParamsRequest;
                fromJSON(_: any): _92.QueryParamsRequest;
                toJSON(_: _92.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_92.QueryParamsRequest>): _92.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _92.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _92.QueryParamsResponse;
                fromJSON(object: any): _92.QueryParamsResponse;
                toJSON(message: _92.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_92.QueryParamsResponse>): _92.QueryParamsResponse;
            };
            QueryClaimsRequest: {
                encode(message: _92.QueryClaimsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _92.QueryClaimsRequest;
                fromJSON(object: any): _92.QueryClaimsRequest;
                toJSON(message: _92.QueryClaimsRequest): unknown;
                fromPartial(object: Partial<_92.QueryClaimsRequest>): _92.QueryClaimsRequest;
            };
            QueryClaimsResponse: {
                encode(message: _92.QueryClaimsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _92.QueryClaimsResponse;
                fromJSON(object: any): _92.QueryClaimsResponse;
                toJSON(message: _92.QueryClaimsResponse): unknown;
                fromPartial(object: Partial<_92.QueryClaimsResponse>): _92.QueryClaimsResponse;
            };
            GenesisState: {
                encode(message: _90.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _90.GenesisState;
                fromJSON(object: any): _90.GenesisState;
                toJSON(message: _90.GenesisState): unknown;
                fromPartial(object: Partial<_90.GenesisState>): _90.GenesisState;
            };
            claimTypeFromJSON(object: any): _89.ClaimType;
            claimTypeToJSON(object: _89.ClaimType): string;
            ClaimType: typeof _89.ClaimType;
            ClaimTypeSDKType: typeof _89.ClaimTypeSDKType;
            Params: {
                encode(_: _89.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _89.Params;
                fromJSON(_: any): _89.Params;
                toJSON(_: _89.Params): unknown;
                fromPartial(_: Partial<_89.Params>): _89.Params;
            };
            Claim: {
                encode(message: _89.Claim, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _89.Claim;
                fromJSON(object: any): _89.Claim;
                toJSON(message: _89.Claim): unknown;
                fromPartial(object: Partial<_89.Claim>): _89.Claim;
            };
        };
    }
    namespace epochs {
        const v1: {
            QueryClientImpl: typeof _184.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                epochInfos(request?: _94.QueryEpochsInfoRequest): Promise<_94.QueryEpochsInfoResponse>;
                currentEpoch(request: _94.QueryCurrentEpochRequest): Promise<_94.QueryCurrentEpochResponse>;
            };
            QueryEpochsInfoRequest: {
                encode(message: _94.QueryEpochsInfoRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _94.QueryEpochsInfoRequest;
                fromJSON(object: any): _94.QueryEpochsInfoRequest;
                toJSON(message: _94.QueryEpochsInfoRequest): unknown;
                fromPartial(object: Partial<_94.QueryEpochsInfoRequest>): _94.QueryEpochsInfoRequest;
            };
            QueryEpochsInfoResponse: {
                encode(message: _94.QueryEpochsInfoResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _94.QueryEpochsInfoResponse;
                fromJSON(object: any): _94.QueryEpochsInfoResponse;
                toJSON(message: _94.QueryEpochsInfoResponse): unknown;
                fromPartial(object: Partial<_94.QueryEpochsInfoResponse>): _94.QueryEpochsInfoResponse;
            };
            QueryCurrentEpochRequest: {
                encode(message: _94.QueryCurrentEpochRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _94.QueryCurrentEpochRequest;
                fromJSON(object: any): _94.QueryCurrentEpochRequest;
                toJSON(message: _94.QueryCurrentEpochRequest): unknown;
                fromPartial(object: Partial<_94.QueryCurrentEpochRequest>): _94.QueryCurrentEpochRequest;
            };
            QueryCurrentEpochResponse: {
                encode(message: _94.QueryCurrentEpochResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _94.QueryCurrentEpochResponse;
                fromJSON(object: any): _94.QueryCurrentEpochResponse;
                toJSON(message: _94.QueryCurrentEpochResponse): unknown;
                fromPartial(object: Partial<_94.QueryCurrentEpochResponse>): _94.QueryCurrentEpochResponse;
            };
            EpochInfo: {
                encode(message: _93.EpochInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _93.EpochInfo;
                fromJSON(object: any): _93.EpochInfo;
                toJSON(message: _93.EpochInfo): unknown;
                fromPartial(object: Partial<_93.EpochInfo>): _93.EpochInfo;
            };
            GenesisState: {
                encode(message: _93.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _93.GenesisState;
                fromJSON(object: any): _93.GenesisState;
                toJSON(message: _93.GenesisState): unknown;
                fromPartial(object: Partial<_93.GenesisState>): _93.GenesisState;
            };
        };
    }
    namespace interchainquery {
        const v1: {
            MsgClientImpl: typeof _190.MsgClientImpl;
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    submitQueryResponse(value: _97.MsgSubmitQueryResponse): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    submitQueryResponse(value: _97.MsgSubmitQueryResponse): {
                        typeUrl: string;
                        value: _97.MsgSubmitQueryResponse;
                    };
                };
                toJSON: {
                    submitQueryResponse(value: _97.MsgSubmitQueryResponse): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    submitQueryResponse(value: any): {
                        typeUrl: string;
                        value: _97.MsgSubmitQueryResponse;
                    };
                };
                fromPartial: {
                    submitQueryResponse(value: _97.MsgSubmitQueryResponse): {
                        typeUrl: string;
                        value: _97.MsgSubmitQueryResponse;
                    };
                };
            };
            AminoConverter: {
                "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse": {
                    aminoType: string;
                    toAmino: ({ chainId, queryId, result, proofOps, height, fromAddress }: _97.MsgSubmitQueryResponse) => {
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
                    }) => _97.MsgSubmitQueryResponse;
                };
            };
            QueryRequestsRequest: {
                encode(message: _98.QueryRequestsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _98.QueryRequestsRequest;
                fromJSON(object: any): _98.QueryRequestsRequest;
                toJSON(message: _98.QueryRequestsRequest): unknown;
                fromPartial(object: Partial<_98.QueryRequestsRequest>): _98.QueryRequestsRequest;
            };
            QueryRequestsResponse: {
                encode(message: _98.QueryRequestsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _98.QueryRequestsResponse;
                fromJSON(object: any): _98.QueryRequestsResponse;
                toJSON(message: _98.QueryRequestsResponse): unknown;
                fromPartial(object: Partial<_98.QueryRequestsResponse>): _98.QueryRequestsResponse;
            };
            GetTxWithProofResponse: {
                encode(message: _98.GetTxWithProofResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _98.GetTxWithProofResponse;
                fromJSON(object: any): _98.GetTxWithProofResponse;
                toJSON(message: _98.GetTxWithProofResponse): unknown;
                fromPartial(object: Partial<_98.GetTxWithProofResponse>): _98.GetTxWithProofResponse;
            };
            MsgSubmitQueryResponse: {
                encode(message: _97.MsgSubmitQueryResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _97.MsgSubmitQueryResponse;
                fromJSON(object: any): _97.MsgSubmitQueryResponse;
                toJSON(message: _97.MsgSubmitQueryResponse): unknown;
                fromPartial(object: Partial<_97.MsgSubmitQueryResponse>): _97.MsgSubmitQueryResponse;
            };
            MsgSubmitQueryResponseResponse: {
                encode(_: _97.MsgSubmitQueryResponseResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _97.MsgSubmitQueryResponseResponse;
                fromJSON(_: any): _97.MsgSubmitQueryResponseResponse;
                toJSON(_: _97.MsgSubmitQueryResponseResponse): unknown;
                fromPartial(_: Partial<_97.MsgSubmitQueryResponseResponse>): _97.MsgSubmitQueryResponseResponse;
            };
            Query: {
                encode(message: _96.Query, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _96.Query;
                fromJSON(object: any): _96.Query;
                toJSON(message: _96.Query): unknown;
                fromPartial(object: Partial<_96.Query>): _96.Query;
            };
            DataPoint: {
                encode(message: _96.DataPoint, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _96.DataPoint;
                fromJSON(object: any): _96.DataPoint;
                toJSON(message: _96.DataPoint): unknown;
                fromPartial(object: Partial<_96.DataPoint>): _96.DataPoint;
            };
            GenesisState: {
                encode(message: _95.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _95.GenesisState;
                fromJSON(object: any): _95.GenesisState;
                toJSON(message: _95.GenesisState): unknown;
                fromPartial(object: Partial<_95.GenesisState>): _95.GenesisState;
            };
        };
    }
    namespace interchainstaking {
        const v1: {
            MsgClientImpl: typeof _191.MsgClientImpl;
            QueryClientImpl: typeof _185.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                zoneInfos(request?: _103.QueryZonesInfoRequest): Promise<_103.QueryZonesInfoResponse>;
                depositAccount(request: _103.QueryDepositAccountForChainRequest): Promise<_103.QueryDepositAccountForChainResponse>;
                delegatorIntent(request: _103.QueryDelegatorIntentRequest): Promise<_103.QueryDelegatorIntentResponse>;
                delegations(request: _103.QueryDelegationsRequest): Promise<_103.QueryDelegationsResponse>;
                receipts(request: _103.QueryReceiptsRequest): Promise<_103.QueryReceiptsResponse>;
                zoneWithdrawalRecords(request: _103.QueryWithdrawalRecordsRequest): Promise<_103.QueryWithdrawalRecordsResponse>;
                withdrawalRecords(request: _103.QueryWithdrawalRecordsRequest): Promise<_103.QueryWithdrawalRecordsResponse>;
                unbondingRecords(request: _103.QueryUnbondingRecordsRequest): Promise<_103.QueryUnbondingRecordsResponse>;
                redelegationRecords(request: _103.QueryRedelegationRecordsRequest): Promise<_103.QueryRedelegationRecordsResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    requestRedemption(value: _101.MsgRequestRedemption): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    signalIntent(value: _101.MsgSignalIntent): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    requestRedemption(value: _101.MsgRequestRedemption): {
                        typeUrl: string;
                        value: _101.MsgRequestRedemption;
                    };
                    signalIntent(value: _101.MsgSignalIntent): {
                        typeUrl: string;
                        value: _101.MsgSignalIntent;
                    };
                };
                toJSON: {
                    requestRedemption(value: _101.MsgRequestRedemption): {
                        typeUrl: string;
                        value: unknown;
                    };
                    signalIntent(value: _101.MsgSignalIntent): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    requestRedemption(value: any): {
                        typeUrl: string;
                        value: _101.MsgRequestRedemption;
                    };
                    signalIntent(value: any): {
                        typeUrl: string;
                        value: _101.MsgSignalIntent;
                    };
                };
                fromPartial: {
                    requestRedemption(value: _101.MsgRequestRedemption): {
                        typeUrl: string;
                        value: _101.MsgRequestRedemption;
                    };
                    signalIntent(value: _101.MsgSignalIntent): {
                        typeUrl: string;
                        value: _101.MsgSignalIntent;
                    };
                };
            };
            AminoConverter: {
                "/quicksilver.interchainstaking.v1.MsgRequestRedemption": {
                    aminoType: string;
                    toAmino: ({ value, destinationAddress, fromAddress }: _101.MsgRequestRedemption) => {
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
                    }) => _101.MsgRequestRedemption;
                };
                "/quicksilver.interchainstaking.v1.MsgSignalIntent": {
                    aminoType: string;
                    toAmino: ({ chainId, intents, fromAddress }: _101.MsgSignalIntent) => {
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
                    }) => _101.MsgSignalIntent;
                };
            };
            AddProtocolDataProposal: {
                encode(message: _104.AddProtocolDataProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _104.AddProtocolDataProposal;
                fromJSON(object: any): _104.AddProtocolDataProposal;
                toJSON(message: _104.AddProtocolDataProposal): unknown;
                fromPartial(object: Partial<_104.AddProtocolDataProposal>): _104.AddProtocolDataProposal;
            };
            AddProtocolDataProposalWithDeposit: {
                encode(message: _104.AddProtocolDataProposalWithDeposit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _104.AddProtocolDataProposalWithDeposit;
                fromJSON(object: any): _104.AddProtocolDataProposalWithDeposit;
                toJSON(message: _104.AddProtocolDataProposalWithDeposit): unknown;
                fromPartial(object: Partial<_104.AddProtocolDataProposalWithDeposit>): _104.AddProtocolDataProposalWithDeposit;
            };
            QueryZonesInfoRequest: {
                encode(message: _103.QueryZonesInfoRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryZonesInfoRequest;
                fromJSON(object: any): _103.QueryZonesInfoRequest;
                toJSON(message: _103.QueryZonesInfoRequest): unknown;
                fromPartial(object: Partial<_103.QueryZonesInfoRequest>): _103.QueryZonesInfoRequest;
            };
            QueryZonesInfoResponse: {
                encode(message: _103.QueryZonesInfoResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryZonesInfoResponse;
                fromJSON(object: any): _103.QueryZonesInfoResponse;
                toJSON(message: _103.QueryZonesInfoResponse): unknown;
                fromPartial(object: Partial<_103.QueryZonesInfoResponse>): _103.QueryZonesInfoResponse;
            };
            QueryDepositAccountForChainRequest: {
                encode(message: _103.QueryDepositAccountForChainRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryDepositAccountForChainRequest;
                fromJSON(object: any): _103.QueryDepositAccountForChainRequest;
                toJSON(message: _103.QueryDepositAccountForChainRequest): unknown;
                fromPartial(object: Partial<_103.QueryDepositAccountForChainRequest>): _103.QueryDepositAccountForChainRequest;
            };
            QueryDepositAccountForChainResponse: {
                encode(message: _103.QueryDepositAccountForChainResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryDepositAccountForChainResponse;
                fromJSON(object: any): _103.QueryDepositAccountForChainResponse;
                toJSON(message: _103.QueryDepositAccountForChainResponse): unknown;
                fromPartial(object: Partial<_103.QueryDepositAccountForChainResponse>): _103.QueryDepositAccountForChainResponse;
            };
            QueryDelegatorIntentRequest: {
                encode(message: _103.QueryDelegatorIntentRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryDelegatorIntentRequest;
                fromJSON(object: any): _103.QueryDelegatorIntentRequest;
                toJSON(message: _103.QueryDelegatorIntentRequest): unknown;
                fromPartial(object: Partial<_103.QueryDelegatorIntentRequest>): _103.QueryDelegatorIntentRequest;
            };
            QueryDelegatorIntentResponse: {
                encode(message: _103.QueryDelegatorIntentResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryDelegatorIntentResponse;
                fromJSON(object: any): _103.QueryDelegatorIntentResponse;
                toJSON(message: _103.QueryDelegatorIntentResponse): unknown;
                fromPartial(object: Partial<_103.QueryDelegatorIntentResponse>): _103.QueryDelegatorIntentResponse;
            };
            QueryDelegationsRequest: {
                encode(message: _103.QueryDelegationsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryDelegationsRequest;
                fromJSON(object: any): _103.QueryDelegationsRequest;
                toJSON(message: _103.QueryDelegationsRequest): unknown;
                fromPartial(object: Partial<_103.QueryDelegationsRequest>): _103.QueryDelegationsRequest;
            };
            QueryDelegationsResponse: {
                encode(message: _103.QueryDelegationsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryDelegationsResponse;
                fromJSON(object: any): _103.QueryDelegationsResponse;
                toJSON(message: _103.QueryDelegationsResponse): unknown;
                fromPartial(object: Partial<_103.QueryDelegationsResponse>): _103.QueryDelegationsResponse;
            };
            QueryReceiptsRequest: {
                encode(message: _103.QueryReceiptsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryReceiptsRequest;
                fromJSON(object: any): _103.QueryReceiptsRequest;
                toJSON(message: _103.QueryReceiptsRequest): unknown;
                fromPartial(object: Partial<_103.QueryReceiptsRequest>): _103.QueryReceiptsRequest;
            };
            QueryReceiptsResponse: {
                encode(message: _103.QueryReceiptsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryReceiptsResponse;
                fromJSON(object: any): _103.QueryReceiptsResponse;
                toJSON(message: _103.QueryReceiptsResponse): unknown;
                fromPartial(object: Partial<_103.QueryReceiptsResponse>): _103.QueryReceiptsResponse;
            };
            QueryWithdrawalRecordsRequest: {
                encode(message: _103.QueryWithdrawalRecordsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryWithdrawalRecordsRequest;
                fromJSON(object: any): _103.QueryWithdrawalRecordsRequest;
                toJSON(message: _103.QueryWithdrawalRecordsRequest): unknown;
                fromPartial(object: Partial<_103.QueryWithdrawalRecordsRequest>): _103.QueryWithdrawalRecordsRequest;
            };
            QueryWithdrawalRecordsResponse: {
                encode(message: _103.QueryWithdrawalRecordsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryWithdrawalRecordsResponse;
                fromJSON(object: any): _103.QueryWithdrawalRecordsResponse;
                toJSON(message: _103.QueryWithdrawalRecordsResponse): unknown;
                fromPartial(object: Partial<_103.QueryWithdrawalRecordsResponse>): _103.QueryWithdrawalRecordsResponse;
            };
            QueryUnbondingRecordsRequest: {
                encode(message: _103.QueryUnbondingRecordsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryUnbondingRecordsRequest;
                fromJSON(object: any): _103.QueryUnbondingRecordsRequest;
                toJSON(message: _103.QueryUnbondingRecordsRequest): unknown;
                fromPartial(object: Partial<_103.QueryUnbondingRecordsRequest>): _103.QueryUnbondingRecordsRequest;
            };
            QueryUnbondingRecordsResponse: {
                encode(message: _103.QueryUnbondingRecordsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryUnbondingRecordsResponse;
                fromJSON(object: any): _103.QueryUnbondingRecordsResponse;
                toJSON(message: _103.QueryUnbondingRecordsResponse): unknown;
                fromPartial(object: Partial<_103.QueryUnbondingRecordsResponse>): _103.QueryUnbondingRecordsResponse;
            };
            QueryRedelegationRecordsRequest: {
                encode(message: _103.QueryRedelegationRecordsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryRedelegationRecordsRequest;
                fromJSON(object: any): _103.QueryRedelegationRecordsRequest;
                toJSON(message: _103.QueryRedelegationRecordsRequest): unknown;
                fromPartial(object: Partial<_103.QueryRedelegationRecordsRequest>): _103.QueryRedelegationRecordsRequest;
            };
            QueryRedelegationRecordsResponse: {
                encode(message: _103.QueryRedelegationRecordsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _103.QueryRedelegationRecordsResponse;
                fromJSON(object: any): _103.QueryRedelegationRecordsResponse;
                toJSON(message: _103.QueryRedelegationRecordsResponse): unknown;
                fromPartial(object: Partial<_103.QueryRedelegationRecordsResponse>): _103.QueryRedelegationRecordsResponse;
            };
            RegisterZoneProposal: {
                encode(message: _102.RegisterZoneProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _102.RegisterZoneProposal;
                fromJSON(object: any): _102.RegisterZoneProposal;
                toJSON(message: _102.RegisterZoneProposal): unknown;
                fromPartial(object: Partial<_102.RegisterZoneProposal>): _102.RegisterZoneProposal;
            };
            RegisterZoneProposalWithDeposit: {
                encode(message: _102.RegisterZoneProposalWithDeposit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _102.RegisterZoneProposalWithDeposit;
                fromJSON(object: any): _102.RegisterZoneProposalWithDeposit;
                toJSON(message: _102.RegisterZoneProposalWithDeposit): unknown;
                fromPartial(object: Partial<_102.RegisterZoneProposalWithDeposit>): _102.RegisterZoneProposalWithDeposit;
            };
            UpdateZoneProposal: {
                encode(message: _102.UpdateZoneProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _102.UpdateZoneProposal;
                fromJSON(object: any): _102.UpdateZoneProposal;
                toJSON(message: _102.UpdateZoneProposal): unknown;
                fromPartial(object: Partial<_102.UpdateZoneProposal>): _102.UpdateZoneProposal;
            };
            UpdateZoneProposalWithDeposit: {
                encode(message: _102.UpdateZoneProposalWithDeposit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _102.UpdateZoneProposalWithDeposit;
                fromJSON(object: any): _102.UpdateZoneProposalWithDeposit;
                toJSON(message: _102.UpdateZoneProposalWithDeposit): unknown;
                fromPartial(object: Partial<_102.UpdateZoneProposalWithDeposit>): _102.UpdateZoneProposalWithDeposit;
            };
            UpdateZoneValue: {
                encode(message: _102.UpdateZoneValue, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _102.UpdateZoneValue;
                fromJSON(object: any): _102.UpdateZoneValue;
                toJSON(message: _102.UpdateZoneValue): unknown;
                fromPartial(object: Partial<_102.UpdateZoneValue>): _102.UpdateZoneValue;
            };
            MsgRequestRedemption: {
                encode(message: _101.MsgRequestRedemption, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _101.MsgRequestRedemption;
                fromJSON(object: any): _101.MsgRequestRedemption;
                toJSON(message: _101.MsgRequestRedemption): unknown;
                fromPartial(object: Partial<_101.MsgRequestRedemption>): _101.MsgRequestRedemption;
            };
            MsgSignalIntent: {
                encode(message: _101.MsgSignalIntent, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _101.MsgSignalIntent;
                fromJSON(object: any): _101.MsgSignalIntent;
                toJSON(message: _101.MsgSignalIntent): unknown;
                fromPartial(object: Partial<_101.MsgSignalIntent>): _101.MsgSignalIntent;
            };
            MsgRequestRedemptionResponse: {
                encode(_: _101.MsgRequestRedemptionResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _101.MsgRequestRedemptionResponse;
                fromJSON(_: any): _101.MsgRequestRedemptionResponse;
                toJSON(_: _101.MsgRequestRedemptionResponse): unknown;
                fromPartial(_: Partial<_101.MsgRequestRedemptionResponse>): _101.MsgRequestRedemptionResponse;
            };
            MsgSignalIntentResponse: {
                encode(_: _101.MsgSignalIntentResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _101.MsgSignalIntentResponse;
                fromJSON(_: any): _101.MsgSignalIntentResponse;
                toJSON(_: _101.MsgSignalIntentResponse): unknown;
                fromPartial(_: Partial<_101.MsgSignalIntentResponse>): _101.MsgSignalIntentResponse;
            };
            Zone: {
                encode(message: _100.Zone, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.Zone;
                fromJSON(object: any): _100.Zone;
                toJSON(message: _100.Zone): unknown;
                fromPartial(object: Partial<_100.Zone>): _100.Zone;
            };
            ICAAccount: {
                encode(message: _100.ICAAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.ICAAccount;
                fromJSON(object: any): _100.ICAAccount;
                toJSON(message: _100.ICAAccount): unknown;
                fromPartial(object: Partial<_100.ICAAccount>): _100.ICAAccount;
            };
            Distribution: {
                encode(message: _100.Distribution, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.Distribution;
                fromJSON(object: any): _100.Distribution;
                toJSON(message: _100.Distribution): unknown;
                fromPartial(object: Partial<_100.Distribution>): _100.Distribution;
            };
            WithdrawalRecord: {
                encode(message: _100.WithdrawalRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.WithdrawalRecord;
                fromJSON(object: any): _100.WithdrawalRecord;
                toJSON(message: _100.WithdrawalRecord): unknown;
                fromPartial(object: Partial<_100.WithdrawalRecord>): _100.WithdrawalRecord;
            };
            UnbondingRecord: {
                encode(message: _100.UnbondingRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.UnbondingRecord;
                fromJSON(object: any): _100.UnbondingRecord;
                toJSON(message: _100.UnbondingRecord): unknown;
                fromPartial(object: Partial<_100.UnbondingRecord>): _100.UnbondingRecord;
            };
            RedelegationRecord: {
                encode(message: _100.RedelegationRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.RedelegationRecord;
                fromJSON(object: any): _100.RedelegationRecord;
                toJSON(message: _100.RedelegationRecord): unknown;
                fromPartial(object: Partial<_100.RedelegationRecord>): _100.RedelegationRecord;
            };
            TransferRecord: {
                encode(message: _100.TransferRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.TransferRecord;
                fromJSON(object: any): _100.TransferRecord;
                toJSON(message: _100.TransferRecord): unknown;
                fromPartial(object: Partial<_100.TransferRecord>): _100.TransferRecord;
            };
            Validator: {
                encode(message: _100.Validator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.Validator;
                fromJSON(object: any): _100.Validator;
                toJSON(message: _100.Validator): unknown;
                fromPartial(object: Partial<_100.Validator>): _100.Validator;
            };
            DelegatorIntent: {
                encode(message: _100.DelegatorIntent, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.DelegatorIntent;
                fromJSON(object: any): _100.DelegatorIntent;
                toJSON(message: _100.DelegatorIntent): unknown;
                fromPartial(object: Partial<_100.DelegatorIntent>): _100.DelegatorIntent;
            };
            ValidatorIntent: {
                encode(message: _100.ValidatorIntent, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.ValidatorIntent;
                fromJSON(object: any): _100.ValidatorIntent;
                toJSON(message: _100.ValidatorIntent): unknown;
                fromPartial(object: Partial<_100.ValidatorIntent>): _100.ValidatorIntent;
            };
            Delegation: {
                encode(message: _100.Delegation, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.Delegation;
                fromJSON(object: any): _100.Delegation;
                toJSON(message: _100.Delegation): unknown;
                fromPartial(object: Partial<_100.Delegation>): _100.Delegation;
            };
            PortConnectionTuple: {
                encode(message: _100.PortConnectionTuple, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.PortConnectionTuple;
                fromJSON(object: any): _100.PortConnectionTuple;
                toJSON(message: _100.PortConnectionTuple): unknown;
                fromPartial(object: Partial<_100.PortConnectionTuple>): _100.PortConnectionTuple;
            };
            Receipt: {
                encode(message: _100.Receipt, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _100.Receipt;
                fromJSON(object: any): _100.Receipt;
                toJSON(message: _100.Receipt): unknown;
                fromPartial(object: Partial<_100.Receipt>): _100.Receipt;
            };
            Params: {
                encode(message: _99.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _99.Params;
                fromJSON(object: any): _99.Params;
                toJSON(message: _99.Params): unknown;
                fromPartial(object: Partial<_99.Params>): _99.Params;
            };
            DelegationsForZone: {
                encode(message: _99.DelegationsForZone, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _99.DelegationsForZone;
                fromJSON(object: any): _99.DelegationsForZone;
                toJSON(message: _99.DelegationsForZone): unknown;
                fromPartial(object: Partial<_99.DelegationsForZone>): _99.DelegationsForZone;
            };
            DelegatorIntentsForZone: {
                encode(message: _99.DelegatorIntentsForZone, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _99.DelegatorIntentsForZone;
                fromJSON(object: any): _99.DelegatorIntentsForZone;
                toJSON(message: _99.DelegatorIntentsForZone): unknown;
                fromPartial(object: Partial<_99.DelegatorIntentsForZone>): _99.DelegatorIntentsForZone;
            };
            GenesisState: {
                encode(message: _99.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _99.GenesisState;
                fromJSON(object: any): _99.GenesisState;
                toJSON(message: _99.GenesisState): unknown;
                fromPartial(object: Partial<_99.GenesisState>): _99.GenesisState;
            };
        };
    }
    namespace mint {
        const v1beta1: {
            QueryClientImpl: typeof _186.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _107.QueryParamsRequest): Promise<_107.QueryParamsResponse>;
                epochProvisions(request?: _107.QueryEpochProvisionsRequest): Promise<_107.QueryEpochProvisionsResponse>;
            };
            QueryParamsRequest: {
                encode(_: _107.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _107.QueryParamsRequest;
                fromJSON(_: any): _107.QueryParamsRequest;
                toJSON(_: _107.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_107.QueryParamsRequest>): _107.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _107.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _107.QueryParamsResponse;
                fromJSON(object: any): _107.QueryParamsResponse;
                toJSON(message: _107.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_107.QueryParamsResponse>): _107.QueryParamsResponse;
            };
            QueryEpochProvisionsRequest: {
                encode(_: _107.QueryEpochProvisionsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _107.QueryEpochProvisionsRequest;
                fromJSON(_: any): _107.QueryEpochProvisionsRequest;
                toJSON(_: _107.QueryEpochProvisionsRequest): unknown;
                fromPartial(_: Partial<_107.QueryEpochProvisionsRequest>): _107.QueryEpochProvisionsRequest;
            };
            QueryEpochProvisionsResponse: {
                encode(message: _107.QueryEpochProvisionsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _107.QueryEpochProvisionsResponse;
                fromJSON(object: any): _107.QueryEpochProvisionsResponse;
                toJSON(message: _107.QueryEpochProvisionsResponse): unknown;
                fromPartial(object: Partial<_107.QueryEpochProvisionsResponse>): _107.QueryEpochProvisionsResponse;
            };
            Minter: {
                encode(message: _106.Minter, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _106.Minter;
                fromJSON(object: any): _106.Minter;
                toJSON(message: _106.Minter): unknown;
                fromPartial(object: Partial<_106.Minter>): _106.Minter;
            };
            DistributionProportions: {
                encode(message: _106.DistributionProportions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _106.DistributionProportions;
                fromJSON(object: any): _106.DistributionProportions;
                toJSON(message: _106.DistributionProportions): unknown;
                fromPartial(object: Partial<_106.DistributionProportions>): _106.DistributionProportions;
            };
            Params: {
                encode(message: _106.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _106.Params;
                fromJSON(object: any): _106.Params;
                toJSON(message: _106.Params): unknown;
                fromPartial(object: Partial<_106.Params>): _106.Params;
            };
            GenesisState: {
                encode(message: _105.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _105.GenesisState;
                fromJSON(object: any): _105.GenesisState;
                toJSON(message: _105.GenesisState): unknown;
                fromPartial(object: Partial<_105.GenesisState>): _105.GenesisState;
            };
        };
    }
    namespace participationrewards {
        const v1: {
            MsgClientImpl: typeof _192.MsgClientImpl;
            QueryClientImpl: typeof _187.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _111.QueryParamsRequest): Promise<_111.QueryParamsResponse>;
                protocolData(request: _111.QueryProtocolDataRequest): Promise<_111.QueryProtocolDataResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    submitClaim(value: _109.MsgSubmitClaim): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    submitClaim(value: _109.MsgSubmitClaim): {
                        typeUrl: string;
                        value: _109.MsgSubmitClaim;
                    };
                };
                toJSON: {
                    submitClaim(value: _109.MsgSubmitClaim): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    submitClaim(value: any): {
                        typeUrl: string;
                        value: _109.MsgSubmitClaim;
                    };
                };
                fromPartial: {
                    submitClaim(value: _109.MsgSubmitClaim): {
                        typeUrl: string;
                        value: _109.MsgSubmitClaim;
                    };
                };
            };
            AminoConverter: {
                "/quicksilver.participationrewards.v1.MsgSubmitClaim": {
                    aminoType: string;
                    toAmino: ({ userAddress, zone, srcZone, claimType, proofs }: _109.MsgSubmitClaim) => {
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
                    }) => _109.MsgSubmitClaim;
                };
            };
            QueryParamsRequest: {
                encode(_: _111.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _111.QueryParamsRequest;
                fromJSON(_: any): _111.QueryParamsRequest;
                toJSON(_: _111.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_111.QueryParamsRequest>): _111.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _111.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _111.QueryParamsResponse;
                fromJSON(object: any): _111.QueryParamsResponse;
                toJSON(message: _111.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_111.QueryParamsResponse>): _111.QueryParamsResponse;
            };
            QueryProtocolDataRequest: {
                encode(message: _111.QueryProtocolDataRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _111.QueryProtocolDataRequest;
                fromJSON(object: any): _111.QueryProtocolDataRequest;
                toJSON(message: _111.QueryProtocolDataRequest): unknown;
                fromPartial(object: Partial<_111.QueryProtocolDataRequest>): _111.QueryProtocolDataRequest;
            };
            QueryProtocolDataResponse: {
                encode(message: _111.QueryProtocolDataResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _111.QueryProtocolDataResponse;
                fromJSON(object: any): _111.QueryProtocolDataResponse;
                toJSON(message: _111.QueryProtocolDataResponse): unknown;
                fromPartial(object: Partial<_111.QueryProtocolDataResponse>): _111.QueryProtocolDataResponse;
            };
            protocolDataTypeFromJSON(object: any): _110.ProtocolDataType;
            protocolDataTypeToJSON(object: _110.ProtocolDataType): string;
            ProtocolDataType: typeof _110.ProtocolDataType;
            ProtocolDataTypeSDKType: typeof _110.ProtocolDataTypeSDKType;
            DistributionProportions: {
                encode(message: _110.DistributionProportions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _110.DistributionProportions;
                fromJSON(object: any): _110.DistributionProportions;
                toJSON(message: _110.DistributionProportions): unknown;
                fromPartial(object: Partial<_110.DistributionProportions>): _110.DistributionProportions;
            };
            Params: {
                encode(message: _110.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _110.Params;
                fromJSON(object: any): _110.Params;
                toJSON(message: _110.Params): unknown;
                fromPartial(object: Partial<_110.Params>): _110.Params;
            };
            KeyedProtocolData: {
                encode(message: _110.KeyedProtocolData, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _110.KeyedProtocolData;
                fromJSON(object: any): _110.KeyedProtocolData;
                toJSON(message: _110.KeyedProtocolData): unknown;
                fromPartial(object: Partial<_110.KeyedProtocolData>): _110.KeyedProtocolData;
            };
            ProtocolData: {
                encode(message: _110.ProtocolData, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _110.ProtocolData;
                fromJSON(object: any): _110.ProtocolData;
                toJSON(message: _110.ProtocolData): unknown;
                fromPartial(object: Partial<_110.ProtocolData>): _110.ProtocolData;
            };
            MsgSubmitClaim: {
                encode(message: _109.MsgSubmitClaim, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _109.MsgSubmitClaim;
                fromJSON(object: any): _109.MsgSubmitClaim;
                toJSON(message: _109.MsgSubmitClaim): unknown;
                fromPartial(object: Partial<_109.MsgSubmitClaim>): _109.MsgSubmitClaim;
            };
            MsgSubmitClaimResponse: {
                encode(_: _109.MsgSubmitClaimResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _109.MsgSubmitClaimResponse;
                fromJSON(_: any): _109.MsgSubmitClaimResponse;
                toJSON(_: _109.MsgSubmitClaimResponse): unknown;
                fromPartial(_: Partial<_109.MsgSubmitClaimResponse>): _109.MsgSubmitClaimResponse;
            };
            Proof: {
                encode(message: _109.Proof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _109.Proof;
                fromJSON(object: any): _109.Proof;
                toJSON(message: _109.Proof): unknown;
                fromPartial(object: Partial<_109.Proof>): _109.Proof;
            };
            GenesisState: {
                encode(message: _108.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _108.GenesisState;
                fromJSON(object: any): _108.GenesisState;
                toJSON(message: _108.GenesisState): unknown;
                fromPartial(object: Partial<_108.GenesisState>): _108.GenesisState;
            };
        };
    }
    namespace tokenfactory {
        const v1beta1: {
            MsgClientImpl: typeof _193.MsgClientImpl;
            QueryClientImpl: typeof _188.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _115.QueryParamsRequest): Promise<_115.QueryParamsResponse>;
                denomAuthorityMetadata(request: _115.QueryDenomAuthorityMetadataRequest): Promise<_115.QueryDenomAuthorityMetadataResponse>;
                denomsFromCreator(request: _115.QueryDenomsFromCreatorRequest): Promise<_115.QueryDenomsFromCreatorResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    createDenom(value: _116.MsgCreateDenom): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    mint(value: _116.MsgMint): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    burn(value: _116.MsgBurn): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    changeAdmin(value: _116.MsgChangeAdmin): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    setDenomMetadata(value: _116.MsgSetDenomMetadata): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    createDenom(value: _116.MsgCreateDenom): {
                        typeUrl: string;
                        value: _116.MsgCreateDenom;
                    };
                    mint(value: _116.MsgMint): {
                        typeUrl: string;
                        value: _116.MsgMint;
                    };
                    burn(value: _116.MsgBurn): {
                        typeUrl: string;
                        value: _116.MsgBurn;
                    };
                    changeAdmin(value: _116.MsgChangeAdmin): {
                        typeUrl: string;
                        value: _116.MsgChangeAdmin;
                    };
                    setDenomMetadata(value: _116.MsgSetDenomMetadata): {
                        typeUrl: string;
                        value: _116.MsgSetDenomMetadata;
                    };
                };
                toJSON: {
                    createDenom(value: _116.MsgCreateDenom): {
                        typeUrl: string;
                        value: unknown;
                    };
                    mint(value: _116.MsgMint): {
                        typeUrl: string;
                        value: unknown;
                    };
                    burn(value: _116.MsgBurn): {
                        typeUrl: string;
                        value: unknown;
                    };
                    changeAdmin(value: _116.MsgChangeAdmin): {
                        typeUrl: string;
                        value: unknown;
                    };
                    setDenomMetadata(value: _116.MsgSetDenomMetadata): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    createDenom(value: any): {
                        typeUrl: string;
                        value: _116.MsgCreateDenom;
                    };
                    mint(value: any): {
                        typeUrl: string;
                        value: _116.MsgMint;
                    };
                    burn(value: any): {
                        typeUrl: string;
                        value: _116.MsgBurn;
                    };
                    changeAdmin(value: any): {
                        typeUrl: string;
                        value: _116.MsgChangeAdmin;
                    };
                    setDenomMetadata(value: any): {
                        typeUrl: string;
                        value: _116.MsgSetDenomMetadata;
                    };
                };
                fromPartial: {
                    createDenom(value: _116.MsgCreateDenom): {
                        typeUrl: string;
                        value: _116.MsgCreateDenom;
                    };
                    mint(value: _116.MsgMint): {
                        typeUrl: string;
                        value: _116.MsgMint;
                    };
                    burn(value: _116.MsgBurn): {
                        typeUrl: string;
                        value: _116.MsgBurn;
                    };
                    changeAdmin(value: _116.MsgChangeAdmin): {
                        typeUrl: string;
                        value: _116.MsgChangeAdmin;
                    };
                    setDenomMetadata(value: _116.MsgSetDenomMetadata): {
                        typeUrl: string;
                        value: _116.MsgSetDenomMetadata;
                    };
                };
            };
            AminoConverter: {
                "/quicksilver.tokenfactory.v1beta1.MsgCreateDenom": {
                    aminoType: string;
                    toAmino: ({ sender, subdenom }: _116.MsgCreateDenom) => {
                        sender: string;
                        subdenom: string;
                    };
                    fromAmino: ({ sender, subdenom }: {
                        sender: string;
                        subdenom: string;
                    }) => _116.MsgCreateDenom;
                };
                "/quicksilver.tokenfactory.v1beta1.MsgMint": {
                    aminoType: string;
                    toAmino: ({ sender, amount }: _116.MsgMint) => {
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
                    }) => _116.MsgMint;
                };
                "/quicksilver.tokenfactory.v1beta1.MsgBurn": {
                    aminoType: string;
                    toAmino: ({ sender, amount }: _116.MsgBurn) => {
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
                    }) => _116.MsgBurn;
                };
                "/quicksilver.tokenfactory.v1beta1.MsgChangeAdmin": {
                    aminoType: string;
                    toAmino: ({ sender, denom, newAdmin }: _116.MsgChangeAdmin) => {
                        sender: string;
                        denom: string;
                        new_admin: string;
                    };
                    fromAmino: ({ sender, denom, new_admin }: {
                        sender: string;
                        denom: string;
                        new_admin: string;
                    }) => _116.MsgChangeAdmin;
                };
                "/quicksilver.tokenfactory.v1beta1.MsgSetDenomMetadata": {
                    aminoType: string;
                    toAmino: ({ sender, metadata }: _116.MsgSetDenomMetadata) => {
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
                    }) => _116.MsgSetDenomMetadata;
                };
            };
            MsgCreateDenom: {
                encode(message: _116.MsgCreateDenom, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgCreateDenom;
                fromJSON(object: any): _116.MsgCreateDenom;
                toJSON(message: _116.MsgCreateDenom): unknown;
                fromPartial(object: Partial<_116.MsgCreateDenom>): _116.MsgCreateDenom;
            };
            MsgCreateDenomResponse: {
                encode(message: _116.MsgCreateDenomResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgCreateDenomResponse;
                fromJSON(object: any): _116.MsgCreateDenomResponse;
                toJSON(message: _116.MsgCreateDenomResponse): unknown;
                fromPartial(object: Partial<_116.MsgCreateDenomResponse>): _116.MsgCreateDenomResponse;
            };
            MsgMint: {
                encode(message: _116.MsgMint, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgMint;
                fromJSON(object: any): _116.MsgMint;
                toJSON(message: _116.MsgMint): unknown;
                fromPartial(object: Partial<_116.MsgMint>): _116.MsgMint;
            };
            MsgMintResponse: {
                encode(_: _116.MsgMintResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgMintResponse;
                fromJSON(_: any): _116.MsgMintResponse;
                toJSON(_: _116.MsgMintResponse): unknown;
                fromPartial(_: Partial<_116.MsgMintResponse>): _116.MsgMintResponse;
            };
            MsgBurn: {
                encode(message: _116.MsgBurn, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgBurn;
                fromJSON(object: any): _116.MsgBurn;
                toJSON(message: _116.MsgBurn): unknown;
                fromPartial(object: Partial<_116.MsgBurn>): _116.MsgBurn;
            };
            MsgBurnResponse: {
                encode(_: _116.MsgBurnResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgBurnResponse;
                fromJSON(_: any): _116.MsgBurnResponse;
                toJSON(_: _116.MsgBurnResponse): unknown;
                fromPartial(_: Partial<_116.MsgBurnResponse>): _116.MsgBurnResponse;
            };
            MsgChangeAdmin: {
                encode(message: _116.MsgChangeAdmin, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgChangeAdmin;
                fromJSON(object: any): _116.MsgChangeAdmin;
                toJSON(message: _116.MsgChangeAdmin): unknown;
                fromPartial(object: Partial<_116.MsgChangeAdmin>): _116.MsgChangeAdmin;
            };
            MsgChangeAdminResponse: {
                encode(_: _116.MsgChangeAdminResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgChangeAdminResponse;
                fromJSON(_: any): _116.MsgChangeAdminResponse;
                toJSON(_: _116.MsgChangeAdminResponse): unknown;
                fromPartial(_: Partial<_116.MsgChangeAdminResponse>): _116.MsgChangeAdminResponse;
            };
            MsgSetDenomMetadata: {
                encode(message: _116.MsgSetDenomMetadata, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgSetDenomMetadata;
                fromJSON(object: any): _116.MsgSetDenomMetadata;
                toJSON(message: _116.MsgSetDenomMetadata): unknown;
                fromPartial(object: Partial<_116.MsgSetDenomMetadata>): _116.MsgSetDenomMetadata;
            };
            MsgSetDenomMetadataResponse: {
                encode(_: _116.MsgSetDenomMetadataResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _116.MsgSetDenomMetadataResponse;
                fromJSON(_: any): _116.MsgSetDenomMetadataResponse;
                toJSON(_: _116.MsgSetDenomMetadataResponse): unknown;
                fromPartial(_: Partial<_116.MsgSetDenomMetadataResponse>): _116.MsgSetDenomMetadataResponse;
            };
            QueryParamsRequest: {
                encode(_: _115.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _115.QueryParamsRequest;
                fromJSON(_: any): _115.QueryParamsRequest;
                toJSON(_: _115.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_115.QueryParamsRequest>): _115.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _115.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _115.QueryParamsResponse;
                fromJSON(object: any): _115.QueryParamsResponse;
                toJSON(message: _115.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_115.QueryParamsResponse>): _115.QueryParamsResponse;
            };
            QueryDenomAuthorityMetadataRequest: {
                encode(message: _115.QueryDenomAuthorityMetadataRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _115.QueryDenomAuthorityMetadataRequest;
                fromJSON(object: any): _115.QueryDenomAuthorityMetadataRequest;
                toJSON(message: _115.QueryDenomAuthorityMetadataRequest): unknown;
                fromPartial(object: Partial<_115.QueryDenomAuthorityMetadataRequest>): _115.QueryDenomAuthorityMetadataRequest;
            };
            QueryDenomAuthorityMetadataResponse: {
                encode(message: _115.QueryDenomAuthorityMetadataResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _115.QueryDenomAuthorityMetadataResponse;
                fromJSON(object: any): _115.QueryDenomAuthorityMetadataResponse;
                toJSON(message: _115.QueryDenomAuthorityMetadataResponse): unknown;
                fromPartial(object: Partial<_115.QueryDenomAuthorityMetadataResponse>): _115.QueryDenomAuthorityMetadataResponse;
            };
            QueryDenomsFromCreatorRequest: {
                encode(message: _115.QueryDenomsFromCreatorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _115.QueryDenomsFromCreatorRequest;
                fromJSON(object: any): _115.QueryDenomsFromCreatorRequest;
                toJSON(message: _115.QueryDenomsFromCreatorRequest): unknown;
                fromPartial(object: Partial<_115.QueryDenomsFromCreatorRequest>): _115.QueryDenomsFromCreatorRequest;
            };
            QueryDenomsFromCreatorResponse: {
                encode(message: _115.QueryDenomsFromCreatorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _115.QueryDenomsFromCreatorResponse;
                fromJSON(object: any): _115.QueryDenomsFromCreatorResponse;
                toJSON(message: _115.QueryDenomsFromCreatorResponse): unknown;
                fromPartial(object: Partial<_115.QueryDenomsFromCreatorResponse>): _115.QueryDenomsFromCreatorResponse;
            };
            Params: {
                encode(message: _114.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _114.Params;
                fromJSON(object: any): _114.Params;
                toJSON(message: _114.Params): unknown;
                fromPartial(object: Partial<_114.Params>): _114.Params;
            };
            GenesisState: {
                encode(message: _113.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _113.GenesisState;
                fromJSON(object: any): _113.GenesisState;
                toJSON(message: _113.GenesisState): unknown;
                fromPartial(object: Partial<_113.GenesisState>): _113.GenesisState;
            };
            GenesisDenom: {
                encode(message: _113.GenesisDenom, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _113.GenesisDenom;
                fromJSON(object: any): _113.GenesisDenom;
                toJSON(message: _113.GenesisDenom): unknown;
                fromPartial(object: Partial<_113.GenesisDenom>): _113.GenesisDenom;
            };
            DenomAuthorityMetadata: {
                encode(message: _112.DenomAuthorityMetadata, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _112.DenomAuthorityMetadata;
                fromJSON(object: any): _112.DenomAuthorityMetadata;
                toJSON(message: _112.DenomAuthorityMetadata): unknown;
                fromPartial(object: Partial<_112.DenomAuthorityMetadata>): _112.DenomAuthorityMetadata;
            };
        };
    }
    const ClientFactory: {
        createRPCMsgClient: ({ rpc }: {
            rpc: import("../helpers").Rpc;
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
                    v1: _189.MsgClientImpl;
                };
                interchainquery: {
                    v1: _190.MsgClientImpl;
                };
                interchainstaking: {
                    v1: _191.MsgClientImpl;
                };
                participationrewards: {
                    v1: _192.MsgClientImpl;
                };
                tokenfactory: {
                    v1beta1: _193.MsgClientImpl;
                };
            };
        }>;
        createRPCQueryClient: ({ rpcEndpoint }: {
            rpcEndpoint: string;
        }) => Promise<{
            cosmos: {
                auth: {
                    v1beta1: {
                        accounts(request?: import("../cosmos/auth/v1beta1/query").QueryAccountsRequest): Promise<import("../cosmos/auth/v1beta1/query").QueryAccountsResponse>;
                        account(request: import("../cosmos/auth/v1beta1/query").QueryAccountRequest): Promise<import("../cosmos/auth/v1beta1/query").QueryAccountResponse>;
                        params(request?: import("../cosmos/auth/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/auth/v1beta1/query").QueryParamsResponse>;
                    };
                };
                authz: {
                    v1beta1: {
                        grants(request: import("../cosmos/authz/v1beta1/query").QueryGrantsRequest): Promise<import("../cosmos/authz/v1beta1/query").QueryGrantsResponse>;
                    };
                };
                bank: {
                    v1beta1: {
                        balance(request: import("../cosmos/bank/v1beta1/query").QueryBalanceRequest): Promise<import("../cosmos/bank/v1beta1/query").QueryBalanceResponse>;
                        allBalances(request: import("../cosmos/bank/v1beta1/query").QueryAllBalancesRequest): Promise<import("../cosmos/bank/v1beta1/query").QueryAllBalancesResponse>;
                        totalSupply(request?: import("../cosmos/bank/v1beta1/query").QueryTotalSupplyRequest): Promise<import("../cosmos/bank/v1beta1/query").QueryTotalSupplyResponse>;
                        supplyOf(request: import("../cosmos/bank/v1beta1/query").QuerySupplyOfRequest): Promise<import("../cosmos/bank/v1beta1/query").QuerySupplyOfResponse>;
                        params(request?: import("../cosmos/bank/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/bank/v1beta1/query").QueryParamsResponse>;
                        denomMetadata(request: import("../cosmos/bank/v1beta1/query").QueryDenomMetadataRequest): Promise<import("../cosmos/bank/v1beta1/query").QueryDenomMetadataResponse>;
                        denomsMetadata(request?: import("../cosmos/bank/v1beta1/query").QueryDenomsMetadataRequest): Promise<import("../cosmos/bank/v1beta1/query").QueryDenomsMetadataResponse>;
                    };
                };
                base: {
                    tendermint: {
                        v1beta1: {
                            getNodeInfo(request?: import("../cosmos/base/tendermint/v1beta1/query").GetNodeInfoRequest): Promise<import("../cosmos/base/tendermint/v1beta1/query").GetNodeInfoResponse>;
                            getSyncing(request?: import("../cosmos/base/tendermint/v1beta1/query").GetSyncingRequest): Promise<import("../cosmos/base/tendermint/v1beta1/query").GetSyncingResponse>;
                            getLatestBlock(request?: import("../cosmos/base/tendermint/v1beta1/query").GetLatestBlockRequest): Promise<import("../cosmos/base/tendermint/v1beta1/query").GetLatestBlockResponse>;
                            getBlockByHeight(request: import("../cosmos/base/tendermint/v1beta1/query").GetBlockByHeightRequest): Promise<import("../cosmos/base/tendermint/v1beta1/query").GetBlockByHeightResponse>;
                            getLatestValidatorSet(request?: import("../cosmos/base/tendermint/v1beta1/query").GetLatestValidatorSetRequest): Promise<import("../cosmos/base/tendermint/v1beta1/query").GetLatestValidatorSetResponse>;
                            getValidatorSetByHeight(request: import("../cosmos/base/tendermint/v1beta1/query").GetValidatorSetByHeightRequest): Promise<import("../cosmos/base/tendermint/v1beta1/query").GetValidatorSetByHeightResponse>;
                        };
                    };
                };
                distribution: {
                    v1beta1: {
                        params(request?: import("../cosmos/distribution/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryParamsResponse>;
                        validatorOutstandingRewards(request: import("../cosmos/distribution/v1beta1/query").QueryValidatorOutstandingRewardsRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryValidatorOutstandingRewardsResponse>;
                        validatorCommission(request: import("../cosmos/distribution/v1beta1/query").QueryValidatorCommissionRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryValidatorCommissionResponse>;
                        validatorSlashes(request: import("../cosmos/distribution/v1beta1/query").QueryValidatorSlashesRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryValidatorSlashesResponse>;
                        delegationRewards(request: import("../cosmos/distribution/v1beta1/query").QueryDelegationRewardsRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryDelegationRewardsResponse>;
                        delegationTotalRewards(request: import("../cosmos/distribution/v1beta1/query").QueryDelegationTotalRewardsRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryDelegationTotalRewardsResponse>;
                        delegatorValidators(request: import("../cosmos/distribution/v1beta1/query").QueryDelegatorValidatorsRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryDelegatorValidatorsResponse>;
                        delegatorWithdrawAddress(request: import("../cosmos/distribution/v1beta1/query").QueryDelegatorWithdrawAddressRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryDelegatorWithdrawAddressResponse>;
                        communityPool(request?: import("../cosmos/distribution/v1beta1/query").QueryCommunityPoolRequest): Promise<import("../cosmos/distribution/v1beta1/query").QueryCommunityPoolResponse>;
                    };
                };
                evidence: {
                    v1beta1: {
                        evidence(request: import("../cosmos/evidence/v1beta1/query").QueryEvidenceRequest): Promise<import("../cosmos/evidence/v1beta1/query").QueryEvidenceResponse>;
                        allEvidence(request?: import("../cosmos/evidence/v1beta1/query").QueryAllEvidenceRequest): Promise<import("../cosmos/evidence/v1beta1/query").QueryAllEvidenceResponse>;
                    };
                };
                feegrant: {
                    v1beta1: {
                        allowance(request: import("../cosmos/feegrant/v1beta1/query").QueryAllowanceRequest): Promise<import("../cosmos/feegrant/v1beta1/query").QueryAllowanceResponse>;
                        allowances(request: import("../cosmos/feegrant/v1beta1/query").QueryAllowancesRequest): Promise<import("../cosmos/feegrant/v1beta1/query").QueryAllowancesResponse>;
                    };
                };
                gov: {
                    v1beta1: {
                        proposal(request: import("../cosmos/gov/v1beta1/query").QueryProposalRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryProposalResponse>;
                        proposals(request: import("../cosmos/gov/v1beta1/query").QueryProposalsRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryProposalsResponse>;
                        vote(request: import("../cosmos/gov/v1beta1/query").QueryVoteRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryVoteResponse>;
                        votes(request: import("../cosmos/gov/v1beta1/query").QueryVotesRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryVotesResponse>;
                        params(request: import("../cosmos/gov/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryParamsResponse>;
                        deposit(request: import("../cosmos/gov/v1beta1/query").QueryDepositRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryDepositResponse>;
                        deposits(request: import("../cosmos/gov/v1beta1/query").QueryDepositsRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryDepositsResponse>;
                        tallyResult(request: import("../cosmos/gov/v1beta1/query").QueryTallyResultRequest): Promise<import("../cosmos/gov/v1beta1/query").QueryTallyResultResponse>;
                    };
                };
                mint: {
                    v1beta1: {
                        params(request?: import("../cosmos/mint/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/mint/v1beta1/query").QueryParamsResponse>;
                        inflation(request?: import("../cosmos/mint/v1beta1/query").QueryInflationRequest): Promise<import("../cosmos/mint/v1beta1/query").QueryInflationResponse>;
                        annualProvisions(request?: import("../cosmos/mint/v1beta1/query").QueryAnnualProvisionsRequest): Promise<import("../cosmos/mint/v1beta1/query").QueryAnnualProvisionsResponse>;
                    };
                };
                params: {
                    v1beta1: {
                        params(request: import("../cosmos/params/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/params/v1beta1/query").QueryParamsResponse>;
                    };
                };
                slashing: {
                    v1beta1: {
                        params(request?: import("../cosmos/slashing/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/slashing/v1beta1/query").QueryParamsResponse>;
                        signingInfo(request: import("../cosmos/slashing/v1beta1/query").QuerySigningInfoRequest): Promise<import("../cosmos/slashing/v1beta1/query").QuerySigningInfoResponse>;
                        signingInfos(request?: import("../cosmos/slashing/v1beta1/query").QuerySigningInfosRequest): Promise<import("../cosmos/slashing/v1beta1/query").QuerySigningInfosResponse>;
                    };
                };
                staking: {
                    v1beta1: {
                        validators(request: import("../cosmos/staking/v1beta1/query").QueryValidatorsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryValidatorsResponse>;
                        validator(request: import("../cosmos/staking/v1beta1/query").QueryValidatorRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryValidatorResponse>;
                        validatorDelegations(request: import("../cosmos/staking/v1beta1/query").QueryValidatorDelegationsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryValidatorDelegationsResponse>;
                        validatorUnbondingDelegations(request: import("../cosmos/staking/v1beta1/query").QueryValidatorUnbondingDelegationsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryValidatorUnbondingDelegationsResponse>;
                        delegation(request: import("../cosmos/staking/v1beta1/query").QueryDelegationRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryDelegationResponse>;
                        unbondingDelegation(request: import("../cosmos/staking/v1beta1/query").QueryUnbondingDelegationRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryUnbondingDelegationResponse>;
                        delegatorDelegations(request: import("../cosmos/staking/v1beta1/query").QueryDelegatorDelegationsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryDelegatorDelegationsResponse>;
                        delegatorUnbondingDelegations(request: import("../cosmos/staking/v1beta1/query").QueryDelegatorUnbondingDelegationsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryDelegatorUnbondingDelegationsResponse>;
                        redelegations(request: import("../cosmos/staking/v1beta1/query").QueryRedelegationsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryRedelegationsResponse>;
                        delegatorValidators(request: import("../cosmos/staking/v1beta1/query").QueryDelegatorValidatorsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryDelegatorValidatorsResponse>;
                        delegatorValidator(request: import("../cosmos/staking/v1beta1/query").QueryDelegatorValidatorRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryDelegatorValidatorResponse>;
                        historicalInfo(request: import("../cosmos/staking/v1beta1/query").QueryHistoricalInfoRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryHistoricalInfoResponse>;
                        pool(request?: import("../cosmos/staking/v1beta1/query").QueryPoolRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryPoolResponse>;
                        params(request?: import("../cosmos/staking/v1beta1/query").QueryParamsRequest): Promise<import("../cosmos/staking/v1beta1/query").QueryParamsResponse>;
                    };
                };
                tx: {
                    v1beta1: {
                        simulate(request: import("../cosmos/tx/v1beta1/service").SimulateRequest): Promise<import("../cosmos/tx/v1beta1/service").SimulateResponse>;
                        getTx(request: import("../cosmos/tx/v1beta1/service").GetTxRequest): Promise<import("../cosmos/tx/v1beta1/service").GetTxResponse>;
                        broadcastTx(request: import("../cosmos/tx/v1beta1/service").BroadcastTxRequest): Promise<import("../cosmos/tx/v1beta1/service").BroadcastTxResponse>;
                        getTxsEvent(request: import("../cosmos/tx/v1beta1/service").GetTxsEventRequest): Promise<import("../cosmos/tx/v1beta1/service").GetTxsEventResponse>;
                    };
                };
                upgrade: {
                    v1beta1: {
                        currentPlan(request?: import("../cosmos/upgrade/v1beta1/query").QueryCurrentPlanRequest): Promise<import("../cosmos/upgrade/v1beta1/query").QueryCurrentPlanResponse>;
                        appliedPlan(request: import("../cosmos/upgrade/v1beta1/query").QueryAppliedPlanRequest): Promise<import("../cosmos/upgrade/v1beta1/query").QueryAppliedPlanResponse>;
                        upgradedConsensusState(request: import("../cosmos/upgrade/v1beta1/query").QueryUpgradedConsensusStateRequest): Promise<import("../cosmos/upgrade/v1beta1/query").QueryUpgradedConsensusStateResponse>;
                        moduleVersions(request: import("../cosmos/upgrade/v1beta1/query").QueryModuleVersionsRequest): Promise<import("../cosmos/upgrade/v1beta1/query").QueryModuleVersionsResponse>;
                    };
                };
            };
            quicksilver: {
                airdrop: {
                    v1: {
                        params(request?: _88.QueryParamsRequest): Promise<_88.QueryParamsResponse>;
                        zoneDrop(request: _88.QueryZoneDropRequest): Promise<_88.QueryZoneDropResponse>;
                        accountBalance(request: _88.QueryAccountBalanceRequest): Promise<_88.QueryAccountBalanceResponse>;
                        zoneDrops(request: _88.QueryZoneDropsRequest): Promise<_88.QueryZoneDropsResponse>;
                        claimRecord(request: _88.QueryClaimRecordRequest): Promise<_88.QueryClaimRecordResponse>;
                        claimRecords(request: _88.QueryClaimRecordsRequest): Promise<_88.QueryClaimRecordsResponse>;
                    };
                };
                claimsmanager: {
                    v1: {
                        params(request?: _92.QueryParamsRequest): Promise<_92.QueryParamsResponse>;
                        claims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
                        lastEpochClaims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
                        userClaims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
                        userLastEpochClaims(request: _92.QueryClaimsRequest): Promise<_92.QueryClaimsResponse>;
                    };
                };
                epochs: {
                    v1: {
                        epochInfos(request?: _94.QueryEpochsInfoRequest): Promise<_94.QueryEpochsInfoResponse>;
                        currentEpoch(request: _94.QueryCurrentEpochRequest): Promise<_94.QueryCurrentEpochResponse>;
                    };
                };
                interchainstaking: {
                    v1: {
                        zoneInfos(request?: _103.QueryZonesInfoRequest): Promise<_103.QueryZonesInfoResponse>;
                        depositAccount(request: _103.QueryDepositAccountForChainRequest): Promise<_103.QueryDepositAccountForChainResponse>;
                        delegatorIntent(request: _103.QueryDelegatorIntentRequest): Promise<_103.QueryDelegatorIntentResponse>;
                        delegations(request: _103.QueryDelegationsRequest): Promise<_103.QueryDelegationsResponse>;
                        receipts(request: _103.QueryReceiptsRequest): Promise<_103.QueryReceiptsResponse>;
                        zoneWithdrawalRecords(request: _103.QueryWithdrawalRecordsRequest): Promise<_103.QueryWithdrawalRecordsResponse>;
                        withdrawalRecords(request: _103.QueryWithdrawalRecordsRequest): Promise<_103.QueryWithdrawalRecordsResponse>;
                        unbondingRecords(request: _103.QueryUnbondingRecordsRequest): Promise<_103.QueryUnbondingRecordsResponse>;
                        redelegationRecords(request: _103.QueryRedelegationRecordsRequest): Promise<_103.QueryRedelegationRecordsResponse>;
                    };
                };
                mint: {
                    v1beta1: {
                        params(request?: _107.QueryParamsRequest): Promise<_107.QueryParamsResponse>;
                        epochProvisions(request?: _107.QueryEpochProvisionsRequest): Promise<_107.QueryEpochProvisionsResponse>;
                    };
                };
                participationrewards: {
                    v1: {
                        params(request?: _111.QueryParamsRequest): Promise<_111.QueryParamsResponse>;
                        protocolData(request: _111.QueryProtocolDataRequest): Promise<_111.QueryProtocolDataResponse>;
                    };
                };
                tokenfactory: {
                    v1beta1: {
                        params(request?: _115.QueryParamsRequest): Promise<_115.QueryParamsResponse>;
                        denomAuthorityMetadata(request: _115.QueryDenomAuthorityMetadataRequest): Promise<_115.QueryDenomAuthorityMetadataResponse>;
                        denomsFromCreator(request: _115.QueryDenomsFromCreatorRequest): Promise<_115.QueryDenomsFromCreatorResponse>;
                    };
                };
            };
        }>;
    };
}
