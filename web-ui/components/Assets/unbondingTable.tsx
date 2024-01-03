import { useUnbondingQuery } from '@/hooks/useQueries';
import { Table, Thead, Tbody, Tr, Th, Td, TableContainer, Text, Box, Flex } from '@chakra-ui/react';

interface UnbondingAssetsTableProps {
  address: string;
  chainName: string;
  isWalletConnected: boolean;
}

const UnbondingAssetsTable: React.FC<UnbondingAssetsTableProps> = ({ address, chainName, isWalletConnected }) => {
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
  const { unbondingData } = useUnbondingQuery(chainName, address);
  const noUnbondingAssets = isWalletConnected && unbondingData?.withdrawals.length === 0;
  return (
    <>
      <Text fontSize="xl" fontWeight="bold" color="white" mb={4}>
        Current Unbonding Assets
      </Text>
      {!isWalletConnected && (
        <Flex
          w="100%"
          backdropFilter="blur(50px)"
          bgColor="rgba(255,255,255,0.1)"
          h="sm"
          p={4}
          borderRadius="lg"
          flexDirection="column"
          justifyContent="center"
          alignItems="center"
          gap={6}
          color="white"
        >
          <Text fontSize="xl" textAlign="center">
            Wallet is not connected! Please connect your wallet to view your unbondng assets.
          </Text>
        </Flex>
      )}
      {noUnbondingAssets && (
        <Flex
          w="100%"
          backdropFilter="blur(50px)"
          bgColor="rgba(255,255,255,0.1)"
          h="sm"
          p={4}
          borderRadius="lg"
          flexDirection="column"
          justifyContent="center"
          alignItems="center"
          gap={6}
          color="white"
        >
          <Text fontSize="xl" textAlign="center">
            You have no unbonding assets.
          </Text>
        </Flex>
      )}
      {isWalletConnected && !noUnbondingAssets && (
        <Box bgColor="rgba(255,255,255,0.1)" p={4} borderRadius="lg">
          <TableContainer h={'200px'} overflowY={'auto'}>
            <Table variant="simple" color="white">
              <Thead boxShadow="0px 0.5px 0px 0px rgba(255,255,255,1)" position={'sticky'} bgColor="#1A1A1A" top="0" zIndex="0">
                <Tr>
                  <Th borderBottomColor={'transparent'} color="complimentary.900">
                    Asset
                  </Th>
                  <Th borderBottomColor={'transparent'} color="complimentary.900">
                    Status
                  </Th>
                  <Th borderBottomColor={'transparent'} color="complimentary.900">
                    Redemption Amount
                  </Th>
                  <Th borderBottomColor={'transparent'} color="complimentary.900">
                    Unstaked On
                  </Th>
                  <Th borderBottomColor={'transparent'} color="complimentary.900">
                    Completion Time
                  </Th>
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
      )}
    </>
  );
};

export default UnbondingAssetsTable;
