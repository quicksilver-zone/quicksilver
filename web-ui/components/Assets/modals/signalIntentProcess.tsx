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
  Spinner,
  Input,
  Grid,
} from '@chakra-ui/react';
import { StdFee } from '@cosmjs/amino';
import { useChain } from '@cosmos-kit/react';
import styled from '@emotion/styled';

import { assets } from 'chain-registry';
import { quicksilver } from 'quicksilverjs';

import React, { useEffect, useState } from 'react';

import { IntentMultiModal } from './intentMultiModal';

import { useTx } from '@/hooks';

const ChakraModalContent = styled(ModalContent)`
  position: relative;
  background: none;
  max-height: 320px;
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
  refetch: () => void;
  selectedOption?: {
    name: string;
    value: string;
    logo: string;
    chainName: string;
    chainId: string;
  };
}

interface Intent {
  valoperAddress: string;
  weight: string;
}

export const SignalIntentModal: React.FC<StakingModalProps> = ({ isOpen, onClose, selectedOption, refetch }) => {
  const [step, setStep] = React.useState(1);
  const getProgressColor = (circleStep: number) => {
    if (step >= circleStep) return 'complimentary.900';
    return 'rgba(255,255,255,0.2)';
  };

  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const { address } = useChain('quicksilver' || '');

  const labels = ['Choose validators', `Set weights`, `Sign & Submit`, `Receive q${selectedOption?.value}`];
  const [isModalOpen, setModalOpen] = useState(false);

  const [selectedValidators, setSelectedValidators] = React.useState<{ name: string; operatorAddress: string }[]>([]);

  const advanceStep = () => {
    if (selectedValidators.length > 0) {
      setStep((prevStep) => prevStep + 1);
    }
  };

  const retreatStep = () => {
    if (step === 3) {
      setStep(1); // If on step 3 and checkbox is checked, go back to step 1
    } else {
      setStep((prevStep) => Math.max(prevStep - 1, 1)); // Otherwise, go to the previous step
    }
  };

  const totalWeights = 1;
  const numberOfValidators = selectedValidators.length;

  const [weights, setWeights] = useState<{ [key: string]: number }>({});

  const [isCustomValid, setIsCustomValid] = useState(true);

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

    setIsCustomValid(newTotalWeight === 100); // Validation for custom weights
  };

  // Calculate default weight per validator
  const weightPerValidator = (totalWeights / numberOfValidators).toFixed(4);

  // Initialize intents array
  let intents: Intent[] = [];

  // Assign default or custom weight to each validator
  selectedValidators.forEach((validator, index) => {
    const customWeight = weights[validator.operatorAddress];
    const weight = customWeight !== undefined ? (customWeight / 100).toFixed(4) : weightPerValidator;
    intents.push({
      valoperAddress: validator.operatorAddress,
      weight: weight,
    });
  });

  // Calculate the total assigned weight
  const totalAssignedWeight = intents.reduce((sum, intent) => sum + parseFloat(intent.weight), 0);

  // If the total weight is not equal to 1, adjust the last validator's weight
  if (totalAssignedWeight !== 1 && intents.length > 0) {
    const lastValidatorWeight = parseFloat(intents[intents.length - 1].weight);
    const remainingWeight = (1 - (totalAssignedWeight - lastValidatorWeight)).toFixed(4);
    intents[intents.length - 1].weight = remainingWeight;
  }

  // Create formatted intents string
  const formattedIntentsString = intents.map((intent) => `${intent.weight}${intent.valoperAddress}`).join(',');

  const remainingWeight = (1 - totalAssignedWeight).toFixed(4);

  // Assign the remaining weight to the last validator
  if (selectedValidators.length > 0) {
    const lastValidator = selectedValidators[selectedValidators.length - 1];
    intents.push({
      valoperAddress: lastValidator.operatorAddress,
      weight: remainingWeight,
    });
  }

  const { signalIntent } = quicksilver.interchainstaking.v1.MessageComposer.withTypeUrl;
  const msgSignalIntent = signalIntent({
    chainId: selectedOption?.chainId ?? '',
    intents: formattedIntentsString,
    fromAddress: address ?? '',
  });

  const mainTokens = assets.find(({ chain_name }) => chain_name === 'quicksilver');
  const mainDenom = mainTokens?.assets[0].base ?? 'uqck';

  const fee: StdFee = {
    amount: [
      {
        denom: mainDenom,
        amount: '5000',
      },
    ],
    gas: '500000',
  };

  const { tx } = useTx('quicksilver' ?? '');

  const handleSignalIntent = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);

    try {
      const result = await tx([msgSignalIntent], {
        fee,
        onSuccess: () => {
          refetch();
          onClose();
        },
      });
    } catch (error) {
      console.error('Transaction failed', error);

      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

  useEffect(() => {
    setSelectedValidators([]);
    setStep(1);
    setIsError(false);
    setIsSigning(false);
  }, [selectedOption?.chainName]);

  const [isCustomWeight, setIsCustomWeight] = useState(false);

  const handleCustomWeightMode = () => {
    setIsCustomWeight(true);
  };

  const handleNextInCustomWeightMode = () => {
    if (isCustomValid) {
      setIsCustomWeight(false);
      advanceStep();
    }
  };

  const handleStepOneButtonClick = () => {
    // Check if only one validator is selected
    if (selectedValidators.length === 1) {
      setStep(3); // Skip directly to step 3
    } else {
      if (selectedValidators.length === 0) {
        setModalOpen(true);
      } else {
        advanceStep();
      }
    }
  };

  return (
    (
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
                    <StatLabel color="rgba(255,255,255,0.5)">SIGNAL INTENT</StatLabel>
                  </Stat>
                  {[1, 2, 3].map((circleStep, index) => (
                    <Flex key={circleStep} align="center" mt={10} mb={circleStep !== 3 ? '48px' : '0'}>
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

              {/* Right Section */}
              <Flex width="67%" flexDirection="column" justifyContent="center" alignItems="center">
                {step === 1 && (
                  <>
                    <Flex maxW="300px" flexDirection={'column'} justifyContent={'left'} alignItems={'center'}>
                      <Text textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                        Choose Validators
                      </Text>
                      <Text mt={2} textAlign={'center'} fontWeight={'light'} fontSize="lg" color="white">
                        Select the validators you would like to split your liquid delegation between.
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
                          {selectedValidators.length} Validators Selected
                        </Text>
                      </>
                    )}
                    <Button
                      mt={4}
                      width="55%"
                      _hover={{
                        bgColor: 'complimentary.500',
                      }}
                      onClick={handleStepOneButtonClick}
                    >
                      {selectedValidators.length > 0 ? 'Next' : 'Choose Validators'}
                    </Button>

                    <IntentMultiModal
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
                        _hover={{ bgColor: 'complimentary.500' }}
                        minW={'100px'}
                        onClick={() => {
                          setWeights({});
                          setIsCustomWeight(false);
                          advanceStep();
                        }}
                      >
                        Equal
                      </Button>
                      {selectedValidators.length > 1 && (
                        <Button
                          minW={'100px'}
                          _hover={{
                            bgColor: 'complimentary.500',
                          }}
                          onClick={handleCustomWeightMode}
                        >
                          Custom
                        </Button>
                      )}
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
                    <Box overflowY="auto" maxH="160px">
                      {' '}
                      {/* Set a maximum height to make the box scrollable */}
                      <Grid
                        templateColumns="repeat(auto-fill, minmax(120px, 1fr))"
                        gap={8}
                        maxWidth="400px" // This ensures that no more than 4 items (120px each) are in a row
                      >
                        {selectedValidators.map((validator, index) => (
                          <Flex key={validator.operatorAddress} flexDirection={'column'} alignItems={'center'}>
                            <Text fontSize="sm" color="white" mb={2}>
                              {validator.name.split(' ').length > 1 && validator.name.length > 9
                                ? `${validator.name.split(' ')[0]}...`
                                : validator.name}
                            </Text>
                            <Input
                              _active={{
                                borderColor: 'complimentary.900',
                              }}
                              _selected={{
                                borderColor: 'complimentary.900',
                              }}
                              _hover={{
                                borderColor: 'complimentary.900',
                              }}
                              _focus={{
                                borderColor: 'complimentary.900',
                                boxShadow: '0 0 0 3px #FF8000',
                              }}
                              color="complimentary.900"
                              type="number"
                              width="55px"
                              placeholder="0"
                              onChange={(e) => handleWeightChange(e, validator.operatorAddress)}
                            />
                          </Flex>
                        ))}
                      </Grid>
                    </Box>
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
                      <Text fontWeight={'bold'} fontSize="lg" w="250px" textAlign={'left'} color="white"></Text>
                      {selectedValidators.length > 0 && (
                        <Flex mt={2} textAlign={'left'} alignItems="baseline" gap="2">
                          <Text mt={2} textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                            {selectedValidators.length === 1 ? 'Selected Validator:' : 'Selected Validators:'}
                          </Text>
                          <Text color="complimentary.900">
                            {selectedValidators.length === 1 ? selectedValidators[0].name : `${selectedValidators.length}`}
                          </Text>
                        </Flex>
                      )}
                      <HStack mt={2} textAlign={'left'} fontWeight={'light'} fontSize="lg" color="white"></HStack>
                      <Text mt={2} textAlign={'left'} fontWeight={'hairline'}>
                        Processing time: ~2 minutes
                      </Text>
                      <Button
                        w="55%"
                        _hover={{
                          bgColor: 'complimentary.500',
                        }}
                        mt={4}
                        onClick={handleSignalIntent}
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
              </Flex>
            </HStack>
          </ModalBody>
        </ChakraModalContent>
      </Modal>
    ) || null
  );
};
export default SignalIntentModal;
