import { Rpc } from "../../../helpers";
import { MsgClaim, MsgClaimResponse } from "./messages";
/** Msg defines the airdrop Msg service. */
export interface Msg {
    claim(request: MsgClaim): Promise<MsgClaimResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    claim(request: MsgClaim): Promise<MsgClaimResponse>;
}
