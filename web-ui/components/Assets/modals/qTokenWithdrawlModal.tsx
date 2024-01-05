import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  useDisclosure,
  useToast,
  Spinner,
} from '@chakra-ui/react';
import { ChainName } from '@cosmos-kit/core';
import { useChain, useManager } from '@cosmos-kit/react';
import { useState, useMemo, useEffect } from 'react';

import { ChooseChain } from '@/components/react/choose-chain';
import { handleSelectChainDropdown, ChainOption } from '@/components/types';
import { ibc } from 'interchain-query';
import { useBalanceQuery, useIbcBalanceQuery } from '@/hooks/useQueries';
import { useTx } from '@/hooks';
import BigNumber from 'bignumber.js';
import { getCoin, getIbcInfo } from '@/utils';
import { StdFee, coins } from '@cosmjs/stargate';
import { store } from '@interchain-ui/react';

interface QDepositModalProps {
  token: string;
}

const QWithdrawModal: React.FC<QDepositModalProps> = ({ token }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();

  const [chainName, setChainName] = useState<ChainName | undefined>('osmosis');
  const { chainRecords, getChainLogo } = useManager();
  const [amount, setAmount] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const chainOptions = useMemo(() => {
    const desiredChains = ['osmosis', 'secretnetwork', 'umee'];
    return chainRecords
      .filter((chainRecord) => desiredChains.includes(chainRecord.name))
      .map((chainRecord) => ({
        chainName: chainRecord?.name,
        label: chainRecord?.chain.pretty_name,
        value: chainRecord?.name,
        icon: getChainLogo(chainRecord.name),
      }));
  }, [chainRecords, getChainLogo]);

  useEffect(() => {
    setChainName(window.localStorage.getItem('selected-chain') || 'osmosis');
  }, []);

  const onChainChange: handleSelectChainDropdown = async (selectedValue: ChainOption | null) => {
    setChainName(selectedValue?.chainName);
    if (selectedValue?.chainName) {
      window?.localStorage.setItem('selected-chain', selectedValue?.chainName);
    } else {
      window?.localStorage.removeItem('selected-chain');
    }
  };

  const chooseChain = <ChooseChain chainName={chainName} chainInfos={chainOptions} onChange={onChainChange} />;

  const fromChain = chainName;
  const toChain = 'quicksilver';

  const { transfer } = ibc.applications.transfer.v1.MessageComposer.withTypeUrl;
  const { address, connect, status, message, wallet } = useChain(fromChain ?? '');
  const { address: qAddress } = useChain('quicksilver');
  const { balance } = useIbcBalanceQuery(fromChain ?? '', address ?? '');
  const { tx } = useTx(fromChain ?? '');
  const qckBalance =
    balance?.balances.find((b) => b.denom === 'ibc/635CB83EF1DFE598B10A3E90485306FD0D47D34217A4BE5FD9977FA010A5367D')?.amount ?? '';

  const onSubmitClick = async () => {
    setIsLoading(true);

    const coin = getCoin(fromChain ?? '');
    const transferAmount = new BigNumber(amount).shiftedBy(6).toString();

    const fee: StdFee = {
      amount: coins('1000', coin.base),
      gas: '300000',
    };

    const sourcePort = 'transfer';
    const sourceChannel = 'channel-0';

    const token = {
      denom: 'ibc/635CB83EF1DFE598B10A3E90485306FD0D47D34217A4BE5FD9977FA010A5367D',
      amount: transferAmount,
    };

    const stamp = Date.now();
    const timeoutInNanos = (stamp + 1.2e6) * 1e6;

    const msg = transfer({
      sourcePort,
      sourceChannel,
      sender: address ?? '',
      receiver: qAddress ?? '',
      token,
      timeoutHeight: undefined,
      //@ts-ignore
      timeoutTimestamp: timeoutInNanos,
      memo: '',
    });

    await tx([msg], {
      fee,
      onSuccess: () => {
        setAmount('');
      },
    });

    setIsLoading(false);
  };

  return (
    <>
      <Button
        _active={{
          transform: 'scale(0.95)',
          color: 'complimentary.800',
        }}
        _hover={{
          bgColor: 'rgba(255,128,0, 0.25)',
          color: 'complimentary.300',
        }}
        color="white"
        flex={1}
        size="sm"
        variant="outline"
        onClick={onOpen}
      >
        Withdraw
      </Button>

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent bgColor="rgb(32,32,32)">
          <ModalHeader color="white">Withdraw {token} Tokens</ModalHeader>
          <ModalCloseButton color={'complimentary.900'} />
          <ModalBody>
            {/* Chain Selection Dropdown */}
            <FormControl>
              <FormLabel color={'white'}>To Chain</FormLabel>
              {chooseChain}
            </FormControl>

            {/* Amount Input */}
            <FormControl mt={4}>
              <FormLabel color="white">Amount</FormLabel>
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
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                color={'white'}
                placeholder="Enter amount"
              />
            </FormControl>
          </ModalBody>

          <ModalFooter>
            <Button
              _active={{
                transform: 'scale(0.95)',
                color: 'complimentary.800',
              }}
              _hover={{
                bgColor: 'rgba(255,128,0, 0.25)',
                color: 'complimentary.300',
              }}
              minW="100px"
              mr={3}
              onClick={onSubmitClick}
              disabled={!amount}
            >
              {isLoading === true && <Spinner size="sm" />}
              {isLoading === false && 'Withdraw'}
            </Button>
            <Button
              _active={{
                transform: 'scale(0.95)',
                color: 'complimentary.800',
              }}
              _hover={{
                bgColor: 'rgba(255,128,0, 0.25)',
                color: 'complimentary.300',
              }}
              color="white"
              variant="ghost"
              onClick={onClose}
            >
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};

export default QWithdrawModal;
