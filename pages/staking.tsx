import Head from "next/head";
import { useState } from "react"
import {
  Box,
  Divider,
  Image,
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
import { BsFillMoonStarsFill, BsFillSunFill, BsArrowDown, BsTrophy, BsCoin, BsClock } from "react-icons/bs";
import { RiStockLine } from "react-icons/ri";
import { Product, Dependency, WalletSection } from "../components";
import { dependencies, products } from "../config";
import { Header } from "../components/react/header";
import { SideHeader } from "../components/react/sideHeader";
import {
  Accordion,
  AccordionItem,
  AccordionButton,
  AccordionPanel,
  AccordionIcon,
} from '@chakra-ui/react'

export default function Staking() {
  const bg = useColorModeValue("primary.light", "primary.dark");
  const buttonTextColor = useColorModeValue("primary.700", "primary.50");
  const invertButtonTextColor = useColorModeValue("primary.50", "primary.700");
  const [selectedOption, setSelectedOption] = useState("Atom");
  const [openItem, setOpenItem] = useState(null);
  const [activeAccordion, setActiveAccordion] = useState(null);
  const handleAccordionChange = (accordionNumber, index) => {
    if (activeAccordion === accordionNumber && openItem === index) {
      setOpenItem(null);
      setActiveAccordion(null);
    } else {
      setOpenItem(index);
      setActiveAccordion(accordionNumber);
    }
  };
  return (
    <>
      <Box
        w="100vw"
        h="100vh"
        bgImage="url('/img/backgroundTest.png')"
        bgSize="cover"
        bgPosition="center center"
        bgAttachment="fixed"
      >
        <Head>
        <title>Staking</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="/img/favicon.png" />
      </Head>
        <Header />
        <SideHeader />
        <Container
  zIndex={2}
        position="relative"
        mt={-7}
        maxW="container.lg" maxH="80vh" h="80vh">
            <Image 
        src="/img/metalmisc2.png" 
        zIndex={-10}
        position="absolute" 
        bottom="-10" 
        left="-10" 
        boxSize="120px" // Set the desired size of the image
    />
          <Flex zIndex={3} direction="column" h="100%">
            {/* Dropdown and Statistic */}
            <Box w="50%" >
            <HStack justifyContent="space-between" w="100%">
        <Menu zIndex={4}>
          <MenuButton
          position="relative"
          zIndex={5}
            maxW="150px"
            minW="150px"
            variant="ghost"
            color="complimentary.900"
            backgroundColor="rgba(255,255,255,0.1)"
            _hover={{
              bgColor: "rgba(255,255,255,0.05)", 
              backdropFilter: "blur(10px)",
              
            }}
            _active={{
              bgColor: "rgba(255,255,255,0.05)",
              backdropFilter: "blur(10px)",
              
            }}
            borderColor={buttonTextColor}
            as={Button} rightIcon={<BsArrowDown />}
          >
            {selectedOption.toUpperCase()} 
          </MenuButton>
          <MenuList
          mt={1}
          bgColor="black"
          borderColor="black"
          >
            {/* Update state when an option is selected */}
            <MenuItem 
            bgColor="black"
            borderRadius="4px"
            color="complimentary.900"
            _hover={{
              bgColor: "rgba(255,255,255,0.25)",
            }}
            onClick={() => setSelectedOption("ATOM")}>Cosmos Hub</MenuItem>
            <MenuItem 
            bgColor="black"
            borderRadius="4px"
            color="complimentary.900"
            _hover={{
              bgColor: "rgba(255,255,255,0.25)",
            }}
            onClick={() => setSelectedOption("OSMO")}>Osmosis</MenuItem>
            <MenuItem 
            bgColor="black"
            borderRadius="4px"
            color="complimentary.900"
            _hover={{
              bgColor: "rgba(255,255,255,0.25)",
            }}
            onClick={() => setSelectedOption("INJ")}>Injective</MenuItem>
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
              position="relative"
               backdropFilter="blur(50px)"
              bgColor="rgba(255,255,255,0.1)" flex="1" borderRadius="10px" p={5}>
              <Tabs isFitted variant='enclosed'>
  <TabList
    mt={"4"}
    mb='1em'
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
      Stake your {selectedOption.toUpperCase()}  tokens in exchange for q{selectedOption.toUpperCase()} which you can deploy around the ecosystem. You can liquid stake half of your balance, if you're going to LP.
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
  <StatNumber>{selectedOption.toUpperCase()} </StatNumber>
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
      Tokens available: 0 {selectedOption.toUpperCase()} 
    </Text>
    <HStack spacing={2}>
      <Button 
       _hover={{
        bgColor: "rgba(255,255,255,0.05)", 
        backdropFilter: "blur(10px)",
        
      }}
      _active={{
        bgColor: "rgba(255,255,255,0.05)",
        backdropFilter: "blur(10px)",
        
      }}
      color="complimentary.900" variant="ghost" w="60px" h="30px">
        half
      </Button>
      <Button 
       _hover={{
        bgColor: "rgba(255,255,255,0.05)", 
        backdropFilter: "blur(10px)",
        
      }}
      _active={{
        bgColor: "rgba(255,255,255,0.05)",
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
    <StatNumber>q{selectedOption.toUpperCase()}:</StatNumber>
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
_hover={{
  bgColor: "#181818"
}}
>Liquid Stake</Button>
</VStack>

    </TabPanel>
    <TabPanel>
    <VStack
      spacing={8}
      align="center"
      >
      <Text
      textAlign="center"
      color="white"
      >
      Unstake your q{selectedOption.toUpperCase()} tokens in exchange for {selectedOption.toUpperCase()}.
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
  <StatLabel>Amount tounstake:</StatLabel>
  <StatNumber>q{selectedOption.toUpperCase()} </StatNumber>
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
      Tokens available: 0 q{selectedOption.toUpperCase()} 
    </Text>
   
      <Button 
       _hover={{
        bgColor: "rgba(255,255,255,0.05)", 
        backdropFilter: "blur(10px)",
        
      }}
      _active={{
        bgColor: "rgba(255,255,255,0.05)",
        backdropFilter: "blur(10px)",
        
      }}
      color="complimentary.900" variant="ghost" w="60px" h="30px">
        max
      </Button>

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
    <StatNumber>{selectedOption.toUpperCase()}:</StatNumber>
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
_hover={{
  bgColor: "complimentary.1000"
}}
>Liquid Stake</Button>
</VStack>

    </TabPanel>
  </TabPanels>
</Tabs>
              </Box>

              <Box w="10px" />

              {/* Right Box */}
              <Flex flex="1" direction="column">
                {/* Top Half (2/3) */}
                <Box 
                position="relative"
                backdropFilter="blur(30px)"
                borderRadius="10px" bgColor="rgba(255,255,255,0.1)" flex="2" p={5}>
                   <Image 
        src="/img/metalmisc3.png" 
        zindex={1}
        position="absolute" 
        top="-40px" 
        right="-65px" 
        boxSize="135px" 
        transform="rotate(25deg)"
    />
                   <Text
                 fontSize="20px"
                 color="white"
                 >About {selectedOption.toUpperCase()} on Quicksilver</Text>
                 <Accordion
        mt={6}
        index={activeAccordion === 1 ? openItem : null}
        onChange={(index) => handleAccordionChange(1, index)}
        allowToggle
      >
    <AccordionItem
    pt={2}
    mb={2}
    borderTop={"none"}
    >
      <h2>
      <Flex 
    borderTopColor={"transparent"} 
    alignItems="center" 
    justifyContent="space-between" 
    width="100%"
    py={2} // padding for top and bottom to mimic button's vertical padding
>
    <Flex flexDirection="row" alignItems="center">
        <Box mr="16px">
            <BsTrophy
                color="#FF8000"
                size="24px"
            />
        </Box>
        <Text fontSize="16px" color={"white"}>Rewards</Text>
    </Flex>
    <Text pr={2} color="complimentary.900">35%</Text>
</Flex>
      </h2>
     
    </AccordionItem>

    <AccordionItem
    pt={2}
    mb={2}
    >
      <h2>
      <Flex 
    borderTopColor={"transparent"} 
    alignItems="center" 
    justifyContent="space-between" 
    width="100%"
    py={2} // padding for top and bottom to mimic button's vertical padding
>
          <Flex  flexDirection="row" flex='1' alignItems="center">
          <Box mr="16px"> {/* Adjusts right margin */}
  <BsCoin
      color="#FF8000"
      size="24px"
  />
</Box>
            <Text 
            fontSize="16px"
            color={"white"}>Fees</Text>
          </Flex>
          <Text
           pr={2}
          color="complimentary.900"
          >Low</Text>
          
        </Flex>
      </h2>
     
    </AccordionItem>
    <AccordionItem
    pt={2}
    mb={2}
    >
      <h2>
      <Flex 
    borderTopColor={"transparent"} 
    alignItems="center" 
    justifyContent="space-between" 
    width="100%"
    py={2} // padding for top and bottom to mimic button's vertical padding
>
          <Flex  flexDirection="row" flex='1' alignItems="center">
          <Box mr="16px"> {/* Adjusts right margin */}
  <BsClock
      color="#FF8000"
      size="24px"
  />
</Box>
            <Text 
            fontSize="16px"
            color={"white"}>Unbonding</Text>
          </Flex>
          <Text
           pr={2}
          color="complimentary.900"
          >21-24 Days*</Text>

        </Flex>
      </h2>
      <AccordionPanel 
  alignItems="center"
  justifyItems="center"
  color="white"
  pb={4}
>
  <VStack spacing={2} width="100%">
    <HStack justifyContent="space-between" width="100%">
      <Text
      color="white"
      >on {selectedOption.toUpperCase()}</Text>
      <Text
      color="complimentary.900"
      >0 {selectedOption.toUpperCase()}</Text>
    </HStack>
    <HStack justifyContent="space-between" width="100%">
      <Text
      color="white"
      >on Quicksilver</Text>
      <Text
      color="complimentary.900"
      >0 {selectedOption.toUpperCase()}</Text>
    </HStack>
  </VStack>
</AccordionPanel>
    </AccordionItem>
    <AccordionItem
    pt={2}
    mb={2}
    borderBottom={"none"}
    >
      <h2>
      <Flex 
    borderTopColor={"transparent"} 
    alignItems="center" 
    justifyContent="space-between" 
    width="100%"
    py={2} // padding for top and bottom to mimic button's vertical padding
>
          <Flex flexDirection="row" flex='1' alignItems="center">
          <Box mr="16px"> {/* Adjusts right margin */}
  <RiStockLine
      color="#FF8000"
      size="24px"
  />
</Box>
            <Text 
            fontSize="16px"
            color={"white"}>Value of 1 q{selectedOption.toUpperCase()}</Text>
          </Flex>
          <Text
           pr={2}
          color="complimentary.900"
          >1 q{selectedOption.toUpperCase()} = 1 {selectedOption.toUpperCase()}</Text>
        </Flex>
      </h2>
     
    </AccordionItem>
  </Accordion>

  <Text mt={3}  color="white" textAlign="center" bgColor="rgba(0,0,0,0.4)" p={5} width="100%" borderRadius={6}>
  Want to learn more about rewards, fees, and unbonding on Quicksilver? Check out the <Link href="https://your-docs-url.com" color="complimentary.900" isExternal>docs</Link>.
</Text>

                </Box>

                <Box h="10px" />
                {/* Bottom Half (1/3) */}
                <Box 
                 position="relative"
                backdropFilter="blur(10px)"
                zindex={10}
                borderRadius="10px" bgColor="rgba(255,255,255,0.1)" flex="1" p={5}>
                 <Text
                 fontSize="20px"
                 color="white"
                 >Assets</Text>
                <Accordion
        mt={6}
        index={activeAccordion === 2 ? openItem : null}
        onChange={(index) => handleAccordionChange(2, index)}
        allowToggle
      >
    <AccordionItem
    mb={4}
    borderTop={"none"}
    >
      <h2>
        <AccordionButton borderTopColor={"transparent"}>
          <Flex p={1} flexDirection="row" flex='1' alignItems="center">
            <Image src="/img/networks/atom.svg" boxSize="35px" mr={2} />
            <Text 
            fontSize="16px"
            color={"white"}>Available to stake</Text>
          </Flex>
          <Text
           pr={2}
          color="complimentary.900"
          >0 {selectedOption.toUpperCase()}</Text>
          <AccordionIcon
          color="complimentary.900"
          />
        </AccordionButton>
      </h2>
      <AccordionPanel 
  alignItems="center"
  justifyItems="center"
  color="white"
  pb={4}
>
  <VStack spacing={2} width="100%">
    <HStack justifyContent="space-between" width="100%">
      <Text
      color="white"
      >on {selectedOption.toUpperCase()}</Text>
      <Text
      color="complimentary.900"
      >0 {selectedOption.toUpperCase()}</Text>
    </HStack>
    <HStack justifyContent="space-between" width="100%">
      <Text
      color="white"
      >on Quicksilver</Text>
      <Text
      color="complimentary.900"
      >0 {selectedOption.toUpperCase()}</Text>
    </HStack>
  </VStack>
</AccordionPanel>
    </AccordionItem>

    <AccordionItem
    pt={4}
      borderBottom={"none"}
    >
      <h2>
      <AccordionButton 
      >
          <Flex 
          p={1}
          flexDirection="row" flex='1' alignItems="center">
            <Image src="/img/networks/q-atom.svg" boxSize="35px" mr={2} />
            <Text 
            fontSize="16px"
            color={"white"}>Liquid Staked</Text>
          </Flex>
          <Text
          pr={2}
          color="complimentary.900"
          >0 q{selectedOption.toUpperCase()}</Text>
          <AccordionIcon
          color="complimentary.900"
          />
        </AccordionButton>
      </h2>
      <AccordionPanel 
  alignItems="center"
  justifyItems="center"
  color="white"
  pb={4}
>
  <VStack spacing={2} width="100%">
    <HStack justifyContent="space-between" width="100%">
      <Text
      color="white"
      >on {selectedOption.toUpperCase()}</Text>
      <Text
      color="complimentary.900"
      >0 q{selectedOption.toUpperCase()}</Text>
    </HStack>
    <HStack justifyContent="space-between" width="100%">
      <Text
      color="white"
      >on Quicksilver</Text>
      <Text
      color="complimentary.900"
      >0 q{selectedOption.toUpperCase()}</Text>
    </HStack>
  </VStack>
</AccordionPanel>
    </AccordionItem>
  </Accordion>
  
                </Box>
              </Flex>
            </Flex>
          </Flex>
        </Container>
      </Box>
    </>
  );
}
