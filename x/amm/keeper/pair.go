package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/b-harvest/roll-kit-tutorial/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddLiquidity(ctx context.Context, fromAddr sdk.AccAddress, coins sdk.Coins) (mintedShare sdk.Coin, err error) {
	coin0, coin1 := coins[0], coins[1]
	denom0, denom1 := coin0.Denom, coin1.Denom

	pair, err := k.GetPairByDenoms(ctx, denom0, denom1)

	if err != nil {
		pairId, err := k.NextPairIDSequence.Next(ctx)
		if err != nil {
			panic(err)
		}
		pair = types.NewPair(pairId, denom0, denom1)
		k.Pairs.Set(ctx, collections.Join(denom0, denom1), pair)
	}

	reserveAddr := types.PairReserveAddress(pair)
	shareDenom := types.ShareDenom(pair)

	reserveBalances := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx := reserveBalances.AmountOf(denom0)
	ry := reserveBalances.AmountOf(denom1)
	x := coin0.Amount
	y := coin1.Amount
	totalShare := k.bankKeeper.GetSupply(ctx, shareDenom).Amount

	var ax, ay, share math.Int
	if totalShare.IsZero() {
		var l math.LegacyDec
		l, err = math.LegacyNewDecFromInt(x.Mul(y)).ApproxSqrt()
		if err != nil {
			return
		}
		share = l.TruncateInt()

		params, err := k.Params.Get(ctx)
		if err != nil {
			panic(err)
		}

		if share.LT(params.MinInitialLiquidity) {
			err = errorsmod.Wrapf(
				types.ErrInsufficientLiquidity, "insufficient initial liquidity: %s", share)
			return sdk.Coin{}, err
		}
		ax = x
		ay = y
	} else {
		share = math.MinInt(totalShare.Mul(x).Quo(rx), totalShare.Mul(y).Quo(ry))
		ax = rx.Mul(share).Quo(totalShare)
		ay = ry.Mul(share).Quo(totalShare)
	}
	if !ax.IsPositive() || !ay.IsPositive() || !share.IsPositive() {
		err = types.ErrInsufficientLiquidity
		return
	}

	amt := sdk.NewCoins(sdk.NewCoin(denom0, ax), sdk.NewCoin(denom1, ay))
	if err = k.bankKeeper.SendCoins(ctx, fromAddr, reserveAddr, amt); err != nil {
		return
	}
	mintedShare = sdk.NewCoin(shareDenom, share)
	if err = k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(mintedShare)); err != nil {
		return
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, fromAddr, sdk.NewCoins(mintedShare)); err != nil {
		return
	}
	return mintedShare, nil
}

func (k Keeper) RemoveLiquidity(ctx context.Context, fromAddr sdk.AccAddress, share sdk.Coin) (withdrawnCoins sdk.Coins, err error) {
	var pairID uint64
	pairID, err = types.ParseShareDenom(share.Denom)
	if err != nil {
		return
	}

	pair, err := k.GetPairById(ctx, pairID)
	if err != nil {
		return
	}

	reserveAddr := types.PairReserveAddress(pair)
	reserveBalances := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx := reserveBalances.AmountOf(pair.Denom0)
	ry := reserveBalances.AmountOf(pair.Denom1)
	totalShare := k.bankKeeper.GetSupply(ctx, share.Denom).Amount

	var wx, wy math.Int
	if share.Amount.Equal(totalShare) {
		wx = rx
		wy = ry
		k.Pairs.Remove(ctx, collections.Join(pair.Denom0, pair.Denom1))
	} else {
		wx = rx.Mul(share.Amount).Quo(totalShare)
		wy = ry.Mul(share.Amount).Quo(totalShare)
	}
	if !wx.IsPositive() && !wy.IsPositive() {
		err = errorsmod.Wrap(types.ErrInsufficientLiquidity, "too small share to remove")
		return
	}

	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, fromAddr, types.ModuleName, sdk.NewCoins(share)); err != nil {
		return
	}
	if err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(share)); err != nil {
		return
	}
	withdrawnCoins = sdk.NewCoins(sdk.NewCoin(pair.Denom0, wx), sdk.NewCoin(pair.Denom1, wy))
	if err = k.bankKeeper.SendCoins(ctx, reserveAddr, fromAddr, withdrawnCoins); err != nil {
		return
	}
	return withdrawnCoins, nil
}

func (k Keeper) SwapExactIn(ctx context.Context, fromAddr sdk.AccAddress, coinIn, minCoinOut sdk.Coin) (coinOut sdk.Coin, err error) {
	pair, err := k.GetPairByDenoms(ctx, coinIn.Denom, minCoinOut.Denom)
	if err != nil {
		err = errorsmod.Wrap(err, "pair not found")
		return
	}

	reserveAddr := types.PairReserveAddress(pair)
	reserveBalances := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx := reserveBalances.AmountOf(pair.Denom0)
	ry := reserveBalances.AmountOf(pair.Denom1)

	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}
	feeRate := params.FeeRate

	var reserveIn, reserveOut math.Int
	amtInWithoutFee := math.LegacyNewDecFromInt(coinIn.Amount).
		MulTruncate(math.LegacyOneDec().Sub(feeRate)).TruncateInt()
	if coinIn.Denom == pair.Denom0 { // x to y
		reserveIn, reserveOut = rx, ry
		coinOut.Denom = pair.Denom1
	} else { // y to x
		reserveIn, reserveOut = ry, rx
		coinOut.Denom = pair.Denom0
	}
	coinOut.Amount = reserveOut.Mul(amtInWithoutFee).Quo(reserveIn.Add(amtInWithoutFee))
	if coinOut.Amount.LT(minCoinOut.Amount) {
		err = errorsmod.Wrapf(
			types.ErrSmallOutCoin, "%s is smaller than %s", coinOut.Amount, minCoinOut.Amount)
		return
	}

	if err = k.bankKeeper.SendCoins(ctx, fromAddr, reserveAddr, sdk.NewCoins(coinIn)); err != nil {
		return
	}
	if err = k.bankKeeper.SendCoins(ctx, reserveAddr, fromAddr, sdk.NewCoins(coinOut)); err != nil {
		return
	}
	return coinOut, nil
}

func (k Keeper) GetPairByDenoms(ctx context.Context, denomA, denomB string) (pair types.Pair, err error) {
	denom0, denom1 := types.SortDenoms(denomA, denomB)
	return k.Pairs.Get(ctx, collections.Join(denom0, denom1))
}

func (k Keeper) GetPairById(ctx context.Context, pairId uint64) (pair types.Pair, err error) {
	primaryKey, err := k.Pairs.Indexes.Id.MatchExact(ctx, pairId)
	if err != nil { // Sanity check
		return types.Pair{}, err
	}
	return k.Pairs.Get(ctx, primaryKey)
}

func (k Keeper) IterateAllPairs(ctx context.Context, cb func(pool types.Pair) bool) {
	err := k.Pairs.Walk(ctx, nil, func(key collections.Pair[string, string], value types.Pair) (stop bool, err error) {
		return cb(value), nil
	})
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetAllPairs(ctx context.Context) []types.Pair {
	pairs := make([]types.Pair, 0)
	k.IterateAllPairs(ctx, func(pair types.Pair) bool {
		pairs = append(pairs, pair)
		return false
	})
	return pairs
}
