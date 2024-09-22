export enum ENVTYPES {
    PROD = "prod",
    TEST = "test"
}

export type Chain = {
    chain_id: string;
    chain_name: string;
    pretty_name: string;
    rpc: string[];
    rest: string[];
    explorer: string;
    show: boolean;
    enable_deposits: boolean;
    enable_withdrawals: boolean;
    lsm_enabled: boolean;
    major_denom: string;
    minor_denom: string;
    exponent: number;
    logo: string;
    qlogo: string;
    is_118: boolean;
}

const quicksilver_mainnet: Chain = {
    chain_id: "quicksilver-2",
    chain_name: "quicksilver",
    pretty_name: "Quicksilver",
    rpc: ["https://quicksilver-2.rpc.quicksilver.zone"],
    rest: ["https://quicksilver-2.lcd.quicksilver.zone"],
    explorer: "https://explorer.quicksilver.zone/tx/{}",
    show: false,
    enable_deposits: true,
    enable_withdrawals: true,
    lsm_enabled: false,
    major_denom: "qck",
    minor_denom: "uqck",
    exponent: 6,
    logo: '/img/networks/qck.svg',
    qlogo: '/img/networks/qck.svg',
    is_118: true,
}

const quicksilver_testnet: Chain = {
    chain_id: "rhye-3",
    chain_name: "quicksilver-testnet",
    pretty_name: "Quicksilver",
    rpc: ["https://rhye-3.rpc.quicksilver.zone"],
    rest: ["https://rhye-3.lcd.quicksilver.zone"],
    explorer: "https://testnet.quicksilver.explorers.guru/transaction/{}",
    show: false,
    enable_deposits: true,
    enable_withdrawals: true,
    lsm_enabled: false,
    major_denom: "qck",
    minor_denom: "uqck",
    exponent: 6,
    logo: '/img/networks/qck.svg',
    qlogo: '/img/networks/qck.svg',
    is_118: true,
}

export const local_chain = new Map<string, Chain>([
    [ENVTYPES.PROD, quicksilver_mainnet],
    [ENVTYPES.TEST, quicksilver_testnet]
])

const test_chains = new Map<string, Chain>([
    ["quicksilver-testnet", quicksilver_testnet],
    ["cosmos-testnet", {
        chain_id: "theta-testnet-001",
        chain_name: "cosmos-testnet",
        pretty_name: "Cosmos",
        rpc: ["https://theta-testnet-001.rpc.quicksilver.zone"],
        rest: ["https://theta-testnet-001.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/cosmos-testnet/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: true,
        major_denom: "atom",
        minor_denom: "uatom",
        exponent: 6,
        logo: '/img/networks/atom.svg',
        qlogo: '/img/networks/qatom.svg',
        is_118: true,
    }]
])

const prod_chains = new Map<string, Chain>([
    ["quicksilver", quicksilver_mainnet],
    ["cosmoshub", {
        chain_id: "cosmoshub-4",
        chain_name: "cosmoshub",
        pretty_name: "Cosmos",
        rpc: ["https://cosmoshub-4.rpc.quicksilver.zone"],
        rest: ["https://cosmoshub-4.lcd.quicksilver.zone"],
        explorer: "https://mintscan.io/cosmos/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: true,
        major_denom: "atom",
        minor_denom: "uatom",
        exponent: 6,
        logo: '/img/networks/atom.svg',
        qlogo: '/img/networks/qatom.svg',
        is_118: true,
    }],[
    "osmosis", {
        chain_id: "osmosis-1",
        chain_name: "osmosis",
        pretty_name: "Osmosis",
        rpc: ["https://osmosis-1.rpc.quicksilver.zone"],
        rest: ["https://osmosis-1.lcd.quicksilver.zone"],
        explorer: "https://mintscan.io/osmosis/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "osmo",
        minor_denom: "uosmo",
        exponent: 6,
        logo: '/img/networks/osmo.svg',
        qlogo: '/img/networks/qosmo.svg',
        is_118: true,
    }],[
    "stargaze", {
        chain_id: "stargaze-1",
        chain_name: "stargaze",
        pretty_name: "Stargaze",
        rpc: ["https://stargaze-1.rpc.quicksilver.zone"],
        rest: ["https://stargaze-1.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/stargaze/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "stars",
        minor_denom: "ustars",
        exponent: 6,
        logo: '/img/networks/stargaze.svg',
        qlogo: '/img/networks/qstars.svg',
        is_118: true,
    }],
    ["juno", {
        chain_id: "juno-1",
        chain_name: "juno",
        pretty_name: "Juno",
        rpc: ["https://juno-1.rpc.quicksilver.zone"],
        rest: ["https://juno-1.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/juno/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "juno",
        minor_denom: "ujuno",
        exponent: 6,
        logo: '/img/networks/juno.svg',
        qlogo: '/img/networks/qjuno.svg',
        is_118: true,
    }],
    ["regen", {
        chain_id: "regen-1",
        chain_name: "regen",
        pretty_name: "Regen",
        rpc: ["https://regen-1.rpc.quicksilver.zone"],
        rest: ["https://regen-1.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/regen/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "regen",
        minor_denom: "uregen",
        exponent: 6,
        logo: '/img/networks/regen.svg',
        qlogo: '/img/networks/qregen.svg',
        is_118: true,
    }],
    // ["terra", {
    //     chain_id: "terra-2",
    //     rpc: "https://terra-2.rpc.quicksilver.zone",
    //     lcd: "https://terra-2.lcd.quicksilver.zone",
    //     explorer: "https://www.mintscan.io/terra/tx/{}",
    //     show: false,
    //     enable_deposits: true,
    //     enable_withdrawals: true,
    //     lsm_enabled: false,
    // }],
    ["saga", {
        chain_id: "ssc-1",
        chain_name: "saga",
        pretty_name: "Saga",
        rpc: ["https://ssc-1.rpc.quicksilver.zone"],
        rest: ["https://ssc-1.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/saga/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "saga",
        minor_denom: "usaga",
        exponent: 6,
        logo: '/img/networks/saga.svg',
        qlogo: '/img/networks/qsaga.svg',
        is_118: true,
    }],
    ["celestia", {
        chain_id: "celestia",
        chain_name: "celestia",
        pretty_name: "Celestia",
        rpc: ["https://celestia.rpc.quicksilver.zone"],
        rest: ["https://celestia.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/celestia/tx/{}",
        show: false,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "tia",
        minor_denom: "utia",
        exponent: 6,
        logo: '/img/networks/tia.svg',
        qlogo: '/img/networks/qtia.svg',
        is_118: true,
    }],
    ["dydx", {
        chain_id: "dydx-mainnet-1",
        chain_name: "dydx",
        pretty_name: "dYdX",
        rpc: ["https://dydx-mainnet-1.rpc.quicksilver.zone"],
        rest: ["https://dydx-mainnet-1.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/dydx/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "dydx",
        minor_denom: "adydx",
        exponent: 18,
        logo: '/img/networks/dydx.svg',
        qlogo: '/img/networks/qdydx.svg',
        is_118: true,
    }],
    ["sommelier", {
        chain_id: "sommelier-3",
        chain_name: "sommelier",
        pretty_name: "Sommelier",
        rpc: ["https://sommelier-3.rpc.quicksilver.zone"],
        rest: ["https://sommelier-3.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/sommelier/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "somm",
        minor_denom: "usomm",
        exponent: 6,
        logo: '/img/networks/somm.svg',
        qlogo: '/img/networks/qsomm.svg',
        is_118: true,
    }],
    ["umee", {
        chain_id: "umee-1",
        chain_name: "umee",
        pretty_name: "Umee",
        rpc: ["https://umee-1.rpc.quicksilver.zone"],
        rest: ["https://umee-1.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/umee/tx/{}",
        show: false,
        enable_deposits: false,
        enable_withdrawals: false,
        lsm_enabled: false,
        major_denom: "umee",
        minor_denom: "uumee",
        exponent: 6,
        logo: '/img/networks/umee.svg',
        qlogo: '/img/networks/qumee.svg',
        is_118: true,
    }],
    ["agoric", {
        chain_id: "agoric-3",
        chain_name: "agoric",
        pretty_name: "Agoric",
        rpc: ["https://agoric-3.rpc.quicksilver.zone"],
        rest: ["https://agoric-3.lcd.quicksilver.zone"],
        explorer: "https://www.mintscan.io/agoric/tx/{}",
        show: true,
        enable_deposits: true,
        enable_withdrawals: true,
        lsm_enabled: false,
        major_denom: "bld",
        minor_denom: "ubld",
        exponent: 6,
        logo: '/img/networks/bld.svg',
        qlogo: '/img/networks/qbld.svg',
        is_118: false,
    }]

]);


export const chains = new Map<string, Map<string, Chain>>([
    [ENVTYPES.PROD, prod_chains],
    [ENVTYPES.TEST, test_chains]
])

export const getEndpoints = (env: string) => {
    return Array.from(chains.get(env)?.entries() ?? []).reduce((acc, [chainname, chain]: [string, Chain]) => ({
        ...acc,
        [chainname]: {
          rpc: chain.rpc,
          rest: chain.rest,
        },
      })
    )
}


