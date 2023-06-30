package tradecase

import (
	c "context"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type TradeRepo interface {
	Store(
		c.Context, string, []byte,
	) error

	Read(
		c.Context, string,
	) ([]byte, error)

	GetByTokens(
		c.Context, string, TokenPair,
	) ([]TradePool, error)

	ListPools(
		c.Context, string,
	) ([]TradePool, error)

	StorePools(
		c.Context, string, []TradePool,
	) error

	AddPool(
		c.Context, TradePool, string,
	) error

	StoreTokens(
		c.Context, string, []Token,
	) error

	AddToken(
		c.Context, string, Token,
	) error

	ListTokens(
		c.Context, string,
	) ([]Token, error)

	GetTokenByAddress(
		c.Context, string, string,
	) (Token, error)

	Clear(
		c.Context, string,
	) error
}

type TradeProvider interface {
	AddToken(
		c.Context, Token,
	) error

	GetToken(
		c.Context, string,
	) (Token, error)

	RemoveToken(
		c.Context, Token,
	) error

	SetTokens(
		c.Context, []Token,
	) error

	ListTokens(
		c.Context,
	) []Token

	Clear()
}

type SmartContract interface {
	AddBaseToken(
		c.Context, Token,
	) (interface{}, error)

	BaseTokensContains(
		c.Context, Token,
	) (bool, error)

	RemoveBaseToken(
		c.Context, Token,
	) (interface{}, error)

	GetBaseTokens(
		c.Context,
	) ([]string, error)

	GetProfit(
		c.Context, TradePair,
	) (int, error)

	AddPair(
		c.Context, TradePair,
	) error

	RemovePair(
		c.Context, string, string,
	) error

	GetPair(
		c.Context, string, string,
	) (TradePair, error)

	GetPairs(
		c.Context, SwapProtocol, TokenPair,
	) ([]TradePair, bool)

	SetPairs(
		c.Context, []TradePair,
	) error

	ListPairs(c.Context) []TradePair

	ClearPairs(c.Context)
}

type Parser interface {
	Parse([]TokenPair) error

	AddPool(
		TradePool,
	) error

	RemovePool(
		TradePool,
	) error

	GetPairPools(
		TokenPair,
	) ([]TradePool, error)

	AddPools(
		[]TradePool,
	)

	ListPools() []TradePool
}
