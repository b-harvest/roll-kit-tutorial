package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/b-harvest/roll-kit-tutorial/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) AddLiquidity(ctx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address %s", err)
	}
	if err := msg.Coins.Validate(); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, err.Error())
	}
	if len(msg.Coins) != 2 {
		return nil, types.ErrWrongCoinNumber
	}

	mintedShare, err := m.Keeper.AddLiquidity(
		ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.Coins,
	)
	if err != nil {
		return nil, err
	}
	return &types.MsgAddLiquidityResponse{
		MintedShare: mintedShare,
	}, nil
}

func (m msgServer) RemoveLiquidity(ctx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %v", err)
	}
	if _, err := types.ParseShareDenom(msg.Share.Denom); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	withdrawnCoins, err := m.Keeper.RemoveLiquidity(
		ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.Share,
	)
	if err != nil {
		return nil, err
	}
	return &types.MsgRemoveLiquidityResponse{
		WithdrawnCoins: withdrawnCoins,
	}, nil
}
