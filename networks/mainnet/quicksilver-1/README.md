# Quicksilver Mainnet joining instructions

** This directory contains the files and instructions for quicksilver-1. These instructions are deprecated and should not be used.**
** Please refer to documentation in the root of this repository.**

## Minimum hardware requirements

- 4 cores (max. clock speed possible)
- 16GB RAM
- 500GB+ of NVMe or SSD disk

## Software requirements

Quicksilver v1.0.0 will be released once a final audit report has been released.

### Install Quicksilver

Requires [Go version v1.19+](https://golang.org/doc/install).

```sh
> git clone https://github.com/ingenuity-build/quicksilver && cd quicksilver
> git fetch origin --tags
> git checkout v1.0.0
> make install
```

#### Verify installation

To verify if the installation was successful, execute the following command:

```sh
> quicksilverd version --long
```

It will display the version of quicksilverd currently installed:

```sh
name: quicksilverd
server_name: quicksilverd
version: 1.0.0
commit: XXX
build_tags: netgo,ledger
go: go version go1.19 linux/amd64
```

## Genesis validators

## Create a validator

1. Init Chain and start your node

   ```sh
   > quicksilverd init <moniker-name> --chain-id=quicksilver-1
   ```

2. Create a local key pair
   **Note: we recommend _only_ using Ledger for mainnet! Key security is important!**

   ```sh
   > ## create a new key:
   > quicksilverd keys add <key-name>
   > ## or use a ledger:
   > quicksilverd key add <key-name> --ledger
   > ## or import an old key:
   > quicksilverd keys show <key-name> -a
   ```

3. Download genesis
   Fetch `genesis.json` into `quicksilverd`'s `config` directory (default: ~/.quicksilverd)

   ```sh
   > curl -s https://raw.githubusercontent.com/ingenuity-build/mainnet/main/genesis/genesis.json > genesis.json
   ```

   **Genesis sha256**

   ```sh
    shasum -a 256 ~/.quicksilverd/config/genesis.json
    8bfc3aa7a81eb8c1a2452bdb8d256b372ecfdd67c634b4f63846f755ef4dd815  /home/<user>/.quicksilverd/config/genesis.json
   ```

4. Start your node and sync to the latest block

5. Create validator

   ```sh
   $ quicksilverd tx staking create-validator \
   --amount 50000000uqck \
   --commission-max-change-rate "0.1" \
   --commission-max-rate "0.20" \
   --commission-rate "0.1" \
   --min-self-delegation "1" \
   --details "a short description lives here" \
   --pubkey=$(quicksilverd tendermint show-validator) \
   --security-contact "youremail@goes.here" \
   --moniker <your_moniker> \
   --chain-id quicksilver-1 \
   --from <key-name>
   ```
