import { HttpEndpoint } from '@cosmjs/stargate';
import { useQuery } from '@tanstack/react-query';
import { quicksilver } from 'quicksilverjs';

import { useQueryHooks } from './useQueryHooks';

const createRPCQueryClient = quicksilver.ClientFactory.createRPCQueryClient;

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
        rpcEndpoint: rpcEndpoint?.toString() || '',
      }),
    enabled: !!rpcEndpoint,
    staleTime: Infinity,
  });

  return {
    rpcQueryClient: rpcQueryClientQuery.data,
  };
};
