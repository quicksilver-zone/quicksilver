import { useQuery } from '@tanstack/react-query';

import { useRpcQueryClient } from './useQsRpcQuery';

export const useBalanceQuery = () => {
  const { rpcQueryClient } = useRpcQueryClient('quicksilver');
  const address = 'quick1uwqjtgjhjctjc45ugy7ev5prprhehc7wje9uwh';
  const balanceQuery = useQuery(
    ['balance', address], // Query key
    async () => {
      if (!rpcQueryClient?.cosmos?.bank?.v1beta1) {
        throw new Error('RPC Client not ready');
      }

      const balance = await rpcQueryClient.cosmos.bank.v1beta1.allBalances({
        address,
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
