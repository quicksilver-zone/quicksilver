import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryZonesInfoRequest, QueryZonesInfoResponse, QueryDepositAccountForChainRequest, QueryDepositAccountForChainResponse, QueryDelegatorIntentRequest, QueryDelegatorIntentResponse, QueryDelegationsRequest, QueryDelegationsResponse, QueryReceiptsRequest, QueryReceiptsResponse, QueryWithdrawalRecordsRequest, QueryWithdrawalRecordsResponse, QueryUnbondingRecordsRequest, QueryUnbondingRecordsResponse, QueryRedelegationRecordsRequest, QueryRedelegationRecordsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** ZoneInfos provides meta data on connected zones. */
    zoneInfos(request?: QueryZonesInfoRequest): Promise<QueryZonesInfoResponse>;
    /** DepositAccount provides data on the deposit address for a connected zone. */
    depositAccount(request: QueryDepositAccountForChainRequest): Promise<QueryDepositAccountForChainResponse>;
    /**
     * DelegatorIntent provides data on the intent of the delegator for the given
     * zone.
     */
    delegatorIntent(request: QueryDelegatorIntentRequest): Promise<QueryDelegatorIntentResponse>;
    /** Delegations provides data on the delegations for the given zone. */
    delegations(request: QueryDelegationsRequest): Promise<QueryDelegationsResponse>;
    /** Delegations provides data on the delegations for the given zone. */
    receipts(request: QueryReceiptsRequest): Promise<QueryReceiptsResponse>;
    /** WithdrawalRecords provides data on the active withdrawals. */
    zoneWithdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse>;
    /** WithdrawalRecords provides data on the active withdrawals. */
    withdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse>;
    /** UnbondingRecords provides data on the active unbondings. */
    unbondingRecords(request: QueryUnbondingRecordsRequest): Promise<QueryUnbondingRecordsResponse>;
    /** RedelegationRecords provides data on the active unbondings. */
    redelegationRecords(request: QueryRedelegationRecordsRequest): Promise<QueryRedelegationRecordsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    zoneInfos(request?: QueryZonesInfoRequest): Promise<QueryZonesInfoResponse>;
    depositAccount(request: QueryDepositAccountForChainRequest): Promise<QueryDepositAccountForChainResponse>;
    delegatorIntent(request: QueryDelegatorIntentRequest): Promise<QueryDelegatorIntentResponse>;
    delegations(request: QueryDelegationsRequest): Promise<QueryDelegationsResponse>;
    receipts(request: QueryReceiptsRequest): Promise<QueryReceiptsResponse>;
    zoneWithdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse>;
    withdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse>;
    unbondingRecords(request: QueryUnbondingRecordsRequest): Promise<QueryUnbondingRecordsResponse>;
    redelegationRecords(request: QueryRedelegationRecordsRequest): Promise<QueryRedelegationRecordsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    zoneInfos(request?: QueryZonesInfoRequest): Promise<QueryZonesInfoResponse>;
    depositAccount(request: QueryDepositAccountForChainRequest): Promise<QueryDepositAccountForChainResponse>;
    delegatorIntent(request: QueryDelegatorIntentRequest): Promise<QueryDelegatorIntentResponse>;
    delegations(request: QueryDelegationsRequest): Promise<QueryDelegationsResponse>;
    receipts(request: QueryReceiptsRequest): Promise<QueryReceiptsResponse>;
    zoneWithdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse>;
    withdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse>;
    unbondingRecords(request: QueryUnbondingRecordsRequest): Promise<QueryUnbondingRecordsResponse>;
    redelegationRecords(request: QueryRedelegationRecordsRequest): Promise<QueryRedelegationRecordsResponse>;
};
