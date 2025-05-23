import { useToast, Text, Box, Link } from '@chakra-ui/react';

import { convertChainName } from '@/utils';
import { chains, env } from '@/config';

export enum ToastType {
  Info = 'info',
  Error = 'error',
  Success = 'success',
  Loading = 'loading',
}

export type CustomToast = {
  title: string;
  type: ToastType;
  message?: string | JSX.Element;
  closable?: boolean;
  duration?: number;
  txHash?: string;
  chainName?: string;
};
export const useToaster = () => {
  const toast = useToast({
    position: 'top-right',
    containerStyle: {
      maxWidth: '300px',
    },
  });

  const customToast = ({ type, title, message, closable = true, duration = 5000, txHash, chainName }: CustomToast) => {
    let description;

    if (type === ToastType.Success && txHash && chainName) {
      let explorerUrl = chains.get(env)?.get(chainName)?.explorer ?? "";
      if (explorerUrl.includes("{}")) {
        explorerUrl = explorerUrl.replace("{}", txHash);
      }
      description = (
        <Box pr="20px">
          <Text fontSize="sm" color="white">
            {message}
          </Text>
          <Link href={explorerUrl} isExternal color="white">
            View in Explorer
          </Link>
        </Box>
      );
    } else {
      description = (
        <Box pr="20px">
          <Text fontSize="sm" color="white">
            {message}
          </Text>
        </Box>
      );
    }

    return toast({
      position: 'bottom-right',
      title,
      duration,
      status: type,
      isClosable: closable,
      description,
    });
  };

  return { ...toast, toast: customToast };
};
  