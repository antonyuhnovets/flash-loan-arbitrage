package pairs

import (
	"fmt"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

func GetTradePairs(
	tradeMap map[entities.TokenPair]map[string]entities.Pool,
) (
	tradePairs []entities.TradePair,
	err error,
) {
	for _, poolMap := range tradeMap {
		n := len(poolMap)

		switch {
		case n < 2:
			continue
		case n == 2:
			pair, _err := MakeTradePair(
				poolMap,
			)
			if err != nil {
				err = _err

				return
			}
			tradePairs = append(
				tradePairs,
				pair,
			)
		case n > 2:
			tradePairs = DeleteDublicates(
				SplitPoolsOnPairs(poolMap),
			)
		}
	}

	return
}

func GetTradeMap(
	pools []entities.Pool,
) (
	trade map[entities.TokenPair]map[string]entities.Pool,
	err error,
) {
	trade = make(map[entities.TokenPair]map[string]entities.Pool)

	for _, pool := range pools {
		_, ok := trade[pool.Pair]
		if !ok {
			trade[pool.Pair] = make(
				map[string]entities.Pool,
				0,
			)
		}
		_, ok = trade[pool.Pair][pool.Address]
		if ok {
			continue
		}
		trade[pool.Pair][pool.Address] = pool
	}

	return
}

func SplitPoolsOnPairs(
	pools map[string]entities.Pool,
) (
	pairs []entities.TradePair,
) {
	addrList := make([]string, 0)

	for addr, pool := range pools {
		addrList = append(
			addrList,
			addr,
		)
		pairs = append(
			pairs,
			entities.TradePair{Pool0: pool},
		)
	}

	for n, pair := range pairs {
		for i := len(addrList); i >= 0; i-- {
			if i == n {
				continue
			}
			pair.Pool1 = pools[addrList[i]]
		}
	}

	return
}

func DeleteDublicates(
	pairs []entities.TradePair,
) (
	out []entities.TradePair,
) {
	out = pairs

	for n, pair := range pairs {
		if n != len(pairs) {
			for i, p := range pairs[n+1:] {
				if !(pair.Pool0 == p.Pool1 &&
					pair.Pool1 == p.Pool0) ||
					!(pair.Pool0 == p.Pool0 &&
						pair.Pool1 == p.Pool1) {
					continue
				} else {
					pairs = append(
						out[:i-1],
						out[i:]...,
					)
				}
			}
		}
	}

	return
}

func MakeTradePair(
	pools map[string]entities.Pool,
) (
	pair entities.TradePair,
	err error,
) {
	for _, pool := range pools {
		if pair.Pool0 != pool {
			pair.Pool0 = pool
		} else if pair.Pool1 != pool {
			pair.Pool1 = pool
		} else {
			err = fmt.Errorf(
				"pair with same pools %v",
				pool,
			)
			return
		}
	}

	return
}

func TokenPairContain(
	pairs []entities.TokenPair,
	searchEl entities.TokenPair,
) (
	index int,
	ok bool,
) {
	ok = false

	for n, pair := range pairs {
		if pair == searchEl {
			ok = true
			index = n

			return
		}
	}

	return
}

func CheckPairProtocol(
	pair entities.TradePair,
	protocol entities.SwapProtocol,
) (
	ok bool,
) {
	ok = false

	if pair.Pool0.Protocol == protocol &&
		pair.Pool1.Protocol == protocol {
		ok = true
	}

	return
}

func CheckPairTokens(
	pair entities.TradePair,
	tokens entities.TokenPair,
) (
	ok bool,
) {
	ok = false

	if pair.Pool0.Pair == tokens &&
		pair.Pool1.Pair == tokens {
		ok = true
	}

	return
}

func GetTokenFromPair(
	pair entities.TokenPair,
	addr string,
) (
	token *entities.Token,
) {
	switch addr {
	case pair.Token0.Address:
		token = &pair.Token0
	case pair.Token1.Address:
		token = &pair.Token1
	default:
		token = nil
	}

	return
}

func PoolContainPair(
	pool entities.Pool,
	pair entities.TokenPair,
) (
	ok bool,
) {
	if (pool.Pair.Token0 == pair.Token0 &&
		pool.Pair.Token1 == pair.Token1) ||
		(pool.Pair.Token0 == pair.Token1 &&
			pool.Pair.Token1 == pair.Token0) {
		ok = true

		return
	}
	ok = false

	return
}

func GetPoolsByPair(
	pools []entities.Pool,
	pair entities.TokenPair,
) (
	out []entities.Pool,
) {
	for _, pool := range pools {
		if PoolContainPair(pool, pair) {
			out = append(out, pool)
		}
	}

	return
}

func PairsFromTokens(tokens []entities.Token) (pairs []entities.TokenPair) {
	m := make(
		map[entities.Token]entities.TokenPair,
	)
	n := make(
		map[entities.TokenPair]int,
	)

	for index, token := range tokens {
		if _, ok := m[token]; !ok {
			for _, t := range tokens {
				if t != token {

					m[token] = entities.TokenPair{
						Token0: token,
						Token1: t,
					}

					n[m[token]] = index
				}
			}
		}

		continue
	}

	for k := range n {
		pairs = append(pairs, k)
	}

	return
}
