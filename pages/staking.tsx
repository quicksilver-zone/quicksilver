import {
  Box,
  Image,
  Container,
  Flex,
  VStack,
  HStack,
  Stat,
  StatLabel,
  StatNumber,
} from '@chakra-ui/react';
import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useState } from 'react';

import { Header } from '@/components';
import { SideHeader } from '@/components';
import { NetworkSelect } from '@/components';
import { StakingBox } from '@/components';
import { InfoBox } from '@/components';
import { AssetsAccordian } from '@/components';
import { useValidatorData } from '@/hooks';
import { useStakingData } from '@/hooks/useStakingData';

const DynamicStakingBox = dynamic(
  () => Promise.resolve(StakingBox),
  { ssr: false },
);

const DynamicNetworkSelect = dynamic(
  () => Promise.resolve(NetworkSelect),
  { ssr: false },
);

export default function Staking() {
  const [selectedOption, setSelectedOption] =
    useState('Atom');
  const [
    selectedChainName,
    setSelectedChainName,
  ] = useState('cosmoshub');
  const [isModalOpen, setModalOpen] =
    useState(false);
  useState(null);

  const { data, isLoading, refetch } =
    useStakingData(selectedChainName);
  console.log(data);
  return (
    <>
      <Box
        w="100vw"
        h="100vh"
        bgImage="url('/img/backgroundTest.png')"
        bgSize="cover"
        bgPosition="center center"
        bgAttachment="fixed"
      >
        <Head>
          <title>Staking</title>
          <meta
            name="viewport"
            content="width=device-width, initial-scale=1.0"
          />
          <link
            rel="icon"
            href="/img/favicon.png"
          />
        </Head>
        <Header />
        <SideHeader />
        <Container
          zIndex={2}
          position="relative"
          mt={-7}
          maxW="container.lg"
          maxH="80vh"
          h="80vh"
        >
          <Image
            alt={''}
            src="/img/metalmisc2.png"
            zIndex={-10}
            position="absolute"
            bottom="-10"
            left="-10"
            boxSize="120px"
          />
          <Flex
            zIndex={3}
            direction="column"
            h="100%"
          >
            {/* Dropdown and Statistic */}
            <Box w="50%">
              <HStack
                justifyContent="space-between"
                w="100%"
              >
                <DynamicNetworkSelect
                  selectedOption={selectedOption}
                  setSelectedOption={
                    setSelectedOption
                  }
                  setSelectedChainName={
                    setSelectedChainName
                  }
                  updateNetworks={refetch}
                />
                <VStack
                  p={1}
                  borderRadius="10px"
                  alignItems="flex-end"
                >
                  <Stat color="complimentary.900">
                    <StatLabel>APY</StatLabel>
                    <StatNumber>35%</StatNumber>
                  </Stat>
                </VStack>
              </HStack>
            </Box>

            {/* Content Boxes */}
            <Flex h="100%">
              {/* Staking Box*/}
              <DynamicStakingBox
                selectedOption={selectedOption}
                isModalOpen={isModalOpen}
                setModalOpen={setModalOpen}
                selectedChainName={
                  selectedChainName
                }
                validators={data?.allValidators}
              />

              <Box w="10px" />

              {/* Right Box */}
              <Flex flex="1" direction="column">
                {/* Top Half (2/3) */}
                <InfoBox
                  selectedOption={selectedOption}
                />

                <Box h="10px" />
                {/* Bottom Half (1/3) */}
                <AssetsAccordian
                  selectedOption={selectedOption}
                />
              </Flex>
            </Flex>
          </Flex>
        </Container>
      </Box>
    </>
  );
}
