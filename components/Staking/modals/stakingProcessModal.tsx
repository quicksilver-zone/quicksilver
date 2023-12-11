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

import { useQueryHooks } from '@/hooks';
import { useZoneQuery } from '@/hooks/useQueries';
import { liquidStakeTx, unbondLiquidStakeTx } from '@/tx/liquidStakeTx';

import { MultiModal } from './validatorSelectionModal';

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
    chainId: string;
  };
}

export const StakingProcessModal: React.FC<StakingModalProps> = ({ isOpen, onClose, selectedOption, tokenAmount }) => {
  const [step, setStep] = React.useState(1);
  const getProgressColor = (circleStep: number) => {
    if (step >= circleStep) return 'complimentary.900';
    return 'rgba(255,255,255,0.2)';
  };

  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const { address, getSigningStargateClient } = useChain(selectedOption?.chainName || '');

  const labels = ['Choose validators', `Set weights`, `Sign & Submit`, `Receive q${selectedOption?.value}`];
  const [isModalOpen, setModalOpen] = useState(false);

  const [selectedValidators, setSelectedValidators] = React.useState<{ name: string; operatorAddress: string }[]>([]);

  const [resp, setResp] = useState('');

  const advanceStep = () => {
    if (selectedValidators.length > 0) {
      setStep((prevStep) => prevStep + 1);
    }
  };

  const retreatStep = () => {
    if (step === 3 && check) {
      setStep(1); // If on step 3 and checkbox is checked, go back to step 1
    } else {
      setStep((prevStep) => Math.max(prevStep - 1, 1)); // Otherwise, go to the previous step
    }
  };

  const toast = useToast();

  const totalWeights = 1;
  const numberOfValidators = selectedValidators.length;

  // Calculate the weight for each validator
  const weightPerValidator = numberOfValidators ? (totalWeights / numberOfValidators).toFixed(4) : '0';

  const [weights, setWeights] = useState<{ [key: string]: number }>({});
  const [totalWeight, setTotalWeight] = useState<string>('0');

  const [isCustomValid, setIsCustomValid] = useState(true);
  const [defaultWeight, setDefaultWeight] = useState(0);

  useEffect(() => {
    // Update the state when selectedValidators changes
    setIsCustomValid(selectedValidators.length === 0);
  }, [selectedValidators]);

  // Modify the handleWeightChange function
  const handleWeightChange = (e: React.ChangeEvent<HTMLInputElement>, validatorName: string) => {
    const value = Number(e.target.value);
    setWeights({
      ...weights,
      [validatorName]: value,
    });

    // Update the total weight as string
    const newTotalWeight = Object.values({ ...weights, [validatorName]: value }).reduce((acc, val) => acc + val, 0);
    setTotalWeight(newTotalWeight.toString());

    setIsCustomValid(newTotalWeight === 100); // Validation for custom weights
  };

  // Calculate defaultWeight as string
  useEffect(() => {
    setDefaultWeight(1 / numberOfValidators);
  }, [numberOfValidators]);

  const [useDefaultWeights, setUseDefaultWeights] = useState(true);

  interface ValidatorsSelect {
    address: string;
    intent: number;
  }

  const intents: ValidatorsSelect[] = selectedValidators.map((validator) => ({
    address: validator.operatorAddress,
    intent: useDefaultWeights ? defaultWeight : weights[validator.operatorAddress] || 0,
  }));

  const { data: zone, isLoading: isZoneLoading, isError: isZoneError } = useZoneQuery(selectedOption?.chainId ?? '');

  const handleLiquidStake = async (event: React.MouseEvent) => {
    const numericAmount = Number(tokenAmount);
    const smallestUnitAmount = numericAmount * Math.pow(10, 6);

    try {
      setIsSigning(true);
      const response = await liquidStakeTx(
        getSigningStargateClient,
        setResp,
        selectedOption?.chainName || '',
        selectedOption?.chainId || '',
        address,
        toast,
        setIsError,
        setIsSigning,
        intents,
        smallestUnitAmount,
        zone,
      )(event);

      // Parse the response
      const parsedResponse = JSON.parse(resp);

      if (parsedResponse && parsedResponse.code === 0) {
        // Successful transaction
        setStep(4);
        setTransactionStatus('Success');
      } else {
        // Unsuccessful transaction
        setIsError(true);
        setTransactionStatus('Failed');
      }
    } catch (error) {
      console.error('Transaction failed', error);
      setIsSigning(false);
      setIsError(true);
      setTransactionStatus('Failed');
    } finally {
      setIsSigning(false);
    }
  };

  //placehoder for transaction status
  const [transactionStatus, setTransactionStatus] = useState('Pending');

  useEffect(() => {
    setSelectedValidators([]);
    setStep(1);
    setIsError(false);
    setIsSigning(false);
    setUseDefaultWeights(true);
  }, [selectedOption?.chainName]);

  const [isCustomWeight, setIsCustomWeight] = useState(false);

  const handleCustomWeightMode = () => {
    setIsCustomWeight(true);
    setUseDefaultWeights(false);
  };

  const handleNextInCustomWeightMode = () => {
    if (isCustomValid) {
      setIsCustomWeight(false);
      advanceStep();
    }
  };

  const [check, setCheck] = useState(false);

  const handleCheck = () => {
    setCheck(!check);
  };

  const handleStepOneButtonClick = () => {
    if (check) {
      // If checkbox is checked, skip directly to step 3
      setStep(3);
    } else {
      // If checkbox is not checked, consider the state of selectedValidators
      if (selectedValidators.length === 0) {
        setModalOpen(true);
      } else {
        advanceStep();
      }
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size={{ base: '3xl', md: '2xl' }}>
      <ModalOverlay />
      <ChakraModalContent h="48%" maxH={'100%'}>
        <ModalBody borderRadius={4} h="48%" maxH={'100%'}>
          <ModalCloseButton color="white" />
          <HStack position={'relative'} h="100%" spacing="48px" align="stretch">
            {/* Left Section */}
            <Flex flexDirection="column" justifyContent="space-between" width="40%" p={4} bg="#1E1C19" height="100%">
              <Box position="relative">
                <Stat>
                  <StatLabel color="rgba(255,255,255,0.5)">LIQUID STAKING</StatLabel>
                  <StatNumber color="white">
                    {tokenAmount} {selectedOption?.value}
                  </StatNumber>
                </Stat>
                {[1, 2, 3, 4].map((circleStep, index) => (
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
                      {circleStep !== 4 && (
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

            {/* Right Section */}
            <Flex width="67%" flexDirection="column" justifyContent="center" alignItems="center">
              {step === 1 && (
                <>
                  <Flex maxW="300px" flexDirection={'column'} justifyContent={'left'} alignItems={'center'}>
                    <Text textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                      Choose Validators
                    </Text>
                    <Text mt={2} textAlign={'center'} fontWeight={'light'} fontSize="lg" color="white">
                      Select up to 8 validators to split your liquid delegation between.
                    </Text>
                  </Flex>
                  {selectedValidators.length > 0 && (
                    <>
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
                      <Text mt={'2'} fontSize={'sm'} fontWeight={'light'}>
                        {selectedValidators.length} / 8 Validators Selected
                      </Text>
                    </>
                  )}
                  <Button
                    mt={4}
                    width="55%"
                    _hover={{
                      bgColor: '#181818',
                    }}
                    onClick={handleStepOneButtonClick}
                  >
                    {check ? 'Skip to Step 3' : selectedValidators.length > 0 ? 'Next' : 'Choose Validators'}
                  </Button>
                  {selectedValidators.length === 0 && (
                    <Flex mt={'6'} flexDir={'row'} gap="3">
                      <Checkbox _selected={{ bgColor: 'transparent' }} isChecked={check} onChange={handleCheck} colorScheme="orange" />
                      <Text>Proceed with existing intent?</Text>
                    </Flex>
                  )}
                  <MultiModal
                    isOpen={isModalOpen}
                    onClose={() => setModalOpen(false)}
                    selectedChainName={selectedOption?.chainName || ''}
                    selectedValidators={selectedValidators}
                    selectedChainId={selectedOption?.chainId || ''}
                    setSelectedValidators={setSelectedValidators}
                  />
                </>
              )}
              {step === 2 && !isCustomWeight && (
                <>
                  <Text textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                    Set Weights
                  </Text>
                  <Text mt={2} textAlign={'center'} fontWeight={'light'} fontSize="lg" color="white">
                    Specifying weights allows you to choose how much of your liquid delegation goes to each validator.
                  </Text>
                  <HStack mt={4} justifyContent={'center'} alignItems={'center'}>
                    <Button
                      onClick={() => {
                        setUseDefaultWeights(true);
                        advanceStep();
                      }}
                    >
                      Default
                    </Button>
                    <Button onClick={handleCustomWeightMode}>Custom</Button>
                  </HStack>
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
              {step === 2 && isCustomWeight && (
                <>
                  <Text textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                    Set Custom Weights
                  </Text>
                  <Text mt={2} textAlign={'center'} fontWeight={'light'} fontSize="lg" color="white">
                    The total weight must equal 100
                  </Text>
                  <Grid mt={2} templateColumns={`repeat(${Math.ceil(Math.sqrt(selectedValidators.length))}, 1fr)`} gap={4}>
                    {selectedValidators.map((validator, index) => (
                      <Flex key={validator.operatorAddress} flexDirection={'column'} alignItems={'center'}>
                        <Text fontSize="sm" color="white" mb={2}>
                          {validator.name}
                        </Text>
                        <Input
                          type="number"
                          width="55px"
                          placeholder=""
                          onChange={(e) => handleWeightChange(e, validator.operatorAddress)}
                        />
                      </Flex>
                    ))}
                  </Grid>
                  <Flex mt={4} justifyContent={'space-between'} width={'100%'} alignItems={'center'}>
                    <Button
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
                      onClick={() => {
                        setIsCustomWeight(false);
                      }}
                    >
                      ←
                    </Button>
                    <Button isDisabled={!isCustomValid} onClick={handleNextInCustomWeightMode}>
                      Next
                    </Button>
                  </Flex>
                </>
              )}

              {step === 3 && (
                <>
                  <Box justifyContent={'center'}>
                    <Text fontWeight={'bold'} fontSize="lg" w="250px" textAlign={'left'} color="white">
                      You’re going to liquid stake {tokenAmount} {selectedOption?.value} on Quicksilver
                    </Text>
                    <HStack mt={2} textAlign={'left'} fontWeight={'light'} fontSize="lg" color="white">
                      <Text fontWeight={'bold'}>Receiving:</Text>
                      <Text color="complimentary.900">
                        {(Number(tokenAmount) / Number(zone?.redemptionRate)).toFixed(2)} q{selectedOption?.value}
                      </Text>
                    </HStack>
                    <Text mt={2} textAlign={'left'} fontWeight={'hairline'}>
                      Processing time: ~2 minutes
                    </Text>
                    <Button
                      w="55%"
                      _hover={{
                        bgColor: '#181818',
                      }}
                      mt={4}
                      onClick={(event) => handleLiquidStake(event)}
                    >
                      {isError ? 'Try Again' : isSigning ? <Spinner /> : 'Confirm'}
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
              {step === 4 && (
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
export default StakingProcessModal;
