import { Box, Button, ButtonGroup, Container, Flex, HStack, SlideFade, Spacer, Spinner, Text } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import axios from 'axios';
import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useEffect, useMemo, useState } from 'react';

import AssetsGrid from '@/components/Assets/assetsGrid';
import StakingIntent from '@/components/Assets/intents';
import MyPortfolio from '@/components/Assets/portfolio';
import QuickBox from '@/components/Assets/quickbox';
import RewardsClaim from '@/components/Assets/rewardsClaim';
import UnbondingAssetsTable from '@/components/Assets/unbondingTable';
import {
  useAPYQuery,
  useBalanceQuery,
  useIntentQuery,
  useLiquidRewardsQuery,
  useQBalanceQuery,
  useTokenPriceQuery,
  useZoneQuery,
} from '@/hooks/useQueries';
import { shiftDigits, toNumber } from '@/utils';

export interface PortfolioItemInterface {
  title: string;
  percentage: string;
  progressBarColor: string;
  amount: string;
  qTokenPrice: number;
}

type NumericRedemptionRates = {
  [key: string]: number;
};

type BalanceRates = {
  [key: string]: string;
};

type APYRates = {
  [key: string]: Number;
};

function Home() {
  const { address, isWalletConnected } = useChain('quicksilver');
  // Function to fetch token price from an API
  const fetchTokenPrice = async (token: any) => {
    try {
      const response = await axios.get(`https://api-osmosis.imperator.co/tokens/v2/price/${token}`);
      return response.data.price;
    } catch (error) {
      console.error('Error fetching token price:', error);
      return null;
    }
  };

  const COSMOSHUB_CHAIN_ID = process.env.NEXT_PUBLIC_COSMOSHUB_CHAIN_ID;
  const OSMOSIS_CHAIN_ID = process.env.NEXT_PUBLIC_OSMOSIS_CHAIN_ID;
  const STARGAZE_CHAIN_ID = process.env.NEXT_PUBLIC_STARGAZE_CHAIN_ID;
  const REGEN_CHAIN_ID = process.env.NEXT_PUBLIC_REGEN_CHAIN_ID;
  const SOMMELIER_CHAIN_ID = process.env.NEXT_PUBLIC_SOMMELIER_CHAIN_ID;

  // Retrieve balance for each token
  const { balance: qAtom, isLoading: qAtomIsLoading, isError: qAtomIsError } = useQBalanceQuery('quicksilver', address ?? '', 'atom');
  const { balance: qOsmo, isLoading: qOsmoIsLoading, isError: qOsmoIsError } = useQBalanceQuery('quicksilver', address ?? '', 'osmo');
  const { balance: qStars, isLoading: qStarsIsLoading, isError: qStarsIsError } = useQBalanceQuery('quicksilver', address ?? '', 'stars');
  const { balance: qRegen, isLoading: qRegenIsLoading, isError: qRegenIsError } = useQBalanceQuery('quicksilver', address ?? '', 'regen');
  const { balance: qSomm, isLoading: qSommIsLoading, isError: qSommIsError } = useQBalanceQuery('quicksilver', address ?? '', 'somm');

  // Retrieve zone data for each token
  const { data: CosmosZone, isLoading: isCosmosZoneLoading, isError: isCosmosZoneError } = useZoneQuery(COSMOSHUB_CHAIN_ID ?? '');
  const { data: OsmoZone, isLoading: isOsmoZoneLoading, isError: isOsmoZoneError } = useZoneQuery(OSMOSIS_CHAIN_ID ?? '');
  const { data: StarZone, isLoading: isStarZoneLoading, isError: isStarZoneError } = useZoneQuery(STARGAZE_CHAIN_ID ?? '');
  const { data: RegenZone, isLoading: isRegenZoneLoading, isError: isRegenZoneError } = useZoneQuery(REGEN_CHAIN_ID ?? '');
  const { data: SommZone, isLoading: isSommZoneLoading, isError: isSommZoneError } = useZoneQuery(SOMMELIER_CHAIN_ID ?? '');
  // Retrieve APY data for each token
  const { APY: cosmosAPY, isLoading: isCosmosAPYLoading, isError: isCosmosAPYError } = useAPYQuery('cosmoshub-4');
  const { APY: osmoAPY, isLoading: isOsmoAPYLoading, isError: isOsmoAPYError } = useAPYQuery('osmosis-1');
  const { APY: starsAPY, isLoading: isStarsAPYLoading, isError: isStarsAPYError } = useAPYQuery('stargaze-1');
  const { APY: regenAPY, isLoading: isRegenAPYLoading, isError: isRegenAPYError } = useAPYQuery('regen-1');
  const { APY: sommAPY, isLoading: isSommAPYLoading, isError: isSommAPYError } = useAPYQuery('sommelier-3');
  const { APY: quickAPY, isLoading: isQuickAPYLoading, isError: isQuickAPYError } = useAPYQuery('quicksilver-2');

  // useMemo hook to cache APY data
  const qAPYRates: APYRates = useMemo(
    () => ({
      qAtom: cosmosAPY,
      qOsmo: osmoAPY,
      qStars: starsAPY,
      qRegen: regenAPY,
      qSomm: sommAPY,
    }),
    [cosmosAPY, osmoAPY, starsAPY, regenAPY, sommAPY],
  );
  // useMemo hook to cache qBalance data
  const qBalances: BalanceRates = useMemo(
    () => ({
      qAtom: shiftDigits(qAtom?.balance?.amount ?? '', -6),
      qOsmo: shiftDigits(qOsmo?.balance?.amount ?? '', -6),
      qStars: shiftDigits(qStars?.balance?.amount ?? '', -6),
      qRegen: shiftDigits(qRegen?.balance?.amount ?? '', -6),
      qSomm: shiftDigits(qSomm?.balance?.amount ?? '', -6),
    }),
    [qAtom, qOsmo, qStars, qRegen, qSomm],
  );

  // useMemo hook to cache redemption rate data
  const redemptionRates: NumericRedemptionRates = useMemo(
    () => ({
      atom: CosmosZone?.redemptionRate ? parseFloat(CosmosZone.redemptionRate) : 1,
      osmo: OsmoZone?.redemptionRate ? parseFloat(OsmoZone.redemptionRate) : 1,
      stars: StarZone?.redemptionRate ? parseFloat(StarZone.redemptionRate) : 1,
      regen: RegenZone?.redemptionRate ? parseFloat(RegenZone.redemptionRate) : 1,
      somm: SommZone?.redemptionRate ? parseFloat(SommZone.redemptionRate) : 1,
    }),
    [CosmosZone, OsmoZone, StarZone, RegenZone, SommZone],
  );

  // State hooks for portfolio items, total portfolio value, and other metrics
  const [portfolioItems, setPortfolioItems] = useState<PortfolioItemInterface[]>([]);
  const [totalPortfolioValue, setTotalPortfolioValue] = useState(0);
  const [averageApy, setAverageAPY] = useState(0);
  const [totalYearlyYield, setTotalYearlyYield] = useState(0);
  // useEffect hook to compute portfolio metrics when dependencies change
  useEffect(() => {
    const updatePortfolioItems = async () => {
      let totalValue = 0;
      let totalYearlyYield = 0;
      let weightedAPY = 0;
      let updatedItems = [];
      // Loop through each token to compute value, APY, and yield
      for (const token of Object.keys(qBalances)) {
        const baseToken = token.replace('q', '').toLowerCase();
        const price = await fetchTokenPrice(baseToken);
        const qTokenPrice = price * Number(redemptionRates[baseToken]);
        const qTokenBalance = qBalances[token];
        const itemValue = Number(qTokenBalance) * qTokenPrice;

        const qTokenAPY = qAPYRates[token] || 0;
        const yearlyYield = itemValue * Number(qTokenAPY);
        // Accumulate total values and compute weighted APY
        totalValue += itemValue;
        totalYearlyYield += yearlyYield;
        weightedAPY += (itemValue / totalValue) * Number(qTokenAPY);

        updatedItems.push({
          title: token.toUpperCase(),
          percentage: 0,
          progressBarColor: 'complimentary.700',
          amount: qTokenBalance,
          qTokenPrice: qTokenPrice || 0,
        });
      }
      // Recalculate percentages for each item based on total value
      updatedItems = updatedItems.map((item) => {
        const itemValue = Number(item.amount) * item.qTokenPrice;
        return {
          ...item,
          percentage: (((itemValue / totalValue) * 100) / 100).toFixed(2),
        };
      });
      // Update state with calculated data
      setPortfolioItems(updatedItems);
      setTotalPortfolioValue(totalValue);
      setAverageAPY(weightedAPY);
      setTotalYearlyYield(totalYearlyYield);
    };

    updatePortfolioItems();
  }, [qBalances, CosmosZone, OsmoZone, StarZone, RegenZone, SommZone, redemptionRates, qAPYRates]);

  const assetsData = useMemo(() => {
    return Object.keys(qBalances).map((token) => {
      return {
        name: token.toUpperCase().replace('Q', 'q'),
        balance: toNumber(qBalances[token], 2).toString(),
        apy: parseFloat(qAPYRates[token]?.toFixed(2)) || 0,
        native: token.replace('q', '').toUpperCase(),
      };
    });
  }, [qBalances, qAPYRates]);

  const { liquidRewards, isLoading } = useLiquidRewardsQuery(address ?? '');

  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container
          flexDir={'column'}
          top={20}
          mt={{ base: 10, md: 10 }}
          zIndex={2}
          position="relative"
          justifyContent="center"
          alignItems="center"
          maxW="6xl"
        >
          <Head>
            <title>Assets</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <link rel="icon" href="/quicksilver/img/favicon.png" />
          </Head>
          <Text pb={2} color="white" fontSize="24px">
            Assets
          </Text>

          <Flex flexDir={{ base: 'column', md: 'row' }} py={6} alignItems="center" justifyContent={'space-between'} gap="4">
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
                  Wallet is not connected! Please connect your wallet to view the Assets Section.
                </Text>
              </Flex>
            )}
            {isWalletConnected && (
              <>
                <Flex
                  position="relative"
                  backdropFilter="blur(50px)"
                  bgColor="rgba(255,255,255,0.1)"
                  borderRadius="10px"
                  p={5}
                  w={{ base: 'full', md: 'sm' }}
                  h="sm"
                  flexDir="column"
                  justifyContent="space-around"
                  alignItems="center"
                >
                  <QuickBox stakingApy={quickAPY} />
                </Flex>

                <Flex
                  alignContent={'center'}
                  position="relative"
                  backdropFilter="blur(50px)"
                  bgColor="rgba(255,255,255,0.1)"
                  borderRadius="10px"
                  p={5}
                  w={{ base: 'full', md: '2xl' }}
                  h="sm"
                >
                  <MyPortfolio
                    portfolioItems={portfolioItems}
                    isWalletConnected={isWalletConnected}
                    totalValue={totalPortfolioValue}
                    averageApy={averageApy}
                    totalYearlyYield={totalYearlyYield}
                  />
                </Flex>
                <Flex
                  alignContent={'center'}
                  position="relative"
                  backdropFilter="blur(50px)"
                  bgColor="rgba(255,255,255,0.1)"
                  borderRadius="10px"
                  p={5}
                  w={{ base: 'full', md: 'lg' }}
                  h="sm"
                >
                  <StakingIntent isWalletConnected={isWalletConnected} address={address ?? ''} />
                </Flex>
              </>
            )}
          </Flex>
          <Spacer />
          <RewardsClaim address={address ?? ''} />
          <Spacer />
          {/* Assets Grid */}
          <AssetsGrid nonNative={liquidRewards} isWalletConnected={isWalletConnected} assets={assetsData} />
          <Spacer />
          {/* Unbonding Table */}
          <Box h="full" w="full" mt="20px">
            <UnbondingAssetsTable isWalletConnected={isWalletConnected} address={address ?? ''} />
          </Box>
          <Box h="40px"></Box>
        </Container>
      </SlideFade>
    </>
  );
}
// disable ssr in order to use useQuery hooks
const DynamicAssetsPage = dynamic(() => Promise.resolve(Home), {
  ssr: false,
});

const AssetsWrapper = () => {
  return <DynamicAssetsPage />;
};

export default AssetsWrapper;
