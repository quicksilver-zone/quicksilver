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
  Divider,
  FormControl,
  FormLabel,
  Input,
  useToast,
  Spinner,
  HStack,
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
import { useFeeEstimation } from '@/hooks/useFeeEstimation';
import { useIbcBalanceQuery } from '@/hooks/useQueries';
import { ibcDenomWithdrawMapping } from '@/state/chains/prod';
import { getCoin, getExponent, getIbcInfo } from '@/utils';

interface QDepositModalProps {
  max: string;
  token: string;
  isOpen: boolean;
  onClose: () => void;
  refetch: () => void;
}

const QWithdrawModal: React.FC<QDepositModalProps> = ({ max, token, isOpen, onClose, refetch }) => {
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

  const fromChain = 'quicksilver';
  const toChain = chainName;

  const { transfer } = ibc.applications.transfer.v1.MessageComposer.withTypeUrl;
  const { address } = useChain(toChain ?? '');
  const { address: qAddress } = useChain('quicksilver');

  const { tx } = useTx(fromChain ?? '');
  const { estimateFee } = useFeeEstimation(fromChain ?? '');

  const onSubmitClick = async () => {
    setIsLoading(true);
    const exp = token === 'qDYDX' ? 18 : 6;

    const transferAmount = new BigNumber(amount).shiftedBy(exp).toString();

    const { source_port, source_channel } = getIbcInfo(fromChain ?? '', toChain ?? '');

    // Function to get the correct IBC denom trace based on chain and token
    type ChainDenomMappingKeys = keyof typeof ibcDenomWithdrawMapping;

    type TokenKeys = keyof (typeof ibcDenomWithdrawMapping)['quicksilver'];

    const getIbcDenom = (chainName: string, token: string) => {
      const chain = chainName as ChainDenomMappingKeys;
      const chainDenoms = ibcDenomWithdrawMapping[chain];

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

    const qckDenom = token === 'qDYDX' ? 'a' + ibcDenom : 'u' + ibcDenom;
    const ibcToken = {
      denom: qckDenom ?? '',
      amount: transferAmount,
    };

    const stamp = Date.now();
    const timeoutInNanos = (stamp + 1.2e6) * 1e6;

    const msg = transfer({
      sourcePort: source_port,
      sourceChannel: source_channel,
      sender: qAddress ?? '',
      receiver: address ?? '',
      token: ibcToken,
      timeoutHeight: undefined,
      //@ts-ignore
      timeoutTimestamp: timeoutInNanos,
    });

    const fee = await estimateFee(qAddress ?? '', [msg]);

    await tx([msg], {
      fee,
      onSuccess: () => {
        setAmount('');
        refetch();
      },
    });

    setIsLoading(false);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent bgColor="rgb(32,32,32)">
        <ModalHeader color="white">
          <Text>Withdraw {token} Tokens</Text> <Divider mt={3} bgColor={'cyan.500'} />
        </ModalHeader>
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
                onChange={(e) => setAmount(e.target.value <= max ? e.target.value : BigNumber(max).toString())}
                max={max}
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
                    onClick={() =>
                      setAmount(
                        BigNumber(parseFloat(max) / 2)
                          .toFixed(6)
                          .toString(),
                      )
                    }
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
                    onClick={() => setAmount(BigNumber(max).toFixed(6).toString())}
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
  );
};

export default QWithdrawModal;
