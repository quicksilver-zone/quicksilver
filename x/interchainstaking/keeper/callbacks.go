package keeper

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	sdkioerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	tmclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/utils"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	tmtypes "github.com/tendermint/tendermint/types"
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
		AddCallback("allbalances", Callback(AllBalancesCallback))

	return a.(Callbacks)
}

// -----------------------------------
// Callback Handlers
// -----------------------------------

func ValsetCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	return SetValidatorsForZone(&k, ctx, zone, args, query.Request)
}

func ValidatorCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	return SetValidatorForZone(&k, ctx, zone, args)
}

func RewardsCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	k.Logger(ctx).Info("rewards callback", "zone", query.ChainId)

	// unmarshal request payload
	rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{}
	if len(query.Request) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (2)")
	}
	err := k.cdc.Unmarshal(query.Request, &rewardsQuery)
	if err != nil {
		return err
	}

	// decrement waitgroup as we have received back the query
	// (initially incremented in AfterEpochEnd)
	zone.WithdrawalWaitgroup--

	k.Logger(ctx).Info("QueryDelegationRewards callback", "wg", zone.WithdrawalWaitgroup, "delegatorAddress", rewardsQuery.DelegatorAddress, "zone", query.ChainId)

	return k.WithdrawDelegationRewardsForResponse(ctx, &zone, rewardsQuery.DelegatorAddress, args)
}

func DelegationsCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{}
	if len(query.Request) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (3)")
	}
	err := k.cdc.Unmarshal(query.Request, &delegationQuery)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info("Delegations callback triggered", "chain", zone.ChainId)

	return k.UpdateDelegationRecordsForAddress(ctx, zone, delegationQuery.DelegatorAddr, args)
}

func DelegationCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	delegation := stakingtypes.Delegation{}
	// delegations _can_ legitimately be nil here, so explicitly DON'T guard against this.
	err := k.cdc.Unmarshal(args, &delegation)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info("Delegation callback", "delegation", delegation, "chain", zone.ChainId)

	if delegation.Shares.IsNil() || delegation.Shares.IsZero() {
		// delegation never gets removed, even with zero shares.
		delegator, validator, err := types.ParseStakingDelegationKey(query.Request)
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
		}
		return nil
	}
	val, found := zone.GetValidatorByValoper(delegation.ValidatorAddress)
	if !found {
		err := fmt.Errorf("unable to get validator: %s", delegation.ValidatorAddress)
		k.Logger(ctx).Error(err.Error())
		return err
	}

	return k.UpdateDelegationRecordForAddress(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress, sdk.NewCoin(zone.BaseDenom, val.SharesToTokens(delegation.Shares)), &zone, true)
}

func PerfBalanceCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	// update account balance first.
	if err := AccountBalanceCallback(k, ctx, response, query); err != nil {
		return err
	}

	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	// initialize performance delegations
	if err := k.UpdatePerformanceDelegations(ctx, zone, response); err != nil {
		k.Logger(ctx).Info(err.Error())
		return err
	}

	return nil
}

func DepositIntervalCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	if !zone.DepositsEnabled {
		return fmt.Errorf("chain id %s does not current allow deposits", query.GetChainId())
	}

	k.Logger(ctx).Info("Deposit interval callback", "zone", zone.ChainId)

	txs := tx.GetTxsEventResponse{}

	if len(args) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (4)")
	}
	err := k.cdc.Unmarshal(args, &txs)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal txs for deposit account", "deposit_address", zone.DepositAddress.GetAddress(), "err", err)
		return err
	}

	for _, txn := range txs.TxResponses {
		req := tx.GetTxRequest{Hash: txn.TxHash}
		hashBytes := k.cdc.MustMarshal(&req)
		_, found = k.GetReceipt(ctx, GetReceiptKey(zone.ChainId, txn.TxHash))
		if found {
			k.Logger(ctx).Info("Found previously handled tx. Ignoring.", "txhash", txn.TxHash)
			continue
		}
		k.ICQKeeper.MakeRequest(ctx, query.ConnectionId, query.ChainId, "tendermint.Tx", hashBytes, sdk.NewInt(-1), types.ModuleName, "deposittx", 0)
	}
	return nil
}

// pulled directly from ibc-go tm light client
// checkTrustedHeader checks that consensus state matches trusted fields of Header
func checkTrustedHeader(header *tmclienttypes.Header, consState *tmclienttypes.ConsensusState) error {
	tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return sdkioerrors.Wrapf(err, "trusted validator set in not tendermint validator set type")
	}

	// assert that trustedVals is NextValidators of last trusted header
	// to do this, we check that trustedVals.Hash() == consState.NextValidatorsHash
	tvalHash := tmTrustedValidators.Hash()
	if !bytes.Equal(consState.NextValidatorsHash, tvalHash) {
		return sdkioerrors.Wrapf(
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
		return sdkioerrors.Wrapf(
			tmclienttypes.ErrInvalidHeaderHeight,
			"header height revision %d does not match trusted header revision %d",
			header.GetHeight().GetRevisionNumber(), header.TrustedHeight.RevisionNumber,
		)
	}

	tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return sdkioerrors.Wrapf(err, "trusted validator set in not tendermint validator set type")
	}

	tmSignedHeader, err := tmtypes.SignedHeaderFromProto(header.SignedHeader)
	if err != nil {
		return sdkioerrors.Wrapf(err, "signed header in not tendermint signed header type")
	}

	tmValidatorSet, err := tmtypes.ValidatorSetFromProto(header.ValidatorSet)
	if err != nil {
		return sdkioerrors.Wrapf(err, "validator set in not tendermint validator set type")
	}

	// assert header height is newer than consensus state
	// if header.GetHeight().LTE(header.TrustedHeight) {
	// 	return sdkioerrors.Wrapf(
	// 		tmclienttypes.ErrInvalidHeader,
	// 		"header height ≤ consensus state height (%s ≤ %s)", header.GetHeight(), header.TrustedHeight,
	// 	)
	// }

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
	err = utils.Verify(
		&signedHeader,
		tmTrustedValidators, tmSignedHeader, tmValidatorSet,
		clientState.TrustingPeriod, currentTimestamp, clientState.MaxClockDrift, clientState.TrustLevel.ToTendermint(),
	)
	if err != nil {
		return sdkioerrors.Wrapf(err, "failed to verify header")
	}
	return nil
}

func DepositTx(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	if !zone.DepositsEnabled {
		return fmt.Errorf("chain id %s does not current allow deposits", query.GetChainId())
	}

	k.Logger(ctx).Info("DepositTx callback", "zone", zone.ChainId)

	res := icqtypes.GetTxWithProofResponse{}
	if len(args) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (6)")
	}
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
		return errors.New("unable to fetch client state")
	}

	/** we can call ClientKeeper.CheckHeaderAndUpdateState() here, but this causes state changes inside the IBCKeeper which feels bad.
	  so instead we copy the above two functions wholesale from ibc-go (this sucks too, but with predicatable behaviour) and validate
	  the inbound header manually. */
	consensusState, found := k.IBCKeeper.ClientKeeper.GetClientConsensusState(ctx, connection.ClientId, res.Header.TrustedHeight)
	if !found {
		return fmt.Errorf("unable to fetch consensus state for trusted height: %s", res.Header.TrustedHeight.String())
	}

	tmclientState, ok := clientState.(*tmclienttypes.ClientState)
	if !ok {
		return errors.New("unable to marshal client state")
	}

	tmconsensusState, ok := consensusState.(*tmclienttypes.ConsensusState)
	if !ok {
		return errors.New("unable to marshal consensus state")
	}

	err = checkValidity(tmclientState, tmconsensusState, res.GetHeader(), ctx.BlockHeader().Time)
	if err != nil {
		k.Logger(ctx).Info("unable to validate header", "header", res.Header)
		return fmt.Errorf("unable to validate header; %w", err)
	}

	tmproof, err := tmtypes.TxProofFromProto(*res.GetProof())
	if err != nil {
		return fmt.Errorf("unable to marshal proof: %w", err)
	}
	err = tmproof.Validate(res.Header.Header.DataHash)
	if err != nil {
		return fmt.Errorf("unable to validate proof: %w", err)
	}

	return k.HandleReceiptTransaction(ctx, res.GetTxResponse(), res.GetTx(), zone)
}

// AccountBalanceCallback is a callback handler for Balance queries.
func AccountBalanceCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	// strip the BalancesPrefix from the request key, as AddressFromBalancesStore expects this to be removed
	// by the prefixIterator. query.Request is a value that Quicksilver always sets, and is not user generated,
	// but lets us be safe here :)
	if len(query.Request) < 2 {
		k.Logger(ctx).Error("unable to unmarshal balance request", "zone", zone.ChainId, "error", "request length is too short")
		return errors.New("account balance icq request must always have a length of at least 2 bytes")
	}
	balancesStore := query.Request[1:]
	accAddr, denom, err := banktypes.AddressAndDenomFromBalancesStore(balancesStore)
	if err != nil {
		return err
	}

	coin, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, args, denom)
	if err != nil {
		return err
	}

	if coin.Denom != denom {
		return fmt.Errorf("received coin denom %s does not match requested denom %s", coin.Denom, denom)
	}

	// Ensure that the coin is valid.
	// Please see https://github.com/ingenuity-build/quicksilver-incognito/issues/80
	if err := coin.Validate(); err != nil {
		k.Logger(ctx).Error("invalid coin for zone", "zone", zone.ChainId, "err", err)
		return err
	}

	address, err := bech32.ConvertAndEncode(zone.AccountPrefix, accAddr)
	if err != nil {
		return err
	}

	return SetAccountBalanceForDenom(k, ctx, zone, address, coin)
}

func AllBalancesCallback(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	balanceQuery := banktypes.QueryAllBalancesRequest{}
	// this shouldn't happen because query.Request comes from Quicksilver
	if len(query.Request) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (7)")
	}
	err := k.cdc.Unmarshal(query.Request, &balanceQuery)
	if err != nil {
		return err
	}

	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	k.Logger(ctx).Info("AllBalances callback", "chain", zone.ChainId)

	switch {
	case zone.DepositAddress != nil && balanceQuery.Address == zone.DepositAddress.Address:
		if zone.DepositAddress.BalanceWaitgroup != 0 {
			zone.DepositAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Error("Zeroing deposit balance waitgroup")
		}
	case zone.WithdrawalAddress != nil && balanceQuery.Address == zone.WithdrawalAddress.Address:
		if zone.WithdrawalAddress.BalanceWaitgroup != 0 {
			zone.WithdrawalAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Error("Zeroing withdrawal balance waitgroup")
		}
	case zone.DelegationAddress != nil && balanceQuery.Address == zone.DelegationAddress.Address:
		if zone.DelegationAddress.BalanceWaitgroup != 0 {
			zone.DelegationAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Error("Zeroing delegation balance waitgroup")
		}
	case zone.PerformanceAddress != nil && balanceQuery.Address == zone.PerformanceAddress.Address:
		if zone.PerformanceAddress.BalanceWaitgroup != 0 {
			zone.PerformanceAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Error("Zeroing performance balance waitgroup")
		}
	}
	k.SetZone(ctx, &zone)

	return k.SetAccountBalance(ctx, zone, balanceQuery.Address, args)
}
