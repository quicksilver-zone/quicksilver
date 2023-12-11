import { Box, Button, ButtonGroup, Container, Flex, HStack, Text } from '@chakra-ui/react';
import Head from 'next/head';

import { Header } from '../components/react/header';
import { SideHeader } from '../components/react/sideHeader';

export default function Home() {
  return (
    <>
      <Container justifyContent="center" alignItems="center" maxW="5xl">
        <Head>
          <title>Quick Silver</title>
          <meta name="viewport" content="width=device-width, initial-scale=1.0" />
          <link rel="icon" href="/quicksilver-app-v2/img/favicon.png" />
        </Head>
        <Flex flexDir={'row'} alignItems="center" justifyContent={'space-between'} gap="4">
          {/* Quick box */}
          <Flex
            position="relative"
            backdropFilter="blur(50px)"
            bgColor="rgba(255,255,255,0.1)" // Slightly more visible background
            borderRadius="lg" // Using standard size
            p={6} // Slightly more padding
            w="sm" // A bit wider for better layout
            h="sm" // A bit taller for better layout
            flexDir="column"
            justifyContent="space-around" // Better distribution of space
            alignItems="center"
          >
            <Flex
              justifyContent="center"
              alignItems="center"
              flexDir="row"
              gap={3} // Slightly more gap for visual spacing
            >
              <Box minW="10px" minH="10px" borderRadius="full" bgColor="grey" />
              <Text fontSize="2xl" fontWeight="bold" textAlign="center">
                QCK
              </Text>
            </Flex>
            <Flex direction="column" align="stretch" gap={2}>
              <HStack justifyContent="center">
                <Text fontSize="md" fontWeight="normal">
                  Staking APY:
                </Text>
                <Text fontSize="md" fontWeight="semibold">
                  12.37%
                </Text>
              </HStack>
              <HStack justifyContent="center">
                <Text fontSize="md" fontWeight="normal">
                  Quicksilver Balance:
                </Text>
                <Text fontSize="md" fontWeight="semibold">
                  10.123456
                </Text>
              </HStack>
            </Flex>
            <ButtonGroup
              spacing={3} // Consistent spacing
            >
              <Button size="md" w="full">
                Withdraw
              </Button>
              <Button size="md" w="full">
                Deposit
              </Button>
            </ButtonGroup>
          </Flex>
          <Flex
            alignContent={'center'}
            position="relative"
            backdropFilter="blur(50px)"
            bgColor="rgba(255,255,255,0.1)"
            borderRadius="10px"
            p={5}
            w="xs"
            h="xs"
          ></Flex>
          <Flex
            alignContent={'center'}
            position="relative"
            backdropFilter="blur(50px)"
            bgColor="rgba(255,255,255,0.1)"
            borderRadius="10px"
            p={5}
            w="xs"
            h="xs"
          ></Flex>
        </Flex>
      </Container>
    </>
  );
}
