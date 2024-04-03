import { ToastId } from '@chakra-ui/react';
import {SkipRouter, TxStatusResponse} from '@skip-router/core';
import { useCallback } from 'react';

import { useToaster, ToastType } from './useToaster';

export function useSkipExecute(skipClient: SkipRouter) {
    if (!skipClient) {
        throw new Error('SkipRouter is not initialized');
    }

    const toaster = useToaster();

    const executeRoute = useCallback(async (route: any, userAddresses: any, refetch: () => void) => {
        // Initialize with null and allow for the type to be null or ToastId
        let broadcastToastId: ToastId | null = null;

        try {
   
            return await skipClient.executeRoute({
                route,
                userAddresses,
                onTransactionCompleted: async (chainID: string, txHash: string, status: TxStatusResponse) => {
                   
                    if (broadcastToastId) {
                        toaster.close(broadcastToastId);
                    }

                   
                    toaster.toast({
                        type: ToastType.Success,
                        title: 'Transaction Successful',
                        message: `Transaction ${txHash} completed on chain ${chainID}`,
                    });

                    refetch();
                },
                onTransactionBroadcast: async (txInfo) => {
     
                    broadcastToastId = toaster.toast({
                        type: ToastType.Loading,
                        title: 'Transaction Broadcasting',
                        message: 'Waiting for transaction to be included in a block',
                        duration: 9999, 
                    });
                },
                onTransactionTracked: async (txInfo) => {
                    if (broadcastToastId) {
                        toaster.close(broadcastToastId);
                    }
              
                },
            });
        } catch (error) {
        
            if (broadcastToastId) {
                toaster.close(broadcastToastId);
            }

            // Show error toast
            console.error('Error executing route:', error);
            toaster.toast({
                type: ToastType.Error,
                title: 'Transaction Failed',
                message: (error as Error).message || 'An unexpected error occurred',
            });
        }
    }, [skipClient, toaster]);

    return executeRoute;
}

export function useSkipMessages(skipClient: SkipRouter) {
    if (!skipClient) {
        throw new Error('SkipRouter is not initialized');
    }
  const skipMessages = useCallback(async (route: any) => {
    return await skipClient.messages({
        sourceAssetDenom: route.sourceAssetDenom,
        sourceAssetChainID: route.sourceAssetChainID,
        destAssetDenom: route.destAssetDenom,
        destAssetChainID: route.destAssetChainID,
        amountIn: route.amountIn,
        amountOut: route.amountOut,
        addressList: route.addressList,
        operations: route.operations,
    });
  }, []);

  return skipMessages;
}