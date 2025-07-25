package types

import (
	"go.uber.org/multierr"

	"github.com/quicksilver-zone/quicksilver/utils"
)

// ValidateBasic performs stateless validation for Proof.
func (p *Proof) ValidateBasic() error {
	errs := make(map[string]error)

	if len(p.Key) == 0 {
		errs["Key"] = ErrUndefinedAttribute
	}

	if len(p.Data) == 0 {
		errs["Data"] = ErrUndefinedAttribute
	}

	if p.ProofOps == nil {
		errs["ProofOps"] = ErrUndefinedAttribute
	}

	if p.Height < 0 {
		errs["Height"] = ErrNegativeAttribute
	}

	if p.ProofType == "" {
		errs["ProofType"] = ErrUndefinedAttribute
	}

	// check for errors and return
	if len(errs) > 0 {
		return multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return nil
}
