import { Center, Grid, GridItem, Icon } from '@chakra-ui/react';
import { useChain, useManager } from '@cosmos-kit/react';
import { MouseEventHandler } from 'react';
import { FiAlertTriangle } from 'react-icons/fi';

import {
  Astronaut,
  Error,
  Connected,
  ConnectedShowAddress,
  ConnectedUserInfo,
  Connecting,
  ConnectStatusWarn,
  CopyAddressBtn,
  Disconnected,
  NotExist,
  Rejected,
  RejectedWarn,
  WalletConnectComponent,
} from '@/components';
import { useDrawerControl } from '@/state/useDrawerController';

export const WalletButton: React.FC<{ chainName: string }> = ({ chainName }) => {
  const { connect, openView, status, message, wallet } = useChain(chainName || 'cosmoshub');

  const { closeDrawer } = useDrawerControl();

  // Events
  const onClickConnect: MouseEventHandler = async (e) => {
    closeDrawer();
    await connect();
  };

  const onClickOpenView: MouseEventHandler = (e) => {
    closeDrawer();
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
      notExist={<NotExist buttonText="Install Wallet" onClick={onClickOpenView} />}
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
    <Center>
      <Grid w="full" maxW="sm" templateColumns="1fr" alignItems="center" justifyContent="center">
        {connectWalletButton}
        {connectWalletWarn && <GridItem>{connectWalletWarn}</GridItem>}
      </Grid>
    </Center>
  );
};
