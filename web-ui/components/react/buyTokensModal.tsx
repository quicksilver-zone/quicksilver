import { Modal, ModalOverlay, ModalContent, ModalBody, ModalCloseButton, VStack, HStack, Box, Heading, Text, Link } from '@chakra-ui/react';
import Image from 'next/image';
import { useState } from 'react';

import KadoIconContent from './kadoIcon';
import KadoModal from './kadoModal';

interface KadoModalProps {
  isOpen: boolean;
  onClose: () => void;
  denom: string;
  zone?: string;
}

export const BuyTokensModal: React.FC<KadoModalProps> = ({ isOpen, onClose, denom, zone }) => {
  const osmosisLink = `https://app.osmosis.zone/?from=USDC&to=${denom.toUpperCase()}`;
  const [isKadoModalOpen, setIsKadoModalOpen] = useState(false);
  const kadoTokens = ['OSMO', 'ATOM', 'STARS', 'BLD', 'REGEN', 'QCK'];

  return (
    <>
      <Modal isOpen={isOpen} onClose={onClose} size={{ base: 'sm', sm: 'sm', md: 'xl' }}>
        <ModalOverlay />
        <ModalContent backgroundColor={'#201c18'} maxH={'100%'}>
          <ModalBody p={8} borderRadius={4} maxH={'100%'}>
            <ModalCloseButton zIndex={1000} color="white" />
            <Text fontSize={'x-large'} fontWeight={'bold'} mb={4}>
              Choose an option
            </Text>
            <VStack spacing={6} align="stretch">
              {kadoTokens.includes(denom.toUpperCase()) && (
                <Box
                  onClick={() => setIsKadoModalOpen(true)}
                  cursor={'pointer'}
                  _hover={{ background: 'rgba(0,0,0,0.3)' }}
                  p={4}
                  bg="rgba(0,0,0,0.5)"
                  borderRadius="md"
                  boxShadow="lg"
                >
                  <HStack spacing={4}>
                    <KadoIconContent stock width={'3em'} height={'3em'} />
                    <VStack align="start" spacing={0}>
                      <Heading size="md" color="white">
                        Kado
                      </Heading>
                      <Text fontSize="sm" color="gray.400">
                        Fiat Onramp
                      </Text>
                    </VStack>
                  </HStack>
                </Box>
              )}

              <Link href={osmosisLink} isExternal>
                <Box
                  cursor={'pointer'}
                  _hover={{ background: 'rgba(0,0,0,0.3)' }}
                  p={4}
                  bg="rgba(0,0,0,0.5)"
                  borderRadius="md"
                  boxShadow="lg"
                >
                  <HStack spacing={4}>
                    <Image alt="Osmosis Icon" src="/img/osmoIcon.svg" width={'48'} height={'48'} />
                    <VStack align="start" spacing={0}>
                      <Heading size="md" color="white">
                        Osmosis
                      </Heading>
                      <Text fontSize="sm" color="gray.400">
                        Trade on Osmosis
                      </Text>
                    </VStack>
                  </HStack>
                </Box>
              </Link>
            </VStack>
          </ModalBody>
        </ModalContent>
      </Modal>
      <KadoModal isOpen={isKadoModalOpen} onClose={() => setIsKadoModalOpen(false)} denom={denom.toUpperCase()} zone={zone} />
    </>
  );
};

export default BuyTokensModal;
