package types

const (
	EventTypeRegisterZone      = "register_zone"
	EventTypeRedemptionRequest = "request_redemption"
	EventTypeCloseICA          = "close_ica_channel"
	EventTypeReopenICA         = "reopen_ica_channel"

	AttributeKeyConnectionID     = "connection_id"
	AttributeKeyRecipientChain   = "chain_id"
	AttributeKeyRecipientAddress = "recipient"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyRedeemAmount     = "redeem_amount"
	AttributeKeySourceAddress    = "source"
	AttributeKeyPortID           = "port_name"
	AttributeKeyChannelID        = "channel_id"

	AttributeValueCategory = ModuleName
)
