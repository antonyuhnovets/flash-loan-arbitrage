package parser

import (
	"fmt"

	uni "github.com/ackermanx/ethclient/uniswap"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type ParserUniV2 struct {
	Protocol SwapProtocol `json:"protocol"`
	Pools    []TradePool  `json:"pools"`
}

func NewUniV2(
	protocol SwapProtocol,
	tokens []TokenPair,
) (
	parser ParserUniV2,
) {
	poolList := make(
		[]TradePool,
		0,
	)

	parser = ParserUniV2{
		protocol,
		poolList,
	}

	return
}

func (pu *ParserUniV2) Parse(pairs []TokenPair) (
	err error,
) {
	for _, pair := range pairs {
		addr, err := getUniPoolAddr(pair)
		if err != nil {
			return err
		}

		pu.AddPool(
			TradePool{
				Protocol: pu.Protocol,
				Address:  addr,
				Pair:     pair,
			},
		)
	}

	return
}

func (pu *ParserUniV2) AddPool(
	pool TradePool,
) (
	err error,
) {
	if _, ok := pu.containPool(pool); ok {
		err = fmt.Errorf(
			"pool %v already added",
			pool,
		)
		return
	}
	pu.Pools = append(pu.Pools, pool)

	return
}

func (pu *ParserUniV2) RemovePool(
	pool TradePool,
) (
	err error,
) {
	index, ok := pu.containPool(pool)
	if !ok {
		err = fmt.Errorf(
			"pool %v not found",
			pool,
		)
		return
	}

	pu.Pools = append(
		pu.Pools[:index],
		pu.Pools[index+1:]...,
	)

	return
}

func (ps *ParserUniV2) Clear() {
	ps.Pools = make([]TradePool, 0)

	return
}

func (pu *ParserUniV2) GetPairPools(
	pair TokenPair,
) (
	pools []TradePool,
	err error,
) {
	index, ok := pu.containPair(pair)
	if !ok {
		pools = nil
		err = fmt.Errorf(
			"pools with %v pair not found",
			pair,
		)

		return
	}

	pools = append(pools, pu.Pools[index])

	return
}

func (pu *ParserUniV2) AddPools(
	pools []TradePool,
) {
	for _, pool := range pools {
		pu.AddPool(pool)
	}

	return
}

func (pu *ParserUniV2) ListPools() (
	listPools []TradePool,
) {
	listPools = pu.Pools

	return
}

func (pu *ParserUniV2) containPool(
	pool TradePool,
) (
	index int,
	ok bool,
) {
	for n, p := range pu.Pools {
		if p == pool {
			index = n
			ok = true

			return
		}
	}
	ok = false

	return
}

func (pu *ParserUniV2) containPair(
	pair TokenPair,
) (
	index int,
	ok bool,
) {
	for n, p := range pu.Pools {
		if p.Pair == pair {
			index = n
			ok = true

			return
		}
	}
	ok = false

	return
}

func getUniPoolAddr(
	pair TokenPair,
) (
	address string,
	err error,
) {
	pAddr, err := uni.CalculatePoolAddressV2(
		pair.Token0.Address,
		pair.Token1.Address,
	)
	if err != nil {
		return
	}

	address = pAddr.Hex()

	return
}
