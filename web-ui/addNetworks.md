# Adding Networks

This document will guide you through the process of adding networks to the web UI.

## Table of Contents

- [.env](#.env)
- [Files](#files)
- [Queries](#queries)
- [Components](#components)

### .env

You will need to add a couple of environment variables to the `.env` file in the root of the project. These variables are used to configure the networks that are available to the user and the endpoints that the web UI will use to interact with the blockchain.

**Example .env file additions**

```bash
NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_NETWORK_NAME="REST_ENDPOINT" # change `NETWORK_NAME`, `REST_ENDPOINT` to the name & endpoint of the network you are adding.
NEXT_PUBLIC_MAINNET_RPC_ENDPOINT_NETWORK_NAME="RPC_ENDPOINT" # change `NETWORK_NAME` and `RPC_ENDPOINT` to the name & endpoint of the network you are adding.

NEXT_PUBLIC_WHITELISTED_DENOM="uatom,ustars,uosmo,usomm,uregen,ujuno,udydx,NEW_UDENOM" # change `NEW_UDENOM` to the appropriate denom for the network you are adding.
NEXT_PUBLIC_WHITELISTED_ZONES="osmosis-1,stargaze-1,regen-1,cosmoshub-4,sommelier-3,juno-1,dydx-mainnet-1,NEW_CHAIN_ID" # change `NEW_CHAIN_ID` to the appropriate chain_id for the network you are adding.

NEXT_PUBLIC_CHAIN_NAME_CHAIN_ID="CHAIN_ID" # change `CHAIN_ID` to the appropriate chain_id for the network you are adding.
```

### Files

There are various files that will require updates to add a new network to the web UI. These files include:

- `hooks/useGrpcQueryClient.ts`
- `hooks/useRpcQueryClient.ts`
- `pages/_app.tsx`
-

### Queries & Signing

The web UI uses tanstack react-query for data fetching and cosmology for query client building. You will need to add a new entry in `hooks/useGrpcQueryClient.ts` & `hooks/useRpcQueryClient.ts` in order to fetch data for your additional network.

Here we rely on the entries in the `.env` file to determine which endpoint to use for the network we are fetching data for so be sure to use the correct endpoint or the query will be broken.

**Example useGrpcQueryClient.ts**

```typescript
const endpoints: { [key: string]: string | undefined } = {
  quicksilver:
    env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_QUICKSILVER : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_QUICKSILVER,
  cosmoshub:
    env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_COSMOSHUB : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_COSMOSHUB,
  sommelier:
    env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_SOMMELIER : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_SOMMELIER,
  stargaze:
    env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_STARGAZE : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_STARGAZE,
  regen: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_REGEN : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_REGEN,
  osmosis: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_OSMOSIS : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_OSMOSIS,
  juno: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_JUNO : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_JUNO,
  dydx: env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_DYDX : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_DYDX,
  NEW_NETWORK:
    env === 'testnet' ? process.env.NEXT_PUBLIC_TESTNET_LCD_ENDPOINT_NEW_NETWORK : process.env.NEXT_PUBLIC_MAINNET_LCD_ENDPOINT_NEW_NETWORK,
};
```

**Example useRpcQueryClient.ts**

```typescript
const endpoints: { [key: string]: string | undefined } = {
  quicksilver: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_QUICKSILVER : process.env.MAINNET_RPC_ENDPOINT_QUICKSILVER,
  cosmoshub: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_COSMOSHUB : process.env.MAINNET_RPC_ENDPOINT_COSMOSHUB,
  sommelier: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_SOMMELIER : process.env.MAINNET_RPC_ENDPOINT_SOMMELIER,
  stargaze: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_STARGAZE : process.env.MAINNET_RP_ENDPOINTC_STARGAZE,
  regen: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_REGEN : process.env.MAINNET_RPC_ENDPOINT_REGEN,
  osmosis: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_OSMOSIS : process.env.MAINNET_RPC_ENDPOINT_OSMOSIS,
  juno: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_JUNO : process.env.MAINNET_RPC_ENDPOINT_JUNO,
  dydx: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_DYDX : process.env.MAINNET_RPC_ENDPOINT_DYDX,
  NEW_NETWORK: env === 'testnet' ? process.env.TESTNET_RPC_ENDPOINT_NEW_NETWORK : process.env.MAINNET_RPC_ENDPOINT_NEW_NETWORK,
};
```

For signing we must update the `pages/_app.tsx` file to include the new network in the `ChainProvider` `endpointOptions` array. This array is used to determine which endpoints to use.

```typescript
            endpointOptions={{
              isLazy: true,
              endpoints: {
                quicksilver: {
                  rpc: [rpcEndpoints.quicksilver ?? ''],
                  rest: [lcdEndpoints.quicksilver ?? ''],
                },
                quicksilvertestnet: {
                  rest: ['https://lcd.test.quicksilver.zone/'],
                  rpc: ['https://rpc.test.quicksilver.zone'],
                },
                cosmoshub: {
                  rpc: [rpcEndpoints.cosmoshub ?? ''],
                  rest: [lcdEndpoints.cosmoshub ?? ''],
                },
                sommelier: {
                  rpc: [rpcEndpoints.sommelier ?? ''],
                  rest: [lcdEndpoints.sommelier ?? ''],
                },
                stargaze: {
                  rpc: [rpcEndpoints.stargaze ?? ''],
                  rest: [lcdEndpoints.stargaze ?? ''],
                },
                regen: {
                  rpc: [rpcEndpoints.regen ?? ''],
                  rest: [lcdEndpoints.regen ?? ''],
                },
                osmosis: {
                  rpc: [rpcEndpoints.osmosis ?? ''],
                  rest: [lcdEndpoints.osmosis ?? ''],
                },
                osmosistestnet: {
                  rpc: [rpcEndpoints.osmosis ?? ''],
                  rest: [lcdEndpoints.osmosis ?? ''],
                },
                umee: {
                  rpc: ['https://rpc-umee-ia.cosmosia.notional.ventures/'],
                  rest: ['https://api-umee-ia.cosmosia.notional.ventures/'],
                },
                dydx: {
                  rpc: [rpcEndpoints.dydx ?? ''],
                  rest: [lcdEndpoints.dydx ?? ''],
                },
                NEW_NETWORK: {
                  rpc: [rpcEndpoints.NEW_NETWORK ?? ''],
                  rest: [lcdEndpoints.NEW_NETWORK ?? ''],
                },
              },
            }}
```

### Components
