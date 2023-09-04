import Head from "next/head";
import { useState } from "react"
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
  VStack,
  Spacer,
  Select,
  HStack,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Switch,
  Input
} from "@chakra-ui/react";
import { Tabs, TabList, TabPanels, Tab, TabPanel } from '@chakra-ui/react'
import { BsFillMoonStarsFill, BsFillSunFill, BsArrowDown } from "react-icons/bs";
import { Product, Dependency, WalletSection } from "../components";
import { dependencies, products } from "../config";
import { Header } from "../components/react/header";
import { SideHeader } from "../components/react/sideHeader";

export default function Staking() {
  const bg = useColorModeValue("primary.light", "primary.dark");
  const buttonTextColor = useColorModeValue("primary.700", "primary.50");
  const invertButtonTextColor = useColorModeValue("primary.50", "primary.700");
  const [selectedOption, setSelectedOption] = useState("Select Token");
  return (
    <>
      <Box
        w="100vw"
        h="100vh"
        bgImage="url('/img/background.png')"
        bgSize="cover"
        bgPosition="center center"
        bgAttachment="fixed"
      >
        <Header />
        <SideHeader />
        <Container 
        maxW="container.lg" maxH="80vh" h="80vh">
          <Flex direction="column" h="100%">
            {/* Dropdown and Statistic */}
            <Box w="50%" >
            <HStack justifyContent="space-between" w="100%">
        <Menu>
          <MenuButton
            maxW="150px"
            minW="150px"
            variant="ghost"
            color="complimentary.900"
            _hover={{
              bgColor: "rgba(0,0,0,0.5)", 
              backdropFilter: "blur(10px)",
              
            }}
            _active={{
              bgColor: "rgba(0,0,0,0.5)", 
              backdropFilter: "blur(10px)",
              
            }}
            borderColor={buttonTextColor}
            opacity={1}
            as={Button} rightIcon={<BsArrowDown />}
          >
            {selectedOption} {/* Display the selected option */}
          </MenuButton>
          <MenuList>
            {/* Update state when an option is selected */}
            <MenuItem onClick={() => setSelectedOption("Atom")}>Cosmos Hub</MenuItem>
            <MenuItem onClick={() => setSelectedOption("Osmo")}>Osmosis</MenuItem>
            <MenuItem onClick={() => setSelectedOption("Inj")}>Injective</MenuItem>
          </MenuList>
        </Menu>
        <VStack 
        p={1}
        borderRadius="10px"
        alignItems="flex-end">
        <Stat
                
                color="complimentary.900"
                >
  <StatLabel>APY</StatLabel>
  <StatNumber>35%</StatNumber>
</Stat>
        </VStack>
      </HStack>
            </Box>

            {/* Content Boxes */}
            <Flex h="100%">
              {/* Left Box */}
              <Box 
               backdropFilter="blur(50px)"
              bgColor="rgba(0,0,0,0.5)" flex="1" borderRadius="10px" p={5}>
              <Tabs isFitted variant='enclosed'>
  <TabList
    mt={"4"}
    mb='1em'
    borderRadius="md"
    overflow="hidden"
    borderBottomColor="transparent"
    bg="rgba(255,255,255,0.1)"
    p={2}
    borderRadius="25px"
  >
      <Tab
  borderRadius="25px"
  flex="1"
  color="white"
  fontWeight="bold"
  transition="background-color 0.2s ease-in-out, color 0.2s ease-in-out, border-color 0.2s ease-in-out" // Added transition property
  _hover={{
    borderBottomColor: "complimentary.900",
  }}
  _selected={{
    bgColor: "rgba(0,0,0,0.5)",
    color: "complimentary.900",
    borderColor: "complimentary.900",
  }}
>
  Stake
</Tab>
<Tab
  borderRadius="25px"
  flex="1"
  color="white"
  fontWeight="bold"
  transition="background-color 0.2s ease-in-out, color 0.2s ease-in-out, border-color 0.2s ease-in-out" // Added transition property
  _hover={{
    borderBottomColor: "complimentary.900",
  }}
  _selected={{
    bgColor: "rgba(0,0,0,0.5)",
    color: "complimentary.900",
    borderColor: "complimentary.900",
  }}
>
  Unstake
</Tab>
  </TabList>
  <TabPanels>
    <TabPanel
    >
      <VStack
      spacing={8}
      align="center"
      >
      <Text
      textAlign="center"
      color="white"
      >
      Stake your ATOM tokens in exchange for qATOM which you can deploy around the ecosystem. You can liquid stake half of your balance, if you're going to LP.
      </Text>
<Flex
flexDirection="column"
w="100%"
>
<Stat
py={4}
                textAlign="left"
                color="white"
                >
  <StatLabel>Amount to stake:</StatLabel>
  <StatNumber>Atom</StatNumber>
</Stat>
<Input
_active={{
  borderColor: "complimentary.900",
}}
_selected={{
  borderColor: "complimentary.900",
}}
_hover={{
  borderColor: "complimentary.900",
}}
_focus={{
  borderColor: "complimentary.900",
  boxShadow: "0 0 0 3px #FF8000",
}}
color="complimentary.900"
textAlign={"right"}
placeholder="amount"
/>
<Flex 
w="100%"
flexDirection="row" py={4} mb={-4} justifyContent="space-between" alignItems="center">
    <Text color="white" fontWeight="light">
      Tokens available: 0 ATOM
    </Text>
    <HStack spacing={2}>
      <Button 
       _hover={{
        bgColor: "rgba(255,255,255,0.5)", 
        backdropFilter: "blur(10px)",
        
      }}
      _active={{
        bgColor: "rgba(255,255,255,0.5)", 
        backdropFilter: "blur(10px)",
        
      }}
      color="complimentary.900" variant="ghost" w="60px" h="30px">
        half
      </Button>
      <Button 
       _hover={{
        bgColor: "rgba(255,255,255,0.5)", 
        backdropFilter: "blur(10px)",
        
      }}
      _active={{
        bgColor: "rgba(255,255,255,0.5)", 
        backdropFilter: "blur(10px)",
        
      }}
      color="complimentary.900" variant="ghost" w="60px" h="30px">
        max
      </Button>
    </HStack>
  </Flex>
</Flex>
<Divider/>
<HStack
  justifyContent="space-between"
  alignItems="left"
  w="100%"
  mt={-8}
>
  <Stat
    textAlign="left"
    color="white"
  >
    <StatLabel>What you'll get</StatLabel>
    <StatNumber>qAtom:</StatNumber>
  </Stat>
  <Spacer /> {/* This will push the next Stat component to the right */}
  <Stat
    py={4}
    textAlign="right"
    color="white"
  >
    <StatNumber
    textColor="complimentary.900"
    >0</StatNumber>
  </Stat>
</HStack>
<Button
width="100%"
>Liquid Stake</Button>
</VStack>

    </TabPanel>
    <TabPanel>
      <p>two!</p>
    </TabPanel>
  </TabPanels>
</Tabs>
              </Box>

              <Box w="10px" />

              {/* Right Box */}
              <Flex flex="1" direction="column">
                {/* Top Half (2/3) */}
                <Box 
                backdropFilter="blur(50px)"
                borderRadius="10px" bgColor="rgba(0,0,0,0.5)" flex="2" p={5}>
                  {/* Content for Top Right Box */}
                </Box>

                <Box h="10px" />
                {/* Bottom Half (1/3) */}
                <Box 
                backdropFilter="blur(50px)"
                borderRadius="10px" flex="1" bgColor="rgba(0,0,0,0.5)" p={5}>
                  {/* Content for Bottom Right Box */}
                </Box>
              </Flex>
            </Flex>
          </Flex>
        </Container>
      </Box>
    </>
  );
}
