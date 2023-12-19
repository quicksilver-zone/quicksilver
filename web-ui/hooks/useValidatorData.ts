import BigNumber from 'bignumber.js';
import { cosmos } from 'interchain-query';
import { useEffect, useMemo, useState } from 'react';

import { parseValidators } from '@/utils/staking';

import { useQueryHooks } from './useQueryHooks';
import { useRpcQueryClient } from './useRpcQueryClient';

(BigInt.prototype as any).toJSON = function () {
  return this.toString();
};

export const useValidatorData = (chainName: string) => {
  const [isLoading, setIsLoading] = useState(false);

  const { rpcQueryClient } = useRpcQueryClient(chainName);

  const { cosmosQuery, isReady } = useQueryHooks(chainName);

  const validatorsQuery = cosmosQuery.staking.v1beta1.useValidators({
    request: {
      status: cosmos.staking.v1beta1.bondStatusToJSON(cosmos.staking.v1beta1.BondStatus.BOND_STATUS_BONDED),
      pagination: {
        key: new Uint8Array(),
        offset: 0n,
        limit: 200n,
        countTotal: true,
        reverse: false,
      },
    },
    options: {
      queryKey: ['validators', chainName],
      enabled: !!rpcQueryClient?.cosmos?.staking?.v1beta1.validator,
      select: ({ validators }) => {
        const sorted = validators.sort((a, b) => new BigNumber(b.tokens).minus(a.tokens).toNumber());
        return parseValidators(sorted);
      },
      onError: (error) => {
        console.error('Error fetching validators:', error);
        validatorsQuery.remove();
        validatorsQuery.refetch();
      },
    },
  });

  const loading = validatorsQuery.isFetching || !isReady;

  useEffect(() => {
    setIsLoading(loading);
  }, [loading]);

  type SingleQueriesData = {
    validators: NonNullable<(typeof validatorsQuery)['data']>;
  };

  const singleQueriesData = useMemo(() => {
    if (validatorsQuery.isFetching || !isReady) return;
    return {
      validators: validatorsQuery.data,
    } as SingleQueriesData;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [validatorsQuery.isFetching, isReady]);

  const refetch = () => {
    validatorsQuery.remove();
    validatorsQuery.refetch();
  };

  return {
    data: singleQueriesData,
    isLoading,
    refetch,
  };
};
