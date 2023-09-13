import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/**
 * Capability defines an implementation of an object capability. The index
 * provided to a Capability must be globally unique.
 */
export interface Capability {
    index: Long;
}
/**
 * Capability defines an implementation of an object capability. The index
 * provided to a Capability must be globally unique.
 */
export interface CapabilitySDKType {
    index: Long;
}
/**
 * Owner defines a single capability owner. An owner is defined by the name of
 * capability and the module name.
 */
export interface Owner {
    module: string;
    name: string;
}
/**
 * Owner defines a single capability owner. An owner is defined by the name of
 * capability and the module name.
 */
export interface OwnerSDKType {
    module: string;
    name: string;
}
/**
 * CapabilityOwners defines a set of owners of a single Capability. The set of
 * owners must be unique.
 */
export interface CapabilityOwners {
    owners: Owner[];
}
/**
 * CapabilityOwners defines a set of owners of a single Capability. The set of
 * owners must be unique.
 */
export interface CapabilityOwnersSDKType {
    owners: OwnerSDKType[];
}
export declare const Capability: {
    encode(message: Capability, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Capability;
    fromJSON(object: any): Capability;
    toJSON(message: Capability): unknown;
    fromPartial(object: Partial<Capability>): Capability;
};
export declare const Owner: {
    encode(message: Owner, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Owner;
    fromJSON(object: any): Owner;
    toJSON(message: Owner): unknown;
    fromPartial(object: Partial<Owner>): Owner;
};
export declare const CapabilityOwners: {
    encode(message: CapabilityOwners, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): CapabilityOwners;
    fromJSON(object: any): CapabilityOwners;
    toJSON(message: CapabilityOwners): unknown;
    fromPartial(object: Partial<CapabilityOwners>): CapabilityOwners;
};
