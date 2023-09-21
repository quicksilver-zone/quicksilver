import { useChain } from '@cosmos-kit/react';
import { useQuery } from '@tanstack/react-query';

import { getCoin } from '@/utils';

import { useRpcQueryClient } from './useQsRpcQuery';

export const useBalanceQuery = (chainName: string, address: string) => {
  const { rpcQueryClient } = useRpcQueryClient(chainName);
  const coin = getCoin(chainName);
  const balanceQuery = useQuery(
    ['balance', address], // Query key
    async () => {
      if (!rpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const balance = await rpcQueryClient.cosmos.bank.v1beta1.balance({
        address: address || '',
        denom: coin.base,
      });

      return balance;
    },
    {
      enabled: !!rpcQueryClient, // Query enabled condition
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
  const { rpcQueryClient } = useRpcQueryClient(chainName);
  const { chain } = useChain(chainName);
  const chainId = chain.chain_id;
  const intentQuery = useQuery(
    ['intent', chainName], // Query key
    async () => {
      if (!rpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const intent =
        await rpcQueryClient.quicksilver.interchainstaking.v1.delegatorIntent({
          chainId: chainId,
          delegatorAddress: address || '',
        });

      return intent;
    },
    {
      enabled: !!rpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    intent: intentQuery.data,
    isLoading: intentQuery.isLoading,
    isError: intentQuery.isError,
  };
};

export const useRewardsQuery = (chainName: string, address: string) => {
  const { rpcQueryClient } = useRpcQueryClient(chainName);
  const { chain } = useChain(chainName);
  const chainId = chain.chain_id;
  const rewardsQuery = useQuery(
    ['rewards', address], // Query key
    async () => {
      if (!rpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const rewards =
        await rpcQueryClient.quicksilver.interchainstaking.v1.withdrawalRecords(
          {
            chainId: chainId,
            delegatorAddress: address || '',
          },
        );

      return rewards;
    },
    {
      enabled: !!rpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    intent: rewardsQuery.data,
    isLoading: rewardsQuery.isLoading,
    isError: rewardsQuery.isError,
  };
};
