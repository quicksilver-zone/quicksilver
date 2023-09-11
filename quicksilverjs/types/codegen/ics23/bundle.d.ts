import * as _82 from "../proofs";
export declare const ics23: {
    hashOpFromJSON(object: any): _82.HashOp;
    hashOpToJSON(object: _82.HashOp): string;
    lengthOpFromJSON(object: any): _82.LengthOp;
    lengthOpToJSON(object: _82.LengthOp): string;
    HashOp: typeof _82.HashOp;
    HashOpSDKType: typeof _82.HashOpSDKType;
    LengthOp: typeof _82.LengthOp;
    LengthOpSDKType: typeof _82.LengthOpSDKType;
    ExistenceProof: {
        encode(message: _82.ExistenceProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.ExistenceProof;
        fromJSON(object: any): _82.ExistenceProof;
        toJSON(message: _82.ExistenceProof): unknown;
        fromPartial(object: Partial<_82.ExistenceProof>): _82.ExistenceProof;
    };
    NonExistenceProof: {
        encode(message: _82.NonExistenceProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.NonExistenceProof;
        fromJSON(object: any): _82.NonExistenceProof;
        toJSON(message: _82.NonExistenceProof): unknown;
        fromPartial(object: Partial<_82.NonExistenceProof>): _82.NonExistenceProof;
    };
    CommitmentProof: {
        encode(message: _82.CommitmentProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.CommitmentProof;
        fromJSON(object: any): _82.CommitmentProof;
        toJSON(message: _82.CommitmentProof): unknown;
        fromPartial(object: Partial<_82.CommitmentProof>): _82.CommitmentProof;
    };
    LeafOp: {
        encode(message: _82.LeafOp, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.LeafOp;
        fromJSON(object: any): _82.LeafOp;
        toJSON(message: _82.LeafOp): unknown;
        fromPartial(object: Partial<_82.LeafOp>): _82.LeafOp;
    };
    InnerOp: {
        encode(message: _82.InnerOp, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.InnerOp;
        fromJSON(object: any): _82.InnerOp;
        toJSON(message: _82.InnerOp): unknown;
        fromPartial(object: Partial<_82.InnerOp>): _82.InnerOp;
    };
    ProofSpec: {
        encode(message: _82.ProofSpec, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.ProofSpec;
        fromJSON(object: any): _82.ProofSpec;
        toJSON(message: _82.ProofSpec): unknown;
        fromPartial(object: Partial<_82.ProofSpec>): _82.ProofSpec;
    };
    InnerSpec: {
        encode(message: _82.InnerSpec, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.InnerSpec;
        fromJSON(object: any): _82.InnerSpec;
        toJSON(message: _82.InnerSpec): unknown;
        fromPartial(object: Partial<_82.InnerSpec>): _82.InnerSpec;
    };
    BatchProof: {
        encode(message: _82.BatchProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.BatchProof;
        fromJSON(object: any): _82.BatchProof;
        toJSON(message: _82.BatchProof): unknown;
        fromPartial(object: Partial<_82.BatchProof>): _82.BatchProof;
    };
    BatchEntry: {
        encode(message: _82.BatchEntry, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.BatchEntry;
        fromJSON(object: any): _82.BatchEntry;
        toJSON(message: _82.BatchEntry): unknown;
        fromPartial(object: Partial<_82.BatchEntry>): _82.BatchEntry;
    };
    CompressedBatchProof: {
        encode(message: _82.CompressedBatchProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.CompressedBatchProof;
        fromJSON(object: any): _82.CompressedBatchProof;
        toJSON(message: _82.CompressedBatchProof): unknown;
        fromPartial(object: Partial<_82.CompressedBatchProof>): _82.CompressedBatchProof;
    };
    CompressedBatchEntry: {
        encode(message: _82.CompressedBatchEntry, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.CompressedBatchEntry;
        fromJSON(object: any): _82.CompressedBatchEntry;
        toJSON(message: _82.CompressedBatchEntry): unknown;
        fromPartial(object: Partial<_82.CompressedBatchEntry>): _82.CompressedBatchEntry;
    };
    CompressedExistenceProof: {
        encode(message: _82.CompressedExistenceProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.CompressedExistenceProof;
        fromJSON(object: any): _82.CompressedExistenceProof;
        toJSON(message: _82.CompressedExistenceProof): unknown;
        fromPartial(object: Partial<_82.CompressedExistenceProof>): _82.CompressedExistenceProof;
    };
    CompressedNonExistenceProof: {
        encode(message: _82.CompressedNonExistenceProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
        decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _82.CompressedNonExistenceProof;
        fromJSON(object: any): _82.CompressedNonExistenceProof;
        toJSON(message: _82.CompressedNonExistenceProof): unknown;
        fromPartial(object: Partial<_82.CompressedNonExistenceProof>): _82.CompressedNonExistenceProof;
    };
};
