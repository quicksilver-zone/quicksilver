import { CloseIcon } from '@chakra-ui/icons';
import { Box, Flex, Text, VStack, Button, HStack, Spinner } from '@chakra-ui/react';
import { assets } from 'chain-registry';
import { GenericAuthorization } from 'interchain-query/cosmos/authz/v1beta1/authz';
import { quicksilver, cosmos } from 'quicksilverjs';
import React, { useState } from 'react';

import { useTx } from '@/hooks';
import { useLiquidEpochQuery } from '@/hooks/useQueries';

interface RewardsClaimInterface {
  address: string;
  onClose: () => void;
}

function transformProofs(proofs: any[]) {
  return proofs.map((proof) => ({
    key: proof.key,
    data: proof.data,
    proofOps: proof.proof_ops
      ? {
          //@ts-ignore
          ops: proof.proof_ops.ops.map((op) => ({
            type: op.type,
            key: op.key,
            data: op.data,
          })),
        }
      : undefined,
    height: proof.height,
    proofType: proof.proof_type,
  }));
}

export const RewardsClaim: React.FC<RewardsClaimInterface> = ({ address, onClose }) => {
  const { tx } = useTx('quicksilver' ?? '');

  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const { liquidEpoch } = useLiquidEpochQuery(address);

  const { submitClaim } = quicksilver.participationrewards.v1.MessageComposer.withTypeUrl;

  const { grant } = cosmos.authz.v1beta1.MessageComposer.withTypeUrl;

  const genericAuth = {
    msg: quicksilver.participationrewards.v1.MsgSubmitClaim.typeUrl,
  };

  const binaryMessage = GenericAuthorization.encode(genericAuth).finish();

  const msgGrant = grant({
    granter: address,
    grantee: 'quick1w5ennfhdqrpyvewf35sv3y3t8yuzwq29mrmyal',
    grant: {
      authorization: {
        typeUrl: cosmos.authz.v1beta1.GenericAuthorization.typeUrl,
        value: binaryMessage,
      },
    },
  });

  const mainTokens = assets.find(({ chain_name }) => chain_name === 'quicksilver');
  const mainDenom = mainTokens?.assets[0].base ?? 'uqck';

  const fee = {
    amount: [
      {
        denom: mainDenom,
        amount: '50',
      },
    ],
    gas: '500000',
  };

  const handleAutoClaimRewards = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);

    try {
      const result = await tx([msgGrant], {
        fee,
        onSuccess: () => {},
      });
    } catch (error) {
      console.error('Transaction failed', error);

      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

  const handleClaimRewards = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);

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

      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

  const [autoClaimEnabled, setAutoClaimEnabled] = useState(true);

  // const handleAutoClaimToggle = () => {
  //   setAutoClaimEnabled(!autoClaimEnabled);
  // };

  const transactionHandler = autoClaimEnabled ? handleAutoClaimRewards : handleClaimRewards;

  return (
    <Box bgColor="rgb(32,32,32)" maxW={'sm'} p="4" borderRadius="lg" mb="4">
      <Flex direction="column" alignItems="flex-end">
        <CloseIcon color="white" cursor="pointer" onClick={onClose} _hover={{ color: 'complimentary.900' }} />
        <VStack alignItems="flex-start" spacing="2">
          <Text fontSize="xl" fontWeight="bold" color="white">
            Cross Chain Claims (XCC) is coming!
          </Text>
          <Text pb={2} color="white" fontSize="md">
            Click the button below to set your authz grant for automatic cross chain claims.
          </Text>
          <HStack gap={8} justifyContent={'space-between'}>
            {/* <Checkbox
              _selected={{ bgColor: 'transparent' }}
              _active={{
                borderColor: 'complimentary.900',
              }}
              _hover={{
                borderColor: 'complimentary.900',
              }}
              _focus={{
                borderColor: 'complimentary.900',
                boxShadow: '0 0 0 3px #FF8000',
              }}
              isChecked={autoClaimEnabled}
              //onChange={handleAutoClaimToggle}
              colorScheme="orange"
            >
              <Text color="white" fontSize="sm">
                Enable Automatic Claiming
              </Text>
            </Checkbox> */}
            {/* <Tooltip
              mr={12}
              label={
                <React.Fragment>
                  <Flex direction="column" p="2" maxW="xs" textAlign={'center'} align="center">
                    <Heading color="red" fontSize="md" fontWeight="bold" pb={2}>
                      Claiming Disabled
                    </Heading>
                    <Text>Reward claiming is disabled until a governance proposal updating the claiming parameter is passed.</Text>
                    <Box p="2">
                      <Text>You may still enable automatic claiming in advance.</Text>
                    </Box>
                  </Flex>
                </React.Fragment>
              }
              isDisabled={autoClaimEnabled}
              placement="top"
              hasArrow
            > */}
              <Box>
                <Button
                  _active={{
                    transform: 'scale(0.95)',
                    color: 'complimentary.800',
                  }}
                  _hover={{
                    bgColor: 'rgba(255,128,0, 0.25)',
                    color: 'complimentary.300',
                  }}
                  minW={'120px'}
                  onClick={transactionHandler}
                  size="sm"
                  alignSelf="end"
                  isDisabled={!autoClaimEnabled}
                >
                  {isError ? 'Try Again' : isSigning ? <Spinner /> : autoClaimEnabled ? 'Auto Claim' : 'Claim Rewards'}
                </Button>
              </Box>
            {/* </Tooltip> */}
          </HStack>
        </VStack>
      </Flex>
    </Box>
  );
};

export default RewardsClaim;
