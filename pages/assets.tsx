import { Box, Button, ButtonGroup, Container, Flex, HStack, SlideFade, Spacer, Text } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import Head from 'next/head';
import { useEffect, useMemo, useState } from 'react';

import { NetworkSelect } from '@/components';
import AssetsGrid from '@/components/Assets/assetsGrid';
import StakingIntent from '@/components/Assets/intents';
import MyPortfolio from '@/components/Assets/portfolio';
import QuickBox from '@/components/Assets/quickbox';
import UnbondingAssetsTable from '@/components/Assets/unbondingTable';
import { useIntentQuery, useQBalanceQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';

export interface PortfolioItem {
  title: string;
  percentage: number;
  progressBarColor: string;
  amount: string;
}

export default function Home() {
  const [selectedOption, setSelectedOption] = useState('cosmoshub');
  const { address } = useChain('quicksilver');
  const { address: qAddress, isWalletConnected } = useChain('quicksilver');

  const { balance: qAtom, isLoading: qAtomIsLoading, isError: qAtomIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'atom');
  const { balance: qOsmo, isLoading: qOsmoIsLoading, isError: qOsmoIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'osmo');
  const { balance: qStars, isLoading: qStarsIsLoading, isError: qStarsIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'stars');
  const { balance: qRegen, isLoading: qRegenIsLoading, isError: qRegenIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'regen');
  const { balance: qSomm, isLoading: qSommIsLoading, isError: qSommIsError } = useQBalanceQuery('quicksilver', qAddress ?? '', 'somm');

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

  // My Portfolio Computation
  const [portfolioItems, setPortfolioItems] = useState<PortfolioItem[]>([]);

  useEffect(() => {
    const nonZeroBalances = Object.entries(qBalances)
      .filter(([_, balance]) => Number(balance) > 0)
      .map(([token, balance]) => ({ token, balance: Number(balance) }));

    const totalBalance = nonZeroBalances.reduce((total, { balance }) => total + balance, 0);

    const items = nonZeroBalances.map(({ token, balance }) => ({
      title: `${token}`,
      percentage: Number(Number(balance / totalBalance).toFixed(2)),
      progressBarColor: 'complimentary.700',
      amount: `${balance}`,
    }));

    setPortfolioItems(items);
  }, [qBalances]);

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
              <QuickBox />
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
              <MyPortfolio portfolioItems={portfolioItems} isWalletConnected={isWalletConnected} />
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
              <StakingIntent address={address ?? ''} />
            </Flex>
          </Flex>
          <Spacer />
          {/* Assets Grid */}
          <AssetsGrid />
          <Spacer />
          {/* Unbonding Table */}
          <Box mt="20px">
            <UnbondingAssetsTable />
          </Box>
          <Box h="40px"></Box>
        </Container>
      </SlideFade>
    </>
  );
}
