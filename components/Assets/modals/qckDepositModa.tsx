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
} from '@chakra-ui/react';
import { ChainName } from '@cosmos-kit/core';
import { useManager } from '@cosmos-kit/react';
import { useState, useMemo, useEffect } from 'react';

import { ChooseChain } from '@/components/react/choose-chain';
import { handleSelectChainDropdown, ChainOption } from '@/components/types';

export function DepositModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();

  const [chainName, setChainName] = useState<ChainName | undefined>('akash');
  const { chainRecords, getChainLogo } = useManager();

  const chainOptions = useMemo(
    () =>
      chainRecords.map((chainRecord) => {
        return {
          chainName: chainRecord?.name,
          label: chainRecord?.chain.pretty_name,
          value: chainRecord?.name,
          icon: getChainLogo(chainRecord.name),
        };
      }),
    [chainRecords, getChainLogo],
  );

  useEffect(() => {
    setChainName(window.localStorage.getItem('selected-chain') || 'akash');
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

  const handleDeposit = () => {
    // Implement deposit logic here
    // Show toast on success or error
    toast({
      title: 'Deposit Successful',
      status: 'success',
      duration: 9000,
      isClosable: true,
    });
    onClose(); // Close modal after deposit
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
          <ModalHeader color="white">Deposit QCK Tokens</ModalHeader>
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
              <Input color={'white'} placeholder="Enter amount" />
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
              onClick={handleDeposit}
            >
              Deposit
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
