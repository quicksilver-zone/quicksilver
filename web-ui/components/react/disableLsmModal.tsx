import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalFooter, ModalBody, ModalCloseButton, Button, Text } from '@chakra-ui/react';

interface DisableLsmModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const DisableLsmModal: React.FC<DisableLsmModalProps> = ({ isOpen, onClose }) => (
  <Modal isOpen={isOpen} onClose={onClose}>
    <ModalOverlay />
    <ModalContent>
      <ModalHeader>Disable LSM</ModalHeader>
      <ModalCloseButton />
      <ModalBody>
        <Text>Are you sure you want to disable LSM?</Text>
      </ModalBody>
      <ModalFooter>
        <Button colorScheme="blue" mr={3} onClick={onClose}>
          Close
        </Button>
        <Button variant="ghost">Confirm</Button>
      </ModalFooter>
    </ModalContent>
  </Modal>
);
