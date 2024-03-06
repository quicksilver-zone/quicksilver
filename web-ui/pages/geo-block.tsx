import { Container, Text, SlideFade, Image, Box, Center, VStack } from '@chakra-ui/react';
import Head from 'next/head';

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
            <title>Geo Block</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <link rel="icon" href="/img/favicon.png" />
          </Head>

          <Center>
            <VStack spacing={4}>
              <Box my={4}>
                <Image
                  src="https://media1.tenor.com/m/eaKunLhdjk8AAAAC/terminator-no.gif"
                  alt="T-1000 No Access"
                  boxSize="400px"
                  objectFit="cover"
                />
              </Box>
              <Text fontSize="48px" fontWeight="bold" color="white">
                YOU WON&apos;T BE BACK
              </Text>
              <Text fontSize="20px" color="white">
                Access to this site is not permitted from your current location in the US or the UK.
              </Text>
            </VStack>
          </Center>
        </Container>
      </SlideFade>
    </>
  );
}
