import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import { Box, Flex, Text, Button, IconButton, VStack, Image, Heading, SlideFade } from '@chakra-ui/react';
import { color } from 'framer-motion';
import { useState } from 'react';

import { useIntentQuery } from '@/hooks/useQueries';

export interface StakingIntentProps {
  address: string;
}

const StakingIntent: React.FC<StakingIntentProps> = ({ address }) => {
  const { intent, isLoading, isError } = useIntentQuery('cosmoshub', address ?? '');

  const validators = [
    { name: 'Validator 1', logo: '/validator1.png', percentage: '30%' },
    { name: 'Validator 2', logo: '/validator2.png', percentage: '40%' },
  ];

  const chains = ['Cosmos', 'Osmosis', 'Stargaze', 'Regen', 'Sommelier'];
  const [currentChainIndex, setCurrentChainIndex] = useState(0);

  const handleLeftArrowClick = () => {
    setCurrentChainIndex((prevIndex) => (prevIndex === 0 ? chains.length - 1 : prevIndex - 1));
  };

  const handleRightArrowClick = () => {
    setCurrentChainIndex((prevIndex) => (prevIndex === chains.length - 1 ? 0 : prevIndex + 1));
  };

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
          <IconButton
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            aria-label="Previous chain"
            icon={<ChevronLeftIcon w={'25px'} h={'25px'} />}
            onClick={handleLeftArrowClick}
          />
          <SlideFade in={true} key={currentChainIndex}>
            <Text fontSize="lg" fontWeight="semibold">
              {chains[currentChainIndex]}
            </Text>
          </SlideFade>
          <IconButton
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
            variant="ghost"
            aria-label="Next chain"
            icon={<ChevronRightIcon w={'25px'} h={'25px'} />}
            onClick={handleRightArrowClick}
          />
        </Flex>

        <VStack spacing={2} align="stretch">
          {validators.map((validator, index) => (
            <Flex key={index} justifyContent="space-between" w="full" alignItems="center">
              <Flex alignItems="center" gap={2}>
                <Image alt="" src={validator.logo} boxSize="24px" borderRadius="full" />
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
