package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/api"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Wallet struct {
	Address    common.Address
	PubKey     *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

func SetWallet(key string) (*Wallet, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}

	return &Wallet{privateKey: privateKey}, nil
}

func (wall *Wallet) setPubKey() error {
	pubKey, err := PullPublicKey(wall.privateKey)
	if err != nil {
		return err
	}
	wall.PubKey = pubKey

	return nil
}

func (wall *Wallet) setAddr() {
	wall.Address = crypto.PubkeyToAddress(*wall.PubKey)
}

func PullPublicKey(pk *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	pubKey := pk.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("wrong key type")
	}

	return publicKeyECDSA, nil

}

type Client struct {
	Client  *ethclient.Client
	Wallet  *Wallet
	ChainID *big.Int
}

func NewClient(url, key string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	cl := &Client{Client: client}

	err = cl.setupWallet(key)
	cl.setChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c *Client) setupWallet(key string) error {
	wall, err := SetWallet(key)
	if err != nil {
		return err
	}

	if err := wall.setPubKey(); err != nil {
		return err
	}

	wall.setAddr()

	return nil
}

func (c *Client) setChainID(ctx context.Context) error {
	chainId, err := c.Client.ChainID(ctx)
	if err != nil {
		return err
	}

	c.ChainID = chainId

	return nil
}

func (c *Client) GetChainID() *big.Int {
	return c.ChainID
}

func (c *Client) DialContract(address common.Address) (*api.Api, error) {
	contract, err := api.NewApi(address, c.Client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

// GetNextTransaction returns the next transaction in the pending transaction queue
// NOTE: this is not an optimized way
func (c *Client) GetNextTransaction() (*bind.TransactOpts, error) {
	// nonce
	nonce, err := c.Client.PendingNonceAt(context.Background(), c.Wallet.Address)
	if err != nil {
		return nil, err
	}

	// sign the transaction
	auth, err := bind.NewKeyedTransactorWithChainID(c.Wallet.privateKey, c.ChainID)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)             // in wei
	auth.GasLimit = uint64(3000000)        // in units
	auth.GasPrice = big.NewInt(1000000000) // in wei

	return auth, nil
}

// func ConnectBlockchain(url, key, provider string) {
// 	// connect to blockchain network

// 	// private key of the deployer
// 	privateKey, err := crypto.HexToECDSA(key)
// 	if err != nil {
// 		log.Println("29 ", err)
// 		panic(err)
// 	}

// 	// extract public key of the deployer from private key
// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		log.Println("37 ")
// 		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
// 	}

// 	// address of the deployer
// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

// 	// chain id of the network
// 	chainID, err := client.NetworkID(context.Background())
// 	if err != nil {
// 		log.Println("47 ", err)
// 		panic(err)
// 	}

// 	// Get Transaction Ops to make a valid Ethereum transaction
// 	auth, err := GetNextTransaction(client, fromAddress, privateKey, chainID)
// 	if err != nil {
// 		log.Println("54 ", err)
// 		panic(err)
// 	}

// 	loanKey, err := crypto.HexToECDSA(provider)
// 	if err != nil {
// 		panic(err)
// 	}

// 	loanPubKey := loanKey.Public()
// 	loanKeyECDSA, ok := loanPubKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		panic("invalid key")
// 	}

// 	loanProvider := crypto.PubkeyToAddress(*loanKeyECDSA)

// 	log.Println(loanProvider)

// 	// deploy the contract
// 	address, tx, FlashLoanApi, err := api.DeployApi(auth, client, loanProvider)
// 	if err != nil {
// 		log.Println("60 ", err)
// 		panic(err)
// 	}

// 	fmt.Printf("Api contract deployed to %s\n", address.Hex())
// 	fmt.Printf("Tx: %s\n", tx.Hash().Hex())

// 	FlashLoanApi.GetBalance(nil, address)

// 	// Set Favorite Number
// 	// Get Transaction Ops to make a valid Ethereum transaction
// 	auth, err = GetNextTransaction(client, fromAddress, privateKey, chainID)
// 	if err != nil {
// 		log.Println("74 ", err)
// 		panic(err)
// 	}

// }
