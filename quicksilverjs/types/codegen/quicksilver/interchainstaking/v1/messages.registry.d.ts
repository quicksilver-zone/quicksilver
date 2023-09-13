import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgRequestRedemption, MsgSignalIntent } from "./messages";
export declare const registry: ReadonlyArray<[string, GeneratedType]>;
export declare const load: (protoRegistry: Registry) => void;
export declare const MessageComposer: {
    encoded: {
        requestRedemption(value: MsgRequestRedemption): {
            typeUrl: string;
            value: Uint8Array;
        };
        signalIntent(value: MsgSignalIntent): {
            typeUrl: string;
            value: Uint8Array;
        };
    };
    withTypeUrl: {
        requestRedemption(value: MsgRequestRedemption): {
            typeUrl: string;
            value: MsgRequestRedemption;
        };
        signalIntent(value: MsgSignalIntent): {
            typeUrl: string;
            value: MsgSignalIntent;
        };
    };
    toJSON: {
        requestRedemption(value: MsgRequestRedemption): {
            typeUrl: string;
            value: unknown;
        };
        signalIntent(value: MsgSignalIntent): {
            typeUrl: string;
            value: unknown;
        };
    };
    fromJSON: {
        requestRedemption(value: any): {
            typeUrl: string;
            value: MsgRequestRedemption;
        };
        signalIntent(value: any): {
            typeUrl: string;
            value: MsgSignalIntent;
        };
    };
    fromPartial: {
        requestRedemption(value: MsgRequestRedemption): {
            typeUrl: string;
            value: MsgRequestRedemption;
        };
        signalIntent(value: MsgSignalIntent): {
            typeUrl: string;
            value: MsgSignalIntent;
        };
    };
};
