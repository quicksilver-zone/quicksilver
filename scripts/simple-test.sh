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
QS_EXEC="docker-compose --ansi never exec -T quicksilver quicksilverd"
TZ_RUN="docker-compose --ansi never run -T testzone icad"
ICQ_RUN="docker-compose --ansi never run -T icq interchain-queries"

VAL_MNEMONIC_1="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
VAL_MNEMONIC_2="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice"
DEMO_MNEMONIC_1="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"
DEMO_MNEMONIC_2="veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
RLY_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
RLY_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

docker-compose down

echo "Removing previous data..."
sudo rm -rf ./$CHAIN_DIR/$CHAINID_1 &> /dev/null
sudo rm -rf ./$CHAIN_DIR/$CHAINID_2 &> /dev/null

# Add directories for both chains, exit if an error occurs
if ! mkdir -p ./$CHAIN_DIR/$CHAINID_1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi
sudo mkdir -p ./$CHAIN_DIR/$CHAINID_1/data/snapshots/metadata.db

if ! mkdir -p ./$CHAIN_DIR/$CHAINID_2 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAINID_1..."
echo "Initializing $CHAINID_2..."

$QS_RUN init test --chain-id $CHAINID_1
$TZ_RUN init test --chain-id $CHAINID_2

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $QS_RUN keys add val1 --recover --keyring-backend=test
echo $VAL_MNEMONIC_2 | $TZ_RUN keys add val2 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_1 | $QS_RUN keys add demowallet1 --recover --keyring-backend=test
echo $DEMO_MNEMONIC_2 | $TZ_RUN keys add demowallet2 --recover --keyring-backend=test
echo $RLY_MNEMONIC_1 | $QS_RUN keys add rly1 --recover --keyring-backend=test
echo $RLY_MNEMONIC_2 | $TZ_RUN keys add rly2 --recover --keyring-backend=test

## Set denoms
sudo sed -i 's/stake/uqck/g' $(pwd)/$CHAIN_DIR/$CHAINID_1/config/genesis.json

VAL_ADDRESS_1=$($QS_RUN keys show val1 --keyring-backend test -a)
DEMO_ADDRESS_1=$($QS_RUN keys show demowallet1 --keyring-backend test -a)
RLY_ADDRESS_1=$($QS_RUN keys show rly1 --keyring-backend test -a)

VAL_ADDRESS_2=$($TZ_RUN keys show val2 --keyring-backend test -a)
DEMO_ADDRESS_2=$($TZ_RUN keys show demowallet2 --keyring-backend test -a)
RLY_ADDRESS_2=$($TZ_RUN keys show rly2 --keyring-backend test -a)

VAL_VALOPER_2=$($TZ_RUN keys show val2 --keyring-backend test --bech=val -a)

$QS_RUN add-genesis-account ${VAL_ADDRESS_1} 100000000000uqck
$QS_RUN add-genesis-account ${DEMO_ADDRESS_1} 100000000000uqck
$QS_RUN add-genesis-account ${RLY_ADDRESS_1} 100000000000uqck

$TZ_RUN add-genesis-account ${VAL_ADDRESS_2} 100000000000stake
$TZ_RUN add-genesis-account ${DEMO_ADDRESS_2} 100000000000stake
$TZ_RUN add-genesis-account ${RLY_ADDRESS_2} 100000000000stake

echo "Creating and collecting gentx..."
$QS_RUN gentx val1 7000000000uqck --chain-id $CHAINID_1 --keyring-backend test
$TZ_RUN gentx val2 7000000000stake --chain-id $CHAINID_2 --keyring-backend test
$QS_RUN collect-gentxs
$TZ_RUN collect-gentxs

echo "Changing defaults and ports in app.toml and config.toml files..."
sudo sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sudo sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sudo sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sudo sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sudo sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sudo sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml

sudo sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sudo sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sudo sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sudo sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sudo sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sudo sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_2/config/app.toml

## add liquidstaking module messages here too.
sudo sed -i -e 's/\"allow_messages\":.*/\"allow_messages\": [\"\/cosmos.bank.v1beta1.MsgSend\", \"\/cosmos.staking.v1beta1.MsgDelegate\"]/g' $CHAIN_DIR/$CHAINID_2/config/genesis.json

docker-compose up --force-recreate -d quicksilver testzone
echo "Chains created"
sleep 5
echo "Restoring keys"
docker-compose run hermes hermes -c /tmp/hermes.toml keys restore --mnemonic "$RLY_MNEMONIC_1" test-1
docker-compose run hermes hermes -c /tmp/hermes.toml keys restore --mnemonic "$RLY_MNEMONIC_2" test-2
sleep 5
echo "Creating transfer channel"
docker-compose run hermes hermes -c /tmp/hermes.toml create channel --port-a transfer --port-b transfer $CHAINID_1 $CHAINID_2
echo "Tranfer channel created"
docker-compose up --force-recreate -d hermes

sudo rm -rf ./icq/keys
echo "Launch and configure interchain query daemon"

ICQ_ADDRESS_1=$($ICQ_RUN keys add test --chain quicksilver | jq .address -r)
ICQ_ADDRESS_2=$($ICQ_RUN keys add test --chain liquidstaking | jq .address -r)

docker-compose exec quicksilver quicksilverd tx bank send val1 $ICQ_ADDRESS_1 1000uqck --chain-id $CHAINID_1 -y --keyring-backend=test
docker-compose exec testzone icad tx bank send val2 $ICQ_ADDRESS_2 1000stake --chain-id $CHAINID_2 -y --keyring-backend=test

docker-compose up --force-recreate -d icq

echo "Register $CHAINID_2 on quicksilver..."
docker-compose exec quicksilver quicksilverd tx interchainstaking register cosmos connection-0 $CHAINID_2 stake --chain-id $CHAINID_1 --from demowallet1 --gas 10000000 --keyring-backend test -y

sleep 30

## TODO: get val2 valoper from keys
docker-compose exec testzone icad tx staking tokenize-share $VAL_VALOPER_2 10000stake $VAL_ADDRESS_2 --keyring-backend test --from val2 --chain-id $CHAINID_2 --gas 400000

# we need to query the deposit account to automate the next step!
DEPOSIT_ACCOUNT=$($QS_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address -r)
#cosmos1rk2z2mlfhtu6l6gulwkzn7p03edr2hrzt2v7xshpppzc0rznk56stzqtkj
docker-compose exec testzone icad tx bank send val2 $DEPOSIT_ACCOUNT 10000${VAL_VALOPER_2}1  --keyring-backend test --chain-id $CHAINID_2
