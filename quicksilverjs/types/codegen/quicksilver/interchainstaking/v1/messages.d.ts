import { Coin, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import { ValidatorIntent, ValidatorIntentSDKType } from "./interchainstaking";
import * as _m0 from "protobufjs/minimal";
/**
 * MsgRequestRedemption represents a message type to request a burn of qAssets
 * for native assets.
 */
export interface MsgRequestRedemption {
    value?: Coin;
    destinationAddress: string;
    fromAddress: string;
}
/**
 * MsgRequestRedemption represents a message type to request a burn of qAssets
 * for native assets.
 */
export interface MsgRequestRedemptionSDKType {
    value?: CoinSDKType;
    destination_address: string;
    from_address: string;
}
/**
 * MsgSignalIntent represents a message type for signalling voting intent for
 * one or more validators.
 */
export interface MsgSignalIntent {
    chainId: string;
    intents: ValidatorIntent[];
    fromAddress: string;
}
/**
 * MsgSignalIntent represents a message type for signalling voting intent for
 * one or more validators.
 */
export interface MsgSignalIntentSDKType {
    chain_id: string;
    intents: ValidatorIntentSDKType[];
    from_address: string;
}
/** MsgRequestRedemptionResponse defines the MsgRequestRedemption response type. */
export interface MsgRequestRedemptionResponse {
}
/** MsgRequestRedemptionResponse defines the MsgRequestRedemption response type. */
export interface MsgRequestRedemptionResponseSDKType {
}
/** MsgSignalIntentResponse defines the MsgSignalIntent response type. */
export interface MsgSignalIntentResponse {
}
/** MsgSignalIntentResponse defines the MsgSignalIntent response type. */
export interface MsgSignalIntentResponseSDKType {
}
export declare const MsgRequestRedemption: {
    encode(message: MsgRequestRedemption, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestRedemption;
    fromJSON(object: any): MsgRequestRedemption;
    toJSON(message: MsgRequestRedemption): unknown;
    fromPartial(object: Partial<MsgRequestRedemption>): MsgRequestRedemption;
};
export declare const MsgSignalIntent: {
    encode(message: MsgSignalIntent, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSignalIntent;
    fromJSON(object: any): MsgSignalIntent;
    toJSON(message: MsgSignalIntent): unknown;
    fromPartial(object: Partial<MsgSignalIntent>): MsgSignalIntent;
};
export declare const MsgRequestRedemptionResponse: {
    encode(_: MsgRequestRedemptionResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestRedemptionResponse;
    fromJSON(_: any): MsgRequestRedemptionResponse;
    toJSON(_: MsgRequestRedemptionResponse): unknown;
    fromPartial(_: Partial<MsgRequestRedemptionResponse>): MsgRequestRedemptionResponse;
};
export declare const MsgSignalIntentResponse: {
    encode(_: MsgSignalIntentResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSignalIntentResponse;
    fromJSON(_: any): MsgSignalIntentResponse;
    toJSON(_: MsgSignalIntentResponse): unknown;
    fromPartial(_: Partial<MsgSignalIntentResponse>): MsgSignalIntentResponse;
};
