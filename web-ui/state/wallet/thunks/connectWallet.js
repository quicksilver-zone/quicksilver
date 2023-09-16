import { createAsyncThunk } from "@reduxjs/toolkit";
import { ProdQuickSilverChainInfo, ProdChainInfos } from '@/state/chains/prod'
import { TestQuickSilverChainInfo, TestChainInfos } from '@/state/chains/test'
import { DevQuickSilverChainInfo, DevChainInfos } from '@/state/chains/dev'
import { getKeplrFromWindow } from '@keplr-wallet/stores';
import { getSigningQuicksilverClient } from "quicksilverjs";
const QuickSilverChains = {
    "preprod": ProdQuickSilverChainInfo,
    "prod": ProdQuickSilverChainInfo,
    "test": TestQuickSilverChainInfo,
    "dev": DevQuickSilverChainInfo,
}

const Chains = {
    "preprod": ProdChainInfos,
    "prod": ProdChainInfos,
    "test": TestChainInfos,
    "dev": DevChainInfos,
}

const ChainInfos = Chains[process.env.NEXT_PUBLIC_CHAIN_ENV]
const QuickSilverChainInfo = QuickSilverChains[process.env.NEXT_PUBLIC_CHAIN_ENV]

const fetchWalletDetails = async (windowWallet, offlineSigner, walletType) => {
    let pubkey = await windowWallet?.getKey(QuickSilverChainInfo.chainId);
    let bech32 = pubkey?.bech32Address;
    if (bech32) {
        let roBalance = await offlineSigner.getAllBalances(bech32);
        localStorage.setItem('WalletType', walletType);
        return {address: bech32, balance: roBalance, connected: true, typeWallet: walletType}
    }
    return {connected: false, typeWallet: "", address: "", balance: ""}
}

const connectToWallet = createAsyncThunk("wallet/connect", async (walletType) => {
    let result = {connected: false, typeWallet: "", address: "", balance: ""}
    try {
        if (walletType === 'keplr') {
            const keplr = await getKeplrFromWindow();
            if (keplr) {
                await keplr.enable(QuickSilverChainInfo.chainId)
                let signer = keplr.getOfflineSignerOnlyAmino(QuickSilverChainInfo.chainId);
                let offlineSigner = await getSigningQuicksilverClient({ rpcEndpoint: QuickSilverChainInfo.rpc, signer: signer });
                localStorage.setItem('ChainId', JSON.stringify(QuickSilverChainInfo.chainId));
                console.log("Enabled for chainid " + QuickSilverChainInfo.chainId)

                await keplr.experimentalSuggestChain(QuickSilverChainInfo)
                let signer1 = keplr.getOfflineSignerOnlyAmino(QuickSilverChainInfo.chainId);
                let offlineSigner1 = await getSigningQuicksilverClient({ rpcEndpoint: QuickSilverChainInfo.rpc, signer: signer1 });
                result = await fetchWalletDetails(keplr, offlineSigner1, walletType)
                localStorage.setItem('ChainId', JSON.stringify(QuickSilverChainInfo.chainId));
                console.log("Added to Keplr for chainid " + QuickSilverChainInfo.chainId)
            }
        } else if (walletType === 'leap') {
            await window.leap.enable(QuickSilverChainInfo.chainId)
            let signer = window.leap.getOfflineSignerOnlyAmino(QuickSilverChainInfo.chainId);
            let offlineSigner = await getSigningQuicksilverClient({ rpcEndpoint: QuickSilverChainInfo.rpc, signer: signer });
            localStorage.setItem('ChainId', JSON.stringify(QuickSilverChainInfo.chainId));
            console.log("Enabled for chainid LEAP" + QuickSilverChainInfo.chainId)

            window.leap.experimentalSuggestChain(QuickSilverChainInfo)

            let signer1 = window.leap.getOfflineSignerOnlyAmino(QuickSilverChainInfo.chainId);
            let offlineSigner1 = await getSigningQuicksilverClient({ rpcEndpoint: QuickSilverChainInfo.rpc, signer: signer1 });
            result = await fetchWalletDetails(window.leap, offlineSigner1, walletType)
            localStorage.setItem('ChainId', JSON.stringify(QuickSilverChainInfo.chainId));
            console.log("Added to Leap for chainid " + QuickSilverChainInfo.chainId)
        }
        else if (walletType === 'cosmostation') {
            const supportedChainIds = await window.cosmostation.cosmos.request({
                method: "cos_supportedChainIds",
            });
            console.log('chain ids', supportedChainIds)
            if (!supportedChainIds['unofficial'].includes(QuickSilverChainInfo.chainId)) {
                await window.cosmostation.cosmos.request({
                    method: "cos_addChain",
                    params: {
                        chainId: QuickSilverChainInfo.chainId,
                        chainName: QuickSilverChainInfo.chainName,
                        addressPrefix: TestQuickSilverChainInfo.bech32Config.bech32PrefixAccAddr,
                        baseDenom: QuickSilverChainInfo.currencies[0].coinMinimalDenom,
                        displayDenom: QuickSilverChainInfo.currencies[0].coinDenom,
                        restURL: QuickSilverChainInfo.rest,
                        coinType: "118", // optional (default: '118')
                        decimals: 6, // optional (default: 6)
                        gasRate: {
                            // optional (default: { average: '0.025', low: '0.0025', tiny: '0.00025' })
                            average: "0.2",
                            low: "0.02",
                            tiny: "0.002",
                        }
                    },
                });
            }
            await window.cosmostation.providers.keplr.enable(QuickSilverChainInfo.chainId)
            let signer = window.cosmostation.providers.keplr.getOfflineSignerOnlyAmino(QuickSilverChainInfo.chainId);
            let offlineSigner = await getSigningQuicksilverClient({ rpcEndpoint: QuickSilverChainInfo.rpc, signer: signer });
            localStorage.setItem('ChainId', JSON.stringify(QuickSilverChainInfo.chainId));
            console.log("Enabled for chainid cosmostation" + QuickSilverChainInfo.chainId);

            window.cosmostation.providers.keplr.experimentalSuggestChain(QuickSilverChainInfo)

            let signer1 = window.cosmostation.providers.keplr.getOfflineSignerOnlyAmino(QuickSilverChainInfo.chainId);
            let offlineSigner1 = await getSigningQuicksilverClient({ rpcEndpoint: QuickSilverChainInfo.rpc, signer: signer1 });
            console.log("Added to Keplr for chainid cosmostation" + QuickSilverChainInfo.chainId)
            await window.cosmostation.cosmos.request({
                method: "cos_addChain",
                params: {
                    chainId: QuickSilverChainInfo.chainId,
                    chainName: QuickSilverChainInfo.chainName,
                    addressPrefix: TestQuickSilverChainInfo.bech32Config.bech32PrefixAccAddr,
                    baseDenom: QuickSilverChainInfo.currencies[0].coinMinimalDenom,
                    displayDenom: QuickSilverChainInfo.currencies[0].coinDenom,
                    restURL: QuickSilverChainInfo.rest,
                    coinType: "118", // optional (default: '118')
                    decimals: 6, // optional (default: 6)
                    gasRate: {
                        // optional (default: { average: '0.025', low: '0.0025', tiny: '0.00025' })
                        average: "0.2",
                        low: "0.02",
                        tiny: "0.002",
                    }
                },
            });
            result = await fetchWalletDetails(window.cosmostation.providers.keplr, offlineSigner1, walletType)
        }
    } catch (e) {
        console.log (e)
    }
    return result
})

export default connectToWallet;
