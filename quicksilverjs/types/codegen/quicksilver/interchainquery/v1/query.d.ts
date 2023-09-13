import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Query, QuerySDKType } from "./interchainquery";
import { Tx, TxSDKType } from "../../../cosmos/tx/v1beta1/tx";
import { TxResponse, TxResponseSDKType } from "../../../cosmos/base/abci/v1beta1/abci";
import { TxProof, TxProofSDKType } from "../../../tendermint/types/types";
import { Header, HeaderSDKType } from "../../../ibc/lightclients/tendermint/v1/tendermint";
import * as _m0 from "protobufjs/minimal";
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryRequestsRequest {
    pagination?: PageRequest;
    chainId: string;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryRequestsRequestSDKType {
    pagination?: PageRequestSDKType;
    chain_id: string;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryRequestsResponse {
    /** params defines the parameters of the module. */
    queries: Query[];
    pagination?: PageResponse;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryRequestsResponseSDKType {
    /** params defines the parameters of the module. */
    queries: QuerySDKType[];
    pagination?: PageResponseSDKType;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxWithProofResponse {
    /** tx is the queried transaction. */
    tx?: Tx;
    /** tx_response is the queried TxResponses. */
    txResponse?: TxResponse;
    /** proof is the tmproto.TxProof for the queried tx */
    proof?: TxProof;
    /** ibc-go header to validate txs */
    header?: Header;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxWithProofResponseSDKType {
    /** tx is the queried transaction. */
    tx?: TxSDKType;
    /** tx_response is the queried TxResponses. */
    tx_response?: TxResponseSDKType;
    /** proof is the tmproto.TxProof for the queried tx */
    proof?: TxProofSDKType;
    /** ibc-go header to validate txs */
    header?: HeaderSDKType;
}
export declare const QueryRequestsRequest: {
    encode(message: QueryRequestsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryRequestsRequest;
    fromJSON(object: any): QueryRequestsRequest;
    toJSON(message: QueryRequestsRequest): unknown;
    fromPartial(object: Partial<QueryRequestsRequest>): QueryRequestsRequest;
};
export declare const QueryRequestsResponse: {
    encode(message: QueryRequestsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryRequestsResponse;
    fromJSON(object: any): QueryRequestsResponse;
    toJSON(message: QueryRequestsResponse): unknown;
    fromPartial(object: Partial<QueryRequestsResponse>): QueryRequestsResponse;
};
export declare const GetTxWithProofResponse: {
    encode(message: GetTxWithProofResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetTxWithProofResponse;
    fromJSON(object: any): GetTxWithProofResponse;
    toJSON(message: GetTxWithProofResponse): unknown;
    fromPartial(object: Partial<GetTxWithProofResponse>): GetTxWithProofResponse;
};
