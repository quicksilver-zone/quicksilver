package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
)

func NewICAAccount(addr, portID string) (*ICAAccount, error) {
	if _, err := addressutils.AccAddressFromBech32(addr, ""); err != nil {
		return nil, err
	}
	return &ICAAccount{Address: addr, WithdrawalAddress: addr, Balance: sdk.Coins{}, PortName: portID}, nil
}

func (a *ICAAccount) SetWithdrawalAddress(addr string) error {
	if _, err := addressutils.AccAddressFromBech32(addr, ""); err != nil {
		return err
	}
	a.WithdrawalAddress = addr
	return nil
}

func (a *ICAAccount) SetBalance(coins sdk.Coins) error {
	if err := coins.Validate(); err != nil {
		return err
	}
	a.Balance = coins
	return nil
}

func (a *ICAAccount) IncrementBalanceWaitgroup() {
	a.BalanceWaitgroup++
}

func (a *ICAAccount) DecrementBalanceWaitgroup() error {
	if a.BalanceWaitgroup <= 0 {
		return errors.New("unable to decrement the balance waitgroup below 0")
	}
	a.BalanceWaitgroup--
	return nil
}
