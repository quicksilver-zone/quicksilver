import * as _m0 from "protobufjs/minimal";
import { isSet } from "../../../helpers";
export interface RegisterZoneProposal {
  title: string;
  description: string;
  connectionId: string;
  baseDenom: string;
  localDenom: string;
  accountPrefix: string;
  multiSend: boolean;
  liquidityModule: boolean;
}
export interface RegisterZoneProposalSDKType {
  title: string;
  description: string;
  connection_id: string;
  base_denom: string;
  local_denom: string;
  account_prefix: string;
  multi_send: boolean;
  liquidity_module: boolean;
}
export interface RegisterZoneProposalWithDeposit {
  title: string;
  description: string;
  connectionId: string;
  baseDenom: string;
  localDenom: string;
  accountPrefix: string;
  multiSend: boolean;
  liquidityModule: boolean;
  deposit: string;
}
export interface RegisterZoneProposalWithDepositSDKType {
  title: string;
  description: string;
  connection_id: string;
  base_denom: string;
  local_denom: string;
  account_prefix: string;
  multi_send: boolean;
  liquidity_module: boolean;
  deposit: string;
}
export interface UpdateZoneProposal {
  title: string;
  description: string;
  chainId: string;
  changes: UpdateZoneValue[];
}
export interface UpdateZoneProposalSDKType {
  title: string;
  description: string;
  chain_id: string;
  changes: UpdateZoneValueSDKType[];
}
export interface UpdateZoneProposalWithDeposit {
  title: string;
  description: string;
  chainId: string;
  changes: UpdateZoneValue[];
  deposit: string;
}
export interface UpdateZoneProposalWithDepositSDKType {
  title: string;
  description: string;
  chain_id: string;
  changes: UpdateZoneValueSDKType[];
  deposit: string;
}
/**
 * ParamChange defines an individual parameter change, for use in
 * ParameterChangeProposal.
 */

export interface UpdateZoneValue {
  key: string;
  value: string;
}
/**
 * ParamChange defines an individual parameter change, for use in
 * ParameterChangeProposal.
 */

export interface UpdateZoneValueSDKType {
  key: string;
  value: string;
}

function createBaseRegisterZoneProposal(): RegisterZoneProposal {
  return {
    title: "",
    description: "",
    connectionId: "",
    baseDenom: "",
    localDenom: "",
    accountPrefix: "",
    multiSend: false,
    liquidityModule: false
  };
}

export const RegisterZoneProposal = {
  encode(message: RegisterZoneProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    if (message.connectionId !== "") {
      writer.uint32(26).string(message.connectionId);
    }

    if (message.baseDenom !== "") {
      writer.uint32(34).string(message.baseDenom);
    }

    if (message.localDenom !== "") {
      writer.uint32(42).string(message.localDenom);
    }

    if (message.accountPrefix !== "") {
      writer.uint32(50).string(message.accountPrefix);
    }

    if (message.multiSend === true) {
      writer.uint32(56).bool(message.multiSend);
    }

    if (message.liquidityModule === true) {
      writer.uint32(64).bool(message.liquidityModule);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegisterZoneProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegisterZoneProposal();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;

        case 2:
          message.description = reader.string();
          break;

        case 3:
          message.connectionId = reader.string();
          break;

        case 4:
          message.baseDenom = reader.string();
          break;

        case 5:
          message.localDenom = reader.string();
          break;

        case 6:
          message.accountPrefix = reader.string();
          break;

        case 7:
          message.multiSend = reader.bool();
          break;

        case 8:
          message.liquidityModule = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): RegisterZoneProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      connectionId: isSet(object.connectionId) ? String(object.connectionId) : "",
      baseDenom: isSet(object.baseDenom) ? String(object.baseDenom) : "",
      localDenom: isSet(object.localDenom) ? String(object.localDenom) : "",
      accountPrefix: isSet(object.accountPrefix) ? String(object.accountPrefix) : "",
      multiSend: isSet(object.multiSend) ? Boolean(object.multiSend) : false,
      liquidityModule: isSet(object.liquidityModule) ? Boolean(object.liquidityModule) : false
    };
  },

  toJSON(message: RegisterZoneProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.connectionId !== undefined && (obj.connectionId = message.connectionId);
    message.baseDenom !== undefined && (obj.baseDenom = message.baseDenom);
    message.localDenom !== undefined && (obj.localDenom = message.localDenom);
    message.accountPrefix !== undefined && (obj.accountPrefix = message.accountPrefix);
    message.multiSend !== undefined && (obj.multiSend = message.multiSend);
    message.liquidityModule !== undefined && (obj.liquidityModule = message.liquidityModule);
    return obj;
  },

  fromPartial(object: Partial<RegisterZoneProposal>): RegisterZoneProposal {
    const message = createBaseRegisterZoneProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.connectionId = object.connectionId ?? "";
    message.baseDenom = object.baseDenom ?? "";
    message.localDenom = object.localDenom ?? "";
    message.accountPrefix = object.accountPrefix ?? "";
    message.multiSend = object.multiSend ?? false;
    message.liquidityModule = object.liquidityModule ?? false;
    return message;
  }

};

function createBaseRegisterZoneProposalWithDeposit(): RegisterZoneProposalWithDeposit {
  return {
    title: "",
    description: "",
    connectionId: "",
    baseDenom: "",
    localDenom: "",
    accountPrefix: "",
    multiSend: false,
    liquidityModule: false,
    deposit: ""
  };
}

export const RegisterZoneProposalWithDeposit = {
  encode(message: RegisterZoneProposalWithDeposit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    if (message.connectionId !== "") {
      writer.uint32(26).string(message.connectionId);
    }

    if (message.baseDenom !== "") {
      writer.uint32(34).string(message.baseDenom);
    }

    if (message.localDenom !== "") {
      writer.uint32(42).string(message.localDenom);
    }

    if (message.accountPrefix !== "") {
      writer.uint32(50).string(message.accountPrefix);
    }

    if (message.multiSend === true) {
      writer.uint32(56).bool(message.multiSend);
    }

    if (message.liquidityModule === true) {
      writer.uint32(64).bool(message.liquidityModule);
    }

    if (message.deposit !== "") {
      writer.uint32(74).string(message.deposit);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegisterZoneProposalWithDeposit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegisterZoneProposalWithDeposit();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;

        case 2:
          message.description = reader.string();
          break;

        case 3:
          message.connectionId = reader.string();
          break;

        case 4:
          message.baseDenom = reader.string();
          break;

        case 5:
          message.localDenom = reader.string();
          break;

        case 6:
          message.accountPrefix = reader.string();
          break;

        case 7:
          message.multiSend = reader.bool();
          break;

        case 8:
          message.liquidityModule = reader.bool();
          break;

        case 9:
          message.deposit = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): RegisterZoneProposalWithDeposit {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      connectionId: isSet(object.connectionId) ? String(object.connectionId) : "",
      baseDenom: isSet(object.baseDenom) ? String(object.baseDenom) : "",
      localDenom: isSet(object.localDenom) ? String(object.localDenom) : "",
      accountPrefix: isSet(object.accountPrefix) ? String(object.accountPrefix) : "",
      multiSend: isSet(object.multiSend) ? Boolean(object.multiSend) : false,
      liquidityModule: isSet(object.liquidityModule) ? Boolean(object.liquidityModule) : false,
      deposit: isSet(object.deposit) ? String(object.deposit) : ""
    };
  },

  toJSON(message: RegisterZoneProposalWithDeposit): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.connectionId !== undefined && (obj.connectionId = message.connectionId);
    message.baseDenom !== undefined && (obj.baseDenom = message.baseDenom);
    message.localDenom !== undefined && (obj.localDenom = message.localDenom);
    message.accountPrefix !== undefined && (obj.accountPrefix = message.accountPrefix);
    message.multiSend !== undefined && (obj.multiSend = message.multiSend);
    message.liquidityModule !== undefined && (obj.liquidityModule = message.liquidityModule);
    message.deposit !== undefined && (obj.deposit = message.deposit);
    return obj;
  },

  fromPartial(object: Partial<RegisterZoneProposalWithDeposit>): RegisterZoneProposalWithDeposit {
    const message = createBaseRegisterZoneProposalWithDeposit();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.connectionId = object.connectionId ?? "";
    message.baseDenom = object.baseDenom ?? "";
    message.localDenom = object.localDenom ?? "";
    message.accountPrefix = object.accountPrefix ?? "";
    message.multiSend = object.multiSend ?? false;
    message.liquidityModule = object.liquidityModule ?? false;
    message.deposit = object.deposit ?? "";
    return message;
  }

};

function createBaseUpdateZoneProposal(): UpdateZoneProposal {
  return {
    title: "",
    description: "",
    chainId: "",
    changes: []
  };
}

export const UpdateZoneProposal = {
  encode(message: UpdateZoneProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    if (message.chainId !== "") {
      writer.uint32(26).string(message.chainId);
    }

    for (const v of message.changes) {
      UpdateZoneValue.encode(v!, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateZoneProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateZoneProposal();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;

        case 2:
          message.description = reader.string();
          break;

        case 3:
          message.chainId = reader.string();
          break;

        case 4:
          message.changes.push(UpdateZoneValue.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): UpdateZoneProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      changes: Array.isArray(object?.changes) ? object.changes.map((e: any) => UpdateZoneValue.fromJSON(e)) : []
    };
  },

  toJSON(message: UpdateZoneProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.chainId !== undefined && (obj.chainId = message.chainId);

    if (message.changes) {
      obj.changes = message.changes.map(e => e ? UpdateZoneValue.toJSON(e) : undefined);
    } else {
      obj.changes = [];
    }

    return obj;
  },

  fromPartial(object: Partial<UpdateZoneProposal>): UpdateZoneProposal {
    const message = createBaseUpdateZoneProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.chainId = object.chainId ?? "";
    message.changes = object.changes?.map(e => UpdateZoneValue.fromPartial(e)) || [];
    return message;
  }

};

function createBaseUpdateZoneProposalWithDeposit(): UpdateZoneProposalWithDeposit {
  return {
    title: "",
    description: "",
    chainId: "",
    changes: [],
    deposit: ""
  };
}

export const UpdateZoneProposalWithDeposit = {
  encode(message: UpdateZoneProposalWithDeposit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    if (message.chainId !== "") {
      writer.uint32(26).string(message.chainId);
    }

    for (const v of message.changes) {
      UpdateZoneValue.encode(v!, writer.uint32(34).fork()).ldelim();
    }

    if (message.deposit !== "") {
      writer.uint32(42).string(message.deposit);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateZoneProposalWithDeposit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateZoneProposalWithDeposit();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;

        case 2:
          message.description = reader.string();
          break;

        case 3:
          message.chainId = reader.string();
          break;

        case 4:
          message.changes.push(UpdateZoneValue.decode(reader, reader.uint32()));
          break;

        case 5:
          message.deposit = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): UpdateZoneProposalWithDeposit {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      chainId: isSet(object.chainId) ? String(object.chainId) : "",
      changes: Array.isArray(object?.changes) ? object.changes.map((e: any) => UpdateZoneValue.fromJSON(e)) : [],
      deposit: isSet(object.deposit) ? String(object.deposit) : ""
    };
  },

  toJSON(message: UpdateZoneProposalWithDeposit): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.chainId !== undefined && (obj.chainId = message.chainId);

    if (message.changes) {
      obj.changes = message.changes.map(e => e ? UpdateZoneValue.toJSON(e) : undefined);
    } else {
      obj.changes = [];
    }

    message.deposit !== undefined && (obj.deposit = message.deposit);
    return obj;
  },

  fromPartial(object: Partial<UpdateZoneProposalWithDeposit>): UpdateZoneProposalWithDeposit {
    const message = createBaseUpdateZoneProposalWithDeposit();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.chainId = object.chainId ?? "";
    message.changes = object.changes?.map(e => UpdateZoneValue.fromPartial(e)) || [];
    message.deposit = object.deposit ?? "";
    return message;
  }

};

function createBaseUpdateZoneValue(): UpdateZoneValue {
  return {
    key: "",
    value: ""
  };
}

export const UpdateZoneValue = {
  encode(message: UpdateZoneValue, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }

    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateZoneValue {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateZoneValue();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;

        case 2:
          message.value = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): UpdateZoneValue {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? String(object.value) : ""
    };
  },

  toJSON(message: UpdateZoneValue): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  fromPartial(object: Partial<UpdateZoneValue>): UpdateZoneValue {
    const message = createBaseUpdateZoneValue();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  }

};