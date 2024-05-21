import { CloseIcon } from '@chakra-ui/icons';
import { Box, Flex, Text, VStack, Button, HStack, Spinner, Checkbox } from '@chakra-ui/react';
import { StdFee } from '@cosmjs/amino';
import { assets } from 'chain-registry';
import { GenericAuthorization } from 'interchain-query/cosmos/authz/v1beta1/authz';
import { quicksilver, cosmos } from 'quicksilverjs';
import React, { useState } from 'react';

import { useTx } from '@/hooks';
import { useFeeEstimation } from '@/hooks/useFeeEstimation';
import { useIncorrectAuthChecker, useLiquidEpochQuery } from '@/hooks/useQueries';

interface RewardsClaimInterface {
  address: string;
  onClose: () => void;
  refetch: () => void;
}
export const RewardsClaim: React.FC<RewardsClaimInterface> = ({ address, onClose, refetch }) => {
  const { tx } = useTx('quicksilver' ?? '');
  const { estimateFee } = useFeeEstimation('quicksilver');
  const { authData } = useIncorrectAuthChecker(address);

  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const { liquidEpoch } = useLiquidEpochQuery(address);

  const { submitClaim } = quicksilver.participationrewards.v1.MessageComposer.withTypeUrl;

  const { grant, revoke } = cosmos.authz.v1beta1.MessageComposer.withTypeUrl;

  const genericAuth = {
    msg: quicksilver.participationrewards.v1.MsgSubmitClaim.typeUrl,
  };

  const utf8Msg = GenericAuthorization.encode(genericAuth).finish();

  const msgGrant = grant({
    granter: address,
    grantee: 'quick1psevptdp90jad76zt9y9x2nga686hutgmasmwd',
    grant: {
      authorization: {
        typeUrl: cosmos.authz.v1beta1.GenericAuthorization.typeUrl,
        value: utf8Msg,
      },
    },
  });

  const msgRevokeBad = revoke({
    granter: address,
    grantee: 'quick1w5ennfhdqrpyvewf35sv3y3t8yuzwq29mrmyal',
    msgTypeUrl: quicksilver.participationrewards.v1.MsgSubmitClaim.typeUrl,
  });

  const handleAutoClaimRewards = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);

    const feeBoth: StdFee = {
      amount: [
        {
          denom: 'uqck',
          amount: '1000000',
        },
      ],
      gas: '2000000',
    };

    const feeSingle = await estimateFee(address, [msgGrant]);
    try {
      if (authData) {
        // Call msgRevokeBad and msgGrant
        await tx([msgRevokeBad, msgGrant], {
          fee: feeBoth,
          onSuccess: () => {
            refetch();
          },
        });
      } else {
        // Call msgGrant
        await tx([msgGrant], {
          fee: feeSingle,
          onSuccess: () => {
            refetch();
          },
        });
      }
    } catch (error) {
      console.error('Transaction failed', error);
      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

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
          //@ts-ignore
          proofs: transformedProofs,
        });
      });
      const fee = await estimateFee(address, msgSubmitClaims);
      await tx(msgSubmitClaims, {
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

  const handleAutoClaimToggle = () => {
    setAutoClaimEnabled(!autoClaimEnabled);
  };

  const transactionHandler = autoClaimEnabled ? handleAutoClaimRewards : handleClaimRewards;

  return (
    <Box bgColor="rgb(32,32,32)" maxW={'sm'} p="4" borderRadius="lg" mb="4">
      <Flex direction="column" alignItems="flex-end">
        <CloseIcon color="white" cursor="pointer" onClick={onClose} _hover={{ color: 'complimentary.900' }} />
        <VStack alignItems="flex-start" spacing="2">
          <Text fontSize="xl" fontWeight="bold" color="white">
            Cross Chain Claims
          </Text>
          <Text pb={2} color="white" fontSize="md">
            Click the button below to claim your cross chain rewards. Click the checkbox to enable automatic claiming.
          </Text>
          <HStack gap={8} justifyContent={'space-between'}>
            <Checkbox
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
              onChange={handleAutoClaimToggle}
              colorScheme="orange"
            >
              <Text color="white" fontSize="sm">
                Enable Automatic Claiming
              </Text>
            </Checkbox>

            {/* 
            // Section for showing a message when claims are disabled. DO NOT DELETE
            <Tooltip
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
                isDisabled={!address}
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
