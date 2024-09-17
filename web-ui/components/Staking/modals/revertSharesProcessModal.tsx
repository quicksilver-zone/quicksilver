import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalBody,
  ModalCloseButton,
  HStack,
  Text,
  Box,
  Circle,
  Flex,
  Button,
  Stat,
  StatLabel,
  StatNumber,
  Spinner,
} from '@chakra-ui/react';
import { StdFee } from '@cosmjs/amino';
import styled from '@emotion/styled';
import chains from 'chain-registry';
import { assets } from 'chain-registry';
import { cosmos } from 'quicksilverjs';
import React, { useEffect, useState } from 'react';

import { useTx } from '@/hooks';
import { useFeeEstimation } from '@/hooks/useFeeEstimation';
import { shiftDigits } from '@/utils';
import { Chain } from '@/config';

const ChakraModalContent = styled(ModalContent)`
  position: relative;
  background: none;
  max-height: 450px;
  &::before,
  &::after {
    z-index: -1;
  }
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    width: 40%;
    background-color: #201c18;
    border-radius: 5px 0 0 5px;
  }
  &::after {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    width: 60%;
    background-color: #1a1a1a;
    border-radius: 0 5px 5px 0;
  }
`;

interface SelectedValidator {
  operatorAddress: string;
  moniker: string;
  tokenAmount: string;
}

interface StakingModalProps {
  isOpen: boolean;
  onClose: () => void;
  children?: React.ReactNode;
  selectedValidator: SelectedValidator;
  selectedOption?: Chain;
  address: string;
  isTokenized: boolean;
  denom: string;
  refetch: () => void;
}

export const RevertSharesProcessModal: React.FC<StakingModalProps> = ({
  isOpen,
  onClose,
  selectedOption,
  selectedValidator,
  address,
  denom,
  refetch,
}) => {
  const [step, setStep] = useState(1);
  const getProgressColor = (circleStep: number) => {
    if (step >= circleStep) return 'complimentary.900';
    return 'rgba(255,255,255,0.2)';
  };

  const [isSigning, setIsSigning] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  const newChainName = selectedOption?.chain_name;


  const labels = ['Revert Shares', `Receive Tokens`];

  const { tx, responseEvents } = useTx(newChainName ?? '');
  const { estimateFee } = useFeeEstimation(newChainName ?? '');
  const [combinedDenom, setCombinedDenom] = useState<string>();

  // prettier-ignore
  useEffect(() => {
  
      const tokenizeSharesEvent = responseEvents?.find(event => event.type === 'tokenize_shares');
    
      if (tokenizeSharesEvent) {
   
        const validatorValue = tokenizeSharesEvent.attributes.find(attr => attr.key === 'validator')?.value;
        const shareRecordIdValue = tokenizeSharesEvent.attributes.find(attr => attr.key === 'share_record_id')?.value;
    
  
        if (validatorValue && shareRecordIdValue) {
          setCombinedDenom(`${validatorValue}/${shareRecordIdValue}`);
        }
      }
    }, [responseEvents]);

  const { redeemTokensForShares } = cosmos.staking.v1beta1.MessageComposer.withTypeUrl;

  const msg = redeemTokensForShares({
    delegatorAddress: address,
    amount: {
      denom: denom ?? combinedDenom,
      amount: selectedValidator.tokenAmount.toString(),
    },
  });

  const handleRevertShares = async (event: React.MouseEvent) => {
    event.preventDefault();
    setIsSigning(true);

    const fee = await estimateFee(address, [msg]);

    try {
      await tx([msg], {
        fee,
        onSuccess: () => {
          refetch();
          setStep(2);
        },
      });
    } catch (error) {
      console.error('Transaction failed', error);

      setIsError(true);
    } finally {
      setIsSigning(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} size={{ base: '3xl', md: '2xl' }}>
      <ModalOverlay />
      <ChakraModalContent h={{ md: '30%', base: '35%' }} maxH={'100%'}>
        <ModalBody borderRadius={4} h="30%" maxH={'100%'}>
          <ModalCloseButton zIndex={1000} color="white" />
          <HStack position={'relative'} h="100%" spacing="48px" align="stretch">
            {/* Left Section */}
            <Flex flexDirection="column" justifyContent="space-between" width="40%" p={4} bg="#1E1C19" height="100%">
              <Box position="relative">
                <Stat>
                  <StatLabel color="rgba(255,255,255,0.5)">REVERT</StatLabel>

                  <StatNumber display={{ base: 'none', md: 'block' }} color="white">
                    {shiftDigits(selectedValidator.tokenAmount, -6)}&nbsp;
                    {selectedOption?.big_denom}
                  </StatNumber>
                </Stat>
                {[1, 2].map((circleStep, index) => (
                  <Flex key={circleStep} align="center" mt={10} mb={circleStep !== 4 ? '48px' : '0'}>
                    <Circle
                      size="36px"
                      bg={getProgressColor(circleStep)}
                      color="white"
                      fontWeight="bold"
                      borderWidth={'2px'}
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                      position="relative"
                      borderColor="rgba(255,255,255,0.5)"
                    >
                      {circleStep}
                      {circleStep !== 2 && (
                        <>
                          <Box
                            width="2px"
                            height="30px"
                            bgColor="rgba(255,255,255,0.01)"
                            position="absolute"
                            bottom="-42px"
                            left="50%"
                            transform="translateX(-50%)"
                          />
                          <Box
                            width="2px"
                            height="30px"
                            bgColor={getProgressColor(circleStep + 1)}
                            position="absolute"
                            bottom="-42px"
                            left="50%"
                            transform="translateX(-50%)"
                          />
                        </>
                      )}
                    </Circle>
                    <Text fontWeight="hairline" ml={3} color="rgba(255,255,255,0.75)">
                      {labels[index]}
                    </Text>
                  </Flex>
                ))}
              </Box>
            </Flex>

            <Flex width="67%" flexDirection="column" justifyContent="center" alignItems="center">
              {step === 1 && (
                <>
                  <Flex maxW="300px" flexDirection={'column'} justifyContent={'flex-start'} alignItems={'center'}>
                    <Text textAlign={'center'} fontWeight={'bold'} fontSize="lg" color="white">
                      You are about to revert your shares back to tokens.
                    </Text>
                    <Text mt={2} textAlign={'left'} fontWeight={'light'} fontSize="lg" color="white">
                      Reverting&nbsp;&nbsp;{shiftDigits(selectedValidator.tokenAmount, -(selectedOption?.exponent ?? 6))}&nbsp; {selectedOption?.big_denom}
                    </Text>
                  </Flex>

                  <Button
                    mt={4}
                    width={{ base: '80%', md: '30%' }}
                    _active={{
                      transform: 'scale(0.95)',
                      color: 'complimentary.800',
                    }}
                    _hover={{
                      bgColor: 'rgba(255,128,0, 0.25)',
                      color: 'complimentary.300',
                    }}
                    onClick={handleRevertShares}
                  >
                    {isError ? 'Try Again' : isSigning ? <Spinner /> : 'Revert'}
                  </Button>
                </>
              )}
              {step === 2 && (
                <>
                  <Text textAlign={'center'} fontWeight={'bold'} fontSize="lg" color="white">
                    Your shares have been successfully reverted back to tokens and should arrive in your wallet.
                  </Text>
                </>
              )}
            </Flex>
          </HStack>
        </ModalBody>
      </ChakraModalContent>
    </Modal>
  );
};
export default RevertSharesProcessModal;
