import { Box, Text, VStack, Heading, Link, Container, SlideFade } from '@chakra-ui/react';
import Head from 'next/head';

const AboutPage = () => {
  return (
    <SlideFade offsetY={'200px'} in={true}>
      <Container
        flexDir={'column'}
        top={20}
        mt={{ base: 10, md: 10 }}
        zIndex={2}
        position="relative"
        justifyContent="center"
        alignItems="center"
        maxW="6xl"
      >
        <Head>
          <title>About </title>
          <meta name="viewport" content="width=device-width, initial-scale=1.0" />
          <link rel="icon" href="/quicksilver/img/favicon.png" />
        </Head>
        <VStack spacing={4} align="stretch" m={8}>
          <Heading as="h1" color="white" size="xl" textAlign="left">
            About Us
          </Heading>
          <Text fontSize="md" color="gray.200">
            QuickSilver is a state-of-the-art platform for liquid staking. We allow users to stake their cryptocurrency in a flexible and
            secure manner. Our mission is to provide a seamless staking experience while maximizing the earning potential for our users.
          </Text>
          <Link href="/quicksilver/privacy-policy" color="orange.400" alignSelf="left">
            Privacy Policy
          </Link>
        </VStack>
      </Container>
    </SlideFade>
  );
};

export default AboutPage;
