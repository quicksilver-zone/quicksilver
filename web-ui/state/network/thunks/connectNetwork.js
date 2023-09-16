import { createAsyncThunk } from "@reduxjs/toolkit";
import { getKeplrFromWindow } from '@keplr-wallet/stores';
import { getSigningQuicksilverClient } from "quicksilverjs";
import { ChainInfos } from "@/state/network/utils";

const fetchNetworkDetails = async (windowWallet, offlineSigner, chainId) => {
    let pubkey = await windowWallet?.getKey(chainId);
    let bech32 = pubkey?.bech32Address;
    if (bech32) {
        let roBalance = await offlineSigner.getAllBalances(bech32);
        return { address: bech32, balance: roBalance, connected: true }
    }
    return { connected: false, address: "", balance: "" }
}

const connectToNetwork = createAsyncThunk("network/connect", async (network, { getState }) => {
    let result = { connected: false, address: "", balance: "" }
    let walletType = (await getState()).typeWallet
    const chainInfo = ChainInfos.find(function (el) { return el.chainId === network })
    try {
        if (walletType === 'keplr') {
            const keplr = await getKeplrFromWindow();
            if (keplr) {
                await keplr.enable(chainInfo.chainId)
                let signer = keplr.getOfflineSignerOnlyAmino(chainInfo.chainId);
                let offlineSigner = await getSigningQuicksilverClient({ rpcEndpoint: chainInfo.rpc, signer: signer });
                localStorage.setItem('ChainId', JSON.stringify(chainInfo.chainId));
                console.log("Enabled for chainid " + chainInfo.chainId)

                await keplr.experimentalSuggestChain(chainInfo)
                let signer1 = keplr.getOfflineSignerOnlyAmino(chainInfo.chainId);
                let offlineSigner1 = await getSigningQuicksilverClient({ rpcEndpoint: chainInfo.rpc, signer: signer1 });
                result = await fetchNetworkDetails(keplr, offlineSigner1, chainInfo.chainId)
                localStorage.setItem('ChainId', JSON.stringify(chainInfo.chainId));
                console.log("Added to Keplr for chainid " + chainInfo.chainId)
            }
        } else if (walletType === 'leap') {
            await window.leap.enable(chainInfo.chainId)
            let signer = window.leap.getOfflineSignerOnlyAmino(chainInfo.chainId);
            let offlineSigner = await getSigningQuicksilverClient({ rpcEndpoint: chainInfo.rpc, signer: signer });
            localStorage.setItem('ChainId', JSON.stringify(chainInfo.chainId));
            console.log("Enabled for chainid LEAP" + chainInfo.chainId)

            window.leap.experimentalSuggestChain(chainInfo)

            let signer1 = window.leap.getOfflineSignerOnlyAmino(chainInfo.chainId);
            let offlineSigner1 = await getSigningQuicksilverClient({ rpcEndpoint: chainInfo.rpc, signer: signer1 });
            result = await fetchNetworkDetails(window.leap, offlineSigner1, chainInfo.chainId)
            localStorage.setItem('ChainId', JSON.stringify(chainInfo.chainId));
            console.log("Added to Leap for chainid " + chainInfo.chainId)
        }
        else if (walletType === 'cosmostation') {
            const supportedChainIds = await window.cosmostation.cosmos.request({
                method: "cos_supportedChainIds",
            });
            console.log('chain ids', supportedChainIds)
            if (!supportedChainIds['unofficial'].includes(chainInfo.chainId)) {
                await window.cosmostation.cosmos.request({
                    method: "cos_addChain",
                    params: {
                        chainId: chainInfo.chainId,
                        chainName: chainInfo.chainName,
                        addressPrefix: chainInfo.bech32Config.bech32PrefixAccAddr,
                        baseDenom: chainInfo.currencies[0].coinMinimalDenom,
                        displayDenom: chainInfo.currencies[0].coinDenom,
                        restURL: chainInfo.rest,
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
            await window.cosmostation.providers.keplr.enable(chainInfo.chainId)
            let signer = window.cosmostation.providers.keplr.getOfflineSignerOnlyAmino(chainInfo.chainId);
            let offlineSigner = await getSigningQuicksilverClient({ rpcEndpoint: chainInfo.rpc, signer: signer });
            localStorage.setItem('ChainId', JSON.stringify(chainInfo.chainId));
            console.log("Enabled for chainid cosmostation" + chainInfo.chainId);

            window.cosmostation.providers.keplr.experimentalSuggestChain(chainInfo)

            let signer1 = window.cosmostation.providers.keplr.getOfflineSignerOnlyAmino(chainInfo.chainId);
            let offlineSigner1 = await getSigningQuicksilverClient({ rpcEndpoint: chainInfo.rpc, signer: signer1 });
            console.log("Added to Keplr for chainid cosmostation" + chainInfo.chainId)
            await window.cosmostation.cosmos.request({
                method: "cos_addChain",
                params: {
                    chainId: chainInfo.chainId,
                    chainName: chainInfo.chainName,
                    addressPrefix: chainInfo.bech32Config.bech32PrefixAccAddr,
                    baseDenom: chainInfo.currencies[0].coinMinimalDenom,
                    displayDenom: chainInfo.currencies[0].coinDenom,
                    restURL: chainInfo.rest,
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
            result = await fetchNetworkDetails(window.cosmostation.providers.keplr, offlineSigner1, chainInfo.chainId)
        }
    } catch (e) {
        console.log(e)
    }
    return {...result, network}
})

export default connectToNetwork;