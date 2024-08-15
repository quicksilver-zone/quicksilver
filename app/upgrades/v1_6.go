package upgrades

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v6migration "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/migrations/v6"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// ============ TESTNET UPGRADE HANDLERS ============

func V010600beta1UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) {
			appKeepers.UpgradeKeeper.Logger(ctx).Info("removing defunct zones")
			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "elgafar-1")
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010600beta0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) {
			appKeepers.UpgradeKeeper.Logger(ctx).Info("migrating capabilities")
			err := v6migration.MigrateICS27ChannelCapability(
				ctx,
				appKeepers.IBCKeeper.Codec(),
				appKeepers.GetKey(capabilitytypes.StoreKey),
				appKeepers.CapabilityKeeper,
				icstypes.ModuleName,
			)
			if err != nil {
				panic(err)
			}

			appKeepers.UpgradeKeeper.Logger(ctx).Info("removing defunct zones")
			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "agoric-3")
			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "archway-1")
			appKeepers.InterchainQueryKeeper.SetLatestHeight(ctx, "provider", 6209948)
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010600rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// no action yet.
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// update block params to 2MB and 150m gas.
func updateBlockParams(ctx sdk.Context, appKeepers *keepers.AppKeepers) {
	ctx.Logger().Info("Updating setting block params; 2MB max_bytes, 150M max_gas")
	ss, found := appKeepers.ParamsKeeper.GetSubspace(baseapp.Paramspace)
	if !found {
		panic("params subspace not found")
	}
	ss.Set(ctx, baseapp.ParamStoreKeyBlockParams, abci.BlockParams{
		MaxBytes: 2072576,
		MaxGas:   150000000,
	})
}

func enableIcaHost(ctx sdk.Context, appKeepers *keepers.AppKeepers) {
	ctx.Logger().Info("Enabling ICAHost")
	appKeepers.ICAHostKeeper.SetParams(ctx, icahosttypes.Params{
		HostEnabled: true,
		AllowMessages: []string{
			"/cosmos.bank.v1beta1.MsgSend",
			"/cosmos.bank.v1beta1.MsgMultiSend",
			"/quicksilver.interchainstaking.v1.MsgSignalIntent",
			"/quicksilver.interchainstaking.v1.MsgRequestRedemption",
			"/quicksilver.participationrewards.v1.MsgSubmitClaim",
			"/cosmos.authz.v1beta1.MsgGrant",
			"/cosmos.authz.v1beta1.MsgRevoke",
			"/ibc.applications.transfer.v1.MsgTransfer",
		},
	})
}

func setTransferChannels(ctx sdk.Context, appKeepers *keepers.AppKeepers, channels map[string]string) {
	ctx.Logger().Info("Set TransferChannel field for zones")
	appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		zone.TransferChannel = channels[zone.ChainId]
		appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
		return false
	})

	appKeepers.ParticipationRewardsKeeper.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeConnection), func(idx int64, _ []byte, data types.ProtocolData) bool {
		pd, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, data.Data)
		if err != nil {
			panic(err)
		}
		pdc, ok := pd.(*types.ConnectionProtocolData)
		if ok {
			tc, ok := channels[pdc.ChainID]
			if ok {
				pdc.TransferChannel = tc

				pdcBlob, err := json.Marshal(pdc)
				if err != nil {
					appKeepers.ParticipationRewardsKeeper.Logger(ctx).Info("Error Marshalling self connection Data")
					panic(err)
				}

				data := types.ProtocolData{
					Type: types.ProtocolDataType_name[int32(types.ProtocolDataTypeConnection)],
					Data: pdcBlob,
				}

				appKeepers.ParticipationRewardsKeeper.SetProtocolData(ctx, pdc.GenerateKey(), &data)
			}
		}
		return false
	})
}

func removeIncorrectLiquidTokenProtocolDatas(ctx sdk.Context, appKeepers *keepers.AppKeepers, channels map[string]string) {
	ctx.Logger().Info("Removing incorrect IBC denom for LiquidAllowedDenomProtocolData")
	appKeepers.ParticipationRewardsKeeper.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeLiquidToken), func(idx int64, _ []byte, data types.ProtocolData) bool {
		pd, err := types.UnmarshalProtocolData(types.ProtocolDataTypeLiquidToken, data.Data)
		if err != nil {
			return false
		}
		token, _ := pd.(*types.LiquidAllowedDenomProtocolData)

		if token.ChainID == ctx.ChainID() {
			return false
		}

		channel, found := appKeepers.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", channels[token.ChainID])
		if !found {
			panic(fmt.Errorf("unable to find channel %s", channels[token.ChainID]))
		}

		// derive the correct ibc denom; if it does not match then remmove it.
		correctIbc := utils.DeriveIbcDenom("transfer", channel.Counterparty.ChannelId, "transfer", channels[token.ChainID], token.QAssetDenom)
		if token.IbcDenom != correctIbc {
			ctx.Logger().Info(fmt.Sprintf("incorrect IBC denom %s for LiquidAllowedDenomProtocolData %s, expected %s. removing", token.IbcDenom, token.QAssetDenom, correctIbc))
			appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, types.GetProtocolDataKey(types.ProtocolDataTypeLiquidToken, pd.GenerateKey()))
		}
		return false
	})
}

func V010601rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		updateBlockParams(ctx, appKeepers)

		enableIcaHost(ctx, appKeepers)

		channels := map[string]string{
			"osmo-test-5":     "channel-39",
			"provider":        "channel-0",
			"elgafar-1":       "channel-1",
			"regen-redwood-1": "channel-2",
		}

		setTransferChannels(ctx, appKeepers, channels)

		removeIncorrectLiquidTokenProtocolDatas(ctx, appKeepers, channels)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// =========== PRODUCTION UPGRADE HANDLER ===========

func V010601UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isMainnet(ctx) || isTest(ctx) {
			appKeepers.UpgradeKeeper.Logger(ctx).Info("migrating capabilities")
			err := v6migration.MigrateICS27ChannelCapability(
				ctx,
				appKeepers.IBCKeeper.Codec(),
				appKeepers.GetKey(capabilitytypes.StoreKey),
				appKeepers.CapabilityKeeper,
				icstypes.ModuleName,
			)
			if err != nil {
				panic(err)
			}

			updateBlockParams(ctx, appKeepers)

			ctx.Logger().Info("Updating agoric-3 zone to set is_118 = false")
			agoricZone, _ := appKeepers.InterchainstakingKeeper.GetZone(ctx, "agoric-3")
			agoricZone.Is_118 = false
			appKeepers.InterchainstakingKeeper.SetZone(ctx, &agoricZone)

			enableIcaHost(ctx, appKeepers)

			channels := map[string]string{
				"osmosis-1":      "channel-2",
				"cosmoshub-4":    "channel-1",
				"stargaze-1":     "channel-0",
				"juno-1":         "channel-86",
				"sommelier-3":    "channel-101",
				"regen-1":        "channel-17",
				"umee-1":         "channel-49",
				"secret-4":       "channel-52",
				"dydx-mainnet-1": "channel-164",
				"agoric-3":       "channel-125",
				"ssc-1":          "channel-170",
			}

			setTransferChannels(ctx, appKeepers, channels)

			removeIncorrectLiquidTokenProtocolDatas(ctx, appKeepers, channels)

			appKeepers.InterchainstakingKeeper.Logger(ctx).Info("setting 5 unsent unbondings for epoch 172 to STATUS_UNBONDING to be picked up by the next end blocker...")
			// remit epoch 172 unbondings that did not send due to channel timeout closure.
			hashes := []string{
				"0000000000000000000000000000000000000000000000000000000000000560", // 1137.426431
				"068d2733ac95552fa0adc35cb1ae48dfaafa49ecb86a2fab0240f2fb600ff96e", // 0.186384
				"0c23dddc51cf16671c0b95bd206980825d115baa2efa985184f543e708171f67", // 25.919839
				"727428acc32d07e2754a2af6eafc949bea6c85de95cba9909e457d271d5dbb83", // 5.043558
				"8c4d9e582303597968111fe06b4acd5fb8dc5ebbeaa1b0a0fc7fa7fd482b1a85", // 303.404702
			}
			for _, hash := range hashes {
				record, found := appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "cosmoshub-4", hash, icstypes.WithdrawStatusSend)
				if !found {
					// do not panic, in case records were updated on previous epoch.
					appKeepers.InterchainstakingKeeper.Logger(ctx).Error(fmt.Sprintf("1: unable to find record for hash %s", hash))
					continue
				}

				// update the record so that it will re-trigger the send.
				appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &record, icstypes.WithdrawStatusUnbond)
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("updated record to STATUS_UNBONDING", "hash", hash)
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010603UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isMainnet(ctx) || isTest(ctx) {
			appKeepers.InterchainstakingKeeper.Logger(ctx).Info("returning escrowed funds for 1 unbonding that failed to send to original user...")

			hashes := []string{
				"ea0d86a3fb4b25fcb13a587e72542f99ebf8c7c3aa255a0922dfa7002a8ee861",
			}
			for _, hash := range hashes {
				record, found := appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "cosmoshub-4", hash, icstypes.WithdrawStatusUnbond)
				if !found {
					panic(fmt.Sprintf("unable to find record for hash %s", hash))
				}

				// move funds from escrowed funds to the original user's account.
				err := appKeepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, addressutils.MustAccAddressFromBech32(record.Delegator, "quick"), sdk.NewCoins(record.BurnAmount))
				if err != nil {
					panic(fmt.Sprintf("unable to send coins from module to account for record %s", hash))
				}

				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("sent funds to original user", "user", record.Delegator, "amount", record.BurnAmount)

				// delete the record so that it won't re-trigger.
				appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, "cosmoshub-4", hash, icstypes.WithdrawStatusUnbond)
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("deleted record", "hash", hash)
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
