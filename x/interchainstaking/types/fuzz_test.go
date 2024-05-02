package types_test

import (
	"testing"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func FuzzParseMemoFields(f *testing.F) {
	if testing.Short() {
		f.Skip("In -short mode")
	}

	// 1. Add the corpa seeds.
	seeds := [][]byte{
		// Valid sequences.
		{
			byte(types.FieldTypeAccountMap), 2, 1, 1,
		},
		{
			byte(types.FieldTypeAccountMap), 2, 1, 1,
			byte(types.FieldTypeReturnToSender), 0,
		},

		// Invalid sequences.
		{
			3, 2, 1, 1,
			byte(types.FieldTypeReturnToSender), 0,
		},
		{
			byte(types.FieldTypeAccountMap), 0,
			byte(types.FieldTypeReturnToSender), 0,
		},
		{
			byte(types.FieldTypeAccountMap), 3, 0, 0,
			byte(types.FieldTypeReturnToSender), 4, 1, 1, 1, 3,
		},
		{byte(types.FieldTypeAccountMap), 1, 0, 0},
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	// 2. Now run the fuzzers.
	f.Fuzz(func(t *testing.T, input []byte) {
		_, _ = types.ParseMemoFields(input)
	})
}
