package cli

import (
	"bytes"
	"testing"
)

func TestGetSignalIntentTxCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			"no args",
			[]string{},
			true,
		},
		{
			"empty args",
			[]string{"", ""},
			true,
		},
		{
			"invalid delegation_intent arg_format",
			[]string{"chainid", "intents"},
			true,
		},
		{
			"invalid delegation_intent content",
			[]string{"chainid", "0.0cosmos1valoper1xxxxxxxxx,0.1cosmosvaloper1yyyyyyyyy,1.1cosmosvaloper1zzzzzzzzz"},
			true,
		},
		{
			"invalid delegation_intent valoperAddress",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,0.1notok1ezyfcl"},
			true,
		},
		{
			"invalid delegation_intent weightOverrun",
			[]string{"chainid", "0.4A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
		},
		{
			"invalid delegation_intent weightUnderrun",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,0.3abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
		},
		{
			"invalid delegation_intent maxWeightOverrun",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,1.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
		},
		/*{
			"valid",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			false,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := []string{"--dry-run"}
			args := append(tt.args, flags...)
			cmd := GetSignalIntentTxCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(nil)
			cmd.SetArgs(args)
			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("GetSignalIntentTxCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
