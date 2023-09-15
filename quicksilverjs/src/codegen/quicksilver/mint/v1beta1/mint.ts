import { Long, isSet, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.mint.v1beta1";
/** Minter represents the minting state. */
export interface Minter {
  /** current epoch provisions */
  epochProvisions: string;
}
export interface MinterProtoMsg {
  typeUrl: "/quicksilver.mint.v1beta1.Minter";
  value: Uint8Array;
}
/** Minter represents the minting state. */
export interface MinterAmino {
  /** current epoch provisions */
  epoch_provisions: string;
}
export interface MinterAminoMsg {
  type: "/quicksilver.mint.v1beta1.Minter";
  value: MinterAmino;
}
/** Minter represents the minting state. */
export interface MinterSDKType {
  epoch_provisions: string;
}
export interface DistributionProportions {
  /**
   * staking defines the proportion of the minted minted_denom that is to be
   * allocated as staking rewards.
   */
  staking: string;
  /**
   * pool_incentives defines the proportion of the minted minted_denom that is
   * to be allocated as pool incentives.
   */
  poolIncentives: string;
  /**
   * participation_rewards defines the proportion of the minted minted_denom
   * that is to be allocated to participation rewards address.
   */
  participationRewards: string;
  /**
   * community_pool defines the proportion of the minted minted_denom that is
   * to be allocated to the community pool.
   */
  communityPool: string;
}
export interface DistributionProportionsProtoMsg {
  typeUrl: "/quicksilver.mint.v1beta1.DistributionProportions";
  value: Uint8Array;
}
export interface DistributionProportionsAmino {
  /**
   * staking defines the proportion of the minted minted_denom that is to be
   * allocated as staking rewards.
   */
  staking: string;
  /**
   * pool_incentives defines the proportion of the minted minted_denom that is
   * to be allocated as pool incentives.
   */
  pool_incentives: string;
  /**
   * participation_rewards defines the proportion of the minted minted_denom
   * that is to be allocated to participation rewards address.
   */
  participation_rewards: string;
  /**
   * community_pool defines the proportion of the minted minted_denom that is
   * to be allocated to the community pool.
   */
  community_pool: string;
}
export interface DistributionProportionsAminoMsg {
  type: "/quicksilver.mint.v1beta1.DistributionProportions";
  value: DistributionProportionsAmino;
}
export interface DistributionProportionsSDKType {
  staking: string;
  pool_incentives: string;
  participation_rewards: string;
  community_pool: string;
}
/** Params holds parameters for the mint module. */
export interface Params {
  /** type of coin to mint */
  mintDenom: string;
  /** epoch provisions from the first epoch */
  genesisEpochProvisions: string;
  /** mint epoch identifier */
  epochIdentifier: string;
  /** number of epochs take to reduce rewards */
  reductionPeriodInEpochs: Long;
  /** reduction multiplier to execute on each period */
  reductionFactor: string;
  /** distribution_proportions defines the proportion of the minted denom */
  distributionProportions: DistributionProportions;
  /** start epoch to distribute minting rewards */
  mintingRewardsDistributionStartEpoch: Long;
}
export interface ParamsProtoMsg {
  typeUrl: "/quicksilver.mint.v1beta1.Params";
  value: Uint8Array;
}
/** Params holds parameters for the mint module. */
export interface ParamsAmino {
  /** type of coin to mint */
  mint_denom: string;
  /** epoch provisions from the first epoch */
  genesis_epoch_provisions: string;
  /** mint epoch identifier */
  epoch_identifier: string;
  /** number of epochs take to reduce rewards */
  reduction_period_in_epochs: string;
  /** reduction multiplier to execute on each period */
  reduction_factor: string;
  /** distribution_proportions defines the proportion of the minted denom */
  distribution_proportions?: DistributionProportionsAmino;
  /** start epoch to distribute minting rewards */
  minting_rewards_distribution_start_epoch: string;
}
export interface ParamsAminoMsg {
  type: "/quicksilver.mint.v1beta1.Params";
  value: ParamsAmino;
}
/** Params holds parameters for the mint module. */
export interface ParamsSDKType {
  mint_denom: string;
  genesis_epoch_provisions: string;
  epoch_identifier: string;
  reduction_period_in_epochs: Long;
  reduction_factor: string;
  distribution_proportions: DistributionProportionsSDKType;
  minting_rewards_distribution_start_epoch: Long;
}
function createBaseMinter(): Minter {
  return {
    epochProvisions: ""
  };
}
export const Minter = {
  typeUrl: "/quicksilver.mint.v1beta1.Minter",
  encode(message: Minter, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.epochProvisions !== "") {
      writer.uint32(10).string(message.epochProvisions);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Minter {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMinter();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.epochProvisions = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Minter {
    const obj = createBaseMinter();
    if (isSet(object.epochProvisions)) obj.epochProvisions = String(object.epochProvisions);
    return obj;
  },
  toJSON(message: Minter): unknown {
    const obj: any = {};
    message.epochProvisions !== undefined && (obj.epochProvisions = message.epochProvisions);
    return obj;
  },
  fromPartial(object: DeepPartial<Minter>): Minter {
    const message = createBaseMinter();
    message.epochProvisions = object.epochProvisions ?? "";
    return message;
  },
  fromSDK(object: MinterSDKType): Minter {
    return {
      epochProvisions: object?.epoch_provisions
    };
  },
  toSDK(message: Minter): MinterSDKType {
    const obj: any = {};
    obj.epoch_provisions = message.epochProvisions;
    return obj;
  },
  fromAmino(object: MinterAmino): Minter {
    return {
      epochProvisions: object.epoch_provisions
    };
  },
  toAmino(message: Minter): MinterAmino {
    const obj: any = {};
    obj.epoch_provisions = message.epochProvisions;
    return obj;
  },
  fromAminoMsg(object: MinterAminoMsg): Minter {
    return Minter.fromAmino(object.value);
  },
  fromProtoMsg(message: MinterProtoMsg): Minter {
    return Minter.decode(message.value);
  },
  toProto(message: Minter): Uint8Array {
    return Minter.encode(message).finish();
  },
  toProtoMsg(message: Minter): MinterProtoMsg {
    return {
      typeUrl: "/quicksilver.mint.v1beta1.Minter",
      value: Minter.encode(message).finish()
    };
  }
};
function createBaseDistributionProportions(): DistributionProportions {
  return {
    staking: "",
    poolIncentives: "",
    participationRewards: "",
    communityPool: ""
  };
}
export const DistributionProportions = {
  typeUrl: "/quicksilver.mint.v1beta1.DistributionProportions",
  encode(message: DistributionProportions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.staking !== "") {
      writer.uint32(10).string(message.staking);
    }
    if (message.poolIncentives !== "") {
      writer.uint32(18).string(message.poolIncentives);
    }
    if (message.participationRewards !== "") {
      writer.uint32(26).string(message.participationRewards);
    }
    if (message.communityPool !== "") {
      writer.uint32(34).string(message.communityPool);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): DistributionProportions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDistributionProportions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.staking = reader.string();
          break;
        case 2:
          message.poolIncentives = reader.string();
          break;
        case 3:
          message.participationRewards = reader.string();
          break;
        case 4:
          message.communityPool = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): DistributionProportions {
    const obj = createBaseDistributionProportions();
    if (isSet(object.staking)) obj.staking = String(object.staking);
    if (isSet(object.poolIncentives)) obj.poolIncentives = String(object.poolIncentives);
    if (isSet(object.participationRewards)) obj.participationRewards = String(object.participationRewards);
    if (isSet(object.communityPool)) obj.communityPool = String(object.communityPool);
    return obj;
  },
  toJSON(message: DistributionProportions): unknown {
    const obj: any = {};
    message.staking !== undefined && (obj.staking = message.staking);
    message.poolIncentives !== undefined && (obj.poolIncentives = message.poolIncentives);
    message.participationRewards !== undefined && (obj.participationRewards = message.participationRewards);
    message.communityPool !== undefined && (obj.communityPool = message.communityPool);
    return obj;
  },
  fromPartial(object: DeepPartial<DistributionProportions>): DistributionProportions {
    const message = createBaseDistributionProportions();
    message.staking = object.staking ?? "";
    message.poolIncentives = object.poolIncentives ?? "";
    message.participationRewards = object.participationRewards ?? "";
    message.communityPool = object.communityPool ?? "";
    return message;
  },
  fromSDK(object: DistributionProportionsSDKType): DistributionProportions {
    return {
      staking: object?.staking,
      poolIncentives: object?.pool_incentives,
      participationRewards: object?.participation_rewards,
      communityPool: object?.community_pool
    };
  },
  toSDK(message: DistributionProportions): DistributionProportionsSDKType {
    const obj: any = {};
    obj.staking = message.staking;
    obj.pool_incentives = message.poolIncentives;
    obj.participation_rewards = message.participationRewards;
    obj.community_pool = message.communityPool;
    return obj;
  },
  fromAmino(object: DistributionProportionsAmino): DistributionProportions {
    return {
      staking: object.staking,
      poolIncentives: object.pool_incentives,
      participationRewards: object.participation_rewards,
      communityPool: object.community_pool
    };
  },
  toAmino(message: DistributionProportions): DistributionProportionsAmino {
    const obj: any = {};
    obj.staking = message.staking;
    obj.pool_incentives = message.poolIncentives;
    obj.participation_rewards = message.participationRewards;
    obj.community_pool = message.communityPool;
    return obj;
  },
  fromAminoMsg(object: DistributionProportionsAminoMsg): DistributionProportions {
    return DistributionProportions.fromAmino(object.value);
  },
  fromProtoMsg(message: DistributionProportionsProtoMsg): DistributionProportions {
    return DistributionProportions.decode(message.value);
  },
  toProto(message: DistributionProportions): Uint8Array {
    return DistributionProportions.encode(message).finish();
  },
  toProtoMsg(message: DistributionProportions): DistributionProportionsProtoMsg {
    return {
      typeUrl: "/quicksilver.mint.v1beta1.DistributionProportions",
      value: DistributionProportions.encode(message).finish()
    };
  }
};
function createBaseParams(): Params {
  return {
    mintDenom: "",
    genesisEpochProvisions: "",
    epochIdentifier: "",
    reductionPeriodInEpochs: Long.ZERO,
    reductionFactor: "",
    distributionProportions: DistributionProportions.fromPartial({}),
    mintingRewardsDistributionStartEpoch: Long.ZERO
  };
}
export const Params = {
  typeUrl: "/quicksilver.mint.v1beta1.Params",
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.mintDenom !== "") {
      writer.uint32(10).string(message.mintDenom);
    }
    if (message.genesisEpochProvisions !== "") {
      writer.uint32(18).string(message.genesisEpochProvisions);
    }
    if (message.epochIdentifier !== "") {
      writer.uint32(26).string(message.epochIdentifier);
    }
    if (!message.reductionPeriodInEpochs.isZero()) {
      writer.uint32(32).int64(message.reductionPeriodInEpochs);
    }
    if (message.reductionFactor !== "") {
      writer.uint32(42).string(message.reductionFactor);
    }
    if (message.distributionProportions !== undefined) {
      DistributionProportions.encode(message.distributionProportions, writer.uint32(50).fork()).ldelim();
    }
    if (!message.mintingRewardsDistributionStartEpoch.isZero()) {
      writer.uint32(56).int64(message.mintingRewardsDistributionStartEpoch);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mintDenom = reader.string();
          break;
        case 2:
          message.genesisEpochProvisions = reader.string();
          break;
        case 3:
          message.epochIdentifier = reader.string();
          break;
        case 4:
          message.reductionPeriodInEpochs = (reader.int64() as Long);
          break;
        case 5:
          message.reductionFactor = reader.string();
          break;
        case 6:
          message.distributionProportions = DistributionProportions.decode(reader, reader.uint32());
          break;
        case 7:
          message.mintingRewardsDistributionStartEpoch = (reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Params {
    const obj = createBaseParams();
    if (isSet(object.mintDenom)) obj.mintDenom = String(object.mintDenom);
    if (isSet(object.genesisEpochProvisions)) obj.genesisEpochProvisions = String(object.genesisEpochProvisions);
    if (isSet(object.epochIdentifier)) obj.epochIdentifier = String(object.epochIdentifier);
    if (isSet(object.reductionPeriodInEpochs)) obj.reductionPeriodInEpochs = Long.fromValue(object.reductionPeriodInEpochs);
    if (isSet(object.reductionFactor)) obj.reductionFactor = String(object.reductionFactor);
    if (isSet(object.distributionProportions)) obj.distributionProportions = DistributionProportions.fromJSON(object.distributionProportions);
    if (isSet(object.mintingRewardsDistributionStartEpoch)) obj.mintingRewardsDistributionStartEpoch = Long.fromValue(object.mintingRewardsDistributionStartEpoch);
    return obj;
  },
  toJSON(message: Params): unknown {
    const obj: any = {};
    message.mintDenom !== undefined && (obj.mintDenom = message.mintDenom);
    message.genesisEpochProvisions !== undefined && (obj.genesisEpochProvisions = message.genesisEpochProvisions);
    message.epochIdentifier !== undefined && (obj.epochIdentifier = message.epochIdentifier);
    message.reductionPeriodInEpochs !== undefined && (obj.reductionPeriodInEpochs = (message.reductionPeriodInEpochs || Long.ZERO).toString());
    message.reductionFactor !== undefined && (obj.reductionFactor = message.reductionFactor);
    message.distributionProportions !== undefined && (obj.distributionProportions = message.distributionProportions ? DistributionProportions.toJSON(message.distributionProportions) : undefined);
    message.mintingRewardsDistributionStartEpoch !== undefined && (obj.mintingRewardsDistributionStartEpoch = (message.mintingRewardsDistributionStartEpoch || Long.ZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<Params>): Params {
    const message = createBaseParams();
    message.mintDenom = object.mintDenom ?? "";
    message.genesisEpochProvisions = object.genesisEpochProvisions ?? "";
    message.epochIdentifier = object.epochIdentifier ?? "";
    if (object.reductionPeriodInEpochs !== undefined && object.reductionPeriodInEpochs !== null) {
      message.reductionPeriodInEpochs = Long.fromValue(object.reductionPeriodInEpochs);
    }
    message.reductionFactor = object.reductionFactor ?? "";
    if (object.distributionProportions !== undefined && object.distributionProportions !== null) {
      message.distributionProportions = DistributionProportions.fromPartial(object.distributionProportions);
    }
    if (object.mintingRewardsDistributionStartEpoch !== undefined && object.mintingRewardsDistributionStartEpoch !== null) {
      message.mintingRewardsDistributionStartEpoch = Long.fromValue(object.mintingRewardsDistributionStartEpoch);
    }
    return message;
  },
  fromSDK(object: ParamsSDKType): Params {
    return {
      mintDenom: object?.mint_denom,
      genesisEpochProvisions: object?.genesis_epoch_provisions,
      epochIdentifier: object?.epoch_identifier,
      reductionPeriodInEpochs: object?.reduction_period_in_epochs,
      reductionFactor: object?.reduction_factor,
      distributionProportions: object.distribution_proportions ? DistributionProportions.fromSDK(object.distribution_proportions) : undefined,
      mintingRewardsDistributionStartEpoch: object?.minting_rewards_distribution_start_epoch
    };
  },
  toSDK(message: Params): ParamsSDKType {
    const obj: any = {};
    obj.mint_denom = message.mintDenom;
    obj.genesis_epoch_provisions = message.genesisEpochProvisions;
    obj.epoch_identifier = message.epochIdentifier;
    obj.reduction_period_in_epochs = message.reductionPeriodInEpochs;
    obj.reduction_factor = message.reductionFactor;
    message.distributionProportions !== undefined && (obj.distribution_proportions = message.distributionProportions ? DistributionProportions.toSDK(message.distributionProportions) : undefined);
    obj.minting_rewards_distribution_start_epoch = message.mintingRewardsDistributionStartEpoch;
    return obj;
  },
  fromAmino(object: ParamsAmino): Params {
    return {
      mintDenom: object.mint_denom,
      genesisEpochProvisions: object.genesis_epoch_provisions,
      epochIdentifier: object.epoch_identifier,
      reductionPeriodInEpochs: Long.fromString(object.reduction_period_in_epochs),
      reductionFactor: object.reduction_factor,
      distributionProportions: object?.distribution_proportions ? DistributionProportions.fromAmino(object.distribution_proportions) : undefined,
      mintingRewardsDistributionStartEpoch: Long.fromString(object.minting_rewards_distribution_start_epoch)
    };
  },
  toAmino(message: Params): ParamsAmino {
    const obj: any = {};
    obj.mint_denom = message.mintDenom;
    obj.genesis_epoch_provisions = message.genesisEpochProvisions;
    obj.epoch_identifier = message.epochIdentifier;
    obj.reduction_period_in_epochs = message.reductionPeriodInEpochs ? message.reductionPeriodInEpochs.toString() : undefined;
    obj.reduction_factor = message.reductionFactor;
    obj.distribution_proportions = message.distributionProportions ? DistributionProportions.toAmino(message.distributionProportions) : undefined;
    obj.minting_rewards_distribution_start_epoch = message.mintingRewardsDistributionStartEpoch ? message.mintingRewardsDistributionStartEpoch.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: ParamsAminoMsg): Params {
    return Params.fromAmino(object.value);
  },
  fromProtoMsg(message: ParamsProtoMsg): Params {
    return Params.decode(message.value);
  },
  toProto(message: Params): Uint8Array {
    return Params.encode(message).finish();
  },
  toProtoMsg(message: Params): ParamsProtoMsg {
    return {
      typeUrl: "/quicksilver.mint.v1beta1.Params",
      value: Params.encode(message).finish()
    };
  }
};