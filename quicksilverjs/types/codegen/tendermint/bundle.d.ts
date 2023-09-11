import * as _117 from "./abci/types";
import * as _118 from "./crypto/keys";
import * as _119 from "./crypto/proof";
import * as _120 from "./libs/bits/types";
import * as _121 from "./p2p/types";
import * as _122 from "./types/block";
import * as _123 from "./types/evidence";
import * as _124 from "./types/params";
import * as _125 from "./types/types";
import * as _126 from "./types/validator";
import * as _127 from "./version/types";
export declare namespace tendermint {
    const abci: {
        checkTxTypeFromJSON(object: any): _117.CheckTxType;
        checkTxTypeToJSON(object: _117.CheckTxType): string;
        responseOfferSnapshot_ResultFromJSON(object: any): _117.ResponseOfferSnapshot_Result;
        responseOfferSnapshot_ResultToJSON(object: _117.ResponseOfferSnapshot_Result): string;
        responseApplySnapshotChunk_ResultFromJSON(object: any): _117.ResponseApplySnapshotChunk_Result;
        responseApplySnapshotChunk_ResultToJSON(object: _117.ResponseApplySnapshotChunk_Result): string;
        evidenceTypeFromJSON(object: any): _117.EvidenceType;
        evidenceTypeToJSON(object: _117.EvidenceType): string;
        CheckTxType: typeof _117.CheckTxType;
        CheckTxTypeSDKType: typeof _117.CheckTxTypeSDKType;
        ResponseOfferSnapshot_Result: typeof _117.ResponseOfferSnapshot_Result;
        ResponseOfferSnapshot_ResultSDKType: typeof _117.ResponseOfferSnapshot_ResultSDKType;
        ResponseApplySnapshotChunk_Result: typeof _117.ResponseApplySnapshotChunk_Result;
        ResponseApplySnapshotChunk_ResultSDKType: typeof _117.ResponseApplySnapshotChunk_ResultSDKType;
        EvidenceType: typeof _117.EvidenceType;
        EvidenceTypeSDKType: typeof _117.EvidenceTypeSDKType;
        Request: {
            encode(message: _117.Request, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.Request;
            fromJSON(object: any): _117.Request;
            toJSON(message: _117.Request): unknown;
            fromPartial(object: Partial<_117.Request>): _117.Request;
        };
        RequestEcho: {
            encode(message: _117.RequestEcho, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestEcho;
            fromJSON(object: any): _117.RequestEcho;
            toJSON(message: _117.RequestEcho): unknown;
            fromPartial(object: Partial<_117.RequestEcho>): _117.RequestEcho;
        };
        RequestFlush: {
            encode(_: _117.RequestFlush, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestFlush;
            fromJSON(_: any): _117.RequestFlush;
            toJSON(_: _117.RequestFlush): unknown;
            fromPartial(_: Partial<_117.RequestFlush>): _117.RequestFlush;
        };
        RequestInfo: {
            encode(message: _117.RequestInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestInfo;
            fromJSON(object: any): _117.RequestInfo;
            toJSON(message: _117.RequestInfo): unknown;
            fromPartial(object: Partial<_117.RequestInfo>): _117.RequestInfo;
        };
        RequestSetOption: {
            encode(message: _117.RequestSetOption, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestSetOption;
            fromJSON(object: any): _117.RequestSetOption;
            toJSON(message: _117.RequestSetOption): unknown;
            fromPartial(object: Partial<_117.RequestSetOption>): _117.RequestSetOption;
        };
        RequestInitChain: {
            encode(message: _117.RequestInitChain, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestInitChain;
            fromJSON(object: any): _117.RequestInitChain;
            toJSON(message: _117.RequestInitChain): unknown;
            fromPartial(object: Partial<_117.RequestInitChain>): _117.RequestInitChain;
        };
        RequestQuery: {
            encode(message: _117.RequestQuery, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestQuery;
            fromJSON(object: any): _117.RequestQuery;
            toJSON(message: _117.RequestQuery): unknown;
            fromPartial(object: Partial<_117.RequestQuery>): _117.RequestQuery;
        };
        RequestBeginBlock: {
            encode(message: _117.RequestBeginBlock, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestBeginBlock;
            fromJSON(object: any): _117.RequestBeginBlock;
            toJSON(message: _117.RequestBeginBlock): unknown;
            fromPartial(object: Partial<_117.RequestBeginBlock>): _117.RequestBeginBlock;
        };
        RequestCheckTx: {
            encode(message: _117.RequestCheckTx, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestCheckTx;
            fromJSON(object: any): _117.RequestCheckTx;
            toJSON(message: _117.RequestCheckTx): unknown;
            fromPartial(object: Partial<_117.RequestCheckTx>): _117.RequestCheckTx;
        };
        RequestDeliverTx: {
            encode(message: _117.RequestDeliverTx, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestDeliverTx;
            fromJSON(object: any): _117.RequestDeliverTx;
            toJSON(message: _117.RequestDeliverTx): unknown;
            fromPartial(object: Partial<_117.RequestDeliverTx>): _117.RequestDeliverTx;
        };
        RequestEndBlock: {
            encode(message: _117.RequestEndBlock, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestEndBlock;
            fromJSON(object: any): _117.RequestEndBlock;
            toJSON(message: _117.RequestEndBlock): unknown;
            fromPartial(object: Partial<_117.RequestEndBlock>): _117.RequestEndBlock;
        };
        RequestCommit: {
            encode(_: _117.RequestCommit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestCommit;
            fromJSON(_: any): _117.RequestCommit;
            toJSON(_: _117.RequestCommit): unknown;
            fromPartial(_: Partial<_117.RequestCommit>): _117.RequestCommit;
        };
        RequestListSnapshots: {
            encode(_: _117.RequestListSnapshots, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestListSnapshots;
            fromJSON(_: any): _117.RequestListSnapshots;
            toJSON(_: _117.RequestListSnapshots): unknown;
            fromPartial(_: Partial<_117.RequestListSnapshots>): _117.RequestListSnapshots;
        };
        RequestOfferSnapshot: {
            encode(message: _117.RequestOfferSnapshot, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestOfferSnapshot;
            fromJSON(object: any): _117.RequestOfferSnapshot;
            toJSON(message: _117.RequestOfferSnapshot): unknown;
            fromPartial(object: Partial<_117.RequestOfferSnapshot>): _117.RequestOfferSnapshot;
        };
        RequestLoadSnapshotChunk: {
            encode(message: _117.RequestLoadSnapshotChunk, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestLoadSnapshotChunk;
            fromJSON(object: any): _117.RequestLoadSnapshotChunk;
            toJSON(message: _117.RequestLoadSnapshotChunk): unknown;
            fromPartial(object: Partial<_117.RequestLoadSnapshotChunk>): _117.RequestLoadSnapshotChunk;
        };
        RequestApplySnapshotChunk: {
            encode(message: _117.RequestApplySnapshotChunk, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.RequestApplySnapshotChunk;
            fromJSON(object: any): _117.RequestApplySnapshotChunk;
            toJSON(message: _117.RequestApplySnapshotChunk): unknown;
            fromPartial(object: Partial<_117.RequestApplySnapshotChunk>): _117.RequestApplySnapshotChunk;
        };
        Response: {
            encode(message: _117.Response, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.Response;
            fromJSON(object: any): _117.Response;
            toJSON(message: _117.Response): unknown;
            fromPartial(object: Partial<_117.Response>): _117.Response;
        };
        ResponseException: {
            encode(message: _117.ResponseException, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseException;
            fromJSON(object: any): _117.ResponseException;
            toJSON(message: _117.ResponseException): unknown;
            fromPartial(object: Partial<_117.ResponseException>): _117.ResponseException;
        };
        ResponseEcho: {
            encode(message: _117.ResponseEcho, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseEcho;
            fromJSON(object: any): _117.ResponseEcho;
            toJSON(message: _117.ResponseEcho): unknown;
            fromPartial(object: Partial<_117.ResponseEcho>): _117.ResponseEcho;
        };
        ResponseFlush: {
            encode(_: _117.ResponseFlush, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseFlush;
            fromJSON(_: any): _117.ResponseFlush;
            toJSON(_: _117.ResponseFlush): unknown;
            fromPartial(_: Partial<_117.ResponseFlush>): _117.ResponseFlush;
        };
        ResponseInfo: {
            encode(message: _117.ResponseInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseInfo;
            fromJSON(object: any): _117.ResponseInfo;
            toJSON(message: _117.ResponseInfo): unknown;
            fromPartial(object: Partial<_117.ResponseInfo>): _117.ResponseInfo;
        };
        ResponseSetOption: {
            encode(message: _117.ResponseSetOption, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseSetOption;
            fromJSON(object: any): _117.ResponseSetOption;
            toJSON(message: _117.ResponseSetOption): unknown;
            fromPartial(object: Partial<_117.ResponseSetOption>): _117.ResponseSetOption;
        };
        ResponseInitChain: {
            encode(message: _117.ResponseInitChain, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseInitChain;
            fromJSON(object: any): _117.ResponseInitChain;
            toJSON(message: _117.ResponseInitChain): unknown;
            fromPartial(object: Partial<_117.ResponseInitChain>): _117.ResponseInitChain;
        };
        ResponseQuery: {
            encode(message: _117.ResponseQuery, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseQuery;
            fromJSON(object: any): _117.ResponseQuery;
            toJSON(message: _117.ResponseQuery): unknown;
            fromPartial(object: Partial<_117.ResponseQuery>): _117.ResponseQuery;
        };
        ResponseBeginBlock: {
            encode(message: _117.ResponseBeginBlock, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseBeginBlock;
            fromJSON(object: any): _117.ResponseBeginBlock;
            toJSON(message: _117.ResponseBeginBlock): unknown;
            fromPartial(object: Partial<_117.ResponseBeginBlock>): _117.ResponseBeginBlock;
        };
        ResponseCheckTx: {
            encode(message: _117.ResponseCheckTx, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseCheckTx;
            fromJSON(object: any): _117.ResponseCheckTx;
            toJSON(message: _117.ResponseCheckTx): unknown;
            fromPartial(object: Partial<_117.ResponseCheckTx>): _117.ResponseCheckTx;
        };
        ResponseDeliverTx: {
            encode(message: _117.ResponseDeliverTx, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseDeliverTx;
            fromJSON(object: any): _117.ResponseDeliverTx;
            toJSON(message: _117.ResponseDeliverTx): unknown;
            fromPartial(object: Partial<_117.ResponseDeliverTx>): _117.ResponseDeliverTx;
        };
        ResponseEndBlock: {
            encode(message: _117.ResponseEndBlock, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseEndBlock;
            fromJSON(object: any): _117.ResponseEndBlock;
            toJSON(message: _117.ResponseEndBlock): unknown;
            fromPartial(object: Partial<_117.ResponseEndBlock>): _117.ResponseEndBlock;
        };
        ResponseCommit: {
            encode(message: _117.ResponseCommit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseCommit;
            fromJSON(object: any): _117.ResponseCommit;
            toJSON(message: _117.ResponseCommit): unknown;
            fromPartial(object: Partial<_117.ResponseCommit>): _117.ResponseCommit;
        };
        ResponseListSnapshots: {
            encode(message: _117.ResponseListSnapshots, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseListSnapshots;
            fromJSON(object: any): _117.ResponseListSnapshots;
            toJSON(message: _117.ResponseListSnapshots): unknown;
            fromPartial(object: Partial<_117.ResponseListSnapshots>): _117.ResponseListSnapshots;
        };
        ResponseOfferSnapshot: {
            encode(message: _117.ResponseOfferSnapshot, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseOfferSnapshot;
            fromJSON(object: any): _117.ResponseOfferSnapshot;
            toJSON(message: _117.ResponseOfferSnapshot): unknown;
            fromPartial(object: Partial<_117.ResponseOfferSnapshot>): _117.ResponseOfferSnapshot;
        };
        ResponseLoadSnapshotChunk: {
            encode(message: _117.ResponseLoadSnapshotChunk, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseLoadSnapshotChunk;
            fromJSON(object: any): _117.ResponseLoadSnapshotChunk;
            toJSON(message: _117.ResponseLoadSnapshotChunk): unknown;
            fromPartial(object: Partial<_117.ResponseLoadSnapshotChunk>): _117.ResponseLoadSnapshotChunk;
        };
        ResponseApplySnapshotChunk: {
            encode(message: _117.ResponseApplySnapshotChunk, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ResponseApplySnapshotChunk;
            fromJSON(object: any): _117.ResponseApplySnapshotChunk;
            toJSON(message: _117.ResponseApplySnapshotChunk): unknown;
            fromPartial(object: Partial<_117.ResponseApplySnapshotChunk>): _117.ResponseApplySnapshotChunk;
        };
        ConsensusParams: {
            encode(message: _117.ConsensusParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ConsensusParams;
            fromJSON(object: any): _117.ConsensusParams;
            toJSON(message: _117.ConsensusParams): unknown;
            fromPartial(object: Partial<_117.ConsensusParams>): _117.ConsensusParams;
        };
        BlockParams: {
            encode(message: _117.BlockParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.BlockParams;
            fromJSON(object: any): _117.BlockParams;
            toJSON(message: _117.BlockParams): unknown;
            fromPartial(object: Partial<_117.BlockParams>): _117.BlockParams;
        };
        LastCommitInfo: {
            encode(message: _117.LastCommitInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.LastCommitInfo;
            fromJSON(object: any): _117.LastCommitInfo;
            toJSON(message: _117.LastCommitInfo): unknown;
            fromPartial(object: Partial<_117.LastCommitInfo>): _117.LastCommitInfo;
        };
        Event: {
            encode(message: _117.Event, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.Event;
            fromJSON(object: any): _117.Event;
            toJSON(message: _117.Event): unknown;
            fromPartial(object: Partial<_117.Event>): _117.Event;
        };
        EventAttribute: {
            encode(message: _117.EventAttribute, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.EventAttribute;
            fromJSON(object: any): _117.EventAttribute;
            toJSON(message: _117.EventAttribute): unknown;
            fromPartial(object: Partial<_117.EventAttribute>): _117.EventAttribute;
        };
        TxResult: {
            encode(message: _117.TxResult, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.TxResult;
            fromJSON(object: any): _117.TxResult;
            toJSON(message: _117.TxResult): unknown;
            fromPartial(object: Partial<_117.TxResult>): _117.TxResult;
        };
        Validator: {
            encode(message: _117.Validator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.Validator;
            fromJSON(object: any): _117.Validator;
            toJSON(message: _117.Validator): unknown;
            fromPartial(object: Partial<_117.Validator>): _117.Validator;
        };
        ValidatorUpdate: {
            encode(message: _117.ValidatorUpdate, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.ValidatorUpdate;
            fromJSON(object: any): _117.ValidatorUpdate;
            toJSON(message: _117.ValidatorUpdate): unknown;
            fromPartial(object: Partial<_117.ValidatorUpdate>): _117.ValidatorUpdate;
        };
        VoteInfo: {
            encode(message: _117.VoteInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.VoteInfo;
            fromJSON(object: any): _117.VoteInfo;
            toJSON(message: _117.VoteInfo): unknown;
            fromPartial(object: Partial<_117.VoteInfo>): _117.VoteInfo;
        };
        Evidence: {
            encode(message: _117.Evidence, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.Evidence;
            fromJSON(object: any): _117.Evidence;
            toJSON(message: _117.Evidence): unknown;
            fromPartial(object: Partial<_117.Evidence>): _117.Evidence;
        };
        Snapshot: {
            encode(message: _117.Snapshot, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _117.Snapshot;
            fromJSON(object: any): _117.Snapshot;
            toJSON(message: _117.Snapshot): unknown;
            fromPartial(object: Partial<_117.Snapshot>): _117.Snapshot;
        };
    };
    const crypto: {
        Proof: {
            encode(message: _119.Proof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _119.Proof;
            fromJSON(object: any): _119.Proof;
            toJSON(message: _119.Proof): unknown;
            fromPartial(object: Partial<_119.Proof>): _119.Proof;
        };
        ValueOp: {
            encode(message: _119.ValueOp, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _119.ValueOp;
            fromJSON(object: any): _119.ValueOp;
            toJSON(message: _119.ValueOp): unknown;
            fromPartial(object: Partial<_119.ValueOp>): _119.ValueOp;
        };
        DominoOp: {
            encode(message: _119.DominoOp, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _119.DominoOp;
            fromJSON(object: any): _119.DominoOp;
            toJSON(message: _119.DominoOp): unknown;
            fromPartial(object: Partial<_119.DominoOp>): _119.DominoOp;
        };
        ProofOp: {
            encode(message: _119.ProofOp, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _119.ProofOp;
            fromJSON(object: any): _119.ProofOp;
            toJSON(message: _119.ProofOp): unknown;
            fromPartial(object: Partial<_119.ProofOp>): _119.ProofOp;
        };
        ProofOps: {
            encode(message: _119.ProofOps, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _119.ProofOps;
            fromJSON(object: any): _119.ProofOps;
            toJSON(message: _119.ProofOps): unknown;
            fromPartial(object: Partial<_119.ProofOps>): _119.ProofOps;
        };
        PublicKey: {
            encode(message: _118.PublicKey, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _118.PublicKey;
            fromJSON(object: any): _118.PublicKey;
            toJSON(message: _118.PublicKey): unknown;
            fromPartial(object: Partial<_118.PublicKey>): _118.PublicKey;
        };
    };
    namespace libs {
        const bits: {
            BitArray: {
                encode(message: _120.BitArray, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _120.BitArray;
                fromJSON(object: any): _120.BitArray;
                toJSON(message: _120.BitArray): unknown;
                fromPartial(object: Partial<_120.BitArray>): _120.BitArray;
            };
        };
    }
    const p2p: {
        NetAddress: {
            encode(message: _121.NetAddress, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _121.NetAddress;
            fromJSON(object: any): _121.NetAddress;
            toJSON(message: _121.NetAddress): unknown;
            fromPartial(object: Partial<_121.NetAddress>): _121.NetAddress;
        };
        ProtocolVersion: {
            encode(message: _121.ProtocolVersion, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _121.ProtocolVersion;
            fromJSON(object: any): _121.ProtocolVersion;
            toJSON(message: _121.ProtocolVersion): unknown;
            fromPartial(object: Partial<_121.ProtocolVersion>): _121.ProtocolVersion;
        };
        DefaultNodeInfo: {
            encode(message: _121.DefaultNodeInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _121.DefaultNodeInfo;
            fromJSON(object: any): _121.DefaultNodeInfo;
            toJSON(message: _121.DefaultNodeInfo): unknown;
            fromPartial(object: Partial<_121.DefaultNodeInfo>): _121.DefaultNodeInfo;
        };
        DefaultNodeInfoOther: {
            encode(message: _121.DefaultNodeInfoOther, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _121.DefaultNodeInfoOther;
            fromJSON(object: any): _121.DefaultNodeInfoOther;
            toJSON(message: _121.DefaultNodeInfoOther): unknown;
            fromPartial(object: Partial<_121.DefaultNodeInfoOther>): _121.DefaultNodeInfoOther;
        };
    };
    const types: {
        ValidatorSet: {
            encode(message: _126.ValidatorSet, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _126.ValidatorSet;
            fromJSON(object: any): _126.ValidatorSet;
            toJSON(message: _126.ValidatorSet): unknown;
            fromPartial(object: Partial<_126.ValidatorSet>): _126.ValidatorSet;
        };
        Validator: {
            encode(message: _126.Validator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _126.Validator;
            fromJSON(object: any): _126.Validator;
            toJSON(message: _126.Validator): unknown;
            fromPartial(object: Partial<_126.Validator>): _126.Validator;
        };
        SimpleValidator: {
            encode(message: _126.SimpleValidator, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _126.SimpleValidator;
            fromJSON(object: any): _126.SimpleValidator;
            toJSON(message: _126.SimpleValidator): unknown;
            fromPartial(object: Partial<_126.SimpleValidator>): _126.SimpleValidator;
        };
        blockIDFlagFromJSON(object: any): _125.BlockIDFlag;
        blockIDFlagToJSON(object: _125.BlockIDFlag): string;
        signedMsgTypeFromJSON(object: any): _125.SignedMsgType;
        signedMsgTypeToJSON(object: _125.SignedMsgType): string;
        BlockIDFlag: typeof _125.BlockIDFlag;
        BlockIDFlagSDKType: typeof _125.BlockIDFlagSDKType;
        SignedMsgType: typeof _125.SignedMsgType;
        SignedMsgTypeSDKType: typeof _125.SignedMsgTypeSDKType;
        PartSetHeader: {
            encode(message: _125.PartSetHeader, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.PartSetHeader;
            fromJSON(object: any): _125.PartSetHeader;
            toJSON(message: _125.PartSetHeader): unknown;
            fromPartial(object: Partial<_125.PartSetHeader>): _125.PartSetHeader;
        };
        Part: {
            encode(message: _125.Part, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.Part;
            fromJSON(object: any): _125.Part;
            toJSON(message: _125.Part): unknown;
            fromPartial(object: Partial<_125.Part>): _125.Part;
        };
        BlockID: {
            encode(message: _125.BlockID, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.BlockID;
            fromJSON(object: any): _125.BlockID;
            toJSON(message: _125.BlockID): unknown;
            fromPartial(object: Partial<_125.BlockID>): _125.BlockID;
        };
        Header: {
            encode(message: _125.Header, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.Header;
            fromJSON(object: any): _125.Header;
            toJSON(message: _125.Header): unknown;
            fromPartial(object: Partial<_125.Header>): _125.Header;
        };
        Data: {
            encode(message: _125.Data, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.Data;
            fromJSON(object: any): _125.Data;
            toJSON(message: _125.Data): unknown;
            fromPartial(object: Partial<_125.Data>): _125.Data;
        };
        Vote: {
            encode(message: _125.Vote, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.Vote;
            fromJSON(object: any): _125.Vote;
            toJSON(message: _125.Vote): unknown;
            fromPartial(object: Partial<_125.Vote>): _125.Vote;
        };
        Commit: {
            encode(message: _125.Commit, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.Commit;
            fromJSON(object: any): _125.Commit;
            toJSON(message: _125.Commit): unknown;
            fromPartial(object: Partial<_125.Commit>): _125.Commit;
        };
        CommitSig: {
            encode(message: _125.CommitSig, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.CommitSig;
            fromJSON(object: any): _125.CommitSig;
            toJSON(message: _125.CommitSig): unknown;
            fromPartial(object: Partial<_125.CommitSig>): _125.CommitSig;
        };
        Proposal: {
            encode(message: _125.Proposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.Proposal;
            fromJSON(object: any): _125.Proposal;
            toJSON(message: _125.Proposal): unknown;
            fromPartial(object: Partial<_125.Proposal>): _125.Proposal;
        };
        SignedHeader: {
            encode(message: _125.SignedHeader, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.SignedHeader;
            fromJSON(object: any): _125.SignedHeader;
            toJSON(message: _125.SignedHeader): unknown;
            fromPartial(object: Partial<_125.SignedHeader>): _125.SignedHeader;
        };
        LightBlock: {
            encode(message: _125.LightBlock, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.LightBlock;
            fromJSON(object: any): _125.LightBlock;
            toJSON(message: _125.LightBlock): unknown;
            fromPartial(object: Partial<_125.LightBlock>): _125.LightBlock;
        };
        BlockMeta: {
            encode(message: _125.BlockMeta, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.BlockMeta;
            fromJSON(object: any): _125.BlockMeta;
            toJSON(message: _125.BlockMeta): unknown;
            fromPartial(object: Partial<_125.BlockMeta>): _125.BlockMeta;
        };
        TxProof: {
            encode(message: _125.TxProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _125.TxProof;
            fromJSON(object: any): _125.TxProof;
            toJSON(message: _125.TxProof): unknown;
            fromPartial(object: Partial<_125.TxProof>): _125.TxProof;
        };
        ConsensusParams: {
            encode(message: _124.ConsensusParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _124.ConsensusParams;
            fromJSON(object: any): _124.ConsensusParams;
            toJSON(message: _124.ConsensusParams): unknown;
            fromPartial(object: Partial<_124.ConsensusParams>): _124.ConsensusParams;
        };
        BlockParams: {
            encode(message: _124.BlockParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _124.BlockParams;
            fromJSON(object: any): _124.BlockParams;
            toJSON(message: _124.BlockParams): unknown;
            fromPartial(object: Partial<_124.BlockParams>): _124.BlockParams;
        };
        EvidenceParams: {
            encode(message: _124.EvidenceParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _124.EvidenceParams;
            fromJSON(object: any): _124.EvidenceParams;
            toJSON(message: _124.EvidenceParams): unknown;
            fromPartial(object: Partial<_124.EvidenceParams>): _124.EvidenceParams;
        };
        ValidatorParams: {
            encode(message: _124.ValidatorParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _124.ValidatorParams;
            fromJSON(object: any): _124.ValidatorParams;
            toJSON(message: _124.ValidatorParams): unknown;
            fromPartial(object: Partial<_124.ValidatorParams>): _124.ValidatorParams;
        };
        VersionParams: {
            encode(message: _124.VersionParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _124.VersionParams;
            fromJSON(object: any): _124.VersionParams;
            toJSON(message: _124.VersionParams): unknown;
            fromPartial(object: Partial<_124.VersionParams>): _124.VersionParams;
        };
        HashedParams: {
            encode(message: _124.HashedParams, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _124.HashedParams;
            fromJSON(object: any): _124.HashedParams;
            toJSON(message: _124.HashedParams): unknown;
            fromPartial(object: Partial<_124.HashedParams>): _124.HashedParams;
        };
        Evidence: {
            encode(message: _123.Evidence, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _123.Evidence;
            fromJSON(object: any): _123.Evidence;
            toJSON(message: _123.Evidence): unknown;
            fromPartial(object: Partial<_123.Evidence>): _123.Evidence;
        };
        DuplicateVoteEvidence: {
            encode(message: _123.DuplicateVoteEvidence, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _123.DuplicateVoteEvidence;
            fromJSON(object: any): _123.DuplicateVoteEvidence;
            toJSON(message: _123.DuplicateVoteEvidence): unknown;
            fromPartial(object: Partial<_123.DuplicateVoteEvidence>): _123.DuplicateVoteEvidence;
        };
        LightClientAttackEvidence: {
            encode(message: _123.LightClientAttackEvidence, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _123.LightClientAttackEvidence;
            fromJSON(object: any): _123.LightClientAttackEvidence;
            toJSON(message: _123.LightClientAttackEvidence): unknown;
            fromPartial(object: Partial<_123.LightClientAttackEvidence>): _123.LightClientAttackEvidence;
        };
        EvidenceList: {
            encode(message: _123.EvidenceList, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _123.EvidenceList;
            fromJSON(object: any): _123.EvidenceList;
            toJSON(message: _123.EvidenceList): unknown;
            fromPartial(object: Partial<_123.EvidenceList>): _123.EvidenceList;
        };
        Block: {
            encode(message: _122.Block, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _122.Block;
            fromJSON(object: any): _122.Block;
            toJSON(message: _122.Block): unknown;
            fromPartial(object: Partial<_122.Block>): _122.Block;
        };
    };
    const version: {
        App: {
            encode(message: _127.App, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _127.App;
            fromJSON(object: any): _127.App;
            toJSON(message: _127.App): unknown;
            fromPartial(object: Partial<_127.App>): _127.App;
        };
        Consensus: {
            encode(message: _127.Consensus, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _127.Consensus;
            fromJSON(object: any): _127.Consensus;
            toJSON(message: _127.Consensus): unknown;
            fromPartial(object: Partial<_127.Consensus>): _127.Consensus;
        };
    };
}
