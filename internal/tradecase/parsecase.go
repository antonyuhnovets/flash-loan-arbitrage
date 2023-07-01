package tradecase

import (
	"context"
	"fmt"
	"log"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type ParseCase struct {
	Parsers map[string]Parser
	Repo    TradeRepo
	Pools   []entities.TradePool
}

func NewParseCase(
	repo TradeRepo,
	parsers map[string]Parser,
) (
	pc ParseCase,
) {
	pc = ParseCase{
		Parsers: parsers,
		Repo:    repo,
		Pools:   make([]entities.TradePool, 0),
	}

	return
}

func (pc *ParseCase) SetProtocol(
	ctx context.Context,
	name string,
	parser Parser,
) (
	err error,
) {
	p, ok := pc.Parsers[name]
	if ok {
		if p == parser {
			err = fmt.Errorf("protocol already added")
			return
		}
		log.Println("changing protocol parser")
	}
	pc.Parsers[name] = parser

	return
}

func (pc *ParseCase) ParseAndStore(
	ctx context.Context,
) (
	err error,
) {
	pairs, err := pc.GetPairs(ctx)
	if err != nil {
		return
	}
	err = pc.ParsePairs(ctx, pairs)
	if err != nil {
		return
	}
	pools, err := pc.GetPools(ctx)
	if err != nil {
		return
	}
	err = pc.Repo.StorePools(ctx, "pools", pools)
	return
}

func (pc *ParseCase) ParsePairs(
	ctx context.Context,
	pairs []entities.TokenPair,
) (
	err error,
) {
	for _, v := range pc.Parsers {
		err = v.Parse(pairs)
		if err != nil {
			return
		}
	}
	return

}

func (pc *ParseCase) GetPools(
	ctx context.Context,
) (
	pools []entities.TradePool,
	err error,
) {
	for _, v := range pc.Parsers {
		pools = append(pools, v.ListPools()...)
	}
	if pools == nil {
		err = fmt.Errorf("Not found pools")
		return
	}
	pc.Pools = pools

	return
}

func (pc *ParseCase) GetPairs(
	ctx context.Context,
) (
	pairs []entities.TokenPair,
	err error,
) {
	tokens, err := pc.Repo.ListTokens(
		ctx,
		"tokens",
	)
	// fmt.Println(tokens)
	if err != nil {
		return
	}

	m := make(
		map[entities.Token]entities.TokenPair,
	)
	n := make(
		map[entities.TokenPair]int,
	)

	for index, token := range tokens {
		if _, ok := m[token]; !ok {
			for _, t := range tokens {
				if t != token {
					m[token] = entities.TokenPair{
						Token0: token,
						Token1: t,
					}
					n[m[token]] = index
				}
			}
		}
		continue
	}
	for k := range n {
		pairs = append(pairs, k)
	}

	return
}

func (pc *ParseCase) Clear(
	ctx context.Context,
) {
	pc.Pools = make([]entities.TradePool, 0)
}

func (pc ParseCase) ClearAll(
	ctx context.Context,
) {
	for _, v := range pc.Parsers {
		v.Clear()
	}
	pc.Clear(ctx)
}

func pairContain(
	pairs []entities.TokenPair,
	searchEl entities.TokenPair,
) (
	index int,
	ok bool,
) {
	ok = false

	for n, pair := range pairs {
		if pair == searchEl {
			ok = true
			index = n

			return
		}
	}

	return
}
