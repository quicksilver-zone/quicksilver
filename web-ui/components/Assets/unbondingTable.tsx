import { Table, Thead, Tbody, Tr, Th, Td, TableContainer, Text, Box } from '@chakra-ui/react';

const UnbondingAssetsTable = () => {
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
  ];

  return (
    <>
      <Text fontSize="xl" fontWeight="bold" color="white" mb={4}>
        Current Unbonding Assets
      </Text>
      <Box bgColor="rgba(255,255,255,0.1)" p={4} borderRadius="lg">
        <TableContainer h={'200px'} overflowY={'auto'}>
          <Table variant="simple" color="white">
            <Thead>
              <Tr>
                <Th color="complimentary.900">Asset</Th>
                <Th color="complimentary.900">Status</Th>
                <Th color="complimentary.900">Redemption Amount</Th>
                <Th color="complimentary.900">Unstaked On</Th>
                <Th color="complimentary.900">Completion Time</Th>
              </Tr>
            </Thead>
            <Tbody>
              {unbondingAssets.map((asset, index) => (
                <Tr key={index}>
                  <Td borderBottomColor={'transparent'} color="complementary.900">
                    {asset.asset}
                  </Td>
                  <Td borderBottomColor={'transparent'}>{asset.status}</Td>
                  <Td borderBottomColor={'transparent'}>{asset.redemptionAmount}</Td>
                  <Td borderBottomColor={'transparent'}>{asset.unstakedOn}</Td>
                  <Td borderBottomColor={'transparent'}>{asset.completionTime}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </TableContainer>
      </Box>
    </>
  );
};

export default UnbondingAssetsTable;
