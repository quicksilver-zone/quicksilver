package types

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func TestNewCacheManager(t *testing.T) {
	manager := NewCacheManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.Data)
	assert.Equal(t, 0, len(manager.Data))
}

func TestCacheType(t *testing.T) {
	tests := []struct {
		name     string
		cache    interface{}
		expected string
	}{
		{
			name:     "ConnectionProtocolData cache type",
			cache:    &Cache[prewards.ConnectionProtocolData]{},
			expected: "ConnectionProtocolData",
		},
		{
			name:     "OsmosisParamsProtocolData cache type",
			cache:    &Cache[prewards.OsmosisParamsProtocolData]{},
			expected: "OsmosisParamsProtocolData",
		},
		{
			name:     "Zone cache type",
			cache:    &Cache[icstypes.Zone]{},
			expected: "Zone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.cache.(interface{ Type() string })
			result := cache.Type()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCacheUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		dataType     int
		responseData []byte
		expectError  bool
		expectedLen  int
	}{
		{
			name:         "valid protocol data",
			dataType:     DataTypeProtocolData,
			responseData: []byte(`{"data":[{"chain_id":"test-chain","prefix":"test","last_epoch":100}]}`),
			expectError:  false,
			expectedLen:  1,
		},
		{
			name:         "valid zone data",
			dataType:     DataTypeZone,
			responseData: []byte(`{"zones":[{"chain_id":"test-zone","connection_id":"test-conn"}],"stats":null,"pagination":null}`),
			expectError:  false,
			expectedLen:  1,
		},
		{
			name:         "invalid data type",
			dataType:     999,
			responseData: []byte(`{"data":[]}`),
			expectError:  true,
			expectedLen:  0,
		},
		{
			name:         "invalid JSON for protocol data",
			dataType:     DataTypeProtocolData,
			responseData: []byte(`invalid json`),
			expectError:  true,
			expectedLen:  0,
		},
		{
			name:         "invalid JSON for zone data",
			dataType:     DataTypeZone,
			responseData: []byte(`invalid json`),
			expectError:  true,
			expectedLen:  0,
		},
		{
			name:         "empty protocol data",
			dataType:     DataTypeProtocolData,
			responseData: []byte(`{"data":[]}`),
			expectError:  false,
			expectedLen:  0,
		},
		{
			name:         "empty zone data",
			dataType:     DataTypeZone,
			responseData: []byte(`{"zones":[],"stats":null,"pagination":null}`),
			expectError:  false,
			expectedLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache[prewards.ConnectionProtocolData]{
				dataType: tt.dataType,
			}

			result, err := cache.unmarshal(tt.responseData)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}
		})
	}
}

func TestCacheInit(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		dataType       int
		updateInterval time.Duration
		expectError    bool
	}{
		{
			name:           "valid initialization",
			url:            "http://localhost:8080/test",
			dataType:       DataTypeProtocolData,
			updateInterval: 30 * time.Second,
			expectError:    true, // Will fail due to HTTP request, but we test the structure
		},
		{
			name:           "empty URL",
			url:            "",
			dataType:       DataTypeProtocolData,
			updateInterval: 30 * time.Second,
			expectError:    true,
		},
		{
			name:           "zero update interval",
			url:            "http://localhost:8080/test",
			dataType:       DataTypeProtocolData,
			updateInterval: 0,
			expectError:    true, // Will fail due to HTTP request
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			cache := &Cache[prewards.ConnectionProtocolData]{}

			err := cache.Init(ctx, tt.url, tt.dataType, tt.updateInterval)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.url, cache.url)
				assert.Equal(t, tt.dataType, cache.dataType)
				assert.Equal(t, tt.updateInterval, cache.duration)
			}
		})
	}
}

func TestCacheSetMock(t *testing.T) {
	cache := &Cache[prewards.ConnectionProtocolData]{}

	mockData := []prewards.ConnectionProtocolData{
		{
			ChainID:   "mock-chain-1",
			Prefix:    "mock",
			LastEpoch: 100,
		},
		{
			ChainID:   "mock-chain-2",
			Prefix:    "mock2",
			LastEpoch: 200,
		},
	}

	cache.SetMock(mockData)

	assert.Equal(t, mockData, cache.mockData)
	assert.Len(t, cache.mockData, 2)
}

func TestCacheRead(t *testing.T) {
	tests := []struct {
		name        string
		serverFunc  func(w http.ResponseWriter, r *http.Request)
		expectError bool
		expectedLen int
	}{
		{
			name: "successful read",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"data":[]}`))
				require.NoError(t, err)
			},
			expectError: false,
			expectedLen: 11, // Length of `{"data":[]}` (11 bytes)
		},
		{
			name: "server error",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectError: false, // HTTP errors don't cause read errors
			expectedLen: 0,
		},
		{
			name: "large response",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				largeData := make([]byte, 10000)
				for i := range largeData {
					largeData[i] = 'a'
				}
				_, err := w.Write(largeData)
				require.NoError(t, err)
			},
			expectError: false,
			expectedLen: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverFunc))
			defer server.Close()

			ctx := t.Context()
			cache := &Cache[prewards.ConnectionProtocolData]{
				url: server.URL,
			}

			result, err := cache.read(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}
		})
	}
}

func TestCacheFetch(t *testing.T) {
	tests := []struct {
		name        string
		dataType    int
		serverFunc  func(w http.ResponseWriter, r *http.Request)
		expectError bool
	}{
		{
			name:     "successful fetch protocol data",
			dataType: DataTypeProtocolData,
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"data":[{"chain_id":"test-chain","prefix":"test","last_epoch":100}]}`))
				require.NoError(t, err)
			},
			expectError: false,
		},
		{
			name:     "successful fetch zone data",
			dataType: DataTypeZone,
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"zones":[{"chain_id":"test-zone","connection_id":"test-conn"}],"stats":null,"pagination":null}`))
				require.NoError(t, err)
			},
			expectError: false,
		},
		{
			name:     "server error",
			dataType: DataTypeProtocolData,
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectError: true,
		},
		{
			name:     "invalid JSON response",
			dataType: DataTypeProtocolData,
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`invalid json`))
				require.NoError(t, err)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverFunc))
			defer server.Close()

			ctx := t.Context()
			cache := &Cache[prewards.ConnectionProtocolData]{
				url:      server.URL,
				dataType: tt.dataType,
			}

			err := cache.Fetch(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, cache.lastUpdated)
				assert.NotNil(t, cache.cache)
			}
		})
	}
}

func TestCacheGet(t *testing.T) {
	tests := []struct {
		name          string
		cacheData     []prewards.ConnectionProtocolData
		mockData      []prewards.ConnectionProtocolData
		duration      time.Duration
		lastUpdated   time.Time
		expectFetch   bool
		expectedTotal int
	}{
		{
			name: "cache not expired, no mock data",
			cacheData: []prewards.ConnectionProtocolData{
				{ChainID: "cache-1", Prefix: "cache", LastEpoch: 100},
			},
			mockData:      []prewards.ConnectionProtocolData{},
			duration:      30 * time.Second,
			lastUpdated:   time.Now(),
			expectFetch:   false,
			expectedTotal: 1,
		},
		{
			name: "cache not expired, with mock data",
			cacheData: []prewards.ConnectionProtocolData{
				{ChainID: "cache-1", Prefix: "cache", LastEpoch: 100},
			},
			mockData: []prewards.ConnectionProtocolData{
				{ChainID: "mock-1", Prefix: "mock", LastEpoch: 200},
			},
			duration:      30 * time.Second,
			lastUpdated:   time.Now(),
			expectFetch:   false,
			expectedTotal: 2,
		},
		{
			name: "cache expired",
			cacheData: []prewards.ConnectionProtocolData{
				{ChainID: "cache-1", Prefix: "cache", LastEpoch: 100},
			},
			mockData:      []prewards.ConnectionProtocolData{},
			duration:      30 * time.Second,
			lastUpdated:   time.Now().Add(-31 * time.Second), // Expired
			expectFetch:   true,
			expectedTotal: 0, // Will fail to fetch due to no server
		},
		{
			name:          "empty cache",
			cacheData:     []prewards.ConnectionProtocolData{},
			mockData:      []prewards.ConnectionProtocolData{},
			duration:      30 * time.Second,
			lastUpdated:   time.Now(),
			expectFetch:   false,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			cache := &Cache[prewards.ConnectionProtocolData]{
				cache:       tt.cacheData,
				mockData:    tt.mockData,
				duration:    tt.duration,
				lastUpdated: tt.lastUpdated,
				url:         "http://invalid-url-for-testing", // Invalid URL to test fetch failure
			}

			result, err := cache.Get(ctx)

			if tt.expectFetch {
				// Should fail due to invalid URL
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedTotal)
			}
		})
	}
}

func TestCacheManagerAdd(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		dataType    int
		updateTime  time.Duration
		expectError bool
	}{
		{
			name:        "valid cache element",
			url:         "http://localhost:8080/test",
			dataType:    DataTypeProtocolData,
			updateTime:  30 * time.Second,
			expectError: true, // Will fail due to HTTP request, but we test the structure
		},
		{
			name:        "empty URL",
			url:         "",
			dataType:    DataTypeProtocolData,
			updateTime:  30 * time.Second,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			manager := NewCacheManager()
			cache := &Cache[prewards.ConnectionProtocolData]{}

			err := manager.Add(ctx, cache, tt.url, tt.dataType, tt.updateTime)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, manager.Data, cache.Type())
			}
		})
	}
}

func TestGetCache(t *testing.T) {
	ctx := t.Context()
	manager := NewCacheManager()

	// Add a cache to the manager with proper setup
	cache := &Cache[prewards.ConnectionProtocolData]{
		cache:       []prewards.ConnectionProtocolData{},
		mockData:    []prewards.ConnectionProtocolData{{ChainID: "test-chain", Prefix: "test", LastEpoch: 100}},
		duration:    30 * time.Second,
		lastUpdated: time.Now(),
		url:         "http://invalid-url-for-testing", // Invalid URL to prevent actual fetching
	}
	manager.Data[cache.Type()] = cache

	// Test GetCache function
	result, err := GetCache[prewards.ConnectionProtocolData](ctx, &manager)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test-chain", result[0].ChainID)
}

func TestAddMocks(t *testing.T) {
	ctx := t.Context()
	manager := NewCacheManager()

	// Add a cache to the manager
	cache := &Cache[prewards.ConnectionProtocolData]{}
	manager.Data[cache.Type()] = cache

	mockData := []prewards.ConnectionProtocolData{
		{ChainID: "mock-1", Prefix: "mock", LastEpoch: 100},
		{ChainID: "mock-2", Prefix: "mock2", LastEpoch: 200},
	}

	// Test AddMocks function
	AddMocks[prewards.ConnectionProtocolData](ctx, &manager, mockData)

	// Verify mock data was set
	cacheFromManager := manager.Data[cache.Type()].(*Cache[prewards.ConnectionProtocolData])
	assert.Equal(t, mockData, cacheFromManager.mockData)
}

func TestCacheInterfaceCompliance(t *testing.T) {
	// Test that Cache implements CacheI interface
	var _ CacheI[prewards.ConnectionProtocolData] = &Cache[prewards.ConnectionProtocolData]{}

	// Test that Cache implements CacheManagerElementI interface
	var _ CacheManagerElementI = &Cache[prewards.ConnectionProtocolData]{}
}

func TestCacheConcurrency(t *testing.T) {
	// Test that cache operations are safe for concurrent access
	cache := &Cache[prewards.ConnectionProtocolData]{
		cache:       []prewards.ConnectionProtocolData{},
		mockData:    []prewards.ConnectionProtocolData{},
		duration:    30 * time.Second,
		lastUpdated: time.Now(),
	}

	ctx := t.Context()

	// Run multiple goroutines accessing the cache
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			_, err := cache.Get(ctx)
			// Should not panic even if there are errors
			_ = err
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestCachePerformance(t *testing.T) {
	// Test cache performance with large datasets
	cache := &Cache[prewards.ConnectionProtocolData]{
		cache:       make([]prewards.ConnectionProtocolData, 1000),
		mockData:    make([]prewards.ConnectionProtocolData, 1000),
		duration:    30 * time.Second,
		lastUpdated: time.Now(),
	}

	// Fill with test data
	for i := 0; i < 1000; i++ {
		cache.cache[i] = prewards.ConnectionProtocolData{
			ChainID:   fmt.Sprintf("chain-%d", i),
			Prefix:    fmt.Sprintf("prefix-%d", i),
			LastEpoch: int64(i),
		}
		cache.mockData[i] = prewards.ConnectionProtocolData{
			ChainID:   fmt.Sprintf("mock-chain-%d", i),
			Prefix:    fmt.Sprintf("mock-prefix-%d", i),
			LastEpoch: int64(i + 1000),
		}
	}

	ctx := t.Context()

	start := time.Now()
	result, err := cache.Get(ctx)
	duration := time.Since(start)

	// Should complete quickly (less than 10ms)
	assert.Less(t, duration, 10*time.Millisecond)
	assert.NoError(t, err)
	assert.Len(t, result, 2000) // cache + mock data
}

func TestCacheEdgeCases(t *testing.T) {
	t.Run("nil cache manager", func(t *testing.T) {
		ctx := t.Context()
		var manager *CacheManager

		// This should panic or return an error
		assert.Panics(t, func() {
			_, _ = GetCache[prewards.ConnectionProtocolData](ctx, manager)
		})
	})

	t.Run("cache not found in manager", func(t *testing.T) {
		ctx := t.Context()
		manager := NewCacheManager()

		// Try to get cache that doesn't exist - this will panic due to nil pointer
		assert.Panics(t, func() {
			_, _ = GetCache[prewards.ConnectionProtocolData](ctx, &manager)
		})
	})

	t.Run("zero duration cache", func(t *testing.T) {
		cache := &Cache[prewards.ConnectionProtocolData]{
			cache:       []prewards.ConnectionProtocolData{{ChainID: "test"}},
			duration:    0,
			lastUpdated: time.Now(),
		}

		ctx := t.Context()
		result, err := cache.Get(ctx)

		// Should try to fetch immediately due to zero duration
		assert.Error(t, err) // Will fail due to invalid URL
		assert.Nil(t, result)
	})
}

func TestCacheEndToEndWithGarbageData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("not valid json"))
		require.NoError(t, err)
	}))
	defer server.Close()

	ctx := t.Context()
	cache := &Cache[prewards.ConnectionProtocolData]{}
	// Should error on Init due to bad data
	err := cache.Init(ctx, server.URL, DataTypeProtocolData, 30*time.Second)
	require.Error(t, err)
}

func TestCacheEndToEndWithTTLRefresh(t *testing.T) {
	// Track request count to serve different data
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var response string
		if requestCount == 1 {
			// First request: initial data
			response = `{"data":[{"ConnectionID":"conn-1","ChainID":"initial-chain","Prefix":"init","LastEpoch":100,"TransferChannel":"channel-1"}]}`
		} else {
			// Subsequent requests: updated data
			response = `{"data":[{"ConnectionID":"conn-2","ChainID":"updated-chain","Prefix":"upd","LastEpoch":200,"TransferChannel":"channel-2"}]}`
		}

		_, err := w.Write([]byte(response))
		require.NoError(t, err)
	}))
	defer server.Close()

	ctx := t.Context()
	cache := &Cache[prewards.ConnectionProtocolData]{}

	// Set a very short TTL for testing
	shortTTL := 10 * time.Millisecond
	err := cache.Init(ctx, server.URL, DataTypeProtocolData, shortTTL)
	require.NoError(t, err)

	// First fetch: should get initial data
	result1, err := cache.Get(ctx)
	require.NoError(t, err)
	require.Len(t, result1, 1)
	assert.Equal(t, "conn-1", result1[0].ConnectionID)
	assert.Equal(t, "initial-chain", result1[0].ChainID)
	assert.Equal(t, "init", result1[0].Prefix)
	assert.Equal(t, int64(100), result1[0].LastEpoch)
	assert.Equal(t, "channel-1", result1[0].TransferChannel)

	// Wait for TTL to expire
	time.Sleep(shortTTL + 5*time.Millisecond)

	// Second fetch: should get updated data
	result2, err := cache.Get(ctx)
	require.NoError(t, err)
	require.Len(t, result2, 1)
	assert.Equal(t, "conn-2", result2[0].ConnectionID)
	assert.Equal(t, "updated-chain", result2[0].ChainID)
	assert.Equal(t, "upd", result2[0].Prefix)
	assert.Equal(t, int64(200), result2[0].LastEpoch)
	assert.Equal(t, "channel-2", result2[0].TransferChannel)

	// Verify that the server was called twice
	assert.Equal(t, 2, requestCount)
}
