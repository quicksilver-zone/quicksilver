import { Status, StatusSDKType, ZoneDrop, ZoneDropSDKType, ClaimRecord, ClaimRecordSDKType } from "./airdrop";
import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsSDKType } from "./params";
import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import * as _m0 from "protobufjs/minimal";
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequestSDKType {
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
    /** params defines the parameters of the module. */
    params?: Params;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponseSDKType {
    /** params defines the parameters of the module. */
    params?: ParamsSDKType;
}
/** QueryZoneDropRequest is the request type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropRequest {
    /** chain_id identifies the zone. */
    chainId: string;
}
/** QueryZoneDropRequest is the request type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropRequestSDKType {
    /** chain_id identifies the zone. */
    chain_id: string;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropResponse {
    zoneDrop?: ZoneDrop;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrop RPC method. */
export interface QueryZoneDropResponseSDKType {
    zone_drop?: ZoneDropSDKType;
}
/**
 * QueryAccountBalanceRequest is the request type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceRequest {
    /** chain_id identifies the zone. */
    chainId: string;
}
/**
 * QueryAccountBalanceRequest is the request type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceRequestSDKType {
    /** chain_id identifies the zone. */
    chain_id: string;
}
/**
 * QueryAccountBalanceResponse is the response type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceResponse {
    accountBalance?: Coin;
}
/**
 * QueryAccountBalanceResponse is the response type for Query/AccountBalance RPC
 * method.
 */
export interface QueryAccountBalanceResponseSDKType {
    account_balance?: CoinSDKType;
}
/** QueryZoneDropsRequest is the request type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsRequest {
    /**
     * status enables to query zone airdrops matching a given status:
     *  - Active
     *  - Future
     *  - Expired
     */
    status: Status;
    pagination?: PageRequest;
}
/** QueryZoneDropsRequest is the request type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsRequestSDKType {
    /**
     * status enables to query zone airdrops matching a given status:
     *  - Active
     *  - Future
     *  - Expired
     */
    status: StatusSDKType;
    pagination?: PageRequestSDKType;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsResponse {
    zoneDrops: ZoneDrop[];
    pagination?: PageResponse;
}
/** QueryZoneDropResponse is the response type for Query/ZoneDrops RPC method. */
export interface QueryZoneDropsResponseSDKType {
    zone_drops: ZoneDropSDKType[];
    pagination?: PageResponseSDKType;
}
/** QueryClaimRecordRequest is the request type for Query/ClaimRecord RPC method. */
export interface QueryClaimRecordRequest {
    chainId: string;
    address: string;
}
/** QueryClaimRecordRequest is the request type for Query/ClaimRecord RPC method. */
export interface QueryClaimRecordRequestSDKType {
    chain_id: string;
    address: string;
}
/**
 * QueryClaimRecordResponse is the response type for Query/ClaimRecord RPC
 * method.
 */
export interface QueryClaimRecordResponse {
    claimRecord?: ClaimRecord;
}
/**
 * QueryClaimRecordResponse is the response type for Query/ClaimRecord RPC
 * method.
 */
export interface QueryClaimRecordResponseSDKType {
    claim_record?: ClaimRecordSDKType;
}
/**
 * QueryClaimRecordsRequest is the request type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsRequest {
    chainId: string;
    pagination?: PageRequest;
}
/**
 * QueryClaimRecordsRequest is the request type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsRequestSDKType {
    chain_id: string;
    pagination?: PageRequestSDKType;
}
/**
 * QueryClaimRecordsResponse is the response type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsResponse {
    claimRecords: ClaimRecord[];
    pagination?: PageResponse;
}
/**
 * QueryClaimRecordsResponse is the response type for Query/ClaimRecords RPC
 * method.
 */
export interface QueryClaimRecordsResponseSDKType {
    claim_records: ClaimRecordSDKType[];
    pagination?: PageResponseSDKType;
}
export declare const QueryParamsRequest: {
    encode(_: QueryParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest;
    fromJSON(_: any): QueryParamsRequest;
    toJSON(_: QueryParamsRequest): unknown;
    fromPartial(_: Partial<QueryParamsRequest>): QueryParamsRequest;
};
export declare const QueryParamsResponse: {
    encode(message: QueryParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse;
    fromJSON(object: any): QueryParamsResponse;
    toJSON(message: QueryParamsResponse): unknown;
    fromPartial(object: Partial<QueryParamsResponse>): QueryParamsResponse;
};
export declare const QueryZoneDropRequest: {
    encode(message: QueryZoneDropRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropRequest;
    fromJSON(object: any): QueryZoneDropRequest;
    toJSON(message: QueryZoneDropRequest): unknown;
    fromPartial(object: Partial<QueryZoneDropRequest>): QueryZoneDropRequest;
};
export declare const QueryZoneDropResponse: {
    encode(message: QueryZoneDropResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropResponse;
    fromJSON(object: any): QueryZoneDropResponse;
    toJSON(message: QueryZoneDropResponse): unknown;
    fromPartial(object: Partial<QueryZoneDropResponse>): QueryZoneDropResponse;
};
export declare const QueryAccountBalanceRequest: {
    encode(message: QueryAccountBalanceRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAccountBalanceRequest;
    fromJSON(object: any): QueryAccountBalanceRequest;
    toJSON(message: QueryAccountBalanceRequest): unknown;
    fromPartial(object: Partial<QueryAccountBalanceRequest>): QueryAccountBalanceRequest;
};
export declare const QueryAccountBalanceResponse: {
    encode(message: QueryAccountBalanceResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAccountBalanceResponse;
    fromJSON(object: any): QueryAccountBalanceResponse;
    toJSON(message: QueryAccountBalanceResponse): unknown;
    fromPartial(object: Partial<QueryAccountBalanceResponse>): QueryAccountBalanceResponse;
};
export declare const QueryZoneDropsRequest: {
    encode(message: QueryZoneDropsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropsRequest;
    fromJSON(object: any): QueryZoneDropsRequest;
    toJSON(message: QueryZoneDropsRequest): unknown;
    fromPartial(object: Partial<QueryZoneDropsRequest>): QueryZoneDropsRequest;
};
export declare const QueryZoneDropsResponse: {
    encode(message: QueryZoneDropsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryZoneDropsResponse;
    fromJSON(object: any): QueryZoneDropsResponse;
    toJSON(message: QueryZoneDropsResponse): unknown;
    fromPartial(object: Partial<QueryZoneDropsResponse>): QueryZoneDropsResponse;
};
export declare const QueryClaimRecordRequest: {
    encode(message: QueryClaimRecordRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordRequest;
    fromJSON(object: any): QueryClaimRecordRequest;
    toJSON(message: QueryClaimRecordRequest): unknown;
    fromPartial(object: Partial<QueryClaimRecordRequest>): QueryClaimRecordRequest;
};
export declare const QueryClaimRecordResponse: {
    encode(message: QueryClaimRecordResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordResponse;
    fromJSON(object: any): QueryClaimRecordResponse;
    toJSON(message: QueryClaimRecordResponse): unknown;
    fromPartial(object: Partial<QueryClaimRecordResponse>): QueryClaimRecordResponse;
};
export declare const QueryClaimRecordsRequest: {
    encode(message: QueryClaimRecordsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordsRequest;
    fromJSON(object: any): QueryClaimRecordsRequest;
    toJSON(message: QueryClaimRecordsRequest): unknown;
    fromPartial(object: Partial<QueryClaimRecordsRequest>): QueryClaimRecordsRequest;
};
export declare const QueryClaimRecordsResponse: {
    encode(message: QueryClaimRecordsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimRecordsResponse;
    fromJSON(object: any): QueryClaimRecordsResponse;
    toJSON(message: QueryClaimRecordsResponse): unknown;
    fromPartial(object: Partial<QueryClaimRecordsResponse>): QueryClaimRecordsResponse;
};
