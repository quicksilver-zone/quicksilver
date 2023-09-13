import * as _1 from "./auth/v1beta1/auth";
import * as _2 from "./auth/v1beta1/genesis";
import * as _3 from "./auth/v1beta1/query";
import * as _4 from "./authz/v1beta1/authz";
import * as _5 from "./authz/v1beta1/event";
import * as _6 from "./authz/v1beta1/genesis";
import * as _7 from "./authz/v1beta1/query";
import * as _8 from "./authz/v1beta1/tx";
import * as _9 from "./bank/v1beta1/authz";
import * as _10 from "./bank/v1beta1/bank";
import * as _11 from "./bank/v1beta1/genesis";
import * as _12 from "./bank/v1beta1/query";
import * as _13 from "./bank/v1beta1/tx";
import * as _14 from "./base/abci/v1beta1/abci";
import * as _15 from "./base/kv/v1beta1/kv";
import * as _16 from "./base/query/v1beta1/pagination";
import * as _17 from "./base/reflection/v1beta1/reflection";
import * as _18 from "./base/reflection/v2alpha1/reflection";
import * as _19 from "./base/snapshots/v1beta1/snapshot";
import * as _20 from "./base/store/v1beta1/commit_info";
import * as _21 from "./base/store/v1beta1/listening";
import * as _22 from "./base/store/v1beta1/snapshot";
import * as _23 from "./base/tendermint/v1beta1/query";
import * as _24 from "./base/v1beta1/coin";
import * as _25 from "./capability/v1beta1/capability";
import * as _26 from "./capability/v1beta1/genesis";
import * as _27 from "./crisis/v1beta1/genesis";
import * as _28 from "./crisis/v1beta1/tx";
import * as _29 from "./crypto/ed25519/keys";
import * as _30 from "./crypto/multisig/keys";
import * as _31 from "./crypto/secp256k1/keys";
import * as _32 from "./crypto/secp256r1/keys";
import * as _33 from "./distribution/v1beta1/distribution";
import * as _34 from "./distribution/v1beta1/genesis";
import * as _35 from "./distribution/v1beta1/query";
import * as _36 from "./distribution/v1beta1/tx";
import * as _37 from "./evidence/v1beta1/evidence";
import * as _38 from "./evidence/v1beta1/genesis";
import * as _39 from "./evidence/v1beta1/query";
import * as _40 from "./evidence/v1beta1/tx";
import * as _41 from "./feegrant/v1beta1/feegrant";
import * as _42 from "./feegrant/v1beta1/genesis";
import * as _43 from "./feegrant/v1beta1/query";
import * as _44 from "./feegrant/v1beta1/tx";
import * as _45 from "./genutil/v1beta1/genesis";
import * as _46 from "./gov/v1beta1/genesis";
import * as _47 from "./gov/v1beta1/gov";
import * as _48 from "./gov/v1beta1/query";
import * as _49 from "./gov/v1beta1/tx";
import * as _50 from "./mint/v1beta1/genesis";
import * as _51 from "./mint/v1beta1/mint";
import * as _52 from "./mint/v1beta1/query";
import * as _53 from "./params/v1beta1/params";
import * as _54 from "./params/v1beta1/query";
import * as _55 from "./slashing/v1beta1/genesis";
import * as _56 from "./slashing/v1beta1/query";
import * as _57 from "./slashing/v1beta1/slashing";
import * as _58 from "./slashing/v1beta1/tx";
import * as _59 from "./staking/v1beta1/authz";
import * as _60 from "./staking/v1beta1/genesis";
import * as _61 from "./staking/v1beta1/query";
import * as _62 from "./staking/v1beta1/staking";
import * as _63 from "./staking/v1beta1/tx";
import * as _64 from "./tx/signing/v1beta1/signing";
import * as _65 from "./tx/v1beta1/service";
import * as _66 from "./tx/v1beta1/tx";
import * as _67 from "./upgrade/v1beta1/query";
import * as _68 from "./upgrade/v1beta1/upgrade";
import * as _69 from "./vesting/v1beta1/tx";
import * as _70 from "./vesting/v1beta1/vesting";
import * as _148 from "./auth/v1beta1/query.rpc.Query";
import * as _149 from "./authz/v1beta1/query.rpc.Query";
import * as _150 from "./bank/v1beta1/query.rpc.Query";
import * as _151 from "./base/tendermint/v1beta1/query.rpc.Service";
import * as _152 from "./distribution/v1beta1/query.rpc.Query";
import * as _153 from "./evidence/v1beta1/query.rpc.Query";
import * as _154 from "./feegrant/v1beta1/query.rpc.Query";
import * as _155 from "./gov/v1beta1/query.rpc.Query";
import * as _156 from "./mint/v1beta1/query.rpc.Query";
import * as _157 from "./params/v1beta1/query.rpc.Query";
import * as _158 from "./slashing/v1beta1/query.rpc.Query";
import * as _159 from "./staking/v1beta1/query.rpc.Query";
import * as _160 from "./tx/v1beta1/service.rpc.Service";
import * as _161 from "./upgrade/v1beta1/query.rpc.Query";
import * as _162 from "./authz/v1beta1/tx.rpc.msg";
import * as _163 from "./bank/v1beta1/tx.rpc.msg";
import * as _164 from "./crisis/v1beta1/tx.rpc.msg";
import * as _165 from "./distribution/v1beta1/tx.rpc.msg";
import * as _166 from "./evidence/v1beta1/tx.rpc.msg";
import * as _167 from "./feegrant/v1beta1/tx.rpc.msg";
import * as _168 from "./gov/v1beta1/tx.rpc.msg";
import * as _169 from "./slashing/v1beta1/tx.rpc.msg";
import * as _170 from "./staking/v1beta1/tx.rpc.msg";
import * as _171 from "./vesting/v1beta1/tx.rpc.msg";
export declare namespace cosmos {
    namespace auth {
        const v1beta1: {
            QueryClientImpl: typeof _148.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                accounts(request?: _3.QueryAccountsRequest): Promise<_3.QueryAccountsResponse>;
                account(request: _3.QueryAccountRequest): Promise<_3.QueryAccountResponse>;
                params(request?: _3.QueryParamsRequest): Promise<_3.QueryParamsResponse>;
            };
            QueryAccountsRequest: {
                encode(message: _3.QueryAccountsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _3.QueryAccountsRequest;
                fromJSON(object: any): _3.QueryAccountsRequest;
                toJSON(message: _3.QueryAccountsRequest): unknown;
                fromPartial(object: Partial<_3.QueryAccountsRequest>): _3.QueryAccountsRequest;
            };
            QueryAccountsResponse: {
                encode(message: _3.QueryAccountsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _3.QueryAccountsResponse;
                fromJSON(object: any): _3.QueryAccountsResponse;
                toJSON(message: _3.QueryAccountsResponse): unknown;
                fromPartial(object: Partial<_3.QueryAccountsResponse>): _3.QueryAccountsResponse;
            };
            QueryAccountRequest: {
                encode(message: _3.QueryAccountRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _3.QueryAccountRequest;
                fromJSON(object: any): _3.QueryAccountRequest;
                toJSON(message: _3.QueryAccountRequest): unknown;
                fromPartial(object: Partial<_3.QueryAccountRequest>): _3.QueryAccountRequest;
            };
            QueryAccountResponse: {
                encode(message: _3.QueryAccountResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _3.QueryAccountResponse;
                fromJSON(object: any): _3.QueryAccountResponse;
                toJSON(message: _3.QueryAccountResponse): unknown;
                fromPartial(object: Partial<_3.QueryAccountResponse>): _3.QueryAccountResponse;
            };
            QueryParamsRequest: {
                encode(_: _3.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _3.QueryParamsRequest;
                fromJSON(_: any): _3.QueryParamsRequest;
                toJSON(_: _3.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_3.QueryParamsRequest>): _3.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _3.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _3.QueryParamsResponse;
                fromJSON(object: any): _3.QueryParamsResponse;
                toJSON(message: _3.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_3.QueryParamsResponse>): _3.QueryParamsResponse;
            };
            GenesisState: {
                encode(message: _2.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _2.GenesisState;
                fromJSON(object: any): _2.GenesisState;
                toJSON(message: _2.GenesisState): unknown;
                fromPartial(object: Partial<_2.GenesisState>): _2.GenesisState;
            };
            BaseAccount: {
                encode(message: _1.BaseAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _1.BaseAccount;
                fromJSON(object: any): _1.BaseAccount;
                toJSON(message: _1.BaseAccount): unknown;
                fromPartial(object: Partial<_1.BaseAccount>): _1.BaseAccount;
            };
            ModuleAccount: {
                encode(message: _1.ModuleAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _1.ModuleAccount;
                fromJSON(object: any): _1.ModuleAccount;
                toJSON(message: _1.ModuleAccount): unknown;
                fromPartial(object: Partial<_1.ModuleAccount>): _1.ModuleAccount;
            };
            Params: {
                encode(message: _1.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _1.Params;
                fromJSON(object: any): _1.Params;
                toJSON(message: _1.Params): unknown;
                fromPartial(object: Partial<_1.Params>): _1.Params;
            };
        };
    }
    namespace authz {
        const v1beta1: {
            MsgClientImpl: typeof _162.MsgClientImpl;
            QueryClientImpl: typeof _149.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                grants(request: _7.QueryGrantsRequest): Promise<_7.QueryGrantsResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    grant(value: _8.MsgGrant): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    exec(value: _8.MsgExec): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    revoke(value: _8.MsgRevoke): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    grant(value: _8.MsgGrant): {
                        typeUrl: string;
                        value: _8.MsgGrant;
                    };
                    exec(value: _8.MsgExec): {
                        typeUrl: string;
                        value: _8.MsgExec;
                    };
                    revoke(value: _8.MsgRevoke): {
                        typeUrl: string;
                        value: _8.MsgRevoke;
                    };
                };
                toJSON: {
                    grant(value: _8.MsgGrant): {
                        typeUrl: string;
                        value: unknown;
                    };
                    exec(value: _8.MsgExec): {
                        typeUrl: string;
                        value: unknown;
                    };
                    revoke(value: _8.MsgRevoke): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    grant(value: any): {
                        typeUrl: string;
                        value: _8.MsgGrant;
                    };
                    exec(value: any): {
                        typeUrl: string;
                        value: _8.MsgExec;
                    };
                    revoke(value: any): {
                        typeUrl: string;
                        value: _8.MsgRevoke;
                    };
                };
                fromPartial: {
                    grant(value: _8.MsgGrant): {
                        typeUrl: string;
                        value: _8.MsgGrant;
                    };
                    exec(value: _8.MsgExec): {
                        typeUrl: string;
                        value: _8.MsgExec;
                    };
                    revoke(value: _8.MsgRevoke): {
                        typeUrl: string;
                        value: _8.MsgRevoke;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.authz.v1beta1.MsgGrant": {
                    aminoType: string;
                    toAmino: ({ granter, grantee, grant }: _8.MsgGrant) => {
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
                    }) => _8.MsgGrant;
                };
                "/cosmos.authz.v1beta1.MsgExec": {
                    aminoType: string;
                    toAmino: ({ grantee, msgs }: _8.MsgExec) => {
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
                    }) => _8.MsgExec;
                };
                "/cosmos.authz.v1beta1.MsgRevoke": {
                    aminoType: string;
                    toAmino: ({ granter, grantee, msgTypeUrl }: _8.MsgRevoke) => {
                        granter: string;
                        grantee: string;
                        msg_type_url: string;
                    };
                    fromAmino: ({ granter, grantee, msg_type_url }: {
                        granter: string;
                        grantee: string;
                        msg_type_url: string;
                    }) => _8.MsgRevoke;
                };
            };
            MsgGrant: {
                encode(message: _8.MsgGrant, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _8.MsgGrant;
                fromJSON(object: any): _8.MsgGrant;
                toJSON(message: _8.MsgGrant): unknown;
                fromPartial(object: Partial<_8.MsgGrant>): _8.MsgGrant;
            };
            MsgExecResponse: {
                encode(message: _8.MsgExecResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _8.MsgExecResponse;
                fromJSON(object: any): _8.MsgExecResponse;
                toJSON(message: _8.MsgExecResponse): unknown;
                fromPartial(object: Partial<_8.MsgExecResponse>): _8.MsgExecResponse;
            };
            MsgExec: {
                encode(message: _8.MsgExec, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _8.MsgExec;
                fromJSON(object: any): _8.MsgExec;
                toJSON(message: _8.MsgExec): unknown;
                fromPartial(object: Partial<_8.MsgExec>): _8.MsgExec;
            };
            MsgGrantResponse: {
                encode(_: _8.MsgGrantResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _8.MsgGrantResponse;
                fromJSON(_: any): _8.MsgGrantResponse;
                toJSON(_: _8.MsgGrantResponse): unknown;
                fromPartial(_: Partial<_8.MsgGrantResponse>): _8.MsgGrantResponse;
            };
            MsgRevoke: {
                encode(message: _8.MsgRevoke, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _8.MsgRevoke;
                fromJSON(object: any): _8.MsgRevoke;
                toJSON(message: _8.MsgRevoke): unknown;
                fromPartial(object: Partial<_8.MsgRevoke>): _8.MsgRevoke;
            };
            MsgRevokeResponse: {
                encode(_: _8.MsgRevokeResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _8.MsgRevokeResponse;
                fromJSON(_: any): _8.MsgRevokeResponse;
                toJSON(_: _8.MsgRevokeResponse): unknown;
                fromPartial(_: Partial<_8.MsgRevokeResponse>): _8.MsgRevokeResponse;
            };
            QueryGrantsRequest: {
                encode(message: _7.QueryGrantsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _7.QueryGrantsRequest;
                fromJSON(object: any): _7.QueryGrantsRequest;
                toJSON(message: _7.QueryGrantsRequest): unknown;
                fromPartial(object: Partial<_7.QueryGrantsRequest>): _7.QueryGrantsRequest;
            };
            QueryGrantsResponse: {
                encode(message: _7.QueryGrantsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _7.QueryGrantsResponse;
                fromJSON(object: any): _7.QueryGrantsResponse;
                toJSON(message: _7.QueryGrantsResponse): unknown;
                fromPartial(object: Partial<_7.QueryGrantsResponse>): _7.QueryGrantsResponse;
            };
            GenesisState: {
                encode(message: _6.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _6.GenesisState;
                fromJSON(object: any): _6.GenesisState;
                toJSON(message: _6.GenesisState): unknown;
                fromPartial(object: Partial<_6.GenesisState>): _6.GenesisState;
            };
            GrantAuthorization: {
                encode(message: _6.GrantAuthorization, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _6.GrantAuthorization;
                fromJSON(object: any): _6.GrantAuthorization;
                toJSON(message: _6.GrantAuthorization): unknown;
                fromPartial(object: Partial<_6.GrantAuthorization>): _6.GrantAuthorization;
            };
            EventGrant: {
                encode(message: _5.EventGrant, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _5.EventGrant;
                fromJSON(object: any): _5.EventGrant;
                toJSON(message: _5.EventGrant): unknown;
                fromPartial(object: Partial<_5.EventGrant>): _5.EventGrant;
            };
            EventRevoke: {
                encode(message: _5.EventRevoke, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _5.EventRevoke;
                fromJSON(object: any): _5.EventRevoke;
                toJSON(message: _5.EventRevoke): unknown;
                fromPartial(object: Partial<_5.EventRevoke>): _5.EventRevoke;
            };
            GenericAuthorization: {
                encode(message: _4.GenericAuthorization, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _4.GenericAuthorization;
                fromJSON(object: any): _4.GenericAuthorization;
                toJSON(message: _4.GenericAuthorization): unknown;
                fromPartial(object: Partial<_4.GenericAuthorization>): _4.GenericAuthorization;
            };
            Grant: {
                encode(message: _4.Grant, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _4.Grant;
                fromJSON(object: any): _4.Grant;
                toJSON(message: _4.Grant): unknown;
                fromPartial(object: Partial<_4.Grant>): _4.Grant;
            };
        };
    }
    namespace bank {
        const v1beta1: {
            MsgClientImpl: typeof _163.MsgClientImpl;
            QueryClientImpl: typeof _150.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                balance(request: _12.QueryBalanceRequest): Promise<_12.QueryBalanceResponse>;
                allBalances(request: _12.QueryAllBalancesRequest): Promise<_12.QueryAllBalancesResponse>;
                totalSupply(request?: _12.QueryTotalSupplyRequest): Promise<_12.QueryTotalSupplyResponse>;
                supplyOf(request: _12.QuerySupplyOfRequest): Promise<_12.QuerySupplyOfResponse>;
                params(request?: _12.QueryParamsRequest): Promise<_12.QueryParamsResponse>;
                denomMetadata(request: _12.QueryDenomMetadataRequest): Promise<_12.QueryDenomMetadataResponse>;
                denomsMetadata(request?: _12.QueryDenomsMetadataRequest): Promise<_12.QueryDenomsMetadataResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    send(value: _13.MsgSend): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    multiSend(value: _13.MsgMultiSend): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    send(value: _13.MsgSend): {
                        typeUrl: string;
                        value: _13.MsgSend;
                    };
                    multiSend(value: _13.MsgMultiSend): {
                        typeUrl: string;
                        value: _13.MsgMultiSend;
                    };
                };
                toJSON: {
                    send(value: _13.MsgSend): {
                        typeUrl: string;
                        value: unknown;
                    };
                    multiSend(value: _13.MsgMultiSend): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    send(value: any): {
                        typeUrl: string;
                        value: _13.MsgSend;
                    };
                    multiSend(value: any): {
                        typeUrl: string;
                        value: _13.MsgMultiSend;
                    };
                };
                fromPartial: {
                    send(value: _13.MsgSend): {
                        typeUrl: string;
                        value: _13.MsgSend;
                    };
                    multiSend(value: _13.MsgMultiSend): {
                        typeUrl: string;
                        value: _13.MsgMultiSend;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.bank.v1beta1.MsgSend": {
                    aminoType: string;
                    toAmino: ({ fromAddress, toAddress, amount }: _13.MsgSend) => {
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
                    }) => _13.MsgSend;
                };
                "/cosmos.bank.v1beta1.MsgMultiSend": {
                    aminoType: string;
                    toAmino: ({ inputs, outputs }: _13.MsgMultiSend) => {
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
                    }) => _13.MsgMultiSend;
                };
            };
            MsgSend: {
                encode(message: _13.MsgSend, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _13.MsgSend;
                fromJSON(object: any): _13.MsgSend;
                toJSON(message: _13.MsgSend): unknown;
                fromPartial(object: Partial<_13.MsgSend>): _13.MsgSend;
            };
            MsgSendResponse: {
                encode(_: _13.MsgSendResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _13.MsgSendResponse;
                fromJSON(_: any): _13.MsgSendResponse;
                toJSON(_: _13.MsgSendResponse): unknown;
                fromPartial(_: Partial<_13.MsgSendResponse>): _13.MsgSendResponse;
            };
            MsgMultiSend: {
                encode(message: _13.MsgMultiSend, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _13.MsgMultiSend;
                fromJSON(object: any): _13.MsgMultiSend;
                toJSON(message: _13.MsgMultiSend): unknown;
                fromPartial(object: Partial<_13.MsgMultiSend>): _13.MsgMultiSend;
            };
            MsgMultiSendResponse: {
                encode(_: _13.MsgMultiSendResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _13.MsgMultiSendResponse;
                fromJSON(_: any): _13.MsgMultiSendResponse;
                toJSON(_: _13.MsgMultiSendResponse): unknown;
                fromPartial(_: Partial<_13.MsgMultiSendResponse>): _13.MsgMultiSendResponse;
            };
            QueryBalanceRequest: {
                encode(message: _12.QueryBalanceRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryBalanceRequest;
                fromJSON(object: any): _12.QueryBalanceRequest;
                toJSON(message: _12.QueryBalanceRequest): unknown;
                fromPartial(object: Partial<_12.QueryBalanceRequest>): _12.QueryBalanceRequest;
            };
            QueryBalanceResponse: {
                encode(message: _12.QueryBalanceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryBalanceResponse;
                fromJSON(object: any): _12.QueryBalanceResponse;
                toJSON(message: _12.QueryBalanceResponse): unknown;
                fromPartial(object: Partial<_12.QueryBalanceResponse>): _12.QueryBalanceResponse;
            };
            QueryAllBalancesRequest: {
                encode(message: _12.QueryAllBalancesRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryAllBalancesRequest;
                fromJSON(object: any): _12.QueryAllBalancesRequest;
                toJSON(message: _12.QueryAllBalancesRequest): unknown;
                fromPartial(object: Partial<_12.QueryAllBalancesRequest>): _12.QueryAllBalancesRequest;
            };
            QueryAllBalancesResponse: {
                encode(message: _12.QueryAllBalancesResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryAllBalancesResponse;
                fromJSON(object: any): _12.QueryAllBalancesResponse;
                toJSON(message: _12.QueryAllBalancesResponse): unknown;
                fromPartial(object: Partial<_12.QueryAllBalancesResponse>): _12.QueryAllBalancesResponse;
            };
            QueryTotalSupplyRequest: {
                encode(message: _12.QueryTotalSupplyRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryTotalSupplyRequest;
                fromJSON(object: any): _12.QueryTotalSupplyRequest;
                toJSON(message: _12.QueryTotalSupplyRequest): unknown;
                fromPartial(object: Partial<_12.QueryTotalSupplyRequest>): _12.QueryTotalSupplyRequest;
            };
            QueryTotalSupplyResponse: {
                encode(message: _12.QueryTotalSupplyResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryTotalSupplyResponse;
                fromJSON(object: any): _12.QueryTotalSupplyResponse;
                toJSON(message: _12.QueryTotalSupplyResponse): unknown;
                fromPartial(object: Partial<_12.QueryTotalSupplyResponse>): _12.QueryTotalSupplyResponse;
            };
            QuerySupplyOfRequest: {
                encode(message: _12.QuerySupplyOfRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QuerySupplyOfRequest;
                fromJSON(object: any): _12.QuerySupplyOfRequest;
                toJSON(message: _12.QuerySupplyOfRequest): unknown;
                fromPartial(object: Partial<_12.QuerySupplyOfRequest>): _12.QuerySupplyOfRequest;
            };
            QuerySupplyOfResponse: {
                encode(message: _12.QuerySupplyOfResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QuerySupplyOfResponse;
                fromJSON(object: any): _12.QuerySupplyOfResponse;
                toJSON(message: _12.QuerySupplyOfResponse): unknown;
                fromPartial(object: Partial<_12.QuerySupplyOfResponse>): _12.QuerySupplyOfResponse;
            };
            QueryParamsRequest: {
                encode(_: _12.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryParamsRequest;
                fromJSON(_: any): _12.QueryParamsRequest;
                toJSON(_: _12.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_12.QueryParamsRequest>): _12.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _12.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryParamsResponse;
                fromJSON(object: any): _12.QueryParamsResponse;
                toJSON(message: _12.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_12.QueryParamsResponse>): _12.QueryParamsResponse;
            };
            QueryDenomsMetadataRequest: {
                encode(message: _12.QueryDenomsMetadataRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryDenomsMetadataRequest;
                fromJSON(object: any): _12.QueryDenomsMetadataRequest;
                toJSON(message: _12.QueryDenomsMetadataRequest): unknown;
                fromPartial(object: Partial<_12.QueryDenomsMetadataRequest>): _12.QueryDenomsMetadataRequest;
            };
            QueryDenomsMetadataResponse: {
                encode(message: _12.QueryDenomsMetadataResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryDenomsMetadataResponse;
                fromJSON(object: any): _12.QueryDenomsMetadataResponse;
                toJSON(message: _12.QueryDenomsMetadataResponse): unknown;
                fromPartial(object: Partial<_12.QueryDenomsMetadataResponse>): _12.QueryDenomsMetadataResponse;
            };
            QueryDenomMetadataRequest: {
                encode(message: _12.QueryDenomMetadataRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryDenomMetadataRequest;
                fromJSON(object: any): _12.QueryDenomMetadataRequest;
                toJSON(message: _12.QueryDenomMetadataRequest): unknown;
                fromPartial(object: Partial<_12.QueryDenomMetadataRequest>): _12.QueryDenomMetadataRequest;
            };
            QueryDenomMetadataResponse: {
                encode(message: _12.QueryDenomMetadataResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _12.QueryDenomMetadataResponse;
                fromJSON(object: any): _12.QueryDenomMetadataResponse;
                toJSON(message: _12.QueryDenomMetadataResponse): unknown;
                fromPartial(object: Partial<_12.QueryDenomMetadataResponse>): _12.QueryDenomMetadataResponse;
            };
            GenesisState: {
                encode(message: _11.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _11.GenesisState;
                fromJSON(object: any): _11.GenesisState;
                toJSON(message: _11.GenesisState): unknown;
                fromPartial(object: Partial<_11.GenesisState>): _11.GenesisState;
            };
            Balance: {
                encode(message: _11.Balance, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _11.Balance;
                fromJSON(object: any): _11.Balance;
                toJSON(message: _11.Balance): unknown;
                fromPartial(object: Partial<_11.Balance>): _11.Balance;
            };
            Params: {
                encode(message: _10.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _10.Params;
                fromJSON(object: any): _10.Params;
                toJSON(message: _10.Params): unknown;
                fromPartial(object: Partial<_10.Params>): _10.Params;
            };
            SendEnabled: {
                encode(message: _10.SendEnabled, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _10.SendEnabled;
                fromJSON(object: any): _10.SendEnabled;
                toJSON(message: _10.SendEnabled): unknown;
                fromPartial(object: Partial<_10.SendEnabled>): _10.SendEnabled;
            };
            Input: {
                encode(message: _10.Input, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _10.Input;
                fromJSON(object: any): _10.Input;
                toJSON(message: _10.Input): unknown;
                fromPartial(object: Partial<_10.Input>): _10.Input;
            };
            Output: {
                encode(message: _10.Output, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _10.Output;
                fromJSON(object: any): _10.Output;
                toJSON(message: _10.Output): unknown;
                fromPartial(object: Partial<_10.Output>): _10.Output;
            };
            Supply: {
                encode(message: _10.Supply, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _10.Supply;
                fromJSON(object: any): _10.Supply;
                toJSON(message: _10.Supply): unknown;
                fromPartial(object: Partial<_10.Supply>): _10.Supply;
            };
            DenomUnit: {
                encode(message: _10.DenomUnit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _10.DenomUnit;
                fromJSON(object: any): _10.DenomUnit;
                toJSON(message: _10.DenomUnit): unknown;
                fromPartial(object: Partial<_10.DenomUnit>): _10.DenomUnit;
            };
            Metadata: {
                encode(message: _10.Metadata, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _10.Metadata;
                fromJSON(object: any): _10.Metadata;
                toJSON(message: _10.Metadata): unknown;
                fromPartial(object: Partial<_10.Metadata>): _10.Metadata;
            };
            SendAuthorization: {
                encode(message: _9.SendAuthorization, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _9.SendAuthorization;
                fromJSON(object: any): _9.SendAuthorization;
                toJSON(message: _9.SendAuthorization): unknown;
                fromPartial(object: Partial<_9.SendAuthorization>): _9.SendAuthorization;
            };
        };
    }
    namespace base {
        namespace abci {
            const v1beta1: {
                TxResponse: {
                    encode(message: _14.TxResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.TxResponse;
                    fromJSON(object: any): _14.TxResponse;
                    toJSON(message: _14.TxResponse): unknown;
                    fromPartial(object: Partial<_14.TxResponse>): _14.TxResponse;
                };
                ABCIMessageLog: {
                    encode(message: _14.ABCIMessageLog, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.ABCIMessageLog;
                    fromJSON(object: any): _14.ABCIMessageLog;
                    toJSON(message: _14.ABCIMessageLog): unknown;
                    fromPartial(object: Partial<_14.ABCIMessageLog>): _14.ABCIMessageLog;
                };
                StringEvent: {
                    encode(message: _14.StringEvent, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.StringEvent;
                    fromJSON(object: any): _14.StringEvent;
                    toJSON(message: _14.StringEvent): unknown;
                    fromPartial(object: Partial<_14.StringEvent>): _14.StringEvent;
                };
                Attribute: {
                    encode(message: _14.Attribute, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.Attribute;
                    fromJSON(object: any): _14.Attribute;
                    toJSON(message: _14.Attribute): unknown;
                    fromPartial(object: Partial<_14.Attribute>): _14.Attribute;
                };
                GasInfo: {
                    encode(message: _14.GasInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.GasInfo;
                    fromJSON(object: any): _14.GasInfo;
                    toJSON(message: _14.GasInfo): unknown;
                    fromPartial(object: Partial<_14.GasInfo>): _14.GasInfo;
                };
                Result: {
                    encode(message: _14.Result, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.Result;
                    fromJSON(object: any): _14.Result;
                    toJSON(message: _14.Result): unknown;
                    fromPartial(object: Partial<_14.Result>): _14.Result;
                };
                SimulationResponse: {
                    encode(message: _14.SimulationResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.SimulationResponse;
                    fromJSON(object: any): _14.SimulationResponse;
                    toJSON(message: _14.SimulationResponse): unknown;
                    fromPartial(object: Partial<_14.SimulationResponse>): _14.SimulationResponse;
                };
                MsgData: {
                    encode(message: _14.MsgData, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.MsgData;
                    fromJSON(object: any): _14.MsgData;
                    toJSON(message: _14.MsgData): unknown;
                    fromPartial(object: Partial<_14.MsgData>): _14.MsgData;
                };
                TxMsgData: {
                    encode(message: _14.TxMsgData, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.TxMsgData;
                    fromJSON(object: any): _14.TxMsgData;
                    toJSON(message: _14.TxMsgData): unknown;
                    fromPartial(object: Partial<_14.TxMsgData>): _14.TxMsgData;
                };
                SearchTxsResult: {
                    encode(message: _14.SearchTxsResult, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _14.SearchTxsResult;
                    fromJSON(object: any): _14.SearchTxsResult;
                    toJSON(message: _14.SearchTxsResult): unknown;
                    fromPartial(object: Partial<_14.SearchTxsResult>): _14.SearchTxsResult;
                };
            };
        }
        namespace kv {
            const v1beta1: {
                Pairs: {
                    encode(message: _15.Pairs, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _15.Pairs;
                    fromJSON(object: any): _15.Pairs;
                    toJSON(message: _15.Pairs): unknown;
                    fromPartial(object: Partial<_15.Pairs>): _15.Pairs;
                };
                Pair: {
                    encode(message: _15.Pair, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _15.Pair;
                    fromJSON(object: any): _15.Pair;
                    toJSON(message: _15.Pair): unknown;
                    fromPartial(object: Partial<_15.Pair>): _15.Pair;
                };
            };
        }
        namespace query {
            const v1beta1: {
                PageRequest: {
                    encode(message: _16.PageRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _16.PageRequest;
                    fromJSON(object: any): _16.PageRequest;
                    toJSON(message: _16.PageRequest): unknown;
                    fromPartial(object: Partial<_16.PageRequest>): _16.PageRequest;
                };
                PageResponse: {
                    encode(message: _16.PageResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _16.PageResponse;
                    fromJSON(object: any): _16.PageResponse;
                    toJSON(message: _16.PageResponse): unknown;
                    fromPartial(object: Partial<_16.PageResponse>): _16.PageResponse;
                };
            };
        }
        namespace reflection {
            const v1beta1: {
                ListAllInterfacesRequest: {
                    encode(_: _17.ListAllInterfacesRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _17.ListAllInterfacesRequest;
                    fromJSON(_: any): _17.ListAllInterfacesRequest;
                    toJSON(_: _17.ListAllInterfacesRequest): unknown;
                    fromPartial(_: Partial<_17.ListAllInterfacesRequest>): _17.ListAllInterfacesRequest;
                };
                ListAllInterfacesResponse: {
                    encode(message: _17.ListAllInterfacesResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _17.ListAllInterfacesResponse;
                    fromJSON(object: any): _17.ListAllInterfacesResponse;
                    toJSON(message: _17.ListAllInterfacesResponse): unknown;
                    fromPartial(object: Partial<_17.ListAllInterfacesResponse>): _17.ListAllInterfacesResponse;
                };
                ListImplementationsRequest: {
                    encode(message: _17.ListImplementationsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _17.ListImplementationsRequest;
                    fromJSON(object: any): _17.ListImplementationsRequest;
                    toJSON(message: _17.ListImplementationsRequest): unknown;
                    fromPartial(object: Partial<_17.ListImplementationsRequest>): _17.ListImplementationsRequest;
                };
                ListImplementationsResponse: {
                    encode(message: _17.ListImplementationsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _17.ListImplementationsResponse;
                    fromJSON(object: any): _17.ListImplementationsResponse;
                    toJSON(message: _17.ListImplementationsResponse): unknown;
                    fromPartial(object: Partial<_17.ListImplementationsResponse>): _17.ListImplementationsResponse;
                };
            };
            const v2alpha1: {
                AppDescriptor: {
                    encode(message: _18.AppDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.AppDescriptor;
                    fromJSON(object: any): _18.AppDescriptor;
                    toJSON(message: _18.AppDescriptor): unknown;
                    fromPartial(object: Partial<_18.AppDescriptor>): _18.AppDescriptor;
                };
                TxDescriptor: {
                    encode(message: _18.TxDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.TxDescriptor;
                    fromJSON(object: any): _18.TxDescriptor;
                    toJSON(message: _18.TxDescriptor): unknown;
                    fromPartial(object: Partial<_18.TxDescriptor>): _18.TxDescriptor;
                };
                AuthnDescriptor: {
                    encode(message: _18.AuthnDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.AuthnDescriptor;
                    fromJSON(object: any): _18.AuthnDescriptor;
                    toJSON(message: _18.AuthnDescriptor): unknown;
                    fromPartial(object: Partial<_18.AuthnDescriptor>): _18.AuthnDescriptor;
                };
                SigningModeDescriptor: {
                    encode(message: _18.SigningModeDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.SigningModeDescriptor;
                    fromJSON(object: any): _18.SigningModeDescriptor;
                    toJSON(message: _18.SigningModeDescriptor): unknown;
                    fromPartial(object: Partial<_18.SigningModeDescriptor>): _18.SigningModeDescriptor;
                };
                ChainDescriptor: {
                    encode(message: _18.ChainDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.ChainDescriptor;
                    fromJSON(object: any): _18.ChainDescriptor;
                    toJSON(message: _18.ChainDescriptor): unknown;
                    fromPartial(object: Partial<_18.ChainDescriptor>): _18.ChainDescriptor;
                };
                CodecDescriptor: {
                    encode(message: _18.CodecDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.CodecDescriptor;
                    fromJSON(object: any): _18.CodecDescriptor;
                    toJSON(message: _18.CodecDescriptor): unknown;
                    fromPartial(object: Partial<_18.CodecDescriptor>): _18.CodecDescriptor;
                };
                InterfaceDescriptor: {
                    encode(message: _18.InterfaceDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.InterfaceDescriptor;
                    fromJSON(object: any): _18.InterfaceDescriptor;
                    toJSON(message: _18.InterfaceDescriptor): unknown;
                    fromPartial(object: Partial<_18.InterfaceDescriptor>): _18.InterfaceDescriptor;
                };
                InterfaceImplementerDescriptor: {
                    encode(message: _18.InterfaceImplementerDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.InterfaceImplementerDescriptor;
                    fromJSON(object: any): _18.InterfaceImplementerDescriptor;
                    toJSON(message: _18.InterfaceImplementerDescriptor): unknown;
                    fromPartial(object: Partial<_18.InterfaceImplementerDescriptor>): _18.InterfaceImplementerDescriptor;
                };
                InterfaceAcceptingMessageDescriptor: {
                    encode(message: _18.InterfaceAcceptingMessageDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.InterfaceAcceptingMessageDescriptor;
                    fromJSON(object: any): _18.InterfaceAcceptingMessageDescriptor;
                    toJSON(message: _18.InterfaceAcceptingMessageDescriptor): unknown;
                    fromPartial(object: Partial<_18.InterfaceAcceptingMessageDescriptor>): _18.InterfaceAcceptingMessageDescriptor;
                };
                ConfigurationDescriptor: {
                    encode(message: _18.ConfigurationDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.ConfigurationDescriptor;
                    fromJSON(object: any): _18.ConfigurationDescriptor;
                    toJSON(message: _18.ConfigurationDescriptor): unknown;
                    fromPartial(object: Partial<_18.ConfigurationDescriptor>): _18.ConfigurationDescriptor;
                };
                MsgDescriptor: {
                    encode(message: _18.MsgDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.MsgDescriptor;
                    fromJSON(object: any): _18.MsgDescriptor;
                    toJSON(message: _18.MsgDescriptor): unknown;
                    fromPartial(object: Partial<_18.MsgDescriptor>): _18.MsgDescriptor;
                };
                GetAuthnDescriptorRequest: {
                    encode(_: _18.GetAuthnDescriptorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetAuthnDescriptorRequest;
                    fromJSON(_: any): _18.GetAuthnDescriptorRequest;
                    toJSON(_: _18.GetAuthnDescriptorRequest): unknown;
                    fromPartial(_: Partial<_18.GetAuthnDescriptorRequest>): _18.GetAuthnDescriptorRequest;
                };
                GetAuthnDescriptorResponse: {
                    encode(message: _18.GetAuthnDescriptorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetAuthnDescriptorResponse;
                    fromJSON(object: any): _18.GetAuthnDescriptorResponse;
                    toJSON(message: _18.GetAuthnDescriptorResponse): unknown;
                    fromPartial(object: Partial<_18.GetAuthnDescriptorResponse>): _18.GetAuthnDescriptorResponse;
                };
                GetChainDescriptorRequest: {
                    encode(_: _18.GetChainDescriptorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetChainDescriptorRequest;
                    fromJSON(_: any): _18.GetChainDescriptorRequest;
                    toJSON(_: _18.GetChainDescriptorRequest): unknown;
                    fromPartial(_: Partial<_18.GetChainDescriptorRequest>): _18.GetChainDescriptorRequest;
                };
                GetChainDescriptorResponse: {
                    encode(message: _18.GetChainDescriptorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetChainDescriptorResponse;
                    fromJSON(object: any): _18.GetChainDescriptorResponse;
                    toJSON(message: _18.GetChainDescriptorResponse): unknown;
                    fromPartial(object: Partial<_18.GetChainDescriptorResponse>): _18.GetChainDescriptorResponse;
                };
                GetCodecDescriptorRequest: {
                    encode(_: _18.GetCodecDescriptorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetCodecDescriptorRequest;
                    fromJSON(_: any): _18.GetCodecDescriptorRequest;
                    toJSON(_: _18.GetCodecDescriptorRequest): unknown;
                    fromPartial(_: Partial<_18.GetCodecDescriptorRequest>): _18.GetCodecDescriptorRequest;
                };
                GetCodecDescriptorResponse: {
                    encode(message: _18.GetCodecDescriptorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetCodecDescriptorResponse;
                    fromJSON(object: any): _18.GetCodecDescriptorResponse;
                    toJSON(message: _18.GetCodecDescriptorResponse): unknown;
                    fromPartial(object: Partial<_18.GetCodecDescriptorResponse>): _18.GetCodecDescriptorResponse;
                };
                GetConfigurationDescriptorRequest: {
                    encode(_: _18.GetConfigurationDescriptorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetConfigurationDescriptorRequest;
                    fromJSON(_: any): _18.GetConfigurationDescriptorRequest;
                    toJSON(_: _18.GetConfigurationDescriptorRequest): unknown;
                    fromPartial(_: Partial<_18.GetConfigurationDescriptorRequest>): _18.GetConfigurationDescriptorRequest;
                };
                GetConfigurationDescriptorResponse: {
                    encode(message: _18.GetConfigurationDescriptorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetConfigurationDescriptorResponse;
                    fromJSON(object: any): _18.GetConfigurationDescriptorResponse;
                    toJSON(message: _18.GetConfigurationDescriptorResponse): unknown;
                    fromPartial(object: Partial<_18.GetConfigurationDescriptorResponse>): _18.GetConfigurationDescriptorResponse;
                };
                GetQueryServicesDescriptorRequest: {
                    encode(_: _18.GetQueryServicesDescriptorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetQueryServicesDescriptorRequest;
                    fromJSON(_: any): _18.GetQueryServicesDescriptorRequest;
                    toJSON(_: _18.GetQueryServicesDescriptorRequest): unknown;
                    fromPartial(_: Partial<_18.GetQueryServicesDescriptorRequest>): _18.GetQueryServicesDescriptorRequest;
                };
                GetQueryServicesDescriptorResponse: {
                    encode(message: _18.GetQueryServicesDescriptorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetQueryServicesDescriptorResponse;
                    fromJSON(object: any): _18.GetQueryServicesDescriptorResponse;
                    toJSON(message: _18.GetQueryServicesDescriptorResponse): unknown;
                    fromPartial(object: Partial<_18.GetQueryServicesDescriptorResponse>): _18.GetQueryServicesDescriptorResponse;
                };
                GetTxDescriptorRequest: {
                    encode(_: _18.GetTxDescriptorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetTxDescriptorRequest;
                    fromJSON(_: any): _18.GetTxDescriptorRequest;
                    toJSON(_: _18.GetTxDescriptorRequest): unknown;
                    fromPartial(_: Partial<_18.GetTxDescriptorRequest>): _18.GetTxDescriptorRequest;
                };
                GetTxDescriptorResponse: {
                    encode(message: _18.GetTxDescriptorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.GetTxDescriptorResponse;
                    fromJSON(object: any): _18.GetTxDescriptorResponse;
                    toJSON(message: _18.GetTxDescriptorResponse): unknown;
                    fromPartial(object: Partial<_18.GetTxDescriptorResponse>): _18.GetTxDescriptorResponse;
                };
                QueryServicesDescriptor: {
                    encode(message: _18.QueryServicesDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.QueryServicesDescriptor;
                    fromJSON(object: any): _18.QueryServicesDescriptor;
                    toJSON(message: _18.QueryServicesDescriptor): unknown;
                    fromPartial(object: Partial<_18.QueryServicesDescriptor>): _18.QueryServicesDescriptor;
                };
                QueryServiceDescriptor: {
                    encode(message: _18.QueryServiceDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.QueryServiceDescriptor;
                    fromJSON(object: any): _18.QueryServiceDescriptor;
                    toJSON(message: _18.QueryServiceDescriptor): unknown;
                    fromPartial(object: Partial<_18.QueryServiceDescriptor>): _18.QueryServiceDescriptor;
                };
                QueryMethodDescriptor: {
                    encode(message: _18.QueryMethodDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _18.QueryMethodDescriptor;
                    fromJSON(object: any): _18.QueryMethodDescriptor;
                    toJSON(message: _18.QueryMethodDescriptor): unknown;
                    fromPartial(object: Partial<_18.QueryMethodDescriptor>): _18.QueryMethodDescriptor;
                };
            };
        }
        namespace snapshots {
            const v1beta1: {
                Snapshot: {
                    encode(message: _19.Snapshot, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _19.Snapshot;
                    fromJSON(object: any): _19.Snapshot;
                    toJSON(message: _19.Snapshot): unknown;
                    fromPartial(object: Partial<_19.Snapshot>): _19.Snapshot;
                };
                Metadata: {
                    encode(message: _19.Metadata, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _19.Metadata;
                    fromJSON(object: any): _19.Metadata;
                    toJSON(message: _19.Metadata): unknown;
                    fromPartial(object: Partial<_19.Metadata>): _19.Metadata;
                };
            };
        }
        namespace store {
            const v1beta1: {
                SnapshotItem: {
                    encode(message: _22.SnapshotItem, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _22.SnapshotItem;
                    fromJSON(object: any): _22.SnapshotItem;
                    toJSON(message: _22.SnapshotItem): unknown;
                    fromPartial(object: Partial<_22.SnapshotItem>): _22.SnapshotItem;
                };
                SnapshotStoreItem: {
                    encode(message: _22.SnapshotStoreItem, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _22.SnapshotStoreItem;
                    fromJSON(object: any): _22.SnapshotStoreItem;
                    toJSON(message: _22.SnapshotStoreItem): unknown;
                    fromPartial(object: Partial<_22.SnapshotStoreItem>): _22.SnapshotStoreItem;
                };
                SnapshotIAVLItem: {
                    encode(message: _22.SnapshotIAVLItem, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _22.SnapshotIAVLItem;
                    fromJSON(object: any): _22.SnapshotIAVLItem;
                    toJSON(message: _22.SnapshotIAVLItem): unknown;
                    fromPartial(object: Partial<_22.SnapshotIAVLItem>): _22.SnapshotIAVLItem;
                };
                StoreKVPair: {
                    encode(message: _21.StoreKVPair, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _21.StoreKVPair;
                    fromJSON(object: any): _21.StoreKVPair;
                    toJSON(message: _21.StoreKVPair): unknown;
                    fromPartial(object: Partial<_21.StoreKVPair>): _21.StoreKVPair;
                };
                CommitInfo: {
                    encode(message: _20.CommitInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _20.CommitInfo;
                    fromJSON(object: any): _20.CommitInfo;
                    toJSON(message: _20.CommitInfo): unknown;
                    fromPartial(object: Partial<_20.CommitInfo>): _20.CommitInfo;
                };
                StoreInfo: {
                    encode(message: _20.StoreInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _20.StoreInfo;
                    fromJSON(object: any): _20.StoreInfo;
                    toJSON(message: _20.StoreInfo): unknown;
                    fromPartial(object: Partial<_20.StoreInfo>): _20.StoreInfo;
                };
                CommitID: {
                    encode(message: _20.CommitID, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _20.CommitID;
                    fromJSON(object: any): _20.CommitID;
                    toJSON(message: _20.CommitID): unknown;
                    fromPartial(object: Partial<_20.CommitID>): _20.CommitID;
                };
            };
        }
        namespace tendermint {
            const v1beta1: {
                ServiceClientImpl: typeof _151.ServiceClientImpl;
                createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                    getNodeInfo(request?: _23.GetNodeInfoRequest): Promise<_23.GetNodeInfoResponse>;
                    getSyncing(request?: _23.GetSyncingRequest): Promise<_23.GetSyncingResponse>;
                    getLatestBlock(request?: _23.GetLatestBlockRequest): Promise<_23.GetLatestBlockResponse>;
                    getBlockByHeight(request: _23.GetBlockByHeightRequest): Promise<_23.GetBlockByHeightResponse>;
                    getLatestValidatorSet(request?: _23.GetLatestValidatorSetRequest): Promise<_23.GetLatestValidatorSetResponse>;
                    getValidatorSetByHeight(request: _23.GetValidatorSetByHeightRequest): Promise<_23.GetValidatorSetByHeightResponse>;
                };
                GetValidatorSetByHeightRequest: {
                    encode(message: _23.GetValidatorSetByHeightRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetValidatorSetByHeightRequest;
                    fromJSON(object: any): _23.GetValidatorSetByHeightRequest;
                    toJSON(message: _23.GetValidatorSetByHeightRequest): unknown;
                    fromPartial(object: Partial<_23.GetValidatorSetByHeightRequest>): _23.GetValidatorSetByHeightRequest;
                };
                GetValidatorSetByHeightResponse: {
                    encode(message: _23.GetValidatorSetByHeightResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetValidatorSetByHeightResponse;
                    fromJSON(object: any): _23.GetValidatorSetByHeightResponse;
                    toJSON(message: _23.GetValidatorSetByHeightResponse): unknown;
                    fromPartial(object: Partial<_23.GetValidatorSetByHeightResponse>): _23.GetValidatorSetByHeightResponse;
                };
                GetLatestValidatorSetRequest: {
                    encode(message: _23.GetLatestValidatorSetRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetLatestValidatorSetRequest;
                    fromJSON(object: any): _23.GetLatestValidatorSetRequest;
                    toJSON(message: _23.GetLatestValidatorSetRequest): unknown;
                    fromPartial(object: Partial<_23.GetLatestValidatorSetRequest>): _23.GetLatestValidatorSetRequest;
                };
                GetLatestValidatorSetResponse: {
                    encode(message: _23.GetLatestValidatorSetResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetLatestValidatorSetResponse;
                    fromJSON(object: any): _23.GetLatestValidatorSetResponse;
                    toJSON(message: _23.GetLatestValidatorSetResponse): unknown;
                    fromPartial(object: Partial<_23.GetLatestValidatorSetResponse>): _23.GetLatestValidatorSetResponse;
                };
                Validator: {
                    encode(message: _23.Validator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.Validator;
                    fromJSON(object: any): _23.Validator;
                    toJSON(message: _23.Validator): unknown;
                    fromPartial(object: Partial<_23.Validator>): _23.Validator;
                };
                GetBlockByHeightRequest: {
                    encode(message: _23.GetBlockByHeightRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetBlockByHeightRequest;
                    fromJSON(object: any): _23.GetBlockByHeightRequest;
                    toJSON(message: _23.GetBlockByHeightRequest): unknown;
                    fromPartial(object: Partial<_23.GetBlockByHeightRequest>): _23.GetBlockByHeightRequest;
                };
                GetBlockByHeightResponse: {
                    encode(message: _23.GetBlockByHeightResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetBlockByHeightResponse;
                    fromJSON(object: any): _23.GetBlockByHeightResponse;
                    toJSON(message: _23.GetBlockByHeightResponse): unknown;
                    fromPartial(object: Partial<_23.GetBlockByHeightResponse>): _23.GetBlockByHeightResponse;
                };
                GetLatestBlockRequest: {
                    encode(_: _23.GetLatestBlockRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetLatestBlockRequest;
                    fromJSON(_: any): _23.GetLatestBlockRequest;
                    toJSON(_: _23.GetLatestBlockRequest): unknown;
                    fromPartial(_: Partial<_23.GetLatestBlockRequest>): _23.GetLatestBlockRequest;
                };
                GetLatestBlockResponse: {
                    encode(message: _23.GetLatestBlockResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetLatestBlockResponse;
                    fromJSON(object: any): _23.GetLatestBlockResponse;
                    toJSON(message: _23.GetLatestBlockResponse): unknown;
                    fromPartial(object: Partial<_23.GetLatestBlockResponse>): _23.GetLatestBlockResponse;
                };
                GetSyncingRequest: {
                    encode(_: _23.GetSyncingRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetSyncingRequest;
                    fromJSON(_: any): _23.GetSyncingRequest;
                    toJSON(_: _23.GetSyncingRequest): unknown;
                    fromPartial(_: Partial<_23.GetSyncingRequest>): _23.GetSyncingRequest;
                };
                GetSyncingResponse: {
                    encode(message: _23.GetSyncingResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetSyncingResponse;
                    fromJSON(object: any): _23.GetSyncingResponse;
                    toJSON(message: _23.GetSyncingResponse): unknown;
                    fromPartial(object: Partial<_23.GetSyncingResponse>): _23.GetSyncingResponse;
                };
                GetNodeInfoRequest: {
                    encode(_: _23.GetNodeInfoRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetNodeInfoRequest;
                    fromJSON(_: any): _23.GetNodeInfoRequest;
                    toJSON(_: _23.GetNodeInfoRequest): unknown;
                    fromPartial(_: Partial<_23.GetNodeInfoRequest>): _23.GetNodeInfoRequest;
                };
                GetNodeInfoResponse: {
                    encode(message: _23.GetNodeInfoResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.GetNodeInfoResponse;
                    fromJSON(object: any): _23.GetNodeInfoResponse;
                    toJSON(message: _23.GetNodeInfoResponse): unknown;
                    fromPartial(object: Partial<_23.GetNodeInfoResponse>): _23.GetNodeInfoResponse;
                };
                VersionInfo: {
                    encode(message: _23.VersionInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.VersionInfo;
                    fromJSON(object: any): _23.VersionInfo;
                    toJSON(message: _23.VersionInfo): unknown;
                    fromPartial(object: Partial<_23.VersionInfo>): _23.VersionInfo;
                };
                Module: {
                    encode(message: _23.Module, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _23.Module;
                    fromJSON(object: any): _23.Module;
                    toJSON(message: _23.Module): unknown;
                    fromPartial(object: Partial<_23.Module>): _23.Module;
                };
            };
        }
        const v1beta1: {
            Coin: {
                encode(message: _24.Coin, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _24.Coin;
                fromJSON(object: any): _24.Coin;
                toJSON(message: _24.Coin): unknown;
                fromPartial(object: Partial<_24.Coin>): _24.Coin;
            };
            DecCoin: {
                encode(message: _24.DecCoin, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _24.DecCoin;
                fromJSON(object: any): _24.DecCoin;
                toJSON(message: _24.DecCoin): unknown;
                fromPartial(object: Partial<_24.DecCoin>): _24.DecCoin;
            };
            IntProto: {
                encode(message: _24.IntProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _24.IntProto;
                fromJSON(object: any): _24.IntProto;
                toJSON(message: _24.IntProto): unknown;
                fromPartial(object: Partial<_24.IntProto>): _24.IntProto;
            };
            DecProto: {
                encode(message: _24.DecProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _24.DecProto;
                fromJSON(object: any): _24.DecProto;
                toJSON(message: _24.DecProto): unknown;
                fromPartial(object: Partial<_24.DecProto>): _24.DecProto;
            };
        };
    }
    namespace capability {
        const v1beta1: {
            GenesisOwners: {
                encode(message: _26.GenesisOwners, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _26.GenesisOwners;
                fromJSON(object: any): _26.GenesisOwners;
                toJSON(message: _26.GenesisOwners): unknown;
                fromPartial(object: Partial<_26.GenesisOwners>): _26.GenesisOwners;
            };
            GenesisState: {
                encode(message: _26.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _26.GenesisState;
                fromJSON(object: any): _26.GenesisState;
                toJSON(message: _26.GenesisState): unknown;
                fromPartial(object: Partial<_26.GenesisState>): _26.GenesisState;
            };
            Capability: {
                encode(message: _25.Capability, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _25.Capability;
                fromJSON(object: any): _25.Capability;
                toJSON(message: _25.Capability): unknown;
                fromPartial(object: Partial<_25.Capability>): _25.Capability;
            };
            Owner: {
                encode(message: _25.Owner, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _25.Owner;
                fromJSON(object: any): _25.Owner;
                toJSON(message: _25.Owner): unknown;
                fromPartial(object: Partial<_25.Owner>): _25.Owner;
            };
            CapabilityOwners: {
                encode(message: _25.CapabilityOwners, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _25.CapabilityOwners;
                fromJSON(object: any): _25.CapabilityOwners;
                toJSON(message: _25.CapabilityOwners): unknown;
                fromPartial(object: Partial<_25.CapabilityOwners>): _25.CapabilityOwners;
            };
        };
    }
    namespace crisis {
        const v1beta1: {
            MsgClientImpl: typeof _164.MsgClientImpl;
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    verifyInvariant(value: _28.MsgVerifyInvariant): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    verifyInvariant(value: _28.MsgVerifyInvariant): {
                        typeUrl: string;
                        value: _28.MsgVerifyInvariant;
                    };
                };
                toJSON: {
                    verifyInvariant(value: _28.MsgVerifyInvariant): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    verifyInvariant(value: any): {
                        typeUrl: string;
                        value: _28.MsgVerifyInvariant;
                    };
                };
                fromPartial: {
                    verifyInvariant(value: _28.MsgVerifyInvariant): {
                        typeUrl: string;
                        value: _28.MsgVerifyInvariant;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.crisis.v1beta1.MsgVerifyInvariant": {
                    aminoType: string;
                    toAmino: ({ sender, invariantModuleName, invariantRoute }: _28.MsgVerifyInvariant) => {
                        sender: string;
                        invariant_module_name: string;
                        invariant_route: string;
                    };
                    fromAmino: ({ sender, invariant_module_name, invariant_route }: {
                        sender: string;
                        invariant_module_name: string;
                        invariant_route: string;
                    }) => _28.MsgVerifyInvariant;
                };
            };
            MsgVerifyInvariant: {
                encode(message: _28.MsgVerifyInvariant, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _28.MsgVerifyInvariant;
                fromJSON(object: any): _28.MsgVerifyInvariant;
                toJSON(message: _28.MsgVerifyInvariant): unknown;
                fromPartial(object: Partial<_28.MsgVerifyInvariant>): _28.MsgVerifyInvariant;
            };
            MsgVerifyInvariantResponse: {
                encode(_: _28.MsgVerifyInvariantResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _28.MsgVerifyInvariantResponse;
                fromJSON(_: any): _28.MsgVerifyInvariantResponse;
                toJSON(_: _28.MsgVerifyInvariantResponse): unknown;
                fromPartial(_: Partial<_28.MsgVerifyInvariantResponse>): _28.MsgVerifyInvariantResponse;
            };
            GenesisState: {
                encode(message: _27.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _27.GenesisState;
                fromJSON(object: any): _27.GenesisState;
                toJSON(message: _27.GenesisState): unknown;
                fromPartial(object: Partial<_27.GenesisState>): _27.GenesisState;
            };
        };
    }
    namespace crypto {
        const ed25519: {
            PubKey: {
                encode(message: _29.PubKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _29.PubKey;
                fromJSON(object: any): _29.PubKey;
                toJSON(message: _29.PubKey): unknown;
                fromPartial(object: Partial<_29.PubKey>): _29.PubKey;
            };
            PrivKey: {
                encode(message: _29.PrivKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _29.PrivKey;
                fromJSON(object: any): _29.PrivKey;
                toJSON(message: _29.PrivKey): unknown;
                fromPartial(object: Partial<_29.PrivKey>): _29.PrivKey;
            };
        };
        const multisig: {
            LegacyAminoPubKey: {
                encode(message: _30.LegacyAminoPubKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _30.LegacyAminoPubKey;
                fromJSON(object: any): _30.LegacyAminoPubKey;
                toJSON(message: _30.LegacyAminoPubKey): unknown;
                fromPartial(object: Partial<_30.LegacyAminoPubKey>): _30.LegacyAminoPubKey;
            };
        };
        const secp256k1: {
            PubKey: {
                encode(message: _31.PubKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _31.PubKey;
                fromJSON(object: any): _31.PubKey;
                toJSON(message: _31.PubKey): unknown;
                fromPartial(object: Partial<_31.PubKey>): _31.PubKey;
            };
            PrivKey: {
                encode(message: _31.PrivKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _31.PrivKey;
                fromJSON(object: any): _31.PrivKey;
                toJSON(message: _31.PrivKey): unknown;
                fromPartial(object: Partial<_31.PrivKey>): _31.PrivKey;
            };
        };
        const secp256r1: {
            PubKey: {
                encode(message: _32.PubKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _32.PubKey;
                fromJSON(object: any): _32.PubKey;
                toJSON(message: _32.PubKey): unknown;
                fromPartial(object: Partial<_32.PubKey>): _32.PubKey;
            };
            PrivKey: {
                encode(message: _32.PrivKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _32.PrivKey;
                fromJSON(object: any): _32.PrivKey;
                toJSON(message: _32.PrivKey): unknown;
                fromPartial(object: Partial<_32.PrivKey>): _32.PrivKey;
            };
        };
    }
    namespace distribution {
        const v1beta1: {
            MsgClientImpl: typeof _165.MsgClientImpl;
            QueryClientImpl: typeof _152.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _35.QueryParamsRequest): Promise<_35.QueryParamsResponse>;
                validatorOutstandingRewards(request: _35.QueryValidatorOutstandingRewardsRequest): Promise<_35.QueryValidatorOutstandingRewardsResponse>;
                validatorCommission(request: _35.QueryValidatorCommissionRequest): Promise<_35.QueryValidatorCommissionResponse>;
                validatorSlashes(request: _35.QueryValidatorSlashesRequest): Promise<_35.QueryValidatorSlashesResponse>;
                delegationRewards(request: _35.QueryDelegationRewardsRequest): Promise<_35.QueryDelegationRewardsResponse>;
                delegationTotalRewards(request: _35.QueryDelegationTotalRewardsRequest): Promise<_35.QueryDelegationTotalRewardsResponse>;
                delegatorValidators(request: _35.QueryDelegatorValidatorsRequest): Promise<_35.QueryDelegatorValidatorsResponse>;
                delegatorWithdrawAddress(request: _35.QueryDelegatorWithdrawAddressRequest): Promise<_35.QueryDelegatorWithdrawAddressResponse>;
                communityPool(request?: _35.QueryCommunityPoolRequest): Promise<_35.QueryCommunityPoolResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    setWithdrawAddress(value: _36.MsgSetWithdrawAddress): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    withdrawDelegatorReward(value: _36.MsgWithdrawDelegatorReward): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    withdrawValidatorCommission(value: _36.MsgWithdrawValidatorCommission): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    fundCommunityPool(value: _36.MsgFundCommunityPool): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    setWithdrawAddress(value: _36.MsgSetWithdrawAddress): {
                        typeUrl: string;
                        value: _36.MsgSetWithdrawAddress;
                    };
                    withdrawDelegatorReward(value: _36.MsgWithdrawDelegatorReward): {
                        typeUrl: string;
                        value: _36.MsgWithdrawDelegatorReward;
                    };
                    withdrawValidatorCommission(value: _36.MsgWithdrawValidatorCommission): {
                        typeUrl: string;
                        value: _36.MsgWithdrawValidatorCommission;
                    };
                    fundCommunityPool(value: _36.MsgFundCommunityPool): {
                        typeUrl: string;
                        value: _36.MsgFundCommunityPool;
                    };
                };
                toJSON: {
                    setWithdrawAddress(value: _36.MsgSetWithdrawAddress): {
                        typeUrl: string;
                        value: unknown;
                    };
                    withdrawDelegatorReward(value: _36.MsgWithdrawDelegatorReward): {
                        typeUrl: string;
                        value: unknown;
                    };
                    withdrawValidatorCommission(value: _36.MsgWithdrawValidatorCommission): {
                        typeUrl: string;
                        value: unknown;
                    };
                    fundCommunityPool(value: _36.MsgFundCommunityPool): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    setWithdrawAddress(value: any): {
                        typeUrl: string;
                        value: _36.MsgSetWithdrawAddress;
                    };
                    withdrawDelegatorReward(value: any): {
                        typeUrl: string;
                        value: _36.MsgWithdrawDelegatorReward;
                    };
                    withdrawValidatorCommission(value: any): {
                        typeUrl: string;
                        value: _36.MsgWithdrawValidatorCommission;
                    };
                    fundCommunityPool(value: any): {
                        typeUrl: string;
                        value: _36.MsgFundCommunityPool;
                    };
                };
                fromPartial: {
                    setWithdrawAddress(value: _36.MsgSetWithdrawAddress): {
                        typeUrl: string;
                        value: _36.MsgSetWithdrawAddress;
                    };
                    withdrawDelegatorReward(value: _36.MsgWithdrawDelegatorReward): {
                        typeUrl: string;
                        value: _36.MsgWithdrawDelegatorReward;
                    };
                    withdrawValidatorCommission(value: _36.MsgWithdrawValidatorCommission): {
                        typeUrl: string;
                        value: _36.MsgWithdrawValidatorCommission;
                    };
                    fundCommunityPool(value: _36.MsgFundCommunityPool): {
                        typeUrl: string;
                        value: _36.MsgFundCommunityPool;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress": {
                    aminoType: string;
                    toAmino: ({ delegatorAddress, withdrawAddress }: _36.MsgSetWithdrawAddress) => {
                        delegator_address: string;
                        withdraw_address: string;
                    };
                    fromAmino: ({ delegator_address, withdraw_address }: {
                        delegator_address: string;
                        withdraw_address: string;
                    }) => _36.MsgSetWithdrawAddress;
                };
                "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward": {
                    aminoType: string;
                    toAmino: ({ delegatorAddress, validatorAddress }: _36.MsgWithdrawDelegatorReward) => {
                        delegator_address: string;
                        validator_address: string;
                    };
                    fromAmino: ({ delegator_address, validator_address }: {
                        delegator_address: string;
                        validator_address: string;
                    }) => _36.MsgWithdrawDelegatorReward;
                };
                "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission": {
                    aminoType: string;
                    toAmino: ({ validatorAddress }: _36.MsgWithdrawValidatorCommission) => {
                        validator_address: string;
                    };
                    fromAmino: ({ validator_address }: {
                        validator_address: string;
                    }) => _36.MsgWithdrawValidatorCommission;
                };
                "/cosmos.distribution.v1beta1.MsgFundCommunityPool": {
                    aminoType: string;
                    toAmino: ({ amount, depositor }: _36.MsgFundCommunityPool) => {
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
                    }) => _36.MsgFundCommunityPool;
                };
            };
            MsgSetWithdrawAddress: {
                encode(message: _36.MsgSetWithdrawAddress, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgSetWithdrawAddress;
                fromJSON(object: any): _36.MsgSetWithdrawAddress;
                toJSON(message: _36.MsgSetWithdrawAddress): unknown;
                fromPartial(object: Partial<_36.MsgSetWithdrawAddress>): _36.MsgSetWithdrawAddress;
            };
            MsgSetWithdrawAddressResponse: {
                encode(_: _36.MsgSetWithdrawAddressResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgSetWithdrawAddressResponse;
                fromJSON(_: any): _36.MsgSetWithdrawAddressResponse;
                toJSON(_: _36.MsgSetWithdrawAddressResponse): unknown;
                fromPartial(_: Partial<_36.MsgSetWithdrawAddressResponse>): _36.MsgSetWithdrawAddressResponse;
            };
            MsgWithdrawDelegatorReward: {
                encode(message: _36.MsgWithdrawDelegatorReward, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgWithdrawDelegatorReward;
                fromJSON(object: any): _36.MsgWithdrawDelegatorReward;
                toJSON(message: _36.MsgWithdrawDelegatorReward): unknown;
                fromPartial(object: Partial<_36.MsgWithdrawDelegatorReward>): _36.MsgWithdrawDelegatorReward;
            };
            MsgWithdrawDelegatorRewardResponse: {
                encode(_: _36.MsgWithdrawDelegatorRewardResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgWithdrawDelegatorRewardResponse;
                fromJSON(_: any): _36.MsgWithdrawDelegatorRewardResponse;
                toJSON(_: _36.MsgWithdrawDelegatorRewardResponse): unknown;
                fromPartial(_: Partial<_36.MsgWithdrawDelegatorRewardResponse>): _36.MsgWithdrawDelegatorRewardResponse;
            };
            MsgWithdrawValidatorCommission: {
                encode(message: _36.MsgWithdrawValidatorCommission, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgWithdrawValidatorCommission;
                fromJSON(object: any): _36.MsgWithdrawValidatorCommission;
                toJSON(message: _36.MsgWithdrawValidatorCommission): unknown;
                fromPartial(object: Partial<_36.MsgWithdrawValidatorCommission>): _36.MsgWithdrawValidatorCommission;
            };
            MsgWithdrawValidatorCommissionResponse: {
                encode(_: _36.MsgWithdrawValidatorCommissionResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgWithdrawValidatorCommissionResponse;
                fromJSON(_: any): _36.MsgWithdrawValidatorCommissionResponse;
                toJSON(_: _36.MsgWithdrawValidatorCommissionResponse): unknown;
                fromPartial(_: Partial<_36.MsgWithdrawValidatorCommissionResponse>): _36.MsgWithdrawValidatorCommissionResponse;
            };
            MsgFundCommunityPool: {
                encode(message: _36.MsgFundCommunityPool, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgFundCommunityPool;
                fromJSON(object: any): _36.MsgFundCommunityPool;
                toJSON(message: _36.MsgFundCommunityPool): unknown;
                fromPartial(object: Partial<_36.MsgFundCommunityPool>): _36.MsgFundCommunityPool;
            };
            MsgFundCommunityPoolResponse: {
                encode(_: _36.MsgFundCommunityPoolResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _36.MsgFundCommunityPoolResponse;
                fromJSON(_: any): _36.MsgFundCommunityPoolResponse;
                toJSON(_: _36.MsgFundCommunityPoolResponse): unknown;
                fromPartial(_: Partial<_36.MsgFundCommunityPoolResponse>): _36.MsgFundCommunityPoolResponse;
            };
            QueryParamsRequest: {
                encode(_: _35.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryParamsRequest;
                fromJSON(_: any): _35.QueryParamsRequest;
                toJSON(_: _35.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_35.QueryParamsRequest>): _35.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _35.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryParamsResponse;
                fromJSON(object: any): _35.QueryParamsResponse;
                toJSON(message: _35.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_35.QueryParamsResponse>): _35.QueryParamsResponse;
            };
            QueryValidatorOutstandingRewardsRequest: {
                encode(message: _35.QueryValidatorOutstandingRewardsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryValidatorOutstandingRewardsRequest;
                fromJSON(object: any): _35.QueryValidatorOutstandingRewardsRequest;
                toJSON(message: _35.QueryValidatorOutstandingRewardsRequest): unknown;
                fromPartial(object: Partial<_35.QueryValidatorOutstandingRewardsRequest>): _35.QueryValidatorOutstandingRewardsRequest;
            };
            QueryValidatorOutstandingRewardsResponse: {
                encode(message: _35.QueryValidatorOutstandingRewardsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryValidatorOutstandingRewardsResponse;
                fromJSON(object: any): _35.QueryValidatorOutstandingRewardsResponse;
                toJSON(message: _35.QueryValidatorOutstandingRewardsResponse): unknown;
                fromPartial(object: Partial<_35.QueryValidatorOutstandingRewardsResponse>): _35.QueryValidatorOutstandingRewardsResponse;
            };
            QueryValidatorCommissionRequest: {
                encode(message: _35.QueryValidatorCommissionRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryValidatorCommissionRequest;
                fromJSON(object: any): _35.QueryValidatorCommissionRequest;
                toJSON(message: _35.QueryValidatorCommissionRequest): unknown;
                fromPartial(object: Partial<_35.QueryValidatorCommissionRequest>): _35.QueryValidatorCommissionRequest;
            };
            QueryValidatorCommissionResponse: {
                encode(message: _35.QueryValidatorCommissionResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryValidatorCommissionResponse;
                fromJSON(object: any): _35.QueryValidatorCommissionResponse;
                toJSON(message: _35.QueryValidatorCommissionResponse): unknown;
                fromPartial(object: Partial<_35.QueryValidatorCommissionResponse>): _35.QueryValidatorCommissionResponse;
            };
            QueryValidatorSlashesRequest: {
                encode(message: _35.QueryValidatorSlashesRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryValidatorSlashesRequest;
                fromJSON(object: any): _35.QueryValidatorSlashesRequest;
                toJSON(message: _35.QueryValidatorSlashesRequest): unknown;
                fromPartial(object: Partial<_35.QueryValidatorSlashesRequest>): _35.QueryValidatorSlashesRequest;
            };
            QueryValidatorSlashesResponse: {
                encode(message: _35.QueryValidatorSlashesResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryValidatorSlashesResponse;
                fromJSON(object: any): _35.QueryValidatorSlashesResponse;
                toJSON(message: _35.QueryValidatorSlashesResponse): unknown;
                fromPartial(object: Partial<_35.QueryValidatorSlashesResponse>): _35.QueryValidatorSlashesResponse;
            };
            QueryDelegationRewardsRequest: {
                encode(message: _35.QueryDelegationRewardsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegationRewardsRequest;
                fromJSON(object: any): _35.QueryDelegationRewardsRequest;
                toJSON(message: _35.QueryDelegationRewardsRequest): unknown;
                fromPartial(object: Partial<_35.QueryDelegationRewardsRequest>): _35.QueryDelegationRewardsRequest;
            };
            QueryDelegationRewardsResponse: {
                encode(message: _35.QueryDelegationRewardsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegationRewardsResponse;
                fromJSON(object: any): _35.QueryDelegationRewardsResponse;
                toJSON(message: _35.QueryDelegationRewardsResponse): unknown;
                fromPartial(object: Partial<_35.QueryDelegationRewardsResponse>): _35.QueryDelegationRewardsResponse;
            };
            QueryDelegationTotalRewardsRequest: {
                encode(message: _35.QueryDelegationTotalRewardsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegationTotalRewardsRequest;
                fromJSON(object: any): _35.QueryDelegationTotalRewardsRequest;
                toJSON(message: _35.QueryDelegationTotalRewardsRequest): unknown;
                fromPartial(object: Partial<_35.QueryDelegationTotalRewardsRequest>): _35.QueryDelegationTotalRewardsRequest;
            };
            QueryDelegationTotalRewardsResponse: {
                encode(message: _35.QueryDelegationTotalRewardsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegationTotalRewardsResponse;
                fromJSON(object: any): _35.QueryDelegationTotalRewardsResponse;
                toJSON(message: _35.QueryDelegationTotalRewardsResponse): unknown;
                fromPartial(object: Partial<_35.QueryDelegationTotalRewardsResponse>): _35.QueryDelegationTotalRewardsResponse;
            };
            QueryDelegatorValidatorsRequest: {
                encode(message: _35.QueryDelegatorValidatorsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegatorValidatorsRequest;
                fromJSON(object: any): _35.QueryDelegatorValidatorsRequest;
                toJSON(message: _35.QueryDelegatorValidatorsRequest): unknown;
                fromPartial(object: Partial<_35.QueryDelegatorValidatorsRequest>): _35.QueryDelegatorValidatorsRequest;
            };
            QueryDelegatorValidatorsResponse: {
                encode(message: _35.QueryDelegatorValidatorsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegatorValidatorsResponse;
                fromJSON(object: any): _35.QueryDelegatorValidatorsResponse;
                toJSON(message: _35.QueryDelegatorValidatorsResponse): unknown;
                fromPartial(object: Partial<_35.QueryDelegatorValidatorsResponse>): _35.QueryDelegatorValidatorsResponse;
            };
            QueryDelegatorWithdrawAddressRequest: {
                encode(message: _35.QueryDelegatorWithdrawAddressRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegatorWithdrawAddressRequest;
                fromJSON(object: any): _35.QueryDelegatorWithdrawAddressRequest;
                toJSON(message: _35.QueryDelegatorWithdrawAddressRequest): unknown;
                fromPartial(object: Partial<_35.QueryDelegatorWithdrawAddressRequest>): _35.QueryDelegatorWithdrawAddressRequest;
            };
            QueryDelegatorWithdrawAddressResponse: {
                encode(message: _35.QueryDelegatorWithdrawAddressResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryDelegatorWithdrawAddressResponse;
                fromJSON(object: any): _35.QueryDelegatorWithdrawAddressResponse;
                toJSON(message: _35.QueryDelegatorWithdrawAddressResponse): unknown;
                fromPartial(object: Partial<_35.QueryDelegatorWithdrawAddressResponse>): _35.QueryDelegatorWithdrawAddressResponse;
            };
            QueryCommunityPoolRequest: {
                encode(_: _35.QueryCommunityPoolRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryCommunityPoolRequest;
                fromJSON(_: any): _35.QueryCommunityPoolRequest;
                toJSON(_: _35.QueryCommunityPoolRequest): unknown;
                fromPartial(_: Partial<_35.QueryCommunityPoolRequest>): _35.QueryCommunityPoolRequest;
            };
            QueryCommunityPoolResponse: {
                encode(message: _35.QueryCommunityPoolResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _35.QueryCommunityPoolResponse;
                fromJSON(object: any): _35.QueryCommunityPoolResponse;
                toJSON(message: _35.QueryCommunityPoolResponse): unknown;
                fromPartial(object: Partial<_35.QueryCommunityPoolResponse>): _35.QueryCommunityPoolResponse;
            };
            DelegatorWithdrawInfo: {
                encode(message: _34.DelegatorWithdrawInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.DelegatorWithdrawInfo;
                fromJSON(object: any): _34.DelegatorWithdrawInfo;
                toJSON(message: _34.DelegatorWithdrawInfo): unknown;
                fromPartial(object: Partial<_34.DelegatorWithdrawInfo>): _34.DelegatorWithdrawInfo;
            };
            ValidatorOutstandingRewardsRecord: {
                encode(message: _34.ValidatorOutstandingRewardsRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.ValidatorOutstandingRewardsRecord;
                fromJSON(object: any): _34.ValidatorOutstandingRewardsRecord;
                toJSON(message: _34.ValidatorOutstandingRewardsRecord): unknown;
                fromPartial(object: Partial<_34.ValidatorOutstandingRewardsRecord>): _34.ValidatorOutstandingRewardsRecord;
            };
            ValidatorAccumulatedCommissionRecord: {
                encode(message: _34.ValidatorAccumulatedCommissionRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.ValidatorAccumulatedCommissionRecord;
                fromJSON(object: any): _34.ValidatorAccumulatedCommissionRecord;
                toJSON(message: _34.ValidatorAccumulatedCommissionRecord): unknown;
                fromPartial(object: Partial<_34.ValidatorAccumulatedCommissionRecord>): _34.ValidatorAccumulatedCommissionRecord;
            };
            ValidatorHistoricalRewardsRecord: {
                encode(message: _34.ValidatorHistoricalRewardsRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.ValidatorHistoricalRewardsRecord;
                fromJSON(object: any): _34.ValidatorHistoricalRewardsRecord;
                toJSON(message: _34.ValidatorHistoricalRewardsRecord): unknown;
                fromPartial(object: Partial<_34.ValidatorHistoricalRewardsRecord>): _34.ValidatorHistoricalRewardsRecord;
            };
            ValidatorCurrentRewardsRecord: {
                encode(message: _34.ValidatorCurrentRewardsRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.ValidatorCurrentRewardsRecord;
                fromJSON(object: any): _34.ValidatorCurrentRewardsRecord;
                toJSON(message: _34.ValidatorCurrentRewardsRecord): unknown;
                fromPartial(object: Partial<_34.ValidatorCurrentRewardsRecord>): _34.ValidatorCurrentRewardsRecord;
            };
            DelegatorStartingInfoRecord: {
                encode(message: _34.DelegatorStartingInfoRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.DelegatorStartingInfoRecord;
                fromJSON(object: any): _34.DelegatorStartingInfoRecord;
                toJSON(message: _34.DelegatorStartingInfoRecord): unknown;
                fromPartial(object: Partial<_34.DelegatorStartingInfoRecord>): _34.DelegatorStartingInfoRecord;
            };
            ValidatorSlashEventRecord: {
                encode(message: _34.ValidatorSlashEventRecord, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.ValidatorSlashEventRecord;
                fromJSON(object: any): _34.ValidatorSlashEventRecord;
                toJSON(message: _34.ValidatorSlashEventRecord): unknown;
                fromPartial(object: Partial<_34.ValidatorSlashEventRecord>): _34.ValidatorSlashEventRecord;
            };
            GenesisState: {
                encode(message: _34.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _34.GenesisState;
                fromJSON(object: any): _34.GenesisState;
                toJSON(message: _34.GenesisState): unknown;
                fromPartial(object: Partial<_34.GenesisState>): _34.GenesisState;
            };
            Params: {
                encode(message: _33.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.Params;
                fromJSON(object: any): _33.Params;
                toJSON(message: _33.Params): unknown;
                fromPartial(object: Partial<_33.Params>): _33.Params;
            };
            ValidatorHistoricalRewards: {
                encode(message: _33.ValidatorHistoricalRewards, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.ValidatorHistoricalRewards;
                fromJSON(object: any): _33.ValidatorHistoricalRewards;
                toJSON(message: _33.ValidatorHistoricalRewards): unknown;
                fromPartial(object: Partial<_33.ValidatorHistoricalRewards>): _33.ValidatorHistoricalRewards;
            };
            ValidatorCurrentRewards: {
                encode(message: _33.ValidatorCurrentRewards, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.ValidatorCurrentRewards;
                fromJSON(object: any): _33.ValidatorCurrentRewards;
                toJSON(message: _33.ValidatorCurrentRewards): unknown;
                fromPartial(object: Partial<_33.ValidatorCurrentRewards>): _33.ValidatorCurrentRewards;
            };
            ValidatorAccumulatedCommission: {
                encode(message: _33.ValidatorAccumulatedCommission, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.ValidatorAccumulatedCommission;
                fromJSON(object: any): _33.ValidatorAccumulatedCommission;
                toJSON(message: _33.ValidatorAccumulatedCommission): unknown;
                fromPartial(object: Partial<_33.ValidatorAccumulatedCommission>): _33.ValidatorAccumulatedCommission;
            };
            ValidatorOutstandingRewards: {
                encode(message: _33.ValidatorOutstandingRewards, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.ValidatorOutstandingRewards;
                fromJSON(object: any): _33.ValidatorOutstandingRewards;
                toJSON(message: _33.ValidatorOutstandingRewards): unknown;
                fromPartial(object: Partial<_33.ValidatorOutstandingRewards>): _33.ValidatorOutstandingRewards;
            };
            ValidatorSlashEvent: {
                encode(message: _33.ValidatorSlashEvent, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.ValidatorSlashEvent;
                fromJSON(object: any): _33.ValidatorSlashEvent;
                toJSON(message: _33.ValidatorSlashEvent): unknown;
                fromPartial(object: Partial<_33.ValidatorSlashEvent>): _33.ValidatorSlashEvent;
            };
            ValidatorSlashEvents: {
                encode(message: _33.ValidatorSlashEvents, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.ValidatorSlashEvents;
                fromJSON(object: any): _33.ValidatorSlashEvents;
                toJSON(message: _33.ValidatorSlashEvents): unknown;
                fromPartial(object: Partial<_33.ValidatorSlashEvents>): _33.ValidatorSlashEvents;
            };
            FeePool: {
                encode(message: _33.FeePool, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.FeePool;
                fromJSON(object: any): _33.FeePool;
                toJSON(message: _33.FeePool): unknown;
                fromPartial(object: Partial<_33.FeePool>): _33.FeePool;
            };
            CommunityPoolSpendProposal: {
                encode(message: _33.CommunityPoolSpendProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.CommunityPoolSpendProposal;
                fromJSON(object: any): _33.CommunityPoolSpendProposal;
                toJSON(message: _33.CommunityPoolSpendProposal): unknown;
                fromPartial(object: Partial<_33.CommunityPoolSpendProposal>): _33.CommunityPoolSpendProposal;
            };
            DelegatorStartingInfo: {
                encode(message: _33.DelegatorStartingInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.DelegatorStartingInfo;
                fromJSON(object: any): _33.DelegatorStartingInfo;
                toJSON(message: _33.DelegatorStartingInfo): unknown;
                fromPartial(object: Partial<_33.DelegatorStartingInfo>): _33.DelegatorStartingInfo;
            };
            DelegationDelegatorReward: {
                encode(message: _33.DelegationDelegatorReward, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.DelegationDelegatorReward;
                fromJSON(object: any): _33.DelegationDelegatorReward;
                toJSON(message: _33.DelegationDelegatorReward): unknown;
                fromPartial(object: Partial<_33.DelegationDelegatorReward>): _33.DelegationDelegatorReward;
            };
            CommunityPoolSpendProposalWithDeposit: {
                encode(message: _33.CommunityPoolSpendProposalWithDeposit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _33.CommunityPoolSpendProposalWithDeposit;
                fromJSON(object: any): _33.CommunityPoolSpendProposalWithDeposit;
                toJSON(message: _33.CommunityPoolSpendProposalWithDeposit): unknown;
                fromPartial(object: Partial<_33.CommunityPoolSpendProposalWithDeposit>): _33.CommunityPoolSpendProposalWithDeposit;
            };
        };
    }
    namespace evidence {
        const v1beta1: {
            MsgClientImpl: typeof _166.MsgClientImpl;
            QueryClientImpl: typeof _153.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                evidence(request: _39.QueryEvidenceRequest): Promise<_39.QueryEvidenceResponse>;
                allEvidence(request?: _39.QueryAllEvidenceRequest): Promise<_39.QueryAllEvidenceResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    submitEvidence(value: _40.MsgSubmitEvidence): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    submitEvidence(value: _40.MsgSubmitEvidence): {
                        typeUrl: string;
                        value: _40.MsgSubmitEvidence;
                    };
                };
                toJSON: {
                    submitEvidence(value: _40.MsgSubmitEvidence): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    submitEvidence(value: any): {
                        typeUrl: string;
                        value: _40.MsgSubmitEvidence;
                    };
                };
                fromPartial: {
                    submitEvidence(value: _40.MsgSubmitEvidence): {
                        typeUrl: string;
                        value: _40.MsgSubmitEvidence;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.evidence.v1beta1.MsgSubmitEvidence": {
                    aminoType: string;
                    toAmino: ({ submitter, evidence }: _40.MsgSubmitEvidence) => {
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
                    }) => _40.MsgSubmitEvidence;
                };
            };
            MsgSubmitEvidence: {
                encode(message: _40.MsgSubmitEvidence, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _40.MsgSubmitEvidence;
                fromJSON(object: any): _40.MsgSubmitEvidence;
                toJSON(message: _40.MsgSubmitEvidence): unknown;
                fromPartial(object: Partial<_40.MsgSubmitEvidence>): _40.MsgSubmitEvidence;
            };
            MsgSubmitEvidenceResponse: {
                encode(message: _40.MsgSubmitEvidenceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _40.MsgSubmitEvidenceResponse;
                fromJSON(object: any): _40.MsgSubmitEvidenceResponse;
                toJSON(message: _40.MsgSubmitEvidenceResponse): unknown;
                fromPartial(object: Partial<_40.MsgSubmitEvidenceResponse>): _40.MsgSubmitEvidenceResponse;
            };
            QueryEvidenceRequest: {
                encode(message: _39.QueryEvidenceRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _39.QueryEvidenceRequest;
                fromJSON(object: any): _39.QueryEvidenceRequest;
                toJSON(message: _39.QueryEvidenceRequest): unknown;
                fromPartial(object: Partial<_39.QueryEvidenceRequest>): _39.QueryEvidenceRequest;
            };
            QueryEvidenceResponse: {
                encode(message: _39.QueryEvidenceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _39.QueryEvidenceResponse;
                fromJSON(object: any): _39.QueryEvidenceResponse;
                toJSON(message: _39.QueryEvidenceResponse): unknown;
                fromPartial(object: Partial<_39.QueryEvidenceResponse>): _39.QueryEvidenceResponse;
            };
            QueryAllEvidenceRequest: {
                encode(message: _39.QueryAllEvidenceRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _39.QueryAllEvidenceRequest;
                fromJSON(object: any): _39.QueryAllEvidenceRequest;
                toJSON(message: _39.QueryAllEvidenceRequest): unknown;
                fromPartial(object: Partial<_39.QueryAllEvidenceRequest>): _39.QueryAllEvidenceRequest;
            };
            QueryAllEvidenceResponse: {
                encode(message: _39.QueryAllEvidenceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _39.QueryAllEvidenceResponse;
                fromJSON(object: any): _39.QueryAllEvidenceResponse;
                toJSON(message: _39.QueryAllEvidenceResponse): unknown;
                fromPartial(object: Partial<_39.QueryAllEvidenceResponse>): _39.QueryAllEvidenceResponse;
            };
            GenesisState: {
                encode(message: _38.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _38.GenesisState;
                fromJSON(object: any): _38.GenesisState;
                toJSON(message: _38.GenesisState): unknown;
                fromPartial(object: Partial<_38.GenesisState>): _38.GenesisState;
            };
            Equivocation: {
                encode(message: _37.Equivocation, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _37.Equivocation;
                fromJSON(object: any): _37.Equivocation;
                toJSON(message: _37.Equivocation): unknown;
                fromPartial(object: Partial<_37.Equivocation>): _37.Equivocation;
            };
        };
    }
    namespace feegrant {
        const v1beta1: {
            MsgClientImpl: typeof _167.MsgClientImpl;
            QueryClientImpl: typeof _154.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                allowance(request: _43.QueryAllowanceRequest): Promise<_43.QueryAllowanceResponse>;
                allowances(request: _43.QueryAllowancesRequest): Promise<_43.QueryAllowancesResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    grantAllowance(value: _44.MsgGrantAllowance): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    revokeAllowance(value: _44.MsgRevokeAllowance): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    grantAllowance(value: _44.MsgGrantAllowance): {
                        typeUrl: string;
                        value: _44.MsgGrantAllowance;
                    };
                    revokeAllowance(value: _44.MsgRevokeAllowance): {
                        typeUrl: string;
                        value: _44.MsgRevokeAllowance;
                    };
                };
                toJSON: {
                    grantAllowance(value: _44.MsgGrantAllowance): {
                        typeUrl: string;
                        value: unknown;
                    };
                    revokeAllowance(value: _44.MsgRevokeAllowance): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    grantAllowance(value: any): {
                        typeUrl: string;
                        value: _44.MsgGrantAllowance;
                    };
                    revokeAllowance(value: any): {
                        typeUrl: string;
                        value: _44.MsgRevokeAllowance;
                    };
                };
                fromPartial: {
                    grantAllowance(value: _44.MsgGrantAllowance): {
                        typeUrl: string;
                        value: _44.MsgGrantAllowance;
                    };
                    revokeAllowance(value: _44.MsgRevokeAllowance): {
                        typeUrl: string;
                        value: _44.MsgRevokeAllowance;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.feegrant.v1beta1.MsgGrantAllowance": {
                    aminoType: string;
                    toAmino: ({ granter, grantee, allowance }: _44.MsgGrantAllowance) => {
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
                    }) => _44.MsgGrantAllowance;
                };
                "/cosmos.feegrant.v1beta1.MsgRevokeAllowance": {
                    aminoType: string;
                    toAmino: ({ granter, grantee }: _44.MsgRevokeAllowance) => {
                        granter: string;
                        grantee: string;
                    };
                    fromAmino: ({ granter, grantee }: {
                        granter: string;
                        grantee: string;
                    }) => _44.MsgRevokeAllowance;
                };
            };
            MsgGrantAllowance: {
                encode(message: _44.MsgGrantAllowance, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _44.MsgGrantAllowance;
                fromJSON(object: any): _44.MsgGrantAllowance;
                toJSON(message: _44.MsgGrantAllowance): unknown;
                fromPartial(object: Partial<_44.MsgGrantAllowance>): _44.MsgGrantAllowance;
            };
            MsgGrantAllowanceResponse: {
                encode(_: _44.MsgGrantAllowanceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _44.MsgGrantAllowanceResponse;
                fromJSON(_: any): _44.MsgGrantAllowanceResponse;
                toJSON(_: _44.MsgGrantAllowanceResponse): unknown;
                fromPartial(_: Partial<_44.MsgGrantAllowanceResponse>): _44.MsgGrantAllowanceResponse;
            };
            MsgRevokeAllowance: {
                encode(message: _44.MsgRevokeAllowance, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _44.MsgRevokeAllowance;
                fromJSON(object: any): _44.MsgRevokeAllowance;
                toJSON(message: _44.MsgRevokeAllowance): unknown;
                fromPartial(object: Partial<_44.MsgRevokeAllowance>): _44.MsgRevokeAllowance;
            };
            MsgRevokeAllowanceResponse: {
                encode(_: _44.MsgRevokeAllowanceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _44.MsgRevokeAllowanceResponse;
                fromJSON(_: any): _44.MsgRevokeAllowanceResponse;
                toJSON(_: _44.MsgRevokeAllowanceResponse): unknown;
                fromPartial(_: Partial<_44.MsgRevokeAllowanceResponse>): _44.MsgRevokeAllowanceResponse;
            };
            QueryAllowanceRequest: {
                encode(message: _43.QueryAllowanceRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _43.QueryAllowanceRequest;
                fromJSON(object: any): _43.QueryAllowanceRequest;
                toJSON(message: _43.QueryAllowanceRequest): unknown;
                fromPartial(object: Partial<_43.QueryAllowanceRequest>): _43.QueryAllowanceRequest;
            };
            QueryAllowanceResponse: {
                encode(message: _43.QueryAllowanceResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _43.QueryAllowanceResponse;
                fromJSON(object: any): _43.QueryAllowanceResponse;
                toJSON(message: _43.QueryAllowanceResponse): unknown;
                fromPartial(object: Partial<_43.QueryAllowanceResponse>): _43.QueryAllowanceResponse;
            };
            QueryAllowancesRequest: {
                encode(message: _43.QueryAllowancesRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _43.QueryAllowancesRequest;
                fromJSON(object: any): _43.QueryAllowancesRequest;
                toJSON(message: _43.QueryAllowancesRequest): unknown;
                fromPartial(object: Partial<_43.QueryAllowancesRequest>): _43.QueryAllowancesRequest;
            };
            QueryAllowancesResponse: {
                encode(message: _43.QueryAllowancesResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _43.QueryAllowancesResponse;
                fromJSON(object: any): _43.QueryAllowancesResponse;
                toJSON(message: _43.QueryAllowancesResponse): unknown;
                fromPartial(object: Partial<_43.QueryAllowancesResponse>): _43.QueryAllowancesResponse;
            };
            GenesisState: {
                encode(message: _42.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _42.GenesisState;
                fromJSON(object: any): _42.GenesisState;
                toJSON(message: _42.GenesisState): unknown;
                fromPartial(object: Partial<_42.GenesisState>): _42.GenesisState;
            };
            BasicAllowance: {
                encode(message: _41.BasicAllowance, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _41.BasicAllowance;
                fromJSON(object: any): _41.BasicAllowance;
                toJSON(message: _41.BasicAllowance): unknown;
                fromPartial(object: Partial<_41.BasicAllowance>): _41.BasicAllowance;
            };
            PeriodicAllowance: {
                encode(message: _41.PeriodicAllowance, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _41.PeriodicAllowance;
                fromJSON(object: any): _41.PeriodicAllowance;
                toJSON(message: _41.PeriodicAllowance): unknown;
                fromPartial(object: Partial<_41.PeriodicAllowance>): _41.PeriodicAllowance;
            };
            AllowedMsgAllowance: {
                encode(message: _41.AllowedMsgAllowance, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _41.AllowedMsgAllowance;
                fromJSON(object: any): _41.AllowedMsgAllowance;
                toJSON(message: _41.AllowedMsgAllowance): unknown;
                fromPartial(object: Partial<_41.AllowedMsgAllowance>): _41.AllowedMsgAllowance;
            };
            Grant: {
                encode(message: _41.Grant, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _41.Grant;
                fromJSON(object: any): _41.Grant;
                toJSON(message: _41.Grant): unknown;
                fromPartial(object: Partial<_41.Grant>): _41.Grant;
            };
        };
    }
    namespace genutil {
        const v1beta1: {
            GenesisState: {
                encode(message: _45.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _45.GenesisState;
                fromJSON(object: any): _45.GenesisState;
                toJSON(message: _45.GenesisState): unknown;
                fromPartial(object: Partial<_45.GenesisState>): _45.GenesisState;
            };
        };
    }
    namespace gov {
        const v1beta1: {
            MsgClientImpl: typeof _168.MsgClientImpl;
            QueryClientImpl: typeof _155.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                proposal(request: _48.QueryProposalRequest): Promise<_48.QueryProposalResponse>;
                proposals(request: _48.QueryProposalsRequest): Promise<_48.QueryProposalsResponse>;
                vote(request: _48.QueryVoteRequest): Promise<_48.QueryVoteResponse>;
                votes(request: _48.QueryVotesRequest): Promise<_48.QueryVotesResponse>;
                params(request: _48.QueryParamsRequest): Promise<_48.QueryParamsResponse>;
                deposit(request: _48.QueryDepositRequest): Promise<_48.QueryDepositResponse>;
                deposits(request: _48.QueryDepositsRequest): Promise<_48.QueryDepositsResponse>;
                tallyResult(request: _48.QueryTallyResultRequest): Promise<_48.QueryTallyResultResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    submitProposal(value: _49.MsgSubmitProposal): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    vote(value: _49.MsgVote): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    voteWeighted(value: _49.MsgVoteWeighted): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    deposit(value: _49.MsgDeposit): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    submitProposal(value: _49.MsgSubmitProposal): {
                        typeUrl: string;
                        value: _49.MsgSubmitProposal;
                    };
                    vote(value: _49.MsgVote): {
                        typeUrl: string;
                        value: _49.MsgVote;
                    };
                    voteWeighted(value: _49.MsgVoteWeighted): {
                        typeUrl: string;
                        value: _49.MsgVoteWeighted;
                    };
                    deposit(value: _49.MsgDeposit): {
                        typeUrl: string;
                        value: _49.MsgDeposit;
                    };
                };
                toJSON: {
                    submitProposal(value: _49.MsgSubmitProposal): {
                        typeUrl: string;
                        value: unknown;
                    };
                    vote(value: _49.MsgVote): {
                        typeUrl: string;
                        value: unknown;
                    };
                    voteWeighted(value: _49.MsgVoteWeighted): {
                        typeUrl: string;
                        value: unknown;
                    };
                    deposit(value: _49.MsgDeposit): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    submitProposal(value: any): {
                        typeUrl: string;
                        value: _49.MsgSubmitProposal;
                    };
                    vote(value: any): {
                        typeUrl: string;
                        value: _49.MsgVote;
                    };
                    voteWeighted(value: any): {
                        typeUrl: string;
                        value: _49.MsgVoteWeighted;
                    };
                    deposit(value: any): {
                        typeUrl: string;
                        value: _49.MsgDeposit;
                    };
                };
                fromPartial: {
                    submitProposal(value: _49.MsgSubmitProposal): {
                        typeUrl: string;
                        value: _49.MsgSubmitProposal;
                    };
                    vote(value: _49.MsgVote): {
                        typeUrl: string;
                        value: _49.MsgVote;
                    };
                    voteWeighted(value: _49.MsgVoteWeighted): {
                        typeUrl: string;
                        value: _49.MsgVoteWeighted;
                    };
                    deposit(value: _49.MsgDeposit): {
                        typeUrl: string;
                        value: _49.MsgDeposit;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.gov.v1beta1.MsgSubmitProposal": {
                    aminoType: string;
                    toAmino: ({ content, initialDeposit, proposer }: _49.MsgSubmitProposal) => {
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
                    }) => _49.MsgSubmitProposal;
                };
                "/cosmos.gov.v1beta1.MsgVote": {
                    aminoType: string;
                    toAmino: ({ proposalId, voter, option }: _49.MsgVote) => {
                        proposal_id: string;
                        voter: string;
                        option: number;
                    };
                    fromAmino: ({ proposal_id, voter, option }: {
                        proposal_id: string;
                        voter: string;
                        option: number;
                    }) => _49.MsgVote;
                };
                "/cosmos.gov.v1beta1.MsgVoteWeighted": {
                    aminoType: string;
                    toAmino: ({ proposalId, voter, options }: _49.MsgVoteWeighted) => {
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
                    }) => _49.MsgVoteWeighted;
                };
                "/cosmos.gov.v1beta1.MsgDeposit": {
                    aminoType: string;
                    toAmino: ({ proposalId, depositor, amount }: _49.MsgDeposit) => {
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
                    }) => _49.MsgDeposit;
                };
            };
            MsgSubmitProposal: {
                encode(message: _49.MsgSubmitProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgSubmitProposal;
                fromJSON(object: any): _49.MsgSubmitProposal;
                toJSON(message: _49.MsgSubmitProposal): unknown;
                fromPartial(object: Partial<_49.MsgSubmitProposal>): _49.MsgSubmitProposal;
            };
            MsgSubmitProposalResponse: {
                encode(message: _49.MsgSubmitProposalResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgSubmitProposalResponse;
                fromJSON(object: any): _49.MsgSubmitProposalResponse;
                toJSON(message: _49.MsgSubmitProposalResponse): unknown;
                fromPartial(object: Partial<_49.MsgSubmitProposalResponse>): _49.MsgSubmitProposalResponse;
            };
            MsgVote: {
                encode(message: _49.MsgVote, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgVote;
                fromJSON(object: any): _49.MsgVote;
                toJSON(message: _49.MsgVote): unknown;
                fromPartial(object: Partial<_49.MsgVote>): _49.MsgVote;
            };
            MsgVoteResponse: {
                encode(_: _49.MsgVoteResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgVoteResponse;
                fromJSON(_: any): _49.MsgVoteResponse;
                toJSON(_: _49.MsgVoteResponse): unknown;
                fromPartial(_: Partial<_49.MsgVoteResponse>): _49.MsgVoteResponse;
            };
            MsgVoteWeighted: {
                encode(message: _49.MsgVoteWeighted, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgVoteWeighted;
                fromJSON(object: any): _49.MsgVoteWeighted;
                toJSON(message: _49.MsgVoteWeighted): unknown;
                fromPartial(object: Partial<_49.MsgVoteWeighted>): _49.MsgVoteWeighted;
            };
            MsgVoteWeightedResponse: {
                encode(_: _49.MsgVoteWeightedResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgVoteWeightedResponse;
                fromJSON(_: any): _49.MsgVoteWeightedResponse;
                toJSON(_: _49.MsgVoteWeightedResponse): unknown;
                fromPartial(_: Partial<_49.MsgVoteWeightedResponse>): _49.MsgVoteWeightedResponse;
            };
            MsgDeposit: {
                encode(message: _49.MsgDeposit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgDeposit;
                fromJSON(object: any): _49.MsgDeposit;
                toJSON(message: _49.MsgDeposit): unknown;
                fromPartial(object: Partial<_49.MsgDeposit>): _49.MsgDeposit;
            };
            MsgDepositResponse: {
                encode(_: _49.MsgDepositResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _49.MsgDepositResponse;
                fromJSON(_: any): _49.MsgDepositResponse;
                toJSON(_: _49.MsgDepositResponse): unknown;
                fromPartial(_: Partial<_49.MsgDepositResponse>): _49.MsgDepositResponse;
            };
            QueryProposalRequest: {
                encode(message: _48.QueryProposalRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryProposalRequest;
                fromJSON(object: any): _48.QueryProposalRequest;
                toJSON(message: _48.QueryProposalRequest): unknown;
                fromPartial(object: Partial<_48.QueryProposalRequest>): _48.QueryProposalRequest;
            };
            QueryProposalResponse: {
                encode(message: _48.QueryProposalResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryProposalResponse;
                fromJSON(object: any): _48.QueryProposalResponse;
                toJSON(message: _48.QueryProposalResponse): unknown;
                fromPartial(object: Partial<_48.QueryProposalResponse>): _48.QueryProposalResponse;
            };
            QueryProposalsRequest: {
                encode(message: _48.QueryProposalsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryProposalsRequest;
                fromJSON(object: any): _48.QueryProposalsRequest;
                toJSON(message: _48.QueryProposalsRequest): unknown;
                fromPartial(object: Partial<_48.QueryProposalsRequest>): _48.QueryProposalsRequest;
            };
            QueryProposalsResponse: {
                encode(message: _48.QueryProposalsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryProposalsResponse;
                fromJSON(object: any): _48.QueryProposalsResponse;
                toJSON(message: _48.QueryProposalsResponse): unknown;
                fromPartial(object: Partial<_48.QueryProposalsResponse>): _48.QueryProposalsResponse;
            };
            QueryVoteRequest: {
                encode(message: _48.QueryVoteRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryVoteRequest;
                fromJSON(object: any): _48.QueryVoteRequest;
                toJSON(message: _48.QueryVoteRequest): unknown;
                fromPartial(object: Partial<_48.QueryVoteRequest>): _48.QueryVoteRequest;
            };
            QueryVoteResponse: {
                encode(message: _48.QueryVoteResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryVoteResponse;
                fromJSON(object: any): _48.QueryVoteResponse;
                toJSON(message: _48.QueryVoteResponse): unknown;
                fromPartial(object: Partial<_48.QueryVoteResponse>): _48.QueryVoteResponse;
            };
            QueryVotesRequest: {
                encode(message: _48.QueryVotesRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryVotesRequest;
                fromJSON(object: any): _48.QueryVotesRequest;
                toJSON(message: _48.QueryVotesRequest): unknown;
                fromPartial(object: Partial<_48.QueryVotesRequest>): _48.QueryVotesRequest;
            };
            QueryVotesResponse: {
                encode(message: _48.QueryVotesResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryVotesResponse;
                fromJSON(object: any): _48.QueryVotesResponse;
                toJSON(message: _48.QueryVotesResponse): unknown;
                fromPartial(object: Partial<_48.QueryVotesResponse>): _48.QueryVotesResponse;
            };
            QueryParamsRequest: {
                encode(message: _48.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryParamsRequest;
                fromJSON(object: any): _48.QueryParamsRequest;
                toJSON(message: _48.QueryParamsRequest): unknown;
                fromPartial(object: Partial<_48.QueryParamsRequest>): _48.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _48.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryParamsResponse;
                fromJSON(object: any): _48.QueryParamsResponse;
                toJSON(message: _48.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_48.QueryParamsResponse>): _48.QueryParamsResponse;
            };
            QueryDepositRequest: {
                encode(message: _48.QueryDepositRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryDepositRequest;
                fromJSON(object: any): _48.QueryDepositRequest;
                toJSON(message: _48.QueryDepositRequest): unknown;
                fromPartial(object: Partial<_48.QueryDepositRequest>): _48.QueryDepositRequest;
            };
            QueryDepositResponse: {
                encode(message: _48.QueryDepositResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryDepositResponse;
                fromJSON(object: any): _48.QueryDepositResponse;
                toJSON(message: _48.QueryDepositResponse): unknown;
                fromPartial(object: Partial<_48.QueryDepositResponse>): _48.QueryDepositResponse;
            };
            QueryDepositsRequest: {
                encode(message: _48.QueryDepositsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryDepositsRequest;
                fromJSON(object: any): _48.QueryDepositsRequest;
                toJSON(message: _48.QueryDepositsRequest): unknown;
                fromPartial(object: Partial<_48.QueryDepositsRequest>): _48.QueryDepositsRequest;
            };
            QueryDepositsResponse: {
                encode(message: _48.QueryDepositsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryDepositsResponse;
                fromJSON(object: any): _48.QueryDepositsResponse;
                toJSON(message: _48.QueryDepositsResponse): unknown;
                fromPartial(object: Partial<_48.QueryDepositsResponse>): _48.QueryDepositsResponse;
            };
            QueryTallyResultRequest: {
                encode(message: _48.QueryTallyResultRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryTallyResultRequest;
                fromJSON(object: any): _48.QueryTallyResultRequest;
                toJSON(message: _48.QueryTallyResultRequest): unknown;
                fromPartial(object: Partial<_48.QueryTallyResultRequest>): _48.QueryTallyResultRequest;
            };
            QueryTallyResultResponse: {
                encode(message: _48.QueryTallyResultResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _48.QueryTallyResultResponse;
                fromJSON(object: any): _48.QueryTallyResultResponse;
                toJSON(message: _48.QueryTallyResultResponse): unknown;
                fromPartial(object: Partial<_48.QueryTallyResultResponse>): _48.QueryTallyResultResponse;
            };
            voteOptionFromJSON(object: any): _47.VoteOption;
            voteOptionToJSON(object: _47.VoteOption): string;
            proposalStatusFromJSON(object: any): _47.ProposalStatus;
            proposalStatusToJSON(object: _47.ProposalStatus): string;
            VoteOption: typeof _47.VoteOption;
            VoteOptionSDKType: typeof _47.VoteOptionSDKType;
            ProposalStatus: typeof _47.ProposalStatus;
            ProposalStatusSDKType: typeof _47.ProposalStatusSDKType;
            WeightedVoteOption: {
                encode(message: _47.WeightedVoteOption, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.WeightedVoteOption;
                fromJSON(object: any): _47.WeightedVoteOption;
                toJSON(message: _47.WeightedVoteOption): unknown;
                fromPartial(object: Partial<_47.WeightedVoteOption>): _47.WeightedVoteOption;
            };
            TextProposal: {
                encode(message: _47.TextProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.TextProposal;
                fromJSON(object: any): _47.TextProposal;
                toJSON(message: _47.TextProposal): unknown;
                fromPartial(object: Partial<_47.TextProposal>): _47.TextProposal;
            };
            Deposit: {
                encode(message: _47.Deposit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.Deposit;
                fromJSON(object: any): _47.Deposit;
                toJSON(message: _47.Deposit): unknown;
                fromPartial(object: Partial<_47.Deposit>): _47.Deposit;
            };
            Proposal: {
                encode(message: _47.Proposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.Proposal;
                fromJSON(object: any): _47.Proposal;
                toJSON(message: _47.Proposal): unknown;
                fromPartial(object: Partial<_47.Proposal>): _47.Proposal;
            };
            TallyResult: {
                encode(message: _47.TallyResult, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.TallyResult;
                fromJSON(object: any): _47.TallyResult;
                toJSON(message: _47.TallyResult): unknown;
                fromPartial(object: Partial<_47.TallyResult>): _47.TallyResult;
            };
            Vote: {
                encode(message: _47.Vote, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.Vote;
                fromJSON(object: any): _47.Vote;
                toJSON(message: _47.Vote): unknown;
                fromPartial(object: Partial<_47.Vote>): _47.Vote;
            };
            DepositParams: {
                encode(message: _47.DepositParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.DepositParams;
                fromJSON(object: any): _47.DepositParams;
                toJSON(message: _47.DepositParams): unknown;
                fromPartial(object: Partial<_47.DepositParams>): _47.DepositParams;
            };
            VotingParams: {
                encode(message: _47.VotingParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.VotingParams;
                fromJSON(object: any): _47.VotingParams;
                toJSON(message: _47.VotingParams): unknown;
                fromPartial(object: Partial<_47.VotingParams>): _47.VotingParams;
            };
            TallyParams: {
                encode(message: _47.TallyParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _47.TallyParams;
                fromJSON(object: any): _47.TallyParams;
                toJSON(message: _47.TallyParams): unknown;
                fromPartial(object: Partial<_47.TallyParams>): _47.TallyParams;
            };
            GenesisState: {
                encode(message: _46.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _46.GenesisState;
                fromJSON(object: any): _46.GenesisState;
                toJSON(message: _46.GenesisState): unknown;
                fromPartial(object: Partial<_46.GenesisState>): _46.GenesisState;
            };
        };
    }
    namespace mint {
        const v1beta1: {
            QueryClientImpl: typeof _156.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _52.QueryParamsRequest): Promise<_52.QueryParamsResponse>;
                inflation(request?: _52.QueryInflationRequest): Promise<_52.QueryInflationResponse>;
                annualProvisions(request?: _52.QueryAnnualProvisionsRequest): Promise<_52.QueryAnnualProvisionsResponse>;
            };
            QueryParamsRequest: {
                encode(_: _52.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _52.QueryParamsRequest;
                fromJSON(_: any): _52.QueryParamsRequest;
                toJSON(_: _52.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_52.QueryParamsRequest>): _52.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _52.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _52.QueryParamsResponse;
                fromJSON(object: any): _52.QueryParamsResponse;
                toJSON(message: _52.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_52.QueryParamsResponse>): _52.QueryParamsResponse;
            };
            QueryInflationRequest: {
                encode(_: _52.QueryInflationRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _52.QueryInflationRequest;
                fromJSON(_: any): _52.QueryInflationRequest;
                toJSON(_: _52.QueryInflationRequest): unknown;
                fromPartial(_: Partial<_52.QueryInflationRequest>): _52.QueryInflationRequest;
            };
            QueryInflationResponse: {
                encode(message: _52.QueryInflationResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _52.QueryInflationResponse;
                fromJSON(object: any): _52.QueryInflationResponse;
                toJSON(message: _52.QueryInflationResponse): unknown;
                fromPartial(object: Partial<_52.QueryInflationResponse>): _52.QueryInflationResponse;
            };
            QueryAnnualProvisionsRequest: {
                encode(_: _52.QueryAnnualProvisionsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _52.QueryAnnualProvisionsRequest;
                fromJSON(_: any): _52.QueryAnnualProvisionsRequest;
                toJSON(_: _52.QueryAnnualProvisionsRequest): unknown;
                fromPartial(_: Partial<_52.QueryAnnualProvisionsRequest>): _52.QueryAnnualProvisionsRequest;
            };
            QueryAnnualProvisionsResponse: {
                encode(message: _52.QueryAnnualProvisionsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _52.QueryAnnualProvisionsResponse;
                fromJSON(object: any): _52.QueryAnnualProvisionsResponse;
                toJSON(message: _52.QueryAnnualProvisionsResponse): unknown;
                fromPartial(object: Partial<_52.QueryAnnualProvisionsResponse>): _52.QueryAnnualProvisionsResponse;
            };
            Minter: {
                encode(message: _51.Minter, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _51.Minter;
                fromJSON(object: any): _51.Minter;
                toJSON(message: _51.Minter): unknown;
                fromPartial(object: Partial<_51.Minter>): _51.Minter;
            };
            Params: {
                encode(message: _51.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _51.Params;
                fromJSON(object: any): _51.Params;
                toJSON(message: _51.Params): unknown;
                fromPartial(object: Partial<_51.Params>): _51.Params;
            };
            GenesisState: {
                encode(message: _50.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _50.GenesisState;
                fromJSON(object: any): _50.GenesisState;
                toJSON(message: _50.GenesisState): unknown;
                fromPartial(object: Partial<_50.GenesisState>): _50.GenesisState;
            };
        };
    }
    namespace params {
        const v1beta1: {
            QueryClientImpl: typeof _157.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request: _54.QueryParamsRequest): Promise<_54.QueryParamsResponse>;
            };
            QueryParamsRequest: {
                encode(message: _54.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _54.QueryParamsRequest;
                fromJSON(object: any): _54.QueryParamsRequest;
                toJSON(message: _54.QueryParamsRequest): unknown;
                fromPartial(object: Partial<_54.QueryParamsRequest>): _54.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _54.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _54.QueryParamsResponse;
                fromJSON(object: any): _54.QueryParamsResponse;
                toJSON(message: _54.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_54.QueryParamsResponse>): _54.QueryParamsResponse;
            };
            ParameterChangeProposal: {
                encode(message: _53.ParameterChangeProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _53.ParameterChangeProposal;
                fromJSON(object: any): _53.ParameterChangeProposal;
                toJSON(message: _53.ParameterChangeProposal): unknown;
                fromPartial(object: Partial<_53.ParameterChangeProposal>): _53.ParameterChangeProposal;
            };
            ParamChange: {
                encode(message: _53.ParamChange, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _53.ParamChange;
                fromJSON(object: any): _53.ParamChange;
                toJSON(message: _53.ParamChange): unknown;
                fromPartial(object: Partial<_53.ParamChange>): _53.ParamChange;
            };
        };
    }
    namespace slashing {
        const v1beta1: {
            MsgClientImpl: typeof _169.MsgClientImpl;
            QueryClientImpl: typeof _158.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                params(request?: _56.QueryParamsRequest): Promise<_56.QueryParamsResponse>;
                signingInfo(request: _56.QuerySigningInfoRequest): Promise<_56.QuerySigningInfoResponse>;
                signingInfos(request?: _56.QuerySigningInfosRequest): Promise<_56.QuerySigningInfosResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    unjail(value: _58.MsgUnjail): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    unjail(value: _58.MsgUnjail): {
                        typeUrl: string;
                        value: _58.MsgUnjail;
                    };
                };
                toJSON: {
                    unjail(value: _58.MsgUnjail): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    unjail(value: any): {
                        typeUrl: string;
                        value: _58.MsgUnjail;
                    };
                };
                fromPartial: {
                    unjail(value: _58.MsgUnjail): {
                        typeUrl: string;
                        value: _58.MsgUnjail;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.slashing.v1beta1.MsgUnjail": {
                    aminoType: string;
                    toAmino: ({ validatorAddr }: _58.MsgUnjail) => {
                        validator_addr: string;
                    };
                    fromAmino: ({ validator_addr }: {
                        validator_addr: string;
                    }) => _58.MsgUnjail;
                };
            };
            MsgUnjail: {
                encode(message: _58.MsgUnjail, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _58.MsgUnjail;
                fromJSON(object: any): _58.MsgUnjail;
                toJSON(message: _58.MsgUnjail): unknown;
                fromPartial(object: Partial<_58.MsgUnjail>): _58.MsgUnjail;
            };
            MsgUnjailResponse: {
                encode(_: _58.MsgUnjailResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _58.MsgUnjailResponse;
                fromJSON(_: any): _58.MsgUnjailResponse;
                toJSON(_: _58.MsgUnjailResponse): unknown;
                fromPartial(_: Partial<_58.MsgUnjailResponse>): _58.MsgUnjailResponse;
            };
            ValidatorSigningInfo: {
                encode(message: _57.ValidatorSigningInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _57.ValidatorSigningInfo;
                fromJSON(object: any): _57.ValidatorSigningInfo;
                toJSON(message: _57.ValidatorSigningInfo): unknown;
                fromPartial(object: Partial<_57.ValidatorSigningInfo>): _57.ValidatorSigningInfo;
            };
            Params: {
                encode(message: _57.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _57.Params;
                fromJSON(object: any): _57.Params;
                toJSON(message: _57.Params): unknown;
                fromPartial(object: Partial<_57.Params>): _57.Params;
            };
            QueryParamsRequest: {
                encode(_: _56.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _56.QueryParamsRequest;
                fromJSON(_: any): _56.QueryParamsRequest;
                toJSON(_: _56.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_56.QueryParamsRequest>): _56.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _56.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _56.QueryParamsResponse;
                fromJSON(object: any): _56.QueryParamsResponse;
                toJSON(message: _56.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_56.QueryParamsResponse>): _56.QueryParamsResponse;
            };
            QuerySigningInfoRequest: {
                encode(message: _56.QuerySigningInfoRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _56.QuerySigningInfoRequest;
                fromJSON(object: any): _56.QuerySigningInfoRequest;
                toJSON(message: _56.QuerySigningInfoRequest): unknown;
                fromPartial(object: Partial<_56.QuerySigningInfoRequest>): _56.QuerySigningInfoRequest;
            };
            QuerySigningInfoResponse: {
                encode(message: _56.QuerySigningInfoResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _56.QuerySigningInfoResponse;
                fromJSON(object: any): _56.QuerySigningInfoResponse;
                toJSON(message: _56.QuerySigningInfoResponse): unknown;
                fromPartial(object: Partial<_56.QuerySigningInfoResponse>): _56.QuerySigningInfoResponse;
            };
            QuerySigningInfosRequest: {
                encode(message: _56.QuerySigningInfosRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _56.QuerySigningInfosRequest;
                fromJSON(object: any): _56.QuerySigningInfosRequest;
                toJSON(message: _56.QuerySigningInfosRequest): unknown;
                fromPartial(object: Partial<_56.QuerySigningInfosRequest>): _56.QuerySigningInfosRequest;
            };
            QuerySigningInfosResponse: {
                encode(message: _56.QuerySigningInfosResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _56.QuerySigningInfosResponse;
                fromJSON(object: any): _56.QuerySigningInfosResponse;
                toJSON(message: _56.QuerySigningInfosResponse): unknown;
                fromPartial(object: Partial<_56.QuerySigningInfosResponse>): _56.QuerySigningInfosResponse;
            };
            GenesisState: {
                encode(message: _55.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _55.GenesisState;
                fromJSON(object: any): _55.GenesisState;
                toJSON(message: _55.GenesisState): unknown;
                fromPartial(object: Partial<_55.GenesisState>): _55.GenesisState;
            };
            SigningInfo: {
                encode(message: _55.SigningInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _55.SigningInfo;
                fromJSON(object: any): _55.SigningInfo;
                toJSON(message: _55.SigningInfo): unknown;
                fromPartial(object: Partial<_55.SigningInfo>): _55.SigningInfo;
            };
            ValidatorMissedBlocks: {
                encode(message: _55.ValidatorMissedBlocks, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _55.ValidatorMissedBlocks;
                fromJSON(object: any): _55.ValidatorMissedBlocks;
                toJSON(message: _55.ValidatorMissedBlocks): unknown;
                fromPartial(object: Partial<_55.ValidatorMissedBlocks>): _55.ValidatorMissedBlocks;
            };
            MissedBlock: {
                encode(message: _55.MissedBlock, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _55.MissedBlock;
                fromJSON(object: any): _55.MissedBlock;
                toJSON(message: _55.MissedBlock): unknown;
                fromPartial(object: Partial<_55.MissedBlock>): _55.MissedBlock;
            };
        };
    }
    namespace staking {
        const v1beta1: {
            MsgClientImpl: typeof _170.MsgClientImpl;
            QueryClientImpl: typeof _159.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                validators(request: _61.QueryValidatorsRequest): Promise<_61.QueryValidatorsResponse>;
                validator(request: _61.QueryValidatorRequest): Promise<_61.QueryValidatorResponse>;
                validatorDelegations(request: _61.QueryValidatorDelegationsRequest): Promise<_61.QueryValidatorDelegationsResponse>;
                validatorUnbondingDelegations(request: _61.QueryValidatorUnbondingDelegationsRequest): Promise<_61.QueryValidatorUnbondingDelegationsResponse>;
                delegation(request: _61.QueryDelegationRequest): Promise<_61.QueryDelegationResponse>;
                unbondingDelegation(request: _61.QueryUnbondingDelegationRequest): Promise<_61.QueryUnbondingDelegationResponse>;
                delegatorDelegations(request: _61.QueryDelegatorDelegationsRequest): Promise<_61.QueryDelegatorDelegationsResponse>;
                delegatorUnbondingDelegations(request: _61.QueryDelegatorUnbondingDelegationsRequest): Promise<_61.QueryDelegatorUnbondingDelegationsResponse>;
                redelegations(request: _61.QueryRedelegationsRequest): Promise<_61.QueryRedelegationsResponse>;
                delegatorValidators(request: _61.QueryDelegatorValidatorsRequest): Promise<_61.QueryDelegatorValidatorsResponse>;
                delegatorValidator(request: _61.QueryDelegatorValidatorRequest): Promise<_61.QueryDelegatorValidatorResponse>;
                historicalInfo(request: _61.QueryHistoricalInfoRequest): Promise<_61.QueryHistoricalInfoResponse>;
                pool(request?: _61.QueryPoolRequest): Promise<_61.QueryPoolResponse>;
                params(request?: _61.QueryParamsRequest): Promise<_61.QueryParamsResponse>;
            };
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    createValidator(value: _63.MsgCreateValidator): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    editValidator(value: _63.MsgEditValidator): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    delegate(value: _63.MsgDelegate): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    beginRedelegate(value: _63.MsgBeginRedelegate): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                    undelegate(value: _63.MsgUndelegate): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    createValidator(value: _63.MsgCreateValidator): {
                        typeUrl: string;
                        value: _63.MsgCreateValidator;
                    };
                    editValidator(value: _63.MsgEditValidator): {
                        typeUrl: string;
                        value: _63.MsgEditValidator;
                    };
                    delegate(value: _63.MsgDelegate): {
                        typeUrl: string;
                        value: _63.MsgDelegate;
                    };
                    beginRedelegate(value: _63.MsgBeginRedelegate): {
                        typeUrl: string;
                        value: _63.MsgBeginRedelegate;
                    };
                    undelegate(value: _63.MsgUndelegate): {
                        typeUrl: string;
                        value: _63.MsgUndelegate;
                    };
                };
                toJSON: {
                    createValidator(value: _63.MsgCreateValidator): {
                        typeUrl: string;
                        value: unknown;
                    };
                    editValidator(value: _63.MsgEditValidator): {
                        typeUrl: string;
                        value: unknown;
                    };
                    delegate(value: _63.MsgDelegate): {
                        typeUrl: string;
                        value: unknown;
                    };
                    beginRedelegate(value: _63.MsgBeginRedelegate): {
                        typeUrl: string;
                        value: unknown;
                    };
                    undelegate(value: _63.MsgUndelegate): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    createValidator(value: any): {
                        typeUrl: string;
                        value: _63.MsgCreateValidator;
                    };
                    editValidator(value: any): {
                        typeUrl: string;
                        value: _63.MsgEditValidator;
                    };
                    delegate(value: any): {
                        typeUrl: string;
                        value: _63.MsgDelegate;
                    };
                    beginRedelegate(value: any): {
                        typeUrl: string;
                        value: _63.MsgBeginRedelegate;
                    };
                    undelegate(value: any): {
                        typeUrl: string;
                        value: _63.MsgUndelegate;
                    };
                };
                fromPartial: {
                    createValidator(value: _63.MsgCreateValidator): {
                        typeUrl: string;
                        value: _63.MsgCreateValidator;
                    };
                    editValidator(value: _63.MsgEditValidator): {
                        typeUrl: string;
                        value: _63.MsgEditValidator;
                    };
                    delegate(value: _63.MsgDelegate): {
                        typeUrl: string;
                        value: _63.MsgDelegate;
                    };
                    beginRedelegate(value: _63.MsgBeginRedelegate): {
                        typeUrl: string;
                        value: _63.MsgBeginRedelegate;
                    };
                    undelegate(value: _63.MsgUndelegate): {
                        typeUrl: string;
                        value: _63.MsgUndelegate;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.staking.v1beta1.MsgCreateValidator": {
                    aminoType: string;
                    toAmino: ({ description, commission, minSelfDelegation, delegatorAddress, validatorAddress, pubkey, value }: _63.MsgCreateValidator) => {
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
                    }) => _63.MsgCreateValidator;
                };
                "/cosmos.staking.v1beta1.MsgEditValidator": {
                    aminoType: string;
                    toAmino: ({ description, validatorAddress, commissionRate, minSelfDelegation }: _63.MsgEditValidator) => {
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
                    }) => _63.MsgEditValidator;
                };
                "/cosmos.staking.v1beta1.MsgDelegate": {
                    aminoType: string;
                    toAmino: ({ delegatorAddress, validatorAddress, amount }: _63.MsgDelegate) => {
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
                    }) => _63.MsgDelegate;
                };
                "/cosmos.staking.v1beta1.MsgBeginRedelegate": {
                    aminoType: string;
                    toAmino: ({ delegatorAddress, validatorSrcAddress, validatorDstAddress, amount }: _63.MsgBeginRedelegate) => {
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
                    }) => _63.MsgBeginRedelegate;
                };
                "/cosmos.staking.v1beta1.MsgUndelegate": {
                    aminoType: string;
                    toAmino: ({ delegatorAddress, validatorAddress, amount }: _63.MsgUndelegate) => {
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
                    }) => _63.MsgUndelegate;
                };
            };
            MsgCreateValidator: {
                encode(message: _63.MsgCreateValidator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgCreateValidator;
                fromJSON(object: any): _63.MsgCreateValidator;
                toJSON(message: _63.MsgCreateValidator): unknown;
                fromPartial(object: Partial<_63.MsgCreateValidator>): _63.MsgCreateValidator;
            };
            MsgCreateValidatorResponse: {
                encode(_: _63.MsgCreateValidatorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgCreateValidatorResponse;
                fromJSON(_: any): _63.MsgCreateValidatorResponse;
                toJSON(_: _63.MsgCreateValidatorResponse): unknown;
                fromPartial(_: Partial<_63.MsgCreateValidatorResponse>): _63.MsgCreateValidatorResponse;
            };
            MsgEditValidator: {
                encode(message: _63.MsgEditValidator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgEditValidator;
                fromJSON(object: any): _63.MsgEditValidator;
                toJSON(message: _63.MsgEditValidator): unknown;
                fromPartial(object: Partial<_63.MsgEditValidator>): _63.MsgEditValidator;
            };
            MsgEditValidatorResponse: {
                encode(_: _63.MsgEditValidatorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgEditValidatorResponse;
                fromJSON(_: any): _63.MsgEditValidatorResponse;
                toJSON(_: _63.MsgEditValidatorResponse): unknown;
                fromPartial(_: Partial<_63.MsgEditValidatorResponse>): _63.MsgEditValidatorResponse;
            };
            MsgDelegate: {
                encode(message: _63.MsgDelegate, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgDelegate;
                fromJSON(object: any): _63.MsgDelegate;
                toJSON(message: _63.MsgDelegate): unknown;
                fromPartial(object: Partial<_63.MsgDelegate>): _63.MsgDelegate;
            };
            MsgDelegateResponse: {
                encode(_: _63.MsgDelegateResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgDelegateResponse;
                fromJSON(_: any): _63.MsgDelegateResponse;
                toJSON(_: _63.MsgDelegateResponse): unknown;
                fromPartial(_: Partial<_63.MsgDelegateResponse>): _63.MsgDelegateResponse;
            };
            MsgBeginRedelegate: {
                encode(message: _63.MsgBeginRedelegate, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgBeginRedelegate;
                fromJSON(object: any): _63.MsgBeginRedelegate;
                toJSON(message: _63.MsgBeginRedelegate): unknown;
                fromPartial(object: Partial<_63.MsgBeginRedelegate>): _63.MsgBeginRedelegate;
            };
            MsgBeginRedelegateResponse: {
                encode(message: _63.MsgBeginRedelegateResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgBeginRedelegateResponse;
                fromJSON(object: any): _63.MsgBeginRedelegateResponse;
                toJSON(message: _63.MsgBeginRedelegateResponse): unknown;
                fromPartial(object: Partial<_63.MsgBeginRedelegateResponse>): _63.MsgBeginRedelegateResponse;
            };
            MsgUndelegate: {
                encode(message: _63.MsgUndelegate, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgUndelegate;
                fromJSON(object: any): _63.MsgUndelegate;
                toJSON(message: _63.MsgUndelegate): unknown;
                fromPartial(object: Partial<_63.MsgUndelegate>): _63.MsgUndelegate;
            };
            MsgUndelegateResponse: {
                encode(message: _63.MsgUndelegateResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _63.MsgUndelegateResponse;
                fromJSON(object: any): _63.MsgUndelegateResponse;
                toJSON(message: _63.MsgUndelegateResponse): unknown;
                fromPartial(object: Partial<_63.MsgUndelegateResponse>): _63.MsgUndelegateResponse;
            };
            bondStatusFromJSON(object: any): _62.BondStatus;
            bondStatusToJSON(object: _62.BondStatus): string;
            BondStatus: typeof _62.BondStatus;
            BondStatusSDKType: typeof _62.BondStatusSDKType;
            HistoricalInfo: {
                encode(message: _62.HistoricalInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.HistoricalInfo;
                fromJSON(object: any): _62.HistoricalInfo;
                toJSON(message: _62.HistoricalInfo): unknown;
                fromPartial(object: Partial<_62.HistoricalInfo>): _62.HistoricalInfo;
            };
            CommissionRates: {
                encode(message: _62.CommissionRates, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.CommissionRates;
                fromJSON(object: any): _62.CommissionRates;
                toJSON(message: _62.CommissionRates): unknown;
                fromPartial(object: Partial<_62.CommissionRates>): _62.CommissionRates;
            };
            Commission: {
                encode(message: _62.Commission, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.Commission;
                fromJSON(object: any): _62.Commission;
                toJSON(message: _62.Commission): unknown;
                fromPartial(object: Partial<_62.Commission>): _62.Commission;
            };
            Description: {
                encode(message: _62.Description, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.Description;
                fromJSON(object: any): _62.Description;
                toJSON(message: _62.Description): unknown;
                fromPartial(object: Partial<_62.Description>): _62.Description;
            };
            Validator: {
                encode(message: _62.Validator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.Validator;
                fromJSON(object: any): _62.Validator;
                toJSON(message: _62.Validator): unknown;
                fromPartial(object: Partial<_62.Validator>): _62.Validator;
            };
            ValAddresses: {
                encode(message: _62.ValAddresses, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.ValAddresses;
                fromJSON(object: any): _62.ValAddresses;
                toJSON(message: _62.ValAddresses): unknown;
                fromPartial(object: Partial<_62.ValAddresses>): _62.ValAddresses;
            };
            DVPair: {
                encode(message: _62.DVPair, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.DVPair;
                fromJSON(object: any): _62.DVPair;
                toJSON(message: _62.DVPair): unknown;
                fromPartial(object: Partial<_62.DVPair>): _62.DVPair;
            };
            DVPairs: {
                encode(message: _62.DVPairs, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.DVPairs;
                fromJSON(object: any): _62.DVPairs;
                toJSON(message: _62.DVPairs): unknown;
                fromPartial(object: Partial<_62.DVPairs>): _62.DVPairs;
            };
            DVVTriplet: {
                encode(message: _62.DVVTriplet, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.DVVTriplet;
                fromJSON(object: any): _62.DVVTriplet;
                toJSON(message: _62.DVVTriplet): unknown;
                fromPartial(object: Partial<_62.DVVTriplet>): _62.DVVTriplet;
            };
            DVVTriplets: {
                encode(message: _62.DVVTriplets, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.DVVTriplets;
                fromJSON(object: any): _62.DVVTriplets;
                toJSON(message: _62.DVVTriplets): unknown;
                fromPartial(object: Partial<_62.DVVTriplets>): _62.DVVTriplets;
            };
            Delegation: {
                encode(message: _62.Delegation, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.Delegation;
                fromJSON(object: any): _62.Delegation;
                toJSON(message: _62.Delegation): unknown;
                fromPartial(object: Partial<_62.Delegation>): _62.Delegation;
            };
            UnbondingDelegation: {
                encode(message: _62.UnbondingDelegation, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.UnbondingDelegation;
                fromJSON(object: any): _62.UnbondingDelegation;
                toJSON(message: _62.UnbondingDelegation): unknown;
                fromPartial(object: Partial<_62.UnbondingDelegation>): _62.UnbondingDelegation;
            };
            UnbondingDelegationEntry: {
                encode(message: _62.UnbondingDelegationEntry, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.UnbondingDelegationEntry;
                fromJSON(object: any): _62.UnbondingDelegationEntry;
                toJSON(message: _62.UnbondingDelegationEntry): unknown;
                fromPartial(object: Partial<_62.UnbondingDelegationEntry>): _62.UnbondingDelegationEntry;
            };
            RedelegationEntry: {
                encode(message: _62.RedelegationEntry, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.RedelegationEntry;
                fromJSON(object: any): _62.RedelegationEntry;
                toJSON(message: _62.RedelegationEntry): unknown;
                fromPartial(object: Partial<_62.RedelegationEntry>): _62.RedelegationEntry;
            };
            Redelegation: {
                encode(message: _62.Redelegation, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.Redelegation;
                fromJSON(object: any): _62.Redelegation;
                toJSON(message: _62.Redelegation): unknown;
                fromPartial(object: Partial<_62.Redelegation>): _62.Redelegation;
            };
            Params: {
                encode(message: _62.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.Params;
                fromJSON(object: any): _62.Params;
                toJSON(message: _62.Params): unknown;
                fromPartial(object: Partial<_62.Params>): _62.Params;
            };
            DelegationResponse: {
                encode(message: _62.DelegationResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.DelegationResponse;
                fromJSON(object: any): _62.DelegationResponse;
                toJSON(message: _62.DelegationResponse): unknown;
                fromPartial(object: Partial<_62.DelegationResponse>): _62.DelegationResponse;
            };
            RedelegationEntryResponse: {
                encode(message: _62.RedelegationEntryResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.RedelegationEntryResponse;
                fromJSON(object: any): _62.RedelegationEntryResponse;
                toJSON(message: _62.RedelegationEntryResponse): unknown;
                fromPartial(object: Partial<_62.RedelegationEntryResponse>): _62.RedelegationEntryResponse;
            };
            RedelegationResponse: {
                encode(message: _62.RedelegationResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.RedelegationResponse;
                fromJSON(object: any): _62.RedelegationResponse;
                toJSON(message: _62.RedelegationResponse): unknown;
                fromPartial(object: Partial<_62.RedelegationResponse>): _62.RedelegationResponse;
            };
            Pool: {
                encode(message: _62.Pool, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _62.Pool;
                fromJSON(object: any): _62.Pool;
                toJSON(message: _62.Pool): unknown;
                fromPartial(object: Partial<_62.Pool>): _62.Pool;
            };
            QueryValidatorsRequest: {
                encode(message: _61.QueryValidatorsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorsRequest;
                fromJSON(object: any): _61.QueryValidatorsRequest;
                toJSON(message: _61.QueryValidatorsRequest): unknown;
                fromPartial(object: Partial<_61.QueryValidatorsRequest>): _61.QueryValidatorsRequest;
            };
            QueryValidatorsResponse: {
                encode(message: _61.QueryValidatorsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorsResponse;
                fromJSON(object: any): _61.QueryValidatorsResponse;
                toJSON(message: _61.QueryValidatorsResponse): unknown;
                fromPartial(object: Partial<_61.QueryValidatorsResponse>): _61.QueryValidatorsResponse;
            };
            QueryValidatorRequest: {
                encode(message: _61.QueryValidatorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorRequest;
                fromJSON(object: any): _61.QueryValidatorRequest;
                toJSON(message: _61.QueryValidatorRequest): unknown;
                fromPartial(object: Partial<_61.QueryValidatorRequest>): _61.QueryValidatorRequest;
            };
            QueryValidatorResponse: {
                encode(message: _61.QueryValidatorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorResponse;
                fromJSON(object: any): _61.QueryValidatorResponse;
                toJSON(message: _61.QueryValidatorResponse): unknown;
                fromPartial(object: Partial<_61.QueryValidatorResponse>): _61.QueryValidatorResponse;
            };
            QueryValidatorDelegationsRequest: {
                encode(message: _61.QueryValidatorDelegationsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorDelegationsRequest;
                fromJSON(object: any): _61.QueryValidatorDelegationsRequest;
                toJSON(message: _61.QueryValidatorDelegationsRequest): unknown;
                fromPartial(object: Partial<_61.QueryValidatorDelegationsRequest>): _61.QueryValidatorDelegationsRequest;
            };
            QueryValidatorDelegationsResponse: {
                encode(message: _61.QueryValidatorDelegationsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorDelegationsResponse;
                fromJSON(object: any): _61.QueryValidatorDelegationsResponse;
                toJSON(message: _61.QueryValidatorDelegationsResponse): unknown;
                fromPartial(object: Partial<_61.QueryValidatorDelegationsResponse>): _61.QueryValidatorDelegationsResponse;
            };
            QueryValidatorUnbondingDelegationsRequest: {
                encode(message: _61.QueryValidatorUnbondingDelegationsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorUnbondingDelegationsRequest;
                fromJSON(object: any): _61.QueryValidatorUnbondingDelegationsRequest;
                toJSON(message: _61.QueryValidatorUnbondingDelegationsRequest): unknown;
                fromPartial(object: Partial<_61.QueryValidatorUnbondingDelegationsRequest>): _61.QueryValidatorUnbondingDelegationsRequest;
            };
            QueryValidatorUnbondingDelegationsResponse: {
                encode(message: _61.QueryValidatorUnbondingDelegationsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryValidatorUnbondingDelegationsResponse;
                fromJSON(object: any): _61.QueryValidatorUnbondingDelegationsResponse;
                toJSON(message: _61.QueryValidatorUnbondingDelegationsResponse): unknown;
                fromPartial(object: Partial<_61.QueryValidatorUnbondingDelegationsResponse>): _61.QueryValidatorUnbondingDelegationsResponse;
            };
            QueryDelegationRequest: {
                encode(message: _61.QueryDelegationRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegationRequest;
                fromJSON(object: any): _61.QueryDelegationRequest;
                toJSON(message: _61.QueryDelegationRequest): unknown;
                fromPartial(object: Partial<_61.QueryDelegationRequest>): _61.QueryDelegationRequest;
            };
            QueryDelegationResponse: {
                encode(message: _61.QueryDelegationResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegationResponse;
                fromJSON(object: any): _61.QueryDelegationResponse;
                toJSON(message: _61.QueryDelegationResponse): unknown;
                fromPartial(object: Partial<_61.QueryDelegationResponse>): _61.QueryDelegationResponse;
            };
            QueryUnbondingDelegationRequest: {
                encode(message: _61.QueryUnbondingDelegationRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryUnbondingDelegationRequest;
                fromJSON(object: any): _61.QueryUnbondingDelegationRequest;
                toJSON(message: _61.QueryUnbondingDelegationRequest): unknown;
                fromPartial(object: Partial<_61.QueryUnbondingDelegationRequest>): _61.QueryUnbondingDelegationRequest;
            };
            QueryUnbondingDelegationResponse: {
                encode(message: _61.QueryUnbondingDelegationResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryUnbondingDelegationResponse;
                fromJSON(object: any): _61.QueryUnbondingDelegationResponse;
                toJSON(message: _61.QueryUnbondingDelegationResponse): unknown;
                fromPartial(object: Partial<_61.QueryUnbondingDelegationResponse>): _61.QueryUnbondingDelegationResponse;
            };
            QueryDelegatorDelegationsRequest: {
                encode(message: _61.QueryDelegatorDelegationsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorDelegationsRequest;
                fromJSON(object: any): _61.QueryDelegatorDelegationsRequest;
                toJSON(message: _61.QueryDelegatorDelegationsRequest): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorDelegationsRequest>): _61.QueryDelegatorDelegationsRequest;
            };
            QueryDelegatorDelegationsResponse: {
                encode(message: _61.QueryDelegatorDelegationsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorDelegationsResponse;
                fromJSON(object: any): _61.QueryDelegatorDelegationsResponse;
                toJSON(message: _61.QueryDelegatorDelegationsResponse): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorDelegationsResponse>): _61.QueryDelegatorDelegationsResponse;
            };
            QueryDelegatorUnbondingDelegationsRequest: {
                encode(message: _61.QueryDelegatorUnbondingDelegationsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorUnbondingDelegationsRequest;
                fromJSON(object: any): _61.QueryDelegatorUnbondingDelegationsRequest;
                toJSON(message: _61.QueryDelegatorUnbondingDelegationsRequest): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorUnbondingDelegationsRequest>): _61.QueryDelegatorUnbondingDelegationsRequest;
            };
            QueryDelegatorUnbondingDelegationsResponse: {
                encode(message: _61.QueryDelegatorUnbondingDelegationsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorUnbondingDelegationsResponse;
                fromJSON(object: any): _61.QueryDelegatorUnbondingDelegationsResponse;
                toJSON(message: _61.QueryDelegatorUnbondingDelegationsResponse): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorUnbondingDelegationsResponse>): _61.QueryDelegatorUnbondingDelegationsResponse;
            };
            QueryRedelegationsRequest: {
                encode(message: _61.QueryRedelegationsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryRedelegationsRequest;
                fromJSON(object: any): _61.QueryRedelegationsRequest;
                toJSON(message: _61.QueryRedelegationsRequest): unknown;
                fromPartial(object: Partial<_61.QueryRedelegationsRequest>): _61.QueryRedelegationsRequest;
            };
            QueryRedelegationsResponse: {
                encode(message: _61.QueryRedelegationsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryRedelegationsResponse;
                fromJSON(object: any): _61.QueryRedelegationsResponse;
                toJSON(message: _61.QueryRedelegationsResponse): unknown;
                fromPartial(object: Partial<_61.QueryRedelegationsResponse>): _61.QueryRedelegationsResponse;
            };
            QueryDelegatorValidatorsRequest: {
                encode(message: _61.QueryDelegatorValidatorsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorValidatorsRequest;
                fromJSON(object: any): _61.QueryDelegatorValidatorsRequest;
                toJSON(message: _61.QueryDelegatorValidatorsRequest): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorValidatorsRequest>): _61.QueryDelegatorValidatorsRequest;
            };
            QueryDelegatorValidatorsResponse: {
                encode(message: _61.QueryDelegatorValidatorsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorValidatorsResponse;
                fromJSON(object: any): _61.QueryDelegatorValidatorsResponse;
                toJSON(message: _61.QueryDelegatorValidatorsResponse): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorValidatorsResponse>): _61.QueryDelegatorValidatorsResponse;
            };
            QueryDelegatorValidatorRequest: {
                encode(message: _61.QueryDelegatorValidatorRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorValidatorRequest;
                fromJSON(object: any): _61.QueryDelegatorValidatorRequest;
                toJSON(message: _61.QueryDelegatorValidatorRequest): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorValidatorRequest>): _61.QueryDelegatorValidatorRequest;
            };
            QueryDelegatorValidatorResponse: {
                encode(message: _61.QueryDelegatorValidatorResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryDelegatorValidatorResponse;
                fromJSON(object: any): _61.QueryDelegatorValidatorResponse;
                toJSON(message: _61.QueryDelegatorValidatorResponse): unknown;
                fromPartial(object: Partial<_61.QueryDelegatorValidatorResponse>): _61.QueryDelegatorValidatorResponse;
            };
            QueryHistoricalInfoRequest: {
                encode(message: _61.QueryHistoricalInfoRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryHistoricalInfoRequest;
                fromJSON(object: any): _61.QueryHistoricalInfoRequest;
                toJSON(message: _61.QueryHistoricalInfoRequest): unknown;
                fromPartial(object: Partial<_61.QueryHistoricalInfoRequest>): _61.QueryHistoricalInfoRequest;
            };
            QueryHistoricalInfoResponse: {
                encode(message: _61.QueryHistoricalInfoResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryHistoricalInfoResponse;
                fromJSON(object: any): _61.QueryHistoricalInfoResponse;
                toJSON(message: _61.QueryHistoricalInfoResponse): unknown;
                fromPartial(object: Partial<_61.QueryHistoricalInfoResponse>): _61.QueryHistoricalInfoResponse;
            };
            QueryPoolRequest: {
                encode(_: _61.QueryPoolRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryPoolRequest;
                fromJSON(_: any): _61.QueryPoolRequest;
                toJSON(_: _61.QueryPoolRequest): unknown;
                fromPartial(_: Partial<_61.QueryPoolRequest>): _61.QueryPoolRequest;
            };
            QueryPoolResponse: {
                encode(message: _61.QueryPoolResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryPoolResponse;
                fromJSON(object: any): _61.QueryPoolResponse;
                toJSON(message: _61.QueryPoolResponse): unknown;
                fromPartial(object: Partial<_61.QueryPoolResponse>): _61.QueryPoolResponse;
            };
            QueryParamsRequest: {
                encode(_: _61.QueryParamsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryParamsRequest;
                fromJSON(_: any): _61.QueryParamsRequest;
                toJSON(_: _61.QueryParamsRequest): unknown;
                fromPartial(_: Partial<_61.QueryParamsRequest>): _61.QueryParamsRequest;
            };
            QueryParamsResponse: {
                encode(message: _61.QueryParamsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _61.QueryParamsResponse;
                fromJSON(object: any): _61.QueryParamsResponse;
                toJSON(message: _61.QueryParamsResponse): unknown;
                fromPartial(object: Partial<_61.QueryParamsResponse>): _61.QueryParamsResponse;
            };
            GenesisState: {
                encode(message: _60.GenesisState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _60.GenesisState;
                fromJSON(object: any): _60.GenesisState;
                toJSON(message: _60.GenesisState): unknown;
                fromPartial(object: Partial<_60.GenesisState>): _60.GenesisState;
            };
            LastValidatorPower: {
                encode(message: _60.LastValidatorPower, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _60.LastValidatorPower;
                fromJSON(object: any): _60.LastValidatorPower;
                toJSON(message: _60.LastValidatorPower): unknown;
                fromPartial(object: Partial<_60.LastValidatorPower>): _60.LastValidatorPower;
            };
            authorizationTypeFromJSON(object: any): _59.AuthorizationType;
            authorizationTypeToJSON(object: _59.AuthorizationType): string;
            AuthorizationType: typeof _59.AuthorizationType;
            AuthorizationTypeSDKType: typeof _59.AuthorizationTypeSDKType;
            StakeAuthorization: {
                encode(message: _59.StakeAuthorization, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _59.StakeAuthorization;
                fromJSON(object: any): _59.StakeAuthorization;
                toJSON(message: _59.StakeAuthorization): unknown;
                fromPartial(object: Partial<_59.StakeAuthorization>): _59.StakeAuthorization;
            };
            StakeAuthorization_Validators: {
                encode(message: _59.StakeAuthorization_Validators, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _59.StakeAuthorization_Validators;
                fromJSON(object: any): _59.StakeAuthorization_Validators;
                toJSON(message: _59.StakeAuthorization_Validators): unknown;
                fromPartial(object: Partial<_59.StakeAuthorization_Validators>): _59.StakeAuthorization_Validators;
            };
        };
    }
    namespace tx {
        namespace signing {
            const v1beta1: {
                signModeFromJSON(object: any): _64.SignMode;
                signModeToJSON(object: _64.SignMode): string;
                SignMode: typeof _64.SignMode;
                SignModeSDKType: typeof _64.SignModeSDKType;
                SignatureDescriptors: {
                    encode(message: _64.SignatureDescriptors, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _64.SignatureDescriptors;
                    fromJSON(object: any): _64.SignatureDescriptors;
                    toJSON(message: _64.SignatureDescriptors): unknown;
                    fromPartial(object: Partial<_64.SignatureDescriptors>): _64.SignatureDescriptors;
                };
                SignatureDescriptor: {
                    encode(message: _64.SignatureDescriptor, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _64.SignatureDescriptor;
                    fromJSON(object: any): _64.SignatureDescriptor;
                    toJSON(message: _64.SignatureDescriptor): unknown;
                    fromPartial(object: Partial<_64.SignatureDescriptor>): _64.SignatureDescriptor;
                };
                SignatureDescriptor_Data: {
                    encode(message: _64.SignatureDescriptor_Data, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _64.SignatureDescriptor_Data;
                    fromJSON(object: any): _64.SignatureDescriptor_Data;
                    toJSON(message: _64.SignatureDescriptor_Data): unknown;
                    fromPartial(object: Partial<_64.SignatureDescriptor_Data>): _64.SignatureDescriptor_Data;
                };
                SignatureDescriptor_Data_Single: {
                    encode(message: _64.SignatureDescriptor_Data_Single, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _64.SignatureDescriptor_Data_Single;
                    fromJSON(object: any): _64.SignatureDescriptor_Data_Single;
                    toJSON(message: _64.SignatureDescriptor_Data_Single): unknown;
                    fromPartial(object: Partial<_64.SignatureDescriptor_Data_Single>): _64.SignatureDescriptor_Data_Single;
                };
                SignatureDescriptor_Data_Multi: {
                    encode(message: _64.SignatureDescriptor_Data_Multi, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _64.SignatureDescriptor_Data_Multi;
                    fromJSON(object: any): _64.SignatureDescriptor_Data_Multi;
                    toJSON(message: _64.SignatureDescriptor_Data_Multi): unknown;
                    fromPartial(object: Partial<_64.SignatureDescriptor_Data_Multi>): _64.SignatureDescriptor_Data_Multi;
                };
            };
        }
        const v1beta1: {
            ServiceClientImpl: typeof _160.ServiceClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                simulate(request: _65.SimulateRequest): Promise<_65.SimulateResponse>;
                getTx(request: _65.GetTxRequest): Promise<_65.GetTxResponse>;
                broadcastTx(request: _65.BroadcastTxRequest): Promise<_65.BroadcastTxResponse>;
                getTxsEvent(request: _65.GetTxsEventRequest): Promise<_65.GetTxsEventResponse>;
            };
            Tx: {
                encode(message: _66.Tx, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.Tx;
                fromJSON(object: any): _66.Tx;
                toJSON(message: _66.Tx): unknown;
                fromPartial(object: Partial<_66.Tx>): _66.Tx;
            };
            TxRaw: {
                encode(message: _66.TxRaw, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.TxRaw;
                fromJSON(object: any): _66.TxRaw;
                toJSON(message: _66.TxRaw): unknown;
                fromPartial(object: Partial<_66.TxRaw>): _66.TxRaw;
            };
            SignDoc: {
                encode(message: _66.SignDoc, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.SignDoc;
                fromJSON(object: any): _66.SignDoc;
                toJSON(message: _66.SignDoc): unknown;
                fromPartial(object: Partial<_66.SignDoc>): _66.SignDoc;
            };
            TxBody: {
                encode(message: _66.TxBody, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.TxBody;
                fromJSON(object: any): _66.TxBody;
                toJSON(message: _66.TxBody): unknown;
                fromPartial(object: Partial<_66.TxBody>): _66.TxBody;
            };
            AuthInfo: {
                encode(message: _66.AuthInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.AuthInfo;
                fromJSON(object: any): _66.AuthInfo;
                toJSON(message: _66.AuthInfo): unknown;
                fromPartial(object: Partial<_66.AuthInfo>): _66.AuthInfo;
            };
            SignerInfo: {
                encode(message: _66.SignerInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.SignerInfo;
                fromJSON(object: any): _66.SignerInfo;
                toJSON(message: _66.SignerInfo): unknown;
                fromPartial(object: Partial<_66.SignerInfo>): _66.SignerInfo;
            };
            ModeInfo: {
                encode(message: _66.ModeInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.ModeInfo;
                fromJSON(object: any): _66.ModeInfo;
                toJSON(message: _66.ModeInfo): unknown;
                fromPartial(object: Partial<_66.ModeInfo>): _66.ModeInfo;
            };
            ModeInfo_Single: {
                encode(message: _66.ModeInfo_Single, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.ModeInfo_Single;
                fromJSON(object: any): _66.ModeInfo_Single;
                toJSON(message: _66.ModeInfo_Single): unknown;
                fromPartial(object: Partial<_66.ModeInfo_Single>): _66.ModeInfo_Single;
            };
            ModeInfo_Multi: {
                encode(message: _66.ModeInfo_Multi, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.ModeInfo_Multi;
                fromJSON(object: any): _66.ModeInfo_Multi;
                toJSON(message: _66.ModeInfo_Multi): unknown;
                fromPartial(object: Partial<_66.ModeInfo_Multi>): _66.ModeInfo_Multi;
            };
            Fee: {
                encode(message: _66.Fee, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _66.Fee;
                fromJSON(object: any): _66.Fee;
                toJSON(message: _66.Fee): unknown;
                fromPartial(object: Partial<_66.Fee>): _66.Fee;
            };
            orderByFromJSON(object: any): _65.OrderBy;
            orderByToJSON(object: _65.OrderBy): string;
            broadcastModeFromJSON(object: any): _65.BroadcastMode;
            broadcastModeToJSON(object: _65.BroadcastMode): string;
            OrderBy: typeof _65.OrderBy;
            OrderBySDKType: typeof _65.OrderBySDKType;
            BroadcastMode: typeof _65.BroadcastMode;
            BroadcastModeSDKType: typeof _65.BroadcastModeSDKType;
            GetTxsEventRequest: {
                encode(message: _65.GetTxsEventRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.GetTxsEventRequest;
                fromJSON(object: any): _65.GetTxsEventRequest;
                toJSON(message: _65.GetTxsEventRequest): unknown;
                fromPartial(object: Partial<_65.GetTxsEventRequest>): _65.GetTxsEventRequest;
            };
            GetTxsEventResponse: {
                encode(message: _65.GetTxsEventResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.GetTxsEventResponse;
                fromJSON(object: any): _65.GetTxsEventResponse;
                toJSON(message: _65.GetTxsEventResponse): unknown;
                fromPartial(object: Partial<_65.GetTxsEventResponse>): _65.GetTxsEventResponse;
            };
            BroadcastTxRequest: {
                encode(message: _65.BroadcastTxRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.BroadcastTxRequest;
                fromJSON(object: any): _65.BroadcastTxRequest;
                toJSON(message: _65.BroadcastTxRequest): unknown;
                fromPartial(object: Partial<_65.BroadcastTxRequest>): _65.BroadcastTxRequest;
            };
            BroadcastTxResponse: {
                encode(message: _65.BroadcastTxResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.BroadcastTxResponse;
                fromJSON(object: any): _65.BroadcastTxResponse;
                toJSON(message: _65.BroadcastTxResponse): unknown;
                fromPartial(object: Partial<_65.BroadcastTxResponse>): _65.BroadcastTxResponse;
            };
            SimulateRequest: {
                encode(message: _65.SimulateRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.SimulateRequest;
                fromJSON(object: any): _65.SimulateRequest;
                toJSON(message: _65.SimulateRequest): unknown;
                fromPartial(object: Partial<_65.SimulateRequest>): _65.SimulateRequest;
            };
            SimulateResponse: {
                encode(message: _65.SimulateResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.SimulateResponse;
                fromJSON(object: any): _65.SimulateResponse;
                toJSON(message: _65.SimulateResponse): unknown;
                fromPartial(object: Partial<_65.SimulateResponse>): _65.SimulateResponse;
            };
            GetTxRequest: {
                encode(message: _65.GetTxRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.GetTxRequest;
                fromJSON(object: any): _65.GetTxRequest;
                toJSON(message: _65.GetTxRequest): unknown;
                fromPartial(object: Partial<_65.GetTxRequest>): _65.GetTxRequest;
            };
            GetTxResponse: {
                encode(message: _65.GetTxResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _65.GetTxResponse;
                fromJSON(object: any): _65.GetTxResponse;
                toJSON(message: _65.GetTxResponse): unknown;
                fromPartial(object: Partial<_65.GetTxResponse>): _65.GetTxResponse;
            };
        };
    }
    namespace upgrade {
        const v1beta1: {
            QueryClientImpl: typeof _161.QueryClientImpl;
            createRpcQueryExtension: (base: import("@cosmjs/stargate").QueryClient) => {
                currentPlan(request?: _67.QueryCurrentPlanRequest): Promise<_67.QueryCurrentPlanResponse>;
                appliedPlan(request: _67.QueryAppliedPlanRequest): Promise<_67.QueryAppliedPlanResponse>;
                upgradedConsensusState(request: _67.QueryUpgradedConsensusStateRequest): Promise<_67.QueryUpgradedConsensusStateResponse>;
                moduleVersions(request: _67.QueryModuleVersionsRequest): Promise<_67.QueryModuleVersionsResponse>;
            };
            Plan: {
                encode(message: _68.Plan, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _68.Plan;
                fromJSON(object: any): _68.Plan;
                toJSON(message: _68.Plan): unknown;
                fromPartial(object: Partial<_68.Plan>): _68.Plan;
            };
            SoftwareUpgradeProposal: {
                encode(message: _68.SoftwareUpgradeProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _68.SoftwareUpgradeProposal;
                fromJSON(object: any): _68.SoftwareUpgradeProposal;
                toJSON(message: _68.SoftwareUpgradeProposal): unknown;
                fromPartial(object: Partial<_68.SoftwareUpgradeProposal>): _68.SoftwareUpgradeProposal;
            };
            CancelSoftwareUpgradeProposal: {
                encode(message: _68.CancelSoftwareUpgradeProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _68.CancelSoftwareUpgradeProposal;
                fromJSON(object: any): _68.CancelSoftwareUpgradeProposal;
                toJSON(message: _68.CancelSoftwareUpgradeProposal): unknown;
                fromPartial(object: Partial<_68.CancelSoftwareUpgradeProposal>): _68.CancelSoftwareUpgradeProposal;
            };
            ModuleVersion: {
                encode(message: _68.ModuleVersion, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _68.ModuleVersion;
                fromJSON(object: any): _68.ModuleVersion;
                toJSON(message: _68.ModuleVersion): unknown;
                fromPartial(object: Partial<_68.ModuleVersion>): _68.ModuleVersion;
            };
            QueryCurrentPlanRequest: {
                encode(_: _67.QueryCurrentPlanRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryCurrentPlanRequest;
                fromJSON(_: any): _67.QueryCurrentPlanRequest;
                toJSON(_: _67.QueryCurrentPlanRequest): unknown;
                fromPartial(_: Partial<_67.QueryCurrentPlanRequest>): _67.QueryCurrentPlanRequest;
            };
            QueryCurrentPlanResponse: {
                encode(message: _67.QueryCurrentPlanResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryCurrentPlanResponse;
                fromJSON(object: any): _67.QueryCurrentPlanResponse;
                toJSON(message: _67.QueryCurrentPlanResponse): unknown;
                fromPartial(object: Partial<_67.QueryCurrentPlanResponse>): _67.QueryCurrentPlanResponse;
            };
            QueryAppliedPlanRequest: {
                encode(message: _67.QueryAppliedPlanRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryAppliedPlanRequest;
                fromJSON(object: any): _67.QueryAppliedPlanRequest;
                toJSON(message: _67.QueryAppliedPlanRequest): unknown;
                fromPartial(object: Partial<_67.QueryAppliedPlanRequest>): _67.QueryAppliedPlanRequest;
            };
            QueryAppliedPlanResponse: {
                encode(message: _67.QueryAppliedPlanResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryAppliedPlanResponse;
                fromJSON(object: any): _67.QueryAppliedPlanResponse;
                toJSON(message: _67.QueryAppliedPlanResponse): unknown;
                fromPartial(object: Partial<_67.QueryAppliedPlanResponse>): _67.QueryAppliedPlanResponse;
            };
            QueryUpgradedConsensusStateRequest: {
                encode(message: _67.QueryUpgradedConsensusStateRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryUpgradedConsensusStateRequest;
                fromJSON(object: any): _67.QueryUpgradedConsensusStateRequest;
                toJSON(message: _67.QueryUpgradedConsensusStateRequest): unknown;
                fromPartial(object: Partial<_67.QueryUpgradedConsensusStateRequest>): _67.QueryUpgradedConsensusStateRequest;
            };
            QueryUpgradedConsensusStateResponse: {
                encode(message: _67.QueryUpgradedConsensusStateResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryUpgradedConsensusStateResponse;
                fromJSON(object: any): _67.QueryUpgradedConsensusStateResponse;
                toJSON(message: _67.QueryUpgradedConsensusStateResponse): unknown;
                fromPartial(object: Partial<_67.QueryUpgradedConsensusStateResponse>): _67.QueryUpgradedConsensusStateResponse;
            };
            QueryModuleVersionsRequest: {
                encode(message: _67.QueryModuleVersionsRequest, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryModuleVersionsRequest;
                fromJSON(object: any): _67.QueryModuleVersionsRequest;
                toJSON(message: _67.QueryModuleVersionsRequest): unknown;
                fromPartial(object: Partial<_67.QueryModuleVersionsRequest>): _67.QueryModuleVersionsRequest;
            };
            QueryModuleVersionsResponse: {
                encode(message: _67.QueryModuleVersionsResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _67.QueryModuleVersionsResponse;
                fromJSON(object: any): _67.QueryModuleVersionsResponse;
                toJSON(message: _67.QueryModuleVersionsResponse): unknown;
                fromPartial(object: Partial<_67.QueryModuleVersionsResponse>): _67.QueryModuleVersionsResponse;
            };
        };
    }
    namespace vesting {
        const v1beta1: {
            MsgClientImpl: typeof _171.MsgClientImpl;
            registry: readonly [string, import("@cosmjs/proto-signing").GeneratedType][];
            load: (protoRegistry: import("@cosmjs/proto-signing").Registry) => void;
            MessageComposer: {
                encoded: {
                    createVestingAccount(value: _69.MsgCreateVestingAccount): {
                        typeUrl: string;
                        value: Uint8Array;
                    };
                };
                withTypeUrl: {
                    createVestingAccount(value: _69.MsgCreateVestingAccount): {
                        typeUrl: string;
                        value: _69.MsgCreateVestingAccount;
                    };
                };
                toJSON: {
                    createVestingAccount(value: _69.MsgCreateVestingAccount): {
                        typeUrl: string;
                        value: unknown;
                    };
                };
                fromJSON: {
                    createVestingAccount(value: any): {
                        typeUrl: string;
                        value: _69.MsgCreateVestingAccount;
                    };
                };
                fromPartial: {
                    createVestingAccount(value: _69.MsgCreateVestingAccount): {
                        typeUrl: string;
                        value: _69.MsgCreateVestingAccount;
                    };
                };
            };
            AminoConverter: {
                "/cosmos.vesting.v1beta1.MsgCreateVestingAccount": {
                    aminoType: string;
                    toAmino: ({ fromAddress, toAddress, amount, endTime, delayed }: _69.MsgCreateVestingAccount) => {
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
                    }) => _69.MsgCreateVestingAccount;
                };
            };
            BaseVestingAccount: {
                encode(message: _70.BaseVestingAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _70.BaseVestingAccount;
                fromJSON(object: any): _70.BaseVestingAccount;
                toJSON(message: _70.BaseVestingAccount): unknown;
                fromPartial(object: Partial<_70.BaseVestingAccount>): _70.BaseVestingAccount;
            };
            ContinuousVestingAccount: {
                encode(message: _70.ContinuousVestingAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _70.ContinuousVestingAccount;
                fromJSON(object: any): _70.ContinuousVestingAccount;
                toJSON(message: _70.ContinuousVestingAccount): unknown;
                fromPartial(object: Partial<_70.ContinuousVestingAccount>): _70.ContinuousVestingAccount;
            };
            DelayedVestingAccount: {
                encode(message: _70.DelayedVestingAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _70.DelayedVestingAccount;
                fromJSON(object: any): _70.DelayedVestingAccount;
                toJSON(message: _70.DelayedVestingAccount): unknown;
                fromPartial(object: Partial<_70.DelayedVestingAccount>): _70.DelayedVestingAccount;
            };
            Period: {
                encode(message: _70.Period, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _70.Period;
                fromJSON(object: any): _70.Period;
                toJSON(message: _70.Period): unknown;
                fromPartial(object: Partial<_70.Period>): _70.Period;
            };
            PeriodicVestingAccount: {
                encode(message: _70.PeriodicVestingAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _70.PeriodicVestingAccount;
                fromJSON(object: any): _70.PeriodicVestingAccount;
                toJSON(message: _70.PeriodicVestingAccount): unknown;
                fromPartial(object: Partial<_70.PeriodicVestingAccount>): _70.PeriodicVestingAccount;
            };
            PermanentLockedAccount: {
                encode(message: _70.PermanentLockedAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _70.PermanentLockedAccount;
                fromJSON(object: any): _70.PermanentLockedAccount;
                toJSON(message: _70.PermanentLockedAccount): unknown;
                fromPartial(object: Partial<_70.PermanentLockedAccount>): _70.PermanentLockedAccount;
            };
            MsgCreateVestingAccount: {
                encode(message: _69.MsgCreateVestingAccount, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _69.MsgCreateVestingAccount;
                fromJSON(object: any): _69.MsgCreateVestingAccount;
                toJSON(message: _69.MsgCreateVestingAccount): unknown;
                fromPartial(object: Partial<_69.MsgCreateVestingAccount>): _69.MsgCreateVestingAccount;
            };
            MsgCreateVestingAccountResponse: {
                encode(_: _69.MsgCreateVestingAccountResponse, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _69.MsgCreateVestingAccountResponse;
                fromJSON(_: any): _69.MsgCreateVestingAccountResponse;
                toJSON(_: _69.MsgCreateVestingAccountResponse): unknown;
                fromPartial(_: Partial<_69.MsgCreateVestingAccountResponse>): _69.MsgCreateVestingAccountResponse;
            };
        };
    }
    const ClientFactory: {
        createRPCMsgClient: ({ rpc }: {
            rpc: import("../helpers").Rpc;
        }) => Promise<{
            cosmos: {
                authz: {
                    v1beta1: _162.MsgClientImpl;
                };
                bank: {
                    v1beta1: _163.MsgClientImpl;
                };
                crisis: {
                    v1beta1: _164.MsgClientImpl;
                };
                distribution: {
                    v1beta1: _165.MsgClientImpl;
                };
                evidence: {
                    v1beta1: _166.MsgClientImpl;
                };
                feegrant: {
                    v1beta1: _167.MsgClientImpl;
                };
                gov: {
                    v1beta1: _168.MsgClientImpl;
                };
                slashing: {
                    v1beta1: _169.MsgClientImpl;
                };
                staking: {
                    v1beta1: _170.MsgClientImpl;
                };
                vesting: {
                    v1beta1: _171.MsgClientImpl;
                };
            };
        }>;
        createRPCQueryClient: ({ rpcEndpoint }: {
            rpcEndpoint: string;
        }) => Promise<{
            cosmos: {
                auth: {
                    v1beta1: {
                        accounts(request?: _3.QueryAccountsRequest): Promise<_3.QueryAccountsResponse>;
                        account(request: _3.QueryAccountRequest): Promise<_3.QueryAccountResponse>;
                        params(request?: _3.QueryParamsRequest): Promise<_3.QueryParamsResponse>;
                    };
                };
                authz: {
                    v1beta1: {
                        grants(request: _7.QueryGrantsRequest): Promise<_7.QueryGrantsResponse>;
                    };
                };
                bank: {
                    v1beta1: {
                        balance(request: _12.QueryBalanceRequest): Promise<_12.QueryBalanceResponse>;
                        allBalances(request: _12.QueryAllBalancesRequest): Promise<_12.QueryAllBalancesResponse>;
                        totalSupply(request?: _12.QueryTotalSupplyRequest): Promise<_12.QueryTotalSupplyResponse>;
                        supplyOf(request: _12.QuerySupplyOfRequest): Promise<_12.QuerySupplyOfResponse>;
                        params(request?: _12.QueryParamsRequest): Promise<_12.QueryParamsResponse>;
                        denomMetadata(request: _12.QueryDenomMetadataRequest): Promise<_12.QueryDenomMetadataResponse>;
                        denomsMetadata(request?: _12.QueryDenomsMetadataRequest): Promise<_12.QueryDenomsMetadataResponse>;
                    };
                };
                base: {
                    tendermint: {
                        v1beta1: {
                            getNodeInfo(request?: _23.GetNodeInfoRequest): Promise<_23.GetNodeInfoResponse>;
                            getSyncing(request?: _23.GetSyncingRequest): Promise<_23.GetSyncingResponse>;
                            getLatestBlock(request?: _23.GetLatestBlockRequest): Promise<_23.GetLatestBlockResponse>;
                            getBlockByHeight(request: _23.GetBlockByHeightRequest): Promise<_23.GetBlockByHeightResponse>;
                            getLatestValidatorSet(request?: _23.GetLatestValidatorSetRequest): Promise<_23.GetLatestValidatorSetResponse>;
                            getValidatorSetByHeight(request: _23.GetValidatorSetByHeightRequest): Promise<_23.GetValidatorSetByHeightResponse>;
                        };
                    };
                };
                distribution: {
                    v1beta1: {
                        params(request?: _35.QueryParamsRequest): Promise<_35.QueryParamsResponse>;
                        validatorOutstandingRewards(request: _35.QueryValidatorOutstandingRewardsRequest): Promise<_35.QueryValidatorOutstandingRewardsResponse>;
                        validatorCommission(request: _35.QueryValidatorCommissionRequest): Promise<_35.QueryValidatorCommissionResponse>;
                        validatorSlashes(request: _35.QueryValidatorSlashesRequest): Promise<_35.QueryValidatorSlashesResponse>;
                        delegationRewards(request: _35.QueryDelegationRewardsRequest): Promise<_35.QueryDelegationRewardsResponse>;
                        delegationTotalRewards(request: _35.QueryDelegationTotalRewardsRequest): Promise<_35.QueryDelegationTotalRewardsResponse>;
                        delegatorValidators(request: _35.QueryDelegatorValidatorsRequest): Promise<_35.QueryDelegatorValidatorsResponse>;
                        delegatorWithdrawAddress(request: _35.QueryDelegatorWithdrawAddressRequest): Promise<_35.QueryDelegatorWithdrawAddressResponse>;
                        communityPool(request?: _35.QueryCommunityPoolRequest): Promise<_35.QueryCommunityPoolResponse>;
                    };
                };
                evidence: {
                    v1beta1: {
                        evidence(request: _39.QueryEvidenceRequest): Promise<_39.QueryEvidenceResponse>;
                        allEvidence(request?: _39.QueryAllEvidenceRequest): Promise<_39.QueryAllEvidenceResponse>;
                    };
                };
                feegrant: {
                    v1beta1: {
                        allowance(request: _43.QueryAllowanceRequest): Promise<_43.QueryAllowanceResponse>;
                        allowances(request: _43.QueryAllowancesRequest): Promise<_43.QueryAllowancesResponse>;
                    };
                };
                gov: {
                    v1beta1: {
                        proposal(request: _48.QueryProposalRequest): Promise<_48.QueryProposalResponse>;
                        proposals(request: _48.QueryProposalsRequest): Promise<_48.QueryProposalsResponse>;
                        vote(request: _48.QueryVoteRequest): Promise<_48.QueryVoteResponse>;
                        votes(request: _48.QueryVotesRequest): Promise<_48.QueryVotesResponse>;
                        params(request: _48.QueryParamsRequest): Promise<_48.QueryParamsResponse>;
                        deposit(request: _48.QueryDepositRequest): Promise<_48.QueryDepositResponse>;
                        deposits(request: _48.QueryDepositsRequest): Promise<_48.QueryDepositsResponse>;
                        tallyResult(request: _48.QueryTallyResultRequest): Promise<_48.QueryTallyResultResponse>;
                    };
                };
                mint: {
                    v1beta1: {
                        params(request?: _52.QueryParamsRequest): Promise<_52.QueryParamsResponse>;
                        inflation(request?: _52.QueryInflationRequest): Promise<_52.QueryInflationResponse>;
                        annualProvisions(request?: _52.QueryAnnualProvisionsRequest): Promise<_52.QueryAnnualProvisionsResponse>;
                    };
                };
                params: {
                    v1beta1: {
                        params(request: _54.QueryParamsRequest): Promise<_54.QueryParamsResponse>;
                    };
                };
                slashing: {
                    v1beta1: {
                        params(request?: _56.QueryParamsRequest): Promise<_56.QueryParamsResponse>;
                        signingInfo(request: _56.QuerySigningInfoRequest): Promise<_56.QuerySigningInfoResponse>;
                        signingInfos(request?: _56.QuerySigningInfosRequest): Promise<_56.QuerySigningInfosResponse>;
                    };
                };
                staking: {
                    v1beta1: {
                        validators(request: _61.QueryValidatorsRequest): Promise<_61.QueryValidatorsResponse>;
                        validator(request: _61.QueryValidatorRequest): Promise<_61.QueryValidatorResponse>;
                        validatorDelegations(request: _61.QueryValidatorDelegationsRequest): Promise<_61.QueryValidatorDelegationsResponse>;
                        validatorUnbondingDelegations(request: _61.QueryValidatorUnbondingDelegationsRequest): Promise<_61.QueryValidatorUnbondingDelegationsResponse>;
                        delegation(request: _61.QueryDelegationRequest): Promise<_61.QueryDelegationResponse>;
                        unbondingDelegation(request: _61.QueryUnbondingDelegationRequest): Promise<_61.QueryUnbondingDelegationResponse>;
                        delegatorDelegations(request: _61.QueryDelegatorDelegationsRequest): Promise<_61.QueryDelegatorDelegationsResponse>;
                        delegatorUnbondingDelegations(request: _61.QueryDelegatorUnbondingDelegationsRequest): Promise<_61.QueryDelegatorUnbondingDelegationsResponse>;
                        redelegations(request: _61.QueryRedelegationsRequest): Promise<_61.QueryRedelegationsResponse>;
                        delegatorValidators(request: _61.QueryDelegatorValidatorsRequest): Promise<_61.QueryDelegatorValidatorsResponse>;
                        delegatorValidator(request: _61.QueryDelegatorValidatorRequest): Promise<_61.QueryDelegatorValidatorResponse>;
                        historicalInfo(request: _61.QueryHistoricalInfoRequest): Promise<_61.QueryHistoricalInfoResponse>;
                        pool(request?: _61.QueryPoolRequest): Promise<_61.QueryPoolResponse>;
                        params(request?: _61.QueryParamsRequest): Promise<_61.QueryParamsResponse>;
                    };
                };
                tx: {
                    v1beta1: {
                        simulate(request: _65.SimulateRequest): Promise<_65.SimulateResponse>;
                        getTx(request: _65.GetTxRequest): Promise<_65.GetTxResponse>;
                        broadcastTx(request: _65.BroadcastTxRequest): Promise<_65.BroadcastTxResponse>;
                        getTxsEvent(request: _65.GetTxsEventRequest): Promise<_65.GetTxsEventResponse>;
                    };
                };
                upgrade: {
                    v1beta1: {
                        currentPlan(request?: _67.QueryCurrentPlanRequest): Promise<_67.QueryCurrentPlanResponse>;
                        appliedPlan(request: _67.QueryAppliedPlanRequest): Promise<_67.QueryAppliedPlanResponse>;
                        upgradedConsensusState(request: _67.QueryUpgradedConsensusStateRequest): Promise<_67.QueryUpgradedConsensusStateResponse>;
                        moduleVersions(request: _67.QueryModuleVersionsRequest): Promise<_67.QueryModuleVersionsResponse>;
                    };
                };
            };
        }>;
    };
}
