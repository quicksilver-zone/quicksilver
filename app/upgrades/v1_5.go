package upgrades

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/quicksilver-zone/quicksilver/app/keepers"
	"github.com/quicksilver-zone/quicksilver/utils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
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

// =========== PRODUCTION UPGRADE HANDLER ===========

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

		// add claims metadata

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
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
			appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx, newRecords[key])
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
					distMap := map[string]uint64{}
					for _, dist := range dist1 {
						distMap[dist.Valoper] = dist.Amount
					}

					for _, dist := range dist2 {
						if _, ok = distMap[dist.Valoper]; !ok {
							distMap[dist.Valoper] = 0
						}
						distMap[dist.Valoper] += dist.Amount
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
			appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx, newRecords[key])
		}

		return false
	})
}
