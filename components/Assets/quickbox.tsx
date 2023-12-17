import { Box, Flex, Text, Button, VStack, useColorModeValue, HStack } from '@chakra-ui/react';

const QuickBox = () => {
  return (
    <Flex direction="column" p={5} borderRadius="lg" align="center" justify="space-around" w="full" h="full">
      <VStack spacing={6}>
        {' '}
        <HStack>
          <Box w={6} h={6} borderRadius="full" bg="gray.300" />
          <Text fontSize="3xl" fontWeight="bold">
            QCK
          </Text>
        </HStack>
        <HStack>
          <Text fontSize="2xl" fontWeight="bold">
            12.34%
          </Text>
          <Text fontSize="md" fontWeight="normal">
            STAKING APY
          </Text>
        </HStack>
        <VStack spacing={1} alignItems="flex-start" w="full">
          <HStack gap={2}>
            <Text fontSize="sm">QUICKSILVER BALANCE</Text>
            <Text fontSize="lg" fontWeight="semibold">
              10.12
            </Text>
          </HStack>
          <HStack gap={2}>
            <Text fontSize="sm">NON-NATIVE BALANCE</Text>
            <Text fontSize="lg" fontWeight="semibold">
              10.12
            </Text>
          </HStack>
        </VStack>
        <Button color={'white'} w="full" variant="outline">
          Deposit
        </Button>
        <Button color={'white'} w="full" variant="outline">
          Withdraw
        </Button>
      </VStack>
    </Flex>
  );
};

export default QuickBox;
