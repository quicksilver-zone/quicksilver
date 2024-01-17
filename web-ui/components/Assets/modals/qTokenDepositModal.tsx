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
  useToast,
  Spinner,
} from '@chakra-ui/react';
import { ibc } from '@chalabi/quicksilverjs';
import { StdFee, coins } from '@cosmjs/stargate';
import { ChainName } from '@cosmos-kit/core';
import { useChain, useManager } from '@cosmos-kit/react';
import BigNumber from 'bignumber.js';
import { useState, useMemo, useEffect } from 'react';

import { ChooseChain } from '@/components/react/choose-chain';
import { handleSelectChainDropdown, ChainOption } from '@/components/types';
import { useTx } from '@/hooks';
import { useIbcBalanceQuery } from '@/hooks/useQueries';
import { ibcDenomDepositMapping } from '@/state/chains/prod';
import { getCoin, getIbcInfo } from '@/utils';

export interface QDepositModalProps {
  token: string;
}

const QDepositModal: React.FC<QDepositModalProps> = ({ token }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();

  const [chainName, setChainName] = useState<ChainName | undefined>('osmosis');
  const { chainRecords, getChainLogo } = useManager();
  const [amount, setAmount] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const chainOptions = useMemo(() => {
    const desiredChains = ['osmosis', 'umee'];
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

    const { sourcePort, sourceChannel } = getIbcInfo(fromChain ?? '', toChain ?? '');

    // Function to get the correct IBC denom trace based on chain and token
    type ChainDenomMappingKeys = keyof typeof ibcDenomDepositMapping;

    type TokenKeys = keyof (typeof ibcDenomDepositMapping)['osmosis'];

    const getIbcDenom = (chainName: string, token: string) => {
      const chain = chainName as ChainDenomMappingKeys;
      const chainDenoms = ibcDenomDepositMapping[chain];

      if (token in chainDenoms) {
        return chainDenoms[token as TokenKeys];
      }

      return undefined;
    };

    const ibcDenom = getIbcDenom(fromChain ?? '', token);
    if (!ibcDenom) {
      toast({
        title: 'Error',
        description: `No IBC denom trace found for ${token} on chain ${fromChain}`,
        status: 'error',
        duration: 9000,
        isClosable: true,
      });
      setIsLoading(false);
      return;
    }

    const ibcToken = {
      denom: ibcDenom ?? '',
      amount: transferAmount,
    };

    const stamp = Date.now();
    const timeoutInNanos = (stamp + 1.2e6) * 1e6;

    const msg = transfer({
      sourcePort,
      sourceChannel,
      sender: address ?? '',
      receiver: qAddress ?? '',
      token: ibcToken,
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
        Deposit
      </Button>

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent bgColor="rgb(32,32,32)">
          <ModalHeader color="white">Deposit {token} Tokens</ModalHeader>
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
};

export default QDepositModal;
