package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/antonyuhnovets/flash-loan-arbitrage/docs"

	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Flash Loan Arbitrage bot
// @description interacting with smart contract, collecting tokens, perform trade
// @version     1.0
// @BasePath 	/v1
// @license.name  MIT
func NewRouter(h *gin.Engine, l logger.Interface, t tradecase.TradeCase) {
	// Options
	h.Use(gin.Logger())
	h.Use(gin.Recovery())

	// Swagger
	// docs.SwaggerInfo.BasePath = "/v1"
	// sh := swaggerFiles.NewHandler()
	// swaggerHandler := ginSwagger.WrapHandler(sh)
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	h.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	h.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	h.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	handler := h.Group("/v1")
	{
		NewTradecaseRouter(handler, t, l)
	}
}
