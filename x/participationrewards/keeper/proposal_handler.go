package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/types"
)

// HandleAddProtocolDataProposal is a handler for executing a passed add protocol data proposal.
func HandleAddProtocolDataProposal(ctx sdk.Context, k *Keeper, p *types.AddProtocolDataProposal) error {
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

	protocolData := types.NewProtocolData(p.Type, p.Data)

	k.SetProtocolData(ctx, pd.GenerateKey(), protocolData)

	return nil
}
