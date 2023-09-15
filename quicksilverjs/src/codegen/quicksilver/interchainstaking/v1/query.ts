import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../../cosmos/base/query/v1beta1/pagination";
import { Zone, ZoneAmino, ZoneSDKType, DelegatorIntent, DelegatorIntentAmino, DelegatorIntentSDKType, Delegation, DelegationAmino, DelegationSDKType, Receipt, ReceiptAmino, ReceiptSDKType, WithdrawalRecord, WithdrawalRecordAmino, WithdrawalRecordSDKType, UnbondingRecord, UnbondingRecordAmino, UnbondingRecordSDKType, RedelegationRecord, RedelegationRecordAmino, RedelegationRecordSDKType } from "./interchainstaking";
import { Long, isSet, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.interchainstaking.v1";
export interface QueryZonesInfoRequest {
  pagination: PageRequest;
}
export interface QueryZonesInfoRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryZonesInfoRequest";
  value: Uint8Array;
}
export interface QueryZonesInfoRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryZonesInfoRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryZonesInfoRequest";
  value: QueryZonesInfoRequestAmino;
}
export interface QueryZonesInfoRequestSDKType {
  pagination: PageRequestSDKType;
}
export interface QueryZonesInfoResponse {
  zones: Zone[];
  pagination: PageResponse;
}
export interface QueryZonesInfoResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryZonesInfoResponse";
  value: Uint8Array;
}
export interface QueryZonesInfoResponseAmino {
  zones: ZoneAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryZonesInfoResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryZonesInfoResponse";
  value: QueryZonesInfoResponseAmino;
}
export interface QueryZonesInfoResponseSDKType {
  zones: ZoneSDKType[];
  pagination: PageResponseSDKType;
}
/**
 * QueryDepositAccountForChainRequest is the request type for the
 * Query/InterchainAccountAddress RPC
 */
export interface QueryDepositAccountForChainRequest {
  chainId: string;
}
export interface QueryDepositAccountForChainRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest";
  value: Uint8Array;
}
/**
 * QueryDepositAccountForChainRequest is the request type for the
 * Query/InterchainAccountAddress RPC
 */
export interface QueryDepositAccountForChainRequestAmino {
  chain_id: string;
}
export interface QueryDepositAccountForChainRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest";
  value: QueryDepositAccountForChainRequestAmino;
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
export interface QueryDepositAccountForChainResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse";
  value: Uint8Array;
}
/**
 * QueryDepositAccountForChainResponse the response type for the
 * Query/InterchainAccountAddress RPC
 */
export interface QueryDepositAccountForChainResponseAmino {
  deposit_account_address: string;
}
export interface QueryDepositAccountForChainResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse";
  value: QueryDepositAccountForChainResponseAmino;
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
export interface QueryDelegatorIntentRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest";
  value: Uint8Array;
}
export interface QueryDelegatorIntentRequestAmino {
  chain_id: string;
  delegator_address: string;
}
export interface QueryDelegatorIntentRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest";
  value: QueryDelegatorIntentRequestAmino;
}
export interface QueryDelegatorIntentRequestSDKType {
  chain_id: string;
  delegator_address: string;
}
export interface QueryDelegatorIntentResponse {
  intent: DelegatorIntent;
}
export interface QueryDelegatorIntentResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse";
  value: Uint8Array;
}
export interface QueryDelegatorIntentResponseAmino {
  intent?: DelegatorIntentAmino;
}
export interface QueryDelegatorIntentResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse";
  value: QueryDelegatorIntentResponseAmino;
}
export interface QueryDelegatorIntentResponseSDKType {
  intent: DelegatorIntentSDKType;
}
export interface QueryDelegationsRequest {
  chainId: string;
  pagination: PageRequest;
}
export interface QueryDelegationsRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegationsRequest";
  value: Uint8Array;
}
export interface QueryDelegationsRequestAmino {
  chain_id: string;
  pagination?: PageRequestAmino;
}
export interface QueryDelegationsRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryDelegationsRequest";
  value: QueryDelegationsRequestAmino;
}
export interface QueryDelegationsRequestSDKType {
  chain_id: string;
  pagination: PageRequestSDKType;
}
export interface QueryDelegationsResponse {
  delegations: Delegation[];
  tvl: Long;
  pagination: PageResponse;
}
export interface QueryDelegationsResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegationsResponse";
  value: Uint8Array;
}
export interface QueryDelegationsResponseAmino {
  delegations: DelegationAmino[];
  tvl: string;
  pagination?: PageResponseAmino;
}
export interface QueryDelegationsResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryDelegationsResponse";
  value: QueryDelegationsResponseAmino;
}
export interface QueryDelegationsResponseSDKType {
  delegations: DelegationSDKType[];
  tvl: Long;
  pagination: PageResponseSDKType;
}
export interface QueryReceiptsRequest {
  chainId: string;
  pagination: PageRequest;
}
export interface QueryReceiptsRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryReceiptsRequest";
  value: Uint8Array;
}
export interface QueryReceiptsRequestAmino {
  chain_id: string;
  pagination?: PageRequestAmino;
}
export interface QueryReceiptsRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryReceiptsRequest";
  value: QueryReceiptsRequestAmino;
}
export interface QueryReceiptsRequestSDKType {
  chain_id: string;
  pagination: PageRequestSDKType;
}
export interface QueryReceiptsResponse {
  receipts: Receipt[];
  pagination: PageResponse;
}
export interface QueryReceiptsResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryReceiptsResponse";
  value: Uint8Array;
}
export interface QueryReceiptsResponseAmino {
  receipts: ReceiptAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryReceiptsResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryReceiptsResponse";
  value: QueryReceiptsResponseAmino;
}
export interface QueryReceiptsResponseSDKType {
  receipts: ReceiptSDKType[];
  pagination: PageResponseSDKType;
}
export interface QueryWithdrawalRecordsRequest {
  chainId: string;
  delegatorAddress: string;
  pagination: PageRequest;
}
export interface QueryWithdrawalRecordsRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsRequest";
  value: Uint8Array;
}
export interface QueryWithdrawalRecordsRequestAmino {
  chain_id: string;
  delegator_address: string;
  pagination?: PageRequestAmino;
}
export interface QueryWithdrawalRecordsRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsRequest";
  value: QueryWithdrawalRecordsRequestAmino;
}
export interface QueryWithdrawalRecordsRequestSDKType {
  chain_id: string;
  delegator_address: string;
  pagination: PageRequestSDKType;
}
export interface QueryWithdrawalRecordsResponse {
  withdrawals: WithdrawalRecord[];
  pagination: PageResponse;
}
export interface QueryWithdrawalRecordsResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsResponse";
  value: Uint8Array;
}
export interface QueryWithdrawalRecordsResponseAmino {
  withdrawals: WithdrawalRecordAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryWithdrawalRecordsResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsResponse";
  value: QueryWithdrawalRecordsResponseAmino;
}
export interface QueryWithdrawalRecordsResponseSDKType {
  withdrawals: WithdrawalRecordSDKType[];
  pagination: PageResponseSDKType;
}
export interface QueryUnbondingRecordsRequest {
  chainId: string;
  pagination: PageRequest;
}
export interface QueryUnbondingRecordsRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsRequest";
  value: Uint8Array;
}
export interface QueryUnbondingRecordsRequestAmino {
  chain_id: string;
  pagination?: PageRequestAmino;
}
export interface QueryUnbondingRecordsRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsRequest";
  value: QueryUnbondingRecordsRequestAmino;
}
export interface QueryUnbondingRecordsRequestSDKType {
  chain_id: string;
  pagination: PageRequestSDKType;
}
export interface QueryUnbondingRecordsResponse {
  Unbondings: UnbondingRecord[];
  pagination: PageResponse;
}
export interface QueryUnbondingRecordsResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsResponse";
  value: Uint8Array;
}
export interface QueryUnbondingRecordsResponseAmino {
  Unbondings: UnbondingRecordAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryUnbondingRecordsResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsResponse";
  value: QueryUnbondingRecordsResponseAmino;
}
export interface QueryUnbondingRecordsResponseSDKType {
  Unbondings: UnbondingRecordSDKType[];
  pagination: PageResponseSDKType;
}
export interface QueryRedelegationRecordsRequest {
  chainId: string;
  pagination: PageRequest;
}
export interface QueryRedelegationRecordsRequestProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsRequest";
  value: Uint8Array;
}
export interface QueryRedelegationRecordsRequestAmino {
  chain_id: string;
  pagination?: PageRequestAmino;
}
export interface QueryRedelegationRecordsRequestAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsRequest";
  value: QueryRedelegationRecordsRequestAmino;
}
export interface QueryRedelegationRecordsRequestSDKType {
  chain_id: string;
  pagination: PageRequestSDKType;
}
export interface QueryRedelegationRecordsResponse {
  Redelegations: RedelegationRecord[];
  pagination: PageResponse;
}
export interface QueryRedelegationRecordsResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsResponse";
  value: Uint8Array;
}
export interface QueryRedelegationRecordsResponseAmino {
  Redelegations: RedelegationRecordAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryRedelegationRecordsResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsResponse";
  value: QueryRedelegationRecordsResponseAmino;
}
export interface QueryRedelegationRecordsResponseSDKType {
  Redelegations: RedelegationRecordSDKType[];
  pagination: PageResponseSDKType;
}
function createBaseQueryZonesInfoRequest(): QueryZonesInfoRequest {
  return {
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryZonesInfoRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryZonesInfoRequest",
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
    const obj = createBaseQueryZonesInfoRequest();
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryZonesInfoRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryZonesInfoRequest>): QueryZonesInfoRequest {
    const message = createBaseQueryZonesInfoRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryZonesInfoRequestSDKType): QueryZonesInfoRequest {
    return {
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryZonesInfoRequest): QueryZonesInfoRequestSDKType {
    const obj: any = {};
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryZonesInfoRequestAmino): QueryZonesInfoRequest {
    return {
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryZonesInfoRequest): QueryZonesInfoRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryZonesInfoRequestAminoMsg): QueryZonesInfoRequest {
    return QueryZonesInfoRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZonesInfoRequestProtoMsg): QueryZonesInfoRequest {
    return QueryZonesInfoRequest.decode(message.value);
  },
  toProto(message: QueryZonesInfoRequest): Uint8Array {
    return QueryZonesInfoRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryZonesInfoRequest): QueryZonesInfoRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryZonesInfoRequest",
      value: QueryZonesInfoRequest.encode(message).finish()
    };
  }
};
function createBaseQueryZonesInfoResponse(): QueryZonesInfoResponse {
  return {
    zones: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryZonesInfoResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryZonesInfoResponse",
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
    const obj = createBaseQueryZonesInfoResponse();
    if (Array.isArray(object?.zones)) obj.zones = object.zones.map((e: any) => Zone.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
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
  fromPartial(object: DeepPartial<QueryZonesInfoResponse>): QueryZonesInfoResponse {
    const message = createBaseQueryZonesInfoResponse();
    message.zones = object.zones?.map(e => Zone.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryZonesInfoResponseSDKType): QueryZonesInfoResponse {
    return {
      zones: Array.isArray(object?.zones) ? object.zones.map((e: any) => Zone.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryZonesInfoResponse): QueryZonesInfoResponseSDKType {
    const obj: any = {};
    if (message.zones) {
      obj.zones = message.zones.map(e => e ? Zone.toSDK(e) : undefined);
    } else {
      obj.zones = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryZonesInfoResponseAmino): QueryZonesInfoResponse {
    return {
      zones: Array.isArray(object?.zones) ? object.zones.map((e: any) => Zone.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryZonesInfoResponse): QueryZonesInfoResponseAmino {
    const obj: any = {};
    if (message.zones) {
      obj.zones = message.zones.map(e => e ? Zone.toAmino(e) : undefined);
    } else {
      obj.zones = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryZonesInfoResponseAminoMsg): QueryZonesInfoResponse {
    return QueryZonesInfoResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZonesInfoResponseProtoMsg): QueryZonesInfoResponse {
    return QueryZonesInfoResponse.decode(message.value);
  },
  toProto(message: QueryZonesInfoResponse): Uint8Array {
    return QueryZonesInfoResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryZonesInfoResponse): QueryZonesInfoResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryZonesInfoResponse",
      value: QueryZonesInfoResponse.encode(message).finish()
    };
  }
};
function createBaseQueryDepositAccountForChainRequest(): QueryDepositAccountForChainRequest {
  return {
    chainId: ""
  };
}
export const QueryDepositAccountForChainRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest",
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
    const obj = createBaseQueryDepositAccountForChainRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    return obj;
  },
  toJSON(message: QueryDepositAccountForChainRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryDepositAccountForChainRequest>): QueryDepositAccountForChainRequest {
    const message = createBaseQueryDepositAccountForChainRequest();
    message.chainId = object.chainId ?? "";
    return message;
  },
  fromSDK(object: QueryDepositAccountForChainRequestSDKType): QueryDepositAccountForChainRequest {
    return {
      chainId: object?.chain_id
    };
  },
  toSDK(message: QueryDepositAccountForChainRequest): QueryDepositAccountForChainRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAmino(object: QueryDepositAccountForChainRequestAmino): QueryDepositAccountForChainRequest {
    return {
      chainId: object.chain_id
    };
  },
  toAmino(message: QueryDepositAccountForChainRequest): QueryDepositAccountForChainRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    return obj;
  },
  fromAminoMsg(object: QueryDepositAccountForChainRequestAminoMsg): QueryDepositAccountForChainRequest {
    return QueryDepositAccountForChainRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDepositAccountForChainRequestProtoMsg): QueryDepositAccountForChainRequest {
    return QueryDepositAccountForChainRequest.decode(message.value);
  },
  toProto(message: QueryDepositAccountForChainRequest): Uint8Array {
    return QueryDepositAccountForChainRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryDepositAccountForChainRequest): QueryDepositAccountForChainRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainRequest",
      value: QueryDepositAccountForChainRequest.encode(message).finish()
    };
  }
};
function createBaseQueryDepositAccountForChainResponse(): QueryDepositAccountForChainResponse {
  return {
    depositAccountAddress: ""
  };
}
export const QueryDepositAccountForChainResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse",
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
    const obj = createBaseQueryDepositAccountForChainResponse();
    if (isSet(object.depositAccountAddress)) obj.depositAccountAddress = String(object.depositAccountAddress);
    return obj;
  },
  toJSON(message: QueryDepositAccountForChainResponse): unknown {
    const obj: any = {};
    message.depositAccountAddress !== undefined && (obj.depositAccountAddress = message.depositAccountAddress);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryDepositAccountForChainResponse>): QueryDepositAccountForChainResponse {
    const message = createBaseQueryDepositAccountForChainResponse();
    message.depositAccountAddress = object.depositAccountAddress ?? "";
    return message;
  },
  fromSDK(object: QueryDepositAccountForChainResponseSDKType): QueryDepositAccountForChainResponse {
    return {
      depositAccountAddress: object?.deposit_account_address
    };
  },
  toSDK(message: QueryDepositAccountForChainResponse): QueryDepositAccountForChainResponseSDKType {
    const obj: any = {};
    obj.deposit_account_address = message.depositAccountAddress;
    return obj;
  },
  fromAmino(object: QueryDepositAccountForChainResponseAmino): QueryDepositAccountForChainResponse {
    return {
      depositAccountAddress: object.deposit_account_address
    };
  },
  toAmino(message: QueryDepositAccountForChainResponse): QueryDepositAccountForChainResponseAmino {
    const obj: any = {};
    obj.deposit_account_address = message.depositAccountAddress;
    return obj;
  },
  fromAminoMsg(object: QueryDepositAccountForChainResponseAminoMsg): QueryDepositAccountForChainResponse {
    return QueryDepositAccountForChainResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDepositAccountForChainResponseProtoMsg): QueryDepositAccountForChainResponse {
    return QueryDepositAccountForChainResponse.decode(message.value);
  },
  toProto(message: QueryDepositAccountForChainResponse): Uint8Array {
    return QueryDepositAccountForChainResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryDepositAccountForChainResponse): QueryDepositAccountForChainResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryDepositAccountForChainResponse",
      value: QueryDepositAccountForChainResponse.encode(message).finish()
    };
  }
};
function createBaseQueryDelegatorIntentRequest(): QueryDelegatorIntentRequest {
  return {
    chainId: "",
    delegatorAddress: ""
  };
}
export const QueryDelegatorIntentRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest",
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
    const obj = createBaseQueryDelegatorIntentRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.delegatorAddress)) obj.delegatorAddress = String(object.delegatorAddress);
    return obj;
  },
  toJSON(message: QueryDelegatorIntentRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryDelegatorIntentRequest>): QueryDelegatorIntentRequest {
    const message = createBaseQueryDelegatorIntentRequest();
    message.chainId = object.chainId ?? "";
    message.delegatorAddress = object.delegatorAddress ?? "";
    return message;
  },
  fromSDK(object: QueryDelegatorIntentRequestSDKType): QueryDelegatorIntentRequest {
    return {
      chainId: object?.chain_id,
      delegatorAddress: object?.delegator_address
    };
  },
  toSDK(message: QueryDelegatorIntentRequest): QueryDelegatorIntentRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.delegator_address = message.delegatorAddress;
    return obj;
  },
  fromAmino(object: QueryDelegatorIntentRequestAmino): QueryDelegatorIntentRequest {
    return {
      chainId: object.chain_id,
      delegatorAddress: object.delegator_address
    };
  },
  toAmino(message: QueryDelegatorIntentRequest): QueryDelegatorIntentRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.delegator_address = message.delegatorAddress;
    return obj;
  },
  fromAminoMsg(object: QueryDelegatorIntentRequestAminoMsg): QueryDelegatorIntentRequest {
    return QueryDelegatorIntentRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDelegatorIntentRequestProtoMsg): QueryDelegatorIntentRequest {
    return QueryDelegatorIntentRequest.decode(message.value);
  },
  toProto(message: QueryDelegatorIntentRequest): Uint8Array {
    return QueryDelegatorIntentRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryDelegatorIntentRequest): QueryDelegatorIntentRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentRequest",
      value: QueryDelegatorIntentRequest.encode(message).finish()
    };
  }
};
function createBaseQueryDelegatorIntentResponse(): QueryDelegatorIntentResponse {
  return {
    intent: DelegatorIntent.fromPartial({})
  };
}
export const QueryDelegatorIntentResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse",
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
    const obj = createBaseQueryDelegatorIntentResponse();
    if (isSet(object.intent)) obj.intent = DelegatorIntent.fromJSON(object.intent);
    return obj;
  },
  toJSON(message: QueryDelegatorIntentResponse): unknown {
    const obj: any = {};
    message.intent !== undefined && (obj.intent = message.intent ? DelegatorIntent.toJSON(message.intent) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryDelegatorIntentResponse>): QueryDelegatorIntentResponse {
    const message = createBaseQueryDelegatorIntentResponse();
    if (object.intent !== undefined && object.intent !== null) {
      message.intent = DelegatorIntent.fromPartial(object.intent);
    }
    return message;
  },
  fromSDK(object: QueryDelegatorIntentResponseSDKType): QueryDelegatorIntentResponse {
    return {
      intent: object.intent ? DelegatorIntent.fromSDK(object.intent) : undefined
    };
  },
  toSDK(message: QueryDelegatorIntentResponse): QueryDelegatorIntentResponseSDKType {
    const obj: any = {};
    message.intent !== undefined && (obj.intent = message.intent ? DelegatorIntent.toSDK(message.intent) : undefined);
    return obj;
  },
  fromAmino(object: QueryDelegatorIntentResponseAmino): QueryDelegatorIntentResponse {
    return {
      intent: object?.intent ? DelegatorIntent.fromAmino(object.intent) : undefined
    };
  },
  toAmino(message: QueryDelegatorIntentResponse): QueryDelegatorIntentResponseAmino {
    const obj: any = {};
    obj.intent = message.intent ? DelegatorIntent.toAmino(message.intent) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryDelegatorIntentResponseAminoMsg): QueryDelegatorIntentResponse {
    return QueryDelegatorIntentResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDelegatorIntentResponseProtoMsg): QueryDelegatorIntentResponse {
    return QueryDelegatorIntentResponse.decode(message.value);
  },
  toProto(message: QueryDelegatorIntentResponse): Uint8Array {
    return QueryDelegatorIntentResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryDelegatorIntentResponse): QueryDelegatorIntentResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegatorIntentResponse",
      value: QueryDelegatorIntentResponse.encode(message).finish()
    };
  }
};
function createBaseQueryDelegationsRequest(): QueryDelegationsRequest {
  return {
    chainId: "",
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryDelegationsRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegationsRequest",
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
    const obj = createBaseQueryDelegationsRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryDelegationsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryDelegationsRequest>): QueryDelegationsRequest {
    const message = createBaseQueryDelegationsRequest();
    message.chainId = object.chainId ?? "";
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryDelegationsRequestSDKType): QueryDelegationsRequest {
    return {
      chainId: object?.chain_id,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryDelegationsRequest): QueryDelegationsRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryDelegationsRequestAmino): QueryDelegationsRequest {
    return {
      chainId: object.chain_id,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryDelegationsRequest): QueryDelegationsRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryDelegationsRequestAminoMsg): QueryDelegationsRequest {
    return QueryDelegationsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDelegationsRequestProtoMsg): QueryDelegationsRequest {
    return QueryDelegationsRequest.decode(message.value);
  },
  toProto(message: QueryDelegationsRequest): Uint8Array {
    return QueryDelegationsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryDelegationsRequest): QueryDelegationsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegationsRequest",
      value: QueryDelegationsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryDelegationsResponse(): QueryDelegationsResponse {
  return {
    delegations: [],
    tvl: Long.ZERO,
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryDelegationsResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegationsResponse",
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
    const obj = createBaseQueryDelegationsResponse();
    if (Array.isArray(object?.delegations)) obj.delegations = object.delegations.map((e: any) => Delegation.fromJSON(e));
    if (isSet(object.tvl)) obj.tvl = Long.fromValue(object.tvl);
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
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
  fromPartial(object: DeepPartial<QueryDelegationsResponse>): QueryDelegationsResponse {
    const message = createBaseQueryDelegationsResponse();
    message.delegations = object.delegations?.map(e => Delegation.fromPartial(e)) || [];
    if (object.tvl !== undefined && object.tvl !== null) {
      message.tvl = Long.fromValue(object.tvl);
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryDelegationsResponseSDKType): QueryDelegationsResponse {
    return {
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => Delegation.fromSDK(e)) : [],
      tvl: object?.tvl,
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryDelegationsResponse): QueryDelegationsResponseSDKType {
    const obj: any = {};
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? Delegation.toSDK(e) : undefined);
    } else {
      obj.delegations = [];
    }
    obj.tvl = message.tvl;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryDelegationsResponseAmino): QueryDelegationsResponse {
    return {
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => Delegation.fromAmino(e)) : [],
      tvl: Long.fromString(object.tvl),
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryDelegationsResponse): QueryDelegationsResponseAmino {
    const obj: any = {};
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? Delegation.toAmino(e) : undefined);
    } else {
      obj.delegations = [];
    }
    obj.tvl = message.tvl ? message.tvl.toString() : undefined;
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryDelegationsResponseAminoMsg): QueryDelegationsResponse {
    return QueryDelegationsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDelegationsResponseProtoMsg): QueryDelegationsResponse {
    return QueryDelegationsResponse.decode(message.value);
  },
  toProto(message: QueryDelegationsResponse): Uint8Array {
    return QueryDelegationsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryDelegationsResponse): QueryDelegationsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryDelegationsResponse",
      value: QueryDelegationsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryReceiptsRequest(): QueryReceiptsRequest {
  return {
    chainId: "",
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryReceiptsRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryReceiptsRequest",
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
    const obj = createBaseQueryReceiptsRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryReceiptsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryReceiptsRequest>): QueryReceiptsRequest {
    const message = createBaseQueryReceiptsRequest();
    message.chainId = object.chainId ?? "";
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryReceiptsRequestSDKType): QueryReceiptsRequest {
    return {
      chainId: object?.chain_id,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryReceiptsRequest): QueryReceiptsRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryReceiptsRequestAmino): QueryReceiptsRequest {
    return {
      chainId: object.chain_id,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryReceiptsRequest): QueryReceiptsRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryReceiptsRequestAminoMsg): QueryReceiptsRequest {
    return QueryReceiptsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryReceiptsRequestProtoMsg): QueryReceiptsRequest {
    return QueryReceiptsRequest.decode(message.value);
  },
  toProto(message: QueryReceiptsRequest): Uint8Array {
    return QueryReceiptsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryReceiptsRequest): QueryReceiptsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryReceiptsRequest",
      value: QueryReceiptsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryReceiptsResponse(): QueryReceiptsResponse {
  return {
    receipts: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryReceiptsResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryReceiptsResponse",
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
    const obj = createBaseQueryReceiptsResponse();
    if (Array.isArray(object?.receipts)) obj.receipts = object.receipts.map((e: any) => Receipt.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
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
  fromPartial(object: DeepPartial<QueryReceiptsResponse>): QueryReceiptsResponse {
    const message = createBaseQueryReceiptsResponse();
    message.receipts = object.receipts?.map(e => Receipt.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryReceiptsResponseSDKType): QueryReceiptsResponse {
    return {
      receipts: Array.isArray(object?.receipts) ? object.receipts.map((e: any) => Receipt.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryReceiptsResponse): QueryReceiptsResponseSDKType {
    const obj: any = {};
    if (message.receipts) {
      obj.receipts = message.receipts.map(e => e ? Receipt.toSDK(e) : undefined);
    } else {
      obj.receipts = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryReceiptsResponseAmino): QueryReceiptsResponse {
    return {
      receipts: Array.isArray(object?.receipts) ? object.receipts.map((e: any) => Receipt.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryReceiptsResponse): QueryReceiptsResponseAmino {
    const obj: any = {};
    if (message.receipts) {
      obj.receipts = message.receipts.map(e => e ? Receipt.toAmino(e) : undefined);
    } else {
      obj.receipts = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryReceiptsResponseAminoMsg): QueryReceiptsResponse {
    return QueryReceiptsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryReceiptsResponseProtoMsg): QueryReceiptsResponse {
    return QueryReceiptsResponse.decode(message.value);
  },
  toProto(message: QueryReceiptsResponse): Uint8Array {
    return QueryReceiptsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryReceiptsResponse): QueryReceiptsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryReceiptsResponse",
      value: QueryReceiptsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryWithdrawalRecordsRequest(): QueryWithdrawalRecordsRequest {
  return {
    chainId: "",
    delegatorAddress: "",
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryWithdrawalRecordsRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsRequest",
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
    const obj = createBaseQueryWithdrawalRecordsRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.delegatorAddress)) obj.delegatorAddress = String(object.delegatorAddress);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryWithdrawalRecordsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryWithdrawalRecordsRequest>): QueryWithdrawalRecordsRequest {
    const message = createBaseQueryWithdrawalRecordsRequest();
    message.chainId = object.chainId ?? "";
    message.delegatorAddress = object.delegatorAddress ?? "";
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryWithdrawalRecordsRequestSDKType): QueryWithdrawalRecordsRequest {
    return {
      chainId: object?.chain_id,
      delegatorAddress: object?.delegator_address,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryWithdrawalRecordsRequest): QueryWithdrawalRecordsRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.delegator_address = message.delegatorAddress;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryWithdrawalRecordsRequestAmino): QueryWithdrawalRecordsRequest {
    return {
      chainId: object.chain_id,
      delegatorAddress: object.delegator_address,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryWithdrawalRecordsRequest): QueryWithdrawalRecordsRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.delegator_address = message.delegatorAddress;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryWithdrawalRecordsRequestAminoMsg): QueryWithdrawalRecordsRequest {
    return QueryWithdrawalRecordsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryWithdrawalRecordsRequestProtoMsg): QueryWithdrawalRecordsRequest {
    return QueryWithdrawalRecordsRequest.decode(message.value);
  },
  toProto(message: QueryWithdrawalRecordsRequest): Uint8Array {
    return QueryWithdrawalRecordsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryWithdrawalRecordsRequest): QueryWithdrawalRecordsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsRequest",
      value: QueryWithdrawalRecordsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryWithdrawalRecordsResponse(): QueryWithdrawalRecordsResponse {
  return {
    withdrawals: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryWithdrawalRecordsResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsResponse",
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
    const obj = createBaseQueryWithdrawalRecordsResponse();
    if (Array.isArray(object?.withdrawals)) obj.withdrawals = object.withdrawals.map((e: any) => WithdrawalRecord.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
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
  fromPartial(object: DeepPartial<QueryWithdrawalRecordsResponse>): QueryWithdrawalRecordsResponse {
    const message = createBaseQueryWithdrawalRecordsResponse();
    message.withdrawals = object.withdrawals?.map(e => WithdrawalRecord.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryWithdrawalRecordsResponseSDKType): QueryWithdrawalRecordsResponse {
    return {
      withdrawals: Array.isArray(object?.withdrawals) ? object.withdrawals.map((e: any) => WithdrawalRecord.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryWithdrawalRecordsResponse): QueryWithdrawalRecordsResponseSDKType {
    const obj: any = {};
    if (message.withdrawals) {
      obj.withdrawals = message.withdrawals.map(e => e ? WithdrawalRecord.toSDK(e) : undefined);
    } else {
      obj.withdrawals = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryWithdrawalRecordsResponseAmino): QueryWithdrawalRecordsResponse {
    return {
      withdrawals: Array.isArray(object?.withdrawals) ? object.withdrawals.map((e: any) => WithdrawalRecord.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryWithdrawalRecordsResponse): QueryWithdrawalRecordsResponseAmino {
    const obj: any = {};
    if (message.withdrawals) {
      obj.withdrawals = message.withdrawals.map(e => e ? WithdrawalRecord.toAmino(e) : undefined);
    } else {
      obj.withdrawals = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryWithdrawalRecordsResponseAminoMsg): QueryWithdrawalRecordsResponse {
    return QueryWithdrawalRecordsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryWithdrawalRecordsResponseProtoMsg): QueryWithdrawalRecordsResponse {
    return QueryWithdrawalRecordsResponse.decode(message.value);
  },
  toProto(message: QueryWithdrawalRecordsResponse): Uint8Array {
    return QueryWithdrawalRecordsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryWithdrawalRecordsResponse): QueryWithdrawalRecordsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryWithdrawalRecordsResponse",
      value: QueryWithdrawalRecordsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryUnbondingRecordsRequest(): QueryUnbondingRecordsRequest {
  return {
    chainId: "",
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryUnbondingRecordsRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsRequest",
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
    const obj = createBaseQueryUnbondingRecordsRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryUnbondingRecordsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryUnbondingRecordsRequest>): QueryUnbondingRecordsRequest {
    const message = createBaseQueryUnbondingRecordsRequest();
    message.chainId = object.chainId ?? "";
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryUnbondingRecordsRequestSDKType): QueryUnbondingRecordsRequest {
    return {
      chainId: object?.chain_id,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryUnbondingRecordsRequest): QueryUnbondingRecordsRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryUnbondingRecordsRequestAmino): QueryUnbondingRecordsRequest {
    return {
      chainId: object.chain_id,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryUnbondingRecordsRequest): QueryUnbondingRecordsRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryUnbondingRecordsRequestAminoMsg): QueryUnbondingRecordsRequest {
    return QueryUnbondingRecordsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryUnbondingRecordsRequestProtoMsg): QueryUnbondingRecordsRequest {
    return QueryUnbondingRecordsRequest.decode(message.value);
  },
  toProto(message: QueryUnbondingRecordsRequest): Uint8Array {
    return QueryUnbondingRecordsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryUnbondingRecordsRequest): QueryUnbondingRecordsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsRequest",
      value: QueryUnbondingRecordsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryUnbondingRecordsResponse(): QueryUnbondingRecordsResponse {
  return {
    Unbondings: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryUnbondingRecordsResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsResponse",
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
    const obj = createBaseQueryUnbondingRecordsResponse();
    if (Array.isArray(object?.Unbondings)) obj.Unbondings = object.Unbondings.map((e: any) => UnbondingRecord.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
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
  fromPartial(object: DeepPartial<QueryUnbondingRecordsResponse>): QueryUnbondingRecordsResponse {
    const message = createBaseQueryUnbondingRecordsResponse();
    message.Unbondings = object.Unbondings?.map(e => UnbondingRecord.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryUnbondingRecordsResponseSDKType): QueryUnbondingRecordsResponse {
    return {
      Unbondings: Array.isArray(object?.Unbondings) ? object.Unbondings.map((e: any) => UnbondingRecord.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryUnbondingRecordsResponse): QueryUnbondingRecordsResponseSDKType {
    const obj: any = {};
    if (message.Unbondings) {
      obj.Unbondings = message.Unbondings.map(e => e ? UnbondingRecord.toSDK(e) : undefined);
    } else {
      obj.Unbondings = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryUnbondingRecordsResponseAmino): QueryUnbondingRecordsResponse {
    return {
      Unbondings: Array.isArray(object?.Unbondings) ? object.Unbondings.map((e: any) => UnbondingRecord.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryUnbondingRecordsResponse): QueryUnbondingRecordsResponseAmino {
    const obj: any = {};
    if (message.Unbondings) {
      obj.Unbondings = message.Unbondings.map(e => e ? UnbondingRecord.toAmino(e) : undefined);
    } else {
      obj.Unbondings = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryUnbondingRecordsResponseAminoMsg): QueryUnbondingRecordsResponse {
    return QueryUnbondingRecordsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryUnbondingRecordsResponseProtoMsg): QueryUnbondingRecordsResponse {
    return QueryUnbondingRecordsResponse.decode(message.value);
  },
  toProto(message: QueryUnbondingRecordsResponse): Uint8Array {
    return QueryUnbondingRecordsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryUnbondingRecordsResponse): QueryUnbondingRecordsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryUnbondingRecordsResponse",
      value: QueryUnbondingRecordsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryRedelegationRecordsRequest(): QueryRedelegationRecordsRequest {
  return {
    chainId: "",
    pagination: PageRequest.fromPartial({})
  };
}
export const QueryRedelegationRecordsRequest = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsRequest",
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
    const obj = createBaseQueryRedelegationRecordsRequest();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.pagination)) obj.pagination = PageRequest.fromJSON(object.pagination);
    return obj;
  },
  toJSON(message: QueryRedelegationRecordsRequest): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<QueryRedelegationRecordsRequest>): QueryRedelegationRecordsRequest {
    const message = createBaseQueryRedelegationRecordsRequest();
    message.chainId = object.chainId ?? "";
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryRedelegationRecordsRequestSDKType): QueryRedelegationRecordsRequest {
    return {
      chainId: object?.chain_id,
      pagination: object.pagination ? PageRequest.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryRedelegationRecordsRequest): QueryRedelegationRecordsRequestSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageRequest.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryRedelegationRecordsRequestAmino): QueryRedelegationRecordsRequest {
    return {
      chainId: object.chain_id,
      pagination: object?.pagination ? PageRequest.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryRedelegationRecordsRequest): QueryRedelegationRecordsRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryRedelegationRecordsRequestAminoMsg): QueryRedelegationRecordsRequest {
    return QueryRedelegationRecordsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryRedelegationRecordsRequestProtoMsg): QueryRedelegationRecordsRequest {
    return QueryRedelegationRecordsRequest.decode(message.value);
  },
  toProto(message: QueryRedelegationRecordsRequest): Uint8Array {
    return QueryRedelegationRecordsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryRedelegationRecordsRequest): QueryRedelegationRecordsRequestProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsRequest",
      value: QueryRedelegationRecordsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryRedelegationRecordsResponse(): QueryRedelegationRecordsResponse {
  return {
    Redelegations: [],
    pagination: PageResponse.fromPartial({})
  };
}
export const QueryRedelegationRecordsResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsResponse",
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
    const obj = createBaseQueryRedelegationRecordsResponse();
    if (Array.isArray(object?.Redelegations)) obj.Redelegations = object.Redelegations.map((e: any) => RedelegationRecord.fromJSON(e));
    if (isSet(object.pagination)) obj.pagination = PageResponse.fromJSON(object.pagination);
    return obj;
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
  fromPartial(object: DeepPartial<QueryRedelegationRecordsResponse>): QueryRedelegationRecordsResponse {
    const message = createBaseQueryRedelegationRecordsResponse();
    message.Redelegations = object.Redelegations?.map(e => RedelegationRecord.fromPartial(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    }
    return message;
  },
  fromSDK(object: QueryRedelegationRecordsResponseSDKType): QueryRedelegationRecordsResponse {
    return {
      Redelegations: Array.isArray(object?.Redelegations) ? object.Redelegations.map((e: any) => RedelegationRecord.fromSDK(e)) : [],
      pagination: object.pagination ? PageResponse.fromSDK(object.pagination) : undefined
    };
  },
  toSDK(message: QueryRedelegationRecordsResponse): QueryRedelegationRecordsResponseSDKType {
    const obj: any = {};
    if (message.Redelegations) {
      obj.Redelegations = message.Redelegations.map(e => e ? RedelegationRecord.toSDK(e) : undefined);
    } else {
      obj.Redelegations = [];
    }
    message.pagination !== undefined && (obj.pagination = message.pagination ? PageResponse.toSDK(message.pagination) : undefined);
    return obj;
  },
  fromAmino(object: QueryRedelegationRecordsResponseAmino): QueryRedelegationRecordsResponse {
    return {
      Redelegations: Array.isArray(object?.Redelegations) ? object.Redelegations.map((e: any) => RedelegationRecord.fromAmino(e)) : [],
      pagination: object?.pagination ? PageResponse.fromAmino(object.pagination) : undefined
    };
  },
  toAmino(message: QueryRedelegationRecordsResponse): QueryRedelegationRecordsResponseAmino {
    const obj: any = {};
    if (message.Redelegations) {
      obj.Redelegations = message.Redelegations.map(e => e ? RedelegationRecord.toAmino(e) : undefined);
    } else {
      obj.Redelegations = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryRedelegationRecordsResponseAminoMsg): QueryRedelegationRecordsResponse {
    return QueryRedelegationRecordsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryRedelegationRecordsResponseProtoMsg): QueryRedelegationRecordsResponse {
    return QueryRedelegationRecordsResponse.decode(message.value);
  },
  toProto(message: QueryRedelegationRecordsResponse): Uint8Array {
    return QueryRedelegationRecordsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryRedelegationRecordsResponse): QueryRedelegationRecordsResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.QueryRedelegationRecordsResponse",
      value: QueryRedelegationRecordsResponse.encode(message).finish()
    };
  }
};