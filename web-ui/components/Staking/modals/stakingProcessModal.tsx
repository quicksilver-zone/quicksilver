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
  Input,
  Grid,
  Checkbox,
} from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import styled from '@emotion/styled';
import { bech32 } from 'bech32';
import { cosmos } from 'quicksilverjs';
import React, { useEffect, useState } from 'react';

import { Chain, env, chains as configChains, local_chain } from '@/config';
import { useTx } from '@/hooks';
import { useFeeEstimation } from '@/hooks/useFeeEstimation';
import { useZoneQuery } from '@/hooks/useQueries';

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
  selectedOption?: Chain;
  address: string;
  refetch: () => void;
}

export const StakingProcessModal: React.FC<StakingModalProps> = ({ isOpen, onClose, selectedOption, tokenAmount, address, refetch }) => {
  const [step, setStep] = useState(1);
  const getProgressColor = (circleStep: number) => {
    if (step >= circleStep) return 'complimentary.700';
    return 'rgba(255,255,255,0.2)';
  };

  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const newChainName = selectedOption?.chain_name;
  

  const labels = ['Choose validators', `Set weights`, `Sign & Submit`, `Receive q${selectedOption?.major_denom.toUpperCase()}`];
  const [isModalOpen, setModalOpen] = useState(false);

  const [selectedValidators, setSelectedValidators] = useState<{ name: string; operatorAddress: string }[]>([]);

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

  const numberOfValidators = selectedValidators.length;

  // Calculate the weight for each validator

  const [weights, setWeights] = useState<{ [key: string]: number }>({});

  const [isCustomValid, setIsCustomValid] = useState(true);
  const [defaultWeight, setDefaultWeight] = useState(0);

  useEffect(() => {
    // Update the state when selectedValidators changes
    setIsCustomValid(selectedValidators.length === 0);
  }, [selectedValidators]);

  // Modify the handleWeightChange function
  const handleWeightChange = (e: React.ChangeEvent<HTMLInputElement>, validatorName: string) => {
    const value = Number(e.target.value);
    setWeights((prevWeights) => ({
      ...prevWeights,
      [validatorName]: value,
    }));

    // Update the total weight as string
    const newTotalWeight = Object.values({ ...weights, [validatorName]: value }).reduce((acc, val) => acc + val, 0);

    setIsCustomValid(newTotalWeight === 100);
  };

  const calculateIntents = () => {
    return selectedValidators.map((validator) => {
      // For each validator, calculate the weight based on whether default weights are used
      const weight = useDefaultWeights ? defaultWeight : weights[validator.operatorAddress] || 0;

      return {
        address: validator.operatorAddress,
        intent: weight.toFixed(4),
      };
    });
  };

  // Calculate defaultWeight as string
  useEffect(() => {
    setDefaultWeight(1 / numberOfValidators);
    //eslint-disable-next-line react-hooks/exhaustive-deps
  }, [numberOfValidators]);

  const [useDefaultWeights, setUseDefaultWeights] = useState(true);

  useEffect(() => {
    if (!useDefaultWeights && selectedValidators.length > 0) {
      const totalWeight = calculateIntents().reduce((acc, intent) => acc + parseFloat(intent.intent), 0);
      if (totalWeight !== 1) {
        const lastValidator = selectedValidators[selectedValidators.length - 1];
        setWeights((prevWeights) => ({
          ...prevWeights,
          [lastValidator.operatorAddress]: (1 - (totalWeight - (weights[lastValidator.operatorAddress] ?? 0) / 100)) * 100,
        }));
      }
    }
    //eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedValidators, weights, useDefaultWeights]);

  interface ValidatorsSelect {
    address: string;
    intent: number;
  }

  const intents: ValidatorsSelect[] = selectedValidators.map((validator) => {
    const weightAsFraction = useDefaultWeights ? defaultWeight : (weights[validator.operatorAddress] ?? 0) / 100;

    return {
      address: validator.operatorAddress,
      intent: weightAsFraction,
    };
  });

  const { data: zone, isLoading: isZoneLoading } = useZoneQuery(selectedOption?.chain_id ?? '');
  const chain = configChains.get(env)?.get(selectedOption?.chain_name ?? '');
  const valToByte = (val: number) => {
    if (val > 1) {
      val = 1;
    }
    if (val < 0) {
      val = 0;
    }
    return Math.abs(val * 200);
  };

  const addValidator = (valAddr: string, weight: number) => {
    let { words } = bech32.decode(valAddr);
    let wordsBuffer = new Uint8Array(bech32.fromWords(words));
    let weightByteBuffer = new Uint8Array([valToByte(weight)]);
    return Buffer.concat([weightByteBuffer, wordsBuffer]);
  };

  let memoBuffer = Buffer.alloc(0);

  if (intents.length > 0) {
    intents.forEach((val) => {
      memoBuffer = Buffer.concat([memoBuffer, addValidator(val.address, val.intent)]);
    });
    memoBuffer = Buffer.concat([Buffer.from([0x02, memoBuffer.length]), memoBuffer]);
  }
  const { address: qsAddress } = useChain(local_chain.get(env)?.chain_name ?? '');

  if (chain?.is_118 == false && qsAddress) {
    let { words } = bech32.decode(qsAddress ?? '');
    let wordsBuffer = Buffer.from(new Uint8Array(bech32.fromWords(words)));
    
    memoBuffer = Buffer.concat([memoBuffer, Buffer.from([0x00, wordsBuffer.length]), wordsBuffer]);
  }

  let memo = memoBuffer.length > 0 ? memoBuffer.toString('base64') : '';

  const parsedAmount = parseFloat(tokenAmount ?? '0');
  let numericAmount = BigInt(0);
  if (!Number.isNaN(parsedAmount)) {
    numericAmount = BigInt(Math.floor(parsedAmount * Math.pow(10, Number(zone?.decimals ?? '6'))));
  }

  if (numericAmount <= 0) {
    numericAmount = BigInt(0);
  }

  const { send } = cosmos.bank.v1beta1.MessageComposer.withTypeUrl;

  const msgSend = send({
    fromAddress: address ?? '',
    toAddress: zone?.depositAddress?.address ?? '',
    amount: [{ denom: zone?.baseDenom ?? '', amount: numericAmount.toString() }],
  });

  const { tx } = useTx(newChainName ?? '');
  
  const { estimateFee } = useFeeEstimation(newChainName ?? '');

  const handleLiquidStake = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);
    setTransactionStatus('Pending');

    console.log(msgSend)
    const feeAmountQuery = await estimateFee(address, [msgSend]);
    console.log(feeAmountQuery)
    
    try {
      await tx([msgSend], {
        memo,
        fee: feeAmountQuery,
        onSuccess: () => {
          refetch();
          setStep(4);
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

  //placehoder for transaction status
  const [transactionStatus, setTransactionStatus] = useState('Pending');

  useEffect(() => {
    setSelectedValidators([]);
    setStep(1);
    setIsError(false);
    setIsSigning(false);
    setUseDefaultWeights(true);
  }, [selectedOption?.chain_name]);

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
    // Check if only one validator is selected
    if (selectedValidators.length === 1) {
      setUseDefaultWeights(true);
      setStep(3); // Skip directly to step 3
    } else if (check) {
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

  const handleEqualWeightAssignment = () => {
    const numberOfValidators = selectedValidators.length;
    const equalWeight = (1 / numberOfValidators).toFixed(4);

    // Update the state with new weights
    setDefaultWeight(Number(equalWeight));
    setUseDefaultWeights(true);
    advanceStep();
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size={{ base: '3xl', md: '2xl' }}>
      <ModalOverlay />
      <ChakraModalContent position={'absolute'} h={{ base: '65%', md: '50%' }} maxH={'550px'}>
        <ModalBody borderRadius={4} h="50%" maxH={'100%'}>
          <ModalCloseButton zIndex={1000} color="white" />
          <HStack position={'relative'} h="100%" spacing="48px" align="stretch">
            {/* Left Section */}
            <Flex flexDirection="column" justifyContent="space-between" width="40%" p={4} bg="#1E1C19" height="100%">
              <Box position="relative">
                <Stat>
                  <StatLabel color="rgba(255,255,255,0.5)">LIQUID STAKING</StatLabel>
                  <StatNumber color="white">
                    {tokenAmount} {selectedOption?.major_denom.toUpperCase()}
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
                          bgColor: 'rgba(255, 128, 0, 0.5)',
                        }}
                        bgColor="rgba(255, 128, 0, 0.25)"
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
                    width={{ base: '80%', md: '55%' }}
                    _active={{
                      transform: 'scale(0.95)',
                      color: 'complimentary.800',
                    }}
                    _hover={{
                      bgColor: 'rgba(255,128,0, 0.25)',
                      color: 'complimentary.300',
                    }}
                    onClick={handleStepOneButtonClick}
                  >
                    {check ? 'Sign & Submit' : selectedValidators.length > 0 ? 'Next' : 'Choose Validators'}
                  </Button>
                  {selectedValidators.length === 0 && (
                    <Flex mt={'6'} flexDir={'row'} gap="3">
                      <Checkbox
                        _selected={{ bgColor: 'transparent' }}
                        _active={{
                          borderColor: 'complimentary.700',
                        }}
                        _hover={{
                          borderColor: 'complimentary.700',
                        }}
                        _focus={{
                          borderColor: 'complimentary.700',
                          boxShadow: '0 0 0 3px #FF8000',
                        }}
                        isChecked={check}
                        onChange={handleCheck}
                        colorScheme="orange"
                      />
                      <Text>Proceed with existing intent?</Text>
                    </Flex>
                  )}
                  <MultiModal
                    isOpen={isModalOpen}
                    onClose={() => setModalOpen(false)}
                    selectedChainName={selectedOption?.chain_name || ''}
                    selectedValidators={selectedValidators}
                    selectedChainId={selectedOption?.chain_id || ''}
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
                      _active={{
                        transform: 'scale(0.95)',
                        color: 'complimentary.800',
                      }}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                        color: 'complimentary.300',
                      }}
                      minW={'100px'}
                      onClick={handleEqualWeightAssignment}
                    >
                      Equal
                    </Button>
                    {selectedValidators.length > 1 && (
                      <Button
                        minW={'100px'}
                        _active={{
                          transform: 'scale(0.95)',
                          color: 'complimentary.800',
                        }}
                        _hover={{
                          bgColor: 'rgba(255,128,0, 0.25)',
                          color: 'complimentary.300',
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
                    fontSize={'2xl'}
                    bgColor="none"
                    _hover={{
                      bgColor: 'none',
                      color: 'complimentary.700',
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
                  <Grid mt={2} templateColumns={`repeat(${Math.ceil(Math.sqrt(selectedValidators.length))}, 1fr)`} gap={8}>
                    {selectedValidators.map((validator, index) => (
                      <Flex key={validator.operatorAddress} flexDirection={'column'} alignItems={'center'}>
                        <Text fontSize="sm" color="white" mb={2}>
                          {validator.name.split(' ').length > 1 && validator.name.length > 9
                            ? `${validator.name.split(' ')[0]}...`
                            : validator.name}
                        </Text>
                        <Input
                          _active={{
                            borderColor: 'complimentary.700',
                          }}
                          _selected={{
                            borderColor: 'complimentary.700',
                          }}
                          _hover={{
                            borderColor: 'complimentary.700',
                          }}
                          _focus={{
                            borderColor: 'complimentary.700',
                            boxShadow: '0 0 0 3px #FF8000',
                          }}
                          color="complimentary.700"
                          type="number"
                          width="55px"
                          placeholder="0"
                          onChange={(e) => handleWeightChange(e, validator.operatorAddress)}
                        />
                      </Flex>
                    ))}
                  </Grid>
                  <Flex mt={4} justifyContent={'space-between'} width={'100%'} alignItems={'center'}>
                    <Button
                      bgColor="none"
                      fontSize={'2xl'}
                      _hover={{
                        bgColor: 'none',
                        color: 'complimentary.700',
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
                    <Button
                      _active={{
                        transform: 'scale(0.95)',
                        color: 'complimentary.800',
                      }}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                        color: 'complimentary.300',
                      }}
                      isDisabled={!isCustomValid}
                      onClick={handleNextInCustomWeightMode}
                    >
                      Next
                    </Button>
                  </Flex>
                </>
              )}

              {step === 3 && (
                <>
                  <Box justifyContent={'center'}>
                    <Text
                      display={{ base: 'none', md: 'block' }}
                      fontWeight={'bold'}
                      fontSize="lg"
                      w="250px"
                      textAlign={'left'}
                      color="white"
                    >
                      You’re going to liquid stake {tokenAmount} {selectedOption?.major_denom.toUpperCase()} on Quicksilver
                    </Text>
                    <Text
                      display={{ base: 'block', md: 'none' }}
                      fontWeight={'bold'}
                      fontSize="md"
                      w="250px"
                      textAlign={'left'}
                      color="white"
                    >
                      Liquid staking {tokenAmount} {selectedOption?.major_denom.toUpperCase()}
                    </Text>
                    {selectedValidators.length > 0 && (
                      <Flex mt={2} textAlign={'left'} alignItems="baseline" gap="2">
                        <Text mt={2} textAlign={'left'} fontWeight={'bold'} fontSize="lg" color="white">
                          {selectedValidators.length === 1 ? 'Selected Validator:' : 'Selected Validators:'}
                        </Text>
                        <Text color="complimentary.700">
                          {selectedValidators.length === 1 ? selectedValidators[0].name : `${selectedValidators.length} / 8`}
                        </Text>
                      </Flex>
                    )}
                    <HStack mt={2} textAlign={'left'} fontWeight={'light'} fontSize="lg" color="white">
                      <Text fontWeight={'bold'}>Receiving:</Text>
                      <Text color="complimentary.700">
                        {!isZoneLoading ? (
                          `${(Number(tokenAmount) / Number(zone?.redemptionRate || 1)).toFixed(2)} q${selectedOption?.major_denom.toUpperCase()}`
                        ) : (
                          <Spinner thickness="2px" speed="0.65s" emptyColor="gray.200" color="complimentary.700" size="sm" />
                        )}
                      </Text>
                    </HStack>
                    <Text mt={2} textAlign={'left'} fontWeight={'hairline'}>
                      Processing time: ~2 minutes
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
                      onClick={(event) => handleLiquidStake(event)}
                    >
                      {isError ? 'Try Again' : isSigning ? <Spinner /> : 'Confirm'}
                    </Button>
                  </Box>
                  <Button
                    position={'absolute'}
                    bottom={3}
                    fontSize={'2xl'}
                    left={'51%'}
                    bgColor="none"
                    _hover={{
                      bgColor: 'none',
                      color: 'complimentary.700',
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
                        Your q{selectedOption?.major_denom.toUpperCase()} will arrive to your wallet in a few minutes.
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
export default StakingProcessModal;
