# Changelog

## Unreleased

- Cosmos-SDK 0.47 upgrade
- IBC-go v7.9.2 upgrade
- CometBFT v0.37.15 upgrade

## Released

### v1.7.7

- determinsitic config filtering by @joe-bowman in #1813
- minor deps bumps
- v1.7.7 release by @joe-bowman in #1821

### v1.7.6

- add archway and celestia to web-app by @joe-bowman in #1771
- Arham/auto register host zone by @arhamchordia in #1784
- remove old ics code that is unused by @joe-bowman in #1736
- Numerous dependencies upgrades
- migrate all proto instances of sdk.Int to math.Int by @joe-bowman in #1799
- Remove deprecated airdrop module by @joe-bowman in #1803
- Webui 1.4 by @joe-bowman in #1801
- Update ledger libraries for work with Ledger OC up to 2.4.1 by @alibabaedge in #1804
- add topN endpoint to supply module (disabled by default) by @joe-bowman in #1800
- Deps and perf by @joe-bowman in #1808
- v1.7.6 upgradehandler by @joe-bowman in #1809
- add a default TTL to non-periodic queries, and GC stale queries by @joe-bowman in #1812

### v1.7.5

- Icq v1.1.2 by @joe-bowman in #1762
- fix git describe commands to fetch versions by @joe-bowman in #1763
- fix: ensure we dont attempt to send funds from deposit to delegate after we refund user by @joe-bowman in #1764
- update rr logic to exclude queued unbondings; fixes #1760 by @joe-bowman in #1765
- use bech32 addresses for mapped address queries; resolves #1719 by @joe-bowman in #1766
- add upgrade handler for v1.7.5 by @joe-bowman in #1767

### v1.7.4

- Fixes issue with balanceWaitgroup decrement < 0
- Revert cometbft-db to v0.12, to resolve ICA issue
- v1.7.4 upgrade handler
- remove old withdrawal not removed by v1.7.2 handler
- requeue failed cosmos unbonding
- remove two old icq queries that can never be completed

### v1.7.3

- Fixes handling of Celestia Inclusion proofs, refactor inclusion proofs for scalability, to support additional inclusion proof types later 
- Improve unbonding handling, by adding completion time field to unbonding record, and populating on acknowledgement; this will allow us to garbage collect unbonding records in c. 7 epochs time
- Improve RR computation reliability, by making RR logic ignore unbonding and escrowed tokens; as unbonded tokens have already been deducted at a given rate, it doesn't make sense include these in the next calc.
- chore(deps): bump codecov/codecov-action from 4 to 5 by @dependabot in #1743
- fix nil reference panic in abci.go by @joe-bowman in #1745
- add v1.7.2 upgrade handler; this removes a duplicate record and updates RRs
- chore(deps): bump github.com/cometbft/cometbft-db from 0.14.1 to 1.0.1 by @dependabot in #1751
- chore(deps): bump google.golang.org/grpc from 1.67.1 to 1.68.0 by @dependabot in #1754
- chore(deps): bump github.com/spf13/cast from 1.6.0 to 1.7.0 by @dependabot in #1752
- bump math v1.3.0 -> v1.4.0 by @joe-bowman in #1750

### v1.6.4

- Web UI updates by @joe-bowman in #1718
- ensure we use the correct zone denom where evaluating pool claims by @joe-bowman in #1721
- Fix: broken links detected in issue #1716 by @juliogarciape in #1717
- fix: stylecheck lint error on chainId in osmosis submodule by @ajansari95 in #1724
- update: CODEOWNERS by @ajansari95 in #1725
- cleanup: remove ic-test by @ajansari95 in #1727
- fix: cleanup and bumps by @ajansari95 in #1730
- fix: use %q instead of manual "%s" by @odeke-em in #1723
- fix: stuck unbondings and UpdateRedemption by @ajansari95 in #1732
- remove unused error message that causes us to fail to redelegate properly for amounts > int64 by @joe-bowman in #1735
- belts and braces around scaling factors to avoid panic by @joe-bowman in #1734

### v1.6.3

- refactor(x/participationrewards/keeper): combine GetProtocolData+UnmarshalProtocolData by @odeke-em in #1681
- test(x/interchainquery): directly assert AppModuleBasic.DefaultGenesis return non-nil JSON by @odeke-em in #1700
- chore(deps): bump oven-sh/setup-bun from 1 to 2 by @dependabot in #1671
- Unbond recovery by @joe-bowman in #1704
- chore(deps): bump the go_modules group across 1 directory with 2 updates by @dependabot in #1710
- Release/v1.6.3 by @joe-bowman in #1713

### v1.6.2

- fix(x/airdrop/keeper): hoist out sdk.Dec from constant strings by @odeke-em in #1677
- ci(icq-relayer): run golangci on pull request related to icq-relayer by @tropicaldog in #1674
- fix: fix some broken links in CONTRIBUTING.md and PAID_BOUNTIES.md by @odeke-em in #1682
- add ibc v6 migration to upgrade handler
- fix v1.6.1 upgrade handler test
- dont error on nil response, it is valid if node has been pruned. log the error and try again. query will be re-raised if balance remains non-zero
- fix: do not throw error on race condition (#1694)
- fix: dont fail acknowledgement on delgation record race condition. fixes #1693
- ensure that src delegation is icq updated in all cases
- pfm can return nil; transfer_middleware should fall through if this is the case; fixes #1695
- add no-op handler for v1.6.2

### v1.5.7

- Security hotfix to fix a bug that permitted repeated claims for the same tokens.

### v1.5.6

- Upgrade handler to fix 4x unbondings.

### v1.5.5

- v1.5.5 upgrade handler
    - Migrate one user killer queen rewards
    - Add Saga and Dydx denoms on Osmosis for claims
    - Remove stale withdrawal record for 0 tokens
    - Re-emit failed 8x epoch 148 unbonding distributions
    - Re-queue 2x unbondings stuck due to previously fixed bug in #1347
- Fix bug where wrong supply was used to determine distribution proportions for non-staking tokens.
- Set per chain thresholds for min delegations to avoid many tiny delegation requests which are inefficient

### v1.5.4

- Add missing queries CLI by @tropicaldog in #1233
- Update deploytestweb.yaml by @joe-bowman in #1341
- add kq migration to v1.6 upgrade handler by @joe-bowman in #1342
- fix(x/mint/simulation): use proper value of maxInt64 by @odeke-em in #1339
- update osmosis Dockerfile by @joe-bowman in #1343
- when flushing tokens, consider inflight unbondings (#996) by @joe-bowman in #1347
- Base signalling intent off the most recent claims by @joe-bowman in #1345
- Add Gov messages for validator deny list. by @tropicaldog in #1329
- upgrade pfm to v5.2.2 by @joe-bowman in #1356
- Non staking denom rewards by @joe-bowman in #1355

### v1.5.3

- Fix: remove usages of int64

### v1.5.2

- Mainline support for pebbledb.

###v1.5.1

- Fixes issue with stale acknowledgements causing a balance underflow.
- Fixes computed IBC denoms for cross chain claims.
- Fix a panic in the participation rewards epoch when there is no allocation to distribute.

### v1.5.0

- Fix race condition in Delegation record updates and withdrawal acknowledgements #1162
- Improve handling of requeued unbondings #1201
- Fix potential panic in claims manager checks #1217
- Remove Crescent claims logic #1143
- Remove unused datapoints functionality from ICQ #1170
- Fix issue with delayed redelegation acknowledgements being garbage collected #1140
- Add MsgCancelQueuedRedemption #1122
- Add unbonding statistics to /zones #1123
- Auto onboard protocol data for new zones #1121
- Handle unstakable tokens #1118
- Performance fixes #993 #1054
- Handle migration of vesting accounts with existing delegations #1175
- Fix panic is validator reduced to zero VP through slashing
- Compensate users that unbonded at a very low RR #1259
- Migrate 2x Notional multisigs to new addresses #1184

### v1.4.7

- Fix undelegation allocations in epoch
- Set redemtpion rate and trigger new RR update in upgrade handler
- Fix bug that permits undeliverable MsgSend to be generated when no rewards present
- Migrate testnet rewards account to original address

### v1.4.6

- Unbonding
- Claims
- Signalling Intent
- Participation Rewards
- LSM

### v1.2.17

- v1.2.17 fixes the ongoing redemption rate calculation issues present on cosmoshub-4
- Bump go to v1.20
- Bump cosmos-sdk and ibc-go
- Change some Info logs to Debug, to reduce log noise.

### v1.2.16

- fix: dep-bump : comet bft by @ajansari95 in #526

### v1.2.15

- fix: handle redelegation ack by @ajansari95 in #516

### v1.2.14

- barberry patch fix by @muku314115 in #453

### v1.2.13

- fix: security fix

### v1.2.12

- dependency upgrade for huckelberry to ibc-go:v5.2.1 by @ajansari95 in #437
- This release includes a fix for the huckleberry security advisory. Credits to Felix Wilhelm (@felixwilhelm) of Jump Crypto for the discovery and responsible disclosure via the cosmos bug bounty program.

### v1.2.11

- fix: check receipts and using safesub by @ajansari95 in #433

### v1.2.10

- fix: proper mergify filename (backport #384) by @mergify in #385
- test: interchainstaking zones cli query by @aljo242 in #382
- chore: minor updates bump (backport #403) by @mergify in #404
- backport: epochly flush logic by @ajansari95 in #424

### v1.2.9

- Add message_per_tx param to zone registration and update proposals
- Add feature to configure max-tx size per zone
- Set ICATimeout to 6h
- Add logic to clear pending delegations when a channel is reopened
- Bump cosmos-sdk to v0.46.11
- Bump comet-bft to 0.34.27
- Upgrade Handler:
  - set messages per tx a maximum of (max_gas per tx on host chain per block / 1m) for every zone
  - flush pending delegations
  - chore: backport lint fixes to main by @aljo242 in #363
- configurable msg per tx by @joe-bowman in #361
- changelog added by @ajansari95 in #364
- Feature/add msgs per tx to proposals by @joe-bowman in #365

### v1.2.7

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

- fix deposit address setWithdrawalAddress callback by @ajansari95 in #284

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
