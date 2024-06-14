import { Text, VStack, Heading, Link, Container, SlideFade } from '@chakra-ui/react';
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
          <title>About - Quicksilver Zone</title>
          <meta name="viewport" content="width=device-width, initial-scale=1.0" />
          <meta name="description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
          <meta
            name="keywords"
            content="staking, Quicksilver, crypto, staking, earn rewards, DeFi, blockchain, liquid staking, lst, quicksilver zone, cosmos, Cosmos-SDK, cosmoshub, osmosis, stride, stride zone, cosmos liquid staking, Persistence "
          />
          <meta name="author" content="Quicksilver Zone" />
          <link rel="icon" href="/img/favicon-main.png" />

          <meta property="og:title" content="About - Quicksilver Zone" />
          <meta property="og:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
          <meta property="og:url" content="https://app.quicksilver.zone/about" />
          <meta property="og:image" content="https://app.quicksilver.zone/img/banner.png" />
          <meta property="og:type" content="website" />
          <meta property="og:site_name" content="Quicksilver Protocol" />

          <meta name="twitter:card" content="summary_large_image" />
          <meta name="twitter:title" content="About - Quicksilver Zone" />
          <meta name="twitter:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
          <meta name="twitter:image" content="https://app.quicksilver.zone/img/banner.png" />
          <meta name="twitter:site" content="@quicksilverzone" />

          <script type="application/ld+json">
            {JSON.stringify({
              '@context': 'https://schema.org',
              '@type': 'WebPage',
              name: 'About - Quicksilver Zone',
              description: 'STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED',
              url: 'https://app.quicksilver.zone/about',
              image: 'https://app.quicksilver.zone/img/banner.png',
              publisher: {
                '@type': 'Organization',
                name: 'Quicksilver Protocol',
                logo: {
                  '@type': 'ImageObject',
                  url: 'https://app.quicksilver.zone/img/logo.png',
                },
              },
            })}
          </script>
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
