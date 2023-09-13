import * as _m0 from "protobufjs/minimal";
import { Long, isSet, bytesFromBase64, base64FromBytes } from "../../../../helpers";
/** Snapshot contains Tendermint state sync snapshot info. */

export interface Snapshot {
  height: Long;
  format: number;
  chunks: number;
  hash: Uint8Array;
  metadata?: Metadata;
}
/** Snapshot contains Tendermint state sync snapshot info. */

export interface SnapshotSDKType {
  height: Long;
  format: number;
  chunks: number;
  hash: Uint8Array;
  metadata?: MetadataSDKType;
}
/** Metadata contains SDK-specific snapshot metadata. */

export interface Metadata {
  /** SHA-256 chunk hashes */
  chunkHashes: Uint8Array[];
}
/** Metadata contains SDK-specific snapshot metadata. */

export interface MetadataSDKType {
  /** SHA-256 chunk hashes */
  chunk_hashes: Uint8Array[];
}

function createBaseSnapshot(): Snapshot {
  return {
    height: Long.UZERO,
    format: 0,
    chunks: 0,
    hash: new Uint8Array(),
    metadata: undefined
  };
}

export const Snapshot = {
  encode(message: Snapshot, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.height.isZero()) {
      writer.uint32(8).uint64(message.height);
    }

    if (message.format !== 0) {
      writer.uint32(16).uint32(message.format);
    }

    if (message.chunks !== 0) {
      writer.uint32(24).uint32(message.chunks);
    }

    if (message.hash.length !== 0) {
      writer.uint32(34).bytes(message.hash);
    }

    if (message.metadata !== undefined) {
      Metadata.encode(message.metadata, writer.uint32(42).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Snapshot {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSnapshot();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.height = (reader.uint64() as Long);
          break;

        case 2:
          message.format = reader.uint32();
          break;

        case 3:
          message.chunks = reader.uint32();
          break;

        case 4:
          message.hash = reader.bytes();
          break;

        case 5:
          message.metadata = Metadata.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): Snapshot {
    return {
      height: isSet(object.height) ? Long.fromValue(object.height) : Long.UZERO,
      format: isSet(object.format) ? Number(object.format) : 0,
      chunks: isSet(object.chunks) ? Number(object.chunks) : 0,
      hash: isSet(object.hash) ? bytesFromBase64(object.hash) : new Uint8Array(),
      metadata: isSet(object.metadata) ? Metadata.fromJSON(object.metadata) : undefined
    };
  },

  toJSON(message: Snapshot): unknown {
    const obj: any = {};
    message.height !== undefined && (obj.height = (message.height || Long.UZERO).toString());
    message.format !== undefined && (obj.format = Math.round(message.format));
    message.chunks !== undefined && (obj.chunks = Math.round(message.chunks));
    message.hash !== undefined && (obj.hash = base64FromBytes(message.hash !== undefined ? message.hash : new Uint8Array()));
    message.metadata !== undefined && (obj.metadata = message.metadata ? Metadata.toJSON(message.metadata) : undefined);
    return obj;
  },

  fromPartial(object: Partial<Snapshot>): Snapshot {
    const message = createBaseSnapshot();
    message.height = object.height !== undefined && object.height !== null ? Long.fromValue(object.height) : Long.UZERO;
    message.format = object.format ?? 0;
    message.chunks = object.chunks ?? 0;
    message.hash = object.hash ?? new Uint8Array();
    message.metadata = object.metadata !== undefined && object.metadata !== null ? Metadata.fromPartial(object.metadata) : undefined;
    return message;
  }

};

function createBaseMetadata(): Metadata {
  return {
    chunkHashes: []
  };
}

export const Metadata = {
  encode(message: Metadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.chunkHashes) {
      writer.uint32(10).bytes(v!);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Metadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMetadata();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.chunkHashes.push(reader.bytes());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromJSON(object: any): Metadata {
    return {
      chunkHashes: Array.isArray(object?.chunkHashes) ? object.chunkHashes.map((e: any) => bytesFromBase64(e)) : []
    };
  },

  toJSON(message: Metadata): unknown {
    const obj: any = {};

    if (message.chunkHashes) {
      obj.chunkHashes = message.chunkHashes.map(e => base64FromBytes(e !== undefined ? e : new Uint8Array()));
    } else {
      obj.chunkHashes = [];
    }

    return obj;
  },

  fromPartial(object: Partial<Metadata>): Metadata {
    const message = createBaseMetadata();
    message.chunkHashes = object.chunkHashes?.map(e => e) || [];
    return message;
  }

};