package keeper

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// ___________________________________________________________________________________________________

// Callbacks wrapper struct for interchainstaking keeper
type Callback func(Keeper, sdk.Context, []byte, icqtypes.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]Callback
}

var _ icqtypes.QueryCallbacks = Callbacks{}

func (k Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]Callback)}
}

//callback handler
func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id] = fn.(Callback)
	return c
}

func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	a := c.
		AddCallback("valset", Callback(ValsetCallback)).
		AddCallback("validator", Callback(ValidatorCallback)).
		AddCallback("rewards", Callback(RewardsCallback)).
		AddCallback("delegations", Callback(DelegationsCallback)).
		AddCallback("delegation", Callback(DelegationCallback)).
		AddCallback("distributerewards", Callback(DistributeRewardsFromWithdrawAccount)).
		AddCallback("depositinterval", Callback(DepositIntervalCallback)).
		AddCallback("perfbalance", Callback(PerfBalanceCallback)).
		AddCallback("accountbalance", Callback(AccountBalanceCallback)).
		AddCallback("allbalances", Callback(AllBalancesCallback))

	return a.(Callbacks)
}

// -----------------------------------
// Callback Handlers
// -----------------------------------

func ValsetCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	SetValidatorsForZone(k, ctx, zone, args)
	return nil
}

func ValidatorCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	k.Logger(ctx).Info("Received provable payload", "data", args)
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	SetValidatorForZone(k, ctx, zone, args)
	return nil
}

func RewardsCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	// unmarshal request payload
	rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{}
	err := k.cdc.Unmarshal(query.Request, &rewardsQuery)
	if err != nil {
		return err
	}

	// decrement waitgroup as we have received back the query (initially incremented in L93).
	zone.WithdrawalWaitgroup--

	k.Logger(ctx).Info("QueryDelegationRewards callback", "wg", zone.WithdrawalWaitgroup, "delegatorAddress", rewardsQuery.DelegatorAddress)

	return k.WithdrawDelegationRewardsForResponse(ctx, &zone, rewardsQuery.DelegatorAddress, args)
}

func DelegationsCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{}
	err := k.cdc.Unmarshal(query.Request, &delegationQuery)
	if err != nil {
		return err
	}

	return k.UpdateDelegationRecordsForAddress(ctx, &zone, delegationQuery.DelegatorAddr, args)
}

func DelegationCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	delegation := stakingtypes.Delegation{}
	err := k.cdc.Unmarshal(args, &delegation)
	if err != nil {
		return err
	}

	if delegation.Shares.IsNil() || delegation.Shares.IsZero() {
		// delegation never gets removed, even with zero shares.
		delegator, validator, err := parseDelegationKey(query.Request)
		if err != nil {
			return err
		}
		validatorAddress, err := bech32.ConvertAndEncode(zone.GetAccountPrefix()+"valoper", validator)
		if err != nil {
			return err
		}
		delegatorAddress, err := bech32.ConvertAndEncode(zone.GetAccountPrefix(), delegator)
		if err != nil {
			return err
		}
		if delegation, ok := k.GetDelegation(ctx, &zone, delegatorAddress, validatorAddress); ok {
			k.RemoveDelegation(ctx, &zone, delegation)
			ica, err := zone.GetDelegationAccountByAddress(delegatorAddress)
			if err != nil {
				return err
			}
			ica.DelegatedBalance = ica.DelegatedBalance.Sub(delegation.Amount)
			k.SetRegisteredZone(ctx, zone)
		}
		return nil
	}
	val, err := zone.GetValidatorByValoper(delegation.ValidatorAddress)
	if err != nil {
		k.Logger(ctx).Error("unable to get validator", "address", delegation.ValidatorAddress)
		return err
	}

	return k.UpdateDelegationRecordForAddress(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress, sdk.NewCoin(zone.BaseDenom, val.SharesToTokens(delegation.Shares)), &zone, true)
}

func PerfBalanceCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	// initialize performance delegations
	if err := k.InitPerformanceDelegations(ctx, zone, response); err != nil {
		k.Logger(ctx).Info(err.Error())
		return err
	}

	return nil
}

func DepositIntervalCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	txs := tx.GetTxsEventResponse{}

	err := k.cdc.Unmarshal(args, &txs)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal txs for deposit account", "deposit_address", zone.DepositAddress.GetAddress(), "err", err)
		return err
	}

	// TODO: use pagination.GetTotal() to dispatch the correct number of requests now; rather than iteratively.
	if len(txs.GetTxs()) == types.TxRetrieveCount {
		req := tx.GetTxsEventRequest{}
		err := k.cdc.Unmarshal(query.Request, &req)
		if err != nil {
			return err
		}
		req.Pagination.Offset = req.Pagination.Offset + req.Pagination.Limit

		k.ICQKeeper.MakeRequest(ctx, query.ConnectionId, query.ChainId, "cosmos.tx.v1beta1.Service/GetTxsEvent", k.cdc.MustMarshal(&req), sdk.NewInt(-1), types.ModuleName, "depositinterval", 0)
	}

	for i, tx := range txs.TxResponses {
		k.HandleReceiptTransaction(ctx, tx, txs.Txs[i], zone)
	}
	return nil
}

// setAccountCb is a callback handler for Balance queries.
func AccountBalanceCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	balancesStore := []byte(query.Request[1:])
	accAddr, _, err := banktypes.AddressAndDenomFromBalancesStore(balancesStore)
	if err != nil {
		return err
	}

	coin := sdk.Coin{}
	err = k.cdc.Unmarshal(args, &coin)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal balance info for zone", "zone", zone.ChainId, "err", err)
		return err
	}

	if coin.IsNil() {
		denom := ""

		for i := 0; i < len(query.Request)-len(accAddr); i++ {
			if bytes.Equal(query.Request[i:i+len(accAddr)], accAddr) {
				denom = string(query.Request[i+len(accAddr):])
				break
			}

		}
		// if balance is nil, the response sent back is nil, so we don't receive the denom. Override that now.
		coin = sdk.NewCoin(denom, sdk.ZeroInt())
	}

	address, err := bech32.ConvertAndEncode(zone.AccountPrefix, accAddr)
	if err != nil {
		return err
	}

	return SetAccountBalanceForDenom(k, ctx, zone, address, coin)
}

func AllBalancesCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {

	balanceQuery := bankTypes.QueryAllBalancesRequest{}
	err := k.cdc.Unmarshal(query.Request, &balanceQuery)
	if err != nil {
		return err
	}

	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	//
	if zone.DepositAddress.BalanceWaitgroup != 0 {
		zone.DepositAddress.BalanceWaitgroup = 0
		k.Logger(ctx).Error("Zeroing deposit balance waitgroup")
		k.SetRegisteredZone(ctx, zone)
	}

	return k.SetAccountBalance(ctx, zone, balanceQuery.Address, args)
}
