import { Box, Center, Grid, GridItem, Icon, Stack, useColorModeValue } from '@chakra-ui/react';
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
  ChainCard,
} from '../components';

import { defaultChainName as chainName } from '@/config';

export const WalletSection = () => {
  const { connect, openView, status, username, address, message, wallet, chain: chainInfo } = useChain(chainName);
  const { getChainLogo } = useManager();

  const chain = {
    chainName,
    label: chainInfo.pretty_name,
    value: chainName,
    icon: getChainLogo(chainName),
  };

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

  const userInfo = username && <ConnectedUserInfo username={username} icon={<Astronaut />} />;
  const addressBtn = <CopyAddressBtn walletStatus={status} connected={<ConnectedShowAddress address={address} isLoading={false} />} />;

  return (
    <Center py={16}>
      <Box w="full" maxW={{ base: 52, md: 64 }}>
        {connectWalletButton}
      </Box>
      {connectWalletWarn && <GridItem>{connectWalletWarn}</GridItem>}
    </Center>
  );
};
