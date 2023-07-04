package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ingenuity-build/quicksilver/osmosis-types/gamm"
	umeetypes "github.com/ingenuity-build/quicksilver/umee-types/leverage/types"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

const (
	ValidatorSelectionRewardsCallbackID       = "validatorselectionrewards"
	OsmosisPoolUpdateCallbackID               = "osmosispoolupdate"
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
	var pd gamm.PoolI
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
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisPool, fmt.Sprintf("%d", poolID))
	if !ok {
		return fmt.Errorf("unable to find protocol data for osmosispools/%d", poolID)
	}
	ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, data.Data)
	if err != nil {
		return err
	}
	pool, ok := ipool.(*types.OsmosisPoolProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for osmosispools/%d", poolID)
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
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeUmeeReserves, denom)
	if !ok {
		return fmt.Errorf("unable to find protocol data for umeereserves/%s", denom)
	}
	ireserves, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeReserves, data.Data)
	if err != nil {
		return err
	}
	reserves, ok := ireserves.(*types.UmeeReservesProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for umeereserves/%s", denom)
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
	totalBorrows := sdk.ZeroDec()
	if err := totalBorrows.Unmarshal(response); err != nil {
		return err
	}

	if query.Request[0] != umeetypes.KeyPrefixAdjustedTotalBorrow[0] {
		return errors.New("query request has unexpected prefix")
	}

	denom := umeetypes.DenomFromKey(query.Request, umeetypes.KeyPrefixAdjustedTotalBorrow)
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeUmeeTotalBorrows, denom)
	if !ok {
		return fmt.Errorf("unable to find protocol data for umee-types total borrows/%s", denom)
	}
	iborrows, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeTotalBorrows, data.Data)
	if err != nil {
		return err
	}
	borrows, ok := iborrows.(*types.UmeeTotalBorrowsProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for umee-types total borrows/%s", denom)
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
	interestScalar := sdk.ZeroDec()
	if err := interestScalar.Unmarshal(response); err != nil {
		return err
	}

	if query.Request[0] != umeetypes.KeyPrefixInterestScalar[0] {
		return errors.New("query request has unexpected prefix")
	}

	denom := umeetypes.DenomFromKey(query.Request, umeetypes.KeyPrefixInterestScalar)
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeUmeeInterestScalar, denom)
	if !ok {
		return fmt.Errorf("unable to find protocol data for interestscalar/%s", denom)
	}
	iinterest, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeInterestScalar, data.Data)
	if err != nil {
		return err
	}
	interest, ok := iinterest.(*types.UmeeInterestScalarProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for interestscalar/%s", denom)
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
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeUmeeUTokenSupply, denom)
	if !ok {
		return fmt.Errorf("unable to find protocol data for umee-types utoken supply/%s", denom)
	}
	isupply, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeUTokenSupply, data.Data)
	if err != nil {
		return err
	}
	supply, ok := isupply.(*types.UmeeUTokenSupplyProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for umee-types utoken supply/%s", denom)
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

	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeUmeeLeverageModuleBalance, denom)
	if !ok {
		return fmt.Errorf("unable to find protocol data for umee-types leverage module/%s", denom)
	}
	ibalance, err := types.UnmarshalProtocolData(types.ProtocolDataTypeUmeeLeverageModuleBalance, data.Data)
	if err != nil {
		return err
	}
	balance, ok := ibalance.(*types.UmeeLeverageModuleBalanceProtocolData)
	if !ok {
		return fmt.Errorf("unable to unmarshal protocol data for umee-types leverage module/%s", denom)
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
	data, ok := k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, query.ChainId)
	if !ok {
		return fmt.Errorf("unable to find protocol data for connection/%s", query.ChainId)
	}
	k.Logger(ctx).Debug("epoch callback called")
	iConnectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, data.Data)
	connectionData, _ := iConnectionData.(*types.ConnectionProtocolData)

	if err != nil {
		return err
	}

	blockResponse := tmservice.GetLatestBlockResponse{}
	// block response is never expected to be nil
	if len(args) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (1)")
	}
	err = k.cdc.Unmarshal(args, &blockResponse)
	if err != nil {
		return err
	}
	k.Logger(ctx).Debug("got block response", "block", blockResponse)

	if blockResponse.SdkBlock == nil {
		// v0.45 and below
		//nolint:staticcheck // SA1019 ignore this!
		connectionData.LastEpoch = blockResponse.Block.Header.Height
	} else {
		// v0.46 and above
		connectionData.LastEpoch = blockResponse.SdkBlock.Header.Height
	}

	heightInBytes := sdk.Uint64ToBigEndian(uint64(connectionData.LastEpoch))
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
