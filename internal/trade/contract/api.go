package contract

import (
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/api"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	eth "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
)

type Contract struct {
	address     string
	contractApi *api.Api
	bind.SignerFn
}

func New(addr string, cli *eth.Client) (
	out *Contract,
	err error,
) {
	ca, err := api.NewApi(
		eth.ToAddress(addr),
		cli.Client,
	)
	out = &Contract{
		address:     addr,
		contractApi: ca,
	}

	return
}

type Api interface {
	Address() string
	Api() API
}

func (ca *Contract) Api() API {
	return &_api{ca.contractApi}
}

func (ca *Contract) Address() string {
	return ca.address
}

type _api struct {
	a *api.Api
}

func (a *_api) Api() *api.Api {
	return a.a
}

func (a *_api) Transactor() *api.ApiTransactor {
	return &a.a.ApiTransactor
}

func (a *_api) Caller() *api.ApiCaller {
	return &a.a.ApiCaller
}

func (a *_api) Filterer() *api.ApiFilterer {
	return &a.a.ApiFilterer
}

type API interface {
	Api() *api.Api
	Transactor() *api.ApiTransactor
	Caller() *api.ApiCaller
	Filterer() *api.ApiFilterer
}
