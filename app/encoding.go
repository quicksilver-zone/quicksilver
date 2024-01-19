package app

import (
	"testing"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig(tb testing.TB) EncodingConfig {
	tb.Helper()

	tempApp := NewQuicksilver(log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{}, "abcd",
		simtestutil.NewAppOptionsWithFlagHome(tb.TempDir()), true, false, GetWasmOpts(EmptyAppOptions{}))
	return makeEncodingConfig(tempApp)
}

func makeEncodingConfig(tempApp *Quicksilver) EncodingConfig {
	encodingConfig := EncodingConfig{
		InterfaceRegistry: tempApp.InterfaceRegistry(),
		Codec:             tempApp.AppCodec(),
		TxConfig:          tempApp.GetTxConfig(),
		Amino:             tempApp.LegacyAmino(),
	}
	return encodingConfig
}
