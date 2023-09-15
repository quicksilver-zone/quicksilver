import { Status, StatusSDKType, ZoneDrop, ZoneDropSDKType, ClaimRecord, ClaimRecordSDKType } from "./airdrop";
import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsSDKType } from "./params";
import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import * as fm from "../../../grpc-gateway";
import { QueryParamsRequest, QueryParamsRequestSDKType, QueryParamsResponse, QueryParamsResponseSDKType, QueryZoneDropRequest, QueryZoneDropRequestSDKType, QueryZoneDropResponse, QueryZoneDropResponseSDKType, QueryAccountBalanceRequest, QueryAccountBalanceRequestSDKType, QueryAccountBalanceResponse, QueryAccountBalanceResponseSDKType, QueryZoneDropsRequest, QueryZoneDropsRequestSDKType, QueryZoneDropsResponse, QueryZoneDropsResponseSDKType, QueryClaimRecordRequest, QueryClaimRecordRequestSDKType, QueryClaimRecordResponse, QueryClaimRecordResponseSDKType, QueryClaimRecordsRequest, QueryClaimRecordsRequestSDKType, QueryClaimRecordsResponse, QueryClaimRecordsResponseSDKType } from "./query";
export class Query {
  /** Params returns the total set of airdrop parameters. */
  static params(request: QueryParamsRequest, initRequest?: fm.InitReq): Promise<QueryParamsResponse> {
    return fm.fetchReq(`/quicksilver/airdrop/v1/params?${fm.renderURLSearchParams({
      ...request
    }, [])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** ZoneDrop returns the details of the specified zone airdrop. */
  static zoneDrop(request: QueryZoneDropRequest, initRequest?: fm.InitReq): Promise<QueryZoneDropResponse> {
    return fm.fetchReq(`/quicksilver/airdrop/v1/zonedrop/${request["chain_id"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** AccountBalance returns the module account balance of the specified zone. */
  static accountBalance(request: QueryAccountBalanceRequest, initRequest?: fm.InitReq): Promise<QueryAccountBalanceResponse> {
    return fm.fetchReq(`/quicksilver/airdrop/v1/accountbalance/${request["chain_id"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** ZoneDrops returns all zone airdrops of the specified status. */
  static zoneDrops(request: QueryZoneDropsRequest, initRequest?: fm.InitReq): Promise<QueryZoneDropsResponse> {
    return fm.fetchReq(`/quicksilver/airdrop/v1/zonedrops/${request["status"]}?${fm.renderURLSearchParams({
      ...request
    }, ["status"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /**
   * ClaimRecord returns the claim record that corresponds to the given zone and
   * address.
   */
  static claimRecord(request: QueryClaimRecordRequest, initRequest?: fm.InitReq): Promise<QueryClaimRecordResponse> {
    return fm.fetchReq(`/quicksilver/airdrop/v1/claimrecord/${request["chain_id"]}/${request["address"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id", "address"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** ClaimRecords returns all the claim records of the given zone. */
  static claimRecords(request: QueryClaimRecordsRequest, initRequest?: fm.InitReq): Promise<QueryClaimRecordsResponse> {
    return fm.fetchReq(`/quicksilver/airdrop/v1/claimrecords/${request["chain_id"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
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
  /** Params returns the total set of airdrop parameters. */
  async params(req: QueryParamsRequest, headers?: HeadersInit): Promise<QueryParamsResponse> {
    return Query.params(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** ZoneDrop returns the details of the specified zone airdrop. */
  async zoneDrop(req: QueryZoneDropRequest, headers?: HeadersInit): Promise<QueryZoneDropResponse> {
    return Query.zoneDrop(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** AccountBalance returns the module account balance of the specified zone. */
  async accountBalance(req: QueryAccountBalanceRequest, headers?: HeadersInit): Promise<QueryAccountBalanceResponse> {
    return Query.accountBalance(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** ZoneDrops returns all zone airdrops of the specified status. */
  async zoneDrops(req: QueryZoneDropsRequest, headers?: HeadersInit): Promise<QueryZoneDropsResponse> {
    return Query.zoneDrops(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /**
   * ClaimRecord returns the claim record that corresponds to the given zone and
   * address.
   */
  async claimRecord(req: QueryClaimRecordRequest, headers?: HeadersInit): Promise<QueryClaimRecordResponse> {
    return Query.claimRecord(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** ClaimRecords returns all the claim records of the given zone. */
  async claimRecords(req: QueryClaimRecordsRequest, headers?: HeadersInit): Promise<QueryClaimRecordsResponse> {
    return Query.claimRecords(req, {
      headers,
      pathPrefix: this.url
    });
  }
}