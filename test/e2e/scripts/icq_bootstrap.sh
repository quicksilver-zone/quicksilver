#!/bin/bash

set -ex

# initialize icq relayer configuration
mkdir -p /root/.icq/
touch /root/.icq/config.yaml

tee /root/.icq/config.yaml <<EOF
default_chain: '$QUICK_A_E2E_CHAIN_ID'
chains:
  '$QUICK_A_E2E_CHAIN_ID':
    key: default
    chain-id: '$QUICK_A_E2E_CHAIN_ID'
    rpc-addr: https://rpc.quicksilver.zone:443
    grpc-addr: https://grpc.quicksilver.zone:443
    account-prefix: quick
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uqck
    min-gas-amount: 0
    key-directory: /root/joe/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct
  '$QUICK_B_E2E_CHAIN_ID':
    key: default
    chain-id: '$QUICK_B_E2E_CHAIN_ID'
    rpc-addr: https://osmosis-1.technofractal.com:443
    grpc-addr: https://gprc.osmosis-1.technofractal.com:443
    account-prefix: osmo
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uosmo
    min-gas-amount: 0
    key-directory: /root/joe/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct
EOF

interchain-queries run