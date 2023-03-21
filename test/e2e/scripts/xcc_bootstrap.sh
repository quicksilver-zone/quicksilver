#!/bin/sh

set -ex

# initialize icq relayer configuration
mkdir -p /root/.xcc/
touch /root/.xcc/config.yaml

tee /root/.xcc/config.yaml <<EOF
source_chain: '$QUICK_A_E2E_CHAIN_ID'
chains:
  quickquick-1: 'http://$QUICK_B_E2E_VAL_HOST:26657'
EOF

xcc -a serve -f /root/.xcc/config.yaml

