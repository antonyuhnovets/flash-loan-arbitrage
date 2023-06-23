package contract

import (
	c "context"
	"log"
	"math/big"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"

	b "github.com/ethereum/go-ethereum/accounts/abi/bind"
	cm "github.com/ethereum/go-ethereum/common"
	t "github.com/ethereum/go-ethereum/core/types"
)

type contractApi interface {
	AddBaseToken(opts *b.TransactOpts, token cm.Address) error
	BaseTokensContains(*b.CallOpts, cm.Address) (bool, error)
	GetBaseTokens(*b.CallOpts) ([]cm.Address, error)
	RemoveBaseToken(opts *b.TransactOpts, token cm.Address) error
	GetProfit(*b.CallOpts, cm.Address, cm.Address) (struct {
		Profit    *big.Int
		BaseToken cm.Address
	}, error)
	Withdraw(*b.CallOpts) (*t.Transaction, error)
	FlashArbitrage(*b.CallOpts, cm.Address, cm.Address) (*t.Transaction, error)
}

type FlashArbContract struct {
	Address cm.Address
	api     contractApi
	tokens  map[string]Token
}

func NewContract(address cm.Address, api contractApi, tokens []Token) *FlashArbContract {
	tokenMap := map[string]Token{}
	for _, token := range tokens {
		tokenMap[token.Address] = token
	}
	return &FlashArbContract{
		Address: address,
		api:     api,
		tokens:  tokenMap,
	}
}

func (fac *FlashArbContract) AddBaseToken(ctx c.Context, token Token) error {
	fac.add(token)

	return fac.api.AddBaseToken(
		&b.TransactOpts{Context: ctx},
		cm.HexToAddress(token.Address),
	)
}

func (fac *FlashArbContract) BaseTokensContains(ctx c.Context, token Token) (bool, error) {
	ok, err := fac.api.BaseTokensContains(
		&b.CallOpts{Context: ctx},
		cm.HexToAddress(token.Address),
	)
	if ok {
		fac.add(token)
	}

	return ok, err
}

func (fac *FlashArbContract) RemoveBaseToken(ctx c.Context, token Token) error {
	fac.remove(token)

	return fac.api.RemoveBaseToken(
		&b.TransactOpts{Context: ctx},
		cm.HexToAddress(token.Address),
	)
}

func (fac *FlashArbContract) GetBaseTokens(ctx c.Context) ([]Token, error) {
	out, err := fac.api.GetBaseTokens(
		&b.CallOpts{Context: ctx},
	)
	if err != nil {
		return nil, err
	}

	outTokens := make([]Token, 0)
	for _, addr := range out {
		t, ok := fac.tokens[addr.String()]
		if !ok {
			log.Printf("unknown %s token", addr.String())
			fac.tokens[addr.String()] = Token{}
		}
		outTokens = append(outTokens, t)
	}

	return outTokens, nil
}

func (fac *FlashArbContract) add(token Token) {
	if t, ok := fac.tokens[token.Address]; (!ok || t == Token{}) {
		fac.tokens[token.Address] = token
	}
}

func (fac *FlashArbContract) remove(token Token) {
	if _, ok := fac.tokens[token.Address]; ok {
		delete(fac.tokens, token.Address)
	}
}
