import * as _m0 from "protobufjs/minimal";
export interface RegisterZoneProposal {
    title: string;
    description: string;
    connectionId: string;
    baseDenom: string;
    localDenom: string;
    accountPrefix: string;
    multiSend: boolean;
    liquidityModule: boolean;
}
export interface RegisterZoneProposalSDKType {
    title: string;
    description: string;
    connection_id: string;
    base_denom: string;
    local_denom: string;
    account_prefix: string;
    multi_send: boolean;
    liquidity_module: boolean;
}
export interface RegisterZoneProposalWithDeposit {
    title: string;
    description: string;
    connectionId: string;
    baseDenom: string;
    localDenom: string;
    accountPrefix: string;
    multiSend: boolean;
    liquidityModule: boolean;
    deposit: string;
}
export interface RegisterZoneProposalWithDepositSDKType {
    title: string;
    description: string;
    connection_id: string;
    base_denom: string;
    local_denom: string;
    account_prefix: string;
    multi_send: boolean;
    liquidity_module: boolean;
    deposit: string;
}
export interface UpdateZoneProposal {
    title: string;
    description: string;
    chainId: string;
    changes: UpdateZoneValue[];
}
export interface UpdateZoneProposalSDKType {
    title: string;
    description: string;
    chain_id: string;
    changes: UpdateZoneValueSDKType[];
}
export interface UpdateZoneProposalWithDeposit {
    title: string;
    description: string;
    chainId: string;
    changes: UpdateZoneValue[];
    deposit: string;
}
export interface UpdateZoneProposalWithDepositSDKType {
    title: string;
    description: string;
    chain_id: string;
    changes: UpdateZoneValueSDKType[];
    deposit: string;
}
/**
 * ParamChange defines an individual parameter change, for use in
 * ParameterChangeProposal.
 */
export interface UpdateZoneValue {
    key: string;
    value: string;
}
/**
 * ParamChange defines an individual parameter change, for use in
 * ParameterChangeProposal.
 */
export interface UpdateZoneValueSDKType {
    key: string;
    value: string;
}
export declare const RegisterZoneProposal: {
    encode(message: RegisterZoneProposal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): RegisterZoneProposal;
    fromJSON(object: any): RegisterZoneProposal;
    toJSON(message: RegisterZoneProposal): unknown;
    fromPartial(object: Partial<RegisterZoneProposal>): RegisterZoneProposal;
};
export declare const RegisterZoneProposalWithDeposit: {
    encode(message: RegisterZoneProposalWithDeposit, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): RegisterZoneProposalWithDeposit;
    fromJSON(object: any): RegisterZoneProposalWithDeposit;
    toJSON(message: RegisterZoneProposalWithDeposit): unknown;
    fromPartial(object: Partial<RegisterZoneProposalWithDeposit>): RegisterZoneProposalWithDeposit;
};
export declare const UpdateZoneProposal: {
    encode(message: UpdateZoneProposal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateZoneProposal;
    fromJSON(object: any): UpdateZoneProposal;
    toJSON(message: UpdateZoneProposal): unknown;
    fromPartial(object: Partial<UpdateZoneProposal>): UpdateZoneProposal;
};
export declare const UpdateZoneProposalWithDeposit: {
    encode(message: UpdateZoneProposalWithDeposit, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateZoneProposalWithDeposit;
    fromJSON(object: any): UpdateZoneProposalWithDeposit;
    toJSON(message: UpdateZoneProposalWithDeposit): unknown;
    fromPartial(object: Partial<UpdateZoneProposalWithDeposit>): UpdateZoneProposalWithDeposit;
};
export declare const UpdateZoneValue: {
    encode(message: UpdateZoneValue, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateZoneValue;
    fromJSON(object: any): UpdateZoneValue;
    toJSON(message: UpdateZoneValue): unknown;
    fromPartial(object: Partial<UpdateZoneValue>): UpdateZoneValue;
};
