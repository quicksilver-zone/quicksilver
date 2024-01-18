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
  Spinner,
} from '@chakra-ui/react';
import { coins, StdFee } from '@cosmjs/amino';
import styled from '@emotion/styled';
import chains from 'chain-registry';
import { assets } from 'chain-registry';
import React, { useEffect, useState } from 'react';
import { cosmos } from 'stridejs';

import { useTx } from '@/hooks';
import { useZoneQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';

const ChakraModalContent = styled(ModalContent)`
  position: relative;
  background: none;
  max-height: 400px;
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
  address: string;
  isTokenized: boolean;
}

export const TransferProcessModal: React.FC<StakingModalProps> = ({
  isOpen,
  onClose,
  selectedOption,
  selectedValidator,
  address,
  isTokenized,
}) => {
  useEffect(() => {
    if (isTokenized === true) {
      setStep(2);
    }
    if (isTokenized === undefined) {
      setStep(1);
    }
  }, [isTokenized, selectedValidator]);
  const [step, setStep] = useState(1);
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

  const { data: zone } = useZoneQuery(selectedOption?.chainId ?? '');
  const labels = ['Tokenize Shares', `Transfer`, `Receive q${selectedOption?.value}`];
  const [transactionStatus, setTransactionStatus] = useState('Pending');
  function truncateString(str: string, num: number) {
    if (str.length > num) {
      return str.slice(0, num) + '...';
    } else {
      return str;
    }
  }

  const { tokenizeShares } = cosmos.staking.v1beta1.MessageComposer.withTypeUrl;

  const msg = tokenizeShares({
    delegatorAddress: address,
    validatorAddress: selectedValidator.operatorAddress,
    amount: {
      denom: 'u' + selectedOption?.value.toLowerCase(),
      amount: selectedValidator.tokenAmount.toString(),
    },
    tokenizedShareOwner: address,
  });

  const mainTokens = assets.find(({ chain_name }) => chain_name === newChainName);
  const fees = chains.chains.find(({ chain_name }) => chain_name === newChainName)?.fees?.fee_tokens;
  const mainDenom = mainTokens?.assets[0].base ?? '';
  const fixedMinGasPrice = fees?.find(({ denom }) => denom === mainDenom)?.high_gas_price ?? '';
  const feeAmount = shiftDigits(fixedMinGasPrice, 6);

  const fee: StdFee = {
    amount: [
      {
        denom: mainDenom,
        amount: feeAmount.toString(),
      },
    ],
    gas: '200000',
  };

  const { tx } = useTx(newChainName ?? '');

  const hanleTokenizeShares = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);
    setTransactionStatus('Pending');
    try {
      const result = await tx([msg], {
        fee,
        onSuccess: () => {
          setStep(2);
          setTransactionStatus('Success');
        },
      });
    } catch (error) {
      console.error('Transaction failed', error);
      setTransactionStatus('Failed');
      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };
  const { send } = cosmos.bank.v1beta1.MessageComposer.withTypeUrl;

  let numericAmount = selectedValidator.tokenAmount;
  if (isNaN(Number(numericAmount)) || Number(numericAmount) <= 0) {
    numericAmount = '0';
  }
  const msgSend = send({
    fromAddress: address ?? '',
    toAddress: zone?.depositAddress?.address ?? '',
    amount: coins(numericAmount, zone?.baseDenom ?? ''),
  });

  const handleSend = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);
    setTransactionStatus('Pending');
    try {
      const result = await tx([msgSend], {
        fee,
        onSuccess: () => {
          setStep(3);
          setTransactionStatus('Success');
        },
      });
    } catch (error) {
      console.error('Transaction failed', error);
      setTransactionStatus('Failed');
      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

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
                    _active={{
                      transform: 'scale(0.95)',
                      color: 'complimentary.800',
                    }}
                    _hover={{
                      bgColor: 'rgba(255,128,0, 0.25)',
                      color: 'complimentary.300',
                    }}
                    onClick={hanleTokenizeShares}
                  >
                    {isError ? 'Try Again' : isSigning ? <Spinner /> : 'Tokenize Shares'}
                  </Button>
                </>
              )}
              {step === 2 && (
                <>
                  <Text textAlign={'center'} fontWeight={'bold'} fontSize="lg" color="white">
                    Send your tokenized shares to Quicksilver
                  </Text>
                  <Button
                    mt={4}
                    width="55%"
                    _active={{
                      transform: 'scale(0.95)',
                      color: 'complimentary.800',
                    }}
                    _hover={{
                      bgColor: 'rgba(255,128,0, 0.25)',
                      color: 'complimentary.300',
                    }}
                    onClick={handleSend}
                  >
                    {isError ? 'Try Again' : isSigning ? <Spinner /> : 'Transfer'}
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
                        _active={{
                          transform: 'scale(0.95)',
                          color: 'complimentary.800',
                        }}
                        _hover={{
                          bgColor: 'rgba(255,128,0, 0.25)',
                          color: 'complimentary.300',
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
