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
} from '@chakra-ui/react';

import { Key, useState } from 'react';

import SignalIntentModal from './modals/signalIntentProcess';

import { useIntentQuery, useValidatorLogos, useValidatorsQuery } from '@/hooks/useQueries';
import { networks as prodNetworks, testNetworks as devNetworks } from '@/state/chains/prod';
import { truncateString } from '@/utils';

export interface StakingIntentProps {
  address: string;
  isWalletConnected: boolean;
}

const StakingIntent: React.FC<StakingIntentProps> = ({ address, isWalletConnected }) => {
  const networks = process.env.NEXT_PUBLIC_CHAIN_ENV === 'mainnet' ? prodNetworks : devNetworks;

  const chains = ['Cosmos', 'Osmosis', 'Stargaze', 'Regen', 'Sommelier'];
  const [currentChainIndex, setCurrentChainIndex] = useState(0);

  const [isSignalIntentModalOpen, setIsSignalIntentModalOpen] = useState(false);
  const openSignalIntentModal = () => setIsSignalIntentModalOpen(true);
  const closeSignalIntentModal = () => setIsSignalIntentModalOpen(false);

  const currentNetwork = networks[currentChainIndex];

  const { validatorsData } = useValidatorsQuery(currentNetwork.chainName);
  const { data: validatorLogos } = useValidatorLogos(currentNetwork.chainName, validatorsData || []);

  const { intent, isLoading, isError, refetch } = useIntentQuery(currentNetwork.chainName, address ?? '');

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
    intent?.data?.intent.intents.map((validatorIntent: { valoper_address: string; weight: string }) => {
      const validatorDetails = validatorsMap[validatorIntent.valoper_address];
      return {
        moniker: validatorDetails?.moniker,
        logoUrl: validatorDetails?.logoUrl,
        percentage: `${(parseFloat(validatorIntent.weight) * 100).toFixed(2)}%`,
      };
    }) || [];

  const handleLeftArrowClick = () => {
    setCurrentChainIndex((prevIndex) => (prevIndex === 0 ? networks.length - 1 : prevIndex - 1));
  };

  const handleRightArrowClick = () => {
    setCurrentChainIndex((prevIndex) => (prevIndex === networks.length - 1 ? 0 : prevIndex + 1));
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
        <Spinner w={'200px'} h="200px" color="complimentary.900" />
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
          <Button color="GrayText" _hover={{ color: 'complimentary.900' }} variant="link" onClick={openSignalIntentModal}>
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

        <Flex borderBottom="1px" borderBottomColor="complimentary.900" alignItems="center" justifyContent="space-between">
          <IconButton
            variant="ghost"
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            aria-label="Previous chain"
            icon={<ChevronLeftIcon w={'25px'} h={'25px'} />}
            onClick={handleLeftArrowClick}
          />
          <SlideFade in={true} key={currentChainIndex}>
            <Text fontSize="lg" fontWeight="semibold">
              {chains[currentChainIndex]}
            </Text>
          </SlideFade>
          <IconButton
            _active={{
              transform: 'scale(0.75)',
              color: 'complimentary.800',
            }}
            _hover={{ bgColor: 'transparent', color: 'complimentary.900' }}
            variant="ghost"
            aria-label="Next chain"
            icon={<ChevronRightIcon w={'25px'} h={'25px'} />}
            onClick={handleRightArrowClick}
          />
        </Flex>

        <VStack pb={4} overflowY="auto" gap={4} spacing={2} align="stretch" maxH="250px">
          {validatorsWithDetails.map(
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
                      startColor="complimentary.900"
                      endColor="complimentary.100"
                    />
                  )}
                  {validator.moniker ? (
                    <Text fontSize="md">{truncateString(validator.moniker, 20)}</Text>
                  ) : (
                    <SkeletonText
                      display="inline-block"
                      verticalAlign="middle"
                      startColor="complimentary.900"
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
          )}
        </VStack>
      </VStack>
    </Box>
  );
};

export default StakingIntent;
