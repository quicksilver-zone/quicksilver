#!/bin/bash -i
set -e

source vars.sh

if [[ ! -f $QS_BIN ]]; then
    echo "Run 'make init' before running this command."
fi

QS_KEY=$($QS_BIN --home $QS_HOME keys show validator --output=json | jq .address -r)

echo -n "Enter your validator name: "
read VAL_NAME

## Upgrade node to validator
$QS_BIN tx staking create-validator \
  --from=validator \
  --amount=5000000uqck \
  --moniker=$VAL_NAME \
  --chain-id=$CHAIN_ID \
  --commission-rate=0.1 \
  --commission-max-rate=0.5 \
  --commission-max-change-rate=0.1 \
  --min-self-delegation=1 \
  --home $QS_HOME \
  --pubkey=$($QS_BIN --home $QS_HOME tendermint show-validator)
