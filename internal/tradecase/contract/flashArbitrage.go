package contract

import (
	c "context"
	"fmt"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/trade"

	cm "github.com/ethereum/go-ethereum/common"
)

type FlashArbContract struct {
	Address    cm.Address
	api        API
	trade      Trade
	tradePairs []entities.TradePair
}

func NewContract(
	address cm.Address,
	api flashArb,
	pairs []entities.TradePair,
) (
	contract *FlashArbContract,
) {
	a := &Api{api}

	contract = &FlashArbContract{
		Address:    address,
		api:        a,
		trade:      TradeApi(a),
		tradePairs: pairs,
	}

	return
}

func (fc *FlashArbContract) API() (
	out flashArb,
) {
	out = fc.api.API()

	return
}

func (fc *FlashArbContract) Trade() (
	out Trade,
) {
	out = fc.trade

	return
}

func (fc *FlashArbContract) AddPair(
	ctx c.Context,
	pair entities.TradePair,
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
	ctx c.Context,
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
	ctx c.Context,
	pool0, pool1 string,
) (
	pair entities.TradePair,
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
	ctx c.Context,
) (
	vals []entities.TradePair,
) {
	vals = append(vals, fc.tradePairs...)

	fmt.Println(vals)

	return
}

func (fc *FlashArbContract) ClearPairs(
	ctx c.Context,
) {
	new := make([]entities.TradePair, 0)
	fc.tradePairs = new

}

func (fc *FlashArbContract) SetPairs(
	ctx c.Context,
	pairs []entities.TradePair,
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
	protocol entities.SwapProtocol,
	tokens entities.TokenPair,
) (
	pairs []entities.TradePair,
	ok bool,
) {
	for _, pair := range fc.tradePairs {
		if trade.CheckPairTokens(
			pair,
			tokens,
		) && trade.CheckPairProtocol(
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
