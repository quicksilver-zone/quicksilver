import '../styles/globals.css';

import '@interchain-ui/react/styles';
import { Chain } from '@chain-registry/types';
import { Box, Center, ChakraProvider, Image } from '@chakra-ui/react';
import { Registry } from '@cosmjs/proto-signing';
import { SigningStargateClientOptions, AminoTypes, GasPrice } from '@cosmjs/stargate';
import { SignerOptions } from '@cosmos-kit/core';
import { wallets as cosmostationWallets } from '@cosmos-kit/cosmostation';
import { wallets as keplrWallets } from '@cosmos-kit/keplr';
import { wallets as leapWallets } from '@cosmos-kit/leap';
import { ChainProvider, ThemeCustomizationProps } from '@cosmos-kit/react';
import { ThemeProvider, useTheme } from '@interchain-ui/react';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { SpeedInsights } from '@vercel/speed-insights/react';
import { chains, assets } from 'chain-registry';
import { ibcAminoConverters, ibcProtoRegistry } from 'interchain-query';
import type { AppProps } from 'next/app';
import { quicksilverProtoRegistry, quicksilverAminoConverters, cosmosAminoConverters, cosmosProtoRegistry } from 'quicksilverjs';

import { DynamicHeaderSection, SideHeader } from '@/components';
import { defaultTheme, Chain as configChain, chains as configChains, env } from '@/config';
import { LiveZonesProvider } from '@/state/LiveZonesContext';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 2,
      refetchOnWindowFocus: false,
    },
  },
});

function QuickApp({ Component, pageProps }: AppProps) {
  const { themeClass } = useTheme();
  const signerOptions: SignerOptions = {
    signingStargate: (chain: Chain | string): SigningStargateClientOptions | undefined => {
      const mergedRegistry = new Registry([...quicksilverProtoRegistry, ...ibcProtoRegistry, ...cosmosProtoRegistry]);
      const mergedAminoTypes = new AminoTypes({
        ...cosmosAminoConverters,
        ...quicksilverAminoConverters,
        ...ibcAminoConverters,
      });
      switch (true) {
        case chain === 'quicksilver':
        case typeof chain != "string" && chain.chain_id === 'quicksilver-2':
          return {
            aminoTypes: mergedAminoTypes,
            registry: mergedRegistry,
            gasPrice: GasPrice.fromString('0.0025uqck'),
          }
        default:
          return {
            aminoTypes: mergedAminoTypes,
            registry: mergedRegistry,
          }
      };

    },
  };

  const walletConnectToken = process.env.NEXT_PUBLIC_WALLET_CONNECT_TOKEN;

  const modalThemeOverrides: ThemeCustomizationProps = {
    modalContentStyles: {
      backgroundColor: 'rgba(0,0,0,0.75)',
      opacity: 0,
    },
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
        focusedBg: {
          light: 'rgba(0,0,0,0.75)',
          dark: 'rgba(32,32,32,0.9)',
        },
        disabledBg: {
          light: 'rgba(0,0,0,0.75)',
          dark: 'rgba(32,32,32,0.9)',
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
        borderColor: { light: 'black', dark: 'black' },
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
        focusedBorderColor: { light: '#FFFFFF', dark: '#FFFFFF' },
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
          light: '#F0F0F0',
          dark: '#FF8000',
        },
      },
      'connect-modal-qr-code-error': {
        bg: {
          light: '#FFEEEE',
          dark: '#FFFFFF',
        },
      },
      'connect-modal-qr-code-error-button': {
        bg: {
          light: '#FFCCCC',
          dark: '#552222',
        },
      },
    },
  };

  return (
    <LiveZonesProvider>
      <ChakraProvider theme={defaultTheme}>
        <ThemeProvider>
          <ChainProvider
            endpointOptions={{
              isLazy: true,
              endpoints: Array.from(configChains.get(env)?.entries() ?? []).reduce((acc, [chainname, chain]: [string, configChain]) => ({
                ...acc,
                [chainname]: {
                  rpc: chain.rpc,
                  rest: chain.rest,
                },
              }), {}),
            }}
            logLevel="NONE"
            modalTheme={modalThemeOverrides}
            chains={chains}
            assetLists={assets}
            // @ts-ignore
            wallets={[...keplrWallets, ...cosmostationWallets, ...leapWallets]}
            walletConnectOptions={{
              signClient: {
                projectId: walletConnectToken ?? '41a0749c331d209190beeac1c2530c90',
                relayUrl: 'wss://relay.walletconnect.org',
                metadata: {
                  name: 'Quicksilver',
                  description: 'Quicksilver App',
                  url: 'https://app.qucksilver.zone/',
                  icons: [],
                },
              },
            }}
            //@ts-ignore
            signerOptions={signerOptions}
          >
            <QueryClientProvider client={queryClient}>
              <ReactQueryDevtools initialIsOpen={true} />

              <main id="main" className={themeClass}>
                <DynamicHeaderSection chainName="quicksilver" />
                <Box display={{ base: 'none', menu: 'block' }}>
                  <SideHeader />
                </Box>
                <Box w="100vw" h="100vh">
                  <Center>
                    <Component {...pageProps} />
                  </Center>
                  <Image
                    zIndex={5}
                    alt="quick logo"
                    w={'230px'}
                    position={'fixed'}
                    bottom={1}
                    right={4}
                    display={{ base: 'none', xl: 'block' }}
                    src="/img/quicksilverWord.png"
                  />

                  <SpeedInsights />
                </Box>
              </main>
            </QueryClientProvider>
          </ChainProvider>
        </ThemeProvider>
      </ChakraProvider>
    </LiveZonesProvider>
  );
}

export default QuickApp;
