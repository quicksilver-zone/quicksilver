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
import React, { Dispatch, SetStateAction } from 'react';
import { BsArrowDown } from 'react-icons/bs';

const networks = [
  {
    value: 'ATOM',
    logo: '/img/networks/atom.svg',
    name: 'Cosmos Hub',
    chainName: 'cosmoshub',
  },
  {
    value: 'OSMO',
    logo: '/img/networks/osmosis.svg',
    name: 'Osmosis',
    chainName: 'osmosis',
  },
  {
    value: 'STARS',
    logo: '/img/networks/stargaze.svg',
    name: 'Stargaze',
    chainName: 'stargaze',
  },
  {
    value: 'REGEN',
    logo: '/img/networks/regen.svg',
    name: 'Regen',
    chainName: 'regen',
  },
  {
    value: 'SOMM',
    logo: '/img/networks/sommelier.png',
    name: 'Sommelier',
    chainName: 'sommelier',
  },
];

interface CustomMenuProps {
  buttonTextColor?: string;
  selectedOption: string;
  selectedChainName: string;
  setSelectedOption: (value: string) => void;
  setSelectedChainName: (value: string) => void;
}

export const NetworkSelect: React.FC<CustomMenuProps> = ({
  buttonTextColor = 'white',
  selectedOption,
  setSelectedOption,
  setSelectedChainName,
  selectedChainName,
}) => {
  const handleOptionClick = (option: string) => {
    setSelectedOption(option);
    const selectedNetwork = networks.find((net) => net.value === option);
    if (selectedNetwork) {
      setSelectedChainName(selectedNetwork.chainName);
    }
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
        {selectedOption.toUpperCase()}
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
            onClick={() => handleOptionClick(network.value)}
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
