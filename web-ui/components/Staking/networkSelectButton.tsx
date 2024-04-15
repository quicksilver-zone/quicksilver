// CustomMenu.tsx

import { ChevronDownIcon } from '@chakra-ui/icons';
import { Menu, MenuButton, MenuList, MenuItem, Button, Flex, Image, Text, useDisclosure } from '@chakra-ui/react';
import axios from 'axios';
import React, { useEffect, useState } from 'react';

import { networks as prodNetworks, testNetworks as devNetworks } from '@/state/chains/prod';

const networks = process.env.NEXT_PUBLIC_CHAIN_ENV === 'mainnet' ? prodNetworks : devNetworks;

interface CustomMenuProps {
  buttonTextColor?: string;
  selectedOption: (typeof networks)[0];
  setSelectedNetwork: (network: (typeof networks)[0]) => void;
}

type Network = {
  value: string;
  logo: string;
  qlogo: string;
  name: string;
  chainName: string;
  chainId: string;
};

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

  const { isOpen } = useDisclosure();

  const fetchLiveZones = async () => {
    try {
      const response = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_API}/quicksilver/interchainstaking/v1/zones`);
      const liveZones = response.data.zones.map((zone: { chain_id: any }) => zone.chain_id);
      return liveZones;
    } catch (error) {
      console.error('Failed to fetch live zones:', error);
      return [];
    }
  };

  const [liveNetworks, setLiveNetworks] = useState<Network[]>([]);

  useEffect(() => {
    const getLiveZones = async () => {
      const liveZones = await fetchLiveZones();
      const filteredNetworks = networks.filter((network) => liveZones.includes(network.chainId));
      setLiveNetworks(filteredNetworks);
    };

    getLiveZones();
  }, []);

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
        _active={{
          bgColor: 'rgba(255,128,0, 0.25)',
        }}
        _focus={{
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
        {liveNetworks.map((network) => (
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
              <Image alt={network.name.toLowerCase()} px={4} borderRadius={'full'} h="40px" src={network.logo} />
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
