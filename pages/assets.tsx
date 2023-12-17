import { Box, Button, ButtonGroup, Container, Flex, HStack, SlideFade, Spacer, Text } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import Head from 'next/head';
import { useState } from 'react';

import { NetworkSelect } from '@/components';
import AssetsGrid from '@/components/Assets/assetsGrid';
import StakingIntent from '@/components/Assets/intents';
import MyPortfolio from '@/components/Assets/portfolio';
import QuickBox from '@/components/Assets/quickbox';
import UnbondingAssetsTable from '@/components/Assets/unbondingTable';
import { useIntentQuery } from '@/hooks/useQueries';

export default function Home() {
  const [selectedOption, setSelectedOption] = useState('cosmoshub');

  const { address } = useChain('quicksilver');
  const { intent, isLoading, isError } = useIntentQuery('cosmoshub', address ?? '');
  console.log(intent);
  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container flexDir={'column'} top={20} zIndex={2} position="relative" justifyContent="center" alignItems="center" maxW="6xl">
          <Head>
            <title>Quick Silver</title>
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
              <MyPortfolio />
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
              <StakingIntent />
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
