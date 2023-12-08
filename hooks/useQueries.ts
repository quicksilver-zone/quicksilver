import { useChain } from '@cosmos-kit/react';
import { Zone } from '@hoangdv2429/quicksilverjs/dist/codegen/quicksilver/interchainstaking/v1/interchainstaking';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { cosmos } from 'interchain-query';

import { getCoin } from '@/utils';
import { parseValidators } from '@/utils/staking';

import { useGrpcQueryClient } from './useGrpcQueryClient';

const BigNumber = require('bignumber.js');
const Long = require('long');

export const useBalanceQuery = (chainName: string, address: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);
  const coin = getCoin(chainName);
  const balanceQuery = useQuery(
    ['balance', address],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const balance = await grpcQueryClient.cosmos.bank.v1beta1.balance({
        address: address || '',
        denom: coin.base,
      });

      return balance;
    },
    {
      enabled: !!grpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    balance: balanceQuery.data,
    isLoading: balanceQuery.isLoading,
    isError: balanceQuery.isError,
  };
};

export const useQBalanceQuery = (chainName: string, address: string, qAsset: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);
  const balanceQuery = useQuery(
    ['balance', qAsset],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const balance = await grpcQueryClient.cosmos.bank.v1beta1.balance({
        address: address || '',
        denom: 'uq' + qAsset,
      });

      return balance;
    },
    {
      enabled: !!grpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    balance: balanceQuery.data,
    isLoading: balanceQuery.isLoading,
    isError: balanceQuery.isError,
  };
};

export const useIntentQuery = (chainName: string, address: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);
  const { chain } = useChain(chainName);
  const chainId = chain.chain_id;
  const intentQuery = useQuery(
    ['intent', chainName],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const intent = await grpcQueryClient.quicksilver.interchainstaking.v1.delegatorIntent({
        chainId: chainId,
        delegatorAddress: address || '',
      });

      return intent;
    },
    {
      enabled: !!grpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    intent: intentQuery.data,
    isLoading: intentQuery.isLoading,
    isError: intentQuery.isError,
  };
};

export const useValidatorsQuery = (chainName: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);

  const fetchValidators = async (nextKey = new Uint8Array()) => {
    if (!grpcQueryClient) {
      throw new Error('RPC Client not ready');
    }

    const validators = await grpcQueryClient.cosmos.staking.v1beta1.validators({
      status: cosmos.staking.v1beta1.bondStatusToJSON(cosmos.staking.v1beta1.BondStatus.BOND_STATUS_BONDED),
      pagination: {
        key: nextKey,
        offset: Long.fromNumber(0),
        limit: Long.fromNumber(100),
        countTotal: true,
        reverse: false,
      },
    });
    return validators;
  };

  const validatorQuery = useQuery(
    ['validators', chainName],
    async () => {
      let allValidators: any[] = [];
      let nextKey = new Uint8Array();

      do {
        const response = await fetchValidators(nextKey);
        allValidators = allValidators.concat(response.validators);
        nextKey = response.pagination.next_key;
      } while (nextKey && nextKey.length > 0);
      const sorted = allValidators.sort((a, b) => new BigNumber(b.tokens).minus(a.tokens).toNumber());
      return parseValidators(sorted);
    },
    {
      enabled: !!grpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    validatorsData: validatorQuery.data,
    isLoading: validatorQuery.isLoading,
    isError: validatorQuery.isError,
  };
};

const fetchAPY = async (chainId: any) => {
  const res = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_DATA_API}/apr`);
  const { chains } = res.data;
  if (!chains) {
      return 0;
  }
  const chainInfo = chains.find((chain: { chain_id: any; }) => chain.chain_id === chainId);
  return chainInfo ? chainInfo.apr : 0;
};

export const useAPYQuery = (chainId: any) => {
  const query = useQuery(
      ['APY', chainId],
      () => fetchAPY(chainId),
      {
          staleTime: Infinity,
          enabled: !!chainId,
      }
  );

  return {
      APY: query.data,
      isLoading: query.isLoading,
      isError: query.isError,
  };
};

export const useZoneQuery = (chainId: string) => {
  return useQuery<Zone, Error>(
    ['zone', chainId],
    async () => {
      const res = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_API}/quicksilver/interchainstaking/v1/zones`);
      const { zones } = res.data;

      if (!zones || zones.length === 0) {
        throw new Error('Failed to query zones');
      }

      const apiZone = zones.find((z: { chain_id: string }) => z.chain_id === chainId);
      if (!apiZone) {
        throw new Error(`No zone with chain id ${chainId} found`);
      }

      // Parse or map the API zone data to your Zone interface
      const parsedZone: Zone = {
        connectionId: apiZone.connection_id,
        chainId: apiZone.chain_id,
        depositAddress: apiZone.deposit_address,
        withdrawalAddress: apiZone.withdrawal_address,
        performanceAddress: apiZone.performance_address,
        delegationAddress: apiZone.delegation_address,
        accountPrefix: apiZone.account_prefix,
        localDenom: apiZone.local_denom,
        baseDenom: apiZone.base_denom,
        redemptionRate: apiZone.redemption_rate,
        lastRedemptionRate: apiZone.last_redemption_rate,
        validators: apiZone.validators,
        aggregateIntent: apiZone.aggregate_intent,
        multiSend: apiZone.multi_send,
        liquidityModule: apiZone.liquidity_module,
        withdrawalWaitgroup: apiZone.withdrawal_waitgroup,
        ibcNextValidatorsHash: apiZone.ibc_next_validators_hash,
        validatorSelectionAllocation: apiZone.validator_selection_allocation,
        holdingsAllocation: apiZone.holdings_allocation,
        lastEpochHeight: apiZone.last_epoch_height,
        tvl: apiZone.tvl,
        unbondingPeriod: apiZone.unbonding_period,
        messagesPerTx: apiZone.messages_per_tx,
        // ... other fields as needed
      };

      return parsedZone;
    },
    {
      enabled: !!chainId,
    }
  );
};
