import * as _m0 from "protobufjs/minimal";
import { isSet, Long, bytesFromBase64, base64FromBytes } from "../../../../helpers";
/** SnapshotItem is an item contained in a rootmulti.Store snapshot. */

export interface SnapshotItem {
  store?: SnapshotStoreItem;
  iavl?: SnapshotIAVLItem;
}
/** SnapshotItem is an item contained in a rootmulti.Store snapshot. */

export interface SnapshotItemSDKType {
  store?: SnapshotStoreItemSDKType;
  iavl?: SnapshotIAVLItemSDKType;
}
/** SnapshotStoreItem contains metadata about a snapshotted store. */

export interface SnapshotStoreItem {
  name: string;
}
/** SnapshotStoreItem contains metadata about a snapshotted store. */

export interface SnapshotStoreItemSDKType {
  name: string;
}
/** SnapshotIAVLItem is an exported IAVL node. */

export interface SnapshotIAVLItem {
  key: Uint8Array;
  value: Uint8Array;
  version: Long;
  height: number;
}
/** SnapshotIAVLItem is an exported IAVL node. */

export interface SnapshotIAVLItemSDKType {
  key: Uint8Array;
  value: Uint8Array;
  version: Long;
  height: number;
}

function createBaseSnapshotItem(): SnapshotItem {
  return {
    store: undefined,
    iavl: undefined
  };
}

export const SnapshotItem = {
  encode(message: SnapshotItem, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.store !== undefined) {
      SnapshotStoreItem.encode(message.store, writer.uint32(10).fork()).ldelim();
    }

    if (message.iavl !== undefined) {
      SnapshotIAVLItem.encode(message.iavl, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotItem {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSnapshotItem();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.store = SnapshotStoreItem.decode(reader, reader.uint32());
          break;

        case 2:
          message.iavl = SnapshotIAVLItem.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): SnapshotItem {
    return {
      store: isSet(object.store) ? SnapshotStoreItem.fromJSON(object.store) : undefined,
      iavl: isSet(object.iavl) ? SnapshotIAVLItem.fromJSON(object.iavl) : undefined
    };
  },

  toJSON(message: SnapshotItem): unknown {
    const obj: any = {};
    message.store !== undefined && (obj.store = message.store ? SnapshotStoreItem.toJSON(message.store) : undefined);
    message.iavl !== undefined && (obj.iavl = message.iavl ? SnapshotIAVLItem.toJSON(message.iavl) : undefined);
    return obj;
  },

  fromPartial(object: Partial<SnapshotItem>): SnapshotItem {
    const message = createBaseSnapshotItem();
    message.store = object.store !== undefined && object.store !== null ? SnapshotStoreItem.fromPartial(object.store) : undefined;
    message.iavl = object.iavl !== undefined && object.iavl !== null ? SnapshotIAVLItem.fromPartial(object.iavl) : undefined;
    return message;
  }

};

function createBaseSnapshotStoreItem(): SnapshotStoreItem {
  return {
    name: ""
  };
}

export const SnapshotStoreItem = {
  encode(message: SnapshotStoreItem, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotStoreItem {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSnapshotStoreItem();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): SnapshotStoreItem {
    return {
      name: isSet(object.name) ? String(object.name) : ""
    };
  },

  toJSON(message: SnapshotStoreItem): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  fromPartial(object: Partial<SnapshotStoreItem>): SnapshotStoreItem {
    const message = createBaseSnapshotStoreItem();
    message.name = object.name ?? "";
    return message;
  }

};

function createBaseSnapshotIAVLItem(): SnapshotIAVLItem {
  return {
    key: new Uint8Array(),
    value: new Uint8Array(),
    version: Long.ZERO,
    height: 0
  };
}

export const SnapshotIAVLItem = {
  encode(message: SnapshotIAVLItem, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key.length !== 0) {
      writer.uint32(10).bytes(message.key);
    }

    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
    }

    if (!message.version.isZero()) {
      writer.uint32(24).int64(message.version);
    }

    if (message.height !== 0) {
      writer.uint32(32).int32(message.height);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotIAVLItem {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSnapshotIAVLItem();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.key = reader.bytes();
          break;

        case 2:
          message.value = reader.bytes();
          break;

        case 3:
          message.version = (reader.int64() as Long);
          break;

        case 4:
          message.height = reader.int32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): SnapshotIAVLItem {
    return {
      key: isSet(object.key) ? bytesFromBase64(object.key) : new Uint8Array(),
      value: isSet(object.value) ? bytesFromBase64(object.value) : new Uint8Array(),
      version: isSet(object.version) ? Long.fromValue(object.version) : Long.ZERO,
      height: isSet(object.height) ? Number(object.height) : 0
    };
  },

  toJSON(message: SnapshotIAVLItem): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = base64FromBytes(message.key !== undefined ? message.key : new Uint8Array()));
    message.value !== undefined && (obj.value = base64FromBytes(message.value !== undefined ? message.value : new Uint8Array()));
    message.version !== undefined && (obj.version = (message.version || Long.ZERO).toString());
    message.height !== undefined && (obj.height = Math.round(message.height));
    return obj;
  },

  fromPartial(object: Partial<SnapshotIAVLItem>): SnapshotIAVLItem {
    const message = createBaseSnapshotIAVLItem();
    message.key = object.key ?? new Uint8Array();
    message.value = object.value ?? new Uint8Array();
    message.version = object.version !== undefined && object.version !== null ? Long.fromValue(object.version) : Long.ZERO;
    message.height = object.height ?? 0;
    return message;
  }

};