package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrWrongCoinNumber = errorsmod.Register(ModuleName, 2, "wrong number of coins")
)
