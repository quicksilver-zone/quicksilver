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
} from '@chakra-ui/react';
import { Spinner } from '@interchain-ui/react';
import React, { useEffect } from 'react';

import { useStakingData } from '@/hooks/useStakingData';
import { ExtendedValidator as Validator } from '@/utils';

export const ValidatorsTable: React.FC<{
  validators: Validator[];
  onValidatorClick: (validatorName: string) => void;
  selectedValidators: string[];
}> = ({ validators, onValidatorClick, selectedValidators }) => {
  const [searchTerm, setSearchTerm] = React.useState<string>('');
  const [sortedValidators, setSortedValidators] = React.useState<Validator[]>(
    [],
  );
  const [sortBy, setSortBy] = React.useState<string | null>(null);
  const [sortOrder, setSortOrder] = React.useState<'asc' | 'desc'>('asc');

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(event.target.value.toLowerCase());
  };

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
    <Box borderRadius={'6px'} maxH="xl" minH="xl">
      <Box
        borderRadius={'6px'}
        maxH="lg"
        minH="lg"
        px={4}
        overflowX="scroll"
        sx={{
          '&::-webkit-scrollbar': {
            width: '8px',
          },
          '&::-webkit-scrollbar-thumb': {
            backgroundColor: 'complimentary.900',
            borderRadius: '4px',
          },
          '&::-webkit-scrollbar-track': {
            backgroundColor: 'primary.900',
          },
          scrollbarWidth: 'thin',
          scrollbarColor: 'complimentary.900 primary.900',
        }}
      >
        <Table
          border="1px solid rgba(255,128,0, 0.25)"
          variant="simple"
          size="md"
        >
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
          <Tbody>
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
      </Box>
      <Flex justifyContent={'flex-end'} w="100%" flexDirection={'row'}>
        <Input
          type="text"
          mr="16px"
          mt={4}
          borderColor="complimentary.1000"
          placeholder="Validator Moniker"
          onChange={handleSearchChange}
          w="28%"
          mb={4}
          borderRadius={'4px'}
          _hover={{
            borderColor: 'complimentary.900',
          }}
          _selected={{
            borderColor: 'complimentary.900',
          }}
          _focus={{
            borderColor: 'complimentary.900',
          }}
        />
      </Flex>
    </Box>
  );
};

interface MultiModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
  selectedChainName: string;
}

export const MultiModal: React.FC<MultiModalProps> = ({
  isOpen,
  onClose,
  selectedChainName,
}) => {
  const [selectedValidators, setSelectedValidators] = React.useState<string[]>(
    [],
  );

  const { data, isLoading, refetch } = useStakingData(selectedChainName);

  useEffect(() => {
    if (isLoading) return;
    refetch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedChainName]);

  const validators = data?.validators;

  const handleValidatorClick = (validatorName: string) => {
    // If the validator is already selected, remove it, else add to the selected list
    setSelectedValidators((prevState) =>
      prevState.includes(validatorName)
        ? prevState.filter((v) => v !== validatorName)
        : [...prevState, validatorName],
    );
  };

  console.log('Multi stake Modal', selectedChainName);

  const handleQuickSelect = (count: number) => {
    if (!data || !validators) return;

    // Get the top N validators
    const topValidators = validators
      .slice(0, count)
      .map((validator) => validator.name);

    setSelectedValidators(topValidators);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="2xl">
      {/* Set the size here */}
      <ModalOverlay />
      <ModalContent borderRadius={'10px'} maxHeight="70vh" bgColor="#1A1A1A">
        <ModalHeader bgColor="#1A1A1A" p={0}>
          <Accordion allowToggle>
            <AccordionItem border="none">
              <h2>
                <AccordionButton
                  _hover={{
                    bgColor: 'transparent',
                  }}
                  p={6}
                >
                  <Box pr={4}>
                    <Text ml={4} color="white" fontSize="24px" textAlign="left">
                      Validator Selection
                    </Text>
                  </Box>
                  <AccordionIcon color="complimentary.900" />
                </AccordionButton>
              </h2>
              <AccordionPanel justifyContent="center" mt={-2}>
                <Box>
                  <Text
                    maxW="85%"
                    color="white"
                    fontSize="16px"
                    letterSpacing={'wider'}
                  >
                    Choose which validator(s) you would like to liquid stake to.
                    You can select from the list below or utilize the quick
                    select to pick the highest ranked validators. To learn more
                    about raninkings click here.
                  </Text>
                </Box>
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
              <Spinner size="large" color="complimentary.900" />
            </Box>
          ) : (
            <Box mt={4}>
              <ValidatorsTable
                validators={validators || []}
                onValidatorClick={handleValidatorClick}
                selectedValidators={selectedValidators}
              />
            </Box>
          )}
          <Box
            mt={4}
            borderRadius="10px"
            w="100%"
            h="100px"
            display="flex"
            justifyContent="space-between"
            alignItems="center"
          >
            <Flex
              flexDirection="column"
              alignItems="center"
              justifyContent="center"
              w="300px"
            >
              <Button
                h="30px"
                w="150px"
                _hover={{
                  bgColor: '#181818',
                }}
              >
                Liquid Stake
              </Button>
            </Flex>

            <Box
              bg="rgba(255,255,255,0.1)"
              borderRadius="10px"
              w="300px"
              h="85px"
              mr={2}
              display="flex"
              flexDirection="column"
              justifyContent="space-between"
            >
              <Text
                ml={5}
                mt={1}
                fontSize="18"
                color="white"
                textDecoration="underline"
              >
                Quick Select
              </Text>
              <Flex
                w="100%"
                h="50%"
                pb={4}
                pr={4}
                pl={4}
                justifyContent="space-between"
                alignItems="center"
              >
                <Button
                  w="60px"
                  _hover={{
                    bgColor: '#181818',
                  }}
                  h="30px"
                  onClick={() => handleQuickSelect(5)}
                >
                  Top 5
                </Button>
                <Button
                  w="60px"
                  _hover={{
                    bgColor: '#181818',
                  }}
                  h="30px"
                  onClick={() => handleQuickSelect(10)}
                >
                  Top 10
                </Button>
                <Button
                  w="60px"
                  _hover={{
                    bgColor: '#181818',
                  }}
                  h="30px"
                  onClick={() => handleQuickSelect(20)}
                >
                  Top 20
                </Button>
              </Flex>
            </Box>
          </Box>
        </ModalBody>
        <ModalFooter></ModalFooter>
      </ModalContent>
    </Modal>
  );
};
