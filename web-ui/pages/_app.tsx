import '../styles/globals.css';

import '@interchain-ui/react/styles';
import { Chain } from '@chain-registry/types';
import { Box, Center, ChakraProvider, Image } from '@chakra-ui/react';
import { Registry } from '@cosmjs/proto-signing';
import { SigningStargateClientOptions, AminoTypes, GasPrice } from '@cosmjs/stargate';
import { SignerOptions } from '@cosmos-kit/core';
import { wallets as cosmostationWallets } from '@cosmos-kit/cosmostation';
import { wallets as keplrWallets } from '@cosmos-kit/keplr-extension';
import { wallets as leapWallets } from '@cosmos-kit/leap-extension';
import { ChainProvider } from '@cosmos-kit/react';
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
        case chain === 'injective':
        case typeof chain != "string" && chain.chain_id === 'injective-1':
          return {
            aminoTypes: mergedAminoTypes,
            registry: mergedRegistry,
            gasPrice: GasPrice.fromString('500000000inj'),
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

  return (
    <LiveZonesProvider>
      <ChakraProvider theme={defaultTheme}>
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
                  url: 'https://app.quicksilver.zone/',
                  icons: [],
                },
              },
            }}
            //@ts-ignore
            signerOptions={signerOptions}
          >
            <QueryClientProvider client={queryClient}>
              <ReactQueryDevtools initialIsOpen={true} />

              <main id="main">
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
      </ChakraProvider>
    </LiveZonesProvider>
  );
}

export default QuickApp;
