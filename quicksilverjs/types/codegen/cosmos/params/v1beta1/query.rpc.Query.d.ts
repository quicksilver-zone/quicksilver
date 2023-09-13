import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /**
     * Params queries a specific parameter of a module, given its subspace and
     * key.
     */
    params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
};
