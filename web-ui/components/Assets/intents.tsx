import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import {
  Box,
  Flex,
  Text,
  Button,
  IconButton,
  VStack,
  Image,
  Heading,
  SlideFade,
  Spinner,
  SkeletonCircle,
  SkeletonText,
  Center,
  Fade,
} from '@chakra-ui/react';
import { Key, useCallback, useState } from 'react';

import { chains, env, Chain } from '@/config';
import { useIntentQuery, useValidatorLogos, useValidatorsQuery } from '@/hooks/useQueries';
import { truncateString } from '@/utils';

import SignalIntentModal from './modals/signalIntentProcess';

export interface StakingIntentProps {
  address: string;
  isWalletConnected: boolean;
}

const StakingIntent: React.FC<StakingIntentProps> = ({ address, isWalletConnected }) => {
  const networks: Map<string, Chain> = chains.get(env) ?? new Map();
  const chain_list = Array.from(networks).filter(([_, network]) => network.show).map(([key, _]) => key);


  const [isBottomVisible, setIsBottomVisible] = useState(true);

  const handleScroll = useCallback((event: React.UIEvent<HTMLDivElement>) => {
    const target = event.currentTarget;
    const isBottom = target.scrollHeight - target.scrollTop <= target.clientHeight;
    setIsBottomVisible(!isBottom);
  }, []);

  const [isSignalIntentModalOpen, setIsSignalIntentModalOpen] = useState(false);
  const openSignalIntentModal = () => setIsSignalIntentModalOpen(true);
  const closeSignalIntentModal = () => setIsSignalIntentModalOpen(false);

  const [currentChainName, setCurrentChainName] = useState(chain_list[0]);



  const currentNetwork = networks?.get(currentChainName)
  const prev = () => {
    const index = chain_list?.indexOf(currentChainName) - 1
    return index < 0 ? chain_list.length - 1 : index
  }

  const next = () => {
    const index = chain_list?.indexOf(currentChainName) + 1
    return index > chain_list.length - 1 ? 0 : index
  }

  const { validatorsData } = useValidatorsQuery(currentChainName);
  const { data: validatorLogos } = useValidatorLogos(currentChainName, validatorsData || []);

  const { intent, refetch } = useIntentQuery(currentChainName, address ?? '');

  interface ValidatorDetails {
    moniker: string;
    logoUrl: string | undefined;
  }

  interface ValidatorMap {
    [valoper_address: string]: ValidatorDetails;
  }

  const validatorsMap: ValidatorMap =
    validatorsData?.reduce((map: ValidatorMap, validatorInfo) => {
      map[validatorInfo.address] = {
        moniker: validatorInfo.name,
        logoUrl: validatorLogos?.[validatorInfo.address],
      };
      return map;
    }, {}) || {};

  const validatorsWithDetails =
    intent?.data?.intent.intents
      .filter((validatorIntent: { valoper_address: string; weight: string }) => parseFloat(validatorIntent.weight) > 0)
      .map((validatorIntent: { valoper_address: string; weight: string }) => {
        const validatorDetails = validatorsMap[validatorIntent.valoper_address];
        return {
          moniker: validatorDetails?.moniker,
          logoUrl: validatorDetails?.logoUrl,
          percentage: `${(parseFloat(validatorIntent.weight) * 100).toFixed(2)}%`,
        };
      }) || [];

  const handleLeftArrowClick = () => {
    setCurrentChainName(chain_list[prev()]);
  };

  const handleRightArrowClick = () => {
    setCurrentChainName(chain_list[next()]);
  };

  if (!isWalletConnected) {
    return (
      <Flex direction="column" p={5} borderRadius="lg" align="center" justify="space-around" w="full" h="full">
        <Text fontSize="xl" textAlign="center">
          Wallet is not connected. Please connect your wallet to interact with your QCK tokens.
        </Text>
      </Flex>
    );
  }

  if (!intent) {
    return (
      <Flex
        w="100%"
        h="100%"
        p={4}
        borderRadius="lg"
        flexDirection="column"
        justifyContent="center"
        alignItems="center"
        gap={6}
        color="white"
      >
        <Spinner w={'200px'} h="200px" color="complimentary.700" />
      </Flex>
    );
  }

  return (
    <Box w="full" color="white" borderRadius="lg" p={4} gap={6}>
      <VStack spacing={4} align="stretch">
        <Flex gap={6} justifyContent="space-between" alignItems="center">
          <Heading fontSize="lg" fontWeight="bold" textTransform="uppercase">
            Stake Intent
          </Heading>
          <Button color="GrayText" _hover={{ color: 'complimentary.700' }} variant="link" onClick={openSignalIntentModal}>
            Edit Intent
            <ChevronRightIcon />
          </Button>
          <SignalIntentModal
            refetch={refetch}
            selectedOption={currentNetwork}
            isOpen={isSignalIntentModalOpen}
            onClose={closeSignalIntentModal}
          />
        </Flex>

        <Flex borderBottom="1px" borderBottomColor="complimentary.700" alignItems="center" justifyContent="space-between">
          <IconButton
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            color="GrayText"
            aria-label="Previous chain"
            icon={<ChevronLeftIcon w={'25px'} h={'25px'} />}
            onClick={handleLeftArrowClick}
          />
          <SlideFade in={true} key={currentChainName}>
            <Text fontSize="lg" fontWeight="semibold">
              {chains.get(env)?.get(currentChainName)?.pretty_name}
            </Text>
          </SlideFade>
          <IconButton
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            color="GrayText"
            _hover={{ bgColor: 'transparent', color: 'complimentary.700' }}
            variant="ghost"
            aria-label="Next chain"
            icon={<ChevronRightIcon w={'25px'} h={'25px'} />}
            onClick={handleRightArrowClick}
          />
        </Flex>

        <VStack
          onScroll={handleScroll}
          pb={4}
          overflowY="auto"
          className="custom-scrollbar"
          gap={4}
          spacing={2}
          align="stretch"
          maxH="210px"
        >
          {(validatorsWithDetails.length > 0 &&
            validatorsWithDetails.map(
              (validator: { logoUrl: string; moniker: string; percentage: string }, index: Key | null | undefined) => (
                <Flex key={index} justifyContent="space-between" w="full" alignItems="center">
                  <Flex alignItems="center" gap={2}>
                    {validator.logoUrl ? (
                      <Image
                        borderRadius={'full'}
                        src={validator.logoUrl}
                        alt={validator.moniker}
                        boxSize="26px"
                        objectFit="cover"
                        marginRight="8px"
                      />
                    ) : (
                      <SkeletonCircle
                        boxSize="26px"
                        objectFit="cover"
                        marginRight="8px"
                        display="inline-block"
                        verticalAlign="middle"
                        startColor="complimentary.700"
                        endColor="complimentary.100"
                      />
                    )}
                    {validator.moniker ? (
                      <Text fontSize="md">{truncateString(validator.moniker, 18)}</Text>
                    ) : (
                      <SkeletonText
                        display="inline-block"
                        verticalAlign="middle"
                        startColor="complimentary.700"
                        endColor="complimentary.100"
                        noOfLines={1}
                        width="100px"
                      />
                    )}
                  </Flex>
                  <Text fontSize="lg" fontWeight="bold">
                    {validator.percentage}
                  </Text>
                </Flex>
              ),
            )) || (
            <Center mt={6}>
              <Text fontSize="xl">No intent set</Text>
            </Center>
          )}
          {isBottomVisible && validatorsWithDetails.length > 5 && (
            <Fade in={isBottomVisible}>
              <Box
                borderRadius="lg"
                position="absolute"
                bottom="0"
                left="0"
                right="0"
                height="110px"
                bgGradient="linear(to top, #1A1A1A, transparent)"
                pointerEvents="none"
                zIndex="10"
              />
            </Fade>
          )}
        </VStack>
      </VStack>
    </Box>
  );
};

export default StakingIntent;
