package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func IsTransaction(t interface{}) (
	tx *types.Transaction,
	ok bool,
) {
	tx, ok = t.(*types.Transaction)

	return
}

func GetCallMsg(from, to string) ethereum.CallMsg {
	toAddr := ToAddress(to)
	msg := ethereum.CallMsg{
		From:     ToAddress(from),
		To:       &toAddr,
		Gas:      21000,
		GasPrice: big.NewInt(18000000000),
		Value:    big.NewInt(666),
		Data:     nil,
	}
	return msg
}

func ToAddress(s string) common.Address {
	return common.HexToAddress(s)
}

func ToHash(s string) common.Hash {
	return common.HexToHash(s)
}

func FromAddress(a common.Address) string {
	return a.Hex()
}

func FromBigInt(b *big.Int) int {
	return int(b.Int64())
}

func ToBigInt(b int) *big.Int {
	return big.NewInt(int64(b))
}

func CallOpts(opts ...interface{}) (
	b *bind.CallOpts,
) {
	b = &bind.CallOpts{}

	for _, opt := range opts {
		switch t := opt.(type) {
		case bool:
			b.Pending = t
		case context.Context:
			b.Context = t
		case *big.Int:
			b.BlockNumber = t
		case common.Address:
			b.From = t
		}
	}

	return
}

func PullPublicKey(pk *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	pubKey := pk.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("wrong key type")
	}

	return publicKeyECDSA, nil

}
