import { Button, Container, Flex, Icon, useColorMode, Box, Image, Spacer, useColorModeValue, VStack, IconButton, ButtonGroup, Tooltip } from "@chakra-ui/react";
import { BsFillMoonStarsFill, BsFillSunFill } from "react-icons/bs";
import { HamburgerIcon } from '@chakra-ui/icons'
import { WalletButton } from "../wallet-button";
import { useRouter } from 'next/router';
import Link from 'next/link'; // Import the Link component from next/link

export function SideHeader() {
    const router = useRouter(); // Initialize the useRouter hook
    
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
            px={2}
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
                    h="75px"
                    src="/img/networks/quicksilver.svg"
                    onClick={() => router.push('/')}
                    cursor="pointer"
                />
                <Spacer/>
                <VStack
                    justifyContent="center"
                    alignItems="center"
                    spacing={8}
                >
                    <Tooltip label="Staking" placement="right">
                    <Box
                            w="55px"
                            h="55px"
                            onClick={() => router.push('/staking')}
                            cursor="pointer"
                            borderRadius="100px"
                            _hover={{ 
                                boxShadow: "0 0 15px 5px rgba(255, 0, 0, 1)",
                                backdropFilter: "blur(70px)",
                                bgColor: "rgba(255, 0, 0, 0.5)",
                            }}
                        >
                            <Image alt="Staking" h="55px" src="/img/networks/quicksilver.svg" />
                        </Box>
                    </Tooltip>

                    <Tooltip label="Governance" placement="right">
                    <Box
                            w="55px"
                            h="55px"
                            onClick={() => router.push('/governance')}
                            cursor="pointer"
                            borderRadius="100px"
                            _hover={{ 
                                boxShadow: "0 0 15px 5px rgba(128, 0, 128, 1)",
                                backdropFilter: "blur(70px)",
                                bgColor: "rgba(128, 0, 128, 0.5)",
                            }}
                        >
                            <Image alt="Governance" h="55px" src="/img/networks/quicksilver.svg" />
                        </Box>
                    </Tooltip>

                    <Tooltip label="Assets" placement="right">
                        <Box
                            w="55px"
                            h="55px"
                            onClick={() => router.push('/assets')}
                            cursor="pointer"
                            borderRadius="100px"
                            _hover={{ 
                                boxShadow: "0 0 15px 5px rgba(0, 0, 255, 1)",
                                backdropFilter: "blur(70px)",
                                bgColor: "rgba(0, 0, 255, 0.5)",
                            }}
                        >
                            <Image alt="Assets" h="55px" src="/img/networks/quicksilver.svg" />
                        </Box>
                    </Tooltip>

                    <Tooltip label="DeFi" placement="right">
                        <Box
                            w="55px"
                            h="55px"
                            onClick={() => router.push('/defi')}
                            cursor="pointer"
                            borderRadius="100px"
                            _hover={{ 
                                boxShadow: "0 0 15px 5px rgba(255, 128, 0, 1)",
                                backdropFilter: "blur(70px)",
                                bgColor: "rgba(255, 128, 0, 0.5)",
                            }}
                        >
                            <Image alt="DeFi" h="55px" src="/img/networks/quicksilver.svg" />
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
