package trade

import (
	"context"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/pairs"
)

type TradeCase struct {
	Repo     Repository
	Provider TradeProvider
	Contract SmartContract
}

func New(
	r Repository,
	p TradeProvider,
	c SmartContract,
) (
	tc *TradeCase,
) {
	tc = &TradeCase{
		Repo:     r,
		Provider: p,
		Contract: c,
	}

	return
}

func (tc *TradeCase) SetContract() (
	contract SmartContract,
) {
	tc.Contract = contract

	return
}

func (tc *TradeCase) SetProvider() (
	provider TradeProvider,
) {
	tc.Provider = provider

	return
}

func (tc *TradeCase) SetRepo() (
	repo Repository,
) {
	tc.Repo = repo

	return
}

func (tc *TradeCase) SetTokens(
	ctx context.Context,
	where string,
) (
	err error,
) {
	tokens, err := tc.Repo.ListTokens(ctx, where)
	if err != nil {
		return
	}
	err = tc.Provider.SetTokens(ctx, tokens)

	return
}

func (tc *TradeCase) SetProfitablePairs(
	ctx context.Context, where string,
) (
	err error,
) {
	pools, err := tc.Repo.ListPools(
		ctx,
		where,
	)
	if err != nil {
		return
	}

	tradeMap, err := pairs.GetTradeMap(
		pools,
	)
	if err != nil {
		return
	}

	tradePairs, err := pairs.GetTradePairs(
		tradeMap,
	)
	if err != nil {
		return
	}

	prof, ok, err := tc.GetProfitable(
		ctx,
		tradePairs,
	)
	if err != nil || !ok {
		return
	}

	err = tc.Contract.SetPairs(
		ctx,
		prof,
	)
	if err != nil {
		return
	}

	return
}

func (tc *TradeCase) GetProfitable(ctx context.Context, from []entities.TradePair) (
	out []entities.TradePair,
	ok bool,
	err error,
) {
	ok = false
	for _, pair := range from {
		prof, _, _err := tc.GetProfit(
			ctx,
			pair.Pool0.Address,
			pair.Pool1.Address,
		)
		if err != nil {
			err = _err

			return
		}
		if prof > 0 {
			out = append(out, pair)
			ok = true
		}
	}
	return
}
