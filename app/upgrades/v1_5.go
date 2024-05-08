package upgrades

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prkeeper "github.com/quicksilver-zone/quicksilver/x/participationrewards/keeper"
	prtypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// =========== TESTNET UPGRADE HANDLER ===========

func V010500rc1UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// 993/1229 - pre-populate zone/denom mapping.
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			appKeepers.InterchainstakingKeeper.SetLocalDenomZoneMapping(ctx, zone)
			return false
		})

		// migrate notional vesting accounts to new addresses - source addresses are not prod multisigs, but test vesting accounts with delegations.
		migrations := map[string]string{
			"quick190yw7mfa8d8lgj9m4nyfh808s9pv7vz6cufff0": "quick1h0sqndv2y4xty6uk0sv4vckgyc5aa7n5at7fll",
			"quick14rptnkqsvwtumvezug6uvd537kxql8up3863cf": "quick1n4g6037cjm0e0v2nvwj2ngau7pk758wtwk6lwq",
		}

		if err := migrateVestingAccounts(ctx, appKeepers, migrations, migrateVestingAccountWithActions); err != nil {
			panic(err)
		}

		// initialise new withdrawal record sequence number
		appKeepers.InterchainstakingKeeper.InitWithdrawalRecordSequence(ctx)

		collateRequeuedWithdrawals(ctx, appKeepers)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010503rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// update all withdrawal records' distributors to use new Amount field.
		appKeepers.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(index int64, record icstypes.WithdrawalRecord) (stop bool) {
			if record.Distribution == nil {
				// skip records that are queued.
				return false
			}

			newDist := make([]*icstypes.Distribution, 0, len(record.Distribution))
			for _, d := range record.Distribution {
				d.Amount = math.NewIntFromUint64(d.XAmount)
				d.XAmount = 0
				newDist = append(newDist, d)
			}
			record.Distribution = newDist
			_ = appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)

			return false
		})

		// for all claims, update to use new Amount field
		appKeepers.ClaimsManagerKeeper.IterateAllClaims(ctx, func(index int64, key []byte, data cmtypes.Claim) (stop bool) {
			data.Amount = math.NewIntFromUint64(data.XAmount)
			appKeepers.ClaimsManagerKeeper.SetClaim(ctx, &data)
			return false
		})

		appKeepers.ClaimsManagerKeeper.IterateAllLastEpochClaims(ctx, func(index int64, key []byte, data cmtypes.Claim) (stop bool) {
			data.Amount = math.NewIntFromUint64(data.XAmount)
			appKeepers.ClaimsManagerKeeper.SetClaim(ctx, &data)
			return false
		})

		// for all redelegation records, migrate
		appKeepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(index int64, key []byte, record icstypes.RedelegationRecord) (stop bool) {
			record.Amount = math.NewInt(record.XAmount)
			appKeepers.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)
			return false
		})

		appKeepers.InterchainstakingKeeper.SetConnectionForPort(ctx, "connection-4", "icacontroller-osmo-test-5.performance")

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// =========== PRODUCTION UPGRADE HANDLER ===========

func V010505UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		appKeepers.BankKeeper.Logger(ctx).Info("migrating killer queen incentives...")
		// migrate killer queen incentives
		migrations := map[string]string{
			"quick1qfyntnmlvznvrkk9xqppmcxqcluv7wd74nmyus": "quick1898x2jpjfelg4jvl4hqm9a9vugyctfdcl9t64x",
		}
		err := migrateVestingAccounts(ctx, appKeepers, migrations, migratePeriodicVestingAccount)
		if err != nil {
			panic(err)
		}

		appKeepers.InterchainstakingKeeper.Logger(ctx).Info("setting minimum dust delegations thresholds...")

		// Update dust threshold configuration for all zones
		thresholds := map[string]math.Int{
			"osmosis-1":      math.NewInt(1_000_000),
			"cosmoshub-4":    math.NewInt(500_000),
			"stargaze-1":     math.NewInt(5_000_000),
			"juno-1":         math.NewInt(2_000_000),
			"sommelier-3":    math.NewInt(5_000_000),
			"regen-1":        math.NewInt(5_000_000),
			"dydx-mainnet-1": math.NewInt(500_000_000_000_000_000),
			"ssc-1":          math.NewInt(1_000_000),
		}
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			threshold, ok := thresholds[zone.ChainId]
			// if threshold not exist => get default value
			if !ok {
				threshold = math.NewInt(1_000_000)
			}
			zone.DustThreshold = threshold
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

		// iterate over all withdrawal records with zero BurnAmount and delete them; expecting a single record!
		appKeepers.InterchainstakingKeeper.Logger(ctx).Info("cleaning up unsatisfiable unbonding records with zero burnAmount...")
		appKeepers.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(_ int64, record icstypes.WithdrawalRecord) (stop bool) {
			if record.BurnAmount.IsZero() {
				appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("deleted withdrawal record...", "hash", record.Txhash)
			}
			return false
		})

		appKeepers.InterchainstakingKeeper.Logger(ctx).Info("setting 8 unsent unbondings for epoch 148 to STATUS_UNBONDING to be picked up by the next end blocker...")
		// remit epoch 148 unbondings
		hashes148 := []string{
			"1a9986d29cdff28014c59cc9ab055f03aebbe65ea5893f88e3582cec08a56a75",
			"350824b3097b4ab6bc17af19b87f352bcd4c0363f147b291ca7f84330e8b8aa9",
			"3ae014e377ba3f1955f58327ae9c38b7355d6e8e07e063f867d1c66901462dd5",
			"4f7884e67bb903fafa13252dcdd0ab358b80967424ccb94e571b70f4cefe4711",
			"8c61d25230f00c6f83c6e49d7af3f9ab35f3816f3b6740a0f07cf030c1a314d4",
			"a62019759694e50e16966460e78406111f08c31a673d33911b4617183d79d222",
			"b02040cca1e9dc7863425466eca4634a00671576889b3d98555404c03c75ab4f",
			"b8278c2deb36608886e31b02e027b09082f28ea881e185e9b255608fe6ef7399",
		}
		for _, hash := range hashes148 {
			record, found := appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "cosmoshub-4", hash, icstypes.WithdrawStatusSend)
			if !found {
				// do not panic, in case records were updated on 15/04 epoch.
				appKeepers.InterchainstakingKeeper.Logger(ctx).Error(fmt.Sprintf("1: unable to find record for hash %s", hash))
				continue
			}

			// update the record so that it will re-trigger the send.
			appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &record, icstypes.WithdrawStatusUnbond)
			appKeepers.InterchainstakingKeeper.Logger(ctx).Info("updated record to STATUS_UNBONDING", "hash", hash)
		}

		appKeepers.InterchainstakingKeeper.Logger(ctx).Info("setting 2 unsent unbondings from previous epochs to STATUS_QUEUED to be retriggered on the next epoch...")

		// requeue previous unbondings
		hashesPrevious := []string{
			"05348d2b80b9611e28c2b4782ab497596ed817a7d1e5f62e242ce02d36680604",
			"e0ee0bee23176ce635cb9412200d6cee3e157a6615fb3d5bddf57cd133592308",
		}

		for _, hash := range hashesPrevious {
			record, found := appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "cosmoshub-4", hash, icstypes.WithdrawStatusSend)
			if !found {
				// do not panic, in case records were updated on 15/04 epoch.
				appKeepers.InterchainstakingKeeper.Logger(ctx).Error(fmt.Sprintf("2: unable to find record for hash %s", hash))
				continue
			}

			// update the record to requeue. These records were subject to an
			record.Amount = nil
			record.CompletionTime = time.Time{}
			record.Requeued = true
			record.Distribution = nil

			appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &record, icstypes.WithdrawStatusQueued)
			appKeepers.InterchainstakingKeeper.Logger(ctx).Info("updated record to STATUS_QUEUED", "hash", hash)
		}

		// add qdydx and qsaga to claims for osmosis.
		appKeepers.ParticipationRewardsKeeper.Logger(ctx).Info("adding qdydx and qsaga ibc denoms to osmosis claimable tokens")
		channel, found := appKeepers.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", "channel-2")
		if !found {
			panic("unable to find channel osmosis transfer channel")
		}

		err = addProtocolData(ctx, appKeepers.ParticipationRewardsKeeper, prtypes.ProtocolDataTypeLiquidToken, &prtypes.LiquidAllowedDenomProtocolData{
			ChainID:               "osmosis-1",
			RegisteredZoneChainID: "dydx-mainnet-1",
			QAssetDenom:           "aqdydx",
			IbcDenom:              utils.DeriveIbcDenom("transfer", "channel-2", "transfer", channel.Counterparty.ChannelId, "aqdydx"),
		})
		if err != nil {
			panic("unable to add aqdydx on osmosis")
		}

		err = addProtocolData(ctx, appKeepers.ParticipationRewardsKeeper, prtypes.ProtocolDataTypeLiquidToken, &prtypes.LiquidAllowedDenomProtocolData{
			ChainID:               "osmosis-1",
			RegisteredZoneChainID: "ssc-1",
			QAssetDenom:           "uqsaga",
			IbcDenom:              utils.DeriveIbcDenom("transfer", "channel-2", "transfer", channel.Counterparty.ChannelId, "uqsaga"),
		})
		if err != nil {
			panic("unable to add uqsaga on osmosis")
		}

		appKeepers.UpgradeKeeper.Logger(ctx).Info("upgrade v1.5.5 completed!")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010504UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(_ int64, zone *icstypes.Zone) (stop bool) {
			// ship everything in local corresponding account to withdrawal account address to fee collector.
			account, err := addressutils.AccAddressFromBech32(zone.WithdrawalAddress.Address, "")
			if err != nil {
				panic(err)
			}

			accountBalance := appKeepers.BankKeeper.GetAllBalances(ctx, account)
			err = appKeepers.BankKeeper.SendCoinsFromAccountToModule(ctx, account, authtypes.FeeCollectorName, accountBalance)
			if err != nil {
				panic(err)
			}

			return false
		})
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010503UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// fix paid up juno records wher ack was never received.
		appKeepers.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, "juno-1", 4, func(index int64, record icstypes.WithdrawalRecord) (stop bool) {
			// burn burnAmount!
			if err := appKeepers.BankKeeper.BurnCoins(ctx, icstypes.EscrowModuleAccount, sdk.NewCoins(record.BurnAmount)); err != nil {
				// if we can't burn the coins, fail.
				panic(err)
			}
			appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)
			return false
		})

		// update all withdrawal records' distributors to use new Amount field.
		appKeepers.InterchainstakingKeeper.IterateWithdrawalRecords(ctx, func(index int64, record icstypes.WithdrawalRecord) (stop bool) {
			if record.Distribution == nil {
				// skip records that are queued.
				return false
			}

			newDist := make([]*icstypes.Distribution, 0, len(record.Distribution))
			for _, d := range record.Distribution {
				d.Amount = math.NewIntFromUint64(d.XAmount)
				d.XAmount = 0
				newDist = append(newDist, d)
			}
			record.Distribution = newDist
			if err := appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record); err != nil {
				panic(err)
			}

			return false
		})

		// for all claims, update to use new Amount field
		appKeepers.ClaimsManagerKeeper.IterateAllClaims(ctx, func(index int64, key []byte, data cmtypes.Claim) (stop bool) {
			data.Amount = math.NewIntFromUint64(data.XAmount)
			appKeepers.ClaimsManagerKeeper.SetClaim(ctx, &data)
			return false
		})

		appKeepers.ClaimsManagerKeeper.IterateAllLastEpochClaims(ctx, func(index int64, key []byte, data cmtypes.Claim) (stop bool) {
			data.Amount = math.NewIntFromUint64(data.XAmount)
			appKeepers.ClaimsManagerKeeper.SetLastEpochClaim(ctx, &data)
			return false
		})

		// for all redelegation records, migrate
		appKeepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(index int64, key []byte, record icstypes.RedelegationRecord) (stop bool) {
			record.Amount = math.NewInt(record.XAmount)
			appKeepers.InterchainstakingKeeper.SetRedelegationRecord(ctx, record)
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010501UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// update claims denoms
		if err := resetLiquidClaimsDenoms(ctx, appKeepers); err != nil {
			panic(err)
		}

		// iterate over all withdrawal records on cosmoshub with status 0 (undefined), and set to status 2 (queued)
		appKeepers.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, "cosmoshub-4", 0, func(index int64, record icstypes.WithdrawalRecord) (stop bool) {
			appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &record, icstypes.WithdrawStatusQueued)
			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010500UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// 993/1229 - pre-populate zone/denom mapping.
		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			appKeepers.InterchainstakingKeeper.SetLocalDenomZoneMapping(ctx, zone)
			return false
		})

		// migrate notional vesting accounts to new addresses
		migrations := map[string]string{
			"quick1a7n7z45gs0dut2syvkszffgwmgps6scqen3e5l": "quick1h0sqndv2y4xty6uk0sv4vckgyc5aa7n5at7fll",
			"quick1m0anwr4kcz0y9s65czusun2ahw35g3humv4j7f": "quick1n4g6037cjm0e0v2nvwj2ngau7pk758wtwk6lwq",
		}

		if err := migrateVestingAccounts(ctx, appKeepers, migrations, migrateVestingAccountWithActions); err != nil {
			panic(err)
		}

		// initialise new withdrawal record sequence number
		appKeepers.InterchainstakingKeeper.InitWithdrawalRecordSequence(ctx)

		// collate requeued withdrawal records
		collateRequeuedWithdrawals(ctx, appKeepers)

		if err := reimburseUsersWithdrawnOnLowRR(ctx, appKeepers); err != nil {
			panic(err)
		}

		// add claims metadata
		if err := initialiseClaimsMetaData(ctx, appKeepers); err != nil {
			panic(err)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func addProtocolData(ctx sdk.Context, keeper *prkeeper.Keeper, prtype prtypes.ProtocolDataType, data prtypes.ProtocolDataI) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	protocolData := prtypes.ProtocolData{
		Type: prtypes.ProtocolDataType_name[int32(prtype)],
		Data: jsonData,
	}

	keeper.SetProtocolData(ctx, data.GenerateKey(), &protocolData)
	return nil
}

// v1.5.0 had the port/channel and counterparty port/channel args flipped so calculated the IBC denom of
// qAssets as if they had come _from_ each chain, not _to_ each chain.
func resetLiquidClaimsDenoms(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	prk := appKeepers.ParticipationRewardsKeeper

	// transfer channels
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
	}
	var err error
	// ProtocolDataTypeConnection
	appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		// add liquid tokens for qasset on osmosis, secret, umee and the host zone itself.
		chainsToAdd := []string{"osmosis-1", "secret-4", "umee-1", zone.ChainId}
		for _, chain := range chainsToAdd {
			channel, found := appKeepers.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", channels[chain])
			if !found {
				err = fmt.Errorf("unable to find channel %s", channels[chain])
				return true
			}

			err = addProtocolData(ctx, prk, prtypes.ProtocolDataTypeLiquidToken, &prtypes.LiquidAllowedDenomProtocolData{
				ChainID:               chain,
				RegisteredZoneChainID: zone.ChainId,
				QAssetDenom:           zone.LocalDenom,
				IbcDenom:              utils.DeriveIbcDenom("transfer", channel.Counterparty.ChannelId, "transfer", channels[chain], zone.LocalDenom),
			})
			if err != nil {
				return true
			}
		}

		return false
	})

	return err
}

// initialiseClaimsMetaData does what it says on the tin, and bootstraps the claims metadata.
// note to future self: the logic to derive the ibc denom has reversed params, and this should
// not be used in future to determine qasset ibc denoms. v1.5.1 handler fixes this.
func initialiseClaimsMetaData(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	prk := appKeepers.ParticipationRewardsKeeper

	// transfer channels
	channels := map[string]string{
		"osmosis-1":   "channel-2",
		"cosmoshub-4": "channel-1",
		"stargaze-1":  "channel-0",
		"juno-1":      "channel-86",
		"sommelier-3": "channel-101",
		"regen-1":     "channel-17",
		"umee-1":      "channel-49",
		"secret-4":    "channel-52",
	}
	var err error
	// ProtocolDataTypeConnection
	appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		// add connection for each zone
		err = addProtocolData(ctx, prk, prtypes.ProtocolDataTypeConnection, &prtypes.ConnectionProtocolData{
			ConnectionID: zone.ConnectionId,
			ChainID:      zone.ChainId,
			Prefix:       zone.AccountPrefix,
		})
		if err != nil {
			return true
		}

		// add local (QS) denom for each chain
		err = addProtocolData(ctx, prk, prtypes.ProtocolDataTypeLiquidToken, &prtypes.LiquidAllowedDenomProtocolData{
			ChainID:               ctx.ChainID(),
			RegisteredZoneChainID: zone.ChainId,
			QAssetDenom:           zone.LocalDenom,
			IbcDenom:              zone.LocalDenom,
		})
		if err != nil {
			return true
		}

		// add liquid tokens for qasset on osmosis, secret, umee and the host zone itself.
		chainsToAdd := []string{"osmosis-1", "secret-4", "umee-1", zone.ChainId}
		for _, chain := range chainsToAdd {
			channel, found := appKeepers.IBCKeeper.ChannelKeeper.GetChannel(ctx, "transfer", channels[chain])
			if !found {
				err = fmt.Errorf("unable to find channel %s", channels[chain])
				return true
			}

			err = addProtocolData(ctx, prk, prtypes.ProtocolDataTypeLiquidToken, &prtypes.LiquidAllowedDenomProtocolData{
				ChainID:               chain,
				RegisteredZoneChainID: zone.ChainId,
				QAssetDenom:           zone.LocalDenom,
				IbcDenom:              utils.DeriveIbcDenom("transfer", channels[chain], "transfer", channel.Counterparty.ChannelId, zone.LocalDenom),
			})
			if err != nil {
				return true
			}
		}

		return false
	})

	if err != nil {
		return err
	}

	// osmosis params
	err = addProtocolData(ctx, prk, prtypes.ProtocolDataTypeOsmosisParams, &prtypes.OsmosisParamsProtocolData{
		ChainID:   "osmosis-1",
		BaseChain: "osmosis-1",
		BaseDenom: "uosmo",
	})
	if err != nil {
		return err
	}

	osmoPools := []*prtypes.OsmosisPoolProtocolData{
		// incentivised pools
		{
			PoolID:   903,
			PoolName: "qSTARS/STARS",
			PoolData: []byte(`{"address":"osmo1cxlrfu8r0v3cyqj78fuvlsmhjdgna0r7tum8cpd0g3x7w7pte8fsfvcs84","id":903,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/903","amount":"36839899550979704528582"},"pool_liquidity":[{"denom":"ibc/46C83BB054E12E189882B5284542DB605D94C99827E367C9192CF0579CD5BC83","amount":"238430845376"},{"denom":"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4","amount":"721246955212"}],"scaling_factors":[1,1],"scaling_factor_controller":""}`),
			PoolType: prtypes.PoolTypeStableSwap,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/46C83BB054E12E189882B5284542DB605D94C99827E367C9192CF0579CD5BC83": {Denom: "uqstars", ChainID: "stargaze-1"},
				"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4": {Denom: "ustars", ChainID: "stargaze-1"},
			},
			IsIncentivized: true,
		},
		{
			PoolID:   944,
			PoolName: "ATOM/qATOM",
			PoolData: []byte(`{"address":"osmo1awr39mc2hrkt8gq8gt3882ru40ay45k8a3yg69nyypqe9g0ryycs66lhkh","id":944,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/944","amount":"6108537302303463956540"},"pool_liquidity":[{"denom":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","amount":"42678069500"},{"denom":"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC","amount":"70488173547"}],"scaling_factors":[1202853876,1000000000],"scaling_factor_controller":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"}`),
			PoolType: prtypes.PoolTypeStableSwap,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {Denom: "uatom", ChainID: "cosmoshub-4"},
				"ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC": {Denom: "uqatom", ChainID: "cosmoshub-4"},
			},
			IsIncentivized: true,
		},
		{
			PoolID:   948,
			PoolName: "REGEN/qREGEN",
			PoolData: []byte(`{"address":"osmo1hylqy4uu5el36wykhzzhj786eh8rx4epyvg6nrtl503wjufz8z3sdptdzw","id":948,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/948","amount":"205748905065147865419005"},"pool_liquidity":[{"denom":"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076","amount":"258144972321"},{"denom":"ibc/79A676508A2ECA1021EDDC7BB9CF70CEEC9514C478DA526A5A8B3E78506C2206","amount":"171893185502"}],"scaling_factors":[1264169382,1000000000],"scaling_factor_controller":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"}`),
			PoolType: prtypes.PoolTypeStableSwap,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076": {Denom: "uregen", ChainID: "regen-1"},
				"ibc/79A676508A2ECA1021EDDC7BB9CF70CEEC9514C478DA526A5A8B3E78506C2206": {Denom: "uqregen", ChainID: "regen-1"},
			},
			IsIncentivized: true,
		},
		{
			PoolID:   956,
			PoolName: "qOSMO/OSMO",
			PoolData: []byte(`{"address":"osmo1q023e9m4d3ffvr96xwaeraa62yfvufkufkr7yf7lmacgkuspsuqsga4xp2","id":956,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/956","amount":"2443712227421775628369"},"pool_liquidity":[{"denom":"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC","amount":"5378958960"},{"denom":"uosmo","amount":"4563791740"}],"scaling_factors":[1000000000,1141529049],"scaling_factor_controller":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"}`),
			PoolType: prtypes.PoolTypeStableSwap,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC": {Denom: "uqosmo", ChainID: "osmosis-1"},
				"uosmo": {Denom: "", ChainID: "osmosis-1"},
			},
			IsIncentivized: true,
		},
		{
			PoolID:   1087,
			PoolName: "SOMM/qSOMM",
			PoolData: []byte(`{"address":"osmo1unwajz776rcsvaaehrq82qldwfw4zeqp7jgty09cw4lytuwfw3pqvs0cmt","id":1087,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000"},"future_pool_governor":"168h","total_shares":{"denom":"gamm/pool/1087","amount":"2103498772356422991482414"},"pool_liquidity":[{"denom":"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E","amount":"150351552009"},{"denom":"ibc/EAF76AD1EEF7B16D167D87711FB26ABE881AC7D9F7E6D0CF313D5FA530417208","amount":"273114101299"}],"scaling_factors":[1032934412,1000000000],"scaling_factor_controller":"osmo16x03wcp37kx5e8ehckjxvwcgk9j0cqnhm8m3yy"}`),
			PoolType: prtypes.PoolTypeStableSwap,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E": {Denom: "usomm", ChainID: "sommelier-3"},
				"ibc/EAF76AD1EEF7B16D167D87711FB26ABE881AC7D9F7E6D0CF313D5FA530417208": {Denom: "uqsomm", ChainID: "sommelier-3"},
			},
			IsIncentivized: true,
		},
		// price pools
		{
			PoolID:   1,
			PoolName: "ATOM/OSMO",
			PoolData: []byte(`{"address":"osmo1mw0ac6rwlp5r8wapwk3zs6g29h8fcscxqakdzw9emkne6c8wjp9q0t3v8t","id":1,"pool_params":{"swap_fee":"0.002000000000000000","exit_fee":"0.000000000000000000","smooth_weight_change_params":null},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/1","amount":"57132468094739651740591169"},"pool_assets":[{"token":{"denom":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","amount":"776914339293"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"6458328512048"},"weight":"536870912000000"}],"total_weight":"1073741824000000"}`),
			PoolType: prtypes.PoolTypeBalancer,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": {Denom: "uatom", ChainID: "cosmoshub-4"},
				"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
			},
			IsIncentivized: false,
		},
		{
			PoolID:   627,
			PoolName: "SOMM/OSMO",
			PoolData: []byte(`{"address":"osmo19qawwfrlkz9upglmpqj6akgz9ap7v2mnd05pxzgmxw3ywz58wnvqtet2mg","id":627,"pool_params":{"swap_fee":"0.002000000000000000","exit_fee":"0.000000000000000000","smooth_weight_change_params":null},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/627","amount":"65307069985087982755662"},"pool_assets":[{"token":{"denom":"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E","amount":"324082699777"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"35173517987"},"weight":"536870912000000"}],"total_weight":"1073741824000000"}`),
			PoolType: prtypes.PoolTypeBalancer,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E": {Denom: "usomm", ChainID: "sommelier-3"},
				"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
			},
			IsIncentivized: false,
		},
		{
			PoolID:   497,
			PoolName: "JUNO/OSMO",
			PoolData: []byte(`{"address":"osmo1h7yfu7x4qsv2urnkl4kzydgxegdfyjdry5ee4xzj98jwz0uh07rqdkmprr","id":497,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000","smooth_weight_change_params":null},"future_pool_governor":"","total_shares":{"denom":"gamm/pool/497","amount":"162333695811959156414166"},"pool_assets":[{"token":{"denom":"ibc/46B44899322F3CD854D2D46DEEF881958467CDD4B3B10086DA49296BBED94BED","amount":"611902251727"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"157137434465"},"weight":"536870912000000"}],"total_weight":"1073741824000000"}`),
			PoolType: prtypes.PoolTypeBalancer,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/46B44899322F3CD854D2D46DEEF881958467CDD4B3B10086DA49296BBED94BED": {Denom: "ujuno", ChainID: "juno-1"},
				"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
			},
			IsIncentivized: false,
		},
		{
			PoolID:   42,
			PoolName: "REGEN/OSMO",
			PoolData: []byte(`{"address":"osmo1txawpctjs6phpqsnkx2r5qud7yvekw93394anhuzz4dquy5jggssgqtn0l","id":42, "pool_params":{"swap_fee":"0.002000000000000000","exit_fee":"0.000000000000000000","smooth_weight_change_params":null},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/42","amount":"26709992201368381268416"},"pool_assets":[{"token":{"denom":"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076","amount":"1051059067702"},"weight":"536870912000000"},{"token":{"denom":"uosmo","amount":"36676123833"},"weight":"536870912000000"}],"total_weight":"1073741824000000"}`),
			PoolType: prtypes.PoolTypeBalancer,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/1DCC8A6CB5689018431323953344A9F6CC4D0BFB261E88C9F7777372C10CD076": {Denom: "uregen", ChainID: "regen-1"},
				"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
			},
			IsIncentivized: false,
		},
		{
			PoolID:   604,
			PoolName: "STARS/OSMO",
			PoolData: []byte(`{"address":"osmo1thscstwxp87g0ygh7le3h92f9ff4sel9y9d2eysa25p43yf43rysk7jp93","id":604,"pool_params":{"swap_fee":"0.003000000000000000","exit_fee":"0.000000000000000000","smooth_weight_change_params":null},"future_pool_governor":"24h","total_shares":{"denom":"gamm/pool/604","amount":"80971873633391327384"},"pool_assets":[{"token":{"denom":"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4","amount":"19005838130969"},"weight":"21474836480"},{"token":{"denom":"uosmo","amount":"465805350468"},"weight":"21474836480"}],"total_weight":"42949672960"}`),
			PoolType: prtypes.PoolTypeBalancer,
			Denoms: map[string]prtypes.DenomWithZone{
				"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4": {Denom: "ustars", ChainID: "stargaze-1"},
				"uosmo": {Denom: "uosmo", ChainID: "osmosis-1"},
			},
			IsIncentivized: false,
		},
	}
	// osmosis pools
	for _, pool := range osmoPools {
		err = addProtocolData(ctx, prk, prtypes.ProtocolDataTypeOsmosisPool, pool)
		if err != nil {
			return err
		}
	}

	// enable params
	params := prk.GetParams(ctx)
	params.ClaimsEnabled = true
	prk.SetParams(ctx, params)
	return nil
}

// for epochs 137 through 144 the redemption rate was adversely affected on cosmoshub-4 and unbonding users received less than they ought to have done.
// in order to compensate them, the portion of qAtoms they essentially didn't receive atoms for, we will re-mint, and create a new queued unbonding record.
// the below users will then receive the appropriate amount of atoms to make them whole. We have to remint the qatoms so that when the unbonding is complete
// there exists the requisite number of qatoms to burn.
func reimburseUsersWithdrawnOnLowRR(ctx sdk.Context, appKeepers *keepers.AppKeepers) error {
	users := map[string]struct {
		recipient string
		amount    sdk.Coin
	}{
		"quick143de92kvypafazd200r7fw4pwqjhnlsm724edv": {recipient: "cosmos143de92kvypafazd200r7fw4pwqjhnlsm4w9t57", amount: sdk.NewCoin("uqatom", sdk.NewInt(54013979))},
		"quick14jy373j0rr5pmpy33e7jlkujc0ve3rdx546ln5": {recipient: "cosmos14jy373j0rr5pmpy33e7jlkujc0ve3rdxl32d2x", amount: sdk.NewCoin("uqatom", sdk.NewInt(3416073))},
		"quick14xyjk9rnc24my8lchp04f3c0fvzrjgl07grk5y": {recipient: "cosmos14xyjk9rnc24my8lchp04f3c0fvzrjgl04vnydk", amount: sdk.NewCoin("uqatom", sdk.NewInt(227822))},
		"quick15xq28alrsk6plt4dp7ag7pjvtyangmx635826c": {recipient: "cosmos15xq28alrsk6plt4dp7ag7pjvtyangmx66shcr2", amount: sdk.NewCoin("uqatom", sdk.NewInt(12303565))},
		"quick1776mt7n23mwcat5vx0cr3x00qgug3d49f45ymy": {recipient: "cosmos1776mt7n23mwcat5vx0cr3x00qgug3d49z3ykzk", amount: sdk.NewCoin("uqatom", sdk.NewInt(538085))},
		"quick1af74tzu8j679405llklm3yanpkhneaq7mnf3xa": {recipient: "cosmos1af74tzu8j679405llklm3yanpkhneaq7sherl0", amount: sdk.NewCoin("uqatom", sdk.NewInt(161997))},
		"quick1alaq3havngy0h5sezl98a8xc0jx7xhad74p9nx": {recipient: "cosmos1alaq3havngy0h5sezl98a8xc0jx7xhad433h25", amount: sdk.NewCoin("uqatom", sdk.NewInt(8776885))},
		"quick1cl6qj3wmf7eynyta7h7a0lud9jemsj6dcaqhxz": {recipient: "cosmos1cl6qj3wmf7eynyta7h7a0lud9jemsj6dnes9ls", amount: sdk.NewCoin("uqatom", sdk.NewInt(2522168))},
		"quick1e4cnw86pl73k2sfv7uwauflfl42qzncna2j9tt": {recipient: "cosmos1e4cnw86pl73k2sfv7uwauflfl42qzncnkwzhje", amount: sdk.NewCoin("uqatom", sdk.NewInt(43791))},
		"quick1jd463sarmhp4zyd27jc9zzedmu8tyqdzhmfu4v": {recipient: "cosmos1jd463sarmhp4zyd27jc9zzedmu8tyqdzulewv7", amount: sdk.NewCoin("uqatom", sdk.NewInt(324773))},
		"quick1jjwf2052uy7fvl8tl65lgxnyr7mggc7vpeq29y": {recipient: "cosmos1jjwf2052uy7fvl8tl65lgxnyr7mggc7v2ascuk", amount: sdk.NewCoin("uqatom", sdk.NewInt(383443))},
		"quick1lcqquw54wdq07qeu2sx643cp5ppqzy8t5mqf34": {recipient: "cosmos1lcqquw54wdq07qeu2sx643cp5ppqzy8tllsmg8", amount: sdk.NewCoin("uqatom", sdk.NewInt(90487))},
		"quick1lz6udrmecnjsqhv48fd8ytd8truvdhd2hq6ytn": {recipient: "cosmos1lz6udrmecnjsqhv48fd8ytd8truvdhd2uy2kjp", amount: sdk.NewCoin("uqatom", sdk.NewInt(778236))},
		"quick1m0e7wr3k4h6xtc97psr66e7njkmv0e9a4l95k9": {recipient: "cosmos1m0e7wr3k4h6xtc97psr66e7njkmv0e9a7m4x0h", amount: sdk.NewCoin("uqatom", sdk.NewInt(316020))},
		"quick1m6lxmqfgf3s4vu0ktl78w2sz28e86v60sgckht": {recipient: "cosmos1m6lxmqfgf3s4vu0ktl78w2sz28e86v60mvgywe", amount: sdk.NewCoin("uqatom", sdk.NewInt(76642954))},
		"quick1mf40cxs57a4px5hj5ul0ute2ej5ec6e26xuzcw": {recipient: "cosmos1mf40cxs57a4px5hj5ul0ute2ej5ec6e23zvspu", amount: sdk.NewCoin("uqatom", sdk.NewInt(5590497))},
		"quick1pwpz0acvw0mc0clr4kknedt94efhwzj8zydzvk": {recipient: "cosmos1pwpz0acvw0mc0clr4kknedt94efhwzj8fqas4y", amount: sdk.NewCoin("uqatom", sdk.NewInt(462513))},
		"quick1pzzdvazgat8t9epvh2n5xn6wk4zcfc549xj5q9": {recipient: "cosmos1pzzdvazgat8t9epvh2n5xn6wk4zcfc54wzzxeh", amount: sdk.NewCoin("uqatom", sdk.NewInt(13187))},
		"quick1q0u34n7dujy3mlataslm2qlups9yxqwfwn35d0": {recipient: "cosmos1q0u34n7dujy3mlataslm2qlups9yxqwf9hpx5a", amount: sdk.NewCoin("uqatom", sdk.NewInt(5382251))},
		"quick1qltxuz7zak8rgx30xvenh6muwrkf8z2d8ffmat": {recipient: "cosmos1qltxuz7zak8rgx30xvenh6muwrkf8z2dvdefye", amount: sdk.NewCoin("uqatom", sdk.NewInt(2564863))},
		"quick1qsk66jfz02x9r6433xdj5ptkpfp07ytk7ephk3": {recipient: "cosmos1qsk66jfz02x9r6433xdj5ptkpfp07ytk4a390r", amount: sdk.NewCoin("uqatom", sdk.NewInt(91503))},
		"quick1r83cmscpqhj36pltqt8msqkcxsnpkl4zqqk8xa": {recipient: "cosmos1r83cmscpqhj36pltqt8msqkcxsnpkl4ztyx4l0", amount: sdk.NewCoin("uqatom", sdk.NewInt(170790))},
		"quick1snvzr84cv8esmlwpcfqg26tfxndn3xwda889w3": {recipient: "cosmos1snvzr84cv8esmlwpcfqg26tfxndn3xwdkrhhhr", amount: sdk.NewCoin("uqatom", sdk.NewInt(1400544))},
		"quick1t3cwpvu4nrk2zqt9tmhsgkk4ra465q8eqvdljz": {recipient: "cosmos1t3cwpvu4nrk2zqt9tmhsgkk4ra465q8etgadts", amount: sdk.NewCoin("uqatom", sdk.NewInt(1668204))},
		"quick1uxpfv475505ylmwhxt8qmz6ewpur5hzhtkhat6": {recipient: "cosmos1uxpfv475505ylmwhxt8qmz6ewpur5hzhqj80jg", amount: sdk.NewCoin("uqatom", sdk.NewInt(41000))},
		"quick1vlfa0p6qm69hyu2zxcfy9zzuqhwkqwzn5tq6zh": {recipient: "cosmos1vlfa0p6qm69hyu2zxcfy9zzuqhwkqwznl0sgm9", amount: sdk.NewCoin("uqatom", sdk.NewInt(754997))},
		"quick1yr8fgts6d76g0u847zkng2e9l9nk4stw5dkzpu": {recipient: "cosmos1yr8fgts6d76g0u847zkng2e9l9nk4stwlfxscw", amount: sdk.NewCoin("uqatom", sdk.NewInt(1598954))},
	}

	for _, delegator := range utils.Keys(users) {

		// mint the coins
		if err := appKeepers.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(users[delegator].amount)); err != nil {
			return err
		}

		// send them to the escrow module account
		if err := appKeepers.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(users[delegator].amount)); err != nil {
			return err
		}

		if err := appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx,
			icstypes.WithdrawalRecord{
				ChainId:      "cosmoshub-4",
				Delegator:    delegator,
				Recipient:    users[delegator].recipient,
				BurnAmount:   users[delegator].amount,
				Txhash:       fmt.Sprintf("%064d", appKeepers.InterchainstakingKeeper.GetNextWithdrawalRecordSequence(ctx)),
				Requeued:     false,
				Acknowledged: false,
				Distribution: nil,
				EpochNumber:  145,
			},
		); err != nil {
			return err
		}
	}
	return nil
}

// collateRequeuedWithdrawals will iterate, per zone, over requeued queued and active withdrawal records and
// collate them into a single record for a delegator/recipient/epoch tuple.
func collateRequeuedWithdrawals(ctx sdk.Context, appKeepers *keepers.AppKeepers) {
	appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(_ int64, zone *icstypes.Zone) (stop bool) {
		newRecords := map[string]icstypes.WithdrawalRecord{}

		appKeepers.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icstypes.WithdrawStatusQueued, func(_ int64, record icstypes.WithdrawalRecord) (stop bool) {
			if !record.Requeued {
				return false
			}

			// this is a requeued record.
			mapKey := fmt.Sprintf("%s/%s", record.Delegator, record.Recipient)
			newRecord, ok := newRecords[mapKey]
			if !ok {
				newRecord = icstypes.WithdrawalRecord{
					ChainId:        record.ChainId,
					Delegator:      record.Delegator,
					Distribution:   nil,
					Recipient:      record.Recipient,
					Amount:         nil,
					BurnAmount:     record.BurnAmount,
					Txhash:         fmt.Sprintf("%064d", appKeepers.InterchainstakingKeeper.GetNextWithdrawalRecordSequence(ctx)),
					Status:         icstypes.WithdrawStatusQueued,
					CompletionTime: time.Time{},
					Requeued:       true,
					Acknowledged:   false,
					EpochNumber:    record.EpochNumber,
				}
			} else {
				newRecord.BurnAmount = newRecord.BurnAmount.Add(record.BurnAmount)
			}
			newRecords[mapKey] = newRecord

			// delete old record
			appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)

			return false
		})

		for _, key := range utils.Keys(newRecords) {
			if err := appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx, newRecords[key]); err != nil {
				panic(err)
			}
		}

		newRecords = map[string]icstypes.WithdrawalRecord{}

		appKeepers.InterchainstakingKeeper.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, icstypes.WithdrawStatusUnbond, func(_ int64, record icstypes.WithdrawalRecord) (stop bool) {
			if !record.Requeued || !record.Acknowledged {
				return false
			}

			// this is a requeued AND acknowledged record.
			mapKey := fmt.Sprintf("%s/%s/%d", record.Delegator, record.Recipient, record.EpochNumber)
			newRecord, ok := newRecords[mapKey]
			if !ok {
				newRecord = icstypes.WithdrawalRecord{
					ChainId:        record.ChainId,
					Delegator:      record.Delegator,
					Distribution:   record.Distribution,
					Recipient:      record.Recipient,
					Amount:         record.Amount,
					BurnAmount:     record.BurnAmount,
					Txhash:         fmt.Sprintf("%064d", appKeepers.InterchainstakingKeeper.GetNextWithdrawalRecordSequence(ctx)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: record.CompletionTime,
					Requeued:       true,
					Acknowledged:   true,
					EpochNumber:    record.EpochNumber,
				}
			} else {
				newRecord.BurnAmount = newRecord.BurnAmount.Add(record.BurnAmount)
				newRecord.Amount = newRecord.Amount.Add(record.Amount...)
				// update completion time if incoming is later.
				if record.CompletionTime.After(newRecord.CompletionTime) {
					newRecord.CompletionTime = record.CompletionTime
				}
				// merge distributions
				newRecord.Distribution = func(dist1, dist2 []*icstypes.Distribution) []*icstypes.Distribution {
					distMap := map[string]math.Int{}
					for _, dist := range dist1 {
						distMap[dist.Valoper] = dist.Amount
					}

					for _, dist := range dist2 {
						if _, ok = distMap[dist.Valoper]; !ok {
							distMap[dist.Valoper] = math.ZeroInt()
						}
						distMap[dist.Valoper] = distMap[dist.Valoper].Add(dist.Amount)
					}

					out := make([]*icstypes.Distribution, 0, len(distMap))
					for _, key := range utils.Keys(distMap) {
						out = append(out, &icstypes.Distribution{Valoper: key, Amount: distMap[key]})
					}

					return out
				}(newRecord.Distribution, record.Distribution)
			}

			newRecords[mapKey] = newRecord

			// delete old record
			appKeepers.InterchainstakingKeeper.DeleteWithdrawalRecord(ctx, record.ChainId, record.Txhash, record.Status)

			return false
		})

		for _, key := range utils.Keys(newRecords) {
			if err := appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx, newRecords[key]); err != nil {
				panic(err)
			}
		}

		return false
	})
}
