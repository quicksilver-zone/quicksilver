### Keyring

Before creating a validator, you must create an operator key. This will be used to identify your validator in the Quicksilver network. 

```go
$ quicksilverd keys add <key-name> [flags]
```

By default, quicksilver will store keys in your OS-backed keyring. You can change this behavior by specifying the `--keyring-backend` flag.

To import an existing key via a mnemonic - for example if you generated and submitted a genesis transaction, you can provide a `--recover` flag and the `keys add` command will prompt you for the BIP39 mnemonic.

**SECURITY NOTE:** _Keep separate mnemonics and keys for testnet purposes. Never reuse mnemonics or keys associated with live wallets or mainnets. It poses a great security risk to do so!_

Visit the Cosmos SDK's keyring [documentation](https://docs.cosmos.network/v0.43/run-node/keyring.html) for more information.

For a secure keyring setup, using Ledger, you can follow this guide by a community member (approved by our dev team):

[https://github.com/rishisidhu/Quicksilver-guides/blob/main/generating_quicksilver_address.md](https://github.com/rishisidhu/Quicksilver-guides/blob/main/generating_quicksilver_address.md)
