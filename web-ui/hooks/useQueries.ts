import { useChain } from '@cosmos-kit/react';
import {SkipRouter, SKIP_API_URL} from '@skip-router/core';
import { useQueries, useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { cosmos } from 'interchain-query';
import { QueryAllBalancesResponse } from 'quicksilverjs/dist/codegen/cosmos/bank/v1beta1/query';
import { Zone } from 'quicksilverjs/dist/codegen/quicksilver/interchainstaking/v1/interchainstaking';

import { useGrpcQueryClient } from './useGrpcQueryClient';

import { getCoin, getLogoUrls } from '@/utils';
import { ExtendedValidator, parseValidators } from '@/utils/staking';


type WithdrawalRecord = {
  chain_id: string;
  delegator: string;
  distribution: { valoper: string; amount: string }[];
  recipient: string;
  amount: { denom: string; amount: string }[];
  burn_amount: { denom: string; amount: string };
  txhash: string;
  status: number;
  completion_time: string;
  requeued: boolean;
  acknowledged: boolean;
  epoch_number: string;
};

type WithdrawalsResponse = {
  withdrawals: WithdrawalRecord[];
  pagination: any; 
};


type Amount = {
  denom: string;
  amount: string;
};


interface Asset {
  Type: string;
  Amount: AssetAmount[];
}


type Errors = {
  Errors: any; 
};


type LiquidRewardsData = {
  messages: any[]; 
  assets: {
    [key: string]: [
      {
        Type: string;
        Amount: Amount[];
      }
    ];
  };
  errors: Errors;
};


type UseLiquidRewardsQueryReturnType = {
  liquidRewards: LiquidRewardsData | undefined;
  isLoading: boolean;
  isError: boolean;
  refetch: () => void;
};

interface ProofOp {
  type: string;
  key: Uint8Array; 
  data: Uint8Array; 
}

interface Proof {
  key: Uint8Array;  
  data: Uint8Array; 
  proofOps: {
    ops: ProofOp[];
  };
  height: Long; 
  proofTypes: string;
}

interface Message {
  user_address: string;
  zone: string;
  src_zone: string;
  claim_type: number;
  proofs: Proof[];
 
}

interface AssetAmount {
  denom: string;
  amount: string;
}

interface LiquidEpochData {
  messages: Message[];
  assets: { [key: string]: Asset[] };
  errors: Record<string, unknown>; 
}


interface UseLiquidEpochQueryReturnType {
  liquidEpoch: LiquidEpochData | undefined;
  isLoading: boolean;
  isError: boolean;
}

const skipClient = new SkipRouter({
  apiURL: SKIP_API_URL,
});

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
      enabled: !!grpcQueryClient && !!address,
      staleTime: Infinity,
    },
  );

  return {
    balance: balanceQuery.data,
    isLoading: balanceQuery.isLoading,
    isError: balanceQuery.isError,
  };
};

export const useIncorrectAuthChecker = (address: string) => {
  const authQuery = useQuery(
    ['authWrong', address],
    async () => {
      if (!address) {
        throw new Error('Address is undefined or null');
      }

      try {
        const url = `https://lcd.quicksilver.zone/cosmos/authz/v1beta1/grants?granter=${address}&grantee=quick1w5ennfhdqrpyvewf35sv3y3t8yuzwq29mrmyal&msgTypeUrl=/quicksilver.participationrewards.v1.MsgSubmitClaim`;
        const response = await axios.get(url);
        return { data: response.data, error: null };
      } catch (error) {
        // Capture and return error
        return { data: null, error: error };
      }
    },
    {
      enabled: !!address,
      staleTime: Infinity,
    },
  );

  return {  
    authData: authQuery.data?.data,
    authError: authQuery.data?.error,
    isLoading: authQuery.isLoading,
    isError: authQuery.isError,
  };
};

export const useAuthChecker = (address: string) => {
  const authQuery = useQuery(
    ['auth', address],
    async () => {
      if (!address) {
        throw new Error('Address is undefined or null');
      }

      try {
        const url = `https://lcd.quicksilver.zone/cosmos/authz/v1beta1/grants?granter=${address}&grantee=quick1psevptdp90jad76zt9y9x2nga686hutgmasmwd&msgTypeUrl=/quicksilver.participationrewards.v1.MsgSubmitClaim`;
        const response = await axios.get(url);
        return { data: response.data, error: null };
      } catch (error) {
        // Capture and return error
        return { data: null, error: error };
      }
    },
    {
      enabled: !!address,
      staleTime: Infinity,
    },
  );

  return {
    authData: authQuery.data?.data,
    authError: authQuery.data?.error,
    isLoading: authQuery.isLoading,
    isError: authQuery.isError,
  };
};

export const useParamsQuery = (chainName: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);

  const paramsQuery = useQuery(
    ['params'],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const params = await grpcQueryClient.cosmos.mint.v1beta1.annualProvisions({


      });

      return params;
    },
    {
      enabled: !!grpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    params: paramsQuery.data,
    isLoading: paramsQuery.isLoading,
    isError: paramsQuery.isError,
  };

}

export const useAllBalancesQuery = (chainName: string, address: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);

  const balancesQuery = useQuery(
    ['balances', address],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }
      const next_key = new Uint8Array()
      const balance = await grpcQueryClient.cosmos.bank.v1beta1.allBalances({
        address: address || '',
        pagination: {
          key: next_key,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          countTotal: true,
          reverse: false,
        },
      });

      return balance;
    },
    {
      enabled: !!grpcQueryClient && !!address,
      staleTime: Infinity,
    },
  );

  return {
    balance: balancesQuery.data,
    isLoading: balancesQuery.isLoading,
    isError: balancesQuery.isError,
  };
};

export const useIbcBalanceQuery = (chainName: string, address: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);
  const balanceQuery = useQuery(
    ['balance', address],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }
      const next_key = new Uint8Array()
      const balance = await grpcQueryClient.cosmos.bank.v1beta1.allBalances({
        address: address || '',
        pagination: {
          key: next_key,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          countTotal: true,
          reverse: false,
        },
      });

      return balance;
    },
    {
      enabled: !!grpcQueryClient && !!address,
      staleTime: Infinity,
    },
  );

  return {
    balance: balanceQuery.data,
    isLoading: balanceQuery.isLoading,
    isError: balanceQuery.isError,
  };
};


export const useTokenPriceQuery = (tokenSymbol: string) => {
  const fetchTokenPrice = async () => {
    if (!tokenSymbol) {
      throw new Error('Token symbol is required');
    }

    const response = await axios.get(`https://api-osmosis.imperator.co/tokens/v2/price/${tokenSymbol}`);
    return response.data;
  };

  return useQuery(['tokenPrice', tokenSymbol], fetchTokenPrice, {
    enabled: !!tokenSymbol,
    staleTime: 300000, 
  });
};

export const useQBalanceQuery = (chainName: string, address: string, qAsset: string, liveNetworks?: string[], chainId?: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);

  const isLive = liveNetworks?.includes(chainId ?? '');

  const balanceQuery = useQuery(
    ['balance', qAsset],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }
      const denom = qAsset === 'dydx' ? 'aq'+ qAsset : 'uq' + qAsset;

      const balance = await grpcQueryClient.cosmos.bank.v1beta1.balance({
        address: address || '',
        denom: denom,
      });

      return balance;
    },
    {
      enabled: !!grpcQueryClient && !!address && isLive,
      staleTime: Infinity,
    },
  );

  return {
    balance: balanceQuery.data,
    isLoading: balanceQuery.isLoading,
    isError: balanceQuery.isError,
    refetch: balanceQuery.refetch,
  };
};

export const useQBalancesQuery = (chainName: string, address: string, grpcQueryClient: { cosmos: { bank: { v1beta1: { allBalances: (arg0: { address: string; pagination: { key: Uint8Array; offset: any; limit: any; countTotal: boolean; reverse: boolean; }; }) => any; }; }; }; } | undefined) => {


  const allQBalanceQuery = useQuery(
    ['balances', address],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }



      const next_key = new Uint8Array();
      const balance = await grpcQueryClient.cosmos.bank.v1beta1.allBalances({
        address: address || '',
        pagination: {
          key: next_key,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          countTotal: true,
          reverse: false,
        },
      });

      return balance;
    },
    {
      enabled: !!grpcQueryClient && !!address,
      staleTime: 0,
    },
  );

  const sortAndFindQAssets = (balances: QueryAllBalancesResponse) => {
    return balances.balances?.filter(b => 
        (b.denom.startsWith('uq') || b.denom.startsWith('aq')) &&
        !b.denom.startsWith('uqck') &&
        !b.denom.includes('ibc/') 
      )
      .sort((a, b) => a.denom.localeCompare(b.denom));
  };


  return {
    qbalance: sortAndFindQAssets(allQBalanceQuery.data ?? {} as QueryAllBalancesResponse),
    qIsLoading: allQBalanceQuery.isLoading,
    qIsError: allQBalanceQuery.isError,
    qRefetch: allQBalanceQuery.refetch,
  };
};

export const useIntentQuery = (chainName: string, address: string) => {
  const { grpcQueryClient } = useGrpcQueryClient('quicksilver');
  const { chain } = useChain(chainName);
  const env = process.env.NEXT_PUBLIC_CHAIN_ENV;
  const baseApiUrl = env === 'testnet' ? 'https://lcd.test.quicksilver.zone' : 'https://lcd.quicksilver.zone';
  
  // Determine the chain ID based on the chain name
  let chainId = chain.chain_id;
  if (chainName === 'osmosistestnet') {
    chainId = 'osmo-test-5';
  } else if (chainName === 'cosmoshubtestnet') {
    chainId = 'provider';
  } else if (chainName === 'stargazetestnet') {
    chainId = 'elgafar-1';
  } else if (chainName === 'osmo-test-5') {
    chainId = 'osmosistestnet';
  } else {
    chainId = chain.chain_id;
  }

  const intentQuery = useQuery(
    ['intent', chainName, address], 
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }

      const intent = await axios.get(`${baseApiUrl}/quicksilver/interchainstaking/v1/zones/${chainId}/delegator_intent/${address}`)
      return intent;
    },
    {
      enabled: !!grpcQueryClient && !!address, 
      staleTime: Infinity,
      cacheTime: 0, 
    },
  );

  return {
    intent: intentQuery.data,
    isLoading: intentQuery.isLoading,
    isError: intentQuery.isError,
    refetch: intentQuery.refetch,
  };
};

export const useLiquidRewardsQuery = (address: string): UseLiquidRewardsQueryReturnType => {
  const liquidRewardsQuery = useQuery(
    ['liquidRewards', address],
    async () => {
      if (!address) {
        throw new Error('Address is not avaialble');
      }

      const response = await axios.get<LiquidRewardsData>(`https://claim.quicksilver.zone/${address}/current`);
      return response.data;
    },
    {
      enabled:!!address,
      staleTime: 0,
    },
  );

  return {
    liquidRewards: liquidRewardsQuery.data,
    isLoading: liquidRewardsQuery.isLoading,
    isError: liquidRewardsQuery.isError,
    refetch: liquidRewardsQuery.refetch,
  };

}

export const useLiquidEpochQuery = (address: string): UseLiquidEpochQueryReturnType => {
  const liquidEpochQuery = useQuery(
    ['liquidEpoch', address],
    async () => {
      if (!address) {
        throw new Error('Address is not available');
      }

      const response = await axios.get<LiquidEpochData>(`https://claim.quicksilver.zone/${address}/epoch`);


      if (response.data.messages.length === 0) {
        console.error('No messages found'); 
      }

      return response.data;
    },
    {
      enabled: !!address,
      staleTime: Infinity,
    },
  );

  return {
    liquidEpoch: liquidEpochQuery.data,
    isLoading: liquidEpochQuery.isLoading,
    isError: liquidEpochQuery.isError,
  };
};

export const useUnbondingQuery = (chainName: string, address: string) => {
  const env = process.env.NEXT_PUBLIC_CHAIN_ENV;
  const baseApiUrl = env === 'testnet' ? 'https://lcd.test.quicksilver.zone' : 'https://lcd.quicksilver.zone';
  
  const { chain } = useChain(chainName);
  let chainId = chain.chain_id;
  if (chainName === 'osmosistestnet') {
    chainId = 'osmo-test-5';
  } else if (chainName === 'stargazetestnet') {
    chainId = 'elgafar-1';
  } else if (chainName === 'osmo-test-5') {
    chainId = 'osmosistestnet';
 
  } else {

    chainId = chain.chain_id;
  }
  const unbondingQuery = useQuery(
    ['unbond', chainName, address],
    async () => {
      const url = `${baseApiUrl}/quicksilver/interchainstaking/v1/zones/${chainId}/withdrawal_records/${address}`;
      const response = await axios.get<WithdrawalsResponse>(url);
      return response.data; 
    },
    {
      enabled: !!chainId && !!address, 
      staleTime: Infinity,
    },
  );

  return {
    unbondingData: unbondingQuery.data,
    isLoading: unbondingQuery.isLoading,
    isError: unbondingQuery.isError,
  };
};


export const useValidatorsQuery = (chainName: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);

  const fetchValidators = async (key = new Uint8Array()) => {
    if (!grpcQueryClient) {
      throw new Error('RPC Client not ready');
    }

    const validators = await grpcQueryClient.cosmos.staking.v1beta1.validators({
      status: cosmos.staking.v1beta1.bondStatusToJSON(cosmos.staking.v1beta1.BondStatus.BOND_STATUS_BONDED),
      pagination: {
        key: key,
        offset: Long.fromNumber(0),
        limit: Long.fromNumber(500),
        countTotal: true,
        reverse: false,
      },
    });

    return validators;
  };


  //TODO: migrate this to use evince cache endpoint.
  const validatorQuery = useQuery(
    ['validators', chainName],
    async () => {
      let allValidators: any[] = [];
      let next_key = new Uint8Array();

      do {
        const response = await fetchValidators(next_key);
        allValidators = allValidators.concat(response.validators);
        next_key = response.pagination.next_key ?? new Uint8Array();
      } while (next_key && next_key.length > 0);
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

export const useTokenPrices = (tokens: string[]) => {
  const fetchTokenPrices = async () => {
    return Promise.all(
      tokens.map(async (token) => {
        try {
          const response = await axios.get(`https://api-osmosis.imperator.co/tokens/v2/price/${token}`);
          return { token, price: response.data.price };
        } catch (error) {
          console.error(`Error fetching price for ${token}:`, error);
          return { token, price: null };
        }
      })
    );
  };

  return useQuery(['tokenPrices', ...tokens], fetchTokenPrices, {
    enabled: !!tokens,
    staleTime: Infinity, 
  });
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

const fetchAPYs = async () => {
  const res = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_DATA_API}/apr`);
  const { chains } = res.data;
  if (!chains) {
    return {};
  }
  const apys = chains.reduce((acc: { [x: string]: any; }, chain: { chain_id: string | number; apr: any; }) => {
    acc[chain.chain_id] = chain.apr;
    return acc;
  }, {});
  return apys;
};



export const useAPYQuery = (chainId: any, liveNetworks?: string[] ) => {
  const isLive = liveNetworks?.some(network => network === chainId);
  const query = useQuery(
      ['APY', chainId],
      () => fetchAPY(chainId),
      {
          staleTime: Infinity,
          enabled: !!chainId && isLive,
      }
  );

  return {
      APY: query.data,
      isLoading: query.isLoading,
      isError: query.isError,
  };
};

export const useAPYsQuery = () => {
  const query = useQuery(
    ['APY'],
    () => fetchAPYs(),
    {
      staleTime: Infinity,
      enabled: true,
    }
  );

  return {
    APYs: query.data,
    APYsLoading: query.isLoading,
    APYsError: query.isError,
    APYsRefetch : query.refetch,
  };
};

function parseZone(apiZone: any): Zone {

  return {
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
    decimals: apiZone.decimals,
    returnToSender: apiZone.return_to_sender,
    unbondingEnabled: apiZone.unbonding_enabled,
    depositsEnabled: apiZone.deposits_enabled,
    is118: apiZone.is118,
    subzoneInfo: apiZone.subzoneInfo,
  };
}

export function useZonesData(networks: { chainId: string }[]) {
  return useQueries({
    queries: networks.map(({ chainId }) => ({
      queryKey: ['zone', chainId],
      queryFn: async () => {
        const response = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_API}/quicksilver/interchainstaking/v1/zones`);
        const zones: any[] = response.data.zones; 
        const apiZone = zones.find(z => z.chain_id === chainId);
        if (!apiZone) {
          throw new Error(`No zone with chain id ${chainId} found`);
        }
        return parseZone(apiZone); 
      },
      enabled: !!chainId,
    }))
  });
}

export const useZoneQuery = (chainId: string, liveNetworks?: string[]) => {
  const isLive = liveNetworks?.some(network => network === chainId);
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
    decimals: apiZone.decimals,
    returnToSender: apiZone.return_to_sender,
    unbondingEnabled: apiZone.unbonding_enabled,
    depositsEnabled: apiZone.deposits_enabled,
    is118: apiZone.is118,
    subzoneInfo: apiZone.subzoneInfo,
      };

      return parsedZone;
    },
    {
      enabled: !!chainId && isLive
    }
  );
};

export const useRedemptionRatesQuery = () => {
  const query = useQuery(
    ['zones'],
    async () => {
      const res = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_API}/quicksilver/interchainstaking/v1/zones`);
      const { zones } = res.data;

      if (!zones || zones.length === 0) {
        throw new Error('Failed to query zones');
      }
      

      const rates = zones.reduce((acc: { [x: string]: { current: number; last: number; }; }, zone: { chain_id: string | number; redemption_rate: string; last_redemption_rate: string; }) => {
        acc[zone.chain_id] = {
          current: parseFloat(zone.redemption_rate),
          last: parseFloat(zone.last_redemption_rate),
        };
        return acc;
      }, {});

      return rates;
    },
    {
      staleTime: Infinity,
      enabled: true,
    }
  );

  return {
    redemptionRates: query.data,
    redemptionLoading: query.isLoading,
    redemptionError: query.isError,
    redemptionRefetch: query.refetch,
  };
};

export const useValidatorLogos = (
  chainName: string,
  validators: ExtendedValidator[]
) => {
  const { data, isLoading } = useQuery({
    queryKey: ['validatorLogos', chainName],
    queryFn: () => getLogoUrls(validators, chainName),
    enabled: validators.length > 0,
    staleTime: Infinity,
  });

  return { data, isLoading };
};

export const useMissedBlocks = (chainName: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);

  const fetchMissedBlocks = async () => {
    if (!grpcQueryClient) {
      throw new Error('RPC Client not ready');
    }
  
    let allMissedBlocks: any[] = [];
    let next_key = new Uint8Array();
  
    do {
      const response = await grpcQueryClient.cosmos.slashing.v1beta1.signingInfos({
        pagination: {
          key: next_key,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          countTotal: true,
          reverse: false,
        },
      });
  
      // Filter out entries without an address
      const filteredMissedBlocks = response.info.filter(block => {
        const hasAddress = block.address && block.address.trim() !== '';
        const notTombstoned = !block.tombstoned;
    
        return hasAddress && notTombstoned;
      });
      
      allMissedBlocks = allMissedBlocks.concat(filteredMissedBlocks);
      next_key = response.pagination?.next_key ?? new Uint8Array();
    } while (next_key && next_key.length > 0);
  
    return allMissedBlocks;
  };
  
  const missedBlocksQuery = useQuery({
    queryKey: ['missedBlocks', chainName],
    queryFn: fetchMissedBlocks,
    enabled: !!grpcQueryClient,
    staleTime: Infinity,
    onError: (error) => {
      console.error('Error in fetching Missed Blocks:', error);
    },
  });

  return {
    missedBlocksData: missedBlocksQuery.data,
    isLoading: missedBlocksQuery.isLoading,
    isError: missedBlocksQuery.isError,
  };
};

interface DefiData {
    assetPair: string;
    provider: string;
    action: string;
    apy: number;
    tvl: number;
    link: string;
    id: string;
}
export const useDefiData = () => {
  const query = useQuery<DefiData[]>(
    ['defi'],
    async () => {
      const res = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_DATA_API}/defi`);
      if (!res.data || res.data.length === 0) {
        throw new Error('Failed to query defi');
      }
      return res.data;
    },
    {
      staleTime: Infinity,
    }
  );

  return {
    defi: query.data,
    isLoading: query.isLoading,
    isError: query.isError,
  };
};

export const useGovernanceQuery = (chainName: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);
  const governanceQuery = useQuery(
    ['governance', chainName],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }
      const next_key = new Uint8Array()
      const governance = await grpcQueryClient.cosmos.gov.v1beta1.proposals({
        proposalStatus: cosmos.gov.v1.ProposalStatus.PROPOSAL_STATUS_UNSPECIFIED,
        pagination: {
          key: next_key,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          countTotal: true,
          reverse: true,
        },
        voter: '',
        depositor: '',
      });

      return governance;
    },
    {
      enabled: !!grpcQueryClient,
      staleTime: Infinity,
    },
  );

  return {
    governance: governanceQuery.data,
    isLoading: governanceQuery.isLoading,
    isError: governanceQuery.isError,
  };

}

export const useNativeStakeQuery = (chainName: string, address: string) => {
  const { grpcQueryClient } = useGrpcQueryClient(chainName);
  const delegationQuery = useQuery(
    ['delegations', address],
    async () => {
      if (!grpcQueryClient) {
        throw new Error('RPC Client not ready');
      }
      const next_key = new Uint8Array()
      const balance = await grpcQueryClient.cosmos.staking.v1beta1.delegatorDelegations({
        delegator_addr: address || '',
        pagination: {
          key: next_key,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          countTotal: true,
          reverse: false,
        },
      });
  
      return balance;
    },
    {
      enabled: !!grpcQueryClient && !!address,
      staleTime: Infinity,
    },
  );

  return {
    delegations: delegationQuery.data,
    delegationsIsLoading: delegationQuery.isLoading,
    delegationsIsError: delegationQuery.isError,
  };
}

export const useSkipAssets = (chainId: string) => {

  const assetsQuery = useQuery(
    ['assets', chainId],
    async () => {
      const assets = await skipClient.assets({
        chainID: chainId,
        includeEvmAssets: true,
        includeCW20Assets: true,
        includeSvmAssets: true,
      });

      return assets;
    },
    {
      staleTime: Infinity,
    },
  );

  return {
    assets: assetsQuery.data,
    assetsIsLoading: assetsQuery.isLoading,
    assetsIsError: assetsQuery.isError,
  };
};

export const useRecommendations = (source_asset_denom: string, source_asset_chain_id: string, dest_chain_id: string) => {
  const recommendationsQuery = useQuery(
    ['requests', source_asset_chain_id],
    async () => {
      const recommendations = await skipClient.recommendAssets([{
        sourceAssetDenom: source_asset_denom,
        sourceAssetChainID: source_asset_chain_id,
        destChainID: dest_chain_id,
      }]);
      return recommendations;
    },
    {
      staleTime: Infinity,
    },
  );

  return {
    recommendations: recommendationsQuery.data,
    recommendationsIsLoading: recommendationsQuery.isLoading,
    recommendationsIsError: recommendationsQuery.isError,
  };
};