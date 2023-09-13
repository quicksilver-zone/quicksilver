import { Delegation, DelegationSDKType, DelegatorIntent, DelegatorIntentSDKType, Zone, ZoneSDKType, Receipt, ReceiptSDKType, PortConnectionTuple, PortConnectionTupleSDKType, WithdrawalRecord, WithdrawalRecordSDKType } from "./interchainstaking";
import * as _m0 from "protobufjs/minimal";
import { Long, isSet } from "../../../helpers";
export interface Params {
  depositInterval: Long;
  validatorsetInterval: Long;
  commissionRate: string;
}
export interface ParamsSDKType {
  deposit_interval: Long;
  validatorset_interval: Long;
  commission_rate: string;
}
export interface DelegationsForZone {
  chainId: string;
  delegations: Delegation[];
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
export interface DelegatorIntentsForZoneSDKType {
  chain_id: string;
  delegation_intent: DelegatorIntentSDKType[];
  snapshot: boolean;
}
/** GenesisState defines the interchainstaking module's genesis state. */

export interface GenesisState {
  params?: Params;
  zones: Zone[];
  receipts: Receipt[];
  delegations: DelegationsForZone[];
  delegatorIntents: DelegatorIntentsForZone[];
  portConnections: PortConnectionTuple[];
  withdrawalRecords: WithdrawalRecord[];
}
/** GenesisState defines the interchainstaking module's genesis state. */

export interface GenesisStateSDKType {
  params?: ParamsSDKType;
  zones: ZoneSDKType[];
  receipts: ReceiptSDKType[];
  delegations: DelegationsForZoneSDKType[];
  delegator_intents: DelegatorIntentsForZoneSDKType[];
  port_connections: PortConnectionTupleSDKType[];
  withdrawal_records: WithdrawalRecordSDKType[];
}

function createBaseParams(): Params {
  return {
    depositInterval: Long.UZERO,
    validatorsetInterval: Long.UZERO,
    commissionRate: ""
  };
}

export const Params = {
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

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): Params {
    return {
      depositInterval: isSet(object.depositInterval) ? Long.fromValue(object.depositInterval) : Long.UZERO,
      validatorsetInterval: isSet(object.validatorsetInterval) ? Long.fromValue(object.validatorsetInterval) : Long.UZERO,
      commissionRate: isSet(object.commissionRate) ? String(object.commissionRate) : ""
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.depositInterval !== undefined && (obj.depositInterval = (message.depositInterval || Long.UZERO).toString());
    message.validatorsetInterval !== undefined && (obj.validatorsetInterval = (message.validatorsetInterval || Long.UZERO).toString());
    message.commissionRate !== undefined && (obj.commissionRate = message.commissionRate);
    return obj;
  },

  fromPartial(object: Partial<Params>): Params {
    const message = createBaseParams();
    message.depositInterval = object.depositInterval !== undefined && object.depositInterval !== null ? Long.fromValue(object.depositInterval) : Long.UZERO;
    message.validatorsetInterval = object.validatorsetInterval !== undefined && object.validatorsetInterval !== null ? Long.fromValue(object.validatorsetInterval) : Long.UZERO;
    message.commissionRate = object.commissionRate ?? "";
    return message;
  }

};

function createBaseDelegationsForZone(): DelegationsForZone {
  return {
    chainId: "",
    delegations: []
  };
}

export const DelegationsForZone = {
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
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => Delegation.fromJSON(e)) : []
    };
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

  fromPartial(object: Partial<DelegationsForZone>): DelegationsForZone {
    const message = createBaseDelegationsForZone();
    message.chainId = object.chainId ?? "";
    message.delegations = object.delegations?.map(e => Delegation.fromPartial(e)) || [];
    return message;
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
    return {
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      delegationIntent: Array.isArray(object?.delegationIntent) ? object.delegationIntent.map((e: any) => DelegatorIntent.fromJSON(e)) : [],
      snapshot: isSet(object.snapshot) ? Boolean(object.snapshot) : false
    };
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

  fromPartial(object: Partial<DelegatorIntentsForZone>): DelegatorIntentsForZone {
    const message = createBaseDelegatorIntentsForZone();
    message.chainId = object.chainId ?? "";
    message.delegationIntent = object.delegationIntent?.map(e => DelegatorIntent.fromPartial(e)) || [];
    message.snapshot = object.snapshot ?? false;
    return message;
  }

};

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    zones: [],
    receipts: [],
    delegations: [],
    delegatorIntents: [],
    portConnections: [],
    withdrawalRecords: []
  };
}

export const GenesisState = {
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

    for (const v of message.delegatorIntents) {
      DelegatorIntentsForZone.encode(v!, writer.uint32(42).fork()).ldelim();
    }

    for (const v of message.portConnections) {
      PortConnectionTuple.encode(v!, writer.uint32(50).fork()).ldelim();
    }

    for (const v of message.withdrawalRecords) {
      WithdrawalRecord.encode(v!, writer.uint32(58).fork()).ldelim();
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
          message.delegatorIntents.push(DelegatorIntentsForZone.decode(reader, reader.uint32()));
          break;

        case 6:
          message.portConnections.push(PortConnectionTuple.decode(reader, reader.uint32()));
          break;

        case 7:
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
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
      zones: Array.isArray(object?.zones) ? object.zones.map((e: any) => Zone.fromJSON(e)) : [],
      receipts: Array.isArray(object?.receipts) ? object.receipts.map((e: any) => Receipt.fromJSON(e)) : [],
      delegations: Array.isArray(object?.delegations) ? object.delegations.map((e: any) => DelegationsForZone.fromJSON(e)) : [],
      delegatorIntents: Array.isArray(object?.delegatorIntents) ? object.delegatorIntents.map((e: any) => DelegatorIntentsForZone.fromJSON(e)) : [],
      portConnections: Array.isArray(object?.portConnections) ? object.portConnections.map((e: any) => PortConnectionTuple.fromJSON(e)) : [],
      withdrawalRecords: Array.isArray(object?.withdrawalRecords) ? object.withdrawalRecords.map((e: any) => WithdrawalRecord.fromJSON(e)) : []
    };
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

  fromPartial(object: Partial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.zones = object.zones?.map(e => Zone.fromPartial(e)) || [];
    message.receipts = object.receipts?.map(e => Receipt.fromPartial(e)) || [];
    message.delegations = object.delegations?.map(e => DelegationsForZone.fromPartial(e)) || [];
    message.delegatorIntents = object.delegatorIntents?.map(e => DelegatorIntentsForZone.fromPartial(e)) || [];
    message.portConnections = object.portConnections?.map(e => PortConnectionTuple.fromPartial(e)) || [];
    message.withdrawalRecords = object.withdrawalRecords?.map(e => WithdrawalRecord.fromPartial(e)) || [];
    return message;
  }

};