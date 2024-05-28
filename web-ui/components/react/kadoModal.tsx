import { Modal, ModalOverlay, ModalContent, ModalBody, ModalCloseButton, Flex, Box, Fade, useBreakpointValue } from '@chakra-ui/react';
import { useChains } from '@cosmos-kit/react';
import { useState, useEffect } from 'react';

import KadoIconContent from './kadoIcon';

interface KadoModalProps {
  isOpen: boolean;
  onClose: () => void;
  denom?: string;
  zone?: string;
}

export const KadoModal: React.FC<KadoModalProps> = ({ isOpen, onClose, denom, zone }) => {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!isOpen) {
      setIsLoading(true);
    }
  }, [isOpen]);

  const isMobile = useBreakpointValue({ base: true, sm: true, md: false, lg: false, xl: false });

  const onRevCurrency = denom ? denom : 'QCK';

  const networks = ['OSMOSIS', 'JUNO', 'COSMOS HUB', 'STARGAZE', 'REGEN', 'AGORIC', 'QUICKSILVER'];

  const network = zone ? zone : 'QUICKSILVER';

  const chains = useChains(['osmosis', 'juno', 'cosmoshub', 'stargaze', 'regen', 'agoric', 'quicksilver']);

  const addresses = [
    `QUICKSILVER:${chains.quicksilver.address}`,
    `OSMOSIS:${chains.osmosis.address}`,
    `COSMOS HUB:${chains.cosmoshub.address}`,
    `STARGAZE:${chains.stargaze.address}`,
    `AGORIC:${chains.agoric.address}`,
    `REGEN:${chains.regen.address}`,
    `JUNO:${chains.juno.address}`,
  ];
  const offFromAddress = addresses.find((address) => address.startsWith(network));

  return (
    <Modal closeOnOverlayClick={false} isOpen={isOpen} onClose={onClose} size={{ base: 'sm', sm: 'sm', md: 'xl' }}>
      <ModalOverlay />
      <ModalContent backgroundColor={'#0b121f'} maxH={'100%'}>
        <ModalBody borderRadius={4} maxH={'100%'}>
          <ModalCloseButton zIndex={1000} color="white" />

          <Flex p={4} justifyContent={'center'} alignItems={'center'} position="relative">
            {isLoading && (
              <Box
                width={isMobile ? '380px' : '480px'}
                height={isMobile ? '620px' : '620px'}
                display="flex"
                justifyContent="center"
                alignItems="center"
                borderRadius="20px"
                position="relative"
                top="0"
                left="0"
              >
                <KadoIconContent width={'8em'} height={'8em'} showAnimation />
              </Box>
            )}
            <Fade in={!isLoading}>
              <iframe
                src={`https://app.kado.money/?apiKey=5fef3eb4-2c88-4645-9f92-519e9b5a9fcc&primaryColor=%23FF8000&secondaryColor=%181515&theme=dark&network=${network}&onRevCurrency=${onRevCurrency}&networkList=${networks}&onToAddressMulti=${addresses}&offFromAddress=${offFromAddress}`}
                width={isMobile ? '380px' : '480px'}
                height={isMobile ? '620px' : '620px'}
                style={{ display: isLoading ? 'none' : 'block' }}
                allow="clipboard-write; payment; accelerometer; gyroscope; camera; geolocation; autoplay; fullscreen;"
                onLoad={() => setIsLoading(false)}
              ></iframe>
            </Fade>
          </Flex>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

export default KadoModal;
