package interchainstaking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// InitGenesis initializes the interchainstaking module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	// set registered zones info from genesis
	for _, zone := range genState.Zones {
		// explicit memory referencing
		zone := zone
		k.SetZone(ctx, &zone)
	}

	for _, pc := range genState.PortConnections {
		k.SetConnectionForPort(ctx, pc.ConnectionId, pc.PortId)
	}

	for _, delegationForZone := range genState.Delegations {
		zone, found := k.GetZone(ctx, delegationForZone.ChainId)
		if !found {
			panic("unable to find zone for delegation")
		}
		for _, delegation := range delegationForZone.Delegations {
			k.SetDelegation(ctx, &zone, *delegation)
		}
	}

	for _, perfDelegationForZone := range genState.PerformanceDelegations {
		zone, found := k.GetZone(ctx, perfDelegationForZone.ChainId)
		if !found {
			panic("unable to find zone for delegation")
		}
		for _, delegation := range perfDelegationForZone.Delegations {
			k.SetPerformanceDelegation(ctx, &zone, *delegation)
		}
	}

	for _, delegatorIntentsForZone := range genState.DelegatorIntents {
		zone, found := k.GetZone(ctx, delegatorIntentsForZone.ChainId)
		if !found {
			panic("unable to find zone for delegation")
		}
		for _, delegatorIntent := range delegatorIntentsForZone.DelegationIntent {
			k.SetDelegatorIntent(ctx, &zone, *delegatorIntent, false)
		}
	}

	for _, receipt := range genState.Receipts {
		k.SetReceipt(ctx, receipt)
	}

	for _, withdrawal := range genState.WithdrawalRecords {
		k.SetWithdrawalRecord(ctx, withdrawal)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:                 k.GetParams(ctx),
		Zones:                  k.AllZones(ctx),
		Receipts:               k.AllReceipts(ctx),
		Delegations:            ExportDelegationsPerZone(ctx, k),
		PerformanceDelegations: ExportPerformanceDelegationsPerZone(ctx, k),
		DelegatorIntents:       ExportDelegatorIntentsPerZone(ctx, k),
		PortConnections:        k.AllPortConnections(ctx),
		WithdrawalRecords:      k.AllWithdrawalRecords(ctx),
	}
}

func ExportDelegationsPerZone(ctx sdk.Context, k *keeper.Keeper) []types.DelegationsForZone {
	delegationsForZones := make([]types.DelegationsForZone, 0)
	k.IterateZones(ctx, func(_ int64, zone *types.Zone) (stop bool) {
		delegationsForZones = append(delegationsForZones, types.DelegationsForZone{ChainId: zone.ChainId, Delegations: k.GetAllDelegationsAsPointer(ctx, zone)})
		return false
	})
	return delegationsForZones
}

func ExportPerformanceDelegationsPerZone(ctx sdk.Context, k *keeper.Keeper) []types.DelegationsForZone {
	delegationsForZones := make([]types.DelegationsForZone, 0)
	k.IterateZones(ctx, func(_ int64, zone *types.Zone) (stop bool) {
		delegationsForZones = append(delegationsForZones, types.DelegationsForZone{ChainId: zone.ChainId, Delegations: k.GetAllPerformanceDelegationsAsPointer(ctx, zone)})
		return false
	})
	return delegationsForZones
}

func ExportDelegatorIntentsPerZone(ctx sdk.Context, k *keeper.Keeper) []types.DelegatorIntentsForZone {
	delegatorIntentsForZones := make([]types.DelegatorIntentsForZone, 0)
	k.IterateZones(ctx, func(_ int64, zone *types.Zone) (stop bool) {
		// export current epoch intents
		delegatorIntentsForZones = append(delegatorIntentsForZones,
			types.DelegatorIntentsForZone{ChainId: zone.ChainId, DelegationIntent: k.AllDelegatorIntentsAsPointer(ctx, zone, false), Snapshot: false},
			// export last epoch intents
			types.DelegatorIntentsForZone{ChainId: zone.ChainId, DelegationIntent: k.AllDelegatorIntentsAsPointer(ctx, zone, true), Snapshot: true},
		)
		return false
	})
	return delegatorIntentsForZones
}
