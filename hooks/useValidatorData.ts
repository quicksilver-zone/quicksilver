import { useQuery } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import BigNumber from 'bignumber.js';
import { cosmos } from 'interchain-query';
import { useEffect, useMemo, useState } from 'react';

import { parseValidators } from '@/utils/staking';

import { useGrpcQueryClient } from './useGrpcQueryClient';
import { useQueryHooks } from './useQueryHooks';
import { useRpcQueryClient } from './useRpcQueryClient';

(BigInt.prototype as any).toJSON = function () {
  return this.toString();
};

export const useValidatorData = (chainName: string, address: string) => {
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
