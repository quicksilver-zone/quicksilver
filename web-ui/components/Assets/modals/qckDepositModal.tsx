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
  useDisclosure,
  Divider,
  Text,
  Spinner,
} from '@chakra-ui/react';
import { ChainName } from '@cosmos-kit/core';
import { useChain, useManager } from '@cosmos-kit/react';
import BigNumber from 'bignumber.js';
import { ibc } from 'interchain-query';
import { useState, useMemo, useEffect } from 'react';

import { ChooseChain } from '@/components/react/choose-chain';
import { handleSelectChainDropdown, ChainOption, ChooseChainInfo } from '@/components/types';
import { useTx } from '@/hooks';
import { useFeeEstimation } from '@/hooks/useFeeEstimation';
import { getIbcInfo } from '@/utils';

export function DepositModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  const [chainName, setChainName] = useState<ChainName | undefined>('osmosis');
  const { chainRecords, getChainLogo } = useManager();
  const [amount, setAmount] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const chainOptions = useMemo(() => {
    return chainRecords
      .filter((chainRecord) => chainRecord.name === 'osmosis')
      .map((chainRecord) => ({
        chainName: chainRecord?.name,
        label: chainRecord?.chain?.pretty_name,
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

  const chooseChain = <ChooseChain chainName={chainName} chainInfos={chainOptions as ChooseChainInfo[]} onChange={onChainChange} />;

  const fromChain = chainName;
  const toChain = 'quicksilver';

  const { transfer } = ibc.applications.transfer.v1.MessageComposer.withTypeUrl;
  const { address } = useChain(fromChain ?? '');
  const { address: qAddress } = useChain('quicksilver');
  const { tx } = useTx(fromChain ?? '');
  const { estimateFee } = useFeeEstimation(fromChain ?? '');

  const onSubmitClick = async () => {
    setIsLoading(true);

    const transferAmount = new BigNumber(amount).shiftedBy(6).toString();

    const { source_port, source_channel } = getIbcInfo(fromChain ?? '', toChain ?? '');

    const token = {
      denom: 'ibc/635CB83EF1DFE598B10A3E90485306FD0D47D34217A4BE5FD9977FA010A5367D',
      amount: transferAmount,
    };

    const stamp = Date.now();
    const timeoutInNanos = (stamp + 1.2e6) * 1e6;

    const msg = transfer({
      sourcePort: source_port,
      sourceChannel: source_channel,
      sender: address ?? '',
      receiver: qAddress ?? '',
      token,
      timeoutHeight: undefined,
      //@ts-ignore
      timeoutTimestamp: timeoutInNanos,
    });
    const fee = await estimateFee(address ?? '', [msg]);
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
        Deposit
      </Button>

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent bgColor="rgb(32,32,32)">
          <ModalHeader color="white"><Text>Deposit QCK Tokens</Text>  <Divider mt={3} bgColor={'cyan.500'} /></ModalHeader>
          <ModalCloseButton color={'complimentary.900'} />
          <ModalBody>
            {/* Chain Selection Dropdown */}
            <FormControl>
              <FormLabel color={'white'}>From Chain</FormLabel>
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
              isDisabled={!amount}
            >
              {isLoading === true && <Spinner size="sm" />}
              {isLoading === false && 'Deposit'}
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
