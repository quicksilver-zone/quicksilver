package types

const (
	EventTypeRegisterZone      = "register_zone"
	EventTypeRedemptionRequest = "request_redemption"
	EventTypeSetIntent         = "set_intent"
	EventTypeCloseICA          = "close_ica_channel"
	EventTypeReopenICA         = "reopen_ica_channel"

	AttributeKeyConnectionID     = "connection_id"
	AttributeKeyChainId          = "chain_id"
	AttributeKeyRecipientAddress = "recipient"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyRedeemAmount     = "redeem_amount"
	AttributeKeySourceAddress    = "source"
	AttributeKeyChannelId        = "channel_id"
	AttributeKeyPortId           = "port_name"
	AttributeKeyUser             = "user_address"

	AttributeValueCategory = ModuleName
)
