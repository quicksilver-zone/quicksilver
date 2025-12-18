# ADR-003: Zone Offboarding

## Status

Implemented

## Context

When a liquid staking zone needs to be wound down (e.g., when a chain is sunsetting), there must be a safe and orderly process to:
1. Stop new deposits and user-initiated redemptions
2. Return funds to users who have pending redemptions
3. Unbond all staked tokens from the remote chain
4. Transfer unbonded funds to a multisig for final distribution

This ADR describes the governance-controlled zone offboarding process.

## Decision

Implement a multi-step offboarding process using three new governance messages:
- `MsgGovSetZoneOffboarding`
- `MsgGovCancelAllPendingRedemptions`
- `MsgGovForceUnbondAllDelegations`

## Execution Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ZONE OFFBOARDING FLOW                               │
└─────────────────────────────────────────────────────────────────────────────┘

Step 1: Enable Offboarding
┌─────────────────────────────────────────┐
│  MsgGovSetZoneOffboarding               │
│  (chain_id, is_offboarding=true)        │
├─────────────────────────────────────────┤
│  - Sets zone.IsOffboarding = true       │
│  - Sets zone.DepositsEnabled = false    │
│  - Sets zone.UnbondingEnabled = false   │
│  - Freezes redemption rate updates      │
└─────────────────────────────────────────┘
                    │
                    ▼
Step 2: Cancel Pending Redemptions
┌─────────────────────────────────────────┐
│  MsgGovCancelAllPendingRedemptions      │
│  (chain_id)                             │
├─────────────────────────────────────────┤
│  - Iterates all QUEUED withdrawals      │
│  - Iterates all UNBOND withdrawals      │
│  - Refunds qAssets from escrow to users │
│  - Deletes withdrawal records           │
│  - Returns: cancelled_count,            │
│             refunded_amounts            │
└─────────────────────────────────────────┘
                    │
                    ▼
Step 3: Force Unbond All Delegations
┌─────────────────────────────────────────┐
│  MsgGovForceUnbondAllDelegations        │
│  (chain_id)                             │
├─────────────────────────────────────────┤
│  - Gets all delegations for zone        │
│  - Creates MsgUndelegate for each       │
│  - Submits via ICA to remote chain      │
│  - Returns: unbonding_count,            │
│             total_unbonded              │
└─────────────────────────────────────────┘
                    │
                    ▼
Step 4: Wait for Unbonding Period
┌─────────────────────────────────────────┐
│  ~21 days (chain-dependent)             │
├─────────────────────────────────────────┤
│  - Tokens unbonding on remote chain     │
│  - Funds move to DelegationAddress      │
│    after unbonding completes            │
└─────────────────────────────────────────┘
                    │
                    ▼
Step 5: Transfer to Multisig
┌─────────────────────────────────────────┐
│  MsgGovExecuteICATx (existing)          │
│  (chain_id, MsgSend to multisig)        │
├─────────────────────────────────────────┤
│  - Sends unbonded tokens from           │
│    DelegationAddress to multisig        │
│  - Multisig handles distribution        │
└─────────────────────────────────────────┘
```

## What Gets Blocked

| Operation | Status | Description |
|-----------|--------|-------------|
| New Deposits | Blocked | Users cannot deposit native tokens |
| User Redemptions | Blocked | Users cannot request new unstaking |
| Redemption Rate Updates | Blocked | RR frozen at last value |
| Governance Actions | Allowed | Force unbond, cancel redemptions, ICA txs |

## Governance Messages

### 1. MsgGovSetZoneOffboarding

Enables or disables offboarding mode for a zone.

```protobuf
message MsgGovSetZoneOffboarding {
  string chain_id = 1;
  bool is_offboarding = 2;
  string authority = 3;  // gov module address
}
```

**Effects when enabling offboarding:**
- `zone.IsOffboarding = true`
- `zone.DepositsEnabled = false`
- `zone.UnbondingEnabled = false`

### 2. MsgGovCancelAllPendingRedemptions

Cancels all pending redemptions and refunds qAssets to users.

**Precondition:** Zone must be in offboarding mode.

```protobuf
message MsgGovCancelAllPendingRedemptions {
  string chain_id = 1;
  string authority = 2;  // gov module address
}

message MsgGovCancelAllPendingRedemptionsResponse {
  uint64 cancelled_count = 1;
  repeated cosmos.base.v1beta1.Coin refunded_amounts = 2;
}
```

**Behavior:**
- Iterates all withdrawal records in `QUEUED` status
- Iterates all withdrawal records in `UNBOND` status
- For each record:
  - Refunds `BurnAmount` (qAssets) from escrow module to user
  - Deletes the withdrawal record
- Emits cancellation events for each record

### 3. MsgGovForceUnbondAllDelegations

Initiates unbonding of all staked tokens via ICA.

**Precondition:** Zone must be in offboarding mode.

```protobuf
message MsgGovForceUnbondAllDelegations {
  string chain_id = 1;
  string authority = 2;  // gov module address
}

message MsgGovForceUnbondAllDelegationsResponse {
  uint64 unbonding_count = 1;
  cosmos.base.v1beta1.Coin total_unbonded = 2;
}
```

**Behavior:**
- Gets all delegations for the zone
- Creates `MsgUndelegate` for each delegation
- Submits transactions via ICA with memo `offboard/<block_height>`
- Increments withdrawal waitgroup for tracking

## Example Governance Proposals

### Step 1: Enable Offboarding

```json
{
  "messages": [
    {
      "@type": "/quicksilver.interchainstaking.v1.MsgGovSetZoneOffboarding",
      "chain_id": "omniflixhub-1",
      "is_offboarding": true,
      "authority": "quick10d07y265gmmuvt4z0w9aw880jnsr700jvss730"
    }
  ],
  "title": "Enable Offboarding for OmniFlix Zone",
  "summary": "This proposal enables offboarding for the OmniFlix zone due to chain sunset."
}
```

### Step 2: Cancel Pending Redemptions

```json
{
  "messages": [
    {
      "@type": "/quicksilver.interchainstaking.v1.MsgGovCancelAllPendingRedemptions",
      "chain_id": "omniflixhub-1",
      "authority": "quick10d07y265gmmuvt4z0w9aw880jnsr700jvss730"
    }
  ],
  "title": "Cancel Pending Redemptions for OmniFlix Zone",
  "summary": "This proposal cancels all pending redemptions and refunds qFLIX to users."
}
```

### Step 3: Force Unbond All Delegations

```json
{
  "messages": [
    {
      "@type": "/quicksilver.interchainstaking.v1.MsgGovForceUnbondAllDelegations",
      "chain_id": "omniflixhub-1",
      "authority": "quick10d07y265gmmuvt4z0w9aw880jnsr700jvss730"
    }
  ],
  "title": "Force Unbond All Delegations for OmniFlix Zone",
  "summary": "This proposal initiates unbonding of all staked FLIX tokens."
}
```

### Step 5: Transfer to Multisig (after unbonding period)

```json
{
  "messages": [
    {
      "@type": "/quicksilver.interchainstaking.v1.MsgGovExecuteICATx",
      "chain_id": "omniflixhub-1",
      "address": "<delegation_address>",
      "msgs": [
        {
          "@type": "/cosmos.bank.v1beta1.MsgSend",
          "from_address": "<delegation_address>",
          "to_address": "<multisig_address>",
          "amount": [{"denom": "uflix", "amount": "<total_amount>"}]
        }
      ],
      "authority": "quick10d07y265gmmuvt4z0w9aw880jnsr700jvss730"
    }
  ],
  "title": "Transfer Unbonded FLIX to Multisig",
  "summary": "This proposal transfers unbonded FLIX to the community multisig for distribution."
}
```

## Important Notes

1. **Order matters**: Steps must be executed in order. Steps 2 and 3 require offboarding to be enabled first.

2. **Unbonding period**: After Step 3, you must wait for the chain's unbonding period (~21 days typically) before funds can be transferred.

3. **qAsset holders**: Users holding qAssets after offboarding is complete will need an alternative redemption mechanism (e.g., multisig-managed distribution based on final redemption rate).

4. **Reversibility**: While offboarding can technically be disabled via another governance proposal, the process is designed to be a one-way operation for sunsetting zones.

5. **ICA Channel**: The ICA channels must remain open for Steps 3 and 5 to work. If channels are closed, they must be reopened via `MsgGovReopenChannel`.

## File Changes

### Proto Files
- `proto/quicksilver/interchainstaking/v1/interchainstaking.proto`: Added `is_offboarding` field to `Zone`
- `proto/quicksilver/interchainstaking/v1/proposals.proto`: Added 3 new message types
- `proto/quicksilver/interchainstaking/v1/messages.proto`: Added 3 new RPC methods

### Go Files
- `x/interchainstaking/types/msgs.go`: Message validation methods
- `x/interchainstaking/types/codec.go`: Message registration
- `x/interchainstaking/types/events.go`: Event types and attribute keys
- `x/interchainstaking/types/ibc_packet.go`: Offboarding memo type
- `x/interchainstaking/keeper/msg_server.go`: Handler implementations
- `x/interchainstaking/keeper/ibc_packet_handlers.go`: Offboarding unbond handling, RR skip

## Consequences

### Positive
- Safe and orderly wind-down of zones
- Users with pending redemptions get their qAssets back immediately
- Governance-controlled process ensures community oversight
- Reuses existing ICA infrastructure

### Negative
- Requires multiple governance proposals
- qAsset holders after offboarding need separate distribution mechanism
- Unbonding period delay before funds can be transferred

### Neutral
- Process is largely irreversible by design
- Requires ICA channels to remain operational

## References

- [Quicksilver ICS Module](../x/interchainstaking/)
- [ICA (Interchain Accounts)](https://ibc.cosmos.network/main/apps/interchain-accounts/overview.html)
