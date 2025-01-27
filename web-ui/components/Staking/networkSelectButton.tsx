// CustomMenu.tsx

import { ChevronDownIcon } from '@chakra-ui/icons';
import { Menu, MenuButton, MenuList, MenuItem, Button, Flex, Image, Text, useDisclosure } from '@chakra-ui/react';
import axios from 'axios';
import React, { useEffect, useState } from 'react';

import { Chain, local_chain, env, chains} from '@/config';

interface CustomMenuProps {
  buttonTextColor?: string;
  selectedOption: Chain|undefined;
  setSelectedNetwork: (network: Chain) => void;
}

export const NetworkSelect: React.FC<CustomMenuProps> = ({ buttonTextColor = 'white', selectedOption, setSelectedNetwork }) => {
  const handleOptionClick = (network: Chain) => {
    setSelectedNetwork(network);
  };

  function RotateIcon({ isOpen }: { isOpen: boolean }) {
    return (
      <ChevronDownIcon
        color="complimentary.700"
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
      const response = await axios.get(`${local_chain.get(env)?.rest[0]}/quicksilver/interchainstaking/v1/zones`);
      const liveZones = response.data.zones.map((zone: { chain_id: any }) => zone.chain_id);
      return liveZones;
    } catch (error) {
      console.error('Failed to fetch live zones:', error);
      return [];
    }
  };

  const [liveNetworks, setLiveNetworks] = useState<Chain[]>(Array.from(chains.get(env)?.values() ?? []));

  useEffect(() => {
    const getLiveZones = async () => {
      const liveZones = await fetchLiveZones();
      const filteredNetworks = Array.from(chains.get(env)?.values() ?? []).filter((network) => liveZones.includes(network.chain_id) && network.show == true);
      setLiveNetworks(filteredNetworks);
    };

    getLiveZones();
  }, []);

  return (
    <Menu>
      <MenuButton
        borderRadius={10}
        position="relative"
        zIndex={5}
        maxW="175px"
        minW="175px"
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
        leftIcon={<Image alt={selectedOption?.pretty_name} src={selectedOption?.logo} borderRadius={'full'} boxSize="24px" mr={1} />}
      >
        {selectedOption?.pretty_name.toUpperCase()}
      </MenuButton>
      <MenuList
  borderColor="rgba(35,35,35,1)"
  mt={1}
  pl={2}
  pr={2}
  bgColor="rgba(35,35,35,1)"
  display="grid"
  gridTemplateColumns="repeat(3, 1fr)" 
  gap={2} // optional: adds spacing between items
  zIndex={10}
>
  {liveNetworks.map((network) => (
    <MenuItem
      key={network.chain_id}
      py={4}
      bgColor="rgba(35,35,35,1)"
      borderRadius="4px"
      color="white"
      _hover={{
        bgColor: 'rgba(255,128,0, 0.25)',
      }}
      onClick={() => handleOptionClick(network)}
    >
      <Flex justifyContent="center" alignItems="center">
        <Image
          alt={network.chain_name}
          px={4}
          borderRadius="full"
          h="40px"
          src={network.logo}
        />
        <Text color="white" fontSize="20px" textAlign="center">
          {network.pretty_name}
        </Text>
      </Flex>
    </MenuItem>
  ))}
</MenuList>

    </Menu>
  );
};
