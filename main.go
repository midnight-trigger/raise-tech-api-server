package main

import (
	"net/http"
	"os"

	"github.com/midnight-trigger/todo/api"

	"github.com/midnight-trigger/todo/configs"
	"github.com/midnight-trigger/todo/infra"
	"github.com/midnight-trigger/todo/infra/mysql"
	"github.com/midnight-trigger/todo/logger"

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
	//requiredAuthGroup.Use(middleware.JWTWithConfig(jwt.GetMiddlewareJWTConfig()))

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
