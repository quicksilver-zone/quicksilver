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
  Flex,
  Button,
  Spacer,
  Input,
  Spinner,
  HStack,
  VStack,
  InputGroup,
  InputLeftElement,
  TableContainer,
} from '@chakra-ui/react';
import { Icon } from '@interchain-ui/react';
import { InputIcon } from '@radix-ui/react-icons';
import React, { useEffect } from 'react';
import { FaSearch } from 'react-icons/fa';

import { useValidatorData } from '@/hooks/useValidatorData';
import { ParsedValidator as Validator } from '@/utils';

export const ValidatorsTable: React.FC<{
  validators: Validator[];
  onValidatorClick: (validatorName: string) => void;
  selectedValidators: string[];
  searchTerm?: string;
}> = ({ validators, onValidatorClick, selectedValidators, searchTerm }) => {
  const [sortedValidators, setSortedValidators] = React.useState<Validator[]>(
    [],
  );
  const [sortBy, setSortBy] = React.useState<string | null>(null);
  const [sortOrder, setSortOrder] = React.useState<'asc' | 'desc'>('asc');

  const handleSort = (column: string) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('asc');
    }
  };

  React.useEffect(() => {
    let filteredValidators = [...validators];

    if (searchTerm) {
      // Split into two arrays: matches and non-matches
      const matches = filteredValidators.filter((validator) =>
        validator.name.toLowerCase().includes(searchTerm),
      );

      const nonMatches = filteredValidators.filter(
        (validator) => !validator.name.toLowerCase().includes(searchTerm),
      );

      // Concatenate them so matches come first
      filteredValidators = [...matches, ...nonMatches];
    }

    if (searchTerm) {
      filteredValidators = validators.filter((validator) =>
        validator.name.toLowerCase().includes(searchTerm),
      );
    }

    switch (sortBy) {
      case 'moniker':
        filteredValidators.sort((a, b) => {
          let aMoniker = a.name || '';
          let bMoniker = b.name || '';
          return sortOrder === 'asc'
            ? aMoniker.localeCompare(bMoniker)
            : bMoniker.localeCompare(aMoniker);
        });
        break;
      case 'commission':
        filteredValidators.sort((a, b) => {
          let aRate = a.commission || '0';
          let bRate = b.commission || '0';
          return sortOrder === 'asc'
            ? parseFloat(aRate) - parseFloat(bRate)
            : parseFloat(bRate) - parseFloat(aRate);
        });
        break;
      default:
        break;
    }

    setSortedValidators(filteredValidators);
  }, [validators, searchTerm, sortBy, sortOrder]);

  return (
    <Box borderRadius={'6px'} maxH="xl" minH="lg">
      <Box
        borderRadius={'6px'}
        maxH="120px"
        minH="md"
        px={4}
        pb={0}
        overflowY="scroll"
        sx={{
          '&::-webkit-scrollbar': {
            width: '8px',
          },
          '&::-webkit-scrollbar-thumb': {
            backgroundColor: 'complimentary.900',
            borderRadius: '4px',
          },
          '&::-webkit-scrollbar-track': {
            backgroundColor: 'rgba(255,128,0, 0.25)',
            borderRadius: '10px',
          },
        }}
      >
        <TableContainer>
          <Table
            mb={2}
            border="1px solid rgba(255,128,0, 0.25)"
            variant="simple"
            height="lg"
          >
            <TableCaption>All validators</TableCaption>
            <Thead>
              <Tr>
                <Th
                  border="1px solid rgba(255,128,0, 0.25)"
                  color="white"
                  fontSize={'16px'}
                  onClick={() => handleSort('moniker')}
                  _hover={{
                    backgroundColor: 'rgba(255,128,0, 0.25)',
                    cursor: 'pointer',
                  }}
                >
                  Moniker
                </Th>
                <Th
                  border="1px solid rgba(255,128,0, 0.25)"
                  color="white"
                  fontSize={'16px'}
                  onClick={() => handleSort('commission')}
                  _hover={{
                    backgroundColor: 'rgba(255,128,0, 0.25)',
                    cursor: 'pointer',
                  }}
                >
                  Commission
                </Th>
                <Th
                  border="1px solid rgba(255,128,0, 0.25)"
                  color="white"
                  fontSize={'16px'}
                >
                  Missed
                </Th>
                <Th
                  border="1px solid rgba(255,128,0, 0.25)"
                  color="white"
                  fontSize={'16px'}
                >
                  Rank
                </Th>
              </Tr>
            </Thead>
            <Tbody borderRadius={'10px'}>
              {sortedValidators.map((validator, index) => (
                <Tr
                  cursor="pointer"
                  key={index}
                  _hover={{
                    bgColor: 'rgba(255,128,0, 0.1)',
                  }}
                  onClick={() => onValidatorClick(validator.name || '')} // Add click handler
                  backgroundColor={
                    selectedValidators.includes(validator.name || '')
                      ? 'rgba(255, 128, 0, 0.25)'
                      : 'transparent'
                  } // Change background color if selected
                >
                  <Td border="1px solid rgba(255,128,0, 0.25)" color="white">
                    {(validator.name.length || 0) > 20
                      ? validator.name.substring(0, 14) || '' + '...'
                      : validator.name || ''}
                  </Td>
                  <Td border="1px solid rgba(255,128,0, 0.25)" color="white">
                    {validator.commission
                      ? (
                          (parseFloat(validator.commission || '') / 1e18) *
                          100
                        ).toFixed(2) + '%'
                      : 'N/A'}
                  </Td>
                  <Td border="1px solid rgba(255,128,0, 0.25)"></Td>
                  <Td border="1px solid rgba(255,128,0, 0.25)"></Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </TableContainer>
      </Box>
      <Flex
        width="100%"
        justifyContent="center"
        alignItems="center"
        mt={4}
        mb={2}
      ></Flex>
    </Box>
  );
};

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
    // If the validator is already selected, remove it, else add to the selected list
    setSelectedValidators((prevState: string[]) =>
      prevState.includes(validatorName)
        ? prevState.filter((v: string) => v !== validatorName)
        : [...prevState, validatorName],
    );
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
                  <InputLeftElement
                    pointerEvents="none" // Makes the icon non-clickable and allows the input to be focused when clicking on the icon
                  >
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
                onValidatorClick={handleValidatorClick}
                selectedValidators={selectedValidators}
                searchTerm={searchTerm}
              />
              <Box w="100%" justifyContent={'center'} alignItems={'center'}>
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
