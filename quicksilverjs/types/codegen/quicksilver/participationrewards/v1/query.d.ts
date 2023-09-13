import { Params, ParamsSDKType } from "./participationrewards";
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
/** QueryProtocolDataRequest is the request type for querying Protocol Data. */
export interface QueryProtocolDataRequest {
    type: string;
    key: string;
}
/** QueryProtocolDataRequest is the request type for querying Protocol Data. */
export interface QueryProtocolDataRequestSDKType {
    type: string;
    key: string;
}
/** QueryProtocolDataResponse is the response type for querying Protocol Data. */
export interface QueryProtocolDataResponse {
    /** data defines the underlying protocol data. */
    data: Uint8Array[];
}
/** QueryProtocolDataResponse is the response type for querying Protocol Data. */
export interface QueryProtocolDataResponseSDKType {
    /** data defines the underlying protocol data. */
    data: Uint8Array[];
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
export declare const QueryProtocolDataRequest: {
    encode(message: QueryProtocolDataRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryProtocolDataRequest;
    fromJSON(object: any): QueryProtocolDataRequest;
    toJSON(message: QueryProtocolDataRequest): unknown;
    fromPartial(object: Partial<QueryProtocolDataRequest>): QueryProtocolDataRequest;
};
export declare const QueryProtocolDataResponse: {
    encode(message: QueryProtocolDataResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryProtocolDataResponse;
    fromJSON(object: any): QueryProtocolDataResponse;
    toJSON(message: QueryProtocolDataResponse): unknown;
    fromPartial(object: Partial<QueryProtocolDataResponse>): QueryProtocolDataResponse;
};
