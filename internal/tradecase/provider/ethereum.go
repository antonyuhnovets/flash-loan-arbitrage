package provider

import (
	c "context"
	"fmt"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"

	eth "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
)

type TradeProvider struct {
	Client  *eth.Client
	Wallets []*eth.Wallet
	Tokens  []Token `json:"tokens"`
}

func NewTradeProvider(
	ctx c.Context,
	url, pk string,
	tokens ...Token,
) (
	provider *TradeProvider,
	err error,
) {
	cl, err := eth.NewClient(url)
	if err != nil {
		return
	}

	wall := eth.NewWallet()
	err = wall.Setup(pk)
	if err != nil {
		return
	}

	err = cl.Setup(ctx, wall)
	if err != nil {
		return
	}

	provider = &TradeProvider{
		Client:  cl,
		Wallets: []*eth.Wallet{wall},
		Tokens:  tokens,
	}

	return
}

func (tp *TradeProvider) GetClient(ctx c.Context) (
	cl interface{},
) {
	cl = tp.Client

	return
}

func (tp *TradeProvider) Ballance(ctx c.Context) (
	ball int,
	err error,
) {
	ball, err = tp.Client.GetBallance(ctx)

	return
}

func (tp *TradeProvider) SwitchWallet(addr string) (
	err error,
) {
	for _, wall := range tp.Wallets {
		if wall.Address.String() == addr {
			tp.Client.UseWallet(wall)

			return
		}
	}
	err = fmt.Errorf("no wallet with address %s", addr)

	return
}

func (tp *TradeProvider) AddWallet(pk string) (
	err error,
) {
	wall := eth.NewWallet()
	err = wall.Setup(pk)
	if err != nil {
		return
	}

	tp.Wallets = append(tp.Wallets, wall)

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
		_, ok := tp.containToken(token.Address)
		if ok {
			continue
		}
		err = tp.AddToken(ctx, token)
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
