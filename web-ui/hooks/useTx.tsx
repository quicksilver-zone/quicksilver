import { Box, color, Link, ToastId, useToast, Text } from '@chakra-ui/react';
import { isDeliverTxSuccess, StdFee } from '@cosmjs/stargate';
import { useChain } from '@cosmos-kit/react';
import { cosmos } from 'interchain-query';
import { TxRaw } from 'interchain-query/cosmos/tx/v1beta1/tx';

import { useToaster, ToastType, type CustomToast } from './useToaster';
import { ChainName } from '@cosmos-kit/core';
import { on } from 'events';
import { isExternal } from 'util/types';

interface TxResult {
  success: boolean;
  txHash?: string;
  error?: string;
}

interface Msg {
  typeUrl: string;
  value: any;
}

export interface TxOptions {
  fee?: StdFee | null;
  memo?: string;
  toast?: Partial<CustomToast>;
  onSuccess?: () => void;
}

export enum TxStatus {
  Failed = 'Transaction Failed',
  Successful = 'Transaction Successful',
  Broadcasting = 'Transaction Broadcasting',
}

const txRaw = cosmos.tx.v1beta1.TxRaw;

export const useTx = (chainName: string) => {
  const { address, getSigningStargateClient, estimateFee } = useChain(chainName);

  const toaster = useToaster();

  const tx = async (msgs: Msg[], options: TxOptions): Promise<TxResult> => {
    if (!address) {
      toaster.toast({
        type: ToastType.Error,
        title: 'Wallet not connected',
        message: 'Please connect your wallet',
      });
      return { success: false, error: 'Wallet not connected' };
    }

    let signed: TxRaw;
    let client: Awaited<ReturnType<typeof getSigningStargateClient>>;

    try {
      let fee: StdFee;
      if (options?.fee) {
        fee = options.fee;
        client = await getSigningStargateClient();
      } else {
        const [_fee, _client] = await Promise.all([estimateFee(msgs), getSigningStargateClient()]);
        fee = _fee;
        client = _client;
      }
      signed = await client.sign(address, msgs, fee, '');
    } catch (e: any) {
      console.error(e);
      toaster.toast({
        title: TxStatus.Failed,
        message: e?.message || 'An unexpected error has occured',
        type: ToastType.Error,
      });
      return { success: false, error: 'An unexpected error has occured' };
    }

    let broadcastToastId: ToastId;

    broadcastToastId = toaster.toast({
      title: TxStatus.Broadcasting,
      message: 'Waiting for transaction to be included in the block',
      type: ToastType.Loading,
      duration: 999999,
    });

    if (client && signed) {
      await client
        .broadcastTx(Uint8Array.from(txRaw.encode(signed).finish()))
        .then((res) => {
          if (isDeliverTxSuccess(res)) {
            if (options.onSuccess) options.onSuccess();

            toaster.toast({
              title: options.toast?.title || TxStatus.Successful,
              type: options.toast?.type || ToastType.Success,
              message: options.toast?.message,
              chainName: chainName,
              txHash: res?.transactionHash,
            });
            return { success: true, txHash: res.transactionHash };
          } else {
            toaster.toast({
              title: TxStatus.Failed,
              message: res?.rawLog,
              type: ToastType.Error,
              duration: 10000,
            });
            return { success: false, error: res.rawLog };
          }
        })
        .catch((err) => {
          toaster.toast({
            title: TxStatus.Failed,
            message: err?.message,
            type: ToastType.Error,
            duration: 10000,
          });
          return { success: false, error: err.message };
        })
        .finally(() => toaster.close(broadcastToastId));
      return { success: false, error: 'Client not initialized or transaction not signed' };
    } else {
      toaster.close(broadcastToastId);
      return { success: false, error: 'Client not initialized or transaction not signed' };
    }
  };

  return { tx };
};
