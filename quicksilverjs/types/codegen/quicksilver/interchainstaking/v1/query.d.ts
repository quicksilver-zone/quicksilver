import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Zone, ZoneSDKType, DelegatorIntent, DelegatorIntentSDKType, Delegation, DelegationSDKType, Receipt, ReceiptSDKType, WithdrawalRecord, WithdrawalRecordSDKType, UnbondingRecord, UnbondingRecordSDKType, RedelegationRecord, RedelegationRecordSDKType } from "./interchainstaking";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
export interface QueryZonesInfoRequest {
    pagination?: PageRequest;
}
export interface QueryZonesInfoRequestSDKType {
    pagination?: PageRequestSDKType;
}
export interface QueryZonesInfoResponse {
    zones: Zone[];
    pagination?: PageResponse;
}
export interface QueryZonesInfoResponseSDKType {
    zones: ZoneSDKType[];
    pagination?: PageResponseSDKType;
}
/**
 * QueryDepositAccountForChainRequest is the request type for the
 * Query/InterchainAccountAddress RPC
 */
export interface QueryDepositAccountForChainRequest {
    chainId: string;
}
/**
 * QueryDepositAccountForChainRequest is the request type for the
 * Query/InterchainAccountAddress RPC
 */
export interface QueryDepositAccountForChainRequestSDKType {
    chain_id: string;
}
/**
 * QueryDepositAccountForChainResponse the response type for the
 * Query/InterchainAccountAddress RPC
 */
export interface QueryDepositAccountForChainResponse {
    depositAccountAddress: string;
}
/**
 * QueryDepositAccountForChainResponse the response type for the
 * Query/InterchainAccountAddress RPC
 */
export interface QueryDepositAccountForChainResponseSDKType {
    deposit_account_address: string;
}
export interface QueryDelegatorIntentRequest {
    chainId: string;
    delegatorAddress: string;
}
export interface QueryDelegatorIntentRequestSDKType {
    chain_id: string;
    delegator_address: string;
}
export interface QueryDelegatorIntentResponse {
    intent?: DelegatorIntent;
}
export interface QueryDelegatorIntentResponseSDKType {
    intent?: DelegatorIntentSDKType;
}
export interface QueryDelegationsRequest {
    chainId: string;
    pagination?: PageRequest;
}
export interface QueryDelegationsRequestSDKType {
    chain_id: string;
    pagination?: PageRequestSDKType;
}
export interface QueryDelegationsResponse {
    delegations: Delegation[];
    tvl: Long;
    pagination?: PageResponse;
}
export interface QueryDelegationsResponseSDKType {
    delegations: DelegationSDKType[];
    tvl: Long;
    pagination?: PageResponseSDKType;
}
export interface QueryReceiptsRequest {
    chainId: string;
    pagination?: PageRequest;
}
export interface QueryReceiptsRequestSDKType {
    chain_id: string;
    pagination?: PageRequestSDKType;
}
export interface QueryReceiptsResponse {
    receipts: Receipt[];
    pagination?: PageResponse;
}
export interface QueryReceiptsResponseSDKType {
    receipts: ReceiptSDKType[];
    pagination?: PageResponseSDKType;
}
export interface QueryWithdrawalRecordsRequest {
    chainId: string;
    delegatorAddress: string;
    pagination?: PageRequest;
}
export interface QueryWithdrawalRecordsRequestSDKType {
    chain_id: string;
    delegator_address: string;
    pagination?: PageRequestSDKType;
}
export interface QueryWithdrawalRecordsResponse {
    withdrawals: WithdrawalRecord[];
    pagination?: PageResponse;
}
export interface QueryWithdrawalRecordsResponseSDKType {
    withdrawals: WithdrawalRecordSDKType[];
    pagination?: PageResponseSDKType;
}
export interface QueryUnbondingRecordsRequest {
    chainId: string;
    pagination?: PageRequest;
}
export interface QueryUnbondingRecordsRequestSDKType {
    chain_id: string;
    pagination?: PageRequestSDKType;
}
export interface QueryUnbondingRecordsResponse {
    Unbondings: UnbondingRecord[];
    pagination?: PageResponse;
}
export interface QueryUnbondingRecordsResponseSDKType {
    Unbondings: UnbondingRecordSDKType[];
    pagination?: PageResponseSDKType;
}
export interface QueryRedelegationRecordsRequest {
    chainId: string;
    pagination?: PageRequest;
}
export interface QueryRedelegationRecordsRequestSDKType {
    chain_id: string;
    pagination?: PageRequestSDKType;
}
export interface QueryRedelegationRecordsResponse {
    Redelegations: RedelegationRecord[];
    pagination?: PageResponse;
}
export interface QueryRedelegationRecordsResponseSDKType {
    Redelegations: RedelegationRecordSDKType[];
    pagination?: PageResponseSDKType;
}
export declare const QueryZonesInfoRequest: {
    encode(message: QueryZonesInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryZonesInfoRequest;
    fromJSON(object: any): QueryZonesInfoRequest;
    toJSON(message: QueryZonesInfoRequest): unknown;
    fromPartial(object: Partial<QueryZonesInfoRequest>): QueryZonesInfoRequest;
};
export declare const QueryZonesInfoResponse: {
    encode(message: QueryZonesInfoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryZonesInfoResponse;
    fromJSON(object: any): QueryZonesInfoResponse;
    toJSON(message: QueryZonesInfoResponse): unknown;
    fromPartial(object: Partial<QueryZonesInfoResponse>): QueryZonesInfoResponse;
};
export declare const QueryDepositAccountForChainRequest: {
    encode(message: QueryDepositAccountForChainRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDepositAccountForChainRequest;
    fromJSON(object: any): QueryDepositAccountForChainRequest;
    toJSON(message: QueryDepositAccountForChainRequest): unknown;
    fromPartial(object: Partial<QueryDepositAccountForChainRequest>): QueryDepositAccountForChainRequest;
};
export declare const QueryDepositAccountForChainResponse: {
    encode(message: QueryDepositAccountForChainResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDepositAccountForChainResponse;
    fromJSON(object: any): QueryDepositAccountForChainResponse;
    toJSON(message: QueryDepositAccountForChainResponse): unknown;
    fromPartial(object: Partial<QueryDepositAccountForChainResponse>): QueryDepositAccountForChainResponse;
};
export declare const QueryDelegatorIntentRequest: {
    encode(message: QueryDelegatorIntentRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegatorIntentRequest;
    fromJSON(object: any): QueryDelegatorIntentRequest;
    toJSON(message: QueryDelegatorIntentRequest): unknown;
    fromPartial(object: Partial<QueryDelegatorIntentRequest>): QueryDelegatorIntentRequest;
};
export declare const QueryDelegatorIntentResponse: {
    encode(message: QueryDelegatorIntentResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegatorIntentResponse;
    fromJSON(object: any): QueryDelegatorIntentResponse;
    toJSON(message: QueryDelegatorIntentResponse): unknown;
    fromPartial(object: Partial<QueryDelegatorIntentResponse>): QueryDelegatorIntentResponse;
};
export declare const QueryDelegationsRequest: {
    encode(message: QueryDelegationsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegationsRequest;
    fromJSON(object: any): QueryDelegationsRequest;
    toJSON(message: QueryDelegationsRequest): unknown;
    fromPartial(object: Partial<QueryDelegationsRequest>): QueryDelegationsRequest;
};
export declare const QueryDelegationsResponse: {
    encode(message: QueryDelegationsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegationsResponse;
    fromJSON(object: any): QueryDelegationsResponse;
    toJSON(message: QueryDelegationsResponse): unknown;
    fromPartial(object: Partial<QueryDelegationsResponse>): QueryDelegationsResponse;
};
export declare const QueryReceiptsRequest: {
    encode(message: QueryReceiptsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryReceiptsRequest;
    fromJSON(object: any): QueryReceiptsRequest;
    toJSON(message: QueryReceiptsRequest): unknown;
    fromPartial(object: Partial<QueryReceiptsRequest>): QueryReceiptsRequest;
};
export declare const QueryReceiptsResponse: {
    encode(message: QueryReceiptsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryReceiptsResponse;
    fromJSON(object: any): QueryReceiptsResponse;
    toJSON(message: QueryReceiptsResponse): unknown;
    fromPartial(object: Partial<QueryReceiptsResponse>): QueryReceiptsResponse;
};
export declare const QueryWithdrawalRecordsRequest: {
    encode(message: QueryWithdrawalRecordsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryWithdrawalRecordsRequest;
    fromJSON(object: any): QueryWithdrawalRecordsRequest;
    toJSON(message: QueryWithdrawalRecordsRequest): unknown;
    fromPartial(object: Partial<QueryWithdrawalRecordsRequest>): QueryWithdrawalRecordsRequest;
};
export declare const QueryWithdrawalRecordsResponse: {
    encode(message: QueryWithdrawalRecordsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryWithdrawalRecordsResponse;
    fromJSON(object: any): QueryWithdrawalRecordsResponse;
    toJSON(message: QueryWithdrawalRecordsResponse): unknown;
    fromPartial(object: Partial<QueryWithdrawalRecordsResponse>): QueryWithdrawalRecordsResponse;
};
export declare const QueryUnbondingRecordsRequest: {
    encode(message: QueryUnbondingRecordsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryUnbondingRecordsRequest;
    fromJSON(object: any): QueryUnbondingRecordsRequest;
    toJSON(message: QueryUnbondingRecordsRequest): unknown;
    fromPartial(object: Partial<QueryUnbondingRecordsRequest>): QueryUnbondingRecordsRequest;
};
export declare const QueryUnbondingRecordsResponse: {
    encode(message: QueryUnbondingRecordsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryUnbondingRecordsResponse;
    fromJSON(object: any): QueryUnbondingRecordsResponse;
    toJSON(message: QueryUnbondingRecordsResponse): unknown;
    fromPartial(object: Partial<QueryUnbondingRecordsResponse>): QueryUnbondingRecordsResponse;
};
export declare const QueryRedelegationRecordsRequest: {
    encode(message: QueryRedelegationRecordsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryRedelegationRecordsRequest;
    fromJSON(object: any): QueryRedelegationRecordsRequest;
    toJSON(message: QueryRedelegationRecordsRequest): unknown;
    fromPartial(object: Partial<QueryRedelegationRecordsRequest>): QueryRedelegationRecordsRequest;
};
export declare const QueryRedelegationRecordsResponse: {
    encode(message: QueryRedelegationRecordsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryRedelegationRecordsResponse;
    fromJSON(object: any): QueryRedelegationRecordsResponse;
    toJSON(message: QueryRedelegationRecordsResponse): unknown;
    fromPartial(object: Partial<QueryRedelegationRecordsResponse>): QueryRedelegationRecordsResponse;
};
