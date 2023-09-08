import '../styles/globals.css';
import { Box, ChakraProvider } from '@chakra-ui/react';
import { SignerOptions, WalletViewProps } from '@cosmos-kit/core';
import { wallets as cosmostationWallets } from '@cosmos-kit/cosmostation';
import { wallets as keplrWallets } from '@cosmos-kit/keplr';
import { wallets as leapWallets } from '@cosmos-kit/leap';
import { ChainProvider } from '@cosmos-kit/react';
import {
  QueryClientProvider,
  QueryClient,
} from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { chains, assets } from 'chain-registry';
import type { AppProps } from 'next/app';

import { defaultTheme } from '@/config';
import '@interchain-ui/react/styles';

const ConnectedView = ({
  onClose,
  onReturn,
  wallet,
}: WalletViewProps) => {
  const {
    walletInfo: { prettyName },
    username,
    address,
  } = wallet;
 
  return (
  <Box
  bgColor="complimentary.900"
  >

  </Box>
  );
};

const ConnectingView = ({
  onClose,
  onReturn,
  wallet,
}: WalletViewProps) => {
  const {
    walletInfo: { prettyName },
    username,
    address,
  } = wallet;
 
  return <div>{`${prettyName}/${username}/${address}`}</div>;
};

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 2,
      refetchOnWindowFocus: false,
    },
  },
});

function CreateCosmosApp({
  Component,
  pageProps,
}: AppProps) {
  const signerOptions: SignerOptions = {
    // signingStargate: () => {
    //   return getSigningCosmosClientOptions();
    // }
  };

  return (
    <ChakraProvider theme={defaultTheme}>
      <ChainProvider
        chains={chains}
        assetLists={assets}
        wallets={[
          ...keplrWallets,
          ...cosmostationWallets,
          ...leapWallets,
        ]}
        walletConnectOptions={{
          signClient: {
            projectId:
              'a8510432ebb71e6948cfd6cde54b70f7',
            relayUrl:
              'wss://relay.walletconnect.org',
            metadata: {
              name: 'CosmosKit Template',
              description:
                'CosmosKit dapp template',
              url: 'https://docs.cosmoskit.com/',
              icons: [],
            },
          },
        }}
        signerOptions={signerOptions}
      >
        <QueryClientProvider client={queryClient}>
          <ReactQueryDevtools
            initialIsOpen={true}
          />
          <Component {...pageProps} />
        </QueryClientProvider>
      </ChainProvider>
    </ChakraProvider>
  );
}

export default CreateCosmosApp;
