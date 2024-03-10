#!/bin/bash -i
set -e

source vars.sh

if [[ ! -f $QS_BIN ]]; then
    echo "Run 'make init' before running this command."
fi

QS_KEY=$($QS_BIN --home $QS_HOME keys show validator --output=json | jq .address -r)

echo "Balance of $QS_KEY on $CHAIN_ID:"
$QS_BIN --home $QS_HOME q bank balances $QS_KEY --node http://seed.${CHAIN_ID}.quicksilver.zone:26657 --output=json | jq .balances
echo
