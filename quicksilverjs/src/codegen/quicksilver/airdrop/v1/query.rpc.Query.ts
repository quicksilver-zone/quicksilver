import { Rpc } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
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
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.zoneDrop = this.zoneDrop.bind(this);
    this.accountBalance = this.accountBalance.bind(this);
    this.zoneDrops = this.zoneDrops.bind(this);
    this.claimRecord = this.claimRecord.bind(this);
    this.claimRecords = this.claimRecords.bind(this);
  }

  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.airdrop.v1.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  zoneDrop(request: QueryZoneDropRequest): Promise<QueryZoneDropResponse> {
    const data = QueryZoneDropRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.airdrop.v1.Query", "ZoneDrop", data);
    return promise.then(data => QueryZoneDropResponse.decode(new _m0.Reader(data)));
  }

  accountBalance(request: QueryAccountBalanceRequest): Promise<QueryAccountBalanceResponse> {
    const data = QueryAccountBalanceRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.airdrop.v1.Query", "AccountBalance", data);
    return promise.then(data => QueryAccountBalanceResponse.decode(new _m0.Reader(data)));
  }

  zoneDrops(request: QueryZoneDropsRequest): Promise<QueryZoneDropsResponse> {
    const data = QueryZoneDropsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.airdrop.v1.Query", "ZoneDrops", data);
    return promise.then(data => QueryZoneDropsResponse.decode(new _m0.Reader(data)));
  }

  claimRecord(request: QueryClaimRecordRequest): Promise<QueryClaimRecordResponse> {
    const data = QueryClaimRecordRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.airdrop.v1.Query", "ClaimRecord", data);
    return promise.then(data => QueryClaimRecordResponse.decode(new _m0.Reader(data)));
  }

  claimRecords(request: QueryClaimRecordsRequest): Promise<QueryClaimRecordsResponse> {
    const data = QueryClaimRecordsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.airdrop.v1.Query", "ClaimRecords", data);
    return promise.then(data => QueryClaimRecordsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },

    zoneDrop(request: QueryZoneDropRequest): Promise<QueryZoneDropResponse> {
      return queryService.zoneDrop(request);
    },

    accountBalance(request: QueryAccountBalanceRequest): Promise<QueryAccountBalanceResponse> {
      return queryService.accountBalance(request);
    },

    zoneDrops(request: QueryZoneDropsRequest): Promise<QueryZoneDropsResponse> {
      return queryService.zoneDrops(request);
    },

    claimRecord(request: QueryClaimRecordRequest): Promise<QueryClaimRecordResponse> {
      return queryService.claimRecord(request);
    },

    claimRecords(request: QueryClaimRecordsRequest): Promise<QueryClaimRecordsResponse> {
      return queryService.claimRecords(request);
    }

  };
};