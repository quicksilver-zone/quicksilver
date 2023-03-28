package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeRegisterZone = "RegisterZone"
	ProposalTypeUpdateZone   = "UpdateZone"
)

var (
	_ govv1beta1.Content = &RegisterZoneProposal{}
	_ govv1beta1.Content = &UpdateZoneProposal{}
)

func NewRegisterZoneProposal(title string, description string, connectionID string, baseDenom string, localDenom string, accountPrefix string, multiSend bool, liquidityModule bool) *RegisterZoneProposal {
	return &RegisterZoneProposal{Title: title, Description: description, ConnectionId: connectionID, BaseDenom: baseDenom, LocalDenom: localDenom, AccountPrefix: accountPrefix, MultiSend: multiSend, LiquidityModule: liquidityModule}
}

func (m RegisterZoneProposal) GetDescription() string { return m.Description }
func (m RegisterZoneProposal) GetTitle() string       { return m.Title }
func (m RegisterZoneProposal) ProposalRoute() string  { return RouterKey }
func (m RegisterZoneProposal) ProposalType() string   { return ProposalTypeRegisterZone }

// ValidateBasic runs basic stateless validity checks
func (m RegisterZoneProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(m)
	if err != nil {
		return err
	}

	// check valid connection id
	if len(m.ConnectionId) < 12 {
		return fmt.Errorf("invalid length connection string: %s", m.ConnectionId)
	}

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
		return errors.New("account prefix must be at least 2 characters") // ki is shortest to date.
	}

	// validate messages_per_tx
	if m.MessagesPerTx < 1 {
		return errors.New("messages_per_tx must be a positive non-zero integer")
	}

	if m.LiquidityModule {
		return errors.New("liquidity module is unsupported")
	}
	return nil
}

// String implements the Stringer interface.
func (m RegisterZoneProposal) String() string {
	return fmt.Sprintf(`Interchain Staking  Zone Registration Proposal:
  Title:                            %s
  Description:                      %s
  Connection Id:                    %s
  Base Denom:                       %s
  Local Denom:                      %s
  Multi Send Enabled:               %t
  Liquidity Staking Module Enabled: %t
  Messages per Tx:                  %d
`, m.Title, m.Description, m.ConnectionId, m.BaseDenom, m.LocalDenom, m.MultiSend, m.LiquidityModule, m.MessagesPerTx)
}

func NewUpdateZoneProposal(title string, description string, chainID string, changes []*UpdateZoneValue) *UpdateZoneProposal {
	return &UpdateZoneProposal{Title: title, Description: description, ChainId: chainID, Changes: changes}
}

func (m UpdateZoneProposal) GetDescription() string { return m.Description }
func (m UpdateZoneProposal) GetTitle() string       { return m.Title }
func (m UpdateZoneProposal) ProposalRoute() string  { return RouterKey }
func (m UpdateZoneProposal) ProposalType() string   { return ProposalTypeUpdateZone }

// ValidateBasic runs basic stateless validity checks
func (m UpdateZoneProposal) ValidateBasic() error {
	return govv1beta1.ValidateAbstract(m)
}

// String implements the Stringer interface.
func (m UpdateZoneProposal) String() string {
	b := new(strings.Builder)
	fmt.Fprintf(b, `Interchain Staking Zone Update Proposal:
  Title:       %s
  Description: %s
  Changes:\n
`, m.Title, m.Description)
	for _, change := range m.Changes {
		fmt.Fprintf(b, `
	  Key:   %s
	  Value: %s
	  -----------------------
	`, change.Key, change.Value)
	}
	return b.String()
}

func (v UpdateZoneValue) Validate() error {
	return nil
}
