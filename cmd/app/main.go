package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rkweber-max/checkout-backend/pkg/config"
	"github.com/rkweber-max/checkout-backend/pkg/database"
	"go.uber.org/fx"

	productHandler "github.com/rkweber-max/checkout-backend/internal/product/handler"
	productRepo "github.com/rkweber-max/checkout-backend/internal/product/repository"
	productService "github.com/rkweber-max/checkout-backend/internal/product/service"

	checkoutHandler "github.com/rkweber-max/checkout-backend/internal/checkout/handler"
	checkoutService "github.com/rkweber-max/checkout-backend/internal/checkout/service"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			newGinEngine,
			database.NewPostgresDB,
			productRepo.NewProductRepository,
			productService.NewProductService,
			productHandler.NewProductHandler,
			checkoutService.NewCheckoutService,
			checkoutHandler.NewCheckoutHandler,
		),
		fx.Invoke(registerRoutes),
	).Run()
}

func newGinEngine() *gin.Engine {
	return gin.New()
}

func registerRoutes(
	router *gin.Engine,
	productHandler *productHandler.ProductHandler,
	checkoutHandler *checkoutHandler.CheckoutHandler,
	config *config.Config,
) {
	router.Use(gin.Logger(), gin.Recovery())

	r := router.Group("/api")
	{
		r.POST("/products", productHandler.CreateProduct)
		r.GET("/products", productHandler.GetAllProducts)
		r.GET("/products/:id", productHandler.GetProductByID)
		r.PUT("/products/:id", productHandler.UpdateProduct)
		r.DELETE("/products/:id", productHandler.DeleteProduct)
		r.POST("/checkout", checkoutHandler.Checkout)
	}

	log.Printf("🚀 Server running on port %s", config.AppPort)
	router.Run(":" + config.AppPort)
}
