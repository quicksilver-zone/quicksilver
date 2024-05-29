package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

var benchSink any = nil

func BenchmarkIntentsFromString(b *testing.B) {
	str := "0.3cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0,0.3cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf,0.4cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll"

	want := []*types.ValidatorIntent{
		{
			ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
			Weight:         sdk.MustNewDecFromStr("0.3"),
		},
		{
			ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			Weight:         sdk.MustNewDecFromStr("0.3"),
		},
		{
			ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll",
			Weight:         sdk.MustNewDecFromStr("0.4"),
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got, err := types.IntentsFromString(str)
		if err != nil {
			b.Fatal(err)
		}
		if diff := cmp.Diff(got, want); diff != "" {
			b.Fatalf("Mismatch: got - want +\n%s", diff)
		}
		benchSink = got
	}

	if benchSink == nil {
		b.Fatal("Benchmark did not run!")
	}

	benchSink = nil
}

func BenchmarkSetForValoper(b *testing.B) {
	v1 := addressutils.GenerateValAddressForTest().String()
	v2 := addressutils.GenerateValAddressForTest().String()
	v3 := addressutils.GenerateValAddressForTest().String()
	v4 := addressutils.GenerateValAddressForTest().String()
	v5 := addressutils.GenerateValAddressForTest().String()
	v6 := addressutils.GenerateValAddressForTest().String()
	v7 := addressutils.GenerateValAddressForTest().String()
	v8 := addressutils.GenerateValAddressForTest().String()
	v9 := addressutils.GenerateValAddressForTest().String()
	v10 := addressutils.GenerateValAddressForTest().String()

	intents := types.ValidatorIntents{
		{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(10, 1)},
		{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(20, 1)},
		// Deliberately missing v3.
		{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(10, 1)},
		{ValoperAddress: v5, Weight: sdk.NewDecWithPrec(10, 1)},
		{ValoperAddress: v6, Weight: sdk.NewDecWithPrec(10, 1)},
		{ValoperAddress: v7, Weight: sdk.NewDecWithPrec(10, 1)},
		{ValoperAddress: v8, Weight: sdk.NewDecWithPrec(10, 1)},
	}

	for i := 0; i < 1000; i++ {
		vi := addressutils.GenerateValAddressForTest().String()
		intents = append(intents, &types.ValidatorIntent{ValoperAddress: vi, Weight: sdk.NewDecWithPrec(0, 1)})
	}

	// Deliberately add v9 and v10 at the very end of the slice to trigger O(n) behavior on first search.
	intents = append(intents, &types.ValidatorIntent{ValoperAddress: v9, Weight: sdk.NewDecWithPrec(10, 1)})
	intents = append(intents, &types.ValidatorIntent{ValoperAddress: v10, Weight: sdk.NewDecWithPrec(10, 1)})

	require.Equal(b, sdk.NewDecWithPrec(10, 1), intents.MustGetForValoper(v1).Weight)
	require.Equal(b, sdk.NewDecWithPrec(20, 1), intents.MustGetForValoper(v2).Weight)
	require.Equal(b, sdk.NewDecWithPrec(10, 1), intents.MustGetForValoper(v4).Weight)
	require.Equal(b, sdk.NewDecWithPrec(10, 1), intents.MustGetForValoper(v9).Weight)
	require.Equal(b, sdk.NewDecWithPrec(10, 1), intents.MustGetForValoper(v10).Weight)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		intents = intents.SetForValoper(v1, &types.ValidatorIntent{ValoperAddress: v1, Weight: sdk.NewDecWithPrec(20, 1)})
		intents = intents.SetForValoper(v2, &types.ValidatorIntent{ValoperAddress: v2, Weight: sdk.NewDecWithPrec(10, 1)})
		intents = intents.SetForValoper(v4, &types.ValidatorIntent{ValoperAddress: v4, Weight: sdk.NewDecWithPrec(15, 1)})
		intents = intents.SetForValoper(v10, &types.ValidatorIntent{ValoperAddress: v10, Weight: sdk.NewDecWithPrec(5, 1)})

		require.Equal(b, sdk.NewDecWithPrec(20, 1), intents.MustGetForValoper(v1).Weight)
		require.Equal(b, sdk.NewDecWithPrec(10, 1), intents.MustGetForValoper(v2).Weight)
		require.Equal(b, sdk.NewDecWithPrec(15, 1), intents.MustGetForValoper(v4).Weight)
		require.Equal(b, sdk.NewDecWithPrec(5, 1), intents.MustGetForValoper(v10).Weight)

		inv3 := intents.MustGetForValoper(v3)
		require.Equal(b, sdk.ZeroDec(), inv3.Weight)
		benchSink = intents
	}

	if benchSink == nil {
		b.Fatal("Benchmark did not run!")
	}

	benchSink = nil
}
