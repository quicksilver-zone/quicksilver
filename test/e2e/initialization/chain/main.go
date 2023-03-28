package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ingenuity-build/quicksilver/test/e2e/initialization"
)

func main() {
	var (
		valConfig    []*initialization.NodeConfig
		dataDir      string
		chainID      string
		config       string
		votingPeriod time.Duration
		forkHeight   int
	)

	flag.StringVar(&dataDir, "data-dir", "", "chain data directory")
	flag.StringVar(&chainID, "chain-id", "", "chain ID")
	flag.StringVar(&config, "config", "", "serialized config")
	flag.DurationVar(&votingPeriod, "voting-period", 30000000000, "voting period")
	flag.IntVar(&forkHeight, "fork-height", 0, "fork height")

	flag.Parse()

	err := json.Unmarshal([]byte(config), &valConfig)
	if err != nil {
		panic(err)
	}

	if len(dataDir) == 0 {
		panic("data-dir is required")
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		panic(err)
	}

	createdChain, err := initialization.InitChain(chainID, dataDir, valConfig, votingPeriod, forkHeight)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(createdChain)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("%v/%v-encode", dataDir, chainID)
	if err = os.WriteFile(fileName, b, 0o777); err != nil { //nolint:gosec
		panic(err)
	}
}
