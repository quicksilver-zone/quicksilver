package types

// NewParams creates a new claimsmanager Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams default claimsmanager params.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates params.
func (p *Params) Validate() error {
	return nil
}
