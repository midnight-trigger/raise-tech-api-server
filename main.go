package main

import (
	"net/http"
	"os"

	"github.com/midnight-trigger/raise-tech-api-server/api"
	"github.com/midnight-trigger/raise-tech-api-server/configs"
	"github.com/midnight-trigger/raise-tech-api-server/infra"
	"github.com/midnight-trigger/raise-tech-api-server/infra/mysql"
	"github.com/midnight-trigger/raise-tech-api-server/logger"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(context echo.Context) bool {
			if context.Request().URL.String() == "/health" {
				return true
			}
			return false
		},
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	configs.Init("")
	infra.Init()
	logger.Init("")

	defer mysql.Orm().Close()

	requiredAuthGroup := e.Group("")

	api.RegisterRoutes(e)
	api.RegisterAuthRoutes(requiredAuthGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Infof("Listening on port %s", port)

	http.Handle("/", e)

	e.Logger.Fatal(e.Start(":" + port))
}
