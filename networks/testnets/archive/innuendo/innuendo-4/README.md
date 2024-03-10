## Innuendo-4 Instructions

Innuendo-4 is a replacement for the innuendo-3 long term testnet for Quicksilver; it was created with the intent of onboarding public testnets against a fresh chain, clean of existing state, as close to mainnet scenario as possible. It starts with v0.10.5 of the `quicksilverd` binary which can be downloaded from the releases section of this repository, or as a docker container from quicksilverzone/quicksilver:v0.10.5.

If you sync from genesis, you will need to ensure you follow the upgrade paths can detailed in the 'Innuendo-4' release, [here](https://github.com/ingenuity-build/testnets/releases/tag/v0.10.5)

Source code will be released for self-building binaries prior to mainnet launch. Code is remaining incognito until then for strategic purposes.

Genesis time is set to 2022-12-08T19:00:00Z. The network was started asyncronously. Accounts and validators from previous testnets were dropped. You will need to request tokens from Discord. Genesis can be downloaded from [here](genesis.json).

You may download a recent snapshot of innuendo-4 from here: https://storage.googleapis.com/innuendo-4-snapshots/innuendo-4.tgz (302Mb). This snapshot is using v0.10.8 of quicksilverd.

The current peers are available: 
```
b9b8bb23e61d53ff3b293485d04ea567ebcd7933@65.108.65.94:26656,a94cf3e93cec8eef6d67c2972e4af5eae1a118b2@65.108.2.27:26656,926ce3f8ce4cda6f1a5ee97a937a44f59ff28fbf@65.108.13.176:26656
```
