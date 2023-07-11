package parser

import (
	"os"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	prs "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/parser"
)

type Parser struct {
	*prs.Parser
	*prs.Protocol
}

func NewParser(sp entities.SwapProtocol) (
	p *Parser,
) {
	return &Parser{
		prs.New(),
		prs.NewProtocol(sp, NewProtoResolver()),
	}
}

func (p *Parser) SetProtocol(sp entities.SwapProtocol) {
	p.Protocol.SetProtocol(sp)
}

func (p *Parser) Parse(pairs []entities.TokenPair) (
	err error,
) {
	for _, pair := range pairs {
		addr, e := p.Protocol.GetPoolAddress(pair)
		if err != nil {
			err = e

			return
		}
		p.Parser.AddPool(entities.TradePool{
			Protocol: p.Protocol.GetProtocol(),
			Address:  addr,
			Pair:     pair,
		})
	}
	return
}

type protocols struct {
	uni2   prs.UniV2
	uni3   prs.UniV3
	sushi2 prs.SushiV2
}

func NewProtoResolver() *protocols {
	u2 := entities.SwapProtocol{
		Name:       "Uniswap-V2",
		Factory:    os.Getenv("UNI_V2_FACTORY_ADDRESS"),
		SwapRouter: os.Getenv("UNI_V2_SWAP_ROUTER_ADDRESS"),
	}

	s2 := entities.SwapProtocol{
		Name:       "Sushiswap-V2",
		Factory:    os.Getenv("SUSHI_V2_FACTORY_ADDRESS"),
		SwapRouter: os.Getenv("SUSHI_V2_SWAP_ROUTER_ADDRESS"),
	}

	u3 := entities.SwapProtocol{
		Name:       "Uniswap-V3",
		Factory:    os.Getenv("UNI_V3_FACTORY_ADDRESS"),
		SwapRouter: os.Getenv("UNI_V3_SWAP_ROUTER_ADDRESS"),
	}

	return &protocols{
		uni2:   prs.UniV2{Sp: u2},
		uni3:   prs.UniV3{Sp: u3},
		sushi2: prs.SushiV2{Sp: s2},
	}
}

func (proto *protocols) Resolve(sp entities.SwapProtocol) (
	pp prs.ProtocolParser,
) {
	switch sp {
	case proto.uni2.Sp:
		pp = proto.uni2
	case proto.uni3.Sp:
		pp = proto.uni3
	case proto.sushi2.Sp:
		pp = proto.sushi2
	}

	return
}
