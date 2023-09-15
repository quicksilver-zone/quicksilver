import { Long, isSet, bytesFromBase64, base64FromBytes, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.interchainquery.v1";
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
export interface QueryProtoMsg {
  typeUrl: "/quicksilver.interchainquery.v1.Query";
  value: Uint8Array;
}
export interface QueryAmino {
  id: string;
  connection_id: string;
  chain_id: string;
  query_type: string;
  request: Uint8Array;
  /** change these to uint64 in v0.5.0 */
  period: string;
  last_height: string;
  callback_id: string;
  ttl: string;
  last_emission: string;
}
export interface QueryAminoMsg {
  type: "/quicksilver.interchainquery.v1.Query";
  value: QueryAmino;
}
export interface QuerySDKType {
  id: string;
  connection_id: string;
  chain_id: string;
  query_type: string;
  request: Uint8Array;
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
export interface DataPointProtoMsg {
  typeUrl: "/quicksilver.interchainquery.v1.DataPoint";
  value: Uint8Array;
}
export interface DataPointAmino {
  id: string;
  /** change these to uint64 in v0.5.0 */
  remote_height: string;
  local_height: string;
  value: Uint8Array;
}
export interface DataPointAminoMsg {
  type: "/quicksilver.interchainquery.v1.DataPoint";
  value: DataPointAmino;
}
export interface DataPointSDKType {
  id: string;
  remote_height: string;
  local_height: string;
  value: Uint8Array;
}
function createBaseQuery(): Query {
  return {
    id: "",
    connectionId: "",
    chainId: "",
    queryType: "",
    request: new Uint8Array(),
    period: "",
    lastHeight: "",
    callbackId: "",
    ttl: Long.UZERO,
    lastEmission: ""
  };
}
export const Query = {
  typeUrl: "/quicksilver.interchainquery.v1.Query",
  encode(message: Query, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.connectionId !== "") {
      writer.uint32(18).string(message.connectionId);
    }
    if (message.chainId !== "") {
      writer.uint32(26).string(message.chainId);
    }
    if (message.queryType !== "") {
      writer.uint32(34).string(message.queryType);
    }
    if (message.request.length !== 0) {
      writer.uint32(42).bytes(message.request);
    }
    if (message.period !== "") {
      writer.uint32(50).string(message.period);
    }
    if (message.lastHeight !== "") {
      writer.uint32(58).string(message.lastHeight);
    }
    if (message.callbackId !== "") {
      writer.uint32(66).string(message.callbackId);
    }
    if (!message.ttl.isZero()) {
      writer.uint32(72).uint64(message.ttl);
    }
    if (message.lastEmission !== "") {
      writer.uint32(82).string(message.lastEmission);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Query {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.connectionId = reader.string();
          break;
        case 3:
          message.chainId = reader.string();
          break;
        case 4:
          message.queryType = reader.string();
          break;
        case 5:
          message.request = reader.bytes();
          break;
        case 6:
          message.period = reader.string();
          break;
        case 7:
          message.lastHeight = reader.string();
          break;
        case 8:
          message.callbackId = reader.string();
          break;
        case 9:
          message.ttl = (reader.uint64() as Long);
          break;
        case 10:
          message.lastEmission = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Query {
    const obj = createBaseQuery();
    if (isSet(object.id)) obj.id = String(object.id);
    if (isSet(object.connectionId)) obj.connectionId = String(object.connectionId);
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.queryType)) obj.queryType = String(object.queryType);
    if (isSet(object.request)) obj.request = bytesFromBase64(object.request);
    if (isSet(object.period)) obj.period = String(object.period);
    if (isSet(object.lastHeight)) obj.lastHeight = String(object.lastHeight);
    if (isSet(object.callbackId)) obj.callbackId = String(object.callbackId);
    if (isSet(object.ttl)) obj.ttl = Long.fromValue(object.ttl);
    if (isSet(object.lastEmission)) obj.lastEmission = String(object.lastEmission);
    return obj;
  },
  toJSON(message: Query): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.connectionId !== undefined && (obj.connectionId = message.connectionId);
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.queryType !== undefined && (obj.queryType = message.queryType);
    message.request !== undefined && (obj.request = base64FromBytes(message.request !== undefined ? message.request : new Uint8Array()));
    message.period !== undefined && (obj.period = message.period);
    message.lastHeight !== undefined && (obj.lastHeight = message.lastHeight);
    message.callbackId !== undefined && (obj.callbackId = message.callbackId);
    message.ttl !== undefined && (obj.ttl = (message.ttl || Long.UZERO).toString());
    message.lastEmission !== undefined && (obj.lastEmission = message.lastEmission);
    return obj;
  },
  fromPartial(object: DeepPartial<Query>): Query {
    const message = createBaseQuery();
    message.id = object.id ?? "";
    message.connectionId = object.connectionId ?? "";
    message.chainId = object.chainId ?? "";
    message.queryType = object.queryType ?? "";
    message.request = object.request ?? new Uint8Array();
    message.period = object.period ?? "";
    message.lastHeight = object.lastHeight ?? "";
    message.callbackId = object.callbackId ?? "";
    if (object.ttl !== undefined && object.ttl !== null) {
      message.ttl = Long.fromValue(object.ttl);
    }
    message.lastEmission = object.lastEmission ?? "";
    return message;
  },
  fromSDK(object: QuerySDKType): Query {
    return {
      id: object?.id,
      connectionId: object?.connection_id,
      chainId: object?.chain_id,
      queryType: object?.query_type,
      request: object?.request,
      period: object?.period,
      lastHeight: object?.last_height,
      callbackId: object?.callback_id,
      ttl: object?.ttl,
      lastEmission: object?.last_emission
    };
  },
  toSDK(message: Query): QuerySDKType {
    const obj: any = {};
    obj.id = message.id;
    obj.connection_id = message.connectionId;
    obj.chain_id = message.chainId;
    obj.query_type = message.queryType;
    obj.request = message.request;
    obj.period = message.period;
    obj.last_height = message.lastHeight;
    obj.callback_id = message.callbackId;
    obj.ttl = message.ttl;
    obj.last_emission = message.lastEmission;
    return obj;
  },
  fromAmino(object: QueryAmino): Query {
    return {
      id: object.id,
      connectionId: object.connection_id,
      chainId: object.chain_id,
      queryType: object.query_type,
      request: object.request,
      period: object.period,
      lastHeight: object.last_height,
      callbackId: object.callback_id,
      ttl: Long.fromString(object.ttl),
      lastEmission: object.last_emission
    };
  },
  toAmino(message: Query): QueryAmino {
    const obj: any = {};
    obj.id = message.id;
    obj.connection_id = message.connectionId;
    obj.chain_id = message.chainId;
    obj.query_type = message.queryType;
    obj.request = message.request;
    obj.period = message.period;
    obj.last_height = message.lastHeight;
    obj.callback_id = message.callbackId;
    obj.ttl = message.ttl ? message.ttl.toString() : undefined;
    obj.last_emission = message.lastEmission;
    return obj;
  },
  fromAminoMsg(object: QueryAminoMsg): Query {
    return Query.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryProtoMsg): Query {
    return Query.decode(message.value);
  },
  toProto(message: Query): Uint8Array {
    return Query.encode(message).finish();
  },
  toProtoMsg(message: Query): QueryProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainquery.v1.Query",
      value: Query.encode(message).finish()
    };
  }
};
function createBaseDataPoint(): DataPoint {
  return {
    id: "",
    remoteHeight: "",
    localHeight: "",
    value: new Uint8Array()
  };
}
export const DataPoint = {
  typeUrl: "/quicksilver.interchainquery.v1.DataPoint",
  encode(message: DataPoint, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.remoteHeight !== "") {
      writer.uint32(18).string(message.remoteHeight);
    }
    if (message.localHeight !== "") {
      writer.uint32(26).string(message.localHeight);
    }
    if (message.value.length !== 0) {
      writer.uint32(34).bytes(message.value);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): DataPoint {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataPoint();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.remoteHeight = reader.string();
          break;
        case 3:
          message.localHeight = reader.string();
          break;
        case 4:
          message.value = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): DataPoint {
    const obj = createBaseDataPoint();
    if (isSet(object.id)) obj.id = String(object.id);
    if (isSet(object.remoteHeight)) obj.remoteHeight = String(object.remoteHeight);
    if (isSet(object.localHeight)) obj.localHeight = String(object.localHeight);
    if (isSet(object.value)) obj.value = bytesFromBase64(object.value);
    return obj;
  },
  toJSON(message: DataPoint): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.remoteHeight !== undefined && (obj.remoteHeight = message.remoteHeight);
    message.localHeight !== undefined && (obj.localHeight = message.localHeight);
    message.value !== undefined && (obj.value = base64FromBytes(message.value !== undefined ? message.value : new Uint8Array()));
    return obj;
  },
  fromPartial(object: DeepPartial<DataPoint>): DataPoint {
    const message = createBaseDataPoint();
    message.id = object.id ?? "";
    message.remoteHeight = object.remoteHeight ?? "";
    message.localHeight = object.localHeight ?? "";
    message.value = object.value ?? new Uint8Array();
    return message;
  },
  fromSDK(object: DataPointSDKType): DataPoint {
    return {
      id: object?.id,
      remoteHeight: object?.remote_height,
      localHeight: object?.local_height,
      value: object?.value
    };
  },
  toSDK(message: DataPoint): DataPointSDKType {
    const obj: any = {};
    obj.id = message.id;
    obj.remote_height = message.remoteHeight;
    obj.local_height = message.localHeight;
    obj.value = message.value;
    return obj;
  },
  fromAmino(object: DataPointAmino): DataPoint {
    return {
      id: object.id,
      remoteHeight: object.remote_height,
      localHeight: object.local_height,
      value: object.value
    };
  },
  toAmino(message: DataPoint): DataPointAmino {
    const obj: any = {};
    obj.id = message.id;
    obj.remote_height = message.remoteHeight;
    obj.local_height = message.localHeight;
    obj.value = message.value;
    return obj;
  },
  fromAminoMsg(object: DataPointAminoMsg): DataPoint {
    return DataPoint.fromAmino(object.value);
  },
  fromProtoMsg(message: DataPointProtoMsg): DataPoint {
    return DataPoint.decode(message.value);
  },
  toProto(message: DataPoint): Uint8Array {
    return DataPoint.encode(message).finish();
  },
  toProtoMsg(message: DataPoint): DataPointProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainquery.v1.DataPoint",
      value: DataPoint.encode(message).finish()
    };
  }
};