import { Container, Text, SlideFade, Image, Box, VStack, SkeletonCircle } from '@chakra-ui/react';
import {
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
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
  Spinner,
  Tooltip,
  Link,
} from '@chakra-ui/react';
import { Coin, StdFee, coins } from '@cosmjs/amino';
import { useChain } from '@cosmos-kit/react';
import { bech32 } from 'bech32';
import { assets, chains } from 'chain-registry';
import Head from 'next/head';
import { useRouter } from 'next/router';
import { cosmos, quicksilver } from 'quicksilverjs';
import React, { useEffect, useState } from 'react';

import { Chain, chains as configChains, env } from '@/config';
import { useTx } from '@/hooks';
import { useBalanceQuery, useQBalanceQuery, useValidatorLogos, useValidatorsQuery, useZoneQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';

function StakingWithValidator() {
  const router = useRouter();
  const { chainId, valoperAddress } = router.query;
  const networks: Map<string, Chain> = configChains.get(env) ?? new Map();
  const selectedNetwork: Chain | undefined = Array.from(networks.values())?.find((network: Chain) => network.chain_id === chainId);
  return { selectedNetwork, valoperAddress };
}

export default function Home() {
  const pageInfo = StakingWithValidator();
  const selectedNetwork = pageInfo.selectedNetwork;
  const valoperAddress = pageInfo.valoperAddress;

  const chainIdToChainName = (chainId: string) => {
    switch (chainId) {
      case 'cosmoshub-4':
        return 'cosmos';
      case 'stargaze-1':
        return 'stargaze';
      case 'osmosis-1':
        return 'osmosis';
      case 'regen-1':
        return 'regen';
      case 'sommelier-1':
        return 'sommelier';
      case 'juno-1':
        return 'juno';
      case 'saga-1':
        return 'saga';
      case 'dydx-1':
        return 'dydx';
      default:
        return '';
    }
  };
  let chainId = selectedNetwork?.chain_id;

  let chainName = chainIdToChainName(chainId ?? '');

  const isValidValoperAddress = () => {
    if (typeof valoperAddress === 'string') {
      try {
        bech32.decode(valoperAddress);
        return true;
      } catch {
        return false;
      }
    }
    return false;
  };

  const validValoperAddress = typeof valoperAddress === 'string' ? valoperAddress : undefined;

  return (
    <>
      <SlideFade offsetY={'200px'} in={true} style={{ width: '100%' }}>
        <Container
          mt={-12}
          flexDir={'column'}
          top={20}
          zIndex={2}
          position="relative"
          justifyContent="center"
          alignItems="center"
          maxW="5xl"
        >
          <Head>
            <title>{chainName.charAt(0).toUpperCase() + chainName.slice(1)} Liquid Staking</title>
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <meta name="description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
            <meta name="keywords" content={`liquid staking on ${chainName}, ${chainName}, Quicksilver`} />
            <meta name="author" content="Quicksilver Zone" />
            <link rel="icon" href="/img/favicon-main.png" />

            <meta property="og:title" content="Staking - Quicksilver Zone" />
            <meta property="og:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
            <meta property="og:url" content={`https://app.quicksilver.zone/${selectedNetwork}/${valoperAddress}`} />
            <meta property="og:image" content="https://app.quicksilver.zone/img/banner.png" />
            <meta property="og:type" content="website" />
            <meta property="og:site_name" content="Quicksilver Protocol" />

            <meta name="twitter:card" content="summary_large_image" />
            <meta name="twitter:title" content="Staking - Quicksilver Zone" />
            <meta name="twitter:description" content="STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED" />
            <meta name="twitter:image" content="https://app.quicksilver.zone/img/banner.png" />
            <meta name="twitter:site" content="@quicksilverzone" />

            <script type="application/ld+json">
              {JSON.stringify({
                '@context': 'https://schema.org',
                '@type': 'WebPage',
                name: 'Staking - Quicksilver Zone',
                description: 'STAKING SIMPLIFIED | LIQUIDITY AMPLIFIED',
                url: `https://app.quicksilver.zone/${selectedNetwork}/${valoperAddress}`,
                image: 'https://app.quicksilver.zone/img/banner.png',
                publisher: {
                  '@type': 'Organization',
                  name: 'Quicksilver Protocol',
                  logo: {
                    '@type': 'ImageObject',
                    url: 'https://app.quicksilver.zone/img/logo.png',
                  },
                },
              })}
            </script>
          </Head>
          {selectedNetwork && validValoperAddress && isValidValoperAddress() ? (
            <StakingBox valoperAddress={validValoperAddress} selectedOption={selectedNetwork} />
          ) : (
            <Box
              maxW="md"
              mx="auto"
              position="relative"
              backdropFilter="blur(50px)"
              bgColor="rgba(255,255,255,0.1)"
              flex="1"
              borderRadius="10px"
              p={5}
            >
              <Stat>
                <StatLabel textAlign={'center'} color="white" fontSize="lg" pb={5}>
                  You are attempting to delegate directly to a validator but there is an issue with the url.
                </StatLabel>
                <StatLabel textAlign={'center'} color="red" fontSize="lg">
                  Error:
                </StatLabel>
                <StatNumber textAlign={'center'} color="white" fontSize="lg" pb={5}>
                  {selectedNetwork
                    ? 'Validator address not found. Please check the address and try again.'
                    : 'The specified network was not found. Please check the URL and try again.'}
                </StatNumber>
              </Stat>
              <Text textAlign={'center'}>
                If you believe there is an issue, please contact us on{' '}
                <Link color="complimentary.900" href="https://discord.gg/4QXEJQcv" isExternal>
                  Discord
                </Link>
              </Text>
            </Box>
          )}
        </Container>
      </SlideFade>
    </>
  );
}

type StakingBoxProps = {
  selectedOption: Chain;
  valoperAddress: string;
};

export const StakingBox = ({ selectedOption, valoperAddress }: StakingBoxProps) => {
  const [activeTabIndex, setActiveTabIndex] = useState(0);
  const [tokenAmount, setTokenAmount] = useState<string>('0');

  const { address } = useChain(selectedOption.chain_name);

  const { address: qAddress } = useChain('quicksilver');
  const { balance, isLoading } = useBalanceQuery(selectedOption.chain_name, address ?? '');

  const { balance: qBalance } = useQBalanceQuery('quicksilver', qAddress ?? '', selectedOption.major_denom.toLowerCase());

  const qAssets = qBalance?.balance.amount || '';

  const baseBalance = shiftDigits(balance?.balance?.amount || '0', -selectedOption.exponent);

  const { data: zone, isLoading: isZoneLoading, isError: isZoneError } = useZoneQuery(selectedOption.chain_id);

  useEffect(() => {
    setTokenAmount('0');
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedOption.chain_name]);

  const truncateToThreeDecimals = (num: number) => {
    return Math.trunc(num * 1000) / 1000;
  };

  const truncatedBalance = truncateToThreeDecimals(Number(baseBalance));

  const maxStakingAmount = truncateToThreeDecimals(truncatedBalance ? truncatedBalance - 0.005 : 0);

  const maxHalfStakingAmount = maxStakingAmount / 2;

  const [inputError, setInputError] = useState(false);

  const exponent = qBalance?.balance.denom === 'aqdydx' ? -18 : -6;
  const qAssetsExponent = shiftDigits(qAssets, exponent);
  const qAssetsDisplay = qAssetsExponent.includes('.') ? qAssetsExponent.substring(0, qAssetsExponent.indexOf('.') + 3) : qAssetsExponent;

  const maxUnstakingAmount = truncateToThreeDecimals(Number(qAssetsDisplay));
  const halfUnstakingAmount = maxUnstakingAmount / 2;

  const [isSigning, setIsSigning] = useState<boolean>(false);

  const env = process.env.NEXT_PUBLIC_CHAIN_ENV;
  const quicksilverChainName = env === 'testnet' ? 'quicksilvertestnet' : 'quicksilver';

  const { requestRedemption } = quicksilver.interchainstaking.v1.MessageComposer.withTypeUrl;
  const numericAmount = Number(tokenAmount);
  const smallestUnitAmount = numericAmount * Math.pow(10, 6);
  const value: Coin = { amount: smallestUnitAmount.toFixed(0), denom: zone?.localDenom ?? '' };

  // Create the message only executes if the unstake button is clickable
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
      await tx([msgRequestRedemption], {
        fee,
        onSuccess: () => {},
      });
    } catch (error) {
      console.error('Transaction failed', error);
    } finally {
      setIsSigning(false);
    }
  };

  const handleTabsChange = (index: number) => {
    setActiveTabIndex(index);
    setTokenAmount('');
  };

  const { validatorsData, isLoading: validatorsDataLoading, isError: validatorsDataError } = useValidatorsQuery(selectedOption.chain_name);

  const moniker = validatorsData?.find((validator) => validator.address === valoperAddress)?.name ?? '';

  const intents = {
    address: valoperAddress,
    intent: 1,
  };

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
    let wordsUint8Array = new Uint8Array(bech32.fromWords(words));
    let weightByte = valToByte(weight);
    return Buffer.concat([Buffer.from([weightByte]), wordsUint8Array]);
  };

  let memoBuffer = Buffer.alloc(0);

  if (intents) {
    memoBuffer = Buffer.concat([memoBuffer, addValidator(intents.address, intents.intent)]);

    memoBuffer = Buffer.concat([Buffer.from([0x02, memoBuffer.length]), memoBuffer]);
  }

  let memo = memoBuffer.length > 0 && valoperAddress ? memoBuffer.toString('base64') : '';

  const { send } = cosmos.bank.v1beta1.MessageComposer.withTypeUrl;
  // Create the message only executes if the liquid stake button is clickable
  const msgSend = send({
    fromAddress: address ?? '',
    toAddress: zone?.depositAddress?.address ?? '',
    amount: coins(smallestUnitAmount.toFixed(0), zone?.baseDenom ?? ''),
  });

  const mainTokens = assets.find(({ chain_name }) => chain_name === selectedOption.chain_name);
  const fees = chains.find(({ chain_name }) => chain_name === selectedOption.chain_name)?.fees?.fee_tokens;
  const mainDenom = mainTokens?.assets[0].base ?? '';
  const fixedMinGasPrice = fees?.find(({ denom }) => denom === mainDenom)?.average_gas_price ?? '';
  const feeAmount = shiftDigits(fixedMinGasPrice, 6);

  const stakeFee: StdFee = {
    amount: [
      {
        denom: mainDenom,
        amount: feeAmount.toString(),
      },
    ],
    gas: '500000',
  };

  const { tx: sendTx } = useTx(selectedOption.chain_name);

  const handleLiquidStake = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);

    try {
      await sendTx([msgSend], {
        memo,
        fee: stakeFee,
        onSuccess: () => {},
      });
    } catch (error) {
      console.error('Transaction failed', error);
    } finally {
      setIsSigning(false);
    }
  };

  const isWalletConnected = !!address;
  const isLiquidStakeDisabled = Number(tokenAmount) === 0 || !isWalletConnected || Number(tokenAmount) < 0.1;

  let liquidStakeTooltip = '';
  if (!isWalletConnected) {
    liquidStakeTooltip = 'Connect your wallet to stake';
  } else if (Number(tokenAmount) < 0.1) {
    liquidStakeTooltip = 'Minimum amount to stake is 0.1';
  }
  const { data: logos, isLoading: isFetchingLogos } = useValidatorLogos(selectedOption?.chain_name ?? '', validatorsData || []);
  const validatorLogo = logos ? logos[valoperAddress] : undefined;
  return (
    <Box
      maxW="500px"
      mx="auto"
      position="relative"
      backdropFilter="blur(50px)"
      bgColor="rgba(255,255,255,0.1)"
      flex="1"
      borderRadius="10px"
      p={5}
    >
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
                  Stake your {selectedOption.major_denom.toUpperCase()} tokens in exchange for q{selectedOption.major_denom.toUpperCase()}
                </Text>
                <Stat py={4} textAlign="left" color="white">
                  <StatLabel fontSize={'lg'} py={1}>
                    Validator:
                  </StatLabel>
                  {moniker ? (
                    <StatNumber>
                      {!validatorLogo && (
                        <SkeletonCircle
                          boxSize="26px"
                          objectFit="cover"
                          marginRight="8px"
                          display="inline-block"
                          verticalAlign="middle"
                          startColor="complimentary.900"
                          endColor="complimentary.100"
                        />
                      )}
                      {validatorLogo && (
                        <Image
                          borderRadius={'full'}
                          src={validatorLogo}
                          alt={moniker}
                          boxSize="35px"
                          objectFit="cover"
                          marginRight="8px"
                          display="inline-block"
                          verticalAlign="middle"
                          mt={-2}
                        />
                      )}
                      {moniker}
                    </StatNumber>
                  ) : (
                    <SkeletonText mt="4" h="20px" startColor="complimentary.900" endColor="complimentary.500" noOfLines={1} w="150px" />
                  )}
                </Stat>

                <>
                  <Flex flexDirection="column" w="100%">
                    <Stat py={4} textAlign="left" color="white">
                      <StatLabel>Amount to stake:</StatLabel>
                      <StatNumber>{selectedOption.major_denom.toUpperCase()} </StatNumber>
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
                                {address ? (
                                  balance?.balance?.amount && Number(balance?.balance?.amount) !== 0 ? (
                                    `${truncatedBalance} ${selectedOption.major_denom.toUpperCase()}`
                                  ) : (
                                    <Link href={`https://app.osmosis.zone/?from=USDC&to=${selectedOption.major_denom.toUpperCase()}`} isExternal>
                                      Get {selectedOption.major_denom.toUpperCase()} tokens here
                                    </Link>
                                  )
                                ) : (
                                  '0'
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
                      <StatNumber>q{selectedOption.major_denom.toUpperCase()}:</StatNumber>
                    </Stat>
                    <Spacer /> {/* This pushes the next Stat component to the right */}
                    <Stat py={4} textAlign="right" color="white">
                      <StatNumber textColor="complimentary.900">
                        {!isZoneLoading ? (
                          (Number(tokenAmount) * Number(zone?.redemptionRate || 1)).toFixed(2)
                        ) : (
                          <Spinner thickness="2px" speed="0.65s" emptyColor="gray.200" color="complimentary.900" size="sm" />
                        )}
                      </StatNumber>
                    </Stat>
                  </HStack>
                  <Tooltip hasArrow label={liquidStakeTooltip} isDisabled={!isLiquidStakeDisabled}>
                    <Button
                      width="100%"
                      _hover={{
                        bgColor: 'complimentary.1000',
                      }}
                      isDisabled={Number(tokenAmount) === 0 || !address || Number(tokenAmount) < 0.1}
                      onClick={handleLiquidStake}
                    >
                      {isSigning ? (
                        <Spinner thickness="2px" speed="0.65s" emptyColor="gray.200" color="complimentary.900" size="sm" />
                      ) : (
                        'Liquid Stake'
                      )}
                    </Button>
                  </Tooltip>
                </>
              </VStack>
            </TabPanel>
          </SlideFade>
          <SlideFade offsetY="200px" in={activeTabIndex === 1}>
            <TabPanel>
              <VStack spacing={8} align="center">
                <Text fontWeight="light" textAlign="center" color="white">
                  Unstake your q{selectedOption.major_denom.toUpperCase()} tokens in exchange for {selectedOption.major_denom.toUpperCase()}.
                </Text>
                <Flex flexDirection="column" w="100%">
                  <Stat py={4} textAlign="left" color="white">
                    <StatLabel>Amount to unstake:</StatLabel>
                    <StatNumber>q{selectedOption.major_denom.toUpperCase()} </StatNumber>
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
                                ? `${qAssetsDisplay} ${selectedOption.major_denom.toUpperCase()}`
                                : `No q${selectedOption.major_denom.toUpperCase()}`
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
                    <StatNumber>{selectedOption.major_denom.toUpperCase()}:</StatNumber>
                  </Stat>
                  <Spacer /> {/* This pushes the next Stat component to the right */}
                  <Stat py={4} textAlign="right" color="white">
                    <StatNumber textColor="complimentary.900">
                      {!isZoneLoading ? (
                        (Number(tokenAmount) * Number(zone?.redemptionRate || 1)).toFixed(2)
                      ) : (
                        <Spinner thickness="2px" speed="0.65s" emptyColor="gray.200" color="complimentary.900" size="sm" />
                      )}
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
