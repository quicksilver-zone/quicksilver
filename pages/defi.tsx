import { Box, Container } from '@chakra-ui/react';
import Head from 'next/head';

export default function Home() {
  return (
    <>
      <Container justifyContent="center" alignItems="center" maxW="5xl">
        <Head>
          <title>DeFi</title>
          <meta name="viewport" content="width=device-width, initial-scale=1.0" />
          <link rel="icon" href="/quicksilver-app-v2/img/favicon.png" />
        </Head>
      </Container>
    </>
  );
}
