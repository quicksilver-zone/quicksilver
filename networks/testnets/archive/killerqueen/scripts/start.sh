#!/bin/bash -i
set -e

source vars.sh

if [[ ! -z $(pgrep quicksilverd) ]]; then
    echo "quicksilverd is already running; you should run 'make stop' to kill the existing process"
fi

if [[ ! -f $QS_BIN ]]; then
    echo "Run 'make init' before running this command."
fi

$QS_BIN --home $QS_HOME start > qs.log 2>&1 &
