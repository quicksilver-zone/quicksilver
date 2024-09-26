import { useQuery } from '@tanstack/react-query';
import { cosmos } from 'interchain-query';
import { chains, env } from '@/config';



const createRPCQueryClient = cosmos.ClientFactory.createRPCQueryClient;

export const useRpcQueryClient = (chainName: string) => {
  let rpcEndpoint = chains.get(env)?.get(chainName)?.rpc[0] ?? '';

const rpcQueryClientQuery = useQuery({
    queryKey: ['rpcQueryClient', rpcEndpoint],
    queryFn: () => {
      return createRPCQueryClient({ rpcEndpoint: rpcEndpoint || '' });
    },
    enabled: !!rpcEndpoint,
    staleTime: Infinity,
  });

 

  return {
    rpcQueryClient: rpcQueryClientQuery.data,
  };
};
