package parser

import (
	"fmt"

	uni "github.com/ackermanx/ethclient/uniswap"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type ParserUniV2 struct {
	protocol SwapProtocol
	pools    map[TokenPair][]TradePool
}

func New(protocol SwapProtocol, tokens []TokenPair) ParserUniV2 {
	poolMap := make(map[TokenPair][]TradePool, 0)

	for _, token := range tokens {
		poolMap[token] = make([]TradePool, 0)
	}

	return ParserUniV2{
		protocol,
		poolMap,
	}
}

func (pu *ParserUniV2) Parse() error {
	for pair, pool := range pu.pools {
		addr, err := getPoolAddr(pair)
		if err != nil {
			return err
		}

		pool = append(pool, TradePool{
			SwapProtocol: pu.protocol,
			Address:      addr,
			Pair:         pair,
		})
	}

	return nil
}

func (pu *ParserUniV2) AddPair(pair TokenPair) error {
	if pair, ok := pu.pools[pair]; ok {
		return fmt.Errorf("pair %v already added", pair)
	}
	pu.pools[pair] = make([]TradePool, 0)

	return nil
}

func (pu *ParserUniV2) RemovePair(pair TokenPair) error {
	if pair, ok := pu.pools[pair]; !ok {
		return fmt.Errorf("pair %v not found", pair)
	}
	delete(pu.pools, pair)

	return nil
}

func (pu *ParserUniV2) GetPairPools(pair TokenPair) ([]TradePool, error) {
	pools, ok := pu.pools[pair]
	if !ok {
		return nil, fmt.Errorf("pair %v not found", pair)
	}

	return pools, nil
}

func (pu *ParserUniV2) AddPairPools(pair TokenPair, pools []TradePool) {
	pools, ok := pu.pools[pair]
	if !ok {
		pu.pools[pair] = pools
		return
	}
	pu.pools[pair] = append(pu.pools[pair], pools...)

	return
}

func (pu *ParserUniV2) ListPools() map[TokenPair][]TradePool {
	return pu.pools
}

func getPoolAddr(pair TokenPair) (string, error) {
	pAddr, err := uni.CalculatePoolAddressV2(
		pair.Token0.Address,
		pair.Token1.Address,
	)
	if err != nil {
		return "", err
	}

	return pAddr.String(), nil
}
