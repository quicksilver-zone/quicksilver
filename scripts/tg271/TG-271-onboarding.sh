#!/bin/bash

THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
SCRIPT_DIR=$( realpath -- "${THIS_DIR}/.." )
export DC="-f ${THIS_DIR}/docker-compose.yml"
. ${SCRIPT_DIR}/vars.sh

export CHAINID_1=gaia-1

docker-compose $DC down

if [[ "$1" == "-r" ]]; then
  #!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

CHAINID_0=qstest-1
CHAINID_1=gaia-1

docker-compose down

echo "Removing previous data..."
rm -rf ./${CHAIN_DIR}/$CHAINID_0 &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_0}a &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_0}b &> /dev/null

rm -rf ./${CHAIN_DIR}/$CHAINID_1 &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_1}a &> /dev/null
rm -rf ./${CHAIN_DIR}/${CHAINID_1}b &> /dev/null

rm -rf ./${CHAIN_DIR}/hermes &> /dev/null
rm -rf ./${CHAIN_DIR}/icq &> /dev/null
rm -rf ./${CHAIN_DIR}/icq2 &> /dev/null

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

#relayers
if ! mkdir -p ./${CHAIN_DIR}/hermes 2>/dev/null; then
    echo "Failed to create hermes folder. Aborting..."
    exit 1
fi

cp $THIS_DIR/hermes.toml ./${CHAIN_DIR}/hermes/config.toml

if ! mkdir -p ./${CHAIN_DIR}/icq 2>/dev/null; then
    echo "Failed to create icq folder. Aborting..."
    exit 1
fi

cp ../config/icq.yaml ./${CHAIN_DIR}/icq/config.yaml

echo "Initializing $CHAINID_0..."
$QS1_RUN init test --chain-id $CHAINID_0
echo "Initializing ${CHAINID_0}a..."
$QS2_RUN init test --chain-id $CHAINID_0
echo "Initializing ${CHAINID_0}b..."
$QS3_RUN init test --chain-id $CHAINID_0

echo "Initializing $CHAINID_1..."
$GAIA1_RUN init test --chain-id $CHAINID_1
echo "Initializing ${CHAINID_1}a..."
$GAIA2_RUN init test --chain-id $CHAINID_1
echo "Initializing ${CHAINID_1}b..."
$GAIA3_RUN init test --chain-id $CHAINID_1

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $QS1_RUN keys add val1 --recover --keyring-backend=test

echo $VAL_MNEMONIC_2 | $GAIA1_RUN keys add val2 --recover --keyring-backend=test
echo $VAL_MNEMONIC_3 | $GAIA2_RUN keys add val3 --recover --keyring-backend=test
echo $VAL_MNEMONIC_4 | $GAIA3_RUN keys add val4 --recover --keyring-backend=test

echo $VAL_MNEMONIC_6 | $QS2_RUN keys add val6 --recover --keyring-backend=test
echo $VAL_MNEMONIC_7 | $QS3_RUN keys add val7 --recover --keyring-backend=test


echo $DEMO_MNEMONIC_1 | $QS1_RUN keys add demowallet1 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_2 | $QS1_RUN keys add demowallet2 --recover --keyring-backend=test

echo $DEMO_MNEMONIC_2 | $GAIA1_RUN keys add demowallet2 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_3 | $GAIA2_RUN keys add demowallet3 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_4 | $GAIA3_RUN keys add demowallet4 --recover --keyring-backend=test

echo $DEMO_MNEMONIC_6 | $QS2_RUN keys add demowallet6 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_7 | $QS3_RUN keys add demowallet7 --recover --keyring-backend=test

echo $RLY_MNEMONIC_1 | $QS1_RUN keys add rly1 --recover --keyring-backend=test
echo $RLY_MNEMONIC_2 | $GAIA1_RUN keys add rly2 --recover --keyring-backend=test

## Set denoms
${SED} 's/stake/uqck/g' $(pwd)/${CHAIN_DIR}/${CHAINID_0}/config/genesis.json
${SED} 's/stake/uqck/g' $(pwd)/${CHAIN_DIR}/${CHAINID_0}a/config/genesis.json
${SED} 's/stake/uqck/g' $(pwd)/${CHAIN_DIR}/${CHAINID_0}b/config/genesis.json

${SED} 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}/config/genesis.json
${SED} 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}a/config/genesis.json
${SED} 's/stake/uatom/g' $(pwd)/${CHAIN_DIR}/${CHAINID_1}b/config/genesis.json

VAL_ADDRESS_1=$($QS1_RUN keys show val1 --keyring-backend test -a)
DEMO_ADDRESS_1=$($QS1_RUN keys show demowallet1 --keyring-backend test -a)
RLY_ADDRESS_1=$($QS1_RUN keys show rly1 --keyring-backend test -a)

VAL_ADDRESS_2=$($GAIA1_RUN keys show val2 --keyring-backend test -a)
DEMO_ADDRESS_2=$($GAIA1_RUN keys show demowallet2 --keyring-backend test -a)
RLY_ADDRESS_2=$($GAIA1_RUN keys show rly2 --keyring-backend test -a)

VAL_ADDRESS_3=$($GAIA2_RUN keys show val3 --keyring-backend test -a)
DEMO_ADDRESS_3=$($GAIA2_RUN keys show demowallet3 --keyring-backend test -a)

VAL_ADDRESS_4=$($GAIA3_RUN keys show val4 --keyring-backend test -a)
DEMO_ADDRESS_4=$($GAIA3_RUN keys show demowallet4 --keyring-backend test -a)

VAL_VALOPER_2=$($GAIA1_RUN keys show val2 --keyring-backend test --bech=val -a)
VAL_VALOPER_3=$($GAIA2_RUN keys show val3 --keyring-backend test --bech=val -a)
VAL_VALOPER_4=$($GAIA3_RUN keys show val4 --keyring-backend test --bech=val -a)

VAL_ADDRESS_6=$($QS2_RUN keys show val6 --keyring-backend test -a)
DEMO_ADDRESS_6=$($QS2_RUN keys show demowallet6 --keyring-backend test -a)

VAL_ADDRESS_7=$($QS3_RUN keys show val7 --keyring-backend test -a)
DEMO_ADDRESS_7=$($QS3_RUN keys show demowallet7 --keyring-backend test -a)

VAL_VALOPER_6=$($QS2_RUN keys show val6 --keyring-backend test --bech=val -a)
VAL_VALOPER_7=$($QS3_RUN keys show val7 --keyring-backend test --bech=val -a)

$QS1_RUN add-genesis-account ${VAL_ADDRESS_1} 100000000000uqck
$QS1_RUN add-genesis-account ${DEMO_ADDRESS_1} 100000000000uqck
$QS1_RUN add-genesis-account ${RLY_ADDRESS_1} 100000000000uqck

$QS1_RUN add-genesis-account ${VAL_ADDRESS_6} 100000000000uqck
$QS1_RUN add-genesis-account ${VAL_ADDRESS_7} 100000000000uqck
$QS1_RUN add-genesis-account ${DEMO_ADDRESS_6} 100000000000uqck
$QS1_RUN add-genesis-account ${DEMO_ADDRESS_7} 100000000000uqck

$QS2_RUN add-genesis-account ${VAL_ADDRESS_6} 100000000000uqck
$QS3_RUN add-genesis-account ${VAL_ADDRESS_7} 100000000000uqck

$GAIA1_RUN add-genesis-account ${VAL_ADDRESS_2} 100000000000uatom
$GAIA1_RUN add-genesis-account ${VAL_ADDRESS_3} 100000000000uatom
$GAIA1_RUN add-genesis-account ${VAL_ADDRESS_4} 100000000000uatom
$GAIA1_RUN add-genesis-account ${DEMO_ADDRESS_2} 100000000000uatom
$GAIA1_RUN add-genesis-account ${DEMO_ADDRESS_3} 100000000000uatom
$GAIA1_RUN add-genesis-account ${DEMO_ADDRESS_4} 100000000000uatom
$GAIA1_RUN add-genesis-account ${RLY_ADDRESS_2} 100000000000uatom

$GAIA2_RUN add-genesis-account ${VAL_ADDRESS_3} 100000000000uatom
$GAIA3_RUN add-genesis-account ${VAL_ADDRESS_4} 100000000000uatom

echo "Creating and collecting gentx..."
$QS1_RUN gentx val1 7000000000uqck --chain-id $CHAINID_0 --keyring-backend test
$QS2_RUN gentx val6 7000000000uqck --chain-id $CHAINID_0 --keyring-backend test
$QS3_RUN gentx val7 7000000000uqck --chain-id $CHAINID_0 --keyring-backend test

$GAIA1_RUN gentx val2 7000000000uatom --commission-rate 0.33 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_1 --keyring-backend test
$GAIA2_RUN gentx val3 6000000000uatom --commission-rate 0.23 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_1 --keyring-backend test
$GAIA3_RUN gentx val4 5000000000uatom --commission-rate 0.13 --commission-max-rate 0.5 --commission-max-change-rate 0.1 --chain-id $CHAINID_1 --keyring-backend test

cp ./${CHAIN_DIR}/${CHAINID_0}a/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_0}/config/gentx/
cp ./${CHAIN_DIR}/${CHAINID_0}b/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_0}/config/gentx/

$QS1_RUN collect-gentxs

cp ./${CHAIN_DIR}/${CHAINID_1}a/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_1}/config/gentx/
cp ./${CHAIN_DIR}/${CHAINID_1}b/config/gentx/*.json ./${CHAIN_DIR}/${CHAINID_1}/config/gentx/

$GAIA1_RUN collect-gentxs

node1=$($GAIA1_RUN tendermint show-node-id)@gaia:26656
node2=$($GAIA2_RUN tendermint show-node-id)@gaia2:26656
node3=$($GAIA3_RUN tendermint show-node-id)@gaia3:26656

node5=$($QS1_RUN tendermint show-node-id)@quicksilver:26656
node6=$($QS2_RUN tendermint show-node-id)@quicksilver2:26656
node7=$($QS3_RUN tendermint show-node-id)@quicksilver3:26656

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

jq '.consensus_params.block.time_iota_ms = "200"'  ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.consensus_params.block.time_iota_ms = "200"'  ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}

## add the message types ICA should allow
jq '.app_state.interchainaccounts.host_genesis_state.params.allow_messages = ["/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend", "/cosmos.staking.v1beta1.MsgDelegate", "/cosmos.staking.v1beta1.MsgRedeemTokensforShares", "/cosmos.staking.v1beta1.MsgTokenizeShares", "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", "/ibc.applications.transfer.v1.MsgTransfer", "/cosmos.staking.v1beta1.MsgUndelegate"]' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
jq '.app.state.mint.minter.inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
jq '.app.state.mint.params.max_inflation = "2.530000000000000000"' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
jq '.app_state.staking.params.unbonding_time = "300s"' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}
cp ./${CHAIN_DIR}/${CHAINID_1}{,a}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_1}{,b}/config/genesis.json

## set the 'epoch' epoch to 5m interval
jq '.app_state.epochs.epochs = [{"identifier": "epoch","start_time": "0001-01-01T00:00:00Z","duration": "240s","current_epoch": "0","current_epoch_start_time": "0001-01-01T00:00:00Z","epoch_counting_started": false,"current_epoch_start_height": "0"}]' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.interchainstaking.params.deposit_interval = 25' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.mint.params.epoch_identifier = "epoch"' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.gov.deposit_params.min_deposit = [{"denom": "uqck", "amount": "100"}]' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.gov.deposit_params.max_deposit_period = "10s"' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}
jq '.app_state.gov.voting_params.voting_period = "10s"' ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_0}/config/genesis.json{.new,}

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

docker-compose down

cat << EOF > ${THIS_DIR}/wallets.sh
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

VAL_VALOPER_2=$VAL_VALOPER_2
VAL_VALOPER_3=$VAL_VALOPER_3
VAL_VALOPER_4=$VAL_VALOPER_4
EOF

else
  echo "Copying previously generated state."
  rm -rf ${CHAIN_DIR}/${CHAINID_0}
  rm -rf ${CHAIN_DIR}/${CHAINID_0}a
  rm -rf ${CHAIN_DIR}/${CHAINID_0}b
  rm -rf ${CHAIN_DIR}/${CHAINID_1}
  rm -rf ${CHAIN_DIR}/${CHAINID_1}a
  rm -rf ${CHAIN_DIR}/${CHAINID_1}b
  rm -rf ${CHAIN_DIR}/hermes &> /dev/null
  rm -rf ${CHAIN_DIR}/icq &> /dev/null
  rm -rf ${CHAIN_DIR}/rly &> /dev/null
  
  TIME=${TIME}
  jq ".genesis_time = \"$TIME\"" ./${CHAIN_DIR}/backup/${CHAINID_0}/config/genesis.json > ./${CHAIN_DIR}/backup/${CHAINID_0}/config/genesis.json.new && mv ./${CHAIN_DIR}/backup/${CHAINID_0}/config/genesis.json{.new,}
  jq ".genesis_time = \"$TIME\"" ./${CHAIN_DIR}/backup/${CHAINID_0}a/config/genesis.json > ./${CHAIN_DIR}/backup/${CHAINID_0}a/config/genesis.json.new && mv ./${CHAIN_DIR}/backup/${CHAINID_0}a/config/genesis.json{.new,}
  jq ".genesis_time = \"$TIME\"" ./${CHAIN_DIR}/backup/${CHAINID_0}b/config/genesis.json > ./${CHAIN_DIR}/backup/${CHAINID_0}b/config/genesis.json.new && mv ./${CHAIN_DIR}/backup/${CHAINID_0}b/config/genesis.json{.new,}

  cp -fr ${CHAIN_DIR}/backup/${CHAINID_0} ${CHAIN_DIR}/${CHAINID_0}
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_0}a ${CHAIN_DIR}/${CHAINID_0}a
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_0}b ${CHAIN_DIR}/${CHAINID_0}b
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_1} ${CHAIN_DIR}/${CHAINID_1}
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_1}a ${CHAIN_DIR}/${CHAINID_1}a
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_1}b ${CHAIN_DIR}/${CHAINID_1}b
  mkdir ${CHAIN_DIR}/hermes ${CHAIN_DIR}/icq
  cp ../config/icq.yaml ./${CHAIN_DIR}/icq/config.yaml
  cp ./hermes.toml ./${CHAIN_DIR}/hermes/config.toml
  cp -rf ../config/rly ./${CHAIN_DIR}/rly
fi

source ${THIS_DIR}/wallets.sh

#############################################################################################################################

docker-compose $DC up --force-recreate -d quicksilver quicksilver2 quicksilver3 gaia gaia2 gaia3
echo "Chains created"
sleep 3
echo "Restoring keys"
echo "$RLY_MNEMONIC_1" | $HERMES_RUN keys add --mnemonic-file /dev/fd/0 --chain $CHAINID_0
echo "$RLY_MNEMONIC_2" | $HERMES_RUN keys add --mnemonic-file /dev/fd/0 --chain $CHAINID_1
sleep 3
#echo "Creating IBC connection"
echo "Creating connection & transfer channel"
$HERMES_RUN create channel --a-chain $CHAINID_0 --b-chain $CHAINID_1 --a-port transfer --b-port transfer --new-client-connection --yes
#$HERMES_RUN create connection --a-chain $CHAINID_0 --b-chain $CHAINID_1
#$HERMES_RUN create channel --port-a transfer --port-b transfer $CHAINID_0 connection-0
echo "Tranfer channel created"
docker-compose $DC up --force-recreate -d hermes
RLY_ADDRESS_3=$($RLY_RUN keys show qstest-1 testkey)
RLY_ADDRESS_4=$($RLY_RUN keys show lstest-1 testkey)

## TG-271 - send to delegate account before we register it!
$GAIA1_EXEC tx bank send val2 cosmos1j5u9y5gm95f4sudpupu8zv7jqlnh8wzlnn88w7upx5exj4ekr0fsmx3jup 1uatom --chain-id $CHAINID_1 -y --keyring-backend=test -b block
sleep 3
$QS1_EXEC tx bank send val1 $RLY_ADDRESS_3 1000uqck --chain-id $CHAINID_0 -y --keyring-backend=test
$GAIA1_EXEC tx bank send val2 $RLY_ADDRESS_4 1000uatom --chain-id $CHAINID_1 -y --keyring-backend=test

docker-compose $DC up --force-recreate -d relayer

rm -rf ./icq/keys
echo "Launch and configure interchain query daemon"

ICQ_ADDRESS_1=$($ICQ_RUN keys add test --chain qstest-1 | jq .address -r)
ICQ_ADDRESS_2=$($ICQ_RUN keys add test --chain lstest-1 | jq .address -r)

sleep 3

$QS1_EXEC tx bank send val1 $ICQ_ADDRESS_1 1000uqck --chain-id $CHAINID_0 -y --keyring-backend=test
$GAIA1_EXEC tx bank send val2 $ICQ_ADDRESS_2 1000uatom --chain-id $CHAINID_1 -y --keyring-backend=test

docker-compose $DC up --force-recreate -d icq

#echo "Register $CHAINID_1 on quicksilver..."
cat $THIS_DIR/../registerzone.json | jq . -c | $QS1_EXEC tx gov submit-proposal /dev/fd/0 --from demowallet1 --chain-id $CHAINID_0 --gas 2000000 -y --keyring-backend=test
sleep 3
$QS1_EXEC tx gov vote 1 yes --from val1 --chain-id $CHAINID_0 -y --keyring-backend=test
$QS2_EXEC tx gov vote 1 yes --from val6 --chain-id $CHAINID_0 -y --keyring-backend=test
$QS3_EXEC tx gov vote 1 yes --from val7 --chain-id $CHAINID_0 -y --keyring-backend=test
sleep 10
docker-compose $DC restart hermes
sleep 5

sleep 3
DEPOSIT_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
while [[ "$DEPOSIT_ACCOUNT" == "null" ]]; do
  sleep 2
  DEPOSIT_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
done

PERFORMANCE_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].performance_address.address -r)
while [[ "$PERFORMANCE_ACCOUNT" == "null" ]]; do
  sleep 2
  PERFORMANCE_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].performance_address.address -r)
done

$GAIA1_EXEC tx bank send val2 $PERFORMANCE_ACCOUNT 40000uatom --chain-id $CHAINID_1 -y --keyring-backend=test

sleep 3

$GAIA1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 333333uatom --chain-id $CHAINID_1 -y --keyring-backend=test
sleep 5
$GAIA1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 20000000uatom --chain-id $CHAINID_1 -y --keyring-backend=test --note MgTUzEjWVVYoDZBarqFL1akb38mxlgTsqdZ/sFxTJBNf+tv6rtckvn3T
sleep 5
$GAIA1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 33000000uatom --chain-id $CHAINID_1 -y --keyring-backend=test
sleep 5


