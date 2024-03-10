#!/bin/bash -i
set -e

source vars.sh

if [[ ! -f $QS_BIN ]]; then
    echo "Run 'make init' before running this command."
fi

QS_KEY=$($QS_BIN --home $QS_HOME keys show validator --output=json | jq .address -r)

## Unjail validator
$QS_BIN tx slashing unjail \
  --from=validator \
  --chain-id=$CHAIN_ID \
  --home $QS_HOME
