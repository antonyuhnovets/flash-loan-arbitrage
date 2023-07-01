package parser

import (
	"fmt"

	sushi "github.com/ebadiere/go-defi/sushiswap"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type ParserSushiV2 struct {
	Protocol SwapProtocol `json:"protocol"`
	Pools    []TradePool  `json:"pools"`
}

func NewSushiV2(
	protocol SwapProtocol,
	tokens []TokenPair,
) (
	parser ParserSushiV2,
) {
	poolList := make(
		[]TradePool,
		0,
	)

	parser = ParserSushiV2{
		protocol,
		poolList,
	}

	return
}

func (ps *ParserSushiV2) Parse(pairs []TokenPair) (
	err error,
) {
	for _, pair := range pairs {
		addr, err := getSushiPoolAddr(pair)
		if err != nil {
			return err
		}

		ps.AddPool(
			TradePool{
				Protocol: ps.Protocol,
				Address:  addr,
				Pair:     pair,
			},
		)
	}

	return
}

func (ps *ParserSushiV2) AddPool(
	pool TradePool,
) (
	err error,
) {
	if _, ok := ps.containPool(pool); ok {
		err = fmt.Errorf(
			"pool %v already added",
			pool,
		)
		return
	}
	ps.Pools = append(ps.Pools, pool)

	return
}

func (ps *ParserSushiV2) RemovePool(
	pool TradePool,
) (
	err error,
) {
	index, ok := ps.containPool(pool)
	if !ok {
		err = fmt.Errorf(
			"pool %v not found",
			pool,
		)
		return
	}

	ps.Pools = append(
		ps.Pools[:index],
		ps.Pools[index+1:]...,
	)

	return
}

func (ps *ParserSushiV2) GetPairPools(
	pair TokenPair,
) (
	pools []TradePool,
	err error,
) {
	index, ok := ps.containPair(pair)
	if !ok {
		pools = nil
		err = fmt.Errorf(
			"pools with %v pair not found",
			pair,
		)

		return
	}

	pools = append(pools, ps.Pools[index])

	return
}

func (ps *ParserSushiV2) AddPools(
	pools []TradePool,
) {
	for _, pool := range pools {
		ps.AddPool(pool)
	}

	return
}

func (ps *ParserSushiV2) Clear() {
	ps.Pools = make([]TradePool, 0)

	return
}

func (ps *ParserSushiV2) ListPools() (
	listPools []TradePool,
) {
	listPools = ps.Pools

	return
}

func (ps *ParserSushiV2) containPool(
	pool TradePool,
) (
	index int,
	ok bool,
) {
	for n, p := range ps.Pools {
		if p == pool {
			index = n
			ok = true

			return
		}
	}
	ok = false

	return
}

func (ps *ParserSushiV2) containPair(
	pair TokenPair,
) (
	index int,
	ok bool,
) {
	for n, p := range ps.Pools {
		if p.Pair == pair {
			index = n
			ok = true

			return
		}
	}
	ok = false

	return
}

func getSushiPoolAddr(
	pair TokenPair,
) (
	address string,
	err error,
) {
	pAddr := sushi.GeneratePairAddress(
		common.HexToAddress(pair.Token0.Address),
		common.HexToAddress(pair.Token1.Address),
	)

	address = pAddr.Hex()

	return
}
