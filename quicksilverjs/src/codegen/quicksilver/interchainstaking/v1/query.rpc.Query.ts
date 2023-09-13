import { Rpc } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
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
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.zoneInfos = this.zoneInfos.bind(this);
    this.depositAccount = this.depositAccount.bind(this);
    this.delegatorIntent = this.delegatorIntent.bind(this);
    this.delegations = this.delegations.bind(this);
    this.receipts = this.receipts.bind(this);
    this.zoneWithdrawalRecords = this.zoneWithdrawalRecords.bind(this);
    this.withdrawalRecords = this.withdrawalRecords.bind(this);
    this.unbondingRecords = this.unbondingRecords.bind(this);
    this.redelegationRecords = this.redelegationRecords.bind(this);
  }

  zoneInfos(request: QueryZonesInfoRequest = {
    pagination: undefined
  }): Promise<QueryZonesInfoResponse> {
    const data = QueryZonesInfoRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "ZoneInfos", data);
    return promise.then(data => QueryZonesInfoResponse.decode(new _m0.Reader(data)));
  }

  depositAccount(request: QueryDepositAccountForChainRequest): Promise<QueryDepositAccountForChainResponse> {
    const data = QueryDepositAccountForChainRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "DepositAccount", data);
    return promise.then(data => QueryDepositAccountForChainResponse.decode(new _m0.Reader(data)));
  }

  delegatorIntent(request: QueryDelegatorIntentRequest): Promise<QueryDelegatorIntentResponse> {
    const data = QueryDelegatorIntentRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "DelegatorIntent", data);
    return promise.then(data => QueryDelegatorIntentResponse.decode(new _m0.Reader(data)));
  }

  delegations(request: QueryDelegationsRequest): Promise<QueryDelegationsResponse> {
    const data = QueryDelegationsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "Delegations", data);
    return promise.then(data => QueryDelegationsResponse.decode(new _m0.Reader(data)));
  }

  receipts(request: QueryReceiptsRequest): Promise<QueryReceiptsResponse> {
    const data = QueryReceiptsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "Receipts", data);
    return promise.then(data => QueryReceiptsResponse.decode(new _m0.Reader(data)));
  }

  zoneWithdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse> {
    const data = QueryWithdrawalRecordsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "ZoneWithdrawalRecords", data);
    return promise.then(data => QueryWithdrawalRecordsResponse.decode(new _m0.Reader(data)));
  }

  withdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse> {
    const data = QueryWithdrawalRecordsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "WithdrawalRecords", data);
    return promise.then(data => QueryWithdrawalRecordsResponse.decode(new _m0.Reader(data)));
  }

  unbondingRecords(request: QueryUnbondingRecordsRequest): Promise<QueryUnbondingRecordsResponse> {
    const data = QueryUnbondingRecordsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "UnbondingRecords", data);
    return promise.then(data => QueryUnbondingRecordsResponse.decode(new _m0.Reader(data)));
  }

  redelegationRecords(request: QueryRedelegationRecordsRequest): Promise<QueryRedelegationRecordsResponse> {
    const data = QueryRedelegationRecordsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.interchainstaking.v1.Query", "RedelegationRecords", data);
    return promise.then(data => QueryRedelegationRecordsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    zoneInfos(request?: QueryZonesInfoRequest): Promise<QueryZonesInfoResponse> {
      return queryService.zoneInfos(request);
    },

    depositAccount(request: QueryDepositAccountForChainRequest): Promise<QueryDepositAccountForChainResponse> {
      return queryService.depositAccount(request);
    },

    delegatorIntent(request: QueryDelegatorIntentRequest): Promise<QueryDelegatorIntentResponse> {
      return queryService.delegatorIntent(request);
    },

    delegations(request: QueryDelegationsRequest): Promise<QueryDelegationsResponse> {
      return queryService.delegations(request);
    },

    receipts(request: QueryReceiptsRequest): Promise<QueryReceiptsResponse> {
      return queryService.receipts(request);
    },

    zoneWithdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse> {
      return queryService.zoneWithdrawalRecords(request);
    },

    withdrawalRecords(request: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponse> {
      return queryService.withdrawalRecords(request);
    },

    unbondingRecords(request: QueryUnbondingRecordsRequest): Promise<QueryUnbondingRecordsResponse> {
      return queryService.unbondingRecords(request);
    },

    redelegationRecords(request: QueryRedelegationRecordsRequest): Promise<QueryRedelegationRecordsResponse> {
      return queryService.redelegationRecords(request);
    }

  };
};