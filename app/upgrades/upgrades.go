package upgrades

import (
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/ingenuity-build/quicksilver/app/keepers"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	tokenfactorytypes "github.com/ingenuity-build/quicksilver/x/tokenfactory/types"
)

func Upgrades() []Upgrade {
	return []Upgrade{
		{
			UpgradeName:          V010300UpgradeName,
			CreateUpgradeHandler: NoOpHandler,
			StoreUpgrades:        storetypes.StoreUpgrades{},
		},
		{
			UpgradeName:          V010400UpgradeName,
			CreateUpgradeHandler: V010400UpgradeHandler,
			StoreUpgrades:        storetypes.StoreUpgrades{},
		},
		{
			UpgradeName:          V010400rc6UpgradeName,
			CreateUpgradeHandler: V010400rc6UpgradeHandler,
			StoreUpgrades:        storetypes.StoreUpgrades{},
		},
		{
			UpgradeName:          V010400rc7UpgradeName,
			CreateUpgradeHandler: NoOpHandler,
			StoreUpgrades:        storetypes.StoreUpgrades{},
		},
		{
			UpgradeName:          V010400rc8UpgradeName,
			CreateUpgradeHandler: V010400rc8UpgradeHandler,
			StoreUpgrades:        storetypes.StoreUpgrades{},
		},
	}
}

func NoOpHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010400UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// upgrade zones
		keepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone icstypes.Zone) (stop bool) {
			zone.DepositsEnabled = true
			zone.ReturnToSender = false
			zone.UnbondingEnabled = false
			zone.Decimals = 6
			keepers.InterchainstakingKeeper.SetZone(ctx, &zone)
			return false
		})

		// upgrade receipts
		time := ctx.BlockTime()
		for _, r := range keepers.InterchainstakingKeeper.AllReceipts(ctx) {
			r.FirstSeen = &time
			r.Completed = &time
			keepers.InterchainstakingKeeper.SetReceipt(ctx, r)
		}
		if isTestnet(ctx) || isTest(ctx) {

			keepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "uni-5")

			// burn uqjunox
			addr1, err := utils.AccAddressFromBech32("quick17v9kk34km3w6hdjs2sn5s5qjdu2zrm0m3rgtmq", "quick")
			if err != nil {
				return nil, err
			}
			addr2, err := utils.AccAddressFromBech32("quick16x03wcp37kx5e8ehckjxvwcgk9j0cqnhcccnty", "quick")
			if err != nil {
				return nil, err
			}

			err = keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(1600000))))
			if err != nil {
				return nil, err
			}

			err = keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, addr2, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(200000000))))
			if err != nil {
				return nil, err
			}

			err = keepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.EscrowModuleAccount, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(400000))))
			if err != nil {
				return nil, err
			}

			err = keepers.BankKeeper.BurnCoins(ctx, tokenfactorytypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqjunox", sdkmath.NewInt(202000000))))
			if err != nil {
				return nil, err
			}
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010400rc6UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) {
			keepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "regen-redwood-1")
			// re-register regen-redwood-1 with new connection
			regenProp := icstypes.NewRegisterZoneProposal("register regen-redwood-1 zone",
				"register regen-redwood-1  (regen-testnet) zone with multisend and lsm disabled",
				"connection-8",
				"uregen",
				"uqregen",
				"regen",
				false,
				true,
				true,
				false,
				6)
			err := icskeeper.HandleRegisterZoneProposal(ctx, keepers.InterchainstakingKeeper, regenProp)
			if err != nil {
				return nil, err
			}
		}

		// remove expired failed redelegation records
		keepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(_ int64, key []byte, record icstypes.RedelegationRecord) (stop bool) {
			if record.CompletionTime.Equal(time.Time{}) {
				keepers.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, record.ChainId, record.Source, record.Destination, record.EpochNumber)
			}
			return false
		})

		// remove and refund failed unbondings
		keepers.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(index int64, record icstypes.WithdrawalRecord) (stop bool) {
			if record.Status == icskeeper.WithdrawStatusUnbond && record.CompletionTime.Equal(time.Time{}) {
				delegatorAcc, err := utils.AccAddressFromBech32(record.Delegator, "quick")
				if err != nil {
					panic(err)
				}
				if err = keepers.InterchainstakingKeeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.EscrowModuleAccount, delegatorAcc, sdk.NewCoins(record.BurnAmount)); err != nil {
					panic(err)
				}
				keepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
			}
			return false
		})

		if isTestnet(ctx) || isDevnet(ctx) {
			keepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zoneInfo icstypes.Zone) (stop bool) {
				keepers.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, zoneInfo)
				return false
			})
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010400rc8UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// remove expired failed redelegation records
		keepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone icstypes.Zone) (stop bool) {
			keepers.InterchainstakingKeeper.IterateAllDelegations(ctx, &zone, func(delegation icstypes.Delegation) (stop bool) {
				if delegation.RedelegationEnd < 0 {
					delegation.RedelegationEnd = 0
					keepers.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
				}
				return false
			})
			return false
		})

		keepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(_ int64, key []byte, record icstypes.RedelegationRecord) (stop bool) {
			if record.CompletionTime.Unix() <= 0 {
				keepers.InterchainstakingKeeper.Logger(ctx).Info("Removing delegation record", "chainid", record.ChainId, "source", record.Source, "destination", record.Destination, "epoch", record.EpochNumber)
				keepers.InterchainstakingKeeper.DeleteRedelegationRecord(ctx, record.ChainId, record.Source, record.Destination, record.EpochNumber)
			}
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
