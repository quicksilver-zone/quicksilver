package keeper_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func FuzzCalcUserScores(f *testing.F) {
	f.Skip("In -short")

	files, err := filepath.Glob(filepath.Join("testdata", "fuzz-corpus-CalcUserValidatorSectionAllocations-*"))
	if err != nil {
		f.Fatal(err)
	}

	type corpusData struct {
		Ctx       sdk.Context     `json:"ctx"`
		Zone      *icstypes.Zone  `json:"zone"`
		ZoneScore types.ZoneScore `json:"zs"`
	}

	// 1. Add the corpus.
	for _, filename := range files {
		corpusData, err := os.ReadFile(filename)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(corpusData)
	}

	ste := new(KeeperTestSuite)
	ste.SetupTest()
	appA := ste.GetQuicksilverApp(ste.chainA)

	// 2. Run the fuzzers.
	f.Fuzz(func(t *testing.T, input []byte) {
		cj := new(corpusData)
		if err := json.Unmarshal(input, cj); err != nil {
			t.Fatal(err)
		}

		_ = appA.ParticipationRewardsKeeper.CalcUserValidatorSelectionAllocations(cj.Ctx, cj.Zone, cj.ZoneScore)
	})
}
