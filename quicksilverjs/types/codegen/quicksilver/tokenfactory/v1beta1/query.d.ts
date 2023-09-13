import { Params, ParamsSDKType } from "./params";
import { DenomAuthorityMetadata, DenomAuthorityMetadataSDKType } from "./authorityMetadata";
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
/**
 * QueryDenomAuthorityMetadataRequest defines the request structure for the
 * DenomAuthorityMetadata gRPC query.
 */
export interface QueryDenomAuthorityMetadataRequest {
    denom: string;
}
/**
 * QueryDenomAuthorityMetadataRequest defines the request structure for the
 * DenomAuthorityMetadata gRPC query.
 */
export interface QueryDenomAuthorityMetadataRequestSDKType {
    denom: string;
}
/**
 * QueryDenomAuthorityMetadataResponse defines the response structure for the
 * DenomAuthorityMetadata gRPC query.
 */
export interface QueryDenomAuthorityMetadataResponse {
    authorityMetadata?: DenomAuthorityMetadata;
}
/**
 * QueryDenomAuthorityMetadataResponse defines the response structure for the
 * DenomAuthorityMetadata gRPC query.
 */
export interface QueryDenomAuthorityMetadataResponseSDKType {
    authority_metadata?: DenomAuthorityMetadataSDKType;
}
/**
 * QueryDenomsFromCreatorRequest defines the request structure for the
 * DenomsFromCreator gRPC query.
 */
export interface QueryDenomsFromCreatorRequest {
    creator: string;
}
/**
 * QueryDenomsFromCreatorRequest defines the request structure for the
 * DenomsFromCreator gRPC query.
 */
export interface QueryDenomsFromCreatorRequestSDKType {
    creator: string;
}
/**
 * QueryDenomsFromCreatorRequest defines the response structure for the
 * DenomsFromCreator gRPC query.
 */
export interface QueryDenomsFromCreatorResponse {
    denoms: string[];
}
/**
 * QueryDenomsFromCreatorRequest defines the response structure for the
 * DenomsFromCreator gRPC query.
 */
export interface QueryDenomsFromCreatorResponseSDKType {
    denoms: string[];
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
export declare const QueryDenomAuthorityMetadataRequest: {
    encode(message: QueryDenomAuthorityMetadataRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDenomAuthorityMetadataRequest;
    fromJSON(object: any): QueryDenomAuthorityMetadataRequest;
    toJSON(message: QueryDenomAuthorityMetadataRequest): unknown;
    fromPartial(object: Partial<QueryDenomAuthorityMetadataRequest>): QueryDenomAuthorityMetadataRequest;
};
export declare const QueryDenomAuthorityMetadataResponse: {
    encode(message: QueryDenomAuthorityMetadataResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDenomAuthorityMetadataResponse;
    fromJSON(object: any): QueryDenomAuthorityMetadataResponse;
    toJSON(message: QueryDenomAuthorityMetadataResponse): unknown;
    fromPartial(object: Partial<QueryDenomAuthorityMetadataResponse>): QueryDenomAuthorityMetadataResponse;
};
export declare const QueryDenomsFromCreatorRequest: {
    encode(message: QueryDenomsFromCreatorRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDenomsFromCreatorRequest;
    fromJSON(object: any): QueryDenomsFromCreatorRequest;
    toJSON(message: QueryDenomsFromCreatorRequest): unknown;
    fromPartial(object: Partial<QueryDenomsFromCreatorRequest>): QueryDenomsFromCreatorRequest;
};
export declare const QueryDenomsFromCreatorResponse: {
    encode(message: QueryDenomsFromCreatorResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDenomsFromCreatorResponse;
    fromJSON(object: any): QueryDenomsFromCreatorResponse;
    toJSON(message: QueryDenomsFromCreatorResponse): unknown;
    fromPartial(object: Partial<QueryDenomsFromCreatorResponse>): QueryDenomsFromCreatorResponse;
};
