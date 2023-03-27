# Changelog

## Released

### v1.2.8
- Add feature to configure max-tx size per zone
- Set `ICATimeout` to 6h
- Add logic to clear pending delegations when a channel is reopened
- Bump `cosmos-sdk` to `v0.46.11`
- Bump `comet-bft` to `0.34.27`
- Upgrade Handler:
  - set messages per tx a maximum of (`max_gas per tx on hostchain per block` / `1m`) for every zone
  - flush outstanding delegations 

### v1.2.7
- fix upgrade handler

### v1.2.6
- fix release height

### v1.2.5
- Backport `GovRepenChannel` tx and `GovCloseChannel` proposal
- Add Logic to handle pending delegations on `regen-1` post upgrade once channel are open.
- Upgrade Handler:
  - Fix `EpochProvisions` value 
  - Set DistributionProportions params

### v1.2.4
- Bump Tendermint to v0.34.26 (informalsystems/tendermint)
- Bump Cosmos-SDK to v0.46.10
- Add v0.46.x ICA Callback handler support
- Add adjacent block verification to Tx validation callback
- Add first_seen and completion_time to receipts to allow for better monitoring of tx processing throughput
- Add callback to remove redelegation records in event of transient redelegation
- Upgrade handler:
  - remove redelegation records with nil-timestamp (failed records)
  - update existing receipts timestamps to upgrade block timestamp

### v1.2.3
- check for nil allocation
- support v0.46 balance callbacks

### v1.2.2
- fix deposit address setWithdrawalAddress

### v1.2.1
- Do not remove delegation records on absence from DelegatorDelegations response
- Add tests for delegations callback changes
- Fix deposit address onboarding race condition
- Fix GetTxEvents pagination and sort issue
- Ensure Undelegate receipt doesn't trigger entire valset refresh
- Add uncapped OverrideRedemptionRate function; add tests; add upgrade handler for block 325001

### v1.2.0

- Remove v1.1.0 migration and epoch skipping logic
- Tidy debug logging
- Fix inverted comparison for unbonding clipping
- Bump ibc-go to v5.2.0

### v1.1.0

- Add escrow account to fix unbonding bug
- Filter zero intents
- Remove burn permissions from mint and participation reward accounts
- Remove pre-v1.0.0 upgrade handlers
- Improve error responses for failed redelegation/withdrawal callbacks
- Direct poolincentives to airdrop module account
- Replace redemption rate hard cap with soft cap
- Add chain restart migration to remove cosmoshub-4, close ICA channels and burn minted qAtoms

### v1.0.0

- First production release
- Bump ibc-go to v5.1.0
- Bump tendermint to v0.34.24
- Bump cosmos-sdk to v0.46.7
