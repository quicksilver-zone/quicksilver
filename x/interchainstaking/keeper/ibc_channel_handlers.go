package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k Keeper) HandleChannelOpenAck(ctx sdk.Context, portID string, channelID string, connectionID string) error {
	chainID, err := k.GetChainID(ctx, connectionID)
	if err != nil {
		ctx.Logger().Error(
			"Unable to obtain chain for given connection and port",
			"ConnectionID", connectionID,
			"PortID", portID,
			"Error", err,
		)
		return nil
	}

	// get zone
	zone, found := k.GetZone(ctx, chainID)
	if !found {
		ctx.Logger().Error(fmt.Sprintf("expected to find zone info for %v", chainID))
		return fmt.Errorf("unable to find zone for chainID: %s", chainID)
	}

	// get interchain account address
	address, found := k.ICAControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		ctx.Logger().Error(fmt.Sprintf("expected to find an address for %s/%s", connectionID, portID))
		return nil
	}

	ctx.Logger().Info("Found matching address", "chain", zone.ChainId, "address", address, "port", portID)
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
	case len(portParts) == 3 && portParts[1] == types.ICASuffixDelegate:
		delegationAccounts := zone.GetDelegationAccounts()
		// check for duplicate address
		for _, existing := range delegationAccounts {
			if existing.Address == address {
				ctx.Logger().Error("unexpectedly found existing address: " + address)
				return nil
			}
		}
		account, err := types.NewICAAccount(address, portID, zone.BaseDenom)
		if err != nil {
			return err
		}

		// append delegation account address
		//nolint:gocritic
		zone.DelegationAddresses = append(delegationAccounts, account)

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
		ctx.Logger().Error("unexpected channel on portID: " + portID)
		return fmt.Errorf("unexpected channel on portID %s", portID)
	}
	k.SetZone(ctx, &zone)
	return nil
}
