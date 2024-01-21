import '../styles/globals.css';
import { Chain } from '@chain-registry/types';
import { Box, ChakraProvider, Flex } from '@chakra-ui/react';
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
import type { AppProps } from 'next/app';
import {
  quicksilverProtoRegistry,
  quicksilverAminoConverters,
  cosmosAminoConverters,
  cosmosProtoRegistry,
  ibcAminoConverters,
  ibcProtoRegistry,
} from 'quicksilverjs';

import { Header, SideHeader } from '@/components';
import { defaultTheme } from '@/config';

import '@interchain-ui/react/styles';

function QuickApp({ Component, pageProps }: AppProps) {
  const signerOptions: SignerOptions = {
    //@ts-ignore
    signingStargate: (chain: Chain): SigningStargateClientOptions | undefined => {
      //@ts-ignore
      const mergedRegistry = new Registry([...quicksilverProtoRegistry, ...ibcProtoRegistry, ...cosmosProtoRegistry]);

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

  function isWalletClientAvailable(walletName: string) {
    if (typeof window === 'undefined') {
      return false;
    }

    switch (walletName) {
      case 'Keplr':
        return typeof window.keplr !== 'undefined';
      case 'Cosmostation':
        return typeof window.cosmostation !== 'undefined';
      case 'Leap':
        return typeof window.leap !== 'undefined';
      default:
        return false;
    }
  }

  const availableWallets = [];

  if (isWalletClientAvailable('Keplr')) {
    availableWallets.push(...keplrWallets);
  }
  if (isWalletClientAvailable('Cosmostation')) {
    availableWallets.push(...cosmostationWallets);
  }
  if (isWalletClientAvailable('Leap')) {
    availableWallets.push(...leapWallets);
  }

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
        wallets={availableWallets}
        walletConnectOptions={{
          signClient: {
            projectId: '41a0749c331d209190beeac1c2530c90',
            relayUrl: 'wss://relay.walletconnect.org',
            metadata: {
              name: 'Quicksilver',
              description: 'CosmosKit dapp template',
              url: 'https://docs.cosmoskit.com/',
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
