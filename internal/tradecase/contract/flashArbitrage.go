package contract

import (
	"context"
	c "context"
	"fmt"
	"math/big"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"

	b "github.com/ethereum/go-ethereum/accounts/abi/bind"
	cm "github.com/ethereum/go-ethereum/common"
	t "github.com/ethereum/go-ethereum/core/types"
)

type contractApi interface {
	AddBaseToken(
		*b.TransactOpts, cm.Address,
	) (*t.Transaction, error)

	BaseTokensContains(
		*b.CallOpts, cm.Address,
	) (bool, error)

	GetBaseTokens(
		*b.CallOpts,
	) ([]cm.Address, error)

	RemoveBaseToken(
		*b.TransactOpts, cm.Address,
	) (*t.Transaction, error)

	GetProfit(
		*b.CallOpts, cm.Address, cm.Address,
	) (struct {
		Profit    *big.Int
		BaseToken cm.Address
	}, error)

	Withdraw(
		*b.TransactOpts,
	) (*t.Transaction, error)

	FlashArbitrage(
		*b.TransactOpts, cm.Address, cm.Address,
	) (*t.Transaction, error)
}

type FlashArbContract struct {
	Address    cm.Address
	api        contractApi
	tradePairs []TradePair
}

func NewContract(
	address cm.Address,
	api contractApi,
	pairs []TradePair,
) (
	contract *FlashArbContract,
) {
	contract = &FlashArbContract{
		Address:    address,
		api:        api,
		tradePairs: pairs,
	}

	return
}

func (fc *FlashArbContract) AddBaseToken(
	ctx c.Context,
	token Token,
) (
	out interface{},
	err error,
) {
	out, err = fc.api.AddBaseToken(
		&b.TransactOpts{Context: ctx},
		cm.HexToAddress(token.Address),
	)

	return
}

func (fc *FlashArbContract) BaseTokensContains(
	ctx c.Context,
	token Token,
) (
	ok bool,
	err error,
) {
	ok, err = fc.api.BaseTokensContains(
		&b.CallOpts{Context: ctx},
		cm.HexToAddress(token.Address),
	)

	return
}

func (fc *FlashArbContract) RemoveBaseToken(
	ctx c.Context,
	token Token,
) (
	out interface{},
	err error,
) {
	out, err = fc.api.RemoveBaseToken(
		&b.TransactOpts{
			Context: ctx,
		},
		cm.HexToAddress(
			token.Address,
		),
	)

	return
}

func (fc *FlashArbContract) GetBaseTokens(
	ctx c.Context,
) (
	baseAddr []string,
	err error,
) {
	out, err := fc.api.GetBaseTokens(
		&b.CallOpts{Context: ctx},
	)
	if err != nil {
		return
	}

	for _, addr := range out {
		baseAddr = append(baseAddr, addr.Hex())
	}

	return
}

func (fc *FlashArbContract) GetProfit(
	ctx context.Context,
	pair TradePair,
) (
	profit int,
	err error,
) {
	p, err := fc.api.GetProfit(
		&b.CallOpts{Context: ctx},
		cm.HexToAddress(pair.Pool0.Address),
		cm.HexToAddress(pair.Pool1.Address),
	)
	if err != nil {
		return
	}

	profit = int(p.Profit.Int64())

	return
}

func (fc *FlashArbContract) AddPair(
	ctx context.Context,
	pair TradePair,
) (
	err error,
) {
	index, ok := fc.containPair(
		pair.Pool0.Address,
		pair.Pool1.Address,
	)
	if ok {
		err = fmt.Errorf(
			"already added with index %v",
			index,
		)
		return
	}
	fc.tradePairs = append(fc.tradePairs, pair)
	fmt.Print(fc.tradePairs)
	return
}

func (fc *FlashArbContract) RemovePair(
	ctx context.Context,
	pool0, pool1 string,
) (
	err error,
) {
	index, ok := fc.containPair(pool0, pool1)
	if !ok {
		err = fmt.Errorf(
			"pair not found",
		)
		return
	}
	fc.tradePairs = append(
		fc.tradePairs[:index],
		fc.tradePairs[index+1:]...,
	)

	return
}

func (fc *FlashArbContract) GetPair(
	ctx context.Context,
	pool0, pool1 string,
) (
	pair TradePair,
	err error,
) {

	index, ok := fc.containPair(pool0, pool1)
	if !ok {
		err = fmt.Errorf(
			"pair not found",
		)
		return
	}
	pair = fc.tradePairs[index]

	return
}

func (fc *FlashArbContract) ListPairs(
	ctx context.Context,
) (
	vals []TradePair,
) {
	for _, v := range fc.tradePairs {
		vals = append(vals, v)
	}
	fmt.Println(vals)

	return
}

func (fc *FlashArbContract) ClearPairs(
	ctx context.Context,
) {
	new := make([]TradePair, 0)
	fc.tradePairs = new

	return
}

func (fc *FlashArbContract) SetPairs(
	ctx context.Context,
	pairs []TradePair,
) (
	err error,
) {
	for _, pair := range pairs {
		_, ok := fc.containPair(
			pair.Pool0.Address,
			pair.Pool1.Address,
		)
		if ok {
			continue
		}
		err = fc.AddPair(ctx, pair)

		return
	}

	return
}

func (fc *FlashArbContract) GetPairs(
	ctx c.Context,
	protocol SwapProtocol,
	tokens TokenPair,
) (
	pairs []TradePair,
	ok bool,
) {
	for _, pair := range fc.tradePairs {
		if checkPairTokens(
			pair,
			tokens,
		) && checkPairProtocol(
			pair,
			protocol,
		) {
			pairs = append(
				pairs,
				pair,
			)
			ok = true
		}
	}

	return
}

func (fc *FlashArbContract) containPair(
	pool0, pool1 string,
) (
	index int,
	ok bool,
) {
	ok = false

	for n, pair := range fc.tradePairs {
		if (pair.Pool0.Address == pool0 &&
			pair.Pool1.Address == pool1) ||
			(pair.Pool0.Address == pool1 &&
				pair.Pool1.Address == pool0) {
			index = n
			ok = true

			return
		}
	}

	return
}

func checkPairProtocol(
	pair TradePair,
	protocol SwapProtocol,
) (
	ok bool,
) {
	ok = false

	if pair.Pool0.Protocol == protocol &&
		pair.Pool1.Protocol == protocol {
		ok = true
	}

	return
}

func checkPairTokens(
	pair TradePair,
	tokens TokenPair,
) (
	ok bool,
) {
	ok = false

	if pair.Pool0.Pair == tokens &&
		pair.Pool1.Pair == tokens {
		ok = true
	}

	return
}
