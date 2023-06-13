package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/antonyuhnovets/flash-loan-arbitrage/api"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ConnectBlockchain(url, key, provider string) {
	// connect to blockchain network
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Println("22 ", err)
		panic(err)
	}

	// private key of the deployer
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Println("29 ", err)
		panic(err)
	}

	// extract public key of the deployer from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("37 ")
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// address of the deployer
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// chain id of the network
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println("47 ", err)
		panic(err)
	}

	// Get Transaction Ops to make a valid Ethereum transaction
	auth, err := GetNextTransaction(client, fromAddress, privateKey, chainID)
	if err != nil {
		log.Println("54 ", err)
		panic(err)
	}

	loanKey, err := crypto.HexToECDSA(provider)
	if err != nil {
		panic(err)
	}

	loanPubKey := loanKey.Public()
	loanKeyECDSA, ok := loanPubKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	loanProvider := crypto.PubkeyToAddress(*loanKeyECDSA)

	log.Println(loanProvider)

	// deploy the contract
	address, tx, FlashLoanApi, err := api.DeployApi(auth, client, loanProvider)
	if err != nil {
		log.Println("60 ", err)
		panic(err)
	}

	fmt.Printf("Api contract deployed to %s\n", address.Hex())
	fmt.Printf("Tx: %s\n", tx.Hash().Hex())

	FlashLoanApi.GetBalance(nil, address)

	// Set Favorite Number
	// Get Transaction Ops to make a valid Ethereum transaction
	auth, err = GetNextTransaction(client, fromAddress, privateKey, chainID)
	if err != nil {
		log.Println("74 ", err)
		panic(err)
	}

}

// GetNextTransaction returns the next transaction in the pending transaction queue
// NOTE: this is not an optimized way
func GetNextTransaction(client *ethclient.Client, fromAddress common.Address, privateKey *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
	// nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	// sign the transaction
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)             // in wei
	auth.GasLimit = uint64(3000000)        // in units
	auth.GasPrice = big.NewInt(1000000000) // in wei

	return auth, nil
}
