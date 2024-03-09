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

		if err := reimburseUsersWithdrawnOnLowRR(ctx, appKeepers); err != nil {
			panic(err)
		}

		// add claims metadata

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
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

		appKeepers.InterchainstakingKeeper.SetWithdrawalRecord(ctx,
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
		)
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
