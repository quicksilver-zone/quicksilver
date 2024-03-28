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
  HStack,
  useDisclosure,
  useToast,
  Spinner,
  InputGroup,
  InputRightElement,
} from '@chakra-ui/react';
import { StdFee, coins } from '@cosmjs/stargate';
import { ChainName } from '@cosmos-kit/core';
import { useChain, useManager } from '@cosmos-kit/react';
import BigNumber from 'bignumber.js';
import { ibc } from 'interchain-query';
import { useState, useMemo, useEffect } from 'react';

import { ChooseChain } from '@/components/react/choose-chain';
import { handleSelectChainDropdown, ChainOption, ChooseChainInfo } from '@/components/types';
import { useTx } from '@/hooks';
import { useIbcBalanceQuery } from '@/hooks/useQueries';
import { ibcDenomDepositMapping } from '@/state/chains/prod';
import { getCoin, getIbcInfo } from '@/utils';

export interface QDepositModalProps {
  token: string;
  isOpen: boolean;
  onClose: () => void;
  interchainDetails: { [chainId: string]: number };
}

const QDepositModal: React.FC<QDepositModalProps> = ({ token, isOpen, onClose, interchainDetails }) => {
  const toast = useToast();
  const { chainRecords, getChainLogo } = useManager();
  const [chainName, setChainName] = useState<ChainName | undefined>('osmosis');
  const [amount, setAmount] = useState<string>('');
  const [maxAmount, setMaxAmount] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const chainOptions = useMemo(() => {
    const availableChains = Object.keys(interchainDetails);
    return chainRecords
      .filter((chainRecord) => availableChains.includes(chainRecord.name))
      .map((chainRecord) => ({
        chainName: chainRecord.name,
        label: chainRecord?.chain?.pretty_name || chainRecord.name,
        value: chainRecord.name,
        icon: getChainLogo(chainRecord.name),
      }));
  }, [chainRecords, getChainLogo, interchainDetails]);

  useEffect(() => {
    const storedChainName = window.localStorage.getItem('selected-chain');
    const defaultChainName = chainOptions[0]?.chainName || 'osmosis';
    const initialChainName = storedChainName || defaultChainName;
    setChainName(initialChainName);
    setMaxAmount(interchainDetails[initialChainName]?.toString() || '0');
  }, [chainOptions, interchainDetails]);

  const onChainChange: handleSelectChainDropdown = (selectedValue: ChainOption | null) => {
    if (selectedValue?.chainName) {
      setChainName(selectedValue.chainName);
      setMaxAmount(interchainDetails[selectedValue.chainName]?.toString() || '0');
      window.localStorage.setItem('selected-chain', selectedValue.chainName);
    }
  };

  const chooseChain = <ChooseChain chainName={chainName} chainInfos={chainOptions as ChooseChainInfo[]} onChange={onChainChange} />;

  const fromChain = chainName;
  const toChain = 'quicksilver';

  const { transfer } = ibc.applications.transfer.v1.MessageComposer.withTypeUrl;
  const { address, connect, status, message, wallet } = useChain(fromChain ?? '');
  const { address: qAddress } = useChain('quicksilver');

  const { tx } = useTx(fromChain ?? '');

  const onSubmitClick = async () => {
    setIsLoading(true);

    const coin = getCoin(fromChain ?? '');
    const transferAmount = new BigNumber(amount).shiftedBy(6).toString();

    const fee: StdFee = {
      amount: coins('1000', coin.base),
      gas: '300000',
    };

    const { source_port, source_channel } = getIbcInfo(fromChain ?? '', toChain ?? '');

    // Function to get the correct IBC denom trace based on chain and token
    type ChainDenomMappingKeys = keyof typeof ibcDenomDepositMapping;

    type TokenKeys = keyof (typeof ibcDenomDepositMapping)['osmosis'];

    const getIbcDenom = (chainName: string, token: string) => {
      const chain = chainName as ChainDenomMappingKeys;
      const chainDenoms = ibcDenomDepositMapping[chain];

      if (chainDenoms && token in chainDenoms) {
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
      sourcePort: source_port,
      sourceChannel: source_channel,
      sender: address ?? '',
      receiver: qAddress ?? '',
      token: ibcToken,
      timeoutHeight: undefined,
      //@ts-ignore
      timeoutTimestamp: timeoutInNanos,
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
            <InputGroup>
              <Input
                type="number"
                pr="4.5rem" // Padding to ensure text doesn't overlap with buttons
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
                onChange={(e) => setAmount(e.target.value <= maxAmount ? e.target.value : maxAmount)}
                max={maxAmount}
                color={'white'}
                placeholder="Enter amount"
              />
              <InputRightElement>
                <HStack mr={14} spacing={1}>
                  <Button
                    variant={'ghost'}
                    color="complimentary.900"
                    h="1.75rem"
                    size="xs"
                    _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
                    _hover={{ bgColor: 'transparent', color: 'complimentary.400' }}
                    onClick={() => setAmount((parseFloat(maxAmount) / 2).toString())}
                  >
                    Half
                  </Button>
                  <Button
                    variant={'ghost'}
                    color="complimentary.900"
                    _active={{ transform: 'scale(0.95)', color: 'complimentary.800' }}
                    _hover={{ bgColor: 'transparent', color: 'complimentary.400' }}
                    h="1.75rem"
                    size="xs"
                    onClick={() => setAmount(maxAmount)}
                  >
                    Max
                  </Button>
                </HStack>
              </InputRightElement>
            </InputGroup>
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
            isDisabled={!amount || !address}
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
  );
};

export default QDepositModal;
