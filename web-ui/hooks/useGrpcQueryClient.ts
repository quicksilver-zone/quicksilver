import { quicksilver } from 'quicksilverjs';
import { useQuery } from '@tanstack/react-query';


const createGrpcGateWayClient = quicksilver.ClientFactory.createGrpcGateWayClient;

export const useGrpcQueryClient = (chainName: string) => {


  let grpcEndpoint: string | undefined;
  const env = process.env.NEXT_PUBLIC_CHAIN_ENV; 



  const endpoints: { [key: string]: string | undefined } = {
    quicksilver: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_QUICKSILVER : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_QUICKSILVER,
    cosmoshub: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_COSMOSHUB : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_COSMOSHUB,
    sommelier: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_SOMMELIER : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_SOMMELIER,
    stargaze: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_STARGAZE : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_STARGAZE,
    regen: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_REGEN : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_REGEN,
    osmosis: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_OSMOSIS : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_OSMOSIS,
    juno: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_JUNO : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_JUNO,
  };


  grpcEndpoint = endpoints[chainName];



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
