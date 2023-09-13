import { Delegation, DelegationSDKType, DelegatorIntent, DelegatorIntentSDKType, Zone, ZoneSDKType, Receipt, ReceiptSDKType, PortConnectionTuple, PortConnectionTupleSDKType, WithdrawalRecord, WithdrawalRecordSDKType } from "./interchainstaking";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
export interface Params {
    depositInterval: Long;
    validatorsetInterval: Long;
    commissionRate: string;
}
export interface ParamsSDKType {
    deposit_interval: Long;
    validatorset_interval: Long;
    commission_rate: string;
}
export interface DelegationsForZone {
    chainId: string;
    delegations: Delegation[];
}
export interface DelegationsForZoneSDKType {
    chain_id: string;
    delegations: DelegationSDKType[];
}
export interface DelegatorIntentsForZone {
    chainId: string;
    delegationIntent: DelegatorIntent[];
    snapshot: boolean;
}
export interface DelegatorIntentsForZoneSDKType {
    chain_id: string;
    delegation_intent: DelegatorIntentSDKType[];
    snapshot: boolean;
}
/** GenesisState defines the interchainstaking module's genesis state. */
export interface GenesisState {
    params?: Params;
    zones: Zone[];
    receipts: Receipt[];
    delegations: DelegationsForZone[];
    delegatorIntents: DelegatorIntentsForZone[];
    portConnections: PortConnectionTuple[];
    withdrawalRecords: WithdrawalRecord[];
}
/** GenesisState defines the interchainstaking module's genesis state. */
export interface GenesisStateSDKType {
    params?: ParamsSDKType;
    zones: ZoneSDKType[];
    receipts: ReceiptSDKType[];
    delegations: DelegationsForZoneSDKType[];
    delegator_intents: DelegatorIntentsForZoneSDKType[];
    port_connections: PortConnectionTupleSDKType[];
    withdrawal_records: WithdrawalRecordSDKType[];
}
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial(object: Partial<Params>): Params;
};
export declare const DelegationsForZone: {
    encode(message: DelegationsForZone, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelegationsForZone;
    fromJSON(object: any): DelegationsForZone;
    toJSON(message: DelegationsForZone): unknown;
    fromPartial(object: Partial<DelegationsForZone>): DelegationsForZone;
};
export declare const DelegatorIntentsForZone: {
    encode(message: DelegatorIntentsForZone, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelegatorIntentsForZone;
    fromJSON(object: any): DelegatorIntentsForZone;
    toJSON(message: DelegatorIntentsForZone): unknown;
    fromPartial(object: Partial<DelegatorIntentsForZone>): DelegatorIntentsForZone;
};
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
