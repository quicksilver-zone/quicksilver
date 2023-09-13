import { DecCoin, DecCoinSDKType, Coin, CoinSDKType } from "../../base/v1beta1/coin";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/** Params defines the set of params for the distribution module. */
export interface Params {
    communityTax: string;
    baseProposerReward: string;
    bonusProposerReward: string;
    withdrawAddrEnabled: boolean;
}
/** Params defines the set of params for the distribution module. */
export interface ParamsSDKType {
    community_tax: string;
    base_proposer_reward: string;
    bonus_proposer_reward: string;
    withdraw_addr_enabled: boolean;
}
/**
 * ValidatorHistoricalRewards represents historical rewards for a validator.
 * Height is implicit within the store key.
 * Cumulative reward ratio is the sum from the zeroeth period
 * until this period of rewards / tokens, per the spec.
 * The reference count indicates the number of objects
 * which might need to reference this historical entry at any point.
 * ReferenceCount =
 *    number of outstanding delegations which ended the associated period (and
 *    might need to read that record)
 *  + number of slashes which ended the associated period (and might need to
 *  read that record)
 *  + one per validator for the zeroeth period, set on initialization
 */
export interface ValidatorHistoricalRewards {
    cumulativeRewardRatio: DecCoin[];
    referenceCount: number;
}
/**
 * ValidatorHistoricalRewards represents historical rewards for a validator.
 * Height is implicit within the store key.
 * Cumulative reward ratio is the sum from the zeroeth period
 * until this period of rewards / tokens, per the spec.
 * The reference count indicates the number of objects
 * which might need to reference this historical entry at any point.
 * ReferenceCount =
 *    number of outstanding delegations which ended the associated period (and
 *    might need to read that record)
 *  + number of slashes which ended the associated period (and might need to
 *  read that record)
 *  + one per validator for the zeroeth period, set on initialization
 */
export interface ValidatorHistoricalRewardsSDKType {
    cumulative_reward_ratio: DecCoinSDKType[];
    reference_count: number;
}
/**
 * ValidatorCurrentRewards represents current rewards and current
 * period for a validator kept as a running counter and incremented
 * each block as long as the validator's tokens remain constant.
 */
export interface ValidatorCurrentRewards {
    rewards: DecCoin[];
    period: Long;
}
/**
 * ValidatorCurrentRewards represents current rewards and current
 * period for a validator kept as a running counter and incremented
 * each block as long as the validator's tokens remain constant.
 */
export interface ValidatorCurrentRewardsSDKType {
    rewards: DecCoinSDKType[];
    period: Long;
}
/**
 * ValidatorAccumulatedCommission represents accumulated commission
 * for a validator kept as a running counter, can be withdrawn at any time.
 */
export interface ValidatorAccumulatedCommission {
    commission: DecCoin[];
}
/**
 * ValidatorAccumulatedCommission represents accumulated commission
 * for a validator kept as a running counter, can be withdrawn at any time.
 */
export interface ValidatorAccumulatedCommissionSDKType {
    commission: DecCoinSDKType[];
}
/**
 * ValidatorOutstandingRewards represents outstanding (un-withdrawn) rewards
 * for a validator inexpensive to track, allows simple sanity checks.
 */
export interface ValidatorOutstandingRewards {
    rewards: DecCoin[];
}
/**
 * ValidatorOutstandingRewards represents outstanding (un-withdrawn) rewards
 * for a validator inexpensive to track, allows simple sanity checks.
 */
export interface ValidatorOutstandingRewardsSDKType {
    rewards: DecCoinSDKType[];
}
/**
 * ValidatorSlashEvent represents a validator slash event.
 * Height is implicit within the store key.
 * This is needed to calculate appropriate amount of staking tokens
 * for delegations which are withdrawn after a slash has occurred.
 */
export interface ValidatorSlashEvent {
    validatorPeriod: Long;
    fraction: string;
}
/**
 * ValidatorSlashEvent represents a validator slash event.
 * Height is implicit within the store key.
 * This is needed to calculate appropriate amount of staking tokens
 * for delegations which are withdrawn after a slash has occurred.
 */
export interface ValidatorSlashEventSDKType {
    validator_period: Long;
    fraction: string;
}
/** ValidatorSlashEvents is a collection of ValidatorSlashEvent messages. */
export interface ValidatorSlashEvents {
    validatorSlashEvents: ValidatorSlashEvent[];
}
/** ValidatorSlashEvents is a collection of ValidatorSlashEvent messages. */
export interface ValidatorSlashEventsSDKType {
    validator_slash_events: ValidatorSlashEventSDKType[];
}
/** FeePool is the global fee pool for distribution. */
export interface FeePool {
    communityPool: DecCoin[];
}
/** FeePool is the global fee pool for distribution. */
export interface FeePoolSDKType {
    community_pool: DecCoinSDKType[];
}
/**
 * CommunityPoolSpendProposal details a proposal for use of community funds,
 * together with how many coins are proposed to be spent, and to which
 * recipient account.
 */
export interface CommunityPoolSpendProposal {
    title: string;
    description: string;
    recipient: string;
    amount: Coin[];
}
/**
 * CommunityPoolSpendProposal details a proposal for use of community funds,
 * together with how many coins are proposed to be spent, and to which
 * recipient account.
 */
export interface CommunityPoolSpendProposalSDKType {
    title: string;
    description: string;
    recipient: string;
    amount: CoinSDKType[];
}
/**
 * DelegatorStartingInfo represents the starting info for a delegator reward
 * period. It tracks the previous validator period, the delegation's amount of
 * staking token, and the creation height (to check later on if any slashes have
 * occurred). NOTE: Even though validators are slashed to whole staking tokens,
 * the delegators within the validator may be left with less than a full token,
 * thus sdk.Dec is used.
 */
export interface DelegatorStartingInfo {
    previousPeriod: Long;
    stake: string;
    height: Long;
}
/**
 * DelegatorStartingInfo represents the starting info for a delegator reward
 * period. It tracks the previous validator period, the delegation's amount of
 * staking token, and the creation height (to check later on if any slashes have
 * occurred). NOTE: Even though validators are slashed to whole staking tokens,
 * the delegators within the validator may be left with less than a full token,
 * thus sdk.Dec is used.
 */
export interface DelegatorStartingInfoSDKType {
    previous_period: Long;
    stake: string;
    height: Long;
}
/**
 * DelegationDelegatorReward represents the properties
 * of a delegator's delegation reward.
 */
export interface DelegationDelegatorReward {
    validatorAddress: string;
    reward: DecCoin[];
}
/**
 * DelegationDelegatorReward represents the properties
 * of a delegator's delegation reward.
 */
export interface DelegationDelegatorRewardSDKType {
    validator_address: string;
    reward: DecCoinSDKType[];
}
/**
 * CommunityPoolSpendProposalWithDeposit defines a CommunityPoolSpendProposal
 * with a deposit
 */
export interface CommunityPoolSpendProposalWithDeposit {
    title: string;
    description: string;
    recipient: string;
    amount: string;
    deposit: string;
}
/**
 * CommunityPoolSpendProposalWithDeposit defines a CommunityPoolSpendProposal
 * with a deposit
 */
export interface CommunityPoolSpendProposalWithDepositSDKType {
    title: string;
    description: string;
    recipient: string;
    amount: string;
    deposit: string;
}
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial(object: Partial<Params>): Params;
};
export declare const ValidatorHistoricalRewards: {
    encode(message: ValidatorHistoricalRewards, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorHistoricalRewards;
    fromJSON(object: any): ValidatorHistoricalRewards;
    toJSON(message: ValidatorHistoricalRewards): unknown;
    fromPartial(object: Partial<ValidatorHistoricalRewards>): ValidatorHistoricalRewards;
};
export declare const ValidatorCurrentRewards: {
    encode(message: ValidatorCurrentRewards, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorCurrentRewards;
    fromJSON(object: any): ValidatorCurrentRewards;
    toJSON(message: ValidatorCurrentRewards): unknown;
    fromPartial(object: Partial<ValidatorCurrentRewards>): ValidatorCurrentRewards;
};
export declare const ValidatorAccumulatedCommission: {
    encode(message: ValidatorAccumulatedCommission, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorAccumulatedCommission;
    fromJSON(object: any): ValidatorAccumulatedCommission;
    toJSON(message: ValidatorAccumulatedCommission): unknown;
    fromPartial(object: Partial<ValidatorAccumulatedCommission>): ValidatorAccumulatedCommission;
};
export declare const ValidatorOutstandingRewards: {
    encode(message: ValidatorOutstandingRewards, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorOutstandingRewards;
    fromJSON(object: any): ValidatorOutstandingRewards;
    toJSON(message: ValidatorOutstandingRewards): unknown;
    fromPartial(object: Partial<ValidatorOutstandingRewards>): ValidatorOutstandingRewards;
};
export declare const ValidatorSlashEvent: {
    encode(message: ValidatorSlashEvent, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorSlashEvent;
    fromJSON(object: any): ValidatorSlashEvent;
    toJSON(message: ValidatorSlashEvent): unknown;
    fromPartial(object: Partial<ValidatorSlashEvent>): ValidatorSlashEvent;
};
export declare const ValidatorSlashEvents: {
    encode(message: ValidatorSlashEvents, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorSlashEvents;
    fromJSON(object: any): ValidatorSlashEvents;
    toJSON(message: ValidatorSlashEvents): unknown;
    fromPartial(object: Partial<ValidatorSlashEvents>): ValidatorSlashEvents;
};
export declare const FeePool: {
    encode(message: FeePool, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): FeePool;
    fromJSON(object: any): FeePool;
    toJSON(message: FeePool): unknown;
    fromPartial(object: Partial<FeePool>): FeePool;
};
export declare const CommunityPoolSpendProposal: {
    encode(message: CommunityPoolSpendProposal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): CommunityPoolSpendProposal;
    fromJSON(object: any): CommunityPoolSpendProposal;
    toJSON(message: CommunityPoolSpendProposal): unknown;
    fromPartial(object: Partial<CommunityPoolSpendProposal>): CommunityPoolSpendProposal;
};
export declare const DelegatorStartingInfo: {
    encode(message: DelegatorStartingInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelegatorStartingInfo;
    fromJSON(object: any): DelegatorStartingInfo;
    toJSON(message: DelegatorStartingInfo): unknown;
    fromPartial(object: Partial<DelegatorStartingInfo>): DelegatorStartingInfo;
};
export declare const DelegationDelegatorReward: {
    encode(message: DelegationDelegatorReward, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelegationDelegatorReward;
    fromJSON(object: any): DelegationDelegatorReward;
    toJSON(message: DelegationDelegatorReward): unknown;
    fromPartial(object: Partial<DelegationDelegatorReward>): DelegationDelegatorReward;
};
export declare const CommunityPoolSpendProposalWithDeposit: {
    encode(message: CommunityPoolSpendProposalWithDeposit, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): CommunityPoolSpendProposalWithDeposit;
    fromJSON(object: any): CommunityPoolSpendProposalWithDeposit;
    toJSON(message: CommunityPoolSpendProposalWithDeposit): unknown;
    fromPartial(object: Partial<CommunityPoolSpendProposalWithDeposit>): CommunityPoolSpendProposalWithDeposit;
};
