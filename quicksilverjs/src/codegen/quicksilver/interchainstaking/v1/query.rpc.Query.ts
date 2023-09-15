import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Zone, ZoneSDKType, DelegatorIntent, DelegatorIntentSDKType, Delegation, DelegationSDKType, Receipt, ReceiptSDKType, WithdrawalRecord, WithdrawalRecordSDKType, UnbondingRecord, UnbondingRecordSDKType, RedelegationRecord, RedelegationRecordSDKType } from "./interchainstaking";
import * as fm from "../../../grpc-gateway";
import { QueryZonesInfoRequest, QueryZonesInfoRequestSDKType, QueryZonesInfoResponse, QueryZonesInfoResponseSDKType, QueryDepositAccountForChainRequest, QueryDepositAccountForChainRequestSDKType, QueryDepositAccountForChainResponse, QueryDepositAccountForChainResponseSDKType, QueryDelegatorIntentRequest, QueryDelegatorIntentRequestSDKType, QueryDelegatorIntentResponse, QueryDelegatorIntentResponseSDKType, QueryDelegationsRequest, QueryDelegationsRequestSDKType, QueryDelegationsResponse, QueryDelegationsResponseSDKType, QueryReceiptsRequest, QueryReceiptsRequestSDKType, QueryReceiptsResponse, QueryReceiptsResponseSDKType, QueryWithdrawalRecordsRequest, QueryWithdrawalRecordsRequestSDKType, QueryWithdrawalRecordsResponse, QueryWithdrawalRecordsResponseSDKType, QueryUnbondingRecordsRequest, QueryUnbondingRecordsRequestSDKType, QueryUnbondingRecordsResponse, QueryUnbondingRecordsResponseSDKType, QueryRedelegationRecordsRequest, QueryRedelegationRecordsRequestSDKType, QueryRedelegationRecordsResponse, QueryRedelegationRecordsResponseSDKType } from "./query";
export class Query {
  /** ZoneInfos provides meta data on connected zones. */
  static zoneInfos(request: QueryZonesInfoRequest, initRequest?: fm.InitReq): Promise<QueryZonesInfoResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/zones?${fm.renderURLSearchParams({
      ...request
    }, [])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** DepositAccount provides data on the deposit address for a connected zone. */
  static depositAccount(request: QueryDepositAccountForChainRequest, initRequest?: fm.InitReq): Promise<QueryDepositAccountForChainResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/zones/${request["chain_id"]}/deposit_address?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /**
   * DelegatorIntent provides data on the intent of the delegator for the given
   * zone.
   */
  static delegatorIntent(request: QueryDelegatorIntentRequest, initRequest?: fm.InitReq): Promise<QueryDelegatorIntentResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/zones/${request["chain_id"]}/delegator_intent/${request["delegator_address"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id", "delegator_address"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** Delegations provides data on the delegations for the given zone. */
  static delegations(request: QueryDelegationsRequest, initRequest?: fm.InitReq): Promise<QueryDelegationsResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/zones/${request["chain_id"]}/delegations?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** Delegations provides data on the delegations for the given zone. */
  static receipts(request: QueryReceiptsRequest, initRequest?: fm.InitReq): Promise<QueryReceiptsResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/zones/${request["chain_id"]}/receipts?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** WithdrawalRecords provides data on the active withdrawals. */
  static zoneWithdrawalRecords(request: QueryWithdrawalRecordsRequest, initRequest?: fm.InitReq): Promise<QueryWithdrawalRecordsResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/zones/${request["chain_id"]}/withdrawal_records/${request["delegator_address"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id", "delegator_address"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** WithdrawalRecords provides data on the active withdrawals. */
  static withdrawalRecords(request: QueryWithdrawalRecordsRequest, initRequest?: fm.InitReq): Promise<QueryWithdrawalRecordsResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/withdrawal_records?${fm.renderURLSearchParams({
      ...request
    }, [])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** UnbondingRecords provides data on the active unbondings. */
  static unbondingRecords(request: QueryUnbondingRecordsRequest, initRequest?: fm.InitReq): Promise<QueryUnbondingRecordsResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/unbonding_records?${fm.renderURLSearchParams({
      ...request
    }, [])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** RedelegationRecords provides data on the active unbondings. */
  static redelegationRecords(request: QueryRedelegationRecordsRequest, initRequest?: fm.InitReq): Promise<QueryRedelegationRecordsResponse> {
    return fm.fetchReq(`/quicksilver/interchainstaking/v1/redelegation_records?${fm.renderURLSearchParams({
      ...request
    }, [])}`, {
      ...initRequest,
      method: "GET"
    });
  }
}
export class QueryClientImpl {
  private readonly url: string;
  constructor(url: string) {
    this.url = url;
  }
  /** ZoneInfos provides meta data on connected zones. */
  async zoneInfos(req: QueryZonesInfoRequest, headers?: HeadersInit): Promise<QueryZonesInfoResponse> {
    return Query.zoneInfos(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** DepositAccount provides data on the deposit address for a connected zone. */
  async depositAccount(req: QueryDepositAccountForChainRequest, headers?: HeadersInit): Promise<QueryDepositAccountForChainResponse> {
    return Query.depositAccount(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /**
   * DelegatorIntent provides data on the intent of the delegator for the given
   * zone.
   */
  async delegatorIntent(req: QueryDelegatorIntentRequest, headers?: HeadersInit): Promise<QueryDelegatorIntentResponse> {
    return Query.delegatorIntent(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** Delegations provides data on the delegations for the given zone. */
  async delegations(req: QueryDelegationsRequest, headers?: HeadersInit): Promise<QueryDelegationsResponse> {
    return Query.delegations(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** Delegations provides data on the delegations for the given zone. */
  async receipts(req: QueryReceiptsRequest, headers?: HeadersInit): Promise<QueryReceiptsResponse> {
    return Query.receipts(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** WithdrawalRecords provides data on the active withdrawals. */
  async zoneWithdrawalRecords(req: QueryWithdrawalRecordsRequest, headers?: HeadersInit): Promise<QueryWithdrawalRecordsResponse> {
    return Query.zoneWithdrawalRecords(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** WithdrawalRecords provides data on the active withdrawals. */
  async withdrawalRecords(req: QueryWithdrawalRecordsRequest, headers?: HeadersInit): Promise<QueryWithdrawalRecordsResponse> {
    return Query.withdrawalRecords(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** UnbondingRecords provides data on the active unbondings. */
  async unbondingRecords(req: QueryUnbondingRecordsRequest, headers?: HeadersInit): Promise<QueryUnbondingRecordsResponse> {
    return Query.unbondingRecords(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** RedelegationRecords provides data on the active unbondings. */
  async redelegationRecords(req: QueryRedelegationRecordsRequest, headers?: HeadersInit): Promise<QueryRedelegationRecordsResponse> {
    return Query.redelegationRecords(req, {
      headers,
      pathPrefix: this.url
    });
  }
}