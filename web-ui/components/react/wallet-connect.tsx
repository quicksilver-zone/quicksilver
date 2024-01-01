import { Button, Icon, Stack, Text, useColorModeValue } from '@chakra-ui/react';
import { WalletStatus } from '@cosmos-kit/core';
import { useChain, useWalletClient } from '@cosmos-kit/react';
import React, { MouseEventHandler, ReactNode, useEffect } from 'react';
import { FiAlertTriangle } from 'react-icons/fi';
import { IoWallet } from 'react-icons/io5';

import { ConnectWalletType } from '../types';

export const ConnectWalletButton = ({ buttonText, isLoading, isDisabled, icon, onClickConnectBtn }: ConnectWalletType) => {
  const buttonTextColor = useColorModeValue('primary.700', 'primary.50');
  const invertButtonTextColor = useColorModeValue('primary.50', 'primary.700');
  const { openView } = useChain('quicksilver');
  const { status, client } = useWalletClient();

  useEffect(() => {
    if (status === 'Done') {
      client?.enable?.(['cosmoshub-4', 'osmosis-1', 'regen-1', 'sommelier-3', 'stargaze-1']);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

  return (
    <Button
      variant="solid"
      minW="fit-content"
      isLoading={isLoading}
      isDisabled={isDisabled}
      bgColor="complimentary.900"
      color={invertButtonTextColor}
      sx={{
        position: 'relative',
        boxShadow: '0 6px 8px rgba(0, 0, 0, 0.3)',
        borderRadius: 100,
        _hover: {
          boxShadow: '0 0 10px #FF9933, 0 0 15px #FFB266',
          transform: 'scale(1.1)',
          '&::before': {
            content: '""',
            position: 'absolute',
            top: '-2px',
            right: '-2px',
            bottom: '-2px',
            left: '-2px',
            zIndex: -1,
            borderRadius: 'inherit',
            background: 'linear-gradient(60deg, #FF8000, #FF9933, #FFB266, #FFD9B3, #FFE6CC)',
            backgroundSize: '350% 350%',
            animation: 'animatedgradient 12s linear infinite',
          },
        },
        _active: {
          transform: 'translateY(4px)',
        },
        '@keyframes animatedgradient': {
          '0%': { backgroundPosition: '0% 50%' },
          '50%': { backgroundPosition: '100% 50%' },
          '100%': { backgroundPosition: '0% 50%' },
        },
      }}
      onClick={openView}
    >
      <Icon as={icon ? icon : IoWallet} mr={2} />
      {buttonText ? buttonText : 'Connect'}
    </Button>
  );
};

export const Disconnected = ({ buttonText, onClick }: { buttonText: string; onClick: MouseEventHandler<HTMLButtonElement> }) => {
  return <ConnectWalletButton buttonText={buttonText} onClickConnectBtn={onClick} />;
};

export const Connected = ({ buttonText, onClick }: { buttonText: string; onClick: MouseEventHandler<HTMLButtonElement> }) => {
  return <ConnectWalletButton buttonText={buttonText} onClickConnectBtn={onClick} />;
};

export const Connecting = () => {
  return <ConnectWalletButton isLoading={true} />;
};

export const Rejected = ({
  buttonText,
  wordOfWarning,
  onClick,
}: {
  buttonText: string;
  wordOfWarning?: string;
  onClick: MouseEventHandler<HTMLButtonElement>;
}) => {
  const bg = useColorModeValue('orange.200', 'orange.300');

  return (
    <Stack>
      <ConnectWalletButton buttonText={buttonText} isDisabled={false} onClickConnectBtn={onClick} />
      {wordOfWarning && (
        <Stack isInline={true} borderRadius="md" bg={bg} color="blackAlpha.900" p={4} spacing={1}>
          <Icon as={FiAlertTriangle} mt={1} />
          <Text>
            <Text fontWeight="semibold" as="span">
              Warning:&ensp;
            </Text>
            {wordOfWarning}
          </Text>
        </Stack>
      )}
    </Stack>
  );
};

export const Error = ({
  buttonText,
  wordOfWarning,
  onClick,
}: {
  buttonText: string;
  wordOfWarning?: string;
  onClick: MouseEventHandler<HTMLButtonElement>;
}) => {
  const bg = useColorModeValue('orange.200', 'orange.300');

  return (
    <Stack>
      <ConnectWalletButton buttonText={buttonText} isDisabled={false} onClickConnectBtn={onClick} />
      {wordOfWarning && (
        <Stack isInline={true} borderRadius="md" bg={bg} color="blackAlpha.900" p={4} spacing={1}>
          <Icon as={FiAlertTriangle} mt={1} />
          <Text>
            <Text fontWeight="semibold" as="span">
              Warning:&ensp;
            </Text>
            {wordOfWarning}
          </Text>
        </Stack>
      )}
    </Stack>
  );
};

export const NotExist = ({ buttonText, onClick }: { buttonText: string; onClick: MouseEventHandler<HTMLButtonElement> }) => {
  return <ConnectWalletButton buttonText={buttonText} isDisabled={false} onClickConnectBtn={onClick} />;
};

export const WalletConnectComponent = ({
  walletStatus,
  disconnect,
  connecting,
  connected,
  rejected,
  error,
  notExist,
}: {
  walletStatus: WalletStatus;
  disconnect: ReactNode;
  connecting: ReactNode;
  connected: ReactNode;
  rejected: ReactNode;
  error: ReactNode;
  notExist: ReactNode;
}) => {
  switch (walletStatus) {
    case WalletStatus.Disconnected:
      return <>{disconnect}</>;
    case WalletStatus.Connecting:
      return <>{connecting}</>;
    case WalletStatus.Connected:
      return <>{connected}</>;
    case WalletStatus.Rejected:
      return <>{rejected}</>;
    case WalletStatus.Error:
      return <>{error}</>;
    case WalletStatus.NotExist:
      return <>{notExist}</>;
    default:
      return <>{disconnect}</>;
  }
};
