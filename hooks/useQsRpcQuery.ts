import { useQuery } from '@tanstack/react-query';
import { quicksilver } from 'quicksilverjs';

import { useQueryHooks } from './useQueryHooks';

const createRPCQueryClient = quicksilver.ClientFactory.createRPCQueryClient;

export const useRpcQueryClient = (chainName: string) => {
  const { rpcEndpoint } = useQueryHooks(chainName);

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
