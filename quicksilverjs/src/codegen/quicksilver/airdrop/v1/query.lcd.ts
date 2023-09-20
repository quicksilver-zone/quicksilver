import { Status, StatusSDKType, ZoneDrop, ZoneDropSDKType, ClaimRecord, ClaimRecordSDKType } from "./airdrop";
import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsSDKType } from "./params";
import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import { setPaginationParams } from "../../../helpers";
import { LCDClient } from "@cosmology/lcd";
import { QueryParamsRequest, QueryParamsRequestSDKType, QueryParamsResponse, QueryParamsResponseSDKType, QueryZoneDropRequest, QueryZoneDropRequestSDKType, QueryZoneDropResponse, QueryZoneDropResponseSDKType, QueryAccountBalanceRequest, QueryAccountBalanceRequestSDKType, QueryAccountBalanceResponse, QueryAccountBalanceResponseSDKType, QueryZoneDropsRequest, QueryZoneDropsRequestSDKType, QueryZoneDropsResponse, QueryZoneDropsResponseSDKType, QueryClaimRecordRequest, QueryClaimRecordRequestSDKType, QueryClaimRecordResponse, QueryClaimRecordResponseSDKType, QueryClaimRecordsRequest, QueryClaimRecordsRequestSDKType, QueryClaimRecordsResponse, QueryClaimRecordsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;
  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.params = this.params.bind(this);
    this.zoneDrop = this.zoneDrop.bind(this);
    this.accountBalance = this.accountBalance.bind(this);
    this.zoneDrops = this.zoneDrops.bind(this);
    this.claimRecord = this.claimRecord.bind(this);
    this.claimRecords = this.claimRecords.bind(this);
  }
  /* Params returns the total set of airdrop parameters. */
  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `quicksilver/airdrop/v1/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* ZoneDrop returns the details of the specified zone airdrop. */
  async zoneDrop(params: QueryZoneDropRequest): Promise<QueryZoneDropResponseSDKType> {
    const endpoint = `quicksilver/airdrop/v1/zonedrop/${params.chainId}`;
    return await this.req.get<QueryZoneDropResponseSDKType>(endpoint);
  }
  /* AccountBalance returns the module account balance of the specified zone. */
  async accountBalance(params: QueryAccountBalanceRequest): Promise<QueryAccountBalanceResponseSDKType> {
    const endpoint = `quicksilver/airdrop/v1/accountbalance/${params.chainId}`;
    return await this.req.get<QueryAccountBalanceResponseSDKType>(endpoint);
  }
  /* ZoneDrops returns all zone airdrops of the specified status. */
  async zoneDrops(params: QueryZoneDropsRequest): Promise<QueryZoneDropsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/airdrop/v1/zonedrops/${params.status}`;
    return await this.req.get<QueryZoneDropsResponseSDKType>(endpoint, options);
  }
  /* ClaimRecord returns the claim record that corresponds to the given zone and
   address. */
  async claimRecord(params: QueryClaimRecordRequest): Promise<QueryClaimRecordResponseSDKType> {
    const endpoint = `quicksilver/airdrop/v1/claimrecord/${params.chainId}/${params.address}`;
    return await this.req.get<QueryClaimRecordResponseSDKType>(endpoint);
  }
  /* ClaimRecords returns all the claim records of the given zone. */
  async claimRecords(params: QueryClaimRecordsRequest): Promise<QueryClaimRecordsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/airdrop/v1/claimrecords/${params.chainId}`;
    return await this.req.get<QueryClaimRecordsResponseSDKType>(endpoint, options);
  }
}