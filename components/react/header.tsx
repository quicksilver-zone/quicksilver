import { Button, Container, Flex, Icon, useColorMode, Box, Image, Spacer, useColorModeValue } from "@chakra-ui/react";
import { BsFillMoonStarsFill, BsFillSunFill } from "react-icons/bs";
import { WalletButton } from "../wallet-button";

export function Header() {
    const { colorMode, toggleColorMode } = useColorMode();
    const buttonTextColor = useColorModeValue("primary.700", "primary.50")
    
    return (
        <Box
            w="100%" 
            borderRadius={0}
            maxH="125px"
            zIndex={10}
            top="0"
            position="sticky"
            px={10}
            bgColor="transparent"

        >
            <Flex
                maxW="100%"
                mx="auto"
                align="center"
                zIndex={10}
                position="sticky"
                top="0"
                justifyContent="space-between"
                py={1}
            >
                <Image
                h="85px"
               
                />
                <Flex
                alignItems="center"
                justifyContent="center"
                >
                <WalletButton/>
               
                </Flex>
            </Flex>
        </Box>
    )
}
