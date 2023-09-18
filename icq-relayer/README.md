# Interchain Queries Relayer

The Interchain Queries (ICQ) Relayer watches for events emitted by the ICQ module. It makes lookups against external chains, and returns query results and proofs such that the ICQ module is able to verify proofs and trigger the appropriate downstream action.

## Configuration

The ICQ Relayer configuration is controlled by a single YAML file, the default path of which is $HOME/.icq/config.

```yaml
default_chain: quicksilver-1
chains:
  quicksilver-1:
    key: default
    chain-id: quicksilver-1
    rpc-addr: https://rpc.quicksilver.zone:443
    grpc-addr: https://grpc.quicksilver.zone:443
    account-prefix: quick
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uqck
    min-gas-amount: 0
    key-directory: /home/joe/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct
  osmosis-1:
    key: default
    chain-id: osmosis-1
    rpc-addr: https://osmosis-1.technofractal.com:443
    grpc-addr: https://gprc.osmosis-1.technofractal.com:443
    account-prefix: osmo
    keyring-backend: test
    gas-adjustment: 1.2
    gas-prices: 0.01uosmo
    min-gas-amount: 0
    key-directory: /home/joe/.icq/keys
    debug: false
    timeout: 20s
    block-timeout: 10s
    output-format: json
    sign-mode: direct

```

## Changelog
### v0.8.2
- Improved efficiency
- More detailed metrics
- Don't panic on failed txs
- Don't query client headers that are going to be rejected

### v0.8.0
- Improved error handling
- Add metrics

### v0.6.2
- Fix default chain instantiation on first run

### v0.6.1
- Fix wg.Wait() deferal bug (#4)

### v0.6.0

- Add structured logging.
- Update Quicksilver to v0.9.0

### v0.5.0

- Upgrade SDK to v0.46
- Upgrade Quicksilver to v0.8.0

