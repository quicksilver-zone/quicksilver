import { Container, Text, SlideFade } from '@chakra-ui/react';
import Head from 'next/head';

import AirdropSection from '@/components/Airdrop/airdropSection';

export default function Home() {
  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container
          mt={12}
          flexDir={'column'}
          top={20}
          zIndex={2}
          position="relative"
          justifyContent="center"
          alignItems="center"
          maxW="5xl"
        >
          <Head>
            <title>Airdrop</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <link rel="icon" href="/quicksilver/img/favicon.png" />
          </Head>
          <Text pb={2} color="white" fontSize="24px">
            Airdrop
          </Text>
          <AirdropSection />
        </Container>
      </SlideFade>
    </>
  );
}
