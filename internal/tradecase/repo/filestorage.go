package repo

import (
	c "context"
	"encoding/json"
	"fmt"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	fs "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/filestorage"
)

type Storage struct {
	fst *fs.FileStorage
}

func NewStorage(files map[string]string) (
	st *Storage,
	err error,
) {
	s := fs.NewStorage()

	err = s.Setup(files)
	if err != nil {
		return
	}

	st = &Storage{s}

	return
}

func (s *Storage) GetStorage() tradecase.Storage {
	return s.fst
}

func (s *Storage) GetByTokens(
	ctx c.Context,
	where string,
	tokens entities.TokenPair,
) (
	pools []entities.Pool,
	err error,
) {
	// var res []entities.TradePool

	err = s.fst.Read(
		ctx,
		where,
		&pools,
	)
	if err != nil {
		return
	}

	return
}

func (s *Storage) AddPool(
	ctx c.Context,
	pool entities.Pool,
	where string,
) (
	err error,
) {
	b, err := json.Marshal(pool)
	if err != nil {
		return
	}

	err = s.fst.ContinueFile(
		ctx,
		where,
		[]byte(string(b[0:])+"\n"),
	)
	if err != nil {
		if err != nil {
			fmt.Printf("/tradecase/repo/filestorage:95 - AddPool continue file\n%s", err)
			return
		}
		return
	}
	err = s.fst.Store(
		ctx,
		where,
		[]byte("]"),
	)
	if err != nil {
		fmt.Printf("/tradecase/repo/filestorage:106 - AddPool store \n%s", err)
	}

	return
}

func (s *Storage) StorePools(
	ctx c.Context,
	where string,
	pools []entities.Pool,
) (
	err error,
) {
	for _, pool := range pools {
		err = s.AddPool(
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

func (s *Storage) ListPools(
	ctx c.Context,
	where string,
) (
	pools []entities.Pool,
	err error,
) {
	err = s.fst.Read(
		ctx,
		where,
		&pools,
	)

	return
}

func (s *Storage) RemovePool(
	ctx c.Context,
	where string,
	pool entities.Pool,
) (
	err error,
) {
	b, err := json.Marshal(pool)
	if err != nil {
		return
	}

	err = s.fst.Remove(ctx, where, b)

	return
}

func (s *Storage) RemovePools(
	ctx c.Context,
	where string,
	pools []entities.Pool,
) (
	out []entities.Pool,
	err error,
) {
	for _, pool := range pools {
		err = s.RemovePool(ctx, where, pool)
		if err != nil {
			return
		}
	}

	out, err = s.ListPools(ctx, where)

	return
}

func (s *Storage) AddToken(
	ctx c.Context,
	where string,
	token entities.Token,
) (
	err error,
) {
	b, err := json.Marshal(token)
	if err != nil {
		return
	}

	err = s.fst.ContinueFile(
		ctx,
		where,
		[]byte(string(b[0:])+"\n"),
	)
	if err != nil {
		return
	}
	err = s.fst.Store(
		ctx,
		where,
		[]byte("]"),
	)

	return
}

func (s *Storage) RemoveToken(
	ctx c.Context,
	where string,
	token entities.Token,
) (
	err error,
) {
	b, err := json.Marshal(token)
	if err != nil {
		return
	}

	err = s.fst.Remove(ctx, where, b)

	return
}

func (s *Storage) RemoveTokens(
	ctx c.Context,
	where string,
	tokens []entities.Token,
) (
	out []entities.Token,
	err error,
) {
	for _, token := range tokens {
		err = s.RemoveToken(ctx, where, token)
		if err != nil {
			return
		}
	}

	out, err = s.ListTokens(ctx, where)

	return
}

func (s *Storage) StoreTokens(
	ctx c.Context,
	where string,
	tokens []entities.Token,
) (
	err error,
) {
	for _, token := range tokens {
		err = s.AddToken(
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

func (s *Storage) ListTokens(
	ctx c.Context,
	where string,
) (
	tokens []entities.Token,
	err error,
) {
	err = s.fst.Read(ctx, where, tokens)

	return
}

func (s *Storage) GetTokenByAddress(
	ctx c.Context,
	where, address string,
) (
	token entities.Token,
	err error,
) {
	tokens, err := s.ListTokens(
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

	err = fmt.Errorf(
		"token with address %s not found",
		address,
	)

	return
}

func (s *Storage) ClearAll(ctx c.Context) (
	err error,
) {
	err = s.fst.ClearAll(ctx)

	return
}
