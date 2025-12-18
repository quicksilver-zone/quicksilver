package types

const (
	EventTypeRegisterZone                = "register_zone"
	EventTypeRedemptionRequest           = "request_redemption"
	EventTypeRedemptionCancellation      = "cancel_redemption"
	EventTypeRedemptionRequeue           = "requeue_redemption"
	EventTypeUpdateRedemption            = "update_redemption"
	EventTypeSetIntent                   = "set_intent"
	EventTypeCloseICA                    = "close_ica_channel"
	EventTypeReopenICA                   = "reopen_ica_channel"
	EventTypeSetLsmCaps                  = "lsm_set_caps"
	EventTypeAddValidatorDenyList        = "add_validator_deny_list"
	EventTypeRemoveValidatorDenyList     = "remove_validator_deny_list"
	EventTypeSetZoneOffboarding          = "set_zone_offboarding"
	EventTypeCancelAllPendingRedemptions = "cancel_all_pending_redemptions"
	EventTypeForceUnbondAllDelegations   = "force_unbond_all_delegations"
	EventTypeOffboardingUnbondAck        = "offboarding_unbond_ack"

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
	AttributeKeyHash             = "hash"
	AttributeKeyNewStatus        = "new_status"

	AttributeLsmValidatorCap     = "lsm_validator_cap"
	AttributeLsmValidatorBondCap = "lsm_validator_bond_cap"
	AttributeLsmGlobalCap        = "lsm_global_cap"

	AttributeValueCategory = ModuleName

	AttributeKeyOperatorAddress = "operator_address"
	AttributeKeyIsOffboarding   = "is_offboarding"
	AttributeKeyCancelledCount  = "cancelled_count"
	AttributeKeyRefundedAmounts = "refunded_amounts"
	AttributeKeyUnbondingCount  = "unbonding_count"
	AttributeKeyTotalUnbonded   = "total_unbonded"
	AttributeKeyValidator       = "validator"
	AttributeKeyAmount          = "amount"
	AttributeKeyCompletionTime  = "completion_time"
)
