import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  Text,
  Spacer,
  Spinner,
  Flex,
  Box,
} from '@chakra-ui/react';
import { StdFee } from '@cosmjs/amino';
import { useChain } from '@cosmos-kit/react';
import { assets, chains } from 'chain-registry';
import { quicksilver, cosmos } from 'quicksilverjs';
import { GenericAuthorization } from 'quicksilverjs/dist/codegen/cosmos/authz/v1beta1/authz';
import { useState } from 'react';

import { useTx } from '@/hooks';
import { useAuthChecker, useIncorrectAuthChecker } from '@/hooks/useQueries';

interface AccountControlModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const AccountControlModal: React.FC<AccountControlModalProps> = ({ isOpen, onClose }) => {
  const [authzSection, setAuthzSection] = useState<boolean>(false);
  const [lsmSection, setLsmSection] = useState<boolean>(false);

  const { address: lsmAddress } = useChain('cosmoshub');
  const { tx: lsmTx } = useTx('cosmoshub');

  const { address: authAddress } = useChain('quicksilver');
  const { tx: authTx } = useTx('quicksilver');

  const { authData: incorrectAccount } = useIncorrectAuthChecker(authAddress ?? '');
  const { authData: correctAccount } = useAuthChecker(authAddress ?? '');

  const { enableTokenizeShares } = cosmos.staking.v1beta1.MessageComposer.withTypeUrl;
  const { disableTokenizeShares } = cosmos.staking.v1beta1.MessageComposer.withTypeUrl;

  const msgEnable = enableTokenizeShares({
    delegatorAddress: lsmAddress ?? '',
  });

  const msgDisable = disableTokenizeShares({
    delegatorAddress: lsmAddress ?? '',
  });

  const [isSigningEnable, setIsSigningEnable] = useState<boolean>(false);
  const [isSigningDisable, setIsSingingDisable] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const mainTokens = assets.find(({ chain_name }) => chain_name === chain_name);
  const fees = chains.find(({ chain_name }) => chain_name === chain_name)?.fees?.fee_tokens;
  const mainDenom = mainTokens?.assets[0].base ?? '';
  const fixedMinGasPrice = fees?.find(({ denom }) => denom === mainDenom)?.high_gas_price ?? '';
  const feeAmount = Number(fixedMinGasPrice) * 750000;

  const fee: StdFee = {
    amount: [
      {
        denom: mainDenom,
        amount: feeAmount.toString(),
      },
    ],
    gas: '750000', // test txs were using well in excess of 600k
  };

  const handleDisable = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSingingDisable(true);

    try {
      await lsmTx([msgDisable], {
        fee,
        onSuccess: () => {
          onClose();
        },
      });
    } catch (error) {
      console.error('Transaction failed', error);

      setIsError(true);
    } finally {
      setIsSingingDisable(false);
    }
  };

  const handleEnable = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigningEnable(true);

    try {
      await lsmTx([msgEnable], {
        fee,
        onSuccess: () => {
          onClose();
        },
      });
    } catch (error) {
      console.error('Transaction failed', error);

      setIsError(true);
    } finally {
      setIsSigningEnable(false);
    }
  };

  const { grant, revoke } = cosmos.authz.v1beta1.MessageComposer.withTypeUrl;

  const genericAuth = {
    msg: quicksilver.participationrewards.v1.MsgSubmitClaim.typeUrl,
  };

  const binaryMessage = GenericAuthorization.encode(genericAuth).finish();
  const msgGrant = grant({
    granter: authAddress ?? '',
    grantee: 'quick1psevptdp90jad76zt9y9x2nga686hutgmasmwd',
    grant: {
      authorization: {
        typeUrl: cosmos.authz.v1beta1.GenericAuthorization.typeUrl,
        value: binaryMessage,
      },
    },
  });

  const revokeGrant = revoke({
    granter: authAddress ?? '',
    grantee: 'quick1psevptdp90jad76zt9y9x2nga686hutgmasmwd',

    msgTypeUrl: quicksilver.participationrewards.v1.MsgSubmitClaim.typeUrl,
  });

  const msgRevokeBad = revoke({
    granter: authAddress ?? '',
    grantee: 'quick1w5ennfhdqrpyvewf35sv3y3t8yuzwq29mrmyal',
    msgTypeUrl: quicksilver.participationrewards.v1.MsgSubmitClaim.typeUrl,
  });

  const handleAutoClaimRewards = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigningEnable(true);

    try {
      await authTx([msgGrant], {
        fee,
        onSuccess: () => {},
      });
    } catch (error) {
      console.error('Transaction failed', error);

      setIsError(true);
    } finally {
      setIsSigningEnable(false);
    }
  };

  const handleRemoveAutoClaim = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSingingDisable(true);

    try {
      if (incorrectAccount) {
        // Call msgRevokeBad
        await authTx([msgRevokeBad], {
          fee,
          onSuccess: () => {},
        });
      }
      // Continue with msgGrant
      if (correctAccount) {
        // Call msgRevokeBad
        await authTx([revokeGrant], {
          fee,
          onSuccess: () => {},
        });
      }
    } catch (error) {
      console.error('Transaction failed', error);
      setIsError(true);
    } finally {
      setIsSingingDisable(false);
    }
  };

  const handleAuthzSection = () => {
    setAuthzSection(!authzSection);
  };

  const handleLsmSection = () => {
    setLsmSection(!lsmSection);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered size="lg">
      <ModalOverlay />

      {!authzSection && !lsmSection && (
        <ModalContent bg={'#1a1a1a'}>
          <ModalHeader color={'white'}>Account Controls</ModalHeader>
          <ModalCloseButton color={'white'} />
          <ModalBody px={4} py={2}>
            <Flex gap={18} flexDirection={'row'} justifyContent="space-between">
              <Flex flexDirection={'column'} width="50%">
                <Text textAlign={'center'} color="white" mb={4}>
                  LSM Controls
                </Text>
                <Button
                  onClick={handleLsmSection}
                  _active={{
                    transform: 'scale(0.95)',
                    color: 'complimentary.800',
                  }}
                  _hover={{
                    bgColor: 'rgba(255,128,0, 0.25)',
                    color: 'complimentary.300',
                  }}
                  size="sm"
                >
                  Control LSM
                </Button>
              </Flex>

              <Box mt={4} width="1px" bg="orange" mx={2} />

              <Flex flexDirection={'column'} width="50%">
                <Text textAlign={'center'} color="white" mb={4}>
                  Authz Controls
                </Text>
                <Button
                  onClick={handleAuthzSection}
                  _active={{
                    transform: 'scale(0.95)',
                    color: 'complimentary.800',
                  }}
                  _hover={{
                    bgColor: 'rgba(255,128,0, 0.25)',
                    color: 'complimentary.300',
                  }}
                  size="sm"
                >
                  Control Authz
                </Button>
              </Flex>
            </Flex>
          </ModalBody>
          <ModalFooter>{/* Buttons or any other footer content */}</ModalFooter>
        </ModalContent>
      )}

      {/* Authz Section */}
      {authzSection && (
        <ModalContent bg={'#1a1a1a'}>
          <ModalHeader color={'white'}>XCC Authz Controls</ModalHeader>
          <ModalCloseButton color={'white'} />
          <ModalBody px={4} py={2}>
            <Text color="white" lineHeight="tall">
              Disable or reenable the ability to auto claim your cross chain rewards.
            </Text>
            <Spacer h={4} />
            <Text color="white" lineHeight="tall">
              Disabling this feature will prevent you from automatically claiming your cross chain rewards.
            </Text>
          </ModalBody>
          <ModalFooter>
            <Flex justify="space-between" width="full">
              <Button
                onClick={() => setAuthzSection(false)}
                _hover={{
                  bgColor: 'rgba(255,255,255,0.05)',
                  backdropFilter: 'blur(10px)',
                }}
                _active={{
                  bgColor: 'rgba(255,255,255,0.05)',
                  backdropFilter: 'blur(10px)',
                }}
                color="white"
                variant="ghost"
                minW={'100px'}
              >
                Back
              </Button>
              <Flex>
                <Button
                  isDisabled={!authAddress}
                  mr={3}
                  onClick={handleRemoveAutoClaim}
                  _hover={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  _active={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  color="red"
                  bgColor="rgba(176, 54, 54, 0.4)"
                  minW={'100px'}
                >
                  {isError ? 'Try Again' : isSigningDisable ? <Spinner /> : 'Revoke'}
                </Button>
                <Button
                  minW={'100px'}
                  isDisabled={!authAddress}
                  onClick={handleAutoClaimRewards}
                  _hover={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  _active={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  borderColor={'green'}
                  color="green"
                  variant="outline"
                >
                  {isError ? 'Try Again' : isSigningEnable ? <Spinner /> : 'Grant'}
                </Button>
              </Flex>
            </Flex>
          </ModalFooter>
        </ModalContent>
      )}

      {/* LSM Section */}
      {lsmSection && (
        <ModalContent bg={'#1a1a1a'}>
          <ModalHeader color={'white'}>Liquid Staking Module Controls</ModalHeader>
          <ModalCloseButton color={'white'} />
          <ModalBody px={4} py={2}>
            <Text color="white" lineHeight="tall">
              If your wallet is compromised, hackers can easily tokenize your staked assets and steal them. Disabling LSM prevents this from
              happening.
            </Text>
            <Spacer h={4} />
            <Text color="white" lineHeight="tall">
              Remember that you will not be able to directly stake your LSM-supported assets, such as Atom, unless you re-enable LSM.
            </Text>
          </ModalBody>
          <ModalFooter>
            <Flex justify="space-between" width="full">
              <Button
                onClick={() => setLsmSection(false)}
                _hover={{
                  bgColor: 'rgba(255,255,255,0.05)',
                  backdropFilter: 'blur(10px)',
                }}
                _active={{
                  bgColor: 'rgba(255,255,255,0.05)',
                  backdropFilter: 'blur(10px)',
                }}
                color="white"
                variant="ghost"
                minW={'100px'}
              >
                Back
              </Button>
              <Flex>
                <Button
                  isDisabled={!lsmAddress}
                  mr={3}
                  onClick={handleDisable}
                  _hover={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  _active={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  color="red"
                  bgColor="rgba(176, 54, 54, 0.4)"
                  minW={'100px'}
                >
                  {isError ? 'Try Again' : isSigningDisable ? <Spinner /> : 'Disable'}
                </Button>
                <Button
                  minW={'100px'}
                  isDisabled={!lsmAddress}
                  onClick={handleEnable}
                  _hover={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  _active={{
                    bgColor: 'rgba(255,255,255,0.05)',
                    backdropFilter: 'blur(10px)',
                  }}
                  borderColor={'green'}
                  color="green"
                  variant="outline"
                >
                  {isError ? 'Try Again' : isSigningEnable ? <Spinner /> : 'Enable'}
                </Button>
              </Flex>
            </Flex>
          </ModalFooter>
        </ModalContent>
      )}
    </Modal>
  );
};
