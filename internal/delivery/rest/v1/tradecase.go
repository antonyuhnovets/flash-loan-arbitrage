package v1

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade"

	eth "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
	log "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

type tradecaseRoutes struct {
	t trade.TradeCase
	l log.Interface
}

// @Summary     Add Pairs
// @Description Add list of pairs to contract tradecase memory
// @ID          addPairs
// @Tags  	    Contract: pairs
// @Accept      json
// @Produce     json
// @Param       request body listPairs true "Add pairs"
// @Success     201 {object} listPairs
// @Failure     500 {object} responseErr
// @Router      /contract/pairs [post]
func (tr *tradecaseRoutes) AddPairs(
	c *gin.Context,
) {
	req := listPairs{
		Pairs: make([]entities.TradePair, 0),
	}

	err := c.BindJSON(&req)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - AddPairs",
			),
		)
		return
	}

	err = tr.t.Contract.SetPairs(c, req.Pairs)
	if err != nil {
		errorInternalServer(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - AddPairs",
			),
		)
		return
	}

	respondCreated(c, req)
}

// @Summary     List Pairs
// @Description Get full list of pairs from contract tradecase memory
// @ID          listPairs
// @Tags  	    Contract: pairs
// @Accept      json
// @Produce     json
// @Success     200 {object} listPairs
// @Router      /contract/pairs [get]
func (tr *tradecaseRoutes) ListPairs(
	c *gin.Context,
) {
	pairLst := listPairs{
		Pairs: make([]entities.TradePair, 0),
	}

	pairs := tr.t.Contract.ListPairs(c)
	pairLst.Pairs = append(pairLst.Pairs, pairs...)

	respondOk(c, pairLst)
}

// @Summary     Get Pairs
// @Description Get list of pool pairs from contract tradecase memory by token pair & protocol
// @ID          getPairs
// @Tags  	    Contract: pairs
// @Accept      json
// @Produce     json
// @Param       request body tokenPair true "Get pairs"
// @Success     200 {object} listPairs
// @Failure     400 {object} responseErr
// @Failure     404 {object} responseErr
// @Router      /contract/pairs/find [post]
func (tr *tradecaseRoutes) GetPairs(
	c *gin.Context,
) {
	req := tokenPair{
		Protocol:  entities.SwapProtocol{},
		TokenPair: entities.TokenPair{},
	}

	err := c.Bind(&req)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - GetPairs",
			),
		)
		return
	}

	lst := listPairs{
		Pairs: make([]entities.TradePair, 0),
	}
	pairs, ok := tr.t.Contract.GetPairs(
		c,
		req.Protocol,
		req.TokenPair,
	)
	if !ok {
		errorNotFound(
			c, "no pairs found",
			Log(
				tr.l.Error,
				"pairs not found",
				"rest - v1 - GetPairs",
			),
		)
		return
	}
	lst.Pairs = pairs

	respondOk(c, lst)
}

// @Summary     Add Tokens
// @Description Add list of tokens
// @ID          addTokens
// @Tags  	    Provider: tokens
// @Accept      json
// @Produce     json
// @Param       request body listTokens true "Add tokens"
// @Success     201 {object} listTokens
// @Failure     400 {object} responseErr
// @Failure     409 {object} responseErr
// @Router      /provider/tokens [post]
func (tr *tradecaseRoutes) AddTokens(c *gin.Context) {
	var req listTokens

	err := c.BindJSON(&req)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - AddTokens",
			),
		)
		return
	}

	for _, token := range req.Tokens {
		err := tr.t.Provider.AddToken(c, token)
		if err != nil {
			errorConflict(
				c, err.Error(),
				Log(
					tr.l.Error,
					err,
					"rest - v1 - AddTokens",
				),
			)
			return
		}

	}

	respondCreated(c, req)
}

// @Summary     List Tokens
// @Description Request token list
// @ID          listTokens
// @Tags  	    Provider: tokens
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Router      /provider/tokens [get]
func (tr *tradecaseRoutes) ListTokens(
	c *gin.Context,
) {
	tokens := tr.t.Provider.ListTokens(c)

	respondOk(c, tokens)
}

// @Summary     List base Tokens
// @Description Request base token list from deployed contract memory
// @ID          getTokens
// @Tags  	    Trade: base tokens
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Failure     503 {object} responseErr
// @Router      /contract/tokens/base [get]
func (tr *tradecaseRoutes) GetBaseTokens(
	c *gin.Context,
) {
	tokens, err := tr.t.Contract.Api().Caller().GetBaseTokens(
		eth.CallOpts(),
	)
	if err != nil {
		errorServiceUnavailable(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - GetBaseTokens",
			),
		)
		return
	}

	respondOk(c, tokens)
}

// @Summary     Load Tokens
// @Description Load unknown tokens from storage
// @ID          loadTokens
// @Tags  	    Trade: setup case
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Failure     507 {object} responseErr
// @Router      /trade/tokens [get]
func (tr *tradecaseRoutes) LoadTokens(
	c *gin.Context,
) {
	err := tr.t.SetTokens(c, "tokens")
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - LoadTokens",
			),
		)
		return
	}

	respondOk(c, tr.t.Provider.ListTokens(c))
}

// @Summary     Load Pairs
// @Description Load profitable pool pairs from storage
// @ID          loadPairs
// @Tags  	    Trade: setup case
// @Accept      json
// @Produce     json
// @Success     200 {object} listPairs
// @Failure     503 {object} responseErr
// @Router      /trade/pairs [get]
func (tr *tradecaseRoutes) LoadPairs(
	c *gin.Context,
) {
	err := tr.t.SetProfitablePairs(c, "pools")
	if err != nil {
		errorServiceUnavailable(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - LoadPairs",
			),
		)
		return
	}

	respondOk(c, tr.t.Contract.ListPairs(c))
}

// @Summary     Add Base Token
// @Description Add base token to contract
// @ID          addBase
// @Tags  	    Trade: base tokens
// @Accept      json
// @Produce     json
// @Param		token query string true "Add base token"
// @Success     201 {object} response
// @Failure     503 {object} responseErr
// @Router      /trade/tokens/base [post]
func (tr *tradecaseRoutes) AddBase(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr := c.Query("token")

	tx, err := tr.t.AddBaseToken(
		ctx,
		addr,
	)
	if err != nil {
		errorServiceUnavailable(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - AddBase",
			),
		)
	}

	res := response{tx}

	respondCreated(c, res)
}

// @Summary     Remove Base Token
// @Description Remove base token from contract
// @ID          rmBase
// @Tags  	    Trade: base tokens
// @Accept      json
// @Produce     json
// @Param		token query string false "Remove base token"
// @Success     202 {object} response
// @Failure     503 {object} responseErr
// @Router      /trade/tokens/base [delete]
func (tr *tradecaseRoutes) RmBase(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr := c.Query("token")

	tx, err := tr.t.RmBaseToken(
		ctx,
		addr,
	)
	if err != nil {
		errorServiceUnavailable(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - RmBase",
			),
		)
	}
	res := response{tx}

	respondAccepted(c, res)
}

// @Summary     Withdraw
// @Description Withdraw tokens from contract
// @ID          withdraw
// @Tags  	    Trade: core
// @Accept      json
// @Produce     json
// @Success     202 {object} response
// @Failure     503 {object} responseErr
// @Router      /trade/core/withdraw [get]
func (tr *tradecaseRoutes) Withdraw(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	tx, err := tr.t.Withdraw(ctx)
	if err != nil {
		errorServiceUnavailable(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - Withdraw",
			),
		)
	}

	res := response{tx}

	respondAccepted(c, res)
}

// @Summary     CheckProfit
// @Description Find out if trade with given pools is profitable
// @ID          checkProfit
// @Tags  	    Trade: core
// @Accept      json
// @Produce     json
// @Param		pool0 query string true "Swap pool 0"
// @Param		pool1 query string true "Swap pool 1"
// @Success     202 {object} response
// @Failure     503 {object} responseErr
// @Router      /trade/core/profit-check [get]
func (tr *tradecaseRoutes) CheckProfit(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	pool0 := c.Query("pool0")
	pool1 := c.Query("pool1")

	profit, baseToken, err := tr.t.GetProfit(ctx, pool0, pool1)
	if err != nil {
		errorServiceUnavailable(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - Withdraw",
			),
		)
	}
	var res response
	if profit > 0 {
		res = response{map[string]string{
			"profit":     fmt.Sprintln(profit),
			"base token": baseToken,
		}}
	} else {
		res = response{false}
	}

	respondAccepted(c, res)
}

// @Summary     DoArbitrage
// @Description Call flash arbitrage func from contract with give pools
// @ID          doArbitrage
// @Tags  	    Trade: core
// @Accept      json
// @Produce     json
// @Param		pool0 query string true "Swap pool 0"
// @Param		pool1 query string true "Swap pool 1"
// @Success     202 {object} response
// @Failure     502 {object} responseErr
// @Router      /trade/core/flash-arbitrage [get]
func (tr *tradecaseRoutes) DoArbitrage(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	pool0 := c.Query("pool0")
	pool1 := c.Query("pool1")

	tx, err := tr.t.Arbitrage(ctx, pool0, pool1)
	if err != nil {
		errorBadGateway(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - DoArbitrage",
			),
		)
	}
	res := response{tx}
	respondAccepted(c, res)
}

// @Summary     Replace Tx with add base
// @Description Replace transaction by hash with add base token tx
// @ID          replaceTxAdd
// @Tags  	    Trade: core
// @Accept      json
// @Produce     json
// @Param		hash query string true "Tx hash"
// @Param		token query string true "Token"
// @Success     202 {object} response
// @Failure     502 {object} responseErr
// @Router      /replace-tx-add [post]
func (tr *tradecaseRoutes) ReplaceTxAddBase(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	addr := c.Query("hash")
	tx, err := tr.t.ReplaceTxWithAddBaseToken(ctx, addr, c.Query("token"))
	if err != nil {
		errorBadGateway(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - ReplaceTxAddBase",
			),
		)
	}

	res := response{tx}

	respondAccepted(c, res)
}

// @Summary     Replace Tx with rm base token
// @Description Replace transaction by hash with remove base token tx
// @ID          replaceTxRm
// @Tags  	    Trade: core
// @Accept      json
// @Produce     json
// @Param		hash query string true "Tx hash"
// @Param		token query string true "Token"
// @Success     202 {object} response
// @Failure     502 {object} responseErr
// @Router      /replace-tx-rm [delete]
func (tr *tradecaseRoutes) ReplaceTxRmBase(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	addr := c.Query("hash")
	tx, err := tr.t.ReplaceTxWithAddBaseToken(ctx, addr, c.Query("token"))
	if err != nil {
		errorBadGateway(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - ReplaceTxRmBase",
			),
		)
	}

	res := response{tx}

	respondAccepted(c, res)
}

func NewTradecaseRouter(
	h *gin.RouterGroup,
	t trade.TradeCase,
	l log.Interface,
) {
	routes := &tradecaseRoutes{t, l}

	NewTradeRouter(h, *routes)
	NewProviderRouter(h, *routes)
	NewContractRouter(h, *routes)
}

func NewContractRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("contract")
	{
		handler.POST(
			"/pairs",
			tr.AddPairs,
		)
		handler.GET(
			"/pairs",
			tr.ListPairs,
		)
		handler.POST(
			"/pairs/find",
			tr.GetPairs,
		)
		handler.GET(
			"/tokens/base",
			tr.GetBaseTokens,
		)
	}
}

func NewProviderRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("provider")
	{
		handler.POST(
			"/tokens",
			tr.AddTokens,
		)
		handler.GET(
			"/tokens",
			tr.ListTokens,
		)
	}
}

func NewTradeRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("trade")
	{
		handler.GET(
			"/core/withdraw",
			tr.Withdraw,
		)
		handler.GET(
			"/core/profit-check",
			tr.CheckProfit,
		)
		handler.GET(
			"/core/flash-arbitrage",
			tr.DoArbitrage,
		)
		handler.POST(
			"/replace-tx-add",
			tr.ReplaceTxAddBase,
		)
		handler.DELETE(
			"/replace-tx-add",
			tr.ReplaceTxRmBase,
		)
		handler.POST(
			"/tokens/base",
			tr.AddBase,
		)
		handler.DELETE(
			"/tokens/base",
			tr.RmBase,
		)
		handler.GET(
			"/tokens",
			tr.LoadTokens,
		)
		handler.GET(
			"/pairs",
			tr.LoadPairs,
		)
	}
}
