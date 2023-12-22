import { shiftDigits } from '@/utils';
import { Box, SimpleGrid, VStack, Text, Button, Divider, useColorModeValue, HStack, Flex, Grid, GridItem } from '@chakra-ui/react';
import React from 'react';
import QDepositModal from './modals/qTokenDepositModal';
import QWithdrawModal from './modals/qTokenWithdrawlModal';
interface AssetCardProps {
  assetName: string;
  balance: string;
  apy: number;
  nativeAssetName: string;
}

interface AssetGridProps {
  assets: Array<{
    name: string;
    balance: string;
    apy: number;
    native: string;
  }>;
}

const AssetCard: React.FC<AssetCardProps> = ({ assetName, balance, apy, nativeAssetName }) => {
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
              {balance}
            </Text>
          </GridItem>
          <GridItem>
            <Text fontSize="md" textAlign="left">
              NON-NATIVE:
            </Text>
          </GridItem>
          <GridItem>
            <Text fontSize="md" textAlign="right" fontWeight="semibold">
              {balance}
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

const AssetsGrid: React.FC<AssetGridProps> = ({ assets }) => {
  return (
    <>
      <Text fontSize="xl" fontWeight="bold" color="white" mb={4}>
        qAssets
      </Text>
      <Box overflowX="auto" w="full">
        <Flex gap="8">
          {assets.map((asset, index) => (
            <Box key={index} minW="350px">
              {' '}
              <AssetCard assetName={asset.name} nativeAssetName={asset.native} balance={asset.balance} apy={asset.apy} />
            </Box>
          ))}
        </Flex>
      </Box>
    </>
  );
};
export default AssetsGrid;
