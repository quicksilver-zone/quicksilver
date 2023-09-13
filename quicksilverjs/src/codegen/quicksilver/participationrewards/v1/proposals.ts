import * as _m0 from "protobufjs/minimal";
import { isSet, bytesFromBase64, base64FromBytes } from "../../../helpers";
export interface AddProtocolDataProposal {
  title: string;
  description: string;
  type: string;
  data: string;
  key: string;
}
export interface AddProtocolDataProposalSDKType {
  title: string;
  description: string;
  type: string;
  data: string;
  key: string;
}
export interface AddProtocolDataProposalWithDeposit {
  title: string;
  description: string;
  protocol: string;
  type: string;
  key: string;
  data: Uint8Array;
  deposit: string;
}
export interface AddProtocolDataProposalWithDepositSDKType {
  title: string;
  description: string;
  protocol: string;
  type: string;
  key: string;
  data: Uint8Array;
  deposit: string;
}

function createBaseAddProtocolDataProposal(): AddProtocolDataProposal {
  return {
    title: "",
    description: "",
    type: "",
    data: "",
    key: ""
  };
}

export const AddProtocolDataProposal = {
  encode(message: AddProtocolDataProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    if (message.type !== "") {
      writer.uint32(34).string(message.type);
    }

    if (message.data !== "") {
      writer.uint32(42).string(message.data);
    }

    if (message.key !== "") {
      writer.uint32(50).string(message.key);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddProtocolDataProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddProtocolDataProposal();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;

        case 2:
          message.description = reader.string();
          break;

        case 4:
          message.type = reader.string();
          break;

        case 5:
          message.data = reader.string();
          break;

        case 6:
          message.key = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): AddProtocolDataProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      type: isSet(object.type) ? String(object.type) : "",
      data: isSet(object.data) ? String(object.data) : "",
      key: isSet(object.key) ? String(object.key) : ""
    };
  },

  toJSON(message: AddProtocolDataProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.type !== undefined && (obj.type = message.type);
    message.data !== undefined && (obj.data = message.data);
    message.key !== undefined && (obj.key = message.key);
    return obj;
  },

  fromPartial(object: Partial<AddProtocolDataProposal>): AddProtocolDataProposal {
    const message = createBaseAddProtocolDataProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.type = object.type ?? "";
    message.data = object.data ?? "";
    message.key = object.key ?? "";
    return message;
  }

};

function createBaseAddProtocolDataProposalWithDeposit(): AddProtocolDataProposalWithDeposit {
  return {
    title: "",
    description: "",
    protocol: "",
    type: "",
    key: "",
    data: new Uint8Array(),
    deposit: ""
  };
}

export const AddProtocolDataProposalWithDeposit = {
  encode(message: AddProtocolDataProposalWithDeposit, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    if (message.protocol !== "") {
      writer.uint32(26).string(message.protocol);
    }

    if (message.type !== "") {
      writer.uint32(34).string(message.type);
    }

    if (message.key !== "") {
      writer.uint32(42).string(message.key);
    }

    if (message.data.length !== 0) {
      writer.uint32(50).bytes(message.data);
    }

    if (message.deposit !== "") {
      writer.uint32(58).string(message.deposit);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddProtocolDataProposalWithDeposit {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddProtocolDataProposalWithDeposit();

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
          message.protocol = reader.string();
          break;

        case 4:
          message.type = reader.string();
          break;

        case 5:
          message.key = reader.string();
          break;

        case 6:
          message.data = reader.bytes();
          break;

        case 7:
          message.deposit = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): AddProtocolDataProposalWithDeposit {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      protocol: isSet(object.protocol) ? String(object.protocol) : "",
      type: isSet(object.type) ? String(object.type) : "",
      key: isSet(object.key) ? String(object.key) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(),
      deposit: isSet(object.deposit) ? String(object.deposit) : ""
    };
  },

  toJSON(message: AddProtocolDataProposalWithDeposit): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.protocol !== undefined && (obj.protocol = message.protocol);
    message.type !== undefined && (obj.type = message.type);
    message.key !== undefined && (obj.key = message.key);
    message.data !== undefined && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    message.deposit !== undefined && (obj.deposit = message.deposit);
    return obj;
  },

  fromPartial(object: Partial<AddProtocolDataProposalWithDeposit>): AddProtocolDataProposalWithDeposit {
    const message = createBaseAddProtocolDataProposalWithDeposit();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.protocol = object.protocol ?? "";
    message.type = object.type ?? "";
    message.key = object.key ?? "";
    message.data = object.data ?? new Uint8Array();
    message.deposit = object.deposit ?? "";
    return message;
  }

};