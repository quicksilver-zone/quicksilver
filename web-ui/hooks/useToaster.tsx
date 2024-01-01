import { useToast, Text, Box, Link } from '@chakra-ui/react';

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

    if (type === ToastType.Success && txHash) {
      const mintscanUrl = `https://www.mintscan.io/${chainName}/txs/${txHash}`;
      description = (
        <Box pr="20px">
          <Text fontSize="sm" color="white">
            {message}
          </Text>
          <Link href={mintscanUrl} isExternal color="complimentary.900">
            View on Mintscan
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
