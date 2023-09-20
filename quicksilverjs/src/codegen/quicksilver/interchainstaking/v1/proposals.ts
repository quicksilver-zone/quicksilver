import { Long, isSet, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.interchainstaking.v1";
export interface RegisterZoneProposal {
  title: string;
  description: string;
  connectionId: string;
  baseDenom: string;
  localDenom: string;
  accountPrefix: string;
  multiSend: boolean;
  liquidityModule: boolean;
  messagesPerTx: Long;
}
export interface RegisterZoneProposalProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.RegisterZoneProposal";
  value: Uint8Array;
}
export interface RegisterZoneProposalAmino {
  title: string;
  description: string;
  connection_id: string;
  base_denom: string;
  local_denom: string;
  account_prefix: string;
  multi_send: boolean;
  liquidity_module: boolean;
  messages_per_tx: string;
}
export interface RegisterZoneProposalAminoMsg {
  type: "/quicksilver.interchainstaking.v1.RegisterZoneProposal";
  value: RegisterZoneProposalAmino;
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
  messages_per_tx: Long;
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
  messagesPerTx: Long;
}
export interface RegisterZoneProposalWithDepositProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.RegisterZoneProposalWithDeposit";
  value: Uint8Array;
}
export interface RegisterZoneProposalWithDepositAmino {
  title: string;
  description: string;
  connection_id: string;
  base_denom: string;
  local_denom: string;
  account_prefix: string;
  multi_send: boolean;
  liquidity_module: boolean;
  deposit: string;
  messages_per_tx: string;
}
export interface RegisterZoneProposalWithDepositAminoMsg {
  type: "/quicksilver.interchainstaking.v1.RegisterZoneProposalWithDeposit";
  value: RegisterZoneProposalWithDepositAmino;
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
  messages_per_tx: Long;
}
export interface UpdateZoneProposal {
  title: string;
  description: string;
  chainId: string;
  changes: UpdateZoneValue[];
}
export interface UpdateZoneProposalProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneProposal";
  value: Uint8Array;
}
export interface UpdateZoneProposalAmino {
  title: string;
  description: string;
  chain_id: string;
  changes: UpdateZoneValueAmino[];
}
export interface UpdateZoneProposalAminoMsg {
  type: "/quicksilver.interchainstaking.v1.UpdateZoneProposal";
  value: UpdateZoneProposalAmino;
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
export interface UpdateZoneProposalWithDepositProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneProposalWithDeposit";
  value: Uint8Array;
}
export interface UpdateZoneProposalWithDepositAmino {
  title: string;
  description: string;
  chain_id: string;
  changes: UpdateZoneValueAmino[];
  deposit: string;
}
export interface UpdateZoneProposalWithDepositAminoMsg {
  type: "/quicksilver.interchainstaking.v1.UpdateZoneProposalWithDeposit";
  value: UpdateZoneProposalWithDepositAmino;
}
export interface UpdateZoneProposalWithDepositSDKType {
  title: string;
  description: string;
  chain_id: string;
  changes: UpdateZoneValueSDKType[];
  deposit: string;
}
/**
 * UpdateZoneValue defines an individual parameter change, for use in
 * UpdateZoneProposal.
 */
export interface UpdateZoneValue {
  key: string;
  value: string;
}
export interface UpdateZoneValueProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneValue";
  value: Uint8Array;
}
/**
 * UpdateZoneValue defines an individual parameter change, for use in
 * UpdateZoneProposal.
 */
export interface UpdateZoneValueAmino {
  key: string;
  value: string;
}
export interface UpdateZoneValueAminoMsg {
  type: "/quicksilver.interchainstaking.v1.UpdateZoneValue";
  value: UpdateZoneValueAmino;
}
/**
 * UpdateZoneValue defines an individual parameter change, for use in
 * UpdateZoneProposal.
 */
export interface UpdateZoneValueSDKType {
  key: string;
  value: string;
}
export interface MsgGovReopenChannel {
  title: string;
  description: string;
  connectionId: string;
  portId: string;
  authority: string;
}
export interface MsgGovReopenChannelProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel";
  value: Uint8Array;
}
export interface MsgGovReopenChannelAmino {
  title: string;
  description: string;
  connection_id: string;
  port_id: string;
  authority: string;
}
export interface MsgGovReopenChannelAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel";
  value: MsgGovReopenChannelAmino;
}
export interface MsgGovReopenChannelSDKType {
  title: string;
  description: string;
  connection_id: string;
  port_id: string;
  authority: string;
}
/** MsgGovReopenChannelResponse defines the MsgGovReopenChannel response type. */
export interface MsgGovReopenChannelResponse {}
export interface MsgGovReopenChannelResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannelResponse";
  value: Uint8Array;
}
/** MsgGovReopenChannelResponse defines the MsgGovReopenChannel response type. */
export interface MsgGovReopenChannelResponseAmino {}
export interface MsgGovReopenChannelResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgGovReopenChannelResponse";
  value: MsgGovReopenChannelResponseAmino;
}
/** MsgGovReopenChannelResponse defines the MsgGovReopenChannel response type. */
export interface MsgGovReopenChannelResponseSDKType {}
export interface MsgGovCloseChannel {
  title: string;
  description: string;
  channelId: string;
  portId: string;
  authority: string;
}
export interface MsgGovCloseChannelProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel";
  value: Uint8Array;
}
export interface MsgGovCloseChannelAmino {
  title: string;
  description: string;
  channel_id: string;
  port_id: string;
  authority: string;
}
export interface MsgGovCloseChannelAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel";
  value: MsgGovCloseChannelAmino;
}
export interface MsgGovCloseChannelSDKType {
  title: string;
  description: string;
  channel_id: string;
  port_id: string;
  authority: string;
}
/** MsgGovCloseChannelResponse defines the MsgGovCloseChannel response type. */
export interface MsgGovCloseChannelResponse {}
export interface MsgGovCloseChannelResponseProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannelResponse";
  value: Uint8Array;
}
/** MsgGovCloseChannelResponse defines the MsgGovCloseChannel response type. */
export interface MsgGovCloseChannelResponseAmino {}
export interface MsgGovCloseChannelResponseAminoMsg {
  type: "/quicksilver.interchainstaking.v1.MsgGovCloseChannelResponse";
  value: MsgGovCloseChannelResponseAmino;
}
/** MsgGovCloseChannelResponse defines the MsgGovCloseChannel response type. */
export interface MsgGovCloseChannelResponseSDKType {}
function createBaseRegisterZoneProposal(): RegisterZoneProposal {
  return {
    title: "",
    description: "",
    connectionId: "",
    baseDenom: "",
    localDenom: "",
    accountPrefix: "",
    multiSend: false,
    liquidityModule: false,
    messagesPerTx: Long.ZERO
  };
}
export const RegisterZoneProposal = {
  typeUrl: "/quicksilver.interchainstaking.v1.RegisterZoneProposal",
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
    if (!message.messagesPerTx.isZero()) {
      writer.uint32(72).int64(message.messagesPerTx);
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
        case 9:
          message.messagesPerTx = (reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): RegisterZoneProposal {
    const obj = createBaseRegisterZoneProposal();
    if (isSet(object.title)) obj.title = String(object.title);
    if (isSet(object.description)) obj.description = String(object.description);
    if (isSet(object.connectionId)) obj.connectionId = String(object.connectionId);
    if (isSet(object.baseDenom)) obj.baseDenom = String(object.baseDenom);
    if (isSet(object.localDenom)) obj.localDenom = String(object.localDenom);
    if (isSet(object.accountPrefix)) obj.accountPrefix = String(object.accountPrefix);
    if (isSet(object.multiSend)) obj.multiSend = Boolean(object.multiSend);
    if (isSet(object.liquidityModule)) obj.liquidityModule = Boolean(object.liquidityModule);
    if (isSet(object.messagesPerTx)) obj.messagesPerTx = Long.fromValue(object.messagesPerTx);
    return obj;
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
    message.messagesPerTx !== undefined && (obj.messagesPerTx = (message.messagesPerTx || Long.ZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<RegisterZoneProposal>): RegisterZoneProposal {
    const message = createBaseRegisterZoneProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.connectionId = object.connectionId ?? "";
    message.baseDenom = object.baseDenom ?? "";
    message.localDenom = object.localDenom ?? "";
    message.accountPrefix = object.accountPrefix ?? "";
    message.multiSend = object.multiSend ?? false;
    message.liquidityModule = object.liquidityModule ?? false;
    if (object.messagesPerTx !== undefined && object.messagesPerTx !== null) {
      message.messagesPerTx = Long.fromValue(object.messagesPerTx);
    }
    return message;
  },
  fromSDK(object: RegisterZoneProposalSDKType): RegisterZoneProposal {
    return {
      title: object?.title,
      description: object?.description,
      connectionId: object?.connection_id,
      baseDenom: object?.base_denom,
      localDenom: object?.local_denom,
      accountPrefix: object?.account_prefix,
      multiSend: object?.multi_send,
      liquidityModule: object?.liquidity_module,
      messagesPerTx: object?.messages_per_tx
    };
  },
  toSDK(message: RegisterZoneProposal): RegisterZoneProposalSDKType {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.connection_id = message.connectionId;
    obj.base_denom = message.baseDenom;
    obj.local_denom = message.localDenom;
    obj.account_prefix = message.accountPrefix;
    obj.multi_send = message.multiSend;
    obj.liquidity_module = message.liquidityModule;
    obj.messages_per_tx = message.messagesPerTx;
    return obj;
  },
  fromAmino(object: RegisterZoneProposalAmino): RegisterZoneProposal {
    return {
      title: object.title,
      description: object.description,
      connectionId: object.connection_id,
      baseDenom: object.base_denom,
      localDenom: object.local_denom,
      accountPrefix: object.account_prefix,
      multiSend: object.multi_send,
      liquidityModule: object.liquidity_module,
      messagesPerTx: Long.fromString(object.messages_per_tx)
    };
  },
  toAmino(message: RegisterZoneProposal): RegisterZoneProposalAmino {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.connection_id = message.connectionId;
    obj.base_denom = message.baseDenom;
    obj.local_denom = message.localDenom;
    obj.account_prefix = message.accountPrefix;
    obj.multi_send = message.multiSend;
    obj.liquidity_module = message.liquidityModule;
    obj.messages_per_tx = message.messagesPerTx ? message.messagesPerTx.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: RegisterZoneProposalAminoMsg): RegisterZoneProposal {
    return RegisterZoneProposal.fromAmino(object.value);
  },
  fromProtoMsg(message: RegisterZoneProposalProtoMsg): RegisterZoneProposal {
    return RegisterZoneProposal.decode(message.value);
  },
  toProto(message: RegisterZoneProposal): Uint8Array {
    return RegisterZoneProposal.encode(message).finish();
  },
  toProtoMsg(message: RegisterZoneProposal): RegisterZoneProposalProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.RegisterZoneProposal",
      value: RegisterZoneProposal.encode(message).finish()
    };
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
    deposit: "",
    messagesPerTx: Long.ZERO
  };
}
export const RegisterZoneProposalWithDeposit = {
  typeUrl: "/quicksilver.interchainstaking.v1.RegisterZoneProposalWithDeposit",
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
    if (!message.messagesPerTx.isZero()) {
      writer.uint32(80).int64(message.messagesPerTx);
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
        case 10:
          message.messagesPerTx = (reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): RegisterZoneProposalWithDeposit {
    const obj = createBaseRegisterZoneProposalWithDeposit();
    if (isSet(object.title)) obj.title = String(object.title);
    if (isSet(object.description)) obj.description = String(object.description);
    if (isSet(object.connectionId)) obj.connectionId = String(object.connectionId);
    if (isSet(object.baseDenom)) obj.baseDenom = String(object.baseDenom);
    if (isSet(object.localDenom)) obj.localDenom = String(object.localDenom);
    if (isSet(object.accountPrefix)) obj.accountPrefix = String(object.accountPrefix);
    if (isSet(object.multiSend)) obj.multiSend = Boolean(object.multiSend);
    if (isSet(object.liquidityModule)) obj.liquidityModule = Boolean(object.liquidityModule);
    if (isSet(object.deposit)) obj.deposit = String(object.deposit);
    if (isSet(object.messagesPerTx)) obj.messagesPerTx = Long.fromValue(object.messagesPerTx);
    return obj;
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
    message.messagesPerTx !== undefined && (obj.messagesPerTx = (message.messagesPerTx || Long.ZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<RegisterZoneProposalWithDeposit>): RegisterZoneProposalWithDeposit {
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
    if (object.messagesPerTx !== undefined && object.messagesPerTx !== null) {
      message.messagesPerTx = Long.fromValue(object.messagesPerTx);
    }
    return message;
  },
  fromSDK(object: RegisterZoneProposalWithDepositSDKType): RegisterZoneProposalWithDeposit {
    return {
      title: object?.title,
      description: object?.description,
      connectionId: object?.connection_id,
      baseDenom: object?.base_denom,
      localDenom: object?.local_denom,
      accountPrefix: object?.account_prefix,
      multiSend: object?.multi_send,
      liquidityModule: object?.liquidity_module,
      deposit: object?.deposit,
      messagesPerTx: object?.messages_per_tx
    };
  },
  toSDK(message: RegisterZoneProposalWithDeposit): RegisterZoneProposalWithDepositSDKType {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.connection_id = message.connectionId;
    obj.base_denom = message.baseDenom;
    obj.local_denom = message.localDenom;
    obj.account_prefix = message.accountPrefix;
    obj.multi_send = message.multiSend;
    obj.liquidity_module = message.liquidityModule;
    obj.deposit = message.deposit;
    obj.messages_per_tx = message.messagesPerTx;
    return obj;
  },
  fromAmino(object: RegisterZoneProposalWithDepositAmino): RegisterZoneProposalWithDeposit {
    return {
      title: object.title,
      description: object.description,
      connectionId: object.connection_id,
      baseDenom: object.base_denom,
      localDenom: object.local_denom,
      accountPrefix: object.account_prefix,
      multiSend: object.multi_send,
      liquidityModule: object.liquidity_module,
      deposit: object.deposit,
      messagesPerTx: Long.fromString(object.messages_per_tx)
    };
  },
  toAmino(message: RegisterZoneProposalWithDeposit): RegisterZoneProposalWithDepositAmino {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.connection_id = message.connectionId;
    obj.base_denom = message.baseDenom;
    obj.local_denom = message.localDenom;
    obj.account_prefix = message.accountPrefix;
    obj.multi_send = message.multiSend;
    obj.liquidity_module = message.liquidityModule;
    obj.deposit = message.deposit;
    obj.messages_per_tx = message.messagesPerTx ? message.messagesPerTx.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: RegisterZoneProposalWithDepositAminoMsg): RegisterZoneProposalWithDeposit {
    return RegisterZoneProposalWithDeposit.fromAmino(object.value);
  },
  fromProtoMsg(message: RegisterZoneProposalWithDepositProtoMsg): RegisterZoneProposalWithDeposit {
    return RegisterZoneProposalWithDeposit.decode(message.value);
  },
  toProto(message: RegisterZoneProposalWithDeposit): Uint8Array {
    return RegisterZoneProposalWithDeposit.encode(message).finish();
  },
  toProtoMsg(message: RegisterZoneProposalWithDeposit): RegisterZoneProposalWithDepositProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.RegisterZoneProposalWithDeposit",
      value: RegisterZoneProposalWithDeposit.encode(message).finish()
    };
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
  typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneProposal",
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
    const obj = createBaseUpdateZoneProposal();
    if (isSet(object.title)) obj.title = String(object.title);
    if (isSet(object.description)) obj.description = String(object.description);
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (Array.isArray(object?.changes)) obj.changes = object.changes.map((e: any) => UpdateZoneValue.fromJSON(e));
    return obj;
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
  fromPartial(object: DeepPartial<UpdateZoneProposal>): UpdateZoneProposal {
    const message = createBaseUpdateZoneProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.chainId = object.chainId ?? "";
    message.changes = object.changes?.map(e => UpdateZoneValue.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: UpdateZoneProposalSDKType): UpdateZoneProposal {
    return {
      title: object?.title,
      description: object?.description,
      chainId: object?.chain_id,
      changes: Array.isArray(object?.changes) ? object.changes.map((e: any) => UpdateZoneValue.fromSDK(e)) : []
    };
  },
  toSDK(message: UpdateZoneProposal): UpdateZoneProposalSDKType {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.chain_id = message.chainId;
    if (message.changes) {
      obj.changes = message.changes.map(e => e ? UpdateZoneValue.toSDK(e) : undefined);
    } else {
      obj.changes = [];
    }
    return obj;
  },
  fromAmino(object: UpdateZoneProposalAmino): UpdateZoneProposal {
    return {
      title: object.title,
      description: object.description,
      chainId: object.chain_id,
      changes: Array.isArray(object?.changes) ? object.changes.map((e: any) => UpdateZoneValue.fromAmino(e)) : []
    };
  },
  toAmino(message: UpdateZoneProposal): UpdateZoneProposalAmino {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.chain_id = message.chainId;
    if (message.changes) {
      obj.changes = message.changes.map(e => e ? UpdateZoneValue.toAmino(e) : undefined);
    } else {
      obj.changes = [];
    }
    return obj;
  },
  fromAminoMsg(object: UpdateZoneProposalAminoMsg): UpdateZoneProposal {
    return UpdateZoneProposal.fromAmino(object.value);
  },
  fromProtoMsg(message: UpdateZoneProposalProtoMsg): UpdateZoneProposal {
    return UpdateZoneProposal.decode(message.value);
  },
  toProto(message: UpdateZoneProposal): Uint8Array {
    return UpdateZoneProposal.encode(message).finish();
  },
  toProtoMsg(message: UpdateZoneProposal): UpdateZoneProposalProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneProposal",
      value: UpdateZoneProposal.encode(message).finish()
    };
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
  typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneProposalWithDeposit",
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
    const obj = createBaseUpdateZoneProposalWithDeposit();
    if (isSet(object.title)) obj.title = String(object.title);
    if (isSet(object.description)) obj.description = String(object.description);
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (Array.isArray(object?.changes)) obj.changes = object.changes.map((e: any) => UpdateZoneValue.fromJSON(e));
    if (isSet(object.deposit)) obj.deposit = String(object.deposit);
    return obj;
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
  fromPartial(object: DeepPartial<UpdateZoneProposalWithDeposit>): UpdateZoneProposalWithDeposit {
    const message = createBaseUpdateZoneProposalWithDeposit();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.chainId = object.chainId ?? "";
    message.changes = object.changes?.map(e => UpdateZoneValue.fromPartial(e)) || [];
    message.deposit = object.deposit ?? "";
    return message;
  },
  fromSDK(object: UpdateZoneProposalWithDepositSDKType): UpdateZoneProposalWithDeposit {
    return {
      title: object?.title,
      description: object?.description,
      chainId: object?.chain_id,
      changes: Array.isArray(object?.changes) ? object.changes.map((e: any) => UpdateZoneValue.fromSDK(e)) : [],
      deposit: object?.deposit
    };
  },
  toSDK(message: UpdateZoneProposalWithDeposit): UpdateZoneProposalWithDepositSDKType {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.chain_id = message.chainId;
    if (message.changes) {
      obj.changes = message.changes.map(e => e ? UpdateZoneValue.toSDK(e) : undefined);
    } else {
      obj.changes = [];
    }
    obj.deposit = message.deposit;
    return obj;
  },
  fromAmino(object: UpdateZoneProposalWithDepositAmino): UpdateZoneProposalWithDeposit {
    return {
      title: object.title,
      description: object.description,
      chainId: object.chain_id,
      changes: Array.isArray(object?.changes) ? object.changes.map((e: any) => UpdateZoneValue.fromAmino(e)) : [],
      deposit: object.deposit
    };
  },
  toAmino(message: UpdateZoneProposalWithDeposit): UpdateZoneProposalWithDepositAmino {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.chain_id = message.chainId;
    if (message.changes) {
      obj.changes = message.changes.map(e => e ? UpdateZoneValue.toAmino(e) : undefined);
    } else {
      obj.changes = [];
    }
    obj.deposit = message.deposit;
    return obj;
  },
  fromAminoMsg(object: UpdateZoneProposalWithDepositAminoMsg): UpdateZoneProposalWithDeposit {
    return UpdateZoneProposalWithDeposit.fromAmino(object.value);
  },
  fromProtoMsg(message: UpdateZoneProposalWithDepositProtoMsg): UpdateZoneProposalWithDeposit {
    return UpdateZoneProposalWithDeposit.decode(message.value);
  },
  toProto(message: UpdateZoneProposalWithDeposit): Uint8Array {
    return UpdateZoneProposalWithDeposit.encode(message).finish();
  },
  toProtoMsg(message: UpdateZoneProposalWithDeposit): UpdateZoneProposalWithDepositProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneProposalWithDeposit",
      value: UpdateZoneProposalWithDeposit.encode(message).finish()
    };
  }
};
function createBaseUpdateZoneValue(): UpdateZoneValue {
  return {
    key: "",
    value: ""
  };
}
export const UpdateZoneValue = {
  typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneValue",
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
    const obj = createBaseUpdateZoneValue();
    if (isSet(object.key)) obj.key = String(object.key);
    if (isSet(object.value)) obj.value = String(object.value);
    return obj;
  },
  toJSON(message: UpdateZoneValue): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },
  fromPartial(object: DeepPartial<UpdateZoneValue>): UpdateZoneValue {
    const message = createBaseUpdateZoneValue();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
  fromSDK(object: UpdateZoneValueSDKType): UpdateZoneValue {
    return {
      key: object?.key,
      value: object?.value
    };
  },
  toSDK(message: UpdateZoneValue): UpdateZoneValueSDKType {
    const obj: any = {};
    obj.key = message.key;
    obj.value = message.value;
    return obj;
  },
  fromAmino(object: UpdateZoneValueAmino): UpdateZoneValue {
    return {
      key: object.key,
      value: object.value
    };
  },
  toAmino(message: UpdateZoneValue): UpdateZoneValueAmino {
    const obj: any = {};
    obj.key = message.key;
    obj.value = message.value;
    return obj;
  },
  fromAminoMsg(object: UpdateZoneValueAminoMsg): UpdateZoneValue {
    return UpdateZoneValue.fromAmino(object.value);
  },
  fromProtoMsg(message: UpdateZoneValueProtoMsg): UpdateZoneValue {
    return UpdateZoneValue.decode(message.value);
  },
  toProto(message: UpdateZoneValue): Uint8Array {
    return UpdateZoneValue.encode(message).finish();
  },
  toProtoMsg(message: UpdateZoneValue): UpdateZoneValueProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.UpdateZoneValue",
      value: UpdateZoneValue.encode(message).finish()
    };
  }
};
function createBaseMsgGovReopenChannel(): MsgGovReopenChannel {
  return {
    title: "",
    description: "",
    connectionId: "",
    portId: "",
    authority: ""
  };
}
export const MsgGovReopenChannel = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
  encode(message: MsgGovReopenChannel, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.connectionId !== "") {
      writer.uint32(26).string(message.connectionId);
    }
    if (message.portId !== "") {
      writer.uint32(34).string(message.portId);
    }
    if (message.authority !== "") {
      writer.uint32(42).string(message.authority);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgGovReopenChannel {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGovReopenChannel();
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
          message.portId = reader.string();
          break;
        case 5:
          message.authority = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgGovReopenChannel {
    const obj = createBaseMsgGovReopenChannel();
    if (isSet(object.title)) obj.title = String(object.title);
    if (isSet(object.description)) obj.description = String(object.description);
    if (isSet(object.connectionId)) obj.connectionId = String(object.connectionId);
    if (isSet(object.portId)) obj.portId = String(object.portId);
    if (isSet(object.authority)) obj.authority = String(object.authority);
    return obj;
  },
  toJSON(message: MsgGovReopenChannel): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.connectionId !== undefined && (obj.connectionId = message.connectionId);
    message.portId !== undefined && (obj.portId = message.portId);
    message.authority !== undefined && (obj.authority = message.authority);
    return obj;
  },
  fromPartial(object: DeepPartial<MsgGovReopenChannel>): MsgGovReopenChannel {
    const message = createBaseMsgGovReopenChannel();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.connectionId = object.connectionId ?? "";
    message.portId = object.portId ?? "";
    message.authority = object.authority ?? "";
    return message;
  },
  fromSDK(object: MsgGovReopenChannelSDKType): MsgGovReopenChannel {
    return {
      title: object?.title,
      description: object?.description,
      connectionId: object?.connection_id,
      portId: object?.port_id,
      authority: object?.authority
    };
  },
  toSDK(message: MsgGovReopenChannel): MsgGovReopenChannelSDKType {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.connection_id = message.connectionId;
    obj.port_id = message.portId;
    obj.authority = message.authority;
    return obj;
  },
  fromAmino(object: MsgGovReopenChannelAmino): MsgGovReopenChannel {
    return {
      title: object.title,
      description: object.description,
      connectionId: object.connection_id,
      portId: object.port_id,
      authority: object.authority
    };
  },
  toAmino(message: MsgGovReopenChannel): MsgGovReopenChannelAmino {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.connection_id = message.connectionId;
    obj.port_id = message.portId;
    obj.authority = message.authority;
    return obj;
  },
  fromAminoMsg(object: MsgGovReopenChannelAminoMsg): MsgGovReopenChannel {
    return MsgGovReopenChannel.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgGovReopenChannelProtoMsg): MsgGovReopenChannel {
    return MsgGovReopenChannel.decode(message.value);
  },
  toProto(message: MsgGovReopenChannel): Uint8Array {
    return MsgGovReopenChannel.encode(message).finish();
  },
  toProtoMsg(message: MsgGovReopenChannel): MsgGovReopenChannelProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannel",
      value: MsgGovReopenChannel.encode(message).finish()
    };
  }
};
function createBaseMsgGovReopenChannelResponse(): MsgGovReopenChannelResponse {
  return {};
}
export const MsgGovReopenChannelResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannelResponse",
  encode(_: MsgGovReopenChannelResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgGovReopenChannelResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGovReopenChannelResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(_: any): MsgGovReopenChannelResponse {
    const obj = createBaseMsgGovReopenChannelResponse();
    return obj;
  },
  toJSON(_: MsgGovReopenChannelResponse): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<MsgGovReopenChannelResponse>): MsgGovReopenChannelResponse {
    const message = createBaseMsgGovReopenChannelResponse();
    return message;
  },
  fromSDK(_: MsgGovReopenChannelResponseSDKType): MsgGovReopenChannelResponse {
    return {};
  },
  toSDK(_: MsgGovReopenChannelResponse): MsgGovReopenChannelResponseSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: MsgGovReopenChannelResponseAmino): MsgGovReopenChannelResponse {
    return {};
  },
  toAmino(_: MsgGovReopenChannelResponse): MsgGovReopenChannelResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgGovReopenChannelResponseAminoMsg): MsgGovReopenChannelResponse {
    return MsgGovReopenChannelResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgGovReopenChannelResponseProtoMsg): MsgGovReopenChannelResponse {
    return MsgGovReopenChannelResponse.decode(message.value);
  },
  toProto(message: MsgGovReopenChannelResponse): Uint8Array {
    return MsgGovReopenChannelResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgGovReopenChannelResponse): MsgGovReopenChannelResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgGovReopenChannelResponse",
      value: MsgGovReopenChannelResponse.encode(message).finish()
    };
  }
};
function createBaseMsgGovCloseChannel(): MsgGovCloseChannel {
  return {
    title: "",
    description: "",
    channelId: "",
    portId: "",
    authority: ""
  };
}
export const MsgGovCloseChannel = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
  encode(message: MsgGovCloseChannel, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.channelId !== "") {
      writer.uint32(26).string(message.channelId);
    }
    if (message.portId !== "") {
      writer.uint32(34).string(message.portId);
    }
    if (message.authority !== "") {
      writer.uint32(42).string(message.authority);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgGovCloseChannel {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGovCloseChannel();
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
          message.channelId = reader.string();
          break;
        case 4:
          message.portId = reader.string();
          break;
        case 5:
          message.authority = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgGovCloseChannel {
    const obj = createBaseMsgGovCloseChannel();
    if (isSet(object.title)) obj.title = String(object.title);
    if (isSet(object.description)) obj.description = String(object.description);
    if (isSet(object.channelId)) obj.channelId = String(object.channelId);
    if (isSet(object.portId)) obj.portId = String(object.portId);
    if (isSet(object.authority)) obj.authority = String(object.authority);
    return obj;
  },
  toJSON(message: MsgGovCloseChannel): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.portId !== undefined && (obj.portId = message.portId);
    message.authority !== undefined && (obj.authority = message.authority);
    return obj;
  },
  fromPartial(object: DeepPartial<MsgGovCloseChannel>): MsgGovCloseChannel {
    const message = createBaseMsgGovCloseChannel();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.channelId = object.channelId ?? "";
    message.portId = object.portId ?? "";
    message.authority = object.authority ?? "";
    return message;
  },
  fromSDK(object: MsgGovCloseChannelSDKType): MsgGovCloseChannel {
    return {
      title: object?.title,
      description: object?.description,
      channelId: object?.channel_id,
      portId: object?.port_id,
      authority: object?.authority
    };
  },
  toSDK(message: MsgGovCloseChannel): MsgGovCloseChannelSDKType {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.channel_id = message.channelId;
    obj.port_id = message.portId;
    obj.authority = message.authority;
    return obj;
  },
  fromAmino(object: MsgGovCloseChannelAmino): MsgGovCloseChannel {
    return {
      title: object.title,
      description: object.description,
      channelId: object.channel_id,
      portId: object.port_id,
      authority: object.authority
    };
  },
  toAmino(message: MsgGovCloseChannel): MsgGovCloseChannelAmino {
    const obj: any = {};
    obj.title = message.title;
    obj.description = message.description;
    obj.channel_id = message.channelId;
    obj.port_id = message.portId;
    obj.authority = message.authority;
    return obj;
  },
  fromAminoMsg(object: MsgGovCloseChannelAminoMsg): MsgGovCloseChannel {
    return MsgGovCloseChannel.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgGovCloseChannelProtoMsg): MsgGovCloseChannel {
    return MsgGovCloseChannel.decode(message.value);
  },
  toProto(message: MsgGovCloseChannel): Uint8Array {
    return MsgGovCloseChannel.encode(message).finish();
  },
  toProtoMsg(message: MsgGovCloseChannel): MsgGovCloseChannelProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannel",
      value: MsgGovCloseChannel.encode(message).finish()
    };
  }
};
function createBaseMsgGovCloseChannelResponse(): MsgGovCloseChannelResponse {
  return {};
}
export const MsgGovCloseChannelResponse = {
  typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannelResponse",
  encode(_: MsgGovCloseChannelResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgGovCloseChannelResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGovCloseChannelResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(_: any): MsgGovCloseChannelResponse {
    const obj = createBaseMsgGovCloseChannelResponse();
    return obj;
  },
  toJSON(_: MsgGovCloseChannelResponse): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<MsgGovCloseChannelResponse>): MsgGovCloseChannelResponse {
    const message = createBaseMsgGovCloseChannelResponse();
    return message;
  },
  fromSDK(_: MsgGovCloseChannelResponseSDKType): MsgGovCloseChannelResponse {
    return {};
  },
  toSDK(_: MsgGovCloseChannelResponse): MsgGovCloseChannelResponseSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: MsgGovCloseChannelResponseAmino): MsgGovCloseChannelResponse {
    return {};
  },
  toAmino(_: MsgGovCloseChannelResponse): MsgGovCloseChannelResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgGovCloseChannelResponseAminoMsg): MsgGovCloseChannelResponse {
    return MsgGovCloseChannelResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgGovCloseChannelResponseProtoMsg): MsgGovCloseChannelResponse {
    return MsgGovCloseChannelResponse.decode(message.value);
  },
  toProto(message: MsgGovCloseChannelResponse): Uint8Array {
    return MsgGovCloseChannelResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgGovCloseChannelResponse): MsgGovCloseChannelResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.MsgGovCloseChannelResponse",
      value: MsgGovCloseChannelResponse.encode(message).finish()
    };
  }
};