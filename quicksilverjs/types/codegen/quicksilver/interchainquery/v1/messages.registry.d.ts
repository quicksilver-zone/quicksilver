import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgSubmitQueryResponse } from "./messages";
export declare const registry: ReadonlyArray<[string, GeneratedType]>;
export declare const load: (protoRegistry: Registry) => void;
export declare const MessageComposer: {
    encoded: {
        submitQueryResponse(value: MsgSubmitQueryResponse): {
            typeUrl: string;
            value: Uint8Array;
        };
    };
    withTypeUrl: {
        submitQueryResponse(value: MsgSubmitQueryResponse): {
            typeUrl: string;
            value: MsgSubmitQueryResponse;
        };
    };
    toJSON: {
        submitQueryResponse(value: MsgSubmitQueryResponse): {
            typeUrl: string;
            value: unknown;
        };
    };
    fromJSON: {
        submitQueryResponse(value: any): {
            typeUrl: string;
            value: MsgSubmitQueryResponse;
        };
    };
    fromPartial: {
        submitQueryResponse(value: MsgSubmitQueryResponse): {
            typeUrl: string;
            value: MsgSubmitQueryResponse;
        };
    };
};
