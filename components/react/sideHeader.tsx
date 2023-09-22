import { HamburgerIcon } from '@chakra-ui/icons';
import {
  Flex,
  Box,
  Image,
  Spacer,
  VStack,
  IconButton,
  Tooltip,
} from '@chakra-ui/react';
import { useRouter } from 'next/router';
import { useState, useEffect } from 'react';

export const SideHeader = () => {
  const router = useRouter();
  const [selectedPage, setSelectedPage] = useState('');
  useEffect(() => {
    // Function to handle route changes
    const handleRouteChange = (url: string) => {
      const path = url.split('/quicksilver-app-v2/')[1];
      setSelectedPage(path);
    };

    // Add the route change listener
    router.events.on('routeChangeComplete', handleRouteChange);

    // Cleanup the listener when the component is unmounted
    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [router]);

  const commonBoxShadowColor = 'rgba(255, 128, 0, 0.25)';

  return (
    <Box
      w="fit-content"
      h="95vh"
      backdropFilter="blur(10px)"
      borderRadius={10}
      zIndex={10}
      top={6}
      left="6"
      position="fixed"
      bgColor="rgba(214, 219, 220, 0.1)"
    >
      <Flex
        direction="column"
        align="center"
        zIndex={10}
        justifyContent="space-between"
        py={4}
        height="100%"
      >
        <Image
          alt="logo"
          mt="-10px"
          h="75px"
          src="/quicksilver-app-v2/img/networks/quicksilver.svg"
          onClick={() => router.push('/')}
          cursor="pointer"
        />
        <Spacer />
        <VStack justifyContent="center" alignItems="center" spacing={16}>
          <Tooltip
            borderLeft="4px solid rgba(255, 128, 0, 0.9)"
            label="Staking"
            placement="right"
          >
            <Box
              w="55px"
              h="55px"
              onClick={() => router.push('/staking')}
              cursor="pointer"
              borderRadius="100px"
              boxShadow={
                selectedPage === 'staking'
                  ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`
                  : ''
              }
              _hover={{
                boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
              }}
            >
              <Image
                filter={
                  selectedPage === 'staking'
                    ? 'contrast(100%)'
                    : 'contrast(50%)'
                }
                _hover={{
                  filter: 'contrast(100%)',
                }}
                alt="Staking"
                h="55px"
                src="/quicksilver-app-v2/img/test.png"
              />
            </Box>
          </Tooltip>

          <Tooltip
            borderLeft="4px solid rgba(255, 128, 0, 0.9)"
            label="Governance"
            placement="right"
          >
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
              }}
            >
              <Image
                filter={
                  selectedPage === 'governance'
                    ? 'contrast(100%)'
                    : 'contrast(50%)'
                }
                _hover={{
                  filter: 'contrast(100%)',
                }}
                alt="Governance"
                h="55px"
                src="/quicksilver-app-v2/img/test2.png"
              />
            </Box>
          </Tooltip>

          <Tooltip
            borderLeft="4px solid rgba(255, 128, 0, 0.9)"
            label="Assets"
            placement="right"
          >
            <Box
              w="55px"
              h="55px"
              onClick={() => router.push('/assets')}
              cursor="pointer"
              borderRadius="100px"
              boxShadow={
                selectedPage === 'assets'
                  ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`
                  : ''
              }
              _hover={{
                boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
              }}
            >
              <Image
                filter={
                  selectedPage === 'assets' ? 'contrast(100%)' : 'contrast(50%)'
                }
                _hover={{
                  filter: 'contrast(100%)',
                }}
                alt="Assets"
                h="55px"
                src="/quicksilver-app-v2/img/test3.png"
              />
            </Box>
          </Tooltip>

          <Tooltip
            borderLeft="4px solid rgba(255, 128, 0, 0.9)"
            label="DeFi"
            placement="right"
          >
            <Box
              w="55px"
              h="55px"
              onClick={() => router.push('/defi')}
              cursor="pointer"
              borderRadius="100px"
              boxShadow={
                selectedPage === 'defi'
                  ? `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`
                  : ''
              }
              _hover={{
                boxShadow: `0 0 15px 5px ${commonBoxShadowColor}, inset 0 0 50px 5px ${commonBoxShadowColor}`,
              }}
            >
              <Image
                filter={
                  selectedPage === 'defi' ? 'contrast(100%)' : 'contrast(50%)'
                }
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
        <Spacer />
        <IconButton icon={<HamburgerIcon />} aria-label="DeFi" />
      </Flex>
    </Box>
  );
};
