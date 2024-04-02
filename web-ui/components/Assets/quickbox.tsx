import { Flex, Text, VStack, HStack, Spinner, Button, Stat, StatLabel, StatNumber, useDisclosure } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';

import { DepositModal } from './modals/qckDepositModal';
import { WithdrawModal } from './modals/qckWithdrawModal';
import RewardsModal from './modals/rewardsModal';

import { defaultChainName } from '@/config';
import { useBalanceQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';


interface QuickBoxProps {
  stakingApy?: number;
}

const QuickBox: React.FC<QuickBoxProps> = ({ stakingApy }) => {
  const { address } = useChain(defaultChainName);
  const { balance, isLoading } = useBalanceQuery(defaultChainName, address ?? '');
  const tokenBalance = Number(shiftDigits(balance?.balance?.amount ?? '', -6))
    .toFixed(2)
    .toString();
  const { isOpen, onOpen, onClose } = useDisclosure();

  if (!address) {
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
  const percentageValue = (decimalValue * 100).toFixed(0);
  const percentageString = percentageValue.toString();

  return (
    <Flex direction="column" py={8} borderRadius="lg" align="center" justify="space-around" w="full" h="full">
      <VStack w="75%" spacing={8}>
        <HStack borderBottom="1px" borderBottomColor="complimentary.900" w="full" justify="space-between">
          <Text fontWeight="bold" fontSize={'xl'} isTruncated>
            QCK
          </Text>
          <HStack>
            <Text fontSize="md" fontWeight="bold" isTruncated>
              {percentageString}%
            </Text>
            <Text fontSize="xs" fontWeight="light" isTruncated>
              APY
            </Text>
          </HStack>
        </HStack>

        <VStack>
          <Stat color={'white'}>
            <StatLabel fontSize={'lg'}>Quicksilver Balance</StatLabel>
            <StatNumber textAlign={'center'} color={'complimentary.900'} fontSize={'lg'}>
              {tokenBalance} QCK
            </StatNumber>
          </Stat>
        </VStack>
        <DepositModal />
        <WithdrawModal />
        <Button
          _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
          _hover={{ bgColor: 'rgba(255,128,0, 0.25)', color: 'complimentary.300' }}
          color="white"
          w="full"
          variant="outline"
          onClick={onOpen}
        >
          Unwind
        </Button>
      </VStack>
      <RewardsModal address={address} isOpen={isOpen} onClose={onClose} />
    </Flex>
  );
};

export default QuickBox;
