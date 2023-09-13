import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Zone, ZoneSDKType, DelegatorIntent, DelegatorIntentSDKType, Delegation, DelegationSDKType, Receipt, ReceiptSDKType, WithdrawalRecord, WithdrawalRecordSDKType, UnbondingRecord, UnbondingRecordSDKType, RedelegationRecord, RedelegationRecordSDKType } from "./interchainstaking";
import * as _m0 from "protobufjs/minimal";
import { isSet, Long } from "../../../helpers";
export interface QueryZonesInfoRequest {
  pagination?: PageRequest;
}
export interface QueryZonesInfoRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryZonesInfoResponse {
  zones: Zone[];
  pagination?: PageResponse;
}
export interface QueryZonesInfoResponseSDKType {
  zones: ZoneSDKType[];
  pagination?: PageResponseSDKType;
}
/**
 * QueryDepositAccountForChainRequest is the request type for the
 * Query/InterchainAccountAddress RPC
 */

export interface QueryDepositAccountForChainRequest {
  chainId: string;
}
/**
 * QueryDepositAccountForChainRequest is the request type for the
 * Query/InterchainAccountAddress RPC
 */

export interface QueryDepositAccountForChainRequestSDKType {
  chain_id: string;
}
/**
 * QueryDepositAccountForChainResponse the response type for the
 * Query/InterchainAccountAddress RPC
 */

export interface QueryDepositAccountForChainResponse {
  depositAccountAddress: string;
}
/**
 * QueryDepositAccountForChainResponse the response type for the
 * Query/InterchainAccountAddress RPC
 */

export interface QueryDepositAccountForChainResponseSDKType {
  deposit_account_address: string;
}
export interface QueryDelegatorIntentRequest {
  chainId: string;
  delegatorAddress: string;
}
export interface QueryDelegatorIntentRequestSDKType {
  chain_id: string;
  delegator_address: string;
}
export interface QueryDelegatorIntentResponse {
  intent?: DelegatorIntent;
}
export interface QueryDelegatorIntentResponseSDKType {
  intent?: DelegatorIntentSDKType;
}
export interface QueryDelegationsRequest {
  chainId: string;
  pagination?: PageRequest;
}
export interface QueryDelegationsRequestSDKType {
  chain_id: string;
  pagination?: PageRequestSDKType;
}
export interface QueryDelegationsResponse {
  delegations: Delegation[];
  tvl: Long;
  pagination?: PageResponse;
}
export interface QueryDelegationsResponseSDKType {
  delegations: DelegationSDKType[];
  tvl: Long;
  pagination?: PageResponseSDKType;
}
export interface QueryReceiptsRequest {
  chainId: string;
  pagination?: PageRequest;
}
export interface QueryReceiptsRequestSDKType {
  chain_id: string;
  pagination?: PageRequestSDKType;
}
export interface QueryReceiptsResponse {
  receipts: Receipt[];
  pagination?: PageResponse;
}
export interface QueryReceiptsResponseSDKType {
  receipts: ReceiptSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryWithdrawalRecordsRequest {
  chainId: string;
  delegatorAddress: string;
  pagination?: PageRequest;
}
export interface QueryWithdrawalRecordsRequestSDKType {
  chain_id: string;
  delegator_address: string;
  pagination?: PageRequestSDKType;
}
export interface QueryWithdrawalRecordsResponse {
  withdrawals: WithdrawalRecord[];
  pagination?: PageResponse;
}
export interface QueryWithdrawalRecordsResponseSDKType {
  withdrawals: WithdrawalRecordSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryUnbondingRecordsRequest {
  chainId: string;
  pagination?: PageRequest;
}
export interface QueryUnbondingRecordsRequestSDKType {
  chain_id: string;
  pagination?: PageRequestSDKType;
}
export interface QueryUnbondingRecordsResponse {
  Unbondings: UnbondingRecord[];
  pagination?: PageResponse;
}
export interface QueryUnbondingRecordsResponseSDKType {
  Unbondings: UnbondingRecordSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryRedelegationRecordsRequest {
  chainId: string;
  pagination?: PageRequest;
}
export interface QueryRedelegationRecordsRequestSDKType {
  chain_id: string;
  pagination?: PageRequestSDKType;
}
export interface QueryRedelegationRecordsResponse {
  Redelegations: RedelegationRecord[];
  pagination?: PageResponse;
}
export interface QueryRedelegationRecordsResponseSDKType {
  Redelegations: RedelegationRecordSDKType[];
  pagination?: PageResponseSDKType;
}

function createBaseQueryZonesInfoRequest(): QueryZonesInfoRequest {
  return {
    pagination: undefined
  };
}

export const QueryZonesInfoRequest = {
  encode(message: QueryZonesInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryZonesInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZonesInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryZonesInfoRequest {
    return {
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryZonesInfoRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryZonesInfoRequest>): QueryZonesInfoRequest {
    const message = createBaseQueryZonesInfoRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryZonesInfoResponse(): QueryZonesInfoResponse {
  return {
    zones: [],
    pagination: undefined
  };
}

export const QueryZonesInfoResponse = {
  encode(message: QueryZonesInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.zones) {
      Zone.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryZonesInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZonesInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.zones.push(Zone.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryZonesInfoResponse {
    return {
      zones: Array.isArray(object?.zones) ? object.zones.map((e: any) => Zone.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryZonesInfoResponse): unknown {
    const obj: any = {};

    if (message.zones) {
      obj.zones = message.zones.map(e => e ? Zone.toJSON(e) : undefined);
    } else {
      obj.zones = [];
    }

    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryZonesInfoResponse>): QueryZonesInfoResponse {
    const message = createBaseQueryZonesInfoResponse();
    message.zones = object.zones?.map(e => Zone.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryDepositAccountForChainRequest(): QueryDepositAccountForChainRequest {
  return {
    chainId: ""
  };
}

export const QueryDepositAccountForChainRequest = {
  encode(message: QueryDepositAccountForChainRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDepositAccountForChainRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDepositAccountForChainRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryDepositAccountForChainRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : ""
    };
  },

  toJSON(message: QueryDepositAccountForChainRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    return obj;
  },

  fromPartial(object: Partial<QueryDepositAccountForChainRequest>): QueryDepositAccountForChainRequest {
    const message = createBaseQueryDepositAccountForChainRequest();
    message.chainId = object.chainId ?? "";
    return message;
  }

};

function createBaseQueryDepositAccountForChainResponse(): QueryDepositAccountForChainResponse {
  return {
    depositAccountAddress: ""
  };
}

export const QueryDepositAccountForChainResponse = {
  encode(message: QueryDepositAccountForChainResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.depositAccountAddress !== "") {
      writer.uint32(10).string(message.depositAccountAddress);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDepositAccountForChainResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDepositAccountForChainResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.depositAccountAddress = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryDepositAccountForChainResponse {
    return {
      depositAccountAddress: isSet(object.depositAccountAddress) ? String(object.depositAccountAddress) : ""
    };
  },

  toJSON(message: QueryDepositAccountForChainResponse): unknown {
    const obj: any = {};
    message.depositAccountAddress !== undefined && (obj.depositAccountAddress = message.depositAccountAddress);
    return obj;
  },

  fromPartial(object: Partial<QueryDepositAccountForChainResponse>): QueryDepositAccountForChainResponse {
    const message = createBaseQueryDepositAccountForChainResponse();
    message.depositAccountAddress = object.depositAccountAddress ?? "";
    return message;
  }

};

function createBaseQueryDelegatorIntentRequest(): QueryDelegatorIntentRequest {
  return {
    chainId: "",
    delegatorAddress: ""
  };
}

export const QueryDelegatorIntentRequest = {
  encode(message: QueryDelegatorIntentRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    if (message.delegatorAddress !== "") {
      writer.uint32(18).string(message.delegatorAddress);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegatorIntentRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelegatorIntentRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        case 2:
          message.delegatorAddress = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryDelegatorIntentRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      delegatorAddress: isSet(object.delegatorAddress) ? String(object.delegatorAddress) : ""
    };
  },

  toJSON(message: QueryDelegatorIntentRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
    return obj;
  },

  fromPartial(object: Partial<QueryDelegatorIntentRequest>): QueryDelegatorIntentRequest {
    const message = createBaseQueryDelegatorIntentRequest();
    message.chainId = object.chainId ?? "";
    message.delegatorAddress = object.delegatorAddress ?? "";
    return message;
  }

};

function createBaseQueryDelegatorIntentResponse(): QueryDelegatorIntentResponse {
  return {
    intent: undefined
  };
}

export const QueryDelegatorIntentResponse = {
  encode(message: QueryDelegatorIntentResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.intent !== undefined) {
      DelegatorIntent.encode(message.intent, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegatorIntentResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelegatorIntentResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.intent = DelegatorIntent.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryDelegatorIntentResponse {
    return {
      intent: isSet(object.intent) ? DelegatorIntent.fromJSON(object.intent) : undefined
    };
  },

  toJSON(message: QueryDelegatorIntentResponse): unknown {
    const obj: any = {};
    message.intent !== undefined && (obj.intent = message.intent ? DelegatorIntent.toJSON(message.intent) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryDelegatorIntentResponse>): QueryDelegatorIntentResponse {
    const message = createBaseQueryDelegatorIntentResponse();
    message.intent = object.intent !== undefined && object.intent !== null ? DelegatorIntent.fromPartial(object.intent) : undefined;
    return message;
  }

};

function createBaseQueryDelegationsRequest(): QueryDelegationsRequest {
  return {
    chainId: "",
    pagination: undefined
  };
}

export const QueryDelegationsRequest = {
  encode(message: QueryDelegationsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegationsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelegationsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryDelegationsRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryDelegationsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryDelegationsRequest>): QueryDelegationsRequest {
    const message = createBaseQueryDelegationsRequest();
    message.chainId = object.chainId ?? "";
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryDelegationsResponse(): QueryDelegationsResponse {
  return {
    delegations: [],
    tvl: Long.ZERO,
    pagination: undefined
  };
}

export const QueryDelegationsResponse = {
  encode(message: QueryDelegationsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.delegations) {
      Delegation.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (!message.tvl.isZero()) {
      writer.uint32(16).int64(message.tvl);
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelegationsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelegationsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.delegations.push(Delegation.decode(reader, reader.uint32()));
          break;

        case 2:
          message.tvl = (reader.int64() as Long);
          break;

        case 3:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryDelegationsResponse {
    return {
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => Delegation.fromJSON(e)) : [],
      tvl: isSet(object.tvl) ? Long.fromValue(object.tvl) : Long.ZERO,
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryDelegationsResponse): unknown {
    const obj: any = {};

    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? Delegation.toJSON(e) : undefined);
    } else {
      obj.delegations = [];
    }

    message.tvl !== undefined && (obj.tvl = (message.tvl || Long.ZERO).toString());
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryDelegationsResponse>): QueryDelegationsResponse {
    const message = createBaseQueryDelegationsResponse();
    message.delegations = object.delegations?.map(e => Delegation.fromPartial(e)) || [];
    message.tvl = object.tvl !== undefined && object.tvl !== null ? Long.fromValue(object.tvl) : Long.ZERO;
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryReceiptsRequest(): QueryReceiptsRequest {
  return {
    chainId: "",
    pagination: undefined
  };
}

export const QueryReceiptsRequest = {
  encode(message: QueryReceiptsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryReceiptsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryReceiptsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryReceiptsRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryReceiptsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryReceiptsRequest>): QueryReceiptsRequest {
    const message = createBaseQueryReceiptsRequest();
    message.chainId = object.chainId ?? "";
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryReceiptsResponse(): QueryReceiptsResponse {
  return {
    receipts: [],
    pagination: undefined
  };
}

export const QueryReceiptsResponse = {
  encode(message: QueryReceiptsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.receipts) {
      Receipt.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryReceiptsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryReceiptsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.receipts.push(Receipt.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryReceiptsResponse {
    return {
      receipts: Array.isArray(object?.receipts) ? object.receipts.map((e: any) => Receipt.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryReceiptsResponse): unknown {
    const obj: any = {};

    if (message.receipts) {
      obj.receipts = message.receipts.map(e => e ? Receipt.toJSON(e) : undefined);
    } else {
      obj.receipts = [];
    }

    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryReceiptsResponse>): QueryReceiptsResponse {
    const message = createBaseQueryReceiptsResponse();
    message.receipts = object.receipts?.map(e => Receipt.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryWithdrawalRecordsRequest(): QueryWithdrawalRecordsRequest {
  return {
    chainId: "",
    delegatorAddress: "",
    pagination: undefined
  };
}

export const QueryWithdrawalRecordsRequest = {
  encode(message: QueryWithdrawalRecordsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    if (message.delegatorAddress !== "") {
      writer.uint32(18).string(message.delegatorAddress);
    }

    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryWithdrawalRecordsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryWithdrawalRecordsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        case 2:
          message.delegatorAddress = reader.string();
          break;

        case 3:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryWithdrawalRecordsRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      delegatorAddress: isSet(object.delegatorAddress) ? String(object.delegatorAddress) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryWithdrawalRecordsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryWithdrawalRecordsRequest>): QueryWithdrawalRecordsRequest {
    const message = createBaseQueryWithdrawalRecordsRequest();
    message.chainId = object.chainId ?? "";
    message.delegatorAddress = object.delegatorAddress ?? "";
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryWithdrawalRecordsResponse(): QueryWithdrawalRecordsResponse {
  return {
    withdrawals: [],
    pagination: undefined
  };
}

export const QueryWithdrawalRecordsResponse = {
  encode(message: QueryWithdrawalRecordsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.withdrawals) {
      WithdrawalRecord.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryWithdrawalRecordsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryWithdrawalRecordsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.withdrawals.push(WithdrawalRecord.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryWithdrawalRecordsResponse {
    return {
      withdrawals: Array.isArray(object?.withdrawals) ? object.withdrawals.map((e: any) => WithdrawalRecord.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryWithdrawalRecordsResponse): unknown {
    const obj: any = {};

    if (message.withdrawals) {
      obj.withdrawals = message.withdrawals.map(e => e ? WithdrawalRecord.toJSON(e) : undefined);
    } else {
      obj.withdrawals = [];
    }

    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryWithdrawalRecordsResponse>): QueryWithdrawalRecordsResponse {
    const message = createBaseQueryWithdrawalRecordsResponse();
    message.withdrawals = object.withdrawals?.map(e => WithdrawalRecord.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryUnbondingRecordsRequest(): QueryUnbondingRecordsRequest {
  return {
    chainId: "",
    pagination: undefined
  };
}

export const QueryUnbondingRecordsRequest = {
  encode(message: QueryUnbondingRecordsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUnbondingRecordsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUnbondingRecordsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        case 3:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryUnbondingRecordsRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryUnbondingRecordsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryUnbondingRecordsRequest>): QueryUnbondingRecordsRequest {
    const message = createBaseQueryUnbondingRecordsRequest();
    message.chainId = object.chainId ?? "";
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryUnbondingRecordsResponse(): QueryUnbondingRecordsResponse {
  return {
    Unbondings: [],
    pagination: undefined
  };
}

export const QueryUnbondingRecordsResponse = {
  encode(message: QueryUnbondingRecordsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Unbondings) {
      UnbondingRecord.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUnbondingRecordsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUnbondingRecordsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.Unbondings.push(UnbondingRecord.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryUnbondingRecordsResponse {
    return {
      Unbondings: Array.isArray(object?.Unbondings) ? object.Unbondings.map((e: any) => UnbondingRecord.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryUnbondingRecordsResponse): unknown {
    const obj: any = {};

    if (message.Unbondings) {
      obj.Unbondings = message.Unbondings.map(e => e ? UnbondingRecord.toJSON(e) : undefined);
    } else {
      obj.Unbondings = [];
    }

    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryUnbondingRecordsResponse>): QueryUnbondingRecordsResponse {
    const message = createBaseQueryUnbondingRecordsResponse();
    message.Unbondings = object.Unbondings?.map(e => UnbondingRecord.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryRedelegationRecordsRequest(): QueryRedelegationRecordsRequest {
  return {
    chainId: "",
    pagination: undefined
  };
}

export const QueryRedelegationRecordsRequest = {
  encode(message: QueryRedelegationRecordsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }

    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRedelegationRecordsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRedelegationRecordsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;

        case 3:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryRedelegationRecordsRequest {
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryRedelegationRecordsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryRedelegationRecordsRequest>): QueryRedelegationRecordsRequest {
    const message = createBaseQueryRedelegationRecordsRequest();
    message.chainId = object.chainId ?? "";
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryRedelegationRecordsResponse(): QueryRedelegationRecordsResponse {
  return {
    Redelegations: [],
    pagination: undefined
  };
}

export const QueryRedelegationRecordsResponse = {
  encode(message: QueryRedelegationRecordsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Redelegations) {
      RedelegationRecord.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRedelegationRecordsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRedelegationRecordsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.Redelegations.push(RedelegationRecord.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): QueryRedelegationRecordsResponse {
    return {
      Redelegations: Array.isArray(object?.Redelegations) ? object.Redelegations.map((e: any) => RedelegationRecord.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined
    };
  },

  toJSON(message: QueryRedelegationRecordsResponse): unknown {
    const obj: any = {};

    if (message.Redelegations) {
      obj.Redelegations = message.Redelegations.map(e => e ? RedelegationRecord.toJSON(e) : undefined);
    } else {
      obj.Redelegations = [];
    }

    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial(object: Partial<QueryRedelegationRecordsResponse>): QueryRedelegationRecordsResponse {
    const message = createBaseQueryRedelegationRecordsResponse();
    message.Redelegations = object.Redelegations?.map(e => RedelegationRecord.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};