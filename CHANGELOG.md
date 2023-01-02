# Changelog

## Released

### v1.1.1

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
