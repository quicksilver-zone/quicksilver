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

import { Header, SideHeader } from '@/components';
import { defaultTheme } from '@/config';

import '@interchain-ui/react/styles';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 2,
      refetchOnWindowFocus: false,
    },
  },
});

function CreateCosmosApp({ Component, pageProps }: AppProps) {
  const signerOptions: SignerOptions = {
    signingStargate: (chain: Chain): SigningStargateClientOptions | undefined => {
      const mergedRegistry = new Registry([...cosmosProtoRegistry, ...quicksilverProtoRegistry]);

      const mergedAminoTypes = new AminoTypes({
        ...cosmosAminoConverters,
        ...quicksilverAminoConverters,
      });

      return {
        aminoTypes: mergedAminoTypes,
        registry: mergedRegistry,
      };
    },
  };

  return (
    <ChakraProvider theme={defaultTheme}>
      <ChainProvider
        chains={chains}
        assetLists={assets}
        wallets={[...keplrWallets, ...cosmostationWallets, ...leapWallets]}
        walletConnectOptions={{
          signClient: {
            projectId: 'a8510432ebb71e6948cfd6cde54b70f7',
            relayUrl: 'wss://relay.walletconnect.org',
            metadata: {
              name: 'CosmosKit Template',
              description: 'CosmosKit dapp template',
              url: 'https://docs.cosmoskit.com/',
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
