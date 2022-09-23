package bindings

// QuickSilverQuery contains quicksilver custom queries.
type QuickSilverQuery struct {
	/// Given a subdenom minted by a contract via `QuickSilverMsg::MintTokens`,
	/// returns the full denom as used by `BankMsg::Send`.
	FullDenom *FullDenom `json:"full_denom,omitempty"`
	/// Return current spot price swapping In for Out on given pool ID.
	/// Warning: this can easily be manipulated via sandwich attacks, do not use as price oracle.
	/// Returns the admin of a denom, if the denom is a Token Factory denom.
	DenomAdmin *DenomAdmin `json:"denom_admin,omitempty"`
}

type FullDenom struct {
	CreatorAddr string `json:"creator_addr"`
	Subdenom    string `json:"subdenom"`
}

type DenomAdmin struct {
	Subdenom string `json:"subdenom"`
}

type DenomAdminResponse struct {
	Admin string `json:"admin"`
}

type FullDenomResponse struct {
	Denom string `json:"denom"`
}
