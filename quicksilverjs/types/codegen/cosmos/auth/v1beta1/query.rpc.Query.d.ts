import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryAccountsRequest, QueryAccountsResponse, QueryAccountRequest, QueryAccountResponse, QueryParamsRequest, QueryParamsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Accounts returns all the existing accounts */
    accounts(request?: QueryAccountsRequest): Promise<QueryAccountsResponse>;
    /** Account returns account details based on address. */
    account(request: QueryAccountRequest): Promise<QueryAccountResponse>;
    /** Params queries all parameters. */
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    accounts(request?: QueryAccountsRequest): Promise<QueryAccountsResponse>;
    account(request: QueryAccountRequest): Promise<QueryAccountResponse>;
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    accounts(request?: QueryAccountsRequest): Promise<QueryAccountsResponse>;
    account(request: QueryAccountRequest): Promise<QueryAccountResponse>;
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
};
