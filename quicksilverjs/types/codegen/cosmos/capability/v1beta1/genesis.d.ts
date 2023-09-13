import { CapabilityOwners, CapabilityOwnersSDKType } from "./capability";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/** GenesisOwners defines the capability owners with their corresponding index. */
export interface GenesisOwners {
    /** index is the index of the capability owner. */
    index: Long;
    /** index_owners are the owners at the given index. */
    indexOwners?: CapabilityOwners;
}
/** GenesisOwners defines the capability owners with their corresponding index. */
export interface GenesisOwnersSDKType {
    /** index is the index of the capability owner. */
    index: Long;
    /** index_owners are the owners at the given index. */
    index_owners?: CapabilityOwnersSDKType;
}
/** GenesisState defines the capability module's genesis state. */
export interface GenesisState {
    /** index is the capability global index. */
    index: Long;
    /**
     * owners represents a map from index to owners of the capability index
     * index key is string to allow amino marshalling.
     */
    owners: GenesisOwners[];
}
/** GenesisState defines the capability module's genesis state. */
export interface GenesisStateSDKType {
    /** index is the capability global index. */
    index: Long;
    /**
     * owners represents a map from index to owners of the capability index
     * index key is string to allow amino marshalling.
     */
    owners: GenesisOwnersSDKType[];
}
export declare const GenesisOwners: {
    encode(message: GenesisOwners, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisOwners;
    fromJSON(object: any): GenesisOwners;
    toJSON(message: GenesisOwners): unknown;
    fromPartial(object: Partial<GenesisOwners>): GenesisOwners;
};
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
