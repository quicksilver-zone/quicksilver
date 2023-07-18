package keeper

import (
	"errors"
	"fmt"
	"github.com/ingenuity-build/quicksilver/third-party-chains/osmosis-types"
	osmosislockuptypes "github.com/ingenuity-build/quicksilver/third-party-chains/osmosis-types/lockup"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
	cmtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var (
	tier1 = "0.05"
	tier2 = "0.10"
	tier3 = "0.15"
	tier4 = "0.22"
	tier5 = "0.30"
)

func (k *Keeper) HandleClaim(ctx sdk.Context, cr types.ClaimRecord, action types.Action, proofs []*cmtypes.Proof) (uint64, error) {
	// action already completed, nothing to claim
	if _, exists := cr.ActionsCompleted[int32(action)]; exists {
		return 0, fmt.Errorf("%s already completed", types.Action_name[int32(action)])
	}

	switch action {
	case types.ActionInitialClaim:
		return k.handleInitial(ctx, &cr, action)
	case types.ActionDepositT1:
		return k.handleDeposit(ctx, &cr, action, sdk.MustNewDecFromStr(tier1))
	case types.ActionDepositT2:
		return k.handleDeposit(ctx, &cr, action, sdk.MustNewDecFromStr(tier2))
	case types.ActionDepositT3:
		return k.handleDeposit(ctx, &cr, action, sdk.MustNewDecFromStr(tier3))
	case types.ActionDepositT4:
		return k.handleDeposit(ctx, &cr, action, sdk.MustNewDecFromStr(tier4))
	case types.ActionDepositT5:
		return k.handleDeposit(ctx, &cr, action, sdk.MustNewDecFromStr(tier5))
	case types.ActionStakeQCK:
		return k.handleBondedDelegation(ctx, &cr, action)
	case types.ActionSignalIntent:
		return k.handleZoneIntent(ctx, &cr, action)
	case types.ActionQSGov:
		return k.handleGovernanceParticipation(ctx, &cr, action)
	case types.ActionGbP:
		// TODO: implement handler once GbP is implemented
	case types.ActionOsmosis:
		return k.handleOsmosisLP(ctx, &cr, action, proofs)
	default:
		return 0, fmt.Errorf("undefined action [%d]", action)
	}

	return 0, fmt.Errorf("handler not implemented for [%d] %s", action, types.Action_name[int32(action)])
}

// ------------
// # Handlers #
// ------------

// handleInitial.
func (k *Keeper) handleInitial(ctx sdk.Context, cr *types.ClaimRecord, action types.Action) (uint64, error) {
	return k.completeClaim(ctx, cr, action)
}

// handleDeposit.
func (k *Keeper) handleDeposit(ctx sdk.Context, cr *types.ClaimRecord, action types.Action, threshold sdk.Dec) (uint64, error) {
	if err := k.verifyDeposit(ctx, *cr, threshold); err != nil {
		return 0, err
	}

	return k.completeClaim(ctx, cr, action)
}

// handleBondedDelegation.
func (k *Keeper) handleBondedDelegation(ctx sdk.Context, cr *types.ClaimRecord, action types.Action) (uint64, error) {
	if err := k.verifyBondedDelegation(ctx, cr.Address); err != nil {
		return 0, err
	}

	return k.completeClaim(ctx, cr, action)
}

// handleZoneIntent.
func (k *Keeper) handleZoneIntent(ctx sdk.Context, cr *types.ClaimRecord, action types.Action) (uint64, error) {
	if err := k.verifyZoneIntent(ctx, cr.ChainId, cr.Address); err != nil {
		return 0, err
	}

	return k.completeClaim(ctx, cr, action)
}

// handleZoneIntent.
func (k *Keeper) handleGovernanceParticipation(ctx sdk.Context, cr *types.ClaimRecord, action types.Action) (uint64, error) {
	if err := k.verifyGovernanceParticipation(ctx, cr.Address); err != nil {
		return 0, err
	}

	return k.completeClaim(ctx, cr, action)
}

// handleOsmosisLP.
func (k *Keeper) handleOsmosisLP(ctx sdk.Context, cr *types.ClaimRecord, action types.Action, proofs []*cmtypes.Proof) (uint64, error) {
	if len(proofs) == 0 {
		return 0, errors.New("expects at least one LP proof")
	}
	if err := k.verifyOsmosisLP(ctx, proofs, *cr); err != nil {
		return 0, err
	}

	return k.completeClaim(ctx, cr, action)
}

// -------------
// # Verifiers #
// -------------

// verifyDeposit.
func (k *Keeper) verifyDeposit(ctx sdk.Context, cr types.ClaimRecord, threshold sdk.Dec) error {
	addr, err := sdk.AccAddressFromBech32(cr.Address)
	if err != nil {
		return err
	}

	zone, ok := k.icsKeeper.GetZone(ctx, cr.ChainId)
	if !ok {
		return fmt.Errorf("zone not found for %s", cr.ChainId)
	}

	// obtain all deposit receipts for this user on this zone
	receipts, err := k.icsKeeper.UserZoneReceipts(ctx, &zone, addr)
	if err != nil {
		return fmt.Errorf("unable to obtain zone receipts for %s on zone %s: %w", cr.Address, cr.ChainId, err)
	}

	// sum gross deposits amount
	gdAmount := sdk.NewInt(0)
	for _, rcpt := range receipts {
		gdAmount = gdAmount.Add(rcpt.Amount.AmountOf(zone.BaseDenom))
	}

	// calculate target amount
	tAmount := threshold.MulInt64(int64(cr.BaseValue)).TruncateInt()

	if gdAmount.LT(tAmount) {
		return fmt.Errorf("insufficient deposit amount, expects %v got %v", tAmount, gdAmount)
	}

	return nil
}

// verifyBondedDelegation indicates if the given address has an active bonded.
// delegation of QCK on the Quicksilver zone.
func (k *Keeper) verifyBondedDelegation(ctx sdk.Context, address string) error {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	amount := k.stakingKeeper.GetDelegatorBonded(ctx, addr)
	if !amount.IsPositive() {
		return fmt.Errorf("no bonded delegation for %s", addr)
	}
	return nil
}

// verifyZoneIntent indicates if the given address has intent set for the given
// zone (chainID).
func (k *Keeper) verifyZoneIntent(ctx sdk.Context, chainID, address string) error {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	zone, ok := k.icsKeeper.GetZone(ctx, chainID)
	if !ok {
		return fmt.Errorf("zone %s not found", chainID)
	}

	intent, ok := k.icsKeeper.GetDelegatorIntent(ctx, &zone, addr.String(), false)
	if !ok || len(intent.Intents) == 0 {
		return fmt.Errorf("intent not found or no intents set for %s", addr)
	}

	return nil
}

// verifyGovernanceParticipation indicates if the given address has voted on
// any governance proposals on the Quicksilver zone.
func (k *Keeper) verifyGovernanceParticipation(ctx sdk.Context, address string) error {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	voted := false
	k.govKeeper.IterateProposals(ctx, func(proposal govv1.Proposal) (stop bool) {
		_, found := k.govKeeper.GetVote(ctx, proposal.Id, addr)
		if found {
			voted = true
			return true
		}
		return false
	})

	if !voted {
		return fmt.Errorf("no governance votes by %s", addr)
	}

	return nil
}

// verifyOsmosisLP utilizes cross-chain-verification (XCV) to indicate if the
// given address provides any liquidity of the zones qAssets on the Osmosis
// chain.
//
// It utilizes Osmosis query:
//
//	rpc LockedByID(LockedRequest) returns (LockedResponse);
func (k *Keeper) verifyOsmosisLP(ctx sdk.Context, proofs []*cmtypes.Proof, cr types.ClaimRecord) error {
	// get Osmosis zone
	var osmoZone *icstypes.Zone
	k.icsKeeper.IterateZones(ctx, func(_ int64, zone *icstypes.Zone) (stop bool) {
		if zone.AccountPrefix == "osmo" {
			osmoZone = zone
			return true
		}
		return false
	})
	if osmoZone == nil {
		return errors.New("unable to find Osmosis zone")
	}

	uAmount := sdk.ZeroInt()
	dupCheck := make(map[string]struct{})
	for i, p := range proofs {
		proof := p

		// check for duplicate proof submission
		if _, exists := dupCheck[string(proof.Key)]; exists {
			return fmt.Errorf("duplicate proof submitted, %s", proof.Key)
		}
		dupCheck[string(proof.Key)] = struct{}{}

		// validate proof tx
		if err := k.ValidateProofOps(
			ctx,
			&k.icsKeeper.IBCKeeper,
			osmoZone.ConnectionId,
			osmoZone.ChainId,
			proof.Height,
			proof.ProofType,
			proof.Key,
			proof.Data,
			proof.ProofOps,
		); err != nil {
			return fmt.Errorf("proofs [%d]: %w", i, err)
		}

		var lock osmosislockuptypes.PeriodLock
		err := k.cdc.Unmarshal(proof.Data, &lock)
		if err != nil {
			return fmt.Errorf("unable to unmarshal locked response: %w", err)
		}

		// verify proof lock owner address is claim record address
		if lock.Owner != cr.Address {
			return fmt.Errorf("invalid lock owner, expected %s got %s", cr.Address, lock.Owner)
		}

		// verify pool is for the relevant zone
		// and sum user amounts
		amount, err := k.verifyPoolAndGetAmount(ctx, lock, cr)
		if err != nil {
			return err
		}
		uAmount = uAmount.Add(amount)
	}

	// calculate target amount
	dThreshold := sdk.MustNewDecFromStr(tier4)
	if err := k.verifyDeposit(ctx, cr, dThreshold); err != nil {
		return fmt.Errorf("%w, must reach at least %s of %d", err, tier4, cr.BaseValue)
	}
	tAmount := dThreshold.MulInt64(int64(cr.BaseValue / 2)).TruncateInt()

	// check liquidity threshold
	if uAmount.LT(tAmount) {
		return fmt.Errorf("insufficient liquidity, expects at least %s, got %s", tAmount, uAmount)
	}

	return nil
}

func (k *Keeper) verifyPoolAndGetAmount(ctx sdk.Context, lock osmosislockuptypes.PeriodLock, cr types.ClaimRecord) (sdkmath.Int, error) {
	return osmosistypes.DetermineApplicableTokensInPool(ctx, k.prKeeper, lock, cr.ChainId)
}

// -----------
// # Helpers #
// -----------

func (k *Keeper) completeClaim(ctx sdk.Context, cr *types.ClaimRecord, action types.Action) (uint64, error) {
	// update ClaimRecord and obtain total claim amount
	claimAmount, err := k.getClaimAmountAndUpdateRecord(ctx, cr, action)
	if err != nil {
		return 0, err
	}

	// send coins to address
	coins, err := k.sendCoins(ctx, *cr, claimAmount)
	if err != nil {
		return 0, err
	}

	// emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, cr.Address),
			sdk.NewAttribute("zone", cr.ChainId),
			sdk.NewAttribute(sdk.AttributeKeyAction, action.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
		),
	})

	return claimAmount, nil
}

// getClaimAmountAndUpdateRecord calculates and returns the claimable amount
// after updating the relevant claim record.
func (k *Keeper) getClaimAmountAndUpdateRecord(ctx sdk.Context, cr *types.ClaimRecord, action types.Action) (uint64, error) {
	var claimAmount uint64

	// check and initialize ActionsCompleted map
	if cr.ActionsCompleted == nil {
		cr.ActionsCompleted = make(map[int32]*types.CompletedAction)
	}

	// The concept here is to intuitively claim all outstanding deposit tiers
	// that are below the current deposit claim (improved user experience).
	//
	// ActionDepositT5: t5amount
	// ActionDepositT4: t4amount
	// ActionDepositT3: t3amount  <-- eg. for T3
	// ActionDepositT2: t2amount  <-- add to claimAmount if not CompletedAction
	// ActionDepositT1: t1amount  <-- add to claimAmount if not CompletedAction
	//
	// For any given deposit action above ActionDepositT1, sum the claimable
	// amounts of non completed deposit actions and mark them as complete.
	// Then, if no errors occurred, update the ClaimRecord state.

	// check for summable ActionDeposit (T2-T5, T1 has nothing below it to sum)
	if action > types.ActionDepositT1 && action <= types.ActionDepositT5 {
		// check ActionDeposits from T1 to the target tier
		// this also ensures that for any completed ActionDeposit tier, all
		// tiers below are guaranteed to be completed as well.
		for a := types.ActionDepositT1; a <= action; a++ {
			if _, exists := cr.ActionsCompleted[int32(a)]; !exists {
				// obtain claimable amount per deposit action
				claimable, err := k.GetClaimableAmountForAction(ctx, cr.ChainId, cr.Address, a)
				if err != nil {
					return 0, err
				}

				// update claim record (transient, not yet written to state)
				cr.ActionsCompleted[int32(a)] = &types.CompletedAction{
					CompleteTime: ctx.BlockTime(),
					ClaimAmount:  claimable,
				}

				// sum total claimable
				claimAmount += claimable
			}
		}
	} else {
		// obtain claimable amount
		claimable, err := k.GetClaimableAmountForAction(ctx, cr.ChainId, cr.Address, action)
		if err != nil {
			return 0, err
		}

		// set claim amount
		claimAmount = claimable

		// update claim record
		cr.ActionsCompleted[int32(action)] = &types.CompletedAction{
			CompleteTime: ctx.BlockTime(),
			ClaimAmount:  claimAmount,
		}
	}

	// set claim record (persistent)
	if err := k.SetClaimRecord(ctx, *cr); err != nil {
		return 0, err
	}

	return claimAmount, nil
}

func (k *Keeper) sendCoins(ctx sdk.Context, cr types.ClaimRecord, amount uint64) (sdk.Coins, error) {
	coins := sdk.NewCoins(
		sdk.NewCoin(k.BondDenom(ctx), sdk.NewIntFromUint64(amount)),
	)

	addr, err := sdk.AccAddressFromBech32(cr.Address)
	if err != nil {
		return sdk.NewCoins(), err
	}

	if err := k.bankKeeper.SendCoins(ctx, k.GetZoneDropAccountAddress(cr.ChainId), addr, coins); err != nil {
		return sdk.NewCoins(), err
	}

	return coins, nil
}
