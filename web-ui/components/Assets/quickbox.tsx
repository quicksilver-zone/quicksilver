import { Box, Flex, Text, Button, VStack, useColorModeValue, HStack, SkeletonCircle, Spinner } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';

import { defaultChainName } from '@/config';
import { useAPYQuery, useBalanceQuery, useParamsQuery, useZoneQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';
import { BsCoin } from 'react-icons/bs';

import { DepositModal } from './modals/qckDepositModal';
import { WithdrawModal } from './modals/qckWithdrawModal';

const QuickBox = () => {
  const { address, isWalletConnected } = useChain(defaultChainName);
  const { params } = useParamsQuery(defaultChainName);
  console.log(params);
  const { balance, isLoading } = useBalanceQuery(defaultChainName, address ?? '');
  const tokenBalance = Number(shiftDigits(balance?.balance?.amount ?? '', -6))
    .toFixed(2)
    .toString();

  if (!isWalletConnected) {
    return (
      <Flex direction="column" p={5} borderRadius="lg" align="center" justify="space-around" w="full" h="full">
        <Text fontSize="xl" textAlign="center">
          Wallet is not connected. Please connect your wallet to interact with your QCK tokens.
        </Text>
      </Flex>
    );
  }

  if (!balance) {
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
    <Flex direction="column" p={5} borderRadius="lg" align="center" justify="space-around" w="full" h="full">
      <VStack spacing={6}>
        {' '}
        <HStack>
          <BsCoin color="#FF8000" size={30} />
          <Text fontSize="3xl" fontWeight="bold">
            QCK
          </Text>
        </HStack>
        <HStack>
          <Text fontSize="2xl" fontWeight="bold"></Text>
          <Text fontSize="md" fontWeight="normal">
            STAKING APY:
          </Text>
          <Text fontSize="lg" fontWeight="semibold">
            39.12%
          </Text>
        </HStack>
        <VStack spacing={1} alignItems="flex-start" w="full">
          <HStack gap={2}>
            <Text fontSize="sm">ON QUICKSILVER:</Text>
            {isLoading === true && !balance && <SkeletonCircle size="2" startColor="complimentary.900" endColor="complimentary.400" />}
            {!isLoading && balance && (
              <Text fontSize="lg" fontWeight="semibold">
                {tokenBalance}
              </Text>
            )}
          </HStack>
        </VStack>
        <DepositModal />
        <WithdrawModal />
      </VStack>
    </Flex>
  );
};

export default QuickBox;
