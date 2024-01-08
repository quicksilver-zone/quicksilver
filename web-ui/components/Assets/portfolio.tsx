import { shiftDigits } from '@/utils';
import { Box, Flex, Text, Icon, VStack, HStack, Heading, Spinner, Tooltip, Grid } from '@chakra-ui/react';

interface PortfolioItemInterface {
  title: string;
  percentage: string;
  progressBarColor: string;
  amount: string;
  qTokenPrice: number;
}

interface MyPortfolioProps {
  portfolioItems: PortfolioItemInterface[];
  isWalletConnected: boolean;
  totalValue: number;
  averageApy: number;
  totalYearlyYield: number;
}
const MyPortfolio: React.FC<MyPortfolioProps> = ({ portfolioItems, isWalletConnected, totalValue, averageApy, totalYearlyYield }) => {
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

  if (!totalValue) {
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
                    {shiftDigits(averageApy.toFixed(2), 2)}%
                  </Text>
                </HStack>
                <Text textAlign="center">
                  <Text as="span" fontSize="md" fontWeight="light">
                    Yearly Yield:{' '}
                  </Text>
                  <Text as="span" fontSize="md" fontWeight="medium">
                    ${totalYearlyYield.toFixed(2)}
                  </Text>
                </Text>
              </VStack>
            </Flex>
          </VStack>
        </Flex>

        <Flex justifyContent="flex-start" borderRadius={6} alignItems="flex-start" gap={4}>
          <VStack alignSelf="stretch" h="158px" mr={4} overflowY="auto" borderRadius={6} alignItems="center" gap={3}>
            {portfolioItems
              .filter((item) => Number(item.amount) > 0)
              .map((item) => (
                <PortfolioItem
                  key={item.title}
                  title={item.title}
                  percentage={Number(item.percentage)}
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
  <Grid templateColumns="2fr 6fr 1fr" gap={4} alignItems="center" width="100%">
    <HStack spacing={-5}>
      <Tooltip label={`Price: ${qTokenPrice.toFixed(2)}`} placement="top">
        <Text textAlign="left" minWidth="50px">
          {Number(amount).toFixed(1)}
        </Text>
      </Tooltip>
      <Text textAlign={'left'} fontSize="md" fontWeight="medium">
        {title}
      </Text>
    </HStack>
    <Box w="100%" h="8px" pos="relative">
      <Box w="100%" h="8px" pos="absolute" bg="complimentary.100" borderRadius="md" />
      <Box w={`${percentage * 100}%`} h="8px" pos="absolute" bg={progressBarColor} borderRadius="md" />
    </Box>
    <Tooltip label={`Value: $${(qTokenPrice * Number(amount)).toFixed(2)}`}>
      <Text textAlign="right" minWidth="50px">
        {`${(percentage * 100).toFixed(0)}%`}
      </Text>
    </Tooltip>
  </Grid>
);

export default MyPortfolio;
