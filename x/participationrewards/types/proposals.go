package types

import (
	"encoding/json"
	"fmt"

	"github.com/ingenuity-build/multierror"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddProtocolData = "AddProtocolData"
)

var _ govv1beta1.Content = &AddProtocolDataProposal{}

func NewAddProtocolDataProposal(title, description, datatype, _, key string, data json.RawMessage) *AddProtocolDataProposal {
	return &AddProtocolDataProposal{Title: title, Description: description, Type: datatype, Data: data, Key: key}
}

func (m *AddProtocolDataProposal) GetDescription() string { return m.Description }
func (m *AddProtocolDataProposal) GetTitle() string       { return m.Title }
func (*AddProtocolDataProposal) ProposalRoute() string    { return RouterKey }
func (*AddProtocolDataProposal) ProposalType() string     { return ProposalTypeAddProtocolData }

// ValidateBasic runs basic stateless validity checks.
func (m *AddProtocolDataProposal) ValidateBasic() error {
	if err := govv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	errors := make(map[string]error)

	if m.Type == "" {
		errors["Type"] = ErrUndefinedAttribute
	}

	// Key is now a deprecated field and unused.
	// if len(m.Key) == 0 {
	// 	errors["Key"] = ErrUndefinedAttribute
	// }

	if len(m.Data) == 0 {
		errors["Data"] = ErrUndefinedAttribute
	}

	pd, err := UnmarshalProtocolData(ProtocolDataType(ProtocolDataType_value[m.Type]), m.Data)
	if err != nil {
		errors["Data"] = err
	} else {
		if err = pd.ValidateBasic(); err != nil {
			errors["Data"] = err
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

// String implements the Stringer interface.
func (m *AddProtocolDataProposal) String() string {
	return fmt.Sprintf(`Add Protocol Data Proposal:
Title:			%s
Description:	%s
Type:			%s
Data:			%s
`, m.Title, m.Description, m.Type, m.Data)
}
