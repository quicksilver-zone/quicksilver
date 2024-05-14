package upgrades

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	abci "github.com/tendermint/tendermint/abci/types"

	v6migration "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/migrations/v6"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils"
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

func V010601rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("updating setting block params")
		appKeepers.ParamsKeeper.
			Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable()).
			Set(ctx, baseapp.ParamStoreKeyBlockParams, abci.BlockParams{
				MaxBytes: 2072576,
				MaxGas:   150000000,
			})

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

		channels := map[string]string{
			"osmo-test-5":     "channel-39",
			"provider":        "channel-0",
			"elgafar-1":       "channel-1",
			"regen-redwood-1": "channel-2",
		}

		ctx.Logger().Info("Set TransferChannel field for zones")
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			zone.TransferChannel = channels[zone.ChainId]
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

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
				appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, pd.GenerateKey())
			}
			return false
		})
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
		ctx.Logger().Info("Updating setting block params; 2MB max_bytes, 150M max_gas")
		appKeepers.ParamsKeeper.
			Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable()).
			Set(ctx, baseapp.ParamStoreKeyBlockParams, abci.BlockParams{
				MaxBytes: 2072576,
				MaxGas:   150000000,
			})

		ctx.Logger().Info("Updating agoric-3 zone to set is_118 = false")
		agoricZone, _ := appKeepers.InterchainstakingKeeper.GetZone(ctx, "agoric-3")
		agoricZone.Is_118 = false
		appKeepers.InterchainstakingKeeper.SetZone(ctx, &agoricZone)

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

		ctx.Logger().Info("Set TransferChannel field for zones")
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			zone.TransferChannel = channels[zone.ChainId]
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

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

			// derive the correct ibc denom; if it does not match then remove it.
			correctIbc := utils.DeriveIbcDenom("transfer", channel.Counterparty.ChannelId, "transfer", channels[token.ChainID], token.QAssetDenom)
			if token.IbcDenom != correctIbc {
				ctx.Logger().Info(fmt.Sprintf("incorrect IBC denom %s for LiquidAllowedDenomProtocolData %s, expected %s. removing", token.IbcDenom, token.QAssetDenom, correctIbc))
				appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, pd.GenerateKey())
			}
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
