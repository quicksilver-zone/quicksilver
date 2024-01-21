import { Box, Container, Flex, VStack, HStack, Stat, StatLabel, StatNumber, SlideFade, SkeletonCircle } from '@chakra-ui/react';

import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useState } from 'react';

import { NetworkSelect } from '@/components';
import { StakingBox } from '@/components';
import { InfoBox } from '@/components';
import { AssetsAccordian } from '@/components';
import { useAPYQuery } from '@/hooks/useQueries';
import { networks as prodNetworks, testNetworks as devNetworks } from '@/state/chains/prod';

const DynamicStakingBox = dynamic(() => Promise.resolve(StakingBox), {
  ssr: false,
});

const DynamicInfoBox = dynamic(() => Promise.resolve(InfoBox), {
  ssr: false,
});

const DynamicAssetBox = dynamic(() => Promise.resolve(AssetsAccordian), {
  ssr: false,
});

const networks = process.env.NEXT_PUBLIC_CHAIN_ENV === 'mainnet' ? prodNetworks : devNetworks;

export default function Staking() {
  const [selectedNetwork, setSelectedNetwork] = useState(networks[0]);

  let newChainId;
  if (selectedNetwork.chainId === 'provider') {
    newChainId = 'cosmoshub-4';
  } else if (selectedNetwork.chainId === 'elgafar-1') {
    newChainId = 'stargaze-1';
  } else if (selectedNetwork.chainId === 'osmo-test-5') {
    newChainId = 'osmosis-1';
  } else if (selectedNetwork.chainId === 'regen-redwood-1') {
    newChainId = 'regen-1';
  } else {
    // Default case
    newChainId = selectedNetwork.chainId;
  }
  const { APY, isLoading, isError } = useAPYQuery(newChainId);
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

  return (
    <>
      <Head>
        <title>Staking</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="/quicksilver/img/favicon.png" />
      </Head>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container top={20} zIndex={2} position="relative" maxW="container.lg" maxH="80vh" h="80vh" mt={{ base: '50px', md: '30px' }}>
          {/* <Image
            alt={''}
            src="/quicksilver/img/metalmisc2.png"
            zIndex={-10}
            position="absolute"
            bottom="-10"
            left="-10"
            boxSize="120px"
          /> */}
          <Flex zIndex={3} direction="column" h="100%">
            {/* Dropdown and Statistic */}
            <Box w="50%">
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
            <Flex h="100%" maxH={'2xl'} flexDir={{ base: 'column', md: 'row' }} gap={{ base: '2', md: '0' }}>
              {/* Staking Box*/}
              <DynamicStakingBox
                selectedOption={selectedNetwork}
                isStakingModalOpen={isStakingModalOpen}
                setStakingModalOpen={setStakingModalOpen}
                isTransferModalOpen={isTransferModalOpen}
                setTransferModalOpen={setTransferModalOpen}
                setBalance={setBalance}
                setQBalance={setQBalance}
              />
              <Box w="10px" display={{ base: 'none', md: 'block' }} />

              {/* Right Box */}
              <Flex flex="1" direction="column">
                {/* Top Half (2/3) */}
                <DynamicInfoBox selectedOption={selectedNetwork} displayApr={displayApr} />

                <Box h="10px" />
                {/* Bottom Half (1/3) */}
                <DynamicAssetBox selectedOption={selectedNetwork} balance={balance} qBalance={qBalance} />
              </Flex>
            </Flex>
          </Flex>
        </Container>
      </SlideFade>
    </>
  );
}
