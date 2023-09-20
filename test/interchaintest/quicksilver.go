package interchaintest

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/stretchr/testify/require"

	query "github.com/cosmos/cosmos-sdk/types/query"
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
)

type ValidatorDelegation struct {
	Address string
	Percent   float64
}

func EncodeValidators(t *testing.T, validators []ValidatorDelegation) string {
	var out []byte
	for _, val := range validators {
		out = append(out, byte(val.Percent/100))

		_, bz, err := bech32.Decode(val.Address)
		require.NoError(t, err)

		converted, err := bech32.ConvertBits(bz, 5, 8, true)
		require.NoError(t, err)
		
		out = append(out, converted...)
	}

	return base64.StdEncoding.EncodeToString(out)
}

type QueryValidatorsResponse struct {
	Validators []Validator `json:"validators"`
	Pagination *query.PageResponse `json:"pagination,omitempty"`
}

type Validator struct {
	// operator_address defines the address of the validator's operator; bech encoded in JSON.
	OperatorAddress string `json:"operator_address,omitempty"`
	// consensus_pubkey is the consensus public key of the validator, as a Protobuf Any.
	ConsensusPubkey *types1.Any `json:"consensus_pubkey,omitempty"`
	// jailed defined whether the validator has been jailed from bonded status or not.
	Jailed bool `json:"jailed,omitempty"`
	// status is the validator status (bonded/unbonding/unbonded).
	Status string `json:"status,omitempty"`
	// tokens define the delegated tokens (incl. self-delegation).
	Tokens github_com_cosmos_cosmos_sdk_types.Int `json:"tokens"`
	// delegator_shares defines total shares issued to a validator's delegators.
	DelegatorShares github_com_cosmos_cosmos_sdk_types.Dec `json:"delegator_shares"`
	// description defines the description terms for the validator.
	Description Description `json:"description"`
	// unbonding_height defines, if unbonding, the height at which this validator has begun unbonding.
	UnbondingHeight int64 `json:"unbonding_height,omitempty"`
	// unbonding_time defines, if unbonding, the min time for the validator to complete unbonding.
	UnbondingTime time.Time `json:"unbonding_time"`
	// commission defines the commission parameters.
	Commission Commission `json:"commission"`
	// min_self_delegation is the validator's self declared minimum self delegation.
	//
	// Since: cosmos-sdk 0.46
	MinSelfDelegation github_com_cosmos_cosmos_sdk_types.Int `json:"min_self_delegation"`
	// strictly positive if this validator's unbonding has been stopped by external modules
	UnbondingOnHoldRefCount int64 `json:"unbonding_on_hold_ref_count,omitempty"`
	// list of unbonding ids, each uniquely identifing an unbonding of this validator
	UnbondingIds []uint64 `json:"unbonding_ids,omitempty"`
}

// Description defines a validator description.
type Description struct {
	// moniker defines a human-readable name for the validator.
	Moniker string `json:"moniker,omitempty"`
	// identity defines an optional identity signature (ex. UPort or Keybase).
	Identity string `json:"identity,omitempty"`
	// website defines an optional website link.
	Website string `json:"website,omitempty"`
	// security_contact defines an optional email for security contact.
	SecurityContact string `json:"security_contact,omitempty"`
	// details define other optional details.
	Details string `json:"details,omitempty"`
}

// CommissionRates defines the initial commission rates to be used for creating
// a validator.
type CommissionRates struct {
	// rate is the commission rate charged to delegators, as a fraction.
	Rate github_com_cosmos_cosmos_sdk_types.Dec `json:"rate"`
	// max_rate defines the maximum commission rate which validator can ever charge, as a fraction.
	MaxRate github_com_cosmos_cosmos_sdk_types.Dec `json:"max_rate"`
	// max_change_rate defines the maximum daily increase of the validator commission, as a fraction.
	MaxChangeRate github_com_cosmos_cosmos_sdk_types.Dec `json:"max_change_rate"`
}

// Commission defines commission parameters for a given validator.
type Commission struct {
	// commission_rates defines the initial commission rates to be used for creating a validator.
	CommissionRates `json:"commission_rates"`
	// update_time is the last time the commission rate was changed.
	UpdateTime time.Time `json:"update_time"`
}