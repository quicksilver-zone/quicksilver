import * as _m0 from "protobufjs/minimal";
export interface AddProtocolDataProposal {
    title: string;
    description: string;
    type: string;
    data: string;
    key: string;
}
export interface AddProtocolDataProposalSDKType {
    title: string;
    description: string;
    type: string;
    data: string;
    key: string;
}
export interface AddProtocolDataProposalWithDeposit {
    title: string;
    description: string;
    protocol: string;
    type: string;
    key: string;
    data: Uint8Array;
    deposit: string;
}
export interface AddProtocolDataProposalWithDepositSDKType {
    title: string;
    description: string;
    protocol: string;
    type: string;
    key: string;
    data: Uint8Array;
    deposit: string;
}
export declare const AddProtocolDataProposal: {
    encode(message: AddProtocolDataProposal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): AddProtocolDataProposal;
    fromJSON(object: any): AddProtocolDataProposal;
    toJSON(message: AddProtocolDataProposal): unknown;
    fromPartial(object: Partial<AddProtocolDataProposal>): AddProtocolDataProposal;
};
export declare const AddProtocolDataProposalWithDeposit: {
    encode(message: AddProtocolDataProposalWithDeposit, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): AddProtocolDataProposalWithDeposit;
    fromJSON(object: any): AddProtocolDataProposalWithDeposit;
    toJSON(message: AddProtocolDataProposalWithDeposit): unknown;
    fromPartial(object: Partial<AddProtocolDataProposalWithDeposit>): AddProtocolDataProposalWithDeposit;
};
