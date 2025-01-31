package umeetypes

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	partcipationrewardstypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type ParticipationRewardsKeeper interface {
	GetProtocolData(ctx sdk.Context, pdType partcipationrewardstypes.ProtocolDataType, key string) (partcipationrewardstypes.ProtocolData, bool)
}

// ExchangeUToken converts an sdk.Coin containing a uToken to its value in a base
// token.
func ExchangeUToken(ctx sdk.Context, uToken sdk.Coin, prKeeper ParticipationRewardsKeeper) (sdk.Coin, error) {
	if err := uToken.Validate(); err != nil {
		return sdk.Coin{}, err
	}

	tokenDenom := types.ToTokenDenom(uToken.Denom)
	if tokenDenom == "" {
		return sdk.Coin{}, nil
	}

	exchangeRate, err := DeriveExchangeRate(ctx, tokenDenom, prKeeper)
	if err != nil {
		return sdk.Coin{}, err
	}
	tokenAmount := sdk.NewDecFromInt(uToken.Amount).Mul(exchangeRate).TruncateInt()
	return sdk.NewCoin(tokenDenom, tokenAmount), nil
}

// DeriveExchangeRate calculated the token:uToken exchange rate of a base token denom.
func DeriveExchangeRate(ctx sdk.Context, denom string, prKeeper ParticipationRewardsKeeper) (sdk.Dec, error) {
	// Get reserves
	reservesPD, ok := prKeeper.GetProtocolData(ctx, partcipationrewardstypes.ProtocolDataTypeUmeeReserves, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	reservesData, err := partcipationrewardstypes.UnmarshalProtocolData(partcipationrewardstypes.ProtocolDataTypeUmeeReserves, reservesPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	reserves, _ := reservesData.(*partcipationrewardstypes.UmeeReservesProtocolData)

	intamount, err := reserves.GetReserveAmount()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	reserveAmount := sdk.NewDecFromInt(intamount)

	// get leverage module balance
	balancePD, ok := prKeeper.GetProtocolData(ctx, partcipationrewardstypes.ProtocolDataTypeUmeeLeverageModuleBalance, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	balanceData, err := partcipationrewardstypes.UnmarshalProtocolData(partcipationrewardstypes.ProtocolDataTypeUmeeLeverageModuleBalance, balancePD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	balance, _ := balanceData.(*partcipationrewardstypes.UmeeLeverageModuleBalanceProtocolData)

	intamount, err = balance.GetModuleBalance()
	if err != nil {
		return sdk.ZeroDec(), err
	}
	moduleBalance := sdk.NewDecFromInt(intamount)

	// get interest scalar
	interestPD, ok := prKeeper.GetProtocolData(ctx, partcipationrewardstypes.ProtocolDataTypeUmeeInterestScalar, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	interestData, err := partcipationrewardstypes.UnmarshalProtocolData(partcipationrewardstypes.ProtocolDataTypeUmeeInterestScalar, interestPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	interest, _ := interestData.(*partcipationrewardstypes.UmeeInterestScalarProtocolData)
	interestScalar, err := interest.GetInterestScalar()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// get total borrowed
	borrowsPD, ok := prKeeper.GetProtocolData(ctx, partcipationrewardstypes.ProtocolDataTypeUmeeTotalBorrows, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	borrowsData, err := partcipationrewardstypes.UnmarshalProtocolData(partcipationrewardstypes.ProtocolDataTypeUmeeTotalBorrows, borrowsPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	borrows, _ := borrowsData.(*partcipationrewardstypes.UmeeTotalBorrowsProtocolData)
	borrowAmount, err := borrows.GetTotalBorrows()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	totalBorrowed := borrowAmount.Mul(interestScalar)

	// get UToken supply
	uTokenPD, ok := prKeeper.GetProtocolData(ctx, partcipationrewardstypes.ProtocolDataTypeUmeeUTokenSupply, types.ToUTokenDenom(denom))
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	uTokenData, err := partcipationrewardstypes.UnmarshalProtocolData(partcipationrewardstypes.ProtocolDataTypeUmeeUTokenSupply, uTokenPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	utokens, _ := uTokenData.(*partcipationrewardstypes.UmeeUTokenSupplyProtocolData)
	uTokenSupply, err := utokens.GetUTokenSupply()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// Derive effective token supply
	tokenSupply := moduleBalance.Add(totalBorrowed).Sub(reserveAmount)

	// Handle uToken supply == 0 case
	if !uTokenSupply.IsPositive() {
		return sdk.OneDec(), nil
	}

	// Derive exchange rate
	return tokenSupply.QuoInt(uTokenSupply), nil
}
