package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "amm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// KVStore keys
var (
	ParamsKey = collections.NewPrefix(0)
)
