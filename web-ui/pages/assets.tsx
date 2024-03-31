import { Box, Container, Flex, SlideFade, Spacer, Text, Center } from '@chakra-ui/react';
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
import { useLiveZones } from '@/state/LiveZonesContext';
import { shiftDigits, truncateToTwoDecimals } from '@/utils';

export interface PortfolioItemInterface {
  title: string;
  percentage: string;
  progressBarColor: string;
  amount: string;
  qTokenPrice: number;
}

interface RedemptionRate {
  current: number;
  last: number;
}

interface RedemptionRates {
  cosmoshub: RedemptionRate;
  osmosis: RedemptionRate;
  stargaze: RedemptionRate;
  regen: RedemptionRate;
  sommelier: RedemptionRate;
  juno: RedemptionRate;
  dydx: RedemptionRate;
  [key: string]: RedemptionRate;
}

type BalanceRates = {
  [key: string]: string;
};

type APYRates = {
  [key: string]: Number;
};

function Home() {
  const { address } = useChain('quicksilver');
  const tokens = ['atom', 'osmo', 'stars', 'regen', 'somm', 'juno', 'dydx'];

  const { data: tokenPrices, isLoading: isLoadingPrices } = useTokenPrices(tokens);

  // TODO: Use live chain ids from .env
  const COSMOSHUB_CHAIN_ID = process.env.NEXT_PUBLIC_COSMOSHUB_CHAIN_ID;
  const OSMOSIS_CHAIN_ID = process.env.NEXT_PUBLIC_OSMOSIS_CHAIN_ID;
  const STARGAZE_CHAIN_ID = process.env.NEXT_PUBLIC_STARGAZE_CHAIN_ID;
  const REGEN_CHAIN_ID = process.env.NEXT_PUBLIC_REGEN_CHAIN_ID;
  const SOMMELIER_CHAIN_ID = process.env.NEXT_PUBLIC_SOMMELIER_CHAIN_ID;
  const JUNO_CHAIN_ID = process.env.NEXT_PUBLIC_JUNO_CHAIN_ID;
  const DYDX_CHAIN_ID = process.env.NEXT_PUBLIC_DYDX_CHAIN_ID;

  const chainIds = [
    { name: 'cosmoshub', chainId: COSMOSHUB_CHAIN_ID, denom: 'atom' },
    { name: 'osmosis', chainId: OSMOSIS_CHAIN_ID, denom: 'osmo' },
    { name: 'stargaze', chainId: STARGAZE_CHAIN_ID, denom: 'stars' },
    { name: 'regen', chainId: REGEN_CHAIN_ID, denom: 'regen' },
    { name: 'sommelier', chainId: SOMMELIER_CHAIN_ID, denom: 'somm' },
    { name: 'juno', chainId: JUNO_CHAIN_ID, denom: 'juno' },
    { name: 'dydx', chainId: DYDX_CHAIN_ID, denom: 'dydx' },
  ];

  const tokenToZoneMapping: { [key: string]: string } = {
    qAtom: 'cosmoshub',
    qOsmo: 'osmosis',
    qStars: 'stargaze',
    qJuno: 'juno',
    qSomm: 'sommelier',
    qRegen: 'regen',
    qDydx: 'dydx',
  };

  // Dynamic retrieval of balance and zone data

  const balances: BalanceRates = {};
  const zones: RedemptionRates = {
    cosmoshub: {
      current: 0,
      last: 0,
    },
    osmosis: {
      current: 0,
      last: 0,
    },
    stargaze: {
      current: 0,
      last: 0,
    },
    regen: {
      current: 0,
      last: 0,
    },
    sommelier: {
      current: 0,
      last: 0,
    },
    juno: {
      current: 0,
      last: 0,
    },
    dydx: {
      current: 0,
      last: 0,
    },
  };
  const apys: APYRates = {};
  const isLoadingBalances: Record<string, boolean> = {};
  const isLoadingZones: Record<string, boolean> = {};
  const isLoadingApys: Record<string, boolean> = {};

  chainIds.forEach(({ name, chainId, denom }) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    const { balance, isLoading: isLoadingBalance } = useQBalanceQuery('quicksilver', address ?? '', denom);
    // eslint-disable-next-line react-hooks/rules-of-hooks
    const { data: zone, isLoading: isLoadingZone } = useZoneQuery(chainId ?? '');
    // eslint-disable-next-line react-hooks/rules-of-hooks
    const { APY: apy, isLoading: isLoadingApy } = useAPYQuery(chainId ?? '');
    let shift = 6;
    if (denom === 'dydx') {
      shift = 18;
    }
    balances[name as keyof typeof balances] = shiftDigits(balance?.balance.amount ?? '', -shift) || '0';
    zones[name as keyof typeof zones] = (zone as unknown as RedemptionRate) || ({ current: 0, last: 0 } as RedemptionRate);
    apys[name as keyof typeof apys] = apy || 0;
    isLoadingBalances[name] = isLoadingBalance;
    isLoadingZones[name] = isLoadingZone;
    isLoadingApys[name] = isLoadingApy;
  });

  console.log({ zones });
  console.log({ balances });
  console.log({ apys });

  // Example of how to access: balances['cosmoshub'], zones['cosmoshub']

  // To check if all data is loaded, you can iterate over isLoadingBalances and isLoadingZones
  const isLoadingAllData = Object.values({ ...isLoadingBalances, ...isLoadingZones, ...isLoadingApys, isLoadingPrices }).some(
    (isLoading) => isLoading,
  );

  // Retrieve list of zones that are enabled for liquid staking || Will use the above instead
  const { liveNetworks } = useLiveZones();

  // TODO: Figure out how to cycle through live networks and retrieve data for each with less lines of code
  // Retrieve balance for each token
  // Depending on whether the chain exists in liveNetworks or not, the query will be enabled/disabled
  const { balance: qAtom, isLoading: isLoadingQABalance } = useQBalanceQuery('quicksilver', address ?? '', 'atom');
  const { balance: qOsmo, isLoading: isLoadingQOBalance } = useQBalanceQuery('quicksilver', address ?? '', 'osmo');
  const { balance: qStars, isLoading: isLoadingQSBalance } = useQBalanceQuery('quicksilver', address ?? '', 'stars');
  const { balance: qRegen, isLoading: isLoadingQRBalance } = useQBalanceQuery('quicksilver', address ?? '', 'regen');
  const { balance: qSomm, isLoading: isLoadingQSOBalance } = useQBalanceQuery('quicksilver', address ?? '', 'somm');
  const { balance: qJuno, isLoading: isLoadingQJBalance } = useQBalanceQuery('quicksilver', address ?? '', 'juno');
  const { balance: qDydx, isLoading: isLoadingQDBalance } = useQBalanceQuery('quicksilver', address ?? '', 'dydx');

  // Retrieve zone data for each token
  const { data: CosmosZone, isLoading: isLoadingCosmosZone } = useZoneQuery(COSMOSHUB_CHAIN_ID ?? '');
  const { data: OsmoZone, isLoading: isLoadingOsmoZone } = useZoneQuery(OSMOSIS_CHAIN_ID ?? '');
  const { data: StarZone, isLoading: isLoadingStarZone } = useZoneQuery(STARGAZE_CHAIN_ID ?? '');
  const { data: RegenZone, isLoading: isLoadingRegenZone } = useZoneQuery(REGEN_CHAIN_ID ?? '');
  const { data: SommZone, isLoading: isLoadingSommZone } = useZoneQuery(SOMMELIER_CHAIN_ID ?? '');
  const { data: JunoZone, isLoading: isLoadingJunoZone } = useZoneQuery(JUNO_CHAIN_ID ?? '');
  const { data: DydxZone, isLoading: isLoadingDydxZone } = useZoneQuery(DYDX_CHAIN_ID ?? '');
  // Retrieve APY data for each token
  const { APY: cosmosAPY, isLoading: isLoadingCosmosApy } = useAPYQuery('cosmoshub-4');
  const { APY: osmoAPY, isLoading: isLoadingOsmoApy } = useAPYQuery('osmosis-1');
  const { APY: starsAPY, isLoading: isLoadingStarsApy } = useAPYQuery('stargaze-1');
  const { APY: regenAPY, isLoading: isLoadingRegenApy } = useAPYQuery('regen-1');
  const { APY: sommAPY, isLoading: isLoadingSommApy } = useAPYQuery('sommelier-3');
  const { APY: quickAPY } = useAPYQuery('quicksilver-2');
  const { APY: junoAPY, isLoading: isLoadingJunoApy } = useAPYQuery('juno-1');
  const { APY: dydxAPY, isLoading: isLoadingDydxApy } = useAPYQuery('dydx-mainnet-1');

  const isLoadingAll =
    isLoadingPrices ||
    isLoadingQABalance ||
    isLoadingQOBalance ||
    isLoadingQSBalance ||
    isLoadingQRBalance ||
    isLoadingQSOBalance ||
    isLoadingQJBalance ||
    isLoadingQDBalance ||
    isLoadingCosmosZone ||
    isLoadingOsmoZone ||
    isLoadingStarZone ||
    isLoadingRegenZone ||
    isLoadingSommZone ||
    isLoadingJunoZone ||
    isLoadingDydxZone ||
    isLoadingCosmosApy ||
    isLoadingOsmoApy ||
    isLoadingStarsApy ||
    isLoadingRegenApy ||
    isLoadingSommApy ||
    isLoadingJunoApy ||
    isLoadingDydxApy;

  // useMemo hook to cache APY data
  const qAPYRates: APYRates = useMemo(
    () => ({
      qAtom: cosmosAPY ?? 0,
      qOsmo: osmoAPY ?? 0,
      qStars: starsAPY ?? 0,
      qRegen: regenAPY ?? 0,
      qSomm: sommAPY ?? 0,
      qJuno: junoAPY ?? 0,
      qDydx: dydxAPY ?? 0,
    }),
    [cosmosAPY, osmoAPY, starsAPY, regenAPY, sommAPY, junoAPY, dydxAPY],
  );
  // useMemo hook to cache qBalance data
  const qBalances: BalanceRates = useMemo(
    () => ({
      qAtom: shiftDigits(qAtom?.balance?.amount ?? '000000', -6),
      qOsmo: shiftDigits(qOsmo?.balance?.amount ?? '000000', -6),
      qStars: shiftDigits(qStars?.balance?.amount ?? '000000', -6),
      qRegen: shiftDigits(qRegen?.balance?.amount ?? '000000', -6),
      qSomm: shiftDigits(qSomm?.balance?.amount ?? '000000', -6),
      qJuno: shiftDigits(qJuno?.balance?.amount ?? '000000', -6),
      qDydx: shiftDigits(qDydx?.balance?.amount ?? '000000', -18),
    }),
    [qAtom, qOsmo, qStars, qRegen, qSomm, qJuno, qDydx],
  );

  // useMemo hook to cache redemption rate data
  const redemptionRates: RedemptionRates = useMemo(
    () => ({
      cosmoshub: {
        current: CosmosZone?.redemptionRate ? parseFloat(CosmosZone.redemptionRate) : 1,
        last: CosmosZone?.lastRedemptionRate ? parseFloat(CosmosZone.lastRedemptionRate) : 1,
      },
      osmosis: {
        current: OsmoZone?.redemptionRate ? parseFloat(OsmoZone.redemptionRate) : 1,
        last: OsmoZone?.lastRedemptionRate ? parseFloat(OsmoZone.lastRedemptionRate) : 1,
      },
      stargaze: {
        current: StarZone?.redemptionRate ? parseFloat(StarZone.redemptionRate) : 1,
        last: StarZone?.lastRedemptionRate ? parseFloat(StarZone.lastRedemptionRate) : 1,
      },
      regen: {
        current: RegenZone?.redemptionRate ? parseFloat(RegenZone.redemptionRate) : 1,
        last: RegenZone?.lastRedemptionRate ? parseFloat(RegenZone.lastRedemptionRate) : 1,
      },
      sommelier: {
        current: SommZone?.redemptionRate ? parseFloat(SommZone.redemptionRate) : 1,
        last: SommZone?.lastRedemptionRate ? parseFloat(SommZone.lastRedemptionRate) : 1,
      },
      juno: {
        current: JunoZone?.redemptionRate ? parseFloat(JunoZone.redemptionRate) : 1,
        last: JunoZone?.lastRedemptionRate ? parseFloat(JunoZone.lastRedemptionRate) : 1,
      },
      dydx: {
        current: DydxZone?.redemptionRate ? parseFloat(DydxZone.redemptionRate) : 1,
        last: DydxZone?.lastRedemptionRate ? parseFloat(DydxZone.lastRedemptionRate) : 1,
      },
    }),
    [CosmosZone, OsmoZone, StarZone, RegenZone, SommZone, JunoZone, DydxZone],
  );

  // State hooks for portfolio items, total portfolio value, and other metrics
  const [portfolioItems, setPortfolioItems] = useState<PortfolioItemInterface[]>([]);
  const [totalPortfolioValue, setTotalPortfolioValue] = useState(0);
  const [averageApy, setAverageAPY] = useState(0);
  const [totalYearlyYield, setTotalYearlyYield] = useState(0);

  // useEffect hook to compute portfolio metrics when dependencies change
  // TODO: cache the computation and make it faster
  const computedValues = useMemo(() => {
    if (isLoadingAll) {
      return { updatedItems: [], totalValue: 0, weightedAPY: 0, totalYearlyYield: 0 };
    }

    let totalValue = 0;
    let totalYearlyYield = 0;
    let weightedAPY = 0;
    let updatedItems = [];

    for (const token of Object.keys(qBalances)) {
      const zone = tokenToZoneMapping[token];
      const baseToken = token.replace('q', '').toLowerCase();
      const tokenPriceInfo = tokenPrices?.find((priceInfo) => priceInfo.token === baseToken);
      const qTokenPrice = tokenPriceInfo ? tokenPriceInfo.price * Number(redemptionRates[zone].current) : 0;
      const qTokenBalance = qBalances[token];
      const itemValue = Number(qTokenBalance) * qTokenPrice;

      const qTokenAPY = qAPYRates[token] || 0;
      const yearlyYield = itemValue * Number(qTokenAPY);
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

    updatedItems = updatedItems.map((item) => {
      const itemValue = Number(item.amount) * item.qTokenPrice;
      return {
        ...item,
        percentage: (((itemValue / totalValue) * 100) / 100).toFixed(2),
      };
    });

    return { updatedItems, totalValue, weightedAPY, totalYearlyYield };
  }, [isLoadingAll, qBalances, tokenPrices, redemptionRates, qAPYRates, address]);

  useEffect(() => {
    if (!isLoadingAll) {
      setPortfolioItems(computedValues.updatedItems);
      setTotalPortfolioValue(computedValues.totalValue);
      setAverageAPY(computedValues.weightedAPY);
      setTotalYearlyYield(computedValues.totalYearlyYield);
    }
  }, [computedValues, isLoadingAll]);

  const assetsData = useMemo(() => {
    return Object.keys(qBalances).map((token) => {
      const zone = tokenToZoneMapping[token];
      return {
        name: token.toUpperCase(),
        balance: truncateToTwoDecimals(Number(qBalances[token])).toString(),
        apy: parseFloat(qAPYRates[token]?.toFixed(2)) || 0,
        native: token.replace('q', '').toUpperCase(),
        redemptionRates: redemptionRates[zone]?.last.toString() || '0',
      };
    });
  }, [qBalances, qAPYRates, redemptionRates]);

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

  if (!address) {
    return (
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Center>
          <Flex height="100vh" mt={{ base: '-20px' }} alignItems="center" justifyContent="center">
            <Container
              p={4}
              m={0}
              textAlign={'left'}
              flexDir="column"
              position="relative"
              justifyContent="flex-start"
              alignItems="flex-start"
              maxW="5xl"
            >
              <Head>
                <title>Assets</title>
                <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                <link rel="icon" href="/img/favicon-main.png" />
              </Head>
              <Text pb={2} color="white" fontSize="24px">
                Assets
              </Text>
              <Flex py={6} alignItems="center" alignContent={'center'} justifyContent={'space-between'} gap="4">
                <Flex
                  backdropFilter="blur(50px)"
                  bgColor="rgba(255,255,255,0.1)"
                  borderRadius="10px"
                  p={12}
                  maxW="5xl"
                  h="md"
                  justifyContent="center"
                  alignItems="center"
                >
                  <Text fontSize="xl">Please connect your wallet to interact with your qAssets.</Text>
                </Flex>
              </Flex>
            </Container>
          </Flex>
        </Center>
      </SlideFade>
    );
  }

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
            <link rel="icon" href="/img/favicon-main.png" />
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
              w={{ base: 'full', md: 'md' }}
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
              px={5}
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
          <AssetsGrid address={address} nonNative={liquidRewards} isWalletConnected={address !== undefined} assets={assetsData} />
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
