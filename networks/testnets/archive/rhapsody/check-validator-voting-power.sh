#!/bin/bash -i
set -e

source vars.sh

if [[ ! -f $QS_BIN ]]; then
        echo "Run 'make init' before running this command."
fi

vp=$($QS_BIN status | jq '.ValidatorInfo.VotingPower')
if [[ $vp = "0" ]]
then
        status="JAILED"
else
        status="OK"
fi

echo "Voting Power: ${vp} [${status}]"
