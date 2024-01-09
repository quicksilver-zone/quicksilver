import { Box, Flex, Text, Button, VStack, useColorModeValue, HStack, SkeletonCircle, Spinner } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';

import { defaultChainName } from '@/config';
import { useAPYQuery, useBalanceQuery, useParamsQuery, useZoneQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';
import { BsCoin } from 'react-icons/bs';

import { DepositModal } from './modals/qckDepositModal';
import { WithdrawModal } from './modals/qckWithdrawModal';

interface QuickBoxProps {
  stakingApy?: number;
}

const QuickBox: React.FC<QuickBoxProps> = ({ stakingApy }) => {
  const { address, isWalletConnected } = useChain(defaultChainName);

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

  const decimalValue = parseFloat(stakingApy?.toString() ?? '0');
  const percentageValue = decimalValue * 100;
  const percentageString = percentageValue.toString();

  const truncatedPercentage = percentageString.slice(0, percentageString.indexOf('.') + 3);

  const quickStakingApy = () => {
    if (stakingApy) {
      return (
        <Text fontSize="lg" fontWeight="semibold">
          {truncatedPercentage}%
        </Text>
      );
    } else {
      return (
        <Box display="inline-block">
          <SkeletonCircle size="8" startColor="complimentary.900" endColor="complimentary.400" />
        </Box>
      );
    }
  };

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
          {quickStakingApy()}
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
