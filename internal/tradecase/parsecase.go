package tradecase

import (
	"context"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type ParseCase struct {
	Parser
	Repository
}

func NewParseCase(
	repo Repository,
	parse Parser,
) (
	pc ParseCase,
) {
	pc = ParseCase{
		parse, repo,
	}

	return
}

func (pc *ParseCase) SwitchProtocol(
	ctx context.Context,
	sp entities.SwapProtocol,
) (
	err error,
) {
	pc.Parser.SetProtocol(sp)

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
	err = pc.StorePools(ctx, "pools", pools)

	return
}

func (pc *ParseCase) ParsePairs(
	ctx context.Context,
	pairs []entities.TokenPair,
) (
	err error,
) {

	err = pc.Parser.Parse(pairs)

	return

}

func (pc *ParseCase) GetPools(
	ctx context.Context,
) (
	pools []entities.Pool,
	err error,
) {
	pools = pc.Parser.ListPools()
	err = nil

	return
}

func (pc *ParseCase) GetPairs(
	ctx context.Context,
) (
	pairs []entities.TokenPair,
	err error,
) {
	tokens, err := pc.Repository.ListTokens(
		ctx,
		"tokens",
	)
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
	pc.Parser.Clear()
}
