import { Box, Container, SlideFade, Text } from '@chakra-ui/react';
import dynamic from 'next/dynamic';
import Head from 'next/head';

import { Header } from '@/components';
import { SideHeader } from '@/components';
import { VotingSection } from '@/components';

const DynamicVotingSection = dynamic(() => Promise.resolve(VotingSection), {
  ssr: false,
});

export default function Home() {
  const chainName = 'quicksilver';

  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container justifyContent="center" alignItems="center" maxW="5xl" mt={12}>
          <Head>
            <title>Governance</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <link rel="icon" href="/quicksilver-app-v2/img/favicon.png" />
          </Head>
          <Box maxHeight="3xl" width="100%" padding={2}>
            <Text pb={2} color="white" fontSize="24px">
              Proposals
            </Text>
            {chainName && <DynamicVotingSection chainName={chainName} />}
          </Box>
        </Container>
      </SlideFade>
    </>
  );
}
