package types

// NewParams creates a new airdrop Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams default ics params.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates params.
func (p *Params) Validate() error {
	return nil
}
