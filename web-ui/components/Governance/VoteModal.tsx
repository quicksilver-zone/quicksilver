import {
  Modal,
  Text,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
  Stack,
  Radio,
  RadioGroup,
  UseDisclosureReturn,
} from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import { cosmos } from 'interchain-query';
import { useState } from 'react';

import { useTx } from '@/hooks';
import { useFeeEstimation } from '@/hooks/useFeeEstimation';

const VoteType = cosmos.gov.v1beta1.VoteOption;
const { vote: composeVoteMessage } = cosmos.gov.v1beta1.MessageComposer.fromJSON;

interface VoteModalProps {
  modalControl: UseDisclosureReturn;
  chainName: string;
  updateVotes: () => void;
  vote: number | undefined;
  title: string;
  proposalId: bigint;
}

export const VoteModal: React.FC<VoteModalProps> = ({ modalControl, chainName, updateVotes, title, vote, proposalId }) => {
  const [option, setOption] = useState<number>();
  const [isLoading, setIsLoading] = useState(false);

  const { tx } = useTx(chainName);
  const { estimateFee } = useFeeEstimation(chainName);
  const { address } = useChain(chainName);

  const { isOpen, onClose } = modalControl;

  const checkIfDisable = (option: number) => option === vote;

  const closeModal = () => {
    onClose();
    setOption(undefined);
  };

  const handleConfirmClick = async () => {
    if (!address || !option) return;
    setIsLoading(true);

    const msg = composeVoteMessage({
      "option": option,
      "proposalId": proposalId.toString(),
      "voter": address,
    });

    console.log(msg)

    const fee = await estimateFee(address, [msg]);

    await tx([msg], {
      fee,
      onSuccess: () => {
        updateVotes();
        closeModal();
      },
    });

    setIsLoading(false);
  };

  return (
    <Modal isOpen={isOpen} onClose={closeModal} isCentered>
      <ModalOverlay />
      <>
        <ModalContent bgColor="#1A1A1A">
          <ModalHeader color="white" mr={4}>
            {title}
          </ModalHeader>
          <ModalCloseButton color="white" />
          <ModalBody>
            <RadioGroup onChange={(e) => setOption(Number(e))}>
              <Stack>
                <Radio
                  colorScheme="green"
                  size="lg"
                  value={VoteType.VOTE_OPTION_YES.toString()}
                  isDisabled={checkIfDisable(VoteType.VOTE_OPTION_YES)}
                >
                  <Text>Yes</Text>
                </Radio>
                <Radio
                  colorScheme="red"
                  size="lg"
                  value={VoteType.VOTE_OPTION_NO.toString()}
                  isDisabled={checkIfDisable(VoteType.VOTE_OPTION_NO)}
                >
                  <Text>No</Text>
                </Radio>
                <Radio
                  colorScheme="red"
                  size="lg"
                  value={VoteType.VOTE_OPTION_NO_WITH_VETO.toString()}
                  isDisabled={checkIfDisable(VoteType.VOTE_OPTION_NO_WITH_VETO)}
                >
                  <Text>No With Veto</Text>
                </Radio>
                <Radio
                  colorScheme="gray"
                  size="lg"
                  value={VoteType.VOTE_OPTION_ABSTAIN.toString()}
                  isDisabled={checkIfDisable(VoteType.VOTE_OPTION_ABSTAIN)}
                >
                  <Text>Abstain</Text>
                </Radio>
              </Stack>
            </RadioGroup>
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
              onClick={handleConfirmClick}
              isDisabled={!option || isLoading}
              isLoading={isLoading}
            >
              Confirm
            </Button>
          </ModalFooter>
        </ModalContent>
      </>
    </Modal>
  );
};
