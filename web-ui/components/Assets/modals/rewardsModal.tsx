
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
} from '@chakra-ui/react';
import { useChains } from '@cosmos-kit/react';
import { ibc } from 'interchain-query';
import { useCallback, useState } from 'react';
import { FaInfoCircle } from 'react-icons/fa';

import { useTx } from '@/hooks';
import { useFeeEstimation } from '@/hooks/useFeeEstimation';
import { useAllBalancesQuery, useSkipAssets } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';

const RewardsModal = ({
  address,

  isOpen,
  onClose,
}: {
  address: string;

  isOpen: boolean;
  onClose: () => void;
}) => {
  const { balance } = useAllBalancesQuery('quicksilver', address);
  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isBottomVisible, setIsBottomVisible] = useState(true);

  const handleScroll = useCallback((event: React.UIEvent<HTMLDivElement>) => {
    const target = event.currentTarget;
    const isBottom = target.scrollHeight - target.scrollTop <= target.clientHeight;
    setIsBottomVisible(!isBottom);
  }, []);

  const chains = useChains(['cosmoshub', 'osmosis', 'stargaze', 'juno', 'sommelier', 'regen', 'dydx']);

  const { assets: skipAssets } = useSkipAssets('quicksilver-2');

  const balanceData = balance?.balances || [];

  const getMappedTokenDetails = () => {
    const denomArray = skipAssets?.['quicksilver-2'] || [];

    return balanceData
      .map((balanceItem) => {
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

      const [_, channel] = tokenDetail?.trace.split('/') ?? '';
      const sourcePort = 'transfer';
      const sourceChannel = destination === 'osmosis' ? 'channel-2' : channel;
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
      };

      const getChainName = (chainId: string) => {
        return chainIdToName[chainId] || chainId;
      };

      const chain = chains[getChainName(tokenDetail?.originChainId ?? '') ?? ''];
      const receiverAddress = destination === 'osmosis' ? chains.osmosis.address ?? '' : chain?.address ?? '';

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
                            <HStack>
                              <Image w="32px" h="32px" alt={detail?.originDenom} src={detail?.logoURI} />
                              <Text fontSize={'large'}>
                                {detail?.originDenom
                                  ? detail.originDenom.toLowerCase().startsWith('factory/')
                                    ? (() => {
                                        const lastSegment = detail.originDenom.split('/').pop() || '';
                                        return lastSegment.startsWith('u') ? lastSegment.slice(1).toUpperCase() : lastSegment.toUpperCase();
                                      })()
                                    : detail.originDenom.slice(1).toUpperCase()
                                  : ''}
                              </Text>
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
                  onClick={onSubmitClick}
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
                      _hover={{ bg: '#2a2a2a' }}
                      _focus={{ bg: '#2a2a2a' }}
                    >
                      Parent Chains
                    </MenuItem>
                    <MenuItem
                      onClick={() => setDestination('osmosis')}
                      bgColor={'rgb(26,26,26)'}
                      _hover={{ bg: '#2a2a2a' }}
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
