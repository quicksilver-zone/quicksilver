import { useChain } from '@cosmos-kit/react';
import BigNumber from 'bignumber.js';
import {
  cosmos,
  useRpcClient,
  useRpcEndpoint,
  createRpcQueryHooks,
} from 'interchain-query';
import { useEffect, useMemo } from 'react';

import {
  extendValidators,
  parseAnnualProvisions,
  parseUnbondingDays,
  parseValidators,
} from '@/utils/staking';

import { getCoin, getExponent } from '../config';
import { shiftDigits } from '../utils';

(BigInt.prototype as any).toJSON = function () {
  return this.toString();
};

export const useStakingData = (chainName: string) => {
  const { address, getRpcEndpoint } = useChain(chainName);
  console.log('useStakingData', chainName);
  const rpcEndpointQuery = useRpcEndpoint({
    getter: getRpcEndpoint,
    extraKey: chainName,
    options: {
      enabled: !!chainName,
      staleTime: Infinity,
    },
  });

  const rpcClientQuery = useRpcClient({
    rpcEndpoint: rpcEndpointQuery.data || '',
    options: {
      enabled: !!rpcEndpointQuery.data,
      staleTime: Infinity,
    },
  });

  const { cosmos: cosmosQuery } = createRpcQueryHooks({
    rpc: rpcClientQuery.data,
  });

  const isDataQueryEnabled = !!rpcClientQuery.data;

  const validatorsQuery = cosmosQuery.staking.v1beta1.useValidators({
    request: {
      status: cosmos.staking.v1beta1.bondStatusToJSON(
        cosmos.staking.v1beta1.BondStatus.BOND_STATUS_BONDED,
      ),
      pagination: {
        key: new Uint8Array(),
        offset: 0n,
        limit: 200n,
        countTotal: true,
        reverse: false,
      },
    },
    options: {
      enabled: isDataQueryEnabled,
      select: ({ validators }) => {
        const sorted = validators.sort((a, b) =>
          new BigNumber(b.tokens).minus(a.tokens).toNumber(),
        );
        return parseValidators(sorted);
      },
    },
  });

  const unbondingDaysQuery = cosmosQuery.staking.v1beta1.useParams({
    options: {
      enabled: isDataQueryEnabled,
      select: ({ params }) => parseUnbondingDays(params),
    },
  });

  const annualProvisionsQuery = cosmosQuery.mint.v1beta1.useAnnualProvisions({
    options: {
      enabled: isDataQueryEnabled,
      select: parseAnnualProvisions,
      retry: false,
    },
  });

  const poolQuery = cosmosQuery.staking.v1beta1.usePool({
    options: {
      enabled: isDataQueryEnabled,
      select: ({ pool }) => pool,
    },
  });

  const communityTaxQuery = cosmosQuery.distribution.v1beta1.useParams({
    options: {
      enabled: isDataQueryEnabled,
      select: ({ params }) => shiftDigits(params?.communityTax || '0', -18),
    },
  });

  const allQueries = {
    allValidators: validatorsQuery,

    unbondingDays: unbondingDaysQuery,
    annualProvisions: annualProvisionsQuery,
    pool: poolQuery,
    communityTax: communityTaxQuery,
  };

  const queriesWithUnchangingKeys = [
    allQueries.unbondingDays,
    allQueries.annualProvisions,
    allQueries.pool,
    allQueries.communityTax,
    allQueries.allValidators,
  ];

  const updatableQueriesAfterMutation = [allQueries.allValidators];

  useEffect(() => {
    queriesWithUnchangingKeys.forEach((query) => query.remove());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chainName]);

  const isInitialFetching = Object.values(allQueries).some(
    ({ isLoading }) => isLoading,
  );

  const isRefetching = Object.values(allQueries).some(
    ({ isRefetching }) => isRefetching,
  );

  const isLoading = isInitialFetching || isRefetching;

  type AllQueries = typeof allQueries;

  type QueriesData = {
    [Key in keyof AllQueries]: NonNullable<AllQueries[Key]['data']>;
  };

  const data = useMemo(() => {
    if (isLoading) return;

    const queriesData = Object.fromEntries(
      Object.entries(allQueries).map(([key, query]) => [key, query.data]),
    ) as QueriesData;

    const { allValidators, annualProvisions, communityTax, pool } = queriesData;

    const chainMetadata = {
      annualProvisions,
      communityTax,
      pool,
    };

    const extendedAllValidators = extendValidators(
      allValidators,

      chainMetadata,
    );

    return {
      ...queriesData,
      allValidators: extendedAllValidators,
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isLoading]);

  const refetch = () => {
    updatableQueriesAfterMutation.forEach((query) => query.remove());
    updatableQueriesAfterMutation.forEach((query) => query.refetch());
  };

  return { data, isLoading, refetch };
};
