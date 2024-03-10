#!/bin/bash -i
set -e

source vars.sh

if [[ ! -f $QS_BIN ]]; then
    echo "Run 'make init' before running this command."
fi

QS_KEY=$($QS_BIN --home $QS_HOME keys show validator --output=json | jq .address -r)

echo "On the Quicksilver discord:"
echo "   - in the #qck-tap channel, enter: '\$request $QS_KEY rhapsody'"
