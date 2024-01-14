import { ChevronDownIcon, ChevronUpIcon, ExternalLinkIcon } from '@chakra-ui/icons';
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
  Image,
  ButtonGroup,
  HStack,
  Link,
  Tooltip,
  Center,
  Spinner,
  useBreakpointValue,
} from '@chakra-ui/react';
import React, { useState } from 'react';

import { useDefiData } from '@/hooks/useQueries';

type ActionButtonTitle = 'Add Liquidity' | 'Borrow' | 'Lend' | 'Mint Stablecoin' | 'Vaults';

const actionTitles: Record<string, ActionButtonTitle> = {
  'add-liquidity': 'Add Liquidity',
  borrow: 'Borrow',
  lend: 'Lend',
  'mint-stablecoin': 'Mint Stablecoin',
  vaults: 'Vaults',
};

interface DefiData {
  assetPair: string;
  apy: number;
  tvl: number;
  provider: string;
  action: string;
  link: string;
}

type SortOrder = 'asc' | 'desc';
type SortableColumn = 'apy' | 'tvl';

const filterCategories: Record<string, (data: DefiData) => boolean> = {
  All: () => true,
  'Borrowing & Lending': (data: DefiData) => data.action === 'Borrow' || data.action === 'Lend',
  Vaults: (data: DefiData) => data.action === 'Vaults',
  'Liquidity Providers': (data: DefiData) => data.action === 'Add Liquidity',
  'Mint Stable Coins': (data: DefiData) => data.action === 'Mint Stablecoin',
};

const formatApy = (apy: number) => {
  return `${(apy * 100).toFixed(2)}%`;
};

const DefiTable = () => {
  const { defi, isLoading, isError } = useDefiData();

  const [activeFilter, setActiveFilter] = useState<string>('All');
  const filterOptions = Object.keys(filterCategories);

  const handleFilterChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setActiveFilter(event.target.value);
  };

  const handleFilterClick = (filter: string) => {
    setActiveFilter(filter);
  };
  const isMobile = useBreakpointValue({ base: true, sm: true, md: false });
  const filteredData = defi ? defi.filter(filterCategories[activeFilter]) : [];

  type ProviderKey = 'osmosis' | 'ux' | 'shade';

  const providerIcons: Record<ProviderKey, string> = {
    osmosis: '/quicksilver/img/osmoIcon.svg',
    ux: '/quicksilver/img/ux.png',
    shade: '/quicksilver/img/shd.svg',
  };

  const isProviderKey = (key: string): key is ProviderKey => {
    return key in providerIcons;
  };

  const [sortColumn, setSortColumn] = useState<SortableColumn | null>(null);
  const [sortOrder, setSortOrder] = useState<SortOrder>('asc');

  const sortData = (data: DefiData[], column: SortableColumn | null, order: SortOrder) => {
    if (!column) return data;
    return [...data].sort((a, b) => {
      let comparison = 0;
      if (column === 'apy' || column === 'tvl') {
        comparison = a[column] - b[column];
      } else {
        comparison = a[column] > b[column] ? 1 : -1;
      }

      return order === 'asc' ? comparison : -comparison;
    });
  };

  const handleSort = (column: SortableColumn) => {
    if (sortColumn === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortColumn(column);
      setSortOrder('asc');
    }
  };

  const sortedData = sortData(filteredData, sortColumn, sortOrder);

  return (
    <Box backdropFilter="blur(50px)" bgColor="rgba(255,255,255,0.1)" flex="1" borderRadius="10px" p={6} rounded="md">
      {isMobile ? (
        <Select
          _active={{
            borderColor: 'complimentary.900',
          }}
          _selected={{
            borderColor: 'complimentary.900',
          }}
          _hover={{
            borderColor: 'complimentary.900',
          }}
          _focus={{
            borderColor: 'complimentary.900',
            boxShadow: '0 0 0 3px #FF8000',
          }}
          color="complimentary.900"
          textAlign={'center'}
          onChange={handleFilterChange}
          value={activeFilter}
          mb={6}
        >
          {filterOptions.map((filter) => (
            <option key={filter} value={filter}>
              {filter}
            </option>
          ))}
        </Select>
      ) : (
        <Stack direction="row" spacing={2} mb={6} justifyContent="space-between">
          {filterOptions.map((filter) => (
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
      )}
      <Box maxH={'480px'} minH={'480px'} overflow={'auto'}>
        <Table color={'white'} variant="simple">
          <Thead position="sticky">
            <Tr>
              <Th color={'complimentary.900'}>Asset Pair</Th>
              <Th
                textAlign={'center'}
                color={'complimentary.900'}
                isNumeric
                onClick={() => handleSort('apy')}
                style={{ cursor: 'pointer' }}
              >
                APY{' '}
                {sortColumn === 'apy' ? (
                  sortOrder === 'asc' ? (
                    <ChevronUpIcon
                      _hover={{ color: 'complimentary.700', transform: 'scale(1.5)' }}
                      _active={{ color: 'complimentary.500' }}
                    />
                  ) : (
                    <ChevronDownIcon
                      _hover={{ color: 'complimentary.700', transform: 'scale(1.5)' }}
                      _active={{ color: 'complimentary.500' }}
                    />
                  )
                ) : (
                  <ChevronDownIcon
                    _hover={{ color: 'complimentary.700', transform: 'scale(1.5)' }}
                    _active={{ color: 'complimentary.500' }}
                  />
                )}
              </Th>
              <Th
                textAlign={'center'}
                color={'complimentary.900'}
                isNumeric
                style={{ cursor: 'pointer' }}
                onClick={() => handleSort('tvl')}
              >
                TVL{' '}
                {sortColumn === 'tvl' ? (
                  sortOrder === 'asc' ? (
                    <ChevronUpIcon
                      _hover={{ color: 'complimentary.700', transform: 'scale(1.5)' }}
                      _active={{ color: 'complimentary.500' }}
                    />
                  ) : (
                    <ChevronDownIcon
                      _hover={{ color: 'complimentary.700', transform: 'scale(1.5)' }}
                      _active={{ color: 'complimentary.500' }}
                    />
                  )
                ) : (
                  <ChevronDownIcon
                    _hover={{ color: 'complimentary.700', transform: 'scale(1.5)' }}
                    _active={{ color: 'complimentary.500' }}
                  />
                )}
              </Th>
              <Th textAlign={'center'} color={'complimentary.900'}>
                Provider
              </Th>
              <Th textAlign={'center'} color={'complimentary.900'}>
                Action
              </Th>
            </Tr>
          </Thead>

          <Tbody>
            {isLoading && !defi && (
              <Tr>
                <Td colSpan={5}>
                  {' '}
                  {/* Span across all columns */}
                  <Center my={42}>
                    <Spinner size="4xl" color="complimentary.900" />
                  </Center>
                </Td>
              </Tr>
            )}
            {defi &&
              sortedData.map((asset, index) => (
                <Tr _even={{ bg: 'rgba(255, 128, 0, 0.1)' }} key={index} borderBottomColor={'transparent'}>
                  <Td textAlign={'center'} borderBottomColor="transparent">
                    <Flex align="center">
                      <Text>{asset.assetPair}</Text>
                    </Flex>
                  </Td>
                  <Td textAlign={'center'} borderBottom="0" borderBottomColor="transparent" isNumeric>
                    {formatApy(asset.apy)}
                  </Td>
                  <Td textAlign={'center'} borderBottomColor="transparent" isNumeric>
                    ${asset.tvl.toLocaleString()}
                  </Td>
                  <Td borderBottomColor="transparent">
                    {isProviderKey(asset.provider.toLowerCase()) && (
                      <Tooltip label={`${asset.provider}`}>
                        <Center>
                          <Image
                            src={providerIcons[asset.provider.toLowerCase() as ProviderKey]}
                            alt={asset.provider}
                            boxSize="2rem"
                            objectFit="cover"
                          />
                        </Center>
                      </Tooltip>
                    )}
                  </Td>
                  <Td textAlign={'center'} borderBottomColor="transparent">
                    <Link href={asset.link} isExternal={true} _hover={{ textDecoration: 'none' }}>
                      <Button
                        _hover={{
                          bgColor: 'complimentary.1000',
                        }}
                        minW="150px"
                        backgroundColor="rgba(255, 128, 0, 0.8)"
                        rightIcon={<ExternalLinkIcon />}
                        variant="ghost"
                      >
                        {actionTitles[asset.action.toLowerCase().replace(/\s+/g, '-') as keyof typeof actionTitles]}
                      </Button>
                    </Link>
                  </Td>
                </Tr>
              ))}
            {defi && filteredData.length === 0 && (
              <Tr>
                <Td colSpan={5}>
                  {' '}
                  {/* Span across all columns */}
                  <Center my={4}>
                    <Text color="complimentary.900">No entries found for this category, please check back later!</Text>
                  </Center>
                </Td>
              </Tr>
            )}
          </Tbody>
        </Table>
      </Box>
    </Box>
  );
};

export default DefiTable;
