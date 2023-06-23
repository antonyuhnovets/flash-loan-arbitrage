package provider

import (
	c "context"
	"fmt"
	"math/big"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type ethClient interface {
	GetChainID() *big.Int
}

type TradeProvider struct {
	Client    ethClient
	PoolPairs []TradePair `json:"poolPairs"`
}

func NewTradeProvider(
	ctx c.Context, client ethClient, pairs ...TradePair,
) *TradeProvider {

	return &TradeProvider{
		Client:    client,
		PoolPairs: pairs,
	}
}

func (tp *TradeProvider) AddPair(
	ctx c.Context, pool TradePair,
) error {
	if _, ok := tp.FindPair(ctx, pool); ok {
		return fmt.Errorf("already in list")
	}
	tp.PoolPairs = append(tp.PoolPairs, pool)

	return nil
}

func (tp *TradeProvider) FindPair(
	ctx c.Context, pool TradePair,
) (
	int,
	bool,
) {
	for n, p := range tp.PoolPairs {
		if p == pool {
			return n, true
		}
	}

	return 0, false
}

func (tp *TradeProvider) RemovePair(
	ctx c.Context, pool TradePair,
) error {
	n, ok := tp.FindPair(ctx, pool)

	if !ok {
		return fmt.Errorf("not in list")
	}
	tp.PoolPairs = append(
		tp.PoolPairs[:n],
		tp.PoolPairs[n+1:]...,
	)

	return nil
}

func (tp *TradeProvider) GetPairs(
	ctx c.Context, protocol SwapProtocol, tokens TokenPair,
) (
	[]TradePair,
	bool,
) {
	var ok bool
	pairs := make([]TradePair, 0)

	for _, pair := range tp.PoolPairs {
		if checkPairTokens(pair, tokens) &&
			checkPairProtocol(pair, protocol) {

			pairs = append(pairs, pair)

			ok = true
		}
	}

	return pairs, ok
}

func (tp *TradeProvider) ListAllPairs(
	ctx c.Context,
) []TradePair {

	return tp.PoolPairs
}

func checkPairProtocol(
	pair TradePair, protocol SwapProtocol,
) bool {
	if pair.Pool0.SwapProtocol == protocol &&
		pair.Pool1.SwapProtocol == protocol {
		return true
	}

	return false
}

func checkPairTokens(
	pair TradePair, tokens TokenPair,
) bool {
	if pair.Pool0.Pair == tokens &&
		pair.Pool1.Pair == tokens {
		return true
	}

	return false
}
