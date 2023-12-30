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

  const rpcEnndpoints = {
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

  const lcdEnndpoints = {
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

  return (
    <ChakraProvider theme={defaultTheme}>
      <ChainProvider
        endpointOptions={{
          isLazy: true,
          endpoints: {
            quicksilver: {
              rpc: [rpcEnndpoints.quicksilver ?? ''],
              rest: [lcdEnndpoints.quicksilver ?? ''],
            },
            quicksilvertestnet: {
              rest: ['https://lcd.test.quicksilver.zone/'],
              rpc: ['https://rpc.test.quicksilver.zone'],
            },
            cosmoshub: {
              rpc: [rpcEnndpoints.cosmoshub ?? ''],
              rest: [lcdEnndpoints.cosmoshub ?? ''],
            },

            sommelier: {
              rpc: [rpcEnndpoints.sommelier ?? ''],
              rest: [lcdEnndpoints.sommelier ?? ''],
            },
            stargaze: {
              rpc: [rpcEnndpoints.stargaze ?? ''],
              rest: [lcdEnndpoints.stargaze ?? ''],
            },
            regen: {
              rpc: [rpcEnndpoints.regen ?? ''],
              rest: [lcdEnndpoints.regen ?? ''],
            },
            osmosis: {
              rpc: [rpcEnndpoints.osmosis ?? ''],
              rest: [lcdEnndpoints.osmosis ?? ''],
            },
            osmosistestnet: {
              rpc: [rpcEnndpoints.osmosis ?? ''],
              rest: [lcdEnndpoints.osmosis ?? ''],
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
          <Box
            w="100vw"
            h="100vh"
            bgImage="url('https://s3-alpha-sig.figma.com/quicksilver-app-v2/img/555d/db64/f5bf65e93a15603069e8e865d5f6d60d?Expires=1694995200&Signature=fYfmbqDdOGRYtSeEsOkavPhhkaNQK1UFFfICaUaM1k9OVEpACsoWOcK2upjRW7Tfs-pPTJBuQuvcmF9gBjosh5-Al2xTWHYzDlR~CYJNzsXcseIEnVf7H8lCdJqhZY-T0r~lmbJK5-CmbulWfOaubc-wyY3C-oM3b1RanGV1TqmPZto5bbHwf56jDYqK86HedVMXbUCOlzkeBw2R93AkmNDMOdDbKa9rIKqxil64DuQQAfIFxWm1Rc69Jc1-4K-bunsS~kfz8bSET6TIGmR15nCo~ibfISG72YYKAa7zz6XqUY6GKmmG-Yhj9XyyYb7Jy02r5axNei3DRD78SBe~6w__&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4')"
            bgSize="fit"
            bgPosition="right center"
            bgAttachment="fixed"
            bgRepeat="no-repeat"
          >
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
