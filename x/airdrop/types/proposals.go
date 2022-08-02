package types

import (
	fmt "fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterZoneDrop = "RegisterZoneDrop"
)

var (
	_ govtypes.Content = &RegisterZoneDropProposal{}
)

func (m RegisterZoneDropProposal) GetDescription() string { return m.Description }
func (m RegisterZoneDropProposal) GetTitle() string       { return m.Title }
func (m RegisterZoneDropProposal) ProposalRoute() string  { return RouterKey }
func (m RegisterZoneDropProposal) ProposalType() string   { return ProposalTypeRegisterZoneDrop }

// ValidateBasic runs basic stateless validity checks
func (m RegisterZoneDropProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	// validate ZoneDrop

	// validate ClaimRecords

	return nil
}

// String implements the Stringer interface.
func (m RegisterZoneDropProposal) String() string {
	var b strings.Builder

	b.WriteString("Airdrop - ZoneDrop Registration Proposal:\n")
	b.WriteString(fmt.Sprintf("\tTitle:       %s\n", m.Title))
	b.WriteString(fmt.Sprintf("\tDescription: %s\n", m.Description))
	b.WriteString(fmt.Sprintf("\tZoneDrop:\n\n%v\n----------", m.ZoneDrop))
	b.WriteString("\tClaimRecords:\n\n")

	for _, cr := range m.ClaimRecords {
		b.WriteString(fmt.Sprintf("%v\n\n", cr))
	}
	b.WriteString("----------\n")

	return b.String()
}
