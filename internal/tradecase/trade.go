package tradecase

import (
	"context"
	"fmt"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type TradeCase struct {
	Repo     TradeRepo
	Provider TradeProvider
	Contract SmartContract
}

func New(
	r TradeRepo,
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
	repo TradeRepo,
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

	tradeMap, err := getTradeMap(
		pools,
	)
	if err != nil {
		return
	}

	tradePairs, err := getTradePairs(
		tradeMap,
	)
	if err != nil {
		return
	}

	prof, err := tc.GetProfitable(
		ctx,
		tradePairs,
	)
	if err != nil {
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

func (tc *TradeCase) GetProfitable(
	ctx context.Context,
	pairs []TradePair,
) (
	out []TradePair,
	err error,
) {
	for _, pair := range pairs {
		p, _err := tc.Contract.GetProfit(
			ctx,
			pair,
		)
		if err != nil {
			err = _err
			return
		}
		if p > 0 {
			out = append(out, pair)
		}
	}

	return
}

func getTradePairs(
	tradeMap map[TokenPair]map[string]TradePool,
) (
	tradePairs []TradePair,
	err error,
) {
	for _, poolMap := range tradeMap {
		n := len(poolMap)

		switch {
		case n < 2:
			continue
		case n == 2:
			pair, _err := makeTradePair(
				poolMap,
			)
			if err != nil {
				err = _err

				return
			}
			tradePairs = append(
				tradePairs,
				pair,
			)
		case n > 2:
			tradePairs = deleteDublicates(
				splitPoolsOnPairs(poolMap),
			)
		}
	}

	return
}

func getTradeMap(
	pools []TradePool,
) (
	trade map[TokenPair]map[string]TradePool,
	err error,
) {
	trade = make(map[TokenPair]map[string]TradePool)

	for _, pool := range pools {
		_, ok := trade[pool.Pair]
		if !ok {
			trade[pool.Pair] = make(
				map[string]TradePool,
				0,
			)
		}
		_, ok = trade[pool.Pair][pool.Address]
		if ok {
			continue
		}
		trade[pool.Pair][pool.Address] = pool
	}

	return
}

func splitPoolsOnPairs(
	pools map[string]TradePool,
) (
	pairs []TradePair,
) {
	addrList := make([]string, 0)

	for addr, pool := range pools {
		addrList = append(
			addrList,
			addr,
		)
		pairs = append(
			pairs,
			TradePair{Pool0: pool},
		)
	}

	for n, pair := range pairs {
		for i := len(addrList); i >= 0; i-- {
			if i == n {
				continue
			}
			pair.Pool1 = pools[addrList[i]]
		}
	}

	return
}

func deleteDublicates(
	pairs []TradePair,
) (
	out []TradePair,
) {
	out = pairs

	for n, pair := range pairs {
		if n != len(pairs) {
			for i, p := range pairs[n+1:] {
				if !(pair.Pool0 == p.Pool1 &&
					pair.Pool1 == p.Pool0) ||
					!(pair.Pool0 == p.Pool0 &&
						pair.Pool1 == p.Pool1) {
					continue
				} else {
					pairs = append(
						out[:i-1],
						out[i:]...,
					)
				}
			}
		}
	}

	return
}

func makeTradePair(
	pools map[string]TradePool,
) (
	pair TradePair,
	err error,
) {
	for _, pool := range pools {
		if pair.Pool0 != pool {
			pair.Pool0 = pool
		} else if pair.Pool1 != pool {
			pair.Pool1 = pool
		} else {
			err = fmt.Errorf(
				"pair with same pools %v",
				pool,
			)
			return
		}
	}

	return
}
