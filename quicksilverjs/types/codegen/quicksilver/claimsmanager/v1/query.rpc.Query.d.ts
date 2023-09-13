import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryClaimsRequest, QueryClaimsResponse } from "./query";
/** Query provides defines the gRPC querier service. */
export interface Query {
    /** Params returns the total set of participation rewards parameters. */
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    claims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    lastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    userClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    userLastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    claims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    lastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    userClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    userLastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    claims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    lastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    userClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
    userLastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
};
