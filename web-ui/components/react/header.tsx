import { Flex, Box, Image } from '@chakra-ui/react';
import dynamic from 'next/dynamic';

import { WalletButton } from '../wallet-button';

const Header: React.FC<{ chainName: string }> = ({ chainName }) => {
  return (
    <Box w="100%" borderRadius={0} maxH="125px" zIndex={0} top={0} position="fixed" px={10} bgColor="transparent">
      <Flex maxW="100%" mx="auto" align="center" zIndex={0} position="sticky" top="0" justifyContent="space-between" py={1}>
        <Image alt="" h="85px" />
        <Flex alignItems="center" justifyContent="center">
          <WalletButton />
        </Flex>
      </Flex>
    </Box>
  );
};

export const DynamicHeaderSection = dynamic(() => Promise.resolve(Header), {
  ssr: false,
});
