import { useChain } from '@cosmos-kit/react';
import { QueryValidatorsResponse } from 'interchain-query/cosmos/staking/v1beta1/query';
import { useEffect, useState } from 'react';

import { useQueryHooks } from './useQueryHooks';

(BigInt.prototype as any).toJSON = function () {
  return this.toString();
};

const getPagination = (
  limit: bigint,
  reverse: boolean = false,
) => ({
  limit,
  reverse,
  key: new Uint8Array(),
  offset: 0n,
  countTotal: true,
});

export const useValidatorData = (
  chainName: string,
) => {
  const [isLoading, setIsLoading] =
    useState(false);

  const { cosmosQuery, isReady } =
    useQueryHooks(chainName);

  const validatorsQuery =
    cosmosQuery.staking.v1beta1.useValidators({
      request: {
        status: 'BOND_STATUS_BONDED',
        pagination: getPagination(50n, true),
      },
      options: {
        enabled: isReady,
        staleTime: Infinity,
        select: (
          response: QueryValidatorsResponse,
        ) => response.validators,
      },
    });

  useEffect(() => {
    if (validatorsQuery.isFetching) {
      setIsLoading(true);
    } else {
      setIsLoading(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [validatorsQuery.isFetching]);

  const refetch = () => {
    validatorsQuery.remove();
    validatorsQuery.refetch();
  };

  return {
    data: { ...validatorsQuery },
    isLoading,
    refetch,
  };
};
