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

// @Description Request list of protocols
type listProtocols struct {
	Protocols []entities.SwapProtocol `json:"protocols" bson:"protocols"` // list of protocols
} //@name ListProtocols

// @Description Request for searching trade pair
type tokenPair struct {
	Protocol  entities.SwapProtocol `json:"protocol" bson:"protocol"`   // trade protocol
	TokenPair entities.TokenPair    `json:"tokenPair" bson:"tokenPair"` // pair of tokens
} //@name RequestPairs

// @Description Response object
type response struct {
	Body interface{} `json:"body" bson:"body"` // returned data
} //@name Response

// 200 ok
func respondOk(c *gin.Context, body interface{}) {
	c.JSON(200, body)
}

// 201 created
func respondCreated(c *gin.Context, body interface{}) {
	c.JSON(201, body)
}

// 202 accepted
func respondAccepted(c *gin.Context, body interface{}) {
	c.JSON(202, body)
}

// @Description Error response object
type responseErr struct {
	Error string `json:"error" example:"message"`
} //@name ErrorResponse

// 404 client - not found
func errorNotFound(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(404, responseErr{msg})
}

// 400 client - bad request
func errorBadRequest(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(400, responseErr{msg})
}

// 409 client -conflict
func errorConflict(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(409, responseErr{msg})
}

// 507 server - insufficient storage
func errorInufficientStorage(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(507, responseErr{msg})
}

// 500 server - internal server error
func errorInternalServer(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(500, responseErr{msg})
}

// 502 server - bad gateway
func errorBadGateway(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(502, responseErr{msg})
}

// 503 server - service unavailable
func errorServiceUnavailable(c *gin.Context, msg string, log func()) {
	log()
	c.AbortWithStatusJSON(503, responseErr{msg})
}

func Log(f func(interface{}, ...interface{}), i interface{}, msg string) func() {
	return func() {
		f(i, msg)
	}
}
