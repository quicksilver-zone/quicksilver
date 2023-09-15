import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Zone, ZoneSDKType, DelegatorIntent, DelegatorIntentSDKType, Delegation, DelegationSDKType, Receipt, ReceiptSDKType, WithdrawalRecord, WithdrawalRecordSDKType, UnbondingRecord, UnbondingRecordSDKType, RedelegationRecord, RedelegationRecordSDKType } from "./interchainstaking";
import { setPaginationParams } from "../../../helpers";
import { LCDClient } from "@cosmology/lcd";
import { QueryZonesInfoRequest, QueryZonesInfoRequestSDKType, QueryZonesInfoResponse, QueryZonesInfoResponseSDKType, QueryDepositAccountForChainRequest, QueryDepositAccountForChainRequestSDKType, QueryDepositAccountForChainResponse, QueryDepositAccountForChainResponseSDKType, QueryDelegatorIntentRequest, QueryDelegatorIntentRequestSDKType, QueryDelegatorIntentResponse, QueryDelegatorIntentResponseSDKType, QueryDelegationsRequest, QueryDelegationsRequestSDKType, QueryDelegationsResponse, QueryDelegationsResponseSDKType, QueryReceiptsRequest, QueryReceiptsRequestSDKType, QueryReceiptsResponse, QueryReceiptsResponseSDKType, QueryWithdrawalRecordsRequest, QueryWithdrawalRecordsRequestSDKType, QueryWithdrawalRecordsResponse, QueryWithdrawalRecordsResponseSDKType, QueryUnbondingRecordsRequest, QueryUnbondingRecordsRequestSDKType, QueryUnbondingRecordsResponse, QueryUnbondingRecordsResponseSDKType, QueryRedelegationRecordsRequest, QueryRedelegationRecordsRequestSDKType, QueryRedelegationRecordsResponse, QueryRedelegationRecordsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;
  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
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
  /* ZoneInfos provides meta data on connected zones. */
  async zoneInfos(params: QueryZonesInfoRequest = {
    pagination: PageRequest.fromPartial({})
  }): Promise<QueryZonesInfoResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/interchainstaking/v1/zones`;
    return await this.req.get<QueryZonesInfoResponseSDKType>(endpoint, options);
  }
  /* DepositAccount provides data on the deposit address for a connected zone. */
  async depositAccount(params: QueryDepositAccountForChainRequest): Promise<QueryDepositAccountForChainResponseSDKType> {
    const endpoint = `quicksilver/interchainstaking/v1/zones/${params.chainId}/deposit_address`;
    return await this.req.get<QueryDepositAccountForChainResponseSDKType>(endpoint);
  }
  /* DelegatorIntent provides data on the intent of the delegator for the given
   zone. */
  async delegatorIntent(params: QueryDelegatorIntentRequest): Promise<QueryDelegatorIntentResponseSDKType> {
    const endpoint = `quicksilver/interchainstaking/v1/zones/${params.chainId}/delegator_intent/${params.delegatorAddress}`;
    return await this.req.get<QueryDelegatorIntentResponseSDKType>(endpoint);
  }
  /* Delegations provides data on the delegations for the given zone. */
  async delegations(params: QueryDelegationsRequest): Promise<QueryDelegationsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/interchainstaking/v1/zones/${params.chainId}/delegations`;
    return await this.req.get<QueryDelegationsResponseSDKType>(endpoint, options);
  }
  /* Delegations provides data on the delegations for the given zone. */
  async receipts(params: QueryReceiptsRequest): Promise<QueryReceiptsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/interchainstaking/v1/zones/${params.chainId}/receipts`;
    return await this.req.get<QueryReceiptsResponseSDKType>(endpoint, options);
  }
  /* WithdrawalRecords provides data on the active withdrawals. */
  async zoneWithdrawalRecords(params: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/interchainstaking/v1/zones/${params.chainId}/withdrawal_records/${params.delegatorAddress}`;
    return await this.req.get<QueryWithdrawalRecordsResponseSDKType>(endpoint, options);
  }
  /* WithdrawalRecords provides data on the active withdrawals. */
  async withdrawalRecords(params: QueryWithdrawalRecordsRequest): Promise<QueryWithdrawalRecordsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.chainId !== "undefined") {
      options.params.chain_id = params.chainId;
    }
    if (typeof params?.delegatorAddress !== "undefined") {
      options.params.delegator_address = params.delegatorAddress;
    }
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/interchainstaking/v1/withdrawal_records`;
    return await this.req.get<QueryWithdrawalRecordsResponseSDKType>(endpoint, options);
  }
  /* UnbondingRecords provides data on the active unbondings. */
  async unbondingRecords(params: QueryUnbondingRecordsRequest): Promise<QueryUnbondingRecordsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.chainId !== "undefined") {
      options.params.chain_id = params.chainId;
    }
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/interchainstaking/v1/unbonding_records`;
    return await this.req.get<QueryUnbondingRecordsResponseSDKType>(endpoint, options);
  }
  /* RedelegationRecords provides data on the active unbondings. */
  async redelegationRecords(params: QueryRedelegationRecordsRequest): Promise<QueryRedelegationRecordsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.chainId !== "undefined") {
      options.params.chain_id = params.chainId;
    }
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/interchainstaking/v1/redelegation_records`;
    return await this.req.get<QueryRedelegationRecordsResponseSDKType>(endpoint, options);
  }
}