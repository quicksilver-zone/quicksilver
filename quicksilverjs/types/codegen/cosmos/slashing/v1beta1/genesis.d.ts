import { Params, ParamsSDKType, ValidatorSigningInfo, ValidatorSigningInfoSDKType } from "./slashing";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/** GenesisState defines the slashing module's genesis state. */
export interface GenesisState {
    /** params defines all the paramaters of related to deposit. */
    params?: Params;
    /**
     * signing_infos represents a map between validator addresses and their
     * signing infos.
     */
    signingInfos: SigningInfo[];
    /**
     * missed_blocks represents a map between validator addresses and their
     * missed blocks.
     */
    missedBlocks: ValidatorMissedBlocks[];
}
/** GenesisState defines the slashing module's genesis state. */
export interface GenesisStateSDKType {
    /** params defines all the paramaters of related to deposit. */
    params?: ParamsSDKType;
    /**
     * signing_infos represents a map between validator addresses and their
     * signing infos.
     */
    signing_infos: SigningInfoSDKType[];
    /**
     * missed_blocks represents a map between validator addresses and their
     * missed blocks.
     */
    missed_blocks: ValidatorMissedBlocksSDKType[];
}
/** SigningInfo stores validator signing info of corresponding address. */
export interface SigningInfo {
    /** address is the validator address. */
    address: string;
    /** validator_signing_info represents the signing info of this validator. */
    validatorSigningInfo?: ValidatorSigningInfo;
}
/** SigningInfo stores validator signing info of corresponding address. */
export interface SigningInfoSDKType {
    /** address is the validator address. */
    address: string;
    /** validator_signing_info represents the signing info of this validator. */
    validator_signing_info?: ValidatorSigningInfoSDKType;
}
/**
 * ValidatorMissedBlocks contains array of missed blocks of corresponding
 * address.
 */
export interface ValidatorMissedBlocks {
    /** address is the validator address. */
    address: string;
    /** missed_blocks is an array of missed blocks by the validator. */
    missedBlocks: MissedBlock[];
}
/**
 * ValidatorMissedBlocks contains array of missed blocks of corresponding
 * address.
 */
export interface ValidatorMissedBlocksSDKType {
    /** address is the validator address. */
    address: string;
    /** missed_blocks is an array of missed blocks by the validator. */
    missed_blocks: MissedBlockSDKType[];
}
/** MissedBlock contains height and missed status as boolean. */
export interface MissedBlock {
    /** index is the height at which the block was missed. */
    index: Long;
    /** missed is the missed status. */
    missed: boolean;
}
/** MissedBlock contains height and missed status as boolean. */
export interface MissedBlockSDKType {
    /** index is the height at which the block was missed. */
    index: Long;
    /** missed is the missed status. */
    missed: boolean;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
export declare const SigningInfo: {
    encode(message: SigningInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SigningInfo;
    fromJSON(object: any): SigningInfo;
    toJSON(message: SigningInfo): unknown;
    fromPartial(object: Partial<SigningInfo>): SigningInfo;
};
export declare const ValidatorMissedBlocks: {
    encode(message: ValidatorMissedBlocks, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorMissedBlocks;
    fromJSON(object: any): ValidatorMissedBlocks;
    toJSON(message: ValidatorMissedBlocks): unknown;
    fromPartial(object: Partial<ValidatorMissedBlocks>): ValidatorMissedBlocks;
};
export declare const MissedBlock: {
    encode(message: MissedBlock, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MissedBlock;
    fromJSON(object: any): MissedBlock;
    toJSON(message: MissedBlock): unknown;
    fromPartial(object: Partial<MissedBlock>): MissedBlock;
};
