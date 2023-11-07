package types

const (
	EventTypeRegisterZone      = "register_zone"
	EventTypeRedemptionRequest = "request_redemption"
	EventTypeCloseICA          = "close_ica_channel"
	EventTypeReopenICA         = "reopen_ica_channel"
	EventTypeSetLsmCaps        = "lsm_set_caps"

	AttributeKeyConnectionID     = "connection_id"
	AttributeKeyRecipientChain   = "chain_id"
	AttributeKeyRecipientAddress = "recipient"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyRedeemAmount     = "redeem_amount"
	AttributeKeySourceAddress    = "source"
	AttributeKeyPortID           = "port_name"
	AttributeKeyChannelID        = "channel_id"

	AttributeLsmValidatorCap     = "lsm_validator_cap"
	AttributeLsmValidatorBondCap = "lsm_validator_bond_cap"
	AttributeLsmGlobalCap        = "lsm_global_cap"

	AttributeValueCategory = ModuleName
)
