package repo

import (
	c "context"
	"encoding/json"
	"log"
	"os"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type FileStorage struct {
	path string
	f    *os.File
}

func NewFile(path string) (*FileStorage, error) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return &FileStorage{path, f}, nil
}

func (fs *FileStorage) Store(ctx c.Context, item []byte) error {
	index, err := fs.f.Write(item)
	defer fs.f.Close()
	if err != nil {
		return err
	}

	log.Printf("%b bytes saved to file", index)

	return nil
}

func (fs *FileStorage) Read(ctx c.Context) ([]byte, error) {
	var b []byte

	index, err := fs.f.Read(b)
	if err != nil {
		return nil, err
	}

	log.Printf("%b bytes readed from file", index)

	return b, nil
}

func (fs *FileStorage) GetByTokens(ctx c.Context, tokens TokenPair) ([]TradePool, error) {
	res := make([]TradePool, 0)

	out, err := fs.Read(ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &res)
	if err != nil {
		return nil, err
	}

	for n, pool := range res {
		if pool.Pair != tokens {
			res = append(res[:n], res[n+1:]...)
		}
	}

	return res, nil
}

func (fs *FileStorage) StorePool(ctx c.Context, pool TradePool) error {
	b, err := json.Marshal(pool)
	if err != nil {
		return err
	}

	err = fs.Store(ctx, b)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) ListAllPools(ctx c.Context) ([]TradePool, error) {
	pools := make([]TradePool, 0)

	out, err := fs.Read(ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &pools)
	if err != nil {
		return nil, err
	}

	return pools, nil
}
