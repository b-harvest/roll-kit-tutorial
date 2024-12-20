package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultParams(), 0, nil)
}

func NewGenesisState(params Params, pairSequence uint64, pairs []Pair) *GenesisState {
	return &GenesisState{
		Params:       params,
		PairSequence: pairSequence,
		Pairs:        pairs,
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
