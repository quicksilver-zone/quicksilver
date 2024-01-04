import {
  Box,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  VStack,
  Text,
  Flex,
  Stat,
  StatLabel,
  StatNumber,
  Input,
  Divider,
  HStack,
  Button,
  Spacer,
  Skeleton,
  SkeletonText,
  useToast,
  SlideFade,
  Spinner,
  FormControl,
  FormLabel,
  Switch,
  Tooltip,
  Image,
  Icon,
  SkeletonCircle,
} from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import React, { useEffect, useState } from 'react';

import {
  useBalanceQuery,
  useNativeStakeQuery,
  useQBalanceQuery,
  useValidatorLogos,
  useValidatorsQuery,
  useZoneQuery,
} from '@/hooks/useQueries';

import { getExponent } from '@/utils';
import { shiftDigits } from '@/utils';

import StakingProcessModal from './modals/stakingProcessModal';
import { Coin, StdFee } from '@cosmjs/amino';
import { quicksilver } from 'quicksilverjs';
import { useTx } from '@/hooks';

import { InfoOutlineIcon } from '@chakra-ui/icons';
import TransferProcessModal from './modals/transferProcessModal';

type StakingBoxProps = {
  selectedOption: {
    name: string;
    value: string;
    logo: string;
    chainName: string;
    chainId: string;
  };
  isModalOpen: boolean;
  setModalOpen: (isOpen: boolean) => void;
  setBalance: (balance: string) => void;
  setQBalance: (qBalance: string) => void;
};

export const StakingBox = ({ selectedOption, isModalOpen, setModalOpen, setBalance, setQBalance }: StakingBoxProps) => {
  const [activeTabIndex, setActiveTabIndex] = useState(0);
  const [tokenAmount, setTokenAmount] = useState<string>('0');
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
    // Default case
    newChainName = selectedOption?.chainName;
  }
  const { address } = useChain(newChainName);
  const { address: qAddress } = useChain('quicksilver');
  const exp = getExponent(selectedOption.chainName);
  const { balance, isLoading } = useBalanceQuery(selectedOption.chainName, address ?? '');
  const {
    balance: qBalance,
    isLoading: qIsLoading,
    isError: qIsError,
  } = useQBalanceQuery('quicksilver', qAddress ?? '', selectedOption.value.toLowerCase());

  const qAssets = qBalance?.balance.amount || '';

  const baseBalance = shiftDigits(balance?.balance?.amount || '0', -exp);

  const { data: zone, isLoading: isZoneLoading, isError: isZoneError } = useZoneQuery(selectedOption.chainId);

  useEffect(() => {
    setQBalance(qAssets);
  }, [qAssets, setQBalance, selectedOption.chainName]);

  useEffect(() => {
    setBalance(baseBalance);
  }, [baseBalance, setBalance]);

  useEffect(() => {
    setTokenAmount('0');
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedOption.chainName]);

  const truncateToThreeDecimals = (num: number) => {
    return Math.trunc(num * 1000) / 1000;
  };

  const truncatedBalance = truncateToThreeDecimals(Number(baseBalance));

  const maxStakingAmount = truncateToThreeDecimals(truncatedBalance ? truncatedBalance - 0.005 : 0);

  const maxHalfStakingAmount = maxStakingAmount / 2;

  const [inputError, setInputError] = useState(false);

  const qAssetsExponent = shiftDigits(qAssets, -6);
  const qAssetsDisplay = qAssetsExponent.includes('.') ? qAssetsExponent.substring(0, qAssetsExponent.indexOf('.') + 3) : qAssetsExponent;

  const maxUnstakingAmount = truncateToThreeDecimals(Number(qAssetsDisplay));
  const halfUnstakingAmount = maxUnstakingAmount / 2;

  const [isSigning, setIsSigning] = useState<boolean>(false);

  const [isError, setIsError] = useState<boolean>(false);
  const [transactionStatus, setTransactionStatus] = useState('Pending');

  const env = process.env.NEXT_PUBLIC_CHAIN_ENV;
  const quicksilverChainName = env === 'testnet' ? 'quicksilvertestnet' : 'quicksilver';

  const isCalculationDataLoaded = tokenAmount && !isNaN(Number(tokenAmount)) && zone && !isNaN(Number(zone.redemptionRate));

  const { requestRedemption } = quicksilver.interchainstaking.v1.MessageComposer.withTypeUrl;
  const numericAmount = Number(tokenAmount);
  const smallestUnitAmount = numericAmount * Math.pow(10, 6);
  const value: Coin = { amount: smallestUnitAmount.toFixed(0), denom: zone?.localDenom ?? '' };
  const msgRequestRedemption = requestRedemption({
    value: value,
    fromAddress: qAddress ?? '',
    destinationAddress: address ?? '',
  });

  const fee: StdFee = {
    amount: [
      {
        denom: 'uqck',
        amount: '7500',
      },
    ],
    gas: '500000',
  };

  const { tx } = useTx(quicksilverChainName);

  const handleLiquidUnstake = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);
    try {
      const result = await tx([msgRequestRedemption], {
        fee,
        onSuccess: () => {
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

  const handleTabsChange = (index: number) => {
    setActiveTabIndex(index);
    setTokenAmount('');
  };

  const isValidNumber = !isNaN(Number(qAssetsDisplay)) && qAssetsDisplay !== '';

  const { delegations, delegationsIsError, delegationsIsLoading } = useNativeStakeQuery(newChainName, address ?? '');
  const delegationsResponse = delegations?.delegation_responses;
  const nativeStakedAmount = delegationsResponse?.reduce((acc, delegationResponse) => {
    const amount = Number(delegationResponse?.balance?.amount) || 0;
    return acc + amount;
  }, 0);

  const [useNativeStake, setUseNativeStake] = useState(false);

  const handleSwitchChange = (event: { target: { checked: boolean | ((prevState: boolean) => boolean) } }) => {
    setUseNativeStake(event.target.checked);
  };

  const { validatorsData, isLoading: validatorsDataLoading, isError: validatorsDataError } = useValidatorsQuery(newChainName);

  const { data: logos, isLoading: isFetchingLogos } = useValidatorLogos(newChainName, validatorsData || []);
  const [selectedValidator, setSelectedValidator] = useState<string | null>(null);

  return (
    <Box position="relative" backdropFilter="blur(50px)" bgColor="rgba(255,255,255,0.1)" flex="1" borderRadius="10px" p={5}>
      <Tabs isFitted variant="enclosed" onChange={handleTabsChange}>
        <TabList mt={'4'} mb="1em" overflow="hidden" borderBottomColor="transparent" bg="rgba(255,255,255,0.1)" p={2} borderRadius="25px">
          <Tab
            borderRadius="25px"
            flex="1"
            color="white"
            fontWeight="bold"
            transition="background-color 0.2s ease-in-out, color 0.2s ease-in-out, border-color 0.2s ease-in-out"
            _hover={{
              borderBottomColor: 'complimentary.900',
            }}
            _selected={{
              bgColor: 'rgba(0,0,0,0.5)',
              color: 'complimentary.900',
              borderColor: 'complimentary.900',
            }}
          >
            Stake
          </Tab>
          <Tab
            borderRadius="25px"
            flex="1"
            color="white"
            fontWeight="bold"
            transition="background-color 0.2s ease-in-out, color 0.2s ease-in-out, border-color 0.2s ease-in-out"
            _hover={{
              borderBottomColor: 'complimentary.900',
            }}
            _selected={{
              bgColor: 'rgba(0,0,0,0.5)',
              color: 'complimentary.900',
              borderColor: 'complimentary.900',
            }}
          >
            Unstake
          </Tab>
        </TabList>
        <TabPanels>
          <SlideFade offsetY="-80px" in={activeTabIndex === 0}>
            <TabPanel>
              <VStack spacing={8} align="center">
                <Text fontWeight="light" textAlign="center" color="white">
                  Stake your {selectedOption.value.toUpperCase()} tokens in exchange for q{selectedOption.value.toUpperCase()} which you can
                  deploy around the ecosystem. You can liquid stake half of your balance, if you&apos;re going to LP.
                </Text>
                {selectedOption.name === 'Cosmos Hub' && (
                  <Flex textAlign={'left'} justifyContent={'flex-start'}>
                    <HStack>
                      <Text fontWeight="medium" textAlign="center" color="white">
                        Use natively staked&nbsp;
                        <span style={{ color: '#FF8000' }}>{selectedOption.value}</span>?
                      </Text>
                      <Switch
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
                        isDisabled={!nativeStakedAmount}
                        isChecked={useNativeStake}
                        onChange={handleSwitchChange}
                        id="use-natively-staked"
                        colorScheme="orange"
                      />

                      <Tooltip
                        label={
                          !address
                            ? 'Please connect your wallet to enable this option.'
                            : !nativeStakedAmount
                            ? "You don't have any native staked tokens."
                            : `You currently have ${shiftDigits(nativeStakedAmount, -6)} ${
                                selectedOption.value
                              } natively staked to ${delegationsResponse?.length} validators. You can tokenize your shares and transfer them to quicksilver by clicking the switch and selecting a validator.`
                        }
                      >
                        <InfoOutlineIcon color="complimentary.900" />
                      </Tooltip>
                    </HStack>
                  </Flex>
                )}
                {!useNativeStake && (
                  <>
                    <Flex flexDirection="column" w="100%">
                      <Stat py={4} textAlign="left" color="white">
                        <StatLabel>Amount to stake:</StatLabel>
                        <StatNumber>{selectedOption.value.toUpperCase()} </StatNumber>
                      </Stat>
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
                        textAlign={'right'}
                        placeholder={inputError ? 'Invalid Number' : 'amount'}
                        _placeholder={{
                          color: inputError ? 'red.500' : 'grey',
                        }}
                        value={tokenAmount}
                        type="text"
                        onChange={(e) => {
                          // Allow any numeric input
                          const validNumberPattern = /^\d*\.?\d*$/;
                          if (validNumberPattern.test(e.target.value)) {
                            setTokenAmount(e.target.value);
                          }
                        }}
                        onBlur={() => {
                          let inputValue = parseFloat(tokenAmount);
                          if (isNaN(inputValue) || inputValue <= 0) {
                            // Set error for invalid or non-positive numbers
                            setInputError(true);
                            setTokenAmount('');
                          } else if (inputValue > maxStakingAmount) {
                            // Limit the input to the max staking amount
                            setInputError(false);
                            setTokenAmount(maxStakingAmount.toString());
                          } else {
                            // Valid input
                            setInputError(false);
                            setTokenAmount(inputValue.toString());
                          }
                        }}
                      />

                      <Flex w="100%" flexDirection="row" py={4} mb={-4} justifyContent="space-between" alignItems="center">
                        <Flex mb={-4} alignItems="center" justifyContent={'center'} gap={4} flexDirection={'row'}>
                          {address ? (
                            <>
                              <Text color="white" fontWeight="light">
                                Tokens available:{' '}
                              </Text>
                              {isLoading ? (
                                <Skeleton startColor="complimentary.900" endColor="complimentary.400">
                                  <SkeletonText w={'95px'} noOfLines={1} skeletonHeight={'18px'} />
                                </Skeleton>
                              ) : (
                                <Text color="complimentary.900" fontWeight="light">
                                  {address
                                    ? balance?.balance?.amount && Number(balance?.balance?.amount) !== 0
                                      ? `${truncatedBalance} ${selectedOption.value.toUpperCase()}`
                                      : `Get ${selectedOption.value.toUpperCase()} tokens here`
                                    : '0'}
                                </Text>
                              )}
                            </>
                          ) : (
                            <Text color="complimentary.900" fontWeight="light">
                              Connect your wallet to stake
                            </Text>
                          )}
                        </Flex>
                        <HStack mb={-4} spacing={2}>
                          <Button
                            _hover={{
                              bgColor: 'rgba(255,255,255,0.05)',
                              backdropFilter: 'blur(10px)',
                            }}
                            _active={{
                              bgColor: 'rgba(255,255,255,0.05)',
                              backdropFilter: 'blur(10px)',
                            }}
                            color="complimentary.900"
                            variant="ghost"
                            w="60px"
                            h="30px"
                            onClick={() => setTokenAmount(maxHalfStakingAmount.toString())}
                            isDisabled={!balance || Number(balance) < 1}
                          >
                            half
                          </Button>
                          <Button
                            _hover={{
                              bgColor: 'rgba(255,255,255,0.05)',
                              backdropFilter: 'blur(10px)',
                            }}
                            _active={{
                              bgColor: 'rgba(255,255,255,0.05)',
                              backdropFilter: 'blur(10px)',
                            }}
                            color="complimentary.900"
                            variant="ghost"
                            w="60px"
                            h="30px"
                            onClick={() => setTokenAmount(maxStakingAmount.toString())}
                            isDisabled={!balance || Number(balance) < 1}
                          >
                            max
                          </Button>
                        </HStack>
                      </Flex>
                    </Flex>
                    <Divider bgColor="complimentary.900" />
                    <HStack justifyContent="space-between" alignItems="left" w="100%" mt={-8}>
                      <Stat textAlign="left" color="white">
                        <StatLabel>What you&apos;ll get</StatLabel>
                        <StatNumber>q{selectedOption.value.toUpperCase()}:</StatNumber>
                      </Stat>
                      <Spacer /> {/* This pushes the next Stat component to the right */}
                      <Stat py={4} textAlign="right" color="white">
                        <StatNumber textColor="complimentary.900">
                          {(Number(tokenAmount) / (Number(zone?.redemptionRate) || 1)).toFixed(2)}
                        </StatNumber>
                      </Stat>
                    </HStack>
                    <Button
                      width="100%"
                      _hover={{
                        bgColor: 'complimentary.1000',
                      }}
                      onClick={() => setModalOpen(true)}
                      isDisabled={Number(tokenAmount) === 0 || !address}
                    >
                      Liquid stake
                    </Button>
                    <StakingProcessModal
                      tokenAmount={tokenAmount}
                      isOpen={isModalOpen}
                      onClose={() => setModalOpen(false)}
                      selectedOption={selectedOption}
                    />
                  </>
                )}
                {useNativeStake && (
                  <Flex flexDirection="column" w="100%">
                    <VStack spacing={8} align="center">
                      <Box maxH="300px" overflowY="scroll" w="fit-content" mb={8}>
                        {delegationsResponse?.map((delegation, index) => {
                          const validator = validatorsData?.find((v) => v.address === delegation.delegation.validator_address);
                          const isSelected = validator && validator.address === selectedValidator;
                          const validatorLogo = logos[delegation.delegation.validator_address ?? ''];

                          return (
                            <Box
                              borderRadius={'md'}
                              as="button"
                              w="full"
                              onClick={() => setSelectedValidator(validator?.address ?? '')}
                              _hover={{ bg: 'rgba(255, 128, 0, 0.25)' }}
                              bg={isSelected ? 'rgba(255, 128, 0, 0.25)' : 'transparent'}
                              key={index}
                              mb={2}
                            >
                              <Flex py={2} align="center">
                                <Box boxSize="50px" borderRadius="md" ml={4}>
                                  {!validatorLogo ? (
                                    <SkeletonCircle size="8" startColor="complimentary.900" endColor="complimentary.400" />
                                  ) : (
                                    <Image
                                      src={validatorLogo}
                                      alt={validator?.name}
                                      boxSize="50px"
                                      objectFit="cover"
                                      borderRadius={'full'}
                                    />
                                  )}
                                </Box>
                                <VStack align="start" ml={2}>
                                  <Text fontSize="md">{validator ? validator.name : 'Validator'}</Text>
                                  <Text color={'complimentary.900'} fontSize="md">
                                    {shiftDigits(delegation.balance.amount, -6)} {selectedOption.value}
                                  </Text>
                                </VStack>
                              </Flex>
                            </Box>
                          );
                        })}
                      </Box>
                    </VStack>
                    <Button
                      width="100%"
                      _hover={{
                        bgColor: 'complimentary.1000',
                      }}
                      onClick={() => setModalOpen(true)}
                      isDisabled={!selectedValidator || !address}
                    >
                      Transfer Existing Delegation
                    </Button>
                    <TransferProcessModal
                      tokenAmount={tokenAmount}
                      isOpen={isModalOpen}
                      onClose={() => setModalOpen(false)}
                      selectedOption={selectedOption}
                    />
                  </Flex>
                )}
              </VStack>
            </TabPanel>
          </SlideFade>
          <SlideFade offsetY="200px" in={activeTabIndex === 1}>
            <TabPanel>
              <VStack spacing={8} align="center">
                <Text fontWeight="light" textAlign="center" color="white">
                  Unstake your q{selectedOption.value.toUpperCase()} tokens in exchange for {selectedOption.value.toUpperCase()}.
                </Text>
                <Flex flexDirection="column" w="100%">
                  <Stat py={4} textAlign="left" color="white">
                    <StatLabel>Amount to unstake:</StatLabel>
                    <StatNumber>q{selectedOption.value.toUpperCase()} </StatNumber>
                  </Stat>
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
                    textAlign={'right'}
                    placeholder={inputError ? 'Invalid Number' : 'amount'}
                    _placeholder={{
                      color: inputError ? 'red.500' : 'grey',
                    }}
                    value={tokenAmount}
                    type="text"
                    onChange={(e) => {
                      // Allow any numeric input
                      const validNumberPattern = /^\d*\.?\d*$/;
                      if (validNumberPattern.test(e.target.value)) {
                        setTokenAmount(e.target.value);
                      }
                    }}
                    onBlur={() => {
                      let inputValue = parseFloat(tokenAmount);
                      if (isNaN(inputValue) || inputValue <= 0) {
                        // Set error for invalid or non-positive numbers
                        setInputError(true);
                        setTokenAmount('');
                      } else if (inputValue > maxUnstakingAmount) {
                        // Limit the input to the max staking amount
                        setInputError(false);
                        setTokenAmount(maxUnstakingAmount.toString());
                      } else {
                        // Valid input
                        setInputError(false);
                        setTokenAmount(inputValue.toString());
                      }
                    }}
                  />
                  <Flex w="100%" flexDirection="row" py={4} mb={-4} justifyContent="space-between" alignItems="center">
                    {address ? (
                      <Flex mb={-4} alignItems="center" justifyContent={'center'} gap={4} flexDirection={'row'}>
                        <Text color="white" fontWeight="light">
                          Tokens available:{' '}
                        </Text>
                        {isLoading ? (
                          <Skeleton startColor="complimentary.900" endColor="complimentary.400">
                            <SkeletonText w={'95px'} noOfLines={1} skeletonHeight={'18px'} />
                          </Skeleton>
                        ) : (
                          <Text color="complimentary.900" fontWeight="light">
                            {address
                              ? qAssets && Number(qAssets) !== 0
                                ? `${qAssetsDisplay} ${selectedOption.value.toUpperCase()}`
                                : `No q${selectedOption.value.toUpperCase()}`
                              : '0'}
                          </Text>
                        )}
                      </Flex>
                    ) : (
                      <Text color="complimentary.900" fontWeight="light">
                        Connect your wallet to unstake
                      </Text>
                    )}

                    <HStack mb={-4} spacing={2}>
                      <Button
                        _hover={{
                          bgColor: 'rgba(255,255,255,0.05)',
                          backdropFilter: 'blur(10px)',
                        }}
                        _active={{
                          bgColor: 'rgba(255,255,255,0.05)',
                          backdropFilter: 'blur(10px)',
                        }}
                        color="complimentary.900"
                        variant="ghost"
                        w="60px"
                        h="30px"
                        onClick={() => setTokenAmount(halfUnstakingAmount.toString())}
                        isDisabled={!qAssets || Number(qAssets) < 1}
                      >
                        half
                      </Button>
                      <Button
                        _hover={{
                          bgColor: 'rgba(255,255,255,0.05)',
                          backdropFilter: 'blur(10px)',
                        }}
                        _active={{
                          bgColor: 'rgba(255,255,255,0.05)',
                          backdropFilter: 'blur(10px)',
                        }}
                        color="complimentary.900"
                        variant="ghost"
                        w="60px"
                        h="30px"
                        onClick={() => setTokenAmount(maxUnstakingAmount.toString())}
                        isDisabled={!qAssets || Number(qAssets) < 1}
                      >
                        max
                      </Button>
                    </HStack>
                  </Flex>
                </Flex>
                <Divider bgColor="complimentary.900" />
                <HStack justifyContent="space-between" alignItems="left" w="100%" mt={-8}>
                  <Stat textAlign="left" color="white">
                    <StatLabel>What you&apos;ll get</StatLabel>
                    <StatNumber>{selectedOption.value.toUpperCase()}:</StatNumber>
                  </Stat>
                  <Spacer /> {/* This pushes the next Stat component to the right */}
                  <Stat py={4} textAlign="right" color="white">
                    <StatNumber textColor="complimentary.900">
                      {(Number(tokenAmount) * Number(zone?.redemptionRate || 1)).toFixed(2)}
                    </StatNumber>
                  </Stat>
                </HStack>
                <Button
                  width="100%"
                  _hover={{
                    bgColor: 'complimentary.1000',
                  }}
                  onClick={handleLiquidUnstake}
                  isDisabled={Number(tokenAmount) === 0 || !address || isSigning || Number(qBalance?.balance.amount) === 0}
                >
                  {isSigning ? (
                    <Spinner thickness="2px" speed="0.65s" emptyColor="gray.200" color="complimentary.900" size="sm" />
                  ) : (
                    'Unstake'
                  )}
                </Button>
              </VStack>
            </TabPanel>
          </SlideFade>
        </TabPanels>
      </Tabs>
    </Box>
  );
};
