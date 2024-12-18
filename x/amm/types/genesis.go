package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultParams())
}

func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}
