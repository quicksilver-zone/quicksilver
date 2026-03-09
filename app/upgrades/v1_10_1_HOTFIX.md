# v1.10.1 Hotfix — HandleFailedUndelegate BurnAmount Zero Guard

## Summary

Consensus-breaking hotfix for a bug in `HandleFailedUndelegate` that blocks all ICA delegate channel ack relay across multiple zones. The bug causes `SetWithdrawalRecord` to reject withdrawal records with zero `BurnAmount`, which reverts the entire `OnAcknowledgementPacket` state transition. On ordered ICA channels this creates a complete deadlock — no acks can be committed, so no new packets can be delivered either.

## Upgrade Details

- **Upgrade Name**: `v1.10.1`
- **Handler**: `NoOpHandler` (no state migration required)
- **Store Changes**: None
- **Consensus Breaking**: Yes — changes ack callback behavior

No state migration is needed because the bug causes full tx revert on every failed ack relay. No partial writes from the failed ack relay path are committed on-chain. All withdrawal records remain in their original `WithdrawStatusUnbond` state, untouched.

## Root Cause

`HandleFailedUndelegate` in `x/interchainstaking/keeper/ibc_packet_handlers.go` processes failed `MsgUndelegate` acknowledgements. For each withdrawal record (WDR) linked to the unbonding record, it:

1. Finds how much of the WDR was delegated to the failing validator (`relatedAmount` from `WDR.Distribution`)
2. Computes the proportional qAsset burn: `relatedQAsset = floor(relatedAmount * BurnAmount / TotalAmount)`
3. Subtracts `relatedQAsset` from the original WDR's `BurnAmount`
4. Creates a requeued withdrawal record with `BurnAmount = relatedQAsset`

Three edge cases cause `SetWithdrawalRecord` to reject the record with "burnAmount cannot be negative or zero":

### Edge Case 1: Validator not in WDR distribution

When the unbonding record references a validator that is NOT present in the WDR's distribution list, `relatedAmount = 0`, so `relatedQAsset = floor(0 * rr) = 0`. Step 4 creates a requeued record with `BurnAmount = 0`.

This is the **primary trigger on mainnet**. On sommelier-3, the WDR has a 2-validator distribution, but 25 of 36 unbonding records reference validators absent from that distribution.

### Edge Case 2: TruncateInt rounds to zero

When `BurnAmount` is very small relative to `TotalAmount`, the redemption rate `rr` is tiny. Even with a nonzero `relatedAmount`, `floor(relatedAmount * rr)` can truncate to zero.

### Edge Case 3: Rounding overshoot

When `relatedQAsset >= BurnAmount` due to integer truncation across multiple iterations, `SubAmount` produces a zero or negative `BurnAmount` on the original WDR, which also fails `SetWithdrawalRecord` validation.

## Fix

Five guards added to `HandleFailedUndelegate`:

### Guard 1 — Skip if validator not in distribution (line 1058-1061)

```go
if relatedAmount.IsZero() {
    k.Logger(ctx).Info("validator not in distribution; skipping", ...)
    continue
}
```

Early exit when the unbonding record's validator is not present in the WDR's distribution. This prevents unnecessary processing and is the primary guard that unblocks mainnet (25 of 36 sommelier unbonding records reference validators absent from the WDR).

### Guard 2 — Skip if Amount is zero (line 1065-1068)

```go
if amount.IsZero() {
    k.Logger(ctx).Error("withdrawal record amount is zero; skipping", ...)
    continue
}
```

Defensive guard against division by zero when computing `rr = BurnAmount/Amount`. This protects against corrupt state where Amount has wrong denom or is empty.

### Guard 3 — Cap relatedQAsset at BurnAmount (line 1073-1076)

```go
if relatedQAsset.GT(wdr.BurnAmount.Amount) {
    relatedQAsset = wdr.BurnAmount.Amount
}
```

Prevents `SubAmount` from producing a negative or zero `BurnAmount` on the original WDR when rounding overshoots. Defensive against corrupt state where distribution amount exceeds total amount.

### Guard 4 — Delete instead of save when BurnAmount exhausted (line 1087-1090)

```go
if !wdr.BurnAmount.IsPositive() {
    k.DeleteWithdrawalRecord(ctx, wdr.ChainId, wdr.Txhash, wdr.Status)
} else {
    k.SetWithdrawalRecord(ctx, wdr)
}
```

If after subtraction the original WDR's `BurnAmount` is zero (rounding consumed it entirely), delete the record cleanly instead of attempting to save it.

### Guard 5 — Skip requeue when relatedQAsset is zero (line 1099-1102)

```go
if relatedQAsset.IsZero() {
    continue
}
```

When `TruncateInt` rounds to zero (tiny BurnAmount relative to Amount), there is no qAsset value to requeue. Skip creating the requeued record entirely.

## Impact — Channels Unblocked by This Fix

### Directly unblocked (ack relay resumes)

| Zone | Channel | Pending | What clears |
|------|---------|---------|-------------|
| cosmoshub-4 | ch-251 (delegate) | 2,055 | 647 acks drain via guard 3, then 1,408 unreceived packets flow |
| osmosis-1 | ch-252 (delegate) | 942 | 901 acks drain, then 41 unreceived packets flow |
| sommelier-3 | ch-225 (delegate) | 56 | All 56 acks drain (acks only, no unreceived) |

### Not affected by this bug (blocked by expired counterparty clients)

| Zone | Channels | Pending | Additional action needed |
|------|----------|---------|------------------------|
| regen-1 | ch-249 delegate, ch-19 withdrawal, ch-17 transfer | 222 | Revive `07-tendermint-113` on regen-1 |
| celestia | ch-211 delegate, ch-210 performance | 2,878 | Revive `07-tendermint-92` on celestia |
| agoric-3 | ch-240 delegate, ch-177 deposit, ch-227 performance | 954 | Revive `07-tendermint-85` on agoric-3 |
| juno-1 | ch-257 delegate, ch-90 deposit, ch-191 performance | 274 | Revive counterparty client on juno-1 |

### Not actionable

| Zone | Reason |
|------|--------|
| injective-1 | Both sides expired |
| omniflixhub-1 | Sunsetted |

## How Ordered Channel Deadlock Works

ICA channels use `ORDER_ORDERED`. The counterparty enforces sequential packet delivery:

```text
cosmoshub-4 ch-1344 (counterparty of QS ch-251):
  next_recv_seq = 3460

QS ch-251 commitments:
  seqs 2813..3459 → already received by cosmoshub, acks pending on QS
  seqs 3460..4867 → not yet delivered to cosmoshub

Deadlock:
  - cosmoshub won't accept seq 3460 until it has processed seq 3459
  - seq 3459 was already received, but its ack hasn't been committed on QS
  - the ack for seq 3459 triggers HandleFailedUndelegate → BurnAmount bug → reverts
  - ack stays pending → seq 3460 stays blocked → all subsequent seqs blocked
```

After the hotfix, the ack for seq 2813 (and all subsequent acks) processes successfully. The dam breaks, all 647 acks drain, then the 1,408 unreceived packets can be delivered.

## Test Coverage

Five new test functions validate the fix against mainnet-equivalent state and defensive edge cases:

### TestHandleFailedUndelegate_ValidatorNotInDistribution

Reproduces the exact mainnet sommelier scenario: WDR with 2-validator distribution, unbonding record for a validator NOT in the distribution. Without the fix, `SetWithdrawalRecord` rejects the zero `BurnAmount`. With guard 1, processing is skipped early.

### TestHandleFailedUndelegate_TinyBurnAmountTruncatesToZero

Tests the edge case where `BurnAmount=1` and `Amount=1,000,000` — the redemption rate is so small that `floor(relatedAmount * rr) = 0` even though the validator IS in the distribution. Guard 5 correctly skips the requeue.

### TestHandleFailedUndelegate_MixedDistributionMultiWDR

Validates correct behavior when an unbonding record references two WDRs: one where the validator IS in the distribution (normal requeue at 299,999,999 due to TruncateInt), and one where it is NOT (skip requeue). Confirms both paths work within a single `HandleFailedUndelegate` call.

### TestHandleFailedUndelegate_DistributionAmountExceedsTotalAmount

Tests the relatedAmount clamp with corrupt state where distribution amount (150) exceeds WDR total amount (100). Verifies that relatedAmount is capped at amount, preventing `relatedQAsset` from exceeding `BurnAmount` and producing negative values after `SubAmount`. Guard 3 (relatedQAsset cap) provides an additional safety net but is unreachable given this prior clamp.

### TestHandleFailedUndelegate_AmountZeroGuard

Tests Guard 2 (division by zero protection) with corrupt state where WDR has wrong denom for Amount (resulting in 0 for BaseDenom). Verifies the guard skips processing to prevent panic.

All 10 tests pass (4 existing + 1 pre-existing missing-WDR + 5 new):

```text
--- PASS: TestKeeperTestSuite/TestHandleFailedUndelegate (1.36s)
    --- PASS: failed_unbond_-_single_wdr,_single_dist (0.50s)
    --- PASS: failed_unbond_-_multi_related_wdr,_single_dist (0.29s)
    --- PASS: failed_unbond_-_multi_related_wdr,_multi_dist (0.17s)
    --- PASS: failed_unbond_-_multi_related_wdr,_multi_dist,_partial_success (0.21s)
--- PASS: TestHandleFailedUndelegate_MissingWithdrawalRecord (0.32s)
--- PASS: TestHandleFailedUndelegate_ValidatorNotInDistribution (0.30s)
--- PASS: TestHandleFailedUndelegate_TinyBurnAmountTruncatesToZero (0.32s)
--- PASS: TestHandleFailedUndelegate_MixedDistributionMultiWDR (0.35s)
--- PASS: TestHandleFailedUndelegate_DistributionAmountExceedsTotalAmount (0.30s)
--- PASS: TestHandleFailedUndelegate_AmountZeroGuard (0.30s)
```

## Post-Upgrade Runbook

1. Deploy patched binary to all validators
2. Coordinate upgrade at agreed halt-height
3. After upgrade, restart hermes relayer
4. Fund osmosis relayer wallet (+0.5 OSMO minimum)
5. Flush acks on cosmoshub, osmosis, sommelier delegate channels
6. Verify channels drain to clean
7. Separately: initiate governance proposals to revive expired counterparty clients on regen, celestia, agoric, juno

## Files Changed

- `x/interchainstaking/keeper/ibc_packet_handlers.go` — 5 guards in `HandleFailedUndelegate` (+25 lines, -4 lines)
- `x/interchainstaking/keeper/ibc_packet_handlers_test.go` — 5 new test functions (+370 lines)
- `app/upgrades/types.go` — add `V0101001UpgradeName`
- `app/upgrades/upgrades.go` — register `NoOpHandler` for v1.10.1
