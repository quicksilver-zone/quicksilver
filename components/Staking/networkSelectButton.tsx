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
import React from 'react';
import { BsArrowDown } from 'react-icons/bs';

interface CustomMenuProps {
  buttonTextColor?: string;
  selectedOption: string;
  setSelectedOption: (value: string) => void;
}

export const NetworkSelect: React.FC<
  CustomMenuProps
> = ({
  buttonTextColor = 'white',
  selectedOption,
  setSelectedOption,
}) => {
  const handleOptionClick = (option: string) => {
    setSelectedOption(option);
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
      <MenuList mt={1} bgColor="rgba(35,35,35,1)">
        <MenuItem
          py={4}
          bgColor="rgba(35,35,35,1)"
          borderRadius="4px"
          color="complimentary.900"
          _hover={{
            bgColor: 'rgba(255,255,255,0.25)',
          }}
          onClick={() =>
            handleOptionClick('ATOM')
          }
        >
          <Flex
            justifyContent="center"
            alignItems="center"
            flexDirection="row"
          >
            <Image
              alt="atom"
              px={4}
              h="40px"
              src="/img/networks/atom.svg"
            />
            <Text
              color="complimentary.900"
              fontSize="20px"
              textAlign="center"
            >
              Cosmos Hub
            </Text>
          </Flex>
        </MenuItem>
        <MenuItem
          py={4}
          bgColor="rgba(35,35,35,1)"
          borderRadius="4px"
          color="complimentary.900"
          _hover={{
            bgColor: 'rgba(255,255,255,0.25)',
          }}
          onClick={() =>
            handleOptionClick('OSMO')
          }
        >
          <Flex
            justifyContent="center"
            alignItems="center"
            flexDirection="row"
          >
            <Image
              alt="osmosis"
              px={4}
              h="40px"
              src="/img/networks/osmosis.svg"
            />
            <Text
              color="complimentary.900"
              fontSize="20px"
              textAlign="center"
            >
              Osmosis
            </Text>
          </Flex>
        </MenuItem>
        <MenuItem
          py={4}
          bgColor="rgba(35,35,35,1)"
          borderRadius="4px"
          color="complimentary.900"
          _hover={{
            bgColor: 'rgba(255,255,255,0.25)',
          }}
          onClick={() =>
            handleOptionClick('REGEN')
          }
        >
          <Flex
            justifyContent="center"
            alignItems="center"
            flexDirection="row"
          >
            <Image
              alt="regen"
              px={4}
              h="40px"
              src="/img/networks/regen.svg"
            />
            <Text
              color="complimentary.900"
              fontSize="20px"
              textAlign="center"
            >
              Regen
            </Text>
          </Flex>
        </MenuItem>
        <MenuItem
          py={4}
          bgColor="rgba(35,35,35,1)"
          borderRadius="4px"
          color="complimentary.900"
          _hover={{
            bgColor: 'rgba(255,255,255,0.25)',
          }}
          onClick={() =>
            handleOptionClick('STARS')
          }
        >
          <Flex
            justifyContent="center"
            alignItems="center"
            flexDirection="row"
          >
            <Image
              alt="stargaze"
              px={4}
              h="40px"
              src="/img/networks/stargaze.svg"
            />
            <Text
              color="complimentary.900"
              fontSize="20px"
              textAlign="center"
            >
              Stargaze
            </Text>
          </Flex>
        </MenuItem>
        <MenuItem
          py={4}
          bgColor="rgba(35,35,35,1)"
          borderRadius="4px"
          color="complimentary.900"
          _hover={{
            bgColor: 'rgba(255,255,255,0.25)',
          }}
          onClick={() =>
            handleOptionClick('SOMM')
          }
        >
          <Flex
            justifyContent="center"
            alignItems="center"
            flexDirection="row"
          >
            <Image
              alt="somm"
              px={4}
              h="40px"
              src="/img/networks/sommelier.png"
              borderRadius="50%"
            />
            <Text
              color="complimentary.900"
              fontSize="20px"
              textAlign="center"
            >
              Sommelier
            </Text>
          </Flex>
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
