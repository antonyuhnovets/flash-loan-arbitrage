package tradecase

import (
	"context"
)

type TradeCase struct {
	repo     TradeRepo
	provider TradeProvider
	contract SmartContract
}

func New(r TradeRepo, p TradeProvider, c SmartContract) *TradeCase {
	return &TradeCase{
		repo:     r,
		provider: p,
	}
}

func (tc *TradeCase) Trade(ctx context.Context) error {

}
