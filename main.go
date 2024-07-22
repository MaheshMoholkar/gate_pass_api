package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"gate_pass_api/internal/database"
	"gate_pass_api/internal/handlers"
	"gate_pass_api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{
			"error": err.Error(),
		})
	},
}

func initDB() *sql.DB {
	cfg := database.DefaultSQLServerConfig()
	db, err := database.Open(cfg)
	if err != nil {
		panic(err)
	}
	return db
}

// Logging middleware

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")

	if os.Getenv("ASPNETCORE_PORT") != "" {
		port = os.Getenv("ASPNETCORE_PORT")
	}
	address := fmt.Sprintf(":%s", port)

	listenAddr := flag.String("listenAddr", address, "The listen address of the api server")
	flag.Parse()

	db := initDB()
	defer db.Close()

	store := database.New(db)

	var (
		app   = fiber.New(config)
		apiv1 = app.Group("/api/v1")

		// Initialize handlers
		authHandler    = handlers.NewAuthHandler(store)
		userHandler    = handlers.NewUserHandler(store)
		visitorHandler = handlers.NewVisitorHandler(store)
		staffHandler   = handlers.NewStaffHandler(store)
	)

	// Use middleware
	app.Use(middleware.Logger)
	apiv1.Use(middleware.JWTMiddleware())

	// Customize the CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Auth handlers
	app.Post("/auth/login", authHandler.HandleLogin)

	// User handlers
	app.Post("/signup", userHandler.HandleSignUp)
	apiv1.Get("/user", userHandler.HandleGetUser)
	apiv1.Put("/user", userHandler.HandleUpdateUser)
	apiv1.Put("/user/password", userHandler.HandleUpdatePassword)
	apiv1.Delete("/user", userHandler.HandleDeleteUser)

	// Visitor handlers
	apiv1.Get("/visitors", visitorHandler.HandleGetVisitors)
	apiv1.Post("/visitor-entry", visitorHandler.HandleVisitorEntryForm)
	apiv1.Put("/visitor-exit/:mobile_no", visitorHandler.HandleVisitorExit)

	// Staff handlers
	apiv1.Post("/staff-entry", staffHandler.HandleStaffEntryForm)
	apiv1.Put("/staff-exit/:id", staffHandler.HandleStaffExit)

	app.Listen(*listenAddr)
}
