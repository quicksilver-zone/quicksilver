import { Box, Container, SlideFade, Text } from '@chakra-ui/react';
import dynamic from 'next/dynamic';
import Head from 'next/head';

import { VotingSection } from '@/components';

const DynamicVotingSection = dynamic(() => Promise.resolve(VotingSection), {
  ssr: false,
});

export default function Home() {
  const chainName = 'quicksilver';

  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container
          zIndex={2}
          position="relative"
          maxW="container.lg"
          height="100vh" // Full viewport height
          display="flex"
          flexDirection="column"
          justifyContent="center" // Center vertically
          alignItems="center" // Center horizontally
          mt={{ base: '-40px', md: '60px' }} // Ensuring the margin-top is zero
          mb="0" // Ensuring the margin-bottom is zero
        >
          <Head>
            <title>Governance</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <link rel="icon" href="/img/favicon.png" />
          </Head>
          <Box width="100%" padding={2} mt={{ base: 10, md: 5 }}>
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
