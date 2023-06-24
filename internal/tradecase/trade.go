package tradecase

import (
	"context"
	"encoding/json"
	"log"
)

type TradeCase struct {
	repo     TradeRepo
	provider TradeProvider
	contract SmartContract
}

func New(r TradeRepo, p TradeProvider, c SmartContract,
) *TradeCase {

	return &TradeCase{
		repo:     r,
		provider: p,
		contract: c,
	}
}

func (tc *TradeCase) GetRepo() TradeRepo {
	return tc.repo
}

func (tc *TradeCase) GetContracr() SmartContract {
	return tc.contract
}

func (tc *TradeCase) GetProvider() TradeProvider {
	return tc.provider
}

func (tc *TradeCase) Trade(ctx context.Context,
) error {
	tokens, err := tc.contract.GetBaseTokens(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	// err = tc.repo.Clear(ctx)
	// if err != nil {
	// 	return err
	// }
	// b = append(b)
	err = tc.repo.Store(ctx, b)
	log.Println("Storing...")
	if err != nil {
		return err
	}
	err = tc.repo.Store(ctx, []byte("\n"))
	log.Println("Stored")

	out, err := tc.repo.Read(ctx)
	log.Println("Reading...")
	if err != nil {
		return err
	}
	log.Println("Readed")

	log.Println(string(out))

	return nil
}
