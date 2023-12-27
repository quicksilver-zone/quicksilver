import { Coin, decodeCosmosSdkDecFromProto } from '@cosmjs/stargate';
import BigNumber from 'bignumber.js';
import { QueryDelegationTotalRewardsResponse } from 'interchain-query/cosmos/distribution/v1beta1/query';
import { QueryAnnualProvisionsResponse } from 'interchain-query/cosmos/mint/v1beta1/query';
import { QueryDelegatorDelegationsResponse, QueryParamsResponse } from 'interchain-query/cosmos/staking/v1beta1/query';
import { Pool, Validator } from 'interchain-query/cosmos/staking/v1beta1/staking';
import * as bech32 from 'bech32';
import * as CryptoJS from 'crypto-js';

import { decodeUint8Arr, isGreaterThanZero, shiftDigits, toNumber } from '.';
import { Any } from 'interchain-query/google/protobuf/any';

interface ConsensusPubkey {
  '@type': string;
  key: string; 
}


const DAY_TO_SECONDS = 24 * 60 * 60;
const ZERO = '0';

export const calcStakingApr = ({ pool, commission, communityTax, annualProvisions }: ChainMetaData & { commission: string }) => {
  const totalSupply = new BigNumber(pool?.bondedTokens || 0).plus(pool?.notBondedTokens || 0);

  const bondedTokensRatio = new BigNumber(pool?.bondedTokens || 0).div(totalSupply);

  const inflation = new BigNumber(annualProvisions || 0).div(totalSupply);

  const one = new BigNumber(1);

  return inflation
    .multipliedBy(one.minus(communityTax || 0))
    .div(bondedTokensRatio)
    .multipliedBy(one.minus(commission))
    .shiftedBy(2)
    .decimalPlaces(2, BigNumber.ROUND_DOWN)
    .toString();
};

export type ParsedValidator = ReturnType<typeof parseValidators>[0];

function extractValconsPrefix(operatorAddress: string): string {
  const prefixEndIndex = operatorAddress.indexOf('valoper');
  const chainPrefix = operatorAddress.substring(0, prefixEndIndex);
  return `${chainPrefix}valcons`;
}


export const parseValidators = (validators: Validator[]) => {
  return validators.map((validator) => {
    const commissionRate = validator.commission?.commission_rates?.rate || ZERO;
    const commissionPercentage = parseFloat(commissionRate) * 100;

    const valconsPrefix = extractValconsPrefix(validator.operator_address);
    const valconsAddress = getValconsAddress(validator.consensus_pubkey, valconsPrefix);

    return {
      consensusPubkey: validator.consensus_pubkey || '',
      valconsAddress, 
      description: validator.description?.details || '',
      name: validator.description?.moniker || '',
      identity: validator.description?.identity || '',
      address: validator.operator_address || '',
      commission: commissionPercentage.toFixed() + '%',
      votingPower: toNumber(shiftDigits(validator.tokens, -6, 0), 0),
    };
  });
};


function getValconsAddress(consensusPubkeyAny: any, valconsPrefix: string) {
  const consensusPubkey = consensusPubkeyAny as ConsensusPubkey;
  
  if (!consensusPubkey || !consensusPubkey.key) {
    return ''; 
  }

  const consensusPubkeyBytes = new Uint8Array(consensusPubkeyAny!.value);
  const decoded = Buffer.from(consensusPubkeyBytes).toString('base64');

  const bytes = CryptoJS.enc.Base64.parse(decoded);
  
  const valconsWords: number[] = bech32.bech32.toWords(Array.from(new Uint8Array(bytes.words)));
  const valconsAddress: string = bech32.bech32.encode(valconsPrefix, valconsWords);

  return valconsAddress;
}

export type ExtendedValidator = ReturnType<typeof extendValidators>[0];

export type ChainMetaData = {
  annualProvisions: string;
  communityTax: string;
  pool: Pool;
};

export const extendValidators = (validators: ParsedValidator[] = [], chainMetadata: ChainMetaData) => {
  const { annualProvisions, communityTax, pool } = chainMetadata;


  return validators.map((validator) => {
    const apr = annualProvisions
      ? calcStakingApr({
          annualProvisions,
            //@ts-ignore
          commission: validator.commission,
          communityTax,
          pool,
        })
      : null;

    return {
      ...validator,

    };
  });
};

const findAndDecodeReward = (coins: Coin[], denom: string, exponent: number) => {
  const amount = coins.find((coin) => coin.denom === denom)?.amount || ZERO;
  const decodedAmount = decodeCosmosSdkDecFromProto(amount).toString();
  return shiftDigits(decodedAmount, exponent);
};

export type ParsedRewards = ReturnType<typeof parseRewards>;

export const parseRewards = ({ rewards, total }: QueryDelegationTotalRewardsResponse, denom: string, exponent: number) => {
  const totalReward = findAndDecodeReward(total, denom, exponent);

  const rewardsParsed = rewards.map(({ reward, validatorAddress }) => ({
    validatorAddress,
    amount: findAndDecodeReward(reward, denom, exponent),
  }));

  return {
    byValidators: rewardsParsed,
    total: totalReward,
  };
};

export type ParsedDelegations = ReturnType<typeof parseDelegations>;

export const parseDelegations = (delegations: QueryDelegatorDelegationsResponse['delegationResponses'], exponent: number) => {
  return delegations.map(({ balance, delegation }) => ({
    validatorAddress: delegation?.validatorAddress || '',
    amount: shiftDigits(balance?.amount || ZERO, exponent),
  }));
};

export const calcTotalDelegation = (delegations: ParsedDelegations) => {
  if (!delegations) {
    console.error('Delegations are undefined:', delegations);
    return '0'; // Handle this case accordingly
  }

  return delegations.reduce((prev, cur) => prev.plus(cur.amount), new BigNumber(0)).toString();
};
export const parseUnbondingDays = (params: QueryParamsResponse['params']) => {
  return new BigNumber(Number(params?.unbondingTime?.seconds || 0)).div(DAY_TO_SECONDS).decimalPlaces(0).toString();
};

export const parseAnnualProvisions = (data: QueryAnnualProvisionsResponse) => {
  const res = shiftDigits(decodeUint8Arr(data?.annualProvisions), -18);
  return isGreaterThanZero(res) ? res : null;
};
