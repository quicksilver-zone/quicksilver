
import React, { useEffect, useState } from 'react';
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalBody,
    ModalCloseButton,
    Button,
    useDisclosure,
    Link,
    useToast,
    Text,
    SimpleGrid,
    Avatar
} from '@chakra-ui/react'
import { SupportedWallets, WalletConfigs } from '@/state/config';
import { useSelector, useDispatch } from 'react-redux';
import connectToWallet from '@/state/wallet/thunks/connectWallet';
export default function WalletModal() {
    const dispatch = useDispatch()
    const [isKeplrInstalled, setIsKeplrInstalled] = useState(false);
    const [isLeapInstalled, setIsLeapInstalled] = useState(false);
    const [isCosmostationInstalled, setIsCosmostationInstalled] = useState(false);

    const { connecting } = useSelector(state => state.wallet)
    const { isOpen, onOpen, onClose } = useDisclosure()
    const toast = useToast();

    useEffect(() => {
        if (window && window.keplr) {
            setIsKeplrInstalled(true);
        }

        if (window && window.leap) {
            setIsLeapInstalled(true);
        }

        if (window && window.cosmostation) {
            setIsCosmostationInstalled(true);
        }
    }, []);

    const handleConnectWallet = async (walletType) => {
        onClose()
        let isInstallWallet = false
        switch(walletType) {
            case 'keplr':
                isInstallWallet = isKeplrInstalled
                break
            case 'leap':
                isInstallWallet = isLeapInstalled
                break
            case 'cosmostation':
                isInstallWallet = isCosmostationInstalled
                break
        }
        if (!isInstallWallet) {
            toast({
                title: 'Wallet is not installed',
                description: 'Please install ' +  walletType + ' wallet',
                status: 'error',
                duration: 3000,
                isClosable: true
            })
        } else {
            dispatch(connectToWallet(walletType))
        }
    }
    return (
        <>
            <Button
                onClick={onOpen}
                bgColor="#E77728"
                isLoading={connecting}
                fontWeight={500}
                borderRadius={'4px'}
                h='40px'
                fontSize='16px'
                _hover={{
                    opacity: 0.8
                }}
            >
                Connect wallet
            </Button>

            <Modal isOpen={isOpen} onClose={onClose} size='md'>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader bgColor='#0E0E0E' borderTopRadius={'md'}>
                        <Text color='white' fontSize='22px'>
                            Connect Your Wallet
                        </Text>
                        <Text color='#CDCDCD' fontSize='14px' fontWeight={400} mt={2}>
                            Select your preferred wallet below.
                        </Text>
                        <Text color='#CDCDCD' fontSize='14px' fontWeight={400} mt={2}>
                            Don't have a wallet?{' '}
                        <Link color='blue.600' textDecoration='underline' href='#' fontWeight={600}>
                            See supported wallets
                        </Link>
                        </Text>
                    </ModalHeader>
                    <ModalCloseButton />
                    <ModalBody 
                        bgColor={'#181818'} borderBottomRadius={'md'}
                        py={5}
                    >
                        <SimpleGrid columns={2} gap={4}>
                            {SupportedWallets.map(item =>
                                <Button
                                    onClick={() => handleConnectWallet(item)}
                                    isDisabled={connecting}
                                    key={item}
                                    w='full'
                                    size='lg'
                                    justifyContent='left'
                                    colorScheme='blackAlpha'
                                    py={4}
                                    fontSize='16px'
                                    fontWeight={500}
                                    border='solid 1px gray'
                                    leftIcon={<Avatar src={WalletConfigs[item].logo} w='30px' h='30px' />}>
                                    {WalletConfigs[item].name}
                                </Button>)
                            }
                        </SimpleGrid>
                    </ModalBody>
                </ModalContent>
            </Modal>

        </>
    );

}
