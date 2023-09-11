import * as _79 from "./core/client/v1/client";
import * as _80 from "./core/commitment/v1/commitment";
import * as _81 from "./lightclients/tendermint/v1/tendermint";
export declare namespace ibc {
    namespace core {
        namespace client {
            const v1: {
                IdentifiedClientState: {
                    encode(message: _79.IdentifiedClientState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _79.IdentifiedClientState;
                    fromJSON(object: any): _79.IdentifiedClientState;
                    toJSON(message: _79.IdentifiedClientState): unknown;
                    fromPartial(object: Partial<_79.IdentifiedClientState>): _79.IdentifiedClientState;
                };
                ConsensusStateWithHeight: {
                    encode(message: _79.ConsensusStateWithHeight, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _79.ConsensusStateWithHeight;
                    fromJSON(object: any): _79.ConsensusStateWithHeight;
                    toJSON(message: _79.ConsensusStateWithHeight): unknown;
                    fromPartial(object: Partial<_79.ConsensusStateWithHeight>): _79.ConsensusStateWithHeight;
                };
                ClientConsensusStates: {
                    encode(message: _79.ClientConsensusStates, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _79.ClientConsensusStates;
                    fromJSON(object: any): _79.ClientConsensusStates;
                    toJSON(message: _79.ClientConsensusStates): unknown;
                    fromPartial(object: Partial<_79.ClientConsensusStates>): _79.ClientConsensusStates;
                };
                ClientUpdateProposal: {
                    encode(message: _79.ClientUpdateProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _79.ClientUpdateProposal;
                    fromJSON(object: any): _79.ClientUpdateProposal;
                    toJSON(message: _79.ClientUpdateProposal): unknown;
                    fromPartial(object: Partial<_79.ClientUpdateProposal>): _79.ClientUpdateProposal;
                };
                UpgradeProposal: {
                    encode(message: _79.UpgradeProposal, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _79.UpgradeProposal;
                    fromJSON(object: any): _79.UpgradeProposal;
                    toJSON(message: _79.UpgradeProposal): unknown;
                    fromPartial(object: Partial<_79.UpgradeProposal>): _79.UpgradeProposal;
                };
                Height: {
                    encode(message: _79.Height, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _79.Height;
                    fromJSON(object: any): _79.Height;
                    toJSON(message: _79.Height): unknown;
                    fromPartial(object: Partial<_79.Height>): _79.Height;
                };
                Params: {
                    encode(message: _79.Params, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _79.Params;
                    fromJSON(object: any): _79.Params;
                    toJSON(message: _79.Params): unknown;
                    fromPartial(object: Partial<_79.Params>): _79.Params;
                };
            };
        }
        namespace commitment {
            const v1: {
                MerkleRoot: {
                    encode(message: _80.MerkleRoot, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _80.MerkleRoot;
                    fromJSON(object: any): _80.MerkleRoot;
                    toJSON(message: _80.MerkleRoot): unknown;
                    fromPartial(object: Partial<_80.MerkleRoot>): _80.MerkleRoot;
                };
                MerklePrefix: {
                    encode(message: _80.MerklePrefix, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _80.MerklePrefix;
                    fromJSON(object: any): _80.MerklePrefix;
                    toJSON(message: _80.MerklePrefix): unknown;
                    fromPartial(object: Partial<_80.MerklePrefix>): _80.MerklePrefix;
                };
                MerklePath: {
                    encode(message: _80.MerklePath, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _80.MerklePath;
                    fromJSON(object: any): _80.MerklePath;
                    toJSON(message: _80.MerklePath): unknown;
                    fromPartial(object: Partial<_80.MerklePath>): _80.MerklePath;
                };
                MerkleProof: {
                    encode(message: _80.MerkleProof, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _80.MerkleProof;
                    fromJSON(object: any): _80.MerkleProof;
                    toJSON(message: _80.MerkleProof): unknown;
                    fromPartial(object: Partial<_80.MerkleProof>): _80.MerkleProof;
                };
            };
        }
    }
    namespace lightclients {
        namespace tendermint {
            const v1: {
                ClientState: {
                    encode(message: _81.ClientState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _81.ClientState;
                    fromJSON(object: any): _81.ClientState;
                    toJSON(message: _81.ClientState): unknown;
                    fromPartial(object: Partial<_81.ClientState>): _81.ClientState;
                };
                ConsensusState: {
                    encode(message: _81.ConsensusState, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _81.ConsensusState;
                    fromJSON(object: any): _81.ConsensusState;
                    toJSON(message: _81.ConsensusState): unknown;
                    fromPartial(object: Partial<_81.ConsensusState>): _81.ConsensusState;
                };
                Misbehaviour: {
                    encode(message: _81.Misbehaviour, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _81.Misbehaviour;
                    fromJSON(object: any): _81.Misbehaviour;
                    toJSON(message: _81.Misbehaviour): unknown;
                    fromPartial(object: Partial<_81.Misbehaviour>): _81.Misbehaviour;
                };
                Header: {
                    encode(message: _81.Header, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _81.Header;
                    fromJSON(object: any): _81.Header;
                    toJSON(message: _81.Header): unknown;
                    fromPartial(object: Partial<_81.Header>): _81.Header;
                };
                Fraction: {
                    encode(message: _81.Fraction, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
                    decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _81.Fraction;
                    fromJSON(object: any): _81.Fraction;
                    toJSON(message: _81.Fraction): unknown;
                    fromPartial(object: Partial<_81.Fraction>): _81.Fraction;
                };
            };
        }
    }
}
