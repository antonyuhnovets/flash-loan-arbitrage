package ethereum

import (
	"context"
	"fmt"
	"math/big"

	// "github.com/antonyuhnovets/flash-loan-arbitrage/internal/api"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	Client  *ethclient.Client
	Wallet  *Wallet
	ChainID *big.Int
}

func NewClient(url string) (
	cl *Client,
	err error,
) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return
	}
	cl = &Client{Client: client}

	return
}

func (c *Client) Setup(ctx context.Context, wall *Wallet) (
	err error,
) {
	c.UseWallet(wall)
	c.setChainID(c.GetChainID(ctx))

	return
}

func (c *Client) ClientGet() (
	cl interface{},
) {
	cl = c.Client

	return
}

func (c *Client) UseWallet(wall *Wallet) {
	c.Wallet = wall
}

func (c *Client) GetBallance(ctx context.Context) (
	ball int,
	err error,
) {
	b, err := c.Client.BalanceAt(ctx, c.Wallet.Address, c.ChainID)
	ball = int(b.Int64())

	return
}

func (c *Client) setChainID(chainId *big.Int) {
	c.ChainID = chainId
}

func (c *Client) GetChainID(ctx context.Context) *big.Int {
	chainId, err := c.Client.ChainID(ctx)
	if err != nil {
		return nil
	}

	return chainId
}

// func (c *Client) DialContract(address common.Address) (api.Api, error) {
// 	contract, err := api.NewApi(
// 		address,
// 		c.Client,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return contract, nil
// }

// GetNextTransaction returns the next transaction in the pending transaction queue
// NOTE: this is not an optimized way
func (c *Client) GetNextTransaction(ctx context.Context) (
	opts *bind.TransactOpts,
	err error,
) {
	// cont, cancel := context.WithCancel(ctx)
	// defer cancel()

	// c.UpdateChainID(cont)

	gasPrice, err := c.Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// nonce
	nonce, err := c.Client.PendingNonceAt(ctx, c.Wallet.Address)
	if err != nil {
		return
	}

	c.UpdateChainID(ctx)
	// sign the transaction
	auth, err := bind.NewKeyedTransactorWithChainID(
		c.Wallet.privateKey,
		c.ChainID,
	)
	if err != nil {
		return
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice        // in wei

	return auth, nil
}

func (c *Client) SignTx(ctx context.Context, tx *types.Transaction) (
	t *types.Transaction,
	err error,
) {
	// err = c.UpdateChainID(ctx)
	// if err != nil {
	// 	return
	// }
	// s := types.NewEIP155Signer(tx.ChainId())
	s := types.NewLondonSigner(c.ChainID)
	t, err = types.SignTx(tx, s, c.Wallet.privateKey)

	return
}

func (c *Client) UpdateChainID(ctx context.Context) (
	err error,
) {
	// cont, cancel := context.WithCancel(ctx)
	// defer cancel()

	chainID := c.GetChainID(ctx)
	if chainID == nil {
		err = fmt.Errorf("fail getting chain id")

		return
	}
	c.setChainID(chainID)

	return
}

func (c *Client) Transact(ctx context.Context, t *types.Transaction) (
	tx *types.Transaction,
	err error,
) {
	cont, cancel := context.WithCancel(ctx)
	defer cancel()

	c.UpdateChainID(cont)

	c.Client.SendTransaction(cont, t)

	tx, err = c.SignTx(cont, t)

	return
}
