#!/bin/bash -i
set -e

source vars.sh

if [[ ! -f $QS_BIN ]]; then
        echo "Run 'make init' before running this command."
fi

$QS_BIN query staking validators -o json | \
        jq .validators[] | \
        jq -s 'sort_by(.tokens) | reverse' | \
        jq -r '["Validator", "VP"], ["----------------", "------------"], (.[] | [.description.moniker, (.tokens|tonumber/1000000)]) | @tsv' | \
        column -t -s "$(printf '\t')"
