import { useChain } from '@cosmos-kit/react';
import BigNumber from 'bignumber.js';
import {
  cosmos,
  createRpcQueryHooks,
  useRpcClient,
  useRpcEndpoint,
} from 'interchain-query';
import { useEffect, useMemo, useState } from 'react';

import {
  calcTotalDelegation,
  extendValidators,
  parseAnnualProvisions,
  parseDelegations,
  parseRewards,
  parseUnbondingDays,
  parseValidators,
} from '@/utils/staking';

import { useQueryHooks } from './useQueryHooks';
import { useRpcQueryClient } from './useRpcQueryClient';
import { getCoin, getExponent } from '../config';
import { shiftDigits } from '../utils';

(BigInt.prototype as any).toJSON = function () {
  return this.toString();
};

export const useStakingData = (chainName: string) => {
  const [isLoading, setIsLoading] = useState(false);

  const { rpcQueryClient } = useRpcQueryClient(chainName);

  const { cosmosQuery, isReady } = useQueryHooks(chainName);

  const { address } = useChain(chainName);

  const coin = getCoin(chainName);
  const exp = getExponent(chainName);

  const isDataQueryEnabled = !!address;

  const balanceQuery = cosmosQuery.bank.v1beta1.useBalance({
    request: {
      address: address || '',
      denom: coin.base,
    },
    options: {
      queryKey: ['balance', chainName],
      enabled: !!rpcQueryClient,
      select: ({ balance }) => shiftDigits(balance?.amount || '0', -exp),
      onError: (error) => {
        console.error('Error fetching balanceQuery:', error);
        balanceQuery.remove();
        balanceQuery.refetch();
      },
    },
  });

  const myValidatorsQuery = cosmosQuery.staking.v1beta1.useDelegatorValidators({
    request: {
      pagination: {
        key: new Uint8Array(),
        offset: 0n,
        limit: 200n,
        countTotal: true,
        reverse: false,
      },
      delegatorAddr: address || '',
    },
    options: {
      queryKey: ['delegatorValidators', chainName],
      enabled: isDataQueryEnabled,

      onError: (error) => {
        console.error('Error fetching myValidatorsQuery:', error);
        myValidatorsQuery.remove();
        myValidatorsQuery.refetch();
      },
      select: ({ validators }) => parseValidators(validators),
    },
  });

  const rewardsQuery =
    cosmosQuery.distribution.v1beta1.useDelegationTotalRewards({
      request: {
        delegatorAddress: address || '',
      },
      options: {
        queryKey: ['delegationTotalRewards', chainName],
        enabled: isDataQueryEnabled,
        select: (data) => parseRewards(data, coin.base, -exp),
        onError: (error) => {
          console.error('Error fetching rewardsQuery:', error);
          rewardsQuery.remove();
          rewardsQuery.refetch();
        },
      },
    });

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
      queryKey: ['validators', chainName],
      enabled: isDataQueryEnabled,
      select: ({ validators }) => {
        const sorted = validators.sort((a, b) =>
          new BigNumber(b.tokens).minus(a.tokens).toNumber(),
        );
        return parseValidators(sorted);
      },
      onError: (error) => {
        console.error('Error fetching validatorsQuery:', error);
        validatorsQuery.remove();
        validatorsQuery.refetch();
      },
    },
  });

  const delegationsQuery = cosmosQuery.staking.v1beta1.useDelegatorDelegations({
    request: {
      delegatorAddr: address || '',
      pagination: {
        key: new Uint8Array(),
        offset: 0n,
        limit: 100n,
        countTotal: true,
        reverse: false,
      },
    },
    options: {
      queryKey: ['delegatorDelegations', chainName],
      enabled: isDataQueryEnabled,
      select: ({ delegationResponses }) =>
        parseDelegations(delegationResponses, -exp),
      onError: (error) => {
        console.error('Error fetching delegationsQuery:', error);
        delegationsQuery.remove();
        delegationsQuery.refetch();
      },
    },
  });

  const unbondingDaysQuery = cosmosQuery.staking.v1beta1.useParams({
    options: {
      queryKey: ['params', chainName],
      enabled: isDataQueryEnabled,
      select: ({ params }) => parseUnbondingDays(params),
      onError: (error) => {
        console.error('Error fetching unbondingDaysQuery:', error);
        unbondingDaysQuery.remove();
        unbondingDaysQuery.refetch();
      },
    },
  });

  const communityTaxQuery = cosmosQuery.distribution.v1beta1.useParams({
    options: {
      queryKey: ['distributionParams', chainName],
      enabled: isDataQueryEnabled,
      select: ({ params }) => shiftDigits(params?.communityTax || '0', -18),
      onError: (error) => {
        console.error('Error fetching communityTaxQuery:', error);
        communityTaxQuery.remove();
        communityTaxQuery.refetch();
      },
    },
  });

  const allQueries = {
    balance: balanceQuery,
    myValidators: myValidatorsQuery,
    rewards: rewardsQuery,
    delegations: delegationsQuery,
    unbondingDays: unbondingDaysQuery,
    communityTax: communityTaxQuery,
  };

  const queriesWithUnchangingKeys = [
    allQueries.unbondingDays,
    allQueries.communityTax,
  ];

  const updatableQueriesAfterMutation = [
    allQueries.balance,
    allQueries.myValidators,
    allQueries.rewards,

    allQueries.delegations,
  ];

  useEffect(() => {
    queriesWithUnchangingKeys.forEach((query) => query.remove());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [chainName]);

  const loading = balanceQuery.isFetching || !isReady;

  useEffect(() => {
    setIsLoading(loading);
  }, [loading]);

  type AllQueries = typeof allQueries;

  type QueriesData = {
    [Key in keyof AllQueries]: NonNullable<AllQueries[Key]['data']>;
  };

  const data = useMemo(() => {
    if (isLoading) return;

    const queriesData = Object.fromEntries(
      Object.entries(allQueries).map(([key, query]) => [key, query.data]),
    ) as QueriesData;

    const { delegations, rewards } = queriesData;

    if (!rewards || !rewards.byValidators) {
      console.error('Rewards or byValidators is undefined:', rewards);
      return; // Handle this case, perhaps by returning a default value or setting some error state
    }

    const totalDelegated = calcTotalDelegation(delegations);

    return {
      ...queriesData,

      totalDelegated,
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isLoading]);

  const refetch = () => {
    Object.values(allQueries).forEach((query) => {
      query.remove();
      query.refetch();
    });
  };
  return { data, isLoading, refetch };
};
