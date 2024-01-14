package keeper

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) HandleChannelOpenAck(ctx sdk.Context, portID, connectionID string) error {
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

		if zone.DepositAddress == nil {
			zone.DepositAddress, err = types.NewICAAccount(address, portID)
			if err != nil {
				return err
			}

			k.SetAddressZoneMapping(ctx, address, zone.ChainId)

			balanceQuery := banktypes.QueryAllBalancesRequest{Address: address}
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
				sdkmath.NewInt(int64(k.GetParam(ctx, types.KeyDepositInterval))),
				types.ModuleName,
				"allbalances",
				0,
			)
		}

	// withdrawal address
	case len(portParts) == 2 && portParts[1] == types.ICASuffixWithdrawal:
		if zone.WithdrawalAddress == nil {
			zone.WithdrawalAddress, err = types.NewICAAccount(address, portID)
			if err != nil {
				return err
			}
			k.SetAddressZoneMapping(ctx, address, zone.ChainId)
		}

	// delegation addresses
	case len(portParts) == 2 && portParts[1] == types.ICASuffixDelegate:
		if zone.DelegationAddress == nil {
			zone.DelegationAddress, err = types.NewICAAccount(address, portID)
			if err != nil {
				return err
			}

			k.SetAddressZoneMapping(ctx, address, zone.ChainId)
		}

	// performance address
	case len(portParts) == 2 && portParts[1] == types.ICASuffixPerformance:
		if zone.PerformanceAddress == nil {
			ctx.Logger().Info("create performance account")
			zone.PerformanceAddress, err = types.NewICAAccount(address, portID)
			if err != nil {
				return err
			}
			k.SetAddressZoneMapping(ctx, address, zone.ChainId)
		}

		// emit this periodic query the first time, but not subsequently.
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
