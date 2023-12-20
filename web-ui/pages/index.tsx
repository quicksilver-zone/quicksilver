import { Box, Flex } from '@chakra-ui/react';
import Head from 'next/head';

import { Header } from '@/components';
import { SideHeader } from '@/components';
import LiquidMetalSphere from '@/components/ThreeJS/liquidMetalSphere';

export default function Home() {
  return (
    <>
      <Box justifyContent="center" alignItems="center" maxW="5xl">
        <Head>
          <title>Quicksilver</title>
          <meta name="viewport" content="width=device-width, initial-scale=1.0" />
          <link rel="icon" href="/quicksilver/img/favicon.png" />
        </Head>
      </Box>

      <Flex flexDir={'column'}>
        <LiquidMetalSphere />
      </Flex>
    </>
  );
}
