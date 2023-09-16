
import React, { useEffect } from 'react';
import { Box } from '@chakra-ui/react';
import WalletModal from '../wallet/WalletModal';
import WalletPopover from '../wallet/WalletPopover';
import { useSelector, useDispatch } from 'react-redux';
import connectToWallet from '@/state/wallet/thunks/connectWallet';
import { SupportedWallets } from '@/state/config';
export default function Header() {
    const dispatch = useDispatch()
    const {connected, connecting} = useSelector(state => state.wallet)

    useEffect(() => {
        if (window && !connected && !connecting) {
            let oldWalletType = localStorage.getItem('WalletType');
            if (oldWalletType && SupportedWallets.includes(oldWalletType)) {
                dispatch(connectToWallet(oldWalletType))
            }
        }
    }, [])

    return (
        <Box position={'fixed'} top={10} right={10} zIndex={10}>
            {connected ? <WalletPopover /> :
                <WalletModal />}
        </Box>
    )
}