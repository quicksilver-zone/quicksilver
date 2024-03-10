![Freddie Mercury](https://static.miraheze.org/nonciclopediawiki/thumb/8/84/Freddie_Mercury_simpson.png/200px-Freddie_Mercury_simpson.png)

# Rhapsody Testnet - Ended 23/06/22
The Rhapsody testnet is named for the hit song "Bohemian Rhapsody" written by Freddie Mercury, for his band Queen, back in 1975. For inspiration whilst joining this testnet, please feel free to listen to this in the background: [Bohemian Rhapsody on YouTube](https://www.youtube.com/watch?v=fJ9rUzIMcZQ).

All Quicksilver testnets shall be named after songs by Freddie Mercury and/or Queen, for no reason more than the guy was a lyrical and musical genius, and there is a somewhat tenuous link between Quicksilver -> Mercury -> Freddie that is ripe for exploitation. 

We have added a bunch of scripts to aid your journey; you'll need `make`, `bash`, `git`, `jq`, `gcc` and `go` (v1.17) installed, along with some additional dependencies that will depend on your OS (e.g. `stdlibc++`).

Additional tasks will be added during the coming days.


**If you experience any bugs, issues or problems, please raise an issue here:** https://github.com/ingenuity-build/quicksilver

## Details

 - Chain-ID: `rhapsody-5`
 - Launch Date: 2022-06-16
 - Current Version: `v0.3.0`
 - Genesis File: https://raw.githubusercontent.com/ingenuity-build/testnets/main/rhapsody/genesis.json

### Hardware Requirements
Like any Cosmos-SDK chain, the hardware requirements are pretty modest.
 - 4x CPUs; the faster clock speed the better
 - 8GB RAM
 - 100GB Disk (we are using statesync, so disk requirements are low)
 - Permanent Internet connection (traffic will be minimal during testnet; 10Mbps will be plenty - for production at least 100Mbps is expected)

### Nodes
We are running the following nodes:

 - node01.rhapsody-5.quicksilver.zone:26657
 - node02.rhapsody-5.quicksilver.zone:26657
 - node03.rhapsody-5.quicksilver.zone:26657
 - node04.rhapsody-5.quicksilver.zone:26657

Seeds:

 - dd3460ec11f78b4a7c4336f22a356fe00805ab64@seed.rhapsody-5.quicksilver.zone:26656
 - 8603d0778bfe0a8d2f8eaa860dcdc5eb85b55982@seed.qscosmos-2.quicksilver.zone:27676

## Semi-automated Configuration

```
## clone this repo
git clone https://github.com/ingenuity-build/testnets ## this repo
cd testnets/rhapsody

## download and build quicksilverd and gaiad
make init

## show keys
make keys

## follow instructions listed to get funds from the faucet via discord

## check balances
make balances

## start the validator
make start

## view the logs
make logs

## submit a create-validator tx to start validating (enter your validator name when prompted)
make validate 

## view validators
make show-validators

## view rewards
make check-validator-rewards

## view voting power
make check-validator-voting-power

---

## reset state (sync from zero)
make stop

make reset

---

## clean up time! (post-testnet)
make stop 

make clean
```

### What am I doing wrong?!

#### Unfunded account
```
joe@desktop:~/code/testnets/rhapsody$ make validate
Enter your validator name: my_validator
Error: rpc error: code = NotFound desc = rpc error: code = NotFound desc = account quick1fk9qtycszzk32c3hk8xwjwvkhmkc8rv6gg0xzd not found: key not found
```

Solution: Use discord #qck-tap channel to fund your address (might take a few seconds to come through!)

#### Node not running
```
joe@desktop:~/code/testnets/rhapsody$ make validate
Enter your validator name: my_validator
Error: post failed: Post "http://localhost:26657": dial tcp 127.0.0.1:26657: connect: connection refused
...
```

Solution: Your node is not running; run `make start`. If problems persist, take a look at `make logs` and come find someone on discord!


## Manual Configuration

Download and build Quicksilver:

    git clone https://github.com/ingenuity-build/quicksilver.git --branch v0.3.0
    cd quicksilver
    make build

Testnet configuration script (`touch scripts/testnet_conf.sh`):

    #!/bin/bash -i
    
    set -xe
    
    ### CONFIGURATION ###
    
    CHAIN_ID=rhapsody-5
    
    GENESIS_URL="https://raw.githubusercontent.com/ingenuity-build/testnets/main/rhapsody/genesis.json"
    SEEDS="dd3460ec11f78b4a7c4336f22a356fe00805ab64@seed.rhapsody-5.quicksilver.zone:26656"
    
    BINARY=./build/quicksilverd
    NODE_HOME=$HOME/.quicksilverd
    
    # SET this value for your node:
    NODE_MONIKER="Your_Node"
    
    echo  "Initializing $CHAIN_ID..."
    $BINARY config chain-id $CHAIN_ID --home $NODE_HOME
    $BINARY config keyring-backend test --home $NODE_HOME
    $BINARY config broadcast-mode block --home $NODE_HOME
    $BINARY init $NODE_MONIKER --chain-id $CHAIN_ID --home $NODE_HOME
    
    echo "Get genesis file..."
    curl -sSL $GENESIS_URL > $NODE_HOME/config/genesis.json
    
    echo "Set seeds..."
    sed -i -e "/seeds =/ s/= .*/= \"$SEEDS\"/"  $NODE_HOME/config/config.toml

    echo "Enable pruning..."
    sed -i -e "/pruning =/ s/= .*/= \"everything\"/"  $NODE_HOME/config/app.toml

Run this script from the quicksilver repository main directory;

Remember to make it executable:

    chmod +x scripts/testnet_conf.sh

Then simply run:

    ./scripts/testnet_conf.sh

## Running your node
At this point you can run the node on the CLI with `./build/quicksilverd start` to ensure everything is configured correctly. At this point you may configure your system to run Quicksilver as a system service or daemon.

## Upgrade to Validator

### Test Wallet
To run as a validator you will need to create a QCK wallet:

    ./build/quicksilverd keys add $YOUR_TEST_WALLET --keyring-backend=test

If you already have a test wallet you want to use run (and enter your mnemonic):

    ./build/quicksilverd keys add $YOUR_TEST_WALLET --recover --keyring-backend=test

### Faucet

Join our discord server to access the faucets for QCK and ATOM. Make sure you are in the appropriate channel:

 - **qck-tap** for QCK tokens;
 - **atom-tap** for ATOM tokens;

To check the faucet address:

    $faucet_address rhapsody

To check your balance:

    $balance $YOUR_TEST_WALLET rhapsody

To request a faucet grant:

    $request $YOUR_TEST_WALLET rhapsody

### Validator Tx

Then simply run the tx to upgrade to validator status:

    ## Upgrade node to validator
    ./build/quicksilverd tx staking create-validator \
      --from=$YOUR_TEST_WALLET \
      --amount=1000000uqck \
      --moniker=$NODE_MONIKER \
      --chain-id=$CHAIN_ID \
      --commission-rate=0.1 \
      --commission-max-rate=0.5 \
      --commission-max-change-rate=0.1 \
      --min-self-delegation=1 \
      --pubkey=$($BINARY tendermint show-validator)


## Archived Testnets

## Rhapsody-4

 - Chain-ID: `rhapsody-4`
 - Launch Date: 2022-06-14
 - End Date: 2022-06-15
 - Current Version: `v0.2.0`
 - Genesis File: https://raw.githubusercontent.com/ingenuity-build/testnets/main/rhapsody/rhapsody-4/genesis.json

## Rhapsody (phase 1)

 - Chain-ID: `quicktest-3`
 - Launch Date: 2022-05-02
 - End Date: 2022-06-03
 - Current Version: `v0.1.10`
 - Genesis File: https://raw.githubusercontent.com/ingenuity-build/testnets/main/rhapsody/quicktest-3/genesis.json

