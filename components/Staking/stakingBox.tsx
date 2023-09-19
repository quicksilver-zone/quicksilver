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
  Spinner,
  Skeleton,
  SkeletonText,
} from '@chakra-ui/react';
import React, { useEffect, useState } from 'react';

import { useInputBox } from '@/hooks/useInputBox';
import { useStakingData } from '@/hooks/useStakingData';

import { MultiModal } from './modals/multiStakeModal';
import StakingProcessModal from './modals/stakingProcessModal';

type StakingBoxProps = {
  selectedOption: {
    name: string;
    value: string;
    logo: string;
    chainName: string;
  };
  isModalOpen: boolean;
  setModalOpen: (isOpen: boolean) => void;
};

export const StakingBox = ({
  selectedOption,
  isModalOpen,
  setModalOpen,
}: StakingBoxProps): JSX.Element => {
  const [tokenAmount, setTokenAmount] = useState<string>('0');
  const { data, isLoading, refetch } = useStakingData(selectedOption.chainName);
  const balance = data?.balance;
  const { renderInputBox, amount, setAmount } = useInputBox(balance);
  const maxStakingAmount = balance ? parseFloat(data.balance) - 0.00005 : 0;
  const maxHalfStakingAmount = maxStakingAmount / 2;

  useEffect(() => {
    refetch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedOption]);

  return (
    <Box
      position="relative"
      backdropFilter="blur(50px)"
      bgColor="rgba(255,255,255,0.1)"
      flex="1"
      borderRadius="10px"
      p={5}
    >
      <Tabs isFitted variant="enclosed">
        <TabList
          mt={'4'}
          mb="1em"
          overflow="hidden"
          borderBottomColor="transparent"
          bg="rgba(255,255,255,0.1)"
          p={2}
          borderRadius="25px"
        >
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
          <TabPanel>
            <VStack spacing={8} align="center">
              <Text fontWeight="light" textAlign="center" color="white">
                Stake your {selectedOption.value.toUpperCase()} tokens in
                exchange for q{selectedOption.value.toUpperCase()} which you can
                deploy around the ecosystem. You can liquid stake half of your
                balance, if you&apos;re going to LP.
              </Text>
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
                  placeholder="amount"
                  value={tokenAmount}
                  onChange={(e) => {
                    const inputValue = parseFloat(e.target.value);
                    if (inputValue <= maxStakingAmount) {
                      setTokenAmount(e.target.value);
                    }
                  }}
                />
                <Flex
                  w="100%"
                  flexDirection="row"
                  py={4}
                  mb={-4}
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <Flex gap={4} flexDirection={'row'}>
                    <Text color="white" fontWeight="light">
                      Tokens available:{' '}
                    </Text>
                    <Text color="complimentary.900" fontWeight="light">
                      {balance}
                      {'0'}&nbsp;
                      {selectedOption.value.toUpperCase()}
                    </Text>
                  </Flex>
                  <HStack spacing={2}>
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
                      onClick={() =>
                        setTokenAmount(maxHalfStakingAmount.toString())
                      }
                      isDisabled={!balance}
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
                      onClick={() =>
                        setTokenAmount(maxStakingAmount.toString())
                      }
                      isDisabled={!balance}
                    >
                      max
                    </Button>
                  </HStack>
                </Flex>
              </Flex>
              <Divider bgColor="complimentary.900" />
              <HStack
                justifyContent="space-between"
                alignItems="left"
                w="100%"
                mt={-8}
              >
                <Stat textAlign="left" color="white">
                  <StatLabel>What you&apos;ll get</StatLabel>
                  <StatNumber>
                    q{selectedOption.value.toUpperCase()}:
                  </StatNumber>
                </Stat>
                <Spacer />{' '}
                {/* This pushes the next Stat component to the right */}
                <Stat py={4} textAlign="right" color="white">
                  <StatNumber textColor="complimentary.900">0</StatNumber>
                </Stat>
              </HStack>
              <Button
                width="100%"
                _hover={{
                  bgColor: '#181818',
                }}
                onClick={() => setModalOpen(true)}
              >
                Liquid stake
              </Button>
              <StakingProcessModal
                tokenAmount={tokenAmount}
                isOpen={isModalOpen}
                onClose={() => setModalOpen(false)}
                selectedOption={selectedOption}
              />
            </VStack>
          </TabPanel>
          <TabPanel>
            <VStack spacing={8} align="center">
              <Text fontWeight="light" textAlign="center" color="white">
                Unstake your q{selectedOption.value.toUpperCase()} tokens in
                exchange for {selectedOption.value.toUpperCase()}.
              </Text>
              <Flex flexDirection="column" w="100%">
                <Stat py={4} textAlign="left" color="white">
                  <StatLabel>Amount tounstake:</StatLabel>
                  <StatNumber>
                    q{selectedOption.value.toUpperCase()}{' '}
                  </StatNumber>
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
                  placeholder="amount"
                  value={tokenAmount}
                  onChange={(e) => setTokenAmount(e.target.value)}
                />
                <Flex
                  w="100%"
                  flexDirection="row"
                  py={4}
                  mb={-4}
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <Skeleton isLoaded={!isLoading}>
                    <SkeletonText>
                      <Text color="white" fontWeight="light">
                        Tokens available: 0 q
                        {selectedOption.value.toUpperCase()}
                      </Text>
                    </SkeletonText>
                  </Skeleton>

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
                  >
                    max
                  </Button>
                </Flex>
              </Flex>
              <Divider bgColor="complimentary.900" />
              <HStack
                justifyContent="space-between"
                alignItems="left"
                w="100%"
                mt={-8}
              >
                <Stat textAlign="left" color="white">
                  <StatLabel>What you&apos;ll get</StatLabel>
                  <StatNumber>{selectedOption.value.toUpperCase()}:</StatNumber>
                </Stat>
                <Spacer />{' '}
                {/* This pushes the next Stat component to the right */}
                <Stat py={4} textAlign="right" color="white">
                  <StatNumber textColor="complimentary.900">0</StatNumber>
                </Stat>
              </HStack>
              <Button
                width="100%"
                _hover={{
                  bgColor: 'complimentary.1000',
                }}
              >
                Liquid Stake
              </Button>
            </VStack>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Box>
  );
};
