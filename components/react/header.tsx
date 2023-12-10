import { Flex, Box, Image } from '@chakra-ui/react';

import WalletTest from './wallet-test';
import { WalletButton } from '../wallet-button';

export const Header: React.FC<{ chainName: string }> = ({ chainName }) => {
  return (
    <Box w="100%" borderRadius={0} maxH="125px" zIndex={10} mt={'5px'} position="relative" px={10} bgColor="transparent">
      <Flex maxW="100%" mx="auto" align="center" zIndex={10} position="sticky" top="0" justifyContent="space-between" py={1}>
        <Image alt="" h="85px" />
        <Flex alignItems="center" justifyContent="center">
          <WalletButton chainName={chainName} />
        </Flex>
      </Flex>
    </Box>
  );
};
