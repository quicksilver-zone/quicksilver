import { HamburgerIcon, ArrowBackIcon } from '@chakra-ui/icons';
import { Flex, Box, Image, Spacer, VStack, IconButton, Tooltip, ScaleFade, useBreakpointValue } from '@chakra-ui/react';
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

  // Use breakpoint value to determine if the device is mobile
  const isMobile = useBreakpointValue({ base: true, md: false });
  const transitionStyle = 'all 0.3s ease';

  return (
    <Box
      w={isMobile ? 'auto' : 'fit-content'}
      h={{ base: 'auto', md: '95vh' }}
      backdropFilter="blur(10px)"
      borderRadius={100}
      zIndex={10}
      top={6}
      left={6}
      position="fixed"
      bgColor="rgba(214, 219, 220, 0.1)"
    >
      <Flex direction="column" align="center" zIndex={10} justifyContent="space-between" py={4} height="100%">
        <Image
          alt="logo"
          mt="-10px"
          h="75px"
          w={'75px'}
          src="/quicksilver-app-v2/img/networks/quicksilver.svg"
          onClick={() => router.push('/')}
          cursor="pointer"
        />

        {/* Only display additional content if not on mobile */}
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
        {/* Only display the IconButton if not on mobile */}
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
