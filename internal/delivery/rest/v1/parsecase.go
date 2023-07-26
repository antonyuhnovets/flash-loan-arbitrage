package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade"
	log "github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

type parsecaseRoutes struct {
	pc trade.ParseCase
	l  log.Interface
}

// @Summary     Add Tokens
// @Description Add list of tokens to storage
// @ID          storeTokens
// @Tags  	    Storage: tokens
// @Accept      json
// @Produce     json
// @Param       request body listTokens true "Add tokens"
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
				"rest - v1 - AddTokens",
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
				"rest - v1 - AddTokens",
			),
		)
		return
	}

	respondCreated(c, tokens)
}

// @Summary     Get Tokens
// @Description Get list of all tokens from storage
// @ID          getTokensList
// @Tags  	    Storage: tokens
// @Accept      json
// @Produce     json
// @Success     200 {object} listTokens
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
				"rest - v1 - ListTokens",
			),
		)
		return
	}

	res.Tokens = out

	respondOk(c, res)
}

// @Summary     Delete Tokens
// @Description Delete tokens from storage
// @ID          deleteTokens
// @Tags  	    Storage: tokens
// @Accept      json
// @Produce     json
// @Param       request body listTokens true "Delete tokens"
// @Success     200 {object} listTokens
// @Failure     400 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/tokens [delete]
func (pr *parsecaseRoutes) DeleteTokens(
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
				"rest - v1 - DeleteTokens",
			),
		)
		return
	}

	out, err := pr.pc.Repository.RemoveTokens(
		c, "tokens", tokens.Tokens,
	)

	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - DeleteTokens",
			),
		)
		return
	}

	res := &listTokens{Tokens: out}

	respondOk(c, res)
}

// @Summary     Add Pools
// @Description Add list of pools to storage
// @ID          addPools
// @Tags  	    Storage: pools
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
		Pools: make([]entities.Pool, 0),
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
				"rest - v1 - AddPools",
			),
		)
		return
	}

	respondCreated(c, pools)
}

// @Summary     Get Pools
// @Description Get list of all pools from storage
// @ID          getPoolList
// @Tags  	    Storage: pools
// @Accept      json
// @Produce     json
// @Success     200 {object} listPools
// @Failure     507 {object} responseErr
// @Router      /storage/pools [get]
func (pr *parsecaseRoutes) ListPools(
	c *gin.Context,
) {

	out, err := pr.pc.Repository.ListPools(c, "pools")
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - ListPools",
			),
		)
		return
	}

	respondOk(c, listPools{Pools: out})
}

// @Summary     Delete Pools
// @Description Delete pools from storage
// @ID          deletePools
// @Tags  	    Storage: pools
// @Accept      json
// @Produce     json
// @Param       request body listPools true "Delete pools"
// @Success     200 {object} listPools
// @Failure     400 {object} responseErr
// @Failure     507 {object} responseErr
// @Router      /storage/pools [delete]
func (pr *parsecaseRoutes) DeletePools(
	c *gin.Context,
) {
	pools := &listPools{
		Pools: make([]entities.Pool, 0),
	}

	err := c.BindJSON(pools)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - DeletePools",
			),
		)
		return
	}

	out, err := pr.pc.Repository.RemovePools(
		c, "pools", pools.Pools,
	)

	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - DeletePools",
			),
		)
		return
	}

	res := &listPools{Pools: out}

	respondOk(c, res)
}

// @Summary     Get parsed pools
// @Description Get pools from current parser
// @ID          getParsed
// @Tags  	    Parse: setup parser
// @Accept      json
// @Produce     json
// @Success     200 {object} listPools
// @Failure     500 {object} responseErr
// @Router      /parser/pools [get]
func (pr *parsecaseRoutes) ReadParsed(
	c *gin.Context,
) {
	pools := listPools{
		Pools: pr.pc.Parser.ListPools(),
	}

	respondOk(c, pools)
}

// @Summary     Store parsed pools
// @Description Save pools from parser to storage
// @ID          ParseWrite
// @Tags  	    Parse: core
// @Accept      json
// @Produce     json
// @Success     200 {object} listPools
// @Failure     507 {object} responseErr
// @Router      /parser/core/parse-save [get]
func (pr *parsecaseRoutes) StoreParsed(
	c *gin.Context,
) {
	err := pr.pc.ParseAndStore(c)
	if err != nil {
		errorInufficientStorage(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - StoreParsed",
			),
		)
	}

	pools := listPools{
		Pools: pr.pc.Parser.ListPools(),
	}

	respondOk(c, pools)
}

// @Summary     Parse
// @Description Parse and store to parser local memory
// @ID          parse
// @Tags  	    Parse: core
// @Accept      json
// @Produce     json
// @Success     200 {object} listPools
// @Failure     500 {object} responseErr
// @Router      /parser/core/parse [get]
func (pr *parsecaseRoutes) Parse(
	c *gin.Context,
) {
	pools, err := pr.pc.JustParse(c)
	if err != nil {
		errorInternalServer(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - Parse",
			),
		)
	}

	respondOk(c, listPools{pools})
}

// @Summary     Add protocols
// @Description Add parse protocols
// @ID          addProtocols
// @Tags  	    Parse: setup parser
// @Accept      json
// @Produce     json
// @Param       request body listProtocols true "Set protocol"
// @Success     200 {object} listProtocols
// @Failure     400 {object} responseErr
// @Failure     502 {object} responseErr
// @Router      /parser/protocols [post]
func (pr *parsecaseRoutes) AddProtocols(
	c *gin.Context,
) {
	protocols := listProtocols{}

	err := c.Bind(&protocols)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - AddProtocol",
			),
		)
	}

	for _, p := range protocols.Protocols {
		err = pr.pc.AddProtocol(context.Background(), p)
		if err != nil {
			errorBadGateway(
				c, err.Error(),
				Log(
					pr.l.Error,
					err,
					"rest - v1 - AddProtocol",
				),
			)
		}
	}

	respondOk(c, protocols)

}

// @Summary     Get protocols
// @Description Get current parse protocols
// @ID          getProtocols
// @Tags  	    Parse: setup parser
// @Accept      json
// @Produce     json
// @Success     200 {object} listProtocols
// @Failure     500 {object} responseErr
// @Router      /parser/protocols [get]
func (pr *parsecaseRoutes) GetProtocols(
	c *gin.Context,
) {
	p := pr.pc.ListProtocols()

	respondOk(c, p)
}

// @Summary     Remove protocols
// @Description Remove parse protocols
// @ID          rmProtocols
// @Tags  	    Parse: setup parser
// @Accept      json
// @Produce     json
// @Param       request body listProtocols true "Remove protocols"
// @Success     202 {object} listProtocols
// @Failure     400 {object} responseErr
// @Failure     404 {object} responseErr
// @Router      /parser/protocols [delete]
func (pr *parsecaseRoutes) RmProtocols(
	c *gin.Context,
) {
	l := listProtocols{}

	err := c.Bind(&l)
	if err != nil {
		errorBadRequest(
			c, err.Error(),
			Log(
				pr.l.Error,
				err,
				"rest - v1 - RmProtocol",
			),
		)
	}

	for _, proto := range l.Protocols {
		err = pr.pc.RmProtocol(context.Background(), proto)
		if err != nil {
			errorNotFound(
				c, err.Error(),
				Log(
					pr.l.Error,
					err,
					"rest - v1 - AddProtocol",
				),
			)
		}
	}

	respondAccepted(c, l)
}

func NewParsecaseRouter(
	h *gin.RouterGroup,
	t trade.ParseCase,
	l log.Interface,
) {
	routes := &parsecaseRoutes{t, l}

	NewStorageRouter(h, *routes)
	NewParserRouter(h, *routes)
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
		handler.DELETE(
			"/tokens",
			pr.DeleteTokens,
		)
		handler.GET(
			"/pools",
			pr.ListPools,
		)
		handler.POST(
			"/pools",
			pr.AddPools,
		)
		handler.DELETE(
			"/pools",
			pr.DeletePools,
		)
	}
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
		handler.GET(
			"/core/parse-save",
			pr.StoreParsed,
		)
		handler.GET(
			"/core/parse",
			pr.Parse,
		)
		handler.POST(
			"/protocols",
			pr.AddProtocols,
		)
		handler.GET(
			"/protocols",
			pr.GetProtocols,
		)
		handler.DELETE(
			"/protocols",
			pr.RmProtocols,
		)
	}
}
