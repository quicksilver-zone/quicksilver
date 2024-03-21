import { Container, Text, SlideFade, Box, Image, Flex } from '@chakra-ui/react';
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
        <Flex height="100vh" mt={{ base: '-20px' }} alignItems="center" justifyContent="center">
          <Container
            p={4}
            m={0}
            textAlign={'left'} // This sets the text alignment for the container
            flexDir="column"
            position="relative"
            justifyContent="flex-start" // Aligns items to the start of the container, along the cross axis
            alignItems="flex-start" // Aligns items to the start of the container, along the main axis
            maxW="5xl"
          >
            <Text pb={2} color="white" fontSize="24px" alignSelf="flex-start">
              DeFi Portal
            </Text>

            <DefiTable />

            <Image
              display={{ base: 'none', lg: 'block' }}
              src="/img/quicksilverWord.png"
              alt="Quicksilver"
              position="absolute"
              bottom="170"
              left="750"
              h="100px"
              transform="translate(50%, 50%) rotate(90deg)"
            />
          </Container>
        </Flex>
      </SlideFade>
    </>
  );
}
