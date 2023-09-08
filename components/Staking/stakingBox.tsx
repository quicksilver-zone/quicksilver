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
} from '@chakra-ui/react';
import React, { useEffect } from 'react';

import { useValidatorData } from '@/hooks';
import { useStakingData } from '@/hooks/useStakingData';
import { type ExtendedValidator as Validator } from '@/utils';

import { MultiModal } from './modals/multiStakeModal';

type StakingBoxProps = {
  selectedOption: string;
  isModalOpen: boolean;
  setModalOpen: (isOpen: boolean) => void;
  selectedChainName: string;
};

export const StakingBox = ({
  selectedOption,
  isModalOpen,
  setModalOpen,
  selectedChainName,
}: StakingBoxProps): JSX.Element => {
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
              borderBottomColor:
                'complimentary.900',
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
              borderBottomColor:
                'complimentary.900',
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
              <Text
                textAlign="center"
                color="white"
              >
                Stake your{' '}
                {selectedOption.toUpperCase()}{' '}
                tokens in exchange for q
                {selectedOption.toUpperCase()}{' '}
                which you can deploy around the
                ecosystem. You can liquid stake
                half of your balance, if
                you&apos;re going to LP.
              </Text>
              <Flex
                flexDirection="column"
                w="100%"
              >
                <Stat
                  py={4}
                  textAlign="left"
                  color="white"
                >
                  <StatLabel>
                    Amount to stake:
                  </StatLabel>
                  <StatNumber>
                    {selectedOption.toUpperCase()}{' '}
                  </StatNumber>
                </Stat>
                <Input
                  _active={{
                    borderColor:
                      'complimentary.900',
                  }}
                  _selected={{
                    borderColor:
                      'complimentary.900',
                  }}
                  _hover={{
                    borderColor:
                      'complimentary.900',
                  }}
                  _focus={{
                    borderColor:
                      'complimentary.900',
                    boxShadow:
                      '0 0 0 3px #FF8000',
                  }}
                  color="complimentary.900"
                  textAlign={'right'}
                  placeholder="amount"
                />
                <Flex
                  w="100%"
                  flexDirection="row"
                  py={4}
                  mb={-4}
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <Text
                    color="white"
                    fontWeight="light"
                  >
                    Tokens available: 0{' '}
                    {selectedOption.toUpperCase()}
                  </Text>
                  <HStack spacing={2}>
                    <Button
                      _hover={{
                        bgColor:
                          'rgba(255,255,255,0.05)',
                        backdropFilter:
                          'blur(10px)',
                      }}
                      _active={{
                        bgColor:
                          'rgba(255,255,255,0.05)',
                        backdropFilter:
                          'blur(10px)',
                      }}
                      color="complimentary.900"
                      variant="ghost"
                      w="60px"
                      h="30px"
                    >
                      half
                    </Button>
                    <Button
                      _hover={{
                        bgColor:
                          'rgba(255,255,255,0.05)',
                        backdropFilter:
                          'blur(10px)',
                      }}
                      _active={{
                        bgColor:
                          'rgba(255,255,255,0.05)',
                        backdropFilter:
                          'blur(10px)',
                      }}
                      color="complimentary.900"
                      variant="ghost"
                      w="60px"
                      h="30px"
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
                <Stat
                  textAlign="left"
                  color="white"
                >
                  <StatLabel>
                    What you&apos;ll get
                  </StatLabel>
                  <StatNumber>
                    q
                    {selectedOption.toUpperCase()}
                    :
                  </StatNumber>
                </Stat>
                <Spacer />{' '}
                {/* This pushes the next Stat component to the right */}
                <Stat
                  py={4}
                  textAlign="right"
                  color="white"
                >
                  <StatNumber textColor="complimentary.900">
                    0
                  </StatNumber>
                </Stat>
              </HStack>
              <Button
                width="100%"
                _hover={{
                  bgColor: '#181818',
                }}
                onClick={() => setModalOpen(true)}
              >
                Validator Selection
              </Button>
              <MultiModal
                isOpen={isModalOpen}
                onClose={() =>
                  setModalOpen(false)
                }
                selectedChainName={
                  selectedChainName
                }
              />
            </VStack>
          </TabPanel>
          <TabPanel>
            <VStack spacing={8} align="center">
              <Text
                textAlign="center"
                color="white"
              >
                Unstake your q
                {selectedOption.toUpperCase()}{' '}
                tokens in exchange for{' '}
                {selectedOption.toUpperCase()}.
              </Text>
              <Flex
                flexDirection="column"
                w="100%"
              >
                <Stat
                  py={4}
                  textAlign="left"
                  color="white"
                >
                  <StatLabel>
                    Amount tounstake:
                  </StatLabel>
                  <StatNumber>
                    q
                    {selectedOption.toUpperCase()}{' '}
                  </StatNumber>
                </Stat>
                <Input
                  _active={{
                    borderColor:
                      'complimentary.900',
                  }}
                  _selected={{
                    borderColor:
                      'complimentary.900',
                  }}
                  _hover={{
                    borderColor:
                      'complimentary.900',
                  }}
                  _focus={{
                    borderColor:
                      'complimentary.900',
                    boxShadow:
                      '0 0 0 3px #FF8000',
                  }}
                  color="complimentary.900"
                  textAlign={'right'}
                  placeholder="amount"
                />
                <Flex
                  w="100%"
                  flexDirection="row"
                  py={4}
                  mb={-4}
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <Text
                    color="white"
                    fontWeight="light"
                  >
                    Tokens available: 0 q
                    {selectedOption.toUpperCase()}
                  </Text>

                  <Button
                    _hover={{
                      bgColor:
                        'rgba(255,255,255,0.05)',
                      backdropFilter:
                        'blur(10px)',
                    }}
                    _active={{
                      bgColor:
                        'rgba(255,255,255,0.05)',
                      backdropFilter:
                        'blur(10px)',
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
                <Stat
                  textAlign="left"
                  color="white"
                >
                  <StatLabel>
                    What you&apos;ll get
                  </StatLabel>
                  <StatNumber>
                    {selectedOption.toUpperCase()}
                    :
                  </StatNumber>
                </Stat>
                <Spacer />{' '}
                {/* This pushes the next Stat component to the right */}
                <Stat
                  py={4}
                  textAlign="right"
                  color="white"
                >
                  <StatNumber textColor="complimentary.900">
                    0
                  </StatNumber>
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
