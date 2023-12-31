import '../styles/globals.css';
import { Chain } from '@chain-registry/types';
import { Box, ChakraProvider, Container, Fade, Flex } from '@chakra-ui/react';
import { Registry } from '@cosmjs/proto-signing';
import { SigningStargateClientOptions, AminoTypes } from '@cosmjs/stargate';
import { SignerOptions, WalletViewProps } from '@cosmos-kit/core';
import { wallets as cosmostationWallets } from '@cosmos-kit/cosmostation';
import { wallets as keplrWallets } from '@cosmos-kit/keplr';
import { wallets as leapWallets } from '@cosmos-kit/leap';
import { ChainProvider } from '@cosmos-kit/react';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { chains, assets } from 'chain-registry';
import { cosmosAminoConverters, cosmosProtoRegistry } from 'interchain-query';
import type { AppProps } from 'next/app';
import { quicksilverProtoRegistry, quicksilverAminoConverters } from 'quicksilverjs';
import { ibcAminoConverters, ibcProtoRegistry } from 'interchain-query';

import { Header, SideHeader } from '@/components';
import { defaultTheme } from '@/config';
import { useRpcQueryClient } from '@/hooks';

import '@interchain-ui/react/styles';

function CreateCosmosApp({ Component, pageProps }: AppProps) {
  const signerOptions: SignerOptions = {
    //@ts-ignore
    signingStargate: (chain: Chain): SigningStargateClientOptions | undefined => {
      const mergedRegistry = new Registry([...cosmosProtoRegistry, ...quicksilverProtoRegistry, ...ibcProtoRegistry]);

      const mergedAminoTypes = new AminoTypes({
        ...cosmosAminoConverters,
        ...quicksilverAminoConverters,
        ...ibcAminoConverters,
      });

      return {
        aminoTypes: mergedAminoTypes,
        //@ts-ignore
        registry: mergedRegistry,
      };
    },
  };

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: 2,
        refetchOnWindowFocus: false,
      },
    },
  });

  const env = process.env.NEXT_PUBLIC_CHAIN_ENV;

  const rpcEndpoints = {
    quicksilver:
      env === 'testnet'
        ? process.env.NEXT_PUBLIC_TESTNET_RPC_ENDPOINT_QUICKSILVER
        : process.env.NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_QUICKSILVER,
    cosmoshub:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_RPC_ENDPOINT_COSMOSHUB : process.env.NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_COSMOSHUB,
    sommelier:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_RPC_ENDPOINT_SOMMELIER : process.env.NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_SOMMELIER,
    stargaze:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_RPC_ENDPOINT_STARGAZE : process.env.NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_STARGAZE,
    regen: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_RPC_ENDPOINT_REGEN : process.env.NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_REGEN,
    osmosis:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_RPC_ENDPOINT_OSMOSIS : process.env.NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_OSMOSIS,
  };

  const lcdEndpoints = {
    quicksilver:
      env === 'testnet'
        ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_QUICKSILVER
        : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_QUICKSILVER,
    cosmoshub:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_COSMOSHUB : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_COSMOSHUB,
    sommelier:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_SOMMELIER : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_SOMMELIER,
    stargaze:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_STARGAZE : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_STARGAZE,
    regen: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_REGEN : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_REGEN,
    osmosis:
      env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_OSMOSIS : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_OSMOSIS,
  };

  return (
    <ChakraProvider theme={defaultTheme}>
      <ChainProvider
        endpointOptions={{
          isLazy: true,
          endpoints: {
            quicksilver: {
              rpc: [rpcEndpoints.quicksilver ?? ''],
              rest: [lcdEndpoints.quicksilver ?? ''],
            },
            quicksilvertestnet: {
              rest: ['https://lcd.test.quicksilver.zone/'],
              rpc: ['https://rpc.test.quicksilver.zone'],
            },
            cosmoshub: {
              rpc: [rpcEndpoints.cosmoshub ?? ''],
              rest: [lcdEndpoints.cosmoshub ?? ''],
            },
            sommelier: {
              rpc: [rpcEndpoints.sommelier ?? ''],
              rest: [lcdEndpoints.sommelier ?? ''],
            },
            stargaze: {
              rpc: [rpcEndpoints.stargaze ?? ''],
              rest: [lcdEndpoints.stargaze ?? ''],
            },
            regen: {
              rpc: [rpcEndpoints.regen ?? ''],
              rest: [lcdEndpoints.regen ?? ''],
            },
            osmosis: {
              rpc: [rpcEndpoints.osmosis ?? ''],
              rest: [lcdEndpoints.osmosis ?? ''],
            },
            osmosistestnet: {
              rpc: [rpcEndpoints.osmosis ?? ''],
              rest: [lcdEndpoints.osmosis ?? ''],
            },
          },
        }}
        chains={chains}
        assetLists={assets}
        //@ts-ignore
        wallets={[...keplrWallets, ...cosmostationWallets, ...leapWallets]}
        walletConnectOptions={{
          signClient: {
            projectId: 'a8510432ebb71e6948cfd6cde54b70f7',
            relayUrl: 'wss://relay.walletconnect.org',
            metadata: {
              name: 'Quicksilver Dashboard',
              description: 'Interact with the Quicksilver Network',
              url: 'https://docs.quicksilver.zone/',
              icons: [],
            },
          },
        }}
        signerOptions={signerOptions}
      >
        <QueryClientProvider client={queryClient}>
          <ReactQueryDevtools initialIsOpen={true} />
          <Box w="100vw" h="100vh" bgSize="fit" bgPosition="right center" bgAttachment="fixed" bgRepeat="no-repeat">
            <Flex justifyContent={'space-between'} alignItems={'center'}>
              <Header chainName="quicksilver" />
              <SideHeader />
            </Flex>
            <Component {...pageProps} />
          </Box>
        </QueryClientProvider>
      </ChainProvider>
    </ChakraProvider>
  );
}

export default CreateCosmosApp;
