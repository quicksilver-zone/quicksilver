import { Box, Text, Link, useToast } from '@chakra-ui/react';
import { getSigningQuicksilverClient, quicksilver } from '@chalabi/quicksilverjs';
import { ValidatorIntent } from '@chalabi/quicksilverjs/dist/codegen/quicksilver/interchainstaking/v1/interchainstaking';
import { MsgSignalIntent } from '@chalabi/quicksilverjs/dist/codegen/quicksilver/interchainstaking/v1/messages';
import { StdFee } from '@cosmjs/amino';
import { HttpEndpoint, SigningStargateClient } from '@cosmjs/stargate';
import { ChainName, Dispatch } from '@cosmos-kit/core';
import { useChain } from '@cosmos-kit/react';
import { assets, chains } from 'chain-registry';
import { SetStateAction } from 'react';

import { useQueryHooks } from '@/hooks';

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

export const intentTx = (
  getSigningStargateClient: (apiUrl: string) => Promise<SigningStargateClient>,
  setResp: (resp: string) => any,
  chainName: string,
  chainId: string,
  address: string | undefined,
  intents: ValidatorIntent[],
  toast: ReturnType<typeof useToast>,
  setIsError: Dispatch<SetStateAction<boolean>>,
  setIsSigning: Dispatch<SetStateAction<boolean>>,
) => {
  setIsError(false);
  setIsSigning(true);

  return async (event: React.MouseEvent) => {
    event.preventDefault();
    const apiUrl = 'https://rpc.test.quicksilver.zone';
    const stargateClient = await getSigningStargateClient(apiUrl);

    if (!stargateClient || !address) {
      console.error('Stargate client undefined or address undefined.');
      return;
    }

    const { signalIntent } = quicksilver.interchainstaking.v1.MessageComposer.withTypeUrl;
    const msgSignalIntent = signalIntent({
      chainId: chainId,
      intents: intents,
      fromAddress: address,
    });

    const mainTokens = assets.find(({ chain_name }) => chain_name === chainName);
    const mainDenom = mainTokens?.assets[0].base ?? 'uqck';

    const fee: StdFee = {
      amount: [
        {
          denom: mainDenom,
          amount: '5000',
        },
      ],
      gas: '500000',
    };

    try {
      stargateClient.registry.register('/quicksilver.interchainstaking.v1.MsgSignalIntent', MsgSignalIntent);
      const response = await stargateClient.signAndBroadcast(address, [msgSignalIntent], fee);
      setResp(JSON.stringify(response, null, 2));
      setIsSigning(false);
      showSuccessToast(toast, response.transactionHash, chainName);
    } catch (error) {
      console.error('Error signing and sending transaction:', error);
      if (error instanceof Error) {
        setIsSigning(false);
        setIsError(true);
        showErrorToast(toast, error.message);
      }
    }
  };
};
