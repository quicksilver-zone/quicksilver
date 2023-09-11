import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgSubmitClaim } from "./messages";
export declare const registry: ReadonlyArray<[string, GeneratedType]>;
export declare const load: (protoRegistry: Registry) => void;
export declare const MessageComposer: {
    encoded: {
        submitClaim(value: MsgSubmitClaim): {
            typeUrl: string;
            value: Uint8Array;
        };
    };
    withTypeUrl: {
        submitClaim(value: MsgSubmitClaim): {
            typeUrl: string;
            value: MsgSubmitClaim;
        };
    };
    toJSON: {
        submitClaim(value: MsgSubmitClaim): {
            typeUrl: string;
            value: unknown;
        };
    };
    fromJSON: {
        submitClaim(value: any): {
            typeUrl: string;
            value: MsgSubmitClaim;
        };
    };
    fromPartial: {
        submitClaim(value: MsgSubmitClaim): {
            typeUrl: string;
            value: MsgSubmitClaim;
        };
    };
};
