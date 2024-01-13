import { ArrowForwardIcon, InfoIcon } from '@chakra-ui/icons';
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
import { Key, useState } from 'react';

import { useLiquidEpochQuery, useLiquidRewardsQuery } from '@/hooks/useQueries';
import { shiftDigits } from '@/utils';
import { quicksilver } from 'quicksilverjs';
import { useTx } from '@/hooks';
import { StdFee } from '@cosmjs/amino';
import { assets } from 'chain-registry';
import { cosmos } from 'interchain-query';
import { Grant, GenericAuthorization } from 'interchain-query/cosmos/authz/v1beta1/authz';
import { MsgGrant } from 'interchain-query/cosmos/authz/v1beta1/tx';

interface RewardsClaimInterface {
  address: string;
}

function transformProofs(proofs: any[]) {
  return proofs.map((proof) => ({
    key: proof.key, // Convert from base64 to Uint8Array if needed
    data: proof.data, // Convert from base64 to Uint8Array if needed
    proofOps: proof.proof_ops
      ? {
          ops: proof.proof_ops.ops.map((op) => ({
            type: op.type,
            key: op.key, // Convert from base64 to Uint8Array if needed
            data: op.data, // Convert from base64 to Uint8Array if needed
          })),
        }
      : undefined,
    height: proof.height,
    proofType: proof.proof_type,
  }));
}

export const RewardsClaim: React.FC<RewardsClaimInterface> = ({ address }) => {
  const { tx } = useTx('quicksilver' ?? '');

  const [transactionStatus, setTransactionStatus] = useState('Pending');
  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);
  const { liquidRewards, isLoading } = useLiquidRewardsQuery(address);
  const { liquidEpoch, isLoading: isEpochLoading } = useLiquidEpochQuery(address);
  console.log(liquidEpoch);

  const { isOpen, onOpen, onClose } = useDisclosure();

  const { submitClaim } = quicksilver.participationrewards.v1.MessageComposer.withTypeUrl;

  const { grant } = cosmos.authz.v1beta1.MessageComposer.withTypeUrl;

  const msgGrant = grant({
    granter: 'quick1c4vz0535677xpdksxh5um7zqqwfsw7245ppdaj',
    grantee: address,
    grant: Grant.fromPartial({
      authorization: GenericAuthorization.fromPartial({
        msg: '/quicksilver.participationrewards.v1.MsgSubmitClaim',
      }),
    }),
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

  const handleAutoClaimRewards = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);
    setTransactionStatus('Pending');
    try {
      const result = await tx([msgGrant], {
        fee,
        onSuccess: () => {},
      });
    } catch (error) {
      console.error('Transaction failed', error);
      setTransactionStatus('Failed');
      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

  const handleClaimRewards = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);
    setTransactionStatus('Pending');

    if (!liquidEpoch || liquidEpoch.messages.length === 0) {
      console.error('No epoch data available or no messages to claim');
      setIsSigning(false);
      return;
    }

    try {
      const msgSubmitClaims = liquidEpoch.messages.map((message) => {
        const transformedProofs = transformProofs(message.proofs);
        return submitClaim({
          userAddress: message.user_address,
          zone: message.zone,
          srcZone: message.src_zone,
          claimType: message.claim_type,
          proofs: transformedProofs,
        });
      });

      const result = await tx(msgSubmitClaims, {
        fee,
        onSuccess: () => {},
      });
    } catch (error) {
      console.error('Transaction failed', error);
      setTransactionStatus('Failed');
      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

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
