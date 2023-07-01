package v1

import (
	"encoding/json"

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
		Pools: pr.pc.Pools,
	}
	respondOk(c, pools)

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
func (pr *parsecaseRoutes) GetPools(
	c *gin.Context,
) {
	res := listPools{
		Pools: make([]entities.TradePool, 0),
	}

	out, err := pr.pc.Repo.Read(c, "pools")
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

	err = json.Unmarshal(out, &res.Pools)
	if err != nil {
		errorConflict(
			c, err.Error(),
			Log(
				pr.l.Error,
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

	err = pr.pc.Repo.StorePools(c, "pools", pools.Pools)
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

	return
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
			"/pools",
			pr.GetPools,
		)
		handler.POST(
			"/pools",
			pr.AddPools,
		)
	}
}
