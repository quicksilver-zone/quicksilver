import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgClaim } from "./messages";
export declare const registry: ReadonlyArray<[string, GeneratedType]>;
export declare const load: (protoRegistry: Registry) => void;
export declare const MessageComposer: {
    encoded: {
        claim(value: MsgClaim): {
            typeUrl: string;
            value: Uint8Array;
        };
    };
    withTypeUrl: {
        claim(value: MsgClaim): {
            typeUrl: string;
            value: MsgClaim;
        };
    };
    toJSON: {
        claim(value: MsgClaim): {
            typeUrl: string;
            value: unknown;
        };
    };
    fromJSON: {
        claim(value: any): {
            typeUrl: string;
            value: MsgClaim;
        };
    };
    fromPartial: {
        claim(value: MsgClaim): {
            typeUrl: string;
            value: MsgClaim;
        };
    };
};
