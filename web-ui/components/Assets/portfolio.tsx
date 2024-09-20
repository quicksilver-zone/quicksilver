import { Flex, Text, VStack, HStack, Heading, Spinner, SimpleGrid, Center, Image, SkeletonText } from '@chakra-ui/react';
import { Divider } from '@interchain-ui/react';

import { shiftDigits, formatNumber } from '@/utils';

interface PortfolioItemInterface {
  title: string;
  amount: string;
  qTokenPrice: number;
  chainId: string;
}

interface MyPortfolioProps {
  portfolioItems: PortfolioItemInterface[];
  isWalletConnected: boolean;
  totalValue: number;
  averageApy: number;
  totalYearlyYield: number;
  isLoading: boolean;
}
const MyPortfolio: React.FC<MyPortfolioProps> = ({
  portfolioItems,
  isWalletConnected,
  totalValue,
  averageApy,
  totalYearlyYield,
  isLoading,
}) => {
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

  if (isLoading && !portfolioItems.length) {
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
        <Spinner w={'220px'} h="220px" color="complimentary.900" />
      </Flex>
    );
  }

  return (
    <Flex w="100%" h="100%" px={4} mt={5} borderRadius="lg" flexDirection="column" justifyContent="center" alignItems="center" gap={6}>
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
            <SimpleGrid columns={3} spacing={10} w="full">
              <Center>
                <VStack spacing={1}>
                  <Text fontSize="sm" fontWeight="medium" textTransform="uppercase">
                    TOTAL VALUE
                  </Text>
                  <Text fontSize="2xl" fontWeight="bold">
                    ${formatNumber(totalValue)}
                  </Text>
                </VStack>
              </Center>
              <Center>
                <VStack spacing={1}>
                  <Text fontSize="sm" fontWeight="medium" textTransform="uppercase">
                    AVERAGE APY
                  </Text>
                  <Text fontSize="2xl" fontWeight="bold">
                    {isNaN(averageApy) ? '0%' : `${shiftDigits(averageApy.toFixed(2), 2)}%`}
                  </Text>
                </VStack>
              </Center>
              <Center>
                <VStack spacing={1}>
                  <Text fontSize="sm" fontWeight="medium" textTransform="uppercase">
                    Est. Yield
                  </Text>
                  <Text fontSize="2xl" fontWeight="bold">
                    ${formatNumber(totalYearlyYield)}
                  </Text>
                </VStack>
              </Center>
            </SimpleGrid>
          </VStack>
        </Flex>
        {isLoading && (
          <Flex w="100%" justifyContent="center" alignItems="center">
            <SkeletonText noOfLines={1} h={'20px'} w={'120px'} startColor="compplimentary.900" endColor="complimentary.100" />
            <SkeletonText noOfLines={1} h={'20px'} w={'120px'} startColor="rgb(107, 105, 105)" endColor="rgb(168, 168, 168)" />
            <SkeletonText noOfLines={1} h={'20px'} w={'120px'} startColor="compplimentary.900" endColor="complimentary.100" />
            <SkeletonText noOfLines={1} h={'20px'} w={'120px'} startColor="rgb(107, 105, 105)" endColor="rgb(168, 168, 168)" />
          </Flex>
        )}
        {totalValue === 0 && (
          <Flex w="100%" justifyContent="center" alignItems="center">
            <Text fontSize="xl" textAlign="center">
              You have no liquid staked assets.
            </Text>
          </Flex>
        )}
        <Flex w="100%" justifyContent="center" borderRadius={6} alignItems="center" gap={4}>
          <VStack alignSelf="stretch" h="185px" overflowY="auto" className="custom-scrollbar" borderRadius={6} alignItems="center" gap={3}>
            {totalValue > 0 && (
              <SimpleGrid position={'sticky'} bgColor={'rgb(26,26,26)'} top={0} minW="100%" columns={3} spacing={4}>
                <VStack spacing={1}>
                  <Text fontSize="sm" fontWeight="medium">
                    ASSET
                  </Text>
                  <Divider width="100px" />
                </VStack>
                <VStack spacing={1}>
                  <Text fontSize="sm" fontWeight="medium" textAlign={'center'}>
                    AMOUNT
                  </Text>
                  <Divider width="100px" />
                </VStack>
                <VStack spacing={1}>
                  <Text fontSize="sm" fontWeight="medium" textAlign={'center'}>
                    VALUE
                  </Text>
                  <Divider width="100px" />
                </VStack>
              </SimpleGrid>
            )}

            {portfolioItems
              .filter((item) => Number(item.amount) > 0)
              .map((item) => (
                <PortfolioItem
                  key={item.title}
                  title={item.title}
                  amount={item.amount}
                  qTokenPrice={item.qTokenPrice}
                  totalValue={totalValue}
                  index={portfolioItems.indexOf(item)}
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

  amount: string;
  qTokenPrice: number;
  totalValue: number;
  index: number;
}

const PortfolioItem: React.FC<PortfolioItemProps> = ({ title, amount, qTokenPrice, index }) => {
  const tokenValue = Number(amount) * qTokenPrice;

  const imgType = title === 'qATOM' || title === 'qSAGA' ? 'svg' : 'png'; // TODO: why?!

  return (
    <SimpleGrid
      _even={{ bg: 'rgba(255, 128, 0, 0.1)' }}
      _odd={{ bg: 'rgba(180, 173, 167, 0.1)' }}
      key={index}
      textAlign={'center'}
      alignItems={'center'}
      minW="400px"
      columns={3}
      spacing={4}
      py={1}
    >
      <HStack gap={3}>
        <Image alt={`${title}`} ml={2} borderRadius={'full'} src={`/img/networks/${title.toLowerCase()}.${imgType}`} boxSize="33px" />
        <Text>q{title.toLocaleLowerCase().slice(1).toLocaleUpperCase()}</Text>
      </HStack>
      <Text> {formatNumber(parseFloat(amount))}</Text>
      <Text>{tokenValue < 0.01 ? '>$0.01' : '$' + formatNumber(tokenValue)}</Text>
    </SimpleGrid>
  );
};

export default MyPortfolio;
