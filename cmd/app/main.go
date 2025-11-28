package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rkweber-max/checkout-backend/pkg/config"
	"github.com/rkweber-max/checkout-backend/pkg/database"
	"go.uber.org/fx"

	"github.com/rkweber-max/checkout-backend/internal/product/handler"
	"github.com/rkweber-max/checkout-backend/internal/product/repository"
	"github.com/rkweber-max/checkout-backend/internal/product/service"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			newGinEngine,
			database.NewPostgresDB,
			repository.NewProductRepository,
			service.NewProductService,
			handler.NewProductHandler,
		),
		fx.Invoke(registerRoutes),
	).Run()
}

func newGinEngine() *gin.Engine {
	return gin.New()
}

func registerRoutes(
	router *gin.Engine,
	h *handler.ProductHandler,
	config *config.Config,
) {
	router.Use(gin.Logger(), gin.Recovery())

	r := router.Group("/api")
	{
		r.POST("/products", h.CreateProduct)
		r.GET("/products", h.GetAllProducts)
		r.GET("/products/:id", h.GetProductByID)
		r.PUT("/products/:id", h.UpdateProduct)
		r.DELETE("/products/:id", h.DeleteProduct)
	}

	log.Printf("🚀 Server running on port %s", config.AppPort)
	router.Run(":" + config.AppPort)
}
