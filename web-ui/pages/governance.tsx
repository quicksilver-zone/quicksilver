import { Box, Center, Container, Flex, SlideFade, Text } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import dynamic from 'next/dynamic';
import Head from 'next/head';

import { VotingSection } from '@/components';

const DynamicVotingSection = dynamic(() => Promise.resolve(VotingSection), {
  ssr: false,
});

export default function Home() {
  const chainName = 'quicksilver';

  const { address } = useChain(chainName);

  if (!address) {
    return (
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Center>
          <Container
            p={4}
            m={0}
            textAlign={'left'}
            flexDir="column"
            position="relative"
            justifyContent="flex-start"
            alignItems="flex-start"
            h={'auto'}
            maxW="5xl"
          >
            <Head>
              <title>Governance</title>
              <meta name="viewport" content="width=device-width, initial-scale=1.0" />
              <link rel="icon" href="/img/favicon-main.png" />
            </Head>
            <Text pb={2} color="white" fontSize="24px">
              Proposals
            </Text>
            <Flex py={6} alignItems="center" alignContent={'center'} justifyContent={'space-between'} gap="4">
              <Flex
                backdropFilter="blur(50px)"
                bgColor="rgba(255,255,255,0.1)"
                borderRadius="10px"
                p={12}
                maxW="5xl"
                h="md"
                justifyContent="center"
                alignItems="center"
              >
                <Text fontSize="xl">Please connect your wallet to view and vote on proposals.</Text>
              </Flex>
            </Flex>
          </Container>
        </Center>
      </SlideFade>
    );
  }

  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Center>
          <Container
            zIndex={2}
            position="relative"
            maxW="container.lg"
            height="auto"
            display="flex"
            flexDirection="column"
            justifyContent="center"
            alignItems="center"
            mt={-6}
          >
            <Head>
              <title>Governance</title>
              <meta name="viewport" content="width=device-width, initial-scale=1.0" />
              <link rel="icon" href="/img/favicon-main.png" />
            </Head>
            <Box overflow={'none'} width="100%" padding={2} mt={{ base: 10, md: 5 }}>
              <Text display={{ base: 'none', md: 'flex' }} pb={2} color="white" fontSize="24px">
                Proposals
              </Text>
              {chainName && <DynamicVotingSection chainName={chainName} />}
            </Box>
          </Container>
        </Center>
      </SlideFade>
    </>
  );
}
