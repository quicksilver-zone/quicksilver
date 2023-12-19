import { Box, Flex, Text, Icon, VStack, HStack, Stack, Heading, Divider, Progress } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import { IoWallet } from 'react-icons/io5';

import { useQBalanceQuery } from '@/hooks/useQueries';
import { PortfolioItem } from '@/pages/assets';
import { shiftDigits } from '@/utils';

interface MyPortfolioProps {
  portfolioItems: PortfolioItem[];
  isWalletConnected: boolean;
  totalValue: number;
}

const MyPortfolio: React.FC<MyPortfolioProps> = ({ portfolioItems, isWalletConnected, totalValue }) => {
  if (!isWalletConnected) {
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
        <Text fontSize="xl" textAlign="center">
          Wallet is not connected. Please connect your wallet to view your portfolio.
        </Text>
      </Flex>
    );
  }

  return (
    <Flex w="100%" h="100%" p={4} borderRadius="lg" flexDirection="column" justifyContent="center" alignItems="center" gap={6}>
      <Heading color={'white'} alignSelf="stretch" fontSize="lg" fontWeight="bold" textTransform="uppercase" noOfLines={1}>
        My QUICKSILVER Portfolio
      </Heading>

      <VStack alignSelf="stretch" h={'300px'} alignItems="flex-end" gap={4}>
        <Flex
          alignSelf="stretch"
          borderBottom="1px"
          borderBottomColor="complimentary.900"
          justifyContent="flex-start"
          alignItems="center"
          gap={5}
        >
          <VStack flex="1" pt={1} pb={2.5} justifyContent="center" alignItems="flex-start" gap={2}>
            <Flex alignSelf="stretch" justifyContent="space-between" alignItems="center">
              <VStack w="161px" alignItems="flex-start" gap={2}>
                <Text fontSize="sm" fontWeight="medium" textTransform="uppercase">
                  TOTAL
                </Text>
                <Text textAlign="right" fontSize="2xl" fontWeight="bold">
                  ${totalValue.toFixed(2)}
                </Text>
              </VStack>

              <VStack alignItems="flex-end" gap={3}>
                <HStack justifyContent="flex-start" alignItems="flex-start" gap={2.5}>
                  <Text fontSize="md" fontWeight="light">
                    AVG APY:
                  </Text>
                  <Text fontSize="md" fontWeight="medium">
                    6.56%
                  </Text>
                </HStack>
                <Text textAlign="center">
                  <Text as="span" fontSize="md" fontWeight="light">
                    Yearly Yield:{' '}
                  </Text>
                  <Text as="span" fontSize="md" fontWeight="medium">
                    $3,917
                  </Text>
                </Text>
              </VStack>
            </Flex>
          </VStack>
        </Flex>

        <Flex justifyContent="flex-start" borderRadius={6} alignItems="flex-start" gap={4}>
          <VStack alignSelf="stretch" h="158px" overflowY="auto" borderRadius={6} alignItems="flex-start" gap={3}>
            {portfolioItems
              .filter((item) => Number(item.amount) > 0) // Filter out items with 0 amount
              .map((item) => (
                <PortfolioItem
                  key={item.title}
                  title={item.title}
                  percentage={item.percentage}
                  progressBarColor={item.progressBarColor}
                  amount={item.amount}
                  qTokenPrice={item.qTokenPrice}
                />
              ))}
          </VStack>
        </Flex>
      </VStack>
    </Flex>
  );
};

interface PortfolioItemProps {
  title: string;
  percentage: number;
  progressBarColor: string;
  amount: string;
  qTokenPrice: number;
}

const PortfolioItem: React.FC<PortfolioItemProps> = ({ title, percentage, progressBarColor, amount, qTokenPrice }) => (
  <Flex alignSelf="stretch" justifyContent="space-between" gap={16} alignItems="center">
    <HStack h="24px" justifyContent="flex-start" alignItems="center" gap={2.75}>
      <Text>{Number(amount).toFixed(2).toString()}</Text>
      <Text fontSize="md" fontWeight="medium">
        {title}
      </Text>
    </HStack>
    <HStack justifyContent="center" alignItems="center" gap={4}>
      <Box w="121px" h="8px" pos="relative">
        <Box w="121px" h="8px" pos="absolute" bg="complimentary.100" borderRadius="md" />
        <Box w={`${percentage * 100}%`} h="8px" pos="absolute" bg={progressBarColor} borderRadius="md" />
      </Box>
      <Text w="44px" textAlign="right" fontSize="sm" fontWeight="normal">
        {`${percentage * 100}%`}
      </Text>
    </HStack>
  </Flex>
);

export default MyPortfolio;
