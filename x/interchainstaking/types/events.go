package types

const (
	EventTypeRegisterZone            = "register_zone"
	EventTypeRedemptionRequest       = "request_redemption"
	EventTypeRedemptionCancellation  = "cancel_redemption"
	EventTypeRedemptionRequeue       = "requeue_redemption"
	EventTypeSetIntent               = "set_intent"
	EventTypeCloseICA                = "close_ica_channel"
	EventTypeReopenICA               = "reopen_ica_channel"
	EventTypeSetLsmCaps              = "lsm_set_caps"
	EventTypeAddValidatorDenyList    = "add_validator_deny_list"
	EventTypeRemoveValidatorDenyList = "remove_validator_deny_list"

	AttributeKeyConnectionID     = "connection_id"
	AttributeKeyChainID          = "chain_id"
	AttributeKeyRecipientAddress = "recipient"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyReturnedAmount   = "returned_amount"
	AttributeKeyRedeemAmount     = "redeem_amount"
	AttributeKeySourceAddress    = "source"
	AttributeKeyChannelID        = "channel_id"
	AttributeKeyPortID           = "port_name"
	AttributeKeyUser             = "user_address"

	AttributeLsmValidatorCap     = "lsm_validator_cap"
	AttributeLsmValidatorBondCap = "lsm_validator_bond_cap"
	AttributeLsmGlobalCap        = "lsm_global_cap"

	AttributeValueCategory = ModuleName

	AttributeKeyOperatorAddress = "operator_address"
)
