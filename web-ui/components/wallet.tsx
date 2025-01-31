import { Box, Center, GridItem, Icon } from '@chakra-ui/react';
import { useChain } from '@cosmos-kit/react';
import { MouseEventHandler } from 'react';
import { FiAlertTriangle } from 'react-icons/fi';

import { defaultChainName as chainName } from '@/config';

import {
  Error,
  Connected,
  Connecting,
  ConnectStatusWarn,
  Disconnected,
  NotExist,
  Rejected,
  RejectedWarn,
  WalletConnectComponent,
} from '../components';


export const WalletSection = () => {
  const { connect, openView, status, message, wallet } = useChain(chainName);

  // Events
  const onClickConnect: MouseEventHandler = async (e) => {
    e.preventDefault();
    await connect();
  };

  const onClickOpenView: MouseEventHandler = (e) => {
    e.preventDefault();
    openView();
  };

  // Components
  const connectWalletButton = (
    <WalletConnectComponent
      walletStatus={status}
      disconnect={<Disconnected buttonText="Connect Wallet" onClick={onClickConnect} />}
      connecting={<Connecting />}
      connected={<Connected buttonText={'My Wallet'} onClick={onClickOpenView} />}
      rejected={<Rejected buttonText="Reconnect" onClick={onClickConnect} />}
      error={<Error buttonText="Change Wallet" onClick={onClickOpenView} />}
      notExist={<NotExist buttonText="Connect Wallet" onClick={onClickOpenView} />}
    />
  );

  const connectWalletWarn = (
    <ConnectStatusWarn
      walletStatus={status}
      rejected={<RejectedWarn icon={<Icon as={FiAlertTriangle} mt={1} />} wordOfWarning={`${wallet?.prettyName}: ${message}`} />}
      error={<RejectedWarn icon={<Icon as={FiAlertTriangle} mt={1} />} wordOfWarning={`${wallet?.prettyName}: ${message}`} />}
    />
  );

  return (
    <Center py={16}>
      <Box w="full" maxW={{ base: 52, md: 64 }}>
        {connectWalletButton}
      </Box>
      {connectWalletWarn && <GridItem>{connectWalletWarn}</GridItem>}
    </Center>
  );
};
