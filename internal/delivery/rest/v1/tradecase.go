package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

type tradecaseRoutes struct {
	t tradecase.TradeCase
	l logger.Interface
}

// @Description Request list of trading pools
type ListPools struct {
	Pools []entities.TradePool `json:"pools" bson:"pools"` // list of trade pools
} //@name ListPools

// @Description Request for searching trade pair
type requestPairs struct {
	Protocol  entities.SwapProtocol `json:"protocol" bson:"protocol"`   // trade protocol
	TokenPair entities.TokenPair    `json:"tokenPair" bson:"tokenPair"` // pair of tokens
} //@name RequestPairs

// @Summary     Get Pools
// @Description Get list of pools
// @ID          getPools
// @Tags  	    Storage, pools
// @Accept      json
// @Produce     json
// @Success     200 {object} ListPools
// @Failure     400 {object} response
// @Failure     409 {object} response
// @Failure     500 {object} response
// @Router      /storage/pools [get]
func (tr *tradecaseRoutes) GetPools(c *gin.Context) {
	res := &ListPools{}

	out, err := tr.t.GetRepo().Read(c)
	if err != nil {
		tr.l.Error(err, "rest - v1 - GetPools")
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	err = json.Unmarshal(out, res)
	if err != nil {
		tr.l.Error(err, "rest - v1 - AddPools")
		errorResponse(c, http.StatusConflict, err.Error())
	}

	c.JSON(http.StatusOK, res)

	return
}

// @Summary     Add Pools
// @Description Add list of pools
// @ID          addPools
// @Tags  	    Storage, pools
// @Accept      json
// @Produce     json
// @Param       request body []entities.TradePool true "Add pools"
// @Success     200 {object} []entities.TradePool
// @Failure     400 {object} response
// @Failure     409 {object} response
// @Failure     507 {object} response
// @Router      /storage/pools [post]
func (tr *tradecaseRoutes) AddPools(c *gin.Context) {
	var body []byte
	pools := &ListPools{}

	err := c.BindJSON(body)
	if err != nil {
		tr.l.Error(err, "rest - v1 - AddPools")
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	err = json.Unmarshal(body, pools)
	if err != nil {
		tr.l.Error(err, "rest - v1 - AddPools")
		errorResponse(c, http.StatusConflict, err.Error())
	}

	err = tr.t.GetRepo().Store(c, body)
	if err != nil {
		tr.l.Error(err, "rest - v1 - Store")
		errorResponse(c, http.StatusInsufficientStorage, err.Error())

		return
	}
	c.JSON(http.StatusOK, body)

	return
}

// @Summary     Get Pairs
// @Description Get list of pairs
// @ID          getPairs
// @Tags  	    Provider, pairs
// @Accept      json
// @Produce     json
// @Param       request body requestPairs true "Get pairs"
// @Success     200 {object} []entities.TradePair
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Router      /provider/pairs [get]
func (tr *tradecaseRoutes) GetPairs(c *gin.Context) {
	var req requestPairs

	err := c.BindJSON(req)
	if err != nil {
		tr.l.Error(err, "rest - v1 - GetPairs")
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	pairs, ok := tr.t.GetProvider().GetPairs(c, req.Protocol, req.TokenPair)
	if !ok {
		tr.l.Error("pairs not found", "rest - v1 - GetPairs")
		errorResponse(c, http.StatusNotFound, "no pairs found")

		return
	}
	c.JSON(http.StatusOK, pairs)

	return
}

// @Summary     Add Pairs
// @Description Add list of pairs
// @ID          addPairs
// @Tags  	    Provider, pairs
// @Accept      json
// @Produce     json
// @Param       request body []entities.TradePair true "Add pairs"
// @Success     200 {object} []entities.TradePair
// @Failure     400 {object} response
// @Router      /provider/pairs [post]
func (tr *tradecaseRoutes) AddPairs(c *gin.Context) {
	var req []entities.TradePair

	err := c.BindJSON(req)
	if err != nil {
		tr.l.Error(err, "rest - v1 - AddPairs")
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	for _, pair := range req {
		err := tr.t.GetProvider().AddPair(c, pair)
		if err != nil {
			tr.l.Error(err, "rest - v1 - AddPairs")
			errorResponse(c, http.StatusBadRequest, err.Error())

			return
		}
	}
	c.JSON(http.StatusOK, req)

	return
}

// @Summary     Add Tokens
// @Description Add list of tokens
// @ID          addTokens
// @Tags  	    Contract, tokens
// @Accept      json
// @Produce     json
// @Param       request body []entities.Token true "Add tokens"
// @Success     200 {object} []entities.Token
// @Failure     400 {object} response
// @Router      /contract/tokens [post]
func (tr *tradecaseRoutes) AddTokens(c *gin.Context) {
	var req []entities.Token

	err := c.BindJSON(req)
	if err != nil {
		tr.l.Error(err, "rest - v1 - AddTokens")
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	for _, token := range req {
		tr.t.GetContract().Add(token)
	}

	c.JSON(http.StatusOK, req)

	return
}

// @Summary     Get Tokens
// @Description Request token list
// @ID          getTokens
// @Tags  	    Contract, tokens
// @Accept      json
// @Produce     json
// @Success     200 {object} []entities.Token
// @Failure     507 {object} response
// @Router      /contract/tokens [get]
func (tr *tradecaseRoutes) GetTokens(c *gin.Context) {
	tokens, err := tr.t.GetContract().GetBaseTokens(c)
	if err != nil {
		tr.l.Error(err, "rest - v1 - GetTokens")
		errorResponse(c, http.StatusInsufficientStorage, err.Error())
	}

	c.JSON(http.StatusOK, tokens)

	return
}

func NewTradecaseRouter(h *gin.RouterGroup, t tradecase.TradeCase, l logger.Interface) {
	routes := &tradecaseRoutes{t, l}

	NewContractRouter(h, *routes)
	NewStorageRouter(h, *routes)
	NewProviderRouter(h, *routes)
}

func NewContractRouter(h *gin.RouterGroup, tr tradecaseRoutes) {
	handler := h.Group("contract")
	{
		handler.GET("/tokens", tr.GetTokens)
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
	handler := h.Group("storage")
	{
		handler.GET("/pairs", tr.GetPairs)
		handler.POST("/pairs", tr.AddPairs)
	}
}
