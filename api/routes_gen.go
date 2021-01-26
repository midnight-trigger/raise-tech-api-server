// This file was auto-generated.
// DO NOT EDIT MANUALLY!!!
package api

import (
	"github.com/labstack/echo"
	"github.com/midnight-trigger/raise-tech-api-server/api/controller"
)

func RegisterRoutes(e *echo.Echo) {
	HealthCheck(e, &controller.Image{})
	PostImage(e, &controller.Image{})
}
func RegisterAuthRoutes(e *echo.Group) {
}
func HealthCheck(
	e *echo.Echo,
	inter *controller.Image,
) {
	e.GET("health", func(c echo.Context) error {
		res := inter.HealthCheck(c)
		return c.JSON(res.Meta.Code, res)
	})
}
func PostImage(
	e *echo.Echo,
	inter *controller.Image,
) {
	e.POST("api/v1/image", func(c echo.Context) error {
		res := inter.PostImage(c)
		return c.JSON(res.Meta.Code, res)
	})
}
