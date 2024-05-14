import { Container, Text, SlideFade, Box, Image, Flex, Center } from '@chakra-ui/react';
import Head from 'next/head';

import DefiTable from '@/components/Defi/defiBox';

export default function Home() {
  return (
    <>
      <Head>
        <title>DeFi</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="/img/favicon-main.png" />
      </Head>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Center>
          <Container
            p={4}
            textAlign={'left'}
            flexDir="column"
            height="auto"
            display="flex"
            flexDirection="column"
            justifyContent="center"
            alignItems="center"
            maxW="5xl"
          >
            <Text pb={2} color="white" fontSize="24px" alignSelf="flex-start">
              DeFi Portal
            </Text>

            <DefiTable />
          </Container>
        </Center>
      </SlideFade>
    </>
  );
}
