package keeper_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func FuzzCalcUserScores(f *testing.F) {
	if testing.Short() {
		f.Skip("In -short")
	}

	files, err := filepath.Glob(filepath.Join("testdata", "fuzz-corpus-CalcUserValidatorSectionAllocations-*"))
	if err != nil {
		f.Fatal(err)
	}

	type corpusData struct {
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

	// 2. Run the fuzzers.
	f.Fuzz(func(t *testing.T, input []byte) {
		ste.SetT(t)
		ste.SetS(ste)
		ste.SetupTest()
		appA := ste.GetQuicksilverApp(ste.chainA)

		cj := new(corpusData)
		if err := json.Unmarshal(input, cj); err != nil {
			t.Fatal(err)
		}

		_ = appA.ParticipationRewardsKeeper.CalcUserValidatorSelectionAllocations(ste.chainA.GetContext(), cj.Zone, cj.ZoneScore)
	})
}
