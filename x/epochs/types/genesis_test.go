package types

import (
	"strings"
	"testing"
	"time"
)

func TestValidateWithBetterContext(t *testing.T) {
	testCases := []struct {
		name    string
		epochs  []EpochInfo
		wantErr string
	}{
		{
			name:    "empty identifier",
			wantErr: "should NOT be empty",
			epochs:  []EpochInfo{{Identifier: ""}},
		},
		{
			name:    "duplicate identifiers",
			wantErr: `value #2: epoch identifier should be unique, got duplicate "day"`,
			epochs: []EpochInfo{
				{
					Identifier: "day",
					Duration:   time.Hour * 24,
				},
				{Identifier: "day"},
			},
		},
		{
			name:    "invalid duration: 0",
			wantErr: `value #1, Identifier: "day": epoch duration should be >0`,
			epochs: []EpochInfo{
				{
					Identifier: "day",
					StartTime:  time.Time{},
					Duration:   0,
				},
			},
		},
		{
			name:    "invalid duration: -2",
			wantErr: `value #1, Identifier: "day": epoch duration should be >0`,
			epochs: []EpochInfo{
				{
					Identifier: "day",
					Duration:   -2,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			gs := NewGenesisState(tc.epochs)
			err := gs.Validate()
			if err == nil {
				t.Fatal("Expecting a non-nil error")
			}

			if tc.wantErr == "" {
				t.Fatal("WantErr is unexpectedly empty!")
			}
			if g, w := err.Error(), tc.wantErr; !strings.Contains(g, w) {
				t.Errorf("Error mismatch\n\twant substring: %q\n\tGot: %s", w, g)
			}
		})
	}
}
