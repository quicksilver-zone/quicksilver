import { Center, Grid, GridItem, Icon } from '@chakra-ui/react';
import { useChain, useManager, useWallet, useWalletClient } from '@cosmos-kit/react';
import { MouseEventHandler, useEffect } from 'react';
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

export const WalletButton: React.FC<{ chainName: string }> = ({ chainName }) => {
  const { connect, username, address, message, wallet, chain: chainInfo } = useChain(chainName || 'cosmoshub');
  const { getChainLogo } = useManager();

  const chain = {
    chainName,
    label: chainInfo.pretty_name,
    value: chainName,
    icon: getChainLogo(chainName || 'cosmoshub'),
  };

  // Events
  const onClickConnect: MouseEventHandler = async (e) => {
    e.preventDefault();
    await connect();
  };

  const { openView } = useChain('quicksilver');
  const { status: WalletStatus, client } = useWalletClient();
  const { status } = useWallet();

  useEffect(() => {
    if (WalletStatus === 'Done') {
      client?.enable?.(['cosmoshub-4', 'osmosis-1', 'regen-1', 'sommelier-3', 'stargaze-1']);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

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

  const userInfo = username && <ConnectedUserInfo username={username} icon={<Astronaut />} />;
  const addressBtn = <CopyAddressBtn walletStatus={status} connected={<ConnectedShowAddress address={address} isLoading={false} />} />;

  return (
    <Center>
      <Grid w="full" maxW="sm" templateColumns="1fr" alignItems="center" justifyContent="center">
        {connectWalletButton}
        {connectWalletWarn && <GridItem>{connectWalletWarn}</GridItem>}
      </Grid>
    </Center>
  );
};
