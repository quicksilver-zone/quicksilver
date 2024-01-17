import '../styles/globals.css';
import { Chain } from '@chain-registry/types';
import { Box, ChakraProvider, Flex } from '@chakra-ui/react';
import { ibcAminoConverters, ibcProtoRegistry } from '@chalabi/quicksilverjs';
import { Registry } from '@cosmjs/proto-signing';
import { SigningStargateClientOptions, AminoTypes } from '@cosmjs/stargate';
import { SignerOptions } from '@cosmos-kit/core';
import { wallets as cosmostationWallets } from '@cosmos-kit/cosmostation';
import { wallets as keplrWallets } from '@cosmos-kit/keplr';
import { wallets as leapWallets } from '@cosmos-kit/leap';
import { ChainProvider, ThemeCustomizationProps } from '@cosmos-kit/react';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { chains, assets } from 'chain-registry';
import { cosmosAminoConverters, cosmosProtoRegistry } from 'interchain-query';
import type { AppProps } from 'next/app';
import { quicksilverProtoRegistry, quicksilverAminoConverters } from 'quicksilverjs';
import { cosmosAminoConverters as cosmosAminoConvertersStride, cosmosProtoRegistry as cosmosProtoRegistryStride } from 'stridejs';

import { Header, SideHeader } from '@/components';
import { defaultTheme } from '@/config';

import '@interchain-ui/react/styles';

function QuickApp({ Component, pageProps }: AppProps) {
  const signerOptions: SignerOptions = {
    //@ts-ignore
    signingStargate: (chain: Chain): SigningStargateClientOptions | undefined => {
      //@ts-ignore
      const mergedRegistry = new Registry([
        ...cosmosProtoRegistryStride,
        ...quicksilverProtoRegistry,
        ...ibcProtoRegistry,
        ...cosmosProtoRegistry,
      ]);

      const mergedAminoTypes = new AminoTypes({
        ...cosmosAminoConvertersStride,
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
    juno: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_RPC_ENDPOINT_JUNO : process.env.NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_JUNO,
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
    juno: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_JUNO : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_JUNO,
  };

  const modalThemeOverrides: ThemeCustomizationProps = {
    overrides: {
      'connect-modal': {
        bg: {
          light: 'rgba(0,0,0,0.75)',
          dark: 'rgba(32,32,32,0.9)',
        },
        activeBg: {
          light: 'rgba(0,0,0,0.75)',
          dark: 'rgba(32,32,32,0.9)',
        },
        color: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
      },
      'clipboard-copy-text': {
        bg: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
      },
      'connect-modal-qr-code-shadow': {
        bg: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
      },
      button: {
        bg: {
          light: '#FF8000',
          dark: '#FF8000',
        },
      },
      'connect-modal-head-title': {
        bg: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
      },
      'connect-modal-wallet-button-label': {
        bg: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
      },
      'connect-modal-wallet-button-sublogo': {
        bg: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
      },
      'connect-modal-qr-code-loading': {
        bg: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
      },
      'connect-modal-wallet-button': {
        bg: {
          light: 'rgba(55,55,55,0.9)',
          dark: 'rgba(55,55,55,0.9',
        },
        hoverBg: {
          light: '#FF8000',
          dark: '#FF8000',
        },
        hoverBorderColor: {
          light: 'black',
          dark: 'black',
        },
        activeBorderColor: {
          light: '#FFFFFF',
          dark: '#FFFFFF',
        },
        color: {
          light: '#000000',
          dark: '#FFFFFF',
        },
      },
      'connect-modal-qr-code': {
        bg: {
          light: '',
          dark: 'blue',
        },
        color: {
          light: '#000000',
          dark: '#000000',
        },
      },
      'connect-modal-install-button': {
        bg: {
          light: '#F0F0F0', // Example background color for light theme
          dark: '#FF8000', // Example background color for dark theme
        },
        // Other properties for 'connect-modal-install-button' if needed
      },
      'connect-modal-qr-code-error': {
        bg: {
          light: '#FFEEEE', // Example background color for light theme
          dark: '#FFFFFF', // Example background color for dark theme
        },
        // Other properties for 'connect-modal-qr-code-error' if needed
      },
      'connect-modal-qr-code-error-button': {
        bg: {
          light: '#FFCCCC', // Example background color for light theme
          dark: '#552222', // Example background color for dark theme
        },
      },
    },
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
            umee: {
              rpc: ['https://rpc-umee-ia.cosmosia.notional.ventures/'],
              rest: ['https://api-umee-ia.cosmosia.notional.ventures/'],
            },
          },
        }}
        modalTheme={modalThemeOverrides}
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
        //@ts-ignore
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

export default QuickApp;
