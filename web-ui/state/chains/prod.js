export const ProdQuickSilverChainInfo  = {
    chainId: "quicksilver-2",
    chainName: "Quicksilver Protocol",
    rpc: "https://rpc.quicksilver.zone",
    rest: "https://lcd.quicksilver.zone",
    bip44: {
        coinType: 118,
    },
    bech32Config: {
        bech32PrefixAccAddr: "quick",
        bech32PrefixAccPub: "quickpub",
        bech32PrefixValAddr: "quickvaloper",
        bech32PrefixValPub: "quickvaloperpub",
        bech32PrefixConsAddr: "quickvalcons",
        bech32PrefixConsPub: "quickvalconspub",
    },
    currencies: [
        {
            coinDenom: "QCK",
            coinMinimalDenom: "uqck",
            coinDecimals: 6,
            coinGeckoId: "quicksilver",
        },
        {
            coinDenom: "qATOM",
            coinMinimalDenom: "uqatom",
            coinDecimals: 6,
            coinGeckoId: "cosmos",
        },
        {
            coinDenom: "qOSMO",
            coinMinimalDenom: "uqosmo",
            coinDecimals: 6,
            coinGeckoId: "osmosis",
        },
        {
            coinDenom: "qSTARS",
            coinMinimalDenom: "uqstars",
            coinDecimals: 6,
            coinGeckoId: "stargaze",
        },
        {
            coinDenom: "qJUNO",
            coinMinimalDenom: "uqjuno",
            coinDecimals: 6,
            coinGeckoId: "juno",
        },
        {
            coinDenom: "qREGEN",
            coinMinimalDenom: "uqregen",
            coinDecimals: 6,
            coinGeckoId: "regen",
        },
        {
            coinDenom: "qSOMM",
            coinMinimalDenom: "uqsomm",
            coinDecimals: 6,
            coinGeckoId: "sommelier",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "QCK",
            coinMinimalDenom: "uqck",
            coinDecimals: 6,
            coinGeckoId: "quicksilver",
        },
    ],
    stakeCurrency: {
        coinDenom: "QCK",
        coinMinimalDenom: "uqck",
        coinDecimals: 6,
        coinGeckoId: "quicksilver",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.00,
        average: 0.015,
        high: 0.03,
},
}

export const ProdChainInfos = [
    ProdQuickSilverChainInfo,
{
    chainId: "cosmoshub-4",
    chainName: "CosmosHub",
    rpc: "https://rpc.cosmoshub-4.quicksilver.zone",
    rest: "https://lcd.cosmoshub-4.quicksilver.zone",
    
    bip44: {
        coinType: 118,
    },
    bech32Config: {
        bech32PrefixAccAddr: "cosmos",
        bech32PrefixAccPub: "cosmospub",
        bech32PrefixValAddr: "cosmosvaloper",
        bech32PrefixValPub: "cosmosvaloperpub",
        bech32PrefixConsAddr: "cosmosvalcons",
        bech32PrefixConsPub: "cosmosvalconspub",
    },
    currencies: [
        {
            coinDenom: "ATOM",
            coinMinimalDenom: "uatom",
            coinDecimals: 6,
            coinGeckoId: "cosmos",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "ATOM",
            coinMinimalDenom: "uatom",
            coinDecimals: 6,
            coinGeckoId: "cosmos",
        },
    ],
    stakeCurrency: {
        coinDenom: "ATOM",
        coinMinimalDenom: "uatom",
        coinDecimals: 6,
        coinGeckoId: "cosmos",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.00,
        average: 0.015,
        high: 0.03,
    },
  },
  {
    chainId: "regen-1",
    chainName: "Regen Mainnet",
    rpc: "https://rpc.regen-1.quicksilver.zone",
    rest: "https://lcd.regen-1.quicksilver.zone",
    bip44: {
        coinType: 118,
    },
    bech32Config: {
        bech32PrefixAccAddr: "regen",
        bech32PrefixAccPub: "regenpub",
        bech32PrefixValAddr: "regenvaloper",
        bech32PrefixValPub: "regenvaloperpub",
        bech32PrefixConsAddr: "regenvalcons",
        bech32PrefixConsPub: "regenvalconspub",
    },
    currencies: [
        {
            coinDenom: "REGEN",
            coinMinimalDenom: "uregen",
            coinDecimals: 6,
            coinGeckoId: "regen",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "REGEN",
            coinMinimalDenom: "uregen",
            coinDecimals: 6,
            coinGeckoId: "regen",
        },
    ],
    stakeCurrency: {
        coinDenom: "REGEN",
        coinMinimalDenom: "uregen",
        coinDecimals: 6,
        coinGeckoId: "regen",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.00,
        average: 0.015,
        high: 0.03,
    },
  },
  {
    chainId: "osmosis-1",
    chainName: "Osmosis",
    rpc: "https://rpc.osmosis-1.quicksilver.zone",
    rest: "https://lcd.osmosis-1.quicksilver.zone",
    bip44: {
        coinType: 118,
    },
    bech32Config: {
        bech32PrefixAccAddr: "osmo",
        bech32PrefixAccPub: "osmopub",
        bech32PrefixValAddr: "osmovaloper",
        bech32PrefixValPub: "osmovaloperpub",
        bech32PrefixConsAddr: "osmovalcons",
        bech32PrefixConsPub: "osmovalconspub",
    },
    currencies: [
        {
            coinDenom: "OSMO",
            coinMinimalDenom: "uosmo",
            coinDecimals: 6,
            coinGeckoId: "osmosis",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "OSMO",
            coinMinimalDenom: "uosmo",
            coinDecimals: 6,
            coinGeckoId: "osmosis",
        },
    ],
    stakeCurrency: {
        coinDenom: "OSMO",
        coinMinimalDenom: "uosmo",
        coinDecimals: 6,
        coinGeckoId: "osmosis",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.00,
        average: 0.015,
        high: 0.03,
    },
  },
  {
    chainId: "stargaze-1",
    chainName: "Stargaze",
    rpc: "https://rpc.stargaze-1.quicksilver.zone",
    rest: "https://lcd.stargaze-1.quicksilver.zone",
    
    bip44: {
        coinType: 118,
    },
    bech32Config: {
        bech32PrefixAccAddr: "stars",
        bech32PrefixAccPub: "starspub",
        bech32PrefixValAddr: "starsvaloper",
        bech32PrefixValPub: "starsvaloperpub",
        bech32PrefixConsAddr: "starsvalcons",
        bech32PrefixConsPub: "starsvalconspub",
    },
    currencies: [
        {
            coinDenom: "STARS",
            coinMinimalDenom: "ustars",
            coinDecimals: 6,
            coinGeckoId: "stargaze",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "STARS",
            coinMinimalDenom: "ustars",
            coinDecimals: 6,
            coinGeckoId: "stargaze",
        },
    ],
    stakeCurrency: {  
        coinDenom: "STARS",
        coinMinimalDenom: "ustars",
        coinDecimals: 6,
        coinGeckoId: "stargaze",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.01,
        average: 0.015,
        high: 0.03,
    },
  },
  {
    chainId: "juno-1",
    chainName: "Juno",
    rpc: "https://rpc.juno-1.quicksilver.zone",
    rest: "https://lcd.juno-1.quicksilver.zone",
    
    bip44: {
        coinType: 118,
    },
    bech32Config: {
        bech32PrefixAccAddr: "juno",
        bech32PrefixAccPub: "junopub",
        bech32PrefixValAddr: "junovaloper",
        bech32PrefixValPub: "junovaloperpub",
        bech32PrefixConsAddr: "junovalcons",
        bech32PrefixConsPub: "junovalconspub",
    },
    currencies: [
        {
            coinDenom: "JUNO",
            coinMinimalDenom: "ujuno",
            coinDecimals: 6,
            coinGeckoId: "juno-network",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "JUNO",
            coinMinimalDenom: "ujuno",
            coinDecimals: 6,
            coinGeckoId: "juno-network",
        },
    ],
    stakeCurrency: {
            coinDenom: "JUNO",
            coinMinimalDenom: "ujuno",
            coinDecimals: 6,
            coinGeckoId: "juno-network",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.01,
        average: 0.015,
        high: 0.03,
    }
  },
  {
    chainId: "sommelier-3",
    chainName: "Sommelier",
    rpc: "https://rpc.sommelier-3.quicksilver.zone",
    rest: "https://lcd.sommelier-3.quicksilver.zone",
    
    bip44: {
        coinType: 118,
    },
    bech32Config: {
        bech32PrefixAccAddr: "somm",
        bech32PrefixAccPub: "sommpub",
        bech32PrefixValAddr: "sommvaloper",
        bech32PrefixValPub: "sommvaloperpub",
        bech32PrefixConsAddr: "sommvalcons",
        bech32PrefixConsPub: "sommvalconspub",
    },
    currencies: [
        {
            coinDenom: "SOMM",
            coinMinimalDenom: "usomm",
            coinDecimals: 6,
            coinGeckoId: "sommelier",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "SOMM",
            coinMinimalDenom: "usomm",
            coinDecimals: 6,
            coinGeckoId: "sommelier",
        },
    ],
    stakeCurrency: {
        coinDenom: "SOMM",
        coinMinimalDenom: "usomm",
        coinDecimals: 6,
        coinGeckoId: "sommelier",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.01,
        average: 0.015,
        high: 0.03,
    }
  }
]

