import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
export interface Query {
    id: string;
    connectionId: string;
    chainId: string;
    queryType: string;
    request: Uint8Array;
    /** change these to uint64 in v0.5.0 */
    period: string;
    lastHeight: string;
    callbackId: string;
    ttl: Long;
    lastEmission: string;
}
export interface QuerySDKType {
    id: string;
    connection_id: string;
    chain_id: string;
    query_type: string;
    request: Uint8Array;
    /** change these to uint64 in v0.5.0 */
    period: string;
    last_height: string;
    callback_id: string;
    ttl: Long;
    last_emission: string;
}
export interface DataPoint {
    id: string;
    /** change these to uint64 in v0.5.0 */
    remoteHeight: string;
    localHeight: string;
    value: Uint8Array;
}
export interface DataPointSDKType {
    id: string;
    /** change these to uint64 in v0.5.0 */
    remote_height: string;
    local_height: string;
    value: Uint8Array;
}
export declare const Query: {
    encode(message: Query, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Query;
    fromJSON(object: any): Query;
    toJSON(message: Query): unknown;
    fromPartial(object: Partial<Query>): Query;
};
export declare const DataPoint: {
    encode(message: DataPoint, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DataPoint;
    fromJSON(object: any): DataPoint;
    toJSON(message: DataPoint): unknown;
    fromPartial(object: Partial<DataPoint>): DataPoint;
};
