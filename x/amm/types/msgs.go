package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgAddLiquidity(sender sdk.AccAddress, coins sdk.Coins) *MsgAddLiquidity {
	return &MsgAddLiquidity{
		Sender: sender.String(),
		Coins:  coins,
	}
}

func NewMsgRemoveLiquidity(sender sdk.AccAddress, share sdk.Coin) *MsgRemoveLiquidity {
	return &MsgRemoveLiquidity{
		Sender: sender.String(),
		Share:  share,
	}
}

func NewMsgSwapExactIn(sender sdk.AccAddress, coinIn, minCoinOut sdk.Coin) *MsgSwapExactIn {
	return &MsgSwapExactIn{
		Sender:     sender.String(),
		CoinIn:     coinIn,
		MinCoinOut: minCoinOut,
	}
}

func NewMsgSwapExactOut(sender sdk.AccAddress, coinOut, maxCoinIn sdk.Coin) *MsgSwapExactOut {
	return &MsgSwapExactOut{
		Sender:    sender.String(),
		CoinOut:   coinOut,
		MaxCoinIn: maxCoinIn,
	}
}
