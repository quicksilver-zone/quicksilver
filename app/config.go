package app

// DONTCOVER

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	purningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	dbm "github.com/cometbft/cometbft-db"
)

func DefaultConfig() network.Config {
	encCfg := MakeEncodingConfig()

	return network.Config{
		Codec:             encCfg.Marshaler,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor:    NewAppConstructor(encCfg),
		GenesisState:      ModuleBasics.DefaultGenesis(encCfg.Marshaler),
		TimeoutCommit:     1 * time.Second / 2,
		ChainID:           "quicktest-1",
		NumValidators:     1,
		BondDenom:         sdk.DefaultBondDenom,
		MinGasPrices:      fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:     sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:     sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:      sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		CleanupDir:        true,
		SigningAlgo:       string(hd.Secp256k1Type),
		KeyringOptions:    []keyring.Option{},
	}
}

func NewAppConstructor(encCfg EncodingConfig) network.AppConstructor {
	return func(val network.ValidatorI) servertypes.Application {
		valCtx := val.GetCtx()
		appConfig := val.GetAppConfig()

		return NewQuicksilver(
			valCtx.Logger,
			dbm.NewMemDB(),
			nil,
			true,
			map[int64]bool{},
			DefaultNodeHome,
			0,
			encCfg,
			EmptyAppOptions{},
			false,
			false,
			"",
			baseapp.SetPruning(purningtypes.NewPruningOptionsFromString(appConfig.Pruning)),
			baseapp.SetChainID("quicktest-1"),

			// baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}
