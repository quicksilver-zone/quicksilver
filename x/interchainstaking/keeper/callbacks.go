package keeper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	tmclienttypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/light"
	tmtypes "github.com/tendermint/tendermint/types"

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

// callback handler
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
		AddCallback("deposittx", Callback(DepositTx)).
		AddCallback("perfbalance", Callback(PerfBalanceCallback)).
		AddCallback("accountbalance", Callback(AccountBalanceCallback)).
		AddCallback("allbalances", Callback(AllBalancesCallback)).
		AddCallback("epochblock", Callback(SetEpochBlockCallback))

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
	return SetValidatorsForZone(k, ctx, zone, args)
}

// SetEpochBlockCallback records the block height of the registered zone at the epoch boundary.
func SetEpochBlockCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	blockResponse := tmservice.GetLatestBlockResponse{}
	err := k.cdc.Unmarshal(args, &blockResponse)
	if err != nil {
		return err
	}
	zone.LastEpochHeight = blockResponse.Block.Header.Height
	k.SetRegisteredZone(ctx, zone)
	return nil
}

func ValidatorCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	k.Logger(ctx).Info("Received provable payload", "data", args)
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	return SetValidatorForZone(k, ctx, zone, args)
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
			err := k.RemoveDelegation(ctx, &zone, delegation)
			if err != nil {
				return err
			}

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
		req.Pagination.Offset += req.Pagination.Limit

		k.ICQKeeper.MakeRequest(ctx, query.ConnectionId, query.ChainId, "cosmos.tx.v1beta1.Service/GetTxsEvent", k.cdc.MustMarshal(&req), sdk.NewInt(-1), types.ModuleName, "depositinterval", 0)
	}

	for _, txn := range txs.TxResponses {

		req := tx.GetTxRequest{Hash: txn.TxHash}
		hashBytes := k.cdc.MustMarshal(&req)
		k.ICQKeeper.MakeRequest(ctx, query.ConnectionId, query.ChainId, "tendermint.Tx", hashBytes, sdk.NewInt(-1), types.ModuleName, "deposittx", 0)

	}
	return nil
}

// pulled directly from ibc-go tm light client
// checkTrustedHeader checks that consensus state matches trusted fields of Header
func checkTrustedHeader(header *tmclienttypes.Header, consState *tmclienttypes.ConsensusState) error {
	tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return sdkerrors.Wrap(err, "trusted validator set in not tendermint validator set type")
	}

	// assert that trustedVals is NextValidators of last trusted header
	// to do this, we check that trustedVals.Hash() == consState.NextValidatorsHash
	tvalHash := tmTrustedValidators.Hash()
	if !bytes.Equal(consState.NextValidatorsHash, tvalHash) {
		return sdkerrors.Wrapf(
			tmclienttypes.ErrInvalidValidatorSet,
			"trusted validators %s, does not hash to latest trusted validators. Expected: %X, got: %X",
			header.TrustedValidators, consState.NextValidatorsHash, tvalHash,
		)
	}
	return nil
}

// pulled directly from ibc-go tm light client
// checkValidity checks if the Tendermint header is valid.
// CONTRACT: consState.Height == header.TrustedHeight
func checkValidity(
	clientState *tmclienttypes.ClientState, consState *tmclienttypes.ConsensusState,
	header *tmclienttypes.Header, currentTimestamp time.Time,
) error {
	if err := checkTrustedHeader(header, consState); err != nil {
		return err
	}

	// UpdateClient only accepts updates with a header at the same revision
	// as the trusted consensus state
	if header.GetHeight().GetRevisionNumber() != header.TrustedHeight.RevisionNumber {
		return sdkerrors.Wrapf(
			tmclienttypes.ErrInvalidHeaderHeight,
			"header height revision %d does not match trusted header revision %d",
			header.GetHeight().GetRevisionNumber(), header.TrustedHeight.RevisionNumber,
		)
	}

	tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return sdkerrors.Wrap(err, "trusted validator set in not tendermint validator set type")
	}

	tmSignedHeader, err := tmtypes.SignedHeaderFromProto(header.SignedHeader)
	if err != nil {
		return sdkerrors.Wrap(err, "signed header in not tendermint signed header type")
	}

	tmValidatorSet, err := tmtypes.ValidatorSetFromProto(header.ValidatorSet)
	if err != nil {
		return sdkerrors.Wrap(err, "validator set in not tendermint validator set type")
	}

	// assert header height is newer than consensus state
	if header.GetHeight().LTE(header.TrustedHeight) {
		return sdkerrors.Wrapf(
			tmclienttypes.ErrInvalidHeader,
			"header height ≤ consensus state height (%s ≤ %s)", header.GetHeight(), header.TrustedHeight,
		)
	}

	chainID := clientState.GetChainID()
	// If chainID is in revision format, then set revision number of chainID with the revision number
	// of the header we are verifying
	// This is useful if the update is at a previous revision rather than an update to the latest revision
	// of the client.
	// The chainID must be set correctly for the previous revision before attempting verification.
	// Updates for previous revisions are not supported if the chainID is not in revision format.
	if clienttypes.IsRevisionFormat(chainID) {
		chainID, _ = clienttypes.SetRevisionNumber(chainID, header.GetHeight().GetRevisionNumber())
	}

	// Construct a trusted header using the fields in consensus state
	// Only Height, Time, and NextValidatorsHash are necessary for verification
	trustedHeader := tmtypes.Header{
		ChainID:            chainID,
		Height:             int64(header.TrustedHeight.RevisionHeight),
		Time:               consState.Timestamp,
		NextValidatorsHash: consState.NextValidatorsHash,
	}
	signedHeader := tmtypes.SignedHeader{
		Header: &trustedHeader,
	}

	// Verify next header with the passed-in trustedVals
	// - asserts trusting period not passed
	// - assert header timestamp is not past the trusting period
	// - assert header timestamp is past latest stored consensus state timestamp
	// - assert that a TrustLevel proportion of TrustedValidators signed new Commit
	err = light.Verify(
		&signedHeader,
		tmTrustedValidators, tmSignedHeader, tmValidatorSet,
		clientState.TrustingPeriod, currentTimestamp, clientState.MaxClockDrift, clientState.TrustLevel.ToTendermint(),
	)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to verify header")
	}
	return nil
}

func DepositTx(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	res := icqtypes.GetTxWithProofResponse{}
	err := k.cdc.Unmarshal(args, &res)
	if err != nil {
		return err
	}

	_, found = k.GetReceipt(ctx, GetReceiptKey(zone.ChainId, res.GetTxResponse().TxHash))
	if found {
		k.Logger(ctx).Info("Found previously handled tx. Ignoring.", "txhash", res.GetTxResponse().TxHash)
		return nil
	}

	// validate proof
	connection, _ := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, zone.ConnectionId)

	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return fmt.Errorf("unable to fetch client state")
	}

	/** we can call ClientKeeper.CheckHeaderAndUpdateState() here, but this causes state changes inside the IBCKeeper which feels bad.
	  so instead we copy the above two functions wholesale from ibc-go (this sucks too, but with predicatable behaviour) and validate
	  the inbound header manually. */
	consensusState, found := k.IBCKeeper.ClientKeeper.GetClientConsensusState(ctx, connection.ClientId, res.Header.TrustedHeight)
	if !found {
		return fmt.Errorf("unable to fetch consensus state")
	}

	tmclientState, ok := clientState.(*tmclienttypes.ClientState)
	if !ok {
		return fmt.Errorf("unable to marshal client state")
	}

	tmconsensusState, ok := consensusState.(*tmclienttypes.ConsensusState)
	if !ok {
		return fmt.Errorf("unable to marshal consensus state")
	}

	err = checkValidity(tmclientState, tmconsensusState, res.GetHeader(), ctx.BlockHeader().Time)
	if err != nil {
		k.Logger(ctx).Info("unable to validate header", "header", res.Header)
	}

	// _, _, err = clientState.CheckHeaderAndUpdateState(ctx, k.cdc, k.IBCKeeper.ClientKeeper.ClientStore(ctx, connection.ClientId), res.GetHeader())
	// if err != nil {
	// 	k.Logger(ctx).Info("Invalid header", "datahash", hex.EncodeToString(res.Header.Header.DataHash), "err", err)
	// }

	tmproof, err := tmtypes.TxProofFromProto(*res.GetProof())
	if err != nil {
		return fmt.Errorf("unable to marshal proof: %s", err)
	}
	// k.Logger(ctx).Error("hashes", "proof", tmproof.RootHash, "header", hex.EncodeToString(res.Header.Header.DataHash))
	err = tmproof.Validate(res.Header.Header.DataHash)
	if err != nil {
		return fmt.Errorf("unable to validate proof: %s", err)
	}

	k.HandleReceiptTransaction(ctx, res.GetTxResponse(), res.GetTx(), zone)
	return nil
}

// setAccountCb is a callback handler for Balance queries.
func AccountBalanceCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	balancesStore := query.Request[1:]
	accAddr, err := banktypes.AddressFromBalancesStore(balancesStore)
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
		if err := sdk.ValidateDenom(denom); err != nil {
			return err
		}
		coin = sdk.NewCoin(denom, sdk.ZeroInt())
	}

	address, err := bech32.ConvertAndEncode(zone.AccountPrefix, accAddr)
	if err != nil {
		return err
	}

	return SetAccountBalanceForDenom(k, ctx, zone, address, coin)
}

func AllBalancesCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	balanceQuery := banktypes.QueryAllBalancesRequest{}
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
