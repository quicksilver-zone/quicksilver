#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/vars.sh

docker-compose down

echo "Removing previous data..."
rm -rf ./${CHAIN_DIR}/$CHAINID_0 &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_0}a &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_0}b &> /dev/null

rm -rf ./${CHAIN_DIR}/$CHAINID_1 &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_1}a &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_1}b &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_1}c &> /dev/null

rm -rf ./${CHAIN_DIR}/$CHAINID_2 &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_2}a &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_2}b &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_2}c &> /dev/null

rm -rf ./${CHAIN_DIR}/hermes &> /dev/null
rm -rf ./${CHAIN_DIR}/icq &> /dev/null

# Add directories for both chains, exit if an error occurs
#chain_0
if ! mkdir -p ./${CHAIN_DIR}/$CHAINID_0 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_0}a 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_0}b 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

#chain_1
if ! mkdir -p ./${CHAIN_DIR}/$CHAINID_1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_1}a 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_1}b 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/${CHAINID_1}c 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    #chain_2
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
fi

#relayers
if ! mkdir -p ./${CHAIN_DIR}/hermes 2>/dev/null; then
    echo "Failed to create hermes folder. Aborting..."
    exit 1
fi

if ! mkdir -p ./${CHAIN_DIR}/icq 2>/dev/null; then
    echo "Failed to create icq folder. Aborting..."
    exit 1
fi

cp ./scripts/config/icq.yaml ./${CHAIN_DIR}/icq/config.yaml

echo "Initializing $CHAINID_0..."
$QS1_RUN init test --chain-id $CHAINID_0
echo "Initializing ${CHAINID_0}a..."
$QS2_RUN init test --chain-id $CHAINID_0
echo "Initializing ${CHAINID_0}b..."
$QS3_RUN init test --chain-id $CHAINID_0

echo "Initializing $CHAINID_1..."
$TZ1_1_RUN init test --chain-id $CHAINID_1
echo "Initializing ${CHAINID_1}a..."
$TZ1_2_RUN init test --chain-id $CHAINID_1
echo "Initializing ${CHAINID_1}b..."
$TZ1_3_RUN init test --chain-id $CHAINID_1
echo "Initializing ${CHAINID_1}c..."
$TZ1_4_RUN init test --chain-id $CHAINID_1

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    echo "Initializing $CHAINID_2..."
    $TZ2_1_RUN init test --chain-id $CHAINID_2
    echo "Initializing ${CHAINID_2}a..."
    $TZ2_2_RUN init test --chain-id $CHAINID_2
    echo "Initializing ${CHAINID_2}b..."
    $TZ2_3_RUN init test --chain-id $CHAINID_2
    echo "Initializing ${CHAINID_2}c..."
    $TZ2_4_RUN init test --chain-id $CHAINID_2
fi

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $QS1_RUN keys add val1 --recover --keyring-backend=test

echo $VAL_MNEMONIC_2 | $TZ1_1_RUN keys add val2 --recover --keyring-backend=test
echo $VAL_MNEMONIC_3 | $TZ1_2_RUN keys add val3 --recover --keyring-backend=test
echo $VAL_MNEMONIC_4 | $TZ1_3_RUN keys add val4 --recover --keyring-backend=test
echo $VAL_MNEMONIC_5 | $TZ1_4_RUN keys add val5 --recover --keyring-backend=test

echo $VAL_MNEMONIC_6 | $QS2_RUN keys add val6 --recover --keyring-backend=test
echo $VAL_MNEMONIC_7 | $QS3_RUN keys add val7 --recover --keyring-backend=test

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    echo $VAL_MNEMONIC_8 | $TZ2_1_RUN keys add val8 --recover --keyring-backend=test
    echo $VAL_MNEMONIC_9 | $TZ2_2_RUN keys add val9 --recover --keyring-backend=test
    echo $VAL_MNEMONIC_10 | $TZ2_3_RUN keys add val10 --recover --keyring-backend=test
    echo $VAL_MNEMONIC_11 | $TZ2_4_RUN keys add val11 --recover --keyring-backend=test
fi

echo $DEMO_MNEMONIC_1 | $QS1_RUN keys add demowallet1 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_2 | $QS1_RUN keys add demowallet2 --recover --keyring-backend=test

echo $DEMO_MNEMONIC_2 | $TZ1_1_RUN keys add demowallet2 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_3 | $TZ1_2_RUN keys add demowallet3 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_4 | $TZ1_3_RUN keys add demowallet4 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_5 | $TZ1_4_RUN keys add demowallet5 --recover --keyring-backend=test

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    echo $DEMO_MNEMONIC_8 | $TZ2_1_RUN keys add demowallet8 --recover --keyring-backend=test
    echo $DEMO_MNEMONIC_9 | $TZ2_2_RUN keys add demowallet9 --recover --keyring-backend=test
    echo $DEMO_MNEMONIC_10 | $TZ2_3_RUN keys add demowallet10 --recover --keyring-backend=test
    echo $DEMO_MNEMONIC_11 | $TZ2_4_RUN keys add demowallet11 --recover --keyring-backend=test
fi

echo $DEMO_MNEMONIC_6 | $QS2_RUN keys add demowallet6 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_7 | $QS3_RUN keys add demowallet7 --recover --keyring-backend=test

echo $RLY_MNEMONIC_1 | $QS1_RUN keys add rly1 --recover --keyring-backend=test
echo $RLY_MNEMONIC_2 | $TZ1_1_RUN keys add rly2 --recover --keyring-backend=test
if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    echo $RLY_MNEMONIC_3 | $TZ2_1_RUN keys add rly3 --recover --keyring-backend=test
fi

## Set denoms
${SED} 's/stake/uqck/g' $(pwd)/${CHAIN_DIR}/${CHAINID_0}/config/genesis.json
${SED} 's/stake/uqck/g' $(pwd)/${CHAIN_DIR}/${CHAINID_0}a/config/genesis.json
${SED} 's/stake/uqck/g' $(pwd)/${CHAIN_DIR}/${CHAINID_0}b/config/genesis.json

${SED} 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}/config/genesis.json
${SED} 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}a/config/genesis.json
${SED} 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}b/config/genesis.json
${SED} 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}c/config/genesis.json

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    ${SED} 's/stake/uosmo/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}/config/genesis.json
    ${SED} 's/stake/uosmo/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}a/config/genesis.json
    ${SED} 's/stake/uosmo/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}b/config/genesis.json
    ${SED} 's/stake/uosmo/g' $(pwd)/${CHAIN_DIR}/${CHAINID_2}c/config/genesis.json
fi

VAL_ADDRESS_1=$($QS1_RUN keys show val1 --keyring-backend test -a)
DEMO_ADDRESS_1=$($QS1_RUN keys show demowallet1 --keyring-backend test -a)
RLY_ADDRESS_1=$($QS1_RUN keys show rly1 --keyring-backend test -a)

VAL_ADDRESS_2=$($TZ1_1_RUN keys show val2 --keyring-backend test -a)
DEMO_ADDRESS_2=$($TZ1_1_RUN keys show demowallet2 --keyring-backend test -a)
RLY_ADDRESS_2=$($TZ1_1_RUN keys show rly2 --keyring-backend test -a)

VAL_ADDRESS_3=$($TZ1_2_RUN keys show val3 --keyring-backend test -a)
DEMO_ADDRESS_3=$($TZ1_2_RUN keys show demowallet3 --keyring-backend test -a)

VAL_ADDRESS_4=$($TZ1_3_RUN keys show val4 --keyring-backend test -a)
DEMO_ADDRESS_4=$($TZ1_3_RUN keys show demowallet4 --keyring-backend test -a)

VAL_ADDRESS_5=$($TZ1_4_RUN keys show val5 --keyring-backend test -a)
DEMO_ADDRESS_5=$($TZ1_4_RUN keys show demowallet5 --keyring-backend test -a)

VAL_VALOPER_2=$($TZ1_1_RUN keys show val2 --keyring-backend test --bech=val -a)
VAL_VALOPER_3=$($TZ1_2_RUN keys show val3 --keyring-backend test --bech=val -a)
VAL_VALOPER_4=$($TZ1_3_RUN keys show val4 --keyring-backend test --bech=val -a)
VAL_VALOPER_5=$($TZ1_4_RUN keys show val5 --keyring-backend test --bech=val -a)

VAL_ADDRESS_6=$($QS2_RUN keys show val6 --keyring-backend test -a)
DEMO_ADDRESS_6=$($QS2_RUN keys show demowallet6 --keyring-backend test -a)

VAL_ADDRESS_7=$($QS3_RUN keys show val7 --keyring-backend test -a)
DEMO_ADDRESS_7=$($QS3_RUN keys show demowallet7 --keyring-backend test -a)

VAL_VALOPER_6=$($QS2_RUN keys show val6 --keyring-backend test --bech=val -a)
VAL_VALOPER_7=$($QS3_RUN keys show val7 --keyring-backend test --bech=val -a)

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    VAL_ADDRESS_8=$($TZ2_1_RUN keys show val8 --keyring-backend test -a)
    DEMO_ADDRESS_8=$($TZ2_1_RUN keys show demowallet8 --keyring-backend test -a)
    # rly3 ?
    RLY_ADDRESS_3=$($TZ2_1_RUN keys show rly3 --keyring-backend test -a)

    VAL_ADDRESS_9=$($TZ2_2_RUN keys show val9 --keyring-backend test -a)
    DEMO_ADDRESS_9=$($TZ2_2_RUN keys show demowallet9 --keyring-backend test -a)

    VAL_ADDRESS_10=$($TZ2_3_RUN keys show val10 --keyring-backend test -a)
    DEMO_ADDRESS_10=$($TZ2_3_RUN keys show demowallet10 --keyring-backend test -a)

    VAL_ADDRESS_11=$($TZ2_4_RUN keys show val11 --keyring-backend test -a)
    DEMO_ADDRESS_11=$($TZ2_4_RUN keys show demowallet11 --keyring-backend test -a)

    VAL_VALOPER_8=$($TZ2_1_RUN keys show val8 --keyring-backend test --bech=val -a)
    VAL_VALOPER_9=$($TZ2_2_RUN keys show val9 --keyring-backend test --bech=val -a)
    VAL_VALOPER_10=$($TZ2_3_RUN keys show val10 --keyring-backend test --bech=val -a)
    VAL_VALOPER_11=$($TZ2_4_RUN keys show val11 --keyring-backend test --bech=val -a)
fi

$QS1_RUN add-genesis-account ${VAL_ADDRESS_1} 100000000000uqck
$QS1_RUN add-genesis-account ${DEMO_ADDRESS_1} 100000000000uqck
$QS1_RUN add-genesis-account ${RLY_ADDRESS_1} 100000000000uqck

$QS1_RUN add-genesis-account ${VAL_ADDRESS_6} 100000000000uqck
$QS1_RUN add-genesis-account ${VAL_ADDRESS_7} 100000000000uqck
$QS1_RUN add-genesis-account ${DEMO_ADDRESS_6} 100000000000uqck
$QS1_RUN add-genesis-account ${DEMO_ADDRESS_7} 100000000000uqck

$QS2_RUN add-genesis-account ${VAL_ADDRESS_6} 100000000000uqck
$QS3_RUN add-genesis-account ${VAL_ADDRESS_7} 100000000000uqck

$TZ1_1_RUN add-genesis-account ${VAL_ADDRESS_2} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${VAL_ADDRESS_3} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${VAL_ADDRESS_4} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${VAL_ADDRESS_5} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${DEMO_ADDRESS_2} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${DEMO_ADDRESS_3} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${DEMO_ADDRESS_4} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${DEMO_ADDRESS_5} 100000000000uatom
$TZ1_1_RUN add-genesis-account ${RLY_ADDRESS_2} 100000000000uatom

$TZ1_2_RUN add-genesis-account ${VAL_ADDRESS_3} 100000000000uatom
$TZ1_3_RUN add-genesis-account ${VAL_ADDRESS_4} 100000000000uatom
$TZ1_4_RUN add-genesis-account ${VAL_ADDRESS_5} 100000000000uatom

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    $TZ2_1_RUN add-genesis-account ${VAL_ADDRESS_8} 100000000000uosmo
    $TZ2_1_RUN add-genesis-account ${VAL_ADDRESS_9} 100000000000uosmo
    $TZ2_1_RUN add-genesis-account ${VAL_ADDRESS_10} 100000000000uosmo
    $TZ2_1_RUN add-genesis-account ${VAL_ADDRESS_11} 100000000000uosmo
    $TZ2_1_RUN add-genesis-account ${DEMO_ADDRESS_8} 100000000000uosmo
    $TZ2_1_RUN add-genesis-account ${DEMO_ADDRESS_9} 100000000000uosmo
    $TZ2_1_RUN add-genesis-account ${DEMO_ADDRESS_10} 100000000000uosmo
    $TZ2_1_RUN add-genesis-account ${DEMO_ADDRESS_11} 100000000000uosmo
    # rly3 ?
    $TZ2_1_RUN add-genesis-account ${RLY_ADDRESS_3} 100000000000uosmo

    $TZ2_2_RUN add-genesis-account ${VAL_ADDRESS_9} 100000000000uosmo
    $TZ2_3_RUN add-genesis-account ${VAL_ADDRESS_10} 100000000000uosmo
    $TZ2_4_RUN add-genesis-account ${VAL_ADDRESS_11} 100000000000uosmo
fi

echo "Creating and collecting gentx..."
$QS1_RUN gentx val1 7000000000uqck --chain-id $CHAINID_0 --keyring-backend test
$QS2_RUN gentx val6 7000000000uqck --chain-id $CHAINID_0 --keyring-backend test
$QS3_RUN gentx val7 7000000000uqck --chain-id $CHAINID_0 --keyring-backend test

$TZ1_1_RUN gentx val2 7000000000uatom --commission-rate 0.33 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_1 --keyring-backend test
$TZ1_2_RUN gentx val3 6000000000uatom --commission-rate 0.23 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_1 --keyring-backend test
$TZ1_3_RUN gentx val4 5000000000uatom --commission-rate 0.13 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_1 --keyring-backend test
$TZ1_4_RUN gentx val5 4000000000uatom --commission-rate 0.03 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_1 --keyring-backend test

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    $TZ2_1_RUN gentx val8 7000000000uosmo --commission-rate 0.33 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_2 --keyring-backend test
    $TZ2_2_RUN gentx val9 6000000000uosmo --commission-rate 0.23 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_2 --keyring-backend test
    $TZ2_3_RUN gentx val10 5000000000uosmo --commission-rate 0.13 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_2 --keyring-backend test
    $TZ2_4_RUN gentx val11 4000000000uosmo --commission-rate 0.03 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_2 --keyring-backend test
fi

cp ./${CHAIN_DIR}/${CHAINID_0}a/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_0}/config/gentx/
cp ./${CHAIN_DIR}/${CHAINID_0}b/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_0}/config/gentx/

$QS1_RUN collect-gentxs

cp ./${CHAIN_DIR}/${CHAINID_1}a/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_1}/config/gentx/
cp ./${CHAIN_DIR}/${CHAINID_1}b/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_1}/config/gentx/
cp ./${CHAIN_DIR}/${CHAINID_1}c/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_1}/config/gentx/

$TZ1_1_RUN collect-gentxs

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    cp ./${CHAIN_DIR}/${CHAINID_2}a/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_2}/config/gentx/
    cp ./${CHAIN_DIR}/${CHAINID_2}b/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_2}/config/gentx/
    cp ./${CHAIN_DIR}/${CHAINID_2}c/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_2}/config/gentx/

    $TZ2_1_RUN collect-gentxs
fi

node1=$($TZ1_1_RUN tendermint show-node-id)@testzone1-1:26656
node2=$($TZ1_2_RUN tendermint show-node-id)@testzone1-2:26656
node3=$($TZ1_3_RUN tendermint show-node-id)@testzone1-3:26656
node4=$($TZ1_4_RUN tendermint show-node-id)@testzone1-4:26656

node5=$($QS1_RUN tendermint show-node-id)@quicksilver:26656
node6=$($QS2_RUN tendermint show-node-id)@quicksilver2:26656
node7=$($QS3_RUN tendermint show-node-id)@quicksilver3:26656

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    node8=$($TZ2_1_RUN tendermint show-node-id)@testzone2-1:26656
    node9=$($TZ2_2_RUN tendermint show-node-id)@testzone2-2:26656
    node10=$($TZ2_3_RUN tendermint show-node-id)@testzone2-3:26656
    node11=$($TZ2_4_RUN tendermint show-node-id)@testzone2-4:26656
fi

echo "Changing defaults and ports in app.toml and config.toml files..."
${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_0}/config/config.toml
${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_0}/config/config.toml
${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_0}/config/config.toml
${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_0}/config/config.toml
${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node6,$node7\"/g" ${CHAIN_DIR}/${CHAINID_0}b/config/config.toml
${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_0}/config/app.toml
${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_0}/config/app.toml

${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_0}a/config/config.toml
${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_0}a/config/config.toml
${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_0}a/config/config.toml
${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_0}a/config/config.toml
${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node5,$node7\"/g" ${CHAIN_DIR}/${CHAINID_0}a/config/config.toml
${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_0}a/config/app.toml
${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_0}a/config/app.toml

${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_0}b/config/config.toml
${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_0}b/config/config.toml
${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_0}b/config/config.toml
${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_0}b/config/config.toml
${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node5,$node6\"/g" ${CHAIN_DIR}/${CHAINID_0}b/config/config.toml
${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_0}b/config/app.toml
${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_0}b/config/app.toml

${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node2,$node3,$node4\"/g" ${CHAIN_DIR}/${CHAINID_1}/config/config.toml
${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_1}/config/app.toml
${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_1}/config/app.toml

${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_1}a/config/config.toml
${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}a/config/config.toml
${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}a/config/config.toml
${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_1}a/config/config.toml
${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node1,$node3,$node4\"/g" ${CHAIN_DIR}/${CHAINID_1}a/config/config.toml
${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_1}a/config/app.toml
${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_1}a/config/app.toml

${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_1}b/config/config.toml
${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}b/config/config.toml
${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}b/config/config.toml
${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_1}b/config/config.toml
${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node1,$node2,$node4\"/g" ${CHAIN_DIR}/${CHAINID_1}b/config/config.toml
${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_1}b/config/app.toml
${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_1}b/config/app.toml

${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_1}c/config/config.toml
${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}c/config/config.toml
${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_1}c/config/config.toml
${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_1}c/config/config.toml
${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node1,$node2,$node3\"/g" ${CHAIN_DIR}/${CHAINID_1}c/config/config.toml
${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_1}c/config/app.toml
${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_1}c/config/app.toml

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    ${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
    ${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
    ${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
    ${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
    ${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node9,$node10,$node11\"/g" ${CHAIN_DIR}/${CHAINID_2}/config/config.toml
    ${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}/config/app.toml
    ${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}/config/app.toml

    ${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
    ${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
    ${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
    ${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
    ${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node8,$node10,$node11\"/g" ${CHAIN_DIR}/${CHAINID_2}a/config/config.toml
    ${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}a/config/app.toml
    ${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}a/config/app.toml

    ${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
    ${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
    ${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
    ${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
    ${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node8,$node9,$node11\"/g" ${CHAIN_DIR}/${CHAINID_2}b/config/config.toml
    ${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}b/config/app.toml
    ${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}b/config/app.toml

    ${SED} -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
    ${SED} -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
    ${SED} -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
    ${SED} -e 's/index_all_keys = false/index_all_keys = true/g' ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
    ${SED} -e "s/persistent_peers = \"\"/persistent_peers = \"$node8,$node9,$node10\"/g" ${CHAIN_DIR}/${CHAINID_2}c/config/config.toml
    ${SED} -e 's/enable = false/enable = true/g' ${CHAIN_DIR}/${CHAINID_2}c/config/app.toml
    ${SED} -e 's/swagger = false/swagger = true/g' ${CHAIN_DIR}/${CHAINID_2}c/config/app.toml
fi

## add the message types ICA should allow
jq '.app_state.interchainaccounts.host_genesis_state.params.allow_messages = ["/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend", "/cosmos.staking.v1beta1.MsgDelegate", "/cosmos.staking.v1beta1.MsgRedeemTokensforShares", "/cosmos.staking.v1beta1.MsgTokenizeShares", "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", "/ibc.applications.transfer.v1.MsgTransfer", "/cosmos.staking.v1beta1.MsgUndelegate"]' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
jq '.app.state.mint.minter.inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
jq '.app.state.mint.params.max_inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
jq '.app_state.staking.params.unbonding_time = "300s"' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
cp ./${CHAIN_DIR}/${CHAINID_1}{,a}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_1}{,b}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_1}{,c}/config/genesis.json

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    jq '.app_state.interchainaccounts.host_genesis_state.params.allow_messages = ["/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend", "/cosmos.staking.v1beta1.MsgDelegate", "/cosmos.staking.v1beta1.MsgRedeemTokensforShares", "/cosmos.staking.v1beta1.MsgTokenizeShares", "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", "/ibc.applications.transfer.v1.MsgTransfer", "/cosmos.staking.v1beta1.MsgUndelegate"]' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
    jq '.app_state.interchainaccounts.host_genesis_state.params.host_enabled = true' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
    jq '.app_state.interchainaccounts.host_genesis_state.port = "icahost"' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
    jq '.app.state.mint.minter.inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
    jq '.app.state.mint.params.max_inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_2}/config/genesis.json{.new,}
    cp ./${CHAIN_DIR}/${CHAINID_2}{,a}/config/genesis.json
    cp ./${CHAIN_DIR}/${CHAINID_2}{,b}/config/genesis.json
    cp ./${CHAIN_DIR}/${CHAINID_2}{,c}/config/genesis.json
fi

## set the 'epoch' epoch to 5m interval
jq '.app_state.epochs.epochs = [{"identifier": "epoch","start_time": "0001-01-01T00:00:00Z","duration": "450s","current_epoch": "0","current_epoch_start_time": "0001-01-01T00:00:00Z","epoch_counting_started": false,"current_epoch_start_height": "0"}]' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.interchainstaking.params.deposit_interval = 25' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.mint.params.epoch_identifier = "epoch"' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.gov.deposit_params.min_deposit = [{"denom": "uqck", "amount": "100"}]' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.gov.deposit_params.max_deposit_period = "30s"' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.gov.voting_params.voting_period = "20s"' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}

cp ./${CHAIN_DIR}/${CHAINID_0}{,a}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_0}{,b}/config/genesis.json

rm -rf ${CHAIN_DIR}/backup
mkdir ${CHAIN_DIR}/backup
cp -fr ${CHAIN_DIR}/${CHAINID_0} ${CHAIN_DIR}/backup/${CHAINID_0}
cp -fr ${CHAIN_DIR}/${CHAINID_0}a ${CHAIN_DIR}/backup/${CHAINID_0}a
cp -fr ${CHAIN_DIR}/${CHAINID_0}b ${CHAIN_DIR}/backup/${CHAINID_0}b
cp -fr ${CHAIN_DIR}/${CHAINID_1} ${CHAIN_DIR}/backup/${CHAINID_1}
cp -fr ${CHAIN_DIR}/${CHAINID_1}a ${CHAIN_DIR}/backup/${CHAINID_1}a
cp -fr ${CHAIN_DIR}/${CHAINID_1}b ${CHAIN_DIR}/backup/${CHAINID_1}b
cp -fr ${CHAIN_DIR}/${CHAINID_1}c ${CHAIN_DIR}/backup/${CHAINID_1}c

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
    cp -fr ${CHAIN_DIR}/${CHAINID_2} ${CHAIN_DIR}/backup/${CHAINID_2}
    cp -fr ${CHAIN_DIR}/${CHAINID_2}a ${CHAIN_DIR}/backup/${CHAINID_2}a
    cp -fr ${CHAIN_DIR}/${CHAINID_2}b ${CHAIN_DIR}/backup/${CHAINID_2}b
    cp -fr ${CHAIN_DIR}/${CHAINID_2}c ${CHAIN_DIR}/backup/${CHAINID_2}c
fi

docker-compose down

if [ "$IS_MULTI_ZONE_TEST" = true ]; then
cat << EOF > ${SCRIPT_DIR}/wallets.sh
VAL_ADDRESS_1=$VAL_ADDRESS_1
DEMO_ADDRESS_1=$DEMO_ADDRESS_1
RLY_ADDRESS_1=$RLY_ADDRESS_1
VAL_ADDRESS_6=$VAL_ADDRESS_6
DEMO_ADDRESS_6=$DEMO_ADDRESS_6
VAL_ADDRESS_7=$VAL_ADDRESS_7
DEMO_ADDRESS_7=$DEMO_ADDRESS_7
VAL_ADDRESS_2=$VAL_ADDRESS_2
DEMO_ADDRESS_2=$DEMO_ADDRESS_2
RLY_ADDRESS_2=$RLY_ADDRESS_2
VAL_ADDRESS_3=$VAL_ADDRESS_3
DEMO_ADDRESS_3=$DEMO_ADDRESS_3

VAL_ADDRESS_4=$VAL_ADDRESS_4
DEMO_ADDRESS_4=$DEMO_ADDRESS_4

VAL_ADDRESS_5=$VAL_ADDRESS_5
DEMO_ADDRESS_5=$DEMO_ADDRESS_5

VAL_VALOPER_2=$VAL_VALOPER_2
VAL_VALOPER_3=$VAL_VALOPER_3
VAL_VALOPER_4=$VAL_VALOPER_4
VAL_VALOPER_5=$VAL_VALOPER_5

VAL_ADDRESS_8=$VAL_ADDRESS_8
DEMO_ADDRESS_8=$DEMO_ADDRESS_8

VAL_ADDRESS_9=$VAL_ADDRESS_9
DEMO_ADDRESS_9=$DEMO_ADDRESS_9

VAL_ADDRESS_10=$VAL_ADDRESS_10
DEMO_ADDRESS_10=$DEMO_ADDRESS_10

VAL_ADDRESS_11=$VAL_ADDRESS_11
DEMO_ADDRESS_11=$DEMO_ADDRESS_11

VAL_VALOPER_8=$VAL_VALOPER_8
VAL_VALOPER_9=$VAL_VALOPER_9
VAL_VALOPER_10=$VAL_VALOPER_10
VAL_VALOPER_11=$VAL_VALOPER_11
EOF
else
cat << EOF > ${SCRIPT_DIR}/wallets.sh
VAL_ADDRESS_1=$VAL_ADDRESS_1
DEMO_ADDRESS_1=$DEMO_ADDRESS_1
RLY_ADDRESS_1=$RLY_ADDRESS_1
VAL_ADDRESS_6=$VAL_ADDRESS_6
DEMO_ADDRESS_6=$DEMO_ADDRESS_6
VAL_ADDRESS_7=$VAL_ADDRESS_7
DEMO_ADDRESS_7=$DEMO_ADDRESS_7
VAL_ADDRESS_2=$VAL_ADDRESS_2
DEMO_ADDRESS_2=$DEMO_ADDRESS_2
RLY_ADDRESS_2=$RLY_ADDRESS_2
VAL_ADDRESS_3=$VAL_ADDRESS_3
DEMO_ADDRESS_3=$DEMO_ADDRESS_3

VAL_ADDRESS_4=$VAL_ADDRESS_4
DEMO_ADDRESS_4=$DEMO_ADDRESS_4

VAL_ADDRESS_5=$VAL_ADDRESS_5
DEMO_ADDRESS_5=$DEMO_ADDRESS_5

VAL_VALOPER_2=$VAL_VALOPER_2
VAL_VALOPER_3=$VAL_VALOPER_3
VAL_VALOPER_4=$VAL_VALOPER_4
VAL_VALOPER_5=$VAL_VALOPER_5
EOF
fi
