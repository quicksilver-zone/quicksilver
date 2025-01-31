package upgrades

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func V010700UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010702UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isMainnet(ctx) || isTest(ctx) {

			hashes := []struct {
				Zone string
				Hash string
			}{
				{Zone: "cosmoshub-4", Hash: "0c8269f04109a55a152d3cdfd22937b4e5c2746111d579935eef4cd7ffa71f7f"},
				{Zone: "stargaze-1", Hash: "10af0ee10a97f01467039a69cbfb8df05dc3111c975d955ca51adda201f36555"},
				{Zone: "stargaze-1", Hash: "0000000000000000000000000000000000000000000000000000000000000577"},
			}
			for _, hashRecord := range hashes {
				// delete duplicate records.
				appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, hashRecord.Zone, hashRecord.Hash, icstypes.WithdrawStatusUnbond)
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("delete duplicate withdrawal record", "hash", hashRecord.Hash, "zone", hashRecord.Zone)
			}

			// mint 50.699994 uqatom into escrow account
			err := appKeepers.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(50699994))))
			if err != nil {
				panic(err)
			}

			err = appKeepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqatom", sdk.NewInt(50699994))))
			if err != nil {
				panic(err)
			}

			// burn 16463.524950 qstars from escrow account
			err = appKeepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.EscrowModuleAccount, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqstars", sdk.NewInt(16463524950))))
			if err != nil {
				panic(err)
			}

			err = appKeepers.BankKeeper.BurnCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqstars", sdk.NewInt(16463524950))))
			if err != nil {
				panic(err)
			}

			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
				appKeepers.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, zone)
				return false
			})

		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010704UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isMainnet(ctx) || isTest(ctx) {

			// delete previously sent records. this record was not removed in the previous upgrade because it was not in the unbond status.
			hashes := []struct {
				Zone   string
				Hash   string
				Status int32
			}{
				{Zone: "stargaze-1", Hash: "10af0ee10a97f01467039a69cbfb8df05dc3111c975d955ca51adda201f36555", Status: icstypes.WithdrawStatusSend},
			}
			for _, hashRecord := range hashes {
				// delete duplicate records.
				appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, hashRecord.Zone, hashRecord.Hash, hashRecord.Status)
			}

			hashes = []struct {
				Zone   string
				Hash   string
				Status int32
			}{
				{Zone: "cosmoshub-4", Hash: "02c2d4bcb869b9ddf26540c2854c2ca09d70492a3831170da293f4101fda32b3", Status: icstypes.WithdrawStatusUnbond},
			}
			for _, hashRecord := range hashes {
				// requeue duplicate records.
				record, found := appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, hashRecord.Zone, hashRecord.Hash, hashRecord.Status)
				if !found {
					panic("record not found")
				}
				record.SendErrors = 0
				record.Amount = nil
				record.Distribution = nil
				record.CompletionTime = time.Time{}
				record.Acknowledged = false

				appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &record, icstypes.WithdrawStatusQueued)
			}

			// get zone
			zone, found := appKeepers.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
			if !found {
				panic("zone not found")
			}

			appKeepers.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, &zone)
			// set new redemption rate to 1.37
			zone.LastRedemptionRate = sdk.NewDecWithPrec(137, 2)
			appKeepers.InterchainstakingKeeper.SetZone(ctx, &zone)

			// remove 2x old icq queries that will never be satisfied.
			appKeepers.InterchainQueryKeeper.DeleteQuery(ctx, "d611198d85fed38e7486b9402480e561533911059a35258abce1220479b7bb7e")
			appKeepers.InterchainQueryKeeper.DeleteQuery(ctx, "533f78f574edcfe6153c753d7769072e86b6586cf9837fe9fec1ad84354433ec")
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010705UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isMainnet(ctx) || isTest(ctx) {
			// get zone
			zone, found := appKeepers.InterchainstakingKeeper.GetZone(ctx, "cosmoshub-4")
			if !found {
				panic("zone not found")
			}

			appKeepers.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, &zone)
			zone.LastRedemptionRate = sdk.NewDecWithPrec(138, 2) // correct as of 3/12
			appKeepers.InterchainstakingKeeper.SetZone(ctx, &zone)
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010706UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isMainnet(ctx) || isTest(ctx) {
			// unblock 2x stuck unbondings
			hashes := []struct {
				Zone   string
				Hash   string
				Status int32
			}{
				{Zone: "regen-1", Hash: "ee0b5f5c423508c8dd6a501168a77a0b72d5a8aaf1702a64804e522334ff272b", Status: icstypes.WithdrawStatusUnbond},
				{Zone: "sommelier-3", Hash: "a55f1f4deaa501ff5671ef96fbbb5b60e225d4b8db4825ae3706893bb94e052c", Status: icstypes.WithdrawStatusUnbond},
			}
			for _, hashRecord := range hashes {
				record, found := appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, hashRecord.Zone, hashRecord.Hash, hashRecord.Status)
				if !found {
					panic("record not found")
				}
				record.SendErrors = 0
				record.Amount = nil
				record.Distribution = nil
				record.CompletionTime = time.Time{}
				record.Acknowledged = false

				appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &record, icstypes.WithdrawStatusQueued)
			}
		}

		// garbage collect old records
		appKeepers.InterchainstakingKeeper.IterateUnbondingRecords(ctx, func(index int64, record icstypes.UnbondingRecord) (stop bool) {
			if record.CompletionTime.Equal(time.Time{}) || record.CompletionTime.Before(ctx.BlockTime().Add(-time.Hour*24)) { // old records
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("deleting old unbonding record", "chain_id", record.ChainId, "validator", record.Validator, "epoch_number", record.EpochNumber)
				appKeepers.InterchainstakingKeeper.DeleteUnbondingRecord(ctx, record.ChainId, record.Validator, record.EpochNumber)
			}
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
