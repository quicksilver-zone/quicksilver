import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Claim, ClaimSDKType } from "./claimsmanager";
import { setPaginationParams } from "../../../helpers";
import { LCDClient } from "@cosmology/lcd";
import { QueryClaimsRequest, QueryClaimsRequestSDKType, QueryClaimsResponse, QueryClaimsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;
  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.claims = this.claims.bind(this);
    this.lastEpochClaims = this.lastEpochClaims.bind(this);
    this.userClaims = this.userClaims.bind(this);
    this.userLastEpochClaims = this.userLastEpochClaims.bind(this);
  }
  /* Claims returns all zone claims from the current epoch. */
  async claims(params: QueryClaimsRequest): Promise<QueryClaimsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.address !== "undefined") {
      options.params.address = params.address;
    }
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/claimsmanager/v1/claims/${params.chainId}`;
    return await this.req.get<QueryClaimsResponseSDKType>(endpoint, options);
  }
  /* LastEpochClaims returns all zone claims from the last epoch. */
  async lastEpochClaims(params: QueryClaimsRequest): Promise<QueryClaimsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.address !== "undefined") {
      options.params.address = params.address;
    }
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/claimsmanager/v1/previous_epoch_claims/${params.chainId}`;
    return await this.req.get<QueryClaimsResponseSDKType>(endpoint, options);
  }
  /* UserClaims returns all zone claims for a given address from the current
   epoch. */
  async userClaims(params: QueryClaimsRequest): Promise<QueryClaimsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.chainId !== "undefined") {
      options.params.chain_id = params.chainId;
    }
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/claimsmanager/v1/user/${params.address}/claims`;
    return await this.req.get<QueryClaimsResponseSDKType>(endpoint, options);
  }
  /* UserLastEpochClaims returns all zone claims for a given address from the
   last epoch. */
  async userLastEpochClaims(params: QueryClaimsRequest): Promise<QueryClaimsResponseSDKType> {
    const options: any = {
      params: {}
    };
    if (typeof params?.chainId !== "undefined") {
      options.params.chain_id = params.chainId;
    }
    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }
    const endpoint = `quicksilver/claimsmanager/v1/user/${params.address}/previous_epoch_claims`;
    return await this.req.get<QueryClaimsResponseSDKType>(endpoint, options);
  }
}