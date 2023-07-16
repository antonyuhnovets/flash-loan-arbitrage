package repo

import (
	c "context"
	"fmt"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/postgres"
)

type PostgresRepo struct {
	ps *postgres.Storage
	ss postgres.Serializers
}

func New(conf config.Database) (
	pr *PostgresRepo,
	err error,
) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Etc/UTC",
		conf.Host, conf.Username, conf.Password, conf.Name, conf.Port,
	)

	conn, err := postgres.Connect(
		dsn,
		&entities.Token{},
		&entities.Pool{},
	)
	if err != nil {
		return
	}

	srls := postgres.NewSerializers()
	srls.Set()
	srls.RegisterAll()

	pr = &PostgresRepo{ps: conn, ss: srls}

	return
}

func (pr *PostgresRepo) GetStorage() tradecase.Storage {
	return pr.ps
}

func (pr *PostgresRepo) ListPools(
	ctx c.Context, table string,
) (
	pools []entities.Pool,
	err error,
) {
	err = pr.GetStorage().Read(ctx, table, &pools)
	if err != nil {

		return
	}

	return
}

func (pr *PostgresRepo) StorePools(
	ctx c.Context, table string, pools []entities.Pool,
) (
	err error,
) {
	for _, pool := range pools {
		err = pr.AddPool(ctx, pool, table)
		if err != nil {
			return
		}
	}

	return
}

func (pr *PostgresRepo) AddPool(
	ctx c.Context, pool entities.Pool, table string,
) (
	err error,
) {
	err = pr.GetStorage().Store(ctx, table, &pool)

	return
}

func (pr *PostgresRepo) GetByTokens(
	ctx c.Context, table string, pair entities.TokenPair,
) (
	pools []entities.Pool,
	err error,
) {
	err = pr.GetStorage().Read(ctx, table, &pools)
	if err != nil {
		return
	}

	return
}

func (pr *PostgresRepo) RemovePool(
	ctx c.Context, table string, pool entities.Pool,
) (
	err error,
) {
	err = pr.GetStorage().Remove(ctx, table, &pool)

	return
}

func (pr *PostgresRepo) RemovePools(
	ctx c.Context, table string, pools []entities.Pool,
) (
	out []entities.Pool,
	err error,
) {
	for _, pool := range pools {
		err = pr.RemovePool(ctx, table, pool)
		if err != nil {
			return
		}
	}

	out, err = pr.ListPools(ctx, table)

	return
}

func (pr *PostgresRepo) StoreTokens(
	ctx c.Context, table string, tokens []entities.Token,
) (
	err error,
) {
	for _, token := range tokens {
		err = pr.AddToken(ctx, table, token)
		if err != nil {
			return
		}
	}

	return
}

func (pr *PostgresRepo) AddToken(
	ctx c.Context, table string, token entities.Token,
) (
	err error,
) {

	err = pr.GetStorage().Store(ctx, table, &token)

	return
}

func (pr *PostgresRepo) ListTokens(
	ctx c.Context, table string,
) (
	tokens []entities.Token,
	err error,
) {
	err = pr.GetStorage().Read(ctx, table, &tokens)
	if err != nil {
		return
	}

	return
}

func (pr *PostgresRepo) GetTokenByAddress(
	ctx c.Context, table string, address string,
) (
	token entities.Token,
	err error,
) {
	tokens, err := pr.ListTokens(ctx, table)
	if err != nil {
		return
	}

	for _, t := range tokens {
		if t.Address == address {
			token = t
			return
		}
	}

	return
}

func (pr *PostgresRepo) RemoveToken(
	ctx c.Context, table string, token entities.Token,
) (
	err error,
) {

	err = pr.GetStorage().Remove(ctx, table, token)

	return
}

func (pr *PostgresRepo) RemoveTokens(
	ctx c.Context, table string, tokens []entities.Token,
) (
	out []entities.Token,
	err error,
) {
	for _, token := range tokens {
		err = pr.RemoveToken(ctx, table, token)
		if err != nil {
			return
		}
	}

	out, err = pr.ListTokens(ctx, table)

	return
}
