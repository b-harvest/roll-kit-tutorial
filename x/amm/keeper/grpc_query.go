package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/b-harvest/roll-kit-tutorial/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = queryServer{}

func NewQueryServer(k Keeper) types.QueryServer {
	return queryServer{k: k}
}

type queryServer struct{ k Keeper }

// Params queries the parameters of the liquidity module.
func (s queryServer) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params, err := s.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

func (s queryServer) Pairs(c context.Context, req *types.QueryPairsRequest) (*types.QueryPairsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	pairs, pageRes, err := query.CollectionPaginate(
		c,
		s.k.Pairs,
		req.Pagination,
		func(_ collections.Pair[string, string], value types.Pair) (types.Pair, error) {
			return value, nil
		},
	)

	return &types.QueryPairsResponse{Pairs: pairs, Pagination: pageRes}, err
}

func (s queryServer) Pair(c context.Context, req *types.QueryPairRequest) (*types.QueryPairResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	pair, err := s.k.GetPairById(c, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "pair %d not found", req.Id)

	}
	return &types.QueryPairResponse{
		Pair: pair,
	}, nil
}
