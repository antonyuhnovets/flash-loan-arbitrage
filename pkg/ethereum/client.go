package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

// GetNextTransaction returns the next transaction in the pending transaction queue
func (c *Client) GetNextTransaction(ctx context.Context) (
	opts *bind.TransactOpts,
	err error,
) {
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
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice       // in wei

	return auth, nil
}

func (c *Client) UpdateChainID(ctx context.Context) (
	err error,
) {
	chainID := c.GetChainID(ctx)
	if chainID == nil {
		err = fmt.Errorf("fail getting chain id")

		return
	}
	c.setChainID(chainID)

	return
}

// not used
func (c *Client) SignTx(ctx context.Context, tx *types.Transaction) (
	t *types.Transaction,
	err error,
) {
	// s := types.NewEIP155Signer(tx.ChainId())
	s := types.NewLondonSigner(c.ChainID)
	t, err = types.SignTx(tx, s, c.Wallet.privateKey)

	return
}

func (c *Client) ReplaceTx(ctx context.Context, hash string) (
	opts *bind.TransactOpts,
	err error,
) {

	opts, err = c.GetNextTransaction(ctx)
	if err != nil {
		fmt.Println(err)

		return
	}

	trans, _, err := c.Client.TransactionByHash(
		ctx,
		ToHash(hash))

	if err != nil {
		return
	}

	opts.Nonce = ToBigInt(int(trans.Nonce()))

	return
}
