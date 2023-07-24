package parser

import (
	"fmt"
	"os"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	prs "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/parser"
)

type Parser struct {
	*prs.Parser
	*prs.ProtocolManager
}

func NewParser() (
	p *Parser,
) {
	return &Parser{
		prs.New(),
		prs.NewManager(NewProtoResolver()),
	}
}

func (p *Parser) Parse(pairs []entities.TokenPair) (
	err error,
) {
	for _, pair := range pairs {
		m, _err := p.GetPoolAddresses(pair)
		if err != nil {
			err = _err

			return
		}
		for proto, addr := range m {
			if !p.containPool(addr) {
				p.AddPool(entities.Pool{
					Protocol: proto,
					Address:  addr,
					Pair:     pair,
				},
				)
			}
		}
	}
	return
}

func (p *Parser) containPool(addr string) bool {
	for _, pool := range p.ListPools() {
		if pool.Address == addr {
			return true
		}
	}
	return false
}

type Protocols struct {
	uni2   *prs.UniV2
	uni3   *prs.UniV3
	sushi2 *prs.SushiV2
}

func NewProtoResolver() *Protocols {
	u2 := entities.SwapProtocol{
		Name:       "Uniswap-V2",
		Factory:    os.Getenv("UNI_V2_FACTORY_ADDRESS"),
		SwapRouter: os.Getenv("UNI_V2_SWAP_ROUTER_ADDRESS"),
	}

	s2 := entities.SwapProtocol{
		ID:         1,
		Name:       "Sushiswap-V2",
		Factory:    os.Getenv("SUSHI_V2_FACTORY_ADDRESS"),
		SwapRouter: os.Getenv("SUSHI_V2_SWAP_ROUTER_ADDRESS"),
	}

	u3 := entities.SwapProtocol{
		Name:       "Uniswap-V3",
		Factory:    os.Getenv("UNI_V3_FACTORY_ADDRESS"),
		SwapRouter: os.Getenv("UNI_V3_SWAP_ROUTER_ADDRESS"),
	}

	return &Protocols{
		uni2:   &prs.UniV2{SwapProtocol: u2},
		uni3:   &prs.UniV3{SwapProtocol: u3},
		sushi2: &prs.SushiV2{SwapProtocol: s2},
	}
}

func (proto *Protocols) Resolve(sp *prs.Protocol) (
	pp prs.ProtocolParser,
	err error,
) {
	switch sp.GetProtocolData().Name {
	case proto.uni2.SwapProtocol.Name:
		proto.checkAndSetProtoData(sp, "uni2")
		pp = proto.uni2

	case proto.uni3.SwapProtocol.Name:
		proto.checkAndSetProtoData(sp, "uni3")
		pp = proto.uni3

	case proto.sushi2.SwapProtocol.Name:
		proto.checkAndSetProtoData(sp, "sushi2")
		pp = proto.sushi2

	}

	if pp == nil {
		err = fmt.Errorf("protocol %s not found in resolver", sp.Name)
	}

	return
}

func (proto *Protocols) checkAndSetProtoData(p *prs.Protocol, which string) {
	var t interface{}
	data := p.GetProtocolData()

	switch which {
	case "sushi2":
		t = proto.sushi2
	case "uni2":
		t = proto.uni2
	case "uni3":
		t = proto.uni3
	}

	tp := t.(*prs.Protocol)

	if data.Factory != tp.Factory {
		tp.Factory = data.Factory
	}
	if data.SwapRouter != tp.SwapRouter {
		tp.SwapRouter = data.SwapRouter
	}
}
