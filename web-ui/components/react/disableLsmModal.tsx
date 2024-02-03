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
} from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import { assets, chains } from 'chain-registry';
import { StdFee } from 'interchain-query';
import { cosmos } from 'quicksilverjs';
import { MsgDisableTokenizeShares, MsgEnableTokenizeShares } from 'quicksilverjs/dist/codegen/cosmos/staking/v1beta1/lsm';
import { useState } from 'react';

import { useTx } from '@/hooks';

interface DisableLsmModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const DisableLsmModal: React.FC<DisableLsmModalProps> = ({ isOpen, onClose }) => {
  const { address } = useChain('cosmoshub');
  const { tx } = useTx('cosmoshub');

  const msgDisable = MsgDisableTokenizeShares.fromPartial({
    delegator_address: address ?? '',
  });

  const msgEnable = MsgEnableTokenizeShares.fromPartial({
    delegator_address: address ?? '',
  });

  const [isSigningEnable, setIsSigningEnable] = useState<boolean>(false);
  const [isSigningDisable, setIsSingingDisable] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const mainTokens = assets.find(({ chain_name }) => chain_name === chain_name);
  const fees = chains.find(({ chain_name }) => chain_name === chain_name)?.fees?.fee_tokens;
  const mainDenom = mainTokens?.assets[0].base ?? '';
  const fixedMinGasPrice = fees?.find(({ denom }) => denom === mainDenom)?.high_gas_price ?? '';
  const feeAmount = Number(fixedMinGasPrice) * 750000;
  const sendFeeAmount = Number(fixedMinGasPrice) * 100000;

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
      const result = await tx([], {
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
      const result = await tx([], {
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

  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered size="lg">
      <ModalOverlay />
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
            Keep in mind you will not be able to natively stake your LSM-enabled assets like Atom without re-enabling LSM.
          </Text>
        </ModalBody>
        <ModalFooter>
          <Button
            isDisabled={!address}
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
            variant="ghost"
            minW={'100px'}
          >
            {isError ? 'Try Again' : isSigningDisable ? <Spinner /> : 'Disable'}
          </Button>
          <Button
            minW={'100px'}
            isDisabled={!address}
            mr={3}
            onClick={handleEnable}
            _hover={{
              bgColor: 'rgba(255,255,255,0.05)',
              backdropFilter: 'blur(10px)',
            }}
            _active={{
              bgColor: 'rgba(255,255,255,0.05)',
              backdropFilter: 'blur(10px)',
            }}
            color="green"
            variant="ghost"
          >
            {isError ? 'Try Again' : isSigningEnable ? <Spinner /> : 'Enable'}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
