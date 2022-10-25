package types

import "github.com/ingenuity-build/quicksilver/internal/multierror"

func (p Proof) ValidateBasic() error {
	errors := make(map[string]error)

	if len(p.Key) == 0 {
		errors["Key"] = ErrUndefinedAttribute
	}

	if len(p.Data) == 0 {
		errors["Data"] = ErrUndefinedAttribute
	}

	if p.ProofOps == nil {
		errors["ProofOps"] = ErrUndefinedAttribute
	}

	if p.Height < 0 {
		errors["Height"] = ErrNegativeAttribute
	}

	if len(p.ProofType) == 0 {
		errors["ProofType"] = ErrUndefinedAttribute
	}

	// check for errors and return
	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}
