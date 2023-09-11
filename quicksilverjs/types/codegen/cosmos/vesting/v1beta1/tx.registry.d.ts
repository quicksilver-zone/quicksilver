import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgCreateVestingAccount } from "./tx";
export declare const registry: ReadonlyArray<[string, GeneratedType]>;
export declare const load: (protoRegistry: Registry) => void;
export declare const MessageComposer: {
    encoded: {
        createVestingAccount(value: MsgCreateVestingAccount): {
            typeUrl: string;
            value: Uint8Array;
        };
    };
    withTypeUrl: {
        createVestingAccount(value: MsgCreateVestingAccount): {
            typeUrl: string;
            value: MsgCreateVestingAccount;
        };
    };
    toJSON: {
        createVestingAccount(value: MsgCreateVestingAccount): {
            typeUrl: string;
            value: unknown;
        };
    };
    fromJSON: {
        createVestingAccount(value: any): {
            typeUrl: string;
            value: MsgCreateVestingAccount;
        };
    };
    fromPartial: {
        createVestingAccount(value: MsgCreateVestingAccount): {
            typeUrl: string;
            value: MsgCreateVestingAccount;
        };
    };
};
