package contract

import (
	"context"
	"fmt"
	"math/big"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	b "github.com/ethereum/go-ethereum/accounts/abi/bind"
	cm "github.com/ethereum/go-ethereum/common"
	t "github.com/ethereum/go-ethereum/core/types"
)

type Ballance interface {
	Withdraw(
		*b.TransactOpts,
	) (*t.Transaction, error)
}

type BaseTokens interface {
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
}

type Arbitrage interface {
	GetProfit(
		*b.CallOpts, cm.Address, cm.Address,
	) (struct {
		Profit    *big.Int
		BaseToken cm.Address
	}, error)

	FlashArbitrage(
		*b.TransactOpts, cm.Address, cm.Address,
	) (*t.Transaction, error)
}

type flashArb interface {
	Ballance
	BaseTokens
	Arbitrage
}

type API interface {
	API() flashArb
}

type Api struct {
	flashArb
}

func (api *Api) API() flashArb {
	return api.flashArb
}

func TradeApi(api API) Trade {
	return &tradeApi{api: api}
}

type Trade interface {
	GetProfitable(
		context.Context, []entities.TradePair) (
		[]entities.TradePair,
		error,
	)

	DoArbitrage(
		context.Context, entities.TradePair) (
		interface{}, error,
	)

	Withdraw(context.Context) (interface{}, error)
}

type tradeApi struct {
	api API
}

func (tr *tradeApi) GetProfitable(
	ctx context.Context,
	pairs []entities.TradePair,
) (
	out []entities.TradePair,
	err error,
) {
	for _, pair := range pairs {

		p, _err := tr.api.API().GetProfit(
			&b.CallOpts{},
			cm.HexToAddress(pair.Pool0.Address),
			cm.HexToAddress(pair.Pool1.Address),
		)
		if err != nil {
			err = _err
			return
		}

		if p.Profit.Int64() > 0 {
			out = append(out, pair)
		}
	}

	return
}

func (tr *tradeApi) DoArbitrage(
	ctx context.Context,
	pair entities.TradePair,
) (
	tx interface{},
	err error,
) {
	tx, err = tr.api.API().FlashArbitrage(
		&b.TransactOpts{},
		cm.HexToAddress(pair.Pool0.Address),
		cm.HexToAddress(pair.Pool1.Address),
	)

	return
}

func (tr *tradeApi) Withdraw(
	ctx context.Context,
) (
	tx interface{},
	err error,
) {
	fmt.Println("sending tx")
	tx, err = tr.api.API().Withdraw(&b.TransactOpts{})
	if err != nil {
		fmt.Println(tx)
	}
	fmt.Println("tx sended")

	return
}
