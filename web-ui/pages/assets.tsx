import { Box, Container, Flex, SlideFade, Spacer, Text, Image } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useEffect, useMemo, useState } from 'react';

import AssetsGrid from '@/components/Assets/assetsGrid';
import StakingIntent from '@/components/Assets/intents';
import MyPortfolio from '@/components/Assets/portfolio';
import QuickBox from '@/components/Assets/quickbox';
import RewardsClaim from '@/components/Assets/rewardsClaim';
import UnbondingAssetsTable from '@/components/Assets/unbondingTable';
import { useAPYQuery, useAuthChecker, useLiquidRewardsQuery, useQBalanceQuery, useTokenPrices, useZoneQuery } from '@/hooks/useQueries';
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
  const { address } = useChain('quicksilver');
  const tokens = ['atom', 'osmo', 'stars', 'regen', 'somm', 'juno']; // Example tokens

  const { data: tokenPrices, isLoading: isLoadingPrices } = useTokenPrices(tokens);

  const COSMOSHUB_CHAIN_ID = process.env.NEXT_PUBLIC_COSMOSHUB_CHAIN_ID;
  const OSMOSIS_CHAIN_ID = process.env.NEXT_PUBLIC_OSMOSIS_CHAIN_ID;
  const STARGAZE_CHAIN_ID = process.env.NEXT_PUBLIC_STARGAZE_CHAIN_ID;
  const REGEN_CHAIN_ID = process.env.NEXT_PUBLIC_REGEN_CHAIN_ID;
  const SOMMELIER_CHAIN_ID = process.env.NEXT_PUBLIC_SOMMELIER_CHAIN_ID;
  const JUNO_CHAIN_ID = process.env.NEXT_PUBLIC_JUNO_CHAIN_ID;

  // Retrieve balance for each token
  const { balance: qAtom, isLoading: isLoadingQABalance } = useQBalanceQuery('quicksilver', address ?? '', 'atom');
  const { balance: qOsmo, isLoading: isLoadingQOBalance } = useQBalanceQuery('quicksilver', address ?? '', 'osmo');
  const { balance: qStars, isLoading: isLoadingQSBalance } = useQBalanceQuery('quicksilver', address ?? '', 'stars');
  const { balance: qRegen, isLoading: isLoadingQRBalance } = useQBalanceQuery('quicksilver', address ?? '', 'regen');
  const { balance: qSomm, isLoading: isLoadingQSOBalance } = useQBalanceQuery('quicksilver', address ?? '', 'somm');
  const { balance: qJuno, isLoading: isLoadingQJBalance } = useQBalanceQuery('quicksilver', address ?? '', 'juno');

  // Retrieve zone data for each token
  const { data: CosmosZone, isLoading: isLoadingCosmosZone } = useZoneQuery(COSMOSHUB_CHAIN_ID ?? '');
  const { data: OsmoZone, isLoading: isLoadingOsmoZone } = useZoneQuery(OSMOSIS_CHAIN_ID ?? '');
  const { data: StarZone, isLoading: isLoadingStarZone } = useZoneQuery(STARGAZE_CHAIN_ID ?? '');
  const { data: RegenZone, isLoading: isLoadingRegenZone } = useZoneQuery(REGEN_CHAIN_ID ?? '');
  const { data: SommZone, isLoading: isLoadingSommZone } = useZoneQuery(SOMMELIER_CHAIN_ID ?? '');
  const { data: JunoZone, isLoading: isLoadingJunoZone } = useZoneQuery(JUNO_CHAIN_ID ?? '');
  // Retrieve APY data for each token
  const { APY: cosmosAPY, isLoading: isLoadingCosmosApy } = useAPYQuery('cosmoshub-4');
  const { APY: osmoAPY, isLoading: isLoadingOsmoApy } = useAPYQuery('osmosis-1');
  const { APY: starsAPY, isLoading: isLoadingStarsApy } = useAPYQuery('stargaze-1');
  const { APY: regenAPY, isLoading: isLoadingRegenApy } = useAPYQuery('regen-1');
  const { APY: sommAPY, isLoading: isLoadingSommApy } = useAPYQuery('sommelier-3');
  const { APY: quickAPY } = useAPYQuery('quicksilver-2');
  const { APY: junoAPY, isLoading: isLoadingJunoApy } = useAPYQuery('juno-1');

  const isLoadingAll =
    isLoadingPrices ||
    isLoadingQABalance ||
    isLoadingQOBalance ||
    isLoadingQSBalance ||
    isLoadingQRBalance ||
    isLoadingQSOBalance ||
    isLoadingQJBalance ||
    isLoadingCosmosZone ||
    isLoadingOsmoZone ||
    isLoadingStarZone ||
    isLoadingRegenZone ||
    isLoadingSommZone ||
    isLoadingJunoZone ||
    isLoadingCosmosApy ||
    isLoadingOsmoApy ||
    isLoadingStarsApy ||
    isLoadingRegenApy ||
    isLoadingSommApy ||
    isLoadingJunoApy;

  // useMemo hook to cache APY data
  const qAPYRates: APYRates = useMemo(
    () => ({
      qAtom: cosmosAPY,
      qOsmo: osmoAPY,
      qStars: starsAPY,
      qRegen: regenAPY,
      qSomm: sommAPY,
      qJuno: junoAPY,
    }),
    [cosmosAPY, osmoAPY, starsAPY, regenAPY, sommAPY, junoAPY],
  );
  // useMemo hook to cache qBalance data
  const qBalances: BalanceRates = useMemo(
    () => ({
      qAtom: shiftDigits(qAtom?.balance?.amount ?? '', -6),
      qOsmo: shiftDigits(qOsmo?.balance?.amount ?? '', -6),
      qStars: shiftDigits(qStars?.balance?.amount ?? '', -6),
      qRegen: shiftDigits(qRegen?.balance?.amount ?? '', -6),
      qSomm: shiftDigits(qSomm?.balance?.amount ?? '', -6),
      qJuno: shiftDigits(qJuno?.balance?.amount ?? '', -6),
    }),
    [qAtom, qOsmo, qStars, qRegen, qSomm, qJuno],
  );

  // useMemo hook to cache redemption rate data
  const redemptionRates: NumericRedemptionRates = useMemo(
    () => ({
      atom: CosmosZone?.redemptionRate ? parseFloat(CosmosZone.redemptionRate) : 1,
      osmo: OsmoZone?.redemptionRate ? parseFloat(OsmoZone.redemptionRate) : 1,
      stars: StarZone?.redemptionRate ? parseFloat(StarZone.redemptionRate) : 1,
      regen: RegenZone?.redemptionRate ? parseFloat(RegenZone.redemptionRate) : 1,
      somm: SommZone?.redemptionRate ? parseFloat(SommZone.redemptionRate) : 1,
      juno: JunoZone?.redemptionRate ? parseFloat(JunoZone.redemptionRate) : 1,
    }),
    [CosmosZone, OsmoZone, StarZone, RegenZone, SommZone, JunoZone],
  );

  // State hooks for portfolio items, total portfolio value, and other metrics
  const [portfolioItems, setPortfolioItems] = useState<PortfolioItemInterface[]>([]);
  const [totalPortfolioValue, setTotalPortfolioValue] = useState(0);
  const [averageApy, setAverageAPY] = useState(0);
  const [totalYearlyYield, setTotalYearlyYield] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  // useEffect hook to compute portfolio metrics when dependencies change

  useEffect(() => {
    const updatePortfolioItems = async () => {
      // Check if all data is loaded
      if (isLoadingAll) {
        return;
      }

      setIsLoading(true);
      let totalValue = 0;
      let totalYearlyYield = 0;
      let weightedAPY = 0;
      let updatedItems = [];

      // Loop through each token to compute value, APY, and yield
      for (const token of Object.keys(qBalances)) {
        const baseToken = token.replace('q', '').toLowerCase();
        // Find the price for the current token
        const tokenPriceInfo = tokenPrices?.find((priceInfo: { token: string }) => priceInfo.token === baseToken);
        const qTokenPrice = tokenPriceInfo ? tokenPriceInfo.price * Number(redemptionRates[baseToken]) : 0;
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
      setIsLoading(false);
    };

    updatePortfolioItems();
  }, [qBalances, CosmosZone, OsmoZone, StarZone, RegenZone, SommZone, redemptionRates, qAPYRates, tokenPrices, isLoadingAll]);

  const assetsData = useMemo(() => {
    return Object.keys(qBalances).map((token) => {
      return {
        name: token.toUpperCase(),
        balance: toNumber(qBalances[token], 2).toString(),
        apy: parseFloat(qAPYRates[token]?.toFixed(2)) || 0,
        native: token.replace('q', '').toUpperCase(),
      };
    });
  }, [qBalances, qAPYRates]);

  const { liquidRewards } = useLiquidRewardsQuery(address ?? '');
  const { authData, authError } = useAuthChecker(address ?? '');

  const [showRewardsClaim, setShowRewardsClaim] = useState(false);
  const [userClosedRewardsClaim, setUserClosedRewardsClaim] = useState(false);

  useEffect(() => {
    if (!authData && authError && !userClosedRewardsClaim) {
      setShowRewardsClaim(true);
    } else {
      setShowRewardsClaim(false);
    }
  }, [authData, authError, userClosedRewardsClaim]);

  // Function to close the RewardsClaim component
  const closeRewardsClaim = () => {
    setShowRewardsClaim(false);
    setUserClosedRewardsClaim(true);
  };

  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container
          flexDir={'column'}
          top={20}
          mt={{ base: '-30px', md: 10 }}
          zIndex={2}
          position="relative"
          justifyContent="center"
          alignItems="center"
          maxW="6xl"
        >
          <Head>
            <title>Assets</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <link rel="icon" href="/img/favicon.png" />
          </Head>
          <Text pb={2} color="white" fontSize="24px">
            Assets
          </Text>

          <Flex flexDir={{ base: 'column', md: 'row' }} py={6} alignItems="center" justifyContent={'space-between'} gap="4">
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
                isLoading={isLoadingAll}
                portfolioItems={portfolioItems}
                isWalletConnected={address !== undefined}
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
              <StakingIntent isWalletConnected={address !== undefined} address={address ?? ''} />
            </Flex>
          </Flex>

          <Spacer />
          {/* Assets Grid */}
          <AssetsGrid nonNative={liquidRewards} isWalletConnected={address !== undefined} assets={assetsData} />
          <Spacer />
          {/* Unbonding Table */}
          <Box h="full" w="full" mt="20px">
            <UnbondingAssetsTable isWalletConnected={address !== undefined} address={address ?? ''} />
          </Box>
          {/* <Box>
            <Image
              display={{ base: 'none', lg: 'block', md: 'none' }}
              src="/img/quicksilverWord.png"
              alt="Quicksilver"
              position="fixed"
              bottom="150"
              left="1350"
              h={'100px'}
              transform="rotate(90deg)"
            />
          </Box> */}
          <Box h="40px"></Box>
        </Container>
        {showRewardsClaim && (
          <SlideFade in={showRewardsClaim} offsetY="20px" style={{ position: 'fixed', right: '20px', bottom: '20px', zIndex: 10 }}>
            <RewardsClaim address={address ?? ''} onClose={closeRewardsClaim} />
          </SlideFade>
        )}
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
