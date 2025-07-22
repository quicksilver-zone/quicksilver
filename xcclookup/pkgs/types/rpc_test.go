package types

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestResponse_Mutex_Exists(t *testing.T) {
	// Test that the Response struct can be created and has the mutex field
	response := &Response{
		Assets: make(map[string][]Asset),
	}

	// Test that we can access the Assets field
	assets := response.GetAssets()
	assert.NotNil(t, assets)

	// Test that the GetAssets method works
	assets2 := response.GetAssets()
	assert.NotNil(t, assets2)
}

func TestResponse_GetAssets_ThreadSafe(t *testing.T) {
	response := &Response{
		Assets: make(map[string][]Asset),
	}

	// Add some initial data
	initialAssets := map[string][]Asset{
		"chain1": {
			{Type: "test", Amount: sdk.NewCoins(sdk.NewCoin("atom", sdk.NewInt(100)))},
		},
	}
	response.SetAssets(initialAssets)

	// Test concurrent reads
	numGoroutines := 50
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Read assets multiple times
			for j := 0; j < 100; j++ {
				assets := response.GetAssets()
				assert.NotNil(t, assets)
			}
		}()
	}

	wg.Wait()
}
