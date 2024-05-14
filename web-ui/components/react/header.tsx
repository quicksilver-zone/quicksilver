import {
  Flex,
  Spacer,
  Box,
  keyframes,
  Image,
  useDisclosure,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  HStack,
  Link,
  Text,
} from '@chakra-ui/react';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/router';
import { FaGithub, FaDiscord } from 'react-icons/fa';
import { FaXTwitter } from 'react-icons/fa6';

import { WalletButton } from '../wallet-button';

const shadowKeyframes = keyframes`
  0% {
    box-shadow: 0 0 10px 5px #FF8000;
  }
  25% {
    box-shadow: 0 0 10px 5px #FF9933;
  }
  50% {
    box-shadow: 0 0 10px 5px #FFB266;
  }
  75% {
    box-shadow: 0 0 10px 5px #FF9933;
  }
  100% {
    box-shadow: 0 0 10px 5px #FF8000;
  }
`;

const commonBoxShadowColor = 'rgba(255, 128, 0, 0.25)';
const transitionStyle = 'all 0.3s ease';

const Header: React.FC<{ chainName: string }> = ({ chainName }) => {
  const router = useRouter();
  const { isOpen: DrawerIsOpen, onOpen: DrawerOnOpen, onClose: DrawerOnClose } = useDisclosure();
  const handleLogoClick = () => {
    DrawerOnOpen();
  };

  return (
    <Flex alignItems="center" zIndex={50} justifyContent="space-between" position={'relative'} p={4}>
      <Spacer display={{ base: 'none', menu: 'block' }} />
      <Box display={{ base: 'block', menu: 'none' }}>
        <Image
          alt="logo"
          h="70px"
          w="70px"
          padding="3px"
          borderRadius="full"
          src="/img/networks/quicksilver.svg"
          onClick={handleLogoClick}
          cursor="pointer"
          background={'#4a4a4a3f'}
          _hover={{
            animation: `${shadowKeyframes} 3s linear infinite`,
            transform: 'scale(1.05)',
            transition: 'transform 0.3s ease',
          }}
        />
      </Box>

      <WalletButton />
      <Drawer isOpen={DrawerIsOpen} placement="left" onClose={DrawerOnClose}>
        <DrawerOverlay />
        <DrawerContent bgColor="rgba(32,32,32,1)">
          <DrawerCloseButton color="white" />
          <DrawerHeader fontSize="3xl" letterSpacing={4} lineHeight={2} color="white">
            QUICKSILVER
          </DrawerHeader>
          <DrawerBody>
            {[/*'Airdrop', */ 'Assets', 'Defi', 'Governance', 'Staking'].map((item) => (
              <Box key={item} mb={4} position="relative">
                <Text
                  cursor={'pointer'}
                  onClick={() => router.push(`/${item.toLowerCase()}`)}
                  fontSize="xl"
                  fontWeight="medium"
                  color="white"
                  position="relative"
                  _hover={{
                    textDecoration: 'none',
                    color: 'transparent',
                    backgroundClip: 'text',
                    bgGradient: 'linear(to-r, #FF8000, #FF9933, #FFB266, #FFD9B3, #FFE6CC)',
                    _before: {
                      width: '100%',
                    },
                  }}
                  _before={{
                    content: `""`,
                    position: 'absolute',
                    bottom: '-2px',
                    left: '0',
                    width: '0',
                    height: '2px',
                    bgGradient: 'linear(to-r, #FF8000, #FF9933, #FFB266, #FFD9B3, #FFE6CC)',
                    transition: 'width 0.4s ease',
                  }}
                >
                  {item}
                </Text>
              </Box>
            ))}
            <Box mt={12} position="relative"></Box>
            <HStack mt={'50px'} alignContent={'center'} justifyContent={'space-around'}>
              <Link href="https://github.com/quicksilver-zone/quicksilver" isExternal>
                <Box
                  _hover={{
                    cursor: 'pointer',
                    boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                  }}
                >
                  <FaGithub size={'25px'} color="rgb(255, 128, 0)" />
                </Box>
              </Link>
              <Link href="https://discord.com/invite/xrSmYMDVrQ" isExternal>
                <Box
                  _hover={{
                    cursor: 'pointer',
                    boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                    transition: transitionStyle,
                  }}
                >
                  <FaDiscord size={'25px'} color="rgb(255, 128, 0)" />
                </Box>
              </Link>
              <Link href="https://twitter.com/quicksilverzone" isExternal>
                <Box
                  _hover={{
                    cursor: 'pointer',
                    boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                    transition: transitionStyle,
                  }}
                >
                  <FaXTwitter size={'25px'} color="rgb(255, 128, 0)" />
                </Box>
              </Link>
            </HStack>
          </DrawerBody>
        </DrawerContent>
      </Drawer>
    </Flex>
  );
};

export const DynamicHeaderSection = dynamic(() => Promise.resolve(Header), {
  ssr: false,
});
