package tradecase

import (
	c "context"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/contract"
)

type Repository interface {
	PoolRepo

	TokenRepo

	GetStorage() Storage
}

type TokenRepo interface {
	StoreTokens(
		c.Context, string, []entities.Token,
	) error

	AddToken(
		c.Context, string, entities.Token,
	) error

	ListTokens(
		c.Context, string,
	) ([]entities.Token, error)

	GetTokenByAddress(
		c.Context, string, string,
	) (entities.Token, error)

	RemoveToken(
		c.Context, string, entities.Token,
	) error

	RemoveTokens(
		c.Context, string, []entities.Token,
	) ([]entities.Token, error)
}

type PoolRepo interface {
	ListPools(
		c.Context, string,
	) ([]entities.Pool, error)

	StorePools(
		c.Context, string, []entities.Pool,
	) error

	AddPool(
		c.Context, entities.Pool, string,
	) error

	GetByTokens(
		c.Context, string, entities.TokenPair,
	) ([]entities.Pool, error)

	RemovePool(
		c.Context, string, entities.Pool,
	) error

	RemovePools(
		c.Context, string, []entities.Pool,
	) ([]entities.Pool, error)
}

type Storage interface {
	Store(
		c.Context, string, interface{},
	) error

	Read(
		c.Context, string, interface{},
	) error

	Remove(
		c.Context, string, interface{},
	) error

	Clear(
		c.Context, string,
	) error

	ClearAll(
		c.Context,
	) error
}

type TradeProvider interface {
	ClientManager

	ProviderStore
}

type ClientManager interface {
	GetClient(c.Context) Client
}

type Client interface {
	Setup(c.Context, interface{}) error

	ClientGet() interface{}

	UseWallet(interface{})

	GetBallance(c.Context) (int, error)

	GetChainID(c.Context) interface{}

	DialContract(string) (interface{}, error)

	GetNextTransaction(c.Context) (interface{}, error)

	UpdateChainID(c.Context) error

	Transact(c.Context, interface{}) (interface{}, error)
}

type ProviderStore interface {
	AddToken(
		c.Context, entities.Token,
	) error

	GetToken(
		c.Context, string,
	) (entities.Token, error)

	RemoveToken(
		c.Context, entities.Token,
	) error

	SetTokens(
		c.Context, []entities.Token,
	) error

	ListTokens(
		c.Context,
	) []entities.Token

	Clear()
}

type SmartContract interface {
	ContractPairs
	ContractAPI
	Trade() contract.Trade
}

type ContractPairs interface {
	AddPair(
		c.Context, entities.TradePair,
	) error

	RemovePair(
		c.Context, string, string,
	) error

	GetPair(
		c.Context, string, string,
	) (entities.TradePair, error)

	GetPairs(
		c.Context, entities.SwapProtocol, entities.TokenPair,
	) ([]entities.TradePair, bool)

	SetPairs(
		c.Context, []entities.TradePair,
	) error

	ListPairs(c.Context) []entities.TradePair

	ClearPairs(c.Context)
}

type ContractAPI interface {
	contract.API
}

type Parser interface {
	Parse([]entities.TokenPair) error
	ParseStore
	ParseProtocol
}

type ParseStore interface {
	AddPool(
		entities.Pool,
	) error

	RemovePool(
		entities.Pool,
	) error

	GetPairPools(
		entities.TokenPair,
	) ([]entities.Pool, error)

	AddPools(
		[]entities.Pool,
	)

	ListPools() []entities.Pool

	Clear()
}

type ParseProtocol interface {
	SetProtocol(entities.SwapProtocol)
	GetProtocol() entities.SwapProtocol
	GetPoolAddress(entities.TokenPair) (string, error)
}
