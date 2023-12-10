// CustomMenu.tsx

import { ChevronDownIcon } from '@chakra-ui/icons';
import { Menu, MenuButton, MenuList, MenuItem, Button, Flex, Image, Text, useDisclosure } from '@chakra-ui/react';
import React, { Dispatch, SetStateAction, useState } from 'react';
import { BsArrowDown } from 'react-icons/bs';

import { networks } from '@/state/chains/prod';

interface CustomMenuProps {
  buttonTextColor?: string;
  selectedOption: (typeof networks)[0];
  setSelectedNetwork: (network: (typeof networks)[0]) => void;
}

export const NetworkSelect: React.FC<CustomMenuProps> = ({ buttonTextColor = 'white', selectedOption, setSelectedNetwork }) => {
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
        borderRadius={100}
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
      <MenuList borderColor="rgba(35,35,35,1)" mt={1} bgColor="rgba(35,35,35,1)">
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
            <Flex justifyContent="center" alignItems="center" flexDirection="row">
              <Image alt={network.name.toLowerCase()} px={4} h="40px" src={network.logo} />
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
