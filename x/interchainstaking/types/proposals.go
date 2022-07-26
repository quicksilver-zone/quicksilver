package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeRegisterZone = "RegisterZone"
	ProposalTypeUpdateZone   = "UpdateZone"
)

var _ govtypes.Content = &RegisterZoneProposal{}
var _ govtypes.Content = &UpdateZoneProposal{}

func NewRegisterZoneProposal(title string, description string, connection_id string, base_denom string, local_denom string, account_prefix string, multi_send bool, liquidity_module bool) *RegisterZoneProposal {
	return &RegisterZoneProposal{Title: title, Description: description, ConnectionId: connection_id, BaseDenom: base_denom, LocalDenom: local_denom, AccountPrefix: account_prefix, MultiSend: multi_send, LiquidityModule: liquidity_module}
}

func (m RegisterZoneProposal) GetDescription() string { return m.Description }
func (m RegisterZoneProposal) GetTitle() string       { return m.Title }
func (m RegisterZoneProposal) ProposalRoute() string  { return RouterKey }
func (m RegisterZoneProposal) ProposalType() string   { return ProposalTypeRegisterZone }

// ValidateBasic runs basic stateless validity checks
func (m RegisterZoneProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	// check valid connection id
	if m.ConnectionId[0:11] != "connection-" {
		return fmt.Errorf("invalid connection string: %s", m.ConnectionId)
	}

	// validate local denominations
	if err := sdk.ValidateDenom(m.LocalDenom); err != nil {
		return err
	}

	// validate base denom
	if err := sdk.ValidateDenom(m.BaseDenom); err != nil {
		return err
	}

	// validate account prefix
	if len(m.AccountPrefix) < 2 {
		return fmt.Errorf("account prefix must be at least 2 characters") // ki is shortest to date.
	}
	return nil
}

// String implements the Stringer interface.
func (m RegisterZoneProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Interchain Staking  Zone Registration Proposal:
  Title:                            %s
  Description:                      %s
  Connection Id:                    %s
  Base Denom:                       %s
  Local Denom:                      %s
  Multi Send Enabled:               %t
  Liquidity Staking Module Enabled: %t
`, m.Title, m.Description, m.ConnectionId, m.BaseDenom, m.LocalDenom, m.MultiSend, m.LiquidityModule))
	return b.String()
}

func NewUpdateZoneProposal(title string, description string, chain_id string, changes []*UpdateZoneValue) *UpdateZoneProposal {
	return &UpdateZoneProposal{Title: title, Description: description, ChainId: chain_id, Changes: changes}
}

func (m UpdateZoneProposal) GetDescription() string { return m.Description }
func (m UpdateZoneProposal) GetTitle() string       { return m.Title }
func (m UpdateZoneProposal) ProposalRoute() string  { return RouterKey }
func (m UpdateZoneProposal) ProposalType() string   { return ProposalTypeUpdateZone }

// ValidateBasic runs basic stateless validity checks
func (m UpdateZoneProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (m UpdateZoneProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Interchain Staking Zone Update Proposal:
  Title:       %s
  Description: %s
  Changes:\n
`, m.Title, m.Description))
	for _, change := range m.Changes {
		b.WriteString(fmt.Sprintf(`
	  Key:   %s
	  Value: %s
	  -----------------------
	`, change.Key, change.Value))
	}
	return b.String()
}

func (v UpdateZoneValue) Validate() error {

	return nil
}
