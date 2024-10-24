import { Box, Container, Flex, VStack, HStack, Stat, StatLabel, StatNumber, SlideFade, SkeletonCircle } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react-lite';
import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useState } from 'react';

import { NetworkSelect } from '@/components';
import { StakingBox } from '@/components';
import { InfoBox } from '@/components';
import { AssetsAccordian } from '@/components';
import { Chain, chains, env } from '@/config';
import { useAPYQuery } from '@/hooks/useQueries';

const DynamicStakingBox = dynamic(() => Promise.resolve(StakingBox), {
  ssr: false,
});

const DynamicInfoBox = dynamic(() => Promise.resolve(InfoBox), {
  ssr: false,
});

const DynamicAssetBox = dynamic(() => Promise.resolve(AssetsAccordian), {
  ssr: false,
});

const networks: Map<string, Chain> = chains.get(env) ?? new Map();
const chain_list = Array.from(networks).filter(([_, network]) => network.show).map(([key, _]) => key);

export default function Staking() {
  const [selectedNetwork, setSelectedNetwork] = useState(networks.get(chain_list[0])); 
  const { address } = useChain('quicksilver');

  const { APY, isLoading, isError } = useAPYQuery(selectedNetwork?.chain_id); // handle testnets (chain.testnet_for)
  const [balance, setBalance] = useState('');
  const [qBalance, setQBalance] = useState('');

  let displayApr = '';
  if (!isLoading && !isError && APY !== undefined) {
    displayApr = (APY * 100).toFixed(2) + '%';
  } else if (isError) {
    displayApr = 'Error';
  }

  const [isStakingModalOpen, setStakingModalOpen] = useState(false);
  const [isTransferModalOpen, setTransferModalOpen] = useState(false);
  const [isRevertSharesModalOpen, setRevertSharesModalOpen] = useState(false);

  return (
    <>
      <Head>
        <title>Staking - Quicksilver Zone</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta name="description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
        <meta
          name="keywords"
          content="staking, Quicksilver, crypto, staking, earn rewards, DeFi, blockchain, liquid staking, lst, quicksilver zone, cosmos, Cosmos-SDK, cosmoshub, osmosis, stride, stride zone, cosmos liquid staking, Persistence "
        />
        <meta name="author" content="Quicksilver Zone" />
        <link rel="icon" href="/img/favicon-main.png" />

        <meta property="og:title" content="Staking - Quicksilver Zone" />
        <meta property="og:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
        <meta property="og:url" content="https://app.quicksilver.zone/staking" />
        <meta property="og:image" content="https://app.quicksilver.zone/img/banner.png" />
        <meta property="og:type" content="website" />
        <meta property="og:site_name" content="Quicksilver Protocol" />

        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content="Staking - Quicksilver Zone" />
        <meta name="twitter:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
        <meta name="twitter:image" content="https://app.quicksilver.zone/img/banner.png" />
        <meta name="twitter:site" content="@quicksilverzone" />

        <script type="application/ld+json">
          {JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'WebPage',
            name: 'Staking - Quicksilver Zone',
            description: 'STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED',
            url: 'https://app.quicksilver.zone/staking',
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

      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container
          zIndex={2}
          position="relative"
          maxW="container.lg"
          height="auto"
          display="flex"
          flexDirection="column"
          justifyContent="center"
          alignItems="center"
          mt={6}
        >
          <Flex justifyContent={'center'} zIndex={3} direction="column" h="100%">
            {/* Dropdown and Statistic */}
            <Box w={{ base: '100%', lg: '50%' }}>
              <HStack justifyContent="space-between" w="100%">
                <NetworkSelect selectedOption={selectedNetwork} setSelectedNetwork={setSelectedNetwork} />
                <VStack p={1} borderRadius="10px" alignItems="flex-end">
                  <Stat minW={'90px'} color="complimentary.900">
                    <StatLabel>APR</StatLabel>
                    <StatNumber height={'34px'}>
                      {!isLoading && APY !== undefined ? (
                        displayApr
                      ) : (
                        <>
                          <HStack height={'34px'}>
                            <SkeletonCircle size="3" startColor="complimentary.900" endColor="complimentary.400" />{' '}
                            <SkeletonCircle size="2" startColor="complimentary.900" endColor="complimentary.400" />
                            <SkeletonCircle size="3" startColor="complimentary.900" endColor="complimentary.400" />
                          </HStack>
                        </>
                      )}
                    </StatNumber>
                  </Stat>
                </VStack>
              </HStack>
            </Box>

            {/* Content Boxes */}
            <Flex h="100%" maxH={'2xl'} flexDir={{ base: 'column', lg: 'row' }} gap={{ base: '2', lg: '0' }}>
              {/* Staking Box*/}
              {selectedNetwork && (
                <DynamicStakingBox
                  selectedOption={selectedNetwork}
                isStakingModalOpen={isStakingModalOpen}
                setStakingModalOpen={setStakingModalOpen}
                isTransferModalOpen={isTransferModalOpen}
                setTransferModalOpen={setTransferModalOpen}
                isRevertSharesModalOpen={isRevertSharesModalOpen}
                setRevertSharesModalOpen={setRevertSharesModalOpen}
                setBalance={setBalance}
                  setQBalance={setQBalance}
                />
              )}
              <Box w="10px" display={{ base: 'none', lg: 'block' }} />

              {/* Right Box */}
              <Flex flex="1" direction="column">
                {/* Top Half (2/3) */}
                <DynamicInfoBox selectedOption={selectedNetwork} displayApr={displayApr} />

                <Box h="10px" />
                {/* Bottom Half (1/3) */}
                {selectedNetwork && (
                  <DynamicAssetBox address={address ?? ''} selectedOption={selectedNetwork} balance={balance} qBalance={qBalance} />
                )}
              </Flex>
            </Flex>
          </Flex>
        </Container>
      </SlideFade>
    </>
  );
}
