import { Box, SimpleGrid, VStack, Text, Button, Divider, useColorModeValue, HStack, Flex, Grid, GridItem } from '@chakra-ui/react';
import React from 'react';

interface AssetCardProps {
  assetName: string;
  balance: string;
  apy: string;
  nativeAssetName: string;
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
              {apy}
            </Text>
          </HStack>
        </HStack>
        <Divider />
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

      <HStack borderBottom="1px" borderBottomColor="complimentary.900" w="full" pb={4} pt={4} spacing={2}>
        <Button color="white" flex={1} size="sm" variant="outline">
          Deposit
        </Button>
        <Button color="white" flex={1} size="sm" variant="outline">
          Withdraw
        </Button>
      </HStack>
      <HStack w="full" justify="space-between">
        <VStack>
          <Text fontWeight="bold" fontSize={'xl'} isTruncated>
            {nativeAssetName}
          </Text>
          <Divider />
        </VStack>
        <VStack>
          <Text fontSize="md" fontWeight="bold" isTruncated></Text>
          <Text fontSize="xs" fontWeight="light" isTruncated></Text>
        </VStack>
      </HStack>
      <Grid templateColumns="repeat(2, 1fr)" gap={4} w="full">
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
      </Grid>
      <HStack w="full" pb={4} pt={4} spacing={2}>
        <Button colorScheme="teal" flex={1} size="sm" variant="solid">
          Deposit
        </Button>
        <Button colorScheme="red" flex={1} size="sm" variant="solid">
          Withdraw
        </Button>
      </HStack>
    </VStack>
  );
};

const AssetsGrid = () => {
  const assets = [
    { name: 'qATOM', balance: '0.123', apy: '12.34%', native: 'ATOM' },
    { name: 'qREGEN', balance: '0.123', apy: '12.34%', native: 'REGEN' },
    { name: 'qOSMO', balance: '0.123', apy: '12.34%', native: 'OSMO' },
    { name: 'qSTARS', balance: '0.123', apy: '12.34%', native: 'STARS' },
    { name: 'qSOMM', balance: '0.123', apy: '12.34%', native: 'SOMM' },
  ];

  return (
    <>
      <Text fontSize="xl" fontWeight="bold" color="white" mb={4}>
        Assets (qAssets + Native Balance)
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
