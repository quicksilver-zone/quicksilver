import { HttpEndpoint } from '@cosmjs/stargate';
import { useQuery } from '@tanstack/react-query';
import { cosmos } from 'interchain-query';

import { useQueryHooks } from './useQueryHooks';

const createRPCQueryClient = cosmos.ClientFactory.createRPCQueryClient;

export const useRpcQueryClient = (chainName: string) => {
  let rpcEndpoint: string | HttpEndpoint | undefined;
  const solution = useQueryHooks(chainName);

  // Custom logic for setting rpcEndpoint based on the chain name
  if (chainName === 'quicksilver') {
    rpcEndpoint = 'https://rpc.test.quicksilver.zone';
  } else if (chainName === 'cosmoshub') {
    rpcEndpoint = 'https://rpc.sentry-01.theta-testnet.polypore.xyz';
  } else {
    rpcEndpoint = solution.rpcEndpoint;
  }

  const rpcQueryClientQuery = useQuery({
    queryKey: ['rpcQueryClient', rpcEndpoint],
    queryFn: () =>
      createRPCQueryClient({
        rpcEndpoint: rpcEndpoint || '',
      }),
    enabled: !!rpcEndpoint,
    staleTime: Infinity,
  });

  return {
    rpcQueryClient: rpcQueryClientQuery.data,
  };
};
