// CustomMenu.tsx

import { ChevronDownIcon } from '@chakra-ui/icons';
import {
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  Button,
  Flex,
  Image,
  Text,
  useDisclosure,
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
    qlogo: '/quicksilver-app-v2/img/networks/qstars.svg',
    name: 'Stargaze',
    chainName: 'stargaze',
  },
  {
    value: 'REGEN',
    logo: '/quicksilver-app-v2/img/networks/regen.svg',
    qlogo: '/quicksilver-app-v2/img/networks/q-regen.svg',
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

  function RotateIcon({ isOpen }: { isOpen: boolean }) {
    return (
      <ChevronDownIcon
        color="complimentary.900"
        transform={isOpen ? 'rotate(180deg)' : 'none'}
        transition="transform 0.2s"
        h="25px"
        w="25px"
      />
    );
  }

  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Menu>
      <MenuButton
        position="relative"
        zIndex={5}
        maxW="150px"
        minW="150px"
        _hover={{
          bgColor: 'rgba(255,128,0, 0.25)',
        }}
        px={2}
        color="white"
        as={Button}
        variant="outline"
        rightIcon={<RotateIcon isOpen={isOpen} />}
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
            color="white"
            _hover={{
              bgColor: 'rgba(255,128,0, 0.25)',
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
              <Text color="white" fontSize="20px" textAlign="center">
                {network.name}
              </Text>
            </Flex>
          </MenuItem>
        ))}
      </MenuList>
    </Menu>
  );
};
