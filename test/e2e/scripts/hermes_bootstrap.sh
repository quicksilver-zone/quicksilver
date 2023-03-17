#!/bin/bash

set -ex

# initialize Hermes relayer configuration
mkdir -p /root/.hermes/
touch /root/.hermes/config.toml
echo $QUICK_A_E2E_VAL_MNEMONIC > /root/.hermes/QUICK_A_MNEMONIC.txt
echo $QUICK_B_E2E_VAL_MNEMONIC > /root/.hermes/QUICK_B_MNEMONIC.txt
# setup Hermes relayer configuration
tee /root/.hermes/config.toml <<EOF
[global]
log_level = 'info'
[mode]
[mode.clients]
enabled = true
refresh = true
misbehaviour = true
[mode.connections]
enabled = false
[mode.channels]
enabled = true
[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true
[rest]
enabled = true
host = '0.0.0.0'
port = 3031
[telemetry]
enabled = true
host = '127.0.0.1'
port = 3001
[[chains]]
id = '$QUICK_A_E2E_CHAIN_ID'
rpc_addr = 'http://$QUICK_A_E2E_VAL_HOST:26657'
grpc_addr = 'http://$QUICK_A_E2E_VAL_HOST:9090'
websocket_addr = 'ws://$QUICK_A_E2E_VAL_HOST:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'quick'
key_name = 'val01-quick-a'
store_prefix = 'ibc'
max_gas = 6000000
gas_multiplier = 1.2
default_gas = 400000
gas_price = { price = 0.0025, denom = 'e2e-default-feetoken' }
clock_drift = '1m' # to accomdate docker containers
trusting_period = '239seconds'
trust_threshold = { numerator = '1', denominator = '3' }
[[chains]]
id = '$QUICK_B_E2E_CHAIN_ID'
rpc_addr = 'http://$QUICK_B_E2E_VAL_HOST:26657'
grpc_addr = 'http://$QUICK_B_E2E_VAL_HOST:9090'
websocket_addr = 'ws://$QUICK_B_E2E_VAL_HOST:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'quick'
key_name = 'val01-quick-b'
store_prefix = 'ibc'
max_gas = 6000000
gas_multiplier = 1.2
default_gas = 400000
gas_price = { price = 0.0025, denom = 'e2e-default-feetoken' }
clock_drift = '1m' # to accommodate docker containers
trusting_period = '239seconds'
trust_threshold = { numerator = '1', denominator = '3' }
EOF

# import keys
hermes keys add --chain ${QUICK_B_E2E_CHAIN_ID} --key-name "val01-quick-b" --mnemonic-file /root/.hermes/QUICK_B_MNEMONIC.txt
hermes keys add --chain ${QUICK_A_E2E_CHAIN_ID} --key-name "val01-quick-a" --mnemonic-file /root/.hermes/QUICK_A_MNEMONIC.txt

# start Hermes relayer
hermes start