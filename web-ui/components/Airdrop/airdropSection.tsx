import { ChevronDownIcon, ChevronRightIcon, InfoOutlineIcon } from '@chakra-ui/icons';
import {
  Accordion,
  AccordionItem,
  AccordionButton,
  AccordionPanel,
  Box,
  Button,
  Flex,
  Text,
  Progress,
  Tooltip,
  VStack,
  HStack,
  Icon,
  Badge,
  useDisclosure,
} from '@chakra-ui/react';

interface AirdropAccordionItemProps {
  index: number;
  defaultIsOpen?: boolean;
}

const isBeta = true;

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
  return isBeta ? (
    // What to render if isBeta is true
    <Flex
      w="100%"
      h="sm"
      p={4}
      borderRadius="lg"
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
      gap={6}
      color="white"
      position="relative"
      _before={{
        content: '""',
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,

        backgroundSize: 'contain',
        backgroundPosition: 'center',
        backdropFilter: 'blur(10px)',
        filter: 'contrast(0.5)',
        opacity: 0.2,
        borderRadius: 'lg',
      }}
    >
      <Box position="relative" zIndex="docked">
        {' '}
        <Text fontSize="xl" fontWeight="bold">
          The Airdrop page is under construction
        </Text>
        <Text fontSize="xl" fontWeight="bold">
          Please check back later
        </Text>
      </Box>
    </Flex>
  ) : (
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
