import { HamburgerIcon, ArrowBackIcon } from '@chakra-ui/icons';
import {
  Flex,
  Box,
  Image,
  Spacer,
  VStack,
  IconButton,
  Tooltip,
  ScaleFade,
  useDisclosure,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  Link,
  HStack,
} from '@chakra-ui/react';
import { keyframes } from '@emotion/react';
import { useRouter } from 'next/router';
import { useState, useEffect } from 'react';
import { FaDiscord, FaTwitter, FaGithub, FaInfo } from 'react-icons/fa';
import { IoIosDocument } from 'react-icons/io';

import { DrawerControlProvider } from '@/state/chains/drawerControlProvider';

import { WalletButton } from '../wallet-button';


export const SideHeader = () => {
  const router = useRouter();
  const [selectedPage, setSelectedPage] = useState('');

  const [showSocialLinks, setShowSocialLinks] = useState(false);

  useEffect(() => {
    const handleRouteChange = (url: string) => {
      const path = url.split('/')[1];
      setSelectedPage(path);
    };

    router.events.on('routeChangeComplete', handleRouteChange);
    return () => router.events.off('routeChangeComplete', handleRouteChange);
  }, [router]);

  const commonBoxShadowColor = 'rgba(255, 128, 0, 0.25)';
  const toggleSocialLinks = () => setShowSocialLinks(!showSocialLinks);

  const [isMobile, setIsMobile] = useState(typeof window !== 'undefined' && window.innerWidth < 1274);

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 1274);
    };

    window.addEventListener('resize', handleResize);

    // Clean up
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);
  const transitionStyle = 'all 0.3s ease';

  const handleLogoClick = () => {
    if (isMobile) {
      DrawerOnOpen();
    } else {
      router.push('/staking');
    }
  };

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

  const { isOpen: DrawerIsOpen, onOpen: DrawerOnOpen, onClose: DrawerOnClose } = useDisclosure();

  return (
    <Box
      w={isMobile ? 'auto' : 'fit-content'}
      h={isMobile ? 'fit-content' : '95vh'}
      backdropFilter="blur(10px)"
      borderRadius={isMobile ? 'full' : 100}
      zIndex={10}
      top={6}
      left={6}
      position="fixed"
      bgColor="rgba(214, 219, 220, 0.1)"
    >
      <Flex direction="column" align="center" zIndex={10} justifyContent="space-between" height="100%">
        <Image
          alt="logo"
          h="75px"
          w="75px"
          padding="5px"
          borderRadius="full"
          src="/img/networks/quicksilver.svg"
          onClick={handleLogoClick}
          cursor="pointer"
          _hover={{
            ...(isMobile && {
              animation: `${shadowKeyframes} 3s linear infinite`,
              transform: 'scale(1.05)',
              transition: 'transform 0.3s ease',
            }),
          }}
        />
        <DrawerControlProvider closeDrawer={DrawerOnClose}>
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
                    <Link
                      href={`/${item.toLowerCase()}`}
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
                    </Link>
                  </Box>
                ))}
                <Box mt={12} position="relative"></Box>
                <HStack mt={'50px'} alignContent={'center'} justifyContent={'space-around'}>
                  <Box
                    _hover={{
                      cursor: 'pointer',
                      boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                    }}
                  >
                    <FaGithub size={'25px'} color="rgb(255, 128, 0)" />
                  </Box>
                  <Box
                    _hover={{
                      cursor: 'pointer',
                      boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                      transition: transitionStyle,
                    }}
                  >
                    <FaDiscord size={'25px'} color="rgb(255, 128, 0)" />
                  </Box>
                  <Box
                    _hover={{
                      cursor: 'pointer',
                      boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                      transition: transitionStyle,
                    }}
                  >
                    <FaTwitter size={'25px'} color="rgb(255, 128, 0)" />
                  </Box>
                </HStack>
              </DrawerBody>
            </DrawerContent>
          </Drawer>
        </DrawerControlProvider>

        {!isMobile && (
          <>
            <Spacer />
            <ScaleFade initialScale={0.5} in={!showSocialLinks}>
              {!showSocialLinks && (
                <VStack justifyContent="center" alignItems="center" spacing={16}>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Staking" placement="right">
                    <Box
                      w="60px"
                      h="60px"
                      onClick={() => router.push('/staking')}
                      cursor="pointer"
                      borderRadius="100px"
                      boxShadow={
                        selectedPage === 'staking' ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}` : ''
                      }
                      _hover={{
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <Image
                        filter={selectedPage === 'staking' ? 'contrast(100%)' : 'contrast(50%)'}
                        _hover={{
                          filter: 'contrast(100%)',
                        }}
                        alt="Staking"
                        h="60px"
                        w="60px"
                        src="/img/liquid.png"
                      />
                    </Box>
                  </Tooltip>

                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Governance" placement="right">
                    <Box
                      w="60px"
                      h="60px"
                      onClick={() => router.push('/governance')}
                      cursor="pointer"
                      borderRadius="100px"
                      boxShadow={
                        selectedPage === 'governance'
                          ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`
                          : ''
                      }
                      _hover={{
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <Image
                        filter={selectedPage === 'governance' ? 'contrast(100%)' : 'contrast(50%)'}
                        _hover={{
                          filter: 'contrast(100%)',
                        }}
                        alt="Governance"
                        h="60px"
                        w="65px"
                        src="/img/governance.png"
                      />
                    </Box>
                  </Tooltip>

                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Assets" placement="right">
                    <Box
                      w="55px"
                      h="55px"
                      onClick={() => router.push('/assets')}
                      cursor="pointer"
                      borderRadius="100px"
                      boxShadow={
                        selectedPage === 'assets' ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}` : ''
                      }
                      _hover={{
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <Image
                        filter={selectedPage === 'assets' ? 'contrast(100%)' : 'contrast(50%)'}
                        _hover={{
                          filter: 'contrast(100%)',
                        }}
                        alt="Assets"
                        h="55px"
                        src="/img/assets.png"
                      />
                    </Box>
                  </Tooltip>
                  {/*<Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Airdrop" placement="right">
                    <Box
                      w="55px"
                      h="55px"
                      onClick={() => router.push('/airdrop')}
                      cursor="pointer"
                      borderRadius="100px"
                      boxShadow={
                        selectedPage === 'airdrop' ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}` : ''
                      }
                      _hover={{
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <Image
                        filter={selectedPage === 'airdrop' ? 'contrast(100%)' : 'contrast(50%)'}
                        _hover={{
                          filter: 'contrast(100%)',
                        }}
                        alt="DeFi"
                        h="55px"
                        src="/img/airdrop.png"
                      />
                    </Box>
                      </Tooltip>*/}

                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="DeFi" placement="right">
                    <Box
                      w="55px"
                      h="55px"
                      onClick={() => router.push('/defi')}
                      cursor="pointer"
                      borderRadius="100px"
                      boxShadow={
                        selectedPage === 'defi' ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}` : ''
                      }
                      _hover={{
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <Image
                        filter={selectedPage === 'defi' ? 'contrast(100%)' : 'contrast(50%)'}
                        _hover={{
                          filter: 'contrast(100%)',
                        }}
                        alt="DeFi"
                        h="55px"
                        src="/img/defi.png"
                      />
                    </Box>
                  </Tooltip>
                </VStack>
              )}
            </ScaleFade>

            <ScaleFade initialScale={0.5} in={showSocialLinks}>
              {showSocialLinks && (
                <VStack justifyContent="center" alignItems="center" spacing={16}>
                  <Link href="https://quicksilver.zone/" isExternal>
                    <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="About" placement="right">
                      <Box
                        _hover={{
                          cursor: 'pointer',
                          boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        }}
                      >
                        <FaInfo size={'25px'} color="rgb(255, 128, 0)" />
                      </Box>
                    </Tooltip>
                  </Link>
                  <Link href="https://docs.quicksilver.zone/" isExternal>
                    <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Docs" placement="right">
                      <Box
                        _hover={{
                          cursor: 'pointer',
                          boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                          transition: transitionStyle,
                        }}
                      >
                        <IoIosDocument size={'25px'} color="rgba(255, 128, 0, 0.9)" />
                      </Box>
                    </Tooltip>
                  </Link>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Discord" placement="right">
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
                  </Tooltip>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Github" placement="right">
                    <Link href="https://github.com/quicksilver-zone/quicksilver" isExternal>
                      <Box
                        _hover={{
                          cursor: 'pointer',
                          boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                          transition: transitionStyle,
                        }}
                      >
                        <FaGithub size={'25px'} color="rgb(255, 128, 0)" />
                      </Box>
                    </Link>
                  </Tooltip>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Twitter" placement="right">
                    <Link href="https://twitter.com/quicksilverzone" isExternal>
                      <Box
                        _hover={{
                          cursor: 'pointer',
                          boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                          transition: transitionStyle,
                        }}
                      >
                        <FaTwitter size={'25px'} color="rgb(255, 128, 0)" />
                      </Box>
                    </Link>
                  </Tooltip>
                  {/*<Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Privacy Policy" placement="right">
                    <Box
                      onClick={() => router.push('/privacy-policy')}
                      _hover={{
                        cursor: 'pointer',
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 15px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <MdPrivacyTip size={'25px'} color="rgb(255, 128, 0)" />
                    </Box>
                    </Tooltip>*/}
                </VStack>
              )}
            </ScaleFade>
          </>
        )}

        <Spacer />
        {!isMobile && (
          <IconButton
            borderRadius={'100'}
            icon={showSocialLinks ? <ArrowBackIcon /> : <HamburgerIcon />}
            aria-label="Toggle View"
            onClick={toggleSocialLinks}
            mb={4}
            _hover={{
              bgColor: 'complimentary.500',
            }}
          />
        )}
      </Flex>
    </Box>
  );
};

export default SideHeader;
