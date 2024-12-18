package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramstypes.ParamSet = (*Params)(nil)

func NewParams() Params {
	return Params{}
}

func DefaultParams() Params {
	return NewParams()
}
func (params *Params) Validate() error {
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
	return paramstypes.ParamSetPairs{}
}
