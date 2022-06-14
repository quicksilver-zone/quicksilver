#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/vars.sh

docker-compose down

if [[ "$1" == "-r" ]]; then
  echo "Regenerating state."
  $SCRIPT_DIR/setup.sh
else
  echo "Copying previously generated state."
  rm -rf ${CHAIN_DIR}/${CHAINID_1}
  rm -rf ${CHAIN_DIR}/${CHAINID_1}a
  rm -rf ${CHAIN_DIR}/${CHAINID_1}b
  rm -rf ${CHAIN_DIR}/${CHAINID_2}
  rm -rf ${CHAIN_DIR}/${CHAINID_2}a
  rm -rf ${CHAIN_DIR}/${CHAINID_2}b
  rm -rf ${CHAIN_DIR}/${CHAINID_2}c
  rm -rf ${CHAIN_DIR}/hermes &> /dev/null
  rm -rf ${CHAIN_DIR}/icq &> /dev/null

  TIME=$(date --date '-2 minutes' +%Y-%m-%dT%H:%M:00Z -u)
  jq ".genesis_time = \"$TIME\"" ./${CHAIN_DIR}/backup/${CHAINID_1}/config/genesis.json > ./${CHAIN_DIR}/backup/${CHAINID_1}/config/genesis.json.new && mv ./${CHAIN_DIR}/backup/${CHAINID_1}/config/genesis.json{.new,}
  jq ".genesis_time = \"$TIME\"" ./${CHAIN_DIR}/backup/${CHAINID_1}a/config/genesis.json > ./${CHAIN_DIR}/backup/${CHAINID_1}a/config/genesis.json.new && mv ./${CHAIN_DIR}/backup/${CHAINID_1}a/config/genesis.json{.new,}
  jq ".genesis_time = \"$TIME\"" ./${CHAIN_DIR}/backup/${CHAINID_1}b/config/genesis.json > ./${CHAIN_DIR}/backup/${CHAINID_1}b/config/genesis.json.new && mv ./${CHAIN_DIR}/backup/${CHAINID_1}b/config/genesis.json{.new,}

  cp -fr ${CHAIN_DIR}/backup/${CHAINID_1} ${CHAIN_DIR}/${CHAINID_1}
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_1}a ${CHAIN_DIR}/${CHAINID_1}a
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_1}b ${CHAIN_DIR}/${CHAINID_1}b
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2} ${CHAIN_DIR}/${CHAINID_2}
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2}a ${CHAIN_DIR}/${CHAINID_2}a
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2}b ${CHAIN_DIR}/${CHAINID_2}b
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2}c ${CHAIN_DIR}/${CHAINID_2}c
  mkdir ${CHAIN_DIR}/hermes ${CHAIN_DIR}/icq
  cp ./scripts/config/icq.yaml ./${CHAIN_DIR}/icq/config.yaml
fi

source ${SCRIPT_DIR}/wallets.sh

#############################################################################################################################

docker-compose up --force-recreate -d quicksilver quicksilver2 quicksilver3 testzone testzone2 testzone3 testzone4
echo "Chains created"
sleep 10
echo "Restoring keys"
docker-compose run hermes hermes -c /tmp/hermes.toml keys restore --mnemonic "$RLY_MNEMONIC_1" test-1
docker-compose run hermes hermes -c /tmp/hermes.toml keys restore --mnemonic "$RLY_MNEMONIC_2" test-2
sleep 10
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
$QS_EXEC tx interchainstaking register connection-0 uqatom uatom cosmos --from demowallet1 --gas 10000000 --chain-id $CHAINID_1 -y --keyring-backend=test --multi-send --lsm-support

sleep 5

## TODO: get val2 valoper from keys
$TZ_EXEC tx staking tokenize-share $VAL_VALOPER_2 10000000uatom $VAL_ADDRESS_2 --from val2 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test  #1
$TZ2_EXEC tx staking tokenize-share $VAL_VALOPER_3 25000000uatom $VAL_ADDRESS_3 --from val3 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test   #2
$TZ3_EXEC tx staking tokenize-share $VAL_VALOPER_4 65000000uatom $VAL_ADDRESS_4 --from val4 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test  #3

sleep 5
DEPOSIT_ACCOUNT=$($QS_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
while [[ "$DEPOSIT_ACCOUNT" == "null" ]]; do
  sleep 5
  DEPOSIT_ACCOUNT=$($QS_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
done

PERFORMANCE_ACCOUNT=$($QS_EXEC q interchainstaking zones --output=json | jq .zones[0].performance_address.address -r)
while [[ "$PERFORMANCE_ACCOUNT" == "null" ]]; do
  sleep 2
  PERFORMANCE_ACCOUNT=$($QS_EXEC q interchainstaking zones --output=json | jq .zones[0].performance_address.address -r)
done

$TZ_EXEC tx bank send val2 $PERFORMANCE_ACCOUNT 40000uatom --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 5
$TZ_EXEC tx bank send val2 $DEPOSIT_ACCOUNT 10000000${VAL_VALOPER_2}1 --chain-id $CHAINID_2 -y --keyring-backend=test
sleep 10
$TZ2_EXEC tx bank send val3 $DEPOSIT_ACCOUNT 15000000${VAL_VALOPER_3}2 --chain-id $CHAINID_2 -y --keyring-backend=test
sleep 10
$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 333333uatom --chain-id $CHAINID_2 -y --keyring-backend=test
sleep 20
$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 20000000uatom --chain-id $CHAINID_2 -y --keyring-backend=test --note MgTUzEjWVVYoDZBarqFL1akb38mxlgTsqdZ/sFxTJBNf+tv6rtckvn3T
sleep 10
$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 33000000uatom --chain-id $CHAINID_2 -y --keyring-backend=test
sleep 10
$TZ_EXEC tx staking tokenize-share $VAL_VALOPER_2 36000000uatom $VAL_ADDRESS_2 --from demowallet2 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test   #4
$TZ2_EXEC tx bank send val3 $DEPOSIT_ACCOUNT 10000000${VAL_VALOPER_3}2 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 10

$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 20000000${VAL_VALOPER_2}4 --chain-id $CHAINID_2 -y --keyring-backend=test
$TZ3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 25000000${VAL_VALOPER_4}3 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 10

$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 10000000${VAL_VALOPER_2}4 --chain-id $CHAINID_2 -y --keyring-backend=test
$TZ3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 15000000${VAL_VALOPER_4}3 --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 10

$TZ_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 6000000${VAL_VALOPER_2}4 --chain-id $CHAINID_2 -y --keyring-backend=test
$TZ3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 25000000${VAL_VALOPER_4}3 --chain-id $CHAINID_2 -y --keyring-backend=test

