# ADR 001: Multi Denomination Rewards

## Changelog

- 2023-05-11: Initial Draft (@joe-bowman)

## Status

PROPOSED

## Abstract

With the advent of replicated security, delegators now accrue rewards in multiple
denominations. These collected rewards should be distributed proportionally to the
appropriate qAsset holders for the host zone.

## Context

The current implementation of rewards distribution is such that
rewards collected in the host zone's staking denom are delegated, minus a
parameterized fee that is sent, along with the remaining rewards (any non-staking-
denom fees that have been collected) to the FeeCollectorAccount on Quicksilver to
be distributed to QCK token holders.

With rewards being accrued in multiple denominations, this proposal asserts that
these rewards should no longer be directed to QCK token holder, but to the holders
of the qAssets whose underlying capital generated the rewards.

As per the current behavior with staking denomination rewards, it is proposed that
QCK token holders should receive a cut of fees, using the exist CommissionRate
parameter.

## Alternatives

The proposed behavior, in line with the stated mission to mirror the behavior of
native staking as closely as possible, is the behavior that users are likely to
expect. The alternative approach would be to keep multi-denomination rewards distributed
to QCK token holders.

## Decision

> This section describes our response to these forces. It is stated in full
> sentences, with active voice. "We will ..."
> {decision body}

## Consequences

> This section describes the resulting context, after applying the decision. All
> consequences should be listed here, not just the "positive" ones. A particular
> decision may have positive, negative, and neutral consequences, but all of them
> affect the team and project in the future.

### Backwards Compatibility

> All ADRs that introduce backwards incompatibilities must include a section
> describing these incompatibilities and their severity. The ADR must explain
> how the author proposes to deal with these incompatibilities. ADR submissions
> without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

> {positive consequences}

### Negative

> {negative consequences}

### Neutral

> {neutral consequences}

## Further Discussions

> While an ADR is in the DRAFT or PROPOSED stage, this section should contain a
> summary of issues to be solved in future iterations (usually referencing comments
> from a pull-request discussion).
>
> Later, this section can optionally list ideas or improvements the author or
> reviewers found during the analysis of this ADR.

## Test Cases [optional]

Test cases for an implementation are mandatory for ADRs that are affecting consensus
changes. Other ADRs can choose to include links to test cases if applicable.

## References

- {reference link}
