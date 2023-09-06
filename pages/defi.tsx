import { Box, Container } from '@chakra-ui/react';
import Head from 'next/head';

import { Header } from '@/components';
import { SideHeader } from '@/components';

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
        <Container
          justifyContent="center"
          alignItems="center"
          maxW="5xl"
        >
          <Head>
            <title>DeFi</title>
            <meta
              name="viewport"
              content="width=device-width, initial-scale=1.0"
            />
            <link
              rel="icon"
              href="/img/favicon.png"
            />
          </Head>
        </Container>
      </Box>
    </>
  );
}
