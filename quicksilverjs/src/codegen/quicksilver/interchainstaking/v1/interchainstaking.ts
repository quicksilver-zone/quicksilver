import { Coin, CoinAmino, CoinSDKType } from "../../../cosmos/base/v1beta1/coin";
import { Timestamp, TimestampAmino, TimestampSDKType } from "../../../google/protobuf/timestamp";
import { Long, isSet, bytesFromBase64, base64FromBytes, DeepPartial, toTimestamp, fromTimestamp } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
export const protobufPackage = "quicksilver.interchainstaking.v1";
export interface Zone {
  connectionId: string;
  chainId: string;
  depositAddress: ICAAccount;
  withdrawalAddress: ICAAccount;
  performanceAddress: ICAAccount;
  delegationAddress: ICAAccount;
  accountPrefix: string;
  localDenom: string;
  baseDenom: string;
  redemptionRate: string;
  lastRedemptionRate: string;
  validators: Validator[];
  aggregateIntent: ValidatorIntent[];
  multiSend: boolean;
  liquidityModule: boolean;
  withdrawalWaitgroup: number;
  ibcNextValidatorsHash: Uint8Array;
  validatorSelectionAllocation: Long;
  holdingsAllocation: Long;
  /** deprecated remove me. */
  lastEpochHeight: Long;
  tvl: string;
  unbondingPeriod: Long;
  messagesPerTx: Long;
}
export interface ZoneProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.Zone";
  value: Uint8Array;
}
export interface ZoneAmino {
  connection_id: string;
  chain_id: string;
  deposit_address?: ICAAccountAmino;
  withdrawal_address?: ICAAccountAmino;
  performance_address?: ICAAccountAmino;
  delegation_address?: ICAAccountAmino;
  account_prefix: string;
  local_denom: string;
  base_denom: string;
  redemption_rate: string;
  last_redemption_rate: string;
  validators: ValidatorAmino[];
  aggregate_intent: ValidatorIntentAmino[];
  multi_send: boolean;
  liquidity_module: boolean;
  withdrawal_waitgroup: number;
  ibc_next_validators_hash: Uint8Array;
  validator_selection_allocation: string;
  holdings_allocation: string;
  /** deprecated remove me. */
  last_epoch_height: string;
  tvl: string;
  unbonding_period: string;
  messages_per_tx: string;
}
export interface ZoneAminoMsg {
  type: "/quicksilver.interchainstaking.v1.Zone";
  value: ZoneAmino;
}
export interface ZoneSDKType {
  connection_id: string;
  chain_id: string;
  deposit_address: ICAAccountSDKType;
  withdrawal_address: ICAAccountSDKType;
  performance_address: ICAAccountSDKType;
  delegation_address: ICAAccountSDKType;
  account_prefix: string;
  local_denom: string;
  base_denom: string;
  redemption_rate: string;
  last_redemption_rate: string;
  validators: ValidatorSDKType[];
  aggregate_intent: ValidatorIntentSDKType[];
  multi_send: boolean;
  liquidity_module: boolean;
  withdrawal_waitgroup: number;
  ibc_next_validators_hash: Uint8Array;
  validator_selection_allocation: Long;
  holdings_allocation: Long;
  last_epoch_height: Long;
  tvl: string;
  unbonding_period: Long;
  messages_per_tx: Long;
}
export interface ICAAccount {
  address: string;
  /** balance defines the different coins this balance holds. */
  balance: Coin[];
  portName: string;
  withdrawalAddress: string;
  balanceWaitgroup: number;
}
export interface ICAAccountProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.ICAAccount";
  value: Uint8Array;
}
export interface ICAAccountAmino {
  address: string;
  /** balance defines the different coins this balance holds. */
  balance: CoinAmino[];
  port_name: string;
  withdrawal_address: string;
  balance_waitgroup: number;
}
export interface ICAAccountAminoMsg {
  type: "/quicksilver.interchainstaking.v1.ICAAccount";
  value: ICAAccountAmino;
}
export interface ICAAccountSDKType {
  address: string;
  balance: CoinSDKType[];
  port_name: string;
  withdrawal_address: string;
  balance_waitgroup: number;
}
export interface Distribution {
  valoper: string;
  amount: Long;
}
export interface DistributionProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.Distribution";
  value: Uint8Array;
}
export interface DistributionAmino {
  valoper: string;
  amount: string;
}
export interface DistributionAminoMsg {
  type: "/quicksilver.interchainstaking.v1.Distribution";
  value: DistributionAmino;
}
export interface DistributionSDKType {
  valoper: string;
  amount: Long;
}
export interface WithdrawalRecord {
  chainId: string;
  delegator: string;
  distribution: Distribution[];
  recipient: string;
  amount: Coin[];
  burnAmount: Coin;
  txhash: string;
  status: number;
  completionTime: Date;
}
export interface WithdrawalRecordProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.WithdrawalRecord";
  value: Uint8Array;
}
export interface WithdrawalRecordAmino {
  chain_id: string;
  delegator: string;
  distribution: DistributionAmino[];
  recipient: string;
  amount: CoinAmino[];
  burn_amount?: CoinAmino;
  txhash: string;
  status: number;
  completion_time?: Date;
}
export interface WithdrawalRecordAminoMsg {
  type: "/quicksilver.interchainstaking.v1.WithdrawalRecord";
  value: WithdrawalRecordAmino;
}
export interface WithdrawalRecordSDKType {
  chain_id: string;
  delegator: string;
  distribution: DistributionSDKType[];
  recipient: string;
  amount: CoinSDKType[];
  burn_amount: CoinSDKType;
  txhash: string;
  status: number;
  completion_time: Date;
}
export interface UnbondingRecord {
  chainId: string;
  epochNumber: Long;
  validator: string;
  relatedTxhash: string[];
}
export interface UnbondingRecordProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.UnbondingRecord";
  value: Uint8Array;
}
export interface UnbondingRecordAmino {
  chain_id: string;
  epoch_number: string;
  validator: string;
  related_txhash: string[];
}
export interface UnbondingRecordAminoMsg {
  type: "/quicksilver.interchainstaking.v1.UnbondingRecord";
  value: UnbondingRecordAmino;
}
export interface UnbondingRecordSDKType {
  chain_id: string;
  epoch_number: Long;
  validator: string;
  related_txhash: string[];
}
export interface RedelegationRecord {
  chainId: string;
  epochNumber: Long;
  source: string;
  destination: string;
  amount: Long;
  completionTime: Date;
}
export interface RedelegationRecordProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.RedelegationRecord";
  value: Uint8Array;
}
export interface RedelegationRecordAmino {
  chain_id: string;
  epoch_number: string;
  source: string;
  destination: string;
  amount: string;
  completion_time?: Date;
}
export interface RedelegationRecordAminoMsg {
  type: "/quicksilver.interchainstaking.v1.RedelegationRecord";
  value: RedelegationRecordAmino;
}
export interface RedelegationRecordSDKType {
  chain_id: string;
  epoch_number: Long;
  source: string;
  destination: string;
  amount: Long;
  completion_time: Date;
}
export interface TransferRecord {
  sender: string;
  recipient: string;
  amount: Coin;
}
export interface TransferRecordProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.TransferRecord";
  value: Uint8Array;
}
export interface TransferRecordAmino {
  sender: string;
  recipient: string;
  amount?: CoinAmino;
}
export interface TransferRecordAminoMsg {
  type: "/quicksilver.interchainstaking.v1.TransferRecord";
  value: TransferRecordAmino;
}
export interface TransferRecordSDKType {
  sender: string;
  recipient: string;
  amount: CoinSDKType;
}
export interface Validator {
  valoperAddress: string;
  commissionRate: string;
  delegatorShares: string;
  votingPower: string;
  score: string;
  status: string;
  jailed: boolean;
  tombstoned: boolean;
  jailedSince: Date;
}
export interface ValidatorProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.Validator";
  value: Uint8Array;
}
export interface ValidatorAmino {
  valoper_address: string;
  commission_rate: string;
  delegator_shares: string;
  voting_power: string;
  score: string;
  status: string;
  jailed: boolean;
  tombstoned: boolean;
  jailed_since?: Date;
}
export interface ValidatorAminoMsg {
  type: "/quicksilver.interchainstaking.v1.Validator";
  value: ValidatorAmino;
}
export interface ValidatorSDKType {
  valoper_address: string;
  commission_rate: string;
  delegator_shares: string;
  voting_power: string;
  score: string;
  status: string;
  jailed: boolean;
  tombstoned: boolean;
  jailed_since: Date;
}
export interface DelegatorIntent {
  delegator: string;
  intents: ValidatorIntent[];
}
export interface DelegatorIntentProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.DelegatorIntent";
  value: Uint8Array;
}
export interface DelegatorIntentAmino {
  delegator: string;
  intents: ValidatorIntentAmino[];
}
export interface DelegatorIntentAminoMsg {
  type: "/quicksilver.interchainstaking.v1.DelegatorIntent";
  value: DelegatorIntentAmino;
}
export interface DelegatorIntentSDKType {
  delegator: string;
  intents: ValidatorIntentSDKType[];
}
export interface ValidatorIntent {
  valoperAddress: string;
  weight: string;
}
export interface ValidatorIntentProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.ValidatorIntent";
  value: Uint8Array;
}
export interface ValidatorIntentAmino {
  valoper_address: string;
  weight: string;
}
export interface ValidatorIntentAminoMsg {
  type: "/quicksilver.interchainstaking.v1.ValidatorIntent";
  value: ValidatorIntentAmino;
}
export interface ValidatorIntentSDKType {
  valoper_address: string;
  weight: string;
}
export interface Delegation {
  delegationAddress: string;
  validatorAddress: string;
  amount: Coin;
  height: Long;
  redelegationEnd: Long;
}
export interface DelegationProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.Delegation";
  value: Uint8Array;
}
export interface DelegationAmino {
  delegation_address: string;
  validator_address: string;
  amount?: CoinAmino;
  height: string;
  redelegation_end: string;
}
export interface DelegationAminoMsg {
  type: "/quicksilver.interchainstaking.v1.Delegation";
  value: DelegationAmino;
}
export interface DelegationSDKType {
  delegation_address: string;
  validator_address: string;
  amount: CoinSDKType;
  height: Long;
  redelegation_end: Long;
}
export interface PortConnectionTuple {
  connectionId: string;
  portId: string;
}
export interface PortConnectionTupleProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.PortConnectionTuple";
  value: Uint8Array;
}
export interface PortConnectionTupleAmino {
  connection_id: string;
  port_id: string;
}
export interface PortConnectionTupleAminoMsg {
  type: "/quicksilver.interchainstaking.v1.PortConnectionTuple";
  value: PortConnectionTupleAmino;
}
export interface PortConnectionTupleSDKType {
  connection_id: string;
  port_id: string;
}
export interface Receipt {
  chainId: string;
  sender: string;
  txhash: string;
  amount: Coin[];
  firstSeen?: Date;
  completed?: Date;
}
export interface ReceiptProtoMsg {
  typeUrl: "/quicksilver.interchainstaking.v1.Receipt";
  value: Uint8Array;
}
export interface ReceiptAmino {
  chain_id: string;
  sender: string;
  txhash: string;
  amount: CoinAmino[];
  first_seen?: Date;
  completed?: Date;
}
export interface ReceiptAminoMsg {
  type: "/quicksilver.interchainstaking.v1.Receipt";
  value: ReceiptAmino;
}
export interface ReceiptSDKType {
  chain_id: string;
  sender: string;
  txhash: string;
  amount: CoinSDKType[];
  first_seen?: Date;
  completed?: Date;
}
function createBaseZone(): Zone {
  return {
    connectionId: "",
    chainId: "",
    depositAddress: ICAAccount.fromPartial({}),
    withdrawalAddress: ICAAccount.fromPartial({}),
    performanceAddress: ICAAccount.fromPartial({}),
    delegationAddress: ICAAccount.fromPartial({}),
    accountPrefix: "",
    localDenom: "",
    baseDenom: "",
    redemptionRate: "",
    lastRedemptionRate: "",
    validators: [],
    aggregateIntent: [],
    multiSend: false,
    liquidityModule: false,
    withdrawalWaitgroup: 0,
    ibcNextValidatorsHash: new Uint8Array(),
    validatorSelectionAllocation: Long.UZERO,
    holdingsAllocation: Long.UZERO,
    lastEpochHeight: Long.ZERO,
    tvl: "",
    unbondingPeriod: Long.ZERO,
    messagesPerTx: Long.ZERO
  };
}
export const Zone = {
  typeUrl: "/quicksilver.interchainstaking.v1.Zone",
  encode(message: Zone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.connectionId !== "") {
      writer.uint32(10).string(message.connectionId);
    }
    if (message.chainId !== "") {
      writer.uint32(18).string(message.chainId);
    }
    if (message.depositAddress !== undefined) {
      ICAAccount.encode(message.depositAddress, writer.uint32(26).fork()).ldelim();
    }
    if (message.withdrawalAddress !== undefined) {
      ICAAccount.encode(message.withdrawalAddress, writer.uint32(34).fork()).ldelim();
    }
    if (message.performanceAddress !== undefined) {
      ICAAccount.encode(message.performanceAddress, writer.uint32(42).fork()).ldelim();
    }
    if (message.delegationAddress !== undefined) {
      ICAAccount.encode(message.delegationAddress, writer.uint32(50).fork()).ldelim();
    }
    if (message.accountPrefix !== "") {
      writer.uint32(58).string(message.accountPrefix);
    }
    if (message.localDenom !== "") {
      writer.uint32(66).string(message.localDenom);
    }
    if (message.baseDenom !== "") {
      writer.uint32(74).string(message.baseDenom);
    }
    if (message.redemptionRate !== "") {
      writer.uint32(82).string(message.redemptionRate);
    }
    if (message.lastRedemptionRate !== "") {
      writer.uint32(90).string(message.lastRedemptionRate);
    }
    for (const v of message.validators) {
      Validator.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.aggregateIntent) {
      ValidatorIntent.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    if (message.multiSend === true) {
      writer.uint32(112).bool(message.multiSend);
    }
    if (message.liquidityModule === true) {
      writer.uint32(120).bool(message.liquidityModule);
    }
    if (message.withdrawalWaitgroup !== 0) {
      writer.uint32(128).uint32(message.withdrawalWaitgroup);
    }
    if (message.ibcNextValidatorsHash.length !== 0) {
      writer.uint32(138).bytes(message.ibcNextValidatorsHash);
    }
    if (!message.validatorSelectionAllocation.isZero()) {
      writer.uint32(144).uint64(message.validatorSelectionAllocation);
    }
    if (!message.holdingsAllocation.isZero()) {
      writer.uint32(152).uint64(message.holdingsAllocation);
    }
    if (!message.lastEpochHeight.isZero()) {
      writer.uint32(160).int64(message.lastEpochHeight);
    }
    if (message.tvl !== "") {
      writer.uint32(170).string(message.tvl);
    }
    if (!message.unbondingPeriod.isZero()) {
      writer.uint32(176).int64(message.unbondingPeriod);
    }
    if (!message.messagesPerTx.isZero()) {
      writer.uint32(184).int64(message.messagesPerTx);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Zone {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseZone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.connectionId = reader.string();
          break;
        case 2:
          message.chainId = reader.string();
          break;
        case 3:
          message.depositAddress = ICAAccount.decode(reader, reader.uint32());
          break;
        case 4:
          message.withdrawalAddress = ICAAccount.decode(reader, reader.uint32());
          break;
        case 5:
          message.performanceAddress = ICAAccount.decode(reader, reader.uint32());
          break;
        case 6:
          message.delegationAddress = ICAAccount.decode(reader, reader.uint32());
          break;
        case 7:
          message.accountPrefix = reader.string();
          break;
        case 8:
          message.localDenom = reader.string();
          break;
        case 9:
          message.baseDenom = reader.string();
          break;
        case 10:
          message.redemptionRate = reader.string();
          break;
        case 11:
          message.lastRedemptionRate = reader.string();
          break;
        case 12:
          message.validators.push(Validator.decode(reader, reader.uint32()));
          break;
        case 13:
          message.aggregateIntent.push(ValidatorIntent.decode(reader, reader.uint32()));
          break;
        case 14:
          message.multiSend = reader.bool();
          break;
        case 15:
          message.liquidityModule = reader.bool();
          break;
        case 16:
          message.withdrawalWaitgroup = reader.uint32();
          break;
        case 17:
          message.ibcNextValidatorsHash = reader.bytes();
          break;
        case 18:
          message.validatorSelectionAllocation = (reader.uint64() as Long);
          break;
        case 19:
          message.holdingsAllocation = (reader.uint64() as Long);
          break;
        case 20:
          message.lastEpochHeight = (reader.int64() as Long);
          break;
        case 21:
          message.tvl = reader.string();
          break;
        case 22:
          message.unbondingPeriod = (reader.int64() as Long);
          break;
        case 23:
          message.messagesPerTx = (reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Zone {
    const obj = createBaseZone();
    if (isSet(object.connectionId)) obj.connectionId = String(object.connectionId);
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.depositAddress)) obj.depositAddress = ICAAccount.fromJSON(object.depositAddress);
    if (isSet(object.withdrawalAddress)) obj.withdrawalAddress = ICAAccount.fromJSON(object.withdrawalAddress);
    if (isSet(object.performanceAddress)) obj.performanceAddress = ICAAccount.fromJSON(object.performanceAddress);
    if (isSet(object.delegationAddress)) obj.delegationAddress = ICAAccount.fromJSON(object.delegationAddress);
    if (isSet(object.accountPrefix)) obj.accountPrefix = String(object.accountPrefix);
    if (isSet(object.localDenom)) obj.localDenom = String(object.localDenom);
    if (isSet(object.baseDenom)) obj.baseDenom = String(object.baseDenom);
    if (isSet(object.redemptionRate)) obj.redemptionRate = String(object.redemptionRate);
    if (isSet(object.lastRedemptionRate)) obj.lastRedemptionRate = String(object.lastRedemptionRate);
    if (Array.isArray(object?.validators)) obj.validators = object.validators.map((e: any) => Validator.fromJSON(e));
    if (Array.isArray(object?.aggregateIntent)) obj.aggregateIntent = object.aggregateIntent.map((e: any) => ValidatorIntent.fromJSON(e));
    if (isSet(object.multiSend)) obj.multiSend = Boolean(object.multiSend);
    if (isSet(object.liquidityModule)) obj.liquidityModule = Boolean(object.liquidityModule);
    if (isSet(object.withdrawalWaitgroup)) obj.withdrawalWaitgroup = Number(object.withdrawalWaitgroup);
    if (isSet(object.ibcNextValidatorsHash)) obj.ibcNextValidatorsHash = bytesFromBase64(object.ibcNextValidatorsHash);
    if (isSet(object.validatorSelectionAllocation)) obj.validatorSelectionAllocation = Long.fromValue(object.validatorSelectionAllocation);
    if (isSet(object.holdingsAllocation)) obj.holdingsAllocation = Long.fromValue(object.holdingsAllocation);
    if (isSet(object.lastEpochHeight)) obj.lastEpochHeight = Long.fromValue(object.lastEpochHeight);
    if (isSet(object.tvl)) obj.tvl = String(object.tvl);
    if (isSet(object.unbondingPeriod)) obj.unbondingPeriod = Long.fromValue(object.unbondingPeriod);
    if (isSet(object.messagesPerTx)) obj.messagesPerTx = Long.fromValue(object.messagesPerTx);
    return obj;
  },
  toJSON(message: Zone): unknown {
    const obj: any = {};
    message.connectionId !== undefined && (obj.connectionId = message.connectionId);
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.depositAddress !== undefined && (obj.depositAddress = message.depositAddress ? ICAAccount.toJSON(message.depositAddress) : undefined);
    message.withdrawalAddress !== undefined && (obj.withdrawalAddress = message.withdrawalAddress ? ICAAccount.toJSON(message.withdrawalAddress) : undefined);
    message.performanceAddress !== undefined && (obj.performanceAddress = message.performanceAddress ? ICAAccount.toJSON(message.performanceAddress) : undefined);
    message.delegationAddress !== undefined && (obj.delegationAddress = message.delegationAddress ? ICAAccount.toJSON(message.delegationAddress) : undefined);
    message.accountPrefix !== undefined && (obj.accountPrefix = message.accountPrefix);
    message.localDenom !== undefined && (obj.localDenom = message.localDenom);
    message.baseDenom !== undefined && (obj.baseDenom = message.baseDenom);
    message.redemptionRate !== undefined && (obj.redemptionRate = message.redemptionRate);
    message.lastRedemptionRate !== undefined && (obj.lastRedemptionRate = message.lastRedemptionRate);
    if (message.validators) {
      obj.validators = message.validators.map(e => e ? Validator.toJSON(e) : undefined);
    } else {
      obj.validators = [];
    }
    if (message.aggregateIntent) {
      obj.aggregateIntent = message.aggregateIntent.map(e => e ? ValidatorIntent.toJSON(e) : undefined);
    } else {
      obj.aggregateIntent = [];
    }
    message.multiSend !== undefined && (obj.multiSend = message.multiSend);
    message.liquidityModule !== undefined && (obj.liquidityModule = message.liquidityModule);
    message.withdrawalWaitgroup !== undefined && (obj.withdrawalWaitgroup = Math.round(message.withdrawalWaitgroup));
    message.ibcNextValidatorsHash !== undefined && (obj.ibcNextValidatorsHash = base64FromBytes(message.ibcNextValidatorsHash !== undefined ? message.ibcNextValidatorsHash : new Uint8Array()));
    message.validatorSelectionAllocation !== undefined && (obj.validatorSelectionAllocation = (message.validatorSelectionAllocation || Long.UZERO).toString());
    message.holdingsAllocation !== undefined && (obj.holdingsAllocation = (message.holdingsAllocation || Long.UZERO).toString());
    message.lastEpochHeight !== undefined && (obj.lastEpochHeight = (message.lastEpochHeight || Long.ZERO).toString());
    message.tvl !== undefined && (obj.tvl = message.tvl);
    message.unbondingPeriod !== undefined && (obj.unbondingPeriod = (message.unbondingPeriod || Long.ZERO).toString());
    message.messagesPerTx !== undefined && (obj.messagesPerTx = (message.messagesPerTx || Long.ZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<Zone>): Zone {
    const message = createBaseZone();
    message.connectionId = object.connectionId ?? "";
    message.chainId = object.chainId ?? "";
    if (object.depositAddress !== undefined && object.depositAddress !== null) {
      message.depositAddress = ICAAccount.fromPartial(object.depositAddress);
    }
    if (object.withdrawalAddress !== undefined && object.withdrawalAddress !== null) {
      message.withdrawalAddress = ICAAccount.fromPartial(object.withdrawalAddress);
    }
    if (object.performanceAddress !== undefined && object.performanceAddress !== null) {
      message.performanceAddress = ICAAccount.fromPartial(object.performanceAddress);
    }
    if (object.delegationAddress !== undefined && object.delegationAddress !== null) {
      message.delegationAddress = ICAAccount.fromPartial(object.delegationAddress);
    }
    message.accountPrefix = object.accountPrefix ?? "";
    message.localDenom = object.localDenom ?? "";
    message.baseDenom = object.baseDenom ?? "";
    message.redemptionRate = object.redemptionRate ?? "";
    message.lastRedemptionRate = object.lastRedemptionRate ?? "";
    message.validators = object.validators?.map(e => Validator.fromPartial(e)) || [];
    message.aggregateIntent = object.aggregateIntent?.map(e => ValidatorIntent.fromPartial(e)) || [];
    message.multiSend = object.multiSend ?? false;
    message.liquidityModule = object.liquidityModule ?? false;
    message.withdrawalWaitgroup = object.withdrawalWaitgroup ?? 0;
    message.ibcNextValidatorsHash = object.ibcNextValidatorsHash ?? new Uint8Array();
    if (object.validatorSelectionAllocation !== undefined && object.validatorSelectionAllocation !== null) {
      message.validatorSelectionAllocation = Long.fromValue(object.validatorSelectionAllocation);
    }
    if (object.holdingsAllocation !== undefined && object.holdingsAllocation !== null) {
      message.holdingsAllocation = Long.fromValue(object.holdingsAllocation);
    }
    if (object.lastEpochHeight !== undefined && object.lastEpochHeight !== null) {
      message.lastEpochHeight = Long.fromValue(object.lastEpochHeight);
    }
    message.tvl = object.tvl ?? "";
    if (object.unbondingPeriod !== undefined && object.unbondingPeriod !== null) {
      message.unbondingPeriod = Long.fromValue(object.unbondingPeriod);
    }
    if (object.messagesPerTx !== undefined && object.messagesPerTx !== null) {
      message.messagesPerTx = Long.fromValue(object.messagesPerTx);
    }
    return message;
  },
  fromSDK(object: ZoneSDKType): Zone {
    return {
      connectionId: object?.connection_id,
      chainId: object?.chain_id,
      depositAddress: object.deposit_address ? ICAAccount.fromSDK(object.deposit_address) : undefined,
      withdrawalAddress: object.withdrawal_address ? ICAAccount.fromSDK(object.withdrawal_address) : undefined,
      performanceAddress: object.performance_address ? ICAAccount.fromSDK(object.performance_address) : undefined,
      delegationAddress: object.delegation_address ? ICAAccount.fromSDK(object.delegation_address) : undefined,
      accountPrefix: object?.account_prefix,
      localDenom: object?.local_denom,
      baseDenom: object?.base_denom,
      redemptionRate: object?.redemption_rate,
      lastRedemptionRate: object?.last_redemption_rate,
      validators: Array.isArray(object?.validators) ? object.validators.map((e: any) => Validator.fromSDK(e)) : [],
      aggregateIntent: Array.isArray(object?.aggregate_intent) ? object.aggregate_intent.map((e: any) => ValidatorIntent.fromSDK(e)) : [],
      multiSend: object?.multi_send,
      liquidityModule: object?.liquidity_module,
      withdrawalWaitgroup: object?.withdrawal_waitgroup,
      ibcNextValidatorsHash: object?.ibc_next_validators_hash,
      validatorSelectionAllocation: object?.validator_selection_allocation,
      holdingsAllocation: object?.holdings_allocation,
      lastEpochHeight: object?.last_epoch_height,
      tvl: object?.tvl,
      unbondingPeriod: object?.unbonding_period,
      messagesPerTx: object?.messages_per_tx
    };
  },
  toSDK(message: Zone): ZoneSDKType {
    const obj: any = {};
    obj.connection_id = message.connectionId;
    obj.chain_id = message.chainId;
    message.depositAddress !== undefined && (obj.deposit_address = message.depositAddress ? ICAAccount.toSDK(message.depositAddress) : undefined);
    message.withdrawalAddress !== undefined && (obj.withdrawal_address = message.withdrawalAddress ? ICAAccount.toSDK(message.withdrawalAddress) : undefined);
    message.performanceAddress !== undefined && (obj.performance_address = message.performanceAddress ? ICAAccount.toSDK(message.performanceAddress) : undefined);
    message.delegationAddress !== undefined && (obj.delegation_address = message.delegationAddress ? ICAAccount.toSDK(message.delegationAddress) : undefined);
    obj.account_prefix = message.accountPrefix;
    obj.local_denom = message.localDenom;
    obj.base_denom = message.baseDenom;
    obj.redemption_rate = message.redemptionRate;
    obj.last_redemption_rate = message.lastRedemptionRate;
    if (message.validators) {
      obj.validators = message.validators.map(e => e ? Validator.toSDK(e) : undefined);
    } else {
      obj.validators = [];
    }
    if (message.aggregateIntent) {
      obj.aggregate_intent = message.aggregateIntent.map(e => e ? ValidatorIntent.toSDK(e) : undefined);
    } else {
      obj.aggregate_intent = [];
    }
    obj.multi_send = message.multiSend;
    obj.liquidity_module = message.liquidityModule;
    obj.withdrawal_waitgroup = message.withdrawalWaitgroup;
    obj.ibc_next_validators_hash = message.ibcNextValidatorsHash;
    obj.validator_selection_allocation = message.validatorSelectionAllocation;
    obj.holdings_allocation = message.holdingsAllocation;
    obj.last_epoch_height = message.lastEpochHeight;
    obj.tvl = message.tvl;
    obj.unbonding_period = message.unbondingPeriod;
    obj.messages_per_tx = message.messagesPerTx;
    return obj;
  },
  fromAmino(object: ZoneAmino): Zone {
    return {
      connectionId: object.connection_id,
      chainId: object.chain_id,
      depositAddress: object?.deposit_address ? ICAAccount.fromAmino(object.deposit_address) : undefined,
      withdrawalAddress: object?.withdrawal_address ? ICAAccount.fromAmino(object.withdrawal_address) : undefined,
      performanceAddress: object?.performance_address ? ICAAccount.fromAmino(object.performance_address) : undefined,
      delegationAddress: object?.delegation_address ? ICAAccount.fromAmino(object.delegation_address) : undefined,
      accountPrefix: object.account_prefix,
      localDenom: object.local_denom,
      baseDenom: object.base_denom,
      redemptionRate: object.redemption_rate,
      lastRedemptionRate: object.last_redemption_rate,
      validators: Array.isArray(object?.validators) ? object.validators.map((e: any) => Validator.fromAmino(e)) : [],
      aggregateIntent: Array.isArray(object?.aggregate_intent) ? object.aggregate_intent.map((e: any) => ValidatorIntent.fromAmino(e)) : [],
      multiSend: object.multi_send,
      liquidityModule: object.liquidity_module,
      withdrawalWaitgroup: object.withdrawal_waitgroup,
      ibcNextValidatorsHash: object.ibc_next_validators_hash,
      validatorSelectionAllocation: Long.fromString(object.validator_selection_allocation),
      holdingsAllocation: Long.fromString(object.holdings_allocation),
      lastEpochHeight: Long.fromString(object.last_epoch_height),
      tvl: object.tvl,
      unbondingPeriod: Long.fromString(object.unbonding_period),
      messagesPerTx: Long.fromString(object.messages_per_tx)
    };
  },
  toAmino(message: Zone): ZoneAmino {
    const obj: any = {};
    obj.connection_id = message.connectionId;
    obj.chain_id = message.chainId;
    obj.deposit_address = message.depositAddress ? ICAAccount.toAmino(message.depositAddress) : undefined;
    obj.withdrawal_address = message.withdrawalAddress ? ICAAccount.toAmino(message.withdrawalAddress) : undefined;
    obj.performance_address = message.performanceAddress ? ICAAccount.toAmino(message.performanceAddress) : undefined;
    obj.delegation_address = message.delegationAddress ? ICAAccount.toAmino(message.delegationAddress) : undefined;
    obj.account_prefix = message.accountPrefix;
    obj.local_denom = message.localDenom;
    obj.base_denom = message.baseDenom;
    obj.redemption_rate = message.redemptionRate;
    obj.last_redemption_rate = message.lastRedemptionRate;
    if (message.validators) {
      obj.validators = message.validators.map(e => e ? Validator.toAmino(e) : undefined);
    } else {
      obj.validators = [];
    }
    if (message.aggregateIntent) {
      obj.aggregate_intent = message.aggregateIntent.map(e => e ? ValidatorIntent.toAmino(e) : undefined);
    } else {
      obj.aggregate_intent = [];
    }
    obj.multi_send = message.multiSend;
    obj.liquidity_module = message.liquidityModule;
    obj.withdrawal_waitgroup = message.withdrawalWaitgroup;
    obj.ibc_next_validators_hash = message.ibcNextValidatorsHash;
    obj.validator_selection_allocation = message.validatorSelectionAllocation ? message.validatorSelectionAllocation.toString() : undefined;
    obj.holdings_allocation = message.holdingsAllocation ? message.holdingsAllocation.toString() : undefined;
    obj.last_epoch_height = message.lastEpochHeight ? message.lastEpochHeight.toString() : undefined;
    obj.tvl = message.tvl;
    obj.unbonding_period = message.unbondingPeriod ? message.unbondingPeriod.toString() : undefined;
    obj.messages_per_tx = message.messagesPerTx ? message.messagesPerTx.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: ZoneAminoMsg): Zone {
    return Zone.fromAmino(object.value);
  },
  fromProtoMsg(message: ZoneProtoMsg): Zone {
    return Zone.decode(message.value);
  },
  toProto(message: Zone): Uint8Array {
    return Zone.encode(message).finish();
  },
  toProtoMsg(message: Zone): ZoneProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.Zone",
      value: Zone.encode(message).finish()
    };
  }
};
function createBaseICAAccount(): ICAAccount {
  return {
    address: "",
    balance: [],
    portName: "",
    withdrawalAddress: "",
    balanceWaitgroup: 0
  };
}
export const ICAAccount = {
  typeUrl: "/quicksilver.interchainstaking.v1.ICAAccount",
  encode(message: ICAAccount, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    for (const v of message.balance) {
      Coin.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.portName !== "") {
      writer.uint32(26).string(message.portName);
    }
    if (message.withdrawalAddress !== "") {
      writer.uint32(34).string(message.withdrawalAddress);
    }
    if (message.balanceWaitgroup !== 0) {
      writer.uint32(40).uint32(message.balanceWaitgroup);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): ICAAccount {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseICAAccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.balance.push(Coin.decode(reader, reader.uint32()));
          break;
        case 3:
          message.portName = reader.string();
          break;
        case 4:
          message.withdrawalAddress = reader.string();
          break;
        case 5:
          message.balanceWaitgroup = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): ICAAccount {
    const obj = createBaseICAAccount();
    if (isSet(object.address)) obj.address = String(object.address);
    if (Array.isArray(object?.balance)) obj.balance = object.balance.map((e: any) => Coin.fromJSON(e));
    if (isSet(object.portName)) obj.portName = String(object.portName);
    if (isSet(object.withdrawalAddress)) obj.withdrawalAddress = String(object.withdrawalAddress);
    if (isSet(object.balanceWaitgroup)) obj.balanceWaitgroup = Number(object.balanceWaitgroup);
    return obj;
  },
  toJSON(message: ICAAccount): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    if (message.balance) {
      obj.balance = message.balance.map(e => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.balance = [];
    }
    message.portName !== undefined && (obj.portName = message.portName);
    message.withdrawalAddress !== undefined && (obj.withdrawalAddress = message.withdrawalAddress);
    message.balanceWaitgroup !== undefined && (obj.balanceWaitgroup = Math.round(message.balanceWaitgroup));
    return obj;
  },
  fromPartial(object: DeepPartial<ICAAccount>): ICAAccount {
    const message = createBaseICAAccount();
    message.address = object.address ?? "";
    message.balance = object.balance?.map(e => Coin.fromPartial(e)) || [];
    message.portName = object.portName ?? "";
    message.withdrawalAddress = object.withdrawalAddress ?? "";
    message.balanceWaitgroup = object.balanceWaitgroup ?? 0;
    return message;
  },
  fromSDK(object: ICAAccountSDKType): ICAAccount {
    return {
      address: object?.address,
      balance: Array.isArray(object?.balance) ? object.balance.map((e: any) => Coin.fromSDK(e)) : [],
      portName: object?.port_name,
      withdrawalAddress: object?.withdrawal_address,
      balanceWaitgroup: object?.balance_waitgroup
    };
  },
  toSDK(message: ICAAccount): ICAAccountSDKType {
    const obj: any = {};
    obj.address = message.address;
    if (message.balance) {
      obj.balance = message.balance.map(e => e ? Coin.toSDK(e) : undefined);
    } else {
      obj.balance = [];
    }
    obj.port_name = message.portName;
    obj.withdrawal_address = message.withdrawalAddress;
    obj.balance_waitgroup = message.balanceWaitgroup;
    return obj;
  },
  fromAmino(object: ICAAccountAmino): ICAAccount {
    return {
      address: object.address,
      balance: Array.isArray(object?.balance) ? object.balance.map((e: any) => Coin.fromAmino(e)) : [],
      portName: object.port_name,
      withdrawalAddress: object.withdrawal_address,
      balanceWaitgroup: object.balance_waitgroup
    };
  },
  toAmino(message: ICAAccount): ICAAccountAmino {
    const obj: any = {};
    obj.address = message.address;
    if (message.balance) {
      obj.balance = message.balance.map(e => e ? Coin.toAmino(e) : undefined);
    } else {
      obj.balance = [];
    }
    obj.port_name = message.portName;
    obj.withdrawal_address = message.withdrawalAddress;
    obj.balance_waitgroup = message.balanceWaitgroup;
    return obj;
  },
  fromAminoMsg(object: ICAAccountAminoMsg): ICAAccount {
    return ICAAccount.fromAmino(object.value);
  },
  fromProtoMsg(message: ICAAccountProtoMsg): ICAAccount {
    return ICAAccount.decode(message.value);
  },
  toProto(message: ICAAccount): Uint8Array {
    return ICAAccount.encode(message).finish();
  },
  toProtoMsg(message: ICAAccount): ICAAccountProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.ICAAccount",
      value: ICAAccount.encode(message).finish()
    };
  }
};
function createBaseDistribution(): Distribution {
  return {
    valoper: "",
    amount: Long.UZERO
  };
}
export const Distribution = {
  typeUrl: "/quicksilver.interchainstaking.v1.Distribution",
  encode(message: Distribution, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.valoper !== "") {
      writer.uint32(10).string(message.valoper);
    }
    if (!message.amount.isZero()) {
      writer.uint32(16).uint64(message.amount);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Distribution {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDistribution();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.valoper = reader.string();
          break;
        case 2:
          message.amount = (reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Distribution {
    const obj = createBaseDistribution();
    if (isSet(object.valoper)) obj.valoper = String(object.valoper);
    if (isSet(object.amount)) obj.amount = Long.fromValue(object.amount);
    return obj;
  },
  toJSON(message: Distribution): unknown {
    const obj: any = {};
    message.valoper !== undefined && (obj.valoper = message.valoper);
    message.amount !== undefined && (obj.amount = (message.amount || Long.UZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<Distribution>): Distribution {
    const message = createBaseDistribution();
    message.valoper = object.valoper ?? "";
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Long.fromValue(object.amount);
    }
    return message;
  },
  fromSDK(object: DistributionSDKType): Distribution {
    return {
      valoper: object?.valoper,
      amount: object?.amount
    };
  },
  toSDK(message: Distribution): DistributionSDKType {
    const obj: any = {};
    obj.valoper = message.valoper;
    obj.amount = message.amount;
    return obj;
  },
  fromAmino(object: DistributionAmino): Distribution {
    return {
      valoper: object.valoper,
      amount: Long.fromString(object.amount)
    };
  },
  toAmino(message: Distribution): DistributionAmino {
    const obj: any = {};
    obj.valoper = message.valoper;
    obj.amount = message.amount ? message.amount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: DistributionAminoMsg): Distribution {
    return Distribution.fromAmino(object.value);
  },
  fromProtoMsg(message: DistributionProtoMsg): Distribution {
    return Distribution.decode(message.value);
  },
  toProto(message: Distribution): Uint8Array {
    return Distribution.encode(message).finish();
  },
  toProtoMsg(message: Distribution): DistributionProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.Distribution",
      value: Distribution.encode(message).finish()
    };
  }
};
function createBaseWithdrawalRecord(): WithdrawalRecord {
  return {
    chainId: "",
    delegator: "",
    distribution: [],
    recipient: "",
    amount: [],
    burnAmount: Coin.fromPartial({}),
    txhash: "",
    status: 0,
    completionTime: new Date()
  };
}
export const WithdrawalRecord = {
  typeUrl: "/quicksilver.interchainstaking.v1.WithdrawalRecord",
  encode(message: WithdrawalRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (message.delegator !== "") {
      writer.uint32(18).string(message.delegator);
    }
    for (const v of message.distribution) {
      Distribution.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.recipient !== "") {
      writer.uint32(34).string(message.recipient);
    }
    for (const v of message.amount) {
      Coin.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.burnAmount !== undefined) {
      Coin.encode(message.burnAmount, writer.uint32(50).fork()).ldelim();
    }
    if (message.txhash !== "") {
      writer.uint32(58).string(message.txhash);
    }
    if (message.status !== 0) {
      writer.uint32(64).int32(message.status);
    }
    if (message.completionTime !== undefined) {
      Timestamp.encode(toTimestamp(message.completionTime), writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): WithdrawalRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWithdrawalRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.delegator = reader.string();
          break;
        case 3:
          message.distribution.push(Distribution.decode(reader, reader.uint32()));
          break;
        case 4:
          message.recipient = reader.string();
          break;
        case 5:
          message.amount.push(Coin.decode(reader, reader.uint32()));
          break;
        case 6:
          message.burnAmount = Coin.decode(reader, reader.uint32());
          break;
        case 7:
          message.txhash = reader.string();
          break;
        case 8:
          message.status = reader.int32();
          break;
        case 9:
          message.completionTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): WithdrawalRecord {
    const obj = createBaseWithdrawalRecord();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.delegator)) obj.delegator = String(object.delegator);
    if (Array.isArray(object?.distribution)) obj.distribution = object.distribution.map((e: any) => Distribution.fromJSON(e));
    if (isSet(object.recipient)) obj.recipient = String(object.recipient);
    if (Array.isArray(object?.amount)) obj.amount = object.amount.map((e: any) => Coin.fromJSON(e));
    if (isSet(object.burnAmount)) obj.burnAmount = Coin.fromJSON(object.burnAmount);
    if (isSet(object.txhash)) obj.txhash = String(object.txhash);
    if (isSet(object.status)) obj.status = Number(object.status);
    if (isSet(object.completionTime)) obj.completionTime = new Date(object.completionTime);
    return obj;
  },
  toJSON(message: WithdrawalRecord): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.delegator !== undefined && (obj.delegator = message.delegator);
    if (message.distribution) {
      obj.distribution = message.distribution.map(e => e ? Distribution.toJSON(e) : undefined);
    } else {
      obj.distribution = [];
    }
    message.recipient !== undefined && (obj.recipient = message.recipient);
    if (message.amount) {
      obj.amount = message.amount.map(e => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.amount = [];
    }
    message.burnAmount !== undefined && (obj.burnAmount = message.burnAmount ? Coin.toJSON(message.burnAmount) : undefined);
    message.txhash !== undefined && (obj.txhash = message.txhash);
    message.status !== undefined && (obj.status = Math.round(message.status));
    message.completionTime !== undefined && (obj.completionTime = message.completionTime.toISOString());
    return obj;
  },
  fromPartial(object: DeepPartial<WithdrawalRecord>): WithdrawalRecord {
    const message = createBaseWithdrawalRecord();
    message.chainId = object.chainId ?? "";
    message.delegator = object.delegator ?? "";
    message.distribution = object.distribution?.map(e => Distribution.fromPartial(e)) || [];
    message.recipient = object.recipient ?? "";
    message.amount = object.amount?.map(e => Coin.fromPartial(e)) || [];
    if (object.burnAmount !== undefined && object.burnAmount !== null) {
      message.burnAmount = Coin.fromPartial(object.burnAmount);
    }
    message.txhash = object.txhash ?? "";
    message.status = object.status ?? 0;
    message.completionTime = object.completionTime ?? undefined;
    return message;
  },
  fromSDK(object: WithdrawalRecordSDKType): WithdrawalRecord {
    return {
      chainId: object?.chain_id,
      delegator: object?.delegator,
      distribution: Array.isArray(object?.distribution) ? object.distribution.map((e: any) => Distribution.fromSDK(e)) : [],
      recipient: object?.recipient,
      amount: Array.isArray(object?.amount) ? object.amount.map((e: any) => Coin.fromSDK(e)) : [],
      burnAmount: object.burn_amount ? Coin.fromSDK(object.burn_amount) : undefined,
      txhash: object?.txhash,
      status: object?.status,
      completionTime: object.completion_time ?? undefined
    };
  },
  toSDK(message: WithdrawalRecord): WithdrawalRecordSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.delegator = message.delegator;
    if (message.distribution) {
      obj.distribution = message.distribution.map(e => e ? Distribution.toSDK(e) : undefined);
    } else {
      obj.distribution = [];
    }
    obj.recipient = message.recipient;
    if (message.amount) {
      obj.amount = message.amount.map(e => e ? Coin.toSDK(e) : undefined);
    } else {
      obj.amount = [];
    }
    message.burnAmount !== undefined && (obj.burn_amount = message.burnAmount ? Coin.toSDK(message.burnAmount) : undefined);
    obj.txhash = message.txhash;
    obj.status = message.status;
    message.completionTime !== undefined && (obj.completion_time = message.completionTime ?? undefined);
    return obj;
  },
  fromAmino(object: WithdrawalRecordAmino): WithdrawalRecord {
    return {
      chainId: object.chain_id,
      delegator: object.delegator,
      distribution: Array.isArray(object?.distribution) ? object.distribution.map((e: any) => Distribution.fromAmino(e)) : [],
      recipient: object.recipient,
      amount: Array.isArray(object?.amount) ? object.amount.map((e: any) => Coin.fromAmino(e)) : [],
      burnAmount: object?.burn_amount ? Coin.fromAmino(object.burn_amount) : undefined,
      txhash: object.txhash,
      status: object.status,
      completionTime: object.completion_time
    };
  },
  toAmino(message: WithdrawalRecord): WithdrawalRecordAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.delegator = message.delegator;
    if (message.distribution) {
      obj.distribution = message.distribution.map(e => e ? Distribution.toAmino(e) : undefined);
    } else {
      obj.distribution = [];
    }
    obj.recipient = message.recipient;
    if (message.amount) {
      obj.amount = message.amount.map(e => e ? Coin.toAmino(e) : undefined);
    } else {
      obj.amount = [];
    }
    obj.burn_amount = message.burnAmount ? Coin.toAmino(message.burnAmount) : undefined;
    obj.txhash = message.txhash;
    obj.status = message.status;
    obj.completion_time = message.completionTime;
    return obj;
  },
  fromAminoMsg(object: WithdrawalRecordAminoMsg): WithdrawalRecord {
    return WithdrawalRecord.fromAmino(object.value);
  },
  fromProtoMsg(message: WithdrawalRecordProtoMsg): WithdrawalRecord {
    return WithdrawalRecord.decode(message.value);
  },
  toProto(message: WithdrawalRecord): Uint8Array {
    return WithdrawalRecord.encode(message).finish();
  },
  toProtoMsg(message: WithdrawalRecord): WithdrawalRecordProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.WithdrawalRecord",
      value: WithdrawalRecord.encode(message).finish()
    };
  }
};
function createBaseUnbondingRecord(): UnbondingRecord {
  return {
    chainId: "",
    epochNumber: Long.ZERO,
    validator: "",
    relatedTxhash: []
  };
}
export const UnbondingRecord = {
  typeUrl: "/quicksilver.interchainstaking.v1.UnbondingRecord",
  encode(message: UnbondingRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (!message.epochNumber.isZero()) {
      writer.uint32(16).int64(message.epochNumber);
    }
    if (message.validator !== "") {
      writer.uint32(26).string(message.validator);
    }
    for (const v of message.relatedTxhash) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): UnbondingRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnbondingRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.epochNumber = (reader.int64() as Long);
          break;
        case 3:
          message.validator = reader.string();
          break;
        case 4:
          message.relatedTxhash.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): UnbondingRecord {
    const obj = createBaseUnbondingRecord();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.epochNumber)) obj.epochNumber = Long.fromValue(object.epochNumber);
    if (isSet(object.validator)) obj.validator = String(object.validator);
    if (Array.isArray(object?.relatedTxhash)) obj.relatedTxhash = object.relatedTxhash.map((e: any) => String(e));
    return obj;
  },
  toJSON(message: UnbondingRecord): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.epochNumber !== undefined && (obj.epochNumber = (message.epochNumber || Long.ZERO).toString());
    message.validator !== undefined && (obj.validator = message.validator);
    if (message.relatedTxhash) {
      obj.relatedTxhash = message.relatedTxhash.map(e => e);
    } else {
      obj.relatedTxhash = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<UnbondingRecord>): UnbondingRecord {
    const message = createBaseUnbondingRecord();
    message.chainId = object.chainId ?? "";
    if (object.epochNumber !== undefined && object.epochNumber !== null) {
      message.epochNumber = Long.fromValue(object.epochNumber);
    }
    message.validator = object.validator ?? "";
    message.relatedTxhash = object.relatedTxhash?.map(e => e) || [];
    return message;
  },
  fromSDK(object: UnbondingRecordSDKType): UnbondingRecord {
    return {
      chainId: object?.chain_id,
      epochNumber: object?.epoch_number,
      validator: object?.validator,
      relatedTxhash: Array.isArray(object?.related_txhash) ? object.related_txhash.map((e: any) => e) : []
    };
  },
  toSDK(message: UnbondingRecord): UnbondingRecordSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.epoch_number = message.epochNumber;
    obj.validator = message.validator;
    if (message.relatedTxhash) {
      obj.related_txhash = message.relatedTxhash.map(e => e);
    } else {
      obj.related_txhash = [];
    }
    return obj;
  },
  fromAmino(object: UnbondingRecordAmino): UnbondingRecord {
    return {
      chainId: object.chain_id,
      epochNumber: Long.fromString(object.epoch_number),
      validator: object.validator,
      relatedTxhash: Array.isArray(object?.related_txhash) ? object.related_txhash.map((e: any) => e) : []
    };
  },
  toAmino(message: UnbondingRecord): UnbondingRecordAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.epoch_number = message.epochNumber ? message.epochNumber.toString() : undefined;
    obj.validator = message.validator;
    if (message.relatedTxhash) {
      obj.related_txhash = message.relatedTxhash.map(e => e);
    } else {
      obj.related_txhash = [];
    }
    return obj;
  },
  fromAminoMsg(object: UnbondingRecordAminoMsg): UnbondingRecord {
    return UnbondingRecord.fromAmino(object.value);
  },
  fromProtoMsg(message: UnbondingRecordProtoMsg): UnbondingRecord {
    return UnbondingRecord.decode(message.value);
  },
  toProto(message: UnbondingRecord): Uint8Array {
    return UnbondingRecord.encode(message).finish();
  },
  toProtoMsg(message: UnbondingRecord): UnbondingRecordProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.UnbondingRecord",
      value: UnbondingRecord.encode(message).finish()
    };
  }
};
function createBaseRedelegationRecord(): RedelegationRecord {
  return {
    chainId: "",
    epochNumber: Long.ZERO,
    source: "",
    destination: "",
    amount: Long.ZERO,
    completionTime: new Date()
  };
}
export const RedelegationRecord = {
  typeUrl: "/quicksilver.interchainstaking.v1.RedelegationRecord",
  encode(message: RedelegationRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (!message.epochNumber.isZero()) {
      writer.uint32(16).int64(message.epochNumber);
    }
    if (message.source !== "") {
      writer.uint32(26).string(message.source);
    }
    if (message.destination !== "") {
      writer.uint32(34).string(message.destination);
    }
    if (!message.amount.isZero()) {
      writer.uint32(40).int64(message.amount);
    }
    if (message.completionTime !== undefined) {
      Timestamp.encode(toTimestamp(message.completionTime), writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): RedelegationRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRedelegationRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.epochNumber = (reader.int64() as Long);
          break;
        case 3:
          message.source = reader.string();
          break;
        case 4:
          message.destination = reader.string();
          break;
        case 5:
          message.amount = (reader.int64() as Long);
          break;
        case 6:
          message.completionTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): RedelegationRecord {
    const obj = createBaseRedelegationRecord();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.epochNumber)) obj.epochNumber = Long.fromValue(object.epochNumber);
    if (isSet(object.source)) obj.source = String(object.source);
    if (isSet(object.destination)) obj.destination = String(object.destination);
    if (isSet(object.amount)) obj.amount = Long.fromValue(object.amount);
    if (isSet(object.completionTime)) obj.completionTime = new Date(object.completionTime);
    return obj;
  },
  toJSON(message: RedelegationRecord): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.epochNumber !== undefined && (obj.epochNumber = (message.epochNumber || Long.ZERO).toString());
    message.source !== undefined && (obj.source = message.source);
    message.destination !== undefined && (obj.destination = message.destination);
    message.amount !== undefined && (obj.amount = (message.amount || Long.ZERO).toString());
    message.completionTime !== undefined && (obj.completionTime = message.completionTime.toISOString());
    return obj;
  },
  fromPartial(object: DeepPartial<RedelegationRecord>): RedelegationRecord {
    const message = createBaseRedelegationRecord();
    message.chainId = object.chainId ?? "";
    if (object.epochNumber !== undefined && object.epochNumber !== null) {
      message.epochNumber = Long.fromValue(object.epochNumber);
    }
    message.source = object.source ?? "";
    message.destination = object.destination ?? "";
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Long.fromValue(object.amount);
    }
    message.completionTime = object.completionTime ?? undefined;
    return message;
  },
  fromSDK(object: RedelegationRecordSDKType): RedelegationRecord {
    return {
      chainId: object?.chain_id,
      epochNumber: object?.epoch_number,
      source: object?.source,
      destination: object?.destination,
      amount: object?.amount,
      completionTime: object.completion_time ?? undefined
    };
  },
  toSDK(message: RedelegationRecord): RedelegationRecordSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.epoch_number = message.epochNumber;
    obj.source = message.source;
    obj.destination = message.destination;
    obj.amount = message.amount;
    message.completionTime !== undefined && (obj.completion_time = message.completionTime ?? undefined);
    return obj;
  },
  fromAmino(object: RedelegationRecordAmino): RedelegationRecord {
    return {
      chainId: object.chain_id,
      epochNumber: Long.fromString(object.epoch_number),
      source: object.source,
      destination: object.destination,
      amount: Long.fromString(object.amount),
      completionTime: object.completion_time
    };
  },
  toAmino(message: RedelegationRecord): RedelegationRecordAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.epoch_number = message.epochNumber ? message.epochNumber.toString() : undefined;
    obj.source = message.source;
    obj.destination = message.destination;
    obj.amount = message.amount ? message.amount.toString() : undefined;
    obj.completion_time = message.completionTime;
    return obj;
  },
  fromAminoMsg(object: RedelegationRecordAminoMsg): RedelegationRecord {
    return RedelegationRecord.fromAmino(object.value);
  },
  fromProtoMsg(message: RedelegationRecordProtoMsg): RedelegationRecord {
    return RedelegationRecord.decode(message.value);
  },
  toProto(message: RedelegationRecord): Uint8Array {
    return RedelegationRecord.encode(message).finish();
  },
  toProtoMsg(message: RedelegationRecord): RedelegationRecordProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.RedelegationRecord",
      value: RedelegationRecord.encode(message).finish()
    };
  }
};
function createBaseTransferRecord(): TransferRecord {
  return {
    sender: "",
    recipient: "",
    amount: Coin.fromPartial({})
  };
}
export const TransferRecord = {
  typeUrl: "/quicksilver.interchainstaking.v1.TransferRecord",
  encode(message: TransferRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sender !== "") {
      writer.uint32(10).string(message.sender);
    }
    if (message.recipient !== "") {
      writer.uint32(18).string(message.recipient);
    }
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): TransferRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTransferRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = reader.string();
          break;
        case 2:
          message.recipient = reader.string();
          break;
        case 3:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): TransferRecord {
    const obj = createBaseTransferRecord();
    if (isSet(object.sender)) obj.sender = String(object.sender);
    if (isSet(object.recipient)) obj.recipient = String(object.recipient);
    if (isSet(object.amount)) obj.amount = Coin.fromJSON(object.amount);
    return obj;
  },
  toJSON(message: TransferRecord): unknown {
    const obj: any = {};
    message.sender !== undefined && (obj.sender = message.sender);
    message.recipient !== undefined && (obj.recipient = message.recipient);
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    return obj;
  },
  fromPartial(object: DeepPartial<TransferRecord>): TransferRecord {
    const message = createBaseTransferRecord();
    message.sender = object.sender ?? "";
    message.recipient = object.recipient ?? "";
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Coin.fromPartial(object.amount);
    }
    return message;
  },
  fromSDK(object: TransferRecordSDKType): TransferRecord {
    return {
      sender: object?.sender,
      recipient: object?.recipient,
      amount: object.amount ? Coin.fromSDK(object.amount) : undefined
    };
  },
  toSDK(message: TransferRecord): TransferRecordSDKType {
    const obj: any = {};
    obj.sender = message.sender;
    obj.recipient = message.recipient;
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toSDK(message.amount) : undefined);
    return obj;
  },
  fromAmino(object: TransferRecordAmino): TransferRecord {
    return {
      sender: object.sender,
      recipient: object.recipient,
      amount: object?.amount ? Coin.fromAmino(object.amount) : undefined
    };
  },
  toAmino(message: TransferRecord): TransferRecordAmino {
    const obj: any = {};
    obj.sender = message.sender;
    obj.recipient = message.recipient;
    obj.amount = message.amount ? Coin.toAmino(message.amount) : undefined;
    return obj;
  },
  fromAminoMsg(object: TransferRecordAminoMsg): TransferRecord {
    return TransferRecord.fromAmino(object.value);
  },
  fromProtoMsg(message: TransferRecordProtoMsg): TransferRecord {
    return TransferRecord.decode(message.value);
  },
  toProto(message: TransferRecord): Uint8Array {
    return TransferRecord.encode(message).finish();
  },
  toProtoMsg(message: TransferRecord): TransferRecordProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.TransferRecord",
      value: TransferRecord.encode(message).finish()
    };
  }
};
function createBaseValidator(): Validator {
  return {
    valoperAddress: "",
    commissionRate: "",
    delegatorShares: "",
    votingPower: "",
    score: "",
    status: "",
    jailed: false,
    tombstoned: false,
    jailedSince: new Date()
  };
}
export const Validator = {
  typeUrl: "/quicksilver.interchainstaking.v1.Validator",
  encode(message: Validator, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.valoperAddress !== "") {
      writer.uint32(10).string(message.valoperAddress);
    }
    if (message.commissionRate !== "") {
      writer.uint32(18).string(message.commissionRate);
    }
    if (message.delegatorShares !== "") {
      writer.uint32(26).string(message.delegatorShares);
    }
    if (message.votingPower !== "") {
      writer.uint32(34).string(message.votingPower);
    }
    if (message.score !== "") {
      writer.uint32(42).string(message.score);
    }
    if (message.status !== "") {
      writer.uint32(50).string(message.status);
    }
    if (message.jailed === true) {
      writer.uint32(56).bool(message.jailed);
    }
    if (message.tombstoned === true) {
      writer.uint32(64).bool(message.tombstoned);
    }
    if (message.jailedSince !== undefined) {
      Timestamp.encode(toTimestamp(message.jailedSince), writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Validator {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidator();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.valoperAddress = reader.string();
          break;
        case 2:
          message.commissionRate = reader.string();
          break;
        case 3:
          message.delegatorShares = reader.string();
          break;
        case 4:
          message.votingPower = reader.string();
          break;
        case 5:
          message.score = reader.string();
          break;
        case 6:
          message.status = reader.string();
          break;
        case 7:
          message.jailed = reader.bool();
          break;
        case 8:
          message.tombstoned = reader.bool();
          break;
        case 9:
          message.jailedSince = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Validator {
    const obj = createBaseValidator();
    if (isSet(object.valoperAddress)) obj.valoperAddress = String(object.valoperAddress);
    if (isSet(object.commissionRate)) obj.commissionRate = String(object.commissionRate);
    if (isSet(object.delegatorShares)) obj.delegatorShares = String(object.delegatorShares);
    if (isSet(object.votingPower)) obj.votingPower = String(object.votingPower);
    if (isSet(object.score)) obj.score = String(object.score);
    if (isSet(object.status)) obj.status = String(object.status);
    if (isSet(object.jailed)) obj.jailed = Boolean(object.jailed);
    if (isSet(object.tombstoned)) obj.tombstoned = Boolean(object.tombstoned);
    if (isSet(object.jailedSince)) obj.jailedSince = new Date(object.jailedSince);
    return obj;
  },
  toJSON(message: Validator): unknown {
    const obj: any = {};
    message.valoperAddress !== undefined && (obj.valoperAddress = message.valoperAddress);
    message.commissionRate !== undefined && (obj.commissionRate = message.commissionRate);
    message.delegatorShares !== undefined && (obj.delegatorShares = message.delegatorShares);
    message.votingPower !== undefined && (obj.votingPower = message.votingPower);
    message.score !== undefined && (obj.score = message.score);
    message.status !== undefined && (obj.status = message.status);
    message.jailed !== undefined && (obj.jailed = message.jailed);
    message.tombstoned !== undefined && (obj.tombstoned = message.tombstoned);
    message.jailedSince !== undefined && (obj.jailedSince = message.jailedSince.toISOString());
    return obj;
  },
  fromPartial(object: DeepPartial<Validator>): Validator {
    const message = createBaseValidator();
    message.valoperAddress = object.valoperAddress ?? "";
    message.commissionRate = object.commissionRate ?? "";
    message.delegatorShares = object.delegatorShares ?? "";
    message.votingPower = object.votingPower ?? "";
    message.score = object.score ?? "";
    message.status = object.status ?? "";
    message.jailed = object.jailed ?? false;
    message.tombstoned = object.tombstoned ?? false;
    message.jailedSince = object.jailedSince ?? undefined;
    return message;
  },
  fromSDK(object: ValidatorSDKType): Validator {
    return {
      valoperAddress: object?.valoper_address,
      commissionRate: object?.commission_rate,
      delegatorShares: object?.delegator_shares,
      votingPower: object?.voting_power,
      score: object?.score,
      status: object?.status,
      jailed: object?.jailed,
      tombstoned: object?.tombstoned,
      jailedSince: object.jailed_since ?? undefined
    };
  },
  toSDK(message: Validator): ValidatorSDKType {
    const obj: any = {};
    obj.valoper_address = message.valoperAddress;
    obj.commission_rate = message.commissionRate;
    obj.delegator_shares = message.delegatorShares;
    obj.voting_power = message.votingPower;
    obj.score = message.score;
    obj.status = message.status;
    obj.jailed = message.jailed;
    obj.tombstoned = message.tombstoned;
    message.jailedSince !== undefined && (obj.jailed_since = message.jailedSince ?? undefined);
    return obj;
  },
  fromAmino(object: ValidatorAmino): Validator {
    return {
      valoperAddress: object.valoper_address,
      commissionRate: object.commission_rate,
      delegatorShares: object.delegator_shares,
      votingPower: object.voting_power,
      score: object.score,
      status: object.status,
      jailed: object.jailed,
      tombstoned: object.tombstoned,
      jailedSince: object.jailed_since
    };
  },
  toAmino(message: Validator): ValidatorAmino {
    const obj: any = {};
    obj.valoper_address = message.valoperAddress;
    obj.commission_rate = message.commissionRate;
    obj.delegator_shares = message.delegatorShares;
    obj.voting_power = message.votingPower;
    obj.score = message.score;
    obj.status = message.status;
    obj.jailed = message.jailed;
    obj.tombstoned = message.tombstoned;
    obj.jailed_since = message.jailedSince;
    return obj;
  },
  fromAminoMsg(object: ValidatorAminoMsg): Validator {
    return Validator.fromAmino(object.value);
  },
  fromProtoMsg(message: ValidatorProtoMsg): Validator {
    return Validator.decode(message.value);
  },
  toProto(message: Validator): Uint8Array {
    return Validator.encode(message).finish();
  },
  toProtoMsg(message: Validator): ValidatorProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.Validator",
      value: Validator.encode(message).finish()
    };
  }
};
function createBaseDelegatorIntent(): DelegatorIntent {
  return {
    delegator: "",
    intents: []
  };
}
export const DelegatorIntent = {
  typeUrl: "/quicksilver.interchainstaking.v1.DelegatorIntent",
  encode(message: DelegatorIntent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.delegator !== "") {
      writer.uint32(10).string(message.delegator);
    }
    for (const v of message.intents) {
      ValidatorIntent.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): DelegatorIntent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDelegatorIntent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegator = reader.string();
          break;
        case 2:
          message.intents.push(ValidatorIntent.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): DelegatorIntent {
    const obj = createBaseDelegatorIntent();
    if (isSet(object.delegator)) obj.delegator = String(object.delegator);
    if (Array.isArray(object?.intents)) obj.intents = object.intents.map((e: any) => ValidatorIntent.fromJSON(e));
    return obj;
  },
  toJSON(message: DelegatorIntent): unknown {
    const obj: any = {};
    message.delegator !== undefined && (obj.delegator = message.delegator);
    if (message.intents) {
      obj.intents = message.intents.map(e => e ? ValidatorIntent.toJSON(e) : undefined);
    } else {
      obj.intents = [];
    }
    return obj;
  },
  fromPartial(object: DeepPartial<DelegatorIntent>): DelegatorIntent {
    const message = createBaseDelegatorIntent();
    message.delegator = object.delegator ?? "";
    message.intents = object.intents?.map(e => ValidatorIntent.fromPartial(e)) || [];
    return message;
  },
  fromSDK(object: DelegatorIntentSDKType): DelegatorIntent {
    return {
      delegator: object?.delegator,
      intents: Array.isArray(object?.intents) ? object.intents.map((e: any) => ValidatorIntent.fromSDK(e)) : []
    };
  },
  toSDK(message: DelegatorIntent): DelegatorIntentSDKType {
    const obj: any = {};
    obj.delegator = message.delegator;
    if (message.intents) {
      obj.intents = message.intents.map(e => e ? ValidatorIntent.toSDK(e) : undefined);
    } else {
      obj.intents = [];
    }
    return obj;
  },
  fromAmino(object: DelegatorIntentAmino): DelegatorIntent {
    return {
      delegator: object.delegator,
      intents: Array.isArray(object?.intents) ? object.intents.map((e: any) => ValidatorIntent.fromAmino(e)) : []
    };
  },
  toAmino(message: DelegatorIntent): DelegatorIntentAmino {
    const obj: any = {};
    obj.delegator = message.delegator;
    if (message.intents) {
      obj.intents = message.intents.map(e => e ? ValidatorIntent.toAmino(e) : undefined);
    } else {
      obj.intents = [];
    }
    return obj;
  },
  fromAminoMsg(object: DelegatorIntentAminoMsg): DelegatorIntent {
    return DelegatorIntent.fromAmino(object.value);
  },
  fromProtoMsg(message: DelegatorIntentProtoMsg): DelegatorIntent {
    return DelegatorIntent.decode(message.value);
  },
  toProto(message: DelegatorIntent): Uint8Array {
    return DelegatorIntent.encode(message).finish();
  },
  toProtoMsg(message: DelegatorIntent): DelegatorIntentProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.DelegatorIntent",
      value: DelegatorIntent.encode(message).finish()
    };
  }
};
function createBaseValidatorIntent(): ValidatorIntent {
  return {
    valoperAddress: "",
    weight: ""
  };
}
export const ValidatorIntent = {
  typeUrl: "/quicksilver.interchainstaking.v1.ValidatorIntent",
  encode(message: ValidatorIntent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.valoperAddress !== "") {
      writer.uint32(10).string(message.valoperAddress);
    }
    if (message.weight !== "") {
      writer.uint32(18).string(message.weight);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorIntent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidatorIntent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.valoperAddress = reader.string();
          break;
        case 2:
          message.weight = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): ValidatorIntent {
    const obj = createBaseValidatorIntent();
    if (isSet(object.valoper_address)) obj.valoperAddress = String(object.valoper_address);
    if (isSet(object.weight)) obj.weight = String(object.weight);
    return obj;
  },
  toJSON(message: ValidatorIntent): unknown {
    const obj: any = {};
    message.valoperAddress !== undefined && (obj.valoper_address = message.valoperAddress);
    message.weight !== undefined && (obj.weight = message.weight);
    return obj;
  },
  fromPartial(object: DeepPartial<ValidatorIntent>): ValidatorIntent {
    const message = createBaseValidatorIntent();
    message.valoperAddress = object.valoperAddress ?? "";
    message.weight = object.weight ?? "";
    return message;
  },
  fromSDK(object: ValidatorIntentSDKType): ValidatorIntent {
    return {
      valoperAddress: object?.valoper_address,
      weight: object?.weight
    };
  },
  toSDK(message: ValidatorIntent): ValidatorIntentSDKType {
    const obj: any = {};
    obj.valoper_address = message.valoperAddress;
    obj.weight = message.weight;
    return obj;
  },
  fromAmino(object: ValidatorIntentAmino): ValidatorIntent {
    return {
      valoperAddress: object.valoper_address,
      weight: object.weight
    };
  },
  toAmino(message: ValidatorIntent): ValidatorIntentAmino {
    const obj: any = {};
    obj.valoper_address = message.valoperAddress;
    obj.weight = message.weight;
    return obj;
  },
  fromAminoMsg(object: ValidatorIntentAminoMsg): ValidatorIntent {
    return ValidatorIntent.fromAmino(object.value);
  },
  fromProtoMsg(message: ValidatorIntentProtoMsg): ValidatorIntent {
    return ValidatorIntent.decode(message.value);
  },
  toProto(message: ValidatorIntent): Uint8Array {
    return ValidatorIntent.encode(message).finish();
  },
  toProtoMsg(message: ValidatorIntent): ValidatorIntentProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.ValidatorIntent",
      value: ValidatorIntent.encode(message).finish()
    };
  }
};
function createBaseDelegation(): Delegation {
  return {
    delegationAddress: "",
    validatorAddress: "",
    amount: Coin.fromPartial({}),
    height: Long.ZERO,
    redelegationEnd: Long.ZERO
  };
}
export const Delegation = {
  typeUrl: "/quicksilver.interchainstaking.v1.Delegation",
  encode(message: Delegation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.delegationAddress !== "") {
      writer.uint32(10).string(message.delegationAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
    }
    if (!message.height.isZero()) {
      writer.uint32(32).int64(message.height);
    }
    if (!message.redelegationEnd.isZero()) {
      writer.uint32(40).int64(message.redelegationEnd);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Delegation {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDelegation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.delegationAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.amount = Coin.decode(reader, reader.uint32());
          break;
        case 4:
          message.height = (reader.int64() as Long);
          break;
        case 5:
          message.redelegationEnd = (reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Delegation {
    const obj = createBaseDelegation();
    if (isSet(object.delegationAddress)) obj.delegationAddress = String(object.delegationAddress);
    if (isSet(object.validatorAddress)) obj.validatorAddress = String(object.validatorAddress);
    if (isSet(object.amount)) obj.amount = Coin.fromJSON(object.amount);
    if (isSet(object.height)) obj.height = Long.fromValue(object.height);
    if (isSet(object.redelegationEnd)) obj.redelegationEnd = Long.fromValue(object.redelegationEnd);
    return obj;
  },
  toJSON(message: Delegation): unknown {
    const obj: any = {};
    message.delegationAddress !== undefined && (obj.delegationAddress = message.delegationAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toJSON(message.amount) : undefined);
    message.height !== undefined && (obj.height = (message.height || Long.ZERO).toString());
    message.redelegationEnd !== undefined && (obj.redelegationEnd = (message.redelegationEnd || Long.ZERO).toString());
    return obj;
  },
  fromPartial(object: DeepPartial<Delegation>): Delegation {
    const message = createBaseDelegation();
    message.delegationAddress = object.delegationAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Coin.fromPartial(object.amount);
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = Long.fromValue(object.height);
    }
    if (object.redelegationEnd !== undefined && object.redelegationEnd !== null) {
      message.redelegationEnd = Long.fromValue(object.redelegationEnd);
    }
    return message;
  },
  fromSDK(object: DelegationSDKType): Delegation {
    return {
      delegationAddress: object?.delegation_address,
      validatorAddress: object?.validator_address,
      amount: object.amount ? Coin.fromSDK(object.amount) : undefined,
      height: object?.height,
      redelegationEnd: object?.redelegation_end
    };
  },
  toSDK(message: Delegation): DelegationSDKType {
    const obj: any = {};
    obj.delegation_address = message.delegationAddress;
    obj.validator_address = message.validatorAddress;
    message.amount !== undefined && (obj.amount = message.amount ? Coin.toSDK(message.amount) : undefined);
    obj.height = message.height;
    obj.redelegation_end = message.redelegationEnd;
    return obj;
  },
  fromAmino(object: DelegationAmino): Delegation {
    return {
      delegationAddress: object.delegation_address,
      validatorAddress: object.validator_address,
      amount: object?.amount ? Coin.fromAmino(object.amount) : undefined,
      height: Long.fromString(object.height),
      redelegationEnd: Long.fromString(object.redelegation_end)
    };
  },
  toAmino(message: Delegation): DelegationAmino {
    const obj: any = {};
    obj.delegation_address = message.delegationAddress;
    obj.validator_address = message.validatorAddress;
    obj.amount = message.amount ? Coin.toAmino(message.amount) : undefined;
    obj.height = message.height ? message.height.toString() : undefined;
    obj.redelegation_end = message.redelegationEnd ? message.redelegationEnd.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: DelegationAminoMsg): Delegation {
    return Delegation.fromAmino(object.value);
  },
  fromProtoMsg(message: DelegationProtoMsg): Delegation {
    return Delegation.decode(message.value);
  },
  toProto(message: Delegation): Uint8Array {
    return Delegation.encode(message).finish();
  },
  toProtoMsg(message: Delegation): DelegationProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.Delegation",
      value: Delegation.encode(message).finish()
    };
  }
};
function createBasePortConnectionTuple(): PortConnectionTuple {
  return {
    connectionId: "",
    portId: ""
  };
}
export const PortConnectionTuple = {
  typeUrl: "/quicksilver.interchainstaking.v1.PortConnectionTuple",
  encode(message: PortConnectionTuple, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.connectionId !== "") {
      writer.uint32(10).string(message.connectionId);
    }
    if (message.portId !== "") {
      writer.uint32(18).string(message.portId);
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): PortConnectionTuple {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePortConnectionTuple();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.connectionId = reader.string();
          break;
        case 2:
          message.portId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): PortConnectionTuple {
    const obj = createBasePortConnectionTuple();
    if (isSet(object.connectionId)) obj.connectionId = String(object.connectionId);
    if (isSet(object.portId)) obj.portId = String(object.portId);
    return obj;
  },
  toJSON(message: PortConnectionTuple): unknown {
    const obj: any = {};
    message.connectionId !== undefined && (obj.connectionId = message.connectionId);
    message.portId !== undefined && (obj.portId = message.portId);
    return obj;
  },
  fromPartial(object: DeepPartial<PortConnectionTuple>): PortConnectionTuple {
    const message = createBasePortConnectionTuple();
    message.connectionId = object.connectionId ?? "";
    message.portId = object.portId ?? "";
    return message;
  },
  fromSDK(object: PortConnectionTupleSDKType): PortConnectionTuple {
    return {
      connectionId: object?.connection_id,
      portId: object?.port_id
    };
  },
  toSDK(message: PortConnectionTuple): PortConnectionTupleSDKType {
    const obj: any = {};
    obj.connection_id = message.connectionId;
    obj.port_id = message.portId;
    return obj;
  },
  fromAmino(object: PortConnectionTupleAmino): PortConnectionTuple {
    return {
      connectionId: object.connection_id,
      portId: object.port_id
    };
  },
  toAmino(message: PortConnectionTuple): PortConnectionTupleAmino {
    const obj: any = {};
    obj.connection_id = message.connectionId;
    obj.port_id = message.portId;
    return obj;
  },
  fromAminoMsg(object: PortConnectionTupleAminoMsg): PortConnectionTuple {
    return PortConnectionTuple.fromAmino(object.value);
  },
  fromProtoMsg(message: PortConnectionTupleProtoMsg): PortConnectionTuple {
    return PortConnectionTuple.decode(message.value);
  },
  toProto(message: PortConnectionTuple): Uint8Array {
    return PortConnectionTuple.encode(message).finish();
  },
  toProtoMsg(message: PortConnectionTuple): PortConnectionTupleProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.PortConnectionTuple",
      value: PortConnectionTuple.encode(message).finish()
    };
  }
};
function createBaseReceipt(): Receipt {
  return {
    chainId: "",
    sender: "",
    txhash: "",
    amount: [],
    firstSeen: undefined,
    completed: undefined
  };
}
export const Receipt = {
  typeUrl: "/quicksilver.interchainstaking.v1.Receipt",
  encode(message: Receipt, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chainId !== "") {
      writer.uint32(10).string(message.chainId);
    }
    if (message.sender !== "") {
      writer.uint32(18).string(message.sender);
    }
    if (message.txhash !== "") {
      writer.uint32(26).string(message.txhash);
    }
    for (const v of message.amount) {
      Coin.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.firstSeen !== undefined) {
      Timestamp.encode(toTimestamp(message.firstSeen), writer.uint32(42).fork()).ldelim();
    }
    if (message.completed !== undefined) {
      Timestamp.encode(toTimestamp(message.completed), writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },
  decode(input: _m0.Reader | Uint8Array, length?: number): Receipt {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReceipt();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.string();
          break;
        case 2:
          message.sender = reader.string();
          break;
        case 3:
          message.txhash = reader.string();
          break;
        case 4:
          message.amount.push(Coin.decode(reader, reader.uint32()));
          break;
        case 5:
          message.firstSeen = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        case 6:
          message.completed = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Receipt {
    const obj = createBaseReceipt();
    if (isSet(object.chainId)) obj.chainId = String(object.chainId);
    if (isSet(object.sender)) obj.sender = String(object.sender);
    if (isSet(object.txhash)) obj.txhash = String(object.txhash);
    if (Array.isArray(object?.amount)) obj.amount = object.amount.map((e: any) => Coin.fromJSON(e));
    if (isSet(object.firstSeen)) obj.firstSeen = new Date(object.firstSeen);
    if (isSet(object.completed)) obj.completed = new Date(object.completed);
    return obj;
  },
  toJSON(message: Receipt): unknown {
    const obj: any = {};
    message.chainId !== undefined && (obj.chainId = message.chainId);
    message.sender !== undefined && (obj.sender = message.sender);
    message.txhash !== undefined && (obj.txhash = message.txhash);
    if (message.amount) {
      obj.amount = message.amount.map(e => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.amount = [];
    }
    message.firstSeen !== undefined && (obj.firstSeen = message.firstSeen.toISOString());
    message.completed !== undefined && (obj.completed = message.completed.toISOString());
    return obj;
  },
  fromPartial(object: DeepPartial<Receipt>): Receipt {
    const message = createBaseReceipt();
    message.chainId = object.chainId ?? "";
    message.sender = object.sender ?? "";
    message.txhash = object.txhash ?? "";
    message.amount = object.amount?.map(e => Coin.fromPartial(e)) || [];
    message.firstSeen = object.firstSeen ?? undefined;
    message.completed = object.completed ?? undefined;
    return message;
  },
  fromSDK(object: ReceiptSDKType): Receipt {
    return {
      chainId: object?.chain_id,
      sender: object?.sender,
      txhash: object?.txhash,
      amount: Array.isArray(object?.amount) ? object.amount.map((e: any) => Coin.fromSDK(e)) : [],
      firstSeen: object.first_seen ?? undefined,
      completed: object.completed ?? undefined
    };
  },
  toSDK(message: Receipt): ReceiptSDKType {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.sender = message.sender;
    obj.txhash = message.txhash;
    if (message.amount) {
      obj.amount = message.amount.map(e => e ? Coin.toSDK(e) : undefined);
    } else {
      obj.amount = [];
    }
    message.firstSeen !== undefined && (obj.first_seen = message.firstSeen ?? undefined);
    message.completed !== undefined && (obj.completed = message.completed ?? undefined);
    return obj;
  },
  fromAmino(object: ReceiptAmino): Receipt {
    return {
      chainId: object.chain_id,
      sender: object.sender,
      txhash: object.txhash,
      amount: Array.isArray(object?.amount) ? object.amount.map((e: any) => Coin.fromAmino(e)) : [],
      firstSeen: object?.first_seen,
      completed: object?.completed
    };
  },
  toAmino(message: Receipt): ReceiptAmino {
    const obj: any = {};
    obj.chain_id = message.chainId;
    obj.sender = message.sender;
    obj.txhash = message.txhash;
    if (message.amount) {
      obj.amount = message.amount.map(e => e ? Coin.toAmino(e) : undefined);
    } else {
      obj.amount = [];
    }
    obj.first_seen = message.firstSeen;
    obj.completed = message.completed;
    return obj;
  },
  fromAminoMsg(object: ReceiptAminoMsg): Receipt {
    return Receipt.fromAmino(object.value);
  },
  fromProtoMsg(message: ReceiptProtoMsg): Receipt {
    return Receipt.decode(message.value);
  },
  toProto(message: Receipt): Uint8Array {
    return Receipt.encode(message).finish();
  },
  toProtoMsg(message: Receipt): ReceiptProtoMsg {
    return {
      typeUrl: "/quicksilver.interchainstaking.v1.Receipt",
      value: Receipt.encode(message).finish()
    };
  }
};