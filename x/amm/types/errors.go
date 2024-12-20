package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrWrongCoinNumber       = errorsmod.Register(ModuleName, 2, "wrong number of coins")
	ErrInsufficientLiquidity = errorsmod.Register(ModuleName, 3, "insufficient liquidity")
	ErrSmallOutCoin          = errorsmod.Register(ModuleName, 4, "calculated out coin is smaller than the minimum")
	ErrBigInCoin             = errorsmod.Register(ModuleName, 5, "calculated in coin is bigger than the maximum")
)
