import { InfoOutlineIcon } from '@chakra-ui/icons';
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
  SlideFade,
  Spinner,
  Switch,
  Tooltip,
  Image,
  SkeletonCircle,
  Link,
} from '@chakra-ui/react';
import { Coin, StdFee } from '@cosmjs/amino';
import { useChain } from '@cosmos-kit/react';
import { quicksilver } from 'quicksilverjs';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { FaStar } from 'react-icons/fa';

import RevertSharesProcessModal from './modals/revertSharesProcessModal';
import StakingProcessModal from './modals/stakingProcessModal';
import TransferProcessModal from './modals/transferProcessModal';

import { useTx } from '@/hooks';
import {
  useAllBalancesQuery,
  useBalanceQuery,
  useNativeStakeQuery,
  useQBalanceQuery,
  useValidatorLogos,
  useValidatorsQuery,
  useZoneQuery,
} from '@/hooks/useQueries';
import { getExponent, shiftDigits } from '@/utils';

type StakingBoxProps = {
  selectedOption: {
    name: string;
    value: string;
    logo: string;
    chainName: string;
    chainId: string;
  };
  isStakingModalOpen: boolean;
  setStakingModalOpen: (isOpen: boolean) => void;
  isTransferModalOpen: boolean;
  setTransferModalOpen: (isOpen: boolean) => void;
  isRevertSharesModalOpen: boolean;
  setRevertSharesModalOpen: (isOpen: boolean) => void;
  setBalance: (balance: string) => void;
  setQBalance: (qBalance: string) => void;
};

export const StakingBox = ({
  selectedOption,
  isStakingModalOpen,
  setStakingModalOpen,
  isTransferModalOpen,
  setTransferModalOpen,
  isRevertSharesModalOpen,
  setRevertSharesModalOpen,
  setBalance,
  setQBalance,
}: StakingBoxProps) => {
  const [activeTabIndex, setActiveTabIndex] = useState(0);
  const [tokenAmount, setTokenAmount] = useState<string>('0');

  const openStakingModal = () => setStakingModalOpen(true);
  const closeStakingModal = () => setStakingModalOpen(false);

  const openTransferModal = () => setTransferModalOpen(true);
  const closeTransferModal = () => setTransferModalOpen(false);

  const openRevertSharesModal = () => setRevertSharesModalOpen(true);
  const closeRevertSharesModal = () => setRevertSharesModalOpen(false);

  const { address } = useChain(selectedOption.chainName);

  const { address: qAddress } = useChain('quicksilver');
  const exp = getExponent(selectedOption.chainName);

  const { balance, isLoading } = useBalanceQuery(selectedOption.chainName, address ?? '');

  const { balance: allBalances } = useAllBalancesQuery(selectedOption.chainName, address ?? '');

  const { balance: qBalance } = useQBalanceQuery('quicksilver', qAddress ?? '', selectedOption.value.toLowerCase());

  const qAssets = qBalance?.balance.amount || '';

  const baseBalance = shiftDigits(balance?.balance?.amount || '0', -exp);

  const { data: zone, isLoading: isZoneLoading } = useZoneQuery(selectedOption.chainId);

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

  const env = process.env.NEXT_PUBLIC_CHAIN_ENV;
  const quicksilverChainName = env === 'testnet' ? 'quicksilvertestnet' : 'quicksilver';

  const { requestRedemption } = quicksilver.interchainstaking.v1.MessageComposer.withTypeUrl;
  const numericAmount = Number(tokenAmount);
  const smallestUnitAmount = numericAmount * Math.pow(10, 6);
  const value: Coin = { amount: smallestUnitAmount.toFixed(0), denom: zone?.local_denom ?? '' };
  const msgRequestRedemption = requestRedemption({
    value: value,
    from_address: qAddress ?? '',
    destination_address: address ?? '',
  });

  const fee: StdFee = {
    amount: [
      {
        denom: 'uqck',
        amount: '50',
      },
    ],
    gas: '500000',
  };

  const { tx } = useTx(quicksilverChainName);

  const handleLiquidUnstake = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);
    try {
      await tx([msgRequestRedemption], {
        fee,
      });
    } catch (error) {
      console.error('Transaction failed', error);
    } finally {
      setIsSigning(false);
    }
  };

  // import { useToaster, ToastType, type CustomToast } from '@/hooks/useToaster';
  // import useToast from chakra-ui
  // You can use this Toast handler and the below message to show there is an issue with unbonding
  // const toaster = useToaster();

  const handleTabsChange = (index: number) => {
    setActiveTabIndex(index);
    setTokenAmount('');
    // You can use this Toast Msg to show there is an issue with unbonding
    // if (index === 1) {
    //   toaster.toast({
    //     type: ToastType.Error,
    //     title: 'Issues with unbonding',
    //     message: 'Unbondings can be submitted but are currently not being processed and will be queued until the issue is resolved.',
    //   });
    // }
  };

  const { delegations, delegationsIsError, delegationsIsLoading } = useNativeStakeQuery(selectedOption.chainName, address ?? '');
  const delegationsResponse = delegations?.delegation_responses;
  const nativeStakedAmount = delegationsResponse?.reduce((acc: number, delegationResponse: { balance: { amount: any } }) => {
    const amount = Number(delegationResponse?.balance?.amount) || 0;
    return acc + amount;
  }, 0);

  const [useNativeStake, setUseNativeStake] = useState(false);

  const hasTokenizedShares = (balances: any[]) => {
    return balances.some((balance: { denom: string | string[] }) => balance.denom.includes('valoper'));
  };

  const hasTokenized = useMemo(() => hasTokenizedShares(allBalances?.balances || []), [allBalances]);

  const handleSwitchChange = (event: { target: { checked: boolean | ((prevState: boolean) => boolean) } }) => {
    setUseNativeStake(event.target.checked);
  };

  useEffect(() => {
    setUseNativeStake(false);
  }, [selectedOption]);

  const { validatorsData } = useValidatorsQuery(selectedOption.chainName);

  const { data: logos } = useValidatorLogos(selectedOption.chainName, validatorsData || []);
  const [selectedValidator, setSelectedValidator] = useState<string | null>(null);

  const [isBottomVisible, setIsBottomVisible] = useState(true);

  const handleScroll = useCallback((event: React.UIEvent<HTMLDivElement>) => {
    const target = event.target as HTMLDivElement;
    const isBottom = target.scrollHeight - target.scrollTop === target.clientHeight;
    setIsBottomVisible(!isBottom);
  }, []);

  interface SelectedValidator {
    operatorAddress: string;
    moniker: string;
    tokenAmount: string;
    isTokenized: boolean;
    denom: string;
  }

  const [selectedValidatorData, setSelectedValidatorData] = useState<SelectedValidator>({
    operatorAddress: '',
    moniker: '',
    tokenAmount: '',
    isTokenized: false,
    denom: '',
  });

  const isWalletConnected = !!address;
  const isLiquidStakeDisabled = Number(tokenAmount) === 0 || !isWalletConnected || Number(tokenAmount) < 0.1;

  let liquidStakeTooltip = '';
  if (!isWalletConnected) {
    liquidStakeTooltip = 'Connect your wallet to stake';
  } else if (Number(tokenAmount) < 0.1) {
    liquidStakeTooltip = 'Minimum amount to stake is 0.1';
  }

  const safeDelegationsResponse = delegationsResponse || [];
  const safeAllBalances = allBalances?.balances || [];

  // Combine delegationsResponse with valoper entries from allBalances
  const combinedDelegations = safeDelegationsResponse.concat(
    safeAllBalances
      .filter((balance) => balance.denom.includes('valoper'))
      .map((balance) => {
        const [validatorAddress, uniqueId] = balance.denom.split('/');
        return {
          delegation: {
            delegator_address: '',
            validator_address: validatorAddress,
            unique_id: uniqueId,
            shares: '',
          },
          balance: {
            amount: balance.amount,
            denom: balance.denom,
          },
          isTokenized: true,
          denom: balance.denom,
        };
      }),
  );

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
          {/* Staking TabPanel */}
          <SlideFade offsetY="-80px" in={activeTabIndex === 0}>
            <TabPanel>
              <VStack spacing={8} align="center">
                <Text fontWeight="light" textAlign="center" color="white">
                  Stake your {selectedOption.value.toUpperCase()} tokens in exchange for q{selectedOption.value.toUpperCase()} which you can
                  deploy around the ecosystem. You can liquid stake half of your balance, if you&apos;re going to LP.
                </Text>
                {selectedOption.name === 'Cosmos Hub' && (
                  <Flex textAlign={'left'} justifyContent={'flex-start'}>
                    {((nativeStakedAmount ?? 0) > 0 || hasTokenized) && (
                      <Tooltip
                        label={
                          !address
                            ? 'Please connect your wallet to enable this option.'
                            : !nativeStakedAmount && !hasTokenized
                              ? "You don't have any native staked tokens or tokenized shares."
                              : nativeStakedAmount ?? 0 > 0
                                ? `You currently have ${shiftDigits(nativeStakedAmount ?? '', -6)} ${
                                    selectedOption.value
                                  } natively staked to ${delegationsResponse?.length} validators.`
                                : hasTokenized
                                  ? 'You have tokenized shares available for transfer.'
                                  : ''
                        }
                      >
                        <HStack>
                          <Text
                            fontWeight={!nativeStakedAmount && !hasTokenized ? 'hairline' : 'normal'}
                            textAlign="center"
                            color={!nativeStakedAmount && !hasTokenized ? 'whiteAlpha.800' : 'white'}
                          >
                            Use staked&nbsp;
                            <span style={{ color: '#FF8000' }}>{selectedOption.value}</span>?
                          </Text>
                          {delegationsIsLoading && <SkeletonCircle size="4" startColor="complimentary.900" endColor="complimentary.400" />}
                          {!delegationsIsLoading && !delegationsIsError && (
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
                              isDisabled={(!nativeStakedAmount && !hasTokenized) || !logos}
                              isChecked={useNativeStake}
                              onChange={handleSwitchChange}
                              id="use-natively-staked"
                              colorScheme="orange"
                            />
                          )}
                          <InfoOutlineIcon color={!nativeStakedAmount && !hasTokenized ? 'complimentary.1100' : 'complimentary.900'} />
                        </HStack>
                      </Tooltip>
                    )}
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
                          const validNumberPattern = /^\d*\.?\d*$/;
                          if (validNumberPattern.test(e.target.value) || e.target.value === '') {
                            setTokenAmount(e.target.value);
                          }
                        }}
                        onBlur={() => {
                          // Check if the input is a lone period or incomplete number format
                          if (tokenAmount === '.') {
                            setInputError(true);
                            setTokenAmount('');
                          } else {
                            let inputValue = parseFloat(tokenAmount);
                            if (isNaN(inputValue) || inputValue <= 0) {
                              setInputError(true);
                              setTokenAmount('');
                            } else if (inputValue > maxStakingAmount) {
                              setInputError(false);
                              setTokenAmount(maxStakingAmount.toString());
                            } else {
                              setInputError(false);
                              setTokenAmount(inputValue.toString());
                            }
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
                                  {balance?.balance?.amount && Number(balance.balance.amount) > 0 ? (
                                    `${truncatedBalance} ${selectedOption.value.toUpperCase()}`
                                  ) : (
                                    <Link href={`https://app.osmosis.zone/?from=USDC&to=${selectedOption.value.toUpperCase()}`} isExternal>
                                      Get {selectedOption.value.toUpperCase()} tokens here
                                    </Link>
                                  )}
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
                    <HStack pt={2} justifyContent="space-between" alignItems="left" w="100%" mt={-8}>
                      <Stat textAlign="left" color="white">
                        <StatLabel>What you&apos;ll get</StatLabel>
                        <StatNumber>q{selectedOption.value.toUpperCase()}:</StatNumber>
                      </Stat>
                      <Spacer /> {/* This pushes the next Stat component to the right */}
                      <Stat py={4} textAlign="right" color="white">
                        <StatNumber textColor="complimentary.900">
                          {!isZoneLoading ? (
                            (Number(tokenAmount) * Number(zone?.redemption_rate || 1)).toFixed(2)
                          ) : (
                            <Spinner thickness="2px" speed="0.65s" emptyColor="gray.200" color="complimentary.900" size="sm" />
                          )}
                        </StatNumber>
                      </Stat>
                    </HStack>
                    <Tooltip hasArrow label={liquidStakeTooltip} isDisabled={!isLiquidStakeDisabled}>
                      <Button
                        width="100%"
                        _active={{
                          transform: 'scale(0.95)',
                          color: 'complimentary.800',
                        }}
                        _hover={{
                          bgColor: 'rgba(255,128,0, 0.25)',
                          color: 'complimentary.300',
                        }}
                        onClick={openStakingModal}
                        isDisabled={Number(tokenAmount) === 0 || !address || Number(tokenAmount) < 0.1}
                      >
                        Liquid stake
                      </Button>
                    </Tooltip>
                    <StakingProcessModal
                      tokenAmount={tokenAmount}
                      isOpen={isStakingModalOpen}
                      onClose={closeStakingModal}
                      selectedOption={selectedOption}
                      address={address ?? ''}
                    />
                  </>
                )}
                {useNativeStake && (
                  <Flex flexDirection="column" w="100%">
                    <VStack spacing={8} align="center">
                      <Box position="relative" mb={8}>
                        <Box className="custom-scrollbar" maxH="290px" overflowY="scroll" w="fit-content" onScroll={handleScroll}>
                          {/* Combine delegationsResponse with valoper entries from allBalances */}
                          {combinedDelegations.map(
                            // @ts-ignore
                            (
                              delegation: {
                                delegation: { validator_address: string | number; unique_id: any };
                                balance: { amount: string | number };
                                isTokenized: any;
                                denom: any;
                              },
                              index: any,
                            ) => {
                              const validator = validatorsData?.find((v) => v.address === delegation.delegation.validator_address);
                              const uniqueKey = `${delegation.delegation.validator_address}-${delegation.delegation.unique_id}`;
                              const isSelected = uniqueKey === selectedValidator;
                              const validatorLogo = logos[delegation.delegation.validator_address];

                              return (
                                <Box
                                  borderRadius={'md'}
                                  as="button"
                                  w="full"
                                  onClick={() => {
                                    setSelectedValidator(uniqueKey);
                                    setSelectedValidatorData({
                                      operatorAddress: delegation.delegation.validator_address.toString(),
                                      moniker: validator?.name ?? '',
                                      tokenAmount: delegation.balance.amount.toString(),
                                      isTokenized: delegation.isTokenized,
                                      denom: delegation.denom,
                                    });
                                  }}
                                  _hover={{ bg: 'rgba(255, 128, 0, 0.25)' }}
                                  bg={isSelected ? 'rgba(255, 128, 0, 0.25)' : 'transparent'}
                                  key={uniqueKey}
                                  mb={2}
                                >
                                  <Flex py={2} align="center" justify="space-between">
                                    <HStack align="center" ml={4}>
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
                                      <VStack align="start">
                                        <HStack>
                                          <Text fontSize="md">{validator?.name ?? 'Validator'}</Text>
                                          {delegation.isTokenized && (
                                            <Tooltip
                                              label="This share is tokenized and can be transferred to quicksilver."
                                              aria-label="Tokenized Share"
                                            >
                                              <Box>
                                                <FaStar color="#FF8000" />
                                              </Box>
                                            </Tooltip>
                                          )}
                                        </HStack>
                                        <Text color={'complimentary.900'} fontSize="md">
                                          {shiftDigits(delegation.balance.amount, -6)} {selectedOption.value}
                                        </Text>
                                      </VStack>
                                    </HStack>
                                    {isSelected && delegation.isTokenized && (
                                      <Button
                                        onClick={openRevertSharesModal}
                                        _active={{
                                          transform: 'scale(0.95)',
                                          color: 'complimentary.800',
                                        }}
                                        _hover={{
                                          bgColor: 'rgba(255,128,0, 0.25)',
                                          color: 'complimentary.300',
                                        }}
                                        color="white"
                                        size="sm"
                                        variant="outline"
                                        mb={6}
                                        mr={2}
                                      >
                                        Revert
                                      </Button>
                                    )}
                                  </Flex>
                                </Box>
                              );
                            },
                          )}
                        </Box>
                        {isBottomVisible && (
                          <Box
                            position="absolute"
                            bottom="0"
                            left="0"
                            right="0"
                            height="70px"
                            bgGradient="linear(to top, #1A1A1A, transparent)"
                            zIndex="1"
                          />
                        )}
                      </Box>
                    </VStack>
                    <Button
                      width="100%"
                      _hover={{
                        bgColor: 'complimentary.1000',
                      }}
                      onClick={openTransferModal}
                      isDisabled={!selectedValidator || !address}
                    >
                      Transfer Existing Delegation
                    </Button>
                    <TransferProcessModal
                      address={address ?? ''}
                      selectedValidator={selectedValidatorData}
                      isOpen={isTransferModalOpen}
                      onClose={closeTransferModal}
                      selectedOption={selectedOption}
                      isTokenized={selectedValidatorData.isTokenized}
                      denom={selectedValidatorData.denom}
                    />
                    <RevertSharesProcessModal
                      address={address ?? ''}
                      selectedValidator={selectedValidatorData}
                      isOpen={isRevertSharesModalOpen}
                      onClose={closeRevertSharesModal}
                      selectedOption={selectedOption}
                      isTokenized={selectedValidatorData.isTokenized}
                      denom={selectedValidatorData.denom}
                    />
                  </Flex>
                )}
              </VStack>
            </TabPanel>
          </SlideFade>
          {/* Unstake TabPanel */}
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
                                ? `${qAssetsDisplay} q${selectedOption.value.toUpperCase()}`
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
                <HStack pt={2} justifyContent="space-between" alignItems="left" w="100%" mt={-8}>
                  <Stat textAlign="left" color="white">
                    <StatLabel>What you&apos;ll get</StatLabel>
                    <StatNumber>{selectedOption.value.toUpperCase()}:</StatNumber>
                  </Stat>
                  <Spacer /> {/* This pushes the next Stat component to the right */}
                  <Stat py={4} textAlign="right" color="white">
                    <StatNumber textColor="complimentary.900">
                      {!isZoneLoading ? (
                        (Number(tokenAmount) * Number(zone?.redemption_rate || 1)).toFixed(2)
                      ) : (
                        <Spinner thickness="2px" speed="0.65s" emptyColor="gray.200" color="complimentary.900" size="sm" />
                      )}
                    </StatNumber>
                  </Stat>
                </HStack>
                <Button
                  width="100%"
                  _active={{
                    transform: 'scale(0.95)',
                    color: 'complimentary.800',
                  }}
                  _hover={{
                    bgColor: 'rgba(255,128,0, 0.25)',
                    color: 'complimentary.300',
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
