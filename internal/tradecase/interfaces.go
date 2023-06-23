package tradecase

import (
	c "context"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type TradeRepo interface {
	Store(c.Context, []byte) error
	Read(c.Context) ([]byte, error)
	GetByTokens(c.Context, TokenPair) ([]TradePool, error)
	StorePool(c.Context, TradePool) error
	ListAllPools(c.Context) ([]TradePool, error)
}

type TradeProvider interface {
	AddPair(c.Context, TradePair) error
	FindPair(c.Context, TradePair) (int, bool)
	RemovePair(c.Context, TradePair) error
	GetPairs(c.Context, SwapProtocol, TokenPair) ([]TradePair, bool)
	ListAllPairs(c.Context) []TradePair
}

type SmartContract interface {
	AddBaseToken(c.Context, Token) error
	BaseTokensContains(c.Context, Token) (int, bool)
	RemoveBaseToken(c.Context, Token) error
	GetBaseToken(c.Context, string) (Token, error)
	GetBaseTokens(c.Context) []Token
}
