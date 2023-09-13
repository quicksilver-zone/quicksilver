# QuicksilverJS

<p align="center">
  <img src="https://user-images.githubusercontent.com/545047/197298013-6e9ac5ed-3c03-4f52-b6e5-eee87cb76b0e.png" width="120"><br />
    Quicksilver is the Cosmos Liquid Staking Zone.
</p>


## install

```sh
npm install quicksilverjs
```
## Table of contents

- [QuicksilverJS](#quicksilverjs)
  - [Install](#install)
  - [Table of contents](#table-of-contents)
- [Usage](#usage)
    - [RPC Clients](#rpc-clients)
    - [Composing Messages](#composing-messages)
        - quicksilver
            - [Airdrop](#airdop-messages)
            - [Interchain Staking](#interchain-staking-messages)
            - [Interchain Query](#interchain-query-messages)
            - [Participation Rewards](#participation-rewards-messages)
            - [Tokenfactory](#tokenfactory-messages)
        - Cosmos and IBC
            - [IBC](#ibc-messages)
            - [Cosmos](#cosmos-messages)
- [Wallets and Signers](#connecting-with-wallets-and-signing-messages)
    - [Stargate Client](#initializing-the-stargate-client)
    - [Creating Signers](#creating-signers)
    - [Broadcasting Messages](#broadcasting-messages)
- [Advanced Usage](#advanced-usage)
- [Developing](#developing)
- [Credits](#credits)

## Usage

### RPC Clients

```js
import { quicksilver } from 'quicksilverjs';

const { createRPCQueryClient } = quicksilver.ClientFactory; 
const client = await createRPCQueryClient({ rpcEndpoint: RPC_ENDPOINT });

// now you can query the cosmos modules
const balance = await client.cosmos.bank.v1beta1
    .allBalances({ address: 'quick1addresshere' });

// you can also query the quicksilver modules
const data = await client.quicksilver.interchainstaking.v1.delegations({
    chainId
});
```

### Composing Messages

Import the `quicksilver` object from `quicksilverjs`. 

```js
import { quicksilver } from 'quicksilverjs';

const {
    submitQueryResponse
} = quicksilver.interchainquery.v1.MessageComposer.withTypeUrl
```

#### Airdop Messages

```js
const {
    claim
} = quicksilver.airdrop.v1.MessageComposer.withTypeUrl
```

#### Interchain Staking Messages

```js
const {
    requestRedemption,
    signalIntent
} = quicksilver.interchainstaking.v1.MessageComposer.withTypeUrl;
```

#### Interchain Query Messages

```js
const {
    submitQueryResponse
} = quicksilver.interchainquery.v1.MessageComposer.withTypeUrl;
```

#### Participation Rewards Messages

```js
const {
    submitClaim
} = quicksilver.participationrewards.v1.MessageComposer.withTypeUrl
```
#### Tokenfactory Messages

```js
const {
    burn,
    changeAdmin,
    createDenom,
    mint,
    setDenomMetadata
} = quicksilver.tokenfactory.v1beta1.MessageComposer.withTypeUrl;
```

#### IBC Messages

```js
import { ibc } from 'quicksilverjs';

const {
    transfer
} = ibc.applications.transfer.v1.MessageComposer.withTypeUrl
```

#### Cosmos Messages

```js
import { cosmos } from 'quicksilverjs';

const {
    fundCommunityPool,
    setWithdrawAddress,
    withdrawDelegatorReward,
    withdrawValidatorCommission
} = cosmos.distribution.v1beta1.MessageComposer.fromPartial;

const {
    multiSend,
    send
} = cosmos.bank.v1beta1.MessageComposer.fromPartial;

const {
    beginRedelegate,
    createValidator,
    delegate,
    editValidator,
    undelegate
} = cosmos.staking.v1beta1.MessageComposer.fromPartial;

const {
    deposit,
    submitProposal,
    vote,
    voteWeighted
} = cosmos.gov.v1beta1.MessageComposer.fromPartial;
```

## Connecting with Wallets and Signing Messages

⚡️ For web interfaces, we recommend using [cosmos-kit](https://github.com/cosmology-tech/cosmos-kit). Continue below to see how to manually construct signers and clients.

Here are the docs on [creating signers](https://github.com/cosmology-tech/cosmos-kit/tree/main/packages/react#signing-clients) in cosmos-kit that can be used with Keplr and other wallets.

### Initializing the Stargate Client

Use `getSigningQuicksilverClient` to get your `SigningStargateClient`, with the proto/amino messages full-loaded. No need to manually add amino types, just require and initialize the client:

```js
import { getSigningQuicksilverClient } from 'quicksilverjs';

const stargateClient = await getSigningQuicksilverClient({
  rpcEndpoint,
  signer // OfflineSigner
});
```
### Creating Signers

To broadcast messages, you can create signers with a variety of options:

* [cosmos-kit](https://github.com/cosmology-tech/cosmos-kit/tree/main/packages/react#signing-clients) (recommended)
* [keplr](https://docs.keplr.app/api/cosmjs.html)
* [cosmjs](https://gist.github.com/webmaster128/8444d42a7eceeda2544c8a59fbd7e1d9)
### Amino Signer

Likely you'll want to use the Amino, so unless you need proto, you should use this one:

```js
import { getOfflineSignerAmino as getOfflineSigner } from 'cosmjs-utils';
```
### Proto Signer

```js
import { getOfflineSignerProto as getOfflineSigner } from 'cosmjs-utils';
```

WARNING: NOT RECOMMENDED TO USE PLAIN-TEXT MNEMONICS. Please take care of your security and use best practices such as AES encryption and/or methods from 12factor applications.

```js
import { chains } from 'chain-registry';

const mnemonic =
  'unfold client turtle either pilot stock floor glow toward bullet car science';
  const chain = chains.find(({ chain_name }) => chain_name === 'quicksilver');
  const signer = await getOfflineSigner({
    mnemonic,
    chain
  });
```
### Broadcasting Messages

Now that you have your `stargateClient`, you can broadcast messages:

```js
const { send } = cosmos.bank.v1beta1.MessageComposer.withTypeUrl;

const msg = send({
    amount: [
    {
        denom: 'quick',
        amount: '1000'
    }
    ],
    toAddress: address,
    fromAddress: address
});

const fee: StdFee = {
    amount: [
    {
        denom: 'quick',
        amount: '864'
    }
    ],
    gas: '86364'
};
const response = await stargateClient.signAndBroadcast(address, [msg], fee);
```

## Advanced Usage


If you want to manually construct a stargate client

```js
import { OfflineSigner, GeneratedType, Registry } from "@cosmjs/proto-signing";
import { AminoTypes, SigningStargateClient } from "@cosmjs/stargate";

import { 
    cosmosAminoConverters,
    cosmosProtoRegistry,
    ibcProtoRegistry,
    ibcAminoConverters,
    quicksilverAminoConverters,
    quicksilverProtoRegistry
} from 'quicksilverjs';

const signer: OfflineSigner = /* create your signer (see above)  */
const rpcEndpint = 'https://rpc.cosmos.directory/quicksilver'; // or another URL

const protoRegistry: ReadonlyArray<[string, GeneratedType]> = [
    ...cosmosProtoRegistry,
    ...ibcProtoRegistry,
    ...quicksilverProtoRegistry
];

const aminoConverters = {
    ...cosmosAminoConverters,
    ...ibcAminoConverters,
    ...quicksilverAminoConverters
};

const registry = new Registry(protoRegistry);
const aminoTypes = new AminoTypes(aminoConverters);

const stargateClient = await SigningStargateClient.connectWithSigner(rpcEndpoint, signer, {
    registry,
    aminoTypes
});
```

## Developing

When first cloning the repo:

```
yarn
yarn build
```

### Codegen

Contract schemas live in `./contracts`, and protos in `./proto`. Look inside of `scripts/codegen.js` and configure the settings for bundling your SDK and contracts into `quicksilverjs`:

```
yarn codegen
```

### Publishing

Build the types and then publish:

```
yarn build:ts
yarn publish
```
## Credits

🛠 Built by Cosmology — if you like our tools, please consider delegating to [our validator ⚛️](https://cosmology.tech/validator)

Code built with the help of these related projects:

* [@osmonauts/telescope](https://github.com/osmosis-labs/telescope) a "babel for the Cosmos", Telescope is a TypeScript Transpiler for Cosmos Protobufs.
* [cosmos-kit](https://github.com/cosmology-tech/cosmos-kit) A wallet connector for the Cosmos ⚛️

## Disclaimer

AS DESCRIBED IN THE LICENSES, THE SOFTWARE IS PROVIDED “AS IS”, AT YOUR OWN RISK, AND WITHOUT WARRANTIES OF ANY KIND.

No developer or entity involved in creating this software will be liable for any claims or damages whatsoever associated with your use, inability to use, or your interaction with other users of the code, including any direct, indirect, incidental, special, exemplary, punitive or consequential damages, or loss of profits, cryptocurrencies, tokens, or anything else of value.
