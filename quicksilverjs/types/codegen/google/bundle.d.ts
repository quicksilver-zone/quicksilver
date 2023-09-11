import * as _73 from "./api/http";
import * as _74 from "./api/httpbody";
import * as _75 from "./protobuf/any";
import * as _76 from "./protobuf/descriptor";
import * as _77 from "./protobuf/duration";
import * as _78 from "./protobuf/timestamp";
export declare namespace google {
    const api: {
        HttpBody: {
            encode(message: _74.HttpBody, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _74.HttpBody;
            fromJSON(object: any): _74.HttpBody;
            toJSON(message: _74.HttpBody): unknown;
            fromPartial(object: Partial<_74.HttpBody>): _74.HttpBody;
        };
        Http: {
            encode(message: _73.Http, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _73.Http;
            fromJSON(object: any): _73.Http;
            toJSON(message: _73.Http): unknown;
            fromPartial(object: Partial<_73.Http>): _73.Http;
        };
        HttpRule: {
            encode(message: _73.HttpRule, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _73.HttpRule;
            fromJSON(object: any): _73.HttpRule;
            toJSON(message: _73.HttpRule): unknown;
            fromPartial(object: Partial<_73.HttpRule>): _73.HttpRule;
        };
        CustomHttpPattern: {
            encode(message: _73.CustomHttpPattern, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _73.CustomHttpPattern;
            fromJSON(object: any): _73.CustomHttpPattern;
            toJSON(message: _73.CustomHttpPattern): unknown;
            fromPartial(object: Partial<_73.CustomHttpPattern>): _73.CustomHttpPattern;
        };
    };
    const protobuf: {
        Timestamp: {
            encode(message: _78.Timestamp, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _78.Timestamp;
            fromJSON(object: any): _78.Timestamp;
            toJSON(message: _78.Timestamp): unknown;
            fromPartial(object: Partial<_78.Timestamp>): _78.Timestamp;
        };
        Duration: {
            encode(message: _77.Duration, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _77.Duration;
            fromJSON(object: any): _77.Duration;
            toJSON(message: _77.Duration): unknown;
            fromPartial(object: Partial<_77.Duration>): _77.Duration;
        };
        fieldDescriptorProto_TypeFromJSON(object: any): _76.FieldDescriptorProto_Type;
        fieldDescriptorProto_TypeToJSON(object: _76.FieldDescriptorProto_Type): string;
        fieldDescriptorProto_LabelFromJSON(object: any): _76.FieldDescriptorProto_Label;
        fieldDescriptorProto_LabelToJSON(object: _76.FieldDescriptorProto_Label): string;
        fileOptions_OptimizeModeFromJSON(object: any): _76.FileOptions_OptimizeMode;
        fileOptions_OptimizeModeToJSON(object: _76.FileOptions_OptimizeMode): string;
        fieldOptions_CTypeFromJSON(object: any): _76.FieldOptions_CType;
        fieldOptions_CTypeToJSON(object: _76.FieldOptions_CType): string;
        fieldOptions_JSTypeFromJSON(object: any): _76.FieldOptions_JSType;
        fieldOptions_JSTypeToJSON(object: _76.FieldOptions_JSType): string;
        methodOptions_IdempotencyLevelFromJSON(object: any): _76.MethodOptions_IdempotencyLevel;
        methodOptions_IdempotencyLevelToJSON(object: _76.MethodOptions_IdempotencyLevel): string;
        FieldDescriptorProto_Type: typeof _76.FieldDescriptorProto_Type;
        FieldDescriptorProto_TypeSDKType: typeof _76.FieldDescriptorProto_TypeSDKType;
        FieldDescriptorProto_Label: typeof _76.FieldDescriptorProto_Label;
        FieldDescriptorProto_LabelSDKType: typeof _76.FieldDescriptorProto_LabelSDKType;
        FileOptions_OptimizeMode: typeof _76.FileOptions_OptimizeMode;
        FileOptions_OptimizeModeSDKType: typeof _76.FileOptions_OptimizeModeSDKType;
        FieldOptions_CType: typeof _76.FieldOptions_CType;
        FieldOptions_CTypeSDKType: typeof _76.FieldOptions_CTypeSDKType;
        FieldOptions_JSType: typeof _76.FieldOptions_JSType;
        FieldOptions_JSTypeSDKType: typeof _76.FieldOptions_JSTypeSDKType;
        MethodOptions_IdempotencyLevel: typeof _76.MethodOptions_IdempotencyLevel;
        MethodOptions_IdempotencyLevelSDKType: typeof _76.MethodOptions_IdempotencyLevelSDKType;
        FileDescriptorSet: {
            encode(message: _76.FileDescriptorSet, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.FileDescriptorSet;
            fromJSON(object: any): _76.FileDescriptorSet;
            toJSON(message: _76.FileDescriptorSet): unknown;
            fromPartial(object: Partial<_76.FileDescriptorSet>): _76.FileDescriptorSet;
        };
        FileDescriptorProto: {
            encode(message: _76.FileDescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.FileDescriptorProto;
            fromJSON(object: any): _76.FileDescriptorProto;
            toJSON(message: _76.FileDescriptorProto): unknown;
            fromPartial(object: Partial<_76.FileDescriptorProto>): _76.FileDescriptorProto;
        };
        DescriptorProto: {
            encode(message: _76.DescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.DescriptorProto;
            fromJSON(object: any): _76.DescriptorProto;
            toJSON(message: _76.DescriptorProto): unknown;
            fromPartial(object: Partial<_76.DescriptorProto>): _76.DescriptorProto;
        };
        DescriptorProto_ExtensionRange: {
            encode(message: _76.DescriptorProto_ExtensionRange, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.DescriptorProto_ExtensionRange;
            fromJSON(object: any): _76.DescriptorProto_ExtensionRange;
            toJSON(message: _76.DescriptorProto_ExtensionRange): unknown;
            fromPartial(object: Partial<_76.DescriptorProto_ExtensionRange>): _76.DescriptorProto_ExtensionRange;
        };
        DescriptorProto_ReservedRange: {
            encode(message: _76.DescriptorProto_ReservedRange, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.DescriptorProto_ReservedRange;
            fromJSON(object: any): _76.DescriptorProto_ReservedRange;
            toJSON(message: _76.DescriptorProto_ReservedRange): unknown;
            fromPartial(object: Partial<_76.DescriptorProto_ReservedRange>): _76.DescriptorProto_ReservedRange;
        };
        ExtensionRangeOptions: {
            encode(message: _76.ExtensionRangeOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.ExtensionRangeOptions;
            fromJSON(object: any): _76.ExtensionRangeOptions;
            toJSON(message: _76.ExtensionRangeOptions): unknown;
            fromPartial(object: Partial<_76.ExtensionRangeOptions>): _76.ExtensionRangeOptions;
        };
        FieldDescriptorProto: {
            encode(message: _76.FieldDescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.FieldDescriptorProto;
            fromJSON(object: any): _76.FieldDescriptorProto;
            toJSON(message: _76.FieldDescriptorProto): unknown;
            fromPartial(object: Partial<_76.FieldDescriptorProto>): _76.FieldDescriptorProto;
        };
        OneofDescriptorProto: {
            encode(message: _76.OneofDescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.OneofDescriptorProto;
            fromJSON(object: any): _76.OneofDescriptorProto;
            toJSON(message: _76.OneofDescriptorProto): unknown;
            fromPartial(object: Partial<_76.OneofDescriptorProto>): _76.OneofDescriptorProto;
        };
        EnumDescriptorProto: {
            encode(message: _76.EnumDescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.EnumDescriptorProto;
            fromJSON(object: any): _76.EnumDescriptorProto;
            toJSON(message: _76.EnumDescriptorProto): unknown;
            fromPartial(object: Partial<_76.EnumDescriptorProto>): _76.EnumDescriptorProto;
        };
        EnumDescriptorProto_EnumReservedRange: {
            encode(message: _76.EnumDescriptorProto_EnumReservedRange, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.EnumDescriptorProto_EnumReservedRange;
            fromJSON(object: any): _76.EnumDescriptorProto_EnumReservedRange;
            toJSON(message: _76.EnumDescriptorProto_EnumReservedRange): unknown;
            fromPartial(object: Partial<_76.EnumDescriptorProto_EnumReservedRange>): _76.EnumDescriptorProto_EnumReservedRange;
        };
        EnumValueDescriptorProto: {
            encode(message: _76.EnumValueDescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.EnumValueDescriptorProto;
            fromJSON(object: any): _76.EnumValueDescriptorProto;
            toJSON(message: _76.EnumValueDescriptorProto): unknown;
            fromPartial(object: Partial<_76.EnumValueDescriptorProto>): _76.EnumValueDescriptorProto;
        };
        ServiceDescriptorProto: {
            encode(message: _76.ServiceDescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.ServiceDescriptorProto;
            fromJSON(object: any): _76.ServiceDescriptorProto;
            toJSON(message: _76.ServiceDescriptorProto): unknown;
            fromPartial(object: Partial<_76.ServiceDescriptorProto>): _76.ServiceDescriptorProto;
        };
        MethodDescriptorProto: {
            encode(message: _76.MethodDescriptorProto, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.MethodDescriptorProto;
            fromJSON(object: any): _76.MethodDescriptorProto;
            toJSON(message: _76.MethodDescriptorProto): unknown;
            fromPartial(object: Partial<_76.MethodDescriptorProto>): _76.MethodDescriptorProto;
        };
        FileOptions: {
            encode(message: _76.FileOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.FileOptions;
            fromJSON(object: any): _76.FileOptions;
            toJSON(message: _76.FileOptions): unknown;
            fromPartial(object: Partial<_76.FileOptions>): _76.FileOptions;
        };
        MessageOptions: {
            encode(message: _76.MessageOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.MessageOptions;
            fromJSON(object: any): _76.MessageOptions;
            toJSON(message: _76.MessageOptions): unknown;
            fromPartial(object: Partial<_76.MessageOptions>): _76.MessageOptions;
        };
        FieldOptions: {
            encode(message: _76.FieldOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.FieldOptions;
            fromJSON(object: any): _76.FieldOptions;
            toJSON(message: _76.FieldOptions): unknown;
            fromPartial(object: Partial<_76.FieldOptions>): _76.FieldOptions;
        };
        OneofOptions: {
            encode(message: _76.OneofOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.OneofOptions;
            fromJSON(object: any): _76.OneofOptions;
            toJSON(message: _76.OneofOptions): unknown;
            fromPartial(object: Partial<_76.OneofOptions>): _76.OneofOptions;
        };
        EnumOptions: {
            encode(message: _76.EnumOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.EnumOptions;
            fromJSON(object: any): _76.EnumOptions;
            toJSON(message: _76.EnumOptions): unknown;
            fromPartial(object: Partial<_76.EnumOptions>): _76.EnumOptions;
        };
        EnumValueOptions: {
            encode(message: _76.EnumValueOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.EnumValueOptions;
            fromJSON(object: any): _76.EnumValueOptions;
            toJSON(message: _76.EnumValueOptions): unknown;
            fromPartial(object: Partial<_76.EnumValueOptions>): _76.EnumValueOptions;
        };
        ServiceOptions: {
            encode(message: _76.ServiceOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.ServiceOptions;
            fromJSON(object: any): _76.ServiceOptions;
            toJSON(message: _76.ServiceOptions): unknown;
            fromPartial(object: Partial<_76.ServiceOptions>): _76.ServiceOptions;
        };
        MethodOptions: {
            encode(message: _76.MethodOptions, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.MethodOptions;
            fromJSON(object: any): _76.MethodOptions;
            toJSON(message: _76.MethodOptions): unknown;
            fromPartial(object: Partial<_76.MethodOptions>): _76.MethodOptions;
        };
        UninterpretedOption: {
            encode(message: _76.UninterpretedOption, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.UninterpretedOption;
            fromJSON(object: any): _76.UninterpretedOption;
            toJSON(message: _76.UninterpretedOption): unknown;
            fromPartial(object: Partial<_76.UninterpretedOption>): _76.UninterpretedOption;
        };
        UninterpretedOption_NamePart: {
            encode(message: _76.UninterpretedOption_NamePart, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.UninterpretedOption_NamePart;
            fromJSON(object: any): _76.UninterpretedOption_NamePart;
            toJSON(message: _76.UninterpretedOption_NamePart): unknown;
            fromPartial(object: Partial<_76.UninterpretedOption_NamePart>): _76.UninterpretedOption_NamePart;
        };
        SourceCodeInfo: {
            encode(message: _76.SourceCodeInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.SourceCodeInfo;
            fromJSON(object: any): _76.SourceCodeInfo;
            toJSON(message: _76.SourceCodeInfo): unknown;
            fromPartial(object: Partial<_76.SourceCodeInfo>): _76.SourceCodeInfo;
        };
        SourceCodeInfo_Location: {
            encode(message: _76.SourceCodeInfo_Location, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.SourceCodeInfo_Location;
            fromJSON(object: any): _76.SourceCodeInfo_Location;
            toJSON(message: _76.SourceCodeInfo_Location): unknown;
            fromPartial(object: Partial<_76.SourceCodeInfo_Location>): _76.SourceCodeInfo_Location;
        };
        GeneratedCodeInfo: {
            encode(message: _76.GeneratedCodeInfo, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.GeneratedCodeInfo;
            fromJSON(object: any): _76.GeneratedCodeInfo;
            toJSON(message: _76.GeneratedCodeInfo): unknown;
            fromPartial(object: Partial<_76.GeneratedCodeInfo>): _76.GeneratedCodeInfo;
        };
        GeneratedCodeInfo_Annotation: {
            encode(message: _76.GeneratedCodeInfo_Annotation, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _76.GeneratedCodeInfo_Annotation;
            fromJSON(object: any): _76.GeneratedCodeInfo_Annotation;
            toJSON(message: _76.GeneratedCodeInfo_Annotation): unknown;
            fromPartial(object: Partial<_76.GeneratedCodeInfo_Annotation>): _76.GeneratedCodeInfo_Annotation;
        };
        Any: {
            encode(message: _75.Any, writer?: import("protobufjs").Writer): import("protobufjs").Writer;
            decode(input: Uint8Array | import("protobufjs").Reader, length?: number): _75.Any;
            fromJSON(object: any): _75.Any;
            toJSON(message: _75.Any): unknown;
            fromPartial(object: Partial<_75.Any>): _75.Any;
        };
    };
}
