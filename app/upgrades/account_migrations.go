package upgrades

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/app/keepers"
)

type ProcessMigrateAccountStrategy func(ctx sdk.Context, appKeepers *keepers.AppKeepers, from sdk.AccAddress, to sdk.AccAddress) error
