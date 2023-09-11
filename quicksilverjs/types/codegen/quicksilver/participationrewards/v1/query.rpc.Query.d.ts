import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryProtocolDataRequest, QueryProtocolDataResponse } from "./query";
/** Query provides defines the gRPC querier service. */
export interface Query {
    /** Params returns the total set of participation rewards parameters. */
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    protocolData(request: QueryProtocolDataRequest): Promise<QueryProtocolDataResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    protocolData(request: QueryProtocolDataRequest): Promise<QueryProtocolDataResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    protocolData(request: QueryProtocolDataRequest): Promise<QueryProtocolDataResponse>;
};
