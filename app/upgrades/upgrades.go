package upgrades

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func Upgrades() []Upgrade {
	return []Upgrade{
		// testnet upgrades
		{UpgradeName: V010405rc6UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010405rc7UpgradeName, CreateUpgradeHandler: NoOpHandler},

		// v1.2: this needs to be present to support upgrade on mainnet
		{UpgradeName: V010217UpgradeName, CreateUpgradeHandler: NoOpHandler},

		{UpgradeName: V010405UpgradeName, CreateUpgradeHandler: NoOpHandler},
	}
}

// no-op handler for upgrades with no state manipulation.
func NoOpHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010405UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			// add new fields
			zone.DepositsEnabled = true
			zone.ReturnToSender = false
			zone.UnbondingEnabled = true
			zone.Decimals = 6
			zone.Is_118 = true
			if zone.ChainId == "cosmoshub-4" {
				zone.LiquidityModule = true
			}

			// migrate all validators from within the zone struct, to own KV store.
			for _, val := range zone.Validators {
				newVal := icstypes.Validator{
					ValoperAddress:      val.ValoperAddress,
					CommissionRate:      val.CommissionRate,
					DelegatorShares:     val.DelegatorShares,
					VotingPower:         val.VotingPower,
					Score:               val.Score,
					Status:              val.Status,
					Jailed:              val.Jailed,
					Tombstoned:          val.Tombstoned,
					JailedSince:         val.JailedSince,
					ValidatorBondShares: val.ValidatorBondShares,
					LiquidShares:        val.LiquidShares,
				}
				err := appKeepers.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, newVal)
				if err != nil {
					panic(err)
				}

				// trigger a valset refresh to update all vals.
				query := stakingtypes.QueryValidatorsRequest{}
				appKeepers.InterchainstakingKeeper.EmitValSetQuery(ctx, zone.ConnectionId, zone.ChainId, query, math.NewInt(-1))
			}
			zone.Validators = nil
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

		// set lsm caps
		appKeepers.InterchainstakingKeeper.SetLsmCaps(ctx, "cosmoshub-4",
			icstypes.LsmCaps{
				ValidatorCap:     sdk.NewDecWithPrec(100, 2),
				ValidatorBondCap: sdk.NewDec(250),
				GlobalCap:        sdk.NewDecWithPrec(25, 2),
			},
		)

		// migrate vesting accounts for misplaced testnet wallets
		if err := migrateTestnetIncentives(ctx, appKeepers); err != nil {
			panic(fmt.Sprintf("unable to migrate testnet incentives: %v", err))
		}

		// migrate vesting account from ingenuity to notional
		if err := migrateIngenuityMultisigToNotional(ctx, appKeepers); err != nil {
			panic(fmt.Sprintf("unable to migrate ingenuity multisig: %v", err))
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// Migrate the outstanding testnet incentives with misplaced wallets.
// N.B. these assets are only returning to their original testnet addresses.
func migrateTestnetIncentives(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {

	migrations := map[string]string{
		"quick1qlckz3nplj3sf323n4ma7n75fmv60lpclq5ccc": "quick15dhqkz3mxxg4tt3m8uz5yy3mzfckgzzh5hpaqp",
		"quick1edavtxhdfs8luyvedgkjcxjc9dtvks3ve7etku": "quick1dz3y9k9harjal8nyqg3vl570aj7slaemmxgn86",
		"quick1pajjuywnj6w3y6pclp4tj55a7ngz9tp2z4pgep": "quick15sr0uhelt0hw4x7l9zsy4a7hqkaw6jepq4ald9",
		"quick1vhd4n5u8rsmsdgs4h7zsn4h4klsej6n8spvsl3": "quick12fyxjyxt64c2q5y0sdts6m4uxcy4cmff7l0ffx",
		"quick1rufya429ss9nlhdram0xkcu0jejsz5atap0xan": "quick124pvdf300p2wmq6cl8wwy2z0637du6ec0nhxen",
		"quick1f8jp5tr86gn5yvwecr7a4a9zypqf2mg85p96rw": "quick1f708swcmeej2ddfksyvtpaxe07fz0r03f79dlq",
	}
	return migrateVestingAccounts(ctx, appKeepers, migrations)
}

// Migrate the Ingenuity genesis allocation to Notional.
func migrateIngenuityMultisigToNotional(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	// migrate ingenuity multisig to notional multisig.
	migrations := map[string]string{
		"quick1e22za5qrqqp488h5p7vw2pfx8v0y4u444ufeuw": "quick1gxrks2rcj9gthzfgrkjk5lnk0g00cg0cpyntlm",
	}
	return migrateVestingAccounts(ctx, appKeepers, migrations)
}

// Migrate a map of address pairs and migrate from key -> value
func migrateVestingAccounts(ctx sdk.Context, appKeepers *keepers.AppKeepers, migrations map[string]string) error {
	for fromBech32, toBech32 := range migrations {
		from, err := addressutils.AccAddressFromBech32(fromBech32, "quick")
		if err != nil {
			return err
		}
		to, err := addressutils.AccAddressFromBech32(toBech32, "quick")
		if err != nil {
			return err
		}
		err = migratePeriodicVestingAccount(ctx, appKeepers, from, to)
		if err != nil {
			return err
		}
	}
	return nil
}

// Migrate a PeriodicVestingAccount from address A to address B, maintaining periods, amounts and end date.
func migratePeriodicVestingAccount(ctx sdk.Context, appKeepers *keepers.AppKeepers, from sdk.AccAddress, to sdk.AccAddress) error {
	oldAccount := appKeepers.AccountKeeper.GetAccount(ctx, from)
	// if the new account already exists in the account keeper, we should fail.
	if newAccount := appKeepers.AccountKeeper.GetAccount(ctx, to); newAccount != nil {
		return fmt.Errorf("unable to migrate vesting account; destination is already an account")
	}

	oldPva, ok := oldAccount.(*vestingtypes.PeriodicVestingAccount)
	if !ok {
		return fmt.Errorf("from account is not a PeriodicVestingAccount")
	}

	// copy the existing PVA.
	newPva := *oldPva

	// create a new baseVesting account with the address provided.
	newBva := vestingtypes.NewBaseVestingAccount(authtypes.NewBaseAccountWithAddress(to), oldPva.OriginalVesting, oldPva.EndTime)
	// change vesting end time so we are able to negate the token lock.
	// if the endDate has passed, we circumvent the period checking logic.
	oldPva.BaseVestingAccount.EndTime = ctx.BlockTime().Unix() - 1
	newPva.BaseVestingAccount = newBva

	// set the old pva (with the altered date), so we can transfer assets.
	appKeepers.AccountKeeper.SetAccount(ctx, oldPva)
	// set the new pva with the correct period and end dates, and new address.
	appKeepers.AccountKeeper.SetAccount(ctx, &newPva)

	// send coins from old account to new.
	err := appKeepers.BankKeeper.SendCoins(ctx, from, to, appKeepers.BankKeeper.GetAllBalances(ctx, from))
	if err != nil {
		return err
	}

	// delete the old account from the account keeper.
	appKeepers.AccountKeeper.RemoveAccount(ctx, oldPva)
	return nil
}
