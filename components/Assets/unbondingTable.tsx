import { Table, Thead, Tbody, Tr, Th, Td, TableContainer, Text, Box } from '@chakra-ui/react';

const UnbondingAssetsTable = () => {
  // Mock data
  const unbondingAssets = [
    {
      asset: '10 ATOM',
      status: 'Processing',
      redemptionAmount: '10 ATOM',
      unstakedOn: '2023-01-01',
      completionTime: '2023-01-14',
    },
    {
      asset: '10 ATOM',
      status: 'Processing',
      redemptionAmount: '10 ATOM',
      unstakedOn: '2023-01-01',
      completionTime: '2023-01-14',
    },
    {
      asset: '10 ATOM',
      status: 'Processing',
      redemptionAmount: '10 ATOM',
      unstakedOn: '2023-01-01',
      completionTime: '2023-01-14',
    },
    {
      asset: '10 ATOM',
      status: 'Processing',
      redemptionAmount: '10 ATOM',
      unstakedOn: '2023-01-01',
      completionTime: '2023-01-14',
    },
    // Add more mock items
  ];

  return (
    <Box bgColor="rgba(255,255,255,0.1)" p={4} borderRadius="lg">
      <Text fontSize="xl" fontWeight="bold" color="white" mb={4}>
        Current Unbonding Assets
      </Text>

      <TableContainer>
        <Table variant="simple" color="white">
          <Thead>
            <Tr>
              <Th>Asset</Th>
              <Th>Status</Th>
              <Th>Redemption Amount</Th>
              <Th>Unstaked On</Th>
              <Th>Completion Time</Th>
            </Tr>
          </Thead>
          <Tbody>
            {unbondingAssets.map((asset, index) => (
              <Tr key={index}>
                <Td color="complementary.900">{asset.asset}</Td>
                <Td>{asset.status}</Td>
                <Td>{asset.redemptionAmount}</Td>
                <Td>{asset.unstakedOn}</Td>
                <Td>{asset.completionTime}</Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </TableContainer>
    </Box>
  );
};

export default UnbondingAssetsTable;
