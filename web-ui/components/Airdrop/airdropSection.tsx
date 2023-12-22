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
} from '@chakra-ui/react';
import { CheckIcon, ChevronRightIcon, InfoOutlineIcon } from '@chakra-ui/icons';

const AirdropSection = () => {
  const bgColor = useColorModeValue('white', 'gray.800');
  const textColor = useColorModeValue('gray.600', 'gray.200');

  return (
    <VStack spacing={4} align="stretch">
      <Box bg={bgColor} p={5} shadow="md" borderRadius="lg">
        <Flex justifyContent="space-between" alignItems="center" mb={5}>
          <Text fontSize="xl" fontWeight="bold">
            Cosmos Hub Active Airdrop
          </Text>
          <Tooltip label="This is your total allocation" aria-label="A tooltip">
            <Icon as={InfoOutlineIcon} />
          </Tooltip>
        </Flex>
        <Progress colorScheme="orange" size="sm" value={40} mb={5} />
        <Accordion allowToggle>
          {Array.from({ length: 5 }).map((_, index) => (
            <AccordionItem key={index}>
              <h2>
                <AccordionButton _expanded={{ bg: 'orange.100', color: 'orange.800' }}>
                  <Box flex="1" textAlign="left">
                    Step {index + 1}: {index === 0 ? 'CLAIM INITIAL QCK AIRDROP' : 'ACTION DESCRIPTION'}
                  </Box>
                  <AccordionIcon />
                </AccordionButton>
              </h2>
              <AccordionPanel pb={4}>
                {/* Replace this with actual content */}
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
              </AccordionPanel>
            </AccordionItem>
          ))}
        </Accordion>
      </Box>
      <Box bg={bgColor} p={5} shadow="md" borderRadius="lg">
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
              <Text fontSize="md" color={textColor}>
                {chain}
              </Text>
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
