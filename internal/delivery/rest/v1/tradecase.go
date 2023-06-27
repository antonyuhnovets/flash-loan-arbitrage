package v1

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

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
func (tr *tradecaseRoutes) GetPools(c *gin.Context) {
	res := &listPools{}

	out, err := tr.t.GetRepo().Read(c)
	if err != nil {
		errorInufficientStorage(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - GetPools"),
		)
		return
	}

	err = json.Unmarshal(out, res)
	if err != nil {
		errorConflict(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - AddPools"),
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
// @Failure     409 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/pools [post]
func (tr *tradecaseRoutes) AddPools(c *gin.Context) {
	var body []byte
	pools := &listPools{}

	err := c.BindJSON(body)
	if err != nil {
		errorBadRequest(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - AddPools"),
		)
		return
	}

	err = json.Unmarshal(body, pools)
	if err != nil {
		errorConflict(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - AddPools"),
		)
		return
	}

	err = tr.t.GetRepo().Store(c, body)
	if err != nil {
		errorInufficientStorage(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - Store"),
		)
		return
	}
	respondCreated(c, body)

	return
}

// @Summary     Get Pairs
// @Description Get list of pairs
// @ID          getPairs
// @Tags  	    Provider
// @Accept      json
// @Produce     json
// @Param       request body tokenPair true "Get pairs"
// @Success     200 {object} listPairs
// @Failure     400 {object} responseErr
// @Failure     404 {object} responseErr
// @Router      /provider/pairs [get]
func (tr *tradecaseRoutes) GetPairs(c *gin.Context) {
	var req tokenPair

	err := c.BindJSON(req)
	if err != nil {
		errorBadRequest(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - GetPairs"),
		)
		return
	}

	pairs, ok := tr.t.GetProvider().GetPairs(c, req.Protocol, req.TokenPair)
	if !ok {
		errorNotFound(c, "no pairs found",
			Log(tr.l.Error, "pairs not found", "rest - v1 - GetPairs"),
		)
		return
	}
	respondOk(c, pairs)

	return
}

// @Summary     List Pairs
// @Description Get full list of pairs
// @ID          listPairs
// @Tags  	    Provider
// @Accept      json
// @Produce     json
// @Success     200 {object} listPairs
// @Failure     400 {object} responseErr
// @Failure     404 {object} responseErr
// @Router      /provider/pairs/list [get]
func (tr *tradecaseRoutes) ListPairs(c *gin.Context) {

	pairs := tr.t.GetProvider().ListPairs(c)

	respondOk(c, pairs)

	return
}

// @Summary     Add Pairs
// @Description Add list of pairs
// @ID          addPairs
// @Tags  	    Provider
// @Accept      json
// @Produce     json
// @Param       request body listPairs true "Add pairs"
// @Success     201 {object} listPairs
// @Failure     400 {object} responseErr
// @Router      /provider/pairs [post]
func (tr *tradecaseRoutes) AddPairs(c *gin.Context) {
	var req listPairs

	err := c.BindJSON(&req)
	if err != nil {
		errorBadRequest(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - AddPairs"),
		)
		return
	}

	for _, pair := range req.Pairs {
		err := tr.t.GetProvider().AddPair(c, pair)
		if err != nil {
			errorBadRequest(c, err.Error(),
				Log(tr.l.Error, err, "rest - v1 - AddPairs"),
			)
			return
		}
	}
	respondCreated(c, req)

	return
}

// @Summary     Add Tokens
// @Description Add list of tokens
// @ID          addTokens
// @Tags  	    Contract
// @Accept      json
// @Produce     json
// @Param       request body listTokens true "Add tokens"
// @Success     201 {object} listTokens
// @Failure     400 {object} responseErr
// @Router      /contract/tokens [post]
func (tr *tradecaseRoutes) AddTokens(c *gin.Context) {
	var req listTokens

	err := c.BindJSON(&req)
	if err != nil {
		errorBadRequest(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - AddTokens"),
		)
		return
	}

	for _, token := range req.Tokens {
		tr.t.GetContract().Add(token)
	}
	respondCreated(c, req)

	return
}

// @Summary     List Tokens
// @Description Request token list
// @ID          listTokens
// @Tags  	    Contract
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Router      /contract/tokens [get]
func (tr *tradecaseRoutes) ListTokens(c *gin.Context) {
	tokens := tr.t.GetContract().List()
	respondOk(c, tokens)

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
func (tr *tradecaseRoutes) GetBaseTokens(c *gin.Context) {
	tokens, err := tr.t.GetContract().GetBaseTokens(c)
	if err != nil {
		errorInufficientStorage(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - GetTokens"),
		)
		return
	}
	respondOk(c, tokens)

	return
}

// @Summary     Load Tokens
// @Description Load unknown tokens
// @ID          loadTokens
// @Tags  	    Trade
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Failure     507 {object} responseErr
// @Router      /trade/tokens [get]
func (tr *tradecaseRoutes) LoadTokens(c *gin.Context) {
	err := tr.t.SetUnknownTokens(c)
	if err != nil {
		errorInufficientStorage(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - GetTokens"),
		)
		return
	}
	respondOk(c, tr.t.GetContract().List())

	return
}

// @Summary     Load Pairs
// @Description Load profitable pairs
// @ID          loadPairs
// @Tags  	    Trade
// @Accept      json
// @Produce     json
// @Success     200 {object} listPairs
// @Failure     507 {object} responseErr
// @Router      /trade/pairs [get]
func (tr *tradecaseRoutes) LoadPairs(c *gin.Context) {
	err := tr.t.SetProfitablePairs(c)
	if err != nil {
		errorInufficientStorage(c, err.Error(),
			Log(tr.l.Error, err, "rest - v1 - GetTokens"),
		)
		return
	}
	respondOk(c, tr.t.GetProvider().ListPairs(c))

	return
}

func NewTradecaseRouter(h *gin.RouterGroup, t tc.TradeCase, l log.Interface) {
	routes := &tradecaseRoutes{t, l}

	NewContractRouter(h, *routes)
	NewStorageRouter(h, *routes)
	NewProviderRouter(h, *routes)
	NewTradeRouter(h, *routes)
}

func NewContractRouter(h *gin.RouterGroup, tr tradecaseRoutes) {
	handler := h.Group("contract")
	{
		handler.GET("/tokens/base", tr.GetBaseTokens)
		handler.GET("/tokens", tr.ListTokens)
		handler.POST("/tokens", tr.AddTokens)
	}
}

func NewStorageRouter(h *gin.RouterGroup, tr tradecaseRoutes) {
	handler := h.Group("storage")
	{
		handler.GET("/pools", tr.GetPools)
		handler.POST("/pools", tr.AddPools)
	}
}

func NewProviderRouter(h *gin.RouterGroup, tr tradecaseRoutes) {
	handler := h.Group("provider")
	{
		handler.GET("/pairs", tr.GetPairs)
		handler.GET("/pairs/list", tr.ListPairs)
		handler.POST("/pairs", tr.AddPairs)
	}
}

func NewTradeRouter(h *gin.RouterGroup, tr tradecaseRoutes) {
	handler := h.Group("trade")
	{
		handler.GET("/tokens", tr.LoadTokens)
		handler.GET("/pairs", tr.LoadPairs)
	}
}
