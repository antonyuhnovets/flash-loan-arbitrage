package ethereum

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	Address    common.Address
	PubKey     *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

func NewWallet() *Wallet {
	return &Wallet{}
}

func (wall *Wallet) Setup(key string) (
	err error,
) {
	err = wall.setPK(key)
	if err != nil {
		return
	}
	err = wall.setPubKey()
	if err != nil {
		return
	}
	wall.setAddr()

	return
}

func (wall *Wallet) setPK(key string) (
	err error,
) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return
	}
	wall.privateKey = privateKey

	return
}

func (wall *Wallet) setPubKey() (
	err error,
) {
	pubKey, err := PullPublicKey(wall.privateKey)
	if err != nil {
		return
	}
	wall.PubKey = pubKey

	return
}

func (wall *Wallet) setAddr() {
	wall.Address = crypto.PubkeyToAddress(*wall.PubKey)
}
