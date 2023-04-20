package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ingenuity-build/quicksilver/app/keepers"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func Upgrades() []Upgrade {
	return []Upgrade{
		{UpgradeName: V010402rc1UpgradeName, CreateUpgradeHandler: V010402rc1UpgradeHandler},
	}
}

func NoOpHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010402rc1UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isTest(ctx) {
			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zoneInfo *types.Zone) (stop bool) {
				for _, val := range zoneInfo.Validators {
					newVal := types.Validator{
						ValoperAddress:  val.ValoperAddress,
						CommissionRate:  val.CommissionRate,
						DelegatorShares: val.DelegatorShares,
						VotingPower:     val.VotingPower,
						Score:           val.Score,
						Status:          val.Status,
						Jailed:          val.Jailed,
						Tombstoned:      val.Tombstoned,
						JailedSince:     val.JailedSince,
					}
					err := appKeepers.InterchainstakingKeeper.SetValidator(ctx, zoneInfo.ChainId, newVal)
					if err != nil {
						panic(err)
					}
				}
				zoneInfo.Validators = nil
				appKeepers.InterchainstakingKeeper.SetZone(ctx, zoneInfo)
				return false
			})
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// func V010400UpgradeHandler(
//	mm *module.Manager,
//	configurator module.Configurator,
//	appKeepers *keepers.AppKeepers,
// ) upgradetypes.UpgradeHandler {
//	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
//		// upgrade zones
//		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
//			zone.DepositsEnabled = true
//			zone.ReturnToSender = false
//			zone.UnbondingEnabled = false
//			zone.Decimals = 6
//			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
//			return false
//		})
//
//		// upgrade receipts
//		blockTime := ctx.BlockTime()
//		for _, r := range appKeepers.InterchainstakingKeeper.AllReceipts(ctx) {
//			r.FirstSeen = &blockTime
//			r.Completed = &blockTime
//			appKeepers.InterchainstakingKeeper.SetReceipt(ctx, r)
//		}
//		if isTestnet(ctx) || isTest(ctx) {
//
//			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "uni-5")
//
//			// burn uqjunox
//			addr1, err := utils.AccAddressFromBech32("quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq", "quick")
//			if err != nil {
//				return nil, err
//			}
//			addr2, err := utils.AccAddressFromBech32("quick16x03wcp37kx5e8ehckjxvwcgk9j0cqnhcccnty", "quick")
//			if err != nil {
//				return nil, err
//			}
//
//			err = appKeepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(1600000))))
//			if err != nil {
//				return nil, err
//			}
//
//			err = appKeepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr2, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(200000000))))
//			if err != nil {
//				return nil, err
//			}
//
//			err = appKeepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.EscrowModuleAccount, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(400000))))
//			if err != nil {
//				return nil, err
//			}
//
//			err = appKeepers.BankKeeper.BurnCoins(ctx, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(202000000))))
//			if err != nil {
//				return nil, err
//			}
//		}
//		return mm.RunMigrations(ctx, configurator, fromVM)
//	}
//}
//
// func V010400rc6UpgradeHandler(
//	mm *module.Manager,
//	configurator module.Configurator,
//	appKeepers *keepers.AppKeepers,
// ) upgradetypes.UpgradeHandler {
//	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
//		if isTestnet(ctx) {
//			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "regen-redwood-1")
//			// re-register regen-redwood-1 with new connection
//			regenProp := icstypes.NewRegisterZoneProposal("register regen-redwood-1 zone",
//				"register regen-redwood-1  (regen-testnet) zone with multisend and lsm disabled",
//				"connection-8",
//				"uregen",
//				"uqregen",
//				"regen",
//				false,
//				true,
//				true,
//				false,
//				6)
//			err := appKeepers.InterchainstakingKeeper.HandleRegisterZoneProposal(ctx, regenProp)
//			if err != nil {
//				return nil, err
//			}
//		}
//
//		// remove expired failed redelegation records
//		appKeepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(_ int64, key []byte, record icstypes.RedelegationRecord) (stop bool) {
//			if record.CompletionTime.Equal(time.Time{}) {
//				appKeepers.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, record.ChainId, record.Source, record.Destination, record.EpochNumber)
//			}
//			return false
//		})
//
//		// remove and refund failed unbondings
//		appKeepers.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(index int64, record icstypes.WithdrawalRecord) (stop bool) {
//			if record.Status == icskeeper.WithdrawStatusUnbond && record.CompletionTime.Equal(time.Time{}) {
//				delegatorAcc, err := utils.AccAddressFromBech32(record.Delegator, "quick")
//				if err != nil {
//					panic(err)
//				}
//				if err = appKeepers.InterchainstakingKeeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.EscrowModuleAccount, delegatorAcc, sdk.NewCoins(record.BurnAmount)); err != nil {
//					panic(err)
//				}
//				appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
//			}
//			return false
//		})
//
//		if isTestnet(ctx) || isDevnet(ctx) {
//			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zoneInfo *icstypes.Zone) (stop bool) {
//				appKeepers.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, zoneInfo)
//				return false
//			})
//		}
//
//		return mm.RunMigrations(ctx, configurator, fromVM)
//	}
// }
//
// func V010400rc8UpgradeHandler(
//	mm *module.Manager,
//	configurator module.Configurator,
//	appKeepers *keepers.AppKeepers,
// ) upgradetypes.UpgradeHandler {
//	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
//		// remove expired failed redelegation records
//		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
//			appKeepers.InterchainstakingKeeper.IterateAllDelegations(ctx, zone, func(delegation icstypes.Delegation) (stop bool) {
//				if delegation.RedelegationEnd < 0 {
//					delegation.RedelegationEnd = 0
//					appKeepers.InterchainstakingKeeper.SetDelegation(ctx, zone, delegation)
//				}
//				return false
//			})
//			return false
//		})
//
//		appKeepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(_ int64, key []byte, record icstypes.RedelegationRecord) (stop bool) {
//			if record.CompletionTime.Unix() <= 0 {
//				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("Removing delegation record", "chainid", record.ChainId, "source", record.Source, "destination", record.Destination, "epoch", record.EpochNumber)
//				appKeepers.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, record.ChainId, record.Source, record.Destination, record.EpochNumber)
//			}
//			return false
//		})
//
//		return mm.RunMigrations(ctx, configurator, fromVM)
//	}
// }
