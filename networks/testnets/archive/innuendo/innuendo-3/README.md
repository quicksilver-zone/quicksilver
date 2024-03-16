## Innuendo-3 Instructions

Innuendo-3 is a replacement for the defunct Innuendo-2 long term testnet for Quicksilver. It starts with v0.9.1 of the `quicksilverd` binary which can be downloaded from the releases section of this repository, or as a docker container from quicksilverzone/quicksilver:v0.9.0.

Source code will be released for self-building binaries prior to mainnet launch. Code is remaining incognito until then for strategic purposes.

Genesis time is set to 2022-09-26T18:30:00Z. The network will be started asyncronously. Accounts and validators from innuendo-1 remain in place; do not delete any keys or destroy seeds from innuendo-1. The genesis instructions are as follows:

1. Download the new genesis file from https://raw.githubusercontent.com/ingenuity-build/testnets/main/innuendo/genesis.json.

1. Assert the genesis file state is correct:
```
joe@mac innuendo % shasum -a256 genesis.json
6f97a06cdcfddc5774d4ca4fbee936bc8462b72b74c4337753771fecdfebe93f  genesis.json
```

1. Stop your existing `quicksilverd` service (depends on setup, but often `systemctl stop quicksilver`).

1. Download the new binary from from https://github.com/ingenuity-build/testnets/releases/tag/v0.9.1. Alternatively, pull the docker image at quicksilverzone/quicksilver:v0.9.1.

1. Run `quicksilverd tendermint unsafe-reset-all` to reset the state of your deployment to empty. 

1. Restart your `quicksilverd` service (depends on setup, but often `systemctl start quicksilver`).
