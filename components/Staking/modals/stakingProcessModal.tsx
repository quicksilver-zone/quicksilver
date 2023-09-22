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
  useToast,
} from '@chakra-ui/react';
import styled from '@emotion/styled';
import { getSigningQuicksilverClient } from 'quicksilverjs';
import React, { useState } from 'react';

import { liquidStakeTx } from '@/tx/liquidSteakTx';

import { MultiModal } from './validatorSelectionModal';

const ChakraModalContent = styled(ModalContent)`
  position: relative;
  background: none;
  &::before,
  &::after {
    z-index: -1; // Push the pseudo-elements to the background
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

interface StakingModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
  tokenAmount: string;
  selectedOption?: {
    name: string;
    value: string;
    logo: string;
    chainName: string;
  };
}

export const StakingProcessModal: React.FC<StakingModalProps> = ({
  isOpen,
  onClose,
  selectedOption,
  tokenAmount,
}) => {
  const [step, setStep] = React.useState(1);
  const getProgressColor = (circleStep: number) => {
    if (step >= circleStep) return 'complimentary.900';
    return 'rgba(255,255,255,0.2)';
  };

  const labels = [
    'Choose validators',
    `Summary`,
    `Receive q${selectedOption?.value}`,
  ];
  const [isModalOpen, setModalOpen] = useState(false);

  const [selectedValidators, setSelectedValidators] = React.useState<string[]>(
    [],
  );

  const [resp, setResp] = useState('');

  const advanceStep = () => {
    if (selectedValidators.length > 0) {
      setStep((prevStep) => prevStep + 1);
    }
  };

  const retreatStep = () => {
    setStep((prevStep) => Math.max(prevStep - 1, 1));
  };

  const toast = useToast();

  const intents = selectedValidators.map((validator) => ({
    validatorAddress: validator,
    weights: '1',
  }));

  const liquidStake = liquidStakeTx(
    setResp,
    selectedOption?.chainName || '',
    selectedOption?.chainName || '',
    intents,
    toast,
  );

  //placehoder for transaction status
  const [transactionStatus, setTransactionStatus] = useState('Pending');
  return (
    <Modal isOpen={isOpen} onClose={onClose} size="2xl">
      <ModalOverlay />
      <ChakraModalContent h="38%">
        <ModalBody borderRadius={4} h="38%">
          <ModalCloseButton color="white" />
          <HStack h="100%" spacing="48px" align="stretch">
            {/* Left Section */}
            <Flex
              flexDirection="column"
              justifyContent="space-between"
              width="40%"
              p={4}
              bg="#1E1C19"
              height="100%"
            >
              <Box position="relative">
                <Stat>
                  <StatLabel color="rgba(255,255,255,0.5)">
                    LIQUID STAKING
                  </StatLabel>
                  <StatNumber color="white">
                    {tokenAmount} {selectedOption?.value}
                  </StatNumber>
                </Stat>
                {[1, 2, 3].map((circleStep, index) => (
                  <Flex
                    key={circleStep}
                    align="center"
                    mt={10}
                    mb={circleStep !== 3 ? '48px' : '0'}
                  >
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
                    <Text
                      fontWeight="hairline"
                      ml={3}
                      color="rgba(255,255,255,0.75)"
                    >
                      {labels[index]}
                    </Text>
                  </Flex>
                ))}
              </Box>
            </Flex>

            {/* Right Section */}
            <Flex
              width="67%"
              flexDirection="column"
              justifyContent="center"
              alignItems="center"
            >
              {step === 1 && (
                <>
                  <Flex
                    maxW="300px"
                    flexDirection={'column'}
                    justifyContent={'left'}
                    alignItems={'center'}
                  >
                    <Text
                      textAlign={'left'}
                      fontWeight={'bold'}
                      fontSize="lg"
                      color="white"
                    >
                      Choose Validators
                    </Text>
                    <Text
                      mt={2}
                      textAlign={'center'}
                      fontWeight={'light'}
                      fontSize="lg"
                      color="white"
                    >
                      Select up to 8 validators to split your liquid delegation
                      between.
                    </Text>
                  </Flex>
                  {selectedValidators.length > 0 && (
                    <Button
                      mt={2}
                      color="white"
                      _hover={{
                        bgColor: 'rgba(255, 128, 0, 0.25)',
                      }}
                      variant="ghost"
                      width="35%"
                      size="xs"
                      onClick={() => setModalOpen(true)}
                    >
                      Reselect Validators
                    </Button>
                  )}
                  <Button
                    mt={4}
                    width="55%"
                    _hover={{
                      bgColor: '#181818',
                    }}
                    onClick={() => {
                      if (selectedValidators.length === 0) {
                        setModalOpen(true);
                      } else {
                        advanceStep();
                      }
                    }}
                  >
                    {selectedValidators.length > 0
                      ? 'Next'
                      : 'Choose Validators'}
                  </Button>
                  <MultiModal
                    isOpen={isModalOpen}
                    onClose={() => setModalOpen(false)}
                    selectedChainName={selectedOption?.chainName || ''}
                    selectedValidators={selectedValidators}
                    setSelectedValidators={setSelectedValidators}
                  />
                </>
              )}
              {step === 2 && (
                <>
                  <Box justifyContent={'center'}>
                    <Text
                      fontWeight={'bold'}
                      fontSize="lg"
                      w="250px"
                      textAlign={'left'}
                      color="white"
                    >
                      You’re going to liquid stake {tokenAmount}{' '}
                      {selectedOption?.value} on Quicksilver
                    </Text>
                    <HStack
                      mt={2}
                      textAlign={'left'}
                      fontWeight={'light'}
                      fontSize="lg"
                      color="white"
                    >
                      <Text fontWeight={'bold'}>Receiving:</Text>
                      <Text color="complimentary.900">
                        {(Number(tokenAmount) * 0.95).toFixed(2)} q
                        {selectedOption?.value}
                      </Text>
                    </HStack>
                    <Text mt={2} textAlign={'left'} fontWeight={'hairline'}>
                      Processing time: 1 minute
                    </Text>
                    <Button
                      w="55%"
                      _hover={{
                        bgColor: '#181818',
                      }}
                      mt={4}
                      onClick={() =>
                        setStep((prevStep) => (prevStep < 3 ? prevStep + 1 : 3))
                      }
                    >
                      Submit
                    </Button>
                  </Box>
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
                    onClick={retreatStep}
                  >
                    ←
                  </Button>
                </>
              )}

              {step === 3 && (
                <>
                  <Box justifyContent={'center'}>
                    <Text
                      fontWeight={'bold'}
                      fontSize="lg"
                      w="250px"
                      textAlign={'left'}
                      color="white"
                    >
                      Status: {transactionStatus}
                    </Text>
                    <HStack
                      mt={2}
                      textAlign={'left'}
                      fontWeight={'light'}
                      fontSize="lg"
                      color="white"
                    >
                      <Text fontWeight={'bold'}>Transaction details:</Text>
                      <Text color="complimentary.900">Mintscan</Text>
                    </HStack>
                    <Button
                      w="55%"
                      _hover={{
                        bgColor: '#181818',
                      }}
                      mt={4}
                    >
                      Stake Again
                    </Button>
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
export default StakingProcessModal;
