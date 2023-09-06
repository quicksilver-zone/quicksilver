import { Box, Container } from '@chakra-ui/react';
import Head from 'next/head';

import { Header } from '../components/react/header';
import { SideHeader } from '../components/react/sideHeader';

export default function Home() {
  return (
    <>
      <Box
        w="100vw"
        h="100vh"
        bgImage="url('/img/background.png')"
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
            <title>Quick Silver</title>
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
