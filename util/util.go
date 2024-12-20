package util

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SampleAddress(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}
