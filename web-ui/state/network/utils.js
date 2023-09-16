import { ProdChainInfos } from '@/state/chains/prod'
import { TestChainInfos } from '@/state/chains/test'
import { DevChainInfos } from '@/state/chains/dev'

const Chains = {
    "preprod": ProdChainInfos,
    "prod": ProdChainInfos,
    "test": TestChainInfos,
    "dev": DevChainInfos,
}

export const ChainInfos = Chains[process.env.NEXT_PUBLIC_CHAIN_ENV]

export const NetworkConfig = {
    uqatom: {
        logo: "/assets/Cosmos.png",
        name: "Cosmos Hub",
    },
    uqosmo: {
        logo: "/assets/Osmosis.png",
        name: "Osmosis",
    },
    uqstars: {
        logo: "/assets/Stargaze.png",
        name: "Stargaze",
    },
    uqjunox: {
        logo: "/assets/Juno.png",
        name: "Juno",
    },
    uqregen: {
        logo: "/assets/Regen.png",
        name: "Regen",
    },
    uqsomm: {
        logo: "/assets/sommelier.png",
        name: "Sommelier",
    }
}