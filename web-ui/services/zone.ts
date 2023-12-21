import { setupStakingExtension, QueryClient } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import axios from 'axios'

import { ProdZoneInfos } from '@/state/chains/prod'


export const statusList = [
    "BOND_STATUS_BONDED",
    "BOND_STATUS_UNBONDING",
    "BOND_STATUS_UNBONDED"
]


export const getAPY = async (chainId: string) => {
    try {
        const res = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_DATA_API}/apr`)
        const { chains } = res.data
        if (!chains) {
            return 0
        }
        const chainInfo = chains.filter((chain: { chain_id: string; }) => {
            return chain.chain_id === chainId
        })
        return chainInfo.length > 0 ? chainInfo[0].apr : 0
    } catch (e) {
        throw e
    }
}

