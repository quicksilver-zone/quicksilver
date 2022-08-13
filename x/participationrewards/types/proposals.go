package types

import (
	"encoding/json"
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddProtocolData = "AddProtocolData"
)

var _ govtypes.Content = &AddProtocolDataProposal{}

func NewAddProtocolDataProposal(title string, description string, datatype string, protocol string, key string, data json.RawMessage) *AddProtocolDataProposal {
	return &AddProtocolDataProposal{Title: title, Description: description, Type: datatype, Protocol: protocol, Key: key, Data: data}
}

func (m AddProtocolDataProposal) GetDescription() string { return m.Description }
func (m AddProtocolDataProposal) GetTitle() string       { return m.Title }
func (m AddProtocolDataProposal) ProposalRoute() string  { return RouterKey }
func (m AddProtocolDataProposal) ProposalType() string   { return ProposalTypeAddProtocolData }

// ValidateBasic runs basic stateless validity checks
func (m AddProtocolDataProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (m AddProtocolDataProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Add Protocol Data Proposal:
  Title:                            %s
  Description:                      %s
  Protocol:                         %s
  Type:						 %s
  Key:                       %s
  Data:                      %s
`, m.Title, m.Description, m.Protocol, m.Type, m.Key, m.Data))
	return b.String()
}
