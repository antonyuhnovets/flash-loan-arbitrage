package ethereum

import (
	"context"
	"net"

	"github.com/ethereum/go-ethereum/rpc"
)

type ClientIPC struct {
	Client *rpc.Client
}

type ServerIPC struct {
	address string
	Server  *rpc.Server
	L       net.Listener
	Api     []rpc.API
}

func NewEndpointIPC(name string, api []rpc.API) (*ServerIPC, error) {
	l, s, err := rpc.StartIPCEndpoint(name, api)
	if err != nil {
		return nil, err
	}
	return &ServerIPC{
		address: name,
		Server:  s,
		L:       l,
		Api:     api,
	}, nil
}

func (srv *ServerIPC) Start() error {
	err := srv.Server.ServeListener(srv.L)
	if err != nil {
		return err
	}
	return nil
}

func (srv *ServerIPC) AddClient(ctx context.Context) (*ClientIPC, error) {
	cl, err := rpc.DialIPC(ctx, srv.address)
	if err != nil {
		return nil, err
	}

	return &ClientIPC{Client: cl}, nil
}

func (srv *ServerIPC) DialInProc() *ClientIPC {
	cl := rpc.DialInProc(srv.Server)

	return &ClientIPC{cl}
}
