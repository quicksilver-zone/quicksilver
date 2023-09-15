import { Params, ParamsSDKType } from "./participationrewards";
import { LCDClient } from "@cosmology/lcd";
import { QueryParamsRequest, QueryParamsRequestSDKType, QueryParamsResponse, QueryParamsResponseSDKType, QueryProtocolDataRequest, QueryProtocolDataRequestSDKType, QueryProtocolDataResponse, QueryProtocolDataResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;
  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.params = this.params.bind(this);
    this.protocolData = this.protocolData.bind(this);
  }
  /* Params returns the total set of participation rewards parameters. */
  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `quicksilver/participationrewards/v1/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* ProtocolData returns the requested protocol data. */
  async protocolData(params: QueryProtocolDataRequest): Promise<QueryProtocolDataResponseSDKType> {
    const endpoint = `quicksilver/participationrewards/v1/protocoldata/${params.type}/${params.key}`;
    return await this.req.get<QueryProtocolDataResponseSDKType>(endpoint);
  }
}