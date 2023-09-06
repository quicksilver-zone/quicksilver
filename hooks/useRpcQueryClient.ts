import { useQuery } from '@tanstack/react-query';
import { cosmos } from 'interchain-query';
import { RPC } from '@/config'

import { useQueryHooks } from './useQueryHooks';

const createRPCQueryClient =
  cosmos.ClientFactory.createRPCQueryClient;

export const useRpcQueryClient = (
  chainName: string,
) => {
  const { rpcEndpoint } =
    useQueryHooks(chainName);

  const rpcQueryClientQuery = useQuery({
    queryKey: ['rpcQueryClient', RPC],
    queryFn: () =>
      createRPCQueryClient({
        rpcEndpoint: RPC || RPC,
      }),
    enabled: !!rpcEndpoint,
    staleTime: Infinity,
  });

  return {
    rpcQueryClient: rpcQueryClientQuery.data,
  };
};
