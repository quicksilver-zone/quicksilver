#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/vars.sh

docker-compose down

IS_MULTI_ZONE_TEST=true
export IS_MULTI_ZONE_TEST

if [[ "$1" == "-r" ]]; then
  echo "Regenerating state."
  $SCRIPT_DIR/setup.sh
else
  echo "Copying previously generated state."
  rm -rf ${CHAIN_DIR}/${CHAINID_0}
  rm -rf ${CHAIN_DIR}/${CHAINID_0}a
  rm -rf ${CHAIN_DIR}/${CHAINID_0}b
  rm -rf ${CHAIN_DIR}/${CHAINID_1}
  rm -rf ${CHAIN_DIR}/${CHAINID_1}a
  rm -rf ${CHAIN_DIR}/${CHAINID_1}b
  rm -rf ${CHAIN_DIR}/${CHAINID_1}c
  rm -rf ${CHAIN_DIR}/${CHAINID_2}
  rm -rf ${CHAIN_DIR}/${CHAINID_2}a
  rm -rf ${CHAIN_DIR}/${CHAINID_2}b
  rm -rf ${CHAIN_DIR}/${CHAINID_2}c
  rm -rf ${CHAIN_DIR}/hermes &> /dev/null
  rm -rf ${CHAIN_DIR}/icq &> /dev/null
  rm -rf ${CHAIN_DIR}/icq2 &> /dev/null

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
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_1}c ${CHAIN_DIR}/${CHAINID_1}c
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2} ${CHAIN_DIR}/${CHAINID_2}
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2}a ${CHAIN_DIR}/${CHAINID_2}a
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2}b ${CHAIN_DIR}/${CHAINID_2}b
  cp -fr ${CHAIN_DIR}/backup/${CHAINID_2}c ${CHAIN_DIR}/${CHAINID_2}c
  mkdir ${CHAIN_DIR}/hermes ${CHAIN_DIR}/icq
  cp ./scripts/config/icq.yaml ./${CHAIN_DIR}/icq/config.yaml
  cp ./scripts/config/icq2.yaml ./${CHAIN_DIR}/icq2/config.yaml
  cp ./scripts/config/hermes-2.toml ./${CHAIN_DIR}/hermes/config.toml
fi

source ${SCRIPT_DIR}/wallets.sh

#############################################################################################################################

docker-compose up --force-recreate -d quicksilver quicksilver2 quicksilver3 testzone1-1 testzone1-2 testzone1-3 testzone1-4 testzone2-1 testzone2-2 testzone2-3 testzone2-4
echo "Chains created"
sleep 10
echo "Restoring keys"
echo "$RLY_MNEMONIC_1" | $HERMES_RUN keys add --mnemonic-file /dev/fd/0 --chain $CHAINID_0
echo "$RLY_MNEMONIC_2" | $HERMES_RUN keys add --mnemonic-file /dev/fd/0 --chain $CHAINID_1
echo "$RLY_MNEMONIC_3" | $HERMES_RUN keys add --mnemonic-file /dev/fd/0 --chain $CHAINID_2
sleep 10
echo "Creating IBC connection"
$HERMES_RUN create channel --a-chain $CHAINID_0 --b-chain $CHAINID_1 --a-port transfer --b-port transfer --new-client-connection --yes
$HERMES_RUN create channel --a-chain $CHAINID_0 --b-chain $CHAINID_2 --a-port transfer --b-port transfer --new-client-connection --yes

echo "Tranfer channel created"
docker-compose up --force-recreate -d hermes

rm -rf ./icq/keys
rm -rf ./icq2/keys
echo "Launch and configure interchain query daemon"

ICQ_ADDRESS_1=$($ICQ_RUN keys add test --chain qstest-1 | jq .address -r)
ICQ_ADDRESS_2=$($ICQ_RUN keys add test --chain lstest-1 | jq .address -r)

ICQ_ADDRESS_3=$($ICQ2_RUN keys add test --chain qstest-1 | jq .address -r)
ICQ_ADDRESS_4=$($ICQ2_RUN keys add test --chain lstest-2 | jq .address -r)

$QS1_EXEC tx bank send val1 $ICQ_ADDRESS_1 1000uqck --chain-id $CHAINID_0 -y --keyring-backend=test
$TZ1_1_EXEC tx bank send val2 $ICQ_ADDRESS_2 1000uatom --chain-id $CHAINID_1 -y --keyring-backend=test

sleep 2
$QS1_EXEC tx bank send val1 $ICQ_ADDRESS_3 1000uqck --chain-id $CHAINID_0 -y --keyring-backend=test
$TZ2_1_EXEC tx bank send val8 $ICQ_ADDRESS_4 1000uosmo --chain-id $CHAINID_2 -y --keyring-backend=test

docker-compose up --force-recreate -d icq icq2

#echo "Register $CHAINID_1 on quicksilver..."
cat $SCRIPT_DIR/registerzone.json | jq . -c | $QS1_EXEC tx gov submit-proposal /dev/fd/0 --from demowallet1 --chain-id $CHAINID_0 --gas 2000000 -y --keyring-backend=test
sleep 5
$QS1_EXEC tx gov vote 1 yes --from val1 --chain-id $CHAINID_0 -y --keyring-backend=test
$QS2_EXEC tx gov vote 1 yes --from val6 --chain-id $CHAINID_0 -y --keyring-backend=test
$QS3_EXEC tx gov vote 1 yes --from val7 --chain-id $CHAINID_0 -y --keyring-backend=test

#echo "Register $CHAINID_2 on quicksilver..."
cat $SCRIPT_DIR/registerosmo.json | jq . -c | $QS1_EXEC tx gov submit-legacy-proposal register-zone  /dev/fd/0 --from demowallet1 --chain-id $CHAINID_0 --gas 2000000 -y --keyring-backend=test
sleep 5
$QS1_EXEC tx gov vote 2 yes --from val1 --chain-id $CHAINID_0 -y --keyring-backend=test
$QS2_EXEC tx gov vote 2 yes --from val6 --chain-id $CHAINID_0 -y --keyring-backend=test
$QS3_EXEC tx gov vote 2 yes --from val7 --chain-id $CHAINID_0 -y --keyring-backend=test
sleep 15
docker-compose restart hermes

## TODO: get val2 valoper from keys
# $TZ1_1_EXEC tx staking tokenize-share $VAL_VALOPER_2 10000000uatom $VAL_ADDRESS_2 --from val2 --gas 400000 --chain-id $CHAINID_1 -y --keyring-backend=test  #1
# $TZ1_2_EXEC tx staking tokenize-share $VAL_VALOPER_3 25000000uatom $VAL_ADDRESS_3 --from val3 --gas 400000 --chain-id $CHAINID_1 -y --keyring-backend=test  #2
# $TZ1_3_EXEC tx staking tokenize-share $VAL_VALOPER_4 65000000uatom $VAL_ADDRESS_4 --from val4 --gas 400000 --chain-id $CHAINID_1 -y --keyring-backend=test  #3

sleep 5
DEPOSIT_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
while [[ "$DEPOSIT_ACCOUNT" == "null" ]]; do
  sleep 5
  DEPOSIT_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].deposit_address.address -r)
done

PERFORMANCE_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].performance_address.address -r)
while [[ "$PERFORMANCE_ACCOUNT" == "null" ]]; do
  sleep 2
  PERFORMANCE_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[0].performance_address.address -r)
done
sleep 2
$TZ1_1_EXEC tx bank send val2 $PERFORMANCE_ACCOUNT 40000uatom --chain-id $CHAINID_1 -y --keyring-backend=test

sleep 3
# $TZ1_1_EXEC tx bank send val2 $DEPOSIT_ACCOUNT 10000000${VAL_VALOPER_2}1 --chain-id $CHAINID_1 -y --keyring-backend=test
# sleep 10
# $TZ1_2_EXEC tx bank send val3 $DEPOSIT_ACCOUNT 15000000${VAL_VALOPER_3}2 --chain-id $CHAINID_1 -y --keyring-backend=test
# sleep 10
$TZ1_1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 333333uatom --chain-id $CHAINID_1 -y --keyring-backend=test
sleep 5
$TZ1_1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 20000000uatom --chain-id $CHAINID_1 -y --keyring-backend=test --note MgTUzEjWVVYoDZBarqFL1akb38mxlgTsqdZ/sFxTJBNf+tv6rtckvn3T
sleep 5
$TZ1_1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 33000000uatom --chain-id $CHAINID_1 -y --keyring-backend=test
# sleep 10
# $TZ1_1_EXEC tx staking tokenize-share $VAL_VALOPER_2 36000000uatom $VAL_ADDRESS_2 --from demowallet2 --gas 400000 --chain-id $CHAINID_1 -y --keyring-backend=test   #4
# $TZ1_2_EXEC tx bank send val3 $DEPOSIT_ACCOUNT 10000000${VAL_VALOPER_3}2 --chain-id $CHAINID_1 -y --keyring-backend=test

# sleep 10

# $TZ1_1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 20000000${VAL_VALOPER_2}4 --chain-id $CHAINID_1 -y --keyring-backend=test
# $TZ1_3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 25000000${VAL_VALOPER_4}3 --chain-id $CHAINID_1 -y --keyring-backend=test

# sleep 10

# $TZ1_1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 10000000${VAL_VALOPER_2}4 --chain-id $CHAINID_1 -y --keyring-backend=test
# $TZ1_3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 15000000${VAL_VALOPER_4}3 --chain-id $CHAINID_1 -y --keyring-backend=test

# sleep 10

# $TZ1_1_EXEC tx bank send demowallet2 $DEPOSIT_ACCOUNT 6000000${VAL_VALOPER_2}4 --chain-id $CHAINID_1 -y --keyring-backend=test
# $TZ1_3_EXEC tx bank send val4 $DEPOSIT_ACCOUNT 25000000${VAL_VALOPER_4}3 --chain-id $CHAINID_1 -y --keyring-backend=test

#$TZ2_1_EXEC tx staking tokenize-share $VAL_VALOPER_8 10000000uosmo $VAL_ADDRESS_8 --from val8 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test  #1
#$TZ2_2_EXEC tx staking tokenize-share $VAL_VALOPER_9 25000000uosmo $VAL_ADDRESS_9 --from val9 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test  #2
#$TZ2_3_EXEC tx staking tokenize-share $VAL_VALOPER_10 65000000uosmo $VAL_ADDRESS_10 --from val10 --gas 400000 --chain-id $CHAINID_2 -y --keyring-backend=test  #3

sleep 5
DEPOSIT_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[1].deposit_address.address -r)
while [[ "$DEPOSIT_ACCOUNT" == "null" ]]; do
  sleep 5
  DEPOSIT_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[1].deposit_address.address -r)
done

PERFORMANCE_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[1].performance_address.address -r)
while [[ "$PERFORMANCE_ACCOUNT" == "null" ]]; do
  sleep 2
  PERFORMANCE_ACCOUNT=$($QS1_EXEC q interchainstaking zones --output=json | jq .zones[1].performance_address.address -r)
done
sleep 2
$TZ2_1_EXEC tx bank send val8 $PERFORMANCE_ACCOUNT 40000uosmo --chain-id $CHAINID_2 -y --keyring-backend=test

sleep 3
$TZ2_1_EXEC tx bank send demowallet8 $DEPOSIT_ACCOUNT 50000000uatom --chain-id $CHAINID_2 -y --keyring-backend=test
