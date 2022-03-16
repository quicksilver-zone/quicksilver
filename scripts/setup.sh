#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/vars.sh

docker-compose down

echo "Removing previous data..."
rm -rf ./${CHAIN_DIR}/$CHAINID_1 &> /dev/null
rm -rf ./${CHAIN_DIR}/$CHAINID_2 &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_2}a &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_2}b &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_2}c &> /dev/null
rm -rf ./${CHAIN_DIR}/hermes &> /dev/null
rm -rf ./${CHAIN_DIR}/icq &> /dev/null

# Add directories for both chains, exit if an error occurs
if ! mkdir -p ./${CHAIN_DIR}/$CHAINID_1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/$CHAINID_2 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_2}a 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_2}b 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_2}c 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/hermes 2>/dev/null; then
    echo "Failed to create hermes folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/icq 2>/dev/null; then
    echo "Failed to create icq folder. Aborting..."
    exit 1
fi

cp ./scripts/config/icq.yaml ./${CHAIN_DIR}/icq/config.yaml

echo "Initializing $CHAINID_1..."
$QS_RUN init test --chain-id $CHAINID_1
echo "Initializing $CHAINID_2..."
$TZ_RUN init test --chain-id $CHAINID_2
echo "Initializing ${CHAINID_2}a..."
$TZ2_RUN init test --chain-id $CHAINID_2
echo "Initializing ${CHAINID_2}b..."
$TZ3_RUN init test --chain-id $CHAINID_2
echo "Initializing ${CHAINID_2}c..."
$TZ4_RUN init test --chain-id $CHAINID_2

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $QS_RUN keys add val1 --recover --keyring-backend=test
echo $VAL_MNEMONIC_2 | $TZ_RUN keys add val2 --recover --keyring-backend=test
echo $VAL_MNEMONIC_3 | $TZ2_RUN keys add val3 --recover --keyring-backend=test
echo $VAL_MNEMONIC_4 | $TZ3_RUN keys add val4 --recover --keyring-backend=test
echo $VAL_MNEMONIC_5 | $TZ4_RUN keys add val5 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_1 | $QS_RUN keys add demowallet1 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_2 | $QS_RUN keys add demowallet2 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_2 | $TZ_RUN keys add demowallet2 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_3 | $TZ2_RUN keys add demowallet3 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_4 | $TZ3_RUN keys add demowallet4 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_5 | $TZ4_RUN keys add demowallet5 --recover --keyring-backend=test
echo $RLY_MNEMONIC_1 | $QS_RUN keys add rly1 --recover --keyring-backend=test
echo $RLY_MNEMONIC_2 | $TZ_RUN keys add rly2 --recover --keyring-backend=test

## Set denoms
sed -i 's/stake/uqck/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}/config/genesis.json
sed -i 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}/config/genesis.json
sed -i 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}a/config/genesis.json
sed -i 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}b/config/genesis.json
sed -i 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}c/config/genesis.json

VAL_ADDRESS_1=$($QS_RUN keys show val1 --keyring-backend test -a)
DEMO_ADDRESS_1=$($QS_RUN keys show demowallet1 --keyring-backend test -a)
RLY_ADDRESS_1=$($QS_RUN keys show rly1 --keyring-backend test -a)

VAL_ADDRESS_2=$($TZ_RUN keys show val2 --keyring-backend test -a)
DEMO_ADDRESS_2=$($TZ_RUN keys show demowallet2 --keyring-backend test -a)
RLY_ADDRESS_2=$($TZ_RUN keys show rly2 --keyring-backend test -a)

VAL_ADDRESS_3=$($TZ2_RUN keys show val3 --keyring-backend test -a)
DEMO_ADDRESS_3=$($TZ2_RUN keys show demowallet3 --keyring-backend test -a)

VAL_ADDRESS_4=$($TZ3_RUN keys show val4 --keyring-backend test -a)
DEMO_ADDRESS_4=$($TZ3_RUN keys show demowallet4 --keyring-backend test -a)

VAL_ADDRESS_5=$($TZ4_RUN keys show val5 --keyring-backend test -a)
DEMO_ADDRESS_5=$($TZ4_RUN keys show demowallet5 --keyring-backend test -a)

VAL_VALOPER_2=$($TZ_RUN keys show val2 --keyring-backend test --bech=val -a)
VAL_VALOPER_3=$($TZ2_RUN keys show val3 --keyring-backend test --bech=val -a)
VAL_VALOPER_4=$($TZ3_RUN keys show val4 --keyring-backend test --bech=val -a)
VAL_VALOPER_5=$($TZ4_RUN keys show val5 --keyring-backend test --bech=val -a)

$QS_RUN add-genesis-account ${VAL_ADDRESS_1} 100000000000uqck
$QS_RUN add-genesis-account ${DEMO_ADDRESS_1} 100000000000uqck
$QS_RUN add-genesis-account ${RLY_ADDRESS_1} 100000000000uqck

$TZ_RUN add-genesis-account ${VAL_ADDRESS_2} 100000000000uatom
$TZ_RUN add-genesis-account ${VAL_ADDRESS_3} 100000000000uatom
$TZ_RUN add-genesis-account ${VAL_ADDRESS_4} 100000000000uatom
$TZ_RUN add-genesis-account ${VAL_ADDRESS_5} 100000000000uatom
$TZ_RUN add-genesis-account ${DEMO_ADDRESS_2} 100000000000uatom
$TZ_RUN add-genesis-account ${DEMO_ADDRESS_3} 100000000000uatom
$TZ_RUN add-genesis-account ${DEMO_ADDRESS_4} 100000000000uatom
$TZ_RUN add-genesis-account ${DEMO_ADDRESS_5} 100000000000uatom
$TZ_RUN add-genesis-account ${RLY_ADDRESS_2} 100000000000uatom

$TZ2_RUN add-genesis-account ${VAL_ADDRESS_3} 100000000000uatom
$TZ3_RUN add-genesis-account ${VAL_ADDRESS_4} 100000000000uatom
$TZ4_RUN add-genesis-account ${VAL_ADDRESS_5} 100000000000uatom

echo "Creating and collecting gentx..."
$QS_RUN gentx val1 7000000000uqck --chain-id $CHAINID_1 --keyring-backend test
$TZ_RUN gentx val2 7000000000uatom --chain-id $CHAINID_2 --keyring-backend test
$TZ2_RUN gentx val3 7000000000uatom --chain-id $CHAINID_2 --keyring-backend test
$TZ3_RUN gentx val4 7000000000uatom --chain-id $CHAINID_2 --keyring-backend test
$TZ4_RUN gentx val5 7000000000uatom --chain-id $CHAINID_2 --keyring-backend test

$QS_RUN collect-gentxs
cp ./${CHAIN_DIR}/${CHAINID_2}a/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_2}/config/gentx/
cp ./${CHAIN_DIR}/${CHAINID_2}b/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_2}/config/gentx/
cp ./${CHAIN_DIR}/${CHAINID_2}c/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_2}/config/gentx/

$TZ_RUN collect-gentxs

node1=$($TZ_RUN tendermint show-node-id)@testzone:26656
node2=$($TZ2_RUN tendermint show-node-id)@testzone2:26656
node3=$($TZ3_RUN tendermint show-node-id)@testzone3:26656
node4=$($TZ4_RUN tendermint show-node-id)@testzone4:26656

echo "Changing defaults and ports in app.toml and config.toml files..."
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
sed -i -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_1}/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_1}/config/app.toml

sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
sed -i -e "s/persistent_peers = \"\"/persistent_peers = \"$node2,$node3,$node4\"/g" ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
sed -i -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}/config/app.toml

sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
sed -i -e "s/persistent_peers = \"\"/persistent_peers = \"$node1,$node3,$node4\"/g" ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
sed -i -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}a/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}a/config/app.toml

sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
sed -i -e "s/persistent_peers = \"\"/persistent_peers = \"$node1,$node2,$node4\"/g" ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
sed -i -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}b/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}b/config/app.toml

sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
sed -i -e "s/persistent_peers = \"\"/persistent_peers = \"$node1,$node2,$node3\"/g" ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
sed -i -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}c/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}c/config/app.toml

## add the message types ICA should allow
jq '.app_state.interchainaccounts.host_genesis_state.params.allow_messages = ["/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend", "/cosmos.staking.v1beta1.MsgDelegate", "/cosmos.staking.v1beta1.MsgRedeemTokensforShares", "/cosmos.staking.v1beta1.MsgTokenizeShares", "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"]' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
jq '.app.state.mint.minter.inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
jq '.app.state.mint.params.max_inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
cp ./${CHAIN_DIR}/${CHAINID_2}{,a}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_2}{,b}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_2}{,c}/config/genesis.json

## set the 'epoch' epoch to 5m interval
jq '.app_state.epochs.epochs = [{"identifier": "epoch","start_time": "0001-01-01T00:00:00Z","duration": "450s","current_epoch": "0","current_epoch_start_time": "0001-01-01T00:00:00Z","epoch_counting_started": false,"current_epoch_start_height": "0"}]' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}

rm -rf ${CHAIN_DIR}/backup
mkdir ${CHAIN_DIR}/backup
cp -fr ${CHAIN_DIR}/${CHAINID_1} ${CHAIN_DIR}/backup/${CHAINID_1}
cp -fr ${CHAIN_DIR}/${CHAINID_2} ${CHAIN_DIR}/backup/${CHAINID_2}
cp -fr ${CHAIN_DIR}/${CHAINID_2}a ${CHAIN_DIR}/backup/${CHAINID_2}a
cp -fr ${CHAIN_DIR}/${CHAINID_2}b ${CHAIN_DIR}/backup/${CHAINID_2}b
cp -fr ${CHAIN_DIR}/${CHAINID_2}c ${CHAIN_DIR}/backup/${CHAINID_2}c

docker-compose down
