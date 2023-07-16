package v1

import (
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"

	"github.com/gin-gonic/gin"
)

// @Description Request list of trading pairs
type listPairs struct {
	Pairs []entities.TradePair `json:"pairs" bson:"pairs"` // list of trade pools
} //@name ListPairs

// @Description Request list of trading pools
type listPools struct {
	Pools []entities.Pool `json:"pools" bson:"pools" gorm:"type:pools"` // list of trade pools
} //@name ListPools

// @Description Request list of tokens
type listTokens struct {
	Tokens []entities.Token `json:"tokens" bson:"tokens"` // list of tokens
} //@name TokenList

// @Description Token address input
// type tokenAddress struct {
// Address string `json:"address"`
// } // @name TokenAddress

// @Description Request for searching trade pair
type tokenPair struct {
	Protocol  entities.SwapProtocol `json:"protocol" bson:"protocol"`   // trade protocol
	TokenPair entities.TokenPair    `json:"tokenPair" bson:"tokenPair"` // pair of tokens
} //@name RequestPairs

// @Description Transaction response
// type transaction struct {
// 	Transaction string `json:"transaction" bson:"transaction"` // transaction
// } //@name Transaction

// @Description Response object
type response struct {
	Body interface{} `json:"body" bson:"body"` // returned data
} //@name Response

func respondOk(c *gin.Context, body interface{}) {
	c.JSON(200, body)
}

func respondCreated(c *gin.Context, body interface{}) {
	c.JSON(201, body)
}

func respondAccepted(c *gin.Context, body interface{}) {
	c.JSON(202, body)
}

// @Description Error response object
type responseErr struct {
	Error string `json:"error" example:"message"`
} //@name ErrorResponse

func errorNotFound(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(404, responseErr{msg})
}

func errorBadRequest(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(400, responseErr{msg})
}

func errorInufficientStorage(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(507, responseErr{msg})
}

func errorConflict(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(409, responseErr{msg})
}

func Log(f func(interface{}, ...interface{}), i interface{}, msg string) func() {
	return func() {
		f(i, msg)
	}
}
