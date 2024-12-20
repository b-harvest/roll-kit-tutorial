package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/b-harvest/roll-kit-tutorial/x/amm/types"
)

func (k Keeper) GetPairByDenoms(ctx context.Context, denomA, denomB string) (pair types.Pair, err error) {
	denom0, denom1 := types.SortDenoms(denomA, denomB)
	return k.Pairs.Get(ctx, collections.Join(denom0, denom1))
}

func (k Keeper) GetPairById(ctx context.Context, pairId uint64) (pair types.Pair, err error) {
	primaryKey, err := k.Pairs.Indexes.Id.MatchExact(ctx, pairId)
	if err != nil { // Sanity check
		return types.Pair{}, err
	}
	return k.Pairs.Get(ctx, primaryKey)
}

func (k Keeper) IterateAllPairs(ctx context.Context, cb func(pool types.Pair) bool) {
	err := k.Pairs.Walk(ctx, nil, func(key collections.Pair[string, string], value types.Pair) (stop bool, err error) {
		return cb(value), nil
	})
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetAllPairs(ctx context.Context) []types.Pair {
	pairs := make([]types.Pair, 0)
	k.IterateAllPairs(ctx, func(pair types.Pair) bool {
		pairs = append(pairs, pair)
		return false
	})
	return pairs
}
