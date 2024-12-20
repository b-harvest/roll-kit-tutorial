package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/b-harvest/roll-kit-tutorial/x/amm/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	bankKeeper types.BankKeeper

	// State
	Schema             collections.Schema
	NextPairIDSequence collections.Sequence
	Pairs              *collections.IndexedMap[collections.Pair[string, string], types.Pair, types.PairIndexes]
	Params             collections.Item[types.Params]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	bankKeeper types.BankKeeper,
) Keeper {

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:                cdc,
		bankKeeper:         bankKeeper,
		NextPairIDSequence: collections.NewSequence(sb, types.NextPairIDSequenceKey, "next_pair_id"),
		Pairs: collections.NewIndexedMap(sb, types.PairPrefix, "pairs",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.Pair](cdc),
			types.NewPairIndexes(sb),
		),
		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
