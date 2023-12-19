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
  useBreakpointValue,
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
import { MdPrivacyTip } from 'react-icons/md';

export const SideHeader = () => {
  const router = useRouter();
  const [selectedPage, setSelectedPage] = useState('');
  const [showSocialLinks, setShowSocialLinks] = useState(false);

  useEffect(() => {
    const handleRouteChange = (url: string) => {
      const path = url.split('/quicksilver-app-v2/')[1];
      setSelectedPage(path);
    };

    router.events.on('routeChangeComplete', handleRouteChange);
    return () => router.events.off('routeChangeComplete', handleRouteChange);
  }, [router]);

  const commonBoxShadowColor = 'rgba(255, 128, 0, 0.25)';
  const toggleSocialLinks = () => setShowSocialLinks(!showSocialLinks);

  const isMobile = useBreakpointValue({ base: true, md: false });
  const transitionStyle = 'all 0.3s ease';
  const { isOpen, onOpen, onClose } = useDisclosure();

  const handleLogoClick = () => {
    if (isMobile) {
      onOpen();
    } else {
      router.push('/');
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

  return (
    <Box
      w={isMobile ? 'auto' : 'fit-content'}
      h={{ base: 'fit-content', md: '95vh' }}
      backdropFilter="blur(10px)"
      borderRadius={{ base: 'full', md: 100 }}
      zIndex={10}
      top={6}
      left={6}
      position="fixed"
      bgColor="rgba(214, 219, 220, 0.1)"
    >
      <Flex direction="column" align="center" zIndex={10} justifyContent="space-between" py={{ base: 0, md: 4 }} height="100%">
        <Image
          alt="logo"
          mt={{ base: 0, md: '-10px' }}
          h="75px"
          w="75px"
          borderRadius="full"
          src="/quicksilver-app-v2/img/networks/quicksilver.svg"
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

        <Drawer isOpen={isOpen} placement="left" onClose={onClose}>
          <DrawerOverlay />
          <DrawerContent bgColor="black">
            <DrawerCloseButton color="white" />
            <DrawerHeader textDecoration={'underline'} fontSize="3xl" letterSpacing={4} lineHeight={2} color="white">
              Quicksilver
            </DrawerHeader>
            <DrawerBody>
              {['Staking', 'Governance', 'Defi', 'Assets'].map((item) => (
                <Box key={item} mb={4} position="relative">
                  <Link
                    href={`/quicksilver-app-v2/${item.toLowerCase()}`}
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

        {!isMobile && (
          <>
            <Spacer />
            <ScaleFade initialScale={0.5} in={!showSocialLinks}>
              {!showSocialLinks && (
                <VStack justifyContent="center" alignItems="center" spacing={16}>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Staking" placement="right">
                    <Box
                      w="55px"
                      h="55px"
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
                        h="55px"
                        src="/quicksilver-app-v2/img/test.png"
                      />
                    </Box>
                  </Tooltip>

                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Governance" placement="right">
                    <Box
                      w="55px"
                      h="55px"
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
                        w={'60px'}
                        src="/quicksilver-app-v2/img/test2.png"
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
                        src="/quicksilver-app-v2/img/test3.png"
                      />
                    </Box>
                  </Tooltip>

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
                        src="/quicksilver-app-v2/img/test4.png"
                      />
                    </Box>
                  </Tooltip>
                </VStack>
              )}
            </ScaleFade>

            <ScaleFade initialScale={0.5} in={showSocialLinks}>
              {showSocialLinks && (
                <VStack justifyContent="center" alignItems="center" spacing={16}>
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
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Discord" placement="right">
                    <Box
                      _hover={{
                        cursor: 'pointer',
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <FaDiscord size={'25px'} color="rgb(255, 128, 0)" />
                    </Box>
                  </Tooltip>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Github" placement="right">
                    <Box
                      _hover={{
                        cursor: 'pointer',
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <FaGithub size={'25px'} color="rgb(255, 128, 0)" />
                    </Box>
                  </Tooltip>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Twitter" placement="right">
                    <Box
                      _hover={{
                        cursor: 'pointer',
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <FaTwitter size={'25px'} color="rgb(255, 128, 0)" />
                    </Box>
                  </Tooltip>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Privacy Policy" placement="right">
                    <Box
                      _hover={{
                        cursor: 'pointer',
                        boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 15px 5px ${commonBoxShadowColor}`,
                        transition: transitionStyle,
                      }}
                    >
                      <MdPrivacyTip size={'25px'} color="rgb(255, 128, 0)" />
                    </Box>
                  </Tooltip>
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
          />
        )}
      </Flex>
    </Box>
  );
};

export default SideHeader;
