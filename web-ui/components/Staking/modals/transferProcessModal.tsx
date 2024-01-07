import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalBody,
  ModalCloseButton,
  HStack,
  Text,
  Box,
  Circle,
  Flex,
  Button,
  Stat,
  StatLabel,
  StatNumber,
  Toast,
  Spinner,
  useToast,
  Input,
  Grid,
  Checkbox,
} from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import styled from '@emotion/styled';
import React, { useEffect, useState } from 'react';

import { MultiModal } from './validatorSelectionModal';

import { useQueryHooks, useTx } from '@/hooks';
import { useZoneQuery } from '@/hooks/useQueries';
import { liquidStakeTx, unbondLiquidStakeTx } from '@/tx/liquidStakeTx';
import { bech32 } from 'bech32';
import { shiftDigits } from '@/utils';
import { coins, StdFee } from '@cosmjs/amino';
import { assets } from 'chain-registry';
import { cosmos } from 'interchain-query';
import chains from '@chalabi/chain-registry';
import { TxResponse } from 'interchain-query/cosmos/base/abci/v1beta1/abci';

const ChakraModalContent = styled(ModalContent)`
  position: relative;
  background: none;
  &::before,
  &::after {
    z-index: -1;
  }
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    width: 40%;
    background-color: #201c18;
    border-radius: 5px 0 0 5px;
  }
  &::after {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    width: 60%;
    background-color: #1a1a1a;
    border-radius: 0 5px 5px 0;
  }
`;

interface SelectedValidator {
  operatorAddress: string;
  moniker: string;
  tokenAmount: string;
}

interface StakingModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
  selectedValidator: SelectedValidator;
  selectedOption?: {
    name: string;
    value: string;
    logo: string;
    chainName: string;
    chainId: string;
  };
}

export const TransferProcessModal: React.FC<StakingModalProps> = ({ isOpen, onClose, selectedOption, selectedValidator }) => {
  const [step, setStep] = React.useState(1);
  const getProgressColor = (circleStep: number) => {
    if (step >= circleStep) return 'complimentary.900';
    return 'rgba(255,255,255,0.2)';
  };

  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  let newChainName: string | undefined;
  if (selectedOption?.chainId === 'provider') {
    newChainName = 'rsprovidertestnet';
  } else if (selectedOption?.chainId === 'elgafar-1') {
    newChainName = 'stargazetestnet';
  } else if (selectedOption?.chainId === 'osmo-test-5') {
    newChainName = 'osmosistestnet';
  } else if (selectedOption?.chainId === 'regen-redwood-1') {
    newChainName = 'regen';
  } else {
    newChainName = selectedOption?.chainName;
  }

  const { data: zone, isLoading: isZoneLoading, isError: isZoneError } = useZoneQuery(selectedOption?.chainId ?? '');
  const labels = ['Tokenize Shares', `Transfer`, `Receive q${selectedOption?.value}`];
  const [transactionStatus, setTransactionStatus] = useState('Pending');
  function truncateString(str: string, num: number) {
    if (str.length > num) {
      return str.slice(0, num) + '...';
    } else {
      return str;
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} size={{ base: '3xl', md: '2xl' }}>
      <ModalOverlay />
      <ChakraModalContent h="48%" maxH={'100%'}>
        <ModalBody borderRadius={4} h="48%" maxH={'100%'}>
          <ModalCloseButton zIndex={1000} color="white" />
          <HStack position={'relative'} h="100%" spacing="48px" align="stretch">
            {/* Left Section */}
            <Flex flexDirection="column" justifyContent="space-between" width="40%" p={4} bg="#1E1C19" height="100%">
              <Box position="relative">
                <Stat>
                  <StatLabel color="rgba(255,255,255,0.5)">TRANSFER DELEGATION</StatLabel>
                  <StatNumber color="white">{truncateString(selectedValidator.moniker, 13)}</StatNumber>
                  <StatNumber color="white">
                    {shiftDigits(selectedValidator.tokenAmount, -6)}&nbsp;
                    {selectedOption?.value}
                  </StatNumber>
                </Stat>
                {[1, 2, 3].map((circleStep, index) => (
                  <Flex key={circleStep} align="center" mt={10} mb={circleStep !== 4 ? '48px' : '0'}>
                    <Circle
                      size="36px"
                      bg={getProgressColor(circleStep)}
                      color="white"
                      fontWeight="bold"
                      borderWidth={'2px'}
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                      position="relative"
                      borderColor="rgba(255,255,255,0.5)"
                    >
                      {circleStep}
                      {circleStep !== 3 && (
                        <>
                          <Box
                            width="2px"
                            height="30px"
                            bgColor="rgba(255,255,255,0.01)"
                            position="absolute"
                            bottom="-42px"
                            left="50%"
                            transform="translateX(-50%)"
                          />
                          <Box
                            width="2px"
                            height="30px"
                            bgColor={getProgressColor(circleStep + 1)}
                            position="absolute"
                            bottom="-42px"
                            left="50%"
                            transform="translateX(-50%)"
                          />
                        </>
                      )}
                    </Circle>
                    <Text fontWeight="hairline" ml={3} color="rgba(255,255,255,0.75)">
                      {labels[index]}
                    </Text>
                  </Flex>
                ))}
              </Box>
            </Flex>

            <Flex width="67%" flexDirection="column" justifyContent="center" alignItems="center">
              {step === 1 && (
                <>
                  <Flex maxW="300px" flexDirection={'column'} justifyContent={'left'} alignItems={'center'}>
                    <Text textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                      Tokenize Shares
                    </Text>
                    <Text mt={2} textAlign={'center'} fontWeight={'light'} fontSize="lg" color="white">
                      Tokenize your delegation in order to transfer it to Quicksilver
                    </Text>
                  </Flex>

                  <Button
                    mt={4}
                    width="55%"
                    _hover={{
                      bgColor: 'complimentary.500',
                    }}
                  ></Button>
                </>
              )}
              {step === 2 && (
                <>
                  <Text textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                    Send to Quicksilver
                  </Text>
                  <Text mt={2} textAlign={'center'} fontWeight={'light'} fontSize="lg" color="white">
                    Specifying weights allows you to choose how much of your liquid delegation goes to each validator.
                  </Text>

                  <Button
                    _hover={{
                      bgColor: 'complimentary.500',
                    }}
                  >
                    Equal
                  </Button>

                  <Button
                    position={'absolute'}
                    bottom={3}
                    left={'51%'}
                    bgColor="none"
                    _hover={{
                      bgColor: 'none',
                      color: 'complimentary.900',
                    }}
                    _selected={{
                      bgColor: 'none',
                    }}
                    color="white"
                    variant="none"
                  >
                    ←
                  </Button>
                </>
              )}

              {step === 3 && (
                <>
                  <Box justifyContent={'center'}>
                    <Flex maxW="300px" flexDirection={'column'} justifyContent={'left'} alignItems={'center'}>
                      <Text textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                        Transaction {transactionStatus}
                      </Text>
                      <Text mt={2} textAlign={'center'} fontWeight={'light'} fontSize="lg" color="white">
                        Your q{selectedOption?.value} will arrive to your wallet in a few minutes.
                      </Text>
                      <Button
                        w="55%"
                        _hover={{
                          bgColor: '#181818',
                        }}
                        mt={4}
                        onClick={() => setStep(1)}
                      >
                        Stake Again
                      </Button>
                    </Flex>
                  </Box>
                </>
              )}
            </Flex>
          </HStack>
        </ModalBody>
      </ChakraModalContent>
    </Modal>
  );
};
export default TransferProcessModal;
