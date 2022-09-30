package types

import (
	"encoding/json"
	"fmt"
	"strings"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddProtocolData = "AddProtocolData"
)

var _ govv1beta1.Content = &AddProtocolDataProposal{}

func NewAddProtocolDataProposal(title string, description string, datatype string, protocol string, key string, data json.RawMessage) *AddProtocolDataProposal {
	return &AddProtocolDataProposal{Title: title, Description: description, Type: datatype, Protocol: protocol, Key: key, Data: data}
}

func (m AddProtocolDataProposal) GetDescription() string { return m.Description }
func (m AddProtocolDataProposal) GetTitle() string       { return m.Title }
func (m AddProtocolDataProposal) ProposalRoute() string  { return RouterKey }
func (m AddProtocolDataProposal) ProposalType() string   { return ProposalTypeAddProtocolData }

// ValidateBasic runs basic stateless validity checks
func (m AddProtocolDataProposal) ValidateBasic() error {
	if err := govv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.Protocol) == 0 {
		return fmt.Errorf("proposal must specify Protocol")
	}

	if len(m.Type) == 0 {
		return fmt.Errorf("proposal must specify Type")
	}

	if len(m.Key) == 0 {
		return fmt.Errorf("proposal must specify Key")
	}

	if m.Data == nil {
		return fmt.Errorf("proposal must specify Data")
	}

	return nil
}

// String implements the Stringer interface.
func (m AddProtocolDataProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Add Protocol Data Proposal:
Title:			%s
Description:	%s
Protocol:		%s
Type:			%s
Key:			%s
Data:			%s
`, m.Title, m.Description, m.Protocol, m.Type, m.Key, m.Data))
	return b.String()
}
