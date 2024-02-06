import { useChain } from '@cosmos-kit/react';
import { useQueries, useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { cosmos } from 'interchain-query';
import { quicksilver } from 'quicksilverjs';
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

type UseWithdrawalsQueryReturnType = {
  data: WithdrawalsResponse | undefined;
  isLoading: boolean;
  isError: boolean;
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
};

interface ProofOp {
  type: string;
  key: Uint8Array;  // Updated to Uint8Array
  data: Uint8Array; // Updated to Uint8Array
}

interface Proof {
  key: Uint8Array;  // Updated to Uint8Array
  data: Uint8Array; // Updated to Uint8Array
  proof_ops: {
    ops: ProofOp[];
  };
  height: Long; // Assuming height is a number
  proof_type: string;
}

interface Message {
  user_address: string;
  zone: string;
  src_zone: string;
  claim_type: number;
  proofs: Proof[];
  // Remove height and proof_type if they are not needed here
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

// Type for the useLiquidEpochQuery return
interface UseLiquidEpochQueryReturnType {
  liquidEpoch: LiquidEpochData | undefined;
  isLoading: boolean;
  isError: boolean;
}



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

export const useAuthChecker = (address: string) => {
  const authQuery = useQuery(
    ['auth', address],
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
      const nextKey = new Uint8Array()
      const balance = await grpcQueryClient.cosmos.bank.v1beta1.allBalances({
        address: address || '',
        pagination: {
          key: nextKey,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          count_total: true,
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
      const nextKey = new Uint8Array()
      const balance = await grpcQueryClient.cosmos.bank.v1beta1.allBalances({
        address: address || '',
        pagination: {
          key: nextKey,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          count_total: true,
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

      const balance = await grpcQueryClient.cosmos.bank.v1beta1.balance({
        address: address || '',
        denom: 'uq' + qAsset,
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
  };
};

export const useIntentQuery = (chainName: string, address: string) => {
  const { grpcQueryClient } = useGrpcQueryClient('quicksilver');
  const { chain } = useChain(chainName);
  const env = process.env.NEXT_PUBLIC_CHAIN_ENV;
  const baseApiUrl = env === 'testnet' ? 'https://lcd.test.quicksilver.zone' : 'https://lcd.quicksilver.zone';
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
    ['intent', chainName],
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

      const response = await axios.get<LiquidRewardsData>(`https://claim.test.quicksilver.zone/${address}/current`);
      return response.data;
    },
    {
      enabled:!!address,
      staleTime: Infinity,
    },
  );

  return {
    liquidRewards: liquidRewardsQuery.data,
    isLoading: liquidRewardsQuery.isLoading,
    isError: liquidRewardsQuery.isError,
  };

}

export const useLiquidEpochQuery = (address: string): UseLiquidEpochQueryReturnType => {
  const liquidEpochQuery = useQuery(
    ['liquidEpoch', address],
    async () => {
      if (!address) {
        throw new Error('Address is not available');
      }

      const response = await axios.get<LiquidEpochData>(`https://claim.test.quicksilver.zone/${address}/epoch`);


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
        count_total: true,
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
      let nextKey = new Uint8Array();

      do {
        const response = await fetchValidators(nextKey);
        allValidators = allValidators.concat(response.validators);
        nextKey = response.pagination.next_key ?? new Uint8Array();
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

function parseZone(apiZone: any): Zone {

  return {
    connection_id: apiZone.connection_id,
    chain_id: apiZone.chain_id,
    deposit_address: apiZone.deposit_address,
    withdrawal_address: apiZone.withdrawal_address,
    performance_address: apiZone.performance_address,
    delegation_address: apiZone.delegation_address,
    account_prefix: apiZone.account_prefix,
    local_denom: apiZone.local_denom,
    base_denom: apiZone.base_denom,
    redemption_rate: apiZone.redemption_rate,
    last_redemption_rate: apiZone.last_redemption_rate,
    validators: apiZone.validators,
    aggregate_intent: apiZone.aggregate_intent,
    multi_send: apiZone.multi_send,
    liquidity_module: apiZone.liquidity_module,
    withdrawal_waitgroup: apiZone.withdrawal_waitgroup,
    ibc_next_validators_hash: apiZone.ibc_next_validators_hash,
    validator_selection_allocation: apiZone.validator_selection_allocation,
    holdings_allocation: apiZone.holdings_allocation,
    last_epoch_height: apiZone.last_epoch_height,
    tvl: apiZone.tvl,
    unbonding_period: apiZone.unbonding_period,
    messages_per_tx: apiZone.messages_per_tx,
    decimals: apiZone.decimals,
    return_to_sender: apiZone.return_to_sender,
    unbonding_enabled: apiZone.unbonding_enabled,
    deposits_enabled: apiZone.deposits_enabled,
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
        connection_id: apiZone.connection_id,
        chain_id: apiZone.chain_id,
        deposit_address: apiZone.deposit_address,
        withdrawal_address: apiZone.withdrawal_address,
        performance_address: apiZone.performance_address,
        delegation_address: apiZone.delegation_address,
        account_prefix: apiZone.account_prefix,
        local_denom: apiZone.local_denom,
        base_denom: apiZone.base_denom,
        redemption_rate: apiZone.redemption_rate,
        last_redemption_rate: apiZone.last_redemption_rate,
        validators: apiZone.validators,
        aggregate_intent: apiZone.aggregate_intent,
        multi_send: apiZone.multi_send,
        liquidity_module: apiZone.liquidity_module,
        withdrawal_waitgroup: apiZone.withdrawal_waitgroup,
        ibc_next_validators_hash: apiZone.ibc_next_validators_hash,
        validator_selection_allocation: apiZone.validator_selection_allocation,
        holdings_allocation: apiZone.holdings_allocation,
        last_epoch_height: apiZone.last_epoch_height,
        tvl: apiZone.tvl,
        unbonding_period: apiZone.unbonding_period,
        messages_per_tx: apiZone.messages_per_tx,
        decimals: apiZone.decimals,
        return_to_sender: apiZone.return_to_sender,
        unbonding_enabled: apiZone.unbonding_enabled,
        deposits_enabled: apiZone.deposits_enabled,
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
    let nextKey = new Uint8Array();
  
    do {
      const response = await grpcQueryClient.cosmos.slashing.v1beta1.signingInfos({
        pagination: {
          key: nextKey,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          count_total: true,
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
      nextKey = response.pagination?.next_key ?? new Uint8Array();
    } while (nextKey && nextKey.length > 0);
  
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
      const nextKey = new Uint8Array()
      const governance = await grpcQueryClient.cosmos.gov.v1beta1.proposals({
        proposal_status: cosmos.gov.v1.ProposalStatus.PROPOSAL_STATUS_UNSPECIFIED,
        pagination: {
          key: nextKey,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          count_total: true,
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
      const nextKey = new Uint8Array()
      const balance = await grpcQueryClient.cosmos.staking.v1beta1.delegatorDelegations({
        delegator_addr: address || '',
        pagination: {
          key: nextKey,
          offset: Long.fromNumber(0),
          limit: Long.fromNumber(100),
          count_total: true,
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