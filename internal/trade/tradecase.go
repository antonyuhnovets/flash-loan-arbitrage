package trade

import (
	"context"
	"fmt"
	"log"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	eth "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
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

func (tc *TradeCase) GetProfitable(from []entities.TradePair) (
	out []entities.TradePair,
	ok bool,
	err error,
) {
	ok = false
	for _, pair := range from {
		res, _err := tc.Contract.Api().Caller().GetProfit(
			eth.CallOpts(),
			eth.ToAddress(pair.Pool0.Address),
			eth.ToAddress(pair.Pool1.Address),
		)
		if err != nil {
			err = _err
			return
		}
		if int(res.Profit.Int64()) > 0 {
			out = append(out, pair)
			ok = true
		}
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
	// auth := tc.Provider.GetClient(ctx)

	b, err := auth.GetNextTransaction(c)
	if err != nil {
		log.Println(err)

		return
	}

	t, err := tc.Contract.Api().Transactor().Withdraw(b)
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
	// c, cancel := context.WithCancel(ctx)
	// defer cancel()

	ok, err := tc.Contract.Api().Caller().BaseTokensContains(
		eth.CallOpts(ctx),
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

	auth := tc.Provider.GetClient(ctx).(*eth.Client)
	// auth := tc.Provider.GetClient(c)

	b, err := auth.GetNextTransaction(ctx)
	if err != nil {

		return
	}

	t, err := tc.Contract.Api().Transactor().AddBaseToken(
		b, eth.ToAddress(address),
	)
	if err != nil {
		log.Println(err)

		return
	}

	// r, err := bind.WaitMined(c, auth.Client, t)

	// fmt.Println(r)

	tx, isPending, err := auth.Client.TransactionByHash(ctx, t.Hash())
	if err != nil {
		return
	}

	fmt.Println(isPending)

	return
}

func (tc *TradeCase) RmBaseToken(
	ctx context.Context,
	address string,
) (
	tx interface{},
	err error,
) {
	// c, cancel := context.WithCancel(ctx)
	// defer cancel()

	ok, err := tc.Contract.Api().Caller().BaseTokensContains(
		eth.CallOpts(ctx),
		eth.ToAddress(address),
	)
	if err != nil {
		return
	}
	if !ok {
		tx = nil
		err = fmt.Errorf("token %v not found", address)

		return
	}

	log.Println("sending tx")

	auth := tc.Provider.GetClient(ctx).(*eth.Client)
	// auth := tc.Provider.GetClient(c)

	b, _err := auth.GetNextTransaction(ctx)
	if err != nil {
		log.Println(err)
		err = _err

		return
	}

	t, _err := tc.Contract.Api().Transactor().RemoveBaseToken(
		b, eth.ToAddress(address),
	)
	if _err != nil {
		err = _err

		return
	}

	tx, isPending, _err := auth.Client.TransactionByHash(ctx, t.Hash())
	if err != nil {
		err = _err

		return
	}

	fmt.Println(isPending)

	return
}
