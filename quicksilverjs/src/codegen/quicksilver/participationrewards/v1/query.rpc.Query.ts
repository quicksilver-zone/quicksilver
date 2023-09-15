import { Params, ParamsSDKType } from "./participationrewards";
import * as fm from "../../../grpc-gateway";
import { QueryParamsRequest, QueryParamsRequestSDKType, QueryParamsResponse, QueryParamsResponseSDKType, QueryProtocolDataRequest, QueryProtocolDataRequestSDKType, QueryProtocolDataResponse, QueryProtocolDataResponseSDKType } from "./query";
export class Query {
  /** Params returns the total set of participation rewards parameters. */
  static params(request: QueryParamsRequest, initRequest?: fm.InitReq): Promise<QueryParamsResponse> {
    return fm.fetchReq(`/quicksilver/participationrewards/v1/params?${fm.renderURLSearchParams({
      ...request
    }, [])}`, {
      ...initRequest,
      method: "GET"
    });
  }
  /** ProtocolData returns the requested protocol data. */
  static protocolData(request: QueryProtocolDataRequest, initRequest?: fm.InitReq): Promise<QueryProtocolDataResponse> {
    return fm.fetchReq(`/quicksilver/participationrewards/v1/protocoldata/${request["type"]}/${request["key"]}?${fm.renderURLSearchParams({
      ...request
    }, ["type", "key"])}`, {
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
  /** Params returns the total set of participation rewards parameters. */
  async params(req: QueryParamsRequest, headers?: HeadersInit): Promise<QueryParamsResponse> {
    return Query.params(req, {
      headers,
      pathPrefix: this.url
    });
  }
  /** ProtocolData returns the requested protocol data. */
  async protocolData(req: QueryProtocolDataRequest, headers?: HeadersInit): Promise<QueryProtocolDataResponse> {
    return Query.protocolData(req, {
      headers,
      pathPrefix: this.url
    });
  }
}