#!/bin/bash -i
set -e

source vars.sh

if [[ ! -f $QS_BIN ]]; then
        echo "Run 'make init' before running this command."
fi
echo "Checking validator rewards..."

rewards=$($QS_BIN query distribution rewards $($QS_BIN --home $QS_HOME keys show validator -a) -o json | jq '.total | to_entries' | jq -r ".[] | select(.value.denom == \"uqck\") | .value.amount")
if [[ -n "$rewards" ]]; then
        rewards_quick=$(echo "$rewards / 1000000" | bc)

        echo "${rewards_quick} QCK / ${rewards} uqck"
else
        echo "No rewards"
fi

echo "Checking validator commission..."

commission=$($QS_BIN query distribution commission $($QS_BIN --home $QS_HOME keys show validator --bech val -a) -o json | jq '.commission | to_entries' | jq -r ".[] | select(.value.denom == \"uqck\") | .value.amount")
if [[ -n "$commission" ]]; then
        commission_quick=$(echo "$commission / 1000000" | bc)

        echo "${commission_quick} QCK / ${commission} uqck"
else
        echo "No commissions"
fi
