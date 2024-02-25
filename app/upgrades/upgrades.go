package upgrades

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
		{UpgradeName: V010407rc0UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010407rc1UpgradeName, CreateUpgradeHandler: V010407rc1UpgradeHandler},
		{UpgradeName: V010407rc2UpgradeName, CreateUpgradeHandler: V010407rc2UpgradeHandler},
		{UpgradeName: V010500rc0UpgradeName, CreateUpgradeHandler: NoOpHandler},

		// v1.2: this needs to be present to support upgrade on mainnet
		{UpgradeName: V010217UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010405UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010406UpgradeName, CreateUpgradeHandler: V010406UpgradeHandler},
		{UpgradeName: V010407UpgradeName, CreateUpgradeHandler: V010407UpgradeHandler},
		{UpgradeName: V010600UpgradeName, CreateUpgradeHandler: V010600UpgradeHandler},
	}
}

// NoOpHandler no-op handler for upgrades with no state manipulation.
func NoOpHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010406UpgradeHandler(
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
				_ = appKeepers.InterchainstakingKeeper.EmitValSetQuery(ctx, zone.ConnectionId, zone.ChainId, query, math.NewInt(-1))
			}

			appKeepers.InterchainstakingKeeper.SetAddressZoneMapping(ctx, zone.DepositAddress.Address, zone.ChainId)
			appKeepers.InterchainstakingKeeper.SetAddressZoneMapping(ctx, zone.DelegationAddress.Address, zone.ChainId)
			appKeepers.InterchainstakingKeeper.SetAddressZoneMapping(ctx, zone.PerformanceAddress.Address, zone.ChainId)
			appKeepers.InterchainstakingKeeper.SetAddressZoneMapping(ctx, zone.WithdrawalAddress.Address, zone.ChainId)

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

func V010407rc1UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) {
			// remove osmo-test-5 so we can reinstate
			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, "osmo-test-5")
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010407rc2UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isDevnet(ctx) {
			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
				vals := appKeepers.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(vals))}}
				bz := appKeepers.InterchainstakingKeeper.GetCodec().MustMarshal(&delegationQuery)

				appKeepers.InterchainQueryKeeper.MakeRequest(
					ctx,
					zone.ConnectionId,
					zone.ChainId,
					"cosmos.staking.v1beta1.Query/DelegatorDelegations",
					bz,
					sdk.NewInt(-1),
					icstypes.ModuleName,
					"delegations",
					0,
				)

				balancesQuery := banktypes.QueryAllBalancesRequest{Address: zone.DelegationAddress.Address}
				bz = appKeepers.InterchainstakingKeeper.GetCodec().MustMarshal(&balancesQuery)
				appKeepers.InterchainQueryKeeper.MakeRequest(
					ctx,
					zone.ConnectionId,
					zone.ChainId,
					"cosmos.bank.v1beta1.Query/AllBalances",
					bz,
					sdk.NewInt(-1),
					icstypes.ModuleName,
					"delegationaccountbalances",
					0,
				)
				// increment waitgroup; decremented in delegationaccountbalance callback
				zone.WithdrawalWaitgroup++

				appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)

				return false
			})
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010407UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// set redemption rate to value correct as of 2024-02-01T20:00UTC.
		rates := map[string]sdk.Dec{
			"cosmoshub-4": sdk.NewDecFromInt(math.NewInt(219_280_116_789)).QuoInt64(186_283_929_157),
			"stargaze-1":  sdk.NewDecFromInt(math.NewInt(7_883_310_380_922)).QuoInt64(6_142_958_768_078),
			"osmosis-1":   sdk.NewDecFromInt(math.NewInt(363_909_524_952)).QuoInt64(322_912_055_083),
			"sommelier-3": sdk.NewDecFromInt(math.NewInt(657_103_764_225)).QuoInt64(637_871_903_193),
			"regen-1":     sdk.NewDecFromInt(math.NewInt(5_606_819_529_428)).QuoInt64(4_543_207_966_192),
			"juno-1":      sdk.NewDecFromInt(math.NewInt(7_439_000_263)).QuoInt64(7_018_171_980),
		}

		// trigger redemption rate update immediately after upgrade.
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			vals := appKeepers.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
			delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(vals))}}
			bz := appKeepers.InterchainstakingKeeper.GetCodec().MustMarshal(&delegationQuery)

			appKeepers.InterchainQueryKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				"cosmos.staking.v1beta1.Query/DelegatorDelegations",
				bz,
				sdk.NewInt(-1),
				icstypes.ModuleName,
				"delegations",
				0,
			)

			balancesQuery := banktypes.QueryAllBalancesRequest{Address: zone.DelegationAddress.Address}
			bz = appKeepers.InterchainstakingKeeper.GetCodec().MustMarshal(&balancesQuery)
			appKeepers.InterchainQueryKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				"cosmos.bank.v1beta1.Query/AllBalances",
				bz,
				sdk.NewInt(-1),
				icstypes.ModuleName,
				"delegationaccountbalances",
				0,
			)
			// increment waitgroup; decremented in delegationaccountbalance callback
			zone.WithdrawalWaitgroup++
			zone.RedemptionRate = rates[zone.ChainId]
			zone.LastRedemptionRate = rates[zone.ChainId]
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)

			return false
		})

		// migrate testnet user account.
		migrations := map[string]string{
			"quick1k67rz3vn73tzp2tatlka2kn2ngtjdw8gpw8zq2": "quick1plq2mrsn0uw2dkksptr9dsyyk62dkk6t7w79j2",
		}

		if err := migrateVestingAccounts(ctx, appKeepers, migrations); err != nil {
			panic(err)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010600UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	// TODO must add test and refactor current duplicated logic out of app/upgrades
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		migrations := map[string]string{
			"quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l": "quick1h0sqndv2y4xty6uk0sv4vckgyc5aa7n5at7fll",
			"quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f": "quick1n4g6037cjm0e0v2nvwj2ngau7pk758wtwk6lwq",
		}

		for fromBech32, toBech32 := range migrations {
			from, err := addressutils.AccAddressFromBech32(fromBech32, "quick")
			if err != nil {
				return nil, err
			}
			to, err := addressutils.AccAddressFromBech32(toBech32, "quick")
			if err != nil {
				return nil, err
			}
			err = processMigratePeriodicVestingAccount(ctx, appKeepers, from, to)
			if err != nil {
				return nil, err
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// TODO: add logic collect rewards
// processMigratePeriodicVestingAccount moves the unvested from current account ->  new account
// - Unbonding all assets
// - Claim rewards
// - Post migrations
func processMigratePeriodicVestingAccount(ctx sdk.Context, appKeepers *keepers.AppKeepers, from sdk.AccAddress, to sdk.AccAddress) error {
	// Unbond all delagation of account
	unbonded, err := unbondAllDelegation(ctx, ctx.BlockTime(), appKeepers, from)
	if err != nil {
		fmt.Printf("processMigratePeriodicVestingAccount: unbonded all delegation failed: %v", err)
		return err
	}
	fmt.Printf("processMigratePeriodicVestingAccount: unbond all delegation amount: %s", unbonded)

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
	err = appKeepers.BankKeeper.SendCoins(ctx, from, to, appKeepers.BankKeeper.GetAllBalances(ctx, from))
	if err != nil {
		return err
	}

	// delete the old account from the account keeper.
	appKeepers.AccountKeeper.RemoveAccount(ctx, oldPva)
	return nil
}

func unbondAllDelegation(ctx sdk.Context, now time.Time, appKeepers *keepers.AppKeepers, accAddr sdk.AccAddress) (math.Int, error) {
	unbondedAmt := math.ZeroInt()

	// Undelegate all delegations from the account
	for _, delegation := range appKeepers.StakingKeeper.GetAllDelegatorDelegations(ctx, accAddr) {
		validatorValAddr := delegation.GetValidatorAddr()
		_, found := appKeepers.StakingKeeper.GetValidator(ctx, validatorValAddr)
		if !found {
			continue
		}

		_, err := appKeepers.StakingKeeper.Undelegate(ctx, accAddr, validatorValAddr, delegation.GetShares())
		if err != nil {
			return math.ZeroInt(), err
		}
	}

	// Complete unbonding of all account's delegations
	for _, unbondingDelegation := range appKeepers.StakingKeeper.GetAllUnbondingDelegations(ctx, accAddr) {
		validatorStringAddr := unbondingDelegation.ValidatorAddress
		validatorValAddr, _ := sdk.ValAddressFromBech32(validatorStringAddr)

		for i := range unbondingDelegation.Entries {
			unbondingDelegation.Entries[i].CompletionTime = now
			unbondedAmt = unbondedAmt.Add(unbondingDelegation.Entries[i].Balance)
		}

		appKeepers.StakingKeeper.SetUnbondingDelegation(ctx, unbondingDelegation)
		_, err := appKeepers.StakingKeeper.CompleteUnbonding(ctx, accAddr, validatorValAddr)
		if err != nil {
			return math.ZeroInt(), err
		}
	}

	return unbondedAmt, nil
}
