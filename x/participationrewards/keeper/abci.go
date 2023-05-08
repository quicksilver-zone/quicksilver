package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

<<<<<<< HEAD
// BeginBlocker of participationrewards module
func (k Keeper) BeginBlocker(_ sdk.Context) {
=======
// BeginBlocker of participationrewards module.
func (k *Keeper) BeginBlocker(_ sdk.Context) {
>>>>>>> origin/develop
}
