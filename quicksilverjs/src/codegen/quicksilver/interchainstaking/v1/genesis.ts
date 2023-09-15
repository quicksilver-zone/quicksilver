import { Delegation, DelegationAmino, DelegationSDKType, DelegatorIntent, DelegatorIntentAmino, DelegatorIntentSDKType, Zone, ZoneAmino, ZoneSDKType, Receipt, ReceiptAmino, ReceiptSDKType, PortConnectionTuple, PortConnectionTupleAmino, PortConnectionTupleSDKType, WithdrawalRecord, WithdrawalRecordAmino, WithdrawalRecordSDKType } from "./interchainstaking";
import { Long, isSet, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.interchainstaking.v1";
export interface Params_v1 {
  depositInterval: Long;
  validatorsetInterval: Long;
  commissionRate: string;
}
export interface Params_v1ProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.Params_v1";
  value: Uint8Array;
}
export interface Params_v1Amino {
  deposit_interval: string;
  validatorset_interval: string;
  commission_rate: string;
}
export interface Params_v1AminoMsg {
  type: "/quicksilver.interchainstaking.v1.Params_v1";
  value: Params_v1Amino;
}
export interface Params_v1SDKType {
  deposit_interval: Long;
  validatorset_interval: Long;
  commission_rate: string;
}
export interface Params {
  depositInterval: Long;
  validatorsetInterval: Long;
  commissionRate: string;
  unbondingEnabled: boolean;
}
export interface ParamsProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.Params";
  value: Uint8Array;
}
export interface ParamsAmino {
  deposit_interval: string;
  validatorset_interval: string;
  commission_rate: string;
  unbonding_enabled: boolean;
}
export interface ParamsAminoMsg {
  type: "/quicksilver.interchainstaking.v1.Params";
  value: ParamsAmino;
}
export interface ParamsSDKType {
  deposit_interval: Long;
  validatorset_interval: Long;
  commission_rate: string;
  unbonding_enabled: boolean;
}
export interface DelegationsForZone {
  chainId: string;
  delegations: Delegation[];
}
export interface DelegationsForZoneProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.DelegationsForZone";
  value: Uint8Array;
}
export interface DelegationsForZoneAmino {
  chain_id: string;
  delegations: DelegationAmino[];
}
export interface DelegationsForZoneAminoMsg {
  type: "/quicksilver.interchainstaking.v1.DelegationsForZone";
  value: DelegationsForZoneAmino;
}
export interface DelegationsForZoneSDKType {
  chain_id: string;
  delegations: DelegationSDKType[];
}
export interface DelegatorIntentsForZone {
  chainId: string;
  delegationIntent: DelegatorIntent[];
  snapshot: boolean;
}
export interface DelegatorIntentsForZoneProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.DelegatorIntentsForZone";
  value: Uint8Array;
}
export interface DelegatorIntentsForZoneAmino {
  chain_id: string;
  delegation_intent: DelegatorIntentAmino[];
  snapshot: boolean;
}
export interface DelegatorIntentsForZoneAminoMsg {
  type: "/quicksilver.interchainstaking.v1.DelegatorIntentsForZone";
  value: DelegatorIntentsForZoneAmino;
}
export interface DelegatorIntentsForZoneSDKType {
  chain_id: string;
  delegation_intent: DelegatorIntentSDKType[];
  snapshot: boolean;
}
/** GenesisState defines the interchainstaking module's genesis state. */
export interface GenesisState {
  params: Params;
  zones: Zone[];
  receipts: Receipt[];
  delegations: DelegationsForZone[];
  performanceDelegations: DelegationsForZone[];
  delegatorIntents: DelegatorIntentsForZone[];
  portConnections: PortConnectionTuple[];
  withdrawalRecords: WithdrawalRecord[];
}
export interface GenesisStateProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the interchainstaking module's genesis state. */
export interface GenesisStateAmino {
  params?: ParamsAmino;
  zones: ZoneAmino[];
  receipts: ReceiptAmino[];
  delegations: DelegationsForZoneAmino[];
  performance_delegations: DelegationsForZoneAmino[];
  delegator_intents: DelegatorIntentsForZoneAmino[];
  port_connections: PortConnectionTupleAmino[];
  withdrawal_records: WithdrawalRecordAmino[];
}
export interface GenesisStateAminoMsg {
  type: "/quicksilver.interchainstaking.v1.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the interchainstaking module's genesis state. */
export interface GenesisStateSDKType {
  params: ParamsSDKType;
  zones: ZoneSDKType[];
  receipts: ReceiptSDKType[];
  delegations: DelegationsForZoneSDKType[];
  performance_delegations: DelegationsForZoneSDKType[];
  delegator_intents: DelegatorIntentsForZoneSDKType[];
  port_connections: PortConnectionTupleSDKType[];
  withdrawal_records: WithdrawalRecordSDKType[];
}
function createBaseParams_v1(): Params_v1 {
  return {
    depositInterval: Long.UZERO,
    validatorsetInterval: Long.UZERO,
    commissionRate: ""
  };
}
export const Params_v1 = {
  typeUrl: "/quicksilver.interchainstaking.v1.Params_v1",
  encode(message: Params_v1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.depositInterval.isZero()) {
      writer.uint32(8).uint64(message.depositInterval);
    }
    if (!message.validatorsetInterval.isZero()) {
      writer.uint32(16).uint64(message.validatorsetInterval);
    }
    if (message.commissionRate !== "") {
      writer.uint32(26).string(message.commissionRate);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Params_v1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams_v1();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositInterval = (reader.uint64() as Long);
          break;
        case 2:
          message.validatorsetInterval = (reader.uint64() as Long);
          break;
        case 3:
          message.commissionRate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Params_v1 {
    const obj = createBaseParams_v1();
    if (isSet(object.depositInterval)) obj.depositInterval = Long.fromValue(object.depositInterval);
    if (isSet(object.validatorsetInterval)) obj.validatorsetInterval = Long.fromValue(object.validatorsetInterval);
    if (isSet(object.commissionRate)) obj.commissionRate = String(object.commissionRate);
    return obj;
  },
  toJSON(message: Params_v1): unknown {
    const obj: any = {};
    message.depositInterval !== undefined && (obj.depositInterval = (message.depositInterval || Long.UZERO).toString());
    message.validatorsetInterval !== undefined && (obj.validatorsetInterval = (message.validatorsetInterval || Long.UZERO).toString());
    message.commissionRate !== undefined && (obj.commissionRate = message.commissionRate);
    return obj;
  },
  fromPartial(object: DeepPartial<Params_v1>): Params_v1 {
    const message = createBaseParams_v1();
    if (object.depositInterval !== undefined && object.depositInterval !== null) {
      message.depositInterval = Long.fromValue(object.depositInterval);
    }
    if (object.validatorsetInterval !== undefined && object.validatorsetInterval !== null) {
      message.validatorsetInterval = Long.fromValue(object.validatorsetInterval);
    }
    message.commissionRate = object.commissionRate ?? "";
    return message;
  },
  fromSDK(object: Params_v1SDKType): Params_v1 {
    return {
      depositInterval: object?.deposit_interval,
      validatorsetInterval: object?.validatorset_interval,
      commissionRate: object?.commission_rate
    };
  },
  toSDK(message: Params_v1): Params_v1SDKType {
    const obj: any = {};
    obj.deposit_interval = message.depositInterval;
    obj.validatorset_interval = message.validatorsetInterval;
    obj.commission_rate = message.commissionRate;
    return obj;
  },
  fromAmino(object: Params_v1Amino): Params_v1 {
    return {
      depositInterval: Long.fromString(object.deposit_interval),
      validatorsetInterval: Long.fromString(object.validatorset_interval),
      commissionRate: object.commission_rate
    };
  },
  toAmino(message: Params_v1): Params_v1Amino {
    const obj: any = {};
    obj.deposit_interval = message.depositInterval ? message.depositInterval.toString() : undefined;
    obj.validatorset_interval = message.validatorsetInterval ? message.validatorsetInterval.toString() : undefined;
    obj.commission_rate = message.commissionRate;
    return obj;
  },
  fromAminoMsg(object: Params_v1AminoMsg): Params_v1 {
    return Params_v1.fromAmino(object.value);
  },
  fromProtoMsg(message: Params_v1ProtoMsg): Params_v1 {
    return Params_v1.decode(message.value);
  },
  toProto(message: Params_v1): Uint8Array {
    return Params_v1.encode(message).finish();
  },
  toProtoMsg(message: Params_v1): Params_v1ProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.Params_v1",
      value: Params_v1.encode(message).finish()
    };
  }
};
function createBaseParams(): Params {
  return {
    depositInterval: Long.UZERO,
    validatorsetInterval: Long.UZERO,
    commissionRate: "",
    unbondingEnabled: false
  };
}
export const Params = {
  typeUrl: "/quicksilver.interchainstaking.v1.Params",
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.depositInterval.isZero()) {
      writer.uint32(8).uint64(message.depositInterval);
    }
    if (!message.validatorsetInterval.isZero()) {
      writer.uint32(16).uint64(message.validatorsetInterval);
    }
    if (message.commissionRate !== "") {
      writer.uint32(26).string(message.commissionRate);
    }
    if (message.unbondingEnabled === true) {
      writer.uint32(32).bool(message.unbondingEnabled);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.depositInterval = (reader.uint64() as Long);
          break;
        case 2:
          message.validatorsetInterval = (reader.uint64() as Long);
          break;
        case 3:
          message.commissionRate = reader.string();
          break;
        case 4:
          message.unbondingEnabled = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Params {
    const obj = createBaseParams();
    if (isSet(object.depositInterval)) obj.depositInterval = Long.fromValue(object.depositInterval);
    if (isSet(object.validatorsetInterval)) obj.validatorsetInterval = Long.fromValue(object.validatorsetInterval);
    if (isSet(object.commissionRate)) obj.commissionRate = String(object.commissionRate);
    if (isSet(object.unbondingEnabled)) obj.unbondingEnabled = Boolean(object.unbondingEnabled);
    return obj;
  },
  toJSON(message: Params): unknown {
    const obj: any = {};
    message.depositInterval !== undefined && (obj.depositInterval = (message.depositInterval || Long.UZERO).toString());
    message.validatorsetInterval !== undefined && (obj.validatorsetInterval = (message.validatorsetInterval || Long.UZERO).toString());
    message.commissionRate !== undefined && (obj.commissionRate = message.commissionRate);
    message.unbondingEnabled !== undefined && (obj.unbondingEnabled = message.unbondingEnabled);
    return obj;
  },
  fromPartial(object: DeepPartial<Params>): Params {
    const message = createBaseParams();
    if (object.depositInterval !== undefined && object.depositInterval !== null) {
      message.depositInterval = Long.fromValue(object.depositInterval);
    }
    if (object.validatorsetInterval !== undefined && object.validatorsetInterval !== null) {
      message.validatorsetInterval = Long.fromValue(object.validatorsetInterval);
    }
    message.commissionRate = object.commissionRate ?? "";
    message.unbondingEnabled = object.unbondingEnabled ?? false;
    return message;
  },
  fromSDK(object: ParamsSDKType): Params {
    return {
      depositInterval: object?.deposit_interval,
      validatorsetInterval: object?.validatorset_interval,
      commissionRate: object?.commission_rate,
      unbondingEnabled: object?.unbonding_enabled
    };
  },
  toSDK(message: Params): ParamsSDKType {
    const obj: any = {};
    obj.deposit_interval = message.depositInterval;
    obj.validatorset_interval = message.validatorsetInterval;
    obj.commission_rate = message.commissionRate;
    obj.unbonding_enabled = message.unbondingEnabled;
    return obj;
  },
  fromAmino(object: ParamsAmino): Params {
    return {
      depositInterval: Long.fromString(object.deposit_interval),
      validatorsetInterval: Long.fromString(object.validatorset_interval),
      commissionRate: object.commission_rate,
      unbondingEnabled: object.unbonding_enabled
    };
  },
  toAmino(message: Params): ParamsAmino {
    const obj: any = {};
    obj.deposit_interval = message.depositInterval ? message.depositInterval.toString() : undefined;
    obj.validatorset_interval = message.validatorsetInterval ? message.validatorsetInterval.toString() : undefined;
    obj.commission_rate = message.commissionRate;
    obj.unbonding_enabled = message.unbondingEnabled;
    return obj;
  },
  fromAminoMsg(object: ParamsAminoMsg): Params {
    return Params.fromAmino(object.value);
  },
  fromProtoMsg(message: ParamsProtoMsg): Params {
    return Params.decode(message.value);
  },
  toProto(message: Params): Uint8Array {
    return Params.encode(message).finish();
  },
  toProtoMsg(message: Params): ParamsProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.Params",
      value: Params.encode(message).finish()
    };
  }
};
function createBaseDelegationsForZone(): DelegationsForZone {
  return {
    chainId: "",
    delegations: []
  };
}
export const DelegationsForZone = {
  typeUrl: "/quicksilver.interchainstaking.v1.DelegationsForZone",
  encode(message: DelegationsForZone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    for (const v of message.delegations) {
      Delegation.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): DelegationsForZone {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDelegationsForZone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.delegations.push(Delegation.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): DelegationsForZone {
    const obj = createBaseDelegationsForZone();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (Array.isArray(object?.delegations)) obj.delegations = object.delegations.map((e: any) => Delegation.fromJSON(e));
    return obj;
  },
  toJSON(message: DelegationsForZone): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? Delegation.toJSON(e) : undefined);
    } else {
      obj.delegations = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<DelegationsForZone>): DelegationsForZone {
    const message = createBaseDelegationsForZone();
    message.chainId = object.chainId ?? "";
    message.delegations = object.delegations?.map(e => Delegation.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: DelegationsForZoneSDKType): DelegationsForZone {
    return {
      chainId: object?.chain_id,
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => Delegation.fromSDK(e)) : []
    };
  },
  toSDK(message: DelegationsForZone): DelegationsForZoneSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? Delegation.toSDK(e) : undefined);
    } else {
      obj.delegations = [];
    }
    return obj;
  },
  fromAmino(object: DelegationsForZoneAmino): DelegationsForZone {
    return {
      chainId: object.chain_id,
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => Delegation.fromAmino(e)) : []
    };
  },
  toAmino(message: DelegationsForZone): DelegationsForZoneAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? Delegation.toAmino(e) : undefined);
    } else {
      obj.delegations = [];
    }
    return obj;
  },
  fromAminoMsg(object: DelegationsForZoneAminoMsg): DelegationsForZone {
    return DelegationsForZone.fromAmino(object.value);
  },
  fromProtoMsg(message: DelegationsForZoneProtoMsg): DelegationsForZone {
    return DelegationsForZone.decode(message.value);
  },
  toProto(message: DelegationsForZone): Uint8Array {
    return DelegationsForZone.encode(message).finish();
  },
  toProtoMsg(message: DelegationsForZone): DelegationsForZoneProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.DelegationsForZone",
      value: DelegationsForZone.encode(message).finish()
    };
  }
};
function createBaseDelegatorIntentsForZone(): DelegatorIntentsForZone {
  return {
    chainId: "",
    delegationIntent: [],
    snapshot: false
  };
}
export const DelegatorIntentsForZone = {
  typeUrl: "/quicksilver.interchainstaking.v1.DelegatorIntentsForZone",
  encode(message: DelegatorIntentsForZone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    for (const v of message.delegationIntent) {
      DelegatorIntent.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.snapshot === true) {
      writer.uint32(24).bool(message.snapshot);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): DelegatorIntentsForZone {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDelegatorIntentsForZone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.delegationIntent.push(DelegatorIntent.decode(reader, reader.uint32()));
          break;
        case 3:
          message.snapshot = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): DelegatorIntentsForZone {
    const obj = createBaseDelegatorIntentsForZone();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (Array.isArray(object?.delegationIntent)) obj.delegationIntent = object.delegationIntent.map((e: any) => DelegatorIntent.fromJSON(e));
    if (isSet(object.snapshot)) obj.snapshot = Boolean(object.snapshot);
    return obj;
  },
  toJSON(message: DelegatorIntentsForZone): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    if (message.delegationIntent) {
      obj.delegationIntent = message.delegationIntent.map(e => e ? DelegatorIntent.toJSON(e) : undefined);
    } else {
      obj.delegationIntent = [];
    }
    message.snapshot !== undefined && (obj.snapshot = message.snapshot);
    return obj;
  },
  fromPartial(object: DeepPartial<DelegatorIntentsForZone>): DelegatorIntentsForZone {
    const message = createBaseDelegatorIntentsForZone();
    message.chainId = object.chainId ?? "";
    message.delegationIntent = object.delegationIntent?.map(e => DelegatorIntent.fromPartial(e)) || [];
    message.snapshot = object.snapshot ?? false;
    return message;
  },
  fromSDK(object: DelegatorIntentsForZoneSDKType): DelegatorIntentsForZone {
    return {
      chainId: object?.chain_id,
      delegationIntent: Array.isArray(object?.delegation_intent) ? object.delegation_intent.map((e: any) => DelegatorIntent.fromSDK(e)) : [],
      snapshot: object?.snapshot
    };
  },
  toSDK(message: DelegatorIntentsForZone): DelegatorIntentsForZoneSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    if (message.delegationIntent) {
      obj.delegation_intent = message.delegationIntent.map(e => e ? DelegatorIntent.toSDK(e) : undefined);
    } else {
      obj.delegation_intent = [];
    }
    obj.snapshot = message.snapshot;
    return obj;
  },
  fromAmino(object: DelegatorIntentsForZoneAmino): DelegatorIntentsForZone {
    return {
      chainId: object.chain_id,
      delegationIntent: Array.isArray(object?.delegation_intent) ? object.delegation_intent.map((e: any) => DelegatorIntent.fromAmino(e)) : [],
      snapshot: object.snapshot
    };
  },
  toAmino(message: DelegatorIntentsForZone): DelegatorIntentsForZoneAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    if (message.delegationIntent) {
      obj.delegation_intent = message.delegationIntent.map(e => e ? DelegatorIntent.toAmino(e) : undefined);
    } else {
      obj.delegation_intent = [];
    }
    obj.snapshot = message.snapshot;
    return obj;
  },
  fromAminoMsg(object: DelegatorIntentsForZoneAminoMsg): DelegatorIntentsForZone {
    return DelegatorIntentsForZone.fromAmino(object.value);
  },
  fromProtoMsg(message: DelegatorIntentsForZoneProtoMsg): DelegatorIntentsForZone {
    return DelegatorIntentsForZone.decode(message.value);
  },
  toProto(message: DelegatorIntentsForZone): Uint8Array {
    return DelegatorIntentsForZone.encode(message).finish();
  },
  toProtoMsg(message: DelegatorIntentsForZone): DelegatorIntentsForZoneProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.DelegatorIntentsForZone",
      value: DelegatorIntentsForZone.encode(message).finish()
    };
  }
};
function createBaseGenesisState(): GenesisState {
  return {
    params: Params.fromPartial({}),
    zones: [],
    receipts: [],
    delegations: [],
    performanceDelegations: [],
    delegatorIntents: [],
    portConnections: [],
    withdrawalRecords: []
  };
}
export const GenesisState = {
  typeUrl: "/quicksilver.interchainstaking.v1.GenesisState",
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.zones) {
      Zone.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.receipts) {
      Receipt.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.delegations) {
      DelegationsForZone.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.performanceDelegations) {
      DelegationsForZone.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.delegatorIntents) {
      DelegatorIntentsForZone.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.portConnections) {
      PortConnectionTuple.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.withdrawalRecords) {
      WithdrawalRecord.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.zones.push(Zone.decode(reader, reader.uint32()));
          break;
        case 3:
          message.receipts.push(Receipt.decode(reader, reader.uint32()));
          break;
        case 4:
          message.delegations.push(DelegationsForZone.decode(reader, reader.uint32()));
          break;
        case 5:
          message.performanceDelegations.push(DelegationsForZone.decode(reader, reader.uint32()));
          break;
        case 6:
          message.delegatorIntents.push(DelegatorIntentsForZone.decode(reader, reader.uint32()));
          break;
        case 7:
          message.portConnections.push(PortConnectionTuple.decode(reader, reader.uint32()));
          break;
        case 8:
          message.withdrawalRecords.push(WithdrawalRecord.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): GenesisState {
    const obj = createBaseGenesisState();
    if (isSet(object.params)) obj.params = Params.fromJSON(object.params);
    if (Array.isArray(object?.zones)) obj.zones = object.zones.map((e: any) => Zone.fromJSON(e));
    if (Array.isArray(object?.receipts)) obj.receipts = object.receipts.map((e: any) => Receipt.fromJSON(e));
    if (Array.isArray(object?.delegations)) obj.delegations = object.delegations.map((e: any) => DelegationsForZone.fromJSON(e));
    if (Array.isArray(object?.performanceDelegations)) obj.performanceDelegations = object.performanceDelegations.map((e: any) => DelegationsForZone.fromJSON(e));
    if (Array.isArray(object?.delegatorIntents)) obj.delegatorIntents = object.delegatorIntents.map((e: any) => DelegatorIntentsForZone.fromJSON(e));
    if (Array.isArray(object?.portConnections)) obj.portConnections = object.portConnections.map((e: any) => PortConnectionTuple.fromJSON(e));
    if (Array.isArray(object?.withdrawalRecords)) obj.withdrawalRecords = object.withdrawalRecords.map((e: any) => WithdrawalRecord.fromJSON(e));
    return obj;
  },
  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.zones) {
      obj.zones = message.zones.map(e => e ? Zone.toJSON(e) : undefined);
    } else {
      obj.zones = [];
    }
    if (message.receipts) {
      obj.receipts = message.receipts.map(e => e ? Receipt.toJSON(e) : undefined);
    } else {
      obj.receipts = [];
    }
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? DelegationsForZone.toJSON(e) : undefined);
    } else {
      obj.delegations = [];
    }
    if (message.performanceDelegations) {
      obj.performanceDelegations = message.performanceDelegations.map(e => e ? DelegationsForZone.toJSON(e) : undefined);
    } else {
      obj.performanceDelegations = [];
    }
    if (message.delegatorIntents) {
      obj.delegatorIntents = message.delegatorIntents.map(e => e ? DelegatorIntentsForZone.toJSON(e) : undefined);
    } else {
      obj.delegatorIntents = [];
    }
    if (message.portConnections) {
      obj.portConnections = message.portConnections.map(e => e ? PortConnectionTuple.toJSON(e) : undefined);
    } else {
      obj.portConnections = [];
    }
    if (message.withdrawalRecords) {
      obj.withdrawalRecords = message.withdrawalRecords.map(e => e ? WithdrawalRecord.toJSON(e) : undefined);
    } else {
      obj.withdrawalRecords = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    }
    message.zones = object.zones?.map(e => Zone.fromPartial(e)) || [];
    message.receipts = object.receipts?.map(e => Receipt.fromPartial(e)) || [];
    message.delegations = object.delegations?.map(e => DelegationsForZone.fromPartial(e)) || [];
    message.performanceDelegations = object.performanceDelegations?.map(e => DelegationsForZone.fromPartial(e)) || [];
    message.delegatorIntents = object.delegatorIntents?.map(e => DelegatorIntentsForZone.fromPartial(e)) || [];
    message.portConnections = object.portConnections?.map(e => PortConnectionTuple.fromPartial(e)) || [];
    message.withdrawalRecords = object.withdrawalRecords?.map(e => WithdrawalRecord.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: GenesisStateSDKType): GenesisState {
    return {
      params: object.params ? Params.fromSDK(object.params) : undefined,
      zones: Array.isArray(object?.zones) ? object.zones.map((e: any) => Zone.fromSDK(e)) : [],
      receipts: Array.isArray(object?.receipts) ? object.receipts.map((e: any) => Receipt.fromSDK(e)) : [],
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => DelegationsForZone.fromSDK(e)) : [],
      performanceDelegations: Array.isArray(object?.performance_delegations) ? object.performance_delegations.map((e: any) => DelegationsForZone.fromSDK(e)) : [],
      delegatorIntents: Array.isArray(object?.delegator_intents) ? object.delegator_intents.map((e: any) => DelegatorIntentsForZone.fromSDK(e)) : [],
      portConnections: Array.isArray(object?.port_connections) ? object.port_connections.map((e: any) => PortConnectionTuple.fromSDK(e)) : [],
      withdrawalRecords: Array.isArray(object?.withdrawal_records) ? object.withdrawal_records.map((e: any) => WithdrawalRecord.fromSDK(e)) : []
    };
  },
  toSDK(message: GenesisState): GenesisStateSDKType {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toSDK(message.params) : undefined);
    if (message.zones) {
      obj.zones = message.zones.map(e => e ? Zone.toSDK(e) : undefined);
    } else {
      obj.zones = [];
    }
    if (message.receipts) {
      obj.receipts = message.receipts.map(e => e ? Receipt.toSDK(e) : undefined);
    } else {
      obj.receipts = [];
    }
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? DelegationsForZone.toSDK(e) : undefined);
    } else {
      obj.delegations = [];
    }
    if (message.performanceDelegations) {
      obj.performance_delegations = message.performanceDelegations.map(e => e ? DelegationsForZone.toSDK(e) : undefined);
    } else {
      obj.performance_delegations = [];
    }
    if (message.delegatorIntents) {
      obj.delegator_intents = message.delegatorIntents.map(e => e ? DelegatorIntentsForZone.toSDK(e) : undefined);
    } else {
      obj.delegator_intents = [];
    }
    if (message.portConnections) {
      obj.port_connections = message.portConnections.map(e => e ? PortConnectionTuple.toSDK(e) : undefined);
    } else {
      obj.port_connections = [];
    }
    if (message.withdrawalRecords) {
      obj.withdrawal_records = message.withdrawalRecords.map(e => e ? WithdrawalRecord.toSDK(e) : undefined);
    } else {
      obj.withdrawal_records = [];
    }
    return obj;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    return {
      params: object?.params ? Params.fromAmino(object.params) : undefined,
      zones: Array.isArray(object?.zones) ? object.zones.map((e: any) => Zone.fromAmino(e)) : [],
      receipts: Array.isArray(object?.receipts) ? object.receipts.map((e: any) => Receipt.fromAmino(e)) : [],
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => DelegationsForZone.fromAmino(e)) : [],
      performanceDelegations: Array.isArray(object?.performance_delegations) ? object.performance_delegations.map((e: any) => DelegationsForZone.fromAmino(e)) : [],
      delegatorIntents: Array.isArray(object?.delegator_intents) ? object.delegator_intents.map((e: any) => DelegatorIntentsForZone.fromAmino(e)) : [],
      portConnections: Array.isArray(object?.port_connections) ? object.port_connections.map((e: any) => PortConnectionTuple.fromAmino(e)) : [],
      withdrawalRecords: Array.isArray(object?.withdrawal_records) ? object.withdrawal_records.map((e: any) => WithdrawalRecord.fromAmino(e)) : []
    };
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    if (message.zones) {
      obj.zones = message.zones.map(e => e ? Zone.toAmino(e) : undefined);
    } else {
      obj.zones = [];
    }
    if (message.receipts) {
      obj.receipts = message.receipts.map(e => e ? Receipt.toAmino(e) : undefined);
    } else {
      obj.receipts = [];
    }
    if (message.delegations) {
      obj.delegations = message.delegations.map(e => e ? DelegationsForZone.toAmino(e) : undefined);
    } else {
      obj.delegations = [];
    }
    if (message.performanceDelegations) {
      obj.performance_delegations = message.performanceDelegations.map(e => e ? DelegationsForZone.toAmino(e) : undefined);
    } else {
      obj.performance_delegations = [];
    }
    if (message.delegatorIntents) {
      obj.delegator_intents = message.delegatorIntents.map(e => e ? DelegatorIntentsForZone.toAmino(e) : undefined);
    } else {
      obj.delegator_intents = [];
    }
    if (message.portConnections) {
      obj.port_connections = message.portConnections.map(e => e ? PortConnectionTuple.toAmino(e) : undefined);
    } else {
      obj.port_connections = [];
    }
    if (message.withdrawalRecords) {
      obj.withdrawal_records = message.withdrawalRecords.map(e => e ? WithdrawalRecord.toAmino(e) : undefined);
    } else {
      obj.withdrawal_records = [];
    }
    return obj;
  },
  fromAminoMsg(object: GenesisStateAminoMsg): GenesisState {
    return GenesisState.fromAmino(object.value);
  },
  fromProtoMsg(message: GenesisStateProtoMsg): GenesisState {
    return GenesisState.decode(message.value);
  },
  toProto(message: GenesisState): Uint8Array {
    return GenesisState.encode(message).finish();
  },
  toProtoMsg(message: GenesisState): GenesisStateProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};