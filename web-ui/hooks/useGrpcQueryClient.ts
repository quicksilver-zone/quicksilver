import { useQuery } from '@tanstack/react-query';
import { quicksilver } from 'quicksilverjs';
import { chains, env } from '../config';

const createGrpcGateWayClient = quicksilver.ClientFactory.createGrpcGateWayClient;

export const useGrpcQueryClient = (chainName: string) => {

  const grpcEndpoint = chains.get(env)?.get(chainName)?.rest[0];

  const grpcQueryClientQuery = useQuery({
    queryKey: ['grpcQueryClient', grpcEndpoint],
    queryFn: () =>
      createGrpcGateWayClient({
        endpoint: grpcEndpoint?.toString() || '',
      }),
    enabled: !!grpcEndpoint,
    staleTime: Infinity,
  });

  return {
    grpcQueryClient: grpcQueryClientQuery.data,
    grpcQueryClientError: grpcQueryClientQuery.error,
  };
};
