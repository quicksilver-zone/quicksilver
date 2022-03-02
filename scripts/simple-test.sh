#!/bin/bash
set -x

QS_IMAGE=quicksilverzone/quicksilver
QS_VERSION=latest
TZ_IMAGE=quicksilverzone/testzone
TZ_VERSION=latest

CHAIN_DIR=data
CHAINID_1=test-1
CHAINID_2=test-2

QS_RUN="docker-compose --ansi never run -T quicksilver quicksilverd"
TZ_RUN="docker-compose --ansi never run -T testzone icad"
TZ2_RUN="docker-compose --ansi never run -T testzone2 icad"
TZ3_RUN="docker-compose --ansi never run -T testzone3 icad"
TZ4_RUN="docker-compose --ansi never run -T testzone4 icad"

QS_EXEC="docker-compose --ansi never exec -T quicksilver quicksilverd"
TZ_EXEC="docker-compose --ansi never exec -T testzone icad"
TZ2_EXEC="docker-compose --ansi never exec -T testzone2 icad"
TZ3_EXEC="docker-compose --ansi never exec -T testzone3 icad"
TZ4_EXEC="docker-compose --ansi never exec -T testzone4 icad"

ICQ_RUN="docker-compose --ansi never run -T icq interchain-queries"

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
VAL_MNEMONIC_2="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice"
VAL_MNEMONIC_3="convince erupt tongue pet jeans leader boil mosquito unfair move dinosaur wrist ankle clog brown nerve next lunch speak source turtle fault gun fade"
VAL_MNEMONIC_4="cheese alarm easy kick now tattoo forward blast exercise abuse brisk race embrace cook august dwarf axis flat allow cup ripple measure keep flip"
VAL_MNEMONIC_5="ecology thank spot fork trust sorry speed april hood midnight put umbrella detail coin census crash ride fan know cup liar plastic kitten affair"
DEMO_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
DEMO_MNEMONIC_2="veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
DEMO_MNEMONIC_3="snow cancel exhibit neutral cushion what bench bomb season hard mesh method virus enforce hip put voice toilet love head risk ankle toy fiscal"
DEMO_MNEMONIC_4="sustain stumble true ozone note engine unit dignity tip sheriff barrel connect fire ridge wealth echo behind will pledge coin joke mouse ripple battle"
DEMO_MNEMONIC_5="remain season shoot frog include erase august click rookie shine person oxygen pyramid table disagree language blossom island begin theory strike planet acid mad"
RLY_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
RLY_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

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
cp ./${CHAIN_DIR}/${CHAINID_2}{,a}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_2}{,b}/config/genesis.json
cp ./${CHAIN_DIR}/${CHAINID_2}{,c}/config/genesis.json

## set the 'epoch' epoch to 5m interval
jq '.app_state.epochs.epochs = [{"identifier": "epoch","start_time": "0001-01-01T00:00:00Z","duration": "300s","current_epoch": "0","current_epoch_start_time": "0001-01-01T00:00:00Z","epoch_counting_started": false,"current_epoch_start_height": "0"}]' ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/${CHAINID_1}/config/genesis.json{.new,}

docker-compose up --force-recreate -d quicksilver testzone testzone2 testzone3 testzone4
echo "Chains created"
sleep 2
echo "Restoring keys"
docker-compose run hermes hermes -c /tmp/hermes.toml keys restore --mnemonic "$RLY_MNEMONIC_1" test-1
docker-compose run hermes hermes -c /tmp/hermes.toml keys restore --mnemonic "$RLY_MNEMONIC_2" test-2
sleep 5
echo "Creating transfer channel"
docker-compose run hermes hermes -c /tmp/hermes.toml create channel --port-a transfer --port-b transfer $CHAINID_1 $CHAINID_2
echo "Tranfer channel created"
docker-compose up --force-recreate -d hermes

rm -rf ./icq/keys
echo "Launch and configure interchain query daemon"

ICQ_ADDRESS_1=$($ICQ_RUN keys add test --chain quicksilver | jq .address -r)
ICQ_ADDRESS_2=$($ICQ_RUN keys add test --chain liquidstaking | jq .address -r)

$QS_EXEC tx bank send val1 $ICQ_ADDRESS_1 1000uqck --chain-id $CHAINID_1 -y --keyring-backend=test
$TZ_EXEC tx bank send val2 $ICQ_ADDRESS_2 1000uatom --chain-id $CHAINID_2 -y --keyring-backend=test

docker-compose up --force-recreate -d icq

echo "Register $CHAINID_2 on quicksilver..."
$QS_EXEC tx interchainstaking register cosmos connection-0 $CHAINID_2 uatom --from demowallet1 --gas 10000000 --chain-id $CHAINID_1 -y --keyring-backend=test

sleep 15

## TODO: get val2 valoper from keys
$TZ_EXEC tx staking tokenize-share $VAL_VALOPER_2 10000uatom $VAL_ADDRESS_2 --from val2 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test  #1
$TZ2_EXEC tx staking tokenize-share $VAL_VALOPER_3 25000uatom $VAL_ADDRESS_3 --from val3 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test   #2
$TZ3_EXEC tx staking tokenize-share $VAL_VALOPER_4 65000uatom $VAL_ADDRESS_4 --from val4 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test  #3

sleep 5
DEPOSIT_ACCOUNT=$($QS_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
while [[ "$DEPOSIT_ACCOUNT" == "null" ]]; do
  sleep 5
  DEPOSIT_ACCOUNT=$($QS_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
done

sleep 5
$TZ_EXEC tx bank send val2 $DEPOSIT_ACCOUNT 10000${VAL_VALOPER_2}1 --chain-id $CHAINID_2 -y --keyring-backend=test
$TZ2_EXEC tx bank send val3 $DEPOSIT_ACCOUNT 15000${VAL_VALOPER_3}2 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 5
$TZ_EXEC tx staking delegate ${VAL_VALOPER_2} 36000uatom --from demowallet2 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 5

$TZ_EXEC tx staking tokenize-share $VAL_VALOPER_2 36000uatom $VAL_ADDRESS_2 --from demowallet2 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test   #4
$TZ2_EXEC tx bank send val3 $DEPOSIT_ACCOUNT 10000${VAL_VALOPER_3}2 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 10

$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 20000${VAL_VALOPER_2}4 --chain-id $CHAINID_2 -y --keyring-backend=test
$TZ3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 25000${VAL_VALOPER_4}3 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 10

$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 10000${VAL_VALOPER_2}4 --chain-id $CHAINID_2 -y --keyring-backend=test
$TZ3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 15000${VAL_VALOPER_4}3 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 10

$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 6000${VAL_VALOPER_2}4 --chain-id $CHAINID_2 -y --keyring-backend=test
$TZ3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 25000${VAL_VALOPER_4}3 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 30

V2_DELEG=$(docker-compose exec testzone icad q staking delegations-to ${VAL_VALOPER_2} --output=json | jq '.delegation_responses[].balance.amount' -r | sort -n | head -n-1 | awk '{sum+=$0} END{print sum}')
V3_DELEG=$(docker-compose exec testzone icad q staking delegations-to ${VAL_VALOPER_3} --output=json | jq '.delegation_responses[].balance.amount' -r | sort -n | head -n-1 | awk '{sum+=$0} END{print sum}')
V4_DELEG=$(docker-compose exec testzone icad q staking delegations-to ${VAL_VALOPER_4} --output=json | jq '.delegation_responses[].balance.amount' -r | sort -n | head -n-1 | awk '{sum+=$0} END{print sum}')

if [[ ! $V2_DELEG -eq 46000 ]]; then echo "ERROR: val 2 delegation does not match 46000 ($V2_DELEG)"; exit 1; fi
if [[ ! $V3_DELEG -eq 25000 ]]; then echo "ERROR: val 3 delegation does not match 25000 ($V3_DELEG)"; exit 1; fi
if [[ ! $V4_DELEG -eq 65000 ]]; then echo "ERROR: val 4 delegation does not match 65000 ($V4_DELEG)"; exit 1; fi

echo "All tests passed :)"
