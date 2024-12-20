package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
)

const (
	// ModuleName defines the module name
	ModuleName = "amm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// KVStore keys
var (
	ParamsKey             = collections.NewPrefix(0)
	NextPairIDSequenceKey = collections.NewPrefix(1)
	PairPrefix            = collections.NewPrefix(2)
	PairIndexPrefix       = collections.NewPrefix(3)
)

type PairIndexes struct {
	Id *indexes.Unique[uint64, collections.Pair[string, string], Pair]
}

func (pi PairIndexes) IndexesList() []collections.Index[collections.Pair[string, string], Pair] {
	return []collections.Index[collections.Pair[string, string], Pair]{
		pi.Id,
	}
}

func NewPairIndexes(sb *collections.SchemaBuilder) PairIndexes {
	return PairIndexes{
		Id: indexes.NewUnique(
			sb, PairIndexPrefix, "pair_by_id",
			collections.Uint64Key, collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			func(primaryKey collections.Pair[string, string], v Pair) (uint64, error) {
				return v.Id, nil
			},
		),
	}
}
