export const networks = [
    {
      value: 'ATOM',
      logo: '/quicksilver-app-v2/img/networks/atom.svg',
      qlogo: '/quicksilver-app-v2/img/networks/qatom.svg',
      name: 'Cosmos Hub',
      chainName: 'cosmoshub',
      chainId: 'cosmoshub-4',
    },
    {
      value: 'OSMO',
      logo: '/quicksilver-app-v2/img/networks/osmosis.svg',
      qlogo: '/quicksilver-app-v2/img/networks/qosmo.svg',
      name: 'Osmosis',
      chainName: 'osmosis',
      chainId: 'osmosis-1',
    },
    {
      value: 'STARS',
      logo: '/quicksilver-app-v2/img/networks/stargaze.svg',
      qlogo: '/quicksilver-app-v2/img/networks/qstars.svg',
      name: 'Stargaze',
      chainName: 'stargaze',
      chainId: 'stargaze-1',
    },
    {
      value: 'REGEN',
      logo: '/quicksilver-app-v2/img/networks/regen.svg',
      qlogo: '/quicksilver-app-v2/img/networks/regen.svg',
      name: 'Regen',
      chainName: 'regen',
      chainId: 'regen-1',
    },
    {
      value: 'SOMM',
      logo: '/quicksilver-app-v2/img/networks/sommelier.png',
      qlogo: '/quicksilver-app-v2/img/networks/sommelier.png',
      name: 'Sommelier',
      chainName: 'sommelier',
      chainId: 'sommelier-3',
    },
  ];

export const ProdQuickSilverChainInfo = {
    chainId: "quicksilver-2",
    chainName: "Quicksilver Protocol",
    rpc: "https://rpc-quicksilver-ia.cosmosia.notional.ventures",
    rest: "https://api-quicksilver-ia.cosmosia.notional.ventures",
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
        rpc: "https://rpc-cosmoshub-ia.cosmosia.notional.ventures",
        rest: "https://api-cosmoshub-ia.cosmosia.notional.ventures",

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
        rpc: "https://rpc-regen-ia.cosmosia.notional.ventures",
        rest: "https://api-regen-ia.cosmosia.notional.ventures",
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
        rpc: "https://rpc-osmosis-ia.cosmosia.notional.ventures",
        rest: "https://api-osmosis-ia.cosmosia.notional.ventures",
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
        rpc: "https://rpc-stargaze-ia.cosmosia.notional.ventures",
        rest: "https://api-stargaze-ia.cosmosia.notional.ventures",

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
        rpc: "https://rpc-juno-ia.cosmosia.notional.ventures",
        rest: "https://api-juno-ia.cosmosia.notional.ventures",

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
        rpc: "https://sommelier-rpc.polkachu.com",
        rest: "https://sommelier-api.polkachu.com",

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

export const ProdZoneInfos = [
    {
        name: "Cosmoshub",
        connection_id: "connection-1",
        chain_id: "cosmoshub-4",
        deposit_address: {
            address: "cosmos1pl2ld9d0x9ve6flklr9jv9sl69y6yucvaelcljml83l7kyn0m0ksacp9tj",
            port_name: "icacontroller-cosmoshub-4.deposit",
            withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
            balance_waitgroup: 0
        },
        withdrawal_address: {
            address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
            port_name: "icacontroller-cosmoshub-4.withdrawal",
            withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
            balance_waitgroup: 0
        },
        performance_address: {
            address: "cosmos1v8p98feyknvf9crjp3yu4ywr3v4y7j3zd9qqxpxjxgaa5phy2t6qamlsga",
            port_name: "icacontroller-cosmoshub-4.performance",
            withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
            balance_waitgroup: 0
        },
        delegation_address: {
            address: "cosmos1st3fng2vjcpz5lhg46un94zg0vn3nj658wc0chc57z29hx8zqeys6mvxdd",
            port_name: "icacontroller-cosmoshub-4.delegate",
            withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
            balance_waitgroup: 0
        },
        account_prefix: "cosmos",
        local_denom: "uqatom",
        base_denom: "uatom",
        local_logo: "/assets/qAtom.svg",
        base_logo: "/assets/Cosmos.png",
    },
    {
        name: "Osmosis",
        connection_id: "connection-2",
        chain_id: "osmosis-1",
        deposit_address: {
            address: "osmo1hw0z99c3ykmfj5pp40wud9ymfp5qwgre55u0hzjyym6zzmux90es6x7pqh",

            port_name: "icacontroller-osmosis-1.deposit",
            withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
            balance_waitgroup: 0
        },
        withdrawal_address: {
            address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",

            port_name: "icacontroller-osmosis-1.withdrawal",
            withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
            balance_waitgroup: 0
        },
        performance_address: {
            address: "osmo1laez676yk5lhtujzp86equz0vl0wu23uk29t7cne6258nk6a3mdsexkaqg",

            port_name: "icacontroller-osmosis-1.performance",
            withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
            balance_waitgroup: 0
        },
        delegation_address: {
            address: "osmo1trpqkzdzprkregdez3xs5mf9w3me76gyu83dnhc2fjxuy9g70unqs5v62q",

            port_name: "icacontroller-osmosis-1.delegate",
            withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
            balance_waitgroup: 0
        },
        account_prefix: "osmo",
        local_denom: "uqosmo",
        base_denom: "uosmo",
        local_logo: "/assets/qOsmo.svg",
        base_logo: "/assets/Osmosis.png",
    },
    {
        name: "Regen",
        connection_id: "connection-9",
        chain_id: "regen-1",
        deposit_address: {
            address: "regen1cs3t5zcn3mqvjah48sk3lp7depe0679qdzkhqlyurwvmpdvjd9yslk6395",

            port_name: "icacontroller-regen-1.deposit",
            withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
            balance_waitgroup: 0
        },
        withdrawal_address: {
            address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",

            port_name: "icacontroller-regen-1.withdrawal",
            withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
            balance_waitgroup: 0
        },
        performance_address: {
            address: "regen1jeszepf9e64pul32qmc2ynr0eqjm7s8w33gk0zsva6vhpe2qhu5symxer2",

            port_name: "icacontroller-regen-1.performance",
            withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
            balance_waitgroup: 0
        },
        delegation_address: {
            address: "regen1huzrs3q3yq2dhq837k3ssrn9d5ysccr7qn3plms20rh666ydryfqnqhwtn",
            port_name: "icacontroller-regen-1.delegate",
            withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
            balance_waitgroup: 0
        },
        account_prefix: "regen",
        local_denom: "uqregen",
        base_denom: "uregen",
        local_logo: "/assets/qRegen.svg",
        base_logo: "/assets/Regen.png",
    },
    {
        name: "Sommelier",
        connection_id: "connection-54",
        chain_id: "sommelier-3",
        deposit_address: {
            address: "somm1g82kh7waf2lgwtlw6j369m5vahvjuuqdjmjdfh5qstp8tphx8lyqdpyevs",

            port_name: "icacontroller-sommelier-3.deposit",
            withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
            balance_waitgroup: 0
        },
        withdrawal_address: {
            address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",

            port_name: "icacontroller-sommelier-3.withdrawal",
            withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
            balance_waitgroup: 0
        },
        performance_address: {
            address: "somm1ysypkxg8znklg6tw703qsm7g8u3qasx5ndr650xy8xk8t9l2zq4sku2azn",

            port_name: "icacontroller-sommelier-3.performance",
            withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
            balance_waitgroup: 0
        },
        delegation_address: {
            address: "somm1mgk9zw6q5ckus3yr0rfjh86smgpg64sjylq0tqxx4928dmxpwjvqp4y64g",

            port_name: "icacontroller-sommelier-3.delegate",
            withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
            balance_waitgroup: 0
        },
        account_prefix: "somm",
        local_denom: "uqsomm",
        base_denom: "usomm",
        local_logo: "/assets/qSomm.svg",
        base_logo: "/assets/sommelier.png",
    },
    {
        name: "Stargaze",
        connection_id: "connection-0",
        chain_id: "stargaze-1",
        deposit_address: {
            address: "stars16k9qkq57kpwcnzawd8u0utl6u2zh5mr2dz7qp3wy7ywx9m7xkdaqnq5msn",

            port_name: "icacontroller-stargaze-1.deposit",
            withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
            balance_waitgroup: 0
        },
        withdrawal_address: {
            address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",

            port_name: "icacontroller-stargaze-1.withdrawal",
            withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
            balance_waitgroup: 0
        },
        performance_address: {
            address: "stars1rqd5082y3lxlx04g77s5m8v7f4hq33zcjsgvk6lcadk289xn9grqvn9zjm",

            port_name: "icacontroller-stargaze-1.performance",
            withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
            balance_waitgroup: 0
        },
        delegation_address: {
            address: "stars1rqeychen93f9j72r7jue56g2gvaaqtjl7ule7u09f3ytxtlxhe8s3zhd98",
            port_name: "icacontroller-stargaze-1.delegate",
            withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
            balance_waitgroup: 0
        },
        account_prefix: "stars",
        local_denom: "uqstars",
        base_denom: "ustars",
        local_logo: "/assets/qSTAR.svg",
        base_logo: "/assets/Stargaze.png",
    }
]

export const ProdDataMap = {
    uatom: {
        local_logo: "/assets/qAtom.svg",
        base_logo: "/assets/Cosmos.png",
        local_symbol: "qATOM",
        base_symbol: "ATOM",
        pool_id: "944",
        chainlist_prefix: "cosmos",
        decimals: 6,
        network_name: "Cosmos",
        symbol: "ATOM",
        zone: {
            connection_id: "connection-1",
            chain_id: "cosmoshub-4",
            deposit_address: {
                address: "cosmos1pl2ld9d0x9ve6flklr9jv9sl69y6yucvaelcljml83l7kyn0m0ksacp9tj",
                port_name: "icacontroller-cosmoshub-4.deposit",
                withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
                balance_waitgroup: 0
            },
            withdrawal_address: {
                address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
                port_name: "icacontroller-cosmoshub-4.withdrawal",
                withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
                balance_waitgroup: 0
            },
            performance_address: {
                address: "cosmos1v8p98feyknvf9crjp3yu4ywr3v4y7j3zd9qqxpxjxgaa5phy2t6qamlsga",
                port_name: "icacontroller-cosmoshub-4.performance",
                withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
                balance_waitgroup: 0
            },
            delegation_address: {
                address: "cosmos1st3fng2vjcpz5lhg46un94zg0vn3nj658wc0chc57z29hx8zqeys6mvxdd",
                port_name: "icacontroller-cosmoshub-4.delegate",
                withdrawal_address: "cosmos1ets9pef83fltsrgvvwkzw055wck9wagsjkyfqgef09g3jn5sy8espdph07",
                balance_waitgroup: 0
            },
            account_prefix: "cosmos",
            local_denom: "uqatom",
            base_denom: "uatom",
            local_ibc_denom: "ibc/32B1E5958441B955D176EE7691EB25CEEA1002D1A9E4A4A897161114FF6ED008",
            base_ibc_denom: "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9",
        },
        network: {
            chainId: "cosmoshub-4",
            chainName: "CosmosHub",
            rpc: "https://rpc-cosmoshub-ia.cosmosia.notional.ventures",
            rest: "https://api-cosmoshub-ia.cosmosia.notional.ventures",

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
        }
    },
    uosmo: {
        local_logo: "/assets/qOsmo.svg",
        base_logo: "/assets/Osmosis.png",
        network_name: "Osmosis",
        local_symbol: "qOSMO",
        base_symbol: "OSMO",
        pool_id: "956",
        chainlist_prefix: "osmosis",
        decimals: 6,
        zone: {
            connection_id: "connection-2",
            chain_id: "osmosis-1",
            deposit_address: {
                address: "osmo1hw0z99c3ykmfj5pp40wud9ymfp5qwgre55u0hzjyym6zzmux90es6x7pqh",

                port_name: "icacontroller-osmosis-1.deposit",
                withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
                balance_waitgroup: 0
            },
            withdrawal_address: {
                address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
                port_name: "icacontroller-osmosis-1.withdrawal",
                withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
                balance_waitgroup: 0
            },
            performance_address: {
                address: "osmo1laez676yk5lhtujzp86equz0vl0wu23uk29t7cne6258nk6a3mdsexkaqg",
                port_name: "icacontroller-osmosis-1.performance",
                withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
                balance_waitgroup: 0
            },
            delegation_address: {
                address: "osmo1trpqkzdzprkregdez3xs5mf9w3me76gyu83dnhc2fjxuy9g70unqs5v62q",

                port_name: "icacontroller-osmosis-1.delegate",
                withdrawal_address: "osmo1ckchpf8xc822qyy7alfknrvux24zx382jqsnwfyxuxk44ljsefgsa2m9x4",
                balance_waitgroup: 0
            },
            account_prefix: "osmo",
            local_denom: "uqosmo",
            base_denom: "uosmo",
            local_ibc_denom: "ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC",
            base_ibc_denom: "ibc/6CE8E927869E764D11651D0E498FDF532963F6B8BFAC13943C458224DB3F88B9",
        },
        network: {
            chainId: "osmosis-1",
            chainName: "Osmosis",
            rpc: "https://rpc-osmosis-ia.cosmosia.notional.ventures",
            rest: "https://api-osmosis-ia.cosmosia.notional.ventures",
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
        }
    },
    uregen: {
        local_logo: "/assets/qRegen.svg",
        base_logo: "/assets/Regen.png",
        network_name: "Regen",
        local_symbol: "qREGEN",
        base_symbol: "REGEN",
        pool_id: "948",
        chainlist_prefix: "regen",
        decimals: 6,
        zone: {
            connection_id: "connection-9",
            chain_id: "regen-1",
            deposit_address: {
                address: "regen1cs3t5zcn3mqvjah48sk3lp7depe0679qdzkhqlyurwvmpdvjd9yslk6395",

                port_name: "icacontroller-regen-1.deposit",
                withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
                balance_waitgroup: 0
            },
            withdrawal_address: {
                address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",

                port_name: "icacontroller-regen-1.withdrawal",
                withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
                balance_waitgroup: 0
            },
            performance_address: {
                address: "regen1jeszepf9e64pul32qmc2ynr0eqjm7s8w33gk0zsva6vhpe2qhu5symxer2",

                port_name: "icacontroller-regen-1.performance",
                withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
                balance_waitgroup: 0
            },
            delegation_address: {
                address: "regen1huzrs3q3yq2dhq837k3ssrn9d5ysccr7qn3plms20rh666ydryfqnqhwtn",
                port_name: "icacontroller-regen-1.delegate",
                withdrawal_address: "regen1pwxxqkpwa6stxvfesfy03nh5jxdxmlvxcnhmqadkdktgwjt7kphs3mf8f5",
                balance_waitgroup: 0
            },
            account_prefix: "regen",
            local_denom: "uqregen",
            base_denom: "uregen",
            local_ibc_denom: "ibc/B2DCA297A3AFF98480BCCCC962E1D00A3BBE06A37136D3FABD16DC8FB19451E1",
            base_ibc_denom: "ibc/A7E38774F447445DB94A8ED00BEE78EFC43EED7A732D314D3F7F4AB743993E9F",
        },
        network: {
            chainId: "regen-1",
            chainName: "Regen Mainnet",
            rpc: "https://rpc-regen-ia.cosmosia.notional.ventures",
            rest: "https://api-regen-ia.cosmosia.notional.ventures",
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
        }
    },
    usomm: {
        local_logo: "/assets/qSomm.png",
        base_logo: "/assets/sommelier.png",
        network_name: "Sommelier",
        local_symbol: "qSOMM",
        base_symbol: "SOMM",
        pool_id: "1087",
        chainlist_prefix: "sommelier",
        decimals: 6,
        zone: {
            connection_id: "connection-54",
            chain_id: "sommelier-3",
            deposit_address: {
                address: "somm1g82kh7waf2lgwtlw6j369m5vahvjuuqdjmjdfh5qstp8tphx8lyqdpyevs",

                port_name: "icacontroller-sommelier-3.deposit",
                withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
                balance_waitgroup: 0
            },
            withdrawal_address: {
                address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",

                port_name: "icacontroller-sommelier-3.withdrawal",
                withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
                balance_waitgroup: 0
            },
            performance_address: {
                address: "somm1ysypkxg8znklg6tw703qsm7g8u3qasx5ndr650xy8xk8t9l2zq4sku2azn",

                port_name: "icacontroller-sommelier-3.performance",
                withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
                balance_waitgroup: 0
            },
            delegation_address: {
                address: "somm1mgk9zw6q5ckus3yr0rfjh86smgpg64sjylq0tqxx4928dmxpwjvqp4y64g",

                port_name: "icacontroller-sommelier-3.delegate",
                withdrawal_address: "somm1vp5qaedjvccxnmwyppqhu8xvt7w5hsenwvqhw0lx8dxskk4udhaqmmg36m",
                balance_waitgroup: 0
            },
            account_prefix: "somm",
            local_denom: "uqsomm",
            base_denom: "usomm",
            local_ibc_denom: "ibc/9C634C3B5AD926FB709CB6F6F5435B8D5B42C5ED7B47D3ABA433868FB47C5A8B",
            base_ibc_denom: "ibc/BFF8BC09B94E2EA90B64961A181D2383280FFA7847109DE1AB4ECA366466462A",
        },
        network: {
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
    },
    ustars: {
        local_logo: "/assets/qSTAR.svg",
        base_logo: "/assets/Stargaze.png",
        network_name: "Stargaze",
        local_symbol: "qSTARS",
        base_symbol: "STARS",
        chainlist_prefix: "stargaze",
        pool_id: "903",
        decimals: 6,
        zone: {
            connection_id: "connection-0",
            chain_id: "stargaze-1",
            deposit_address: {
                address: "stars16k9qkq57kpwcnzawd8u0utl6u2zh5mr2dz7qp3wy7ywx9m7xkdaqnq5msn",

                port_name: "icacontroller-stargaze-1.deposit",
                withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
                balance_waitgroup: 0
            },
            withdrawal_address: {
                address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",

                port_name: "icacontroller-stargaze-1.withdrawal",
                withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
                balance_waitgroup: 0
            },
            performance_address: {
                address: "stars1rqd5082y3lxlx04g77s5m8v7f4hq33zcjsgvk6lcadk289xn9grqvn9zjm",

                port_name: "icacontroller-stargaze-1.performance",
                withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
                balance_waitgroup: 0
            },
            delegation_address: {
                address: "stars1rqeychen93f9j72r7jue56g2gvaaqtjl7ule7u09f3ytxtlxhe8s3zhd98",
                port_name: "icacontroller-stargaze-1.delegate",
                withdrawal_address: "stars1t0d7kv3mnc42s2y7fvw47fpcvh39qdgyguccts6n9xejtcar8kzqe5qkce",
                balance_waitgroup: 0
            },
            account_prefix: "stars",
            local_denom: "uqstars",
            base_denom: "ustars",
            local_ibc_denom: "ibc/46E27FBBC56A14AD0029678BB34A4164F650AA3711EEDEA0D05E08DB41D13BF0",
            base_ibc_denom: "ibc/49BAE4CD2172833F14000627DA87ED8024AD46A38D6ED33F6239F22B5832F958",
        },
        network: {
            chainId: "stargaze-1",
            chainName: "Stargaze",
            rpc: "https://rpc-stargaze-ia.cosmosia.notional.ventures",
            rest: "https://api-stargaze-ia.cosmosia.notional.ventures",

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
        }
    },
}