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
import { color } from 'framer-motion';
import { useState, useMemo, useEffect } from 'react';

import { ChooseChain } from '@/components/react/choose-chain';
import { handleSelectChainDropdown, ChainOption } from '@/components/types';
import { useTx } from '@/hooks';
import { useIbcBalanceQuery } from '@/hooks/useQueries';
import { getCoin } from '@/utils';
import { StdFee, coins } from '@cosmjs/stargate';
import { ibc } from 'interchain-query';
import BigNumber from 'bignumber.js';

export function WithdrawModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();

  const [chainName, setChainName] = useState<ChainName | undefined>('osmosis');
  const { chainRecords, getChainLogo } = useManager();
  const [amount, setAmount] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const chainOptions = useMemo(() => {
    return chainRecords
      .filter((chainRecord) => chainRecord.name === 'osmosis')
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

  const fromChain = 'quicksilver';
  const toChain = chainName;

  const { transfer } = ibc.applications.transfer.v1.MessageComposer.withTypeUrl;
  const { address, connect, status, message, wallet } = useChain(toChain ?? '');
  const { address: qAddress } = useChain('quicksilver');
  const { balance } = useIbcBalanceQuery(fromChain ?? '', qAddress ?? '');
  const { tx } = useTx(fromChain ?? '');

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
      denom: 'uqck',
      amount: transferAmount,
    };

    const stamp = Date.now();
    const timeoutInNanos = (stamp + 1.2e6) * 1e6;

    const msg = transfer({
      sourcePort,
      sourceChannel,
      sender: qAddress ?? '',
      receiver: address ?? '',
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
        onClick={onOpen}
        _hover={{
          bgColor: 'rgba(255,128,0, 0.25)',
          color: 'complimentary.300',
        }}
        color={'white'}
        w="full"
        variant="outline"
      >
        Withdraw
      </Button>

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent bgColor="rgb(32,32,32)">
          <ModalHeader color="white">Withdraw QCK Tokens</ModalHeader>
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
              mr={3}
              minW="100px"
              onClick={onSubmitClick}
              disabled={Number.isNaN(Number(amount))}
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
}
