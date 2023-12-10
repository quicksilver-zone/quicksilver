import { HttpEndpoint } from '@cosmjs/stargate';
import { quicksilver } from '@hoangdv2429/quicksilverjs';
import { useQuery } from '@tanstack/react-query';

import { useQueryHooks } from './useQueryHooks';

const createGrpcGateWayClient = quicksilver.ClientFactory.createGrpcGateWayClient;

export const useGrpcQueryClient = (chainName: string) => {
  let grpcEndpoint: string | HttpEndpoint | undefined;
  const solution = useQueryHooks(chainName);

  // Custom logic for setting rpcEndpoint based on the chain name
  if (chainName === 'quicksilver') {
    grpcEndpoint = 'https://lcd.quicksilver.zone';
  } else if (chainName === 'cosmoshub') {
    grpcEndpoint = 'https://rest.sentry-01.theta-testnet.polypore.xyz';
  } else {
    grpcEndpoint = solution.rpcEndpoint;
  }

  grpcEndpoint = solution.rpcEndpoint;

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
  };
};
