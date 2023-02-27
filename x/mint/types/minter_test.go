package types_test

import (
	"github.com/ingenuity-build/quicksilver/x/mint/types"
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Benchmarking :)
// previously using sdk.Int operations:
// BenchmarkEpochProvision-4 5000000 220 ns/op
//
// using sdk.Dec operations: (current implementation)
// BenchmarkEpochProvision-4 3000000 429 ns/op
func BenchmarkEpochProvision(b *testing.B) {
	b.ReportAllocs()
	minter := types.InitialMinter()
	params := types.DefaultParams()

	s1 := rand.NewSource(100)
	r1 := rand.New(s1)
	minter.EpochProvisions = sdk.NewDec(r1.Int63n(1000000))

	// run the EpochProvision function b.N times
	for n := 0; n < b.N; n++ {
		minter.EpochProvision(params)
	}
}

// Next epoch provisions benchmarking
// BenchmarkNextEpochProvisions-4 5000000 251 ns/op
func BenchmarkNextEpochProvisions(b *testing.B) {
	b.ReportAllocs()
	minter := types.InitialMinter()
	params := types.DefaultParams()

	// run the NextEpochProvisions function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextEpochProvisions(params)
	}
}

func TestMinterValidate(t *testing.T) {
	testcases := []struct {
		name    string
		minter  types.Minter
		isValid bool
	}{
		{
			"valid",
			types.InitialMinter(),
			true,
		},
		{
			"negative",
			types.Minter{
				EpochProvisions: sdk.NewDec(-1),
			},
			false,
		},
		{
			"nil",
			types.Minter{},
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.minter.Validate()
			if !tc.isValid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
