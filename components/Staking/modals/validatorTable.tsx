import { Box, Table, TableCaption, Tbody, Td, Th, Thead, Tr, Flex, TableContainer } from '@chakra-ui/react';
import React from 'react';

import { ParsedValidator as Validator } from '@/utils';

export const ValidatorsTable: React.FC<{
  validators: Validator[];
  onValidatorClick: (validator: { name: string; operatorAddress: string }) => void;
  selectedValidators: { name: string; operatorAddress: string }[];
  searchTerm?: string;
}> = ({ validators, onValidatorClick, selectedValidators, searchTerm }) => {
  const [sortedValidators, setSortedValidators] = React.useState<Validator[]>([]);
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
      const matches = filteredValidators.filter((validator) => validator.name.toLowerCase().includes(searchTerm));

      const nonMatches = filteredValidators.filter((validator) => !validator.name.toLowerCase().includes(searchTerm));

      // Concatenate them so matches come first
      filteredValidators = [...matches, ...nonMatches];
    }

    if (searchTerm) {
      filteredValidators = validators.filter((validator) => validator.name.toLowerCase().includes(searchTerm));
    }

    switch (sortBy) {
      case 'moniker':
        filteredValidators.sort((a, b) => {
          let aMoniker = a.name || '';
          let bMoniker = b.name || '';
          return sortOrder === 'asc' ? aMoniker.localeCompare(bMoniker) : bMoniker.localeCompare(aMoniker);
        });
        break;
      case 'commission':
        filteredValidators.sort((a, b) => {
          let aRate = a.commission || '';
          let bRate = b.commission || '';
          return sortOrder === 'asc' ? parseFloat(aRate) - parseFloat(bRate) : parseFloat(bRate) - parseFloat(aRate);
        });
        break;
      default:
        break;
    }

    setSortedValidators(filteredValidators);
  }, [validators, searchTerm, sortBy, sortOrder]);

  return (
    <Box borderRadius={'6px'} maxH="xl" minH="sm">
      <Box
        borderRadius={'6px'}
        maxH="120px"
        minH="sm"
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
          <Table mb={2} border="1px solid rgba(255,128,0, 0.25)" variant="simple" height="lg">
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
                <Th border="1px solid rgba(255,128,0, 0.25)" color="white" fontSize={'16px'}>
                  Missed
                </Th>
                <Th border="1px solid rgba(255,128,0, 0.25)" color="white" fontSize={'16px'}>
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
                  onClick={() =>
                    onValidatorClick({
                      name: validator.name || '',
                      operatorAddress: validator.address || '',
                    })
                  }
                  backgroundColor={selectedValidators.some((v) => v.name === validator.name) ? 'rgba(255, 128, 0, 0.25)' : 'transparent'}
                >
                  <Td border="1px solid rgba(255,128,0, 0.25)" color="white">
                    {(validator.name.length || 0) > 20 ? validator.name.substring(0, 14) || '' + '...' : validator.name || ''}
                  </Td>
                  <Td border="1px solid rgba(255,128,0, 0.25)" color="white">
                    {validator.commission ? validator.commission + '%' : 'N/A'}
                  </Td>
                  <Td border="1px solid rgba(255,128,0, 0.25)"></Td>
                  <Td border="1px solid rgba(255,128,0, 0.25)"></Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </TableContainer>
      </Box>
      <Flex width="100%" justifyContent="center" alignItems="center" mt={4} mb={2}></Flex>
    </Box>
  );
};
