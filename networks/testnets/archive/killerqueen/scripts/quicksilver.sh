#!/bin/bash -i
set -e
SRC_DIR=$(pwd)
source vars.sh

### BUILD ###
mkdir -p $DIR

if [[ ! -x $DIR/quicksilver ]]; then
    git clone https://github.com/ingenuity-build/quicksilver --branch=$QS_VERSION $DIR/quicksilver
fi

cd $DIR/quicksilver
make build

### CONFIGURATION ###

GENESIS_URL="https://raw.githubusercontent.com/ingenuity-build/testnets/killerqueen/killerqueen/genesis.json"
SEEDS="dd3460ec11f78b4a7c4336f22a356fe00805ab64@seed.$CHAIN_ID.quicksilver.zone:26656,8603d0778bfe0a8d2f8eaa860dcdc5eb85b55982@seed02.killerqueen-1.quicksilver.zone:27656"

# SET this value for your node:
NODE_MONIKER="$(hostname -f)"

### OPTIONAL STATE ###

# set height
# INTERVAL=1500
# LATEST_HEIGHT=$(curl -s http://seed.$CHAIN_ID.quicksilver.zone:26657/block | jq -r .result.block.header.height);
# BLOCK_HEIGHT=$(($(($LATEST_HEIGHT / $INTERVAL)) * $INTERVAL));
# if [ $BLOCK_HEIGHT -eq 0 ]; then
#   echo "Error: Cannot state sync to block 0; Latest block is $LATEST_HEIGHT and must be at least $INTERVAL; wait a few blocks!"
#   exit 1
# fi

#TRUST_HASH=$(curl -s "http://seed.killerqueen-1.quicksilver.zone:26657/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)
#echo "Trust hash: $TRUST_HASH"
#SYNC_RPC="node02.killerqueen-1.quicksilver.zone:26657,node03.killerqueen-1.quicksilver.zone:26657,node04.killerqueen-1.quicksilver.zone:26657"

echo  "Initializing $CHAIN_ID..."
$QS_BIN config chain-id $CHAIN_ID --home $QS_HOME
$QS_BIN config keyring-backend test --home $QS_HOME
$QS_BIN config broadcast-mode block --home $QS_HOME
$QS_BIN init $NODE_MONIKER --chain-id $CHAIN_ID --home $QS_HOME -o

echo "Get genesis file..."
#curl -sSL $GENESIS_URL > $QS_HOME/config/genesis.json
cp $SRC_DIR/../genesis.json $QS_HOME/config/genesis.json
#echo  "Enabling state sync..."
#sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true| ; \
#s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"$SYNC_RPC\"| ; \
#s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1$BLOCK_HEIGHT| ; \
#s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"$TRUST_HASH\"|" $QS_HOME/config/config.toml

echo "Set seeds..."
sed -i -e "/seeds =/ s/= .*/= \"$SEEDS\"/"  $QS_HOME/config/config.toml

$QS_BIN --home $QS_HOME tendermint unsafe-reset-all

echo ""
echo "----------------------------------------------------------------"
echo " ! PLEASE NOTE ! "
echo ""
echo "To restore an existing mnemonic, run the following command:"
echo ">> $QS_BIN --home $QS_HOME keys add validator --recover"
echo ""
echo "alternatively, run the following command to generate a new seed:"
echo ">> $QS_BIN --home $QS_HOME keys add validator"
echo
echo "Additionally, if you have a priv_validator_key.json file from"
echo "when you submitted your gentx; ensure that you replace the file"
echo "located at $QS_HOME/config/priv_validator_key.json"
echo "----------------------------------------------------------------"

