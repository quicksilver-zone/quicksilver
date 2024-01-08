import {
  Box,
  Flex,
  Text,
  VStack,
  Button,
  Switch,
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalBody,
  Stack,
  ModalCloseButton,
  ModalHeader,
} from '@chakra-ui/react';
import { ArrowForwardIcon, InfoIcon } from '@chakra-ui/icons';
import { useLiquidRewardsQuery } from '@/hooks/useQueries';
import { Key } from 'react';
import { shiftDigits } from '@/utils';

interface RewardsClaimInterface {
  address: string;
}

export const RewardsClaim: React.FC<RewardsClaimInterface> = ({ address }) => {
  const { liquidRewards, isLoading } = useLiquidRewardsQuery(address);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const handleClaimRewards = () => {};

  return (
    <>
      <Text fontSize="xl" fontWeight="bold" color="white" mb={4}>
        Participation Rewards
      </Text>
      <Flex
        flexDirection={['column', 'column', 'row']}
        justifyContent="space-between"
        alignItems="flex-start"
        bgColor="rgba(255,255,255,0.1)"
        p="4"
        borderRadius="lg"
        mb="4"
        gap="6"
      >
        <VStack flex="1" spacing="3.5" alignItems="flex-start">
          <Text color="white" fontSize="base" fontWeight="normal">
            Stake with validators with a high PR score to earn QCK rewards. Automatic claiming of rewards is{' '}
            <Text as="span" textDecoration="underline">
              required
            </Text>{' '}
            for the protocol to consider your validator staking intent.
          </Text>
          <Button leftIcon={<InfoIcon />} variant="link" colorScheme="blue" onClick={onOpen}>
            Learn more about Participation Rewards
          </Button>
        </VStack>

        <Box flex="2" overflowY="auto" maxH="300px" p="4" borderRadius="lg" border="1px" borderColor="white" maxW={'200px'}>
          <Stack spacing={4}>
            {!isLoading &&
              liquidRewards?.assets?.['rhye-2']?.map((assetGroup) =>
                assetGroup.Amount.map((asset, index) => (
                  <Text key={index} color="white" fontSize="sm">
                    {Number(shiftDigits(asset.amount, -6)).toLocaleString()} {asset.denom.toUpperCase().slice(1)}
                  </Text>
                )),
              )}
            {isLoading && <Text>Loading rewards...</Text>}
          </Stack>
        </Box>

        <VStack flex="1" spacing="3.5" alignItems="flex-end">
          <Button size="lg" colorScheme="blue" onClick={handleClaimRewards} isDisabled={isLoading || !liquidRewards}>
            Claim All Rewards
          </Button>
        </VStack>
      </Flex>

      <Modal isOpen={isOpen} onClose={onClose} isCentered>
        <ModalOverlay />
        <ModalContent backgroundColor="gray.800" color="white">
          <ModalHeader>Participation Rewards</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>More information about participation rewards...</Text>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
};

export default RewardsClaim;
