import { Box, Image, Text, Accordion, AccordionItem, Flex, AccordionPanel, VStack, HStack, Link } from '@chakra-ui/react';
import React from 'react';
import { BsTrophy, BsCoin, BsClock } from 'react-icons/bs';
import { RiStockLine } from 'react-icons/ri';

import { useZoneQuery } from '@/hooks/useQueries';

type AssetsAccordianProps = {
  displayApr: string;
  selectedOption: {
    name: string;
    value: string;
    chainId: string;
  };
};

export const InfoBox: React.FC<AssetsAccordianProps> = ({ selectedOption, displayApr }) => {
  const { data: zone, isLoading: isZoneLoading, isError: isZoneError } = useZoneQuery(selectedOption.chainId);
  const redemptionRate = zone?.redemptionRate;
  const unbondingPeriod = (Number(zone?.unbondingPeriod) / 86400000000000).toString() + ' days';
  return (
    <Box zIndex={2} position="relative" backdropFilter="blur(30px)" borderRadius="10px" bgColor="rgba(255,255,255,0.1)" flex="2" p={5}>
      {/* <Image
        alt="embelish"
        src="/quicksilver-app-v2/img/metalmisc3.png"
        zIndex={1}
        position="absolute"
        top="-40px"
        right="-65px"
        boxSize="135px"
        transform="rotate(25deg)"
      /> */}
      <Text fontSize="20px" color="white">
        {selectedOption.value.toUpperCase()}&nbsp;on Quicksilver
      </Text>
      <Accordion mt={6} allowToggle>
        <AccordionItem pt={2} mb={2} borderTop={'none'}>
          <h2>
            <Flex borderTopColor={'transparent'} alignItems="center" justifyContent="space-between" width="100%" py={2}>
              <Flex flexDirection="row" alignItems="center">
                <Box mr="16px">
                  <BsTrophy color="#FF8000" size="24px" />
                </Box>
                <Text fontSize="16px" color={'white'}>
                  Rewards
                </Text>
              </Flex>
              <Text pr={2} color="complimentary.900">
                {displayApr}
              </Text>
            </Flex>
          </h2>
        </AccordionItem>

        <AccordionItem pt={2} mb={2}>
          <h2>
            <Flex borderTopColor={'transparent'} alignItems="center" justifyContent="space-between" width="100%" py={2}>
              <Flex flexDirection="row" flex="1" alignItems="center">
                <Box mr="16px">
                  {' '}
                  {/* Adjusts right margin */}
                  <BsCoin color="#FF8000" size="24px" />
                </Box>
                <Text fontSize="16px" color={'white'}>
                  Fees
                </Text>
              </Flex>
              <Text pr={2} color="complimentary.900">
                Low
              </Text>
            </Flex>
          </h2>
        </AccordionItem>
        <AccordionItem pt={2} mb={2}>
          <h2>
            <Flex borderTopColor={'transparent'} alignItems="center" justifyContent="space-between" width="100%" py={2}>
              <Flex flexDirection="row" flex="1" alignItems="center">
                <Box mr="16px">
                  {' '}
                  {/* Adjusts right margin */}
                  <BsClock color="#FF8000" size="24px" />
                </Box>
                <Text fontSize="16px" color={'white'}>
                  Unbonding
                </Text>
              </Flex>
              <Text pr={2} color="complimentary.900">
                {unbondingPeriod}
              </Text>
            </Flex>
          </h2>
          <AccordionPanel alignItems="center" justifyItems="center" color="white" pb={4}>
            <VStack spacing={2} width="100%">
              <HStack justifyContent="space-between" width="100%">
                <Text color="white">on {selectedOption.value.toUpperCase()}</Text>
                <Text color="complimentary.900">0 {selectedOption.value.toUpperCase()}</Text>
              </HStack>
              <HStack justifyContent="space-between" width="100%">
                <Text color="white">on Quicksilver</Text>
                <Text color="complimentary.900">0 {selectedOption.value.toUpperCase()}</Text>
              </HStack>
            </VStack>
          </AccordionPanel>
        </AccordionItem>
        <AccordionItem pt={2} mb={2} borderBottom={'none'}>
          <h2>
            <Flex borderTopColor={'transparent'} alignItems="center" justifyContent="space-between" width="100%" py={2}>
              <Flex flexDirection="row" flex="1" alignItems="center">
                <Box mr="16px">
                  {' '}
                  {/* Adjusts right margin */}
                  <RiStockLine color="#FF8000" size="24px" />
                </Box>
                <Text fontSize="16px" color={'white'}>
                  Value of 1 q{selectedOption.value.toUpperCase()}
                </Text>
              </Flex>
              <Text pr={2} color="complimentary.900">
                1 q{selectedOption.value.toUpperCase()} = {Number(redemptionRate).toFixed(2).toString()}{' '}
                {selectedOption.value.toUpperCase()}
              </Text>
            </Flex>
          </h2>
        </AccordionItem>
      </Accordion>

      <Text mt={3} color="white" textAlign="center" bgColor="rgba(0,0,0,0.4)" p={5} width="100%" borderRadius={6} fontWeight="light">
        Want to learn more about rewards, fees, and unbonding on Quicksilver?&nbsp;Check out the{' '}
        <Link href="https://your-docs-url.com" color="complimentary.900" isExternal>
          docs
        </Link>
        .
      </Text>
    </Box>
  );
};
