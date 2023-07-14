package upgrades

import (
	"errors"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/ingenuity-build/quicksilver/app/keepers"
	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	epochtypes "github.com/ingenuity-build/quicksilver/x/epochs/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	prtypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func Upgrades() []Upgrade {
	return []Upgrade{
		{UpgradeName: V010402rc1UpgradeName, CreateUpgradeHandler: V010402rc1UpgradeHandler},
		{UpgradeName: V010402rc2UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010402rc3UpgradeName, CreateUpgradeHandler: V010402rc3UpgradeHandler},
		{UpgradeName: V010402rc4UpgradeName, CreateUpgradeHandler: V010402rc4UpgradeHandler},
		{UpgradeName: V010402rc5UpgradeName, CreateUpgradeHandler: V010402rc5UpgradeHandler},
		{UpgradeName: V010402rc6UpgradeName, CreateUpgradeHandler: V010402rc6UpgradeHandler},
		{UpgradeName: V010402rc7UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010403rc0UpgradeName, CreateUpgradeHandler: V010403rc0UpgradeHandler},
		{UpgradeName: V010404beta0UpgradeName, CreateUpgradeHandler: V010404beta0UpgradeHandler},
		{UpgradeName: V010404beta1UpgradeName, CreateUpgradeHandler: NoOpHandler},
		{UpgradeName: V010404beta5UpgradeName, CreateUpgradeHandler: V010404beta5UpgradeHandler},
		{UpgradeName: V010404beta7UpgradeName, CreateUpgradeHandler: V010404beta7UpgradeHandler},
		{UpgradeName: V010404rc0UpgradeName, CreateUpgradeHandler: V010404rc0UpgradeHandler},
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
			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
				for _, val := range zone.Validators {
					newVal := icstypes.Validator{
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
					err := appKeepers.InterchainstakingKeeper.SetValidator(ctx, zone.ChainId, newVal)
					if err != nil {
						panic(err)
					}
				}
				zone.Validators = nil
				appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
				return false
			})
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010402rc3UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isTest(ctx) {
			appKeepers.InterchainstakingKeeper.RemoveZoneAndAssociatedRecords(ctx, OsmosisTestnetChainID)
			pdType, exists := prtypes.ProtocolDataType_value["ProtocolDataTypeConnection"]
			if !exists {
				panic("connection protocol data type not found")
			}

			appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, prtypes.GetProtocolDataKey(prtypes.ProtocolDataType(pdType), []byte("rege-redwood-1")))
			vals := appKeepers.InterchainstakingKeeper.GetValidators(ctx, OsmosisTestnetChainID)
			for _, val := range vals {
				valoper, _ := addressutils.ValAddressFromBech32(val.ValoperAddress, "osmovaloper")
				appKeepers.InterchainstakingKeeper.DeleteValidator(ctx, OsmosisTestnetChainID, valoper)
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010402rc4UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isTest(ctx) {
			pdType, exists := prtypes.ProtocolDataType_value["ProtocolDataTypeLiquidToken"]
			if !exists {
				panic("liquid tokens protocol data type not found")
			}
			appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, prtypes.GetProtocolDataKey(prtypes.ProtocolDataType(pdType), []byte("osmo-test-5/ibc/FBD3AC18A981B89F60F9FE5B21BD7F1DE87A53C3505D5A5E438E2399409CFB6F")))
			appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, prtypes.GetProtocolDataKey(prtypes.ProtocolDataType(pdType), []byte("rhye-1/uqosmo")))
			rcptTime := time.Unix(1682932342, 0)
			rcpt1 := icstypes.Receipt{
				ChainId:   "theta-testnet-001",
				Sender:    "cosmos1e6p7tk969ftlzmz82drp84ruukwge6z6udand8",
				Txhash:    "005AABC399866544CBEC4DC57887A7297289BF40C056A1544D3CE18946DB7DB9",
				Amount:    sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100000000))),
				FirstSeen: &rcptTime,
				Completed: nil,
			}

			rcpt2 := icstypes.Receipt{
				ChainId:   "elgafar-1",
				Sender:    "stars1e6p7tk969ftlzmz82drp84ruukwge6z6g32wxk",
				Txhash:    "01041964B4CDDD3ECA1C9F1EFC039B547C2D30D5B85C55089EB6F7DF311786B6",
				Amount:    sdk.NewCoins(sdk.NewCoin("ustars", sdkmath.NewInt(100000000))),
				FirstSeen: &rcptTime,
				Completed: nil,
			}

			appKeepers.InterchainstakingKeeper.SetReceipt(ctx, rcpt1)
			appKeepers.InterchainstakingKeeper.SetReceipt(ctx, rcpt2)

		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010402rc5UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isTest(ctx) {

			rcptTime := time.Unix(1682932342, 0)

			rcpts := []icstypes.Receipt{
				{
					ChainId:   "theta-testnet-001",
					Sender:    "cosmos1e6p7tk969ftlzmz82drp84ruukwge6z6udand8",
					Txhash:    "005AABC399866544CBEC4DC57887A7297289BF40C056A1544D3CE18946DB7DB9",
					Amount:    sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100000000))),
					FirstSeen: &rcptTime,
					Completed: nil,
				},
				{
					ChainId:   "theta-testnet-001",
					Sender:    "cosmos1e6p7tk969ftlzmz82drp84ruukwge6z6udand8",
					Txhash:    "60DBC8449D74B5782D5918A908F16AFF594E0E4CF28CAD82B9B329610ED01B80",
					Amount:    sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200000000))),
					FirstSeen: &rcptTime,
					Completed: nil,
				},
				{
					ChainId:   "theta-testnet-001",
					Sender:    "cosmos1e6p7tk969ftlzmz82drp84ruukwge6z6udand8",
					Txhash:    "2BB80824D07D3B2FA5B69E23C973D3B4885A4C8407871DDEFC324305748366BA",
					Amount:    sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(150000000))),
					FirstSeen: &rcptTime,
					Completed: nil,
				},
				{
					ChainId:   "elgafar-1",
					Sender:    "stars1e6p7tk969ftlzmz82drp84ruukwge6z6g32wxk",
					Txhash:    "01041964B4CDDD3ECA1C9F1EFC039B547C2D30D5B85C55089EB6F7DF311786B6",
					Amount:    sdk.NewCoins(sdk.NewCoin("ustars", sdkmath.NewInt(100000000))),
					FirstSeen: &rcptTime,
					Completed: nil,
				},
				{
					ChainId:   "elgafar-1",
					Sender:    "stars1e6p7tk969ftlzmz82drp84ruukwge6z6g32wxk",
					Txhash:    "74E497648091F539A47965EC8EDCA36F54329A5FEFC417F5BD28DD2C8297BBAC",
					Amount:    sdk.NewCoins(sdk.NewCoin("ustars", sdkmath.NewInt(200000000))),
					FirstSeen: &rcptTime,
					Completed: nil,
				},
				{
					ChainId:   "uni-6",
					Sender:    "juno1f6g9guyeyzgzjc9l8wg4xl5x0rvxddew0wx2jp",
					Txhash:    "11C89B3342326B8C84B0FDE1C63DC233B51E56D8EA6E1AB2B0BAD8094E77C38B",
					Amount:    sdk.NewCoins(sdk.NewCoin("ujunox", sdkmath.NewInt(200000000))),
					FirstSeen: &rcptTime,
					Completed: nil,
				},
				{
					ChainId:   "regen-redwood-1",
					Sender:    "regen1f6g9guyeyzgzjc9l8wg4xl5x0rvxddewx7wdre",
					Txhash:    "D5D1C2741A613E1303D32A7755DFC68072D23BCA281CE24D2A4A7857A8674D3B",
					Amount:    sdk.NewCoins(sdk.NewCoin("uregen", sdkmath.NewInt(200000000))),
					FirstSeen: &rcptTime,
					Completed: nil,
				},
			}

			for _, rcpt := range rcpts {
				appKeepers.InterchainstakingKeeper.SetReceipt(ctx, rcpt)
			}

		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010402rc6UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isTest(ctx) {
			// for each zone, trigger an icq request to update all delegations.
			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
				vals := appKeepers.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(vals))}}
				bz := appKeepers.InterchainstakingKeeper.GetCodec().MustMarshal(&delegationQuery)

				appKeepers.InterchainstakingKeeper.ICQKeeper.MakeRequest(
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
				return false
			})
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010403rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isTest(ctx) {
			appKeepers.ParticipationRewardsKeeper.IteratePrefixedProtocolDatas(ctx, prtypes.GetPrefixProtocolDataKey(prtypes.ProtocolDataTypeLiquidToken), func(index int64, key []byte, data prtypes.ProtocolData) (stop bool) {
				prefixedKey := append(prtypes.GetPrefixProtocolDataKey(prtypes.ProtocolDataTypeLiquidToken), key...)
				appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, prefixedKey)
				pd, err := prtypes.UnmarshalProtocolData(prtypes.ProtocolDataTypeLiquidToken, data.Data)
				if err != nil {
					panic(err)
				}
				newKey := pd.GenerateKey()
				appKeepers.ParticipationRewardsKeeper.SetProtocolData(ctx, newKey, &data)
				return false
			})
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010404beta0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isTestnet(ctx) || isTest(ctx) {
			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
				zone.Is_118 = true
				appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
				return false
			})
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010404beta5UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if isDevnet(ctx) || isTest(ctx) {
			// 6d3cc69d3276dd59a93a252e1ea15fc1e507c56512266c87c615fac4dcddb5cb
			wr, found := appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "theta-testnet-001", "6d3cc69d3276dd59a93a252e1ea15fc1e507c56512266c87c615fac4dcddb5cb", 3)
			if !found {
				return nil, errors.New("unable to find withdrawal record 6d3cc69d3276dd59a93a252e1ea15fc1e507c56512266c87c615fac4dcddb5cb")
			}
			appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &wr, icstypes.WithdrawStatusQueued)

			// b9c6587af3317bfb4b21a29df3f7e1a00709c25b0590446cceb01b8c6996b656
			wr, found = appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "theta-testnet-001", "b9c6587af3317bfb4b21a29df3f7e1a00709c25b0590446cceb01b8c6996b656", 3)
			if !found {
				return nil, errors.New("unable to find withdrawal record b9c6587af3317bfb4b21a29df3f7e1a00709c25b0590446cceb01b8c6996b656")
			}
			appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &wr, icstypes.WithdrawStatusQueued)

			// 995c6a77a568a7c03906ce6c7d470c11daa7e506f33264360cf1fec71fc774fe
			wr, found = appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "regen-redwood-1", "995c6a77a568a7c03906ce6c7d470c11daa7e506f33264360cf1fec71fc774fe", 4)
			if !found {
				return nil, errors.New("unable to find withdrawal record 995c6a77a568a7c03906ce6c7d470c11daa7e506f33264360cf1fec71fc774fe")
			}
			appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &wr, icstypes.WithdrawStatusUnbond)

			// 95aec506a8281c90cb45395ecc7b562248135f8643e1017db469d847db125fbd
			wr, found = appKeepers.InterchainstakingKeeper.GetWithdrawalRecord(ctx, "uni-6", "95aec506a8281c90cb45395ecc7b562248135f8643e1017db469d847db125fbd", 4)
			if !found {
				return nil, errors.New("unable to find withdrawal record 95aec506a8281c90cb45395ecc7b562248135f8643e1017db469d847db125fbd")
			}
			appKeepers.InterchainstakingKeeper.UpdateWithdrawalRecordStatus(ctx, &wr, icstypes.WithdrawStatusUnbond)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010404beta7UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		const (
			thetaUnbondingPeriod = int64(172800)
			uniUnbondingPeriod   = int64(2419200)
			osmoUnbondingPeriod  = int64(86400)
			regenUnbondingPeriod = int64(1814400)
			epochDurations       = int64(10800)
		)

		appKeepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(idx int64, key []byte, redelegation icstypes.RedelegationRecord) (stop bool) {
			var UnbondingPeriod int64
			switch redelegation.ChainId {
			case "theta-testnet-001":
				UnbondingPeriod = thetaUnbondingPeriod
			case "uni-6":
				UnbondingPeriod = uniUnbondingPeriod
			case "osmo-test-5":
				UnbondingPeriod = osmoUnbondingPeriod
			case "regen-redwood-1":
				UnbondingPeriod = regenUnbondingPeriod
			}

			epochInfo := appKeepers.EpochsKeeper.GetEpochInfo(ctx, epochtypes.EpochIdentifierEpoch)

			if UnbondingPeriod < (epochInfo.CurrentEpoch-redelegation.EpochNumber)*epochDurations {
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("garbage collecting completed redelegations", "key", key, "completion", redelegation.CompletionTime)
				appKeepers.InterchainstakingKeeper.DeleteRedelegationRecordByKey(ctx, append(icstypes.KeyPrefixRedelegationRecord, key...))
			}

			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func V010404rc0UpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		const (
			thetaUnbondingPeriod = int64(172800)
			uniUnbondingPeriod   = int64(2419200)
			osmoUnbondingPeriod  = int64(86400)
			regenUnbondingPeriod = int64(1814400)
			epochDurations       = int64(43200)
		)

		appKeepers.InterchainstakingKeeper.IterateRedelegationRecords(ctx, func(idx int64, key []byte, redelegation icstypes.RedelegationRecord) (stop bool) {
			var UnbondingPeriod int64
			switch redelegation.ChainId {
			case "theta-testnet-001":
				UnbondingPeriod = thetaUnbondingPeriod
			case "uni-6":
				UnbondingPeriod = uniUnbondingPeriod
			case "osmo-test-5":
				UnbondingPeriod = osmoUnbondingPeriod
			case "regen-redwood-1":
				UnbondingPeriod = regenUnbondingPeriod
			}

			epochInfo := appKeepers.EpochsKeeper.GetEpochInfo(ctx, epochtypes.EpochIdentifierEpoch)

			if UnbondingPeriod < (epochInfo.CurrentEpoch-redelegation.EpochNumber)*epochDurations {
				appKeepers.InterchainstakingKeeper.Logger(ctx).Info("garbage collecting completed redelegations", "key", key, "completion", redelegation.CompletionTime)
				appKeepers.InterchainstakingKeeper.DeleteRedelegationRecordByKey(ctx, append(icstypes.KeyPrefixRedelegationRecord, key...))
			}

			return false
		})

		if isTestnet(ctx) || isTest(ctx) {
			appKeepers.ParticipationRewardsKeeper.IteratePrefixedProtocolDatas(ctx, prtypes.GetPrefixProtocolDataKey(prtypes.ProtocolDataTypeLiquidToken), func(index int64, key []byte, data prtypes.ProtocolData) (stop bool) {
				prefixedKey := append(prtypes.GetPrefixProtocolDataKey(prtypes.ProtocolDataTypeLiquidToken), key...)
				appKeepers.ParticipationRewardsKeeper.DeleteProtocolData(ctx, prefixedKey)
				pd, err := prtypes.UnmarshalProtocolData(prtypes.ProtocolDataTypeLiquidToken, data.Data)
				if err != nil {
					panic(err)
				}
				newKey := pd.GenerateKey()
				appKeepers.ParticipationRewardsKeeper.SetProtocolData(ctx, newKey, &data)
				return false
			})
		}

		appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
			zone.Is_118 = true
			appKeepers.InterchainstakingKeeper.SetZone(ctx, zone)
			return false
		})

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
//			appKeepers.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
//				appKeepers.InterchainstakingKeeper.OverrideRedemptionRateNoCap(ctx, zone)
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
