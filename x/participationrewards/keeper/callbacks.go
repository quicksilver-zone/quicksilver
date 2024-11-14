package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/concentrated-liquidity/model"
	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm"
	umeetypes "github.com/quicksilver-zone/quicksilver/third-party-chains/umee-types/leverage/types"
	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

const (
	ValidatorSelectionRewardsCallbackID       = "validatorselectionrewards"
	OsmosisPoolUpdateCallbackID               = "osmosispoolupdate"
	OsmosisClPoolUpdateCallbackID             = "osmosisclpoolupdate"
	SetEpochBlockCallbackID                   = "epochblock"
	UmeeReservesUpdateCallbackID              = "umeereservesupdatecallback"
	UmeeTotalBorrowsUpdateCallbackID          = "umeetotalborrowsupdatecallback"
	UmeeInterestScalarUpdateCallbackID        = "umeeinterestscalarupdatecallback"
	UmeeUTokenSupplyUpdateCallbackID          = "umeeutokensupplyupdatecallback"
	UmeeLeverageModuleBalanceUpdateCallbackID = "umeeleveragemodulebalanceupdatecallback"
)

// Callback wrapper struct for interchainstaking keeper.
type Callback func(sdk.Context, *Keeper, []byte, icqtypes.Query) error

type Callbacks struct {
	k         *Keeper
	callbacks map[string]Callback
}

var _ icqtypes.QueryCallbacks = Callbacks{}

func (k *Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]Callback)}
}

// Call calls callback handler.
func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	if !c.Has(id) {
		return fmt.Errorf("callback %s not found", id)
	}
	return c.callbacks[id](ctx, c.k, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id], _ = fn.(Callback)
	return c
}

func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	a := c.
		AddCallback(ValidatorSelectionRewardsCallbackID, Callback(ValidatorSelectionRewardsCallback)).
		AddCallback(OsmosisPoolUpdateCallbackID, Callback(OsmosisPoolUpdateCallback)).
		AddCallback(OsmosisClPoolUpdateCallbackID, Callback(OsmosisClPoolUpdateCallback)).
		AddCallback(SetEpochBlockCallbackID, Callback(SetEpochBlockCallback)).
		AddCallback(UmeeReservesUpdateCallbackID, Callback(UmeeReservesUpdateCallback)).
		AddCallback(UmeeTotalBorrowsUpdateCallbackID, Callback(UmeeTotalBorrowsUpdateCallback)).
		AddCallback(UmeeInterestScalarUpdateCallbackID, Callback(UmeeInterestScalarUpdateCallback)).
		AddCallback(UmeeUTokenSupplyUpdateCallbackID, Callback(UmeeUTokenSupplyUpdateCallback)).
		AddCallback(UmeeLeverageModuleBalanceUpdateCallbackID, Callback(UmeeLeverageModuleBalanceUpdateCallback))

	return a.(Callbacks)
}

// Callbacks

func ValidatorSelectionRewardsCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	delegatorRewards := distrtypes.QueryDelegationTotalRewardsResponse{}
	err := k.cdc.Unmarshal(response, &delegatorRewards)
	if err != nil {
		return err
	}

	zone, found := k.icsKeeper.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	zs, err := k.getZoneScores(ctx, zone, delegatorRewards)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info(
		"callback zone score",
		"zone", zs.ZoneID,
		"total voting power", zs.TotalVotingPower,
		"validator scores", zs.ValidatorScores,
	)

	// snapshot obtained and used here
	userAllocations := k.CalcUserValidatorSelectionAllocations(ctx, &zone, *zs)

	if err := k.DistributeToUsersFromModule(ctx, userAllocations); err != nil {
		return err
	}

	// create snapshot of current intents for next epoch boundary
	for _, di := range k.icsKeeper.AllDelegatorIntents(ctx, &zone, false) {
		k.icsKeeper.SetDelegatorIntent(ctx, &zone, di, true)
	}

	// set zone ValidatorSelectionAllocation to zero
	zone.ValidatorSelectionAllocation = 0
	k.icsKeeper.SetZone(ctx, &zone)

	return nil
}

func OsmosisPoolUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	var pd gamm.CFMMPoolI
	if err := k.cdc.UnmarshalInterface(response, &pd); err != nil {
		return err
	}

	// check query.Request is at least 9 bytes in length. (0x02 + 8 bytes for uint64)
	if len(query.Request) < 9 {
		return errors.New("query request not sufficient length")
	}
	// assert first character is 0x02 as expected.
	if query.Request[0] != 0x02 {
		return errors.New("query request has unexpected prefix")
	}

	poolID := sdk.BigEndianToUint64(query.Request[1:])
	key := fmt.Sprintf("%d", poolID)
	data, pool, err := GetAndUnmarshalProtocolData[*types.OsmosisPoolProtocolData](ctx, k, key, types.ProtocolDataTypeOsmosisPool)
	if err != nil {
		return err
	}

	pool.PoolData, err = json.Marshal(pd)
	if err != nil {
		return err
	}
	pool.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(pool)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, pool.GenerateKey(), &data)

	return nil
}

func OsmosisClPoolUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	var pd model.Pool
	if err := k.cdc.Unmarshal(response, &pd); err != nil {
		return err
	}

	// check query.Request is at least 2 bytes - 0x03 + poolID
	if len(query.Request) < 2 {
		return errors.New("query request not sufficient length")
	}
	// assert first character is 0x03 as expected (cl pool prefix)
	if query.Request[0] != 0x03 {
		return errors.New("query request has unexpected prefix")
	}

	poolID, err := strconv.ParseInt(string(query.Request[1:]), 10, 64)
	if err != nil {
		return err
	}

	data, pool, err := GetAndUnmarshalProtocolData[*types.OsmosisClPoolProtocolData](ctx, k, fmt.Sprintf("%d", poolID), types.ProtocolDataTypeOsmosisCLPool)
	if err != nil {
		return err
	}

	pool.PoolData, err = json.Marshal(pd)
	if err != nil {
		return err
	}
	pool.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(pool)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, pool.GenerateKey(), &data)

	return nil
}

func UmeeReservesUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	reserveAmount := sdk.ZeroInt()
	if err := reserveAmount.Unmarshal(response); err != nil {
		return err
	}

	if query.Request[0] != umeetypes.KeyPrefixReserveAmount[0] {
		return errors.New("query request has unexpected prefix")
	}

	denom := umeetypes.DenomFromKey(query.Request, umeetypes.KeyPrefixReserveAmount)
	data, reserves, err := GetAndUnmarshalProtocolData[*types.UmeeReservesProtocolData](ctx, k, denom, types.ProtocolDataTypeUmeeReserves)
	if err != nil {
		return err
	}

	reserves.Data, err = json.Marshal(reserveAmount)
	if err != nil {
		return err
	}
	reserves.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(reserves)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, reserves.GenerateKey(), &data)

	return nil
}

func UmeeTotalBorrowsUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	if g, w := query.Request[0], umeetypes.KeyPrefixAdjustedTotalBorrow[0]; g != w {
		return fmt.Errorf("unexpected query request prefix %q, want %q", g, w)
	}

	totalBorrows := sdk.ZeroDec()
	if err := totalBorrows.Unmarshal(response); err != nil {
		return err
	}

	denom := umeetypes.DenomFromKey(query.Request, umeetypes.KeyPrefixAdjustedTotalBorrow)
	data, borrows, err := GetAndUnmarshalProtocolData[*types.UmeeTotalBorrowsProtocolData](ctx, k, denom, types.ProtocolDataTypeUmeeTotalBorrows)
	if err != nil {
		return err
	}

	borrows.Data, err = json.Marshal(totalBorrows)
	if err != nil {
		return err
	}
	borrows.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(borrows)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, borrows.GenerateKey(), &data)

	return nil
}

func UmeeInterestScalarUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	if g, w := query.Request[0], umeetypes.KeyPrefixInterestScalar[0]; g != w {
		return fmt.Errorf("unexpected query request prefix %q, want %q", g, w)
	}

	interestScalar := sdk.ZeroDec()
	if err := interestScalar.Unmarshal(response); err != nil {
		return err
	}

	denom := umeetypes.DenomFromKey(query.Request, umeetypes.KeyPrefixInterestScalar)
	data, interest, err := GetAndUnmarshalProtocolData[*types.UmeeInterestScalarProtocolData](ctx, k, denom, types.ProtocolDataTypeUmeeInterestScalar)
	if err != nil {
		return err
	}

	interest.Data, err = json.Marshal(interestScalar)
	if err != nil {
		return err
	}
	interest.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(interest)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, interest.GenerateKey(), &data)

	return nil
}

func UmeeUTokenSupplyUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	supplyAmount := sdk.ZeroInt()
	if err := supplyAmount.Unmarshal(response); err != nil {
		return err
	}

	if query.Request[0] != umeetypes.KeyPrefixUtokenSupply[0] {
		return errors.New("query request has unexpected prefix")
	}

	denom := umeetypes.DenomFromKey(query.Request, umeetypes.KeyPrefixUtokenSupply)
	data, supply, err := GetAndUnmarshalProtocolData[*types.UmeeUTokenSupplyProtocolData](ctx, k, denom, types.ProtocolDataTypeUmeeUTokenSupply)
	if err != nil {
		return err
	}
	supply.Data, err = json.Marshal(supplyAmount)
	if err != nil {
		return err
	}
	supply.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(supply)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, supply.GenerateKey(), &data)

	return nil
}

func UmeeLeverageModuleBalanceUpdateCallback(ctx sdk.Context, k *Keeper, response []byte, query icqtypes.Query) error {
	if len(query.Request) < 2 {
		k.Logger(ctx).Error("unable to unmarshal balance request, request length is too short")
		return errors.New("account balance icq request must always have a length of at least 2 bytes")
	}

	balancesStore := query.Request[1:]
	_, denom, err := banktypes.AddressAndDenomFromBalancesStore(balancesStore)
	if err != nil {
		return err
	}

	balanceCoin, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, response, denom)
	if err != nil {
		return err
	}
	balanceAmount := balanceCoin.Amount

	data, balance, err := GetAndUnmarshalProtocolData[*types.UmeeLeverageModuleBalanceProtocolData](ctx, k, denom, types.ProtocolDataTypeUmeeLeverageModuleBalance)
	if err != nil {
		return err
	}
	balance.Data, err = json.Marshal(balanceAmount)
	if err != nil {
		return err
	}
	balance.LastUpdated = ctx.BlockTime()
	data.Data, err = json.Marshal(balance)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, balance.GenerateKey(), &data)

	return nil
}

// SetEpochBlockCallback records the block height of the registered zone at the epoch boundary.
func SetEpochBlockCallback(ctx sdk.Context, k *Keeper, args []byte, query icqtypes.Query) error {
	k.Logger(ctx).Debug("epoch callback called")
	data, connectionData, err := GetAndUnmarshalProtocolData[*types.ConnectionProtocolData](ctx, k, query.ChainId, types.ProtocolDataTypeConnection)
	if err != nil {
		return err
	}

	// block response is never expected to be nil
	if len(args) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (1)")
	}

	blockResponse := tmservice.GetLatestBlockResponse{}
	err = k.cdc.Unmarshal(args, &blockResponse)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug("got block response", "block", blockResponse)

	if blockResponse.SdkBlock == nil {
		// v0.45 and below
		// nolint:staticcheck // SA1019 ignore this!
		connectionData.LastEpoch = blockResponse.Block.Header.Height
	} else {
		// v0.46 and above
		connectionData.LastEpoch = blockResponse.SdkBlock.Header.Height
	}

	heightInBytes := sdk.Uint64ToBigEndian(uint64(connectionData.LastEpoch)) //nolint:gosec
	// trigger a client update at the epoch boundary
	k.IcqKeeper.MakeRequest(
		ctx,
		query.ConnectionId,
		query.ChainId,
		"ibc.ClientUpdate",
		heightInBytes,
		sdk.NewInt(-1),
		types.ModuleName,
		"",
		0,
	)

	k.Logger(ctx).Debug("emitted client update", "height", connectionData.LastEpoch)

	data.Data, err = json.Marshal(connectionData)
	if err != nil {
		return err
	}
	k.SetProtocolData(ctx, connectionData.GenerateKey(), &data)
	return nil
}
