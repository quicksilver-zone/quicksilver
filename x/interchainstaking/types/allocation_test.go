package types_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// Delegation Allocation Tests

// Test with valid inputs and expected outputs
func TestValidInputs(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(100),
		"validator2": sdkmath.NewInt(200),
	}
	currentSum := sdkmath.NewInt(300)
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(5, 1)},
	}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(1000))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(550),
		"validator2": sdkmath.NewInt(450),
	}

	result, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.NoError(t, err)

	if !reflect.DeepEqual(result, expectedAllocations) {
		t.Errorf("Expected allocations %v, but got %v", expectedAllocations, result)
	}
}

// Test with minimum input values - fail on zero amount
func TestMinimumInputs(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{}
	currentSum := sdkmath.ZeroInt()
	targetAllocations := types.ValidatorIntents{}
	amount := sdk.Coins{sdk.NewCoin("token", sdkmath.ZeroInt())}

	_, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.Error(t, err)
}

// Test with maximum input values
func TestMaximumInputs(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(1000000000),
		"validator2": sdkmath.NewInt(2000000000),
	}
	currentSum := sdkmath.NewInt(3000000000)
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(5, 1)},
	}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(10000000000))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(5500000000),
		"validator2": sdkmath.NewInt(4500000000),
	}

	result, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.NoError(t, err)

	if !reflect.DeepEqual(result, expectedAllocations) {
		t.Errorf("Expected allocations %v, but got %v", expectedAllocations, result)
	}
}

// Test with empty currentAllocations
func TestEmptyCurrentAllocations(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{}
	currentSum := sdkmath.ZeroInt()
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(5, 1)},
	}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(1000))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(500),
		"validator2": sdkmath.NewInt(500),
	}

	result, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.NoError(t, err)

	if !reflect.DeepEqual(result, expectedAllocations) {
		t.Errorf("Expected allocations %v, but got %v", expectedAllocations, result)
	}
}

// Test with empty targetAllocations - error on empty target list.
func TestEmptyTargetAllocations(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(100),
		"validator2": sdkmath.NewInt(200),
	}
	currentSum := sdkmath.NewInt(300)
	targetAllocations := types.ValidatorIntents{}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(1000))}

	_, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.Error(t, err)
}

// Test non-equal targetAllocations
func TestNonEqualTargetAllocations(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(100),
		"validator2": sdkmath.NewInt(200),
	}
	currentSum := sdkmath.NewInt(300)
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(3, 1)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(7, 1)},
	}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(1000))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(290),
		"validator2": sdkmath.NewInt(710),
	}

	result, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.NoError(t, err)

	if !reflect.DeepEqual(result, expectedAllocations) {
		t.Errorf("Expected allocations %v, but got %v", expectedAllocations, result)
	}
}

// Test with some validators having zero weight
func TestValidInputsWithZeroWeight(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(100),
		"validator2": sdkmath.NewInt(200),
		"validator3": sdkmath.NewInt(0),
	}
	currentSum := sdkmath.NewInt(300)
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator3", Weight: sdk.NewDec(0)},
	}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(1000))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(550),
		"validator2": sdkmath.NewInt(450),
	}

	result, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.NoError(t, err)

	if !reflect.DeepEqual(result, expectedAllocations) {
		t.Errorf("Expected allocations %v, but got %v", expectedAllocations, result)
	}
}

// Test with targetAllocations having more validators than currentAllocations
func TestTargetAllocationsMoreValidators(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(100),
		"validator2": sdkmath.NewInt(200),
	}
	currentSum := sdkmath.NewInt(300)
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator3", Weight: sdk.NewDecWithPrec(5, 1)},
		{ValoperAddress: "validator4", Weight: sdk.NewDecWithPrec(5, 1)},
	}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(1000))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(241),
		"validator2": sdkmath.NewInt(195),
		"validator3": sdkmath.NewInt(282),
		"validator4": sdkmath.NewInt(282),
	}

	result, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, make(map[string]sdkmath.Int))
	require.NoError(t, err)

	if !reflect.DeepEqual(result, expectedAllocations) {
		t.Errorf("Expected allocations %v, but got %v", expectedAllocations, result)
	}
}

// Test with currentAllocations having more validators than targetAllocations
func TestCurrentAllocationsMoreValidators(t *testing.T) {
	currentAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(100),
		"validator2": sdkmath.NewInt(200),
		"validator3": sdkmath.NewInt(300),
	}
	currentSum := sdkmath.NewInt(600)
	targetAllocations := types.ValidatorIntents{
		{ValoperAddress: "validator1", Weight: sdk.NewDecWithPrec(43, 2)},
		{ValoperAddress: "validator2", Weight: sdk.NewDecWithPrec(57, 2)},
	}
	amount := sdk.Coins{sdk.NewCoin("token", sdk.NewInt(1000))}

	expectedAllocations := map[string]sdkmath.Int{
		"validator1": sdkmath.NewInt(453),
		"validator2": sdkmath.NewInt(547),
	}

	result, err := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations.Normalize(), amount, make(map[string]sdkmath.Int))
	require.NoError(t, err)

	if !reflect.DeepEqual(result, expectedAllocations) {
		t.Errorf("Expected allocations %v, but got %v", expectedAllocations, result)
	}
}

// Deltas tests.
