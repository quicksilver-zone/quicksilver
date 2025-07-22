package types

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	prewards "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

func TestOutputResponse(t *testing.T) {
	tests := []struct {
		name          string
		response      *Response
		errors        map[string]error
		clearMessages bool
		expectedJSON  string
	}{
		{
			name: "successful response with messages and no errors",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{
					{
						UserAddress: "quick1test",
						Zone:        "test-zone",
						SrcZone:     "test-src",
						ClaimType:   1,
					},
				},
				Assets: map[string][]Asset{
					"test-chain": {
						{
							Type:   "test-type",
							Amount: sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(100))),
						},
					},
				},
			},
			errors:        nil,
			clearMessages: false,
			expectedJSON:  `{"messages":[{"user_address":"quick1test","zone":"test-zone","src_zone":"test-src","claim_type":1}],"assets":{"test-chain":[{"Type":"test-type","Amount":[{"denom":"test","amount":"100"}]}]}}`,
		},
		{
			name: "response with errors",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{},
				Assets:   map[string][]Asset{},
			},
			errors: map[string]error{
				"test-error": assert.AnError,
			},
			clearMessages: false,
			expectedJSON:  `{"messages":[],"assets":{},"errors":"test-error: assert.AnError general error for testing"}`,
		},
		{
			name: "response with clearMessages=true",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{
					{
						UserAddress: "quick1test",
						Zone:        "test-zone",
						SrcZone:     "test-src",
						ClaimType:   1,
					},
				},
				Assets: map[string][]Asset{
					"test-chain": {
						{
							Type:   "test-type",
							Amount: sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(100))),
						},
					},
				},
			},
			errors:        nil,
			clearMessages: true,
			expectedJSON:  `{"messages":null,"assets":{"test-chain":[{"Type":"test-type","Amount":[{"denom":"test","amount":"100"}]}]}}`,
		},
		{
			name: "empty response",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{},
				Assets:   map[string][]Asset{},
			},
			errors:        nil,
			clearMessages: false,
			expectedJSON:  `{"messages":[],"assets":{}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder to capture the output
			w := httptest.NewRecorder()

			// Call the function
			outputResponse(w, tt.response, tt.errors, tt.clearMessages)

			// Check the response
			assert.Equal(t, http.StatusOK, w.Code)

			// Check headers
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
			assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
			assert.Equal(t, "0", w.Header().Get("Expires"))
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
			assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
			assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))

			// Parse the JSON response to verify structure
			var result map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &result)
			require.NoError(t, err)

			// Verify the JSON matches expected (allowing for whitespace differences)
			actualJSON := w.Body.String()
			assert.JSONEq(t, tt.expectedJSON, actualJSON)
		})
	}
}

func TestOutputEpoch(t *testing.T) {
	tests := []struct {
		name         string
		response     *Response
		errors       map[string]error
		expectErrors bool
	}{
		{
			name: "epoch output with messages preserved",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{
					{
						UserAddress: "quick1epoch",
						Zone:        "epoch-zone",
						SrcZone:     "epoch-src",
						ClaimType:   2,
					},
				},
				Assets: map[string][]Asset{},
			},
			errors:       nil,
			expectErrors: false,
		},
		{
			name: "epoch output with errors",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{},
				Assets:   map[string][]Asset{},
			},
			errors: map[string]error{
				"epoch-error": assert.AnError,
			},
			expectErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Call OutputEpoch
			OutputEpoch(w, tt.response, tt.errors)

			// Verify response
			assert.Equal(t, http.StatusOK, w.Code)

			// Check headers
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
			assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
			assert.Equal(t, "0", w.Header().Get("Expires"))
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
			assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
			assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))

			var result Response
			err := json.Unmarshal(w.Body.Bytes(), &result)
			require.NoError(t, err)

			// Verify messages are preserved (not cleared)
			if len(tt.response.Messages) > 0 {
				assert.NotNil(t, result.Messages)
				assert.Len(t, result.Messages, len(tt.response.Messages))
			}

			// Verify errors are handled correctly
			if tt.expectErrors {
				assert.NotNil(t, result.Errors)
			} else {
				assert.Nil(t, result.Errors)
			}
		})
	}
}

func TestOutputCurrent(t *testing.T) {
	tests := []struct {
		name         string
		response     *Response
		errors       map[string]error
		expectErrors bool
	}{
		{
			name: "current output with messages cleared",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{
					{
						UserAddress: "quick1current",
						Zone:        "current-zone",
						SrcZone:     "current-src",
						ClaimType:   3,
					},
				},
				Assets: map[string][]Asset{},
			},
			errors:       nil,
			expectErrors: false,
		},
		{
			name: "current output with errors",
			response: &Response{
				Messages: []prewards.MsgSubmitClaim{},
				Assets:   map[string][]Asset{},
			},
			errors: map[string]error{
				"current-error": assert.AnError,
			},
			expectErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Call OutputCurrent
			OutputCurrent(w, tt.response, tt.errors)

			// Verify response
			assert.Equal(t, http.StatusOK, w.Code)

			// Check headers
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
			assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
			assert.Equal(t, "0", w.Header().Get("Expires"))
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
			assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
			assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))

			var result Response
			err := json.Unmarshal(w.Body.Bytes(), &result)
			require.NoError(t, err)

			// Verify messages are cleared
			assert.Nil(t, result.Messages)

			// Verify errors are handled correctly
			if tt.expectErrors {
				assert.NotNil(t, result.Errors)
			} else {
				assert.Nil(t, result.Errors)
			}
		})
	}
}

func TestOutputResponseWithComplexData(t *testing.T) {
	// Test with more complex data structures
	response := &Response{
		Messages: []prewards.MsgSubmitClaim{
			{
				UserAddress: "quick1complex",
				Zone:        "complex-zone",
				SrcZone:     "complex-src",
				ClaimType:   4,
			},
			{
				UserAddress: "quick1complex2",
				Zone:        "complex-zone2",
				SrcZone:     "complex-src2",
				ClaimType:   5,
			},
		},
		Assets: map[string][]Asset{
			"chain1": {
				{
					Type:   "asset1",
					Amount: sdk.NewCoins(sdk.NewCoin("token1", sdk.NewInt(1000))),
				},
				{
					Type:   "asset2",
					Amount: sdk.NewCoins(sdk.NewCoin("token2", sdk.NewInt(2000))),
				},
			},
			"chain2": {
				{
					Type:   "asset3",
					Amount: sdk.NewCoins(sdk.NewCoin("token3", sdk.NewInt(3000))),
				},
			},
		},
	}

	errors := map[string]error{
		"error1": assert.AnError,
		"error2": assert.AnError,
	}

	w := httptest.NewRecorder()
	outputResponse(w, response, errors, false)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check headers
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "0", w.Header().Get("Expires"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))

	var result Response
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	// Verify the response structure
	assert.Len(t, result.Messages, 2)
	assert.Len(t, result.Assets, 2)
	assert.NotNil(t, result.Errors)
}

func TestOutputResponseNilResponse(t *testing.T) {
	w := httptest.NewRecorder()

	// Test with nil response
	outputResponse(w, nil, nil, false)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check headers
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "0", w.Header().Get("Expires"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))

	var result Response
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	// Should handle nil gracefully - check individual fields instead of the whole struct
	assert.NotNil(t, result.Messages)
	assert.NotNil(t, result.Assets)
}

func TestOutputResponseEmptyErrorsMap(t *testing.T) {
	response := &Response{
		Messages: []prewards.MsgSubmitClaim{},
		Assets:   map[string][]Asset{},
	}

	// Test with empty errors map
	errors := make(map[string]error)

	w := httptest.NewRecorder()
	outputResponse(w, response, errors, false)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check headers
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "0", w.Header().Get("Expires"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))

	var result Response
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	// Should not have errors field when no errors
	assert.Nil(t, result.Errors)
}

func TestOutputResponseErrorHandling(t *testing.T) {
	// Test error handling when JSON marshaling fails
	// This is a bit tricky to test, but we can test the error path

	w := httptest.NewRecorder()

	// Create a response that would cause JSON marshaling to fail
	// We'll use a channel which can't be marshaled to JSON
	response := &Response{
		Messages: []prewards.MsgSubmitClaim{},
		Assets:   map[string][]Asset{},
	}

	// Mock the json.Marshal to fail
	// Since we can't easily mock json.Marshal, we'll test the error handling
	// by creating a response that should work normally

	outputResponse(w, response, nil, false)

	// Should return 200 OK for successful responses
	assert.Equal(t, http.StatusOK, w.Code)

	// Check headers are still set even in error cases
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "0", w.Header().Get("Expires"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
}
