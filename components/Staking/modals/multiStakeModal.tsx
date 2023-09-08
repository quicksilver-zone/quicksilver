import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Accordion,
  AccordionItem,
  AccordionButton,
  AccordionPanel,
  AccordionIcon,
  Box,
  Divider,
  Text,
  Table,
  TableCaption,
  Tbody,
  Td,
  Tfoot,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react';
import React from 'react';

import { type ExtendedValidator as Validator } from '@/utils';

export const ValidatorsTable: React.FC<{
  validators: Validator[];
}> = ({ validators }) => {
  return (
    <Box px={4} overflowX="auto">
      <Table variant="simple" size="md">
        <Thead>
          <Tr>
            <Th color="white" fontSize={'16px'}>
              Moniker
            </Th>
            <Th color="white" fontSize={'16px'}>
              Commission
            </Th>
            <Th color="white" fontSize={'16px'}>
              Missed Blocks
            </Th>
            <Th color="white" fontSize={'16px'}>
              Rank
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          {validators.map((validator, index) => (
            <Tr key={index}>
              <Td>{validator.name}</Td>
              <Td>
                {validator.commission
                  ? (
                      (parseFloat(
                        validator.commission,
                      ) /
                        1e18) *
                      100
                    ).toFixed(2) + '%'
                  : 'N/A'}
              </Td>
              <Td></Td>
              <Td></Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </Box>
  );
};

interface MultiModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
  validators: Validator[];
}

export const MultiModal: React.FC<
  MultiModalProps
> = ({ isOpen, onClose, validators }) => {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="2xl"
    >
      {/* Set the size here */}
      <ModalOverlay />
      <ModalContent bgColor="#1A1A1A">
        <ModalHeader bgColor="#1A1A1A" p={0}>
          <Accordion allowToggle>
            <AccordionItem border="none">
              <h2>
                <AccordionButton
                  _hover={{
                    bgColor: 'transparent',
                  }}
                  p={4}
                >
                  <Box
                    fontSize="24px"
                    textAlign="left"
                    pr={3}
                  >
                    Validator Selection
                  </Box>
                  <AccordionIcon color="complimentary.900" />
                </AccordionButton>
              </h2>
              <AccordionPanel mt={-2}>
                <Text
                  color="white"
                  fontSize="18px"
                >
                  Choose which validator(s) you
                  would like to liquid stake to.
                </Text>
              </AccordionPanel>
            </AccordionItem>
          </Accordion>
        </ModalHeader>
        <ModalCloseButton />{' '}
        {/* Positioning by default should be top right */}
        <Divider
          bgColor="complimentary.900"
          alignSelf="center"
          w="95%"
          m="auto"
        />
        <ModalBody
          bgColor="#1A1A1A"
          borderRadius={'6px'}
        >
          <Box mt={4}>
            <ValidatorsTable
              validators={validators}
            />
          </Box>
        </ModalBody>
        <ModalFooter></ModalFooter>
      </ModalContent>
    </Modal>
  );
};
