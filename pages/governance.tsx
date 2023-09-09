import { Box, Container } from '@chakra-ui/react';
import dynamic from 'next/dynamic';
import Head from 'next/head';

import { Header } from '@/components';
import { SideHeader } from '@/components';
import { VotingSection } from '@/components';
import { chainName } from '@/config';

const DynamicVotingSection = dynamic(() => Promise.resolve(VotingSection), {
  ssr: false,
});

export default function Home() {
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
        <Header />
        <SideHeader />
        <Container justifyContent="center" alignItems="center" maxW="5xl">
          <Head>
            <title>Governance</title>
            <meta
              name="viewport"
              content="width=device-width, initial-scale=1.0"
            />
            <link rel="icon" href="/img/favicon.png" />
          </Head>
          <Box
            maxHeight="3xl" // Adjust this value based on your preference
            overflowY="auto"
            width="100%"
            padding={2} // Optional: for some spacing inside the box
          >
            {chainName && <DynamicVotingSection chainName={chainName} />}
          </Box>
        </Container>
      </Box>
    </>
  );
}
