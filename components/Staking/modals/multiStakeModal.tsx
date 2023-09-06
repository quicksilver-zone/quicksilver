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
} from '@chakra-ui/react';
import React from 'react';

interface MultiModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
}

export const MultiModal: React.FC<
  MultiModalProps
> = ({ isOpen, onClose }) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
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
                  Ligma figma sigma digma Ligma
                  figma sigma digma Ligma figma
                  sigma digma Ligma figma sigma
                  digma
                </Text>
              </AccordionPanel>
            </AccordionItem>
          </Accordion>
        </ModalHeader>
        <ModalCloseButton />
        <Divider
          bgColor="complimentary.900"
          alignSelf="center"
          w="95%"
        />
        <ModalBody
          bgColor="#1A1A1A"
          borderRadius={'6px'}
        ></ModalBody>
        <ModalFooter></ModalFooter>
      </ModalContent>
    </Modal>
  );
};
