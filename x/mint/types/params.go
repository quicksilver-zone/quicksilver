package types

import (
	"errors"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	epochtypes "github.com/quicksilver-zone/quicksilver/v7/x/epochs/types"
)

// Parameter store keys.
var (
	KeyMintDenom                            = []byte("MintDenom")
	KeyGenesisEpochProvisions               = []byte("GenesisEpochProvisions")
	KeyEpochIdentifier                      = []byte("EpochIdentifier")
	KeyReductionPeriodInEpochs              = []byte("ReductionPeriodInEpochs")
	KeyReductionFactor                      = []byte("ReductionFactor")
	KeyPoolAllocationRatio                  = []byte("PoolAllocationRatio")
	KeyMintingRewardsDistributionStartEpoch = []byte("MintingRewardsDistributionStartEpoch")
)

// ParamKeyTable returns ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, genesisEpochProvisions sdkmath.LegacyDec, epochIdentifier string,
	reductionPeriodInEpochs int64, reductionFactor sdkmath.LegacyDec,
	distrProportions DistributionProportions, mintingRewardsDistributionStartEpoch int64,
) Params {
	return Params{
		MintDenom:                            mintDenom,
		GenesisEpochProvisions:               genesisEpochProvisions,
		EpochIdentifier:                      epochIdentifier,
		ReductionPeriodInEpochs:              reductionPeriodInEpochs,
		ReductionFactor:                      reductionFactor,
		DistributionProportions:              distrProportions,
		MintingRewardsDistributionStartEpoch: mintingRewardsDistributionStartEpoch,
	}
}

// DefaultParams returns default minting module parameters.
func DefaultParams() Params {
	return NewParams(
		sdk.DefaultBondDenom,
		sdkmath.LegacyNewDec(200000000/122),
		"day",
		365,
		sdkmath.LegacyNewDecWithPrec(75, 2),
		DistributionProportions{
			Staking:              sdkmath.LegacyNewDecWithPrec(3, 1), // 0.3
			PoolIncentives:       sdkmath.LegacyNewDecWithPrec(3, 1), // 0.3
			ParticipationRewards: sdkmath.LegacyNewDecWithPrec(3, 1), // 0.3
			CommunityPool:        sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
		},
		0,
	)
}

// Validate validates params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateGenesisEpochProvisions(p.GenesisEpochProvisions); err != nil {
		return err
	}
	if err := epochtypes.ValidateEpochIdentifierInterface(p.EpochIdentifier); err != nil {
		return err
	}
	if err := validateReductionPeriodInEpochs(p.ReductionPeriodInEpochs); err != nil {
		return err
	}
	if err := validateReductionFactor(p.ReductionFactor); err != nil {
		return err
	}
	if err := validateDistributionProportions(p.DistributionProportions); err != nil {
		return err
	}
	return validateMintingRewardsDistributionStartEpoch(p.MintingRewardsDistributionStartEpoch)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyGenesisEpochProvisions, &p.GenesisEpochProvisions, validateGenesisEpochProvisions),
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, epochtypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(KeyReductionPeriodInEpochs, &p.ReductionPeriodInEpochs, validateReductionPeriodInEpochs),
		paramtypes.NewParamSetPair(KeyReductionFactor, &p.ReductionFactor, validateReductionFactor),
		paramtypes.NewParamSetPair(KeyPoolAllocationRatio, &p.DistributionProportions, validateDistributionProportions),
		paramtypes.NewParamSetPair(KeyMintingRewardsDistributionStartEpoch, &p.MintingRewardsDistributionStartEpoch, validateMintingRewardsDistributionStartEpoch),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	return sdk.ValidateDenom(v)
}

func validateGenesisEpochProvisions(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdkmath.LegacyZeroDec()) {
		return errors.New("genesis epoch provision must be non-negative")
	}

	return nil
}

func validateReductionPeriodInEpochs(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateReductionFactor(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdkmath.LegacyOneDec()) {
		return errors.New("reduction factor cannot be greater than 1")
	}

	if v.IsNegative() {
		return errors.New("reduction factor cannot be negative")
	}

	return nil
}

func validateDistributionProportions(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Staking.IsNegative() {
		return errors.New("staking distribution ratio should not be negative")
	}

	if v.PoolIncentives.IsNegative() {
		return errors.New("pool incentives distribution ratio should not be negative")
	}

	if v.ParticipationRewards.IsNegative() {
		return errors.New("participation rewards distribution ratio should not be negative")
	}

	if v.CommunityPool.IsNegative() {
		return errors.New("community pool distribution ratio should not be negative")
	}

	totalProportions := v.Staking.Add(v.PoolIncentives).Add(v.CommunityPool).Add(v.ParticipationRewards)
	if !totalProportions.Equal(sdk.OneDec()) {
		return errors.New("total distributions ratio should be 1")
	}

	return nil
}

func validateMintingRewardsDistributionStartEpoch(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return errors.New("start epoch must be non-negative")
	}

	return nil
}
