package repo

import (
	c "context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
)

type FileStorage struct {
	path string
	// f    *os.File
}

func UseFile(path string) *FileStorage {
	return &FileStorage{
		path: path,
	}
}

func NewFile(
	path string,
) (
	*FileStorage,
	error,
) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return &FileStorage{path}, nil
}

func (fs *FileStorage) Store(
	ctx c.Context, item []byte,
) error {
	f, err := os.OpenFile(fs.path,
		os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	index, err := f.Write(item)
	if err != nil {
		return err
	}

	log.Printf(
		"%b bytes saved to file",
		index,
	)

	return nil
}

func (fs *FileStorage) Read(
	ctx c.Context,
) (
	[]byte,
	error,
) {
	var b []byte
	b, err := os.ReadFile(fs.path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (fs *FileStorage) GetByTokens(
	ctx c.Context, tokens TokenPair,
) (
	[]TradePool,
	error,
) {
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
			res = append(
				res[:n],
				res[n+1:]...,
			)
		}
	}

	return res, nil
}

func (fs *FileStorage) StorePool(
	ctx c.Context, pool TradePool,
) error {
	b, err := json.Marshal(pool)
	if err != nil {
		return err
	}

	data := string(b) + "\n"
	err = fs.Store(ctx, []byte(data))
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) StorePools(
	ctx c.Context, pools []TradePool,
) error {
	for _, pool := range pools {
		err := fs.StorePool(ctx, pool)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *FileStorage) ListPools(
	ctx c.Context,
) (
	[]TradePool,
	error,
) {
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

func (fs *FileStorage) GetTokenByAddress(
	ctx c.Context, addr string,
) (
	Token,
	error,
) {
	pools, err := fs.ListPools(ctx)
	if err != nil {
		return Token{}, err
	}

	for _, pool := range pools {
		if pool.Pair.Token0.Address == addr {
			return pool.Pair.Token0, nil
		} else if pool.Pair.Token1.Address == addr {
			return pool.Pair.Token1, nil
		} else {
			continue
		}
	}

	return Token{}, fmt.Errorf("token with address %s not found", addr)
}

func (fs *FileStorage) Clear(ctx c.Context) error {
	return os.Truncate(fs.path, 0)
}
