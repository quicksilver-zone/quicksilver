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
import { chains, env, tokenToChainIdMap, getChainForMajorDenom, getChainForQDenom } from '@/config';
import { useGrpcQueryClient } from '@/hooks/useGrpcQueryClient';
import {
  useAPYQuery,
  useAPYsQuery,
  useAuthChecker,
  useCurrentInterchainAssetsQuery,
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
  const tokens = Array.from(chains.get(env)?.values() || []).map((chain) => chain.major_denom.toLowerCase());
  //const tokens = ['atom', 'osmo', 'stars', 'regen', 'somm', 'juno', 'dydx', 'saga', 'bld'];

  const { grpcQueryClient } = useGrpcQueryClient('quicksilver');

  const { data: tokenPrices, isLoading: isLoadingPrices } = useTokenPrices(tokens);
  const { qbalance, qIsLoading, qIsError, qRefetch } = useQBalancesQuery('quicksilver-2', address ?? '', grpcQueryClient);
  const { APYs, APYsLoading } = useAPYsQuery();
  const { redemptionRates, redemptionLoading } = useRedemptionRatesQuery();
  const { APY: quickAPY } = useAPYQuery('quicksilver-2');
  const { assets, refetch: interchainAssetsRefetch } = useCurrentInterchainAssetsQuery(address ?? '');
  const { authData, authError, authRefetch } = useAuthChecker(address ?? '');


  const refetchAll = () => {
    qRefetch();
    interchainAssetsRefetch();
  };

  const isLoadingAll = qIsLoading || APYsLoading || redemptionLoading || isLoadingPrices;

  const nonNative = assets?.assets;
  const portfolioItems: PortfolioItemInterface[] = useMemo(() => {
    if (!qbalance || !APYs || !redemptionRates || isLoadingAll || !assets) return [];

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
      const chain = getChainForQDenom(env, denom);
      const normalizedDenom = chain?.major_denom ?? "";
      const chainId = chain?.chain_id;
      const tokenPriceInfo = tokenPrices?.get(normalizedDenom.toLocaleUpperCase());
      const redemptionRate = chainId && redemptionRates[chainId] ? redemptionRates[chainId].current : 1;
      const qTokenPrice = tokenPriceInfo ? tokenPriceInfo * redemptionRate : 0;
      const exp = chain?.exponent;
      const normalizedAmount = shiftDigits(amount, -(exp ?? 6));
      
      return {
        title: 'q' + normalizedDenom?.toUpperCase(),
        amount: normalizedAmount.toString(),
        qTokenPrice: qTokenPrice,
        chainId: chainId ?? '',
      };
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [qbalance, APYs, redemptionRates, isLoadingAll, assets, tokenToChainIdMap, tokenPrices, refetchAll]);

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
  const liveChains = Array.from(chains.get(env)?.values() || []).filter(chain => chain.show).map((chain) => chain.chain_name);

  const assetsData = useMemo(() => {
    return liveChains.map((chain_name) => {
      const chain = chains.get(env)?.get(chain_name);
      const baseToken = chain?.major_denom.toLowerCase();

      const asset = qbalance?.find((a) => a.denom.substring(2).toLowerCase() === baseToken);
      const chainId = chain?.chain_id;
      const apy = chainId && APYs && APYs[chainId] ? APYs[chainId] : 0;

      const redemptionRate = chainId && redemptionRates && redemptionRates[chainId] ? redemptionRates[chainId].last || 1 : 1;
      const exp = chain?.exponent;
      return {
        name: "q"+chain?.major_denom.toUpperCase(),
        balance: asset ? shiftDigits(Number(asset.amount), -(exp || 6)).toString() : '0',
        apy: parseFloat(((apy * 100) / 100).toFixed(4)),
        native: baseToken?.toUpperCase() ?? '',
        redemptionRates: redemptionRate.toString(),
      };
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [liveChains, qbalance, tokenToChainIdMap, APYs, redemptionRates, refetchAll]);

  const showAssetsGrid = qbalance && qbalance.length > 0 && !qIsLoading && !qIsError;

  if (!address) {
    return (
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Center>
          <Flex height="auto" alignItems="center" justifyContent="center">
            <Container
              p={4}
              m={0}
              textAlign={'left'}
              flexDir="column"
              position="relative"
              justifyContent="flex-start"
              alignItems="flex-start"
            >
              <Head>
                <title>Assets - Quicksilver Zone</title>
                <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                <meta name="description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
                <meta
                  name="keywords"
                  content="staking, Quicksilver, crypto, staking, earn rewards, DeFi, blockchain, liquid staking, lst, quicksilver zone, cosmos, Cosmos-SDK, cosmoshub, osmosis, stride, stride zone, cosmos liquid staking, Persistence "
                />
                <meta name="author" content="Quicksilver Zone" />
                <link rel="icon" href="/img/favicon-main.png" />

                <meta property="og:title" content="Assets - Quicksilver Zone" />
                <meta property="og:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
                <meta property="og:url" content="https://app.quicksilver.zone/assets" />
                <meta property="og:image" content="https://app.quicksilver.zone/img/banner.png" />
                <meta property="og:type" content="website" />
                <meta property="og:site_name" content="Quicksilver Protocol" />

                <meta name="twitter:card" content="summary_large_image" />
                <meta name="twitter:title" content="Assets - Quicksilver Zone" />
                <meta name="twitter:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
                <meta name="twitter:image" content="https://app.quicksilver.zone/img/banner.png" />
                <meta name="twitter:site" content="@quicksilverzone" />

                <script type="application/ld+json">
                  {JSON.stringify({
                    '@context': 'https://schema.org',
                    '@type': 'WebPage',
                    name: 'Assets - Quicksilver Zone',
                    description: 'STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED',
                    url: 'https://app.quicksilver.zone/assets',
                    image: 'https://app.quicksilver.zone/img/banner.png',
                    publisher: {
                      '@type': 'Organization',
                      name: 'Quicksilver Protocol',
                      logo: {
                        '@type': 'ImageObject',
                        url: 'https://app.quicksilver.zone/img/logo.png',
                      },
                    },
                  })}
                </script>
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
        <Center h={'100vh'}>
          <Container
            flexDir={'column'}
            height={'100vh'}
            my={'auto'}
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
            <Text mb={-2} color="white" fontSize="24px">
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
                interchainAssets={assets}
                address={address}
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
        </Center>
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
