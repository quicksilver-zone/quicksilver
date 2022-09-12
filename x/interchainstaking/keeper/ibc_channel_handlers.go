package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k Keeper) HandleChannelOpenAck(ctx sdk.Context, portID string, connectionID string) error {
	chainID, err := k.GetChainID(ctx, connectionID)
	if err != nil {
		ctx.Logger().Error(
			"unable to obtain chain for given connection and port",
			"connectionID", connectionID,
			"portID", portID,
			"error", err,
		)
		return fmt.Errorf("unable to obtain chain for %s/%s: %w", connectionID, portID, err)
	}

	// get zone
	zone, found := k.GetZone(ctx, chainID)
	if !found {
		err := fmt.Errorf("unable to obtain zone for chainID %s", chainID)
		ctx.Logger().Error(err.Error())
		return err
	}

	// get interchain account address
	address, found := k.ICAControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		err := fmt.Errorf("expected to find an address for %s/%s", connectionID, portID)
		ctx.Logger().Error(err.Error())
		return err
	}

	ctx.Logger().Info("found matching address", "chain", zone.ChainId, "address", address, "port", portID)
	portParts := strings.Split(portID, ".")

	switch {
	// deposit address
	case len(portParts) == 2 && portParts[1] == types.ICASuffixDeposit:

		zone.DepositAddress, err = types.NewICAAccount(address, portID, zone.BaseDenom)
		if err != nil {
			return err
		}

		balanceQuery := bankTypes.QueryAllBalancesRequest{Address: address}
		bz, err := k.GetCodec().Marshal(&balanceQuery)
		if err != nil {
			return err
		}

		k.ICQKeeper.MakeRequest(
			ctx,
			connectionID,
			chainID,
			"cosmos.bank.v1beta1.Query/AllBalances",
			bz,
			sdk.NewInt(int64(k.GetParam(ctx, types.KeyDepositInterval))),
			types.ModuleName,
			"allbalances",
			0,
		)

	// withdrawal address
	case len(portParts) == 2 && portParts[1] == types.ICASuffixWithdrawal:
		zone.WithdrawalAddress, err = types.NewICAAccount(address, portID, zone.BaseDenom)
		if err != nil {
			return err
		}

	// delegation addresses
	case len(portParts) == 2 && portParts[1] == types.ICASuffixDelegate:
		zone.DelegationAddress, err = types.NewICAAccount(address, portID, zone.BaseDenom)
		if err != nil {
			return err
		}

	// performance address
	case len(portParts) == 2 && portParts[1] == types.ICASuffixPerformance:
		zone.PerformanceAddress, err = types.NewICAAccount(address, portID, zone.BaseDenom)
		if err != nil {
			return err
		}

		if err := k.EmitPerformanceBalanceQuery(ctx, &zone); err != nil {
			k.Logger(ctx).Error("error emitting performance balance query", "error", err)
			return err
		}

	default:
		err := fmt.Errorf("unexpected channel on portID: %s", portID)
		ctx.Logger().Error(err.Error())
		return err
	}
	k.SetZone(ctx, &zone)
	return nil
}
