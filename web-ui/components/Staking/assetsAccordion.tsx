import { Box, Image, Text, Accordion, AccordionItem, Flex, AccordionButton, SkeletonCircle } from '@chakra-ui/react';
import React, { useEffect, useState } from 'react';

import { Chain, env, getChainForFieldValue} from '@/config';
import { useCurrentInterchainAssetsQuery } from '@/hooks/useQueries';

const BigNumber = require('bignumber.js');

type AssetsAccordianProps = {
  selectedOption: Chain;
  balance: string;
  qBalance: string;
  address: string;
};

export const AssetsAccordian: React.FC<AssetsAccordianProps> = ({ selectedOption, balance, qBalance, address }) => {

  const { assets: liquidRewards } = useCurrentInterchainAssetsQuery(address);
  const [updatedQBalance, setUpdatedQBalance] = useState(qBalance);

  useEffect(() => {
    const calculateLiquidRewards = () => {
      let totalAmount = new BigNumber(0);
      
      const denomToFind = getChainForFieldValue(env, "chain_id",selectedOption.chain_id)?.q_denom;
      for (const chain in liquidRewards?.assets) {
        const chainAssets = liquidRewards?.assets[chain];
        chainAssets.forEach((assetGroup) => {
          if (assetGroup.Type === "liquid") {
            assetGroup.Amount.forEach((asset) => {
              if (asset.denom === denomToFind) {
                totalAmount = totalAmount.plus(asset.amount);
              }
            });
          }
        });
      }

      return totalAmount.shiftedBy(-selectedOption.exponent).toString();
    };

    setUpdatedQBalance(calculateLiquidRewards());
  }, [selectedOption, liquidRewards, qBalance]);

  const qAssetsDisplay = updatedQBalance.includes('.') ? updatedQBalance.substring(0, updatedQBalance.indexOf('.') + 3) : updatedQBalance;
  const balanceDisplay = balance.includes('.') ? balance.substring(0, balance.indexOf('.') + 4) : balance;

  const renderQAssets = () => {
    if (qBalance && liquidRewards) {
      return (
        <Text pr={2} color="complimentary.700">
          {qAssetsDisplay}
        </Text>
      );
    } else {
      return (
        <Box mr={2} display="inline-block">
          <SkeletonCircle size="2" startColor="complimentary.700" endColor="complimentary.400" />
        </Box>
      );
    }
  };

  const renderAssets = () => {
    if (balance) {
      return (
        <Text pr={2} color="complimentary.700">
          {balanceDisplay}
        </Text>
      );
    } else {
      return (
        <Box mr={2} display="inline-block">
          <SkeletonCircle size="2" startColor="complimentary.700" endColor="complimentary.400" />
        </Box>
      );
    }
  };

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
                <Image alt={selectedOption.pretty_name} src={selectedOption.logo} borderRadius={'full'} boxSize="35px" mr={2} />
                <Text fontSize="16px" color={'white'}>
                  Available to stake
                </Text>
              </Flex>
              {renderAssets()}
              <Text pr={2} color="complimentary.700">
                {selectedOption.major_denom.toUpperCase()}
              </Text>
            </AccordionButton>
          </h2>
        </AccordionItem>

        <AccordionItem pt={4} borderBottom={'none'}>
          <h2>
            <AccordionButton _hover={{ cursor: 'default' }} borderRadius={'10px'} borderTopColor={'transparent'}>
              <Flex p={1} flexDirection="row" flex="1" alignItems="center">
                <Image alt={`q${selectedOption.major_denom.toUpperCase()}`} borderRadius={'full'} src={selectedOption.qlogo} boxSize="35px" mr={2} />
                <Text fontSize="16px" color={'white'}>
                  Liquid Staked
                </Text>
              </Flex>

              {renderQAssets()}
              <Text pr={2} color="complimentary.700">
                q{selectedOption.major_denom.toUpperCase()}
              </Text>
            </AccordionButton>
          </h2>
        </AccordionItem>
      </Accordion>
    </Box>
  );
};
