import { Box, Link, useToast, Text } from '@chakra-ui/react';
import { SigningStargateClient, StdFee } from '@cosmjs/stargate';
import { ChainName } from '@cosmos-kit/core';

import { Dispatch, SetStateAction } from 'react';
import { ibc } from 'interchain-query';

const showSuccessToast = (toast: ReturnType<typeof useToast>, txHash: string, chainName: ChainName) => {
  const mintscanUrl = `https://www.mintscan.io/${chainName}/txs/${txHash}`;
  toast({
    position: 'bottom-right',
    duration: 5000,
    isClosable: true,
    render: () => (
      <Box color="white" p={3} bg="green.500" borderRadius="md">
        <Text mb={1} fontWeight="bold">
          Transaction Successful
        </Text>
        <Link href={mintscanUrl} isExternal>
          View on Mintscan: {mintscanUrl}
        </Link>
      </Box>
    ),
  });
};

const showErrorToast = (toast: ReturnType<typeof useToast>, errorMsg: string) => {
  toast({
    title: 'Transaction Failed',
    description: `Error: ${errorMsg}`,
    status: 'error',
    duration: 5000,
    isClosable: true,
    position: 'bottom-right',
  });
};

export const ibcWithdrawlTx = async (
  dstAddress: string,
  fromAddress: string,
  getSigningStargateClient: () => Promise<SigningStargateClient>,
  setResp: Dispatch<SetStateAction<string>>,
  toast: ReturnType<typeof useToast>,
  setIsError: Dispatch<SetStateAction<boolean>>,
  setIsSigning: Dispatch<SetStateAction<boolean>>,
  chainName: ChainName,
) => {
  setIsError(false);
  setIsSigning(true);

  try {
    const stargateClient = await getSigningStargateClient();

    if (!stargateClient || !fromAddress) {
      console.error('Stargate client undefined or fromAddress undefined.');
      return;
    }

    const { transfer } = ibc.applications.transfer.v1.MessageComposer.withTypeUrl;

    const msgIbcTransfer = transfer({
      sourcePort: 'transfer',
      sourceChannel: 'channel-0',
      token: {
        denom: 'uqck',
        amount: '7500',
      },
      sender: fromAddress,
      receiver: dstAddress,
      timeoutHeight: {
        revisionNumber: BigInt(0),
        revisionHeight: BigInt(0),
      },
      timeoutTimestamp: BigInt(0),
      memo: '',
    });

    const fee: StdFee = {
      amount: [
        {
          denom: 'uqck',
          amount: '7500',
        },
      ],
      gas: '500000',
    };

    const response = await stargateClient.signAndBroadcast(fromAddress, [msgIbcTransfer], fee);

    // Handle response
    setResp(JSON.stringify(response, null, 2));
    setIsSigning(false);

    if (response.code === 0) {
      showSuccessToast(toast, response.transactionHash, chainName);
    } else {
      setIsError(true);
      showErrorToast(toast, 'Transaction failed');
    }
  } catch (error) {
    console.error('Error in unbonding transaction:', error);
    if (error instanceof Error) {
      setIsSigning(false);
      setIsError(true);
      showErrorToast(toast, error.message);
    }
  }
};
