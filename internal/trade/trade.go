package trade

import (
	"context"
	"fmt"

	eth "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
)

func (tc *TradeCase) AddBaseToken(
	ctx context.Context,
	address string,
) (
	tx interface{},
	err error,
) {
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
	fmt.Println("sending tx")

	auth := tc.Provider.GetClient(ctx).(*eth.Client)

	b, err := auth.GetNextTransaction(ctx)
	if err != nil {

		return
	}

	t, err := tc.Contract.Api().Transactor().AddBaseToken(
		b, eth.ToAddress(address),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

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
	ok, err := tc.Contract.Api().Caller().BaseTokensContains(
		eth.CallOpts(ctx),
		eth.ToAddress(address),
	)
	if err != nil {
		return
	}
	if !ok {
		err = fmt.Errorf("token %v not found", address)

		return
	}

	fmt.Println("sending tx")

	auth := tc.Provider.GetClient(ctx).(*eth.Client)

	b, err := auth.GetNextTransaction(ctx)
	if err != nil {
		fmt.Println(err)

		return
	}

	t, err := tc.Contract.Api().Transactor().RemoveBaseToken(
		b, eth.ToAddress(address),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	tx, isPending, err := auth.Client.TransactionByHash(ctx, t.Hash())
	if err != nil {
		return
	}

	fmt.Println(isPending)

	return
}

func (tc *TradeCase) Withdraw(
	ctx context.Context,
) (
	tx interface{},
	err error,
) {
	fmt.Println("sending tx")
	auth := tc.Provider.GetClient(ctx).(*eth.Client)

	// bal, err := auth.Client.BalanceAt(
	// 	ctx, eth.ToAddress(tc.Contract.Address()), nil)
	// if err != nil {
	// 	return
	// }
	b, err := auth.GetNextTransaction(ctx)
	if err != nil {
		fmt.Println(err)

		return
	}

	// b.Value = bal.Sub

	t, err := tc.Contract.Api().Transactor().Withdraw(b)
	if err != nil {
		fmt.Println(err)

		return
	}

	tx, isPending, err := auth.Client.TransactionByHash(ctx, t.Hash())
	if err != nil {
		return
	}

	fmt.Println(isPending)

	return
}

func (tc *TradeCase) GetProfit(ctx context.Context, pool0, pool1 string) (
	profit int,
	baseToken string,
	err error,
) {
	res, err := tc.Contract.Api().Caller().GetProfit(
		eth.CallOpts(ctx),
		eth.ToAddress(pool0),
		eth.ToAddress(pool1),
	)
	if err != nil {
		return
	}

	profit = eth.FromBigInt(res.Profit)
	baseToken = eth.FromAddress(res.BaseToken)

	return
}

func (tc *TradeCase) Arbitrage(ctx context.Context, pool0, pool1 string) (
	tx interface{},
	err error,
) {
	fmt.Println("sending tx")

	auth := tc.Provider.GetClient(ctx).(*eth.Client)

	b, err := auth.GetNextTransaction(ctx)
	if err != nil {
		return
	}

	t, err := tc.Contract.Api().Transactor().FlashArbitrage(
		b, eth.ToAddress(pool0), eth.ToAddress(pool1),
	)
	if err != nil {
		return
	}

	tx, isPending, err := auth.Client.TransactionByHash(ctx, t.Hash())
	if err != nil {
		return
	}

	fmt.Println(isPending)

	return
}

func (tc *TradeCase) ReplaceTxWithAddBaseToken(
	ctx context.Context,
	token, hash string,
) (
	tx interface{},
	err error,
) {
	auth := tc.Provider.GetClient(ctx).(*eth.Client)
	b, err := auth.ReplaceTx(ctx, hash)
	if err != nil {
		return
	}

	tx, err = tc.Contract.Api().Transactor().AddBaseToken(b, eth.ToAddress(token))

	return
}

func (tc *TradeCase) ReplaceTxWithRmBaseToken(
	ctx context.Context,
	token, hash string,
) (
	tx interface{},
	err error,
) {
	auth := tc.Provider.GetClient(ctx).(*eth.Client)
	b, err := auth.ReplaceTx(ctx, hash)
	if err != nil {
		return
	}

	tx, err = tc.Contract.Api().Transactor().RemoveBaseToken(b, eth.ToAddress(token))

	return
}
