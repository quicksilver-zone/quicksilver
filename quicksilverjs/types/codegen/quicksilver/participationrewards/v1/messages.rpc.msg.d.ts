import { Rpc } from "../../../helpers";
import { MsgSubmitClaim, MsgSubmitClaimResponse } from "./messages";
/** Msg defines the participationrewards Msg service. */
export interface Msg {
    submitClaim(request: MsgSubmitClaim): Promise<MsgSubmitClaimResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    submitClaim(request: MsgSubmitClaim): Promise<MsgSubmitClaimResponse>;
}
