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
