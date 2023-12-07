import { Box, Text, Link, useToast } from '@chakra-ui/react';
import { StdFee } from '@cosmjs/amino';
import { coins, SigningStargateClient } from '@cosmjs/stargate';
import { ChainName, Dispatch } from '@cosmos-kit/core';
import { bech32 } from 'bech32';
import { assets } from 'chain-registry';
import { cosmos } from 'interchain-query';
import { Zone } from 'quicksilverjs/types/codegen/quicksilver/interchainstaking/v1/interchainstaking';
import { SetStateAction } from 'react';

import { useZoneQuery } from '@/hooks/useQueries';

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

interface ValidatorsSelect {
  address: string;
  intent: number;
}

export const liquidStakeTx = (
  getSigningStargateClient: () => Promise<SigningStargateClient>,
  setResp: (resp: string) => any,
  chainName: string,
  chainId: string,
  address: string | undefined,
  toast: ReturnType<typeof useToast>,
  setIsError: Dispatch<SetStateAction<boolean>>,
  setIsSigning: Dispatch<SetStateAction<boolean>>,
  validatorsSelect: ValidatorsSelect[],
  amount: number,
  zone: Zone,
) => {
  setIsError(false);
  setIsSigning(true);
  console.log(validatorsSelect);
  const valToByte = (val: number) => {
    if (val > 1) {
      val = 1;
    }
    if (val < 0) {
      val = 0;
    }
    return Math.abs(val * 200);
  };

  const addValidator = (valAddr: string, weight: number) => {
    let { words } = bech32.decode(valAddr);
    let wordsUint8Array = new Uint8Array(bech32.fromWords(words));
    let weightByte = valToByte(weight);
    return Buffer.concat([Buffer.from([0x02, 0x15]), Buffer.from([weightByte]), wordsUint8Array]);
  };

  let memoBuffer = Buffer.alloc(0);
  validatorsSelect.forEach((val) => {
    memoBuffer = Buffer.concat([memoBuffer, addValidator(val.address, val.intent / 100)]);
  });
  const memo = memoBuffer.toString('base64');

  console.log(amount);

  return async (event: React.MouseEvent) => {
    event.preventDefault();
    const stargateClient = await getSigningStargateClient();

    if (!stargateClient || !address) {
      console.error('Stargate client undefined or address undefined.');
      return;
    }

    const { send } = cosmos.bank.v1beta1.MessageComposer.withTypeUrl;
    const msgSend = send({
      fromAddress: address,
      toAddress: zone?.depositAddress?.address ?? '',
      amount: coins(amount.toFixed(0), zone?.baseDenom ?? ''),
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
      const response = await stargateClient.signAndBroadcast(address, [msgSend], fee, memo);
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
