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
    rpc-addr: 'http://$QUICK_A_E2E_VAL_HOST:26657'
    grpc-addr: 'http://$QUICK_A_E2E_VAL_HOST:9090'
    account-prefix: quick
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uqck
    min-gas-amount: 0
    key-directory: /root/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct
  '$QUICK_B_E2E_CHAIN_ID':
    key: default
    chain-id: '$QUICK_B_E2E_CHAIN_ID'
    rpc-addr: 'http://$QUICK_B_E2E_VAL_HOST:26657'
    grpc-addr: 'http://$QUICK_B_E2E_VAL_HOST:9090'
    account-prefix: quick
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uqck
    min-gas-amount: 0
    key-directory: /root/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct
EOF

interchain-queries run