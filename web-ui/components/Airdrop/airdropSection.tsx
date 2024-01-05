import {
  Accordion,
  AccordionItem,
  AccordionButton,
  AccordionPanel,
  AccordionIcon,
  Box,
  Button,
  Flex,
  Text,
  Progress,
  Tooltip,
  VStack,
  HStack,
  useColorModeValue,
  Icon,
  Badge,
  useDisclosure,
} from '@chakra-ui/react';
import { CheckIcon, ChevronDownIcon, ChevronRightIcon, InfoOutlineIcon } from '@chakra-ui/icons';
import { useAccordionStyles } from '@chakra-ui/accordion';
import { useState } from 'react';

interface AirdropAccordionItemProps {
  index: number;
  defaultIsOpen?: boolean;
}

const AirdropAccordionItem: React.FC<AirdropAccordionItemProps> = ({ index, defaultIsOpen }) => {
  const { isOpen, onToggle } = useDisclosure({ defaultIsOpen });

  return (
    <AccordionItem
      mb="30px"
      borderRadius={isOpen ? '20px 20px 0 0' : '20px'}
      h="80px"
      borderBottomColor={'transparent'}
      zIndex={10 - index}
      borderTopColor={'transparent'}
      bgColor={'rgba(255, 128, 0, 1)'}
      shadow={'md'}
      position="relative"
      _hover={{ bgColor: 'rgba(255, 128, 0, 0.9)' }}
    >
      <h2>
        <AccordionButton onClick={onToggle} h="80px" justifyContent="space-between">
          <Badge borderRadius="full" px={3} py={1.5} mr={3} colorScheme="orange">
            {index + 1}
          </Badge>
          <Box flex="1" textAlign="left" alignItems={'center'} fontSize="lg" fontWeight="semibold">
            <Text fontSize={'2xl'}>{index === 0 ? 'CLAIM INITIAL QCK AIRDROP' : 'ACTION DESCRIPTION'}</Text>
          </Box>
          <ChevronDownIcon />
        </AccordionButton>
      </h2>
      {isOpen && (
        <AccordionPanel p={4} borderBottomRadius="20px" h="120px" pb={'30px'} bgColor="rgba(255, 128, 0, 1)">
          <Text>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
          </Text>
        </AccordionPanel>
      )}
    </AccordionItem>
  );
};

const AirdropSection = () => {
  const { isOpen, onToggle } = useDisclosure();

  return (
    <VStack spacing={4} align="stretch">
      <Box bg={'rgba(255,255,255,0.1)'} p={5} shadow="md" borderRadius="lg">
        <Flex justifyContent="space-between" alignItems="center" mb={5}>
          <Text fontSize="xl" fontWeight="bold">
            Cosmos Hub Active Airdrop
          </Text>
          <Tooltip label="This is your total allocation" aria-label="A tooltip">
            <Icon as={InfoOutlineIcon} />
          </Tooltip>
        </Flex>
        <Progress colorScheme="orange" size="sm" value={40} mb={5} />
        <Accordion allowToggle defaultIndex={[0]}>
          {Array.from({ length: 5 }).map((_, index) => (
            <AirdropAccordionItem key={index} index={index} defaultIsOpen={index === 0} />
          ))}
        </Accordion>
      </Box>
      <Box bg={'rgba(255,255,255,0.1)'} p={5} shadow="md" borderRadius="lg">
        <Flex justifyContent="space-between" alignItems="center" mb={5}>
          <Text fontSize="xl" fontWeight="bold">
            Participate in Other Airdrops
          </Text>
        </Flex>
        <HStack overflow="auto">
          {/* Replace these with actual chain data */}
          {['OSMOSIS', 'JUNO', 'STARGAZE', 'REGEN', 'COSMOS'].map((chain) => (
            <VStack key={chain} spacing={3} align="center">
              <Box w="100px" h="100px" bg="gray.100" borderRadius="md" />
              <Text fontSize="md">{chain}</Text>
              <Button rightIcon={<ChevronRightIcon />} colorScheme="teal" variant="link">
                View Airdrop
              </Button>
            </VStack>
          ))}
        </HStack>
      </Box>
    </VStack>
  );
};

export default AirdropSection;
