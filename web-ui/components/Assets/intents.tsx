import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import { Box, Flex, Text, Button, IconButton, HStack, VStack, Image, Heading, Spacer } from '@chakra-ui/react';

const StakingIntent = () => {
  // Example data - replace with your actual data
  const validators = [
    { name: 'Validator 1', logo: '/validator1.png', percentage: '30%' },
    { name: 'Validator 2', logo: '/validator2.png', percentage: '40%' },
    // ... more validators
  ];

  const chains = ['Chain 1', 'Chain 2']; // Chain names
  let currentChainIndex = 0; // Replace with state logic to switch chains

  const handleLeftArrowClick = () => {
    // Logic to switch to the previous chain
  };

  const handleRightArrowClick = () => {
    // Logic to switch to the next chain
  };

  return (
    <Box w="fit-content" p={0} color="white">
      <Flex justifyContent="space-between" gap={14} alignItems="center">
        <Heading top={0} fontSize="lg" fontWeight="bold" textTransform="uppercase" noOfLines={1}>
          Stake Intent
        </Heading>

        <Button top={0} bgColor={'transparent'} size="sm">
          Edit/Reset Intent
        </Button>
      </Flex>

      <Flex my={4} alignItems="center" justifyContent="center">
        <IconButton bgColor={'transparent'} aria-label="Previous chain" icon={<ChevronLeftIcon />} onClick={handleLeftArrowClick} />
        <Text mx={4}>{chains[currentChainIndex]}</Text>
        <IconButton bgColor={'transparent'} aria-label="Next chain" icon={<ChevronRightIcon />} onClick={handleRightArrowClick} />
      </Flex>

      <VStack spacing={2}>
        {validators.map((validator, index) => (
          <HStack key={index} justifyContent="space-between" w="100%">
            <HStack spacing={3}>
              <Image src={validator.logo} boxSize="24px" borderRadius="full" />
              <Text>{validator.name}</Text>
            </HStack>
            <Text fontWeight="bold" color="complementary.900">
              {validator.percentage}
            </Text>
          </HStack>
        ))}
      </VStack>
    </Box>
  );
};

export default StakingIntent;
