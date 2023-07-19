package parser

import (
	"fmt"
	"math/big"

	uni "github.com/ackermanx/ethclient/uniswap"
	sushi "github.com/ebadiere/go-defi/sushiswap"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
)

type ProtocolManager struct {
	p []*Protocol
	ProtocolResolver
}

func NewManager(pr ProtocolResolver) *ProtocolManager {
	return &ProtocolManager{
		make([]*Protocol, 0),
		pr,
	}
}

func (pm *ProtocolManager) AddProtocol(sp entities.SwapProtocol) (
	err error,
) {
	for _, proto := range pm.p {
		if proto.SwapProtocol == sp {
			err = fmt.Errorf("protocol %s already added", sp.Name)

			return
		}
	}
	pm.p = append(pm.p, NewProtocol(sp))

	return
}

func (pm *ProtocolManager) RemoveProtocol(sp entities.SwapProtocol) (
	err error,
) {
	err = fmt.Errorf("protocol not found")

	for index, proto := range pm.ListProtocols() {
		if proto == sp {
			pm.p = append(pm.p[:index], pm.p[index+1:]...)
			err = nil
		}
	}
	return
}

func (pm *ProtocolManager) GetPoolAddresses(pair entities.TokenPair) (
	out map[entities.SwapProtocol]string,
	err error,
) {
	out = make(map[entities.SwapProtocol]string)

	for _, proto := range pm.p {
		parser, _err := pm.ProtocolResolver.Resolve(proto)
		if err != nil {
			err = _err
			return

		}

		if parser == nil {
			err = fmt.Errorf("parser is nil")
			fmt.Println(err)
			return
		}

		if pair.Token0.Address == "" || pair.Token1.Address == "" {
			err = fmt.Errorf("pair token address is nil")
			fmt.Println(err)
			return
		}

		address, _err := parser.GetPoolAddress(pair)
		if err != nil {
			err = _err

			return
		}

		out[proto.GetProtocolData()] = address
	}

	return
}

func (pm *ProtocolManager) ListProtocols() (
	out []entities.SwapProtocol,
) {
	for _, proto := range pm.p {
		p := proto.GetProtocolData()
		out = append(out, p)
	}

	return
}

type Protocol struct {
	entities.SwapProtocol
}

func NewProtocol(sp entities.SwapProtocol) (
	p *Protocol,
) {
	p = &Protocol{sp}

	return
}

func (pro *Protocol) GetProtocolData() (
	sp entities.SwapProtocol,
) {
	return pro.SwapProtocol
}

type ProtocolResolver interface {
	Resolve(*Protocol) (ProtocolParser, error)
}

type ProtocolParser interface {
	GetPoolAddress(entities.TokenPair) (string, error)
}

type UniV2 Protocol

func (u2 *UniV2) GetPoolAddress(
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

type UniV3 Protocol

func (u3 *UniV3) GetPoolAddress(
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

type SushiV2 Protocol

func (s2 *SushiV2) GetPoolAddress(
	pair entities.TokenPair,
) (
	address string,
	err error,
) {
	sushi.FactoryAddress = ethereum.ToAddress(s2.Factory)
	sushi.Router02Address = ethereum.ToAddress(s2.SwapRouter)

	pAddr := sushi.GeneratePairAddress(
		ethereum.ToAddress(pair.Token0.Address),
		ethereum.ToAddress(pair.Token1.Address),
	)

	address = pAddr.Hex()

	return
}
