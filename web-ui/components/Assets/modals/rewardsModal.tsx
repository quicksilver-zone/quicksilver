import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Divider,
  Tooltip,
  Grid,
  Text,
  Button,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Box,
  Icon,
  Heading,
  Flex,
  useDisclosure,
  Stat,
  StatHelpText,
  StatLabel,
  StatNumber,
  HStack,
  TableContainer,
  Tfoot,
  VStack,
} from '@chakra-ui/react';

import { useLiquidEpochQuery, useLiquidRewardsQuery } from '@/hooks/useQueries';

const RewardsModal = ({
  address,

  isOpen,
  onClose,
}: {
  address: string;

  isOpen: boolean;
  onClose: () => void;
}) => {
  const tokenDetails = [
    { name: 'ATOM', amount: '100', chainId: 'cosmoshub-4' },
    { name: 'OSMO', amount: '200', chainId: 'osmosis-1' },
    { name: 'SCRT', amount: '300', chainId: 'secret-1' },
    { name: 'ATOM', amount: '100', chainId: 'cosmoshub-4' },
    { name: 'OSMO', amount: '200', chainId: 'osmosis-1' },
    { name: 'SCRT', amount: '300', chainId: 'secret-1' },
    { name: 'ATOM', amount: '100', chainId: 'cosmoshub-4' },
    { name: 'OSMO', amount: '200', chainId: 'osmosis-1' },
    { name: 'SCRT', amount: '300', chainId: 'secret-1' },
  ];

  return (
    <Modal size={'xl'} isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent bgColor="rgb(32,32,32)">
        <ModalHeader color="white" fontSize="xl">
          <HStack>
            <Text>Rewards</Text>
          </HStack>
          <Divider mt={3} bgColor={'cyan.500'} />
        </ModalHeader>
        <ModalCloseButton color={'complimentary.900'} />
        <ModalBody>
          <TableContainer maxH="170px" overflowY="auto">
            <Table variant="simple" colorScheme="whiteAlpha" size="sm">
              <Thead position="sticky" top={0} bg="rgb(32,32,32)" zIndex={1}>
                <Tr>
                  <Th color="complimentary.900">Token</Th>
                  <Th color="complimentary.900" isNumeric>
                    Amount
                  </Th>
                </Tr>
              </Thead>
              <Tbody>
                {tokenDetails.map((detail, index) => (
                  <Tr key={index}>
                    <Td color="white">{detail.name}</Td>
                    <Td color="white" isNumeric>
                      {detail.amount}
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </TableContainer>
          <Button
            mt={4}
            _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
            _hover={{ bgColor: 'rgba(255,128,0, 0.25)', color: 'complimentary.300' }}
            color="white"
            size="sm"
            w="160px"
            variant="outline"
          >
            Unwind
          </Button>
        </ModalBody>

        <ModalFooter></ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default RewardsModal;
