import { EncodeObject } from '@cosmjs/proto-signing';
import { GasPrice, calculateFee } from '@cosmjs/stargate';
import { useChain } from '@cosmos-kit/react';

import { getCoin } from '../config';
import { useState } from 'react';

export const useFeeEstimation = (chainName: string) => {
  const { getSigningStargateClient, chain } = useChain(chainName);
  const [error, setError] = useState<string | null>(null)

  const gasPrice = chain.fees?.fee_tokens[0].average_gas_price || 0.025;

  const coin = getCoin(chainName);

  const estimateFee = async (
    address: string,
    messages: EncodeObject[],
    modifier?: number,
    memo?: string
  ) => {
    try {
      const stargateClient = await getSigningStargateClient();
      if (!stargateClient) {
        throw new Error('getSigningStargateClient error');
      }
  
      const gasEstimation = await stargateClient.simulate(
        address,
        messages,
        memo
      );
      if (!gasEstimation) {
        throw new Error('estimate gas error');
      }
  
      const fee = calculateFee(
        Math.round(gasEstimation * (modifier || 1.5)),
        GasPrice.fromString(gasPrice + coin.base)
      );
      return fee;
    } catch (err: any) {
      setError(err?.message || err);
    }
  };
  return { estimateFee, error };
};