package types

const (
	EventTypeRegisterZone      = "register_zone"
	EventTypeRedemptionRequest = "request_redemption"
	EventTypeSetIntent         = "set_intent"
	EventTypeCloseICA          = "close_ica_channel"
	EventTypeReopenICA         = "reopen_ica_channel"

	AttributeKeyConnectionID     = "connection_id"
	AttributeKeyChainID          = "chain_id"
	AttributeKeyRecipientAddress = "recipient"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyRedeemAmount     = "redeem_amount"
	AttributeKeySourceAddress    = "source"
	AttributeKeyChannelID        = "channel_id"
	AttributeKeyPortID           = "port_name"
	AttributeKeyUser             = "user_address"

	AttributeValueCategory = ModuleName
)
