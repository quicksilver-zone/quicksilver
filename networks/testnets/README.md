# Quicksilver Testnets

All Quicksilver testnets have been and will continue to be named after songs by Freddie Mercury and/or Queen, for no reason more than the guy was a lyrical and musical genius, and there is a somewhat tenuous link between Quicksilver -> Mercury -> Freddie that is ripe for exploitation. 

![Freddie Mercury](https://static.miraheze.org/nonciclopediawiki/thumb/8/84/Freddie_Mercury_simpson.png/200px-Freddie_Mercury_simpson.png)

# Rhye-2

## Existing validators

The second installment of the "Rhye" long-running testnet is here.

Given a number of failed clients on `rhye-1`, it was decided to relaunch the testnet as `rhye-2`.

An export of `rhye-1` was taken as of 2023-12-21 22:00 UTC, and all validators set to jailed status, as to permit an asynchronous relaunch of the network, during the holiday period.

In order to onboard your existing `rhye-1` validator to `rhye-2`, you must do the following:

1. Stop your existing `rhye-1` validator.
2. Build, or download the Quicksilver v1.4.5-rc1 binary from https://github.com/quicksilver-zone/quicksilver; alternatively use the docker container at quicksilverzone/quicksilver:v1.4.5-rc1.

3. Replace your genesis file (usually in ~/.quicksilverd/config/genesis.json) with the genesis.json file found in this repository.
4. Assert the genesis file you downloaded matches the hash provided, using `shasum -a256 ~/.quicksilverd/config/genesis.json`.
4. Run `quicksilverd tendermint unsafe-reset-all`.
5. Restart the process
6. Once synced (check https://rpc.test.quicksilver.zone/status), then unjail your validator, using the following command: `quicksilver tx slashing unjail --from <key>`.

## New validators

In order to onboard your new `rhye-2`, you must do the following:

1. Build, or download the Quicksilver v1.4.5-rc1 binary from https://github.com/quicksilver-zone/quicksilver; alternatively use the docker container at quicksilverzone/quicksilver:v1.4.5-rc1.
2. Run `quicksilverd init <moniker>`.
3. Download the genesis file in this repository to ~/.quicksilverd/config/genesis.json.
4. Assert the genesis file you downloaded matches the hash provided, using `shasum -a256 ~/.quicksilverd/config/genesis.json`.
5. Start the process with `quicksilverd start`. You will probably want to run this as a service.
6. Once synced (check https://rpc.test.quicksilver.zone/status), you will want to create a key, get some funds from the faucet, and create your validator:

#### Create a key
```sh
$ quicksilverd keys add <key>
```

For example:
```sh
$ quicksilverd keys add testkey

- address: quick14xcgnfvmd9xzu5em2gr5d0ykepv4m0y4f4z8lk
  name: testkey
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AvE1BMXRvtydR95jRdrGzWOVpmlC1Uf6V5SazxxFTECa"}'
  type: local


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

random spoil vivid negative wedding moon blast own oxygen fish border project cabbage agent belt dress body absent book tiny myself reflect minimum supreme
```

#### Get funds from a faucet

In the #testnet-faucet channel on discord, write `/faucet <address>` - e.g. `/faucet quick14xcgnfvmd9xzu5em2gr5d0ykepv4m0y4f4z8lk` using the example above,.

Alternatively, try `https://quicksilver-testnet.faucetme.pro/`.

#### Create a validator
```sh
$ quicksilverd tx staking create-validator \
--amount=<amount> \
--pubkey=$(quicksilverd tendermint show-validator) \
--moniker="<moniker>" \
--chain-id="rhye-2" \
--commission-rate="<commission>" \
--commission-max-rate="<max-commission>" \
--commission-max-change-rate="<max-commission-rate-change>" \
--min-self-delegation="<min-self-delegation>" \
--fees=<fees> \
--from=<key>
```

## Peers and Seeds

Peers and Seeds can be found [here](./rhye-2/PEERS.md). Please add your nodes to this list through a PR.

## Docs
Found anything missing or inaccurate? [Create an issue](https://github.com/ingenuity-build/testnets/issues) or make a pull request!
