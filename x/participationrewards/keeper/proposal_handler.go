package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// HandleAddProtocolDataProposal is a handler for executing a passed add protocol data proposal.
func HandleAddProtocolDataProposal(ctx sdk.Context, k *Keeper, p *types.AddProtocolDataProposal) error {
	protocolData := types.NewProtocolData(p.Type, p.Data)

	if err := p.ValidateBasic(); err != nil {
		return err
	}

	pdtv, exists := types.ProtocolDataType_value[p.Type]
	if !exists {
		return types.ErrUnknownProtocolDataType
	}

	pd, err := types.UnmarshalProtocolData(types.ProtocolDataType(pdtv), p.Data)
	if err != nil {
		return err
	}

	if err := pd.ValidateBasic(); err != nil {
		return err
	}

	k.SetProtocolData(ctx, pd.GenerateKey(), protocolData)

	return nil
}
