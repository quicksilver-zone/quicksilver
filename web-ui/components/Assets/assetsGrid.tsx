import {
  Box,
  VStack,
  Text,
  Divider,
  HStack,
  Flex,
  Spinner,
  Button,
  useDisclosure,
  Stat,
  StatHelpText,
  StatLabel,
  StatNumber,
  SimpleGrid,
  Skeleton,
} from '@chakra-ui/react';
import React, { useEffect, useState } from 'react';

import { shiftDigits, formatQasset, formatNumber } from '@/utils';

import QDepositModal from './modals/qTokenDepositModal';
import QWithdrawModal from './modals/qTokenWithdrawlModal';


interface AssetCardProps {
  address: string;
  assetName: string;
  balance: string;
  apy: number;
  nativeAssetName: string;
  redemptionRates: string;
  isWalletConnected: boolean;
  nonNative: LiquidRewardsData | undefined;
  liquidRewards: LiquidRewardsData | undefined;
  refetch: () => void;
}

interface AssetGridProps {
  address: string;
  isWalletConnected: boolean;
  assets: Array<{
    name: string;
    balance: string;
    apy: number;
    native: string;
    redemptionRates: string;
  }>;
  nonNative: LiquidRewardsData | undefined;
  liquidRewards: LiquidRewardsData | undefined;
  refetch: () => void;
}

type Amount = {
  denom: string;
  amount: string;
};

type Errors = {
  Errors: any;
};

type LiquidRewardsData = {
  messages: any[];
  assets: {
    [key: string]: [
      {
        Type: string;
        Amount: Amount[];
      },
    ];
  };
  errors: Errors;
};

const AssetCard: React.FC<AssetCardProps> = ({ address, assetName, balance, apy, redemptionRates, liquidRewards, refetch }) => {
  const chainIdToName: { [key: string]: string } = {
    'osmosis-1': 'osmosis',
    'secret-1': 'secretnetwork',
    'umee-1': 'umee',
    'cosmoshub-4': 'cosmoshub',
    'stargaze-1': 'stargaze',
    'sommelier-3': 'sommelier',
    'regen-1': 'regen',
    'juno-1': 'juno',
    'dydx-mainnet-1': 'dydx',
    'ssc-1': 'saga',
  };

  const getChainName = (chainId: string) => {
    return chainIdToName[chainId] || chainId;
  };

  const convertAmount = (amount: string, denom: string) => {
    if (denom.startsWith('a')) {
      return shiftDigits(amount, -18);
    }

    return shiftDigits(amount, -6);
  };

  const [interchainDetails, setInterchainDetails] = useState({});

  useEffect(() => {
    const calculateInterchainBalance = () => {
      if (!liquidRewards || !liquidRewards.assets) return '0';

      let totalAmount = 0;
      const assetDenom = `uq${assetName.toLowerCase().replace('q', '')}`;
      const aAssetDenom = `aq${assetName.toLowerCase().replace('q', '')}`;

      const details: { [key: string]: number } = {};

      Object.keys(liquidRewards.assets).forEach((chainId) => {
        if (chainId !== 'quicksilver-2') {
          liquidRewards.assets[chainId].forEach((asset) => {
            asset.Amount.forEach((amount) => {
              if (amount.denom === assetDenom || amount.denom === aAssetDenom) {
                const convertedAmount = parseFloat(convertAmount(amount.amount, amount.denom));
                totalAmount += convertedAmount;
                details[getChainName(chainId)] = (details[getChainName(chainId)] || 0) + convertedAmount;
              }
            });
          });
        }
      });

      setInterchainDetails(details);
      return totalAmount.toString();
    };

    calculateInterchainBalance();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [liquidRewards, assetName]);

  const interchainBalance = Object.values(interchainDetails as { [key: string]: number })
    .reduce((acc: number, val: number) => acc + val, 0)
    .toString();

  const withdrawDisclosure = useDisclosure();
  const depositDisclosure = useDisclosure();

  if (balance === undefined || balance === null || apy === undefined || apy === null) {
    return (
      <Flex
        w="100%"
        h="100%"
        p={4}
        borderRadius="lg"
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        gap={6}
        color="white"
      >
        <Spinner w={'200px'} h="200px" color="complimentary.900" />
      </Flex>
    );
  }

  return (
    <VStack bg={'rgba(255,255,255,0.1)'} p={4} boxShadow="lg" align="center" spacing={4} borderRadius="lg" maxH="240px" minH="240px">
      <HStack w="full" justify="space-between">
        <Text fontWeight="bold" fontSize={'xl'} isTruncated>
          {assetName}
        </Text>
        <HStack>
          <Text fontSize="md" fontWeight="bold" isTruncated>
            {Number(shiftDigits(apy, 2))}%
          </Text>
          <Text fontSize="xs" fontWeight="light" isTruncated>
            APY
          </Text>
        </HStack>
      </HStack>
      <Divider bgColor={'complimentary.900'} />
      <HStack h="140px" justifyContent={'space-between'} w="full">
        <VStack minH="150px" alignItems="left">
          <Stat color={'white'}>
            <StatLabel fontSize={'lg'}>On Quicksilver</StatLabel>

            {!balance || !liquidRewards ? (
              <Skeleton startColor="complimentary.900" endColor="complimentary.100" height="10px" width="auto" />
            ) : (
              <StatNumber color={'complimentary.900'} fontSize={'md'}>
                {formatNumber(parseFloat(balance))} {assetName}
              </StatNumber>
            )}

            {!balance || !liquidRewards ? (
              <>
                <Skeleton startColor="complimentary.900" endColor="complimentary.100" height="10px" width="auto" mt={2} />
                <Skeleton startColor="complimentary.900" endColor="complimentary.100" height="10px" width="auto" mt={2} />
              </>
            ) : (
              Number(balance) > 0 && (
                <>
                  <StatHelpText mt={2} fontSize={'md'}>
                    Redeem For
                  </StatHelpText>
                  <StatHelpText mt={-2} color={'complimentary.400'} fontSize={'sm'}>
                    {formatNumber(parseFloat(balance) * Number(redemptionRates))} {assetName.replace('q', '')}
                  </StatHelpText>
                </>
              )
            )}
          </Stat>
          <Button
            _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
            _hover={{ bgColor: 'rgba(255,128,0, 0.25)', color: 'complimentary.300' }}
            color="white"
            size="sm"
            w="130px"
            variant="outline"
            onClick={withdrawDisclosure.onOpen}
            isDisabled={Number(balance) === 0}
          >
            Withdraw
          </Button>
          <QWithdrawModal
            refetch={refetch}
            max={balance}
            isOpen={withdrawDisclosure.isOpen}
            onClose={withdrawDisclosure.onClose}
            token={assetName}
          />
        </VStack>

        <VStack minH="150px" alignItems="left">
          <Stat color={'white'}>
            <StatLabel fontSize={'lg'}>Interchain</StatLabel>

            {!balance || !liquidRewards || !interchainBalance ? (
              <Skeleton startColor="complimentary.900" endColor="complimentary.100" height="10px" width="auto" />
            ) : (
              <StatNumber color={'complimentary.900'} fontSize={'md'}>
                {formatNumber(parseFloat(interchainBalance))} {assetName}
              </StatNumber>
            )}

            {!balance || !liquidRewards || !interchainBalance ? (
              <>
                <Skeleton startColor="complimentary.900" endColor="complimentary.100" height="10px" width="auto" mt={2} />
                <Skeleton startColor="complimentary.900" endColor="complimentary.100" height="10px" width="auto" mt={2} />
              </>
            ) : (
              Number(interchainBalance) > 0 && (
                <>
                  <StatHelpText mt={2} fontSize={'md'}>
                    Redeem For
                  </StatHelpText>
                  <StatHelpText mt={-2} color={'complimentary.400'} fontSize={'sm'}>
                    {formatNumber(parseFloat(interchainBalance) / Number(redemptionRates))} {assetName.replace('q', '')}
                  </StatHelpText>
                </>
              )
            )}
          </Stat>
          <Button
            _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
            _hover={{ bgColor: 'rgba(255,128,0, 0.25)', color: 'complimentary.300' }}
            color="white"
            size="sm"
            w="130px"
            variant="outline"
            onClick={depositDisclosure.onOpen}
            isDisabled={Number(interchainBalance) === 0}
          >
            Deposit
          </Button>
          <QDepositModal
            refetch={refetch}
            interchainDetails={interchainDetails}
            isOpen={depositDisclosure.isOpen}
            onClose={depositDisclosure.onClose}
            token={assetName}
          />
        </VStack>
      </HStack>
    </VStack>
  );
};

const AssetsGrid: React.FC<AssetGridProps> = ({ address, assets, isWalletConnected, nonNative, liquidRewards, refetch }) => {
  const scrollRef = React.useRef<HTMLDivElement>(null);
  const [focusedIndex, setFocusedIndex] = useState(0);

  const handleMouseEnter = (index: number) => {
    setFocusedIndex(index);
  };

  // const scrollByOne = (direction: 'left' | 'right') => {
  //   if (!scrollRef.current) return;

  //   const cardWidth = 380;
  //   let newIndex = focusedIndex;

  //   if (direction === 'left' && focusedIndex > 0) {
  //     scrollRef.current.scrollBy({ left: -cardWidth, behavior: 'smooth' });
  //     newIndex = focusedIndex - 1;
  //   } else if (direction === 'right' && focusedIndex < assets.length - 1) {
  //     scrollRef.current.scrollBy({ left: cardWidth, behavior: 'smooth' });
  //     newIndex = focusedIndex + 1;
  //   }

  //   setFocusedIndex(newIndex);
  // };

  // const getZoneName = (qAssetName: string) => {
  //   switch (qAssetName) {
  //     case 'QATOM':
  //       return 'Cosmos';
  //     case 'QOSMO':
  //       return 'Osmosis';
  //     case 'QSTARS':
  //       return 'Stargaze';
  //     case 'QSOMM':
  //       return 'Sommelier';
  //     case 'QREGEN':
  //       return 'Regen';
  //     case 'QJUNO':
  //       return 'Juno';
  //     case 'QDYDX':
  //       return 'DyDx';

  //     default:
  //       return qAssetName;
  //   }
  // };

  return (
    <>
      {/* Carousel controls and title */}
      <Flex justifyContent="space-between" alignItems="center" mb={4}>
        <Text fontSize="xl" fontWeight="bold" color="white">
          qAssets
        </Text>
        {/* <Flex alignItems="center" gap="2">
          <IconButton
            icon={<ChevronLeftIcon />}
            onClick={() => scrollByOne('left')}
            aria-label="Scroll left"
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
            _active={{ transform: 'scale(0.75)', color: 'complimentary.800' }}
            color="white"
            isDisabled={focusedIndex === 0}
            _disabled={{ cursor: 'default' }}
          />
          <Box minWidth="100px" textAlign="center">
            <Text fontSize="md" fontWeight="bold" color="white">
              {getZoneName(assets[focusedIndex]?.name)}
            </Text>
          </Box>
          <IconButton
            icon={<ChevronRightIcon />}
            onClick={() => scrollByOne('right')}
            aria-label="Scroll right"
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
            _active={{ transform: 'scale(0.75)', color: 'complimentary.800' }}
            color="white"
            isDisabled={focusedIndex === assets.length - 1}
            _disabled={{ cursor: 'default' }}
          />
        </Flex> */}
      </Flex>

      {/* Carousel content */}
      {!isWalletConnected ? (
        <Flex
          backdropFilter="blur(50px)"
          bgColor="rgba(255,255,255,0.1)"
          direction="column"
          p={5}
          borderRadius="lg"
          align="center"
          justify="space-around"
          w="full"
          h="200px"
        >
          <Text fontSize="xl" textAlign="center">
            Wallet is not connected! Please connect your wallet to interact with your qAssets.
          </Text>
        </Flex>
      ) : (
        <SimpleGrid columns={{ base: 1, md: 1, lg: 2, xl: 3 }} spacing={8} w="full" py={4} ref={scrollRef}>
          {assets?.map((asset, index) => (
            <Box
              key={index}
              minW="350px"
              transform={focusedIndex === index ? 'translateY(-10px)' : 'none'}
              transition="transform 0.1s"
              onMouseEnter={() => handleMouseEnter(index)}
            >
              <AssetCard
                address={address}
                isWalletConnected={isWalletConnected}
                assetName={formatQasset(asset.name)}
                nativeAssetName={asset.native}
                balance={asset.balance}
                apy={asset.apy}
                nonNative={nonNative}
                redemptionRates={asset.redemptionRates}
                liquidRewards={liquidRewards}
                refetch={refetch}
              />
            </Box>
          ))}
        </SimpleGrid>
      )}
    </>
  );
};

export default AssetsGrid;
