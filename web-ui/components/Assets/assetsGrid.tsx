import { Box, VStack, Text, Divider, HStack, Flex, Grid, GridItem, Spinner } from '@chakra-ui/react';
import React from 'react';

import QDepositModal from './modals/qTokenDepositModal';
import QWithdrawModal from './modals/qTokenWithdrawlModal';

import { shiftDigits } from '@/utils';

interface AssetCardProps {
  assetName: string;
  balance: string;
  apy: number;
  nativeAssetName: string;
  isWalletConnected: boolean;
  nonNative: LiquidRewardsData | undefined;
}

interface AssetGridProps {
  isWalletConnected: boolean;
  assets: Array<{
    name: string;
    balance: string;
    apy: number;
    native: string;
  }>;
  nonNative: LiquidRewardsData | undefined;
}

type Amount = {
  denom: string;
  amount: string;
};

type Errors = {
  Errors: any;
};

type LiquidRewardsData = {
  messages: any[];
  assets: {
    [key: string]: [
      {
        Type: string;
        Amount: Amount[];
      },
    ];
  };
  errors: Errors;
};

const AssetCard: React.FC<AssetCardProps> = ({ assetName, balance, apy, nativeAssetName, isWalletConnected, nonNative }) => {
  const calculateTotalBalance = (nonNative: LiquidRewardsData | undefined, nativeAssetName: string) => {
    if (!nonNative) {
      return '0';
    }
    const chainIds = ['osmosis-1', 'secret-1', 'umee-1', 'cosmoshub-4', 'stargaze-1', 'sommelier-3', 'regen-1'];
    let totalAmount = 0;

    chainIds.forEach((chainId) => {
      const assetsInChain = nonNative?.assets[chainId];
      if (assetsInChain) {
        assetsInChain.forEach((asset: any) => {
          const assetAmount = asset.Amount.find((amount: { denom: string }) => amount.denom === `uq${nativeAssetName.toLowerCase()}`);
          if (assetAmount) {
            totalAmount += parseInt(assetAmount.amount, 10); // assuming amount is a string
          }
        });
      }
    });

    return shiftDigits(totalAmount.toString(), -6); // Adjust the shift as per your data's scale
  };

  const nativeAssets = nonNative?.assets['rhye-2']
    ? nonNative.assets['rhye-2'][0].Amount.find((amount) => amount.denom === `uq${nativeAssetName.toLowerCase()}`)
    : undefined;

  const formattedNonNativeBalance = calculateTotalBalance(nonNative, nativeAssetName);

  const formattedNativebalance = nativeAssets ? shiftDigits(nativeAssets.amount, -6) : '0';

  if (!balance || !apy) {
    return (
      <Flex
        w="100%"
        h="100%"
        p={4}
        borderRadius="lg"
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        gap={6}
        color="white"
      >
        <Spinner w={'200px'} h="200px" color="complimentary.900" />
      </Flex>
    );
  }
  return (
    <VStack bg={'rgba(255,255,255,0.1)'} p={4} boxShadow="lg" align="center" spacing={4} borderRadius="lg">
      <VStack w="full" align="center" alignItems={'center'} spacing={3}>
        <HStack w="full" justify="space-between">
          <Text fontWeight="bold" fontSize={'xl'} isTruncated>
            {assetName}
          </Text>
          <HStack>
            <Text fontSize="xs" fontWeight="light" isTruncated>
              APY:
            </Text>
            <Text fontSize="md" fontWeight="bold" isTruncated>
              {shiftDigits(apy.toFixed(2), 2)}%
            </Text>
          </HStack>
        </HStack>
        <Divider bgColor={'complimentary.900'} />
        <Grid mt={4} templateColumns="repeat(2, 1fr)" gap={4} w="full">
          <GridItem>
            <Text fontSize="md" textAlign="left">
              ON QUICKSILVER:
            </Text>
          </GridItem>
          <GridItem>
            <Text fontSize="md" textAlign="right" fontWeight="semibold">
              {formattedNativebalance}
            </Text>
          </GridItem>
          <GridItem>
            <Text fontSize="md" textAlign="left">
              NON-NATIVE:
            </Text>
          </GridItem>
          <GridItem>
            <Text fontSize="md" textAlign="right" fontWeight="semibold">
              {formattedNonNativeBalance}
            </Text>
          </GridItem>
        </Grid>
      </VStack>

      <HStack w="full" pb={4} pt={4} spacing={2}>
        <QDepositModal token={assetName} />
        <QWithdrawModal token={assetName} />
      </HStack>
    </VStack>
  );
};

const AssetsGrid: React.FC<AssetGridProps> = ({ assets, isWalletConnected, nonNative }) => {
  return (
    <>
      <Text fontSize="xl" fontWeight="bold" color="white" mb={4}>
        qAssets
      </Text>
      {!isWalletConnected && (
        <Flex
          backdropFilter="blur(50px)"
          bgColor="rgba(255,255,255,0.1)"
          direction="column"
          p={5}
          borderRadius="lg"
          align="center"
          justify="space-around"
          w="full"
          h="200px"
        >
          <Text fontSize="xl" textAlign="center">
            Wallet is not connected! Please connect your wallet to interact with your qAssets.
          </Text>
        </Flex>
      )}
      {isWalletConnected && (
        <Box overflowX="auto" w="full">
          <Flex gap="8">
            {assets.map((asset, index) => (
              <Box key={index} minW="350px">
                <AssetCard
                  isWalletConnected={isWalletConnected}
                  assetName={asset.name}
                  nativeAssetName={asset.native}
                  balance={asset.balance}
                  apy={asset.apy}
                  nonNative={nonNative}
                />
              </Box>
            ))}
          </Flex>
        </Box>
      )}
    </>
  );
};
export default AssetsGrid;
