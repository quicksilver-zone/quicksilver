export const TestQuickSilverChainInfo = {
    chainId: "rhye-1",
    chainName: "Quicksilver Testnet",
    rpc: "https://rpc.test.quicksilver.zone",
    rest: "https://lcd.test.quicksilver.zone",
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
            coinDenom: "qMUON",
            coinMinimalDenom: "uqmuon",
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
            coinDenom: "qATOM",
            coinMinimalDenom: "uqatom",
            coinDecimals: 6,
            coinGeckoId: "cosmos",
        },
        {
            coinDenom: "qSTARS",
            coinMinimalDenom: "uqstars",
            coinDecimals: 6,
            coinGeckoId: "stargaze",
        },
        {
            coinDenom: "qJUNOX",
            coinMinimalDenom: "uqjunox",
            coinDecimals: 6,
            coinGeckoId: "juno-network",
        },
        {
            coinDenom: "qREGEN",
            coinMinimalDenom: "uqregen",
            coinDecimals: 6,
            coinGeckoId: "regen",
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

export const TestChainInfos = [
    TestQuickSilverChainInfo,
{
    chainId: "fauxgaia-1",
    chainName: "FauxGaia Testnet",
    rpc: "https://rpc.fauxgaia-1.test.quicksilver.zone",
    rest: "https://lcd.fauxgaia-1.test.quicksilver.zone",
    
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
            coinDenom: "MUON",
            coinMinimalDenom: "umuon",
            coinDecimals: 6,
            coinGeckoId: "cosmos",
        },
    ],
    feeCurrencies: [
        {
            coinDenom: "MUON",
            coinMinimalDenom: "umuon",
            coinDecimals: 6,
            coinGeckoId: "cosmos",
        },
    ],
    stakeCurrency: {
        coinDenom: "MUON",
        coinMinimalDenom: "umuon",
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
    chainId: "osmo-test-4",
    chainName: "Osmosis Testnet",
    rpc: "https://rpc.osmo-test-4.test.quicksilver.zone",
    rest: "https://lcd.osmo-test-4.test.quicksilver.zone",
    
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
    chainId: "regen-redwood-1",
    chainName: "Regen Testnet",
    rpc: "https://rpc.regen-redwood-1.test.quicksilver.zone",
    rest: "https://lcd.regen-redwood-1.test.quicksilver.zone",
    
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
    chainId: "elgafar-1",
    chainName: "Stargaze Testnet",
    rpc: "https://rpc.elgafar-1.test.quicksilver.zone",
    rest: "https://lcd.elgafar-1.test.quicksilver.zone",
    
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
    chainId: "theta-testnet-001",
    chainName: "Cosmos Hub Test",
    rpc: "https://rpc.theta-testnet-001.test.quicksilver.zone",
    rest: "https://lcd.theta-testnet-001.test.quicksilver.zone",
    
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
        low: 0.01,
        average: 0.015,
        high: 0.03,
    },
  }, 
  {
    chainId: "uni-6",
    chainName: "Juno Testnet",
    rpc: "https://rpc.uni-6.test.quicksilver.zone",
    rest: "https://lcd.uni-6.test.quicksilver.zone",
    
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
            coinDenom: "JUNOX",
            coinMinimalDenom: "ujunox",
            coinDecimals: 6,
            coinGeckoId: "juno-network",
        },
    ],
    feeCurrencies: [
        {
           coinDenom: "JUNOX",
            coinMinimalDenom: "ujunox",
            coinDecimals: 6,
            coinGeckoId: "juno-network",
        },
    ],
    stakeCurrency: {
        
        coinDenom: "JUNOX",
            coinMinimalDenom: "ujunox",
            coinDecimals: 6,
            coinGeckoId: "juno-network",
    },
    coinType: 118,
    gasPriceStep: {
        low: 0.01,
        average: 0.015,
        high: 0.03,
    }
}
]

