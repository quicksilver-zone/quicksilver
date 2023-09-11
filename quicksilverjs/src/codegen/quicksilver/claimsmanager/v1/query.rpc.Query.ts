import { Rpc } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryClaimsRequest, QueryClaimsResponse } from "./query";
/** Query provides defines the gRPC querier service. */

export interface Query {
  /** Params returns the total set of participation rewards parameters. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  claims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
  lastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
  userClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
  userLastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.claims = this.claims.bind(this);
    this.lastEpochClaims = this.lastEpochClaims.bind(this);
    this.userClaims = this.userClaims.bind(this);
    this.userLastEpochClaims = this.userLastEpochClaims.bind(this);
  }

  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.claimsmanager.v1.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  claims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
    const data = QueryClaimsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.claimsmanager.v1.Query", "Claims", data);
    return promise.then(data => QueryClaimsResponse.decode(new _m0.Reader(data)));
  }

  lastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
    const data = QueryClaimsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.claimsmanager.v1.Query", "LastEpochClaims", data);
    return promise.then(data => QueryClaimsResponse.decode(new _m0.Reader(data)));
  }

  userClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
    const data = QueryClaimsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.claimsmanager.v1.Query", "UserClaims", data);
    return promise.then(data => QueryClaimsResponse.decode(new _m0.Reader(data)));
  }

  userLastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
    const data = QueryClaimsRequest.encode(request).finish();
    const promise = this.rpc.request("quicksilver.claimsmanager.v1.Query", "UserLastEpochClaims", data);
    return promise.then(data => QueryClaimsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },

    claims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
      return queryService.claims(request);
    },

    lastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
      return queryService.lastEpochClaims(request);
    },

    userClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
      return queryService.userClaims(request);
    },

    userLastEpochClaims(request: QueryClaimsRequest): Promise<QueryClaimsResponse> {
      return queryService.userLastEpochClaims(request);
    }

  };
};