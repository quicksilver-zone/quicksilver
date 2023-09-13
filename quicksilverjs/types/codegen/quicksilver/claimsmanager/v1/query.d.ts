import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsSDKType, Claim, ClaimSDKType } from "./claimsmanager";
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
export interface QueryClaimsRequest {
    chainId: string;
    address: string;
    pagination?: PageRequest;
}
export interface QueryClaimsRequestSDKType {
    chain_id: string;
    address: string;
    pagination?: PageRequestSDKType;
}
export interface QueryClaimsResponse {
    claims: Claim[];
    pagination?: PageResponse;
}
export interface QueryClaimsResponseSDKType {
    claims: ClaimSDKType[];
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
export declare const QueryClaimsRequest: {
    encode(message: QueryClaimsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimsRequest;
    fromJSON(object: any): QueryClaimsRequest;
    toJSON(message: QueryClaimsRequest): unknown;
    fromPartial(object: Partial<QueryClaimsRequest>): QueryClaimsRequest;
};
export declare const QueryClaimsResponse: {
    encode(message: QueryClaimsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryClaimsResponse;
    fromJSON(object: any): QueryClaimsResponse;
    toJSON(message: QueryClaimsResponse): unknown;
    fromPartial(object: Partial<QueryClaimsResponse>): QueryClaimsResponse;
};
