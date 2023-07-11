package v1

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	tc "github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"

	eth "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
	log "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

type tradecaseRoutes struct {
	t tc.TradeCase
	l log.Interface
}

// @Summary     Get Pairs
// @Description Get list of pairs by protocol, tokens
// @ID          getPairs
// @Tags  	    Contract
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

// @Summary     List Pairs
// @Description Get full list of pairs
// @ID          listPairs
// @Tags  	    Contract
// @Accept      json
// @Produce     json
// @Success     200 {object} listPairs
// @Failure     400 {object} responseErr
// @Failure     404 {object} responseErr
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

// @Summary     Add Pairs
// @Description Add list of pairs
// @ID          addPairs
// @Tags  	    Contract
// @Accept      json
// @Produce     json
// @Param       request body listPairs true "Add pairs"
// @Success     201 {object} listPairs
// @Failure     400 {object} responseErr
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

	respondCreated(c, req)
}

// @Summary     Get Tokens
// @Description Request base token list
// @ID          getTokens
// @Tags  	    Contract
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Failure     507 {object} responseErr
// @Router      /contract/tokens/base [get]
func (tr *tradecaseRoutes) GetBaseTokens(
	c *gin.Context,
) {
	tokens, err := tr.t.Contract.API().GetBaseTokens(
		eth.CallOpts(),
	)
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - GetTokens",
			),
		)
		return
	}

	respondOk(c, tokens)
}

// @Summary     Add Tokens
// @Description Add list of tokens
// @ID          addTokens
// @Tags  	    Provider
// @Accept      json
// @Produce     json
// @Param       request body listTokens true "Add tokens"
// @Success     201 {object} listTokens
// @Failure     400 {object} responseErr
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
// @Tags  	    Provider
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

// @Summary     Load Tokens
// @Description Load unknown tokens from storage
// @ID          loadTokens
// @Tags  	    Trade
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Failure     409 {object} responseErr
// @Router      /trade/tokens [get]
func (tr *tradecaseRoutes) LoadTokens(
	c *gin.Context,
) {
	err := tr.t.SetTokens(c, "tokens")
	if err != nil {
		errorConflict(
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
// @Description Load profitable pairs from storage
// @ID          loadPairs
// @Tags  	    Trade
// @Accept      json
// @Produce     json
// @Success     200 {object} listPairs
// @Failure     507 {object} responseErr
// @Router      /trade/pairs [get]
func (tr *tradecaseRoutes) LoadPairs(
	c *gin.Context,
) {
	err := tr.t.SetProfitablePairs(c, "pools")
	if err != nil {
		errorInufficientStorage(
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

// @Summary     Withdraw
// @Description Withdraw tokens from contract
// @ID          withdraw
// @Tags  	    Trade
// @Accept      json
// @Produce     json
// @Success     202 {object} response
// @Failure     507 {object} responseErr
// @Router      /trade/withdraw [get]
func (tr *tradecaseRoutes) Withdraw(
	c *gin.Context,
) {

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	tx, err := tr.t.Withdraw(ctx)
	if err != nil {
		errorInufficientStorage(
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

// @Summary     Add Base Token
// @Description Add base token to contract
// @ID          addBase
// @Tags  	    Trade
// @Accept      json
// @Produce     json
// @Param		token query string true "Add base token"
// @Success     201 {object} response
// @Failure     507 {object} responseErr
// @Router      /trade/token [post]
func (tr *tradecaseRoutes) AddBase(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	addr := c.Query("token")

	tx, err := tr.t.AddBaseToken(
		ctx,
		addr,
	)
	if err != nil {
		errorConflict(
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
// @Tags  	    Trade
// @Accept      json
// @Produce     json
// @Param		token query string false "Remove base token"
// @Success     202 {object} response
// @Failure     507 {object} responseErr
// @Router      /trade/token [delete]
func (tr *tradecaseRoutes) RmBase(
	c *gin.Context,
) {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	addr := c.Query("token")

	fmt.Println(addr)

	tx, err := tr.t.RmBaseToken(
		ctx,
		addr,
	)
	if err != nil {
		errorConflict(
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

func NewTradecaseRouter(
	h *gin.RouterGroup,
	t tc.TradeCase,
	l log.Interface,
) {
	routes := &tradecaseRoutes{t, l}

	NewContractRouter(h, *routes)
	NewProviderRouter(h, *routes)
	NewTradeRouter(h, *routes)
}

func NewContractRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("provider")
	{
		handler.GET(
			"/tokens",
			tr.ListTokens,
		)
		handler.POST(
			"/tokens",
			tr.AddTokens,
		)
	}
}

func NewProviderRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("contract")
	{
		handler.POST(
			"/pairs/find",
			tr.GetPairs,
		)
		handler.GET(
			"/pairs",
			tr.ListPairs,
		)
		handler.POST(
			"/pairs",
			tr.AddPairs,
		)
		handler.GET(
			"/tokens/base",
			tr.GetBaseTokens,
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
			"/tokens",
			tr.LoadTokens,
		)
		handler.GET(
			"/pairs",
			tr.LoadPairs,
		)
		handler.GET(
			"/withdraw",
			tr.Withdraw,
		)
		handler.POST(
			"/token",
			tr.AddBase,
		)
		handler.DELETE(
			"/token",
			tr.RmBase,
		)
	}
}
