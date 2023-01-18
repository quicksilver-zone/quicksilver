# Changelog

## Released

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
