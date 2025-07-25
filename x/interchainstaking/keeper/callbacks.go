package keeper

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"time"

	"google.golang.org/protobuf/encoding/protowire"

	sdkioerrors "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmtypes "github.com/cometbft/cometbft/types"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	tmclienttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	lsmtypes "github.com/quicksilver-zone/quicksilver/third-party-chains/gaia-types/liquid/types"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/utils/proofs"
	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// ___________________________________________________________________________________________________

type Callback func(*Keeper, sdk.Context, []byte, icqtypes.Query) error

// Callbacks wrapper struct for interchainstaking keeper.
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
	return c.callbacks[id](c.k, ctx, args, query)
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
		AddCallback("valset", Callback(ValsetCallback)).
		AddCallback("validator", Callback(ValidatorCallback)).
		AddCallback("rewards", Callback(RewardsCallback)).
		AddCallback("delegations", Callback(DelegationsCallback)).
		AddCallback("delegations_epoch", Callback(DelegationsEpochCallback)).
		AddCallback("delegation", Callback(DelegationCallback)).
		AddCallback("delegation_epoch", Callback(DelegationEpochCallback)).
		AddCallback("distributerewards", Callback(DistributeRewardsFromWithdrawAccount)).
		AddCallback("depositinterval", Callback(DepositIntervalCallback)).
		AddCallback("deposittx", Callback(DepositTxCallback)).
		AddCallback("perfbalance", Callback(PerfBalanceCallback)).
		AddCallback("accountbalance", Callback(AccountBalanceCallback)).
		AddCallback("allbalances", Callback(AllBalancesCallback)).
		AddCallback("delegationaccountbalance", Callback(DelegationAccountBalanceCallback)).
		AddCallback("delegationaccountbalances", Callback(DelegationAccountBalancesCallback)).
		AddCallback("signinginfo", Callback(SigningInfoCallback)).
		AddCallback("lsminfo", Callback(LsmInfoCallback))

	return a.(Callbacks)
}

// -----------------------------------
// Callback Handlers
// -----------------------------------

func ValsetCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	return k.SetValidatorsForZone(ctx, args, query)
}

func ValidatorCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId()) // can't we get rid of this check?
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	return k.SetValidatorForZone(ctx, &zone, args)
}

func RewardsCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	k.Logger(ctx).Debug("rewards callback", "zone", query.ChainId)

	if len(query.Request) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (2)")
	}

	// unmarshal request payload
	rewardsQuery := distrtypes.QueryDelegationTotalRewardsRequest{}
	err := k.cdc.Unmarshal(query.Request, &rewardsQuery)
	if err != nil {
		return err
	}

	// decrement waitgroup as we have received back the query
	// (initially incremented in AfterEpochEnd)
	if err = zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), 1, "rewards callback"); err != nil {
		// given that there _could_ be a backlog of message, we don't want to bail here, else they will remain undeliverable.
		k.Logger(ctx).Error(err.Error())
	}

	k.Logger(ctx).Debug("QueryDelegationRewards callback", "wg", zone.GetWithdrawalWaitgroup(), "delegatorAddress", rewardsQuery.DelegatorAddress, "zone", query.ChainId)

	return k.WithdrawDelegationRewardsForResponse(ctx, &zone, rewardsQuery.DelegatorAddress, args)
}

func DelegationsEpochCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	return delegationsCallback(k, ctx, args, query, true)
}

func DelegationsCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	return delegationsCallback(k, ctx, args, query, false)
}

func delegationsCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query, isEpoch bool) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	if len(query.Request) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (3)")
	}

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{}
	err := k.cdc.Unmarshal(query.Request, &delegationQuery)
	if err != nil {
		return err
	}

	k.Logger(ctx).Debug("Delegations callback triggered", "chain", zone.ChainId)

	return k.UpdateDelegationRecordsForAddress(ctx, zone, delegationQuery.DelegatorAddr, args, isEpoch)
}

func DelegationEpochCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	return delegationCallback(k, ctx, args, query, true)
}

func DelegationCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	return delegationCallback(k, ctx, args, query, false)
}

func delegationCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query, isEpoch bool) error {
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

	k.Logger(ctx).Debug("Delegation callback", "delegation", delegation, "chain", zone.ChainId)

	if delegation.Shares.IsNil() || delegation.Shares.IsZero() {
		// delegation never gets removed, even with zero shares.
		delegator, validator, err := types.ParseStakingDelegationKey(query.Request)
		if err != nil {
			return err
		}
		validatorAddress, err := addressutils.EncodeAddressToBech32(zone.GetValoperPrefix(), validator)
		if err != nil {
			return err
		}
		delegatorAddress, err := addressutils.EncodeAddressToBech32(zone.GetAccountPrefix(), delegator)
		if err != nil {
			return err
		}

		if delegation, ok := k.GetDelegation(ctx, zone.ChainId, delegatorAddress, validatorAddress); ok {
			err := k.RemoveDelegation(ctx, zone.ChainId, delegation)
			if err != nil {
				return err
			}
		}
		return nil
	}
	valAddrBytes, err := addressutils.ValAddressFromBech32(delegation.ValidatorAddress, zone.GetValoperPrefix())
	if err != nil {
		return err
	}
	val, found := k.GetValidator(ctx, zone.ChainId, valAddrBytes)
	if !found {
		err := fmt.Errorf("unable to get validator: %s", delegation.ValidatorAddress)
		k.Logger(ctx).Error(err.Error())
		return err
	}

	return k.UpdateDelegationRecordForAddress(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress, sdk.NewCoin(zone.BaseDenom, val.SharesToTokens(delegation.Shares)), &zone, true, isEpoch)
}

func PerfBalanceCallback(k *Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	// update account balance first.
	if err := AccountBalanceCallback(k, ctx, response, query); err != nil {
		return err
	}

	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	// initialize performance delegations
	if err := k.UpdatePerformanceDelegations(ctx, zone); err != nil {
		k.Logger(ctx).Info(err.Error())
		return err
	}

	return nil
}

func DepositIntervalCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	if !zone.DepositsEnabled {
		return fmt.Errorf("chain id %s does not current allow deposits", query.GetChainId())
	}

	k.Logger(ctx).Debug("Deposit interval callback", "zone", zone.ChainId)

	if len(args) == 0 {
		// this is a no-op, but we still want to return success to avoid blocking icq submissions. Query will be re-raised if balance does not drop.
		k.Logger(ctx).Info("response was nil. balance suggests deposit, but corresponding tx not found by icq. check rpc node history.")
		return nil
	}

	txs := tx.GetTxsEventResponse{}
	err := k.cdc.Unmarshal(args, &txs)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal txs for deposit account", "deposit_address", zone.DepositAddress.GetAddress(), "err", err)
		return err
	}

	for _, txn := range txs.TxResponses {
		req := tx.GetTxRequest{Hash: txn.TxHash}
		hashBytes := k.cdc.MustMarshal(&req)
		_, found = k.GetReceipt(ctx, zone.ChainId, txn.TxHash)
		if found {
			k.Logger(ctx).Debug("Found previously handled tx. Ignoring.", "txhash", txn.TxHash)
			continue
		}
		k.Logger(ctx).Info("Found previously unhandled tx. Processing.", "txhash", txn.TxHash)
		k.ICQKeeper.MakeRequest(ctx, query.ConnectionId, query.ChainId, "tendermint.Tx", hashBytes, sdk.NewInt(-1), types.ModuleName, "deposittx", 0)
	}
	return nil
}

func SigningInfoCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", zone.ChainId)
	}

	k.Logger(ctx).Debug("Validator signing info callback", "zone", zone.ChainId)

	valSigningInfo := slashingtypes.ValidatorSigningInfo{}
	if len(args) == 0 {
		k.Logger(ctx).Error("unable to find signing info for validator", "query", query.Request)
		return nil
	}
	err := k.cdc.Unmarshal(args, &valSigningInfo)
	if err != nil {
		return err
	}
	if valSigningInfo.Tombstoned {
		consAddr, err := addressutils.AddressFromBech32(valSigningInfo.Address, "")
		if err != nil {
			return err
		}
		valAddr, found := k.GetValidatorAddrByConsAddr(ctx, zone.ChainId, consAddr)
		if !found {
			return fmt.Errorf("can not get validator address from consensus address: %s", valSigningInfo.Address)
		}

		k.Logger(ctx).Info("tombstoned validator found", "valoper", valAddr)

		valAddrBytes, err := addressutils.ValAddressFromBech32(valAddr, zone.GetValoperPrefix())
		if err != nil {
			return err
		}
		val, found := k.GetValidator(ctx, zone.ChainId, valAddrBytes)
		// NOTE: this shouldn't be reachable, but keeping here as it doesn't do any harm.
		if !found {
			err := k.SetValidator(ctx, zone.ChainId, types.Validator{
				ValoperAddress: valAddr,
				Jailed:         true,
				Tombstoned:     true,
			})
			if err != nil {
				return err
			}
		} else {
			val.Tombstoned = true
			if err = k.SetValidator(ctx, zone.ChainId, val); err != nil {
				return err
			}
		}
		k.Logger(ctx).Info(fmt.Sprintf("%q on chainID: %q was found to already have been tombstoned, added information", val.ValoperAddress, zone.ChainId))

	}
	return nil
}

func LsmInfoCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", zone.ChainId)
	}

	k.Logger(ctx).Debug("Validator liquid info callback", "zone", zone.ChainId)

	lsmValInfo := lsmtypes.LiquidValidator{}
	if len(args) == 0 {
		k.Logger(ctx).Error("unable to find liquid info for validator", "query", query.Request)
		return nil
	}
	err := k.cdc.Unmarshal(args, &lsmValInfo)
	if err != nil {
		return err
	}

	validatorAddr, err := addressutils.ValAddressFromBech32(lsmValInfo.OperatorAddress, zone.GetValoperPrefix())
	if err != nil {
		return err
	}
	validator, found := k.GetValidator(ctx, query.ChainId, validatorAddr)
	if !found {
		return fmt.Errorf("validator not found: %s", validatorAddr)
	}
	validator.LiquidShares = lsmValInfo.LiquidShares
	if err := k.SetValidator(ctx, query.ChainId, validator); err != nil {
		return err
	}
	k.Logger(ctx).Info("Liquid validator info updated", "validator", validator)
	return nil
}

// pulled directly from ibc-go tm light client
// checkTrustedHeader checks that consensus state matches trusted fields of Header.
func checkTrustedHeader(header *tmclienttypes.Header, consState *tmclienttypes.ConsensusState) error {
	tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	if err != nil {
		return sdkioerrors.Wrap(err, "trusted validator set in not tendermint validator set type")
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

// checkTMStateValidity checks if the Tendermint header is valid.
// CONTRACT: consState.Height == header.TrustedHeight
// pulled directly from ibc-go tm light client.
func checkTMStateValidity(
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
		return sdkioerrors.Wrap(err, "trusted validator set in not tendermint validator set type")
	}

	tmSignedHeader, err := tmtypes.SignedHeaderFromProto(header.SignedHeader)
	if err != nil {
		return sdkioerrors.Wrap(err, "signed header in not tendermint signed header type")
	}

	tmValidatorSet, err := tmtypes.ValidatorSetFromProto(header.ValidatorSet)
	if err != nil {
		return sdkioerrors.Wrap(err, "validator set in not tendermint validator set type")
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
	if header.TrustedHeight.RevisionHeight > math.MaxInt64 {
		return fmt.Errorf("trusted height revision exceeds int64 range: %d", header.TrustedHeight.RevisionHeight)
	}

	trustedHeader := tmtypes.Header{
		ChainID:            chainID,
		Height:             int64(header.TrustedHeight.RevisionHeight), //nolint:gosec
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
		return sdkioerrors.Wrap(err, "failed to verify header")
	}
	return nil
}

// CheckTMHeaderForZone verifies the Tendermint consensus and client states for a given zone. Returns error if unable
// to verify.
func (k *Keeper) CheckTMHeaderForZone(ctx sdk.Context, zone *types.Zone, header *tmclienttypes.Header) error {
	connection, _ := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, zone.ConnectionId)
	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return errors.New("unable to fetch client state")
	}
	/*
	   We can call ClientKeeper.CheckHeaderAndUpdateState() here, but this causes state changes inside the IBCKeeper
	   which feels bad. so instead we copy the above two functions wholesale from ibc-go (this sucks too, but with
	   predictable behaviour) and validate the inbound header manually.
	*/
	consensusState, found := k.IBCKeeper.ClientKeeper.GetClientConsensusState(ctx, connection.ClientId, header.TrustedHeight)
	if !found {
		return fmt.Errorf("unable to fetch consensus state for trusted height: %s", header.TrustedHeight.String())
	}

	tmclientState, ok := clientState.(*tmclienttypes.ClientState)
	if !ok {
		return errors.New("unable to marshal client state")
	}

	tmconsensusState, ok := consensusState.(*tmclienttypes.ConsensusState)
	if !ok {
		return errors.New("unable to marshal consensus state")
	}

	// validate tendermint statefor
	err := checkTMStateValidity(tmclientState, tmconsensusState, header, ctx.BlockHeader().Time)
	if err != nil {
		return fmt.Errorf("unable to validate header; %w", err)
	}

	return nil
}

// DepositTxCallback is a callback that verifies client chain state validity, gets Tx receipt and calls
// HandleReceiptForTransaction.
func DepositTxCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	// check validity
	if len(args) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (6)")
	}

	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	if !zone.DepositsEnabled {
		return fmt.Errorf("chain id %s does not current allow deposits", query.GetChainId())
	}

	k.Logger(ctx).Debug("DepositTx callback", "zone", zone.ChainId)

	res := icqtypes.GetTxWithProofResponse{}
	if err := k.cdc.Unmarshal(args, &res); err != nil {
		k.Logger(ctx).Error("Error unmarshalling args", "err", err)
		return err
	}

	var txBytes []byte

	var inclusionProof proofs.InclusionProof
	err := k.cdc.UnpackAny(res.ProofAny, &inclusionProof)
	if err != nil {
		return err
	}

	if inclusionProof == nil {
		return errors.New("inclusion proof is nil")
	}

	queryRequest := tx.GetTxRequest{}
	if err := k.cdc.Unmarshal(query.Request, &queryRequest); err != nil {
		return err
	}

	txBytes, err = inclusionProof.Validate(res.Header.Header.DataHash, queryRequest.Hash)
	if err != nil {
		return err
	}

	_, found = k.GetReceipt(ctx, zone.ChainId, queryRequest.Hash)
	if found {
		k.Logger(ctx).Info("Found previously handled tx. Ignoring.", "txhash", queryRequest.Hash)
		return nil
	}

	// check client state validity
	if err := k.CheckTMHeaderForZone(ctx, &zone, res.GetHeader()); err != nil {
		return err
	}

	txn, err := TxDecoder(k.cdc)(txBytes)
	if err != nil {
		return err
	}

	txtx, ok := txn.(*tx.Tx)
	if !ok {
		return errors.New("cannot assert type of tx")
	}
	return k.HandleReceiptTransaction(ctx, txtx, queryRequest.Hash, zone)
}

// AccountBalanceCallback is a callback handler for Balance queries.
func AccountBalanceCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
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
	// Please see https://github.com/quicksilver-zone/quicksilver-incognito/issues/80
	if err := coin.Validate(); err != nil {
		k.Logger(ctx).Error("invalid coin for zone", "zone", zone.ChainId, "err", err)
		return err
	}

	address, err := addressutils.EncodeAddressToBech32(zone.AccountPrefix, accAddr)
	if err != nil {
		return err
	}

	return k.SetAccountBalanceForDenom(ctx, &zone, address, coin)
}

// DelegationAccountBalanceCallback is a callback handler for Balance queries.
func DelegationAccountBalanceCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
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
	// Please see https://github.com/quicksilver-zone/quicksilver-incognito/issues/80
	if err := coin.Validate(); err != nil {
		k.Logger(ctx).Debug("invalid coin for zone", "zone", zone.ChainId, "err", err)
		return err
	}
	address, err := addressutils.EncodeAddressToBech32(zone.AccountPrefix, accAddr)
	if err != nil {
		return err
	}

	if zone.DelegationAddress == nil || address != zone.DelegationAddress.Address {
		k.Logger(ctx).Debug("delegation address does not match ")
		return err
	}

	k.Logger(ctx).Info("Received balance response for denom", "denom", coin.Denom)
	if err = zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), 1, "delegationaccountbalance callback"); err != nil {
		// given that there _could_ be a backlog of message, we don't want to bail here, else they will remain undeliverable.
		k.Logger(ctx).Error(err.Error())
	}

	// set the zone amount.
	balance := zone.DelegationAddress.Balance
	if ok, _ := zone.DelegationAddress.Balance.Find(coin.Denom); !ok {
		zone.DelegationAddress.Balance = zone.DelegationAddress.Balance.Add(coin)
	} else {
		for idx, i := range balance {
			if coin.Denom == i.Denom {
				zone.DelegationAddress.Balance[idx].Amount = coin.Amount
				break
			}
		}
	}

	k.SetZone(ctx, &zone)

	// if token is not valid for staking, then send to withdrawal account.
	if valid, _ := zone.ValidateCoinsForZone(sdk.NewCoins(coin), k.GetValidatorAddressesAsMap(ctx, zone.ChainId)); !valid {
		k.Logger(ctx).Info("token is not a valid staking token, so sending to withdrawal account for disbursal", "chain", zone.ChainId, "assets", coin)
		if zone.GetWithdrawalWaitgroup() == 0 {
			k.Logger(ctx).Info("triggering redemption rate calc in lieu of delegation flush")
			if err := k.TriggerRedemptionRate(ctx, &zone); err != nil {
				return err
			}
		}
		return k.SendToWithdrawal(ctx, &zone, zone.DelegationAddress, sdk.NewCoins(coin))
	}

	return k.FlushOutstandingDelegations(ctx, &zone, coin)
}

func DelegationAccountBalancesCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}
	result := banktypes.QueryAllBalancesResponse{}
	k.cdc.MustUnmarshal(args, &result)

	if err := zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), 1, "delegationaccountbalances callback"); err != nil {
		// given that there _could_ be a backlog of message, we don't want to bail here, else they will remain undeliverable.
		k.Logger(ctx).Error(err.Error())
	}

	addressBytes, err := addressutils.AccAddressFromBech32(zone.DelegationAddress.Address, zone.AccountPrefix)
	if err != nil {
		k.Logger(ctx).Error("cannot decode bech32 delegation addr")
		return err
	}
	balances := result.GetBalances().Sort()
	accountBalances := zone.DelegationAddress.Balance.Sort()

	for _, coin := range balances.Add(accountBalances...) { // we want to iterate over all denoms, including ones we currently have values for.

		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			types.BankStoreKey,
			append(banktypes.CreateAccountBalancesPrefix(addressBytes), coin.Denom...),
			sdk.NewInt(-1),
			types.ModuleName,
			"delegationaccountbalance",
			0,
		)

		if err = zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), 1, fmt.Sprintf("delegation account balance for %s", coin.Denom)); err != nil {
			return err
		}
		k.Logger(ctx).Info("Emitting balance request for denom", "denom", coin.Denom, "waitgroup", zone.GetWithdrawalWaitgroup())
	}
	k.SetZone(ctx, &zone)

	return nil
}

func AllBalancesCallback(k *Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	// this shouldn't happen because query.Request comes from Quicksilver
	if len(query.Request) == 0 {
		return errors.New("attempted to unmarshal zero length byte slice (7)")
	}

	balanceQuery := banktypes.QueryAllBalancesRequest{}
	err := k.cdc.Unmarshal(query.Request, &balanceQuery)
	if err != nil {
		return err
	}

	zone, found := k.GetZone(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	k.Logger(ctx).Debug("AllBalances callback", "chain", zone.ChainId)

	switch {
	case zone.DepositAddress != nil && balanceQuery.Address == zone.DepositAddress.Address:
		if zone.DepositAddress.BalanceWaitgroup != 0 {
			zone.DepositAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Info("zeroing deposit balance waitgroup")
		}
	case zone.WithdrawalAddress != nil && balanceQuery.Address == zone.WithdrawalAddress.Address:
		if zone.WithdrawalAddress.BalanceWaitgroup != 0 {
			zone.WithdrawalAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Info("zeroing withdrawal balance waitgroup")
		}
	case zone.DelegationAddress != nil && balanceQuery.Address == zone.DelegationAddress.Address:
		if zone.DelegationAddress.BalanceWaitgroup != 0 {
			zone.DelegationAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Info("zeroing delegation balance waitgroup")
		}
	case zone.PerformanceAddress != nil && balanceQuery.Address == zone.PerformanceAddress.Address:
		if zone.PerformanceAddress.BalanceWaitgroup != 0 {
			zone.PerformanceAddress.BalanceWaitgroup = 0
			k.Logger(ctx).Info("zeroing performance balance waitgroup")
		}
	}
	k.SetZone(ctx, &zone)

	return k.SetAccountBalance(ctx, zone, balanceQuery.Address, args)
}

// TxDecoder.
func TxDecoder(cdc codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		// Make sure txBytes follow ADR-027.
		err := RejectNonADR027TxRaw(txBytes)
		if err != nil {
			return nil, sdkioerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		var raw tx.TxRaw

		err = cdc.Unmarshal(txBytes, &raw)
		if err != nil {
			return nil, err
		}

		var body tx.TxBody

		err = cdc.Unmarshal(raw.BodyBytes, &body)
		if err != nil {
			return nil, sdkioerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		var authInfo tx.AuthInfo

		err = cdc.Unmarshal(raw.AuthInfoBytes, &authInfo)
		if err != nil {
			return nil, sdkioerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		return &tx.Tx{
			Body:       &body,
			AuthInfo:   &authInfo,
			Signatures: raw.Signatures,
		}, nil
	}
}

func RejectNonADR027TxRaw(txBytes []byte) error {
	// Make sure all fields are ordered in ascending order with this variable.
	prevTagNum := protowire.Number(0)

	for len(txBytes) > 0 {
		tagNum, wireType, m := protowire.ConsumeTag(txBytes)
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		// TxRaw only has bytes fields.
		if wireType != protowire.BytesType {
			return fmt.Errorf("expected %d wire type, got %d", protowire.BytesType, wireType)
		}
		// Make sure fields are ordered in ascending order.
		if tagNum < prevTagNum {
			return fmt.Errorf("txRaw must follow ADR-027, got tagNum %d after tagNum %d", tagNum, prevTagNum)
		}
		prevTagNum = tagNum

		// All 3 fields of TxRaw have wireType == 2, so their next component
		// is a varint, so we can safely call ConsumeVarint here.
		// Byte structure: <varint of bytes length><bytes sequence>
		// Inner  fields are verified in `DefaultTxDecoder`
		lengthPrefix, m := protowire.ConsumeVarint(txBytes[m:])
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		// We make sure that this varint is as short as possible.
		n := VarintMinLength(lengthPrefix)
		if n != m {
			return fmt.Errorf("length prefix varint for tagNum %d is not as short as possible, read %d, only need %d", tagNum, m, n)
		}

		// Skip over the bytes that store fieldNumber and wireType bytes.
		_, _, m = protowire.ConsumeField(txBytes)
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		txBytes = txBytes[m:]
	}

	return nil
}

// VarintMinLength returns the minimum number of bytes necessary to encode an
// uint using varint encoding.
func VarintMinLength(n uint64) int {
	switch {
	// Note: 1<<N == 2**N.
	case n < 1<<(7):
		return 1
	case n < 1<<(7*2):
		return 2
	case n < 1<<(7*3):
		return 3
	case n < 1<<(7*4):
		return 4
	case n < 1<<(7*5):
		return 5
	case n < 1<<(7*6):
		return 6
	case n < 1<<(7*7):
		return 7
	case n < 1<<(7*8):
		return 8
	case n < 1<<(7*9):
		return 9
	default:
		return 10
	}
}
