import { Box, Image, Text, Accordion, AccordionItem, Flex, AccordionButton } from '@chakra-ui/react';
import React from 'react';

import { shiftDigits } from '@/utils';

type AssetsAccordianProps = {
  selectedOption: {
    name: string;
    value: string;
    logo: string;
    qlogo: string;
    chainName: string;
  };
  balance: string;
  qBalance: string;
};

export const AssetsAccordian: React.FC<AssetsAccordianProps> = ({ selectedOption, balance, qBalance }) => {
  const qAssets = shiftDigits(qBalance, -6);
  const qAssetsDisplay = qAssets.includes('.') ? qAssets.substring(0, qAssets.indexOf('.') + 3) : qAssets;
  const balanceDisplay = balance.includes('.') ? balance.substring(0, balance.indexOf('.') + 3) : balance;

  return (
    <Box position="relative" backdropFilter="blur(10px)" zIndex={10} borderRadius="10px" bgColor="rgba(255,255,255,0.1)" flex="1" p={5}>
      <Text fontSize="20px" color="white">
        Assets
      </Text>
      <Accordion mt={6} allowToggle>
        <AccordionItem mb={4} borderTop={'none'}>
          <h2>
            <AccordionButton _hover={{ cursor: 'default' }} borderRadius={'10px'} borderTopColor={'transparent'}>
              <Flex p={1} flexDirection="row" flex="1" alignItems="center">
                <Image alt="atom" src={selectedOption.logo} boxSize="35px" mr={2} />
                <Text fontSize="16px" color={'white'}>
                  Available to stake
                </Text>
              </Flex>
              <Text pr={2} color="complimentary.900">
                {balanceDisplay}
              </Text>
              <Text pr={2} color="complimentary.900">
                {selectedOption.value.toUpperCase()}
              </Text>
            </AccordionButton>
          </h2>
        </AccordionItem>

        <AccordionItem pt={4} borderBottom={'none'}>
          <h2>
            <AccordionButton _hover={{ cursor: 'default' }} borderRadius={'10px'} borderTopColor={'transparent'}>
              <Flex p={1} flexDirection="row" flex="1" alignItems="center">
                <Image alt="qAtom" src={selectedOption.qlogo} boxSize="35px" mr={2} />
                <Text fontSize="16px" color={'white'}>
                  Liquid Staked
                </Text>
              </Flex>
              <Text pr={2} color="complimentary.900">
                {qAssetsDisplay}
              </Text>
              <Text pr={2} color="complimentary.900">
                q{selectedOption.value.toUpperCase()}
              </Text>
            </AccordionButton>
          </h2>
        </AccordionItem>
      </Accordion>
    </Box>
  );
};
