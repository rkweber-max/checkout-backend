package main

import (
	"log"

	"github.com/gin-gonic/gin"
	authHandler "github.com/rkweber-max/checkout-backend/internal/handler"
	"github.com/rkweber-max/checkout-backend/internal/middleware"
	"github.com/rkweber-max/checkout-backend/pkg/config"
	"github.com/rkweber-max/checkout-backend/pkg/database"
	"go.uber.org/fx"

	userHandler "github.com/rkweber-max/checkout-backend/internal/user/handler"
	userRepo "github.com/rkweber-max/checkout-backend/internal/user/repository"
	userService "github.com/rkweber-max/checkout-backend/internal/user/service"

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
			authHandler.NewAuthHandler,
			userHandler.NewUserHandler,
			userRepo.NewUserRepository,
			userService.NewUserService,
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
	authHandler *authHandler.AuthHandler,
	userHandler *userHandler.UserHandler,
	productHandler *productHandler.ProductHandler,
	checkoutHandler *checkoutHandler.CheckoutHandler,
	config *config.Config,
) {
	router.Use(gin.Logger(), gin.Recovery())

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/users", userHandler.Create)

		authenticated := api.Group("/")
		authenticated.Use(middleware.JWTAuthMiddleware(config))

		authenticated.GET("/users", userHandler.List)
		authenticated.GET("/users/:id", userHandler.GetByID)
		authenticated.GET("/users/email/:email", userHandler.GetByEmail)
		authenticated.PUT("/users/:id", userHandler.Update)
		authenticated.DELETE("/users/:id", userHandler.Delete)

		authenticated.POST("/products", productHandler.CreateProduct)
		authenticated.GET("/products", productHandler.GetAllProducts)
		authenticated.GET("/products/:id", productHandler.GetProductByID)
		authenticated.PUT("/products/:id", productHandler.UpdateProduct)
		authenticated.DELETE("/products/:id", productHandler.DeleteProduct)

		authenticated.POST("/checkout", checkoutHandler.Checkout)
	}

	log.Printf("🚀 Server running on port %s", config.AppPort)
	router.Run(":" + config.AppPort)
}
