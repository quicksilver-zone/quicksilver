import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Claim, ClaimSDKType } from "./claimsmanager";
import * as fm from "../../../grpc-gateway";
import { QueryClaimsRequest, QueryClaimsRequestSDKType, QueryClaimsResponse, QueryClaimsResponseSDKType } from "./query";
export class Query {
  /** Claims returns all zone claims from the current epoch. */
  static claims(request: QueryClaimsRequest, initRequest?: fm.InitReq): Promise<QueryClaimsResponse> {
    return fm.fetchReq(`/quicksilver/claimsmanager/v1/claims/${request["chain_id"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** LastEpochClaims returns all zone claims from the last epoch. */
  static lastEpochClaims(request: QueryClaimsRequest, initRequest?: fm.InitReq): Promise<QueryClaimsResponse> {
    return fm.fetchReq(`/quicksilver/claimsmanager/v1/previous_epoch_claims/${request["chain_id"]}?${fm.renderURLSearchParams({
      ...request
    }, ["chain_id"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /**
   * UserClaims returns all zone claims for a given address from the current
   * epoch.
   */
  static userClaims(request: QueryClaimsRequest, initRequest?: fm.InitReq): Promise<QueryClaimsResponse> {
    return fm.fetchReq(`/quicksilver/claimsmanager/v1/user/${request["address"]}/claims?${fm.renderURLSearchParams({
      ...request
    }, ["address"])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /**
   * UserLastEpochClaims returns all zone claims for a given address from the
   * last epoch.
   */
  static userLastEpochClaims(request: QueryClaimsRequest, initRequest?: fm.InitReq): Promise<QueryClaimsResponse> {
    return fm.fetchReq(`/quicksilver/claimsmanager/v1/user/${request["address"]}/previous_epoch_claims?${fm.renderURLSearchParams({
      ...request
    }, ["address"])}`, {
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
  /** Claims returns all zone claims from the current epoch. */
  async claims(req: QueryClaimsRequest, headers?: HeadersInit): Promise<QueryClaimsResponse> {
    return Query.claims(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** LastEpochClaims returns all zone claims from the last epoch. */
  async lastEpochClaims(req: QueryClaimsRequest, headers?: HeadersInit): Promise<QueryClaimsResponse> {
    return Query.lastEpochClaims(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /**
   * UserClaims returns all zone claims for a given address from the current
   * epoch.
   */
  async userClaims(req: QueryClaimsRequest, headers?: HeadersInit): Promise<QueryClaimsResponse> {
    return Query.userClaims(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /**
   * UserLastEpochClaims returns all zone claims for a given address from the
   * last epoch.
   */
  async userLastEpochClaims(req: QueryClaimsRequest, headers?: HeadersInit): Promise<QueryClaimsResponse> {
    return Query.userLastEpochClaims(req, {
      headers,
      pathPrefix: this.url
    });
  }
}