package parser

import (
	"math/big"

	uni "github.com/ackermanx/ethclient/uniswap"
	sushi "github.com/ebadiere/go-defi/sushiswap"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
)

type Protocol struct {
	p entities.SwapProtocol
	ProtocolResolver
}

func NewProtocol(sp entities.SwapProtocol, pr ProtocolResolver) (
	p *Protocol,
) {
	p = &Protocol{sp, pr}

	return
}

func (pro *Protocol) GetPoolAddress(pair entities.TokenPair) (
	address string,
	err error,
) {
	address, err = pro.Resolve(pro.p).GetPoolAddress(pair)

	return
}

func (pro *Protocol) GetProtocol() (
	sp entities.SwapProtocol,
) {
	return pro.p
}

func (pro *Protocol) SetProtocol(sp entities.SwapProtocol) {
	pro.p = sp
}

type ProtocolResolver interface {
	Resolve(entities.SwapProtocol) ProtocolParser
}

type ProtocolParser interface {
	GetPoolAddress(entities.TokenPair) (string, error)
}
type UniV2 struct {
	Sp entities.SwapProtocol
}

func (u2 UniV2) GetPoolAddress(
	pair entities.TokenPair,
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

type UniV3 struct {
	Sp entities.SwapProtocol
}

func (u3 UniV3) GetPoolAddress(
	pair entities.TokenPair,
) (
	address string,
	err error,
) {
	pAddr, err := uni.CalculatePoolAddressV3(
		pair.Token0.Address,
		pair.Token1.Address,
		big.NewInt(3000),
	)
	if err != nil {
		return
	}

	address = pAddr.Hex()

	return
}

type SushiV2 struct {
	Sp entities.SwapProtocol
}

func (s2 SushiV2) GetPoolAddress(
	pair entities.TokenPair,
) (
	address string,
	err error,
) {
	pAddr := sushi.GeneratePairAddress(
		ethereum.ToAddress(pair.Token0.Address),
		ethereum.ToAddress(pair.Token1.Address),
	)

	address = pAddr.Hex()

	return
}
