package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tradingsystem/internal/handler"
	"tradingsystem/internal/repository"
	"tradingsystem/internal/service"
	"tradingsystem/pkg/response"
)

func New(db *gorm.DB) *gin.Engine {
	engine := gin.Default()
	engine.LoadHTMLGlob("web/templates/*")

	repos := repository.New(db)
	strategyHandler := handler.NewStrategyHandler(service.NewStrategyService(repos.Strategy))
	marketDatasetHandler := handler.NewMarketDatasetHandler(service.NewMarketDatasetService(repos.MarketDataset))
	backtestHandler := handler.NewBacktestHandler(service.NewBacktestService(repos))

	engine.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/backtests/view")
	})

	engine.GET("/backtests/view", func(c *gin.Context) {
		c.HTML(200, "backtest_result.html", gin.H{
			"title": "Backtest Result",
		})
	})

	api := engine.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			response.Success(c, gin.H{
				"status": "ok",
			})
		})

		strategies := api.Group("/strategies")
		{
			strategies.POST("", strategyHandler.Create)
			strategies.GET("", strategyHandler.List)
			strategies.GET("/:id", strategyHandler.Get)
			strategies.PUT("/:id", strategyHandler.Update)
			strategies.DELETE("/:id", strategyHandler.Delete)
		}

		marketDatasets := api.Group("/market-datasets")
		{
			marketDatasets.POST("/import", marketDatasetHandler.Import)
			marketDatasets.GET("", marketDatasetHandler.List)
			marketDatasets.GET("/:id", marketDatasetHandler.Get)
			marketDatasets.DELETE("/:id", marketDatasetHandler.Delete)
		}

		backtests := api.Group("/backtests")
		{
			backtests.POST("", backtestHandler.Create)
			backtests.GET("", backtestHandler.List)
			backtests.GET("/:id", backtestHandler.Get)
			backtests.GET("/:id/trades", backtestHandler.ListTrades)
		}
	}

	return engine
}
