# Interchain Queries Relayer

The Interchain Queries (ICQ) Relayer watches for events emitted by the ICQ module. It makes lookups against external chains, and returns query results and proofs such that the ICQ module is able to verify proofs and trigger the appropriate downstream action.

## Configuration

The ICQ Relayer configuration is controlled by a single YAML file, the default path of which is $HOME/.icq/config.

The first run of `icq-relayer` will generate a mainnet compatible config file, if one is not present.

## Changelog

### v0.11.0
- Add support for cosmos-sdk v0.50 GetTxsEvents request type
- Make metrics bind port configurable
- Set default config file to be mainnet ready
- Reduce log verbosity
- Add max_msgs_per_tx congig variable
- Dynamic MsgPerTx: Make the MsgPerTx value reduce on failed requests, and slowly return to MaxMsgsPerTx over time on success

### v0.10.0
- Add CometBFT v0.37 compatibility.

### v0.9.0
- Add caching for significant performance improvment.

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

