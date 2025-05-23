import {
  Modal,
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
  Flex,
  Button,
  Input,
  Spinner,
  InputGroup,
  InputLeftElement,
  Link,
} from '@chakra-ui/react';
import React from 'react';
import { FaSearch } from 'react-icons/fa';

import { ValidatorsTable } from '@/components/Staking/modals/validatorTable';
import { useValidatorsQuery } from '@/hooks/useQueries';
import { useValidatorLogos } from '@/hooks/useQueries';

interface MultiModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
  selectedChainName: string;
  selectedChainId: string;
  selectedValidators: { name: string; operatorAddress: string }[];
  setSelectedValidators: React.Dispatch<React.SetStateAction<{ name: string; operatorAddress: string }[]>>;
}

export const IntentMultiModal: React.FC<MultiModalProps> = ({
  isOpen,
  onClose,
  selectedChainName,
  selectedValidators,
  setSelectedValidators,
}) => {
  const [searchTerm, setSearchTerm] = React.useState<string>('');

  const { validatorsData, isLoading } = useValidatorsQuery(selectedChainName);

  const { data: logos } = useValidatorLogos(selectedChainName, validatorsData || []);

  const validators = validatorsData;
  const handleValidatorClick = (validator: { name: string; operatorAddress: string }) => {
    setSelectedValidators((prevState) => {
      const alreadySelected = prevState.some((v) => v.name === validator.name);

      if (!alreadySelected && prevState.length >= 100) {
        alert("You can't select more than 100 validators.");
        return prevState;
      }

      return alreadySelected ? prevState.filter((v) => v.name !== validator.name) : [...prevState, validator];
    });
  };

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(event.target.value.toLowerCase());
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="2xl">
      {/* Set the size here */}

      <ModalContent borderRadius={'10px'} maxHeight="70vh" bgColor="#1A1A1A">
        <ModalHeader borderRadius="10px" bgColor="#1A1A1A" p={0}>
          <Accordion allowToggle>
            <AccordionItem border="none">
              <h2>
                <AccordionButton
                  _hover={{
                    bgColor: 'transparent',
                  }}
                  p={6}
                >
                  <Box h="100%" mb={-4} pr={4}>
                    <Text ml={4} color="white" fontSize="24px" textAlign="left">
                      Validator Selection
                    </Text>
                  </Box>
                  <AccordionIcon mb={-4} color="complimentary.700" />
                </AccordionButton>
              </h2>
              <AccordionPanel textAlign="left" alignContent="center" justifyContent="center" mt={-2}>
                <Text fontWeight="light" pl={6} maxW="95%" color="white" fontSize="16px" letterSpacing={'wider'}>
                  Choose which validator(s) you would like to liquid stake to.
                </Text>
              </AccordionPanel>
            </AccordionItem>
          </Accordion>
        </ModalHeader>
        <ModalCloseButton color="white" />
        <Divider bgColor="complimentary.700" alignSelf="center" w="88%" m="auto" />
        <ModalBody bgColor="#1A1A1A" borderRadius={'6px'} justifyContent="center">
          {isLoading ? (
            <Box minH={'md'} display="flex" justifyContent="center" alignItems="center" height="200px">
              <Spinner h="50px" w="50px" color="complimentary.700" />
            </Box>
          ) : (
            <Box mt={-1}>
              <Flex py={2} px={4} alignContent="center" alignItems="center" justifyContent={'space-between'} w="100%" flexDirection={'row'}>
                <InputGroup>
                  <InputLeftElement pointerEvents="none">
                    <FaSearch color="orange" />
                  </InputLeftElement>
                  <Input
                    type="text"
                    color="white"
                    borderColor="complimentary.1000"
                    placeholder="validator moniker..."
                    fontWeight="light"
                    onChange={handleSearchChange}
                    width="55%"
                    borderRadius={'4px'}
                    _active={{
                      borderColor: 'complimentary.700',
                    }}
                    _selected={{
                      borderColor: 'complimentary.700',
                    }}
                    _hover={{
                      borderColor: 'complimentary.700',
                    }}
                    _focus={{
                      borderColor: 'complimentary.700',
                      boxShadow: '0 0 0 3px #FF8000',
                    }}
                  />
                </InputGroup>
                <Box
                  borderRadius="10px"
                  w="300px"
                  h="50px"
                  mr={2}
                  display="flex"
                  flexDirection="column"
                  justifyContent="space-between"
                ></Box>
              </Flex>
              <ValidatorsTable
                logos={logos || []}
                validators={validators || []}
                onValidatorClick={(validator) => {
                  const isSelected = selectedValidators.some((v) => v.name === validator.name);
                  if (selectedValidators.length < 100 || isSelected) {
                    handleValidatorClick(validator);
                  }
                }}
                selectedValidators={selectedValidators}
                searchTerm={searchTerm}
              />
              <Box mt={-12} w="100%" display="flex" justifyContent="center" alignItems="center">
                <Button
                  onClick={onClose}
                  h="30px"
                  w="25%"
                  _active={{
                    transform: 'scale(0.95)',
                    color: 'complimentary.800',
                  }}
                  _hover={{
                    bgColor: 'rgba(255,128,0, 0.25)',
                    color: 'complimentary.300',
                  }}
                >
                  Return
                </Button>
              </Box>
            </Box>
          )}
          <Text mt={'2'} fontSize={'sm'} fontWeight={'light'}>
            <>
              <Text as="span" color="complimentary.700">
                {selectedValidators.length}
              </Text>
              &nbsp;{'Validators Selected'}
            </>
          </Text>
        </ModalBody>
        <ModalFooter></ModalFooter>
      </ModalContent>
    </Modal>
  );
};
