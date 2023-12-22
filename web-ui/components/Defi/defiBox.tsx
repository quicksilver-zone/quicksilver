import React, { useState } from 'react';
import {
  Box,
  Button,
  Flex,
  Heading,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Text,
  Select,
  Stack,
  useColorModeValue,
  ButtonGroup,
  HStack,
} from '@chakra-ui/react';
import { ChevronDownIcon, ExternalLinkIcon } from '@chakra-ui/icons';
type ActionButtonTitle = 'Add Liquidity' | 'Borrow' | 'Lend' | 'Mint Stablecoin' | 'Vaults';
interface DefiAsset {
  id: string;
  assetPair: string;
  apy: number;
  tvl: string;
  provider: string;
  action: string;
}

const fakeData: DefiAsset[] = [
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 0.56,
    tvl: '$416.87',
    provider: 'Radiyum',
    action: 'Add Liquidity',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 10,
    tvl: '$20006.87',
    provider: 'Osmosis',
    action: 'Mint Stablecoin',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 10,
    tvl: '$20006.87',
    provider: 'Osmosis',
    action: 'Lend',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 0.56,
    tvl: '$416.87',
    provider: 'Radiyum',
    action: 'Borrow',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 0.56,
    tvl: '$416.87',
    provider: 'Radiyum',
    action: 'add-liquidity',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 0.56,
    tvl: '$416.87',
    provider: 'Radiyum',
    action: 'add-liquidity',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 0.56,
    tvl: '$416.87',
    provider: 'Radiyum',
    action: 'Vaults',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 0.56,
    tvl: '$416.87',
    provider: 'Radiyum',
    action: 'add-liquidity',
  },
  {
    id: '1',
    assetPair: 'qATOM - ATOM',
    apy: 0.56,
    tvl: '$416.87',
    provider: 'Radiyum',
    action: 'add-liquidity',
  },
];

const actionTitles: Record<string, ActionButtonTitle> = {
  'add-liquidity': 'Add Liquidity',
  borrow: 'Borrow',
  lend: 'Lend',
  'mint-stablecoin': 'Mint Stablecoin',
  vaults: 'Vaults',
};

const filterCategories: Record<string, (asset: DefiAsset) => boolean> = {
  All: () => true,
  'Borrowing & Lending': (asset: DefiAsset) => asset.action === 'Borrow' || asset.action === 'Lend',
  Vaults: (asset: DefiAsset) => asset.action === 'Vaults',
  'Liquidity Providers': (asset: DefiAsset) => asset.action === 'Add Liquidity',
  'Mint Stable Coins': (asset: DefiAsset) => asset.action === 'Mint Stablecoin',
};

const DefiTable = () => {
  const [activeFilter, setActiveFilter] = useState<string>('All');

  const handleFilterClick = (filter: string) => {
    setActiveFilter(filter);
  };

  const filteredData = fakeData.filter(filterCategories[activeFilter]);

  return (
    <Box backdropFilter="blur(50px)" bgColor="rgba(255,255,255,0.1)" flex="1" borderRadius="10px" p={6} rounded="md">
      <Stack direction="row" spacing={2} mb={6} justifyContent="space-between">
        {Object.keys(filterCategories).map((filter) => (
          <Button
            key={filter}
            onClick={() => handleFilterClick(filter)}
            isActive={activeFilter === filter}
            _active={{
              transform: 'scale(0.95)',
            }}
            _hover={{
              bgColor: 'rgba(255,128,0, 0.25)',
              color: 'complimentary.300',
            }}
            color="white"
            minW={'180px'}
            colorScheme={activeFilter === filter ? 'orange' : 'gray'}
            variant={activeFilter === filter ? 'solid' : 'outline'}
          >
            {filter}
          </Button>
        ))}
      </Stack>
      <Box maxH={'480px'} overflow={'auto'}>
        <Table color={'white'} variant="simple">
          <Thead>
            <Tr>
              <Th color={'complimentary.900'}>
                Asset Pair <ChevronDownIcon />
              </Th>
              <Th color={'complimentary.900'} isNumeric>
                APY <ChevronDownIcon />
              </Th>
              <Th color={'complimentary.900'} isNumeric>
                TVL <ChevronDownIcon />
              </Th>
              <Th color={'complimentary.900'}>Provider</Th>
              <Th color={'complimentary.900'}>Action</Th>
            </Tr>
          </Thead>
          <Tbody>
            {filteredData.map((asset, index) => (
              <Tr _even={{ bg: 'rgba(255, 128, 0, 0.1)' }} key={asset.id} borderBottomColor={'transparent'}>
                <Td borderBottomColor="transparent">
                  <Flex align="center">
                    <Box w="2rem" h="2rem" bg="gray.200" rounded="full" mr={2}></Box>
                    <Text>{asset.assetPair}</Text>
                  </Flex>
                </Td>
                <Td borderBottom="0" borderBottomColor="transparent" isNumeric>
                  {asset.apy}%
                </Td>
                <Td borderBottomColor="transparent" isNumeric>
                  {asset.tvl}
                </Td>
                <Td borderBottomColor="transparent">{asset.provider}</Td>
                <Td borderBottomColor="transparent">
                  <Button backgroundColor="rgba(255, 128, 0, 0.8)" rightIcon={<ExternalLinkIcon />} variant="ghost">
                    {actionTitles[asset.action.toLowerCase().replace(/\s+/g, '-') as keyof typeof actionTitles]}
                  </Button>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </Box>
    </Box>
  );
};

export default DefiTable;
