import { Button, Container, Flex, Icon, useColorMode, Box, Image, Spacer, useColorModeValue, VStack, IconButton, ButtonGroup, Tooltip } from "@chakra-ui/react";
import { BsFillMoonStarsFill, BsFillSunFill } from "react-icons/bs";
import { HamburgerIcon } from '@chakra-ui/icons'
import { WalletButton } from "../wallet-button";
import { useRouter } from 'next/router';
import Link from 'next/link'; 
import { useState, useEffect } from 'react';

export function SideHeader() {
    const router = useRouter(); 
    const [selectedPage, setSelectedPage] = useState("");
    useEffect(() => {
        // Function to handle route changes
        const handleRouteChange = (url: string) => {
            const path = url.split("/")[1]; // Get the path after the first '/'
            setSelectedPage(path);
        };

        // Add the route change listener
        router.events.on('routeChangeComplete', handleRouteChange);

        // Cleanup the listener when the component is unmounted
        return () => {
            router.events.off('routeChangeComplete', handleRouteChange);
        };
    }, [router]);

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
                    mt="-10px"
                    h="75px"
                    src="/img/networks/quicksilver.svg"
                    onClick={() => router.push('/')}
                    cursor="pointer"
                />
                <Spacer/>
                <VStack
                    justifyContent="center"
                    alignItems="center"
                    spacing={16}
                >
                    <Tooltip 
                    borderLeft= "4px solid rgba(255, 0, 0, 0.5)"
                    label="Staking" placement="right">
                    <Box
    w="55px"
    h="55px"
    onClick={() => router.push('/staking')}
    cursor="pointer"
    borderRadius="100px"
    boxShadow={selectedPage === 'staking' ? "0 0 15px 5px rgba(255, 0, 0, 0.25), inset 0 0 50px 5px rgba(255, 0, 0, 0.25)" : ""}
    _hover={{ 
        boxShadow:"0 0 15px 5px rgba(255, 0, 0, 0.25), inset 0 0 50px 5px rgba(255, 0, 0, 0.25)",

    }}
>
                            <Image 
                             filter={selectedPage === 'staking' ? "contrast(100%)" : "contrast(50%)"} 
                             _hover={{ filter: "contrast(100%)" }}
                            alt="Staking" h="55px" src="/img/test.png" />
                        </Box>
                    </Tooltip>

                    <Tooltip 
                    borderLeft= "4px solid rgba(128, 0, 128, 0.5)"
                    label="Governance" placement="right">
                    <Box
        w="55px"
        h="55px"
        onClick={() => router.push('/governance')}
        cursor="pointer"
        borderRadius="100px"
        boxShadow={selectedPage === 'governance' ? "0 0 15px 5px rgba(128, 0, 128, 0.25), inset 0 0 50px 5px rgba(128, 0, 128, 0.25)" : ""}
        _hover={{ 
            boxShadow:"0 0 15px 5px rgba(128, 0, 128, 0.25), inset 0 0 50px 5px rgba(128, 0, 128, 0.25)",
        }}
    >
                            <Image 
                            filter={selectedPage === 'governance' ? "contrast(100%)" : "contrast(50%)"} 
                            _hover={{ filter: "contrast(100%)" }}
                            alt="Governance" h="55px" src="/img/test2.png" />
                        </Box>
                    </Tooltip>

                    <Tooltip 
                    borderLeft= "4px solid rgba(0, 0, 255, 0.5)"
                    label="Assets" placement="right">
                    <Box
        w="55px"
        h="55px"
        onClick={() => router.push('/assets')}
        cursor="pointer"
        borderRadius="100px"
        boxShadow={selectedPage === 'assets' ? "0 0 15px 5px rgba(0, 0, 255, 0.25), inset 0 0 50px 5px rgba(0, 0, 255, 0.25)" : ""}
        _hover={{ 
            boxShadow:"0 0 15px 5px rgba(0, 0, 255, 0.25), inset 0 0 50px 5px rgba(0, 0, 255, 0.25)",
        }}
    >
                            <Image 
                            filter={selectedPage === 'assets' ? "contrast(100%)" : "contrast(50%)"} 
                            _hover={{ filter: "contrast(100%)" }}
                            alt="Assets" h="55px" src="/img/test3.png" />
                        </Box>
                    </Tooltip>

                    <Tooltip 
                    borderLeft= "4px solid rgba(255, 128, 0, 0.5)"
                    label="DeFi" placement="right">
                    <Box
        w="55px"
        h="55px"
        onClick={() => router.push('/defi')}
        cursor="pointer"
        borderRadius="100px"
        boxShadow={selectedPage === 'defi' ? "0 0 15px 5px rgba(255, 128, 0, 0.25), inset 0 0 50px 5px rgba(255, 128, 0, 0.25)" : ""}
        _hover={{ 
            boxShadow:"0 0 15px 5px rgba(255, 128, 0, 0.25), inset 0 0 50px 5px rgba(255, 128, 0, 0.25)",
        }}
    >
                            <Image 
                            filter={selectedPage === 'defi' ? "contrast(100%)" : "contrast(50%)"} 
                            _hover={{ filter: "contrast(100%)" }}
                            alt="DeFi" h="55px" src="/img/test4.png" />
                        </Box>
                    </Tooltip>
                </VStack>
                <Spacer/>
                <IconButton
                    icon={<HamburgerIcon />}
                    aria-label="DeFi"
                />
            </Flex>
        </Box>
    )
}
