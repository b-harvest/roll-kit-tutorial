package types

import (
	"fmt"

	"cosmossdk.io/math"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramstypes.ParamSet = (*Params)(nil)

var (
	KeyFeeRate             = []byte("FeeRate")
	KeyMinInitialLiquidity = []byte("MinInitialLiquidity")

	DefaultFeeRate             = math.LegacyNewDecWithPrec(3, 3)
	DefaultMinInitialLiquidity = math.NewInt(1000)
)

// NewParams creates a new Params instance
func NewParams(feeRate math.LegacyDec, minInitialLiquidity math.Int) Params {
	return Params{
		FeeRate:             feeRate,
		MinInitialLiquidity: minInitialLiquidity,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultFeeRate,
		DefaultMinInitialLiquidity,
	)
}

func (params *Params) Validate() error {
	if err := validateFeeRate(params.FeeRate); err != nil {
		return err
	}
	if err := validateMinInitialLiquidity(params.MinInitialLiquidity); err != nil {
		return err
	}
	return nil
}

func validateFeeRate(v interface{}) error {
	feeRate, ok := v.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if feeRate.IsNegative() {
		return fmt.Errorf("invalid parameter value: %v", feeRate)
	}

	return nil
}

func validateMinInitialLiquidity(v interface{}) error {
	minInitialLiquidity, ok := v.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if minInitialLiquidity.IsNegative() {
		return fmt.Errorf("invalid parameter value: %v", minInitialLiquidity)
	}

	return nil
}

func (params *Params) String() string {
	out, _ := yaml.Marshal(params)
	return string(out)
}

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyFeeRate, &params.FeeRate, validateFeeRate),
		paramstypes.NewParamSetPair(KeyMinInitialLiquidity, &params.MinInitialLiquidity, validateMinInitialLiquidity),
	}
}
