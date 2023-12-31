package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/supply/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/mint keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// Supply returns supply and circulating supply of the staking denom.
func (q Querier) Supply(c context.Context, _ *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if q.endpointEnabled {
		baseDenom := q.stakingKeeper.BondDenom(ctx)
		supply := q.bankKeeper.GetSupply(ctx, baseDenom)
		circulatingSupply := q.CalculateCirculatingSupply(ctx, baseDenom, []string{
			"quick1yxe3vmd2ypjf0fs4cejnmv2559tqq5x5cc5nyh", // foundation account
			"quick1j5cgdlthhstqy2gqnglpjf4fvx3gs24yrcdtrf", // founder
			"quick1puj8yjmgrvn4w8vfswsnx972lucywetd57zalh", // founder
			"quick1d04jsq0kw4797kk4vp53y7hgmy8zdn8x7es279", // founder
			"quick1etqtc49wywy9ptx2gplhj0nrw5hy48hzzc6n20", // founder
			"quick1hdl587g7urer06myjkua86gc63vmq6pcr4d9hl", // founder
			"quick1ghwtkyrdr8lxm6x8dyr0nkqzghny955qe4j6zr", // founder
			"quick1a8dg5fuxtcwt8z6d9earl2sd0tukknx2txjm4j", // founder
			"quick1e22za5qrqqp488h5p7vw2pfx8v0y4u444ufeuw", // ingenuity
		})

		return &types.QuerySupplyResponse{Supply: supply.Amount.Uint64(), CirculatingSupply: circulatingSupply.Uint64()}, nil
	}
	return nil, fmt.Errorf("endpoint disabled")
}
