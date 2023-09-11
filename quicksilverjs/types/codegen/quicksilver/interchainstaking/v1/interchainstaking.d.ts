import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import { Timestamp, TimestampSDKType } from "../../../google/protobuf/timestamp";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
export interface Zone {
    connectionId: string;
    chainId: string;
    depositAddress?: ICAAccount;
    withdrawalAddress?: ICAAccount;
    performanceAddress?: ICAAccount;
    delegationAddress?: ICAAccount;
    accountPrefix: string;
    localDenom: string;
    baseDenom: string;
    redemptionRate: string;
    lastRedemptionRate: string;
    validators: Validator[];
    aggregateIntent: ValidatorIntent[];
    multiSend: boolean;
    liquidityModule: boolean;
    withdrawalWaitgroup: number;
    ibcNextValidatorsHash: Uint8Array;
    validatorSelectionAllocation: Long;
    holdingsAllocation: Long;
    /** deprecated remove me. */
    lastEpochHeight: Long;
    tvl: string;
    unbondingPeriod: Long;
}
export interface ZoneSDKType {
    connection_id: string;
    chain_id: string;
    deposit_address?: ICAAccountSDKType;
    withdrawal_address?: ICAAccountSDKType;
    performance_address?: ICAAccountSDKType;
    delegation_address?: ICAAccountSDKType;
    account_prefix: string;
    local_denom: string;
    base_denom: string;
    redemption_rate: string;
    last_redemption_rate: string;
    validators: ValidatorSDKType[];
    aggregate_intent: ValidatorIntentSDKType[];
    multi_send: boolean;
    liquidity_module: boolean;
    withdrawal_waitgroup: number;
    ibc_next_validators_hash: Uint8Array;
    validator_selection_allocation: Long;
    holdings_allocation: Long;
    /** deprecated remove me. */
    last_epoch_height: Long;
    tvl: string;
    unbonding_period: Long;
}
export interface ICAAccount {
    address: string;
    /** balance defines the different coins this balance holds. */
    balance: Coin[];
    portName: string;
    withdrawalAddress: string;
    balanceWaitgroup: number;
}
export interface ICAAccountSDKType {
    address: string;
    /** balance defines the different coins this balance holds. */
    balance: CoinSDKType[];
    port_name: string;
    withdrawal_address: string;
    balance_waitgroup: number;
}
export interface Distribution {
    valoper: string;
    amount: Long;
}
export interface DistributionSDKType {
    valoper: string;
    amount: Long;
}
export interface WithdrawalRecord {
    chainId: string;
    delegator: string;
    distribution: Distribution[];
    recipient: string;
    amount: Coin[];
    burnAmount?: Coin;
    txhash: string;
    status: number;
    completionTime?: Timestamp;
}
export interface WithdrawalRecordSDKType {
    chain_id: string;
    delegator: string;
    distribution: DistributionSDKType[];
    recipient: string;
    amount: CoinSDKType[];
    burn_amount?: CoinSDKType;
    txhash: string;
    status: number;
    completion_time?: TimestampSDKType;
}
export interface UnbondingRecord {
    chainId: string;
    epochNumber: Long;
    validator: string;
    relatedTxhash: string[];
}
export interface UnbondingRecordSDKType {
    chain_id: string;
    epoch_number: Long;
    validator: string;
    related_txhash: string[];
}
export interface RedelegationRecord {
    chainId: string;
    epochNumber: Long;
    delegator: string;
    validator: string;
    amount: Long;
    completionTime?: Timestamp;
}
export interface RedelegationRecordSDKType {
    chain_id: string;
    epoch_number: Long;
    delegator: string;
    validator: string;
    amount: Long;
    completion_time?: TimestampSDKType;
}
export interface TransferRecord {
    sender: string;
    recipient: string;
    amount?: Coin;
}
export interface TransferRecordSDKType {
    sender: string;
    recipient: string;
    amount?: CoinSDKType;
}
export interface Validator {
    valoperAddress: string;
    commissionRate: string;
    delegatorShares: string;
    votingPower: string;
    score: string;
}
export interface ValidatorSDKType {
    valoper_address: string;
    commission_rate: string;
    delegator_shares: string;
    voting_power: string;
    score: string;
}
export interface DelegatorIntent {
    delegator: string;
    intents: ValidatorIntent[];
}
export interface DelegatorIntentSDKType {
    delegator: string;
    intents: ValidatorIntentSDKType[];
}
export interface ValidatorIntent {
    valoperAddress: string;
    weight: string;
}
export interface ValidatorIntentSDKType {
    valoper_address: string;
    weight: string;
}
export interface Delegation {
    delegationAddress: string;
    validatorAddress: string;
    amount?: Coin;
    height: Long;
    redelegationEnd: Long;
}
export interface DelegationSDKType {
    delegation_address: string;
    validator_address: string;
    amount?: CoinSDKType;
    height: Long;
    redelegation_end: Long;
}
export interface PortConnectionTuple {
    connectionId: string;
    portId: string;
}
export interface PortConnectionTupleSDKType {
    connection_id: string;
    port_id: string;
}
export interface Receipt {
    chainId: string;
    sender: string;
    txhash: string;
    amount: Coin[];
}
export interface ReceiptSDKType {
    chain_id: string;
    sender: string;
    txhash: string;
    amount: CoinSDKType[];
}
export declare const Zone: {
    encode(message: Zone, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Zone;
    fromJSON(object: any): Zone;
    toJSON(message: Zone): unknown;
    fromPartial(object: Partial<Zone>): Zone;
};
export declare const ICAAccount: {
    encode(message: ICAAccount, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ICAAccount;
    fromJSON(object: any): ICAAccount;
    toJSON(message: ICAAccount): unknown;
    fromPartial(object: Partial<ICAAccount>): ICAAccount;
};
export declare const Distribution: {
    encode(message: Distribution, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Distribution;
    fromJSON(object: any): Distribution;
    toJSON(message: Distribution): unknown;
    fromPartial(object: Partial<Distribution>): Distribution;
};
export declare const WithdrawalRecord: {
    encode(message: WithdrawalRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): WithdrawalRecord;
    fromJSON(object: any): WithdrawalRecord;
    toJSON(message: WithdrawalRecord): unknown;
    fromPartial(object: Partial<WithdrawalRecord>): WithdrawalRecord;
};
export declare const UnbondingRecord: {
    encode(message: UnbondingRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UnbondingRecord;
    fromJSON(object: any): UnbondingRecord;
    toJSON(message: UnbondingRecord): unknown;
    fromPartial(object: Partial<UnbondingRecord>): UnbondingRecord;
};
export declare const RedelegationRecord: {
    encode(message: RedelegationRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): RedelegationRecord;
    fromJSON(object: any): RedelegationRecord;
    toJSON(message: RedelegationRecord): unknown;
    fromPartial(object: Partial<RedelegationRecord>): RedelegationRecord;
};
export declare const TransferRecord: {
    encode(message: TransferRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TransferRecord;
    fromJSON(object: any): TransferRecord;
    toJSON(message: TransferRecord): unknown;
    fromPartial(object: Partial<TransferRecord>): TransferRecord;
};
export declare const Validator: {
    encode(message: Validator, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Validator;
    fromJSON(object: any): Validator;
    toJSON(message: Validator): unknown;
    fromPartial(object: Partial<Validator>): Validator;
};
export declare const DelegatorIntent: {
    encode(message: DelegatorIntent, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelegatorIntent;
    fromJSON(object: any): DelegatorIntent;
    toJSON(message: DelegatorIntent): unknown;
    fromPartial(object: Partial<DelegatorIntent>): DelegatorIntent;
};
export declare const ValidatorIntent: {
    encode(message: ValidatorIntent, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorIntent;
    fromJSON(object: any): ValidatorIntent;
    toJSON(message: ValidatorIntent): unknown;
    fromPartial(object: Partial<ValidatorIntent>): ValidatorIntent;
};
export declare const Delegation: {
    encode(message: Delegation, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Delegation;
    fromJSON(object: any): Delegation;
    toJSON(message: Delegation): unknown;
    fromPartial(object: Partial<Delegation>): Delegation;
};
export declare const PortConnectionTuple: {
    encode(message: PortConnectionTuple, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): PortConnectionTuple;
    fromJSON(object: any): PortConnectionTuple;
    toJSON(message: PortConnectionTuple): unknown;
    fromPartial(object: Partial<PortConnectionTuple>): PortConnectionTuple;
};
export declare const Receipt: {
    encode(message: Receipt, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Receipt;
    fromJSON(object: any): Receipt;
    toJSON(message: Receipt): unknown;
    fromPartial(object: Partial<Receipt>): Receipt;
};
