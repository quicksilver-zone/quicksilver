import { Rpc } from "../../../helpers";
import { MsgRequestRedemption, MsgRequestRedemptionResponse, MsgSignalIntent, MsgSignalIntentResponse } from "./messages";
/** Msg defines the interchainstaking Msg service. */
export interface Msg {
    /**
     * RequestRedemption defines a method for requesting burning of qAssets for
     * native assets.
     */
    requestRedemption(request: MsgRequestRedemption): Promise<MsgRequestRedemptionResponse>;
    /**
     * SignalIntent defines a method for signalling voting intent for one or more
     * validators.
     */
    signalIntent(request: MsgSignalIntent): Promise<MsgSignalIntentResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    requestRedemption(request: MsgRequestRedemption): Promise<MsgRequestRedemptionResponse>;
    signalIntent(request: MsgSignalIntent): Promise<MsgSignalIntentResponse>;
}
