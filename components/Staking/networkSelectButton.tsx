// CustomMenu.tsx

import {
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  Button,
  Flex,
  Image,
  Text,
} from '@chakra-ui/react';
import React, { Dispatch, SetStateAction, useState } from 'react';
import { BsArrowDown } from 'react-icons/bs';

const networks = [
  {
    value: 'ATOM',
    logo: '/quicksilver-app-v2/img/networks/atom.svg',
    qlogo: '/quicksilver-app-v2/img/networks/q-atom.svg',
    name: 'Cosmos Hub',
    chainName: 'cosmoshub',
  },
  {
    value: 'OSMO',
    logo: '/quicksilver-app-v2/img/networks/osmosis.svg',
    qlogo: '/quicksilver-app-v2/img/networks/qosmo.svg',
    name: 'Osmosis',
    chainName: 'osmosis',
  },
  {
    value: 'STARS',
    logo: '/quicksilver-app-v2/img/networks/stargaze.svg',
    qlogo: '/quicksilver-app-v2/img/networks/stargaze-2.png',
    name: 'Stargaze',
    chainName: 'stargaze',
  },
  {
    value: 'REGEN',
    logo: '/quicksilver-app-v2/img/networks/regen.svg',
    qlogo: '/quicksilver-app-v2/img/networks/regen.svg',
    name: 'Regen',
    chainName: 'regen',
  },
  {
    value: 'SOMM',
    logo: '/quicksilver-app-v2/img/networks/sommelier.png',
    qlogo: '/quicksilver-app-v2/img/networks/sommelier.png',
    name: 'Sommelier',
    chainName: 'sommelier',
  },
];

interface CustomMenuProps {
  buttonTextColor?: string;
  selectedOption: (typeof networks)[0];
  setSelectedNetwork: (network: (typeof networks)[0]) => void;
}

export const NetworkSelect: React.FC<CustomMenuProps> = ({
  buttonTextColor = 'white',
  selectedOption,
  setSelectedNetwork,
}) => {
  const handleOptionClick = (network: (typeof networks)[0]) => {
    setSelectedNetwork(network);
  };
  return (
    <Menu>
      <MenuButton
        position="relative"
        zIndex={5}
        maxW="150px"
        minW="150px"
        variant="ghost"
        color="complimentary.900"
        backgroundColor="rgba(255,255,255,0.1)"
        _hover={{
          bgColor: 'rgba(255,255,255,0.05)',
          backdropFilter: 'blur(10px)',
        }}
        _active={{
          bgColor: 'rgba(255,255,255,0.05)',
          backdropFilter: 'blur(10px)',
        }}
        borderColor={buttonTextColor}
        as={Button}
        rightIcon={<BsArrowDown />}
      >
        {selectedOption.value.toUpperCase()}
      </MenuButton>
      <MenuList
        borderColor="rgba(35,35,35,1)"
        mt={1}
        bgColor="rgba(35,35,35,1)"
      >
        {networks.map((network) => (
          <MenuItem
            key={network.value}
            py={4}
            bgColor="rgba(35,35,35,1)"
            borderRadius="4px"
            color="complimentary.900"
            _hover={{
              bgColor: 'rgba(255,255,255,0.25)',
            }}
            onClick={() => handleOptionClick(network)}
          >
            <Flex
              justifyContent="center"
              alignItems="center"
              flexDirection="row"
            >
              <Image
                alt={network.name.toLowerCase()}
                px={4}
                h="40px"
                src={network.logo}
              />
              <Text
                color="complimentary.900"
                fontSize="20px"
                textAlign="center"
              >
                {network.name}
              </Text>
            </Flex>
          </MenuItem>
        ))}
      </MenuList>
    </Menu>
  );
};
