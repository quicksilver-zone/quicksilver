import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import { Table, Thead, Tbody, Tr, Th, Td, TableContainer, Text, Box, Flex, IconButton, Spinner } from '@chakra-ui/react';
import { useState } from 'react';

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
  const chains = ['Cosmos', 'Stargaze', 'Osmosis', 'Regen', 'Sommelier', 'Juno', 'Dydx', 'Saga', 'Agoric'];
  const [currentChainIndex, setCurrentChainIndex] = useState(0);

  // Switcher lets us use a pretty name for the chain in the UI, but query the chain by its actual name.
  const currentChainName = chains[currentChainIndex];
  let newChainName: string | undefined;
  if (currentChainName === 'Cosmos') {
    newChainName = 'cosmoshub';
  } else if (currentChainName === 'Osmosis') {
    newChainName = 'osmosis';
  } else if (currentChainName === 'Stargaze') {
    newChainName = 'stargaze';
  } else if (currentChainName === 'Regen') {
    newChainName = 'regen';
  } else if (currentChainName === 'Sommelier') {
    newChainName = 'sommelier';
  } else if (currentChainName === 'Juno') {
    newChainName = 'juno';
  } else if (currentChainName === 'Dydx') {
    newChainName = 'dydx';
  } else if (currentChainName === 'Saga') {
    newChainName = 'saga';
  } else if (currentChainName === 'Agoric') {
    newChainName = 'agoric';
  } else {
    // Default case
    newChainName = currentChainName;
  }

  const { unbondingData, isLoading } = useUnbondingQuery(newChainName, address);

  // Handlers for chain slider
  const handleLeftArrowClick = () => {
    setCurrentChainIndex((prevIndex: number) => (prevIndex === 0 ? chains.length - 1 : prevIndex - 1));
  };

  const handleRightArrowClick = () => {
    setCurrentChainIndex((prevIndex: number) => (prevIndex === chains.length - 1 ? 0 : prevIndex + 1));
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
              _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
              _active={{
                transform: 'scale(0.75)',
                color: 'complimentary.800',
              }}
              color="white"
            />
            <Box minWidth="100px" textAlign="center">
              <Text>{chains[currentChainIndex]}</Text>
            </Box>
            <IconButton
              icon={<ChevronRightIcon />}
              onClick={handleRightArrowClick}
              aria-label="Next chain"
              variant="ghost"
              _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
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
              _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
              _active={{
                transform: 'scale(0.75)',
                color: 'complimentary.800',
              }}
              color="white"
            />
            <Box minWidth="100px" textAlign="center">
              <Text>{chains[currentChainIndex]}</Text>
            </Box>
            <IconButton
              icon={<ChevronRightIcon />}
              onClick={handleRightArrowClick}
              aria-label="Next chain"
              variant="ghost"
              _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
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
            <Spinner size="xl" color="complimentary.900" />
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
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            color="white"
          />
          <Box minWidth="100px" textAlign="center">
            <Text>{chains[currentChainIndex]}</Text>
          </Box>
          <IconButton
            icon={<ChevronRightIcon />}
            onClick={handleRightArrowClick}
            aria-label="Next chain"
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
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
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.900">
                      Burn Amount
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.900" display={hideOnMobile}>
                      Status
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.900">
                      Redemption Amount
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.900" display={hideOnMobile}>
                      Epoch Number
                    </Th>
                    <Th textAlign="center" borderBottomColor={'transparent'} color="complimentary.900" display={hideOnMobile}>
                      Completion Time
                    </Th>
                  </Tr>
                </Thead>
                <Tbody>
                  {unbondingData?.withdrawals.map((withdrawal, index) => {
                    const shiftAmount = formatDenom(withdrawal.burn_amount.denom) === 'qDYDX' ? -18 : -6;
                    return (
                      <Tr _even={{ bg: 'rgba(255, 128, 0, 0.1)' }} key={index}>
                        <Td textAlign="center" borderBottomColor={'transparent'}>
                          {Number(shiftDigits(withdrawal.burn_amount.amount, shiftAmount))} {formatDenom(withdrawal.burn_amount.denom)}
                        </Td>
                        <Td textAlign="center" borderBottomColor={'transparent'} display={hideOnMobile}>
                          {statusCodes.get(withdrawal.status)}
                        </Td>
                        <Td textAlign="center" borderBottomColor={'transparent'}>
                          {withdrawal.amount.map((amt) => `${shiftDigits(amt.amount, shiftAmount)} ${formatDenom(amt.denom)}`).join(', ')}
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
