import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Divider,
  Tooltip,
  Text,
  Button,
  Table,
  Thead,
  Spinner,
  Image,
  Tbody,
  Tr,
  Th,
  Td,
  Icon,
  Flex,
  HStack,
  TableContainer,
  Menu,
  Fade,
  MenuButton,
  MenuItem,
  MenuList,
  Box,
  Checkbox,
} from '@chakra-ui/react';
import { useChain, useChains } from '@cosmos-kit/react';
import { SkipRouter, SKIP_API_URL, UserAddress } from '@skip-router/core';
import { ibc } from 'interchain-query';
import { useCallback, useState } from 'react';
import { FaInfoCircle } from 'react-icons/fa';

import { useTx } from '@/hooks';
import { useFeeEstimation } from '@/hooks/useFeeEstimation';
import { useAllBalancesQuery, useSkipAssets, useSkipReccomendedRoutes, useSkipRoutesData } from '@/hooks/useQueries';
import { useSkipExecute } from '@/hooks/useSkipExecute';
import { shiftDigits } from '@/utils';

const RewardsModal = ({ address, isOpen, onClose }: { address: string; isOpen: boolean; onClose: () => void }) => {
  const chains = useChains([
    'cosmoshub',
    'osmosis',
    'stargaze',
    'juno',
    'sommelier',
    'regen',
    'dydx',
    'saga',
    'stride',
    'noble',
    'neutron',
  ]);

  const { wallet } = useChain('quicksilver');

  const walletMapping: { [key: string]: any } = {
    Keplr: window.keplr?.getOfflineSignerOnlyAmino,
    Cosmostation: window.cosmostation?.providers?.keplr?.getOfflineSignerOnlyAmino,
    Leap: window.leap?.getOfflineSignerOnlyAmino,
  };

  const offlineSigner = wallet ? walletMapping[wallet.prettyName] : undefined;
  const skipClient = new SkipRouter({
    apiURL: SKIP_API_URL,
    getCosmosSigner: offlineSigner,
    endpointOptions: {
      endpoints: {
        'quicksilver-2': { rpc: 'https://rpc.quicksilver.zone/' },
      },
    },
  });
  const { balance, refetch } = useAllBalancesQuery('quicksilver', address);
  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isBottomVisible, setIsBottomVisible] = useState(true);
  const [selectedAssets, setSelectedAssets] = useState<string[]>([]);
  const [selectAll, setSelectAll] = useState<boolean>(false);

  const handleScroll = useCallback((event: React.UIEvent<HTMLDivElement>) => {
    const target = event.currentTarget;
    const isBottom = target.scrollHeight - target.scrollTop <= target.clientHeight;
    setIsBottomVisible(!isBottom);
  }, []);

  const handleCheckboxChange = (denom: string) => {
    setSelectedAssets((prevSelectedAssets) =>
      prevSelectedAssets.includes(denom) ? prevSelectedAssets.filter((item) => item !== denom) : [...prevSelectedAssets, denom],
    );
  };

  const handleSelectAllChange = () => {
    if (selectAll) {
      setSelectedAssets([]);
    } else {
      setSelectedAssets(tokenDetails.map((detail) => detail?.denom ?? ''));
    }
    setSelectAll(!selectAll);
  };

  const { assets: skipAssets } = useSkipAssets('quicksilver-2');
  const balanceData = balance?.balances || [];

  // maps through the balance query to get the token details and assigns them to the correct values for the skip router
  const getMappedTokenDetails = () => {
    const denomArray = skipAssets?.['quicksilver-2'] || [];

    return balanceData
      .map((balanceItem) => {
        // filter out native quick tokens
        if (balanceItem.denom.startsWith('q') || balanceItem.denom.startsWith('aq') || balanceItem.denom.startsWith('uq')) {
          return null;
        }

        const denomDetail = denomArray.find((denomItem) => denomItem.denom === balanceItem.denom);

        if (denomDetail) {
          return {
            amount: balanceItem.amount,
            denom: denomDetail.denom,
            originDenom: denomDetail.originDenom,
            originChainId: denomDetail.originChainID,
            trace: denomDetail.trace,
            logoURI: denomDetail.logoURI,
            decimals: denomDetail.decimals,
          };
        }

        return null;
      })
      .filter(Boolean);
  };

  const tokenDetails = getMappedTokenDetails();

  // maps through the token details to get the route objects for the skip router
  const osmosisRouteObjects = tokenDetails.map((token) => ({
    sourceDenom: token?.denom ?? '',
    sourceChainId: 'quicksilver-2',
    destChainId: 'osmosis-1',
  }));

  const { routes: osmosisRoutes } = useSkipReccomendedRoutes(osmosisRouteObjects);
  // maps through the token details and the route objects to get the specific token details for the skip router
  const osmosisRoutesDataObjects = tokenDetails
    .flatMap((token, index) => {
      return osmosisRoutes[index]?.flatMap((route) => {
        return route.recommendations.map((recommendation) => {
          return {
            amountIn: token?.amount ?? '0',
            sourceDenom: token?.denom ?? '',
            sourceChainId: 'quicksilver-2',
            destDenom: recommendation.asset.denom ?? '',
            destChainId: recommendation.asset.chainID ?? '',
          };
        });
      });
    })
    .filter(Boolean) as { amountIn: string; sourceDenom: string; sourceChainId: string; destDenom: string; destChainId: string }[];
  const { routesData } = useSkipRoutesData(osmosisRoutesDataObjects);

  const executeRoute = useSkipExecute(skipClient);

  // Helper function to get ordered addresses based on requiredChainAddresses
  const getOrderedAddresses = (requiredChainAddresses: string[], allAddresses: UserAddress[]): UserAddress[] => {
    return requiredChainAddresses.map((chainID) => {
      const addressObj = allAddresses.find((addr) => addr.chainID === chainID);
      if (!addressObj) {
        throw new Error(`Address for chainID ${chainID} not found`);
      }
      return addressObj;
    });
  };

  // Helper function to format the token name
  const formatTokenName = (originDenom: string) => {
    if (originDenom.toLowerCase().startsWith('st')) {
      if (originDenom.toLowerCase() === 'stinj') {
        return `st${originDenom.slice(2).toUpperCase()}`;
      }
      return `st${originDenom.slice(3).toUpperCase()}`;
    }
    return originDenom.toLowerCase().startsWith('factory/')
      ? originDenom.split('/').pop()?.replace(/^u/, '').toUpperCase()
      : originDenom.slice(1).toUpperCase();
  };

  // uses all the data gathered to create the ibc transactions for sending assets to osmosis.
  const handleExecuteRoute = async () => {
    setIsSigning(true);

    const allAddresses: UserAddress[] = [
      { chainID: 'quicksilver-2', address },
      { chainID: 'osmosis-1', address: chains.osmosis.address ?? '' },
      { chainID: 'cosmoshub-4', address: chains.cosmoshub.address ?? '' },
      { chainID: 'stargaze-1', address: chains.stargaze.address ?? '' },
      { chainID: 'sommelier-3', address: chains.sommelier.address ?? '' },
      { chainID: 'regen-1', address: chains.regen.address ?? '' },
      { chainID: 'juno-1', address: chains.juno.address ?? '' },
      { chainID: 'dydx-mainnet-1', address: chains.dydx.address ?? '' },
      { chainID: 'stride-1', address: chains.stride.address ?? '' },
      { chainID: 'noble-1', address: chains.noble.address ?? '' },
      { chainID: 'neutron-1', address: chains.neutron.address ?? '' },
    ];

    // Check if any address is undefined
    const undefinedAddresses = allAddresses.filter((addr) => !addr.address);
    if (undefinedAddresses.length > 0) {
      console.error('Some addresses are undefined:', undefinedAddresses);
      setIsSigning(false);
      return;
    }

    // Filter routesData to only include selected assets
    const filteredRoutesData = routesData.filter((route) => selectedAssets.includes(route?.sourceAssetDenom ?? ''));

    // Execute each route in sequence with ordered addresses
    for (const route of filteredRoutesData) {
      const orderedAddresses = getOrderedAddresses(route?.requiredChainAddresses ?? ([] as string[]), allAddresses);
      try {
        await executeRoute(route, orderedAddresses, refetch);
      } catch (error) {
        console.error('Error executing route:', error);
        setIsSigning(false);
        return;
      }
    }

    setIsSigning(false);
  };

  const { tx } = useTx('quicksilver');
  const { transfer } = ibc.applications.transfer.v1.MessageComposer.withTypeUrl;
  const { estimateFee } = useFeeEstimation('quicksilver');

  const onSubmitClick = async () => {
    setIsSigning(true);

    const messages = [];

    for (const tokenDetail of tokenDetails) {
      if (!tokenDetails) {
        setIsSigning(false);
        return;
      }

      if (!selectedAssets.includes(tokenDetail?.denom ?? '')) {
        continue;
      }

      const [_, channel] = tokenDetail?.trace.split('/') ?? '';
      const sourcePort = 'transfer';
      const sourceChannel = channel;
      const senderAddress = address ?? '';

      const ibcToken = {
        denom: tokenDetail?.denom ?? '',
        amount: tokenDetail?.amount ?? '0',
      };

      const stamp = Date.now();
      const timeoutInNanos = (stamp + 1.2e6) * 1e6;

      const chainIdToName: { [key: string]: string } = {
        'osmosis-1': 'osmosis',
        'cosmoshub-4': 'cosmoshub',
        'stargaze-1': 'stargaze',
        'sommelier-3': 'sommelier',
        'regen-1': 'regen',
        'juno-1': 'juno',
        'dydx-mainnet-1': 'dydx',
        'stride-1': 'stride',
        'noble-1': 'noble',
        'neutron-1': 'neutron',
      };

      const getChainName = (chainId: string) => {
        return chainIdToName[chainId] || chainId;
      };

      const chain = chains[getChainName(tokenDetail?.originChainId ?? '') ?? ''];
      const receiverAddress = chain?.address ?? '';

      const msg = transfer({
        sourcePort,
        sourceChannel,
        sender: senderAddress,
        receiver: receiverAddress,
        token: ibcToken,
        timeoutHeight: undefined,
        //@ts-ignore
        timeoutTimestamp: timeoutInNanos,
      });

      messages.push(msg);
    }

    try {
      const fee = await estimateFee(address, messages);

      await tx(messages, {
        fee,
        onSuccess: () => {
          setIsSigning(false);
        },
      });
      setIsSigning(false);
    } catch (error) {
      setIsSigning(false);
    }
  };

  const [destination, setDestination] = useState('');

  return (
    <Modal size={'sm'} isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent bgColor="rgb(32,32,32)">
        <ModalHeader color="white" fontSize="xl">
          <HStack>
            <Text>Rewards</Text>
            <Tooltip
              label="These are tokens from fees, staking rewards, and holding qAsset's. Select a direction then use the unwind button to bridge them."
              fontSize="md"
              placement="right"
            >
              <span>
                <Icon color={'complimentary.900'} as={FaInfoCircle} w={4} h={4} />
              </span>
            </Tooltip>
          </HStack>
          <Divider mt={3} bgColor={'cyan.500'} />
        </ModalHeader>
        <ModalCloseButton color={'complimentary.900'} />
        <ModalBody>
          {tokenDetails.length === 0 && (
            <Text color="white" textAlign={'center'} fontSize="md">
              No rewards available to claim
            </Text>
          )}
          {tokenDetails.length >= 1 && (
            <>
              <Box position="relative">
                <TableContainer className="custom-scrollbar" maxH="260px" overflowY="auto" onScroll={handleScroll} position="relative">
                  <Table variant="simple" colorScheme="whiteAlpha" size="sm">
                    <Thead position="sticky" top={0} bg="rgb(32,32,32)" zIndex={1}>
                      <Tr>
                        <Th color="complimentary.900">
                          <Checkbox
                            isChecked={selectAll}
                            onChange={handleSelectAllChange}
                            colorScheme="complimentary"
                            borderColor="complimentary.900"
                          />
                        </Th>
                        <Th color="complimentary.900">Token</Th>
                        <Th color="complimentary.900" isNumeric>
                          Amount
                        </Th>
                      </Tr>
                    </Thead>
                    <Tbody overflowY={'auto'} maxH="250px">
                      {tokenDetails.map((detail, index) => (
                        <Tr key={index}>
                          <Td color="white">
                            <Checkbox
                              isChecked={selectedAssets.includes(detail?.denom ?? '')}
                              onChange={() => handleCheckboxChange(detail?.denom ?? '')}
                              colorScheme="complimentary"
                              borderColor="complimentary.900"
                            />
                          </Td>
                          <Td color="white">
                            <HStack>
                              <Image w="32px" h="32px" alt={detail?.originDenom} src={detail?.logoURI} />
                              <Text fontSize={'large'}>{formatTokenName(detail?.originDenom ?? '')}</Text>
                            </HStack>
                          </Td>
                          <Td fontSize={'large'} color="white" isNumeric>
                            {Number(shiftDigits(detail?.amount ?? '', -Number(detail?.decimals)))
                              .toFixed(2)
                              .toString()}
                          </Td>
                        </Tr>
                      ))}
                    </Tbody>
                  </Table>
                </TableContainer>
                {isBottomVisible && tokenDetails.length > 6 && (
                  <Fade in={isBottomVisible}>
                    <Box
                      position="absolute"
                      bottom="0"
                      left="0"
                      right="0"
                      height="70px"
                      bgGradient="linear(to top, #1A1A1A, transparent)"
                      pointerEvents="none"
                      zIndex="1"
                    />
                  </Fade>
                )}
              </Box>
              <Flex justifyContent={'center'}>
                <Button
                  mt={4}
                  _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
                  _hover={{ bgColor: 'rgba(255,128,0, 0.25)', color: 'complimentary.300' }}
                  color="white"
                  size="sm"
                  w="160px"
                  variant="outline"
                  onClick={() => (destination === 'osmosis' ? handleExecuteRoute() : onSubmitClick())}
                >
                  {isSigning === true && <Spinner size="sm" />}
                  {isSigning === false && 'Unwind'}
                </Button>
              </Flex>
              <Flex justifyContent={'center'} mt={4}>
                <Menu>
                  <MenuButton
                    as={Button}
                    _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
                    _hover={{ bgColor: 'rgba(255,128,0, 0.25)', color: 'complimentary.300' }}
                    color="white"
                    size="sm"
                    w="130px"
                    bgColor={'rgb(26,26,26)'}
                  >
                    {destination === 'parentChains' ? 'Parent Chains' : destination === 'osmosis' ? 'Osmosis' : 'Direction'}
                  </MenuButton>
                  <MenuList maxW="130px" minWidth="130px" bgColor="#1a1a1a" color="white">
                    <MenuItem
                      onClick={() => setDestination('parentChains')}
                      bgColor={'rgb(26,26,26)'}
                      _hover={{ bg: 'complimentary.400' }}
                      _focus={{ bg: '#2a2a2a' }}
                    >
                      Parent Chains
                    </MenuItem>
                    <MenuItem
                      onClick={() => setDestination('osmosis')}
                      bgColor={'rgb(26,26,26)'}
                      _hover={{ bg: 'complimentary.400' }}
                      _focus={{ bg: '#2a2a2a' }}
                    >
                      Osmosis
                    </MenuItem>
                  </MenuList>
                </Menu>
              </Flex>
            </>
          )}
        </ModalBody>

        <ModalFooter></ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default RewardsModal;
