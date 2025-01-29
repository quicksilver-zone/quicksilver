import { HamburgerIcon, ArrowBackIcon } from '@chakra-ui/icons';
import { Flex, Box, Image, Spacer, VStack, IconButton, Tooltip, ScaleFade, useDisclosure, Link } from '@chakra-ui/react';
import { useRouter } from 'next/router';
import { useState, useEffect } from 'react';
import { FaDiscord, FaGithub, FaInfo } from 'react-icons/fa';
import { FaXTwitter } from 'react-icons/fa6';
import { IoIosDocument } from 'react-icons/io';
import { MdPrivacyTip } from 'react-icons/md';

import { AccountControlModal } from './accountControlModal';

export const SideHeader = () => {
  const router = useRouter();
  const [selectedPage, setSelectedPage] = useState('');
  const { isOpen, onOpen, onClose } = useDisclosure();

  const [showSocialLinks, setShowSocialLinks] = useState(false);

  useEffect(() => {
    const path = router.asPath.split('/')[1];
    setSelectedPage(path);

    const handleRouteChange = (url: string) => {
      const newPath = url.split('/')[1];
      setSelectedPage(newPath);
    };

    router.events.on('routeChangeComplete', handleRouteChange);

    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [router]);

  const commonBoxShadowColor = 'rgba(255, 128, 0, 0.25)';
  const toggleSocialLinks = () => setShowSocialLinks(!showSocialLinks);

  const transitionStyle = 'all 0.3s ease';

  return (
    <Box
      w={'fit-content'}
      h={'95vh'}
      backdropFilter="blur(10px)"
      borderRadius={10}
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
          onClick={() => router.push('/')}
          cursor="pointer"
        />

        <>
          <Spacer />
          <ScaleFade initialScale={0.5} in={!showSocialLinks}>
            {!showSocialLinks && (
              <VStack justifyContent="center" alignItems="center" spacing={16}>
                <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Staking" placement="right">
                  <Box
                    w="50px"
                    h="50px"
                    onClick={() => router.push('/staking')}
                    cursor="pointer"
                    borderRadius="100px"
                    _hover={{
                      transition: transitionStyle,
                    }}
                  >
                    <Image
                      filter={selectedPage === 'staking' ? 'contrast(100%)' : 'contrast(0%)'}
                      _hover={{
                        filter: 'contrast(100%)',
                      }}
                      alt="Staking"
                      h="50px"
                      w="50px"
                      src="/img/liquid.png"
                    />
                  </Box>
                </Tooltip>

                <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Governance" placement="right">
                  <Box
                    w="50px"
                    h="50px"
                    onClick={() => router.push('/governance')}
                    cursor="pointer"
                    borderRadius="100px"
                    _hover={{
                      transition: transitionStyle,
                    }}
                  >
                    <Image
                      filter={selectedPage === 'governance' ? 'contrast(100%)' : 'contrast(0%)'}
                      _hover={{
                        filter: 'contrast(100%)',
                      }}
                      alt="Governance"
                      h="50px"
                      w="50px"
                      src="/img/governance.png"
                    />
                  </Box>
                </Tooltip>

                <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Assets" placement="right">
                  <Box
                    w="50px"
                    h="50px"
                    onClick={() => router.push('/assets')}
                    cursor="pointer"
                    borderRadius="100px"
                    _hover={{
                      transition: transitionStyle,
                    }}
                  >
                    <Image
                      filter={selectedPage === 'assets' ? 'contrast(100%)' : 'contrast(0%)'}
                      _hover={{
                        filter: 'contrast(100%)',
                      }}
                      alt="Assets"
                      h="50px"
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
                    w="50px"
                    h="50px"
                    onClick={() => router.push('/defi')}
                    cursor="pointer"
                    borderRadius="100px"
                    _hover={{
                      transition: transitionStyle,
                    }}
                  >
                    <Image
                      filter={selectedPage === 'defi' ? 'contrast(100%)' : 'contrast(0%)'}
                      _hover={{
                        filter: 'contrast(100%)',
                      }}
                      alt="DeFi"
                      h="50px"
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
                      borderRadius={'full'}
                      _hover={{
                        cursor: 'pointer',
                      }}
                    >
                      <FaInfo size={'25px'} color="rgb(255, 128, 0)" />
                    </Box>
                  </Tooltip>
                </Link>

                <Link href="https://docs.quicksilver.zone/" isExternal>
                  <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Docs" placement="right">
                    <Box
                      borderRadius={'full'}
                      _hover={{
                        cursor: 'pointer',
                        transition: transitionStyle,
                      }}
                    >
                      <IoIosDocument size={'25px'} color="rgba(255, 128, 0, 0.9)" />
                    </Box>
                  </Tooltip>
                </Link>
                <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Account Controls" placement="right">
                  <Box
                    borderRadius={'full'}
                    onClick={onOpen}
                    _hover={{
                      cursor: 'pointer',
                      transition: transitionStyle,
                    }}
                  >
                    <MdPrivacyTip size={'25px'} color="rgb(255, 128, 0)" />
                  </Box>
                </Tooltip>
                <AccountControlModal isOpen={isOpen} onClose={onClose} />
                <Tooltip borderLeft="4px solid rgba(255, 128, 0, 0.9)" label="Discord" placement="right">
                  <Link href="https://discord.com/invite/xrSmYMDVrQ" isExternal>
                    <Box
                      borderRadius={'full'}
                      _hover={{
                        cursor: 'pointer',
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
                      borderRadius={'full'}
                      _hover={{
                        cursor: 'pointer',
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
                      borderRadius={'full'}
                      _hover={{
                        cursor: 'pointer',
                        transition: transitionStyle,
                      }}
                    >
                      <FaXTwitter size={'25px'} color="rgb(255, 128, 0)" />
                    </Box>
                  </Link>
                </Tooltip>
              </VStack>
            )}
          </ScaleFade>
        </>

        <Spacer />

        <IconButton
          borderRadius={10}
          icon={showSocialLinks ? <ArrowBackIcon /> : <HamburgerIcon />}
          aria-label="Toggle View"
          onClick={toggleSocialLinks}
          mb={4}
          _hover={{
            bgColor: 'complimentary.500',
          }}
        />
      </Flex>
    </Box>
  );
};

export default SideHeader;
