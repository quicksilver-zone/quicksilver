package types

const (
	EventTypeRegisterZone      = "register_zone"
	EventTypeRedemptionRequest = "request_redemption"
<<<<<<< HEAD
=======
	EventTypeSetIntent         = "set_intent"
>>>>>>> origin/develop
	EventTypeCloseICA          = "close_ica_channel"
	EventTypeReopenICA         = "reopen_ica_channel"

	AttributeKeyConnectionID     = "connection_id"
	AttributeKeyChainID          = "chain_id"
	AttributeKeyRecipientAddress = "recipient"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyRedeemAmount     = "redeem_amount"
	AttributeKeySourceAddress    = "source"
<<<<<<< HEAD
	AttributeKeyPortID           = "port_name"
	AttributeKeyChannelID        = "channel_id"
=======
	AttributeKeyChannelID        = "channel_id"
	AttributeKeyPortID           = "port_name"
	AttributeKeyUser             = "user_address"
>>>>>>> origin/develop

	AttributeValueCategory = ModuleName
)
