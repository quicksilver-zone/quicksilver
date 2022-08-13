package types

const (
	EventTypeRegisterZone      = "register_zone"
	EventTypeRedemptionRequest = "request_redemption"

	AttributeKeyConnectionID     = "connection_id"
	AttributeKeyRecipientChain   = "chain_id"
	AttributeKeyRecipientAddress = "recipient"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyRedeemAmount     = "redeem_amount"
	AttributeKeySourceAddress    = "source"

	AttributeValueCategory = ModuleName
)
