import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryZoneDropRequest, QueryZoneDropResponse, QueryAccountBalanceRequest, QueryAccountBalanceResponse, QueryZoneDropsRequest, QueryZoneDropsResponse, QueryClaimRecordRequest, QueryClaimRecordResponse, QueryClaimRecordsRequest, QueryClaimRecordsResponse } from "./query";
/** Query provides defines the gRPC querier service. */
export interface Query {
    /** Params returns the total set of airdrop parameters. */
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    /** ZoneDrop returns the details of the specified zone airdrop. */
    zoneDrop(request: QueryZoneDropRequest): Promise<QueryZoneDropResponse>;
    /** AccountBalance returns the module account balance of the specified zone. */
    accountBalance(request: QueryAccountBalanceRequest): Promise<QueryAccountBalanceResponse>;
    /** ZoneDrops returns all zone airdrops of the specified status. */
    zoneDrops(request: QueryZoneDropsRequest): Promise<QueryZoneDropsResponse>;
    /**
     * ClaimRecord returns the claim record that corresponds to the given zone and
     * address.
     */
    claimRecord(request: QueryClaimRecordRequest): Promise<QueryClaimRecordResponse>;
    /** ClaimRecords returns all the claim records of the given zone. */
    claimRecords(request: QueryClaimRecordsRequest): Promise<QueryClaimRecordsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    zoneDrop(request: QueryZoneDropRequest): Promise<QueryZoneDropResponse>;
    accountBalance(request: QueryAccountBalanceRequest): Promise<QueryAccountBalanceResponse>;
    zoneDrops(request: QueryZoneDropsRequest): Promise<QueryZoneDropsResponse>;
    claimRecord(request: QueryClaimRecordRequest): Promise<QueryClaimRecordResponse>;
    claimRecords(request: QueryClaimRecordsRequest): Promise<QueryClaimRecordsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    zoneDrop(request: QueryZoneDropRequest): Promise<QueryZoneDropResponse>;
    accountBalance(request: QueryAccountBalanceRequest): Promise<QueryAccountBalanceResponse>;
    zoneDrops(request: QueryZoneDropsRequest): Promise<QueryZoneDropsResponse>;
    claimRecord(request: QueryClaimRecordRequest): Promise<QueryClaimRecordResponse>;
    claimRecords(request: QueryClaimRecordsRequest): Promise<QueryClaimRecordsResponse>;
};
