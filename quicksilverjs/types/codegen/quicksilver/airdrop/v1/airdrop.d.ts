import { Timestamp, TimestampSDKType } from "../../../google/protobuf/timestamp";
import { Duration, DurationSDKType } from "../../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { Long } from "../../../helpers";
/** Action is used as an enum to denote specific actions or tasks. */
export declare enum Action {
    /** ActionUndefined - Undefined action (per protobuf spec) */
    ActionUndefined = 0,
    /** ActionInitialClaim - Initial claim action */
    ActionInitialClaim = 1,
    /** ActionDepositT1 - Deposit tier 1 (e.g. > 5% of base_value) */
    ActionDepositT1 = 2,
    /** ActionDepositT2 - Deposit tier 2 (e.g. > 10% of base_value) */
    ActionDepositT2 = 3,
    /** ActionDepositT3 - Deposit tier 3 (e.g. > 15% of base_value) */
    ActionDepositT3 = 4,
    /** ActionDepositT4 - Deposit tier 4 (e.g. > 22% of base_value) */
    ActionDepositT4 = 5,
    /** ActionDepositT5 - Deposit tier 5 (e.g. > 30% of base_value) */
    ActionDepositT5 = 6,
    /** ActionStakeQCK - Active QCK delegation */
    ActionStakeQCK = 7,
    /** ActionSignalIntent - Intent is set */
    ActionSignalIntent = 8,
    /** ActionQSGov - Cast governance vote on QS */
    ActionQSGov = 9,
    /** ActionGbP - Governance By Proxy (GbP): cast vote on remote zone */
    ActionGbP = 10,
    /** ActionOsmosis - Provide liquidity on Osmosis */
    ActionOsmosis = 11,
    UNRECOGNIZED = -1
}
/** Action is used as an enum to denote specific actions or tasks. */
export declare enum ActionSDKType {
    /** ActionUndefined - Undefined action (per protobuf spec) */
    ActionUndefined = 0,
    /** ActionInitialClaim - Initial claim action */
    ActionInitialClaim = 1,
    /** ActionDepositT1 - Deposit tier 1 (e.g. > 5% of base_value) */
    ActionDepositT1 = 2,
    /** ActionDepositT2 - Deposit tier 2 (e.g. > 10% of base_value) */
    ActionDepositT2 = 3,
    /** ActionDepositT3 - Deposit tier 3 (e.g. > 15% of base_value) */
    ActionDepositT3 = 4,
    /** ActionDepositT4 - Deposit tier 4 (e.g. > 22% of base_value) */
    ActionDepositT4 = 5,
    /** ActionDepositT5 - Deposit tier 5 (e.g. > 30% of base_value) */
    ActionDepositT5 = 6,
    /** ActionStakeQCK - Active QCK delegation */
    ActionStakeQCK = 7,
    /** ActionSignalIntent - Intent is set */
    ActionSignalIntent = 8,
    /** ActionQSGov - Cast governance vote on QS */
    ActionQSGov = 9,
    /** ActionGbP - Governance By Proxy (GbP): cast vote on remote zone */
    ActionGbP = 10,
    /** ActionOsmosis - Provide liquidity on Osmosis */
    ActionOsmosis = 11,
    UNRECOGNIZED = -1
}
export declare function actionFromJSON(object: any): Action;
export declare function actionToJSON(object: Action): string;
/** Status is used as an enum to denote zone status. */
export declare enum Status {
    StatusUndefined = 0,
    StatusActive = 1,
    StatusFuture = 2,
    StatusExpired = 3,
    UNRECOGNIZED = -1
}
/** Status is used as an enum to denote zone status. */
export declare enum StatusSDKType {
    StatusUndefined = 0,
    StatusActive = 1,
    StatusFuture = 2,
    StatusExpired = 3,
    UNRECOGNIZED = -1
}
export declare function statusFromJSON(object: any): Status;
export declare function statusToJSON(object: Status): string;
/** ZoneDrop represents an airdrop for a specific zone. */
export interface ZoneDrop {
    chainId: string;
    startTime?: Timestamp;
    duration?: Duration;
    decay?: Duration;
    allocation: Long;
    actions: string[];
    isConcluded: boolean;
}
/** ZoneDrop represents an airdrop for a specific zone. */
export interface ZoneDropSDKType {
    chain_id: string;
    start_time?: TimestampSDKType;
    duration?: DurationSDKType;
    decay?: DurationSDKType;
    allocation: Long;
    actions: string[];
    is_concluded: boolean;
}
export interface ClaimRecord_ActionsCompletedEntry {
    key: number;
    value?: CompletedAction;
}
export interface ClaimRecord_ActionsCompletedEntrySDKType {
    key: number;
    value?: CompletedActionSDKType;
}
/**
 * ClaimRecord represents a users' claim (including completed claims) for a
 * given zone.
 */
export interface ClaimRecord {
    chainId: string;
    address: string;
    /** Protobuf3 does not allow enum as map key */
    actionsCompleted?: {
        [key: number]: CompletedAction;
    };
    maxAllocation: Long;
    baseValue: Long;
}
/**
 * ClaimRecord represents a users' claim (including completed claims) for a
 * given zone.
 */
export interface ClaimRecordSDKType {
    chain_id: string;
    address: string;
    /** Protobuf3 does not allow enum as map key */
    actions_completed?: {
        [key: number]: CompletedActionSDKType;
    };
    max_allocation: Long;
    base_value: Long;
}
/** CompletedAction represents a claim action completed by the user. */
export interface CompletedAction {
    completeTime?: Timestamp;
    claimAmount: Long;
}
/** CompletedAction represents a claim action completed by the user. */
export interface CompletedActionSDKType {
    complete_time?: TimestampSDKType;
    claim_amount: Long;
}
export declare const ZoneDrop: {
    encode(message: ZoneDrop, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ZoneDrop;
    fromJSON(object: any): ZoneDrop;
    toJSON(message: ZoneDrop): unknown;
    fromPartial(object: Partial<ZoneDrop>): ZoneDrop;
};
export declare const ClaimRecord_ActionsCompletedEntry: {
    encode(message: ClaimRecord_ActionsCompletedEntry, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ClaimRecord_ActionsCompletedEntry;
    fromJSON(object: any): ClaimRecord_ActionsCompletedEntry;
    toJSON(message: ClaimRecord_ActionsCompletedEntry): unknown;
    fromPartial(object: Partial<ClaimRecord_ActionsCompletedEntry>): ClaimRecord_ActionsCompletedEntry;
};
export declare const ClaimRecord: {
    encode(message: ClaimRecord, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ClaimRecord;
    fromJSON(object: any): ClaimRecord;
    toJSON(message: ClaimRecord): unknown;
    fromPartial(object: Partial<ClaimRecord>): ClaimRecord;
};
export declare const CompletedAction: {
    encode(message: CompletedAction, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): CompletedAction;
    fromJSON(object: any): CompletedAction;
    toJSON(message: CompletedAction): unknown;
    fromPartial(object: Partial<CompletedAction>): CompletedAction;
};
