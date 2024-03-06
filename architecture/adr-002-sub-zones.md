ADR 002: Sub-Zones
Changelog
2023-06-02: Initial Draft (@aljo242)
Status
DRAFT

Abstract

Sub zones on Quicksilver refers to registering a chain on Quicksilver with a tailor-made validator subset and distribution for an enterprise use case.
The primary benefit of sub zones is for users to be able to issue their own qAsset with their own risk parameters and validator distribution.
These assets also know as ‘sub assets’ are not fungible with qAssets in its current form, though with the LSM in place it becomes possible for these assets to become practically fungible.

Context

The current implementation of a `Zone` is:

```proto
message Zone {
  string connection_id = 1;
  string chain_id = 2;
  ICAAccount deposit_address = 3;
  ICAAccount withdrawal_address = 4;
  ICAAccount performance_address = 5;
  ICAAccount delegation_address = 6;
  string account_prefix = 7;
  string local_denom = 8;
  string base_denom = 9;
  string redemption_rate = 10 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string last_redemption_rate = 11 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  repeated Validator validators = 12;
  repeated ValidatorIntent aggregate_intent = 13 [(gogoproto.castrepeated) = "ValidatorIntents"];
  bool multi_send = 14; // deprecated
  bool liquidity_module = 15;
  uint32 withdrawal_waitgroup = 16;
  bytes ibc_next_validators_hash = 17;
  uint64 validator_selection_allocation = 18;
  uint64 holdings_allocation = 19;
  int64 last_epoch_height = 20; // deprecated
  string tvl = 21 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  int64 unbonding_period = 22;
  int64 messages_per_tx = 23;
  int64 decimals = 24;
  bool unbonding_enabled = 25;
  bool deposits_enabled = 26;
  bool return_to_sender = 27;
  bool is_118 = 28;
}
```

Alternatives

Sub-zones could also be implemented as an array of "subzoneID-authority" tuples that are stored themselves on the parent zone.
This would mean that the existing member functions of a `Zone` could not be directly accessed by a sub zone, making things a bit more complicated in the implementation.
This alternative would be more clearer that a sub-zone is actually not a new "zone" meaning chain.

Decision

We will add the following fields to the `Zone` struct:

```proto
message Zone {
  string connection_id = 1;
  string chain_id = 2;
  ICAAccount deposit_address = 3;
  ICAAccount withdrawal_address = 4;
  ICAAccount performance_address = 5;
  ICAAccount delegation_address = 6;
  string account_prefix = 7;
  string local_denom = 8;
  string base_denom = 9;
  string redemption_rate = 10 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string last_redemption_rate = 11 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  repeated Validator validators = 12;
  repeated ValidatorIntent aggregate_intent = 13 [(gogoproto.castrepeated) = "ValidatorIntents"];
  bool multi_send = 14; // deprecated
  bool liquidity_module = 15;
  uint32 withdrawal_waitgroup = 16;
  bytes ibc_next_validators_hash = 17;
  uint64 validator_selection_allocation = 18;
  uint64 holdings_allocation = 19;
  int64 last_epoch_height = 20; // deprecated
  string tvl = 21 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  int64 unbonding_period = 22;
  int64 messages_per_tx = 23;
  int64 decimals = 24;
  bool unbonding_enabled = 25;
  bool deposits_enabled = 26;
  bool return_to_sender = 27;
  bool is_118 = 28;
  SubzoneInfo subzoneInfo = 29;
}

message SubzoneInfo {
      string authority = 1;
      string base_chainID = 2;
}
```
The `SubzoneInfo` `authority` field is the whilelisted Quicksilver account which controls this subzone.  The `BaseChainID`
field is a reference to the "base" or "parent" chain that this zone is a sub-zone of.  For example, if a base zone for
the Cosmos Hub with chainID "gaia-5" exists, an `authority` could propose to create a new `sub-zone` with "gaia-5" as its
base zone.

If `SubzoneInfo` and is non-empty, then the zone is a subzone.  We can add the following helper function:

```go
func (z *Zone) IsSubzone() bool {
    return z.SubZoneInfo != nil
}
```

Subzones now effectively have two "chainIDs": the unique ID created by the `authority` when the zone is created, and the
`BaseChainID` which refers to the base zone.  We can create helper functions to simplify translation:

```go
func (z *Zone) ChainID() string {
    if z.IsSubzone() {
        return z.SubzoneInfo.BaseChainID
    }

    return z.ChainId}

func (z *Zone) ID() string {
    return z.ChainId
}
```

The `zone.ChainID()` function will always return the chainID of the running Cosmos SDK chain that this `zone` is representing.
The `zone.ID()` function will always return the unique identifier for this `zone`.  These helper functions should be used in place of
all direct accesses to the `zone.ChainId` variable.

Consequences
This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

Backwards Compatibility
All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

Positive
{positive consequences}

Negative
{negative consequences}

Neutral
{neutral consequences}

Further Discussions
While an ADR is in the DRAFT or PROPOSED stage, this section should contain a summary of issues to be solved in future iterations (usually referencing comments from a pull-request discussion).

Later, this section can optionally list ideas or improvements the author or reviewers found during the analysis of this ADR.

Test Cases [optional]
Test cases for an implementation are mandatory for ADRs that are affecting consensus changes. Other ADRs can choose to include links to test cases if applicable.

References
{reference link}