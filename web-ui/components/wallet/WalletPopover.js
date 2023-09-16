
import React , {useEffect, useState} from 'react';
import {
    Button,
    Image,
    useDisclosure,
    Popover,
    PopoverTrigger,
    PopoverContent,
    useToast,
    PopoverBody,
    Text,
    VStack
} from '@chakra-ui/react'
import { shortenAddress } from '@/helper/address';
import { WalletConfigs } from '@/state/config';
import { CopyIcon } from '@chakra-ui/icons';
import { useSelector, useDispatch } from 'react-redux';
import { disconnectWallet } from '@/state/wallet/slice';
import connectToWallet from '@/state/wallet/thunks/connectWallet';

export default function WalletPopover() {
    const dispatch = useDispatch()
    const { address, typeWallet } = useSelector(state => state.wallet)

    const toast = useToast()
    const { isOpen, onToggle, onClose } = useDisclosure()

    const handleCopy = () => {
        navigator.clipboard.writeText(address)
        toast({
            title: 'Copied',
            status: 'success',
            duration: 2000,
            isClosable: true
        })
    }

    const handleDisconnectWallet = () => {
        onClose()
        dispatch(disconnectWallet())
        if (window) {
            localStorage.removeItem('WalletType')
        }
    }

    useEffect(() => {
        if (typeWallet && window) {
            window.addEventListener("keplr_keystorechange", () => {
                dispatch(connectToWallet(typeWallet))
            })
        }
    }, [typeWallet])
    return (
        <Popover
            isOpen={isOpen}
            onClose={onClose}
        >
            <PopoverTrigger>
                <Button
                    minW='102px'
                    onClick={onToggle}
                    colorScheme='blackAlpha'
                    fontWeight={500}
                    borderRadius={'4px'}
                    h='40px'
                    fontSize='16px'
                    _hover={{
                        opacity: 0.8
                    }}
                    leftIcon={<Image src={WalletConfigs[typeWallet]?.logo} w='20px' h='20px' />}
                >
                    {shortenAddress(address)}
                </Button>
            </PopoverTrigger>
            <PopoverContent bgColor={'#202020'}>
                <PopoverBody>
                    <VStack gap={8}>
                        <Button
                            onClick={handleCopy}
                            justifyContent={'space-between'}
                            variant='ghost'
                            leftIcon={<Image mr={3} src={WalletConfigs[typeWallet]?.logo} w='40px' h='40px' />}
                            rightIcon={<CopyIcon alignSelf={'end'} />}
                            colorScheme='whiteAlpha'
                            size='lg'
                            w='full'
                            px={0}
                            py={4}
                        >
                            <VStack alignItems={'start'} gap={1} w='full'>
                                <Text fontSize={'14px'}>{shortenAddress(address)}</Text>
                            </VStack>
                        </Button>
                        <Button
                            bgColor="#E77728"
                            onClick={handleDisconnectWallet}
                            w='full'
                            fontSize={'16px'}
                            fontWeight={400}
                            borderRadius={'4px'}
                            _hover={{
                                opacity: 0.8
                            }}
                        >Disconnect Wallet</Button>
                    </VStack>

                </PopoverBody>
            </PopoverContent>
        </Popover>

    );

}
