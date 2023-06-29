package tradecase

import (
	"context"
	"fmt"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type TradeCase struct {
	repo     TradeRepo
	provider TradeProvider
	contract SmartContract
	parser   Parser
}

func New(r TradeRepo, p TradeProvider, c SmartContract,
) *TradeCase {

	return &TradeCase{
		repo:     r,
		provider: p,
		contract: c,
	}
}

func (tc *TradeCase) SetParser(p Parser) {
	tc.parser = p
}

func (tc *TradeCase) GetParser() Parser {
	return tc.parser
}

func (tc *TradeCase) GetContract() SmartContract {
	return tc.contract
}

func (tc *TradeCase) GetProvider() TradeProvider {
	return tc.provider
}

func (tc *TradeCase) GetRepo() TradeRepo {
	return tc.repo
}

func (tc *TradeCase) ParseWrite(ctx context.Context) error {
	tc.parser.Parse()

	poolMap := tc.parser.ListPools()

	for _, pools := range poolMap {
		err := tc.repo.StorePools(ctx, pools)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tc *TradeCase) SetUnknownTokens(ctx context.Context,
) error {
	for adr, token := range tc.contract.Tokens() {
		t, err := tc.repo.GetTokenByAddress(ctx, adr)
		if err != nil {
			return err
		}
		if t != token {
			tc.contract.Add(t)
		}
	}

	return nil
}

func (tc *TradeCase) SetProfitablePairs(ctx context.Context) error {
	pools, err := tc.repo.ListPools(ctx)
	if err != nil {
		return err
	}

	tradeMap, err := getTradeMap(pools)
	if err != nil {
		return err
	}

	tradePairs, err := getTradePairs(tradeMap)
	if err != nil {
		return err
	}

	prof, err := tc.GetProfitable(ctx, tradePairs)
	if err != nil {
		return err
	}

	err = tc.provider.SetPairs(ctx, prof)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TradeCase) GetProfitable(ctx context.Context, pairs []TradePair) ([]TradePair, error) {
	out := make([]TradePair, 0)

	for _, pair := range pairs {
		p, err := tc.contract.GetProfit(ctx, pair)
		if err != nil {
			return nil, err
		}
		if p > 0 {
			out = append(out, pair)
		}
	}

	return out, nil
}

func getTradePairs(tradeMap map[TokenPair]map[string]TradePool) ([]TradePair, error) {
	tradePairs := make([]TradePair, 0)

	for _, poolMap := range tradeMap {
		n := len(poolMap)
		switch {
		case n < 2:
			continue
		case n == 2:
			pair, err := makeTradePair(poolMap)
			if err != nil {
				return nil, err
			}
			tradePairs = append(tradePairs, pair)
		case n > 2:
			tradePairs = deleteDublicates(splitPoolsOnPairs(poolMap))
		}
	}

	return tradePairs, nil
}

func getTradeMap(pools []TradePool) (map[TokenPair]map[string]TradePool, error) {
	trade := make(map[TokenPair]map[string]TradePool, 0)

	for _, pool := range pools {
		if _, ok := trade[pool.Pair]; !ok {
			trade[pool.Pair] = make(map[string]TradePool, 0)
		}
		if _, ok := trade[pool.Pair][pool.Address]; ok {
			return nil, fmt.Errorf("pools with the same address found %s", pool.Address)
		}
		trade[pool.Pair][pool.Address] = pool
	}

	return trade, nil
}

func splitPoolsOnPairs(pools map[string]TradePool) []TradePair {
	pairs := make([]TradePair, 0)
	addrList := make([]string, 0)

	for addr, pool := range pools {
		addrList = append(addrList, addr)
		pairs = append(pairs, TradePair{Pool0: pool})
	}
	for n, pair := range pairs {
		for i := len(addrList); i >= 0; i-- {
			if i == n {
				continue
			}
			pair.Pool1 = pools[addrList[i]]
		}
	}

	return pairs
}

func deleteDublicates(pairs []TradePair) []TradePair {
	out := pairs

	for n, pair := range pairs {
		if n != len(pairs) {
			for i, p := range pairs[n+1:] {
				if !(pair.Pool0 == p.Pool1 && pair.Pool1 == p.Pool0) ||
					!(pair.Pool0 == p.Pool0 && pair.Pool1 == p.Pool1) {
					continue
				} else {
					pairs = append(out[:i-1], out[i:]...)
				}
			}
		}
	}

	return out
}

func makeTradePair(pools map[string]TradePool) (TradePair, error) {
	var pair TradePair

	for _, pool := range pools {
		if pair.Pool0 != pool {
			pair.Pool0 = pool
		} else if pair.Pool1 != pool {
			pair.Pool1 = pool
		} else {
			return TradePair{}, fmt.Errorf("pair with same pools %v", pool)
		}
	}

	return pair, nil
}
