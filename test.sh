# #!/bin/bash
# set -eu

# PATH=build:$PATH

# MONIKER=hiu

# quicksilverd init $MONIKER --chain-id my-test-chain

# quicksilverd keys add my_validator --keyring-backend test

# # Put the generated address in a variable for later use.
# MY_VALIDATOR_ADDRESS=$(quicksilverd keys show my_validator -a --keyring-backend test)

# quicksilverd add-genesis-account $MY_VALIDATOR_ADDRESS 100000000stake

# # Create a gentx.
# quicksilverd gentx my_validator 100000000stake --chain-id my-test-chain --keyring-backend test

# # Add the gentx to the genesis file.
# quicksilverd collect-gentxs

# # quicksilverd start

#!/bin/bash

KEY="test"
CHAINID="quicksilver-testnet-1"
KEYRING="test"
MONIKER="localtestnet"
KEYALGO="secp256k1"
LOGLEVEL="info"

# retrieve all args
WILL_RECOVER=0
WILL_INSTALL=0
WILL_CONTINUE=0
# $# is to check number of arguments
if [ $# -gt 0 ];
then
    # $@ is for getting list of arguments
    for arg in "$@"; do
        case $arg in
        --recover)
            WILL_RECOVER=1
            shift
            ;;
        --install)
            WILL_INSTALL=1
            shift
            ;;
        --continue)
            WILL_CONTINUE=1
            shift
            ;;
        *)
            printf >&2 "wrong argument somewhere"; exit 1;
            ;;
        esac
    done
fi

# continue running if everything is configured
if [ $WILL_CONTINUE -eq 1 ];
then
    # Start the node (remove the --pruning=nothing flag if historical queries are not needed)
    quicksilverd start --pruning=nothing --log_level $LOGLEVEL --minimum-gas-prices=0.0001stake
    exit 1;
fi

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# install quicksilverd if not exist
if [ $WILL_INSTALL -eq 0 ];
then 
    command -v quicksilverd > /dev/null 2>&1 || { echo >&1 "installing quicksilverd"; make install; }
else
    echo >&1 "installing quicksilverd"
    rm -rf $HOME/.quicksilver*
    make install
fi

quicksilverd config keyring-backend $KEYRING
quicksilverd config chain-id $CHAINID

# determine if user wants to recorver or create new
if [ $WILL_RECOVER -eq 0 ];
then
    quicksilverd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO
else
    quicksilverd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover
fi

echo >&1 "\n"

# init chain
quicksilverd init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to stake
cat $HOME/.quicksilver/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="stake"' > $HOME/.quicksilver/config/tmp_genesis.json && mv $HOME/.quicksilver/config/tmp_genesis.json $HOME/.quicksilver/config/genesis.json
cat $HOME/.quicksilver/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="stake"' > $HOME/.quicksilver/config/tmp_genesis.json && mv $HOME/.quicksilver/config/tmp_genesis.json $HOME/.quicksilver/config/genesis.json
cat $HOME/.quicksilver/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="stake"' > $HOME/.quicksilver/config/tmp_genesis.json && mv $HOME/.quicksilver/config/tmp_genesis.json $HOME/.quicksilver/config/genesis.json
cat $HOME/.quicksilver/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="stake"' > $HOME/.quicksilver/config/tmp_genesis.json && mv $HOME/.quicksilver/config/tmp_genesis.json $HOME/.quicksilver/config/genesis.json

# Set gas limit in genesis
# cat $HOME/.quicksilver/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.quicksilver/config/tmp_genesis.json && mv $HOME/.quicksilver/config/tmp_genesis.json $HOME/.quicksilver/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
quicksilverd add-genesis-account $KEY 1000000000000stake --keyring-backend $KEYRING

# Sign genesis transaction
quicksilverd gentx $KEY 1000000stake --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
quicksilverd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
quicksilverd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
# quicksilverd start --pruning=nothing --log_level $LOGLEVEL --minimum-gas-prices=0.0001stake --p2p.laddr tcp://0.0.0.0:2280 --rpc.laddr tcp://0.0.0.0:2281 --grpc.address 0.0.0.0:2282 --grpc-web.address 0.0.0.0:2283
