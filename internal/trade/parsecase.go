package trade

import (
	"context"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/pairs"
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

func (pc *ParseCase) AddProtocol(
	ctx context.Context,
	sp entities.SwapProtocol,
) (
	err error,
) {
	err = pc.Parser.AddProtocol(sp)

	return
}

func (pc *ParseCase) RmProtocol(
	ctx context.Context,
	sp entities.SwapProtocol,
) (
	err error,
) {
	err = pc.RemoveProtocol(sp)

	return
}

func (pc *ParseCase) JustParse(
	ctx context.Context,
) (
	pools []entities.Pool,
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

	pools = pc.Parser.ListPools()

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
	p []entities.TokenPair,
	err error,
) {
	tokens, err := pc.Repository.ListTokens(
		ctx,
		"tokens",
	)
	if err != nil {
		return
	}

	p = pairs.PairsFromTokens(tokens)

	return
}

func (pc *ParseCase) Clear(
	ctx context.Context,
) {
	pc.Parser.Clear()
}
