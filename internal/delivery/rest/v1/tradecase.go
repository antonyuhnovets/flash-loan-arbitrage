package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	tc "github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
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

	return
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

	return
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

	return
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
	tokens, err := tr.t.Contract.GetBaseTokens(c)
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

	return
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

	return
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

	return
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
				"rest - v1 - GetTokens",
			),
		)
		return
	}
	respondOk(c, tr.t.Provider.ListTokens(c))

	return
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
				"rest - v1 - GetTokens",
			),
		)
		return
	}
	respondOk(c, tr.t.Contract.ListPairs(c))

	return
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
	}
}
