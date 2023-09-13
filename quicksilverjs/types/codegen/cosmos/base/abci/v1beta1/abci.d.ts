import { Any, AnySDKType } from "../../../../google/protobuf/any";
import { Event, EventSDKType } from "../../../../tendermint/abci/types";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../../helpers";
/**
 * TxResponse defines a structure containing relevant tx data and metadata. The
 * tags are stringified and the log is JSON decoded.
 */
export interface TxResponse {
    /** The block height */
    height: Long;
    /** The transaction hash. */
    txhash: string;
    /** Namespace for the Code */
    codespace: string;
    /** Response code. */
    code: number;
    /** Result bytes, if any. */
    data: string;
    /**
     * The output of the application's logger (raw string). May be
     * non-deterministic.
     */
    rawLog: string;
    /** The output of the application's logger (typed). May be non-deterministic. */
    logs: ABCIMessageLog[];
    /** Additional information. May be non-deterministic. */
    info: string;
    /** Amount of gas requested for transaction. */
    gasWanted: Long;
    /** Amount of gas consumed by transaction. */
    gasUsed: Long;
    /** The request transaction bytes. */
    tx?: Any;
    /**
     * Time of the previous block. For heights > 1, it's the weighted median of
     * the timestamps of the valid votes in the block.LastCommit. For height == 1,
     * it's genesis time.
     */
    timestamp: string;
}
/**
 * TxResponse defines a structure containing relevant tx data and metadata. The
 * tags are stringified and the log is JSON decoded.
 */
export interface TxResponseSDKType {
    /** The block height */
    height: Long;
    /** The transaction hash. */
    txhash: string;
    /** Namespace for the Code */
    codespace: string;
    /** Response code. */
    code: number;
    /** Result bytes, if any. */
    data: string;
    /**
     * The output of the application's logger (raw string). May be
     * non-deterministic.
     */
    raw_log: string;
    /** The output of the application's logger (typed). May be non-deterministic. */
    logs: ABCIMessageLogSDKType[];
    /** Additional information. May be non-deterministic. */
    info: string;
    /** Amount of gas requested for transaction. */
    gas_wanted: Long;
    /** Amount of gas consumed by transaction. */
    gas_used: Long;
    /** The request transaction bytes. */
    tx?: AnySDKType;
    /**
     * Time of the previous block. For heights > 1, it's the weighted median of
     * the timestamps of the valid votes in the block.LastCommit. For height == 1,
     * it's genesis time.
     */
    timestamp: string;
}
/** ABCIMessageLog defines a structure containing an indexed tx ABCI message log. */
export interface ABCIMessageLog {
    msgIndex: number;
    log: string;
    /**
     * Events contains a slice of Event objects that were emitted during some
     * execution.
     */
    events: StringEvent[];
}
/** ABCIMessageLog defines a structure containing an indexed tx ABCI message log. */
export interface ABCIMessageLogSDKType {
    msg_index: number;
    log: string;
    /**
     * Events contains a slice of Event objects that were emitted during some
     * execution.
     */
    events: StringEventSDKType[];
}
/**
 * StringEvent defines en Event object wrapper where all the attributes
 * contain key/value pairs that are strings instead of raw bytes.
 */
export interface StringEvent {
    type: string;
    attributes: Attribute[];
}
/**
 * StringEvent defines en Event object wrapper where all the attributes
 * contain key/value pairs that are strings instead of raw bytes.
 */
export interface StringEventSDKType {
    type: string;
    attributes: AttributeSDKType[];
}
/**
 * Attribute defines an attribute wrapper where the key and value are
 * strings instead of raw bytes.
 */
export interface Attribute {
    key: string;
    value: string;
}
/**
 * Attribute defines an attribute wrapper where the key and value are
 * strings instead of raw bytes.
 */
export interface AttributeSDKType {
    key: string;
    value: string;
}
/** GasInfo defines tx execution gas context. */
export interface GasInfo {
    /** GasWanted is the maximum units of work we allow this tx to perform. */
    gasWanted: Long;
    /** GasUsed is the amount of gas actually consumed. */
    gasUsed: Long;
}
/** GasInfo defines tx execution gas context. */
export interface GasInfoSDKType {
    /** GasWanted is the maximum units of work we allow this tx to perform. */
    gas_wanted: Long;
    /** GasUsed is the amount of gas actually consumed. */
    gas_used: Long;
}
/** Result is the union of ResponseFormat and ResponseCheckTx. */
export interface Result {
    /**
     * Data is any data returned from message or handler execution. It MUST be
     * length prefixed in order to separate data from multiple message executions.
     */
    data: Uint8Array;
    /** Log contains the log information from message or handler execution. */
    log: string;
    /**
     * Events contains a slice of Event objects that were emitted during message
     * or handler execution.
     */
    events: Event[];
}
/** Result is the union of ResponseFormat and ResponseCheckTx. */
export interface ResultSDKType {
    /**
     * Data is any data returned from message or handler execution. It MUST be
     * length prefixed in order to separate data from multiple message executions.
     */
    data: Uint8Array;
    /** Log contains the log information from message or handler execution. */
    log: string;
    /**
     * Events contains a slice of Event objects that were emitted during message
     * or handler execution.
     */
    events: EventSDKType[];
}
/**
 * SimulationResponse defines the response generated when a transaction is
 * successfully simulated.
 */
export interface SimulationResponse {
    gasInfo?: GasInfo;
    result?: Result;
}
/**
 * SimulationResponse defines the response generated when a transaction is
 * successfully simulated.
 */
export interface SimulationResponseSDKType {
    gas_info?: GasInfoSDKType;
    result?: ResultSDKType;
}
/**
 * MsgData defines the data returned in a Result object during message
 * execution.
 */
export interface MsgData {
    msgType: string;
    data: Uint8Array;
}
/**
 * MsgData defines the data returned in a Result object during message
 * execution.
 */
export interface MsgDataSDKType {
    msg_type: string;
    data: Uint8Array;
}
/**
 * TxMsgData defines a list of MsgData. A transaction will have a MsgData object
 * for each message.
 */
export interface TxMsgData {
    data: MsgData[];
}
/**
 * TxMsgData defines a list of MsgData. A transaction will have a MsgData object
 * for each message.
 */
export interface TxMsgDataSDKType {
    data: MsgDataSDKType[];
}
/** SearchTxsResult defines a structure for querying txs pageable */
export interface SearchTxsResult {
    /** Count of all txs */
    totalCount: Long;
    /** Count of txs in current page */
    count: Long;
    /** Index of current page, start from 1 */
    pageNumber: Long;
    /** Count of total pages */
    pageTotal: Long;
    /** Max count txs per page */
    limit: Long;
    /** List of txs in current page */
    txs: TxResponse[];
}
/** SearchTxsResult defines a structure for querying txs pageable */
export interface SearchTxsResultSDKType {
    /** Count of all txs */
    total_count: Long;
    /** Count of txs in current page */
    count: Long;
    /** Index of current page, start from 1 */
    page_number: Long;
    /** Count of total pages */
    page_total: Long;
    /** Max count txs per page */
    limit: Long;
    /** List of txs in current page */
    txs: TxResponseSDKType[];
}
export declare const TxResponse: {
    encode(message: TxResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxResponse;
    fromJSON(object: any): TxResponse;
    toJSON(message: TxResponse): unknown;
    fromPartial(object: Partial<TxResponse>): TxResponse;
};
export declare const ABCIMessageLog: {
    encode(message: ABCIMessageLog, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ABCIMessageLog;
    fromJSON(object: any): ABCIMessageLog;
    toJSON(message: ABCIMessageLog): unknown;
    fromPartial(object: Partial<ABCIMessageLog>): ABCIMessageLog;
};
export declare const StringEvent: {
    encode(message: StringEvent, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): StringEvent;
    fromJSON(object: any): StringEvent;
    toJSON(message: StringEvent): unknown;
    fromPartial(object: Partial<StringEvent>): StringEvent;
};
export declare const Attribute: {
    encode(message: Attribute, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Attribute;
    fromJSON(object: any): Attribute;
    toJSON(message: Attribute): unknown;
    fromPartial(object: Partial<Attribute>): Attribute;
};
export declare const GasInfo: {
    encode(message: GasInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GasInfo;
    fromJSON(object: any): GasInfo;
    toJSON(message: GasInfo): unknown;
    fromPartial(object: Partial<GasInfo>): GasInfo;
};
export declare const Result: {
    encode(message: Result, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Result;
    fromJSON(object: any): Result;
    toJSON(message: Result): unknown;
    fromPartial(object: Partial<Result>): Result;
};
export declare const SimulationResponse: {
    encode(message: SimulationResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SimulationResponse;
    fromJSON(object: any): SimulationResponse;
    toJSON(message: SimulationResponse): unknown;
    fromPartial(object: Partial<SimulationResponse>): SimulationResponse;
};
export declare const MsgData: {
    encode(message: MsgData, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgData;
    fromJSON(object: any): MsgData;
    toJSON(message: MsgData): unknown;
    fromPartial(object: Partial<MsgData>): MsgData;
};
export declare const TxMsgData: {
    encode(message: TxMsgData, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxMsgData;
    fromJSON(object: any): TxMsgData;
    toJSON(message: TxMsgData): unknown;
    fromPartial(object: Partial<TxMsgData>): TxMsgData;
};
export declare const SearchTxsResult: {
    encode(message: SearchTxsResult, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SearchTxsResult;
    fromJSON(object: any): SearchTxsResult;
    toJSON(message: SearchTxsResult): unknown;
    fromPartial(object: Partial<SearchTxsResult>): SearchTxsResult;
};
