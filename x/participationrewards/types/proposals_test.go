package types

import (
	"encoding/json"
	encoding_json "encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var validLiquidData string = `{
	"chainid": "somechain",
	"localdenom": "lstake",
	"denom": "qstake"
}`

func TestAddProtocolDataProposal_ValidateBasic(t *testing.T) {
	type fields struct {
		Title       string
		Description string
		Protocol    string
		Type        string
		Key         string
		Data        json.RawMessage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"blank",
			fields{},
			true,
		},
		{
			"invalid_protocol",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "",
				Type:        "",
				Key:         "",
				Data:        nil,
			},
			true,
		},
		{
			"invalid_type",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "TestProtocol",
				Type:        "",
				Key:         "",
				Data:        nil,
			},
			true,
		},
		{
			"invalid_key",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "TestProtocol",
				Type:        "TestType",
				Key:         "",
				Data:        nil,
			},
			true,
		},
		{
			"invalid_data",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "TestProtocol",
				Type:        "TestType",
				Key:         "TestKey",
				Data:        nil,
			},
			true,
		},
		{
			"valid_data",
			fields{
				Title:       "Valid Protocol Data",
				Description: "A valid protocol that is valid",
				Protocol:    "ValidProtocol",
				Type:        "liquidtoken",
				Key:         "liquid",
				Data:        []byte(validLiquidData),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AddProtocolDataProposal{
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Protocol:    tt.fields.Protocol,
				Type:        tt.fields.Type,
				Key:         tt.fields.Key,
				Data:        tt.fields.Data,
			}
			err := m.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAddProtocolDataProposal_String(t *testing.T) {
	type fields struct {
		Title       string
		Description string
		Protocol    string
		Type        string
		Key         string
		Data        encoding_json.RawMessage
	}

	tt := fields{
		Title:       "Valid Protocol Data",
		Description: "A valid protocol that is valid",
		Protocol:    "ValidProtocol",
		Type:        "liquidtoken",
		Key:         "liquid",
		Data:        []byte(validLiquidData),
	}

	want := `Add Protocol Data Proposal:
Title:			Valid Protocol Data
Description:	A valid protocol that is valid
Protocol:		ValidProtocol
Type:			liquidtoken
Key:			liquid
Data:			{
	"chainid": "somechain",
	"localdenom": "lstake",
	"denom": "qstake"
}
`

	t.Run("stringer", func(t *testing.T) {
		m := AddProtocolDataProposal{
			Title:       tt.Title,
			Description: tt.Description,
			Protocol:    tt.Protocol,
			Type:        tt.Type,
			Key:         tt.Key,
			Data:        tt.Data,
		}
		got := m.String()
		require.Equal(t, want, got)
	})
}
