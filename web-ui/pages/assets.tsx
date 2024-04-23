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
import { useGrpcQueryClient } from '@/hooks/useGrpcQueryClient';
import {
  useAPYQuery,
  useAPYsQuery,
  useAuthChecker,
  useLiquidRewardsQuery,
  useQBalancesQuery,
  useTokenPrices,
  useRedemptionRatesQuery,
} from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';

export interface PortfolioItemInterface {
  title: string;
  amount: string;
  qTokenPrice: number;
  chainId: string;
}

function Home() {
  const { address } = useChain('quicksilver');
  const tokens = ['atom', 'osmo', 'stars', 'regen', 'somm', 'juno', 'dydx', 'saga', 'bld'];
  const getExponent = (denom: string) => (['qdydx', 'aqdydx'].includes(denom) ? 18 : 6);

  const { grpcQueryClient } = useGrpcQueryClient('quicksilver');

  const { data: tokenPrices, isLoading: isLoadingPrices } = useTokenPrices(tokens);
  const { qbalance, qIsLoading, qIsError, qRefetch } = useQBalancesQuery('quicksilver-2', address ?? '', grpcQueryClient);
  const { APYs, APYsLoading } = useAPYsQuery();
  const { redemptionRates, redemptionLoading } = useRedemptionRatesQuery();
  const { APY: quickAPY } = useAPYQuery('quicksilver-2');
  const { liquidRewards, refetch: liquidRefetch } = useLiquidRewardsQuery(address ?? '');
  const { authData, authError, authRefetch } = useAuthChecker(address ?? '');

  const refetchAll = () => {
    qRefetch();
    liquidRefetch();
  };

  const isLoadingAll = qIsLoading || APYsLoading || redemptionLoading || isLoadingPrices;

  const COSMOSHUB_CHAIN_ID = process.env.NEXT_PUBLIC_COSMOSHUB_CHAIN_ID;
  const OSMOSIS_CHAIN_ID = process.env.NEXT_PUBLIC_OSMOSIS_CHAIN_ID;
  const STARGAZE_CHAIN_ID = process.env.NEXT_PUBLIC_STARGAZE_CHAIN_ID;
  const REGEN_CHAIN_ID = process.env.NEXT_PUBLIC_REGEN_CHAIN_ID;
  const SOMMELIER_CHAIN_ID = process.env.NEXT_PUBLIC_SOMMELIER_CHAIN_ID;
  const JUNO_CHAIN_ID = process.env.NEXT_PUBLIC_JUNO_CHAIN_ID;
  const DYDX_CHAIN_ID = process.env.NEXT_PUBLIC_DYDX_CHAIN_ID;
  const SAGA_CHAIN_ID = process.env.NEXT_PUBLIC_SAGA_CHAIN_ID;
  const AGORIC_CHAIN_ID = process.env.NEXT_PUBLIC_AGORIC_CHAIN_ID;

  const tokenToChainIdMap: { [key: string]: string | undefined } = useMemo(() => {
    return {
      atom: COSMOSHUB_CHAIN_ID,
      osmo: OSMOSIS_CHAIN_ID,
      stars: STARGAZE_CHAIN_ID,
      regen: REGEN_CHAIN_ID,
      somm: SOMMELIER_CHAIN_ID,
      juno: JUNO_CHAIN_ID,
      dydx: DYDX_CHAIN_ID,
      saga: SAGA_CHAIN_ID,
      agoric: AGORIC_CHAIN_ID,
    };
  }, [
    COSMOSHUB_CHAIN_ID,
    OSMOSIS_CHAIN_ID,
    STARGAZE_CHAIN_ID,
    REGEN_CHAIN_ID,
    SOMMELIER_CHAIN_ID,
    JUNO_CHAIN_ID,
    DYDX_CHAIN_ID,
    SAGA_CHAIN_ID,
    AGORIC_CHAIN_ID,
  ]);

  function getChainIdForToken(tokenToChainIdMap: { [x: string]: any }, baseToken: string) {
    return tokenToChainIdMap[baseToken.toLowerCase()] || null;
  }
  const nonNative = liquidRewards?.assets;
  const portfolioItems: PortfolioItemInterface[] = useMemo(() => {
    if (!qbalance || !APYs || !redemptionRates || isLoadingAll || !liquidRewards) return [];

    // Flatten nonNative assets into a single array and accumulate amounts for each denom
    const amountsMap = new Map();
    Object.values(nonNative || {})
      .flat()
      .flatMap((reward) => reward.Amount)
      .forEach(({ denom, amount }) => {
        const currentAmount = amountsMap.get(denom) || 0;
        amountsMap.set(denom, currentAmount + Number(amount));
      });

    // Map over the accumulated results to create portfolio items
    return Array.from(amountsMap.entries()).map(([denom, amount]) => {
      const normalizedDenom = denom.slice(2);
      const chainId = getChainIdForToken(tokenToChainIdMap, normalizedDenom);
      const tokenPriceInfo = tokenPrices?.find((info) => info.token === normalizedDenom);
      const redemptionRate = chainId && redemptionRates[chainId] ? redemptionRates[chainId].current : 1;
      const qTokenPrice = tokenPriceInfo ? tokenPriceInfo.price * redemptionRate : 0;
      const exp = getExponent(denom);
      const normalizedAmount = shiftDigits(amount, -exp);

      return {
        title: 'q' + normalizedDenom.toUpperCase(),
        amount: normalizedAmount.toString(),
        qTokenPrice: qTokenPrice,
        chainId: chainId ?? '',
      };
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [qbalance, APYs, redemptionRates, isLoadingAll, liquidRewards, nonNative, tokenToChainIdMap, tokenPrices, refetchAll]);

  const totalPortfolioValue = useMemo(
    () => portfolioItems.reduce((acc, item) => acc + Number(item.amount) * item.qTokenPrice, 0),
    [portfolioItems],
  );
  const averageApy = useMemo(() => {
    const totalWeightedApy = portfolioItems.reduce(
      (acc, item) => acc + Number(item.amount) * item.qTokenPrice * (APYs[item.chainId] || 0),
      0,
    );
    return totalWeightedApy / totalPortfolioValue || 0;
  }, [portfolioItems, APYs, totalPortfolioValue]);

  const totalYearlyYield = useMemo(
    () => portfolioItems.reduce((acc, item) => acc + Number(item.amount) * item.qTokenPrice * (APYs[item.chainId] || 0), 0),
    [portfolioItems, APYs],
  );

  const [showRewardsClaim, setShowRewardsClaim] = useState(false);
  const [userClosedRewardsClaim, setUserClosedRewardsClaim] = useState(false);

  useEffect(() => {
    if (!authData && authError && !userClosedRewardsClaim) {
      setShowRewardsClaim(true);
    } else {
      setShowRewardsClaim(false);
    }
  }, [authData, authError, userClosedRewardsClaim]);

  const closeRewardsClaim = () => {
    setShowRewardsClaim(false);
    setUserClosedRewardsClaim(true);
  };

  // Data for the assets grid
  // the query return `qbalance` is an array of quicksilver staked assets held by the user
  // assetsData maps over the assets in qbalance and returns the name, balance, apy, native asset denom, and redemption rate.
  const qtokens = useMemo(() => ['qatom', 'qosmo', 'qstars', 'qregen', 'qsomm', 'qjuno', 'qdydx', 'qsaga', 'qbld'], []);

  const assetsData = useMemo(() => {
    return qtokens.map((token) => {
      const baseToken = token.substring(1).toLowerCase();

      const asset = qbalance?.find((a) => a.denom.substring(2).toLowerCase() === baseToken);
      const apyAsset = qtokens.find((a) => a.substring(1).toLowerCase() === baseToken);
      const chainId = apyAsset ? getChainIdForToken(tokenToChainIdMap, baseToken) : undefined;

      const apy = chainId && chainId !== 'dydx-mainnet-1' && APYs && APYs.hasOwnProperty(chainId) ? APYs[chainId] : 0;
      const redemptionRate = chainId && redemptionRates && redemptionRates[chainId] ? redemptionRates[chainId].last || 1 : 1;
      const exp = apyAsset ? getExponent(apyAsset) : 0;

      return {
        name: token.toUpperCase(),
        balance: asset ? shiftDigits(Number(asset.amount), -exp).toString() : '0',
        apy: parseFloat(((apy * 100) / 100).toFixed(4)),
        native: baseToken.toUpperCase(),
        redemptionRates: redemptionRate.toString(),
      };
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [qtokens, qbalance, tokenToChainIdMap, APYs, redemptionRates, refetchAll]);

  const showAssetsGrid = qbalance && qbalance.length > 0 && !qIsLoading && !qIsError;

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
          {showAssetsGrid && (
            <AssetsGrid
              refetch={refetchAll}
              liquidRewards={liquidRewards}
              address={address}
              nonNative={liquidRewards}
              isWalletConnected={address !== undefined}
              assets={assetsData}
            />
          )}

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
            <RewardsClaim refetch={authRefetch} address={address ?? ''} onClose={closeRewardsClaim} />
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
