import { saga } from '@chain-registry/assets';
import { HttpEndpoint } from '@cosmjs/stargate';
import { useQuery } from '@tanstack/react-query';
import { cosmos } from 'interchain-query';



const createRPCQueryClient = cosmos.ClientFactory.createRPCQueryClient;

export const useRpcQueryClient = (chainName: string) => {
  let rpcEndpoint: string | HttpEndpoint | undefined;

  const env = process.env.NEXT_PUBLIC_CHAIN_ENV; 
// Builds the query client with the proper endpoint
  const endpoints: { [key: string]: string | undefined } = {
    quicksilver: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_QUICKSILVER : process.env.MAINNET_RPC_ENDPOINT_QUICKSILVER,
    cosmoshub: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_COSMOSHUB : process.env.MAINNET_RPC_ENDPOINT_COSMOSHUB,
    sommelier: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_SOMMELIER : process.env.MAINNET_RPC_ENDPOINT_SOMMELIER,
    stargaze: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_STARGAZE : process.env.MAINNET_RP_ENDPOINTC_STARGAZE,
    regen: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_REGEN : process.env.MAINNET_RPC_ENDPOINT_REGEN,
    osmosis: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_OSMOSIS : process.env.MAINNET_RPC_ENDPOINT_OSMOSIS,
    juno: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_JUNO : process.env.MAINNET_RPC_ENDPOINT_JUNO,
    dydx: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_DYDX : process.env.MAINNET_RPC_ENDPOINT_DYDX,
    saga: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_SAGA : process.env.MAINNET_RPC_ENDPOINT_SAGA,
  };

  rpcEndpoint = endpoints[chainName];

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
