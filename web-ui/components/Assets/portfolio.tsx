import { Box, Flex, Text, Icon, VStack, HStack, Stack, Heading, Divider, Progress } from '@chakra-ui/react';
import { IoWallet } from 'react-icons/io5'; // Example icon, replace with your own

const MyPortfolio = () => {
  return (
    <Flex
      w="100%"
      h="100%"
      p={6}
      borderRadius="lg"
      border="1px"
      borderColor="rgba(0, 0, 0, 0.50)"
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
      gap={6}
    >
      <Heading alignSelf="stretch" fontSize="lg" fontWeight="bold" textTransform="uppercase" noOfLines={1}>
        My QUICKSILVER Portfolio
      </Heading>

      <VStack alignSelf="stretch" h={250} alignItems="flex-end" gap={4}>
        <Flex alignSelf="stretch" borderBottom="1px" borderBottomColor="black" justifyContent="flex-start" alignItems="center" gap={5}>
          <VStack
            flex="1"
            pt={1}
            pb={2.5}
            borderRight="1px"
            borderRightColor="gray.200"
            justifyContent="center"
            alignItems="flex-start"
            gap={2}
          >
            <Flex alignSelf="stretch" justifyContent="space-between" alignItems="center">
              <VStack w="161px" alignItems="flex-start" gap={2}>
                <Text color="black" fontSize="sm" fontWeight="medium" textTransform="uppercase">
                  TOTAL
                </Text>
                <Text textAlign="right" color="black" fontSize="2xl" fontWeight="bold">
                  $ 1,222.28
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

        {/* Repeat the following structure for each item */}
        <Flex justifyContent="flex-start" alignItems="flex-start" gap={4}>
          <VStack h="150px" justifyContent="center" alignItems="flex-start" gap={10}>
            {/* Item content */}
            <PortfolioItem title="qOSMO" percentage={0.75} progressBarColor="#0066FF" />
            {/* Repeat PortfolioItem for each item */}
          </VStack>
          <Box w="4px" h="113px" bg="#A3A3A3" borderRadius="40px" />
        </Flex>
      </VStack>
    </Flex>
  );
};

const PortfolioItem = ({ title, percentage, progressBarColor }) => (
  <Flex alignSelf="stretch" justifyContent="space-between" alignItems="center">
    <HStack h="24px" justifyContent="flex-start" alignItems="center" gap={2.75}>
      <Box w="24px" h="24px" bg="#DEDEDE" borderRadius="full" />
      <Text color="black" fontSize="md" fontWeight="medium">
        {title}
      </Text>
    </HStack>
    <HStack justifyContent="center" alignItems="center" gap={4}>
      <Box w="121px" h="8px" pos="relative">
        <Box w="121px" h="8px" pos="absolute" bg="#262A46" borderRadius="md" />
        <Box w={`${percentage}%`} h="8px" pos="absolute" bg={progressBarColor} borderRadius="md" />
      </Box>
      <Text w="44px" textAlign="right" color="black" fontSize="sm" fontWeight="normal">
        {percentage}%
      </Text>
    </HStack>
  </Flex>
);

export default MyPortfolio;
