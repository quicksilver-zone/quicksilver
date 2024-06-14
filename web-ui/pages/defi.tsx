import { Container, Text, SlideFade, Center } from '@chakra-ui/react';
import Head from 'next/head';

import DefiTable from '@/components/Defi/defiBox';

export default function Home() {
  return (
    <>
      <Head>
        <title>Defi - Quicksilver Zone</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta name="description" content="STAKING SIMPLIFIED | LQUIDITY AMPLIFIED" />
        <meta
          name="keywords"
          content="staking, Quicksilver, crypto, staking, earn rewards, DeFi, blockchain, liquid staking, lst, quicksilver zone, cosmos, Cosmos-SDK, cosmoshub, osmosis, stride, stride zone, cosmos liquid staking, Persistence "
        />
        <meta name="author" content="Quicksilver Zone" />
        <link rel="icon" href="/img/favicon-main.png" />

        <meta property="og:title" content="Defi - Quicksilver Zone" />
        <meta property="og:description" content="STAKING SIMPLIFIED | LQUIDITY AMPLIFIED" />
        <meta property="og:url" content="https://app.quicksilver.zone/defi" />
        <meta property="og:image" content="https://app.quicksilver.zone/img/banner.png" />
        <meta property="og:type" content="website" />
        <meta property="og:site_name" content="Quicksilver Protocol" />

        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content="Defi - Quicksilver Zone" />
        <meta name="twitter:description" content="STAKING SIMPLIFIED | LQUIDITY AMPLIFIED" />
        <meta name="twitter:image" content="https://app.quicksilver.zone/img/banner.png" />
        <meta name="twitter:site" content="@quicksilverzone" />

        <script type="application/ld+json">
          {JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'WebPage',
            name: 'Defi - Quicksilver Zone',
            description: 'STAKING SIMPLIFIED | LQUIDITY AMPLIFIED',
            url: 'https://app.quicksilver.zone/defi',
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
