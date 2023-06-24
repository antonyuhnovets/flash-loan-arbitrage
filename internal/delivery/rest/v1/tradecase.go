package v1

import (
	"encoding/json"
	"net/http"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
	"github.com/gin-gonic/gin"
)

type tradecaseRoutes struct {
	t tradecase.TradeCase
	l logger.Interface
}

type ListPools struct {
	Pools []entities.TradePool `json:"pools"`
}

type requestPairs struct {
	Protocol  entities.SwapProtocol `json:"protocol"`
	TokenPair entities.TokenPair    `json:"tokenPair"`
}

// @Summary     GetPools
// @Description Get list of pools
// @ID          pools
// @Tags  	    getPools
// @Accept      json
// @Produce     json
// @Success     200 {object} ListPools
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /tradecase/pools [get]
func (tr *tradecaseRoutes) GetPools(c *gin.Context) {
	out, err := tr.t.GetRepo().Read(c)
	if err != nil {
		tr.l.Error(err, "rest - v1 - GetPools")
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}
	res := &ListPools{}
	json.Unmarshal(out, res)

	c.JSON(http.StatusOK, res)

	return
}

// @Summary     AddPools
// @Description Add list of pools
// @ID          pools
// @Tags  	    addPools
// @Accept      json
// @Produce     json
// @Param       request body []entities.TradePool true "Add pools"
// @Success     200 {object} []entities.TradePool
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /tradecase/pools [post]
func (tr *tradecaseRoutes) AddPools(c *gin.Context) {
	var body []byte
	pools := &ListPools{}

	err := json.Unmarshal(body, pools)
	if err != nil {
		tr.l.Error(err, "rest - v1 - AddPools")
		errorResponse(c, http.StatusConflict, err.Error())
	}

	err = c.BindJSON(body)
	if err != nil {
		tr.l.Error(err, "rest - v1 - AddPools")
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
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

// @Summary     GetPairs
// @Description Get list of pairs
// @ID          pairs
// @Tags  	    getPairs
// @Accept      json
// @Produce     json
// @Param       request body requestPairs true "Get pairs"
// @Success     200 {object} []entities.TradePair
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /tradecase/pairs [get]
func (tr *tradecaseRoutes) GetPairs(c *gin.Context) {
	var req requestPairs

	c.BindJSON(req)

	pairs, ok := tr.t.GetProvider().GetPairs(c, req.Protocol, req.TokenPair)
	if !ok {
		tr.l.Error("pairs not found", "rest - v1 - GetPairs")
		errorResponse(c, http.StatusNotFound, "no pairs found")

		return
	}
	c.JSON(http.StatusOK, pairs)

	return
}

// @Summary     AddPairs
// @Description Add list of pairs
// @ID          pairs
// @Tags  	    addPairs
// @Accept      json
// @Produce     json
// @Param       request body []entities.TradePair true "Add pairs"
// @Success     200 {object} []entities.TradePair
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /tradecase/pairs [post]
func (tr *tradecaseRoutes) AddPairs(c *gin.Context) {
	var req []entities.TradePair

	c.BindJSON(req)

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

// @Summary     AddTokens
// @Description Add list of tokens
// @ID          tokens
// @Tags  	    addTokens
// @Accept      json
// @Produce     json
// @Param       request body []entities.Token true "Add tokens"
// @Success     200 {object} []entities.Token
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /tradecase/tokens [post]
func (tr *tradecaseRoutes) AddTokens(c *gin.Context) {
	var req []entities.Token

	c.BindJSON(req)

	for _, token := range req {
		tr.t.GetContracr().Add(token)
	}

	c.JSON(http.StatusOK, req)

	return
}

// @Summary     Get Tokens
// @Description Request token list
// @ID          tokens
// @Tags  	    getTokens
// @Accept      json
// @Produce     json
// @Success     200 {object} []entities.Token
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /tradecase/tokens [get]
func (tr *tradecaseRoutes) GetTokens(c *gin.Context) {
	tokens, err := tr.t.GetContracr().GetBaseTokens(c)
	if err != nil {
		tr.l.Error(err, "rest - v1 - GetTokens")
		errorResponse(c, http.StatusBadGateway, err.Error())
	}

	c.JSON(http.StatusOK, tokens)

	return
}

func NewTradecaseRouter(h *gin.RouterGroup, t tradecase.TradeCase, l logger.Interface) {
	routes := &tradecaseRoutes{t, l}

	handler := h.Group("/tradecase")
	{
		handler.GET("/tokens", routes.GetTokens)
		handler.POST("/tokens", routes.AddTokens)
		handler.GET("/pools", routes.GetPools)
		handler.POST("/pools", routes.AddPools)
		handler.GET("/pairs", routes.GetPairs)
		handler.POST("/pairs", routes.AddPairs)
	}
}
