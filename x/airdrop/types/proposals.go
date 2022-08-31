package types

import (
	fmt "fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterZoneDrop = "RegisterZoneDrop"
)

var _ govtypes.Content = &RegisterZoneDropProposal{}

func (m RegisterZoneDropProposal) GetDescription() string { return m.Description }
func (m RegisterZoneDropProposal) GetTitle() string       { return m.Title }
func (m RegisterZoneDropProposal) ProposalRoute() string  { return RouterKey }
func (m RegisterZoneDropProposal) ProposalType() string   { return ProposalTypeRegisterZoneDrop }

// ValidateBasic runs basic stateless validity checks.
//
// ZoneDrop is validated in HandleRegisterZoneDropProposal.
// ClaimRecords are validated in HandleRegisterZoneDropProposal.
//
// HandleRegisterZoneDropProposal does validation checks as ZoneDrop is related
// to ClaimRecords. ClaimRecords are in compressed []byte slice format and
// must be decompressed in order to be validated.
func (m RegisterZoneDropProposal) ValidateBasic() error {
	return govtypes.ValidateAbstract(m)
}

// String implements the Stringer interface.
func (m RegisterZoneDropProposal) String() string {
	var b strings.Builder

	b.WriteString("Airdrop - ZoneDrop Registration Proposal:\n")
	b.WriteString(fmt.Sprintf("\tTitle:       %s\n", m.Title))
	b.WriteString(fmt.Sprintf("\tDescription: %s\n", m.Description))
	b.WriteString("\tZoneDrop:\n")
	b.WriteString(fmt.Sprintf("\n%v\n", m.ZoneDrop))
	b.WriteString("\n----------\n")
	return b.String()
}
