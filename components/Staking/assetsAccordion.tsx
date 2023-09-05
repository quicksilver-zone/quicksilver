import React from 'react';
import { 
  Box, 
  Image, 
  Text, 
  Accordion, 
  AccordionItem, 
  Flex, 
  AccordionPanel, 
  VStack, 
  HStack, 
  AccordionButton,
  AccordionIcon 
} from "@chakra-ui/react";

const AssetsAccordian = ({ selectedOption, activeAccordion, handleAccordionChange, openItem }) => {
    return (

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
        <AccordionButton 
        borderRadius={"10px"}
        _hover={{
          bgColor: "rgba(0,0,0,0.05)", 
          backdropFilter: "blur(10px)",
          
        }}
        _active={{
          bgColor: "rgba(0,0,0,0.05)",
          backdropFilter: "blur(10px)",
          
        }}
        borderTopColor={"transparent"}>
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
        borderRadius={"10px"}
        _hover={{
          bgColor: "rgba(0,0,0,0.05)", 
          backdropFilter: "blur(10px)",
          
        }}
        _active={{
          bgColor: "rgba(0,0,0,0.05)",
          backdropFilter: "blur(10px)",
          
        }}
        borderTopColor={"transparent"}>
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
                 );
                };
                
                export default AssetsAccordian;