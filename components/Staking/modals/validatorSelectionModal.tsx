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
} from '@chakra-ui/react';
import React, { useEffect } from 'react';
import { FaSearch } from 'react-icons/fa';

import { useValidatorData } from '@/hooks/useValidatorData';

import { ValidatorsTable } from './validatorTable';

interface MultiModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
  selectedChainName: string;
  selectedValidators: string[];
  setSelectedValidators: React.Dispatch<React.SetStateAction<string[]>>;
}

export const MultiModal: React.FC<MultiModalProps> = ({
  isOpen,
  onClose,
  selectedChainName,
  selectedValidators,
  setSelectedValidators,
}) => {
  const [searchTerm, setSearchTerm] = React.useState<string>('');

  const { data, isLoading, refetch } = useValidatorData(selectedChainName);

  useEffect(() => {
    if (isLoading) return;
    refetch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedChainName]);

  const validators = data?.validators;
  const handleValidatorClick = (validatorName: string) => {
    setSelectedValidators((prevState: string[]) => {
      // Check if selecting another validator would exceed the limit of 8
      if (!prevState.includes(validatorName) && prevState.length >= 8) {
        // Show a warning
        alert("You can't select more than 8 validators.");
        return prevState;
      }

      // If the validator is already selected, remove it, else add to the selected list
      return prevState.includes(validatorName)
        ? prevState.filter((v: string) => v !== validatorName)
        : [...prevState, validatorName];
    });
  };

  const handleQuickSelect = (count: number) => {
    if (!data || !validators) return;

    // Get the top N validators
    const topValidators = validators
      .slice(0, count)
      .map((validator) => validator.name);

    setSelectedValidators(topValidators);
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
                  <AccordionIcon mb={-4} color="complimentary.900" />
                </AccordionButton>
              </h2>
              <AccordionPanel
                textAlign="left"
                alignContent="center"
                justifyContent="center"
                mt={-2}
              >
                <Text
                  fontWeight="light"
                  pl={6}
                  maxW="95%"
                  color="white"
                  fontSize="16px"
                  letterSpacing={'wider'}
                >
                  Choose which validator(s) you would like to liquid stake to.
                  You can select from the list below or utilize the quick select
                  to pick the highest ranked validators. To learn more about
                  rainkings click here.
                </Text>
              </AccordionPanel>
            </AccordionItem>
          </Accordion>
        </ModalHeader>
        <ModalCloseButton color="white" />{' '}
        {/* Positioning by default should be top right */}
        <Divider
          bgColor="complimentary.900"
          alignSelf="center"
          w="88%"
          m="auto"
        />
        <ModalBody
          bgColor="#1A1A1A"
          borderRadius={'6px'}
          justifyContent="center"
        >
          {isLoading ? (
            <Box
              display="flex"
              justifyContent="center"
              alignItems="center"
              height="200px"
            >
              <Spinner h="50px" w="50px" color="complimentary.900" />
            </Box>
          ) : (
            <Box mt={-1}>
              <Flex
                py={2}
                px={4}
                alignContent="center"
                alignItems="center"
                justifyContent={'space-between'}
                w="100%"
                flexDirection={'row'}
              >
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
                      borderColor: 'complimentary.900',
                    }}
                    _selected={{
                      borderColor: 'complimentary.900',
                    }}
                    _hover={{
                      borderColor: 'complimentary.900',
                    }}
                    _focus={{
                      borderColor: 'complimentary.900',
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
                >
                  <Flex
                    w="100%"
                    h="50%"
                    pt={6}
                    justifyContent="space-between"
                    alignItems="center"
                    flexDir={'row'}
                  >
                    <Button
                      w="60px"
                      _hover={{
                        bgColor: '#181818',
                      }}
                      h="30px"
                      onClick={() => handleQuickSelect(2)}
                    >
                      Top 2
                    </Button>
                    <Button
                      w="60px"
                      _hover={{
                        bgColor: '#181818',
                      }}
                      h="30px"
                      onClick={() => handleQuickSelect(4)}
                    >
                      Top 4
                    </Button>
                    <Button
                      w="60px"
                      _hover={{
                        bgColor: '#181818',
                      }}
                      h="30px"
                      onClick={() => handleQuickSelect(8)}
                    >
                      Top 8
                    </Button>
                  </Flex>
                </Box>
              </Flex>
              <ValidatorsTable
                validators={validators || []}
                onValidatorClick={(validatorName) => {
                  if (
                    selectedValidators.length < 9 ||
                    selectedValidators.includes(validatorName)
                  ) {
                    handleValidatorClick(validatorName);
                  }
                }}
                selectedValidators={selectedValidators}
                searchTerm={searchTerm}
              />
              <Box
                mt={-12}
                w="100%"
                display="flex"
                justifyContent="center"
                alignItems="center"
              >
                <Button
                  onClick={onClose}
                  h="30px"
                  w="25%"
                  _hover={{
                    bgColor: '#181818',
                  }}
                >
                  Return
                </Button>
              </Box>
            </Box>
          )}
        </ModalBody>
        <ModalFooter></ModalFooter>
      </ModalContent>
    </Modal>
  );
};
