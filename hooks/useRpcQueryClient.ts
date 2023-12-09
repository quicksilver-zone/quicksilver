import { HttpEndpoint } from '@cosmjs/stargate';
import { useQuery } from '@tanstack/react-query';
import { cosmos } from 'interchain-query';

import { useQueryHooks } from './useQueryHooks';

const createRPCQueryClient = cosmos.ClientFactory.createRPCQueryClient;

export const useRpcQueryClient = (chainName: string) => {
  let rpcEndpoint: string | HttpEndpoint | undefined;
  const solution = useQueryHooks(chainName);



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
