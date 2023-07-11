package tradecase

import (
	"context"
	"fmt"
	"log"

	eth "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/trade"
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

	tradeMap, err := trade.GetTradeMap(
		pools,
	)
	if err != nil {
		return
	}

	tradePairs, err := trade.GetTradePairs(
		tradeMap,
	)
	if err != nil {
		return
	}

	prof, err := tc.Contract.Trade().GetProfitable(
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

func (tc *TradeCase) Withdraw(
	ctx context.Context,
) (
	tx interface{},
	err error,
) {
	log.Println("sending tx")
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	auth := tc.Provider.GetClient(ctx).(*eth.Client)

	b, err := auth.GetNextTransaction(c)
	if err != nil {
		log.Println(err)

		return
	}

	t, err := tc.Contract.API().Withdraw(b)
	if err != nil {
		log.Println(err)

		return
	}

	tx, err = auth.Transact(c, t)
	if err != nil {
		log.Println(err)

		return
	}

	return
}

func (tc *TradeCase) AddBaseToken(
	ctx context.Context,
	address string,
) (
	tx interface{},
	err error,
) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	ok, err := tc.Contract.API().BaseTokensContains(
		eth.CallOpts(c),
		eth.ToAddress(address),
	)
	if err != nil {
		return
	}
	if ok {
		err = fmt.Errorf("token %v already added", address)
		return
	}
	log.Println("sending tx")

	auth := tc.Provider.GetClient(c).(*eth.Client)

	b, err := auth.GetNextTransaction(c)
	if err != nil {

		return
	}

	t, err := tc.Contract.API().AddBaseToken(
		b, eth.ToAddress(address),
	)
	if err != nil {
		log.Println(err)

		return
	}

	tx, err = auth.Transact(c, t)
	if err != nil {
		log.Println(err)

		return
	}

	return
}

func (tc *TradeCase) RmBaseToken(
	ctx context.Context,
	address string,
) (
	tx interface{},
	err error,
) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	ok, err := tc.Contract.API().BaseTokensContains(
		eth.CallOpts(c),
		eth.ToAddress(address),
	)
	if err != nil {
		return
	}
	if !ok {
		err = fmt.Errorf("token %v not found", address)

		return
	}
	log.Println("sending tx")

	auth := tc.Provider.GetClient(c).(*eth.Client)

	b, err := auth.GetNextTransaction(c)
	if err != nil {
		log.Println(err)

		return
	}

	t, err := tc.Contract.API().RemoveBaseToken(
		b, eth.ToAddress(address),
	)
	if err != nil {
		log.Println(err)

		return
	}

	tx, err = auth.Transact(c, t)
	if err != nil {
		log.Println(err)

		return
	}

	return
}
