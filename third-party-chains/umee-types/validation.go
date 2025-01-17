package umeetypes

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	claimsmanagertypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
)

type ClaimsManagerKeeper interface {
	GetProtocolData(ctx sdk.Context, pdType claimsmanagertypes.ProtocolDataType, key string) (claimsmanagertypes.ProtocolData, bool)
}

// ExchangeUToken converts an sdk.Coin containing a uToken to its value in a base
// token.
func ExchangeUToken(ctx sdk.Context, uToken sdk.Coin, prKeeper ClaimsManagerKeeper) (sdk.Coin, error) {
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
func DeriveExchangeRate(ctx sdk.Context, denom string, prKeeper ClaimsManagerKeeper) (sdk.Dec, error) {
	// Get reserves
	reservesPD, ok := prKeeper.GetProtocolData(ctx, claimsmanagertypes.ProtocolDataTypeUmeeReserves, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	reservesData, err := claimsmanagertypes.UnmarshalProtocolData(claimsmanagertypes.ProtocolDataTypeUmeeReserves, reservesPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	reserves, _ := reservesData.(*claimsmanagertypes.UmeeReservesProtocolData)

	intamount, err := reserves.GetReserveAmount()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	reserveAmount := sdk.NewDecFromInt(intamount)

	// get leverage module balance
	balancePD, ok := prKeeper.GetProtocolData(ctx, claimsmanagertypes.ProtocolDataTypeUmeeLeverageModuleBalance, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	balanceData, err := claimsmanagertypes.UnmarshalProtocolData(claimsmanagertypes.ProtocolDataTypeUmeeLeverageModuleBalance, balancePD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	balance, _ := balanceData.(*claimsmanagertypes.UmeeLeverageModuleBalanceProtocolData)

	intamount, err = balance.GetModuleBalance()
	if err != nil {
		return sdk.ZeroDec(), err
	}
	moduleBalance := sdk.NewDecFromInt(intamount)

	// get interest scalar
	interestPD, ok := prKeeper.GetProtocolData(ctx, claimsmanagertypes.ProtocolDataTypeUmeeInterestScalar, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	interestData, err := claimsmanagertypes.UnmarshalProtocolData(claimsmanagertypes.ProtocolDataTypeUmeeInterestScalar, interestPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	interest, _ := interestData.(*claimsmanagertypes.UmeeInterestScalarProtocolData)
	interestScalar, err := interest.GetInterestScalar()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// get total borrowed
	borrowsPD, ok := prKeeper.GetProtocolData(ctx, claimsmanagertypes.ProtocolDataTypeUmeeTotalBorrows, denom)
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	borrowsData, err := claimsmanagertypes.UnmarshalProtocolData(claimsmanagertypes.ProtocolDataTypeUmeeTotalBorrows, borrowsPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	borrows, _ := borrowsData.(*claimsmanagertypes.UmeeTotalBorrowsProtocolData)
	borrowAmount, err := borrows.GetTotalBorrows()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	totalBorrowed := borrowAmount.Mul(interestScalar)

	// get UToken supply
	uTokenPD, ok := prKeeper.GetProtocolData(ctx, claimsmanagertypes.ProtocolDataTypeUmeeUTokenSupply, types.ToUTokenDenom(denom))
	if !ok {
		return sdk.ZeroDec(), fmt.Errorf("unable to obtain protocol data for denom=%s", denom)
	}
	uTokenData, err := claimsmanagertypes.UnmarshalProtocolData(claimsmanagertypes.ProtocolDataTypeUmeeUTokenSupply, uTokenPD.Data)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	utokens, _ := uTokenData.(*claimsmanagertypes.UmeeUTokenSupplyProtocolData)
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
