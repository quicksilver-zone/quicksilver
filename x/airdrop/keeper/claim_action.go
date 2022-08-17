package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/ingenuity-build/quicksilver/x/airdrop/types"
)

// VerifyClaimAction validates the claim against the specific action
// conditions. It is the responsibility of the caller to ensure that the given
// action is in bounds by calling action.InBounds().
//
// TODO: we also want to verify that the action was executed on the remote
// chain before we execute the claim...
func (k Keeper) VerifyClaimAction(ctx sdk.Context, cr types.ClaimRecord, action types.Action) error {
	// action already completed, nothing to claim
	if _, exists := cr.ActionsCompleted[int32(action)]; exists {
		return fmt.Errorf("%s already completed", types.Action_name[int32(action)])
	}

	switch action {
	case types.ActionInitialClaim:
		return nil
	case types.ActionDepositT1:
		return k.checkDeposit(ctx, cr, sdk.MustNewDecFromStr("0.05"))
	case types.ActionDepositT2:
		return k.checkDeposit(ctx, cr, sdk.MustNewDecFromStr("0.10"))
	case types.ActionDepositT3:
		return k.checkDeposit(ctx, cr, sdk.MustNewDecFromStr("0.15"))
	case types.ActionDepositT4:
		return k.checkDeposit(ctx, cr, sdk.MustNewDecFromStr("0.22"))
	case types.ActionDepositT5:
		return k.checkDeposit(ctx, cr, sdk.MustNewDecFromStr("0.30"))
	case types.ActionStakeQCK:
		return k.checkBondedDelegation(ctx, cr.Address)
	case types.ActionSignalIntent:
		return k.checkIntentIsSet(ctx, cr.ChainId, cr.Address)
	case types.ActionQSGov:
		return k.checkQSGov(ctx, cr.Address)
	case types.ActionGbP:
		// TODO: implement check once GbP is implemented
	case types.ActionOsmosis:
		// IBC proof based verification on Osmosis remote zone
		// TODO: implement
	default:
		return fmt.Errorf("undefined action [%d]", action)
	}

	return fmt.Errorf("verification not implemented for [%d] %s", action, types.Action_name[int32(action)])
}

// checkDeposit checks
func (k Keeper) checkDeposit(ctx sdk.Context, cr types.ClaimRecord, threshold sdk.Dec) error {
	addr, err := sdk.AccAddressFromBech32(cr.Address)
	if err != nil {
		return err
	}

	zone, ok := k.icsKeeper.GetZone(ctx, cr.ChainId)
	if !ok {
		return fmt.Errorf("zone not found for %s", cr.ChainId)
	}

	// obtain all deposit receipts for this user on this zone
	rcpts, err := k.icsKeeper.UserZoneReceipts(ctx, &zone, addr)
	if err != nil {
		return fmt.Errorf("unable to obtain zone receipts for %s on zone %s: %w", cr.Address, cr.ChainId, err)
	}

	// sum gross deposits amount
	gdAmount := sdk.NewInt(0)
	for _, rcpt := range rcpts {
		gdAmount = gdAmount.Add(rcpt.Amount.AmountOf(zone.BaseDenom))
	}

	// calculate target amount
	tAmount := threshold.MulInt64(int64(cr.BaseValue)).TruncateInt()

	if gdAmount.LT(tAmount) {
		return fmt.Errorf("insufficient deposit amount")
	}

	return nil
}

// checkBondedDelegation indicates if the given address has an active bonded
// delegation of QCK on the Quicksilver zone.
func (k Keeper) checkBondedDelegation(ctx sdk.Context, address string) error {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	amount := k.stakingKeeper.GetDelegatorBonded(ctx, addr)
	if !amount.IsPositive() {
		return fmt.Errorf("ActionStakeQCK: no bonded delegation")
	}
	return nil
}

// checkIntentIsSet indicates if the given address has intent set for the given
// zone (chainID).
func (k Keeper) checkIntentIsSet(ctx sdk.Context, chainID string, address string) error {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	zone, ok := k.icsKeeper.GetZone(ctx, chainID)
	if !ok {
		return fmt.Errorf("zone %s not found", chainID)
	}

	intent, ok := k.icsKeeper.GetIntent(ctx, zone, addr.String(), false)
	if !ok || len(intent.Intents) == 0 {
		return fmt.Errorf("intent not found or no intents set for %s", addr)
	}

	return nil
}

// checkQSGov indicates if the given address has voted on any governance
// proposals on the Quicksilver zone.
func (k Keeper) checkQSGov(ctx sdk.Context, address string) error {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	voted := false
	k.govKeeper.IterateProposals(ctx, func(proposal gov.Proposal) (stop bool) {
		_, found := k.govKeeper.GetVote(ctx, proposal.ProposalId, addr)
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
