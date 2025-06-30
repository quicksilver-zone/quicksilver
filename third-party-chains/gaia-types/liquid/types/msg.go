package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	address "github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

func (msg *MsgUpdateParams) ValidateBasic() error {
	return nil
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.Authority, "")}
}

func (msg *MsgTokenizeShares) ValidateBasic() error {
	return nil
}

func (msg *MsgTokenizeShares) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.TokenizedShareOwner, "")}
}

func (msg *MsgRedeemTokensForShares) ValidateBasic() error {
	return nil
}

func (msg *MsgRedeemTokensForShares) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.DelegatorAddress, "")}
}

func (msg *MsgTransferTokenizeShareRecord) ValidateBasic() error {
	return nil
}

func (msg *MsgTransferTokenizeShareRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.Sender, "")}
}

func (msg *MsgDisableTokenizeShares) ValidateBasic() error {
	return nil
}

func (msg *MsgDisableTokenizeShares) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.DelegatorAddress, "")}
}

func (msg *MsgEnableTokenizeShares) ValidateBasic() error {
	return nil
}

func (msg *MsgEnableTokenizeShares) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.DelegatorAddress, "")}
}

func (msg *MsgWithdrawTokenizeShareRecordReward) ValidateBasic() error {
	return nil
}

func (msg *MsgWithdrawTokenizeShareRecordReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.OwnerAddress, "")}
}

func (msg *MsgWithdrawAllTokenizeShareRecordReward) ValidateBasic() error {
	return nil
}

func (msg *MsgWithdrawAllTokenizeShareRecordReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{address.MustAccAddressFromBech32(msg.OwnerAddress, "")}
}
