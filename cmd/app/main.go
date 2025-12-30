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
		// Public routes
		api.POST("/login", authHandler.Login)

		// Authenticated routes
		authenticated := api.Group("/")
		authenticated.Use(middleware.JWTAuthMiddleware(config))

		// Admin routes
		admin := authenticated.Group("/admin")
		admin.Use(middleware.AuthorizationRole("admin"))
		{
			admin.POST("/users", userHandler.Create)
			admin.GET("/users", userHandler.List)
			admin.GET("/users/:id", userHandler.GetByID)
			admin.GET("/users/email/:email", userHandler.GetByEmail)
			admin.PUT("/users/:id", userHandler.Update)
			admin.DELETE("/users/:id", userHandler.Delete)
		}

		// Customer routes
		customer := authenticated.Group("/customer")
		customer.Use(middleware.AuthorizationRole("customer"))
		{
			customer.POST("/products", productHandler.CreateProduct)
			customer.GET("/products", productHandler.GetAllProducts)
			customer.GET("/products/:id", productHandler.GetProductByID)
			customer.PUT("/products/:id", productHandler.UpdateProduct)
			customer.DELETE("/products/:id", productHandler.DeleteProduct)

			customer.POST("/checkout", checkoutHandler.Checkout)
		}

		// Employees routes
		employee := authenticated.Group("/employee")
		employee.Use(middleware.AuthorizationRole("employee"))
		{
			employee.POST("/checkout", checkoutHandler.Checkout)
		}

		// Shared routes
		sharedCheckout := authenticated.Group("/checkout")
		sharedCheckout.Use(middleware.AuthorizationRole("customer", "employee"))
		{
			sharedCheckout.POST("/", checkoutHandler.Checkout)
		}
	}

	log.Printf("🚀 Server running on port %s", config.AppPort)
	router.Run(":" + config.AppPort)
}
