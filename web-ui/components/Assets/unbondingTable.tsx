import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import { Table, Thead, Tbody, Tr, Th, Td, TableContainer, Text, Box, Flex, IconButton, Spinner } from '@chakra-ui/react';
import { useState } from 'react';

import { Chain, chains, env } from '@/config';
import { useUnbondingQuery } from '@/hooks/useQueries';
import { shiftDigits, formatQasset } from '@/utils';

const statusCodes = new Map<number, string>([
  [2, 'QUEUED'],
  [3, 'UNBONDING'],
  [4, 'SENDING'],
  [5, 'COMPLETED'],
]);

const formatDateAndTime = (dateString: string | number | Date) => {
  const date = new Date(dateString);
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',

    hour12: false,
  };
  return date.toLocaleString(undefined, options);
};

const formatDenom = (denom: string) => {
  return formatQasset(denom.substring(1).toUpperCase());
};

interface UnbondingAssetsTableProps {
  address: string;
  isWalletConnected: boolean;
}

const UnbondingAssetsTable: React.FC<UnbondingAssetsTableProps> = ({ address, isWalletConnected }) => {

  const networks: Map<string, Chain> = chains.get(env) ?? new Map();
  const chain_list = Array.from(networks).filter(([_, network]) => network.show).map(([key, _]) => key);

  const [currentChainName, setCurrentChainName] = useState(chain_list[0]);
  const currentNetwork = networks?.get(currentChainName)

  const prev = () => {
    const index = chain_list?.indexOf(currentChainName) - 1
    return index < 0 ? chain_list.length - 1 : index
  }

  const next = () => {
    const index = chain_list?.indexOf(currentChainName) + 1
    return index > chain_list.length - 1 ? 0 : index
  }

  const { unbondingData, isLoading } = useUnbondingQuery(currentChainName, address);

  // Handlers for chain slider
  const handleLeftArrowClick = () => {
    setCurrentChainName(chain_list[prev()]);
  };

  const handleRightArrowClick = () => {
    setCurrentChainName(chain_list[next()]);
  };

  const hideOnMobile = {
    base: 'none',
    md: 'table-cell',
  };

  const noUnbondingAssets = isWalletConnected && unbondingData?.withdrawals.length === 0;
  if (!isWalletConnected) {
    return (
      <Flex direction="column" gap={4}>
        <Flex justifyContent="space-between" alignItems="center">
          <Text fontSize="xl" fontWeight="bold" color="white">
            Unbonding Assets
          </Text>
          <Flex alignItems="center" gap="2">
            <IconButton
              icon={<ChevronLeftIcon />}
              onClick={handleLeftArrowClick}
              aria-label="Previous chain"
              variant="ghost"
              _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
              _active={{
                transform: 'scale(0.75)',
                color: 'complimentary.800',
              }}
              color="white"
            />
            <Box minWidth="100px" textAlign="center">
              <Text>{currentNetwork?.pretty_name}</Text>
            </Box>
            <IconButton
              icon={<ChevronRightIcon />}
              onClick={handleRightArrowClick}
              aria-label="Next chain"
              variant="ghost"
              _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
              _active={{
                transform: 'scale(0.75)',
                color: 'complimentary.800',
              }}
              color="white"
            />
          </Flex>
        </Flex>
        <Flex
          w="100%"
          backdropFilter="blur(50px)"
          bgColor="rgba(255,255,255,0.1)"
          h="sm"
          p={4}
          borderRadius="lg"
          flexDirection="column"
          justifyContent="center"
          alignItems="center"
          gap={6}
          color="white"
        >
          <Flex justifyContent="center" alignItems="center" h="200px">
            <Text fontSize="xl" textAlign="center">
              Wallet is not connected! Please connect your wallet to view your unbonding assets.
            </Text>
          </Flex>
        </Flex>
      </Flex>
    );
  }
  if (isLoading) {
    return (
      <Flex direction="column" gap={4}>
        <Flex justifyContent="space-between" alignItems="center">
          <Text fontSize="xl" fontWeight="bold" color="white">
            Unbonding Assets
          </Text>
          <Flex alignItems="center" gap="2">
            <IconButton
              icon={<ChevronLeftIcon />}
              onClick={handleLeftArrowClick}
              aria-label="Previous chain"
              variant="ghost"
              _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
              _active={{
                transform: 'scale(0.75)',
                color: 'complimentary.800',
              }}
              color="white"
            />
            <Box minWidth="100px" textAlign="center">
              <Text>{currentNetwork?.pretty_name}</Text>
            </Box>
            <IconButton
              icon={<ChevronRightIcon />}
              onClick={handleRightArrowClick}
              aria-label="Next chain"
              variant="ghost"
              _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
              _active={{
                transform: 'scale(0.75)',
                color: 'complimentary.800',
              }}
              color="white"
            />
          </Flex>
        </Flex>
        <Flex
          w="100%"
          backdropFilter="blur(50px)"
          bgColor="rgba(255,255,255,0.1)"
          h="sm"
          p={4}
          borderRadius="lg"
          flexDirection="column"
          justifyContent="center"
          alignItems="center"
          gap={6}
          color="white"
        >
          <Flex justifyContent="center" alignItems="center" h="200px">
            <Spinner size="xl" color="complimentary.700" />
          </Flex>
        </Flex>
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap={4}>
      <Flex justifyContent="space-between" alignItems="center">
        <Text fontSize="xl" fontWeight="bold" color="white">
          Unbonding Assets
        </Text>
        <Flex alignItems="center" gap="2">
          <IconButton
            icon={<ChevronLeftIcon />}
            onClick={handleLeftArrowClick}
            aria-label="Previous chain"
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            color="white"
          />
          <Box minWidth="100px" textAlign="center">
            <Text>{currentNetwork?.pretty_name}</Text>
          </Box>
          <IconButton
            icon={<ChevronRightIcon />}
            onClick={handleRightArrowClick}
            aria-label="Next chain"
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            color="white"
          />
        </Flex>
      </Flex>
      <Flex
        w="100%"
        backdropFilter="blur(50px)"
        bgColor="rgba(255,255,255,0.1)"
        h="sm"
        borderRadius="lg"
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        gap={6}
        color="white"
      >
        {!isWalletConnected ? (
          <Text fontSize="xl" textAlign="center">
            Wallet is not connected! Please connect your wallet to view your unbonding assets.
          </Text>
        ) : noUnbondingAssets ? (
          <Text fontSize="xl" textAlign="center">
            You have no unbonding assets.
          </Text>
        ) : (
          <Box w="100%" h="100%" p={4} borderRadius="lg">
            <TableContainer h={'200px'} overflowY={'auto'}>
              <Table variant="simple" color="white">
                <Thead boxShadow="0px 0.5px 0px 0px rgba(255,255,255,1)" position={'sticky'} bgColor="#1A1A1A" top="0" zIndex="sticky">
                  <Tr>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.700">
                      Burn Amount
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.700" display={hideOnMobile}>
                      Status
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.700">
                      Redemption Amount
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.700" display={hideOnMobile}>
                      Epoch Number
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.700" display={hideOnMobile}>
                      Completion Time
                    </Th>
                  </Tr>
                </Thead>
                <Tbody>
                  {unbondingData?.withdrawals.map((withdrawal, index) => {
                    
                    const exp = chains.get(env)?.get(currentChainName)?.exponent ?? 6
                    return (
                      <Tr _even={{ bg: 'rgba(255, 128, 0, 0.1)' }} key={index}>
                        <Td textAlign="center" borderBottomColor={'transparent'}>
                          {Number(shiftDigits(withdrawal.burn_amount.amount, -exp))} {formatDenom(withdrawal.burn_amount.denom)}
                        </Td>
                        <Td textAlign="center" borderBottomColor={'transparent'} display={hideOnMobile}>
                          {statusCodes.get(withdrawal.status)}
                        </Td>
                        <Td textAlign="center" borderBottomColor={'transparent'}>
                          {withdrawal.amount.map((amt) => `${shiftDigits(amt.amount, -exp)} ${formatDenom(amt.denom)}`).join(', ')}
                        </Td>
                        <Td textAlign="center" borderBottomColor={'transparent'} display={hideOnMobile}>
                          {withdrawal.epoch_number}
                        </Td>
                        <Td textAlign="center" borderBottomColor={'transparent'} display={hideOnMobile}>
                          {withdrawal.status === 2
                            ? 'Pending'
                            : withdrawal.status === 4
                              ? 'A few moments'
                              : formatDateAndTime(withdrawal.completion_time)}
                        </Td>
                      </Tr>
                    );
                  })}
                </Tbody>
              </Table>
            </TableContainer>
          </Box>
        )}
      </Flex>
    </Flex>
  );
};

export default UnbondingAssetsTable;
