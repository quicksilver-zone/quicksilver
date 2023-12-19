import { Box, Button, ButtonGroup, Container, Flex, HStack, SlideFade, Spacer, Text } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import axios from 'axios';
import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useEffect, useMemo, useState } from 'react';

import { NetworkSelect } from '@/components';
import AssetsGrid from '@/components/Assets/assetsGrid';
import StakingIntent from '@/components/Assets/intents';
import MyPortfolio from '@/components/Assets/portfolio';
import QuickBox from '@/components/Assets/quickbox';
import UnbondingAssetsTable from '@/components/Assets/unbondingTable';
import { useIntentQuery, useQBalanceQuery, useTokenPriceQuery, useZoneQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';

export interface PortfolioItem {
  title: string;
  percentage: number;
  progressBarColor: string;
  amount: string;
  qTokenPrice: number;
}

type NumericRedemptionRates = {
  [key: string]: number;
};

const DynamingPortfolio = dynamic(() => Promise.resolve(MyPortfolio), {
  ssr: false,
});

const DynamicStakingIntent = dynamic(() => Promise.resolve(StakingIntent), {
  ssr: false,
});

const DynamicQuickBox = dynamic(() => Promise.resolve(QuickBox), {
  ssr: false,
});

const DynamicAssetsGrid = dynamic(() => Promise.resolve(AssetsGrid), {
  ssr: false,
});

const DynamicUnbondingTable = dynamic(() => Promise.resolve(UnbondingAssetsTable), {
  ssr: false,
});

export default function Home() {
  const [selectedOption, setSelectedOption] = useState('cosmoshub');
  const { address } = useChain('quicksilver');
  const { address: qAddress, isWalletConnected } = useChain('quicksilver');

  const { balance: qAtom, isLoading: qAtomIsLoading, isError: qAtomIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'atom');
  const { balance: qOsmo, isLoading: qOsmoIsLoading, isError: qOsmoIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'osmo');
  const { balance: qStars, isLoading: qStarsIsLoading, isError: qStarsIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'stars');
  const { balance: qRegen, isLoading: qRegenIsLoading, isError: qRegenIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'regen');
  const { balance: qSomm, isLoading: qSommIsLoading, isError: qSommIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'somm');

  const { data: CosmosZone, isLoading: isCosmosZoneLoading, isError: isCosmosZoneError } = useZoneQuery('cosmoshub-4');
  const { data: OsmoZone, isLoading: isOsmoZoneLoading, isError: isOsmoZoneError } = useZoneQuery('osmosis-1');
  const { data: StarZone, isLoading: isStarZoneLoading, isError: isStarZoneError } = useZoneQuery('stargaze-1');
  const { data: RegenZone, isLoading: isRegenZoneLoading, isError: isRegenZoneError } = useZoneQuery('regen-1');
  const { data: SommZone, isLoading: isSommZoneLoading, isError: isSommZoneError } = useZoneQuery('sommelier-3');

  const qBalances = useMemo(
    () => ({
      qAtom: shiftDigits(qAtom?.balance.amount ?? '', -6),
      qOsmo: shiftDigits(qOsmo?.balance.amount ?? '', -6),
      qStars: shiftDigits(qStars?.balance.amount ?? '', -6),
      qRegen: shiftDigits(qRegen?.balance.amount ?? '', -6),
      qSomm: shiftDigits(qSomm?.balance.amount ?? '', -6),
    }),
    [qAtom, qOsmo, qStars, qRegen, qSomm],
  );

  // Define a function to fetch token price data
  const fetchTokenPrice = async (token: any) => {
    try {
      const response = await axios.get(`https://api-osmosis.imperator.co/tokens/v2/price/${token}`);
      return response.data.price; // Adjust this according to your API response structure
    } catch (error) {
      console.error('Error fetching token price:', error);
      return null;
    }
  };

  const [portfolioItems, setPortfolioItems] = useState<PortfolioItem[]>([]);
  const [totalPortfolioValue, setTotalPortfolioValue] = useState(0);

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

  useEffect(() => {
    const updatePortfolioItems = async () => {
      let totalValue = 0;
      let updatedItems = [];

      for (const token of Object.keys(qBalances)) {
        const baseToken = token.replace('q', '').toLowerCase();
        const price = await fetchTokenPrice(baseToken);
        const qTokenPrice = price * Number(redemptionRates[baseToken]);
        const qTokenBalance = qBalances[token];

        const itemValue = qTokenBalance * qTokenPrice;
        totalValue += itemValue;

        updatedItems.push({
          title: token.toUpperCase(),
          percentage: 0, // Temporarily set to 0, will be updated later
          progressBarColor: 'complimentary.700',
          amount: qTokenBalance,
          qTokenPrice: qTokenPrice || 0,
        });
      }

      // Now, calculate the percentage of each item
      updatedItems = updatedItems.map((item) => ({
        ...item,
        percentage: ((((item.amount * item.qTokenPrice) / totalValue) * 100) / 100).toFixed(2),
      }));

      setPortfolioItems(updatedItems);
      setTotalPortfolioValue(totalValue);
    };

    updatePortfolioItems();
  }, [qBalances, CosmosZone, OsmoZone, StarZone, RegenZone, SommZone, redemptionRates]);

  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container flexDir={'column'} top={20} zIndex={2} position="relative" justifyContent="center" alignItems="center" maxW="6xl">
          <Head>
            <title>Assets</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <link rel="icon" href="/quicksilver-app-v2/img/favicon.png" />
          </Head>
          <Flex flexDir={'row'} py={6} alignItems="center" justifyContent={'space-between'} gap="4">
            {/* Quick box */}
            <Flex
              position="relative"
              backdropFilter="blur(50px)"
              bgColor="rgba(255,255,255,0.1)"
              borderRadius="10px"
              p={5}
              w="md"
              h="sm"
              flexDir="column"
              justifyContent="space-around"
              alignItems="center"
            >
              <DynamicQuickBox />
            </Flex>
            {/* Portfolio box */}
            <Flex
              alignContent={'center'}
              position="relative"
              backdropFilter="blur(50px)"
              bgColor="rgba(255,255,255,0.1)"
              borderRadius="10px"
              p={5}
              w="lg"
              h="sm"
            >
              <DynamingPortfolio portfolioItems={portfolioItems} isWalletConnected={isWalletConnected} totalValue={totalPortfolioValue} />
            </Flex>
            {/* Intent box */}
            <Flex
              alignContent={'center'}
              position="relative"
              backdropFilter="blur(50px)"
              bgColor="rgba(255,255,255,0.1)"
              borderRadius="10px"
              p={5}
              w="lg"
              h="sm"
            >
              <DynamicStakingIntent address={address ?? ''} />
            </Flex>
          </Flex>
          <Spacer />
          {/* Assets Grid */}
          <AssetsGrid />
          <Spacer />
          {/* Unbonding Table */}
          <Box mt="20px">
            <DynamicUnbondingTable />
          </Box>
          <Box h="40px"></Box>
        </Container>
      </SlideFade>
    </>
  );
}

export const AssetsSection = dynamic(() => Promise.resolve(Home), {
  ssr: false,
});
