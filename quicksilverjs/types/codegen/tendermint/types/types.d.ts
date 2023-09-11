import { Proof, ProofSDKType } from "../crypto/proof";
import { Consensus, ConsensusSDKType } from "../version/types";
import { Timestamp, TimestampSDKType } from "../../google/protobuf/timestamp";
import { ValidatorSet, ValidatorSetSDKType } from "./validator";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../helpers";
/** BlockIdFlag indicates which BlcokID the signature is for */
export declare enum BlockIDFlag {
    BLOCK_ID_FLAG_UNKNOWN = 0,
    BLOCK_ID_FLAG_ABSENT = 1,
    BLOCK_ID_FLAG_COMMIT = 2,
    BLOCK_ID_FLAG_NIL = 3,
    UNRECOGNIZED = -1
}
/** BlockIdFlag indicates which BlcokID the signature is for */
export declare enum BlockIDFlagSDKType {
    BLOCK_ID_FLAG_UNKNOWN = 0,
    BLOCK_ID_FLAG_ABSENT = 1,
    BLOCK_ID_FLAG_COMMIT = 2,
    BLOCK_ID_FLAG_NIL = 3,
    UNRECOGNIZED = -1
}
export declare function blockIDFlagFromJSON(object: any): BlockIDFlag;
export declare function blockIDFlagToJSON(object: BlockIDFlag): string;
/** SignedMsgType is a type of signed message in the consensus. */
export declare enum SignedMsgType {
    SIGNED_MSG_TYPE_UNKNOWN = 0,
    /** SIGNED_MSG_TYPE_PREVOTE - Votes */
    SIGNED_MSG_TYPE_PREVOTE = 1,
    SIGNED_MSG_TYPE_PRECOMMIT = 2,
    /** SIGNED_MSG_TYPE_PROPOSAL - Proposals */
    SIGNED_MSG_TYPE_PROPOSAL = 32,
    UNRECOGNIZED = -1
}
/** SignedMsgType is a type of signed message in the consensus. */
export declare enum SignedMsgTypeSDKType {
    SIGNED_MSG_TYPE_UNKNOWN = 0,
    /** SIGNED_MSG_TYPE_PREVOTE - Votes */
    SIGNED_MSG_TYPE_PREVOTE = 1,
    SIGNED_MSG_TYPE_PRECOMMIT = 2,
    /** SIGNED_MSG_TYPE_PROPOSAL - Proposals */
    SIGNED_MSG_TYPE_PROPOSAL = 32,
    UNRECOGNIZED = -1
}
export declare function signedMsgTypeFromJSON(object: any): SignedMsgType;
export declare function signedMsgTypeToJSON(object: SignedMsgType): string;
/** PartsetHeader */
export interface PartSetHeader {
    total: number;
    hash: Uint8Array;
}
/** PartsetHeader */
export interface PartSetHeaderSDKType {
    total: number;
    hash: Uint8Array;
}
export interface Part {
    index: number;
    bytes: Uint8Array;
    proof?: Proof;
}
export interface PartSDKType {
    index: number;
    bytes: Uint8Array;
    proof?: ProofSDKType;
}
/** BlockID */
export interface BlockID {
    hash: Uint8Array;
    partSetHeader?: PartSetHeader;
}
/** BlockID */
export interface BlockIDSDKType {
    hash: Uint8Array;
    part_set_header?: PartSetHeaderSDKType;
}
/** Header defines the structure of a Tendermint block header. */
export interface Header {
    /** basic block info */
    version?: Consensus;
    chainId: string;
    height: Long;
    time?: Timestamp;
    /** prev block info */
    lastBlockId?: BlockID;
    /** hashes of block data */
    lastCommitHash: Uint8Array;
    dataHash: Uint8Array;
    /** hashes from the app output from the prev block */
    validatorsHash: Uint8Array;
    /** validators for the next block */
    nextValidatorsHash: Uint8Array;
    /** consensus params for current block */
    consensusHash: Uint8Array;
    /** state after txs from the previous block */
    appHash: Uint8Array;
    lastResultsHash: Uint8Array;
    /** consensus info */
    evidenceHash: Uint8Array;
    /** original proposer of the block */
    proposerAddress: Uint8Array;
}
/** Header defines the structure of a Tendermint block header. */
export interface HeaderSDKType {
    /** basic block info */
    version?: ConsensusSDKType;
    chain_id: string;
    height: Long;
    time?: TimestampSDKType;
    /** prev block info */
    last_block_id?: BlockIDSDKType;
    /** hashes of block data */
    last_commit_hash: Uint8Array;
    data_hash: Uint8Array;
    /** hashes from the app output from the prev block */
    validators_hash: Uint8Array;
    /** validators for the next block */
    next_validators_hash: Uint8Array;
    /** consensus params for current block */
    consensus_hash: Uint8Array;
    /** state after txs from the previous block */
    app_hash: Uint8Array;
    last_results_hash: Uint8Array;
    /** consensus info */
    evidence_hash: Uint8Array;
    /** original proposer of the block */
    proposer_address: Uint8Array;
}
/** Data contains the set of transactions included in the block */
export interface Data {
    /**
     * Txs that will be applied by state @ block.Height+1.
     * NOTE: not all txs here are valid.  We're just agreeing on the order first.
     * This means that block.AppHash does not include these txs.
     */
    txs: Uint8Array[];
}
/** Data contains the set of transactions included in the block */
export interface DataSDKType {
    /**
     * Txs that will be applied by state @ block.Height+1.
     * NOTE: not all txs here are valid.  We're just agreeing on the order first.
     * This means that block.AppHash does not include these txs.
     */
    txs: Uint8Array[];
}
/**
 * Vote represents a prevote, precommit, or commit vote from validators for
 * consensus.
 */
export interface Vote {
    type: SignedMsgType;
    height: Long;
    round: number;
    blockId?: BlockID;
    timestamp?: Timestamp;
    validatorAddress: Uint8Array;
    validatorIndex: number;
    signature: Uint8Array;
}
/**
 * Vote represents a prevote, precommit, or commit vote from validators for
 * consensus.
 */
export interface VoteSDKType {
    type: SignedMsgTypeSDKType;
    height: Long;
    round: number;
    block_id?: BlockIDSDKType;
    timestamp?: TimestampSDKType;
    validator_address: Uint8Array;
    validator_index: number;
    signature: Uint8Array;
}
/** Commit contains the evidence that a block was committed by a set of validators. */
export interface Commit {
    height: Long;
    round: number;
    blockId?: BlockID;
    signatures: CommitSig[];
}
/** Commit contains the evidence that a block was committed by a set of validators. */
export interface CommitSDKType {
    height: Long;
    round: number;
    block_id?: BlockIDSDKType;
    signatures: CommitSigSDKType[];
}
/** CommitSig is a part of the Vote included in a Commit. */
export interface CommitSig {
    blockIdFlag: BlockIDFlag;
    validatorAddress: Uint8Array;
    timestamp?: Timestamp;
    signature: Uint8Array;
}
/** CommitSig is a part of the Vote included in a Commit. */
export interface CommitSigSDKType {
    block_id_flag: BlockIDFlagSDKType;
    validator_address: Uint8Array;
    timestamp?: TimestampSDKType;
    signature: Uint8Array;
}
export interface Proposal {
    type: SignedMsgType;
    height: Long;
    round: number;
    polRound: number;
    blockId?: BlockID;
    timestamp?: Timestamp;
    signature: Uint8Array;
}
export interface ProposalSDKType {
    type: SignedMsgTypeSDKType;
    height: Long;
    round: number;
    pol_round: number;
    block_id?: BlockIDSDKType;
    timestamp?: TimestampSDKType;
    signature: Uint8Array;
}
export interface SignedHeader {
    header?: Header;
    commit?: Commit;
}
export interface SignedHeaderSDKType {
    header?: HeaderSDKType;
    commit?: CommitSDKType;
}
export interface LightBlock {
    signedHeader?: SignedHeader;
    validatorSet?: ValidatorSet;
}
export interface LightBlockSDKType {
    signed_header?: SignedHeaderSDKType;
    validator_set?: ValidatorSetSDKType;
}
export interface BlockMeta {
    blockId?: BlockID;
    blockSize: Long;
    header?: Header;
    numTxs: Long;
}
export interface BlockMetaSDKType {
    block_id?: BlockIDSDKType;
    block_size: Long;
    header?: HeaderSDKType;
    num_txs: Long;
}
/** TxProof represents a Merkle proof of the presence of a transaction in the Merkle tree. */
export interface TxProof {
    rootHash: Uint8Array;
    data: Uint8Array;
    proof?: Proof;
}
/** TxProof represents a Merkle proof of the presence of a transaction in the Merkle tree. */
export interface TxProofSDKType {
    root_hash: Uint8Array;
    data: Uint8Array;
    proof?: ProofSDKType;
}
export declare const PartSetHeader: {
    encode(message: PartSetHeader, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): PartSetHeader;
    fromJSON(object: any): PartSetHeader;
    toJSON(message: PartSetHeader): unknown;
    fromPartial(object: Partial<PartSetHeader>): PartSetHeader;
};
export declare const Part: {
    encode(message: Part, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Part;
    fromJSON(object: any): Part;
    toJSON(message: Part): unknown;
    fromPartial(object: Partial<Part>): Part;
};
export declare const BlockID: {
    encode(message: BlockID, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BlockID;
    fromJSON(object: any): BlockID;
    toJSON(message: BlockID): unknown;
    fromPartial(object: Partial<BlockID>): BlockID;
};
export declare const Header: {
    encode(message: Header, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Header;
    fromJSON(object: any): Header;
    toJSON(message: Header): unknown;
    fromPartial(object: Partial<Header>): Header;
};
export declare const Data: {
    encode(message: Data, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Data;
    fromJSON(object: any): Data;
    toJSON(message: Data): unknown;
    fromPartial(object: Partial<Data>): Data;
};
export declare const Vote: {
    encode(message: Vote, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Vote;
    fromJSON(object: any): Vote;
    toJSON(message: Vote): unknown;
    fromPartial(object: Partial<Vote>): Vote;
};
export declare const Commit: {
    encode(message: Commit, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Commit;
    fromJSON(object: any): Commit;
    toJSON(message: Commit): unknown;
    fromPartial(object: Partial<Commit>): Commit;
};
export declare const CommitSig: {
    encode(message: CommitSig, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): CommitSig;
    fromJSON(object: any): CommitSig;
    toJSON(message: CommitSig): unknown;
    fromPartial(object: Partial<CommitSig>): CommitSig;
};
export declare const Proposal: {
    encode(message: Proposal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Proposal;
    fromJSON(object: any): Proposal;
    toJSON(message: Proposal): unknown;
    fromPartial(object: Partial<Proposal>): Proposal;
};
export declare const SignedHeader: {
    encode(message: SignedHeader, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SignedHeader;
    fromJSON(object: any): SignedHeader;
    toJSON(message: SignedHeader): unknown;
    fromPartial(object: Partial<SignedHeader>): SignedHeader;
};
export declare const LightBlock: {
    encode(message: LightBlock, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): LightBlock;
    fromJSON(object: any): LightBlock;
    toJSON(message: LightBlock): unknown;
    fromPartial(object: Partial<LightBlock>): LightBlock;
};
export declare const BlockMeta: {
    encode(message: BlockMeta, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BlockMeta;
    fromJSON(object: any): BlockMeta;
    toJSON(message: BlockMeta): unknown;
    fromPartial(object: Partial<BlockMeta>): BlockMeta;
};
export declare const TxProof: {
    encode(message: TxProof, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxProof;
    fromJSON(object: any): TxProof;
    toJSON(message: TxProof): unknown;
    fromPartial(object: Partial<TxProof>): TxProof;
};
