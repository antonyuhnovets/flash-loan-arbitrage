package v1

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	tc "github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	log "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

type tradecaseRoutes struct {
	t tc.TradeCase
	l log.Interface
}

// @Summary     Get Pools
// @Description Get list of pools
// @ID          getPools
// @Tags  	    Storage
// @Accept      json
// @Produce     json
// @Success     200 {object} listPools
// @Failure     409 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/pools [get]
func (tr *tradecaseRoutes) GetPools(
	c *gin.Context,
) {
	res := listPools{
		Pools: make([]entities.TradePool, 0),
	}

	out, err := tr.t.GetRepo().Read(c, "pools")
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - GetPools",
			),
		)
		return
	}

	err = json.Unmarshal(out, &res.Pools)
	if err != nil {
		errorConflict(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - AddPools",
			),
		)
		return
	}
	respondOk(c, res)

	return
}

// @Summary     Add Pools
// @Description Add list of pools
// @ID          addPools
// @Tags  	    Storage
// @Accept      json
// @Produce     json
// @Param       request body listPools true "Add pools"
// @Success     201 {object} listPools
// @Failure     400 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/pools [post]
func (tr *tradecaseRoutes) AddPools(
	c *gin.Context,
) {
	// var body []byte
	pools := &listPools{
		make([]entities.TradePool, 0),
	}

	err := c.BindJSON(pools)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - AddPools",
			),
		)
		return
	}

	err = tr.t.GetRepo().StorePools(c, "pools", pools.Pools)
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				tr.l.Error,
				err,
				"rest - v1 - Store",
			),
		)
		return
	}
	respondCreated(c, pools)

	return
}

// @Summary     Get Pairs
// @Description Get list of pairs
// @ID          getPairs
// @Tags  	    Contract
// @Accept      json
// @Produce     json
// @Param       request body tokenPair true "Get pairs"
// @Success     200 {object} listPairs
// @Failure     400 {object} responseErr
// @Failure     404 {object} responseErr
// @Router      /contract/pairs [get]
func (tr *tradecaseRoutes) GetPairs(
	c *gin.Context,
) {
	var req tokenPair

	err := c.BindJSON(req)
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

	pairs, ok := tr.t.GetContract().GetPairs(
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
	respondOk(c, pairs)

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
// @Router      /contract/pairs/list [get]
func (tr *tradecaseRoutes) ListPairs(
	c *gin.Context,
) {
	pairs := tr.t.GetContract().ListPairs(c)

	respondOk(c, pairs)

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
	var req listPairs

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

	for _, pair := range req.Pairs {
		err := tr.t.GetContract().AddPair(
			c,
			pair,
		)
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
	tokens, err := tr.t.GetContract().GetBaseTokens(c)
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
		err := tr.t.GetProvider().AddToken(c, token)
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
	tokens := tr.t.GetProvider().ListTokens(c)
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
	respondOk(c, tr.t.GetProvider().ListTokens(c))

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
	respondOk(c, tr.t.GetContract().ListPairs(c))

	return
}

// @Summary     Store parsed pools
// @Description Save pools from parser to storage
// @ID          ParseWrite
// @Tags  	    Parse
// @Accept      json
// @Produce     json
// @Success     200 {object} listPools
// @Failure     507 {object} responseErr
// @Router      /parser/pools [get]
func (tr *tradecaseRoutes) ReadParsed(
	c *gin.Context,
) {
	pools := listPools{
		Pools: tr.t.GetParser().ListPools(),
	}
	respondOk(c, pools)

	return
}

func NewTradecaseRouter(
	h *gin.RouterGroup,
	t tc.TradeCase,
	l log.Interface,
) {
	routes := &tradecaseRoutes{t, l}

	NewContractRouter(h, *routes)
	NewStorageRouter(h, *routes)
	NewProviderRouter(h, *routes)
	NewTradeRouter(h, *routes)
	NewParserRouter(h, *routes)
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

func NewStorageRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("storage")
	{
		handler.GET(
			"/pools",
			tr.GetPools,
		)
		handler.POST(
			"/pools",
			tr.AddPools,
		)
	}
}

func NewProviderRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("contract")
	{
		handler.GET(
			"/pairs",
			tr.GetPairs,
		)
		handler.GET(
			"/pairs/list",
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

func NewParserRouter(
	h *gin.RouterGroup,
	tr tradecaseRoutes,
) {
	handler := h.Group("parser")
	{
		handler.GET(
			"/pools",
			tr.ReadParsed,
		)
	}
}
