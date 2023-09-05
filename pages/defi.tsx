import Head from 'next/head';
import {
  Box,
  Divider,
  Grid,
  Heading,
  Text,
  Stack,
  Container,
  Link,
  Button,
  Flex,
  Icon,
  useColorMode,
  useColorModeValue,
  VStack
} from '@chakra-ui/react';
import { BsFillMoonStarsFill, BsFillSunFill } from 'react-icons/bs';
import { Product, Dependency, WalletSection } from '../components';
import { dependencies, products } from '../config';
import { Header } from '../components/react/header';
import { SideHeader } from '../components/react/sideHeader';

export default function Home() {
const bg = useColorModeValue("primary.light", "primary.dark")

  return (
    <>
    <Box
    w="100vw"
    h="100vh"
    bgImage="url('/img/backgroundTest.png')" // Add this line
    bgSize="cover" // Add this line
    bgPosition="center center" // Add this line
    bgAttachment="fixed" // Add this line
    
    >
    <Header />
    <SideHeader/>
    <Container
    justifyContent="center"
    alignItems="center"
      maxW="5xl"
      >
      <Head>
        <title>DeFi</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="/img/favicon.png" />
      </Head>
     
    </Container>
    </Box>
    </>
  );
}
