package tradecase

import (
	c "context"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type TradeRepo interface {
	Store(c.Context, []byte,
	) error

	Read(c.Context) (
		[]byte,
		error,
	)
	GetByTokens(c.Context, TokenPair) (
		[]TradePool,
		error,
	)
	ListPools(c.Context) (
		[]TradePool,
		error,
	)
	StorePool(c.Context, TradePool,
	) error

	GetTokenByAddress(c.Context, string) (
		Token,
		error,
	)

	Clear(c.Context) error
}

type TradeProvider interface {
	AddPair(c.Context, TradePair,
	) error

	FindPair(c.Context, TradePair) (
		int,
		bool,
	)
	RemovePair(c.Context, TradePair,
	) error

	GetPairs(c.Context, SwapProtocol, TokenPair) (
		[]TradePair,
		bool,
	)

	SetPairs(c.Context, []TradePair) error

	ListPairs(c.Context,
	) []TradePair

	Clear()
}

type SmartContract interface {
	AddBaseToken(c.Context, Token) (
		interface{},
		error,
	)
	BaseTokensContains(c.Context, Token) (
		bool,
		error,
	)
	RemoveBaseToken(c.Context, Token) (
		interface{},
		error,
	)
	GetBaseTokens(c.Context) (
		[]Token,
		error,
	)

	GetProfit(c.Context, TradePair) (
		int,
		error,
	)

	Add(Token)

	Remove(Token)

	Get(string) (Token, bool)

	Tokens() map[string]Token

	List() []Token

	Clear()
}
