import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryGrantsRequest, QueryGrantsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Returns list of `Authorization`, granted to the grantee by the granter. */
    grants(request: QueryGrantsRequest): Promise<QueryGrantsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    grants(request: QueryGrantsRequest): Promise<QueryGrantsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    grants(request: QueryGrantsRequest): Promise<QueryGrantsResponse>;
};
