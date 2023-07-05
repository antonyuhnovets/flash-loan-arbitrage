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
	files map[string]string
}

func NewStorage() (
	fs *FileStorage,
) {
	fs = &FileStorage{
		files: make(map[string]string),
	}

	return
}

func (fs *FileStorage) Setup(
	files map[string]string,
) (
	err error,
) {
	for k, v := range files {
		err = fs.NewFile(k, v)
		if err != nil {
			return
		}
	}

	return
}

func (fs *FileStorage) UseFile(
	name, path string,
) {
	fs.files[name] = path
}

func (fs *FileStorage) NewFile(
	name, path string,
) (

	err error,
) {
	f, err := os.Create(path)

	if err != nil {
		return
	}
	defer f.Close()

	fs.UseFile(
		name,
		path,
	)

	fs.Store(
		c.Background(),
		name,
		[]byte("[\n"),
	)
	err = fs.Store(
		c.Background(),
		name,
		[]byte("]"),
	)

	return
}

func (fs *FileStorage) Store(
	ctx c.Context,
	where string,
	item []byte,
) (
	err error,
) {
	f, err := os.OpenFile(
		fs.files[where],
		os.O_APPEND|os.O_RDWR,
		0644,
	)
	if err != nil {
		return
	}

	n, err := f.Write(item)
	log.Println(n)

	return
}

func (fs *FileStorage) Read(
	ctx c.Context,
	where string,
) (
	b []byte,
	err error,
) {
	b, err = os.ReadFile(
		fs.files[where],
	)

	return
}

func (fs *FileStorage) GetByTokens(
	ctx c.Context,
	where string,
	tokens TokenPair,
) (
	pools []TradePool,
	err error,
) {
	var res []TradePool
	out, err := fs.Read(
		ctx,
		where,
	)
	if err != nil {
		return
	}

	err = json.Unmarshal(
		out,
		&res,
	)
	if err != nil {
		return
	}

	for _, pool := range res {
		if pool.Pair != tokens {
			continue
		} else {
			pools = append(pools, pool)
		}
	}

	return
}

func (fs *FileStorage) AddPool(
	ctx c.Context,
	pool TradePool,
	where string,
) (
	err error,
) {
	b, err := json.Marshal(pool)
	if err != nil {
		return
	}

	err = fs.rmCloser(
		ctx,
		where,
		[]byte(string(b[0:])+"\n"),
	)
	if err != nil {
		return
	}
	err = fs.Store(
		ctx,
		where,
		[]byte("]"),
	)

	return
}

func (fs *FileStorage) rmCloser(
	ctx c.Context,
	where string,
	item []byte,
) (
	err error,
) {
	f, err := os.OpenFile(
		fs.files[where],
		os.O_RDWR,
		0644,
	)
	if err != nil {
		return
	}
	b, err := fs.Read(ctx, where)
	if err != nil {
		return
	}

	if string(b[len(b)-3]) != "[" {
		b = append(b[:len(b)-2], []byte(",\n")...)
	} else {
		b = append(b[:len(b)-2], []byte("\n")...)
	}
	b = append(b, item...)

	n, err := f.Write(b)
	log.Println(n)

	return
}

func (fs *FileStorage) StorePools(
	ctx c.Context,
	where string,
	pools []TradePool,
) (
	err error,
) {
	for _, pool := range pools {
		err = fs.AddPool(
			ctx,
			pool,
			where,
		)
		if err != nil {
			return
		}
	}

	return
}

func (fs *FileStorage) ListPools(
	ctx c.Context,
	where string,
) (
	pools []TradePool,
	err error,
) {
	out, err := fs.Read(
		ctx,
		where,
	)
	if err != nil {
		return
	}

	err = json.Unmarshal(
		out,
		&pools,
	)

	return
}

func (fs *FileStorage) AddToken(
	ctx c.Context,
	where string,
	token Token,
) (
	err error,
) {
	b, err := json.Marshal(token)
	if err != nil {
		return
	}

	err = fs.rmCloser(
		ctx,
		where,
		[]byte(string(b[0:])+"\n"),
	)
	if err != nil {
		return
	}
	err = fs.Store(
		ctx,
		where,
		[]byte("]"),
	)

	return
}

func (fs *FileStorage) StoreTokens(
	ctx c.Context,
	where string,
	tokens []Token,
) (
	err error,
) {
	for _, token := range tokens {
		err = fs.AddToken(
			ctx,
			where,
			token,
		)
		if err != nil {
			return
		}
	}

	return
}

func (fs *FileStorage) ListTokens(
	ctx c.Context,
	where string,
) (
	tokens []Token,
	err error,
) {
	b, err := fs.Read(ctx, where)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &tokens)

	return
}

func (fs *FileStorage) GetTokenByAddress(
	ctx c.Context,
	where, address string,
) (
	token Token,
	err error,
) {
	tokens, err := fs.ListTokens(
		ctx,
		where,
	)
	if err != nil {
		return
	}

	for _, t := range tokens {
		if t.Address == address {
			token = t
			return
		}
	}

	// pools, err := fs.ListPools(
	// 	ctx,
	// 	where,
	// )
	// if err != nil {
	// 	return
	// }

	// for _, pool := range pools {
	// 	t := getTokenFromPair(
	// 		pool.Pair,
	// 		address,
	// 	)
	// 	if t != nil {
	// 		token = *t
	// 		return
	// 	} else {
	// 		continue
	// 	}
	// }

	err = fmt.Errorf(
		"token with address %s not found\n",
		address,
	)

	return
}

func (fs *FileStorage) Clear(
	ctx c.Context,
	where string,
) (
	err error,
) {
	err = os.Truncate(
		fs.files[where],
		0,
	)

	return
}

func (fs *FileStorage) ClearAll(
	ctx c.Context,
) (
	err error,
) {
	for k, v := range fs.files {
		err = os.Remove(v)
		if err != nil {
			return
		}
		fs.files[k] = ""
	}

	return
}

func getTokenFromPair(
	pair TokenPair,
	addr string,
) (
	token *Token,
) {
	switch addr {
	case pair.Token0.Address:
		token = &pair.Token0
	case pair.Token1.Address:
		token = &pair.Token1
	default:
		token = nil
	}

	return
}
