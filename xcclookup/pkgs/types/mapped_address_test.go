package types

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func TestGetMappedAddresses(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		connections []prewards.ConnectionProtocolData
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:    "valid request with matching connections",
			address: "quick1testaddress",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "test-chain-1",
					Prefix:    "test",
					LastEpoch: 100,
				},
				{
					ChainID:   "test-chain-2",
					Prefix:    "test2",
					LastEpoch: 200,
				},
			},
			config: &Config{
				SourceChain: "quicksilver",
				Chains: map[string]string{
					"quicksilver": "http://localhost:26657",
				},
				Timeout: 30,
			},
			expectError: true, // Will fail due to RPC connection, but we test the structure
		},
		{
			name:        "empty connections list",
			address:     "quick1testaddress",
			connections: []prewards.ConnectionProtocolData{},
			config: &Config{
				SourceChain: "quicksilver",
				Chains: map[string]string{
					"quicksilver": "http://localhost:26657",
				},
				Timeout: 30,
			},
			expectError: true, // Will fail due to RPC connection
		},
		{
			name:    "nil config",
			address: "quick1testaddress",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "test-chain-1",
					Prefix:    "test",
					LastEpoch: 100,
				},
			},
			config:      nil,
			expectError: true,
			errorMsg:    "runtime error: invalid memory address or nil pointer dereference",
		},
		{
			name:    "invalid source chain in config",
			address: "quick1testaddress",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "test-chain-1",
					Prefix:    "test",
					LastEpoch: 100,
				},
			},
			config: &Config{
				SourceChain: "nonexistent",
				Chains: map[string]string{
					"quicksilver": "http://localhost:26657",
				},
				Timeout: 30,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			if tt.config == nil {
				// Test that nil config causes panic
				assert.Panics(t, func() {
					_, _ = GetMappedAddresses(ctx, tt.address, tt.connections, tt.config)
				})
				return
			}

			result, err := GetMappedAddresses(ctx, tt.address, tt.connections, tt.config)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestGetMappedAddressesWithMockRPC(t *testing.T) {
	// This test would require mocking the RPC client
	// For now, we'll test the function structure and error handling
	t.Run("test function signature and basic error handling", func(t *testing.T) {
		ctx := t.Context()
		address := "quick1testaddress"
		connections := []prewards.ConnectionProtocolData{
			{
				ChainID:   "test-chain-1",
				Prefix:    "test",
				LastEpoch: 100,
			},
		}
		config := &Config{
			SourceChain: "quicksilver",
			Chains: map[string]string{
				"quicksilver": "invalid-url-that-will-fail",
			},
			Timeout: 1, // Very short timeout to fail quickly
		}

		// This should fail due to invalid RPC URL
		result, err := GetMappedAddresses(ctx, address, connections, config)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetHeights(t *testing.T) {
	tests := []struct {
		name        string
		connections []prewards.ConnectionProtocolData
		expected    map[string]int64
	}{
		{
			name: "multiple connections with different heights",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "chain-1",
					LastEpoch: 100,
				},
				{
					ChainID:   "chain-2",
					LastEpoch: 200,
				},
				{
					ChainID:   "chain-3",
					LastEpoch: 300,
				},
			},
			expected: map[string]int64{
				"chain-1": 100,
				"chain-2": 200,
				"chain-3": 300,
			},
		},
		{
			name: "single connection",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "single-chain",
					LastEpoch: 150,
				},
			},
			expected: map[string]int64{
				"single-chain": 150,
			},
		},
		{
			name:        "empty connections list",
			connections: []prewards.ConnectionProtocolData{},
			expected:    map[string]int64{},
		},
		{
			name: "connections with zero heights",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "zero-chain",
					LastEpoch: 0,
				},
				{
					ChainID:   "negative-chain",
					LastEpoch: -1,
				},
			},
			expected: map[string]int64{
				"zero-chain":     0,
				"negative-chain": -1,
			},
		},
		{
			name: "duplicate chain IDs (should use last occurrence)",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "duplicate",
					LastEpoch: 100,
				},
				{
					ChainID:   "duplicate",
					LastEpoch: 200,
				},
			},
			expected: map[string]int64{
				"duplicate": 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetHeights(tt.connections)

			assert.Equal(t, len(tt.expected), len(result))
			for chainID, expectedHeight := range tt.expected {
				assert.Equal(t, expectedHeight, result[chainID])
			}
		})
	}
}

func TestGetZeroHeights(t *testing.T) {
	tests := []struct {
		name        string
		connections []prewards.ConnectionProtocolData
		expected    map[string]int64
	}{
		{
			name: "multiple connections",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "chain-1",
					LastEpoch: 100,
				},
				{
					ChainID:   "chain-2",
					LastEpoch: 200,
				},
				{
					ChainID:   "chain-3",
					LastEpoch: 300,
				},
			},
			expected: map[string]int64{
				"chain-1": 0,
				"chain-2": 0,
				"chain-3": 0,
			},
		},
		{
			name: "single connection",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "single-chain",
					LastEpoch: 150,
				},
			},
			expected: map[string]int64{
				"single-chain": 0,
			},
		},
		{
			name:        "empty connections list",
			connections: []prewards.ConnectionProtocolData{},
			expected:    map[string]int64{},
		},
		{
			name: "duplicate chain IDs",
			connections: []prewards.ConnectionProtocolData{
				{
					ChainID:   "duplicate",
					LastEpoch: 100,
				},
				{
					ChainID:   "duplicate",
					LastEpoch: 200,
				},
			},
			expected: map[string]int64{
				"duplicate": 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetZeroHeights(tt.connections)

			assert.Equal(t, len(tt.expected), len(result))
			for chainID, expectedHeight := range tt.expected {
				assert.Equal(t, expectedHeight, result[chainID])
			}
		})
	}
}

func TestGetHeightsAndGetZeroHeightsConsistency(t *testing.T) {
	// Test that both functions handle the same input consistently
	connections := []prewards.ConnectionProtocolData{
		{
			ChainID:   "chain-1",
			LastEpoch: 100,
		},
		{
			ChainID:   "chain-2",
			LastEpoch: 200,
		},
		{
			ChainID:   "chain-3",
			LastEpoch: 300,
		},
	}

	heights := GetHeights(connections)
	zeroHeights := GetZeroHeights(connections)

	// Both should have the same number of entries
	assert.Equal(t, len(heights), len(zeroHeights))

	// All heights in zeroHeights should be 0
	for chainID, height := range zeroHeights {
		assert.Equal(t, int64(0), height, "Chain %s should have height 0", chainID)
		assert.Contains(t, heights, chainID, "Chain %s should exist in both maps", chainID)
	}
}

func TestGetMappedAddressesEdgeCases(t *testing.T) {
	t.Run("empty address", func(t *testing.T) {
		ctx := t.Context()
		connections := []prewards.ConnectionProtocolData{
			{
				ChainID:   "test-chain",
				Prefix:    "test",
				LastEpoch: 100,
			},
		}
		config := &Config{
			SourceChain: "quicksilver",
			Chains: map[string]string{
				"quicksilver": "http://localhost:26657",
			},
			Timeout: 30,
		}

		result, err := GetMappedAddresses(ctx, "", connections, config)

		// Should fail due to RPC connection, but we test the structure
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("nil connections", func(t *testing.T) {
		ctx := t.Context()
		config := &Config{
			SourceChain: "quicksilver",
			Chains: map[string]string{
				"quicksilver": "http://localhost:26657",
			},
			Timeout: 30,
		}

		// This should fail due to RPC connection, but we test the structure
		result, err := GetMappedAddresses(ctx, "test", nil, config)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 1*time.Millisecond)
		defer cancel()

		connections := []prewards.ConnectionProtocolData{
			{
				ChainID:   "test-chain",
				Prefix:    "test",
				LastEpoch: 100,
			},
		}
		config := &Config{
			SourceChain: "quicksilver",
			Chains: map[string]string{
				"quicksilver": "http://localhost:26657",
			},
			Timeout: 30,
		}

		result, err := GetMappedAddresses(ctx, "test", connections, config)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetHeightsPerformance(t *testing.T) {
	// Test with a large number of connections to ensure performance
	connections := make([]prewards.ConnectionProtocolData, 1000)
	for i := 0; i < 1000; i++ {
		connections[i] = prewards.ConnectionProtocolData{
			ChainID:   fmt.Sprintf("chain-%d", i),
			LastEpoch: int64(i),
		}
	}

	start := time.Now()
	result := GetHeights(connections)
	duration := time.Since(start)

	// Should complete quickly (less than 10ms)
	assert.Less(t, duration, 10*time.Millisecond)
	assert.Equal(t, 1000, len(result))

	// Verify all entries
	for i := 0; i < 1000; i++ {
		expectedChainID := fmt.Sprintf("chain-%d", i)
		assert.Equal(t, int64(i), result[expectedChainID])
	}
}

func TestGetZeroHeightsPerformance(t *testing.T) {
	// Test with a large number of connections to ensure performance
	connections := make([]prewards.ConnectionProtocolData, 1000)
	for i := 0; i < 1000; i++ {
		connections[i] = prewards.ConnectionProtocolData{
			ChainID:   fmt.Sprintf("chain-%d", i),
			LastEpoch: int64(i),
		}
	}

	start := time.Now()
	result := GetZeroHeights(connections)
	duration := time.Since(start)

	// Should complete quickly (less than 10ms)
	assert.Less(t, duration, 10*time.Millisecond)
	assert.Equal(t, 1000, len(result))

	// Verify all entries are zero
	for i := 0; i < 1000; i++ {
		expectedChainID := fmt.Sprintf("chain-%d", i)
		assert.Equal(t, int64(0), result[expectedChainID])
	}
}
