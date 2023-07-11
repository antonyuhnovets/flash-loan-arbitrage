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
	pools []entities.TradePool,
	err error,
) {
	var res []entities.TradePool

	out, err := s.fst.Read(
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

func (s *Storage) AddPool(
	ctx c.Context,
	pool entities.TradePool,
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
		return
	}
	err = s.fst.Store(
		ctx,
		where,
		[]byte("]"),
	)

	return
}

func (s *Storage) StorePools(
	ctx c.Context,
	where string,
	pools []entities.TradePool,
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
	pools []entities.TradePool,
	err error,
) {
	out, err := s.fst.Read(
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
	b, err := s.fst.Read(ctx, where)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &tokens)

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
