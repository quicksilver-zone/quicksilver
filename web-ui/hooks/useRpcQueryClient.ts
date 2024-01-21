import { HttpEndpoint } from '@cosmjs/stargate';
import { useQuery } from '@tanstack/react-query';
import { cosmos } from 'interchain-query';



const createRPCQueryClient = cosmos.ClientFactory.createRPCQueryClient;

export const useRpcQueryClient = (chainName: string) => {
  let rpcEndpoint: string | HttpEndpoint | undefined;

  const env = process.env.NEXT_PUBLIC_CHAIN_ENV; 

  const endpoints: { [key: string]: string | undefined } = {
    quicksilver: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_QUICKSILVER : process.env.MAINNET_RPC_ENDPOINT_QUICKSILVER,
    cosmoshub: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_COSMOSHUB : process.env.MAINNET_RPC_ENDPOINT_COSMOSHUB,
    sommelier: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_SOMMELIER : process.env.MAINNET_RPC_ENDPOINT_SOMMELIER,
    stargaze: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_STARGAZE : process.env.MAINNET_RP_ENDPOINTC_STARGAZE,
    regen: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_REGEN : process.env.MAINNET_RPC_ENDPOINT_REGEN,
    osmosis: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_OSMOSIS : process.env.MAINNET_RPC_ENDPOINT_OSMOSIS,
  };

  rpcEndpoint = endpoints[chainName];

  const rpcQueryClientQuery = useQuery({
    queryKey: ['rpcQueryClient', rpcEndpoint],
    queryFn: () => {

      return createRPCQueryClient({ rpcEndpoint: rpcEndpoint || '' });
    },
    enabled: !!rpcEndpoint,
    staleTime: Infinity,
    onError: (error) => {
      console.error('Error in fetching RPC Query Client:', error);
    }
  });

 

  return {
    rpcQueryClient: rpcQueryClientQuery.data,
  };
};
