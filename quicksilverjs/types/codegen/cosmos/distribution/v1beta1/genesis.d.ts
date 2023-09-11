import { DecCoin, DecCoinSDKType } from "../../base/v1beta1/coin";
import { ValidatorAccumulatedCommission, ValidatorAccumulatedCommissionSDKType, ValidatorHistoricalRewards, ValidatorHistoricalRewardsSDKType, ValidatorCurrentRewards, ValidatorCurrentRewardsSDKType, DelegatorStartingInfo, DelegatorStartingInfoSDKType, ValidatorSlashEvent, ValidatorSlashEventSDKType, Params, ParamsSDKType, FeePool, FeePoolSDKType } from "./distribution";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/**
 * DelegatorWithdrawInfo is the address for where distributions rewards are
 * withdrawn to by default this struct is only used at genesis to feed in
 * default withdraw addresses.
 */
export interface DelegatorWithdrawInfo {
    /** delegator_address is the address of the delegator. */
    delegatorAddress: string;
    /** withdraw_address is the address to withdraw the delegation rewards to. */
    withdrawAddress: string;
}
/**
 * DelegatorWithdrawInfo is the address for where distributions rewards are
 * withdrawn to by default this struct is only used at genesis to feed in
 * default withdraw addresses.
 */
export interface DelegatorWithdrawInfoSDKType {
    /** delegator_address is the address of the delegator. */
    delegator_address: string;
    /** withdraw_address is the address to withdraw the delegation rewards to. */
    withdraw_address: string;
}
/** ValidatorOutstandingRewardsRecord is used for import/export via genesis json. */
export interface ValidatorOutstandingRewardsRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** outstanding_rewards represents the oustanding rewards of a validator. */
    outstandingRewards: DecCoin[];
}
/** ValidatorOutstandingRewardsRecord is used for import/export via genesis json. */
export interface ValidatorOutstandingRewardsRecordSDKType {
    /** validator_address is the address of the validator. */
    validator_address: string;
    /** outstanding_rewards represents the oustanding rewards of a validator. */
    outstanding_rewards: DecCoinSDKType[];
}
/**
 * ValidatorAccumulatedCommissionRecord is used for import / export via genesis
 * json.
 */
export interface ValidatorAccumulatedCommissionRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** accumulated is the accumulated commission of a validator. */
    accumulated?: ValidatorAccumulatedCommission;
}
/**
 * ValidatorAccumulatedCommissionRecord is used for import / export via genesis
 * json.
 */
export interface ValidatorAccumulatedCommissionRecordSDKType {
    /** validator_address is the address of the validator. */
    validator_address: string;
    /** accumulated is the accumulated commission of a validator. */
    accumulated?: ValidatorAccumulatedCommissionSDKType;
}
/**
 * ValidatorHistoricalRewardsRecord is used for import / export via genesis
 * json.
 */
export interface ValidatorHistoricalRewardsRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** period defines the period the historical rewards apply to. */
    period: Long;
    /** rewards defines the historical rewards of a validator. */
    rewards?: ValidatorHistoricalRewards;
}
/**
 * ValidatorHistoricalRewardsRecord is used for import / export via genesis
 * json.
 */
export interface ValidatorHistoricalRewardsRecordSDKType {
    /** validator_address is the address of the validator. */
    validator_address: string;
    /** period defines the period the historical rewards apply to. */
    period: Long;
    /** rewards defines the historical rewards of a validator. */
    rewards?: ValidatorHistoricalRewardsSDKType;
}
/** ValidatorCurrentRewardsRecord is used for import / export via genesis json. */
export interface ValidatorCurrentRewardsRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** rewards defines the current rewards of a validator. */
    rewards?: ValidatorCurrentRewards;
}
/** ValidatorCurrentRewardsRecord is used for import / export via genesis json. */
export interface ValidatorCurrentRewardsRecordSDKType {
    /** validator_address is the address of the validator. */
    validator_address: string;
    /** rewards defines the current rewards of a validator. */
    rewards?: ValidatorCurrentRewardsSDKType;
}
/** DelegatorStartingInfoRecord used for import / export via genesis json. */
export interface DelegatorStartingInfoRecord {
    /** delegator_address is the address of the delegator. */
    delegatorAddress: string;
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** starting_info defines the starting info of a delegator. */
    startingInfo?: DelegatorStartingInfo;
}
/** DelegatorStartingInfoRecord used for import / export via genesis json. */
export interface DelegatorStartingInfoRecordSDKType {
    /** delegator_address is the address of the delegator. */
    delegator_address: string;
    /** validator_address is the address of the validator. */
    validator_address: string;
    /** starting_info defines the starting info of a delegator. */
    starting_info?: DelegatorStartingInfoSDKType;
}
/** ValidatorSlashEventRecord is used for import / export via genesis json. */
export interface ValidatorSlashEventRecord {
    /** validator_address is the address of the validator. */
    validatorAddress: string;
    /** height defines the block height at which the slash event occured. */
    height: Long;
    /** period is the period of the slash event. */
    period: Long;
    /** validator_slash_event describes the slash event. */
    validatorSlashEvent?: ValidatorSlashEvent;
}
/** ValidatorSlashEventRecord is used for import / export via genesis json. */
export interface ValidatorSlashEventRecordSDKType {
    /** validator_address is the address of the validator. */
    validator_address: string;
    /** height defines the block height at which the slash event occured. */
    height: Long;
    /** period is the period of the slash event. */
    period: Long;
    /** validator_slash_event describes the slash event. */
    validator_slash_event?: ValidatorSlashEventSDKType;
}
/** GenesisState defines the distribution module's genesis state. */
export interface GenesisState {
    /** params defines all the paramaters of the module. */
    params?: Params;
    /** fee_pool defines the fee pool at genesis. */
    feePool?: FeePool;
    /** fee_pool defines the delegator withdraw infos at genesis. */
    delegatorWithdrawInfos: DelegatorWithdrawInfo[];
    /** fee_pool defines the previous proposer at genesis. */
    previousProposer: string;
    /** fee_pool defines the outstanding rewards of all validators at genesis. */
    outstandingRewards: ValidatorOutstandingRewardsRecord[];
    validatorAccumulatedCommissions: ValidatorAccumulatedCommissionRecord[];
    /** fee_pool defines the historical rewards of all validators at genesis. */
    validatorHistoricalRewards: ValidatorHistoricalRewardsRecord[];
    /** fee_pool defines the current rewards of all validators at genesis. */
    validatorCurrentRewards: ValidatorCurrentRewardsRecord[];
    /** fee_pool defines the delegator starting infos at genesis. */
    delegatorStartingInfos: DelegatorStartingInfoRecord[];
    /** fee_pool defines the validator slash events at genesis. */
    validatorSlashEvents: ValidatorSlashEventRecord[];
}
/** GenesisState defines the distribution module's genesis state. */
export interface GenesisStateSDKType {
    /** params defines all the paramaters of the module. */
    params?: ParamsSDKType;
    /** fee_pool defines the fee pool at genesis. */
    fee_pool?: FeePoolSDKType;
    /** fee_pool defines the delegator withdraw infos at genesis. */
    delegator_withdraw_infos: DelegatorWithdrawInfoSDKType[];
    /** fee_pool defines the previous proposer at genesis. */
    previous_proposer: string;
    /** fee_pool defines the outstanding rewards of all validators at genesis. */
    outstanding_rewards: ValidatorOutstandingRewardsRecordSDKType[];
    validator_accumulated_commissions: ValidatorAccumulatedCommissionRecordSDKType[];
    /** fee_pool defines the historical rewards of all validators at genesis. */
    validator_historical_rewards: ValidatorHistoricalRewardsRecordSDKType[];
    /** fee_pool defines the current rewards of all validators at genesis. */
    validator_current_rewards: ValidatorCurrentRewardsRecordSDKType[];
    /** fee_pool defines the delegator starting infos at genesis. */
    delegator_starting_infos: DelegatorStartingInfoRecordSDKType[];
    /** fee_pool defines the validator slash events at genesis. */
    validator_slash_events: ValidatorSlashEventRecordSDKType[];
}
export declare const DelegatorWithdrawInfo: {
    encode(message: DelegatorWithdrawInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelegatorWithdrawInfo;
    fromJSON(object: any): DelegatorWithdrawInfo;
    toJSON(message: DelegatorWithdrawInfo): unknown;
    fromPartial(object: Partial<DelegatorWithdrawInfo>): DelegatorWithdrawInfo;
};
export declare const ValidatorOutstandingRewardsRecord: {
    encode(message: ValidatorOutstandingRewardsRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorOutstandingRewardsRecord;
    fromJSON(object: any): ValidatorOutstandingRewardsRecord;
    toJSON(message: ValidatorOutstandingRewardsRecord): unknown;
    fromPartial(object: Partial<ValidatorOutstandingRewardsRecord>): ValidatorOutstandingRewardsRecord;
};
export declare const ValidatorAccumulatedCommissionRecord: {
    encode(message: ValidatorAccumulatedCommissionRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorAccumulatedCommissionRecord;
    fromJSON(object: any): ValidatorAccumulatedCommissionRecord;
    toJSON(message: ValidatorAccumulatedCommissionRecord): unknown;
    fromPartial(object: Partial<ValidatorAccumulatedCommissionRecord>): ValidatorAccumulatedCommissionRecord;
};
export declare const ValidatorHistoricalRewardsRecord: {
    encode(message: ValidatorHistoricalRewardsRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorHistoricalRewardsRecord;
    fromJSON(object: any): ValidatorHistoricalRewardsRecord;
    toJSON(message: ValidatorHistoricalRewardsRecord): unknown;
    fromPartial(object: Partial<ValidatorHistoricalRewardsRecord>): ValidatorHistoricalRewardsRecord;
};
export declare const ValidatorCurrentRewardsRecord: {
    encode(message: ValidatorCurrentRewardsRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorCurrentRewardsRecord;
    fromJSON(object: any): ValidatorCurrentRewardsRecord;
    toJSON(message: ValidatorCurrentRewardsRecord): unknown;
    fromPartial(object: Partial<ValidatorCurrentRewardsRecord>): ValidatorCurrentRewardsRecord;
};
export declare const DelegatorStartingInfoRecord: {
    encode(message: DelegatorStartingInfoRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelegatorStartingInfoRecord;
    fromJSON(object: any): DelegatorStartingInfoRecord;
    toJSON(message: DelegatorStartingInfoRecord): unknown;
    fromPartial(object: Partial<DelegatorStartingInfoRecord>): DelegatorStartingInfoRecord;
};
export declare const ValidatorSlashEventRecord: {
    encode(message: ValidatorSlashEventRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorSlashEventRecord;
    fromJSON(object: any): ValidatorSlashEventRecord;
    toJSON(message: ValidatorSlashEventRecord): unknown;
    fromPartial(object: Partial<ValidatorSlashEventRecord>): ValidatorSlashEventRecord;
};
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: Partial<GenesisState>): GenesisState;
};
