package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// func SignTransactAndSend(c context.Context, tx *types.Transaction) {
// 	ctx, cancel := context.WithCancel(c)
// 	defer cancel()

// 	bind.NewTransactorWithChainID()
// }

func Logs() types.Log {
	return types.Log{}
}
func ToAddress(s string) common.Address {
	return common.HexToAddress(s)
}

func ListAddress(addrs ...common.Address) (
	out []common.Address,
) {
	out = append(out, addrs...)

	return
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
