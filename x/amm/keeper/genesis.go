package keeper

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/b-harvest/roll-kit-tutorial/x/amm/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		panic(err)
	}
	if err := k.NextPairIDSequence.Set(ctx, genState.PairSequence); err != nil {
		panic(err)
	}
	for _, pair := range genState.Pairs {
		if err := k.Pairs.Set(ctx, collections.Join(pair.Denom0, pair.Denom1), pair); err != nil {
			panic(err)
		}
	}

}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	sequence, err := k.NextPairIDSequence.Peek(ctx)
	if err != nil {
		panic(err)
	}

	pairs := k.GetAllPairs(ctx)

	return types.NewGenesisState(
		params,
		sequence,
		pairs,
	)
}
