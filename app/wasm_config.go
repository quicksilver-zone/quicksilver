package app

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	// DefaultInstanceCost is initially set the same as in wasmd.
	DefaultInstanceCost uint64 = 60_000
	// DefaultCompileCost set to a large number for testing.
	DefaultCompileCost uint64 = 100
)

// GasRegisterConfig is defaults plus a custom compile amount.
func GasRegisterConfig() wasmtypes.WasmGasRegisterConfig {
	gasConfig := wasmtypes.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultInstanceCost
	gasConfig.CompileCost = DefaultCompileCost

	return gasConfig
}

func NewWasmGasRegister() wasmtypes.WasmGasRegister {
	return wasmtypes.NewWasmGasRegister(GasRegisterConfig())
}
