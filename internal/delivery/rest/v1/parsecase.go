package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	tc "github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	log "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

type parsecaseRoutes struct {
	pc tc.ParseCase
	l  log.Interface
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
func (pr *parsecaseRoutes) ReadParsed(
	c *gin.Context,
) {
	err := pr.pc.ParseAndStore(c)
	if err != nil {
		errorConflict(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - ReadParsed",
			),
		)
	}

	pools := listPools{
		Pools: pr.pc.Parser.ListPools(),
	}
	respondOk(c, pools)

}

// @Summary     Get Tokens
// @Description Get list of all tokens from storage
// @ID          getTokensList
// @Tags  	    Storage
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
// @Failure     409 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/tokens [get]
func (pr *parsecaseRoutes) ListTokens(
	c *gin.Context,
) {
	res := listTokens{
		Tokens: make([]entities.Token, 0),
	}

	out, err := pr.pc.Repository.ListTokens(c, "tokens")
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - GetPools",
			),
		)
		return
	}

	res.Tokens = out

	respondOk(c, res)
}

// @Summary     Add Tokens
// @Description Add list of tokens to storage
// @ID          addTokens
// @Tags  	    Storage
// @Accept      json
// @Produce     json
// @Param       request body listTokens true "Add pools"
// @Success     201 {object} listTokens
// @Failure     400 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/tokens [post]
func (pr *parsecaseRoutes) AddTokens(
	c *gin.Context,
) {
	tokens := &listTokens{
		Tokens: make([]entities.Token, 0),
	}

	err := c.BindJSON(tokens)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - AddPools",
			),
		)
		return
	}

	err = pr.pc.Repository.StoreTokens(c, "tokens", tokens.Tokens)
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - Store",
			),
		)
		return
	}

	respondCreated(c, tokens)
}

// @Summary     Get Pools
// @Description Get list of all pools from storage
// @ID          getPoolList
// @Tags  	    Storage
// @Accept      json
// @Produce     json
// @Success     200 {object} listPools
// @Failure     409 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/pools [get]
func (pr *parsecaseRoutes) ListPools(
	c *gin.Context,
) {
	res := listPools{
		Pools: make([]entities.TradePool, 0),
	}

	out, err := pr.pc.Repository.ListPools(c, "pools")
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - GetPools",
			),
		)
		return
	}
	res.Pools = out

	respondOk(c, res)
}

// @Summary     Add Pools
// @Description Add list of pools to storage
// @ID          addPools
// @Tags  	    Storage
// @Accept      json
// @Produce     json
// @Param       request body listPools true "Add pools"
// @Success     201 {object} listPools
// @Failure     400 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/pools [post]
func (pr *parsecaseRoutes) AddPools(
	c *gin.Context,
) {
	pools := &listPools{
		make([]entities.TradePool, 0),
	}

	err := c.BindJSON(pools)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - AddPools",
			),
		)
		return
	}

	err = pr.pc.Repository.StorePools(c, "pools", pools.Pools)
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - Store",
			),
		)
		return
	}

	respondCreated(c, pools)
}

func NewParsecaseRouter(
	h *gin.RouterGroup,
	t tc.ParseCase,
	l log.Interface,
) {
	routes := &parsecaseRoutes{t, l}

	NewStorageRouter(h, *routes)
	NewParserRouter(h, *routes)
}

func NewParserRouter(
	h *gin.RouterGroup,
	pr parsecaseRoutes,
) {
	handler := h.Group("parser")
	{
		handler.GET(
			"/pools",
			pr.ReadParsed,
		)
	}
}

func NewStorageRouter(
	h *gin.RouterGroup,
	pr parsecaseRoutes,
) {
	handler := h.Group("storage")
	{
		handler.GET(
			"/tokens",
			pr.ListTokens,
		)
		handler.POST(
			"/tokens",
			pr.AddTokens,
		)
		handler.GET(
			"/pools",
			pr.ListPools,
		)
		handler.POST(
			"/pools",
			pr.AddPools,
		)
	}
}
