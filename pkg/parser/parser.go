package parser

import (
	"fmt"
	"log"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type Parser struct {
	Pools []entities.TradePool
}

func New() *Parser {
	return &Parser{make([]entities.TradePool, 0)}
}

func (p *Parser) AddPool(pool entities.TradePool) error {
	if _, ok := p.containPool(pool); ok {
		return fmt.Errorf("pool already added")
	}
	p.Pools = append(p.Pools, pool)

	return nil
}

func (p *Parser) RemovePool(pool entities.TradePool) error {
	index, ok := p.containPool(pool)
	if !ok {
		return fmt.Errorf("pool not found")
	}
	p.Pools = append(p.Pools[:index], p.Pools[index+1:]...)

	return nil
}

func (p *Parser) Clear() {
	p.Pools = make([]entities.TradePool, 0)
}

func (p *Parser) GetPairPools(
	pair entities.TokenPair,
) (
	pools []entities.TradePool,
	err error,
) {
	index, ok := p.containPair(pair)
	if !ok {
		pools = nil
		err = fmt.Errorf(
			"pools with %v pair not found",
			pair,
		)

		return
	}

	pools = append(pools, p.Pools[index])

	return
}

func (p *Parser) AddPools(
	pools []entities.TradePool,
) {
	for _, pool := range pools {
		err := p.AddPool(pool)
		if err != nil {
			log.Println(err)
		}
	}
}

func (p *Parser) ListPools() (
	listPools []entities.TradePool,
) {
	listPools = p.Pools

	return
}

func (p *Parser) containPool(pool entities.TradePool) (
	index int,
	ok bool,
) {
	for n, pl := range p.Pools {
		if pl == pool {
			index = n
			ok = true

			return
		}
	}
	ok = false

	return
}

func (p *Parser) containPair(
	pair entities.TokenPair,
) (
	index int,
	ok bool,
) {
	for n, p := range p.Pools {
		if p.Pair == pair {
			index = n
			ok = true

			return
		}
	}
	ok = false

	return
}
