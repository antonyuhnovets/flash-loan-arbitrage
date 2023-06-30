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
	Client ethClient
	Tokens []Token `json:"tokens"`
}

func NewTradeProvider(
	ctx c.Context,
	client ethClient,
	tokens ...Token,
) (
	provider *TradeProvider,
) {
	provider = &TradeProvider{
		Client: client,
		Tokens: tokens,
	}

	return
}

func (tp *TradeProvider) AddToken(
	ctx c.Context,
	token Token,
) (
	err error,
) {
	index, ok := tp.containToken(token.Address)
	if ok {
		err = fmt.Errorf(
			"token %v already added with index %v",
			token,
			index,
		)
		return
	}

	tp.Tokens = append(
		tp.Tokens,
		token,
	)

	return
}

func (tp *TradeProvider) GetToken(
	ctx c.Context,
	address string,
) (
	token Token,
	err error,
) {
	index, ok := tp.containToken(address)
	if !ok {
		err = fmt.Errorf(
			"no token with address %s",
			address,
		)
		return
	}
	token = tp.Tokens[index]

	return
}

func (tp *TradeProvider) RemoveToken(
	ctx c.Context,
	token Token,
) (
	err error,
) {
	index, ok := tp.containToken(token.Address)
	if !ok {
		err = fmt.Errorf(
			"no token %v",
			token,
		)
		return
	}

	tp.Tokens = append(
		tp.Tokens[:index],
		tp.Tokens[index+1:]...,
	)

	return
}

func (tp *TradeProvider) SetTokens(
	ctx c.Context,
	tokens []Token,
) (
	err error,
) {
	for _, token := range tokens {
		err = tp.AddToken(
			ctx,
			token,
		)
		if err != nil {
			return
		}
	}

	return
}

func (tp *TradeProvider) ListTokens(
	ctx c.Context,
) (
	tokens []Token,
) {
	tokens = tp.Tokens

	return
}

func (tp *TradeProvider) Clear() {
	new := make(
		[]Token,
		0,
	)
	tp.Tokens = new

	return
}

func (tp *TradeProvider) containToken(
	address string,
) (
	index int,
	ok bool,
) {
	ok = false

	for n, t := range tp.Tokens {
		if t.Address == address {
			index = n
			ok = true

			return
		}
	}
	return
}
