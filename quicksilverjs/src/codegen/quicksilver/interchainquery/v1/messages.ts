import { ProofOps, ProofOpsAmino, ProofOpsSDKType } from "../../../tendermint/crypto/proof";
import { Long, isSet, bytesFromBase64, base64FromBytes, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.interchainquery.v1";
/** MsgSubmitQueryResponse represents a message type to fulfil a query request. */
export interface MsgSubmitQueryResponse {
  chainId: string;
  queryId: string;
  result: Uint8Array;
  proofOps: ProofOps;
  height: Long;
  fromAddress: string;
}
export interface MsgSubmitQueryResponseProtoMsg {
  typeUrl: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse";
  value: Uint8Array;
}
/** MsgSubmitQueryResponse represents a message type to fulfil a query request. */
export interface MsgSubmitQueryResponseAmino {
  chain_id: string;
  query_id: string;
  result: Uint8Array;
  proof_ops?: ProofOpsAmino;
  height: string;
  from_address: string;
}
export interface MsgSubmitQueryResponseAminoMsg {
  type: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse";
  value: MsgSubmitQueryResponseAmino;
}
/** MsgSubmitQueryResponse represents a message type to fulfil a query request. */
export interface MsgSubmitQueryResponseSDKType {
  chain_id: string;
  query_id: string;
  result: Uint8Array;
  proof_ops: ProofOpsSDKType;
  height: Long;
  from_address: string;
}
/**
 * MsgSubmitQueryResponseResponse defines the MsgSubmitQueryResponse response
 * type.
 */
export interface MsgSubmitQueryResponseResponse {}
export interface MsgSubmitQueryResponseResponseProtoMsg {
  typeUrl: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse";
  value: Uint8Array;
}
/**
 * MsgSubmitQueryResponseResponse defines the MsgSubmitQueryResponse response
 * type.
 */
export interface MsgSubmitQueryResponseResponseAmino {}
export interface MsgSubmitQueryResponseResponseAminoMsg {
  type: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse";
  value: MsgSubmitQueryResponseResponseAmino;
}
/**
 * MsgSubmitQueryResponseResponse defines the MsgSubmitQueryResponse response
 * type.
 */
export interface MsgSubmitQueryResponseResponseSDKType {}
function createBaseMsgSubmitQueryResponse(): MsgSubmitQueryResponse {
  return {
    chainId: "",
    queryId: "",
    result: new Uint8Array(),
    proofOps: ProofOps.fromPartial({}),
    height: Long.ZERO,
    fromAddress: ""
  };
}
export const MsgSubmitQueryResponse = {
  typeUrl: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse",
  encode(message: MsgSubmitQueryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (message.queryId !== "") {
      writer.uint32(18).string(message.queryId);
    }
    if (message.result.length !== 0) {
      writer.uint32(26).bytes(message.result);
    }
    if (message.proofOps !== undefined) {
      ProofOps.encode(message.proofOps, writer.uint32(34).fork()).ldelim();
    }
    if (!message.height.isZero()) {
      writer.uint32(40).int64(message.height);
    }
    if (message.fromAddress !== "") {
      writer.uint32(50).string(message.fromAddress);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitQueryResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSubmitQueryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.queryId = reader.string();
          break;
        case 3:
          message.result = reader.bytes();
          break;
        case 4:
          message.proofOps = ProofOps.decode(reader, reader.uint32());
          break;
        case 5:
          message.height = (reader.int64() as Long);
          break;
        case 6:
          message.fromAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): MsgSubmitQueryResponse {
    const obj = createBaseMsgSubmitQueryResponse();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.queryId)) obj.queryId = String(object.queryId);
    if (isSet(object.result)) obj.result = bytesFromBase64(object.result);
    if (isSet(object.proofOps)) obj.proofOps = ProofOps.fromJSON(object.proofOps);
    if (isSet(object.height)) obj.height = Long.fromValue(object.height);
    if (isSet(object.fromAddress)) obj.fromAddress = String(object.fromAddress);
    return obj;
  },
  toJSON(message: MsgSubmitQueryResponse): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.queryId !== undefined && (obj.queryId = message.queryId);
    message.result !== undefined && (obj.result = base64FromBytes(message.result !== undefined ? message.result : new Uint8Array()));
    message.proofOps !== undefined && (obj.proofOps = message.proofOps ? ProofOps.toJSON(message.proofOps) : undefined);
    message.height !== undefined && (obj.height = (message.height || Long.ZERO).toString());
    message.fromAddress !== undefined && (obj.fromAddress = message.fromAddress);
    return obj;
  },
  fromPartial(object: DeepPartial<MsgSubmitQueryResponse>): MsgSubmitQueryResponse {
    const message = createBaseMsgSubmitQueryResponse();
    message.chainId = object.chainId ?? "";
    message.queryId = object.queryId ?? "";
    message.result = object.result ?? new Uint8Array();
    if (object.proofOps !== undefined && object.proofOps !== null) {
      message.proofOps = ProofOps.fromPartial(object.proofOps);
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = Long.fromValue(object.height);
    }
    message.fromAddress = object.fromAddress ?? "";
    return message;
  },
  fromSDK(object: MsgSubmitQueryResponseSDKType): MsgSubmitQueryResponse {
    return {
      chainId: object?.chain_id,
      queryId: object?.query_id,
      result: object?.result,
      proofOps: object.proof_ops ? ProofOps.fromSDK(object.proof_ops) : undefined,
      height: object?.height,
      fromAddress: object?.from_address
    };
  },
  toSDK(message: MsgSubmitQueryResponse): MsgSubmitQueryResponseSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.query_id = message.queryId;
    obj.result = message.result;
    message.proofOps !== undefined && (obj.proof_ops = message.proofOps ? ProofOps.toSDK(message.proofOps) : undefined);
    obj.height = message.height;
    obj.from_address = message.fromAddress;
    return obj;
  },
  fromAmino(object: MsgSubmitQueryResponseAmino): MsgSubmitQueryResponse {
    return {
      chainId: object.chain_id,
      queryId: object.query_id,
      result: object.result,
      proofOps: object?.proof_ops ? ProofOps.fromAmino(object.proof_ops) : undefined,
      height: Long.fromString(object.height),
      fromAddress: object.from_address
    };
  },
  toAmino(message: MsgSubmitQueryResponse): MsgSubmitQueryResponseAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.query_id = message.queryId;
    obj.result = message.result;
    obj.proof_ops = message.proofOps ? ProofOps.toAmino(message.proofOps) : undefined;
    obj.height = message.height ? message.height.toString() : undefined;
    obj.from_address = message.fromAddress;
    return obj;
  },
  fromAminoMsg(object: MsgSubmitQueryResponseAminoMsg): MsgSubmitQueryResponse {
    return MsgSubmitQueryResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSubmitQueryResponseProtoMsg): MsgSubmitQueryResponse {
    return MsgSubmitQueryResponse.decode(message.value);
  },
  toProto(message: MsgSubmitQueryResponse): Uint8Array {
    return MsgSubmitQueryResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSubmitQueryResponse): MsgSubmitQueryResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponse",
      value: MsgSubmitQueryResponse.encode(message).finish()
    };
  }
};
function createBaseMsgSubmitQueryResponseResponse(): MsgSubmitQueryResponseResponse {
  return {};
}
export const MsgSubmitQueryResponseResponse = {
  typeUrl: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse",
  encode(_: MsgSubmitQueryResponseResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSubmitQueryResponseResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSubmitQueryResponseResponse();
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
  fromJSON(_: any): MsgSubmitQueryResponseResponse {
    const obj = createBaseMsgSubmitQueryResponseResponse();
    return obj;
  },
  toJSON(_: MsgSubmitQueryResponseResponse): unknown {
    const obj: any = {};
    return obj;
  },
  fromPartial(_: DeepPartial<MsgSubmitQueryResponseResponse>): MsgSubmitQueryResponseResponse {
    const message = createBaseMsgSubmitQueryResponseResponse();
    return message;
  },
  fromSDK(_: MsgSubmitQueryResponseResponseSDKType): MsgSubmitQueryResponseResponse {
    return {};
  },
  toSDK(_: MsgSubmitQueryResponseResponse): MsgSubmitQueryResponseResponseSDKType {
    const obj: any = {};
    return obj;
  },
  fromAmino(_: MsgSubmitQueryResponseResponseAmino): MsgSubmitQueryResponseResponse {
    return {};
  },
  toAmino(_: MsgSubmitQueryResponseResponse): MsgSubmitQueryResponseResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgSubmitQueryResponseResponseAminoMsg): MsgSubmitQueryResponseResponse {
    return MsgSubmitQueryResponseResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSubmitQueryResponseResponseProtoMsg): MsgSubmitQueryResponseResponse {
    return MsgSubmitQueryResponseResponse.decode(message.value);
  },
  toProto(message: MsgSubmitQueryResponseResponse): Uint8Array {
    return MsgSubmitQueryResponseResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSubmitQueryResponseResponse): MsgSubmitQueryResponseResponseProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainquery.v1.MsgSubmitQueryResponseResponse",
      value: MsgSubmitQueryResponseResponse.encode(message).finish()
    };
  }
};