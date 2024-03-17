package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"webhook-router/handler"
	"os"
)

func main() {
	e := echo.New()
	if os.Getenv("DEBUG") == "True"{
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	api := e.Group("", serverHeader)
	api.GET("/", handler.Default)
	api.POST("/bot", handler.Bot)
	api.POST("/setWebhook", handler.SetWebhook)
	api.POST("/deleteWebhook", handler.DeleteWebhook)

	err := e.Start(":80")
	if err != nil {
		e.Logger.Fatal(err)
	}
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("x-version", "webhook-router/v2.0")
		return next(c)
	}
}
