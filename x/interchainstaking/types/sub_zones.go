package types

// IsSubzone returns true if this zone is a sub-zone.
func (z *Zone) IsSubzone() bool {
	return z.SubzoneInfo != nil
}

// ChainID returns the ID of the running chain for the given zone.
func (z *Zone) ChainID() string {
	if z.IsSubzone() {
		return z.SubzoneInfo.BaseChainID
	}

	return z.ChainId
}

// ID returns the unique identifier for the given zone.
func (z *Zone) ID() string {
	return z.ChainId
}
