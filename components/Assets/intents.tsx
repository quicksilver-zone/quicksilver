import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import { Box, Flex, Text, Button, IconButton, VStack, Image, Heading } from '@chakra-ui/react';

const StakingIntent = () => {
  const validators = [
    { name: 'Validator 1', logo: '/validator1.png', percentage: '30%' },
    { name: 'Validator 2', logo: '/validator2.png', percentage: '40%' },
  ];

  const chains = ['Cosmos', 'Osmosis'];
  let currentChainIndex = 0;

  const handleLeftArrowClick = () => {};

  const handleRightArrowClick = () => {};

  return (
    <Box w="full" color="white" borderRadius="lg" p={4} gap={6}>
      <VStack spacing={4} align="stretch">
        <Flex gap={6} justifyContent="space-between" alignItems="center">
          <Heading fontSize="lg" fontWeight="bold" textTransform="uppercase">
            Stake Intent
          </Heading>
          <Button color="GrayText" variant="link">
            Edit Intent
            <ChevronRightIcon />
          </Button>
        </Flex>

        <Flex borderBottom="1px" borderBottomColor="complimentary.900" alignItems="center" justifyContent="space-between">
          <IconButton variant="ghost" aria-label="Previous chain" icon={<ChevronLeftIcon />} onClick={handleLeftArrowClick} />
          <Text fontSize="lg" fontWeight="semibold">
            {chains[currentChainIndex]}
          </Text>
          <IconButton variant="ghost" aria-label="Next chain" icon={<ChevronRightIcon />} onClick={handleRightArrowClick} />
        </Flex>

        <VStack spacing={2} align="stretch">
          {validators.map((validator, index) => (
            <Flex key={index} justifyContent="space-between" w="full" alignItems="center">
              <Flex alignItems="center" gap={2}>
                <Image src={validator.logo} boxSize="24px" borderRadius="full" />
                <Text fontSize="md">{validator.name}</Text>
              </Flex>
              <Text fontSize="lg" fontWeight="bold">
                {validator.percentage}
              </Text>
            </Flex>
          ))}
        </VStack>
      </VStack>
    </Box>
  );
};

export default StakingIntent;
